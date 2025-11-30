CREATE TABLE favorite_strategies (
    user_id     BIGINT NOT NULL,
    strategy_id BIGINT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, strategy_id),

    CONSTRAINT fk_favorites_user
        FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,

    CONSTRAINT fk_favorites_strategy
        FOREIGN KEY (strategy_id)
        REFERENCES strategies (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);