package nsq

import (
	"github.com/axliupore/judge/config"
	"github.com/axliupore/judge/internal/judge"
	"github.com/axliupore/judge/pkg/cache"
	"github.com/axliupore/judge/pkg/pool"
	"sync"
)

type Server struct {
	conf  *config.Nsq
	pool  *pool.Pool
	cache *cache.Cache
	judge *judge.Server
}

var (
	once   sync.Once
	server *Server
)

func NewServer(conf *config.Nsq, c *cache.Cache) *Server {
	once.Do(func() {
		p, _ := pool.New()
		server = &Server{
			conf:  conf,
			pool:  p,
			cache: c,
			judge: judge.NewServer(),
		}
	})
	return server
}
