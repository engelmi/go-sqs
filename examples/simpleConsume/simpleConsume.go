package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	gosqs "github.com/engelmi/go-sqs"
	"github.com/engelmi/go-sqs/examples"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	handler := func(ctx context.Context, receivedMsg gosqs.IncomingMessage) error {
		defer wg.Done()

		fmt.Println(fmt.Sprintf("Got message '%s'", *receivedMsg.Body))
		return nil
	}

	consumer, err := gosqs.NewConsumer(gosqs.ConsumerConfig{
		QueueConfig: gosqs.QueueConfig{
			Region:   "eu-central-1",
			Endpoint: "http://localhost:9324",
			Queue:    "simple_consume",
		},
		PollTimeout:         10 * time.Second,
		AckTimeout:          2 * time.Second,
		MaxNumberOfMessages: 10,
	}, handler)
	if err != nil {
		panic(fmt.Sprintf("Could not create consumer: %s", err.Error()))
	}
	go consumer.StartListening(context.Background(), nil)

	examples.NewGopher("simple_consume", "PushingGopher-1", 1*time.Second).PushMessage("Hello World!")

	wg.Wait()
}
