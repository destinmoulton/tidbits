-- migrate:up
CREATE TABLE IF NOT EXISTS settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    value TEXT
);

-- migrate:down