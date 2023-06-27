package consumer

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	Ch            *amqp.Channel
	Deliveries    <-chan amqp.Delivery
	CloseConsumer chan struct{}
	// WorkersWg  *sync.WaitGroup
}

func NewConsumer(ch *amqp.Channel) *Consumer {
	return &Consumer{
		Ch:            ch,
		CloseConsumer: make(chan struct{}),
		// WorkersWg: &sync.WaitGroup{},
	}
}

func (c *Consumer) Run(queueName string) error {
	queue, err := c.Ch.QueueDeclare(queueName, false, false, false, false, nil)
	errReceiveHandler(err, "Failed to declare a queue")

	delivery, err := c.Ch.Consume(queue.Name, "", true, false, false, false, nil)
	errReceiveHandler(err, "Failed consume queue")

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
			log.Printf(" Received msg: %s", delivery.Body)
			// case <-c.StopConsumer // todo
			// case <-ctx.Done() // todo
		}
	}

}

func errReceiveHandler(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
