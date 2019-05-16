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

func fetchDocument(uri *URI, timeout time.Duration) (*html.Node, *documentMetadata, error) {
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
	meta := &documentMetadata{
		modified: modified,
	}
	doc, err := html.Parse(utf8)

	return doc, meta, err
}

func findFirstNode(doc *html.Node, node string) *html.Node {
	if doc.Type == html.ElementNode && doc.Data == node {
		return doc
	}
	for n := doc.FirstChild; n != nil; n = n.NextSibling {
		if r := findFirstNode(n, node); r != nil {
			return r
		}
	}
	return nil
}

func href(node *html.Node) string {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}

func collectLinks(node *html.Node, links []string) []string {
	if node.Type == html.ElementNode && node.Data == "a" {
		if link := href(node); link != "" {
			links = append(links, link)
		}
	}
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		links = collectLinks(n, links)
	}
	return links
}
