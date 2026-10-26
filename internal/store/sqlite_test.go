package store

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	path := filepath.Join(t.TempDir(), "test.db")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	schema, err := os.ReadFile("../../db/migrations/000001_create_documents_table.up.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	if _, err := db.Exec(string(schema)); err != nil {
		t.Fatalf("apply migration: %v", err)
	}

	return db
}

func TestSQLiteUpsertDocument(t *testing.T) {
	db := newTestDB(t)
	store := NewSQLite(db)
	ctx := context.Background()

	url := "https://go.dev/doc"

	err := store.UpsertDocument(ctx, Document{
		URL:      url,
		Title:    "First",
		Headings: []string{"foo", "bar", "foobar"},
		Body:     "example body",
	})
	if err != nil {
		t.Fatalf("error upserting document: %v", err)
	}

	err = store.UpsertDocument(ctx, Document{
		URL:      "https://go.dev/doc",
		Title:    "Second",
		Headings: []string{"foo", "bar", "foobar"},
		Body:     "example body",
	})
	if err != nil {
		t.Fatalf("error upserting document: %v", err)
	}

	var count int64
	var title string

	err = db.QueryRow("SELECT COUNT(*) FROM documents").Scan(&count)
	if count != 1 {
		t.Errorf("error counting row: %v", err)
	}

	err = db.QueryRow("SELECT title FROM documents WHERE url = ?", url).Scan(&title)
	if title != "Second" {
		t.Errorf("error querying title: %v", err)
	}
}
