package judge

type Cpp struct {
}

func (cpp *Cpp) ExecArgs() []string {
	return []string{"main"}
}

func (cpp *Cpp) RunArgs() []string {
	return []string{"/usr/bin/g++", "main.cpp", "-o", "main"}
}

func (cpp *Cpp) IsBuild() bool {
	return true
}

func (cpp *Cpp) Language() string {
	return "cpp"
}

func (cpp *Cpp) Env() []string {
	return []string{"PATH=/usr/bin:/bin"}
}

func (cpp *Cpp) RunFile() string {
	return "main.cpp"
}

func (cpp *Cpp) ExecFile() []string { return []string{"main"} }
