package consumer

import (
	"encoding/json"
	"github.com/axliupore/judge/pkg/log"
	"github.com/axliupore/judge/pkg/nsq"
	"github.com/axliupore/judge/pkg/response"
	"time"
)

// Service defines the consumer service structure.
type Service struct{}

// SendResponse sends evaluation response messages to the message queue.
func (c *Service) SendResponse(res *response.Response, resultQueue string) error {
	maxRetries := 3
	retryDelay := 10 * time.Millisecond

	// Serialize the response struct to JSON format.
	message, err := json.Marshal(res)
	if err != nil {
		return err
	}

	// Create a new producer instance.
	// todo 修改为参数传递的
	producer, err := nsq.NewProducer("")
	if err != nil {
		return err
	}
	defer producer.Stop()

	for i := 0; i < maxRetries; i++ {
		// Publish the message to the specified result queue.
		err = producer.Publish(resultQueue, message)
		if err != nil {
			time.Sleep(retryDelay)
			continue
		}
		log.Logger.Infof("Producer send message: %s\n", message)
		break
	}
	return nil
}
