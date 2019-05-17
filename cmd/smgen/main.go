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
	fmt.Println(m)
}
