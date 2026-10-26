package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aaksant/scoped-search-engine/internal/extract"
	"github.com/aaksant/scoped-search-engine/internal/store"
)

func Crawl(
	ctx context.Context,
	sqlite *store.SQLite,
	seed string,
	maxPages int,
) error {
	var queue []string
	visited := make(map[string]bool)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	parsedSeed, err := url.Parse(seed)
	if err != nil {
		return fmt.Errorf("error parsing seed: %v", err)
	}

	resolvedSeed, ok := NormalizeLink(parsedSeed, seed)
	if !ok {
		return fmt.Errorf("error resolving seed")
	}
	if !InScope(resolvedSeed) {
		return fmt.Errorf("resolved seed is not allowed")
	}

	visited[resolvedSeed] = true
	queue = append(queue, resolvedSeed)
	pagesCrawled := 0

	for pagesCrawled < maxPages && len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		fmt.Println(current)

		resp, err := client.Get(current)
		if err != nil {
			fmt.Printf("error fetching page: %v\n", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			continue
		}
		if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
			fmt.Println("Content-Type is not text/html")
			continue
		}
		defer resp.Body.Close()

		page, err := extract.Parse(resp.Body)
		if err != nil {
			fmt.Printf("error parsing page: %v\n", err)
			continue
		}

		doc := store.Document{
			URL:      current,
			Title:    page.Title,
			Headings: page.Headings,
			Body:     page.Body,
		}
		if err := sqlite.UpsertDocument(ctx, doc); err != nil {
			fmt.Printf("error upserting document: %v\n", err)
			continue
		}

		pagesCrawled++
		fmt.Println(page.Title)

		for _, href := range page.Links {
			normalized, ok := NormalizeLink(resp.Request.URL, href)
			if !ok {
				continue
			}
			if !InScope(normalized) {
				continue
			}
			if _, ok := visited[normalized]; ok {
				continue
			}

			visited[normalized] = true
			queue = append(queue, normalized)
		}
	}
	return nil
}
