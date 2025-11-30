CREATE TABLE offers (
    id                       BIGSERIAL PRIMARY KEY,
    strategy_id              BIGINT NOT NULL,
    name                     TEXT NOT NULL,
    status                   offer_status NOT NULL DEFAULT 'active',
    performance_fee_percent  NUMERIC(5,2),   -- %, например 20.00
    management_fee_percent   NUMERIC(5,2),
    registration_fee_amount  NUMERIC(10,2),
    created_at               TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at               TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_offers_strategy
        FOREIGN KEY (strategy_id)
        REFERENCES strategies (id)
        ON UPDATE CASCADE
        ON DELETE CASCADE
);