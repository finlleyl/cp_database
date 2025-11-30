CREATE TABLE import_jobs (
    id             BIGSERIAL PRIMARY KEY,
    type           import_job_type NOT NULL,
    status         import_job_status NOT NULL DEFAULT 'pending',
    file_name      TEXT,
    total_rows     INTEGER NOT NULL DEFAULT 0 CHECK (total_rows >= 0),
    processed_rows INTEGER NOT NULL DEFAULT 0 CHECK (processed_rows >= 0),
    error_rows     INTEGER NOT NULL DEFAULT 0 CHECK (error_rows >= 0),
    started_at     TIMESTAMPTZ,
    finished_at    TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);