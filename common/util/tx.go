package util

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"runtime/debug"

	"github.com/jmoiron/sqlx"
)

type TransactionRunner struct {
	DB *sqlx.DB
}

type TxFunc func(tx *sqlx.Tx) error
type TxOpt func(t *TransactionRunner)

func SetDB(db *sqlx.DB) TxOpt {
	return func(t *TransactionRunner) {
		t.DB = db
	}
}

func NewTransactionRunner(db *sqlx.DB) *TransactionRunner {
	return &TransactionRunner{
		DB: db,
	}
}

func (t *TransactionRunner) WithTx(ctx context.Context, txFunc TxFunc, opts *sql.TxOptions) (err error) {
	tx, err := StartTx(ctx, t.DB, opts)
	if err != nil {
		return err
	}

	defer func() {
		r := recover()
		if r != nil {
			log.Println(string(debug.Stack()))
			mErr := fmt.Errorf("%v", r)
			errRb := RollbackTx(tx)
			if errRb != nil {
				err = errRb
			} else {
				err = mErr
			}
		} else {
			err = CommitTx(tx)
		}

	}()
	err = txFunc(tx)
	PanicIfError(err)
	return nil
}

func StartTx(ctx context.Context, db *sqlx.DB, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return db.BeginTxx(ctx, opts)
}

func RollbackTx(tx *sqlx.Tx) error {
	return tx.Rollback()
}

func CommitTx(tx *sqlx.Tx) error {
	return tx.Commit()
}
