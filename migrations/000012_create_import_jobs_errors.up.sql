
CREATE TABLE import_job_errors (
    id            BIGSERIAL PRIMARY KEY,
    job_id        BIGINT NOT NULL,
    row_number    INTEGER,
    raw_data      JSONB,
    error_message TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_job_errors_job
        FOREIGN KEY (job_id)
        REFERENCES import_jobs (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);