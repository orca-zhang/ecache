package goredis

import (
	// "sync"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/orca-zhang/orcache"
	"github.com/orca-zhang/orcache/dist"
)

var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:         ":6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
}

func TestBind(t *testing.T) {
	dist.Init(Take(rdb, 10000))
	lc1 := orcache.NewLRUCache(1, 100, 10*time.Second)
	lc2 := orcache.NewLRUCache(1, 100, 10*time.Second)
	lc1.Put("1", "1")
	lc2.Put("1", "1")
	lc1.Put("2", "1")
	lc2.Put("2", "1")
	lc1.Put("3", "1")
	lc2.Put("3", "1")

	// bind them into a pool
	dist.Bind("lc", lc1, lc2)

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

func TestDisconnect(t *testing.T) {
	dist.Init(Take(rdb, 10000))
	rdb.Close()

	time.Sleep(5 * time.Second)
}

/*
func TestConcurrent(t *testing.T) {
	lc := orcache.NewLRUCache(4, 1, 2*time.Second).LRU2(1)
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
