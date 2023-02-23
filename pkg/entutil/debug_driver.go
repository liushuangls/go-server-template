package entutil

import (
	"context"
	"fmt"
	"log"
	"time"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

// 实现打印sql执行时间，https://github.com/ent/ent/issues/3055

// DebugDriver is a driver that logs all driver operations.
type DebugDriver struct {
	dialect.Driver // underlying driver.
	log            func(context.Context, ...any)
}

// Debug gets a driver and an optional logging function, and returns
// a new debugged-driver that prints all outgoing operations.
func Debug(d dialect.Driver, logger ...func(...any)) dialect.Driver {
	logf := log.Println
	if len(logger) == 1 {
		logf = logger[0]
	}
	drv := &DebugDriver{
		Driver: d,
		log:    func(_ context.Context, v ...any) { logf(v...) },
	}
	return drv
}

// Exec logs its params and calls the underlying driver Exec method.
func (d *DebugDriver) Exec(ctx context.Context, query string, args, v any) error {
	start := time.Now()
	err := d.Driver.Exec(ctx, query, args, v)
	d.log(ctx, fmt.Sprintf("driver.Exec: query=%v args=%v time=%v\n", query, args, time.Since(start)))
	return err
}

// ExecContext logs its params and calls the underlying driver ExecContext method if it is supported.
func (d *DebugDriver) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	drv, ok := d.Driver.(interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.ExecContext is not supported")
	}
	start := time.Now()
	result, err := drv.ExecContext(ctx, query, args...)
	d.log(ctx, fmt.Sprintf("driver.ExecContext: query=%v args=%v time=%v\n", query, args, time.Since(start)))
	return result, err
}

// Query logs its params and calls the underlying driver Query method.
func (d *DebugDriver) Query(ctx context.Context, query string, args, v any) error {
	start := time.Now()
	err := d.Driver.Query(ctx, query, args, v)
	d.log(ctx, fmt.Sprintf("driver.Query: query=%v args=%v time=%v\n", query, args, time.Since(start)))
	return err
}

// QueryContext logs its params and calls the underlying driver QueryContext method if it is supported.
func (d *DebugDriver) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	drv, ok := d.Driver.(interface {
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.QueryContext is not supported")
	}
	start := time.Now()
	rows, err := drv.QueryContext(ctx, query, args...)
	d.log(ctx, fmt.Sprintf("driver.QueryContext: query=%v args=%v time=%v\n", query, args, time.Since(start)))
	return rows, err
}

// Tx adds an log-id for the transaction and calls the underlying driver Tx command.
func (d *DebugDriver) Tx(ctx context.Context) (dialect.Tx, error) {
	tx, err := d.Driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	d.log(ctx, fmt.Sprintf("driver.Tx(%s): started", id))
	return &DebugTx{tx, id, d.log, ctx}, nil
}

// BeginTx adds an log-id for the transaction and calls the underlying driver BeginTx command if it is supported.
func (d *DebugDriver) BeginTx(ctx context.Context, opts *sql.TxOptions) (dialect.Tx, error) {
	drv, ok := d.Driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	})
	if !ok {
		return nil, fmt.Errorf("Driver.BeginTx is not supported")
	}
	tx, err := drv.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	id := uuid.New().String()
	d.log(ctx, fmt.Sprintf("driver.BeginTx(%s): started", id))
	return &DebugTx{tx, id, d.log, ctx}, nil
}

// DebugTx is a transaction implementation that logs all transaction operations.
type DebugTx struct {
	dialect.Tx                               // underlying transaction.
	id         string                        // transaction logging id.
	log        func(context.Context, ...any) // log function. defaults to fmt.Println.
	ctx        context.Context               // underlying transaction context.
}

// Exec logs its params and calls the underlying transaction Exec method.
func (d *DebugTx) Exec(ctx context.Context, query string, args, v any) error {
	start := time.Now()
	err := d.Tx.Exec(ctx, query, args, v)
	d.log(ctx, fmt.Sprintf("Tx(%s).Exec: query=%v args=%v time=%v\n", d.id, query, args, time.Since(start)))
	return err
}

// ExecContext logs its params and calls the underlying transaction ExecContext method if it is supported.
func (d *DebugTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	drv, ok := d.Tx.(interface {
		ExecContext(context.Context, string, ...any) (sql.Result, error)
	})
	if !ok {
		return nil, fmt.Errorf("Tx.ExecContext is not supported")
	}
	start := time.Now()
	result, err := drv.ExecContext(ctx, query, args...)
	d.log(ctx, fmt.Sprintf("Tx(%s).ExecContext: query=%v args=%v time=%v\n", d.id, query, args, time.Since(start)))
	return result, err
}

// Query logs its params and calls the underlying transaction Query method.
func (d *DebugTx) Query(ctx context.Context, query string, args, v any) error {
	start := time.Now()
	err := d.Tx.Query(ctx, query, args, v)
	d.log(ctx, fmt.Sprintf("Tx(%s).Query: query=%v args=%v time=%v\n", d.id, query, args, time.Since(start)))
	return err
}

// QueryContext logs its params and calls the underlying transaction QueryContext method if it is supported.
func (d *DebugTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	drv, ok := d.Tx.(interface {
		QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	})
	if !ok {
		return nil, fmt.Errorf("Tx.QueryContext is not supported")
	}
	start := time.Now()
	rows, err := drv.QueryContext(ctx, query, args...)
	d.log(ctx, fmt.Sprintf("Tx(%s).QueryContext: query=%v args=%v time=%v\n", d.id, query, args, time.Since(start)))
	return rows, err
}

// Commit logs this step and calls the underlying transaction Commit method.
func (d *DebugTx) Commit() error {
	d.log(d.ctx, fmt.Sprintf("Tx(%s): committed", d.id))
	return d.Tx.Commit()
}
