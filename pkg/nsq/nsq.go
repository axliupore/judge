package nsq

import "github.com/nsqio/go-nsq"

// NewConsumer creates a new consumer instance for NSQ.
func NewConsumer(topic, channel string) (*nsq.Consumer, error) {
	consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// NewProducer creates a new producer instance for NSQ.
func NewProducer(nsqdAddress string) (*nsq.Producer, error) {
	producer, err := nsq.NewProducer(nsqdAddress, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	return producer, nil
}
