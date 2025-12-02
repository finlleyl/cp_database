CREATE OR REPLACE FUNCTION fn_get_strategy_leaderboard(
    p_limit INT
)
RETURNS TABLE (
    strategy_id        BIGINT,
    title              TEXT,
    total_profit       NUMERIC(18,2),
    total_commissions  NUMERIC(18,2),
    active_subscriptions INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        sp.strategy_id,
        sp.title,
        sp.total_profit,
        sp.total_commissions,
        sp.active_subscriptions
    FROM vw_strategy_performance sp
    ORDER BY sp.total_profit DESC
    LIMIT p_limit;
END;
$$ LANGUAGE plpgsql STABLE;