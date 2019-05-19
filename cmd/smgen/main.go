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
		if meta := item.DocumentMeta; meta != nil && !meta.Modified.IsZero() {
			fmt.Println(meta.Modified)
		}
		fmt.Println()
	}
}
