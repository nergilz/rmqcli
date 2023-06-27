package consumer

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	Ch            *amqp.Channel
	Conn          *amqp.Connection
	Deliveries    <-chan amqp.Delivery
	CloseConsumer chan struct{}
	// WorkersWg  *sync.WaitGroup
}

func NewConsumer(conn *amqp.Connection) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("init channel: %s", err.Error())
	}
	// defer ch.Close()

	return &Consumer{
		Ch:            ch,
		Conn:          conn,
		CloseConsumer: make(chan struct{}),
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
		return fmt.Errorf("ch consume: %s", err.Error())
	}

	c.Deliveries = delivery

	// for i := 0; i < c.cfg.Concurrency; i++ {
	// 	c.WorkersWg.Add(1)
	// 	go c.runWorker()
	// }

	go c.runWorker()

	return nil
}

func (c *Consumer) runWorker() {
	// defer c.WorkersWg.Done()

	for {
		select {
		case <-c.CloseConsumer:
			return
		case delivery, isOpen := <-c.Deliveries:
			if !isOpen {
				return
			}
			log.Printf(" Received msg from consumer: %s", delivery.Body)
			// case <-c.StopConsumer // todo
			// case <-ctx.Done() // todo
		}
	}

}

func (c *Consumer) CloseCh() error {
	err := c.Ch.Close()
	if err != nil {
		return fmt.Errorf("close channel: %s", err.Error())
	}
	return nil
}
