package message

import (
	"context"
	"fmt"
)

func (p *persistence) CreateSentInfo(ctx context.Context, messageID, time string) error {
	if err := p.redis.Set(ctx, fmt.Sprintf("message:%s", messageID), time); err != nil {
		return fmt.Errorf("persistence.redis.Set(): %w", err)
	}

	return nil
}
