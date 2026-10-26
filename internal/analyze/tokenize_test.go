package analyze

import (
	"slices"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"punctuation", "Hello, World!", []string{"hello", "world"}},
		{"punctuation", "net/http", []string{"net", "http"}},
		{
			"punctuation",
			"You’ll use Gin (v1.7).",
			[]string{"you", "ll", "use", "gin", "v1", "7"},
		},
		{"empty", "", nil},
		{"whitespace", " ", []string{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Tokenize(test.in)
			if !slices.Equal(got, test.want) {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}
