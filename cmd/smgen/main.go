package main

import (
	"fmt"
	"os"

	"github.com/wtask/sitemap/internal/sitemap"
)

func main() {
	parser, err := sitemap.NewParser()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m := parser.Parse(startURL, depth, numWorkers)
	for _, item := range m {
		fmt.Println(item.URI.String())
		if item.DocumentMetadata != nil && item.DocumentMetadata.Modified != nil {
			fmt.Println(*item.DocumentMetadata.Modified)
		}
		fmt.Println()
	}
}
