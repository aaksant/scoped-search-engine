CREATE TABLE IF NOT EXISTS documents (
    id INTEGER PRIMARY KEY,
    url TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL DEFAULT '', -- default to '' to avoid sql.NullString unwrapping
    headings TEXT NOT NULL DEFAULT '',
    body TEXT NOT NULL DEFAULT '',
    crawled_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
