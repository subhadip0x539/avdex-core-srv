CREATE TABLE IF NOT EXISTS artifacts (
    id TEXT,
    name TEXT NOT NULL,
    image TEXT,
    reference TEXT,
    dex_id TEXT NOT NULL,
    extra JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_artifacts PRIMARY KEY (id)
);