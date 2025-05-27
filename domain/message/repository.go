package message

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, message *Message) error
	FindAllByStatus(ctx context.Context, status Status) ([]Message, error)
	FindByID(ctx context.Context, id string) (*Message, error)
	UpdateAllStatusesByStatus(ctx context.Context, from, to Status) error
	CreateSentInfo(ctx context.Context, messageID, time string) error
}
