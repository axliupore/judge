package handle

import (
	"errors"
	"github.com/axliupore/judge/logic"
	"github.com/axliupore/judge/pkg/request"
	"github.com/axliupore/judge/pkg/status"
	"github.com/axliupore/judge/pkg/xres"
	"github.com/gin-gonic/gin"
)

func JudgeServer(c *gin.Context) {
	r := &request.JudgeRequest{}
	if err := c.ShouldBindJSON(r); err != nil {
		xres.Http(c, status.Code(status.ParamsError), status.ParamsError, nil)
		return
	}
	if err := verify(r); err != nil {
		xres.Http(c, status.Code(status.ParamsError), err.Error(), nil)
		return
	}

	l := logic.NewLogic()
	res, st, err := l.Start(r)
	if err != nil {
		xres.Http(c, status.Code(st), err.Error(), nil)
		return
	}

	xres.Http(c, status.Code(st), "", res)
}

func verify(r *request.JudgeRequest) error {
	if r.Language == "" || r.Content == "" {
		return errors.New("programming language and code content must be filled in")
	}
	return nil
}
