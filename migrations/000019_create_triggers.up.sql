-- 3.1. Функция пересчёта агрегатов для одной стратегии
CREATE OR REPLACE FUNCTION fn_refresh_strategy_stats(p_strategy_id BIGINT)
RETURNS VOID AS $$
DECLARE
    v_total_subscriptions   INTEGER;
    v_active_subscriptions  INTEGER;
    v_total_copied_trades   INTEGER;
    v_total_profit          NUMERIC(18,2);
    v_total_commissions     NUMERIC(18,2);
BEGIN
    -- Подписки по стратегии (через offers)
    SELECT
        COUNT(*)::INT,
        COUNT(*) FILTER (WHERE sub.status = 'active')::INT
    INTO
        v_total_subscriptions,
        v_active_subscriptions
    FROM subscriptions sub
    JOIN offers o ON o.id = sub.offer_id
    WHERE o.strategy_id = p_strategy_id;

    -- Скопированные сделки по стратегии (через subscriptions -> offers)
    SELECT
        COUNT(ct.id)::INT,
        COALESCE(SUM(ct.profit), 0)
    INTO
        v_total_copied_trades,
        v_total_profit
    FROM copied_trades ct
    JOIN subscriptions sub ON sub.id = ct.subscription_id
    JOIN offers o ON o.id = sub.offer_id
    WHERE o.strategy_id = p_strategy_id;

    -- Комиссии по стратегии (через subscriptions -> offers)
    SELECT COALESCE(SUM(c.amount), 0)
    INTO v_total_commissions
    FROM commissions c
    JOIN subscriptions sub ON sub.id = c.subscription_id
    JOIN offers o ON o.id = sub.offer_id
    WHERE o.strategy_id = p_strategy_id;

    -- Обновляем или создаём строку в strategy_stats
    INSERT INTO strategy_stats (
        strategy_id,
        total_subscriptions,
        active_subscriptions,
        total_copied_trades,
        total_profit,
        total_commissions,
        updated_at
    )
    VALUES (
        p_strategy_id,
        COALESCE(v_total_subscriptions, 0),
        COALESCE(v_active_subscriptions, 0),
        COALESCE(v_total_copied_trades, 0),
        COALESCE(v_total_profit, 0),
        COALESCE(v_total_commissions, 0),
        now()
    )
    ON CONFLICT (strategy_id) DO UPDATE
    SET
        total_subscriptions  = EXCLUDED.total_subscriptions,
        active_subscriptions = EXCLUDED.active_subscriptions,
        total_copied_trades  = EXCLUDED.total_copied_trades,
        total_profit         = EXCLUDED.total_profit,
        total_commissions    = EXCLUDED.total_commissions,
        updated_at           = EXCLUDED.updated_at;
END;
$$ LANGUAGE plpgsql;


-- Вспомогательные функции для триггеров:
-- 3.2. Получить strategy_id по subscription_id
CREATE OR REPLACE FUNCTION fn_get_strategy_id_by_subscription(p_subscription_id BIGINT)
RETURNS BIGINT AS $$
DECLARE
    v_strategy_id BIGINT;
BEGIN
    SELECT o.strategy_id
    INTO v_strategy_id
    FROM subscriptions sub
    JOIN offers o ON o.id = sub.offer_id
    WHERE sub.id = p_subscription_id;

    RETURN v_strategy_id;
END;
$$ LANGUAGE plpgsql STABLE;

-- 3.3. Триггер-функция для subscriptions
CREATE OR REPLACE FUNCTION trg_subscriptions_refresh_stats()
RETURNS TRIGGER AS $$
DECLARE
    v_strategy_id BIGINT;
    v_offer_id BIGINT;
BEGIN
    -- Получаем offer_id из OLD или NEW в зависимости от операции
    IF TG_OP = 'DELETE' THEN
        v_offer_id := OLD.offer_id;
    ELSE
        v_offer_id := NEW.offer_id;
    END IF;

    -- Получаем strategy_id напрямую из offers (работает и при DELETE)
    SELECT strategy_id INTO v_strategy_id
    FROM offers
    WHERE id = v_offer_id;

    IF v_strategy_id IS NOT NULL THEN
        PERFORM fn_refresh_strategy_stats(v_strategy_id);
    END IF;

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- 3.4. Триггер-функция для copied_trades
CREATE OR REPLACE FUNCTION trg_copied_trades_refresh_stats()
RETURNS TRIGGER AS $$
DECLARE
    v_strategy_id BIGINT;
    v_subscription_id BIGINT;
BEGIN
    IF TG_OP = 'DELETE' THEN
        v_subscription_id := OLD.subscription_id;
    ELSE
        v_subscription_id := NEW.subscription_id;
    END IF;

    v_strategy_id := fn_get_strategy_id_by_subscription(v_subscription_id);

    IF v_strategy_id IS NOT NULL THEN
        PERFORM fn_refresh_strategy_stats(v_strategy_id);
    END IF;

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

-- 3.5. Триггер-функция для commissions
CREATE OR REPLACE FUNCTION trg_commissions_refresh_stats()
RETURNS TRIGGER AS $$
DECLARE
    v_strategy_id BIGINT;
    v_subscription_id BIGINT;
BEGIN
    IF TG_OP = 'DELETE' THEN
        v_subscription_id := OLD.subscription_id;
    ELSE
        v_subscription_id := NEW.subscription_id;
    END IF;

    v_strategy_id := fn_get_strategy_id_by_subscription(v_subscription_id);

    IF v_strategy_id IS NOT NULL THEN
        PERFORM fn_refresh_strategy_stats(v_strategy_id);
    END IF;

    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;

------------------------------------------------------------
-- 3.6. Создание триггеров для статистики стратегий
------------------------------------------------------------

CREATE TRIGGER subscriptions_refresh_stats_trg
AFTER INSERT OR UPDATE OR DELETE ON subscriptions
FOR EACH ROW EXECUTE FUNCTION trg_subscriptions_refresh_stats();

CREATE TRIGGER copied_trades_refresh_stats_trg
AFTER INSERT OR UPDATE OR DELETE ON copied_trades
FOR EACH ROW EXECUTE FUNCTION trg_copied_trades_refresh_stats();

CREATE TRIGGER commissions_refresh_stats_trg
AFTER INSERT OR UPDATE OR DELETE ON commissions
FOR EACH ROW EXECUTE FUNCTION trg_commissions_refresh_stats();


------------------------------------------------------------
-- 4. АУДИТ-ТРИГГЕРЫ
-- Автоматическое логирование INSERT/UPDATE/DELETE операций
------------------------------------------------------------

-- 4.1. Универсальная функция аудита
CREATE OR REPLACE FUNCTION fn_audit_trigger()
RETURNS TRIGGER AS $$
DECLARE
    v_old_row JSONB;
    v_new_row JSONB;
    v_entity_pk TEXT;
    v_operation audit_operation;
BEGIN
    -- Определяем операцию
    IF TG_OP = 'INSERT' THEN
        v_operation := 'insert';
        v_new_row := to_jsonb(NEW);
        v_old_row := NULL;
        v_entity_pk := NEW.id::TEXT;
    ELSIF TG_OP = 'UPDATE' THEN
        v_operation := 'update';
        v_old_row := to_jsonb(OLD);
        v_new_row := to_jsonb(NEW);
        v_entity_pk := NEW.id::TEXT;
    ELSIF TG_OP = 'DELETE' THEN
        v_operation := 'delete';
        v_old_row := to_jsonb(OLD);
        v_new_row := NULL;
        v_entity_pk := OLD.id::TEXT;
    END IF;

    -- Записываем в audit_log
    INSERT INTO audit_log (
        entity_name,
        entity_pk,
        operation,
        changed_by,
        changed_at,
        old_row,
        new_row
    ) VALUES (
        TG_TABLE_NAME,
        v_entity_pk,
        v_operation,
        current_setting('app.current_user_id', true)::BIGINT,
        now(),
        v_old_row,
        v_new_row
    );

    -- Возвращаем соответствующую запись
    IF TG_OP = 'DELETE' THEN
        RETURN OLD;
    ELSE
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;


------------------------------------------------------------
-- 4.2. Создание аудит-триггеров для ключевых таблиц
------------------------------------------------------------

-- Триггер аудита для users
CREATE TRIGGER users_audit_trg
AFTER INSERT OR UPDATE OR DELETE ON users
FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger();

-- Триггер аудита для accounts
CREATE TRIGGER accounts_audit_trg
AFTER INSERT OR UPDATE OR DELETE ON accounts
FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger();

-- Триггер аудита для strategies
CREATE TRIGGER strategies_audit_trg
AFTER INSERT OR UPDATE OR DELETE ON strategies
FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger();

-- Триггер аудита для offers
CREATE TRIGGER offers_audit_trg
AFTER INSERT OR UPDATE OR DELETE ON offers
FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger();

-- Триггер аудита для subscriptions
CREATE TRIGGER subscriptions_audit_trg
AFTER INSERT OR UPDATE OR DELETE ON subscriptions
FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger();

-- Триггер аудита для trades
CREATE TRIGGER trades_audit_trg
AFTER INSERT OR UPDATE OR DELETE ON trades
FOR EACH ROW EXECUTE FUNCTION fn_audit_trigger();
