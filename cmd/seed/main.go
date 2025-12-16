package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type options struct {
	dsn string

	seed int64

	truncate bool

	users         int
	strategies    int
	offers        int
	subscriptions int

	trades       int
	copiedTrades int
	commissions  int

	favorites int
	audit     int

	importJobs   int
	importErrors int
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func defaultDSN() string {
	user := envOr("POSTGRES_USER", "postgres")
	pass := envOr("POSTGRES_PASSWORD", "postgres")
	db := envOr("POSTGRES_DB", "postgres")
	host := envOr("POSTGRES_HOST", "localhost")
	port := envOr("POSTGRES_PORT", "5432")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, db)
}

func setseedFromInt(seed int64) float64 {
	// pg setseed expects [-1, 1]. Keep it deterministic and "random enough".
	// Map [0..9999] -> [-1..1].
	m := float64(seed % 10000)
	f01 := m / 9999.0
	return -1.0 + 2.0*f01
}

func main() {
	var opt options

	flag.StringVar(&opt.dsn, "dsn", defaultDSN(), "PostgreSQL DSN (e.g. postgres://user:pass@localhost:5432/db?sslmode=disable)")
	flag.Int64Var(&opt.seed, "seed", 42, "Seed for deterministic pseudo-randomness")
	flag.BoolVar(&opt.truncate, "truncate", true, "TRUNCATE tables before inserting")

	flag.IntVar(&opt.users, "users", 1000, "Rows in users")
	flag.IntVar(&opt.strategies, "strategies", 1000, "Rows in strategies")
	flag.IntVar(&opt.offers, "offers", 1000, "Rows in offers")
	flag.IntVar(&opt.subscriptions, "subscriptions", 1000, "Rows in subscriptions")

	flag.IntVar(&opt.trades, "trades", 5000, "Rows in trades")
	flag.IntVar(&opt.copiedTrades, "copied-trades", 5000, "Rows in copied_trades")
	flag.IntVar(&opt.commissions, "commissions", 5000, "Rows in commissions")

	flag.IntVar(&opt.favorites, "favorites", 1000, "Attempts to insert rows in favorite_strategies (duplicates ignored)")
	flag.IntVar(&opt.audit, "audit", 5000, "Target rows in audit_log (triggers will create entries, UPDATE operations will add more if needed)")

	flag.IntVar(&opt.importJobs, "import-jobs", 100, "Rows in import_jobs")
	flag.IntVar(&opt.importErrors, "import-errors", 200, "Rows in import_job_errors")

	flag.Parse()

	if opt.users <= 0 || opt.strategies <= 0 || opt.trades <= 0 {
		fmt.Fprintln(os.Stderr, "users/strategies/trades must be > 0")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	pool, err := pgxpool.New(ctx, opt.dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pgxpool.New: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	tx, err := pool.Begin(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "begin: %v\n", err)
		os.Exit(1)
	}

	// ВАЖНО: ВСЕ выполняется в одной транзакции!
	// При любой ошибке все изменения автоматически откатываются (ROLLBACK).
	// Commit происходит только в самом конце, если все операции успешны.
	// Это означает: либо все данные вставляются, либо ничего не вставляется.
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			// Игнорируем ошибки "transaction already committed" (это нормально, если был Commit)
		}
	}()

	if _, err := tx.Exec(ctx, `SET TIME ZONE 'UTC'`); err != nil {
		fmt.Fprintf(os.Stderr, "set timezone: %v\n", err)
		os.Exit(1)
	}
	if _, err := tx.Exec(ctx, `SELECT setseed($1)`, setseedFromInt(opt.seed)); err != nil {
		fmt.Fprintf(os.Stderr, "setseed: %v\n", err)
		os.Exit(1)
	}

	if opt.truncate {
		// CASCADE is important because there are many FKs with RESTRICT.
		if _, err := tx.Exec(ctx, `
TRUNCATE TABLE
  import_job_errors,
  import_jobs,
  favorite_strategies,
  copied_trades,
  commissions,
  subscriptions,
  offers,
  trades,
  strategy_stats,
  strategies,
  accounts,
  audit_log,
  users
RESTART IDENTITY CASCADE`); err != nil {
			fmt.Fprintf(os.Stderr, "truncate: %v\n", err)
			os.Exit(1)
		}
	}

	// Disable audit triggers during seeding to avoid creating too many audit_log entries
	// We'll enable them back and use UPDATE operations to generate the desired amount
	if _, err := tx.Exec(ctx, `
ALTER TABLE users DISABLE TRIGGER users_audit_trg;
ALTER TABLE accounts DISABLE TRIGGER accounts_audit_trg;
ALTER TABLE strategies DISABLE TRIGGER strategies_audit_trg;
ALTER TABLE offers DISABLE TRIGGER offers_audit_trg;
ALTER TABLE subscriptions DISABLE TRIGGER subscriptions_audit_trg;
ALTER TABLE trades DISABLE TRIGGER trades_audit_trg;
`); err != nil {
		fmt.Fprintf(os.Stderr, "disable audit triggers: %v\n", err)
		os.Exit(1)
	}

	// users
	if _, err := tx.Exec(ctx, `
INSERT INTO users(email, name, role, created_at, updated_at)
SELECT
  'user' || gs::text || '@example.com',
  (ARRAY['Иван','Пётр','Анна','Мария','Олег','Елена','Дмитрий','Наталья','Сергей','Алексей'])[1 + floor(random()*10)::int]
    || ' ' ||
  (ARRAY['Иванов','Петров','Сидоров','Смирнов','Кузнецов','Попов','Соколов','Лебедев','Козлов','Новиков'])[1 + floor(random()*10)::int],
  CASE
    -- гарантируем, что будут и мастер, и инвестор, и both (иначе не из чего строить FK-связи)
    WHEN gs = 1 THEN 'master'
    WHEN gs = 2 THEN 'investor'
    WHEN gs = 3 THEN 'both'
    WHEN random() < 0.40 THEN 'master'
    WHEN random() < 0.80 THEN 'investor'
    ELSE 'both'
  END,
  now() - (random() * interval '365 days'),
  now()
FROM generate_series(1, $1) gs
`, opt.users); err != nil {
		fmt.Fprintf(os.Stderr, "insert users: %v\n", err)
		os.Exit(1)
	}

	// accounts (1 master account for each master/both; 1 investor account for each investor/both)
	if _, err := tx.Exec(ctx, `
INSERT INTO accounts(user_id, name, account_type, currency, created_at, updated_at)
SELECT
  u.id,
  'ACC-' || u.id::text || '-M',
  'master',
  (ARRAY['USD','EUR','RUB'])[1 + floor(random()*3)::int],
  u.created_at + (random() * interval '30 days'),
  now()
FROM users u
WHERE u.role IN ('master','both')
`); err != nil {
		fmt.Fprintf(os.Stderr, "insert master accounts: %v\n", err)
		os.Exit(1)
	}
	if _, err := tx.Exec(ctx, `
INSERT INTO accounts(user_id, name, account_type, currency, created_at, updated_at)
SELECT
  u.id,
  'ACC-' || u.id::text || '-I',
  'investor',
  (ARRAY['USD','EUR','RUB'])[1 + floor(random()*3)::int],
  u.created_at + (random() * interval '30 days'),
  now()
FROM users u
WHERE u.role IN ('investor','both')
`); err != nil {
		fmt.Fprintf(os.Stderr, "insert investor accounts: %v\n", err)
		os.Exit(1)
	}

	// strategies
	if _, err := tx.Exec(ctx, `
WITH master_accounts AS (
  SELECT a.id, a.user_id
  FROM accounts a
  WHERE a.account_type = 'master'
)
INSERT INTO strategies(master_user_id, master_account_id, title, description, status, created_at, updated_at)
SELECT
  ma.user_id,
  ma.id,
  'Strategy #' || gs::text,
  'Synthetic strategy ' || gs::text,
  (ARRAY['preparing','active','archived'])[1 + floor(random()*3)::int]::strategy_status,
  now() - (random() * interval '365 days'),
  now()
FROM generate_series(1, $1) gs
JOIN LATERAL (SELECT * FROM master_accounts ORDER BY random() LIMIT 1) ma ON true
`, opt.strategies); err != nil {
		fmt.Fprintf(os.Stderr, "insert strategies: %v\n", err)
		os.Exit(1)
	}

	// offers
	if _, err := tx.Exec(ctx, `
WITH st AS (SELECT id FROM strategies)
INSERT INTO offers(strategy_id, name, status, performance_fee_percent, management_fee_percent, registration_fee_amount, created_at, updated_at)
SELECT
  s.id,
  'Offer #' || gs::text,
  'active'::offer_status,
  round((5 + random()*25)::numeric, 2),
  round((0 + random()*5)::numeric, 2),
  round((0 + random()*50)::numeric, 2),
  now() - (random() * interval '180 days'),
  now()
FROM generate_series(1, $1) gs
JOIN LATERAL (SELECT * FROM st ORDER BY random() LIMIT 1) s ON true
`, opt.offers); err != nil {
		fmt.Fprintf(os.Stderr, "insert offers: %v\n", err)
		os.Exit(1)
	}

	// subscriptions
	if _, err := tx.Exec(ctx, `
WITH investor_accounts AS (
  SELECT a.id, a.user_id
  FROM accounts a
  WHERE a.account_type = 'investor'
),
offs AS (
  SELECT id FROM offers
)
INSERT INTO subscriptions(investor_user_id, investor_account_id, offer_id, status, created_at, updated_at)
SELECT
  ia.user_id,
  ia.id,
  o.id,
  (ARRAY['preparing','active','suspended','archived'])[1 + floor(random()*4)::int]::subscription_status,
  now() - (random() * interval '180 days'),
  now()
FROM generate_series(1, $1) gs
JOIN LATERAL (SELECT * FROM investor_accounts ORDER BY random() LIMIT 1) ia ON true
JOIN LATERAL (SELECT * FROM offs ORDER BY random() LIMIT 1) o ON true
`, opt.subscriptions); err != nil {
		fmt.Fprintf(os.Stderr, "insert subscriptions: %v\n", err)
		os.Exit(1)
	}

	// trades
	if _, err := tx.Exec(ctx, `
WITH st AS (
  SELECT id, master_account_id
  FROM strategies
),
symbols AS (
  SELECT ARRAY['EURUSD','GBPUSD','USDJPY','XAUUSD','BTCUSD','ETHUSD','AAPL','TSLA','SPY','USOIL']::text[] AS arr
)
INSERT INTO trades(
  strategy_id, master_account_id, symbol, volume_lots, direction,
  open_time, close_time, open_price, close_price, profit, commission, swap, created_at
)
SELECT
  s.id,
  s.master_account_id,
  (SELECT arr[1 + floor(random()*10)::int] FROM symbols),
  calc.vol,
  calc.dir,
  calc.t_open,
  CASE WHEN calc.is_open THEN NULL ELSE calc.t_close END,
  calc.p_open,
  CASE WHEN calc.is_open THEN NULL ELSE calc.p_close END,
  CASE
    WHEN calc.is_open THEN NULL
    ELSE round(((CASE WHEN calc.dir = 'buy'::trade_direction THEN (calc.p_close - calc.p_open) ELSE (calc.p_open - calc.p_close) END) * calc.vol * 100)::numeric, 2)
  END,
  round((random()*5)::numeric, 2),
  round(((random()-0.5)*2)::numeric, 2),
  calc.t_open
FROM (
  SELECT
    gs,
    (random() < 0.10) AS is_open
  FROM generate_series(1, $1) gs
) g
JOIN LATERAL (SELECT * FROM st ORDER BY random() LIMIT 1) s ON true
JOIN LATERAL (
  SELECT
    (CASE WHEN random() < 0.5 THEN 'buy' ELSE 'sell' END)::trade_direction AS dir,
    round((0.01 + random()*2.00)::numeric, 4) AS vol,
    (now() - (random() * interval '90 days')) AS t_open
) base ON true
JOIN LATERAL (
  SELECT
    base.t_open AS t_open,
    base.t_open + (random() * interval '48 hours') AS t_close,
    round((1 + random()*1000)::numeric, 6) AS p_open,
    round((1 + random()*1000 + (random()-0.5)*10)::numeric, 6) AS p_close,
    base.dir AS dir,
    base.vol AS vol,
    g.is_open AS is_open
) calc ON true
`, opt.trades); err != nil {
		fmt.Fprintf(os.Stderr, "insert trades: %v\n", err)
		os.Exit(1)
	}

	// copied_trades (prefer closed trades to have profit)
	if _, err := tx.Exec(ctx, `
WITH subs AS (
  SELECT id, investor_account_id
  FROM subscriptions
),
trs AS (
  SELECT id, volume_lots, profit, commission, swap, open_time, close_time
  FROM trades
  WHERE close_time IS NOT NULL
)
INSERT INTO copied_trades(
  trade_id, subscription_id, investor_account_id, volume_lots,
  profit, commission, swap, open_time, close_time, created_at
)
SELECT
  t.id,
  s.id,
  s.investor_account_id,
  round((t.volume_lots * (0.10 + random()*0.90))::numeric, 4),
  CASE WHEN t.profit IS NULL THEN NULL ELSE round((t.profit * (0.10 + random()*0.90))::numeric, 2) END,
  CASE WHEN t.commission IS NULL THEN NULL ELSE round((abs(t.commission) * (0.10 + random()*0.90))::numeric, 2) END,
  CASE WHEN t.swap IS NULL THEN NULL ELSE round((t.swap * (0.10 + random()*0.90))::numeric, 2) END,
  t.open_time,
  t.close_time,
  now()
FROM generate_series(1, $1) gs
JOIN LATERAL (SELECT * FROM trs ORDER BY random() LIMIT 1) t ON true
JOIN LATERAL (SELECT * FROM subs ORDER BY random() LIMIT 1) s ON true
`, opt.copiedTrades); err != nil {
		fmt.Fprintf(os.Stderr, "insert copied_trades: %v\n", err)
		os.Exit(1)
	}

	// commissions
	if _, err := tx.Exec(ctx, `
WITH subs AS (SELECT id FROM subscriptions)
INSERT INTO commissions(subscription_id, type, amount, period_from, period_to, created_at)
SELECT
  s.id,
  (ARRAY['performance','management','registration'])[1 + floor(random()*3)::int]::commission_type,
  round((random()*200)::numeric, 2),
  now() - (floor(random()*12)::int * interval '30 days'),
  now(),
  now() - (random() * interval '180 days')
FROM generate_series(1, $1) gs
JOIN LATERAL (SELECT * FROM subs ORDER BY random() LIMIT 1) s ON true
`, opt.commissions); err != nil {
		fmt.Fprintf(os.Stderr, "insert commissions: %v\n", err)
		os.Exit(1)
	}

	// favorites
	if opt.favorites > 0 {
		if _, err := tx.Exec(ctx, `
INSERT INTO favorite_strategies(user_id, strategy_id, created_at)
SELECT
  u.id,
  st.id,
  now() - (random() * interval '90 days')
FROM generate_series(1, $1) gs
JOIN LATERAL (SELECT id FROM users ORDER BY random() LIMIT 1) u ON true
JOIN LATERAL (SELECT id FROM strategies ORDER BY random() LIMIT 1) st ON true
ON CONFLICT DO NOTHING
`, opt.favorites); err != nil {
			fmt.Fprintf(os.Stderr, "insert favorites: %v\n", err)
			os.Exit(1)
		}
	}

	// import jobs
	if opt.importJobs > 0 {
		if _, err := tx.Exec(ctx, `
INSERT INTO import_jobs(type, status, file_name, total_rows, processed_rows, error_rows, started_at, finished_at, created_at)
SELECT
  (ARRAY['trades','accounts','statistics'])[1 + floor(random()*3)::int]::import_job_type,
  (ARRAY['pending','running','success','failed'])[1 + floor(random()*4)::int]::import_job_status,
  'file_' || gs::text || '.csv',
  (100 + floor(random()*10000))::int,
  (100 + floor(random()*10000))::int,
  (floor(random()*100))::int,
  now() - (random() * interval '30 days'),
  now() - (random() * interval '1 days'),
  now() - (random() * interval '30 days')
FROM generate_series(1, $1) gs
`, opt.importJobs); err != nil {
			fmt.Fprintf(os.Stderr, "insert import_jobs: %v\n", err)
			os.Exit(1)
		}
	}

	// import job errors
	if opt.importErrors > 0 {
		if _, err := tx.Exec(ctx, `
WITH jobs AS (SELECT id FROM import_jobs)
INSERT INTO import_job_errors(job_id, row_number, raw_data, error_message, created_at)
SELECT
  j.id,
  (1 + floor(random()*10000))::int,
  jsonb_build_object('row', gs, 'value', md5(random()::text)),
  'Synthetic import error',
  now() - (random() * interval '30 days')
FROM generate_series(1, $1) gs
JOIN LATERAL (SELECT * FROM jobs ORDER BY random() LIMIT 1) j ON true
`, opt.importErrors); err != nil {
			fmt.Fprintf(os.Stderr, "insert import_job_errors: %v\n", err)
			os.Exit(1)
		}
	}

	// Ensure strategy_stats is consistent (triggers also do this, but this makes the final state deterministic).
	if _, err := tx.Exec(ctx, `SELECT fn_refresh_strategy_stats(id) FROM strategies`); err != nil {
		fmt.Fprintf(os.Stderr, "refresh strategy_stats: %v\n", err)
		os.Exit(1)
	}

	// Re-enable audit triggers
	if _, err := tx.Exec(ctx, `
ALTER TABLE users ENABLE TRIGGER users_audit_trg;
ALTER TABLE accounts ENABLE TRIGGER accounts_audit_trg;
ALTER TABLE strategies ENABLE TRIGGER strategies_audit_trg;
ALTER TABLE offers ENABLE TRIGGER offers_audit_trg;
ALTER TABLE subscriptions ENABLE TRIGGER subscriptions_audit_trg;
ALTER TABLE trades ENABLE TRIGGER trades_audit_trg;
`); err != nil {
		fmt.Fprintf(os.Stderr, "enable audit triggers: %v\n", err)
		os.Exit(1)
	}

	// Generate audit_log entries via UPDATE operations (triggers will create audit_log entries)
	// This ensures audit_log is populated by triggers, not direct inserts
	if opt.audit > 0 {
		// Set app.current_user_id for audit triggers (they require this session variable)
		var firstUserID int64
		if err := tx.QueryRow(ctx, `SELECT id FROM users LIMIT 1`).Scan(&firstUserID); err != nil {
			fmt.Fprintf(os.Stderr, "get first user id: %v\n", err)
			os.Exit(1)
		}
		// SET LOCAL doesn't support parameters, so we use string formatting
		if _, err := tx.Exec(ctx, fmt.Sprintf(`SET LOCAL app.current_user_id = %d`, firstUserID)); err != nil {
			fmt.Fprintf(os.Stderr, "set app.current_user_id: %v\n", err)
			os.Exit(1)
		}

		// Count existing audit_log entries (should be 0 since triggers were disabled)
		var currentCount int
		if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM audit_log`).Scan(&currentCount); err != nil {
			fmt.Fprintf(os.Stderr, "count audit_log: %v\n", err)
			os.Exit(1)
		}

		needed := opt.audit - currentCount
		if needed > 0 {
			// Use UPDATE operations to trigger audit_log entries
			// Update random records from different tables to generate audit entries
			updatesPerTable := (needed + 5) / 6 // distribute across 6 tables with triggers

			// Update users (will create audit_log entries via trigger)
			if updatesPerTable > 0 {
				if _, err := tx.Exec(ctx, `
UPDATE users SET updated_at = now()
WHERE id IN (
  SELECT id FROM users ORDER BY random() LIMIT $1
)
`, updatesPerTable); err != nil {
					fmt.Fprintf(os.Stderr, "update users for audit: %v\n", err)
					os.Exit(1)
				}
			}

			// Update accounts
			if updatesPerTable > 0 {
				if _, err := tx.Exec(ctx, `
UPDATE accounts SET updated_at = now()
WHERE id IN (
  SELECT id FROM accounts ORDER BY random() LIMIT $1
)
`, updatesPerTable); err != nil {
					fmt.Fprintf(os.Stderr, "update accounts for audit: %v\n", err)
					os.Exit(1)
				}
			}

			// Update strategies
			if updatesPerTable > 0 {
				if _, err := tx.Exec(ctx, `
UPDATE strategies SET updated_at = now()
WHERE id IN (
  SELECT id FROM strategies ORDER BY random() LIMIT $1
)
`, updatesPerTable); err != nil {
					fmt.Fprintf(os.Stderr, "update strategies for audit: %v\n", err)
					os.Exit(1)
				}
			}

			// Update offers
			if updatesPerTable > 0 {
				if _, err := tx.Exec(ctx, `
UPDATE offers SET updated_at = now()
WHERE id IN (
  SELECT id FROM offers ORDER BY random() LIMIT $1
)
`, updatesPerTable); err != nil {
					fmt.Fprintf(os.Stderr, "update offers for audit: %v\n", err)
					os.Exit(1)
				}
			}

			// Update subscriptions
			if updatesPerTable > 0 {
				if _, err := tx.Exec(ctx, `
UPDATE subscriptions SET updated_at = now()
WHERE id IN (
  SELECT id FROM subscriptions ORDER BY random() LIMIT $1
)
`, updatesPerTable); err != nil {
					fmt.Fprintf(os.Stderr, "update subscriptions for audit: %v\n", err)
					os.Exit(1)
				}
			}

			// Update trades
			if updatesPerTable > 0 {
				if _, err := tx.Exec(ctx, `
UPDATE trades SET created_at = now()
WHERE id IN (
  SELECT id FROM trades ORDER BY random() LIMIT $1
)
`, updatesPerTable); err != nil {
					fmt.Fprintf(os.Stderr, "update trades for audit: %v\n", err)
					os.Exit(1)
				}
			}

			// If we still need more entries, do additional updates
			var finalCount int
			if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM audit_log`).Scan(&finalCount); err != nil {
				fmt.Fprintf(os.Stderr, "count audit_log after updates: %v\n", err)
				os.Exit(1)
			}

			stillNeeded := opt.audit - finalCount
			if stillNeeded > 0 {
				// Do more updates on random records
				if _, err := tx.Exec(ctx, `
UPDATE users SET updated_at = now()
WHERE id IN (
  SELECT id FROM users ORDER BY random() LIMIT $1
)
`, stillNeeded); err != nil {
					fmt.Fprintf(os.Stderr, "additional updates for audit: %v\n", err)
					os.Exit(1)
				}
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "commit: %v\n", err)
		os.Exit(1)
	}

	// Tiny summary.
	fmt.Printf("Seed done: users=%d strategies=%d offers=%d subscriptions=%d trades=%d copied_trades=%d commissions=%d audit=%d (seed=%d)\n",
		opt.users, opt.strategies, opt.offers, opt.subscriptions, opt.trades, opt.copiedTrades, opt.commissions, opt.audit, opt.seed)
}
