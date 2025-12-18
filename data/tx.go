package data

import (
	"context"
	"fmt"
	"log/slog"

	ent2 "github.com/liushuangls/go-server-template/data/ent"
)

func withTx(ctx context.Context, db *ent2.Client, fn func(tx *ent2.Tx) error) error {
	tx, err := db.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			if err := tx.Rollback(); err != nil {
				slog.Error("withTx: rolling back transaction", "err", err)
			}
			panic(v)
		}
	}()
	if err := fn(tx); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, err2)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}
