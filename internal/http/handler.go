package http

import (
	"errors"
	"github.com/axliupore/judge/pkg/log"
	"github.com/axliupore/judge/pkg/request"
	"github.com/axliupore/judge/pkg/response"
	"github.com/axliupore/judge/pkg/status"
	"github.com/axliupore/judge/pkg/verify"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sync"
)

func Run(c *gin.Context) {
	r, err := parse[*request.Run](c)
	if err != nil {
		return
	}

	process[*request.Run, *response.Run](c, r, func(r *request.Run) (*response.Run, error) {
		return server.judge.Run(r)
	})
}

func Build(c *gin.Context) {
	r, err := parse[*request.Build](c)
	if err != nil {
		return
	}

	process[*request.Build, *response.Build](c, r, func(r *request.Build) (*response.Build, error) {
		return server.judge.Build(r)
	})
}

func Exec(c *gin.Context) {
	r, err := parse[*request.Exec](c)
	if err != nil {
		return
	}

	process[*request.Exec, *response.Exec](c, r, func(r *request.Exec) (*response.Exec, error) {
		return server.judge.Exec(r)
	})
}

func Delete(c *gin.Context) {
	r, err := parse[*request.Delete](c)
	if err != nil {
		return
	}

	server.judge.Delete(r)
}

// parse parses the request body into a slice of type R.
func parse[R any](c *gin.Context) (r []R, err error) {
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": status.ServerError})
		return
	}

	if err = sonic.Unmarshal(b, &r); err != nil {
		log.Logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": status.ServerError})
		return nil, err
	}

	// Verify the parsed request objects.
	if err = verify.Slice[R](r); err != nil || len(b) == 0 || len(r) == 0 {
		log.Logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": status.ParamsError})
		return nil, errors.New("param error")
	}

	return r, err
}

// process processes a slice of requests of type R and returns a slice of responses of type S.
func process[R any, S any](c *gin.Context, r []R, f func(R) (S, error)) {
	var err error
	res := make([]S, len(r))

	// Use a WaitGroup to wait for all goroutines to finish.
	wg := &sync.WaitGroup{}
	for i, v := range r {
		wg.Add(1)

		// Submit the task to the server's pool.
		if e := server.pool.Submit(func() {
			defer wg.Done()

			var e error
			res[i], e = f(v)
			if e != nil {
				err = e
			}
		}); e != nil {
			log.Logger.Error(e)
			err = e
		}
	}

	// Wait for all tasks to complete.
	wg.Wait()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": status.ServerError})
		return
	}

	d, err := sonic.Marshal(res)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": status.ServerError})
		return
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(http.StatusOK, string(d))
}
