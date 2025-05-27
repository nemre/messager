package message

import (
	"context"
	"fmt"

	"messager/domain/message"
)

func (s *service) Process(ctx context.Context) error {
	if err := s.repository.UpdateAllStatusesByStatus(ctx, message.StatusPending, message.StatusSent); err != nil {
		return fmt.Errorf("service.repository.UpdateAllStatusesByStatus(): %w", err)
	}

	return nil
}
