package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/wtask/sitemap/internal/sitemap"
)

var (
	startURL          *sitemap.URI
	outputFile        string
	numWorkers, depth uint
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
	flag.UintVar(&numWorkers, "w", 1, "Number of allowed concurrent workers to build site map.")
	flag.UintVar(&depth, "d", 0, "Maximum depth of link-junctions from start URL to render site map.")
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
		os.Exit(2)
	}
	if numWorkers == 0 {
		fmt.Fprint(flag.CommandLine.Output(), "Error: unable to start with 0 workers.\n\n")
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
	var err error
	startURL, err = sitemap.NewURI(start)
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Error: %v.\n\n", err)
		printUsage(flag.CommandLine.Output())
		os.Exit(2)
	}
	fmt.Fprintf(flag.CommandLine.Output(), "%+v\n\n", startURL)
}
