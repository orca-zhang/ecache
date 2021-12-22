package dist

import (
	// "sync"
	"testing"
	"time"

	"github.com/orca-zhang/cache"
)

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
	Bind("lc", lc1)
	Bind("lc", lc2)

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
}
