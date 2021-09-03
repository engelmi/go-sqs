package gosqs

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Region   string
	Endpoint string
	Queue    string
}

type ProducerConfig struct {
	Config
	Timeout time.Duration
}

type ConsumerConfig struct {
	Config
	PollTimeout         time.Duration
	AckTimeout          time.Duration
	MaxNumberOfMessages int64
	Logger              logrus.Logger
}
