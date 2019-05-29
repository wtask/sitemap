package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/wtask/sitemap/internal/sitemap"
)

var (
	// startURL - base URL from parser will start
	startURL *sitemap.URI
	// outputFormat - format for generating site map files;
	// index will always be saved as XML
	outputFormat = "xml"
	// mapFilename - base name for site map file, used to generate final names
	mapFilename,
	// indexFilename - base name for site map index file, the same as for map
	indexFilename,
	// outputDir - absolute path to directory where all files will be generated
	outputDir string
	// numWorkers - number of concurrent work instances which fetches and parses html documents
	numWorkers,
	// depth - link fetching depth
	depth uint
	// limitFileSizeBytes - maximum size in bytes of any generated file,
	// when file size is over this value, file is compressed into gzip
	limitFileSizeBytes int64
	// limitMapEntries - max number of entries per map file
	limitMapEntries,
	// limitIndexEntries - maximum number of entries per index file
	limitIndexEntries int
)

func init() {

	cwd, _ := os.Getwd()

	usage := `Generate site map suggested by https://www.sitemaps.org/protocol.html, starting from given URI:

	smgen [options] URI

`
	printUsage := func(out io.Writer) {
		fmt.Fprint(out, usage)
		fmt.Fprint(out, "Options:\n\n")
		flag.PrintDefaults()
		fmt.Fprint(out, "\n")
	}

	help := false
	flag.BoolVar(&help, "h", false, "")
	flag.BoolVar(&help, "help", false, "Print usage help.")
	flag.UintVar(&numWorkers, "num-workers", 1, "Number of allowed concurrent workers to build site map.")
	flag.UintVar(&depth, "depth", 1, "Maximum depth of link-junctions from start URL to render site map.")
	flag.StringVar(&mapFilename, "map-name", "sitemap", "Base name for site map FILE.")
	flag.StringVar(&indexFilename, "index-name", "sitemap_index", "Base name for site map INDEX.")
	flag.StringVar(&outputDir, "output-dir", cwd, "Output directory where site map and index will be generated.")
	flag.Int64Var(
		&limitFileSizeBytes,
		"size-limit",
		50000*1024*1024,
		"Maximum size of any generated file in bytes. If file size is greater than limitation, file is compressed into gzip.",
	)
	flag.IntVar(&limitMapEntries, "map-limit", 50000, "Limit number of entries per site map file.")
	flag.IntVar(&limitIndexEntries, "index-limit", 50000, "Limit number of entries per index file.")

	flag.Parse()

	if help {
		printUsage(flag.CommandLine.Output())
		os.Exit(0)
	}

	start := flag.Arg(0)
	if start == "" {
		fmt.Fprint(flag.CommandLine.Output(), "Error: URI is required.\n\n")
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
	if numWorkers == 0 {
		fmt.Fprint(flag.CommandLine.Output(), "Error: unable to start with 0 workers.\n\n")
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
	if mapFilename == "" {
		fmt.Fprint(flag.CommandLine.Output(), "Error: can not continue when site map filename is not specified.\n\n")
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
	if indexFilename == "" {
		fmt.Fprint(flag.CommandLine.Output(), "Error: can not continue when site map index filename is not specified.\n\n")
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
	if outputDir == "" {
		fmt.Fprint(flag.CommandLine.Output(), "Error: can not continue when output directory is not specified.\n\n")
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
	stat, err := os.Stat(outputDir)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Unable to check output directory: %v.\n\n", err)
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
	if !stat.IsDir() {
		fmt.Fprintf(flag.CommandLine.Output(), "Can not use output directory: %s.\n\n", outputDir)
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
	if limitFileSizeBytes <= 0 {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: invalid file size limitation (%d)\n\n", limitFileSizeBytes)
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
	if limitMapEntries < 1 {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: invalid map entries limitation (%d)\n\n", limitMapEntries)
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
	if limitIndexEntries < 1 {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: invalid index entries limitation (%d)\n\n", limitIndexEntries)
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}

	startURL, err = sitemap.NewURI(start)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: %v.\n\n", err)
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
}
