package sitemap

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func Test_fetchDocument(t *testing.T) {
	cases := []struct {
		path   string
		errMsg string
	}{
		{"/nonexistent.html", "/nonexistent.html, status code: 404"},
		{"/valid.html", ""},
		{"/invalid.html", ""},
		{"/text.txt", "/text.txt, invalid content type: \"text/plain; charset=utf-8\""},
		{"/script.js", "/script.js, invalid content type: \"application/javascript\""},
		{"/json.json", "/json.json, invalid content type: \"application/json\""},
		{"/image.jpeg", "/image.jpeg, invalid content type: \"image/jpeg\""},
	}

	server := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer server.Close()

	for _, c := range cases {
		uri, _ := NewURI(server.URL + c.path)
		t.Log(uri.String())
		doc, meta, err := fetchDocument(uri, 0)
		switch {
		case err != nil:
			if c.errMsg == "" {
				t.Error(c.path, "unexpected error:", err)
			} else if !strings.Contains(err.Error(), c.errMsg) {
				t.Error(c.path, "expected error contain:", c.errMsg, "got error:", err)
			}
			if doc != nil {
				t.Error("Got non-nil document along error")
			}
			if meta != nil {
				t.Error("Got non-nil document metadata along error:", *meta)
			}
		case err == nil:
			if c.errMsg != "" {
				t.Error("expected error contain:", c.errMsg, "got error: nil")
			}
			if doc == nil {
				t.Error("Got nil document without error")
			}
			if meta == nil {
				t.Error("Got nil document metadata without error")
			}
		}

	}
}

func Test_findFirstNode(t *testing.T) {
	content := `
	<!doctype html>

	<html lang="en">
	<head>
		<base href="http://host/path/">
		<meta charset="utf-8">
		<title>Test valid HTML</title>
		<meta name="description" content="Test valid HTML">
		<meta name="author" content="github.com/wtask/sitemap">
	</head>

	<body>
		<h1>Test valid HTML</h1>
	</body>
	</html>	
	`

	doc, _ := html.Parse(bytes.NewReader([]byte(content)))
	head := findFirstNode("head", doc)
	if head == nil {
		t.Fatal("<head> node not found")
	}
	if base := findFirstNode("base", head); base == nil {
		t.Error("<base> node not found")
	}
	if style := findFirstNode("style", head); style != nil {
		t.Error("found nonexistent <style> node")
	}
	if h1 := findFirstNode("h1", doc); h1 == nil {
		t.Error("<h1> node not found")
	}

	nothing := findFirstNode("body", nil)
	if nothing != nil {
		t.Error("Found document tree or element for nil")
	}
}

func Test_attribute(t *testing.T) {
	content := `
	<head>
		<base href="http://host/path/">
		<meta charset="utf-8">
		<title>Test valid HTML</title>
		<meta name="description" content="Test valid HTML">
		<meta name="author" content="github.com/wtask/sitemap">
	</head>
	`
	doc, _ := html.Parse(bytes.NewReader([]byte(content)))
	head := findFirstNode("head", doc)
	if head == nil {
		t.Fatal("<head> node not found")
	}
	base := findFirstNode("base", head)
	if base == nil {
		t.Fatal("<base> node not found")
	}
	expected := "http://host/path/"
	if link := attribute("href", base); link != expected {
		t.Errorf("Expected %q for href value of <base> node, got %q", expected, link)
	}
	title := findFirstNode("title", head)
	if title == nil {
		t.Fatal("<title> node not found")
	}
	expected = ""
	if link := attribute("href", title); link != expected {
		t.Errorf("Expected %q for href value of <title> node, got %q", expected, link)
	}

	nothing := attribute("src", nil)
	if nothing != "" {
		t.Error("Got non-empty value for nil element:", nothing)
	}
}

func Test_collectAttributes(t *testing.T) {
	content := `
	<!doctype html>

	<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>Test valid HTML</title>
		<meta name="description" content="Test valid HTML">
		<meta name="author" content="github.com/wtask/sitemap">
	</head>

	<body>
		<h1>Test valid HTML</h1>
		<a href="http://host/main.html">Home page</a>
		<p>Text paragraph, see <a href="/details.php">details</a>
		</p>
		<a>Invalid anchor</a>
	</body>
	</html>	
	`
	doc, _ := html.Parse(bytes.NewReader([]byte(content)))

	head := findFirstNode("head", doc)
	if head == nil {
		t.Fatal("<head> node not found")
	}
	names := collectAttributes("meta", "name", head, nil)
	if len(names) != 2 {
		t.Error("Unexpected num of names:", len(names))
	}
	expected := []string{"description", "author"}
	if !reflect.DeepEqual(names, expected) {
		t.Error("Expected names:", expected, "got:", names)
	}

	body := findFirstNode("body", doc)
	if body == nil {
		t.Fatal("<body> node not found")
	}
	links := collectAttributes("a", "href", body, nil)
	if len(links) != 2 {
		t.Error("Unexpected num of links:", len(links))
	}
	expected = []string{"http://host/main.html", "/details.php"}
	if !reflect.DeepEqual(links, expected) {
		t.Error("Expected links:", expected, "got:", links)
	}

	nothing := collectAttributes("img", "scr", nil, nil)
	if nothing != nil || len(nothing) != 0 {
		t.Error("Unexpected non-empty collection for nil tree:", nothing)
	}
}
