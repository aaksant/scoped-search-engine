package store

import "context"

type Document struct {
	URL      string
	Title    string
	Headings []string
	Body     string
}

type Store interface {
	UpsertDocument(ctx context.Context, d Document) error
}
