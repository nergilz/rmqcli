package consumer

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type HandlerFoo func(d *amqp.Delivery)

type Consumer struct {
	ctx           context.Context
	ch            *amqp.Channel
	conn          *amqp.Connection
	deliveries    <-chan amqp.Delivery
	closeConsumer chan struct{}
	handlerFunc   HandlerFoo
	reconnect     time.Duration
	// workersWg  *sync.WaitGroup
}

func NewConsumer(ctx context.Context, conn *amqp.Connection, h HandlerFoo) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("init channel: %s", err.Error())
	}

	return &Consumer{
		ctx:           ctx,
		ch:            ch,
		conn:          conn,
		closeConsumer: make(chan struct{}),
		handlerFunc:   h,
		reconnect:     time.Second,
		// workersWg: &sync.WaitGroup{},
	}, nil
}

func (c *Consumer) Consume(queueName string) error {
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

	go c.runWorker(c.ctx)

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
	err := c.ch.Close()
	if err != nil {
		return fmt.Errorf("close channel: %s", err.Error())
	}
	return nil
}

// example:

// var res bool
// go c.runWorker(c.Ctx, boo(&res))
// fmt.Println(res)

// func boo(req *bool) func(delivery *amqp.Delivery) bool {
// 	return func(delivery *amqp.Delivery) bool {
// 		log.Printf(" Received msg from consumer: %s", delivery.Body)
// 		*req = false
// 		return false
// 	}
// }
