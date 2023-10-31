package rabbitmq

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func New(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
	}, nil
}

func (r *RabbitMQ) Publish(ctx context.Context, key string, message []byte) error {
	_, err := r.ch.QueueDeclare(
		key,   // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	return r.ch.PublishWithContext(ctx, "", key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        message,
	})
}

func (r *RabbitMQ) Subscribe(key, consumer string, handler func(message []byte) error) error {
	_, err := r.ch.QueueDeclare(
		key,   // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	msgs, err := r.ch.Consume(
		key,      // queue
		consumer, // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return fmt.Errorf("consume: %w", err)
	}

	for d := range msgs {
		err := handler(d.Body)
		if err != nil {
			log.Println("subscriber handler:", err)
		}

		err = d.Ack(false)
		if err != nil {
			return fmt.Errorf("ack: %w", err)
		}
	}
	return nil
}

func (r *RabbitMQ) Unsubscribe(key string) error {
	return r.ch.Cancel(key, false)
}

func (r *RabbitMQ) Close() error {
	return r.conn.Close()
}
