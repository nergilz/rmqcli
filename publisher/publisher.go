package publisher

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	Ctx  context.Context
	Ch   *amqp.Channel
	Conn *amqp.Connection
}

func NewPublisher(ctx context.Context, conn *amqp.Connection) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	// defer ch.Close()

	return &Publisher{
		Ctx:  ctx,
		Ch:   ch,
		Conn: conn,
	}, nil
}

func (p *Publisher) Run(pub *amqp.Publishing, qName, exchange string) error {
	queue, err := p.Ch.QueueDeclare(qName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	err = p.Ch.PublishWithContext(
		p.Ctx,
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

func (c *Publisher) CloseCh() error {
	err := c.Ch.Close()
	if err != nil {
		return err
	}
	return nil
}
