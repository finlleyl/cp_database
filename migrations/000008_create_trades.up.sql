CREATE TABLE trades (
    id               BIGSERIAL PRIMARY KEY,
    strategy_id      BIGINT NOT NULL,
    master_account_id BIGINT NOT NULL,
    symbol           TEXT NOT NULL,
    volume_lots      NUMERIC(12,4) NOT NULL CHECK (volume_lots > 0),
    direction        trade_direction NOT NULL,
    open_time        TIMESTAMPTZ NOT NULL,
    close_time       TIMESTAMPTZ,
    open_price       NUMERIC(18,6) NOT NULL,
    close_price      NUMERIC(18,6),
    profit           NUMERIC(18,2),
    commission       NUMERIC(18,2),
    swap             NUMERIC(18,2),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_trades_strategy
        FOREIGN KEY (strategy_id)
        REFERENCES strategies (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_trades_master_account
        FOREIGN KEY (master_account_id)
        REFERENCES accounts (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);