package publisher

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	Ch *amqp.Channel
}

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	return &Publisher{
		Ch: ch,
	}, nil
}

func (p *Publisher) Run(ctx context.Context, pub *amqp.Publishing, qName, exchange string) error {
	queue, err := p.Ch.QueueDeclare(qName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	err = p.Ch.PublishWithContext(
		ctx,
		exchange,
		queue.Name,
		false,
		false,
		*pub,
	)
	if err != nil {
		return err
	}

	return nil
}
