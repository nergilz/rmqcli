package main

import (
	"context"
	"log"
	"time"

	"github.com/nergilz/rmqcli"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg := &rmqcli.RmqConfig{
		Url:   "amqp://guest:guest@localhost:5672/",
		Queue: "hello_rmqcli",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cli, err := rmqcli.InitRmqCli(ctx, cfg.Url, handlerFoo)
	if err != nil {
		log.Println("error init cli:", err.Error())
	}

	runPublisher(cli, cfg.Queue)

	runConsumer(cli, cfg.Queue)

	err = cli.CloseConnection()
	if err != nil {
		log.Println("error close connection:", err.Error())
	}
}

func handlerFoo(delivery *amqp.Delivery) {
	log.Printf("Received msg: %s", delivery.Body)
}

func runPublisher(cli *rmqcli.RmqCli, queue string) {
	pub := &amqp.Publishing{ContentType: "text/plain", Body: []byte("test rmq msg by cli v9.5")}

	log.Println("run publisher...")

	err := cli.Publisher.Run(pub, queue, "")
	if err != nil {
		log.Println("error run publish:", err.Error())
	}

	err = cli.Publisher.CloseChannel()
	if err != nil {
		log.Println("error close ch:", err.Error())
	}
}

func runConsumer(cli *rmqcli.RmqCli, queue string) {
	log.Println("run consumer...")

	err := cli.Consumer.Run(queue)
	if err != nil {
		log.Println("error run publish:", err.Error())
	}

	err = cli.Consumer.CloseChannel()
	if err != nil {
		log.Println("error close channel:", err.Error())
	}

	time.Sleep(2 * time.Second)
}
