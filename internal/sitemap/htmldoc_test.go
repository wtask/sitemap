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
		doc, _, err := fetchDocument(uri, 0)
		switch {
		case err != nil:
			if c.errMsg == "" {
				t.Error(c.path, "unexpected error:", err)
			} else if !strings.Contains(err.Error(), c.errMsg) {
				t.Error(c.path, "expected error contains:", c.errMsg, "got error:", err)
			}
			if doc != nil {
				t.Error("got non-nil document along with error")
			}
		case err == nil:
			if c.errMsg != "" {
				t.Error("expected error contains:", c.errMsg, "got error: nil")
			}
			if doc == nil {
				t.Error("got nil document without error")
			}
		}

	}
}

func Test_findNode(t *testing.T) {
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
	head := findFirstNode(doc, "head")
	if head == nil {
		t.Fatal("<head> node not found")
	}
	if base := findFirstNode(head, "base"); base == nil {
		t.Error("<base> node not found")
	}
	if style := findFirstNode(head, "style"); style != nil {
		t.Error("found nonexistent <style> node")
	}
	if h1 := findFirstNode(doc, "h1"); h1 == nil {
		t.Error("<h1> node not found")
	}
}

func Test_href(t *testing.T) {
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
	head := findFirstNode(doc, "head")
	if head == nil {
		t.Fatal("<head> node not found")
	}
	base := findFirstNode(head, "base")
	if base == nil {
		t.Fatal("<base> node not found")
	}
	expected := "http://host/path/"
	if link := href(base); link != expected {
		t.Errorf("Expected %q for href value of <base> node, got %q", expected, link)
	}
	title := findFirstNode(head, "title")
	if title == nil {
		t.Fatal("<title> node not found")
	}
	expected = ""
	if link := href(title); link != expected {
		t.Errorf("Expected %q for href value of <title> node, got %q", expected, link)
	}
}

func Test_collectLinks(t *testing.T) {
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
	</body>
	</html>	
	`
	doc, _ := html.Parse(bytes.NewReader([]byte(content)))
	body := findFirstNode(doc, "body")
	if body == nil {
		t.Fatal("<body> node not found")
	}

	links := collectLinks(body, nil)
	if len(links) != 2 {
		t.Error("Unexpected links num:", len(links))
	}
	if !reflect.DeepEqual(links, []string{"http://host/main.html", "/details.php"}) {
		t.Error("Unexpected links:", links)
	}
}
