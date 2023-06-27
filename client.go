package rmqcli

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RmqConfig struct {
	Url         string
	Queue       string
	Concurrency int
	Exchange    string
	RoutingKey  string
}

type RmqCli struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func InitRmqCli(ctx context.Context, cfg *RmqConfig) (*RmqCli, error) {
	conn, err := amqp.Dial(cfg.Url)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	cli := &RmqCli{
		Conn: conn,
		Ch:   ch,
	}

	return cli, nil
}
