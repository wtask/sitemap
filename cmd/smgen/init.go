package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
)

var (
	startURL             *url.URL
	outputFile           string
	numWorkers, maxDepth int
)

func init() {
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
	flag.IntVar(&numWorkers, "w", 1, "Number of allowed concurrent workers to build site map.")
	flag.IntVar(&maxDepth, "d", 0, "Maximum depth of link-junctions from start URL to render site map.")
	flag.StringVar(&outputFile, "file", "sitemap.xml", "Write site map to given file.")

	flag.Parse()

	if help {
		printUsage(flag.CommandLine.Output())
		os.Exit(0)
	}

	start := flag.Arg(0)
	if start == "" {
		fmt.Fprint(flag.CommandLine.Output(), "Error: URI is required.\n\n")
		printUsage(flag.CommandLine.Output())
		os.Exit(1)
	}
	var err error
	startURL, err = url.ParseRequestURI(start)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: %v.\n\n", err)
		printUsage(flag.CommandLine.Output())
		os.Exit(1)
	}
	fmt.Fprintf(flag.CommandLine.Output(), "%+v\n\n", startURL)
}
