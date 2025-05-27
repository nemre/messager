package message

import (
	"context"
	"fmt"

	"messager/domain/message"
	"messager/infrastructure/database/postgresql"
	"messager/infrastructure/database/redis"
)

type persistence struct {
	postgreSQL postgresql.PostgreSQL
	redis      redis.Redis
}

func New(postgreSQL postgresql.PostgreSQL, redis redis.Redis) (message.Repository, error) {
	p := persistence{
		postgreSQL: postgreSQL,
		redis:      redis,
	}

	if err := p.migrate(context.Background()); err != nil {
		return nil, fmt.Errorf("persistence.migrate(): %w", err)
	}

	return &p, nil
}
