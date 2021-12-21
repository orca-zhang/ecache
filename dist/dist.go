package dist

import (
	"log"
	"math/rand"
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
var lock = &sync.Mutex{}
var cacheMap = make(map[string][]*cache.Cache, 0)
var sender string

func Init(r RedisCli) {
	sender = func() string {
		res := []byte{}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i < 6; i++ {
			res = append(res, byte(r.Intn(26)+'a'))
		}
		return string(res)
	}()
	log.Println(sender)
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
				if len(vs) >= 3 {
					if sender != vs[2] {
						lock.Lock()
						for _, c := range cacheMap[vs[0]] {
							c.Del(vs[1], struct{}{})
						}
						lock.Unlock()
					}
				}
			})
		}
	}()
}

// Bind - to enable distributed consistency
// `pool` is not necessary, it can be used to classify instances that store same items
// but it will be more efficient if it is not empty
func Bind(pool string, caches ...*cache.Cache) error {
	lock.Lock()
	cacheMap[pool] = append(cacheMap[pool], caches...)
	lock.Unlock()
	for _, c := range caches {
		c.Inspect(func(action int, key string, ok int) {
			switch action {
			case cache.PUT:
				OnUpdate(pool, key)
			case cache.DEL:
				OnDel(pool, key)
			}
		})
	}
	return nil
}

// OnDel - delete `key` in `pool` at distributed scale
// `excludeLocal` is `true`, exclude local pool
func OnDel(pool string, key string, excludeLocal ...bool) error {
	exclude := ""
	// exclude local pool
	if len(excludeLocal) > 0 && excludeLocal[0] {
		exclude = sender
	}
	// pub to remote nodes
	r := redisCli
	if r != nil {
		_ = r.Pub(topic, strings.Join([]string{pool, key, exclude}, ":"))
	}
	return nil
}

// OnUpdate - delete `key` in `pool` at distributed scale exclude local
func OnUpdate(pool string, key string) error {
	return OnDel(pool, key, true)
}
