package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/hatchways-community/2e26b1bef5c64db4a4d3e9decab77101/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate mockgen -destination=internal/postgres/database_mock.go -package=postgres -source=internal/postgres/database.go Conn
type Conn interface {
	Close()
	Ping(context.Context) error
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type Database struct {
	*pgxpool.Pool
	sq squirrel.StatementBuilderType
}

func NewDatabase(
	ctx context.Context,
	cfg config.Database,
) (*Database, error) {
	c, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, err
	}

	return &Database{
		Pool: pool,
		sq:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

func (c *Database) Ping(ctx context.Context) error {
	conn, err := c.Pool.Acquire(ctx)
	if err != nil {
		return err
	}

	return conn.Conn().Ping(ctx)
}
