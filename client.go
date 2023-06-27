package rmqcli

import (
	"context"

	"github.com/nergilz/rmqcli/consumer"
	"github.com/nergilz/rmqcli/publisher"
	amqp "github.com/rabbitmq/amqp091-go"
)

// type RmqConfig struct {
// 	Url         string
// 	Queue       string
// 	Concurrency int
// 	Exchange    string
// 	RoutingKey  string
// }

type RmqCli struct {
	Conn      *amqp.Connection
	Consumer  *consumer.Consumer
	Publisher *publisher.Publisher
}

func InitRmqCli(ctx context.Context, url string) (*RmqCli, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	p, err := publisher.NewPublisher(conn)
	if err != nil {
		return nil, err
	}

	c, err := consumer.NewConsumer(conn)
	if err != nil {
		return nil, err
	}

	cli := &RmqCli{
		Conn:      conn,
		Publisher: p,
		Consumer:  c,
	}

	return cli, nil
}
