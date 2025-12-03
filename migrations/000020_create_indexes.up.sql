-- ============================================================================
-- Индексы для оптимизации запросов
-- ============================================================================

-- ============================================================================
-- Таблица: accounts
-- ============================================================================

-- Для GetByUserID: WHERE user_id = $1 ORDER BY created_at DESC
-- Для List: WHERE user_id = $1 ORDER BY created_at DESC
CREATE INDEX idx_accounts_user_id_created_at ON accounts(user_id, created_at DESC);

-- Для List: WHERE account_type = $1
CREATE INDEX idx_accounts_account_type ON accounts(account_type);

-- ============================================================================
-- Таблица: strategies
-- ============================================================================

-- Для GetByAccountID: WHERE master_account_id = $1 ORDER BY created_at DESC
CREATE INDEX idx_strategies_master_account_id ON strategies(master_account_id);

-- Для GetActiveByID и List: WHERE status = ...
CREATE INDEX idx_strategies_status ON strategies(status);

-- Для запросов по master_user_id (связка с users)
CREATE INDEX idx_strategies_master_user_id ON strategies(master_user_id);

-- ============================================================================
-- Таблица: offers
-- ============================================================================

-- Для GetByStrategyID: WHERE strategy_id = $1 ORDER BY created_at DESC
CREATE INDEX idx_offers_strategy_id_created_at ON offers(strategy_id, created_at DESC);

-- Для GetActiveByStrategyID и List: WHERE strategy_id = $1 AND status = ...
CREATE INDEX idx_offers_strategy_id_status ON offers(strategy_id, status);

-- ============================================================================
-- Таблица: subscriptions
-- ============================================================================

-- Для GetByOfferID: WHERE offer_id = $1 ORDER BY created_at DESC
CREATE INDEX idx_subscriptions_offer_id_created_at ON subscriptions(offer_id, created_at DESC);

-- Для List: WHERE investor_user_id = $1
CREATE INDEX idx_subscriptions_investor_user_id ON subscriptions(investor_user_id);

-- Для List и GetActiveByStrategyID: WHERE status = ...
CREATE INDEX idx_subscriptions_status ON subscriptions(status);

-- Для GetActiveByStrategyID и ArchiveByStrategyID (покрывающий индекс)
CREATE INDEX idx_subscriptions_offer_id_status ON subscriptions(offer_id, status);

-- Для запросов по investor_account_id (связка с accounts)
CREATE INDEX idx_subscriptions_investor_account_id ON subscriptions(investor_account_id);

-- ============================================================================
-- Таблица: trades
-- ============================================================================

-- Для GetByStrategyID и List: WHERE strategy_id = $1 AND open_time >= ... ORDER BY open_time DESC
CREATE INDEX idx_trades_strategy_id_open_time ON trades(strategy_id, open_time DESC);

-- Для проверки связей с master_account_id
CREATE INDEX idx_trades_master_account_id ON trades(master_account_id);

-- Для фильтрации по времени открытия
CREATE INDEX idx_trades_open_time ON trades(open_time DESC);

-- Для аналитики по символам
CREATE INDEX idx_trades_symbol ON trades(symbol);

-- ============================================================================
-- Таблица: copied_trades
-- ============================================================================

-- Для GetBySubscriptionID: WHERE subscription_id = $1 ORDER BY open_time DESC
CREATE INDEX idx_copied_trades_subscription_id_open_time ON copied_trades(subscription_id, open_time DESC);

-- Для GetByTradeID: WHERE trade_id = $1 ORDER BY created_at DESC
CREATE INDEX idx_copied_trades_trade_id_created_at ON copied_trades(trade_id, created_at DESC);

-- Для запросов по investor_account_id
CREATE INDEX idx_copied_trades_investor_account_id ON copied_trades(investor_account_id);

-- ============================================================================
-- Таблица: audit_log
-- ============================================================================

-- Для GetByEntity: WHERE entity_name = $1 AND entity_pk = $2 ORDER BY changed_at DESC
CREATE INDEX idx_audit_log_entity_changed_at ON audit_log(entity_name, entity_pk, changed_at DESC);

-- Для List: фильтрация по changed_at (диапазонные запросы)
CREATE INDEX idx_audit_log_changed_at ON audit_log(changed_at DESC);

-- Для List: WHERE changed_by = $1
CREATE INDEX idx_audit_log_changed_by ON audit_log(changed_by) WHERE changed_by IS NOT NULL;

-- Для List и GetStats: WHERE operation = ...
CREATE INDEX idx_audit_log_operation ON audit_log(operation);

-- Для CountByEntity: WHERE entity_name = $1
CREATE INDEX idx_audit_log_entity_name ON audit_log(entity_name);

-- ============================================================================
-- Таблица: commissions
-- ============================================================================

-- Для GetCommissionsBySubscriptionID: WHERE subscription_id = $1 ORDER BY created_at DESC
CREATE INDEX idx_commissions_subscription_id_created_at ON commissions(subscription_id, created_at DESC);

-- Для GetMasterIncome: фильтрация по created_at
CREATE INDEX idx_commissions_created_at ON commissions(created_at);

-- Для GetMasterIncome: группировка/фильтрация по type
CREATE INDEX idx_commissions_type ON commissions(type);

-- ============================================================================
-- Таблица: import_jobs
-- ============================================================================

-- Для ListJobs: WHERE type = $1
CREATE INDEX idx_import_jobs_type ON import_jobs(type);

-- Для ListJobs: WHERE status = $1
CREATE INDEX idx_import_jobs_status ON import_jobs(status);

-- Для ListJobs: ORDER BY created_at DESC
CREATE INDEX idx_import_jobs_created_at ON import_jobs(created_at DESC);

-- ============================================================================
-- Таблица: import_job_errors
-- ============================================================================

-- Для GetJobErrors и CountJobErrors: WHERE job_id = $1 ORDER BY row_number ASC
CREATE INDEX idx_import_job_errors_job_id_row_number ON import_job_errors(job_id, row_number ASC);

-- ============================================================================
-- Таблица: users
-- ============================================================================

-- Для List: WHERE role = $1
CREATE INDEX idx_users_role ON users(role);

-- Для List: ORDER BY created_at DESC
CREATE INDEX idx_users_created_at ON users(created_at DESC);

-- ============================================================================
-- Таблица: favorite_strategies
-- ============================================================================

-- Для обратных запросов по strategy_id (PK уже покрывает user_id + strategy_id)
CREATE INDEX idx_favorite_strategies_strategy_id ON favorite_strategies(strategy_id);

-- ============================================================================
-- Таблица: strategy_stats
-- ============================================================================

-- Для сортировки в представлении vw_strategy_performance и лидерборде
CREATE INDEX idx_strategy_stats_total_profit ON strategy_stats(total_profit DESC);

-- Для фильтрации по активным подпискам
CREATE INDEX idx_strategy_stats_active_subscriptions ON strategy_stats(active_subscriptions DESC);

