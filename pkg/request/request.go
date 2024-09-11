package request

type Run struct {
	Language string `json:"language,omitempty" validate:"required"`
	Code     string `json:"code,omitempty"`
	Input    string `json:"input,omitempty"`
	Nsq      string `json:"nsq,omitempty"`
	Limit
}

type Build struct {
	Language string `json:"language,omitempty" validate:"required"`
	Code     string `json:"code,omitempty"`
	Nsq      string `json:"nsq,omitempty"`
	Limit
}

type Exec struct {
	FileId   string `json:"fileId,omitempty" validate:"required"`
	Language string `json:"language,omitempty" validate:"required"`
	Input    string `json:"input,omitempty"`
	Nsq      string `json:"nsq,omitempty"`
	Limit
}

type Limit struct {
	CpuLimit    uint64 `json:"cpuLimit,omitempty"`
	MemoryLimit uint64 `json:"memoryLimit,omitempty"`
	StackLimit  uint32 `json:"stackLimit,omitempty"`
	ProcLimit   uint32 `json:"procLimit,omitempty"`
}

type Delete struct {
	FileId string `json:"fileId,omitempty" validate:"required"`
}
