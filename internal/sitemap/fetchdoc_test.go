package sitemap

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
