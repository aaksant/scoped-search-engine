package rank

import (
	"math"

	"github.com/aaksant/scoped-search-engine/internal/index"
)

// Empirical defaults, not derived from anything.
// They come from the original Okapi BM25 experiments on information-retrieval test collections,
// where researchers tried values and kept the ones that ranked best on average.

// k1 sets how fast term frequency saturates.
// At average length, the TF part of the score is tf × (k1+1) / (tf + k1),
// which can never exceed k1 + 1.
//
// With k1 = 1.2 that ceiling is 2.2, and a single mention already scores exactly 1.0.
//
// So no number of repeats can be worth more than 2.2 mentions’ worth of credit.
const k1 float64 = 1.2

// b sets how much document length matters.
// At b = 0 length is ignored entirely, and at b = 1 the TF is fully scaled by dl / avgdl.
const b float64 = 0.75

func BM25(idx *index.Index, terms []string) []Result {
	scoreMap := make(map[int64]float64)

	if len(idx.DocLen) == 0 {
		return nil
	}
	var total int
	for _, l := range idx.DocLen {
		total += l
	}
	avgdl := float64(total) / float64(len(idx.DocLen))

	for _, term := range terms {
		postings, ok := idx.Postings[term]
		if !ok {
			continue
		}

		N := float64(len(idx.DocLen))
		df := float64(len(postings))
		idf := math.Log(1 + (N-df+0.5)/(df+0.5))

		for _, posting := range postings {
			tf := float64(posting.TF)
			dl := float64(idx.DocLen[posting.DocID])
			numerator := idf * tf * (k1 + 1)
			denominator := tf + k1*(1-b+(b*dl/avgdl))
			scoreMap[posting.DocID] += numerator / denominator
		}
	}

	var results []Result
	for docID, score := range scoreMap {
		results = append(results, Result{
			DocID: docID,
			Score: score,
		})
	}
	sortResults(results)
	return results
}
