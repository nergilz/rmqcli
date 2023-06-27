package rmqcli

import (
	"context"

	"github.com/nergilz/rmqcli"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	Cfg rmqcli.RmqConfig
	Ch  *amqp.Channel
}

func NewPublisher(cfg *rmqcli.RmqConfig, ch *amqp.Channel) *Publisher {
	return &Publisher{
		Cfg: *cfg,
		Ch:  ch,
	}
}

func (p *Publisher) Run(ctx context.Context, conn *amqp.Connection, pub *amqp.Publishing) error {
	queue, err := p.Ch.QueueDeclare(p.Cfg.Queue, false, false, false, false, nil)
	if err != nil {
		return err
	}

	err = p.Ch.PublishWithContext(
		ctx,
		p.Cfg.Exchange, // exchange
		queue.Name,     // key
		false,
		false,
		*pub,
	)
	if err != nil {
		return err
	}

	return nil
}
