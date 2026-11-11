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

// The pause after every fetch, successful or not, so the
// crawler stays polite to the host. go.dev's robots.txt sets no Crawl-delay,
// so this value is our own choice.
const requestDelay = 750 * time.Millisecond

// Fetches pages breadth-first starting at seed and saves each one to
// the store.
//
// The seed must resolve to a URL that `InScope` accepts. From there it
// follows in-scope links, visiting each URL at most once, and stops when
// `maxPages` pages have been stored or the queue runs out. Only pages that are
// saved count towards `maxPages`, so skipped pages don't use up the budget.
//
// A page is keyed by the URL it finally landed on after redirects, not the
// URL that was requested. A redirect to a URL that was already visited, or
// to one outside the scope, is skipped. This stops different spellings of one
// page, such as "/doc/effective_go.html" and "/doc/effective_go", from being
// stored twice.
//
// Pages are skipped, and the problem is printed, when the fetch fails, the
// status is not 200, the content type is not text/html, parsing fails, or the
// store write fails.
//
// `Crawl` returns an error only when the seed is invalid or out of scope.
// Per-page failures are printed and never returned. `ctx` is passed to the
// store writes only, so cancelling it does not interrupt a fetch or stop the
// loop.
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
		time.Sleep(requestDelay)
		if err != nil {
			fmt.Printf("error fetching page: %v\n", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("skip %s: status %d\n", current, resp.StatusCode)
			resp.Body.Close()
			continue
		}

		final, ok := NormalizeLink(resp.Request.URL, resp.Request.URL.String())
		if !ok || !InScope(final) {
			resp.Body.Close()
			continue
		}
		if final != current {
			if _, ok := visited[final]; ok {
				resp.Body.Close()
				continue // redirect target already crawled
			}
			visited[final] = true
		}

		if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
			fmt.Println("Content-Type is not text/html")
			resp.Body.Close()
			continue
		}

		page, err := extract.Parse(resp.Body)
		resp.Body.Close() // defer leaks because body closed when return
		if err != nil {
			fmt.Printf("error parsing page: %v\n", err)
			continue
		}

		doc := store.Document{
			URL:      final,
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
