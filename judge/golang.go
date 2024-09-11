package judge

type Golang struct {
}

func (golang *Golang) ExecArgs() []string {
	return nil
}

func (golang *Golang) RunArgs() []string {
	return []string{"/usr/bin/go", "run", "main.go"}
}

func (golang *Golang) IsBuild() bool {
	return false
}

func (golang *Golang) Language() string {
	return "golang"
}

func (golang *Golang) Env() []string {
	return []string{"PATH=/usr/bin:/bin", "GOCACHE=/tmp/go-cache", "HOME=/root"}
}

func (golang *Golang) RunFile() string {
	return "main.go"
}

func (golang *Golang) ExecFile() []string { return nil }
