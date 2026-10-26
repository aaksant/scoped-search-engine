package analyze

// Analyze turns raw text into index terms.
func Analyze(text string) []string {
	tokens := Tokenize(text)
	return RemoveStopwords(tokens)
}
