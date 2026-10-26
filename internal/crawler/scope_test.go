package crawler

import "testing"

func TestScope(t *testing.T) {
	links := []string{
		"https://go.dev/doc",
		"https://go.dev/blog",
		"go.dev/ref",
		"https://go.dev",
		"https://go.dev/ref/",
		"https://go.dev/ref/mem",
	}

	for _, l := range links {
		if ok := InScope(l); ok {
			t.Logf("link in scope: %v\n", l)
		}
	}
}
