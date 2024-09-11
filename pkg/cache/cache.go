package cache

import (
	"github.com/dgraph-io/ristretto"
	"sync"
	"time"
)

type Cache struct {
	r *ristretto.Cache
}

var (
	r    *Cache
	once sync.Once
)

func New() (*Cache, error) {
	var err error
	once.Do(func() {
		c, err := ristretto.NewCache(&ristretto.Config{
			NumCounters: 1e7,
			MaxCost:     1 << 30,
			BufferItems: 64,
		})
		if err != nil {
			return
		}

		r = &Cache{r: c}
	})
	return r, err
}

func (c *Cache) Get(key string) (interface{}, bool) {
	return c.r.Get(key)
}

func (c *Cache) Set(key string, value interface{}) bool {
	return c.r.Set(key, value, 0)
}

func (c *Cache) SetTime(key string, value interface{}, ttl time.Duration) bool {
	return c.r.SetWithTTL(key, value, 0, ttl)
}

func (c *Cache) Del(key string) {
	c.r.Del(key)
}

func (c *Cache) Wait() {
	c.r.Wait()
}
