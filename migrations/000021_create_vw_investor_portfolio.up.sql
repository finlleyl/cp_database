-- Создание представления vw_investor_portfolio
-- Это представление используется функцией fn_get_investor_portfolio (миграция 000017)
-- 
-- ВАЖНО: Для новой базы данных рекомендуется применить эту миграцию ДО миграции 000017
-- Для существующей базы: после применения этой миграции функция из 000017 автоматически заработает
-- (функция использует CREATE OR REPLACE, но для работы требуется существующее представление)

CREATE VIEW vw_investor_portfolio AS
SELECT
    sub.id AS subscription_id,
    sub.investor_user_id,
    s.id AS strategy_id,
    s.title AS strategy_title,
    COALESCE(SUM(ct.profit), 0) AS total_profit,
    COUNT(ct.id) AS copied_trades_count
FROM subscriptions sub
JOIN offers o ON o.id = sub.offer_id
JOIN strategies s ON s.id = o.strategy_id
LEFT JOIN copied_trades ct ON ct.subscription_id = sub.id
GROUP BY sub.id, sub.investor_user_id, s.id, s.title;

