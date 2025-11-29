package repository

import "github.com/jmoiron/sqlx"

type (
	MasterRepository struct {
		db *sqlx.DB
	}
)

func NewMasterRepository(db *sqlx.DB) *MasterRepository {
	return &MasterRepository{db: db}
}