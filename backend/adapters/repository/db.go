package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

type AccountDB struct {
	db *pgxpool.Pool
}

type WorkspaceDB struct {
	db *pgxpool.Pool
}

type MemberDB struct {
	db *pgxpool.Pool
}

type CustomerDB struct {
	db *pgxpool.Pool
}

type ThreadDB struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewAccountDB(db *pgxpool.Pool) *AccountDB {
	return &AccountDB{
		db: db,
	}
}

func NewWorkspaceDB(db *pgxpool.Pool) *WorkspaceDB {
	return &WorkspaceDB{
		db: db,
	}
}

func NewMemberDB(db *pgxpool.Pool) *MemberDB {
	return &MemberDB{
		db: db,
	}
}

func NewCustomerDB(db *pgxpool.Pool) *CustomerDB {
	return &CustomerDB{
		db: db,
	}
}

func NewThreadDB(db *pgxpool.Pool, rdb *redis.Client) *ThreadDB {
	return &ThreadDB{
		db:  db,
		rdb: rdb,
	}
}

func debugQuery(query string) {
	slog.Info("db", slog.Any("query", query))
}
