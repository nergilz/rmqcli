package rmqcli

import (
	"context"
	"fmt"
	"time"

	"github.com/nergilz/rmqcli/consumer"
	"github.com/nergilz/rmqcli/declorator"
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
	conn             *amqp.Connection
	Consumer         *consumer.Consumer
	Publisher        *publisher.Publisher
	Declorator       *declorator.Declorator
	reconnectTimeout time.Duration
}

func InitRmqCli(ctx context.Context, url string, h consumer.HandlerFoo) (*RmqCli, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	p, err := publisher.NewPublisher(ctx, conn)
	if err != nil {
		return nil, fmt.Errorf("new publisher: %s", err.Error())
	}

	c, err := consumer.NewConsumer(ctx, conn, h)
	if err != nil {
		return nil, fmt.Errorf("new consumer: %s", err.Error())
	}

	cli := &RmqCli{
		conn:      conn,
		Publisher: p,
		Consumer:  c,
	}

	return cli, nil
}

func (rmq *RmqCli) CloseConnection() error {
	rmq.Consumer.CloseChannel()
	rmq.Publisher.CloseChannel()

	err := rmq.conn.Close()
	if err != nil {
		return fmt.Errorf("close connection: %s", err.Error())
	}
	return nil
}
