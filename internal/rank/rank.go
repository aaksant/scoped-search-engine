package rank

import (
	"cmp"
	"slices"
)

type Result struct {
	DocID int64
	Score float64
}

// Sorts results by score in descending order.
// When scores are equal, results are sorted by DocID in ascending order.
func sortResults(results []Result) {
	slices.SortFunc(results, func(x, y Result) int {
		if c := cmp.Compare(y.Score, x.Score); c != 0 {
			return c
		}
		return cmp.Compare(x.DocID, y.DocID)
	})
}
