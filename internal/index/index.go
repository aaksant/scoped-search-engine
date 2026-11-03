package index

import "github.com/aaksant/scoped-search-engine/internal/store"

// Records that a term occurs in a document.
type Posting struct {
	DocID int64 // ID of the document containing the term.
	TF    int   // Term Frequency is number of occurrences in the document.
}

// Stores an inverted index and document metadata.
type Index struct {
	Postings map[string][]Posting // Maps each term to its postings (one per document).
	DocLen   map[int64]int        // Maps each document ID to its total token count.
}

// Creates a new index.
func New() *Index {
	return &Index{
		Postings: make(map[string][]Posting),
		DocLen:   make(map[int64]int),
	}
}

// Creates a new index and populates it with tokenized document.
func Build(docs []store.Document) *Index {
	idx := New()
	for _, doc := range docs {
		tokens := TransformDocumentToTokens(
			doc.Title,
			doc.Body,
			doc.Headings,
		)
		idx.Add(doc.ID, tokens)
	}
	return idx
}

// Add indexes a document by counting its term occurrences
// and storing its document length.
func (i *Index) Add(docID int64, tokens []string) {
	counts := make(map[string]int)
	for _, term := range tokens {
		counts[term]++
	}

	for term, count := range counts {
		i.Postings[term] = append(i.Postings[term], Posting{DocID: docID, TF: count})
	}

	i.DocLen[docID] = len(tokens)
}
