Сервис копирования сделок (Copy Trading Platform).

## API Документация

**Локально:** После запуска приложения Swagger UI доступен по адресу: http://localhost:8080/swagger/index.html

**Online:** [Открыть Swagger UI](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/finlleyl/cp_database/main/docs/swagger.json)


## Схема базы данных

### Enum типы

| Тип | Значения | Описание |
|-----|----------|----------|
| `strategy_status` | preparing, active, archived, deleted | Статус стратегии |
| `offer_status` | active, archived, deleted | Статус оффера |
| `subscription_status` | preparing, active, archived, deleted, suspended | Статус подписки |
| `trade_direction` | buy, sell | Направление сделки |
| `commission_type` | performance, management, registration | Тип комиссии |
| `import_job_type` | trades, accounts, statistics | Тип импорта |
| `import_job_status` | pending, running, success, failed | Статус задачи импорта |
| `audit_operation` | insert, update, delete | Тип операции аудита |

### ER-диаграмма

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│   users     │───┬───│  accounts   │───────│ strategies  │
└─────────────┘   │   └─────────────┘       └──────┬──────┘
                  │                                │
                  │   ┌─────────────┐              │
                  └───│ favorite_   │──────────────┤
                      │ strategies  │              │
                      └─────────────┘              │
                                            ┌──────┴──────┐
                                            │   offers    │
                                            └──────┬──────┘
                                                   │
┌─────────────┐       ┌─────────────┐       ┌──────┴──────┐
│   trades    │───────│copied_trades│───────│subscriptions│
└─────────────┘       └──────┬──────┘       └──────┬──────┘
                             │                     │
                             │              ┌──────┴──────┐
                             └──────────────│ commissions │
                                            └─────────────┘
```

### Таблицы

#### users
Пользователи системы (мастера и инвесторы).

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | BIGSERIAL | PK |
| email | TEXT | Уникальный email |
| name | TEXT | Имя пользователя |
| role | TEXT | Роль: master, investor, both |
| created_at | TIMESTAMPTZ | Дата создания |
| updated_at | TIMESTAMPTZ | Дата обновления |

#### accounts
Торговые счета пользователей.

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | BIGSERIAL | PK |
| user_id | BIGINT | FK → users.id |
| name | TEXT | Название счёта |
| account_type | TEXT | Тип: master, investor |
| currency | CHAR(3) | Валюта (USD, EUR, RUB и т.д.) |
| created_at | TIMESTAMPTZ | Дата создания |
| updated_at | TIMESTAMPTZ | Дата обновления |

#### strategies
Торговые стратегии мастеров.

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | BIGSERIAL | PK |
| master_user_id | BIGINT | FK → users.id (владелец) |
| master_account_id | BIGINT | FK → accounts.id (мастер-счёт) |
| title | TEXT | Название стратегии |
| description | TEXT | Описание |
| status | strategy_status | Статус стратегии |
| created_at | TIMESTAMPTZ | Дата создания |
| updated_at | TIMESTAMPTZ | Дата обновления |

#### offers
Офферы (условия подписки на стратегию).

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | BIGSERIAL | PK |
| strategy_id | BIGINT | FK → strategies.id |
| name | TEXT | Название оффера |
| status | offer_status | Статус оффера |
| performance_fee_percent | NUMERIC(5,2) | % от прибыли |
| management_fee_percent | NUMERIC(5,2) | % за управление |
| registration_fee_amount | NUMERIC(10,2) | Фикс. плата за регистрацию |
| created_at | TIMESTAMPTZ | Дата создания |
| updated_at | TIMESTAMPTZ | Дата обновления |

#### subscriptions
Подписки инвесторов на офферы.

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | BIGSERIAL | PK |
| investor_user_id | BIGINT | FK → users.id (инвестор) |
| investor_account_id | BIGINT | FK → accounts.id (счёт инвестора) |
| offer_id | BIGINT | FK → offers.id |
| status | subscription_status | Статус подписки |
| created_at | TIMESTAMPTZ | Дата создания |
| updated_at | TIMESTAMPTZ | Дата обновления |

#### favorite_strategies
Избранные стратегии пользователей.

| Колонка | Тип | Описание |
|---------|-----|----------|
| user_id | BIGINT | PK, FK → users.id |
| strategy_id | BIGINT | PK, FK → strategies.id |
| created_at | TIMESTAMPTZ | Дата добавления |

#### trades
Сделки мастера.

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | BIGSERIAL | PK |
| strategy_id | BIGINT | FK → strategies.id |
| master_account_id | BIGINT | FK → accounts.id |
| symbol | TEXT | Торговый инструмент |
| volume_lots | NUMERIC(12,4) | Объём в лотах |
| direction | trade_direction | Направление: buy/sell |
| open_time | TIMESTAMPTZ | Время открытия |
| close_time | TIMESTAMPTZ | Время закрытия |
| open_price | NUMERIC(18,6) | Цена открытия |
| close_price | NUMERIC(18,6) | Цена закрытия |
| profit | NUMERIC(18,2) | Прибыль/убыток |
| commission | NUMERIC(18,2) | Комиссия брокера |
| swap | NUMERIC(18,2) | Своп |
| created_at | TIMESTAMPTZ | Дата создания |

#### copied_trades
Скопированные сделки инвесторов.

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | BIGSERIAL | PK |
| trade_id | BIGINT | FK → trades.id (оригинальная сделка) |
| subscription_id | BIGINT | FK → subscriptions.id |
| investor_account_id | BIGINT | FK → accounts.id |
| volume_lots | NUMERIC(12,4) | Объём в лотах |
| profit | NUMERIC(18,2) | Прибыль/убыток |
| commission | NUMERIC(18,2) | Комиссия брокера |
| swap | NUMERIC(18,2) | Своп |
| open_time | TIMESTAMPTZ | Время открытия |
| close_time | TIMESTAMPTZ | Время закрытия |
| created_at | TIMESTAMPTZ | Дата создания |

#### commissions
Комиссии за использование стратегий.

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | BIGSERIAL | PK |
| subscription_id | BIGINT | FK → subscriptions.id |
| type | commission_type | Тип комиссии |
| amount | NUMERIC(18,2) | Сумма комиссии |
| period_from | TIMESTAMPTZ | Начало периода |
| period_to | TIMESTAMPTZ | Конец периода |
| created_at | TIMESTAMPTZ | Дата создания |

#### import_jobs
Задачи пакетного импорта данных.

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | BIGSERIAL | PK |
| type | import_job_type | Тип импорта |
| status | import_job_status | Статус задачи |
| file_name | TEXT | Имя файла |
| total_rows | INTEGER | Всего строк |
| processed_rows | INTEGER | Обработано строк |
| error_rows | INTEGER | Строк с ошибками |
| started_at | TIMESTAMPTZ | Время начала |
| finished_at | TIMESTAMPTZ | Время завершения |
| created_at | TIMESTAMPTZ | Дата создания |

#### import_job_errors
Ошибки импорта.

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | BIGSERIAL | PK |
| job_id | BIGINT | FK → import_jobs.id |
| row_number | INTEGER | Номер строки |
| raw_data | JSONB | Исходные данные строки |
| error_message | TEXT | Сообщение об ошибке |
| created_at | TIMESTAMPTZ | Дата создания |

#### audit_log
Журнал аудита изменений.

| Колонка | Тип | Описание |
|---------|-----|----------|
| id | BIGSERIAL | PK |
| entity_name | TEXT | Имя таблицы |
| entity_pk | TEXT | Первичный ключ записи |
| operation | audit_operation | Тип операции |
| changed_by | BIGINT | FK → users.id |
| changed_at | TIMESTAMPTZ | Время изменения |
| old_row | JSONB | Старое значение |
| new_row | JSONB | Новое значение |

#### strategy_stats
Агрегированная статистика по стратегиям (обновляется триггерами).

| Колонка | Тип | Описание |
|---------|-----|----------|
| strategy_id | BIGINT | PK, FK → strategies.id |
| total_subscriptions | INTEGER | Всего подписок |
| active_subscriptions | INTEGER | Активных подписок |
| total_copied_trades | INTEGER | Всего скопированных сделок |
| total_profit | NUMERIC(18,2) | Общая прибыль |
| total_commissions | NUMERIC(18,2) | Общие комиссии |
| updated_at | TIMESTAMPTZ | Дата обновления |

### Представления (Views)

#### vw_strategy_performance
Производительность стратегий с агрегированной статистикой.

```sql
SELECT strategy_id, title, status, total_subscriptions, 
       active_subscriptions, total_copied_trades, 
       total_profit, total_commissions, updated_at
FROM strategies s
LEFT JOIN strategy_stats ss ON ss.strategy_id = s.id
```

### Функции

| Функция | Параметры | Описание |
|---------|-----------|----------|
| `fn_get_strategy_leaderboard` | p_limit INT | Топ стратегий по прибыли |
| `fn_get_investor_portfolio` | p_investor_user_id BIGINT | Портфель инвестора |
| `fn_get_strategy_total_profit` | p_strategy_id BIGINT | Общая прибыль стратегии |
| `fn_refresh_strategy_stats` | p_strategy_id BIGINT | Пересчёт статистики стратегии |

### Триггеры

| Триггер | Таблица | Описание |
|---------|---------|----------|
| `subscriptions_refresh_stats_trg` | subscriptions | Обновление статистики при изменении подписок |
| `copied_trades_refresh_stats_trg` | copied_trades | Обновление статистики при изменении сделок |
| `commissions_refresh_stats_trg` | commissions | Обновление статистики при изменении комиссий |
| `users_audit_trg` | users | Аудит изменений пользователей |
| `accounts_audit_trg` | accounts | Аудит изменений счетов |
| `strategies_audit_trg` | strategies | Аудит изменений стратегий |
| `offers_audit_trg` | offers | Аудит изменений офферов |
| `subscriptions_audit_trg` | subscriptions | Аудит изменений подписок |
| `trades_audit_trg` | trades | Аудит изменений сделок |

## Запуск

```bash
make docker-up
make migrate-up
make seed
make run
```
