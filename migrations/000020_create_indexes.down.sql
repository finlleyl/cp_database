-- ============================================================================
-- Удаление индексов
-- ============================================================================

-- accounts
DROP INDEX IF EXISTS idx_accounts_user_id_created_at;
DROP INDEX IF EXISTS idx_accounts_account_type;

-- strategies
DROP INDEX IF EXISTS idx_strategies_master_account_id;
DROP INDEX IF EXISTS idx_strategies_status;
DROP INDEX IF EXISTS idx_strategies_master_user_id;

-- offers
DROP INDEX IF EXISTS idx_offers_strategy_id_created_at;
DROP INDEX IF EXISTS idx_offers_strategy_id_status;

-- subscriptions
DROP INDEX IF EXISTS idx_subscriptions_offer_id_created_at;
DROP INDEX IF EXISTS idx_subscriptions_investor_user_id;
DROP INDEX IF EXISTS idx_subscriptions_status;
DROP INDEX IF EXISTS idx_subscriptions_offer_id_status;
DROP INDEX IF EXISTS idx_subscriptions_investor_account_id;

-- trades
DROP INDEX IF EXISTS idx_trades_strategy_id_open_time;
DROP INDEX IF EXISTS idx_trades_master_account_id;
DROP INDEX IF EXISTS idx_trades_open_time;
DROP INDEX IF EXISTS idx_trades_symbol;

-- copied_trades
DROP INDEX IF EXISTS idx_copied_trades_subscription_id_open_time;
DROP INDEX IF EXISTS idx_copied_trades_trade_id_created_at;
DROP INDEX IF EXISTS idx_copied_trades_investor_account_id;

-- audit_log
DROP INDEX IF EXISTS idx_audit_log_entity_changed_at;
DROP INDEX IF EXISTS idx_audit_log_changed_at;
DROP INDEX IF EXISTS idx_audit_log_changed_by;
DROP INDEX IF EXISTS idx_audit_log_operation;
DROP INDEX IF EXISTS idx_audit_log_entity_name;

-- commissions
DROP INDEX IF EXISTS idx_commissions_subscription_id_created_at;
DROP INDEX IF EXISTS idx_commissions_created_at;
DROP INDEX IF EXISTS idx_commissions_type;

-- import_jobs
DROP INDEX IF EXISTS idx_import_jobs_type;
DROP INDEX IF EXISTS idx_import_jobs_status;
DROP INDEX IF EXISTS idx_import_jobs_created_at;

-- import_job_errors
DROP INDEX IF EXISTS idx_import_job_errors_job_id_row_number;

-- users
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_created_at;

-- favorite_strategies
DROP INDEX IF EXISTS idx_favorite_strategies_strategy_id;

-- strategy_stats
DROP INDEX IF EXISTS idx_strategy_stats_total_profit;
DROP INDEX IF EXISTS idx_strategy_stats_active_subscriptions;

