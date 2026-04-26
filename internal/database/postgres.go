package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PoolConfig struct {
	DSN            string
	MaxConns       int
	MinConns       int
	ConnLifetime   time.Duration
	ConnIdleTime   time.Duration
	ConnectTimeout time.Duration
}

func NewPostgresPool(ctx context.Context, poolcfg PoolConfig) (*pgxpool.Pool, error) {

	ctx, cancel := context.WithTimeout(ctx, poolcfg.ConnectTimeout)
	defer cancel()

	// set up connection pool configuration
	poolConfig, err := pgxpool.ParseConfig(poolcfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("internal.database.postgres: error parsing DSN: %w", err)
	}

	poolConfig.MaxConns = int32(poolcfg.MaxConns)
	poolConfig.MinConns = int32(poolcfg.MinConns)
	poolConfig.MaxConnLifetime = poolcfg.ConnLifetime
	poolConfig.MaxConnIdleTime = poolcfg.ConnIdleTime

	// create the connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("internal.database.postgres: error creating connection pool: %w", err)
	}

	// test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("internal.database.postgres: error testing connection: %w", err)
	}

	return pool, nil
}
