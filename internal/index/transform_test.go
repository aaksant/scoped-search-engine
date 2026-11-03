package index

import (
	"slices"
	"testing"
)

func TestTransformDocumentToTokens(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		body     string
		headings []string
		want     []string
	}{
		{
			name:     "document with all fields",
			title:    "Go Programming",
			body:     "Learn how to write Go code efficiently.",
			headings: []string{"Introduction", "Getting Started"},
			want: []string{
				"go", "programming",
				"go", "programming",
				"go", "programming",
				"introduction", "getting", "started",
				"introduction", "getting", "started",
				"learn", "write", "go", "code", "efficiently",
			},
		},
		{
			name:     "empty headings and body",
			title:    "Quick Title",
			body:     "",
			headings: nil,
			want: []string{
				"quick", "title",
				"quick", "title",
				"quick", "title",
			},
		},
		{
			name:     "empty title, headings, body",
			title:    "",
			body:     "",
			headings: []string{},
			want:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TransformDocumentToTokens(tt.title, tt.body, tt.headings)
			if !slices.Equal(got, tt.want) {
				t.Errorf("got: %q, want %q", got, tt.want)
			}
		})
	}
}
