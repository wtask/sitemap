package sitemap

import (
	"fmt"
	"net/url"
)

type URI struct {
	*url.URL
}

// NewURI - builds URI from given string and returns an error if any of the following conditions is true:
// 	- string does not contain allowed scheme (http or https)
// 	- string does not contain hostname
func NewURI(raw string) (*URI, error) {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return nil, err
	}
	if !(u.Scheme == "http" || u.Scheme == "https") {
		return nil, fmt.Errorf("sitemap.NewURI(): disallowed scheme %q for %q", u.Scheme, raw)
	}
	if u.Hostname() == "" {
		return nil, fmt.Errorf("sitemap.NewURI(): empty host %q", raw)
	}
	if u.Path == "" {
		u.Path = "/"
	}
	return &URI{u}, nil
}
