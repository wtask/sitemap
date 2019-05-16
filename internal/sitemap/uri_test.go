package sitemap

import "testing"

func TestNewURI(t *testing.T) {
	cases := []struct {
		rawURL string
		valid  bool
		uri    string
	}{
		{"localhost", false, ""},
		{"/dir/page.html", false, ""},
		{"../dir/page.html", false, ""},
		{"//localhost", false, ""},
		{"ssh://localhost", false, ""},
		{"ftp://localhost", false, ""},
		{"http://localhost:8080", true, "http://localhost:8080/"},
		{"http://localhost", true, "http://localhost/"},
		{"https://localhost", true, "https://localhost/"},
		{"http://localhost/", true, "http://localhost/"},
		{"http://user:pass@localhost", true, "http://user:pass@localhost/"},
		{"http://localhost/?search=text", true, "http://localhost/?search=text"},
		{"http://localhost/cool search/?q=text", true, "http://localhost/cool%20search/?q=text"},
	}
	for _, c := range cases {
		u, err := NewURI(c.rawURL)
		if err != nil && c.valid {
			t.Errorf("Unexpected error %q for %q", err, c.rawURL)
		}
		if err == nil && !c.valid {
			t.Errorf("Error was expected for %q, but it did not happened", c.rawURL)
		}
		if c.uri != "" && c.uri != u.String() {
			t.Errorf("Expected URI %q, got %q for %q", c.uri, u.String(), c.rawURL)
		}
	}
}
