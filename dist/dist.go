package dist

import (
	"log"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/orca-zhang/cache"
)

const topic = "orca-zhang/cache"

type RedisCli interface {
	// if the redis client is ready
	OK() bool
	// pub a key to channel
	Pub(channel, key string) error
	// sub a key from channel, callback uill tidy the local cache
	Sub(channel string, callback func(payload string)) error
}

var redisCli RedisCli
var m sync.Map

func delAll(pool, key string) {
	if caches, _ := m.Load(pool); caches != nil {
		for _, c := range *(caches.(*[]*cache.Cache)) {
			c.Del(key)
		}
	}
}

func Init(r RedisCli) {
	redisCli = r
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				debug.PrintStack()
			}
		}()

		for {
			for r == nil || !r.OK() {
				time.Sleep(10 * time.Millisecond)
			}
			_ = r.Sub(topic, func(payload string) {
				vs := strings.Split(payload, ":")
				if len(vs) >= 2 {
					delAll(vs[0], vs[1])
				}
			})
		}
	}()
}

// Bind - to enable distributed consistency
// `pool` is not necessary, it can be used to classify instances that store same items
// but it will be more efficient if it is not empty
func Bind(pool string, caches ...*cache.Cache) error {
	c, _ := m.LoadOrStore(pool, &[]*cache.Cache{})
	*(c.(*[]*cache.Cache)) = append(*(c.(*[]*cache.Cache)), caches...)
	return nil
}

// OnDel - delete `key` in `pool` at distributed scale
func OnDel(pool, key string) error {
	// pub to remote nodes
	r := redisCli
	if r != nil && r.Pub(topic, strings.Join([]string{pool, key}, ":")) == nil {
		return nil
	}
	delAll(pool, key)
	return nil
}
