package redigo

import (
	// "sync"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/orca-zhang/cache"
	"github.com/orca-zhang/cache/dist"
)

func init() {
	pool := &redis.Pool{
		// Other pool configuration not shown in this example.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
	dist.Init(Take(pool))
}

func TestBind(t *testing.T) {
	lc1 := cache.NewLRUCache(1, 100, 10*time.Second)
	lc2 := cache.NewLRUCache(1, 100, 10*time.Second)
	lc1.Put("1", "1")
	lc2.Put("1", "1")
	lc1.Put("2", "1")
	lc2.Put("2", "1")
	lc1.Put("3", "1")
	lc2.Put("3", "1")

	// bind them into a pool
	dist.Bind("lc", lc1)
	dist.Bind("lc", lc2)

	time.Sleep(3 * time.Second)

	// try to del a item
	dist.OnDel("lc", "1")

	time.Sleep(3 * time.Second)

	if _, ok := lc1.Get("1"); ok {
		t.Error("case 1 failed")
	}
	if _, ok := lc2.Get("1"); ok {
		t.Error("case 1 failed")
	}
}

/*
func TestConcurrent(t *testing.T) {
	lc := cache.NewLRUCache(4, 1, 2*time.Second).LRU2(1)
	dist.Bind("lc", lc)
	var wg sync.WaitGroup
	for index := 0; index < 10000; index++ {
		wg.Add(2)
		go func() {
			lc.Put("1", "2")
			wg.Done()
		}()
		go func() {
			lc.Get("1")
			wg.Done()
		}()
	}
	for index := 0; index < 100; index++ {
		wg.Add(1)
		go func() {
			time.Sleep(50 * time.Millisecond)
			dist.OnDel("lc", "1")
			wg.Done()
		}()
	}
	wg.Wait()
}*/
