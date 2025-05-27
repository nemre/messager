package postgresql

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var ErrNoRows = errors.New("no rows in result set")

type Row interface {
	Scan(destination ...any) error
}

type row struct {
	row pgx.Row
}

func (r *row) Scan(destination ...any) error {
	err := r.row.Scan(destination...)
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("row.row.Scan(): %w", ErrNoRows)
	}
	if err != nil {
		return fmt.Errorf("row.row.Scan(): %w", err)
	}

	return nil
}
