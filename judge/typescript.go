package judge

type TypeScript struct {
}

func (ts *TypeScript) ExecArgs() []string {
	return nil
}

func (ts *TypeScript) RunArgs() []string {
	return []string{"/usr/bin/node", "main.ts"}
}

func (ts *TypeScript) IsBuild() bool {
	return false
}

func (ts *TypeScript) Language() string {
	return "typescript"
}

func (ts *TypeScript) Env() []string {
	return []string{"PATH=/usr/bin:/bin"}
}

func (ts *TypeScript) RunFile() string {
	return "main.ts"
}

func (ts *TypeScript) ExecFile() []string {
	return nil
}
