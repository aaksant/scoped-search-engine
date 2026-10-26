package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/aaksant/scoped-search-engine/internal/crawler"
	"github.com/aaksant/scoped-search-engine/internal/store"
	_ "modernc.org/sqlite"
)

const maxPages = 10

func main() {
	ctx := context.Background()

	if len(os.Args) < 2 {
		fmt.Println("usage: go run ./cmd/crawl <url>")
		os.Exit(1)
	}
	url := os.Args[1]

	db, err := sql.Open("sqlite", "db/index.db")
	if err != nil {
		fmt.Printf("error opening database connection: %v\n", err)
		os.Exit(1)
	}
	sqlite := store.NewSQLite(db)

	if err := crawler.Crawl(ctx, sqlite, url, maxPages); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
