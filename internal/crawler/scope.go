package crawler

import (
	"net/url"
	"strings"
)

var allowedPrefixes = []string{"/doc", "/blog", "/ref"}

// Reports whether a normalized link is on go.dev under an allowed path prefix.
func InScope(link string) bool {
	parsed, err := url.Parse(link)
	if err != nil {
		return false
	}
	if parsed.Host != "go.dev" {
		return false
	}

	for _, prefix := range allowedPrefixes {
		path := parsed.Path
		// strings.HasPrefix to cover /ref/spec, /ref/mem, etc.
		if path == prefix || strings.HasPrefix(path, prefix+"/") {
			return true
		}
	}
	return false
}
