CREATE TABLE commissions (
    id              BIGSERIAL PRIMARY KEY,
    subscription_id BIGINT NOT NULL,
    type            commission_type NOT NULL,
    amount          NUMERIC(18,2) NOT NULL CHECK (amount >= 0),
    period_from     TIMESTAMPTZ,
    period_to       TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_commissions_subscription
        FOREIGN KEY (subscription_id)
        REFERENCES subscriptions (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);