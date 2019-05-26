package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/wtask/sitemap/internal/sitemap/render"

	"github.com/wtask/sitemap/internal/sitemap"
)

const (
	MaxURIPerFile = 50000
	MaxFileSize   = 50000 * 1024 * 1024 * 1024
)

type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

func main() {
	var l Logger = log.New(os.Stdout, "smgen ", log.Ldate|log.Ltime)

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
	results := saveMap(m, MaxURIPerFile, mapFilename, outputFormat, outputDir, saver)
	numErrors := 0
	files := make([]string, len(results))
	for file, err := range results {
		if err != nil {
			numErrors++
			l.Println("ERR", file, err)
		} else {
			l.Println("OK", file)
		}
		files = append(files, file)
	}
	if numErrors > 0 {
		l.Println("Done with error(s):", numErrors)
		os.Exit(1)
	}

	// TODO Generate index if len(files) > 1 or filesize > MaxFileSize
	
	l.Println("All done")
}

// saveMap - asynchronously saves whole site map into files with no more than `itemsPerFile` in each.
// Returns map of file names and errors if any occurred when file was saving.
func saveMap(
	m []sitemap.MapItem,
	itemsPerFile int,
	basename, extension, outputDir string,
	saver func(filename string, chunk []sitemap.MapItem) error,
) map[string]error {
	numFiles, reminder := len(m)/itemsPerFile, len(m)%itemsPerFile
	if reminder > 0 {
		numFiles++
	}
	wg := sync.WaitGroup{}
	mx := sync.Mutex{} // protects files
	files := make(map[string]error, numFiles)
	filename := fmt.Sprintf("%s.%s", mapFilename, outputFormat) // if single file only
	for i := 0; i < numFiles; i++ {
		if numFiles > 1 {
			filename = fmt.Sprintf("%s%d.%s", mapFilename, i+1, outputFormat)
		}
		filename = filepath.Join(outputDir, filename)
		start := i * itemsPerFile
		end := start + itemsPerFile
		if end > len(m) {
			end = len(m)
		}
		wg.Add(1)
		go func(filename string, chunk []sitemap.MapItem) {
			defer wg.Done()
			err := saver(filename, chunk)
			mx.Lock()
			files[filename] = err
			mx.Unlock()
		}(filename, m[start:end])
	}

	wg.Wait()

	return files
}

func buildMapSaver(format string) (func(filename string, m []sitemap.MapItem) error, error) {
	switch format {
	case "xml":
		return saveXML, nil
	default:
		return nil, fmt.Errorf("format %q is not supported", format)
	}
}

func saveXML(filename string, m []sitemap.MapItem) error {
	if len(m) == 0 {
		return fmt.Errorf("site map is empty")
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("can not open file: %s", err)
	}
	defer f.Close()
	err = render.XMLMap(f, m)
	if err != nil {
		return fmt.Errorf("render site map (%d) as XML failed: %s", len(m), err)
	}
	return nil
}
