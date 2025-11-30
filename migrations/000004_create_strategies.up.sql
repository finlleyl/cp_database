CREATE TABLE strategies (
    id               BIGSERIAL PRIMARY KEY,
    master_user_id   BIGINT NOT NULL,
    master_account_id BIGINT NOT NULL,
    title            TEXT NOT NULL,
    description      TEXT,
    status           strategy_status NOT NULL DEFAULT 'preparing',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_strategies_master_user
        FOREIGN KEY (master_user_id)
        REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_strategies_master_account
        FOREIGN KEY (master_account_id)
        REFERENCES accounts (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);