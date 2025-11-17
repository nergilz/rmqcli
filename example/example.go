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

	runPublisher(ctx, cli, cfg.Queue)

	runConsumer(ctx, cli, cfg.Queue)

	if err := cli.CloseConnection(); err != nil {
		log.Println("error close connection:", err.Error())
	}
}

func runPublisher(ctx context.Context, cli *rmqcli.RmqCli, queue string) {
	pub := &amqp.Publishing{ContentType: "text/plain", Body: []byte("test rmq msg by cli v9.7")}

	log.Println("run publisher...")

	err := cli.Publisher.Publish(ctx, pub, "", queue)
	if err != nil {
		log.Println("error run publish:", err.Error())
	}

	defer cli.Publisher.CloseChannel()
}

func runConsumer(ctx context.Context, cli *rmqcli.RmqCli, queue string) {
	log.Println("run consumer...")

	if err := cli.Consumer.Consume(ctx, queue); err != nil {
		log.Println("error run publish:", err.Error())
	}

	defer cli.Consumer.CloseChannel()

	time.Sleep(time.Second)
}

func handlerFoo(delivery *amqp.Delivery) {
	// do something
	log.Printf("Received msg: %s", delivery.Body)
}
