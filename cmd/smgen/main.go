package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wtask/sitemap/internal/compression"

	"github.com/wtask/sitemap/internal/sitemap/render"

	"github.com/wtask/sitemap/internal/sitemap"
)

const (
	MaxMapEntries     = 50000
	MaxMapSizeBytes   = 50000 * 1024 * 1024
	MaxIndexEntries   = MaxMapEntries
	MaxIndexSizeBytes = MaxMapSizeBytes
)

// mapSaver - func which saves a whole map or it part and returns size of file and error.
type mapSaver func(filename string, chunk []sitemap.MapItem) (int64, error)

type logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

func main() {
	var l logger = log.New(os.Stdout, "smgen ", log.Ldate|log.Ltime)

	l.Printf(
		"Started for %q, depth: %d, workers: %d, output format: %q, output dir: %s\n",
		startURL.String(),
		depth,
		numWorkers,
		outputFormat,
		outputDir,
	)

	parser, err := sitemap.NewParser(
		sitemap.WithErrorHandler(func(e error) {
			l.Println("ERR", e)
		}),
	)
	if err != nil {
		l.Println("Parser can not be started:", err)
		os.Exit(1)
	}

	l.Println("Parser has launched...")
	m := parser.Parse(startURL, depth, numWorkers)
	l.Println("Completed, num of links found:", len(m))
	if len(m) == 0 {
		l.Println("Stop on empty map")
		return
	}

	saver, err := buildMapSaver(outputFormat)
	if err != nil {
		l.Println("ERR", err)
		os.Exit(1)
	}

	l.Println("Started saving site map...")
	numErrors := 0
	files := []string{}
	for file, err := range saveMap(m, MaxMapEntries, MaxMapSizeBytes, mapFilename, outputFormat, outputDir, saver) {
		if err != nil {
			numErrors++
			l.Println("ERR", "MAP", file, err)
		} else {
			l.Printf("OK", "MAP", file)
		}
		files = append(files, file)
	}
	if numErrors > 0 {
		l.Println("Map saving stage done with error(s):", numErrors)
		os.Exit(1)
	}

	if len(files) > 1 {
		l.Println("Started saving index ...")
		numErrors = 0
		for file, err := range ensureIndex(files, MaxIndexEntries, MaxIndexSizeBytes, indexFilename, outputDir) {
			if err != nil {
				numErrors++
				l.Println("ERR", "INDEX", file, err)
			} else {
				l.Println("OK", "INDEX", file)
			}
		}
		if numErrors > 0 {
			l.Println("Index saving stage done with error(s):", numErrors)
			os.Exit(1)
		}
	}

	l.Println("All done")
}

// replaceWithGzip - compress file into gzip and remove origin if there was no error.
func replaceWithGzip(origin, gz string) error {
	err := compression.GzipFile(origin, gz)
	if err == nil {
		os.Remove(origin)
	} else {
		os.Remove(gz)
	}
	return err
}

// saveMap - asynchronously saves whole site map into files with no more than `itemsPerFile` in each.
// If resulting file size is over `maxFileSizeBytes`, source will be replaced with gzip-compressed one.
// Returns the map of file names and errors if any occurred when file was saving or compressing.
func saveMap(
	m []sitemap.MapItem,
	itemsPerFile int,
	maxFileSizeBytes int64,
	basename, extension, outputDir string,
	saver mapSaver,
) map[string]error {
	numFiles, reminder := len(m)/itemsPerFile, len(m)%itemsPerFile
	if reminder > 0 {
		numFiles++
	}
	wg := sync.WaitGroup{}
	mx := sync.Mutex{} // protects files
	files := make(map[string]error, numFiles)
	filename := fmt.Sprintf("%s.%s", basename, extension) // if single file only
	for i := 0; i < numFiles; i++ {
		if numFiles > 1 {
			filename = fmt.Sprintf("%s%d.%s", basename, i+1, extension)
		}
		start := i * itemsPerFile
		end := start + itemsPerFile
		if end > len(m) {
			end = len(m)
		}
		wg.Add(1)
		// save whole map in chunks, compress results if needed
		go func(filename string, chunk []sitemap.MapItem) {
			defer wg.Done()
			filesize, err := saver(filename, chunk)
			if err == nil && filesize > maxFileSizeBytes {
				// TODO We can start new goroutine here
				if err = replaceWithGzip(filename, filename+".gzip"); err == nil {
					filename += ".gzip"
				}
			}
			mx.Lock()
			files[filename] = err
			mx.Unlock()
		}(
			filepath.Join(outputDir, filename),
			m[start:end],
		)
	}

	wg.Wait()

	return files
}

// buildMapSaver - factory method to return map saving according given format.
// Now builds XML-saver only.
func buildMapSaver(format string) (mapSaver, error) {
	switch format {
	case "xml":
		return saveMapXML, nil
	default:
		return nil, fmt.Errorf("format %q is not supported", format)
	}
}

// saveMapXML - saves site map in XML format into the single file.
func saveMapXML(filename string, m []sitemap.MapItem) (int64, error) {
	var size int64
	if len(m) == 0 {
		return size, fmt.Errorf("site map is empty")
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return size, fmt.Errorf("can not open file: %s", err)
	}
	defer f.Close()
	err = render.XMLMap(f, m)
	st, _ := f.Stat()
	if st != nil {
		size = st.Size()
	}
	if err != nil {
		return size, fmt.Errorf("render site map (%d) as XML failed: %s", len(m), err)
	}

	return size, nil
}

// saveIndex - generate single map index and saves it in XML format.
func saveIndexXML(filename string, mapFiles []string) (int64, error) {
	var size int64
	if len(mapFiles) == 0 {
		return 0, fmt.Errorf("list of map files is empty")
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return size, fmt.Errorf("can not open file: %s", err)
	}
	defer f.Close()
	err = render.XMLIndex(f, time.Now().UTC(), mapFiles)
	st, _ := f.Stat()
	if st != nil {
		size = st.Size()
	}
	if err != nil {
		return size, fmt.Errorf("render site index (%d) failed: %s", len(mapFiles), err)
	}

	return size, nil
}

// ensureIndex - create a set of site map index files if needed.
func ensureIndex(
	filenames []string,
	maxIndexEntries int,
	maxIndexSizeBytes int64,
	basename, outputDir string,
) map[string]error {
	if len(filenames) <= 1 {
		return nil
	}
	numFiles, reminder := len(filenames)/maxIndexEntries, len(filenames)%maxIndexEntries
	if reminder > 0 {
		numFiles++
	}
	wg := sync.WaitGroup{}
	mx := sync.Mutex{} // protects files
	files := make(map[string]error, numFiles)
	for i := 0; i < numFiles; i++ {
		filename := fmt.Sprintf("%s%d.xml", basename, i+1)
		start := i * maxIndexEntries
		end := start + maxIndexEntries
		if end > len(filenames) {
			end = len(filenames)
		}
		wg.Add(1)
		go func(filename string, chunk []string) {
			defer wg.Done()
			filesize, err := saveIndexXML(filename, chunk)
			if err == nil && filesize > maxIndexSizeBytes {
				if err = replaceWithGzip(filename, filename+".gzip"); err == nil {
					filename += ".gzip"
				}
			}
			mx.Lock()
			files[filename] = err
			mx.Unlock()
		}(
			filepath.Join(outputDir, filename),
			filenames[start:end],
		)
	}

	wg.Wait()

	return files
}
