package message

import (
	"context"
	"fmt"

	"messager/domain/message"
)

func (p *persistence) Create(ctx context.Context, message *message.Message) error {
	query := `
		INSERT INTO messages (content, phone, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`

	row := p.postgreSQL.QueryRow(ctx, query,
		message.Content, message.Phone, message.Status)

	if err := row.Scan(&message.ID, &message.CreatedAt, &message.UpdatedAt); err != nil {
		return fmt.Errorf("persistence.postgreSQL.QueryRow().Row.Scan(): %w", err)
	}

	return nil
}
