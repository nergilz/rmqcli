package main

import (
	"context"
	"log"
	"time"

	"github.com/nergilz/rmqcli"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RmqConfig struct {
	Url   string
	Queue string
}

func main() {
	cfg := &RmqConfig{
		Url:   "amqp://guest:guest@localhost:5672/",
		Queue: "q_hello",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	body := "test rmq by cli v8.3"

	pub := &amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	}

	cli, err := rmqcli.InitRmqCli(ctx, cfg.Url)
	if err != nil {
		log.Println("error init cli:", err.Error())
	}

	err = cli.Publisher.Run(pub, cfg.Queue, "")
	if err != nil {
		log.Println("error run publish:", err.Error())
	}

	err = cli.Publisher.CloseCh()
	if err != nil {
		log.Println("error close ch:", err.Error())
	}

	time.Sleep(3 * time.Second)

	err = cli.Consumer.Run(cfg.Queue)
	if err != nil {
		log.Println("error run publish:", err.Error())
	}

	err = cli.Consumer.CloseCh()
	if err != nil {
		log.Println("error close ch:", err.Error())
	}

	cli.Conn.Close()

	time.Sleep(3 * time.Second)
}
