package judge

import (
	"github.com/axliupore/judge/judge"
	"github.com/axliupore/judge/pkg/client"
	"github.com/axliupore/judge/pkg/cmd"
	"github.com/axliupore/judge/pkg/log"
	"github.com/axliupore/judge/pkg/request"
	"github.com/axliupore/judge/pkg/response"
	"github.com/axliupore/judge/pkg/status"
	"github.com/jinzhu/copier"
	"sync"
	"time"
)

// Constants for default values
const (
	DefaultCpuLimit    int64 = 5000 * 1000 * 1000
	DefaultMemoryLimit int64 = 512 * 1024 * 1024
	DefaultStackLimit  int32 = 8 * 1024 * 1024
	DefaultProcLimit   int32 = 100
)

// Params struct holds the parameters for the judge operations
type Params struct {
	Code        string // Code to be executed
	Input       string // Input for the execution
	CpuLimit    int64  // CPU limit for the execution
	MemoryLimit int64  // Memory limit for the execution
	StackLimit  int32  // Stack limit for the execution
	ProcLimit   int32  // Process limit for the execution
	FileId      string // File ID for the execution
}

// Server struct holds the client service for judge operations
type Server struct {
	c *client.Service
}

// NewServer initializes and returns a new Server instance
func NewServer() *Server {
	return &Server{
		c: client.NewService(),
	}
}

func (s *Server) Run(r *request.Run) (*response.Run, error) {

	res := &response.Run{}
	c, p, j := buildCmd[*request.Run](r, r.Language)
	if j == nil {
		res.Status = status.ParamsError
		return res, nil
	}

	rsp, err := s.run(c, p, j)
	if err != nil {
		log.Logger.Error(err)
		return nil, err
	}

	if !j.IsBuild() || rsp.Status != status.Accepted {
		s.setOutput(res, rsp, rsp.Files["stderr"], rsp.Files["stdout"])
		return res, err
	}

	p.FileId = rsp.FileIds[j.ExecFile()[0]]

	rbs, err := s.exec(c, p, j)
	if err != nil {
		log.Logger.Error(err)
		return nil, err
	}

	go func() {
		_ = s.c.Delete(p.FileId)
	}()

	s.setOutput(res, rbs, rbs.Files["stderr"], rbs.Files["stdout"])
	return res, err
}

func (s *Server) Build(r *request.Build) (*response.Build, error) {

	res := &response.Build{}
	c, p, j := buildCmd[*request.Build](r, r.Language)
	if j == nil || !j.IsBuild() {
		res.Status = status.ParamsError
		return res, nil
	}

	rsp, err := s.build(c, p, j)
	if err != nil {
		log.Logger.Error(err)
		return nil, err
	}
	fileId := rsp.FileIds[j.ExecFile()[0]]

	go s.delete(fileId)()

	s.setOutput(res, rsp, rsp.Files["stderr"], fileId)
	return res, nil
}

func (s *Server) Exec(r *request.Exec) (*response.Exec, error) {

	res := &response.Exec{}
	c, p, j := buildCmd[*request.Exec](r, r.Language)
	if j == nil {
		res.Status = status.ParamsError
		return res, nil
	}

	rsp, err := s.exec(c, p, j)
	if err != nil {
		log.Logger.Error(err)
		return nil, err
	}

	s.setOutput(res, rsp, rsp.Files["stderr"], rsp.Files["stdout"])
	return res, nil
}

// Delete handles the Delete request and deletes the specified files
func (s *Server) Delete(r []*request.Delete) {

	wg := &sync.WaitGroup{}
	for _, d := range r {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.c.Delete(d.FileId); err != nil {
				log.Logger.Error(err)
			}
		}()
	}
	wg.Wait()
}

func (s *Server) run(r *cmd.Request, p *Params, j judge.Judge) (*cmd.Response, error) {

	r.Cmd[0].Files = []map[string]interface{}{{"content": ""}, {"name": "stdout", "max": 102400}, {"name": "stderr", "max": 102400}}
	r.Cmd[0].CopyIn = map[string]map[string]string{j.RunFile(): {"content": p.Code}}

	return s.c.Send(r)
}

func (s *Server) build(r *cmd.Request, p *Params, j judge.Judge) (*cmd.Response, error) {

	r.Cmd[0].Files = []map[string]interface{}{{"content": ""}, {"name": "stdout", "max": 102400}, {"name": "stderr", "max": 102400}}
	r.Cmd[0].CopyIn = map[string]map[string]string{j.RunFile(): {"content": p.Code}}

	return s.c.Send(r)
}

func (s *Server) exec(r *cmd.Request, p *Params, j judge.Judge) (*cmd.Response, error) {

	r.Cmd[0].Files = []map[string]interface{}{{"content": p.Input}, {"name": "stdout", "max": 102400}, {"name": "stderr", "max": 102400}}
	r.Cmd[0].CopyIn = map[string]map[string]string{j.ExecFile()[0]: {"fileId": p.FileId}}
	r.Cmd[0].Args = j.ExecArgs()
	r.Cmd[0].CopyOutCached = nil

	return s.c.Send(r)
}

// setOutput sets the output for the response
func (s *Server) setOutput(res response.Response, rsp *cmd.Response, stderr string, stdout string) {

	_ = copier.Copy(res, rsp)
	if rsp.Status != status.Accepted {
		res.SetOutput(stderr)
	} else {
		res.SetOutput(stdout)
	}
}

// delete schedules the deletion of a file after a delay
func (s *Server) delete(fileId string) func() {
	return func() {
		time.Sleep(10 * time.Minute)
		_ = s.c.Delete(fileId)
	}
}

// buildCmd builds the command request and returns the request, params, and judge
func buildCmd[R any](r R, l string) (*cmd.Request, *Params, judge.Judge) {

	j := judge.NewJudge(l)
	if j == nil {
		return nil, nil, nil
	}

	p := &Params{}
	_ = copier.Copy(p, r)

	c := &cmd.Cmd{
		Args:          j.RunArgs(),
		Env:           j.Env(),
		CpuLimit:      p.CpuLimit,
		MemoryLimit:   p.MemoryLimit,
		StackLimit:    p.StackLimit,
		ProcLimit:     p.ProcLimit,
		CopyOut:       []string{"stdout", "stderr"},
		CopyOutCached: j.ExecFile(),
	}

	if c.CpuLimit == 0 {
		c.CpuLimit = DefaultCpuLimit
	}
	if c.MemoryLimit == 0 {
		c.MemoryLimit = DefaultMemoryLimit
	}
	if c.StackLimit == 0 {
		c.StackLimit = DefaultStackLimit
	}
	if c.ProcLimit == 0 {
		c.ProcLimit = DefaultProcLimit
	}
	rc := &cmd.Request{Cmd: []*cmd.Cmd{c}}

	return rc, p, j
}
