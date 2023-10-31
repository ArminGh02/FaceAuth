package broker

import "context"

type Broker interface {
	Publish(ctx context.Context, key string, message []byte) error
	Subscribe(key, consumer string, handler func(message []byte) error) error
	Unsubscribe(key string) error
	Close() error
}
