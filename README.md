
# Scoped Search Engine
A search engine built from scratch in Go that crawls web pages within a predefined scope, stores their content in SQLite, and makes the content searchable through an inverted index. 

The project serves as a hands-on exploration of web crawling, information retrieval, indexing, and search ranking.

## Progress 

- [x] extract title, headings, and body
- [x] normalize URLs (mixed absolute and relative, anchors)
- [x] crawl pages from seed URL
- [x] integrate database (`documents` table)
- [x] analyze page content (tokenize, stopwords removal)
- [ ] in-memory inverted index (term → postings)     
- [ ] load documents from the store and build the index (`cmd/index`)
- [ ] save/load an index snapshot
- [ ] query CLI (`cmd/search`): analyze the query, look up terms
- [ ] ranking: TF-IDF, then BM25

## Planned Architecture
```
go-search-engine/
  cmd/
    crawl/      runs crawler, fills the store
    index/      builds index snapshot from the store
    search/     CLI query tool
    server/     HTTP API (later)
  internal/
    crawler/    frontier, fetcher, robots.txt, per-domain rate limit, worker pool
    extract/    HTML -> `Document{URL, Title, Headings, Body}`
    analyze/    lowercase, tokenize, stop words, (stemming)
    index/      inverted index, postings, snapshot save/load
    rank/       TF-IDF, then BM25
    store/      Store interface + SQLite implementation
    api/        handlers, pagination, rate limiting, cache
```

## `internal/crawler/`

The crawler package crawls web pages within a predefined scope, extracts their content, and stores the resulting documents in SQLite. It uses a queue-based approach to discover and visit pages, starting from a seed URL.

The crawler only processes URLs on go.dev under `/doc`, `/blog`, and `/ref` prefixes. The scope check accepts both the prefix itself and its subpaths. For example, `/ref` and `/ref/spec` are allowed, but `/reference` is not.

**WIP: implement concurrent workers** 

### How it works

1.  **Initialize the crawl.** Parse and normalize the seed URL, then verify that it falls within the allowed scope.
2.  **Process the queue.** Use a FIFO queue to visit discovered URLs in breadth-first order.
3.  **Fetch pages.** Send HTTP GET requests and skip pages that cannot be fetched successfully or do not return a `200 OK` response.
4.  **Extract content.** Parse HTML pages to obtain their titles, headings, body text, and links.
5.  **Store documents.** Save extracted content to SQLite, updating existing documents when necessary.
6.  **Discover links.** Normalize links found on each page, filter out URLs outside the allowed scope, and enqueue URLs that have not been discovered before.
7.  **Stop crawling.** Finish when the page limit is reached or no URLs remain in the queue.

## `internal/extract/`

The extract package parses HTML documents and extracts structured content for the search engine. It uses goquery to select HTML elements and returns a Page containing the page title, body text, headings, and links.

The extracted data can then be used by the crawler for storage and by subsequent components for indexing and searching.

### How it works

1.  **Parse the HTML.** Read the input through an `io.Reader` and parse it using `goquery`.
2.  **Extract the title.** Retrieve the text inside the `<title>` element and trim surrounding whitespace.
3.  **Select the content area.** Use `<main>` when available; otherwise, use `<body>`.
4.  **Remove unwanted elements.** Remove scripts, styles, navigation, footers, and other specified elements from the selected content.
5.  **Extract and clean the body text.** Convert the remaining content to text and normalize its whitespace.
6.  **Extract headings.** Collect and clean text from all `<h1>`, `<h2>`, and `<h3>` elements in the document.
7.  **Extract links.** Collect `href` values from anchor elements.
8.  **Return the result.** Return a `Page` containing the extracted information.

### Data structure
```
type Page struct {
	Title    string
	Body     string
	Headings []string
	Links    []string
}
```
|Property|Description  |
|--|--|
| `Title` | Text from the document's `<title>` element |
| `Body` | Cleaned text from the selected content area |
| `Headings` | Non-empty, cleaned text from `<h1>`, `<h2>`, and `<h3>` elements |
| `Links` | `href` attribute values collected from anchor elements. |

## `internal/analyze`
The `analyze` package transforms raw text into terms suitable for indexing. It lowercases text, splits it into tokens, and removes common words (stopwords) that are considered less useful for search. This package prepares extracted document content for the indexing stage of the search engine.

### How it works
1.  **Receive raw text.** Accept a string containing document content or search query text.
2.  **Normalize case.** Convert the entire string to lowercase.
3.  **Tokenize the text.** Split the text wherever a character is neither a Unicode letter nor a Unicode number.
4.  **Remove stopwords.** Exclude tokens that appear in the predefined stopword set.
5.  **Return the terms.** Return the remaining tokens for indexing or search processing.

### Handling stopwords
The package uses a manually curated list rather than adopting a comprehensive, general-purpose English stopword list.

The goal is to remove common words that contribute relatively little to search while preserving terminology that may be meaningful in Go documentation.

Examples of words removed include:

-   Articles: `the`, `a`, `an`
-   Conjunctions: `and`, `or`, `but`
-   Prepositions: `of`, `to`, `in`, `on`, `at`, `by`
-   Common pronouns: `this`, `that`, `you`, `we`
-   Common auxiliary verbs: `is`, `are`, `was`, `were`, `will`, `can`
    
The list also includes common English contraction fragments, such as `s`, `t`, `ll`, and `ve`.

Certain words that might appear in generic stopword lists are intentionally retained because they can be meaningful in programming documentation. Examples include `go`, `for`, `if`, `else`, `ok`, `io`, `os`, and `http`.

### Why use a set?
The stopword list is convenient to maintain as a slice, but checking every token against that slice would require scanning the list repeatedly.

Converting the list into a set allows membership checks to take expected constant time, making stopword filtering more efficient as the text grows.
