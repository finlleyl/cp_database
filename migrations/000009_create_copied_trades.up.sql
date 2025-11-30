CREATE TABLE copied_trades (
    id                  BIGSERIAL PRIMARY KEY,
    trade_id            BIGINT NOT NULL,
    subscription_id     BIGINT NOT NULL,
    investor_account_id BIGINT NOT NULL,
    volume_lots         NUMERIC(12,4) NOT NULL CHECK (volume_lots > 0),
    profit              NUMERIC(18,2),
    commission          NUMERIC(18,2),
    swap                NUMERIC(18,2),
    open_time           TIMESTAMPTZ NOT NULL,
    close_time          TIMESTAMPTZ,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_copied_trades_trade
        FOREIGN KEY (trade_id)
        REFERENCES trades (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_copied_trades_subscription
        FOREIGN KEY (subscription_id)
        REFERENCES subscriptions (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_copied_trades_investor_account
        FOREIGN KEY (investor_account_id)
        REFERENCES accounts (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);