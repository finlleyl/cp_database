CREATE TABLE audit_log (
    id          BIGSERIAL PRIMARY KEY,
    entity_name TEXT NOT NULL,
    entity_pk   TEXT NOT NULL,
    operation   audit_operation NOT NULL,
    changed_by  BIGINT,
    changed_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    old_row     JSONB,
    new_row     JSONB,

    CONSTRAINT fk_audit_log_user
        FOREIGN KEY (changed_by)
        REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);