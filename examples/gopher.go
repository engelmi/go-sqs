package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	gosqs "github.com/engelmi/go-sqs"
)

type Gopher struct {
	name         string
	queueName    string
	recoveryTime time.Duration
}

func NewGopher(queueName string, gopherName string, recoveryTime time.Duration) Gopher {
	return Gopher{
		name:         gopherName,
		queueName:    queueName,
		recoveryTime: recoveryTime,
	}
}

func (g Gopher) StartPushingMessages() {
	p := g.setupProducer()

	for i := 0; ; i++ {
		g.sendMessage(p, fmt.Sprintf("%s: This is message no. %d", g.name, i))
		time.Sleep(g.recoveryTime)
	}
}

func (g Gopher) PushMessage(msg string) {
	p := g.setupProducer()
	g.sendMessage(p, msg)
}

func (g Gopher) setupProducer() gosqs.Producer {
	p, err := gosqs.NewProducer(gosqs.ProducerConfig{
		QueueConfig: gosqs.QueueConfig{
			Region:   "eu-central-1",
			Endpoint: "http://localhost:9324",
			Queue:    g.queueName,
		},
		Timeout: 2 * time.Second,
	})
	if err != nil {
		panic(fmt.Sprintf("Could not create gopher: %s", err.Error()))
	}
	return p
}

func (g Gopher) sendMessage(p gosqs.Producer, msg string) {
	p.Send(context.TODO(), gosqs.OutgoingMessage{
		Payload: []byte(msg),
	})
}

func (g Gopher) Consume() {
	consumer, err := gosqs.NewConsumer(gosqs.ConsumerConfig{
		QueueConfig: gosqs.QueueConfig{
			Region:   "eu-central-1",
			Endpoint: "http://localhost:9324",
			Queue:    g.queueName,
		},
		PollTimeout:         10 * time.Second,
		AckTimeout:          2 * time.Second,
		MaxNumberOfMessages: 10,
	}, func(ctx context.Context, receivedMsg *sqs.Message) error {
		msg := fmt.Sprintf("%s: Got message '%s'", g.name, *receivedMsg.Body)
		fmt.Println(msg)

		time.Sleep(g.recoveryTime)
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("Could not create consumer: %s", err.Error()))
	}
	consumer.StartListening(context.Background(), nil)
}
