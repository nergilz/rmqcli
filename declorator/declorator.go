package declorator

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Declorator struct {
	ch       *amqp.Channel
	Exchange ExchangeSource
	Binding  BindingSource
	Queue    QueueSource
}

type ExchangeSource struct {
	Name string
	Type string
	Args amqp.Table
}

type BindingSource struct {
	ExchangeName string
	RoutingKey   string
	QueueName    string
	Args         amqp.Table
}

type QueueSource struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Args       amqp.Table
}

func NewDeclarator(ch *amqp.Channel) *Declorator {
	return &Declorator{
		ch: ch,
	}
}

func (d *Declorator) Declare() error {
	err := d.ch.ExchangeDeclare(d.Exchange.Name, d.Exchange.Type, true, false, false, false, d.Exchange.Args)
	if err != nil {
		return fmt.Errorf("declorator exchange: %s", err.Error())
	}

	_, err = d.ch.QueueDeclare(d.Queue.Name, d.Queue.Durable, d.Queue.AutoDelete, false, false, d.Exchange.Args)
	if err != nil {
		return fmt.Errorf("declorator queue: %s", err.Error())
	}

	err = d.ch.QueueBind(d.Binding.QueueName, d.Binding.RoutingKey, d.Binding.ExchangeName, false, d.Exchange.Args)
	if err != nil {
		return fmt.Errorf("declorator exchange: %s", err.Error())
	}

	return nil
}

func (d *Declorator) CloseChannel() error {
	err := d.ch.Close()
	if err != nil {
		return fmt.Errorf("declorator close channel: %s", err.Error())
	}
	return nil
}
