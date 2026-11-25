
CREATE TABLE IF NOT EXISTS artifact_sites (
    id TEXT,
    name TEXT NOT NULL,
    artifact_id TEXT NOT NULL,
    url TEXT NOT NULL,

    CONSTRAINT pk_artifact_sites PRIMARY KEY (id),

    CONSTRAINT fk_artifact_sites_artifact_id FOREIGN KEY (artifact_id) REFERENCES artifacts (id) ON DELETE CASCADE,

    CONSTRAINT uq_artifact_sites_artifact_id_url UNIQUE (artifact_id, url)
);