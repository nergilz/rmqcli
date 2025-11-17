package publisher

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch   *amqp.Channel
	conn *amqp.Connection
}

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %s", err.Error())
	}

	return &Publisher{
		ch:   ch,
		conn: conn,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, pub *amqp.Publishing, exchange, key string) error {
	_, err := p.ch.QueueDeclare(
		key,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare: %s", err.Error())
	}

	err = p.ch.PublishWithContext(
		ctx,
		exchange,
		key, // routing_key
		false,
		false,
		*pub,
	)
	if err != nil {
		return fmt.Errorf("publish with context: %s", err.Error())
	}

	return nil
}

func (c *Publisher) CloseChannel() error {
	err := c.ch.Close()
	if err != nil {
		return fmt.Errorf("close publisher ch: %s", err.Error())
	}

	return nil
}
