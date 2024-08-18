package handle

import (
	"encoding/json"
	"github.com/axliupore/judge/logic"
	"github.com/axliupore/judge/pkg/consumer"
	"github.com/axliupore/judge/pkg/log"
	"github.com/axliupore/judge/pkg/request"
	"github.com/axliupore/judge/pkg/response"
	"github.com/axliupore/judge/pkg/status"
	"github.com/nsqio/go-nsq"
)

// MessageHandler handles NSQ messages.
type MessageHandler struct{}

// HandleMessage implements the method to process messages.
func (h *MessageHandler) HandleMessage(m *nsq.Message) error {
	go func() {
		r := &request.JudgeRequest{}
		if err := json.Unmarshal(m.Body, r); err != nil {
			log.Logger.Error(err)
			return
		}

		l := logic.NewLogic()

		rsp := &response.Response{}

		res, st, err := l.ProcessJudgeRequest(r)
		rsp.Code = status.Code(st)
		rsp.Data = res
		if err != nil {
			rsp.Msg = err.Error()
			log.Logger.Error(err)
		}

		c := &consumer.Service{}

		if err = c.SendResponse(rsp, r.Nsq); err != nil {
			log.Logger.Error(err)
		}
	}()
	return nil
}
