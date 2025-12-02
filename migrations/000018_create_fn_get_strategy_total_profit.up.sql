CREATE OR REPLACE FUNCTION fn_get_strategy_total_profit(p_strategy_id BIGINT)
RETURNS NUMERIC AS $$
DECLARE
    v_profit NUMERIC(18,2);
BEGIN
    SELECT COALESCE(SUM(ct.profit), 0)
    INTO v_profit
    FROM copied_trades ct
    JOIN subscriptions sub ON sub.id = ct.subscription_id
    JOIN offers o ON o.id = sub.offer_id
    WHERE o.strategy_id = p_strategy_id;

    RETURN v_profit;
END;
$$ LANGUAGE plpgsql STABLE;