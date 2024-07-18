package logic

import (
	"errors"
	"github.com/axliupore/judge/judge"
	"github.com/axliupore/judge/pkg/client"
	"github.com/axliupore/judge/pkg/exec"
	"github.com/axliupore/judge/pkg/log"
	"github.com/axliupore/judge/pkg/request"
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
	Args        []string // Arguments
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

// Start processes a JudgeRequest and returns a JudgeResponse or an error.
func (l *Logic) Start(r *request.JudgeRequest) (*request.JudgeResponse, string, error) {
	// Create a new judge instance based on the requested language
	j := judge.NewJudge(r.Language)
	if j == nil {
		log.Logger.Error("unsupported programming language")
		return nil, status.ParamsError, errors.New("unsupported programming language")
	}

	// Initialize params with default values
	params := NewParams()

	// Copy request data to params
	if err := copier.Copy(params, r); err != nil {
		log.Logger.Errorf("failed to copy params: %v", err)
		return nil, status.InternalError, nil
	}

	// Set params for running or executing based on judge type
	params.Args = j.RunArgs()
	params.Env = j.Env()
	params.ExecFile = j.ExecFile()
	params.RunFile = j.RunFile()
	params.Language = j.Language()

	// Run the code with the given params
	runRes, err := l.Run(params)
	if err != nil {
		log.Logger.Errorf("failed to Run %v", err)
		return nil, status.InternalError, err
	}

	// Judge the run result
	accept, fileId, output, time, memory := l.JudgeRun(params.ExecFile, runRes)
	defer l.DeleteFile(fileId)

	// Handle non-accepted results
	if !accept {
		if runRes.Status == status.TimeLimitExceeded || runRes.Status == status.MemoryLimitExceeded {
			log.Logger.Errorf("failed run accepted :%v", runRes.Status)
			return nil, runRes.Status, nil
		}
		log.Logger.Errorf("failed run accepted :%v", output)
		return nil, runRes.Status, errors.New(output)
	}

	// If it's not an executable task, return JudgeResponse with output, time, and memory
	if !j.IsExec() {
		return &request.JudgeResponse{
			Output: output,
			Time:   time,
			Memory: memory,
		}, runRes.Status, nil
	}

	// Set params for executing the code
	params.Args = j.ExecArgs()
	params.FileId = fileId

	// Execute the code with the given params
	execRes, err := l.Exec(params)
	if err != nil {
		log.Logger.Errorf("failed exec %v", err)
		return nil, runRes.Status, err
	}

	// Judge the exec result
	st, output, time, memory := l.JudgeExec(execRes)
	if st != status.Accepted {
		log.Logger.Errorf("failed exec %v %v", st, output)
		return nil, st, errors.New(output)
	}

	// Return JudgeResponse with output, time, and memory
	return &request.JudgeResponse{
		Output: output,
		Time:   time,
		Memory: memory,
	}, st, nil
}

// Run sends a run request with given params and returns the run response or an error.
func (l *Logic) Run(params *Params) (*run.Response, error) {
	cmd := run.Cmd{
		Args:        params.Args,
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

// Exec sends an exec request with given params and returns the exec response or an error.
func (l *Logic) Exec(params *Params) (*exec.Response, error) {
	cmd := exec.Cmd{
		Args:        params.Args,
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

// JudgeRun judges the run result and returns whether it's accepted, along with output, time, and memory.
func (l *Logic) JudgeRun(execFile string, res *run.Response) (bool, string, string, int64, int64) {
	switch res.Status {
	case status.Accepted:
		return true, res.FileIds[execFile], res.Files["stdout"], res.Time, res.Memory
	default:
		return false, res.FileIds[execFile], res.Files["stderr"], res.Time, res.Memory
	}
}

// JudgeExec judges the exec result and returns status, output, time, and memory.
func (l *Logic) JudgeExec(res *exec.Response) (string, string, int64, int64) {
	switch res.Status {
	case status.Accepted:
		return res.Status, res.Files["stdout"], res.Time, res.Memory
	default:
		return res.Status, res.Files["stderr"], res.Time, res.Memory
	}
}

// DeleteFile deletes the file identified by fileId.
func (l *Logic) DeleteFile(fileId string) {
	if fileId == "" {
		return
	}
	if err := l.service.DeleteRequest(fileId); err != nil {
		log.Logger.Errorf("deleteFile err: %v", err)
	}
}
