package postgresql

import (
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Rows interface {
	Close()
	Next() bool
	Scan(destination ...any) error
	Err() error
}

type rows struct {
	rows pgx.Rows
}

func (r *rows) Close() {
	r.rows.Close()
}

func (r *rows) Next() bool {
	return r.rows.Next()
}

func (r *rows) Scan(destination ...any) error {
	if err := r.rows.Scan(destination...); err != nil {
		return fmt.Errorf("rows.rows.Scan(): %w", err)
	}

	return nil
}

func (r *rows) Err() error {
	return r.rows.Err()
}
