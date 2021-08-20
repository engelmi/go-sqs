package sqsgo

import (
	"context"
	"encoding/json"
	"go-sqs/internal"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

type Producer interface {
	Send(ctx context.Context, msg OutgoingMessage) (*string, error)
}

type OutgoingMessage struct {
	DeduplicationId *string
	GroupId         *string
	Payload         interface{}
	Attributes      map[string]string
}

type producer struct {
	*internal.Client
	timeout time.Duration
}

func NewProducer(config ProducerConfig) (Producer, error) {
	c, err := internal.NewClient(config.Region, config.Endpoint, config.Queue)
	if err != nil {
		return nil, err
	}

	return &producer{
		Client:  c,
		timeout: config.Timeout,
	}, nil
}

func (p *producer) Send(ctx context.Context, msg OutgoingMessage) (*string, error) {
	attributes := make(map[string]*sqs.MessageAttributeValue)
	for key, value := range msg.Attributes {
		attributes[key] = &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}

	bytes, err := json.Marshal(msg.Payload)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal payload")
	}
	payload := string(bytes)

	input := sqs.SendMessageInput{
		QueueUrl:               aws.String(p.QueueUrl),
		MessageBody:            aws.String(payload),
		MessageAttributes:      attributes,
		MessageGroupId:         msg.GroupId,
		MessageDeduplicationId: msg.DeduplicationId,
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	output, err := p.Sqs.SendMessageWithContext(timeoutCtx, &input)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send message")
	}

	return output.MessageId, nil
}
