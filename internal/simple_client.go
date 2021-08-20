package internal

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

type Client struct {
	Sqs      *sqs.SQS
	QueueUrl string
}

func NewClient(region, endpoint, queue string) (*Client, error) {
	sess, err := session.NewSession(aws.NewConfig().WithRegion(region).WithEndpoint(endpoint))
	if err != nil {
		return nil, err
	}

	c := sqs.New(sess)

	output, err := c.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queue),
	})
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to get url for queue '%s'", queue))

	}
	if output.QueueUrl == nil {
		return nil, errors.New(fmt.Sprintf("QueueUrl was empty for queue '%s'", queue))
	}

	return &Client{
		Sqs:      sqs.New(sess),
		QueueUrl: *output.QueueUrl,
	}, nil
}
