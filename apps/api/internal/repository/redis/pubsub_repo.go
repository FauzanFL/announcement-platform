package redis

import (
	"announcement-api/internal/domain/repository"
	"context"

	"github.com/redis/go-redis/v9"
)

type pubsubRepo struct {
	client *redis.Client
}

func NewPubSubRepository(client *redis.Client) repository.PubSubRepository {
	return &pubsubRepo{client: client}
}

func (r *pubsubRepo) Publish(ctx context.Context, channel string, message []byte) error {
	return r.client.Publish(ctx, channel, message).Err()
}

func (r *pubsubRepo) Subscribe(ctx context.Context, channel string) (<-chan string, func(), error) {
	sub := r.client.Subscribe(ctx, channel)
	redisChannel := sub.Channel()

	out := make(chan string)
	go func() {
		defer close(out)
		for msg := range redisChannel {
			out <- msg.Payload
		}
	}()

	cleanup := func() {
		_ = sub.Close()
	}

	return out, cleanup, nil
}
