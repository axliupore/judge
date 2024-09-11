package cmd

type Cmd struct {
	Args          []string                     `json:"args,omitempty"`
	Env           []string                     `json:"env,omitempty"`
	Files         []map[string]interface{}     `json:"files,omitempty"`
	CpuLimit      int64                        `json:"cpuLimit,omitempty"`
	MemoryLimit   int64                        `json:"memoryLimit,omitempty"`
	ProcLimit     int32                        `json:"procLimit,omitempty"`
	StackLimit    int32                        `json:"stackLimit,omitempty"`
	CopyIn        map[string]map[string]string `json:"copyIn,omitempty"`
	CopyOut       []string                     `json:"copyOut,omitempty"`
	CopyOutCached []string                     `json:"copyOutCached,omitempty"`
}

type Request struct {
	Cmd []*Cmd `json:"cmd"`
}

type Response struct {
	Status     string            `json:"status,omitempty"`
	ExitStatus int32             `json:"exitStatus,omitempty"`
	Time       int64             `json:"time,omitempty"`
	Memory     int64             `json:"memory,omitempty"`
	RunTime    int64             `json:"runTime,omitempty"`
	Files      map[string]string `json:"files,omitempty"`
	FileIds    map[string]string `json:"fileIds,omitempty"`
}
