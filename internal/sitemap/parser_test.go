package sitemap

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

	// fetch 0-level page and parse it links
	c := parser.worker(root, 1, Target{root, 0})

	if c.targets == nil {
		t.Fatal("Unexpected nil targets")
	}

	expected := []string{
		server.URL + "/faq.html",
		server.URL + "/protocol.html",
		// worker made non-unique link list
		// all checks are inside Parser method
		// TODO Try to exclude duplicates on worker level too
		server.URL + "/terms.html",
		server.URL + "/terms.html",
		// due to we have href="#top" inside homepage.html
		server.URL + "/homepage.html",
	}
	actual := []string{}
	for target := range c.targets {
		actual = append(actual, target.URI.String())
		if target.Level != 1 {
			t.Error("Unexpected level:", target.Level, "for URI:", target.URI.String())
		}
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected:", expected, "actual:", actual)
	}
}

func ExampleParser_Parse() {
	// This test example allows you not to sort the results.
	// Otherwise we need to prepend server.URL into expected results
	// and to sort both expected and actual sitemap before compare.

	server := httptest.NewServer(http.FileServer(http.Dir("testdata/simplesite")))
	defer server.Close()

	root, _ := NewURI(server.URL + "/homepage.html")
	serverURI, _ := NewURI(server.URL)
	parser, err := NewParser()
	if err != nil {
		panic("Unexpected NewParser() error: " + err.Error())
	}

	// At first, parser.Parse stores unique URLs inside map.
	// And then method generates slice from map.
	// Therefore, the constant order of items is not guaranteed.
	for _, item := range parser.Parse(root, 1, 2) {
		fmt.Println(
			strings.Replace(item.URI.String(), serverURI.String(), "{test-server-uri}/", -1),
		)
	}

	// Unordered output:
	// {test-server-uri}/homepage.html
	// {test-server-uri}/faq.html
	// {test-server-uri}/protocol.html
	// {test-server-uri}/terms.html
}

func ExampleParser_errorHandler() {
	server := httptest.NewServer(http.FileServer(http.Dir("testdata/simplesite")))
	defer server.Close()

	root, _ := NewURI(server.URL + "/homepage.witherror.html")
	serverURI, _ := NewURI(server.URL)
	parser, err := NewParser(
		WithErrorHandler(func(e error) {
			fmt.Println(
				strings.Replace(e.Error(), serverURI.String(), "{test-server-uri}/", -1),
			)
		}),
	)
	if err != nil {
		panic("Unexpected NewParser() error: " + err.Error())
	}

	parser.Parse(root, 1, 2)

	// Unordered output:
	// cannot fetch {test-server-uri}/notfound.php, status code: 404
}
