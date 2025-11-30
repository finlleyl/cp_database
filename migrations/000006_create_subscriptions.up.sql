CREATE TABLE subscriptions (
    id                   BIGSERIAL PRIMARY KEY,
    investor_user_id     BIGINT NOT NULL,
    investor_account_id  BIGINT NOT NULL,
    offer_id             BIGINT NOT NULL,
    status               subscription_status NOT NULL DEFAULT 'preparing',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_subscriptions_investor_user
        FOREIGN KEY (investor_user_id)
        REFERENCES users (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_subscriptions_investor_account
        FOREIGN KEY (investor_account_id)
        REFERENCES accounts (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,

    CONSTRAINT fk_subscriptions_offer
        FOREIGN KEY (offer_id)
        REFERENCES offers (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);