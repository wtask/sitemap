package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/wtask/sitemap/internal/sitemap"
)

const ()

var (
	startURL                   *sitemap.URI
	outputFormat               = "xml"
	mapFilename, indexFilename string
	outputDir                  string
	numWorkers, depth          uint
)

func init() {

	cwd, _ := os.Getwd()

	usage := `Generate XML site map suggested by https://www.sitemaps.org/protocol.html, starting from given URI:

	smgen [options] URI

`
	printUsage := func(out io.Writer) {
		fmt.Fprint(out, usage)
		fmt.Fprint(out, "Options:\n\n")
		flag.PrintDefaults()
		fmt.Fprint(out, "\n")
	}

	help := false
	flag.BoolVar(&help, "h", false, "Print usage help.")
	flag.UintVar(&numWorkers, "w", 1, "Number of allowed concurrent workers to build site map.")
	flag.UintVar(&depth, "d", 1, "Maximum depth of link-junctions from start URL to render site map.")
	flag.StringVar(&mapFilename, "map", "sitemap", "Site map FILE name [without ext].")
	flag.StringVar(&indexFilename, "index", "sitemap_index", "Site map INDEX filename [without extension].")
	flag.StringVar(&outputDir, "dir", cwd, "Output directory where site map and index will be generated.")

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

	startURL, err = sitemap.NewURI(start)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: %v.\n\n", err)
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
}
