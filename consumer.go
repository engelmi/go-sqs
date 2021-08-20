package sqsgo

import (
	"context"
	"fmt"
	"go-sqs/internal"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Consumer interface {
	StartListening(ctx context.Context, wg *sync.WaitGroup)
	StopListening() error
}

type MessageHandler func(ctx context.Context, msg *sqs.Message) error

type consumer struct {
	*internal.Client
	logger logrus.Logger

	handlerFunc MessageHandler
	cancelFunc  context.CancelFunc

	maxNumberOfMessages int64
	pollTimeout         int64
	ackTimeout          time.Duration
}

func NewConsumer(config ConsumerConfig, handlerFunc MessageHandler) (Consumer, error) {
	client, err := internal.NewClient(config.Region, config.Endpoint, config.Queue)
	if err != nil {
		return nil, err
	}

	return &consumer{
		Client:              client,
		logger:              config.logger,
		handlerFunc:         handlerFunc,
		pollTimeout:         int64(config.PollTimeout.Seconds()),
		ackTimeout:          config.AckTimeout,
		maxNumberOfMessages: config.MaxNumberOfMessages,
	}, nil
}

func (c *consumer) StartListening(ctx context.Context, wg *sync.WaitGroup) {
	cancelCtx, cancel := context.WithCancel(ctx)
	c.cancelFunc = cancel

	if wg == nil {
		wg = &sync.WaitGroup{}
	}
	wg.Add(1)
	for {
		select {
		case <-cancelCtx.Done():
			wg.Done()
			return
		default:
			c.listen(context.Background())
		}
	}
}

func (c *consumer) StopListening() error {
	if c.cancelFunc == nil {
		errors.New("Consumer has not been started yet")
	}
	c.cancelFunc()
	return nil
}

func (c *consumer) listen(ctx context.Context) {
	msgs, err := c.pollMessages(ctx)
	if err != nil {
		c.logger.WithContext(ctx).WithError(err).Error("Error polling sqs messages")
		return
	}

	for _, msg := range msgs {
		c.processMessage(ctx, msg)
	}
}

func (c *consumer) pollMessages(ctx context.Context) ([]*sqs.Message, error) {
	res, err := c.Sqs.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(c.QueueUrl),
		MaxNumberOfMessages:   aws.Int64(c.maxNumberOfMessages),
		WaitTimeSeconds:       aws.Int64(c.pollTimeout),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return res.Messages, nil
}

func (c *consumer) processMessage(ctx context.Context, msg *sqs.Message) {
	err := c.callHandlerFunc(ctx, msg)
	if err != nil {
		c.logger.WithContext(ctx).WithError(err).Error("Error handling sqs message")
		return
	}

	err = c.acknowledgeMessage(ctx, msg)
	if err != nil {
		c.logger.WithContext(ctx).WithError(err).Error("Error acknowledging sqs message")
	}
}

func (c *consumer) callHandlerFunc(ctx context.Context, msg *sqs.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("SQS-Handler function paniced: %v", r))
		}
	}()

	return c.handlerFunc(ctx, msg)
}

func (c *consumer) acknowledgeMessage(ctx context.Context, msg *sqs.Message) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, c.ackTimeout)
	defer cancel()

	_, err := c.Sqs.DeleteMessageWithContext(timeoutCtx, &sqs.DeleteMessageInput{
		QueueUrl:      &c.QueueUrl,
		ReceiptHandle: msg.ReceiptHandle,
	})

	return errors.WithStack(err)
}
