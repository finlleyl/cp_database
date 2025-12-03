package strategy

import (
	"time"

	"github.com/finlleyl/cp_database/internal/domain/common"
)

type Strategy struct {
	ID              int64                 `json:"id" db:"id"`
	MasterUserID    int64                 `json:"master_user_id" db:"master_user_id"`
	MasterAccountID int64                 `json:"master_account_id" db:"master_account_id"`
	Title           string                `json:"title" db:"title"`
	Description     string                `json:"description" db:"description"`
	Status          common.StrategyStatus `json:"status" db:"status"`
	CreatedAt       time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at" db:"updated_at"`
}
