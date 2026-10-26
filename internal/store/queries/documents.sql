-- name: UpsertDocument :exec
INSERT INTO documents (url, title, headings, body)
VALUES (?, ?, ?, ?)
-- Same URL will be crawled more than once,
-- each time database has to handle a row that may already exist.
ON CONFLICT(url) DO UPDATE SET
    title = excluded.title,
    headings = excluded.headings,
    body = excluded.body,
    crawled_at = CURRENT_TIMESTAMP;
