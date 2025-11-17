package consumer

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type HandlerFoo func(d *amqp.Delivery)

type Consumer struct {
	ch            *amqp.Channel
	conn          *amqp.Connection
	deliveries    <-chan amqp.Delivery
	closeConsumer chan struct{}
	handlerFunc   HandlerFoo
	reconnect     time.Duration
	// workersWg  *sync.WaitGroup
}

func NewConsumer(conn *amqp.Connection, h HandlerFoo) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("init channel: %s", err.Error())
	}

	return &Consumer{
		ch:            ch,
		conn:          conn,
		closeConsumer: make(chan struct{}),
		handlerFunc:   h,
		reconnect:     time.Second,
		// workersWg: &sync.WaitGroup{},
	}, nil
}

func (c *Consumer) Consume(ctx context.Context, queueName string) error {
	queue, err := c.ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %s", err.Error())
	}

	delivery, err := c.ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("channel consume: %s", err.Error())
	}

	c.deliveries = delivery

	// for i := 0; i < c.cfg.Concurrency; i++ {
	// 	c.workersWg.Add(1)
	// 	go c.runWorker()
	// }

	go c.runWorker(ctx)

	return nil
}

func (c *Consumer) runWorker(ctx context.Context) {
	// defer c.workersWg.Done()

	ticker := time.NewTicker(c.reconnect)
	defer ticker.Stop()

	for {
		select {
		case <-c.closeConsumer:
			return
		case delivery, isOpen := <-c.deliveries:
			if !isOpen {
				return
			}
			go c.handlerFunc(&delivery)
		case <-ctx.Done():
			return
		case <-ticker.C:
			// case <-c.StopConsumer // todo
		}
	}
}

func (c *Consumer) CloseChannel() error {
	if err := c.ch.Close(); err != nil {
		return fmt.Errorf("close channel: %s", err.Error())
	}

	return nil
}
