package request

type JudgeRequest struct {
	Language    string `json:"language"`
	Content     string `json:"content,omitempty"`
	Input       string `json:"input,omitempty"`
	CpuLimit    *int64 `json:"cpuLimit,omitempty"`
	MemoryLimit *int64 `json:"memoryLimit,omitempty"`
	StackLimit  *int64 `json:"stackLimit,omitempty"`
	ProcLimit   *int32 `json:"procLimit,omitempty"`
	Nsq         string `json:"nsq,omitempty"`
	NotExec     bool   `json:"notExec,omitempty"`
	FileId      string `json:"fileId,omitempty"`
	DeleteFile  bool   `json:"deleteFile,omitempty"`
}
