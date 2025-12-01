CREATE TABLE strategy_stats (
    strategy_id           BIGINT PRIMARY KEY,
    total_subscriptions   INTEGER NOT NULL DEFAULT 0,
    active_subscriptions  INTEGER NOT NULL DEFAULT 0,
    total_copied_trades   INTEGER NOT NULL DEFAULT 0,
    total_profit          NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_commissions     NUMERIC(18,2) NOT NULL DEFAULT 0,
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_strategy_stats_strategy
        FOREIGN KEY (strategy_id)
        REFERENCES strategies (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);