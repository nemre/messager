package message

import (
	"context"
	"errors"
	"fmt"

	"messager/domain/message"
)

func (s *service) Create(ctx context.Context, message message.Message) (*message.Message, error) {
	if err := message.ValidateForCreate(); err != nil {
		return nil, errors.Join(message.NewErrMessageDoesNotValidForCreate(), err)
	}

	if err := s.repository.Create(ctx, &message); err != nil {
		return nil, fmt.Errorf("service.repository.Create(): %w", err)
	}

	return &message, nil
}
