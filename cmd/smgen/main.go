package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/wtask/sitemap/internal/sitemap/render"

	"github.com/wtask/sitemap/internal/sitemap"
)

const (
	MaxURIPerFile = 50000
	MaxFileSize   = 50000 * 1024 * 1024 * 1024
)

func main() {
	logger := log.New(os.Stdout, "smgen ", log.Ldate|log.Ltime)
	logger.Printf(
		"Started for %q, depth: %d, workers: %d, output format: %q\n",
		startURL.String(),
		depth,
		numWorkers,
		outputFormat,
	)

	parser, err := sitemap.NewParser(
		sitemap.WithErrorHandler(func(e error) {
			logger.Println("Error:", e)
		}),
	)
	if err != nil {
		logger.Println("Parser can not be started:", err)
		os.Exit(1)
	}

	logger.Println("Parser has launched...")
	m := parser.Parse(startURL, depth, numWorkers)
	logger.Println("Completed, num of links found:", len(m))
	if len(m) == 0 {
		logger.Println("Stop on empty map")
		return
	}

	// save map
	numFiles, reminder := len(m)/MaxURIPerFile, len(m)%MaxURIPerFile
	if reminder > 0 {
		numFiles++
	}
	saver, err := buildMapSaver(outputFormat)
	if err != nil {
		logger.Println("Error:", err)
		os.Exit(1)
	}
	logger.Println("Started saving site map...")
	var (
		wg        sync.WaitGroup
		numErrors uint32
		files     []string
	)
	files = make([]string, numFiles)
	for i := 0; i < numFiles; i++ {
		filename := ""
		if numFiles > 1 {
			filename = fmt.Sprintf("%s%d.%s", mapFilename, i+1, outputFormat)
		} else {
			filename = fmt.Sprintf("%s.%s", mapFilename, outputFormat)
		}
		start := i * MaxURIPerFile
		end := start + MaxURIPerFile
		if end > len(m) {
			end = len(m)
		}
		files[i] = filepath.Join(outputDir, filename)
		wg.Add(1)
		go func(file string, m []sitemap.MapItem) {
			defer wg.Done()
			if err := saver(file, m); err != nil {
				atomic.AddUint32(&numErrors, 1)
				logger.Println("Error:", err)
			} else {
				logger.Println("OK:", file)
			}
		}(files[i], m[start:end])
	}

	wg.Wait()

	if numErrors > 0 {
		logger.Println("Done with error(s):", numErrors)
		os.Exit(1)
	}
	logger.Println("All done")
}

func buildMapSaver(format string) (func(filename string, m []sitemap.MapItem) error, error) {
	switch format {
	case "xml":
		return saveXML, nil
	default:
		return nil, fmt.Errorf("saving site map into %q is not supported", format)
	}
}

func saveXML(filename string, m []sitemap.MapItem) error {
	if len(m) == 0 {
		return fmt.Errorf("site map is empty, nothing to save into %s", filename)
	}
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("can not open file %q: %s", filename, err)
	}
	defer f.Close()
	if err := render.XMLMap(f, m); err != nil {
		return fmt.Errorf("render site map (%d) into XML failed: %s", len(m), err)
	}
	return nil
}
