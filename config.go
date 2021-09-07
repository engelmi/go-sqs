package gosqs

import (
	"time"

	"github.com/sirupsen/logrus"
)

type QueueConfig struct {
	Region   string
	Endpoint string
	Queue    string
}

type ProducerConfig struct {
	QueueConfig
	Timeout time.Duration
}

type ConsumerConfig struct {
	QueueConfig
	PollTimeout         time.Duration
	AckTimeout          time.Duration
	MaxNumberOfMessages int64
	Logger              logrus.Logger
}
