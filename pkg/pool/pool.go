package pool

import (
	"github.com/panjf2000/ants/v2"
	"sync"
)

const (
	defaultSize = 2000
)

type Pool struct {
	pool *ants.Pool
	wg   sync.WaitGroup
}

var (
	p    *Pool
	once sync.Once
)

func New(size ...int) (*Pool, error) {
	var err error
	once.Do(func() {
		poolSize := defaultSize
		if len(size) > 0 {
			poolSize = size[0]
		}

		pool, err := ants.NewPool(poolSize)
		if err != nil {
			return
		}
		p = &Pool{pool: pool}
	})
	return p, err
}

func (p *Pool) Submit(task func()) error {
	p.wg.Add(1)
	if err := p.pool.Submit(func() {
		defer p.wg.Done()
		task()
	}); err != nil {
		return err
	}
	return nil
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Release() {
	p.pool.Release()
}
