package store

import (
	"context"
	"database/sql"
	"strings"

	"github.com/aaksant/scoped-search-engine/internal/store/sqlc"
)

type SQLite struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewSQLite(conn *sql.DB) *SQLite {
	return &SQLite{
		db:      conn,
		queries: sqlc.New(conn),
	}
}

func (store *SQLite) UpsertDocument(
	ctx context.Context,
	d Document,
) error {
	headings := strings.Join(d.Headings, "\n")
	err := store.queries.UpsertDocument(ctx, sqlc.UpsertDocumentParams{
		Url:      d.URL,
		Title:    d.Title,
		Headings: headings,
		Body:     d.Body,
	})
	if err != nil {
		return err
	}
	return nil
}

func (store *SQLite) ListDocuments(ctx context.Context) ([]Document, error) {
	raw, err := store.queries.ListDocuments(ctx)
	if err != nil {
		return []Document{}, err
	}

	documents := make([]Document, 0, len(raw))
	for _, doc := range raw {
		documents = append(documents, toDocument(doc))
	}
	return documents, nil
}

func toDocument(d sqlc.Document) Document {
	var headings []string
	if d.Headings != "" {
		headings = strings.Split(d.Headings, "\n")
	}

	return Document{
		ID:       d.ID,
		URL:      d.Url,
		Title:    d.Title,
		Headings: headings,
		Body:     d.Body,
	}
}
