package judge

const (
	golang     = "golang"
	cpp        = "cpp"
	c          = "c"
	java       = "java"
	python     = "python"
	javaScript = "javascript"
	typescript = "typescript"
)

// Judge interface defines methods for various programming language judges.
type Judge interface {
	ExecArgs() []string // Arguments for executing the file

	RunArgs() []string // Arguments for running the file

	IsBuild() bool // Whether the file is build

	Language() string // Programming language

	Env() []string // Environment variables

	RunFile() string // Name of the file to run

	ExecFile() []string // Name of the file to execute
}

// NewJudge creates a new Judge instance based on the given judgeType.
func NewJudge(judgeType string) Judge {
	switch judgeType {
	case golang:
		return &Golang{}
	case cpp:
		return &Cpp{}
	case c:
		return &C{}
	case java:
		return &Java{}
	case python:
		return &Python{}
	case javaScript:
		return &JavaScript{}
	case typescript:
		return &TypeScript{}
	default:
		return nil
	}
}
