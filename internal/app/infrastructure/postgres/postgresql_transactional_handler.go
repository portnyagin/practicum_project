package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"github.com/portnyagin/practicum_project/internal/app/repository/basedbhandler"
	"go.uber.org/zap"
	"time"
)

type PostgresqlHandlerTX struct {
	pool *pgxpool.Pool
	log  *infrastructure.Logger
}

func NewPostgresqlHandlerTX(ctx context.Context, dataSource string, log *infrastructure.Logger) (*PostgresqlHandlerTX, error) {
	// Format DSN
	//("postgresql://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname)
	poolConfig, err := pgxpool.ParseConfig(dataSource)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = 5
	poolConfig.MinConns = 2
	poolConfig.MaxConnIdleTime = time.Second * 120
	pool, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	postgresqlHandler := new(PostgresqlHandlerTX)
	postgresqlHandler.pool = pool
	postgresqlHandler.log = log
	return postgresqlHandler, nil
}

func (handler *PostgresqlHandlerTX) NewTx(ctx context.Context) (pgx.Tx, error) {
	return handler.pool.Begin(ctx)
}

func (handler *PostgresqlHandlerTX) getTx(ctx context.Context) (tx pgx.Tx, err error) {
	defer func() {
		if recover() != nil {
			err = errors.New("can't get tx: conversion error")
			handler.log.Error("PostgresqlHandlerTX: can't get tx", zap.Error(err))
		}
	}()
	ctxValue := ctx.Value("tx")
	if ctxValue == nil {
		handler.log.Debug("PostgresqlHandlerTX: can't get tx")
		return nil, errors.New("can't get tx: nil value got")
	}
	// TODO: тест на обработку паники
	tx = ctxValue.(pgx.Tx)
	return tx, err
}

func (handler *PostgresqlHandlerTX) Commit(ctx context.Context) error {
	tx, err := handler.getTx(ctx)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		handler.log.Error("Can't commit transaction", zap.Error(err))
		return err
	}
	return err
}

func (handler *PostgresqlHandlerTX) Rollback(ctx context.Context) error {
	tx, err := handler.getTx(ctx)
	if err != nil {
		return err
	}
	err = tx.Rollback(ctx)
	if err != nil {
		handler.log.Error("Can't commit transaction", zap.Error(err))
		return err
	}
	return err
}

func (handler *PostgresqlHandlerTX) Execute(ctx context.Context, statement string, args ...interface{}) error {
	tx, err := handler.getTx(ctx)
	// Пытаемся получить транзакцию из контекста, если не нашли, работаем без транзакции
	if err == nil {
		if len(args) > 0 {
			_, err = tx.Exec(ctx, statement, args...)
		} else {
			_, err = tx.Exec(ctx, statement)
		}
	} else {
		conn, e := handler.pool.Acquire(ctx)
		if e != nil {
			return e
		}
		defer conn.Release()

		if len(args) > 0 {
			_, e = conn.Exec(ctx, statement, args...)
		} else {
			_, e = conn.Exec(ctx, statement)
		}
		err = e
	}
	return err
}

func (handler *PostgresqlHandlerTX) ExecuteBatch(ctx context.Context, statement string, args [][]interface{}) error {
	var (
		err error
		ct  pgconn.CommandTag
		br  pgx.BatchResults
	)

	batch := &pgx.Batch{}
	if len(args) > 0 {
		for _, argset := range args {
			batch.Queue(statement, argset...)
		}
	} else {
		return nil
	}
	tx, err := handler.getTx(ctx)
	// Пытаемся получить транзакцию из контекста, если не нашли, работаем без транзакции
	if err == nil {
		br = tx.SendBatch(context.Background(), batch)
	} else {
		conn, err := handler.pool.Acquire(ctx)
		if err != nil {
			return err
		}
		defer conn.Release()
		br = conn.SendBatch(context.Background(), batch)
	}
	ct, err = br.Exec()

	if err != nil {
		return err
	}
	fmt.Println(ct.RowsAffected())
	return nil
}

func (handler *PostgresqlHandlerTX) QueryRow(ctx context.Context, statement string, args ...interface{}) (basedbhandler.Row, error) {
	var row pgx.Row
	tx, err := handler.getTx(ctx)
	// Пытаемся получить транзакцию из контекста, если не нашли, работаем без транзакции
	if err == nil {
		if len(args) > 0 {
			row = tx.QueryRow(ctx, statement, args...)
		} else {
			row = tx.QueryRow(ctx, statement)
		}
	} else {
		conn, err := handler.pool.Acquire(ctx)
		if err != nil {
			return nil, err
		}
		defer conn.Release()
		if len(args) > 0 {
			row = conn.QueryRow(ctx, statement, args...)
		} else {
			row = conn.QueryRow(ctx, statement)
		}
	}
	return row, nil
}

func (handler *PostgresqlHandlerTX) Query(ctx context.Context, statement string, args ...interface{}) (basedbhandler.Rows, error) {
	var rows pgx.Rows
	tx, err := handler.getTx(ctx)
	// Пытаемся получить транзакцию из контекста, если не нашли, работаем без транзакции
	if err == nil {
		if len(args) > 0 {
			rows, err = tx.Query(ctx, statement, args...)
		} else {
			rows, err = tx.Query(ctx, statement)
		}
	} else {
		conn, e := handler.pool.Acquire(ctx)
		if e != nil {
			return nil, e
		}
		defer conn.Release()
		if len(args) > 0 {
			rows, e = conn.Query(ctx, statement, args...)
		} else {
			rows, e = conn.Query(ctx, statement)
		}
		err = e
	}
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (handler *PostgresqlHandlerTX) Close() {
	if handler != nil {
		handler.pool.Close()
	}
}
