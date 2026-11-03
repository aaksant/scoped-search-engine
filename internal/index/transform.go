package index

import (
	"slices"
	"strings"

	"github.com/aaksant/scoped-search-engine/internal/analyze"
)

const (
	titleWeight    = 3
	headingsWeight = 2
)

// Transforms given document title, body, and headings into tokens.
func TransformDocumentToTokens(
	title,
	body string,
	headings []string,
) []string {
	// "," as the separator, a heading that contains comma gets cut into three pieces
	joinedHeadings := strings.Join(headings, "\n")
	return slices.Concat(
		slices.Repeat(analyze.Analyze(title), titleWeight),
		slices.Repeat(analyze.Analyze(joinedHeadings), headingsWeight),
		analyze.Analyze(body),
	)
}
