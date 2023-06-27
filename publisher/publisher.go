package publisher

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	Ctx  context.Context
	Ch   *amqp.Channel
	Conn *amqp.Connection
}

func NewPublisher(ctx context.Context, conn *amqp.Connection) (*Publisher, error) {
	return &Publisher{
		Ctx:  ctx,
		Conn: conn,
	}, nil
}

func (p *Publisher) Run(pub *amqp.Publishing, qName, exchange string) error {
	ch, err := p.Conn.Channel()
	if err != nil {
		return fmt.Errorf("open channel: %s", err.Error())
	}
	// defer ch.Close()

	p.Ch = ch

	queue, err := ch.QueueDeclare(qName, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %s", err.Error())
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
		return fmt.Errorf("publish with context: %s", err.Error())
	}

	return nil
}

func (c *Publisher) CloseCh() error {
	err := c.Ch.Close()
	if err != nil {
		return fmt.Errorf("close publisher ch: %s", err.Error())
	}
	return nil
}
