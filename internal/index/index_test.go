package index

import (
	"reflect"
	"testing"
)

func TestIndexAdd(t *testing.T) {
	tests := []struct {
		name string
		docs []struct {
			docID  int64
			tokens []string
		}
		expectedPostings map[string][]Posting
		expectedDocLen   map[int64]int
	}{
		{
			name: "single document with unique tokens",
			docs: []struct {
				docID  int64
				tokens []string
			}{
				{docID: 1, tokens: []string{"hello", "world"}},
			},
			expectedPostings: map[string][]Posting{
				"hello": {{DocID: 1, TF: 1}},
				"world": {{DocID: 1, TF: 1}},
			},
			expectedDocLen: map[int64]int{
				1: 2,
			},
		},
		{
			name: "single document with repeated tokens",
			docs: []struct {
				docID  int64
				tokens []string
			}{
				{docID: 1, tokens: []string{"go", "test", "go", "code"}},
			},
			expectedPostings: map[string][]Posting{
				"go":   {{DocID: 1, TF: 2}},
				"test": {{DocID: 1, TF: 1}},
				"code": {{DocID: 1, TF: 1}},
			},
			expectedDocLen: map[int64]int{
				1: 4,
			},
		},
		{
			name: "multiple documents with overlapping terms",
			docs: []struct {
				docID  int64
				tokens []string
			}{
				{docID: 1, tokens: []string{"search", "engine"}},
				{docID: 2, tokens: []string{"search", "index", "engine"}},
			},
			expectedPostings: map[string][]Posting{
				"search": {
					{DocID: 1, TF: 1},
					{DocID: 2, TF: 1},
				},
				"engine": {
					{DocID: 1, TF: 1},
					{DocID: 2, TF: 1},
				},
				"index": {
					{DocID: 2, TF: 1},
				},
			},
			expectedDocLen: map[int64]int{
				1: 2,
				2: 3,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			indexer := New()
			for _, d := range test.docs {
				indexer.Add(d.docID, d.tokens)
			}
			if !reflect.DeepEqual(indexer.Postings, test.expectedPostings) {
				t.Errorf("Index.Postings: %v, want: %v", indexer.Postings, test.expectedPostings)
			}
			if !reflect.DeepEqual(indexer.DocLen, test.expectedDocLen) {
				t.Errorf("Index.DocLen: %v, want: %v", indexer.DocLen, test.expectedDocLen)
			}
		})
	}
}
