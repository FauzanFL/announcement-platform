package repository

import "context"

type PubSubRepository interface {
	Publish(ctx context.Context, channel string, message []byte) error
	Subscribe(ctx context.Context, channel string) (<-chan string, func(), error)
}
