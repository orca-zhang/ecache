package dist

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/orca-zhang/cache"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:         ":6379",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	Init(GoRedis(rdb))
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
	Bind("lc", lc1, lc2)

	time.Sleep(3 * time.Second)

	// try to del a item
	OnDel("lc", "1")

	time.Sleep(3 * time.Second)

	if _, ok := lc1.Get("1"); ok {
		t.Error("case 1 failed")
	}
	if _, ok := lc2.Get("1"); ok {
		t.Error("case 1 failed")
	}

	// try to del a item for other nodes
	OnDel("lc", "2", true)

	time.Sleep(3 * time.Second)

	if _, ok := lc1.Get("2"); !ok {
		t.Error("case 1 failed")
	}
	if _, ok := lc2.Get("2"); !ok {
		t.Error("case 1 failed")
	}

	lc1.Del("3")

	time.Sleep(3 * time.Second)

	if _, ok := lc1.Get("3"); ok {
		t.Error("case 1 failed")
	}
	if _, ok := lc2.Get("3"); ok {
		t.Error("case 1 failed")
	}
}
