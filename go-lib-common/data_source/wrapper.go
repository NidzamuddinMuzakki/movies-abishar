package data_source

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Exec wrapping multiple queries or single query without transaction.
func Exec(ctx context.Context, db *sqlx.DB, statements ...*Statement) error {
	err := run(ctx, db, statements...)
	fmt.Println(*&statements[0].query, &statements[0].args, err, "nidzam23")
	if err != nil {
		return err
	}

	return nil
}

// ExecTx wrapping multiple queries or single query in a transaction.
func ExecTx(ctx context.Context, db *sqlx.DB, statements ...*Statement) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err = runTx(ctx, tx, statements...); err != nil {
		fmt.Println(*&statements[0].query, &statements[0].args, err, "nidzam23")
		if er := tx.Rollback(); er != nil {
			return er
		}
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
