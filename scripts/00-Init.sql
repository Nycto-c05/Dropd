CREATE TABLE paste_metadata (
    id TEXT PRIMARY KEY,
    idempotency_key TEXT NOT NULL UNIQUE,
    filename TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_paste_created_at
ON paste_metadata (created_at DESC);

CREATE INDEX idx_paste_expires_at
ON paste_metadata (expires_at);
