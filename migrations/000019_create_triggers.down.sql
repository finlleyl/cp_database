-- Удаляем триггеры аудита
DROP TRIGGER IF EXISTS users_audit_trg ON users;
DROP TRIGGER IF EXISTS accounts_audit_trg ON accounts;
DROP TRIGGER IF EXISTS strategies_audit_trg ON strategies;
DROP TRIGGER IF EXISTS offers_audit_trg ON offers;
DROP TRIGGER IF EXISTS subscriptions_audit_trg ON subscriptions;
DROP TRIGGER IF EXISTS trades_audit_trg ON trades;

-- Удаляем функцию аудита
DROP FUNCTION IF EXISTS fn_audit_trigger();

-- Удаляем триггеры статистики
DROP TRIGGER IF EXISTS subscriptions_refresh_stats_trg ON subscriptions;
DROP TRIGGER IF EXISTS copied_trades_refresh_stats_trg ON copied_trades;
DROP TRIGGER IF EXISTS commissions_refresh_stats_trg ON commissions;

-- Удаляем функции статистики
DROP FUNCTION IF EXISTS trg_commissions_refresh_stats();
DROP FUNCTION IF EXISTS trg_copied_trades_refresh_stats();
DROP FUNCTION IF EXISTS trg_subscriptions_refresh_stats();
DROP FUNCTION IF EXISTS fn_get_strategy_id_by_subscription(BIGINT);
DROP FUNCTION IF EXISTS fn_refresh_strategy_stats(BIGINT);

