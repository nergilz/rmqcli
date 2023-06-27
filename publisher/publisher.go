package publisher

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	Ch *amqp.Channel
}

func NewPublisher(ch *amqp.Channel) *Publisher {
	return &Publisher{
		Ch: ch,
	}
}

func (p *Publisher) Run(ctx context.Context, conn *amqp.Connection, pub *amqp.Publishing, qName, exchange string) error {
	queue, err := p.Ch.QueueDeclare(qName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	err = p.Ch.PublishWithContext(
		ctx,
		exchange,   // exchange
		queue.Name, // key
		false,
		false,
		*pub,
	)
	if err != nil {
		return err
	}

	return nil
}
