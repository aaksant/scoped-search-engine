package store

import "context"

type Document struct {
	ID       int64
	URL      string
	Title    string
	Headings []string
	Body     string
}

type Store interface {
	UpsertDocument(ctx context.Context, d Document) error
	ListDocuments(ctx context.Context) ([]Document, error)
}
