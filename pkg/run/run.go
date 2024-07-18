package run

import (
	"encoding/json"
	"fmt"
)

type Cmd struct {
	Args          []string                     `json:"args"`
	Env           []string                     `json:"env"`
	Files         []map[string]interface{}     `json:"files"`
	CpuLimit      int64                        `json:"cpuLimit"`
	MemoryLimit   int64                        `json:"memoryLimit"`
	ProcLimit     int32                        `json:"procLimit"`
	StackLimit    int64                        `json:"stackLimit"`
	CopyIn        map[string]map[string]string `json:"copyIn"`
	CopyOut       []string                     `json:"copyOut"`
	CopyOutCached []string                     `json:"copyOutCached"`
}

type Request struct {
	Cmd []Cmd `json:"cmd"`
}

type Response struct {
	Status     string            `json:"status"`
	ExitStatus int32             `json:"exitStatus"`
	Time       int64             `json:"time"`
	Memory     int64             `json:"memory"`
	RunTime    int64             `json:"runTime"`
	Files      map[string]string `json:"files"`
	FileIds    map[string]string `json:"fileIds"`
}

func (c *Request) String() string {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("Error marshalling Request: %v", err)
	}
	return string(jsonData)
}
