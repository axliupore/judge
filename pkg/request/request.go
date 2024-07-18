package request

type JudgeRequest struct {
	Language    string `json:"language"`
	Content     string `json:"content"`
	Input       string `json:"input,omitempty"`
	CpuLimit    *int64 `json:"cpuLimit,omitempty"`
	MemoryLimit *int64 `json:"memoryLimit,omitempty"`
	StackLimit  *int64 `json:"stackLimit,omitempty"`
	ProcLimit   *int32 `json:"procLimit,omitempty"`
	Nsq         string `json:"nsq,omitempty"`
}

type JudgeResponse struct {
	Output string `json:"output"`
	Time   int64  `json:"time"`
	Memory int64  `json:"memory"`
}