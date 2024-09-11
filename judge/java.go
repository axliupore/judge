package judge

type Java struct {
}

func (java *Java) ExecArgs() []string {
	return []string{"/usr/bin/java", "Main"}
}

func (java *Java) RunArgs() []string {
	return []string{"/usr/bin/javac", "Main.java"}
}

func (java *Java) IsBuild() bool {
	return true
}

func (java *Java) Language() string {
	return "java"
}

func (java *Java) Env() []string {
	return []string{"PATH=/usr/bin:/bin"}
}

func (java *Java) RunFile() string {
	return "Main.java"
}

func (java *Java) ExecFile() []string {
	return []string{"Main.class"}
}
