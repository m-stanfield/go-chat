package database

import (
	"context"
	"database/sql"
	"fmt"
)

type AtomicCallback = func(r Repositories) error

type DBConn interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Repositories interface {
	Atomic(context.Context, AtomicCallback) error
	Database() Database
}

type Database struct {
	db   *sql.DB
	conn DBConn
}

func (r *Database) Close() {
	r.db.Close()
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{db: db, conn: db}
}

func (db *Database) withTx(tx *sql.Tx) *Database {
	return &Database{db: db.db, conn: tx}
}

func (r *Database) Atomic(ctx context.Context, cb func(ds *Database) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("tx err: %w, rb err: %w", err, rbErr)
			}
		} else {
			err = tx.Commit()
		}
	}()
	dbTx := r.withTx(tx)
	err = cb(dbTx)
	return err
}
