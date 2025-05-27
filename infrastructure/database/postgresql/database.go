package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQL interface {
	Close()
	Query(ctx context.Context, query string, arguments ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, arguments ...any) Row
	Exec(ctx context.Context, query string, arguments ...any) error
}

type Config struct {
	Host     string
	Port     uint16
	User     string
	Password string
	Name     string
}

type postgreSQL struct {
	pool *pgxpool.Pool
}

func New(config Config) (PostgreSQL, error) {
	poolConfig, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.Name,
	))
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig(): %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig(): %w", err)
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("pool.Ping(): %w", err)
	}

	return &postgreSQL{
		pool: pool,
	}, nil
}

func (p *postgreSQL) Close() {
	p.pool.Close()
}

func (p *postgreSQL) Query(ctx context.Context, query string, arguments ...any) (Rows, error) {
	queryRows, err := p.pool.Query(ctx, query, arguments...)
	if err != nil {
		return nil, fmt.Errorf("postgreSQL.pool.Query(): %w", err)
	}

	return &rows{
		rows: queryRows,
	}, nil
}

func (p *postgreSQL) QueryRow(ctx context.Context, query string, arguments ...any) Row {
	return &row{
		row: p.pool.QueryRow(ctx, query, arguments...),
	}
}

func (p *postgreSQL) Exec(ctx context.Context, query string, arguments ...any) error {
	_, err := p.pool.Exec(ctx, query, arguments...)
	if err != nil {
		return fmt.Errorf("postgreSQL.pool.Exec(): %w", err)
	}

	return nil
}
