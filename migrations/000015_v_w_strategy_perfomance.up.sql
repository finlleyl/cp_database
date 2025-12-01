CREATE VIEW vw_strategy_performance AS
SELECT
    s.id              AS strategy_id,
    s.title,
    s.status,
    ss.total_subscriptions,
    ss.active_subscriptions,
    ss.total_copied_trades,
    ss.total_profit,
    ss.total_commissions,
    ss.updated_at
FROM strategies s
LEFT JOIN strategy_stats ss ON ss.strategy_id = s.id;