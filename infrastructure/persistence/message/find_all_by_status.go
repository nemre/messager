package message

import (
	"context"
	"fmt"

	"messager/domain/message"
)

func (p *persistence) FindAllByStatus(ctx context.Context, status message.Status) ([]message.Message, error) {
	query := `
		SELECT id, created_at, updated_at, content, phone, status
		FROM messages
		WHERE status = $1
		ORDER BY created_at DESC;
	`
	rows, err := p.postgreSQL.Query(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("persistence.postgreSQL.Query(): %w", err)
	}

	defer rows.Close()

	var records []message.Message

	for rows.Next() {
		var record message.Message

		if err := rows.Scan(&record.ID, &record.CreatedAt, &record.UpdatedAt, &record.Content, &record.Phone, &record.Status); err != nil {
			return nil, fmt.Errorf("persistence.postgreSQL.Query().Rows.Scan(): %w", err)
		}

		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("persistence.postgreSQL.Query().Rows.Err(): %w", err)
	}

	return records, nil
}
