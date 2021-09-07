package main

import (
	"context"
	"fmt"
	"time"

	gosqs "github.com/engelmi/go-sqs"
)

func main() {
	p, err := gosqs.NewProducer(gosqs.ProducerConfig{
		QueueConfig: gosqs.QueueConfig{
			Region:   "eu-central-1",
			Endpoint: "http://localhost:9324",
			Queue:    "simple_produce",
		},
		Timeout: 2 * time.Second,
	})
	if err != nil {
		panic(fmt.Sprintf("Could not create producer: %s", err.Error()))
	}

	p.Send(context.TODO(), gosqs.OutgoingMessage{
		Payload: []byte("This is a message"),
	})
}
