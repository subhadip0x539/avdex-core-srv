
CREATE TABLE IF NOT EXISTS dexs (
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    color TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_dexs
        PRIMARY KEY (id),

    CONSTRAINT uq_dexs_name
        UNIQUE (name)
);
