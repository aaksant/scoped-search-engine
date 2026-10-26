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
	headings := strings.Join(d.Headings, ",")
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
