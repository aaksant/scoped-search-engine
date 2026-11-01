package store

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"slices"
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

	tests := []Document{
		{
			URL:      url,
			Title:    "First",
			Headings: []string{"foo", "bar", "foobar"},
			Body:     "example body",
		},
		{
			URL:      url,
			Title:    "Second",
			Headings: []string{"foo", "bar", "foobar"},
			Body:     "example body",
		},
	}

	for _, doc := range tests {
		err := store.UpsertDocument(ctx, doc)
		if err != nil {
			t.Fatalf("error upserting document: %v", err)
		}
	}

	var count int64
	var title string

	err := db.QueryRow("SELECT COUNT(*) FROM documents").Scan(&count)
	if count != 1 {
		t.Errorf("error counting row: %v", err)
	}
	err = db.QueryRow("SELECT title FROM documents WHERE url = ?", url).Scan(&title)
	if title != "Second" {
		t.Errorf("error querying title: %v", err)
	}
}

func TestSQLiteListDocuments(t *testing.T) {
	db := newTestDB(t)
	store := NewSQLite(db)
	ctx := context.Background()

	tests := []Document{
		{
			URL:      "https://go.dev/doc/a",
			Title:    "First",
			Headings: []string{"Intro", "Setup"},
			Body:     "example body",
		},
		{
			URL:      "https://go.dev/doc/b",
			Title:    "Second",
			Headings: nil,
			Body:     "example body",
		},
		{
			URL:      "https://go.dev/doc/c",
			Title:    "Third",
			Headings: []string{"Errors, panics, and recovery"},
			Body:     "example body",
		},
	}

	for _, doc := range tests {
		err := store.UpsertDocument(ctx, doc)
		if err != nil {
			t.Fatalf("error upserting document: %v", err)
		}
	}

	docs, err := store.ListDocuments(ctx)
	if err != nil {
		t.Fatalf("error listing documents: %v", err)
	}
	if len(docs) != 3 {
		t.Fatalf("got: %v, want: %v", len(docs), 3)
	}

	// order by id asc check
	for i := 1; i < len(docs); i++ {
		current := docs[i].ID
		prev := docs[i-1].ID

		if current <= prev {
			t.Fatalf(
				"ID at position %v: %v, not greater than previous: %v",
				i,
				current,
				prev,
			)
		}
	}

	// additional test
	// round-trip check
	// what we inserted is what we got back
	inserted := make(map[string]Document)
	for _, doc := range tests {
		inserted[doc.URL] = doc
	}

	for _, got := range docs {
		want, ok := inserted[got.URL]
		if !ok {
			t.Fatalf("unexpected URL in results: %q", got.URL)
		}
		if got.Title != want.Title {
			t.Errorf("%s title: got %q, want %q", got.URL, got.Title, want.Title)
		}
		if got.Body != want.Body {
			t.Errorf("%s body: got %q, want %q", got.URL, got.Body, want.Body)
		}
		if !slices.Equal(got.Headings, want.Headings) {
			t.Errorf("%s headings: got %q, want %q", got.URL, got.Headings, want.Headings)
		}
	}
}
