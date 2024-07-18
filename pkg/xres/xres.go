package xres

import (
	"github.com/axliupore/judge/pkg/response"
	"github.com/axliupore/judge/pkg/status"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Http(c *gin.Context, code int, msg string, data interface{}) {
	if msg == "" {
		msg = status.Status(code)
	}
	c.JSON(http.StatusOK, response.Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}
