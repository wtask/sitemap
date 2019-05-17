package sitemap

import (
	"context"
	"fmt"
	"net/http"

	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

// fetchDocument - build http GET request, fetch response body abd parse given HTML into document tree.
// Only "text/html" content type is fetched.
func fetchDocument(uri *URI, timeout time.Duration) (*html.Node, *DocumentMetadata, error) {
	var cancel context.CancelFunc
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()
	ctx := context.Background()
	// TODO overwrite timeout if it is 0, for example set it to max allowed timeout
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}
	// TODO to avoid non text/html responses may to use HEAD first?
	url := uri.String()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("preparing request to %q failed: %s", url, err)
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		return nil, nil, fmt.Errorf("request to %s failed: %s", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("cannot fetch %s, status code: %d", url, resp.StatusCode)
	}

	// TODO Need to check final URI from response: resp.Request.URL.String()
	// If redirect was occurred we get different URL and we have to deal with this situation.

	ctype := resp.Header.Get("Content-Type")
	if !strings.Contains(ctype, "text/html") {
		return nil, nil, fmt.Errorf("%s, invalid content type: %q", url, ctype)
	}

	utf8, err := charset.NewReader(resp.Body, ctype)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to decode %q: %s", url, err)
	}

	var modified *time.Time
	t, err := http.ParseTime(resp.Header.Get("Last-Modified"))
	if err != nil {
		modified = &t
	}
	meta := &DocumentMetadata{
		Modified: modified,
	}
	doc, err := html.Parse(utf8)

	return doc, meta, err
}

// findFirstNode - parses elements tree to find first node for given tag.
// Returns nil if node for tag is not found.
func findFirstNode(tag string, tree *html.Node) *html.Node {
	if tree == nil {
		return nil
	}
	if tree.Type == html.ElementNode && tree.Data == tag {
		return tree
	}
	for n := tree.FirstChild; n != nil; n = n.NextSibling {
		if node := findFirstNode(tag, n); node != nil {
			return node
		}
	}
	return nil
}

// attribute - returns value of given attribute name for single document element.
// Returns empty string if attribute is not found.
func attribute(name string, element *html.Node) string {
	if element == nil {
		return ""
	}
	for _, a := range element.Attr {
		if a.Key == name {
			return a.Val
		}
	}
	return ""
}

// collectAttributes - parses elements tree and collects all attributes values for given tag.
// You can pass nil for values, but always check length of results.
func collectAttributes(tag, attr string, tree *html.Node, values []string) []string {
	if tree == nil {
		return values
	}
	if tree.Type == html.ElementNode && tree.Data == tag {
		if v := attribute(attr, tree); v != "" {
			values = append(values, v)
		}
	}
	for n := tree.FirstChild; n != nil; n = n.NextSibling {
		values = collectAttributes(tag, attr, n, values)
	}
	return values
}
