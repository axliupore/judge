package ws

import (
	"github.com/axliupore/judge/config"
	"github.com/axliupore/judge/pkg/cache"
	"github.com/axliupore/judge/pkg/pool"
	"sync"
)

type Server struct {
	conf *config.Ws
	pool *pool.Pool
	c    *cache.Cache
}

var (
	once   sync.Once
	server *Server
)

func NewServer(conf *config.Ws, c *cache.Cache) *Server {
	once.Do(func() {
		p, _ := pool.New()
		server = &Server{conf: conf, pool: p, c: c}
	})
	return server
}
