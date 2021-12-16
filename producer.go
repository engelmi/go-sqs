package gosqs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/engelmi/go-sqs/internal"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

type Producer interface {
	Send(ctx context.Context, msg OutgoingMessage) (*string, error)
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
	timeoutCtx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	output, err := p.Sqs.SendMessageWithContext(timeoutCtx, p.mapToAwsMessage(msg))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to send message")
	}

	return output.MessageId, nil
}

func (p *producer) mapToAwsMessage(msg OutgoingMessage) *sqs.SendMessageInput {
	attributes := make(map[string]*sqs.MessageAttributeValue)
	for key, value := range msg.Attributes {
		attributes[key] = &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}
	}

	payload := string(msg.Payload)
	return &sqs.SendMessageInput{
		QueueUrl:               aws.String(p.QueueUrl),
		MessageBody:            aws.String(payload),
		MessageAttributes:      attributes,
		MessageGroupId:         msg.GroupId,
		MessageDeduplicationId: msg.DeduplicationId,
	}
}

func MarshalToJson(payload interface{}) ([]byte, error) {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal payload")
	}
	return bytes, nil
}
