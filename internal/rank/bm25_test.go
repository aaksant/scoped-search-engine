package rank

import (
	"math"
	"testing"

	"github.com/aaksant/scoped-search-engine/internal/index"
)

func TestBM25(t *testing.T) {
	t.Run("hand-computed score", func(t *testing.T) {
		idx := &index.Index{
			Postings: map[string][]index.Posting{
				"go": {
					{DocID: 1, TF: 2},
				},
			},
			DocLen: map[int64]int{
				1: 10,
			},
		}

		results := BM25(idx, []string{"go"})

		// N = 1, df = 1, TF = 2, DL = 10, avgDL = 10
		//
		// IDF = ln(1 + (1 - 1 + 0.5) / (1 + 0.5))
		//     = ln(4/3)
		//
		// BM25 = IDF * (2 * (1.2 + 1)) / (2 + 1.2)
		//      = ln(4/3) * 1.375s
		expected := math.Log(4.0/3.0) * 1.375

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].DocID != 1 {
			t.Errorf("expected DocID 1, got %d", results[0].DocID)
		}
		if math.Abs(results[0].Score-expected) > 1e-9 {
			t.Errorf("expected score %.9f, got %.9f", expected, results[0].Score)
		}
	})

	t.Run("length normalization", func(t *testing.T) {
		idx := &index.Index{
			Postings: map[string][]index.Posting{
				"go": {
					{DocID: 1, TF: 2},
					{DocID: 2, TF: 2},
				},
			},
			DocLen: map[int64]int{
				1: 5,
				2: 15,
			},
		}

		results := BM25(idx, []string{"go"})

		// Both documents have the same TF, but document 1 is shorter.
		// Its score should therefore be higher.

		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}
		if results[0].DocID != 1 {
			t.Errorf("expected shorter document (DocID 1) to rank first, got %d", results[0].DocID)
		}
		if results[0].Score <= results[1].Score {
			t.Errorf("expected shorter document to have a higher score, got %.9f <= %.9f",
				results[0].Score, results[1].Score)
		}
	})

	t.Run("term frequency saturation", func(t *testing.T) {
		newIndex := func(tf int) *index.Index {
			return &index.Index{
				Postings: map[string][]index.Posting{
					"go": {
						{DocID: 1, TF: tf},
					},
				},
				DocLen: map[int64]int{
					1: 10,
					2: 10,
				},
			}
		}

		single := BM25(newIndex(1), []string{"go"})
		doubled := BM25(newIndex(2), []string{"go"})

		if len(single) != 1 || len(doubled) != 1 {
			t.Fatal("expected 1 result from each index")
		}

		// Doubling TF should increase the score,
		// but the resulting score should be less than double.

		if doubled[0].Score <= single[0].Score {
			t.Errorf("expected doubled TF to increase score")
		}
		if doubled[0].Score >= single[0].Score*2 {
			t.Errorf("expected TF saturation: doubled score %.9f should be less than %.9f",
				doubled[0].Score, single[0].Score*2)
		}
	})

	t.Run("common term stays positive", func(t *testing.T) {
		idx := &index.Index{
			Postings: map[string][]index.Posting{
				"go": {
					{DocID: 1, TF: 1},
					{DocID: 2, TF: 1},
				},
			},
			DocLen: map[int64]int{
				1: 10,
				2: 10,
			},
		}

		results := BM25(idx, []string{"go"})

		// The term appears in every document.
		// Its IDF should still be positive with BM25's IDF formula.

		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}

		for _, result := range results {
			if result.Score <= 0 {
				t.Errorf("expected positive score for DocID %d, got %.9f",
					result.DocID, result.Score)
			}
		}
	})

	t.Run("empty cases", func(t *testing.T) {
		t.Run("no matching terms", func(t *testing.T) {
			idx := &index.Index{
				Postings: map[string][]index.Posting{
					"go": {
						{DocID: 1, TF: 1},
					},
				},
				DocLen: map[int64]int{
					1: 10,
				},
			}

			results := BM25(idx, []string{"rust"})

			if len(results) != 0 {
				t.Errorf("expected no results, got %d", len(results))
			}
		})

		t.Run("empty index", func(t *testing.T) {
			idx := &index.Index{
				Postings: map[string][]index.Posting{},
				DocLen:   map[int64]int{},
			}

			results := BM25(idx, []string{"go"})

			if len(results) != 0 {
				t.Errorf("expected no results, got %d", len(results))
			}
		})
	})
}
