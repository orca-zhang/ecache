package dist

import (
	// "sync"
	"testing"
	"time"

	"github.com/orca-zhang/cache"
)

type DIYCli struct {
	ok *bool
	c  chan string
}

// if the redis client is ready
func (d *DIYCli) OK() bool {
	return d.ok != nil && *d.ok
}

// pub a payload to channel
func (d *DIYCli) Pub(channel, payload string) error {
	d.c <- payload
	return nil
}

// sub a payload from channel, callback uill tidy the local cache
func (d *DIYCli) Sub(channel string, callback func(payload string)) error {
	for {
		if payload, ok := <-d.c; ok {
			callback(payload)
		} else {
			break
		}
	}
	return nil
}

func Take(ok *bool) RedisCli {
	return &DIYCli{
		ok: ok,
		c:  make(chan string, 100),
	}
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

func TestInit(t *testing.T) {
	// nil Init
	Init(nil)

	// is OK
	OnDel("lc", "1")
}

func TestDIYClient(t *testing.T) {
	ok := false

	Init(Take(&ok))

	time.Sleep(3 * time.Second)

	// mark ready
	ok = true

	// is OK
	OnDel("lc", "1")

	time.Sleep(3 * time.Second)

	lc1 := cache.NewLRUCache(1, 100, 10*time.Second)
	lc1.Put("1", "1")

	if _, ok := lc1.Get("1"); !ok {
		t.Error("case 2 failed")
	}

	// bind them into a pool
	Bind("lc", lc1)
	OnDel("lc", "1")

	time.Sleep(3 * time.Second)

	if _, ok := lc1.Get("1"); ok {
		t.Error("case 2 failed")
	}
}
