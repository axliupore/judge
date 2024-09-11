package judge

type Python struct {
}

func (python *Python) ExecArgs() []string {
	return nil
}

func (python *Python) RunArgs() []string {
	return []string{"/usr/bin/python3", "main.py"}
}

func (python *Python) IsBuild() bool {
	return false
}

func (python *Python) Language() string {
	return "python"
}

func (python *Python) Env() []string {
	return []string{"PATH=/usr/bin:/bin"}
}

func (python *Python) RunFile() string {
	return "main.py"
}

func (python *Python) ExecFile() []string { return nil }
