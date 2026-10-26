package crawler

import (
	"net/url"
	"strings"
)

// Drops fragments, lowercases the host, and trims trailing slashes.
func NormalizeLink(baseURL *url.URL, href string) (string, bool) {
	parsed, err := url.Parse(href)
	if err != nil {
		return "", false
	}
	if href == "" {
		return "", false
	}

	absolute := baseURL.ResolveReference(parsed)
	absolute.Fragment = ""
	absolute.Host = strings.ToLower(absolute.Host)
	if absolute.Scheme != "http" && absolute.Scheme != "https" {
		return "", false
	}
	if absolute.Path != "/" {
		absolute.Path = strings.TrimSuffix(absolute.Path, "/")
	}

	return absolute.String(), true
}

// Resolves relative links against the URL after redirects,
// which avoids wrong paths like /install instead of /doc/install.
func ResolveLinks(baseURL *url.URL, hrefs []string) []string {
	var resolved []string
	for _, href := range hrefs {
		normalized, ok := NormalizeLink(baseURL, href)
		if !ok {
			continue
		}
		resolved = append(resolved, normalized)
	}
	return resolved
}
