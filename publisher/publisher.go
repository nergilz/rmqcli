package publisher

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	Conn *amqp.Connection
}

func NewPublisher(conn *amqp.Connection) (*Publisher, error) {
	return &Publisher{
		Conn: conn,
	}, nil
}

func (p *Publisher) Run(ctx context.Context, pub *amqp.Publishing, qName, exchange string) error {
	ch, err := p.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := ch.QueueDeclare(qName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(
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
