package config

type Config struct {
	Http Http
	Ws   Ws
	Nsq  Nsq
}

type Http struct {
	Port int
}

type Ws struct {
	Port int
}

type Nsq struct {
	Address string
}
