package stats

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/orca-zhang/orcache"
)

func TestLRU2Cache(t *testing.T) {
	lc := cache.NewLRUCache(1, 3, 1*time.Second).LRU2(1)
	Bind("lc", lc)
	lc.Put("1", "1")
	lc.Put("2", "2")
	lc.Put("3", "3")
	lc.Get("2") // l0 -> l1
	lc.Get("3") // l0 -> l1
	if _, ok := lc.Get("2"); ok {
		t.Error("case 4 failed")
	}
	lc.Put("4", "4")
	lc.Put("5", "5")
	if _, ok := lc.Get("1"); !ok {
		t.Error("case 4 failed")
	}
	lc.Put("6", "6")
	lc.Put("7", "7")
	if _, ok := lc.Get("4"); ok {
		t.Error("case 4 failed")
	}
	lc.Del("7")
	lc.Del("8")
	lc.Put("1", "1")
	lc.Put("1", "2")
	lc.Del("1")
	if _, ok := lc.Get("1"); ok {
		t.Error("case 4 failed")
	}
	Stats().Range(func(k, v interface{}) bool {
		fmt.Printf("stats: %s %+v\n", k, v)
		return true
	})
}

func TestConcurrent(t *testing.T) {
	lc := cache.NewLRUCache(4, 1, 2*time.Second).LRU2(1)
	Bind("aaaa", lc)
	var wg sync.WaitGroup
	for index := 0; index < 1000000; index++ {
		wg.Add(3)
		go func() {
			lc.Put("1", "2")
			wg.Done()
		}()
		go func() {
			lc.Get("1")
			wg.Done()
		}()
		go func() {
			lc.Del("1")
			wg.Done()
		}()
	}
	wg.Wait()
	Stats().Range(func(k, v interface{}) bool {
		fmt.Printf("stats: %s %+v\n", k, v)
		return true
	})
}

func TestBindToExistPool(t *testing.T) {
	lcOld := cache.NewLRUCache(1, 3, 1*time.Second).LRU2(1)
	Bind("lc", lcOld)
	lc := cache.NewLRUCache(1, 3, 1*time.Second).LRU2(1)
	Bind("lc", lc)
	lc.Put("1", "1")
	lc.Put("2", "2")
	lc.Put("3", "3")
	lc.Get("2") // l0 -> l1
	lc.Get("3") // l0 -> l1
	if _, ok := lc.Get("2"); ok {
		t.Error("case 4 failed")
	}
	lc.Put("4", "4")
	lc.Put("5", "5")
	if _, ok := lc.Get("1"); !ok {
		t.Error("case 4 failed")
	}
	lc.Put("6", "6")
	lc.Put("7", "7")
	if _, ok := lc.Get("4"); ok {
		t.Error("case 4 failed")
	}
	lc.Del("7")
	lc.Del("8")
	lc.Put("1", "1")
	lc.Put("1", "2")
	lc.Del("1")
	if _, ok := lc.Get("1"); ok {
		t.Error("case 4 failed")
	}
	Stats().Range(func(k, v interface{}) bool {
		fmt.Printf("stats: %s %+v\n", k, v)
		return true
	})
}
