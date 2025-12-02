CREATE OR REPLACE FUNCTION fn_get_investor_portfolio(
    p_investor_user_id BIGINT
)
RETURNS TABLE (
    subscription_id      BIGINT,
    strategy_id          BIGINT,
    strategy_title       TEXT,
    total_profit         NUMERIC(18,2),
    copied_trades_count  BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        v.subscription_id,
        v.strategy_id,
        v.strategy_title,
        v.total_profit,
        v.copied_trades_count
    FROM vw_investor_portfolio v
    WHERE v.investor_user_id = p_investor_user_id;
END;
$$ LANGUAGE plpgsql STABLE;