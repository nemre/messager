package message

import "context"

type Service interface {
	Create(ctx context.Context, message Message) (*Message, error)
	ListByStatus(ctx context.Context, status Status) ([]Message, error)
	Process(ctx context.Context) error
	Sent(ctx context.Context, message Message) error
}
