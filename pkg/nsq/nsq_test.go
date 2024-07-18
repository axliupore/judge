package nsq

import (
	"encoding/json"
	"github.com/axliupore/judge/pkg/request"
	"github.com/nsqio/go-nsq"
	"testing"
	"time"
)

func TestProducer(t *testing.T) {
	producer, err := NewProducer("127.0.0.1:4150")
	if err != nil {
		t.Fatal(err)
	}
	defer producer.Stop()

	req := request.JudgeRequest{
		Language: "golang",
		Content:  "package main\n\nimport \"fmt\"\n\nfunc main() {\n  var a, b int\n  fmt.Scanf(\"%d%d\", &a, &b)\n  fmt.Printf(\"%d\", a + b)\n}\n",
		Input:    "1 2",
		Nsq:      "axliupore",
	}
	mesBody, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	err = producer.Publish("judge_topic", mesBody)
	if err != nil {
		t.Fatal(err)
	}

	consumer, err := NewConsumer("axliupore", "channel")
	if err != nil {
		t.Fatal(err)
	}
	defer consumer.Stop()

	consumer.AddHandler(&MessageHandler{
		t: t,
	})

	err = consumer.ConnectToNSQD("127.0.0.1:4150")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)
}

type MessageHandler struct {
	t *testing.T
}

func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	var req request.JudgeResponse
	err := json.Unmarshal(m.Body, &req)
	if err != nil {
		h.t.Error(err)
		return err
	}
	h.t.Logf("Received message: %+v", req)
	return nil
}
