package response

type Response interface {
	SetOutput(string)
}

type Run struct {
	Status string `json:"status"`
	Output string `json:"output"`
	Time   int64  `json:"time"`
	Memory int64  `json:"memory"`
}

func (r *Run) SetOutput(output string) {
	r.Output = output
}

type Build struct {
	Status string `json:"status"`
	FileId string `json:"fileId"`
	Time   int64  `json:"time"`
	Memory int64  `json:"memory"`
}

func (r *Build) SetOutput(output string) {
	r.FileId = output
}

type Exec struct {
	Status string `json:"status"`
	Output string `json:"output"`
	Time   int64  `json:"time"`
	Memory int64  `json:"memory"`
}

func (r *Exec) SetOutput(output string) {
	r.Output = output
}
