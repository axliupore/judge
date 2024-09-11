package judge

type JavaScript struct {
}

func (js *JavaScript) ExecArgs() []string {
	return nil
}

func (js *JavaScript) RunArgs() []string {
	return []string{"/usr/bin/node", "main.js"}
}

func (js *JavaScript) IsBuild() bool {
	return false
}

func (js *JavaScript) Language() string {
	return "javascript"
}

func (js *JavaScript) Env() []string {
	return []string{"PATH=/usr/bin:/bin"}
}

func (js *JavaScript) RunFile() string {
	return "main.js"
}

func (js *JavaScript) ExecFile() []string {
	return nil
}
