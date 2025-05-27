package message

import (
	"context"
	"fmt"

	"messager/domain/message"
)

func (p *persistence) FindByID(ctx context.Context, id string) (*message.Message, error) {
	query := `
		SELECT id, created_at, updated_at, content, phone, status
		FROM messages
		WHERE id = $1
	`
	row := p.postgreSQL.QueryRow(ctx, query, id)

	var record message.Message

	if err := row.Scan(&record.ID, &record.CreatedAt, &record.UpdatedAt, &record.Content, &record.Phone, &record.Status); err != nil {
		return nil, fmt.Errorf("persistence.postgreSQL.QueryRow().Row.Scan(): %w", err)
	}

	return &record, nil
}
