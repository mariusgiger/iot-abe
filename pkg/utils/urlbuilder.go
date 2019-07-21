package utils

import (
	"net/url"
	"path"
)

// URLBuilder is a helper to modify URL
// Can add component to the path.
// Can add query parameters.
type URLBuilder struct {
	url   *url.URL
	query url.Values
}

// ParseURL initializes new URL builder
func ParseURL(rawURL string) (*URLBuilder, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err // do not wrap, use error "as is"
	}

	return &URLBuilder{
		url:   u,
		query: u.Query(),
	}, nil // OK
}

// MustParseURL initializes new URL builder
// panic if rawURL is not valid and cannot be parsed
func MustParseURL(rawURL string) *URLBuilder {
	ub, err := ParseURL(rawURL)
	if err != nil {
		panic(err) // unlikely
	}
	return ub
}

// SetPath replaces the whole URI path.
// newPath is cleared before assignment
func (ub *URLBuilder) SetPath(newPath string) *URLBuilder {
	ub.url.Path = path.Clean(newPath)
	return ub // self reference
}

// AddPath adds a URI path component to the end.
func (ub *URLBuilder) AddPath(pathToAdd string) *URLBuilder {
	ub.url.Path = path.Join(ub.url.Path, pathToAdd)
	return ub // self reference
}

// SetQuery sets the query parameter.
// previous value (if key exists) replaced
func (ub *URLBuilder) SetQuery(key, value string) *URLBuilder {
	ub.query.Set(key, value)
	return ub // self reference
}

// AddQuery adds the new query parameter.
// new value is added, so it's possible to send multiple values for the same key
// `?key=a&key=b&key=c`
func (ub *URLBuilder) AddQuery(key, value string) *URLBuilder {
	ub.query.Add(key, value)
	return ub // self reference
}

// String converts resulting URL back to string,
func (ub *URLBuilder) String() string {
	ub.url.RawQuery = ub.query.Encode()
	return ub.url.String()
}
