package analyze

import (
	"slices"
	"testing"
)

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			"contains stopwords",
			"The Go Tutorial: You’ll use Gin",
			[]string{"go", "tutorial", "use", "gin"},
		},
		{
			"go idioms",
			"for loop if ok run mem",
			[]string{"for", "loop", "if", "ok", "run", "mem"},
		},
		{
			"all stopwords",
			"the and of",
			[]string{},
		},
		{"empty", "", nil},
		{"whitespace", " ", []string{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Analyze(test.in)
			if !slices.Equal(got, test.want) {
				t.Errorf("got %q, want %q", got, test.want)
			}
		})
	}
}
