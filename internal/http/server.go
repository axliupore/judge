package http

import (
	"fmt"
	"github.com/axliupore/judge/config"
	"github.com/axliupore/judge/internal/judge"
	"github.com/axliupore/judge/pkg/cache"
	"github.com/axliupore/judge/pkg/pool"
	"sync"
)

type Server struct {
	conf  *config.Http
	pool  *pool.Pool
	cache *cache.Cache
	judge *judge.Server
}

var (
	once   sync.Once
	server *Server
)

func NewServer(conf *config.Http, c *cache.Cache) *Server {
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

func (s *Server) Start() {
	if err := Router().Run(fmt.Sprintf(":%d", s.conf.Port)); err != nil {
		panic(err)
	}
}
