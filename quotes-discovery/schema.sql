CREATE TABLE IF NOT EXISTS quotes (
    id String,
    quote String,
    author String,
    created_at DateTime DEFAULT now(),
    updated_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (id)
PRIMARY KEY (id);