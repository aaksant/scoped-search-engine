package rank

import (
	"math"
	"testing"

	"github.com/aaksant/scoped-search-engine/internal/index"
)

func TestTFIDF(t *testing.T) {
	tests := []struct {
		name  string
		idx   *index.Index
		terms []string
		want  []Result
	}{
		{
			name: "common term contributes zero and rare term ranks first",
			idx: &index.Index{
				DocLen: map[int64]int{
					1: 10,
					2: 10,
					3: 10,
				},
				Postings: map[string][]index.Posting{
					"go": {
						{DocID: 1, TF: 2},
						{DocID: 2, TF: 2},
						{DocID: 3, TF: 2},
					},
					"goroutine": {
						{DocID: 2, TF: 3},
					},
				},
			},
			terms: []string{"go", "goroutine"},
			want: []Result{
				{DocID: 2, Score: 3 * math.Log(3.0/1.0)},
				{DocID: 1, Score: 0},
				{DocID: 3, Score: 0},
			},
		},
		{
			name: "equal scores are ordered by DocID ascending",
			idx: &index.Index{
				DocLen: map[int64]int{
					1: 10,
					2: 10,
					3: 10,
				},
				Postings: map[string][]index.Posting{
					"rare": {
						{DocID: 3, TF: 2},
						{DocID: 1, TF: 2},
					},
				},
			},
			terms: []string{"rare"},
			want: []Result{
				{DocID: 1, Score: 2 * math.Log(3.0/2.0)},
				{DocID: 3, Score: 2 * math.Log(3.0/2.0)},
			},
		},
		{
			name: "no matching terms returns empty results",
			idx: &index.Index{
				DocLen: map[int64]int{
					1: 10,
					2: 10,
				},
				Postings: map[string][]index.Posting{
					"go": {
						{DocID: 1, TF: 2},
					},
				},
			},
			terms: []string{"unknown"},
			want:  []Result{},
		},
		{
			name: "scores from multiple matching terms are accumulated",
			idx: &index.Index{
				DocLen: map[int64]int{
					1: 10,
					2: 10,
					3: 10,
				},
				Postings: map[string][]index.Posting{
					"a": {
						{DocID: 1, TF: 2},
						{DocID: 2, TF: 1},
					},
					"b": {
						{DocID: 1, TF: 3},
					},
				},
			},
			terms: []string{"a", "b"},
			want: []Result{
				{
					DocID: 1,
					Score: 2*math.Log(3.0/2.0) + 3*math.Log(3.0/1.0),
				},
				{
					DocID: 2,
					Score: 1 * math.Log(3.0/2.0),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := TFIDF(test.idx, test.terms)
			if len(got) != len(test.want) {
				t.Fatalf("Score() returned %d results, want %d",
					len(got), len(test.want))
			}

			for i := range test.want {
				if got[i].DocID != test.want[i].DocID {
					t.Errorf("result[%d].DocID = %d, want %d",
						i, got[i].DocID, test.want[i].DocID)
				}
				if math.Abs(got[i].Score-test.want[i].Score) > 1e-9 {
					t.Errorf("result[%d].Score = %f, want %f",
						i, got[i].Score, test.want[i].Score)
				}
			}
		})
	}
}
