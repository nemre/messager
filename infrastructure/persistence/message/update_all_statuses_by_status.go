package message

import (
	"context"
	"fmt"

	"messager/domain/message"
)

func (p *persistence) UpdateAllStatusesByStatus(ctx context.Context, from, to message.Status) error {
	query := `
		UPDATE messages
		SET status = $1
		WHERE status = $2
	`

	if err := p.postgreSQL.Exec(ctx, query, to, from); err != nil {
		return fmt.Errorf("persistence.postgreSQL.Exec(): %w", err)
	}

	return nil
}
