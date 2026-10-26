package analyze

// Complete list from https://gist.github.com/sebleier/554280 contains many noises (ab, bx, cz).
//
// For example, removing "go" and "for" isn’t enough
// because the list has many other words that matter in Go docs.
//
// Go keywords and idioms: if, else, for, go, in, ok.
// Standard library package names: io, os, http.
// Common doc terms: run, use, get, new, index, value, name, ref, example.
var stopwords = []string{
	"the", "a", "an", "and", "or", "but", "of", "to", "in", "on", "at", "by",
	"with", "from", "as", "is", "are", "was", "were", "be", "been",
	"this", "that", "these", "those", "it", "its", "you", "your", "we", "our",
	"will", "can",
	"s", "t", "ll", "ve", "re", "d", "m",
}
var stopwordsSet = toSet(stopwords)

// Remove stopwords from given tokens.
func RemoveStopwords(tokens []string) []string {
	var out []string
	for _, token := range tokens {
		if _, ok := stopwordsSet[token]; !ok {
			out = append(out, token)
		}
	}
	return out
}

func toSet(stopwords []string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, word := range stopwords {
		set[word] = struct{}{}
	}
	return set
}
