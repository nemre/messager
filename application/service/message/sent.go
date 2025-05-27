package message

import (
	"context"
	"errors"
	"fmt"
	"time"

	"messager/domain/message"
	entity "messager/domain/message"
	"messager/infrastructure/database/postgresql"
)

func (s *service) Sent(ctx context.Context, message message.Message) error {
	if err := message.ValidateForSent(); err != nil {
		return errors.Join(message.NewErrMessageDoesNotValidForSent(), err)
	}

	foundMessage, err := s.repository.FindByID(ctx, message.ID)
	if foundMessage == nil || errors.Is(err, postgresql.ErrNoRows) {
		return message.NewErrMessageNotFound()
	}
	if err != nil {
		return fmt.Errorf("service.repository.FindByID(): %w", err)
	}

	if foundMessage.Status != entity.StatusSent {
		return message.NewErrMessageStatusDoesNotEligibleForSent()
	}

	id, err := s.client.SendMessage(ctx, *foundMessage)
	if err != nil {
		return fmt.Errorf("service.client.SendMessage(): %w", err)
	}

	if err := s.repository.CreateSentInfo(ctx, id, time.Now().Format(time.RFC3339)); err != nil {
		return fmt.Errorf("service.repository.CreateSentInfo(): %w", err)
	}

	return nil
}
