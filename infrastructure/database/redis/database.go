package redis

import (
	"context"
	"fmt"
	"net"

	rdb "github.com/redis/go-redis/v9"
)

type Redis interface {
	Close() error
	Set(ctx context.Context, key string, value any) error
}

type Config struct {
	Host     string
	Port     uint16
	User     string
	Password string
	DB       uint16
}

type redis struct {
	client *rdb.Client
}

func New(config Config) (Redis, error) {
	client := rdb.NewClient(&rdb.Options{
		Addr:     net.JoinHostPort(config.Host, fmt.Sprintf("%d", config.Port)),
		Username: config.User,
		Password: config.Password,
		DB:       int(config.DB),
	})

	ping := client.Ping(context.Background())
	if ping.Err() != nil {
		return nil, fmt.Errorf("client.Ping(): %w", ping.Err())
	}

	return &redis{
		client: client,
	}, nil
}

func (r *redis) Close() error {
	if err := r.client.Close(); err != nil {
		return fmt.Errorf("redis.client.Close(): %w", err)
	}

	return nil
}

func (r *redis) Set(ctx context.Context, key string, value any) error {
	if err := r.client.Set(ctx, key, value, 0).Err(); err != nil {
		return fmt.Errorf("redis.client.Set(): %w", err)
	}

	return nil
}
