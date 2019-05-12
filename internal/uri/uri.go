package uri

import (
	"fmt"
	"net/url"
)

// FromString - builds URL from given string and returns an error if any of the following conditions is true:
// 	- rawURL does not contain allowed scheme (http or https);
// 	- there is no hostname detected in the source url
func FromString(rawURL string) (*url.URL, error) {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}
	if !(u.Scheme == "http" || u.Scheme == "https") {
		return nil, fmt.Errorf("uri.FromString(): disallowed scheme %q for %q", u.Scheme, rawURL)
	}
	if u.Hostname() == "" {
		return nil, fmt.Errorf("uri.FromString(): empty host %q", rawURL)
	}
	if u.Path == "" {
		u.Path = "/"
	}
	return u, nil
}
