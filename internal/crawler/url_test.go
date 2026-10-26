package crawler

import (
	"net/url"
	"testing"
)

func TestUrl(t *testing.T) {
	baseURL, err := url.Parse("https://go.dev/")
	if err != nil {
		t.Fatal(err)
	}

	relativeURLs := []string{
		"/doc",
		"/learn#guided-learning-journeys",
		"/dl/",
		"#main-content",
		"https://pkg.go.dev/std",
		"https://GitHub.com/golang",
		"mailto:someone@example.com",
		"",
	}

	for _, rel := range relativeURLs {
		normalized, ok := NormalizeLink(baseURL, rel)
		if !ok {
			continue
		}
		t.Logf("normalized relative URL: %v\n", normalized)
	}
}
