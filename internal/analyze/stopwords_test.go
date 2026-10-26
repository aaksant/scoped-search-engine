package analyze

import (
	"slices"
	"testing"
)

func TestStopwords(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{"no stopwords", []string{"for", "loop"}, []string{"for", "loop"}},
		{"1 stopwords", []string{"the", "go", "tutorial"}, []string{"go", "tutorial"}},
		{"full stopwords", []string{"the", "and", "of"}, []string{}},
		{"empty", []string{}, []string{}},
		{"nil", nil, []string{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := RemoveStopwords(test.in)
			if !slices.Equal(got, test.want) {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}
