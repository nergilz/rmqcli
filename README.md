# rmqcli
## Go rabbitmq client

Wrapper abstraction reallization - [amqp091-go](https://github.com/rabbitmq/amqp091-go)

Use-cases for rabbitMq broker
* Binding queue by exchange and routing key
* Create handle func in consumer 
* Simple use consumer and publish 
* Graseful shutdown

## Install
```bash
$ go get "github.com/nergilz/rmqcli"
```

### Examples

#### create publisher
```go
func runPublisher(cli *rmqcli.RmqCli, queue string) {
	pub := &amqp.Publishing{ContentType: "text/plain", Body: []byte("test rmq msg by cli v9.7")}

	log.Println("run publisher...")

	err := cli.Publisher.Publish(pub, "", queue)
	if err != nil {
		log.Println("error run publish:", err.Error())
	}

	defer cli.Publisher.CloseChannel()
}
```

#### create consumer
```go
func runConsumer(cli *rmqcli.RmqCli, queue string) {
	log.Println("run consumer...")

	err := cli.Consumer.Consume(queue)
	if err != nil {
		log.Println("error run publish:", err.Error())
	}

	defer cli.Consumer.CloseChannel()

	time.Sleep(time.Second)
}
```

#### create hendler function
```go
func handlerFoo(delivery *amqp.Delivery) {
	// do something
	log.Printf("Received msg: %s", delivery.Body)
}
```

#### main.go
```go
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
```