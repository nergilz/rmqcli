package rmqcli

import (
	"context"

	"github.com/nergilz/rmqcli/consumer"
	"github.com/nergilz/rmqcli/publisher"
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
	Conn      *amqp.Connection
	Ch        *amqp.Channel
	Consumer  *consumer.Consumer
	Publisher *publisher.Publisher
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
