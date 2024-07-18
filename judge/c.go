package judge

type C struct {
}

func (c *C) ExecArgs() []string {
	return []string{"main"}
}

func (c *C) RunArgs() []string {
	return []string{"/usr/bin/gcc", "main.c", "-o", "main"}
}

func (c *C) IsExec() bool {
	return true
}

func (c *C) Language() string {
	return "c"
}

func (c *C) Env() []string {
	return []string{"PATH=/usr/bin:/bin"}
}

func (c *C) RunFile() string {
	return "main.c"
}

func (c *C) ExecFile() string { return "main" }
