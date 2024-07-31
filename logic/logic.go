package logic

import (
	"errors"
	"github.com/axliupore/judge/judge"
	"github.com/axliupore/judge/pkg/client"
	"github.com/axliupore/judge/pkg/exec"
	"github.com/axliupore/judge/pkg/log"
	"github.com/axliupore/judge/pkg/request"
	"github.com/axliupore/judge/pkg/response"
	"github.com/axliupore/judge/pkg/run"
	"github.com/axliupore/judge/pkg/status"
	"github.com/jinzhu/copier"
)

// Logic struct represents the logic handler for processing judge requests.
type Logic struct {
	service *client.Service
}

// Params struct encapsulates compilation parameters.
type Params struct {
	ExecArgs    []string // Exec Args
	RunArgs     []string // Run Args
	Env         []string // Environment variables
	Content     string   // Content
	Input       string   // Input
	CpuLimit    int64    // CPU limit
	MemoryLimit int64    // Memory limit
	StackLimit  int64    // Stack limit
	ProcLimit   int32    // Process limit
	RunFile     string   // Run file name
	ExecFile    string   // Exec file name
	FileId      string   // File ID
	Language    string   // code language
	IsExec      bool     // Is exec file
}

// Constants for default values
const (
	DefaultCpuLimit    int64 = 5000 * 1000 * 1000
	DefaultMemoryLimit int64 = 512 * 1024 * 1024
	DefaultStackLimit  int64 = 8 * 1024 * 1024
	DefaultProcLimit   int32 = 100
)

// NewParams initializes a new Params struct with default values.
func NewParams() *Params {
	return &Params{
		CpuLimit:    DefaultCpuLimit,
		MemoryLimit: DefaultMemoryLimit,
		StackLimit:  DefaultStackLimit,
		ProcLimit:   DefaultProcLimit,
	}
}

// NewLogic initializes a new Logic struct with a client.Service instance.
func NewLogic() *Logic {
	return &Logic{
		service: client.NewService(),
	}
}

// ProcessJudgeRequest processes a JudgeRequest and returns a JudgeResponse or an error.
func (l *Logic) ProcessJudgeRequest(r *request.JudgeRequest) (*response.JudgeResponse, string, error) {
	params, st, err := l.copyParamsFromRequest(r)
	if err != nil {
		return nil, st, err
	}

	// Run the code with the given params
	runRes, err := l.runCode(params)
	if err != nil {
		log.Logger.Errorf("failed to run code: %v", err)
		return nil, status.InternalError, err
	}

	// Judge the run result
	accept, fileId, output, time, memory := l.evaluateRunResult(params.ExecFile, runRes)

	// Handle non-accepted results
	if !accept {
		if runRes.Status == status.TimeLimitExceeded || runRes.Status == status.MemoryLimitExceeded {
			log.Logger.Errorf("run failed: %v", runRes.Status)
			return nil, runRes.Status, nil
		}
		log.Logger.Errorf("run failed: %v", output)
		return nil, runRes.Status, errors.New(output)
	}

	// If it's not an executable task, return JudgeResponse with output, time, and memory
	if !params.IsExec || r.NotExec {
		return &response.JudgeResponse{
			Output: output,
			Time:   time,
			Memory: memory,
			FileId: fileId,
		}, runRes.Status, nil
	}

	params.FileId = fileId

	return l.executeWithParams(params)
}

// ExecuteJudgeRequest executes a JudgeRequest and returns a JudgeResponse or an error.
func (l *Logic) ExecuteJudgeRequest(r *request.JudgeRequest) (*response.JudgeResponse, string, error) {
	params, st, err := l.copyParamsFromRequest(r)
	if err != nil {
		return nil, st, err
	}

	return l.executeWithParams(params)
}

// removeFile deletes the file identified by fileId.
func (l *Logic) removeFile(fileId string) {
	if fileId == "" {
		return
	}
	if err := l.service.DeleteRequest(fileId); err != nil {
		log.Logger.Errorf("removeFile error: %v", err)
	}
}

func (l *Logic) executeWithParams(params *Params) (*response.JudgeResponse, string, error) {
	defer l.removeFile(params.FileId)
	// Execute the code with the given params
	execRes, err := l.executeCode(params)
	if err != nil {
		log.Logger.Errorf("failed to execute code: %v", err)
		return nil, execRes.Status, err
	}

	// Judge the exec result
	st, output, time, memory := l.evaluateExecResult(execRes)
	if st != status.Accepted {
		log.Logger.Errorf("execution failed: %v %v", st, output)
		return nil, st, errors.New(output)
	}

	// Return JudgeResponse with output, time, and memory
	return &response.JudgeResponse{
		Output: output,
		Time:   time,
		Memory: memory,
	}, st, nil
}

func (l *Logic) copyParamsFromRequest(r *request.JudgeRequest) (*Params, string, error) {
	j := judge.NewJudge(r.Language)
	if j == nil {
		return nil, status.ParamsError, errors.New("unsupported programming language")
	}

	// Initialize params with default values
	p := NewParams()

	// Copy request data to params
	if err := copier.Copy(p, r); err != nil {
		return nil, status.InternalError, err
	}

	p.ExecArgs = j.ExecArgs()
	p.RunArgs = j.RunArgs()
	p.Env = j.Env()
	p.IsExec = j.IsExec()
	p.RunFile = j.RunFile()
	p.ExecFile = j.ExecFile()
	return p, "", nil
}

// runCode sends a run request with given params and returns the run response or an error.
func (l *Logic) runCode(params *Params) (*run.Response, error) {
	cmd := run.Cmd{
		Args:        params.RunArgs,
		Env:         params.Env,
		Files:       []map[string]interface{}{{"content": params.Input}, {"name": "stdout", "max": 10240}, {"name": "stderr", "max": 10240}},
		CpuLimit:    params.CpuLimit,
		MemoryLimit: params.MemoryLimit,
		StackLimit:  params.StackLimit,
		ProcLimit:   params.ProcLimit,
		CopyIn:      map[string]map[string]string{params.RunFile: {"content": params.Content}},
		CopyOut:     []string{"stdout", "stderr"},
	}

	if params.ExecFile != "" {
		cmd.CopyOutCached = []string{params.ExecFile}
	}

	runRequest := &run.Request{
		Cmd: []run.Cmd{cmd},
	}

	return l.service.SendRunRequest(runRequest)
}

// executeCode sends an exec request with given params and returns the exec response or an error.
func (l *Logic) executeCode(params *Params) (*exec.Response, error) {
	cmd := exec.Cmd{
		Args:        params.ExecArgs,
		Env:         params.Env,
		Files:       []map[string]interface{}{{"content": params.Input}, {"name": "stdout", "max": 10240}, {"name": "stderr", "max": 10240}},
		CpuLimit:    params.CpuLimit,
		MemoryLimit: params.MemoryLimit,
		StackLimit:  params.StackLimit,
		ProcLimit:   params.ProcLimit,
		CopyIn:      map[string]map[string]string{params.ExecFile: {"fileId": params.FileId}},
	}

	execRequest := &exec.Request{
		Cmd: []exec.Cmd{cmd},
	}

	return l.service.SendExecRequest(execRequest)
}

// evaluateRunResult judges the run result and returns whether it's accepted, along with output, time, and memory.
func (l *Logic) evaluateRunResult(execFile string, res *run.Response) (bool, string, string, int64, int64) {
	switch res.Status {
	case status.Accepted:
		return true, res.FileIds[execFile], res.Files["stdout"], res.Time, res.Memory
	default:
		return false, res.FileIds[execFile], res.Files["stderr"], res.Time, res.Memory
	}
}

// evaluateExecResult judges the exec result and returns status, output, time, and memory.
func (l *Logic) evaluateExecResult(res *exec.Response) (string, string, int64, int64) {
	switch res.Status {
	case status.Accepted:
		return res.Status, res.Files["stdout"], res.Time, res.Memory
	default:
		return res.Status, res.Files["stderr"], res.Time, res.Memory
	}
}
