package stats

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/orca-zhang/ecache"
)

func TestLRU2Cache(t *testing.T) {
	lc := ecache.NewLRUCache(1, 3, 10*time.Second).LRU2(1)
	Bind("lc", lc)
	lc.Put("1", "1")              // Added
	lc.Put("2", "2")              // Added
	lc.Put("3", "3")              // Added
	lc.Get("2")                   // l0 -> l1 GetHit
	lc.Get("3")                   // l0 -> l1 GetHit, Evicted
	if _, ok := lc.Get("2"); ok { // GetMiss
		t.Error("case 1 failed")
	}
	lc.Put("4", "4")               // Added
	lc.Put("5", "5")               // Added
	if _, ok := lc.Get("1"); !ok { // l0 -> l1 GetHit, Evicted
		t.Error("case 1 failed")
	}

	Stats().Range(func(k, v interface{}) bool {
		fmt.Printf("stats: %s %+v\n", k, v)
		if k == "lc" {
			node := v.(*StatsNode)
			if node.Evicted != 2 {
				t.Error("case 1 failed")
			}
			if node.Updated != 0 {
				t.Error("case 1 failed")
			}
			if node.Added != 5 {
				t.Error("case 1 failed")
			}
			if node.GetMiss != 1 {
				t.Error("case 1 failed")
			}
			if node.GetHit != 3 {
				t.Error("case 1 failed")
			}
			if node.DelMiss != 0 {
				t.Error("case 1 failed")
			}
			if node.DelHit != 0 {
				t.Error("case 1 failed")
			}
		}
		return true
	})

	lc.Put("6", "6")              // Added
	lc.Put("7", "7")              // Added, Evicted
	if _, ok := lc.Get("4"); ok { // GetMiss
		t.Error("case 1 failed")
	}
	lc.Del("7")                   // DelHit
	lc.Del("8")                   // DelMiss
	lc.Put("1", "1")              // Added
	lc.Put("1", "2")              // Updated
	lc.Del("1")                   // DelHit
	if _, ok := lc.Get("1"); ok { // GetMiss
		t.Error("case 1 failed")
	}

	Stats().Range(func(k, v interface{}) bool {
		fmt.Printf("stats: %s %+v\n", k, v)
		if k == "lc" {
			node := v.(*StatsNode)
			if node.Evicted != 3 {
				t.Error("case 1 failed")
			}
			if node.Updated != 1 {
				t.Error("case 1 failed")
			}
			if node.Added != 8 {
				t.Error("case 1 failed")
			}
			if node.GetMiss != 3 {
				t.Error("case 1 failed")
			}
			if node.GetHit != 3 {
				t.Error("case 1 failed")
			}
			if node.DelMiss != 1 {
				t.Error("case 1 failed")
			}
			if node.DelHit != 2 {
				t.Error("case 1 failed")
			}
		}
		return true
	})

	v, _ := Stats().Load("lc")
	node := v.(*StatsNode)
	if node.HitRate()-0.5 > 1e-6 {
		t.Error("case 1 failed")
	}
}

func TestConcurrent(t *testing.T) {
	lc := ecache.NewLRUCache(4, 1, 10*time.Second).LRU2(1)
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
	lcOld := ecache.NewLRUCache(1, 3, 1*time.Second).LRU2(1)
	Bind("lc2", lcOld)
	lc := ecache.NewLRUCache(1, 3, 1*time.Second).LRU2(1)
	Bind("lc2", lc)
	lc.Put("1", "1")
	Stats().Range(func(k, v interface{}) bool {
		fmt.Printf("stats: %s %+v\n", k, v)
		if k == "lc2" {
			node := v.(*StatsNode)
			if node.Added != 1 {
				t.Error("case 3 failed")
			}
		}
		return true
	})
}
