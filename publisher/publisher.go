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
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open channel: %s", err.Error())
	}

	return &Publisher{
		Ctx:  ctx,
		Ch:   ch,
		Conn: conn,
	}, nil
}

func (p *Publisher) Run(pub *amqp.Publishing, qName, exchange string) error {
	_, err := p.Ch.QueueDeclare(qName, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %s", err.Error())
	}

	err = p.Ch.PublishWithContext(
		p.Ctx,
		exchange,
		qName,
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
	err := c.Ch.Close()
	if err != nil {
		return fmt.Errorf("close publisher ch: %s", err.Error())
	}
	return nil
}
