package esync

import "github.com/jackc/pgx/v5/pgxpool"

type SyncDB struct {
	db *pgxpool.Pool
}

func NewSyncDB(db *pgxpool.Pool) *SyncDB {
	return &SyncDB{db: db}
}
