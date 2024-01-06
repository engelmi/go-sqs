package internal

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

type Client struct {
	Sqs      *sqs.SQS
	QueueUrl string
}

func extractAccountIDFromSQSQueueURL(queueURL string) (string, error) {
	// Define the regular expression pattern to match the AWS account ID
	pattern := `https://sqs\.[a-zA-Z0-9-]+\.amazonaws\.com/(\d+)/`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find matches in the queue URL
	matches := re.FindStringSubmatch(queueURL)

	// The AWS account ID is captured in the first submatch
	accountId := matches[1]

	return accountId, nil
}

func NewClient(region, endpoint, queue string) (*Client, error) {
	sess, err := session.NewSession(aws.NewConfig().WithRegion(region).WithEndpoint(endpoint))
	if err != nil {
		return nil, err
	}

	c := sqs.New(sess)

	queueOwnerAWSAccountId, _ := extractAccountIDFromSQSQueueURL(endpoint)

	output, err := c.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName:              aws.String(queue),
		QueueOwnerAWSAccountId: aws.String(queueOwnerAWSAccountId),
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
