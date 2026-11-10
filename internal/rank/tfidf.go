package rank

import (
	"math"

	"github.com/aaksant/scoped-search-engine/internal/index"
)

func TFIDF(idx *index.Index, terms []string) []Result {
	scoreMap := make(map[int64]float64)

	for _, term := range terms {
		postings, ok := idx.Postings[term]
		if !ok {
			continue
		}
		// Computes the IDF value of current term.
		//
		// ln(N / df)
		//
		// N is the number of documents (len(idx.DocLen)) and
		// df is the number of postings for the term.
		idf := math.Log(float64(len(idx.DocLen)) / float64(len(postings)))
		// Σ over query terms of TF(term, doc) * IDF(term)
		for _, posting := range postings {
			scoreMap[posting.DocID] += idf * float64(posting.TF)
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
