package consumer

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type HandlerFoo func(d *amqp.Delivery)

type Consumer struct {
	Ctx           context.Context
	Ch            *amqp.Channel
	Conn          *amqp.Connection
	Deliveries    <-chan amqp.Delivery
	CloseConsumer chan struct{}
	HandlerFunc   HandlerFoo
	Reconnect     time.Duration
	// WorkersWg  *sync.WaitGroup
}

func NewConsumer(ctx context.Context, conn *amqp.Connection, h HandlerFoo) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("init channel: %s", err.Error())
	}

	return &Consumer{
		Ctx:           ctx,
		Ch:            ch,
		Conn:          conn,
		CloseConsumer: make(chan struct{}),
		HandlerFunc:   h,
		Reconnect:     time.Second,
		// WorkersWg: &sync.WaitGroup{},
	}, nil
}

func (c *Consumer) Run(queueName string) error {
	queue, err := c.Ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("queue declare: %s", err.Error())
	}

	delivery, err := c.Ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("channel consume: %s", err.Error())
	}

	c.Deliveries = delivery

	// for i := 0; i < c.cfg.Concurrency; i++ {
	// 	c.WorkersWg.Add(1)
	// 	go c.runWorker()
	// }

	go c.runWorker(c.Ctx)

	// var res bool
	// go c.runWorker(c.Ctx, boo(&res))
	// fmt.Println(res)

	return nil
}

// func foo(delivery *amqp.Delivery) bool {
// 	log.Printf(" Received msg from consumer: %s", delivery.Body)
// 	return true
// }

// func boo(req *bool) func(delivery *amqp.Delivery) bool {
// 	return func(delivery *amqp.Delivery) bool {
// 		log.Printf(" Received msg from consumer: %s", delivery.Body)
// 		*req = false
// 		return false
// 	}
// }

func (c *Consumer) runWorker(ctx context.Context) {
	// defer c.WorkersWg.Done()

	ticker := time.NewTicker(c.Reconnect)
	defer ticker.Stop()

	for {
		select {
		case <-c.CloseConsumer:
			return
		case delivery, isOpen := <-c.Deliveries:
			if !isOpen {
				return
			}
			go c.HandlerFunc(&delivery)
		case <-ctx.Done():
			return
		case <-ticker.C:
			// case <-c.StopConsumer // todo
		}
	}
}

func (c *Consumer) CloseChannel() error {
	err := c.Ch.Close()
	if err != nil {
		return fmt.Errorf("close channel: %s", err.Error())
	}
	return nil
}
