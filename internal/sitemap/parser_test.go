package sitemap

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestParser_worker(t *testing.T) {
	server := httptest.NewServer(http.FileServer(http.Dir("testdata/simplesite")))
	defer server.Close()

	root, _ := NewURI(server.URL + "/homepage.html")
	t.Log(root.String())
	parser, err := NewParser()
	if err != nil {
		t.Fatal("Unexpected NewParser() error:", err)
	}

	var level uint
	// fetch 0-level page and parse it links
	c := parser.worker(root, 1, Target{root, level})

	if c.targets == nil {
		t.Fatal("Unexpected nil targets")
	}

	expected := []string{
		server.URL + "/faq.html",
		server.URL + "/protocol.html",
		server.URL + "/terms.html",
	}
	actual := []string{}
	for target := range c.targets {
		// ignoring DocumentMetadata
		actual = append(actual, target.URI.String())
		if target.Level != level+1 {
			t.Error("Unexpected level:", target.Level, "for URI:", target.URI.String())
		}
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected:", expected, "actual:", actual)
	}
}

func ExampleParser_Parse() {
	// This test example allows you not to sort the results
	// Otherwise we need to prepend server.URL to expected results
	// and sort both expected and actual sitemap before compare.

	server := httptest.NewServer(http.FileServer(http.Dir("testdata/simplesite")))
	defer server.Close()

	root, _ := NewURI(server.URL + "/homepage.html")
	parser, err := NewParser()
	if err != nil {
		panic("Unexpected NewParser() error: " + err.Error())
	}

	// At first, parser.Parse stores unique URLs inside map.
	// And then method generates slice from map.
	// Therefore, the constant order of items is not guaranteed.
	for _, item := range parser.Parse(root, 1, 2) {
		fmt.Println(item.URI.Scheme, item.URI.Hostname(), item.URI.EscapedPath())
	}

	// Unordered output:
	// http 127.0.0.1 /homepage.html
	// http 127.0.0.1 /faq.html
	// http 127.0.0.1 /protocol.html
	// http 127.0.0.1 /terms.html
}
