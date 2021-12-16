package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	gosqs "github.com/engelmi/go-sqs"
	"github.com/engelmi/go-sqs/examples"
	"github.com/sirupsen/logrus"
)

func main() {
	wg := &sync.WaitGroup{}

	handler := func(ctx context.Context, receivedMsg gosqs.IncomingMessage) error {
		fmt.Println(fmt.Sprintf("Got message '%s'", *receivedMsg.Body))
		panic(fmt.Sprintf("Handler panic for message '%s'", *receivedMsg.Body))
	}

	consumer, err := gosqs.NewConsumer(gosqs.ConsumerConfig{
		QueueConfig: gosqs.QueueConfig{
			Region:   "eu-central-1",
			Endpoint: "http://localhost:9324",
			Queue:    "panic_queue",
		},
		PollTimeout:         10 * time.Second,
		AckTimeout:          2 * time.Second,
		MaxNumberOfMessages: 10,
		Logger:              *logrus.New(),
	}, handler)
	if err != nil {
		panic(fmt.Sprintf("Could not create consumer: %s", err.Error()))
	}
	go consumer.StartListening(context.Background(), nil)

	examples.NewGopher("panic_queue", "PushingGopher-1", 1*time.Second).PushMessage("Hello World! No. 1")
	examples.NewGopher("panic_queue", "PushingGopher-1", 1*time.Second).PushMessage("Hello World! No. 2")

	wg.Add(1)
	wg.Wait()
}
