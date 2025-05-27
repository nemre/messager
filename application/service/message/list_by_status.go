package message

import (
	"context"
	"errors"
	"fmt"

	"messager/domain/message"
)

func (s *service) ListByStatus(ctx context.Context, status message.Status) ([]message.Message, error) {
	message := message.Message{
		Status: status,
	}

	if err := message.ValidateForListByStatus(); err != nil {
		return nil, errors.Join(message.NewErrMessageDoesNotValidForListByStatus(), err)
	}

	messages, err := s.repository.FindAllByStatus(ctx, status)
	if err != nil {
		return nil, fmt.Errorf("service.repository.FindAllByStatus(): %w", err)
	}

	return messages, nil
}
