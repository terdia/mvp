package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func OpenDb(cfg config) (*sql.DB, error) {

	conn, err := sql.Open("postgres", cfg.Db.Dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(cfg.Db.MaxOpenConns)
	conn.SetMaxIdleConns(cfg.Db.MaxIdleConns)

	maxIdleTime, err := time.ParseDuration(cfg.Db.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	conn.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//if connection is not established within 5 seconds return error
	err = conn.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
