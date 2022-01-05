package ecache

import (
	"container/list"
	"fmt"
	"sync"
	"testing"
	"time"
)

var on = func(int, string, *Value, int) {}

var inst = NewLRUCache(1, 1, time.Second)

type Elem struct {
	key string
	val string
}

func Test_create(t *testing.T) {
	c := create(5)
	if len(c.hmap) != 0 {
		t.Error("case 1 failed")
	}
}

func Test_put(t *testing.T) {
	c := create(5)
	c.put("1", inst.I("1"), nil, 0, on)
	c.put("2", inst.I("2"), nil, 0, on)
	c.put("1", inst.I("3"), nil, 0, on)
	if len(c.hmap) != 2 {
		t.Error("case 2.1 failed")
	}

	l := list.New()
	l.PushBack(&Elem{"1", "3"})
	l.PushBack(&Elem{"2", "2"})

	e := l.Front()
	for idx := c.dlnk[0][n]; idx != 0; idx = c.dlnk[idx][n] {
		v := e.Value.(*Elem)
		el := c.m[idx-1]
		if el.ts <= 0 {
			continue
		}
		if el.k != v.key {
			t.Error("case 2.2 failed: ", el.k, v.key)
		}
		if (*(el.v.I)).(string) != v.val {
			t.Error("case 2.3 failed: ", (*(el.v.I)).(string), v.val)
		}
		e = e.Next()
	}

	c.put("3", inst.I("4"), nil, 0, on)
	c.put("4", inst.I("5"), nil, 0, on)
	c.put("5", inst.I("6"), nil, 0, on)
	c.put("2", inst.I("7"), nil, 0, on)
	if len(c.hmap) != 5 {
		t.Error("case 3.1 failed")
	}

	l = list.New()
	l.PushBack(&Elem{"2", "7"})
	l.PushBack(&Elem{"5", "6"})
	l.PushBack(&Elem{"4", "5"})
	l.PushBack(&Elem{"3", "4"})
	l.PushBack(&Elem{"1", "3"})

	rl := list.New()
	rl.PushBack(&Elem{"1", "3"})
	rl.PushBack(&Elem{"3", "4"})
	rl.PushBack(&Elem{"4", "5"})
	rl.PushBack(&Elem{"5", "6"})
	rl.PushBack(&Elem{"2", "7"})

	e = l.Front()
	for idx := c.dlnk[0][n]; idx != 0; idx = c.dlnk[idx][n] {
		v := e.Value.(*Elem)
		el := c.m[idx-1]
		if el.ts <= 0 {
			continue
		}
		if el.k != v.key {
			t.Error("case 3.2 failed: ", el.k, v.key)
		}
		if (*(el.v.I)).(string) != v.val {
			t.Error("case 3.3 failed: ", (*(el.v.I)).(string), v.val)
		}
		e = e.Next()
	}

	e = rl.Front()
	for idx := c.dlnk[0][p]; idx != 0; idx = c.dlnk[idx][p] {
		v := e.Value.(*Elem)
		el := c.m[idx-1]
		if el.ts <= 0 {
			continue
		}
		if el.k != v.key {
			t.Error("case 3.4 failed: ", el.k, v.key)
		}
		if (*(el.v.I)).(string) != v.val {
			t.Error("case 3.5 failed: ", (*(el.v.I)).(string), v.val)
		}
		e = e.Next()
	}

	c.put("6", inst.I("8"), nil, 0, on)
	if len(c.hmap) != 5 {
		t.Error("case 4.1 failed")
	}

	l = list.New()
	l.PushBack(&Elem{"6", "8"})
	l.PushBack(&Elem{"2", "7"})
	l.PushBack(&Elem{"5", "6"})
	l.PushBack(&Elem{"4", "5"})
	l.PushBack(&Elem{"3", "4"})

	e = l.Front()
	for idx := c.dlnk[0][n]; idx != 0; idx = c.dlnk[idx][n] {
		v := e.Value.(*Elem)
		el := c.m[idx-1]
		if el.ts <= 0 {
			continue
		}
		if el.k != v.key {
			t.Error("case 4.2 failed: ", el.k, v.key)
		}
		if (*(el.v.I)).(string) != v.val {
			t.Error("case 4.3 failed: ", (*(el.v.I)).(string), v.val)
		}
		e = e.Next()
	}
}

func Test_get(t *testing.T) {
	c := create(2)
	c.put("1", inst.I("1"), nil, 0, on)
	c.put("2", inst.I("2"), nil, 0, on)
	if v, _ := c.get("1"); *(v.v.I) != "1" {
		t.Error("case 1.1 failed")
	}
	c.put("3", inst.I("3"), nil, 0, on)
	if len(c.hmap) != 2 {
		t.Error("case 1.2 failed")
	}

	l := list.New()
	l.PushBack(&Elem{"3", "3"})
	l.PushBack(&Elem{"1", "1"})

	e := l.Front()
	for idx := c.dlnk[0][n]; idx != 0; idx = c.dlnk[idx][n] {
		v := e.Value.(*Elem)
		el := c.m[idx-1]
		if el.k != v.key {
			t.Error("case 1.3 failed: ", el.k, v.key)
		}
		if (*(el.v.I)).(string) != v.val {
			t.Error("case 1.4 failed: ", (*(el.v.I)).(string), v.val)
		}
		e = e.Next()
	}
}

func Test_delete(t *testing.T) {
	c := create(5)
	c.put("3", inst.I("4"), nil, 0, on)
	c.put("4", inst.I("5"), nil, 0, on)
	c.put("5", inst.I("6"), nil, 0, on)
	c.put("2", inst.I("7"), nil, 0, on)
	c.put("6", inst.I("8"), nil, 0, on)
	c.del("5")

	l := list.New()
	l.PushBack(&Elem{"6", "8"})
	l.PushBack(&Elem{"2", "7"})
	l.PushBack(&Elem{"4", "5"})
	l.PushBack(&Elem{"3", "4"})
	/*if len(c.hmap) != 4 {
		t.Error("case 1.1 failed")
	}*/

	e := l.Front()
	for idx := c.dlnk[0][n]; idx != 0; idx = c.dlnk[idx][n] {
		el := c.m[idx-1]
		if el.ts <= 0 {
			continue
		}
		v := e.Value.(*Elem)
		if el.k != v.key {
			t.Error("case 1.2 failed: ", el.k, v.key)
		}
		if (*(el.v.I)).(string) != v.val {
			t.Error("case 1.3 failed: ", (*(el.v.I)).(string), v.val)
		}
		e = e.Next()
	}

	c.del("6")

	l = list.New()
	l.PushBack(&Elem{"2", "7"})
	l.PushBack(&Elem{"4", "5"})
	l.PushBack(&Elem{"3", "4"})
	/*if len(c.hmap) != 3 {
		t.Error("case 2.1 failed")
	}*/

	e = l.Front()
	for idx := c.dlnk[0][n]; idx != 0; idx = c.dlnk[idx][n] {
		el := c.m[idx-1]
		if el.ts <= 0 {
			continue
		}
		v := e.Value.(*Elem)
		if el.k != v.key {
			t.Error("case 2.2 failed: ", el.k, v.key)
		}
		if (*(el.v.I)).(string) != v.val {
			t.Error("case 2.3 failed: ", (*(el.v.I)).(string), v.val)
		}
		e = e.Next()
	}

	c.del("3")

	l = list.New()
	l.PushBack(&Elem{"2", "7"})
	l.PushBack(&Elem{"4", "5"})
	/*if len(c.hmap) != 2 {
		t.Error("case 3.1 failed")
	}*/

	e = l.Front()
	for idx := c.dlnk[0][n]; idx != 0; idx = c.dlnk[idx][n] {
		el := c.m[idx-1]
		if el.ts <= 0 {
			continue
		}
		v := e.Value.(*Elem)
		if el.k != v.key {
			t.Error("case 3.2 failed: ", el.k, v.key)
		}
		if (*(el.v.I)).(string) != v.val {
			t.Error("case 3.3 failed: ", (*(el.v.I)).(string), v.val)
		}
		e = e.Next()
	}
}

func Test_walk(t *testing.T) {
	c := create(5)
	c.put("3", inst.I(4), nil, 0, on)
	c.put("4", inst.I(5), nil, 0, on)
	c.put("5", inst.I(6), nil, 0, on)
	c.put("2", inst.I(7), nil, 0, on)
	c.put("6", inst.I(8), nil, 0, on)

	l := list.New()
	l.PushBack(&Elem{"6", "8"})
	l.PushBack(&Elem{"2", "7"})
	l.PushBack(&Elem{"5", "6"})
	l.PushBack(&Elem{"4", "5"})
	l.PushBack(&Elem{"3", "4"})

	e := l.Front()
	c.walk(
		func(key string, val *Value, ts int64) bool {
			v := e.Value.(*Elem)
			if key != v.key {
				t.Error("case 1.1 failed: ", key, v.key)
			}
			if fmt.Sprint(*(val.I)) != v.val {
				t.Error("case 1.2 failed: ", val.I, v.val)
			}
			e = e.Next()
			return true
		})

	if e != nil {
		t.Error("case 1.3 failed: ", e.Value)
	}

	e = l.Front()
	c.walk(
		func(key string, val *Value, ts int64) bool {
			v := e.Value.(*Elem)
			if key != v.key {
				t.Error("case 1.1 failed: ", key, v.key)
			}
			if fmt.Sprint(*(val.I)) != v.val {
				t.Error("case 1.2 failed: ", val.I, v.val)
			}
			return false
		})
}

func TestHashCode(t *testing.T) {
	if hashCode("12345") != int32(181817270) {
		t.Error("case 1 failed")
	}
	if hashCode("abcdefghijklmnopqrstuvwxyz") != int32(-91651236) {
		t.Error("case 2 failed")
	}
}

func TestNextPowOf2(t *testing.T) {
	if nextPowOf2(0) != 1 {
		t.Error("case 1 failed")
	}
	if nextPowOf2(1) != 1 {
		t.Error("case 2 failed")
	}
	if nextPowOf2(2) != 2 {
		t.Error("case 3 failed")
	}
	if nextPowOf2(3) != 4 {
		t.Error("case 4 failed")
	}
	if nextPowOf2(123) != 128 {
		t.Error("case 5 failed")
	}
	if nextPowOf2(0x7FFFFFFF) != 0x80000000 {
		t.Error("case 6 failed")
	}
}

func TestExpiration(t *testing.T) {
	lc := NewLRUCache(2, 1, time.Second)
	lc.Put("1", "2")
	if v, ok := lc.Get("1"); !ok || v != "2" {
		t.Error("case 1 failed")
	}
	time.Sleep(2 * time.Second)
	if _, ok := lc.Get("1"); ok {
		t.Error("case 2 failed")
	}

	// permanent
	lc2 := NewLRUCache(2, 1, 0)
	lc2.Put("1", "2")
	if v, ok := lc2.Get("1"); !ok || v != "2" {
		t.Error("case 1 failed")
	}
	time.Sleep(time.Second)
	if _, ok := lc2.Get("1"); !ok {
		t.Error("case 2 failed")
	}
}

func TestLRUCache(t *testing.T) {
	lc := NewLRUCache(1, 3, 1*time.Second)
	lc.Put("1", "1")
	lc.Put("2", "2")
	lc.Put("3", "3")
	v, _ := lc.Get("2") // check reuse
	lc.Put("4", "4")
	lc.Put("5", "5")
	lc.Put("6", "6")
	if v != "2" {
		t.Error("case 3 failed")
	}
}

func TestWalk(t *testing.T) {
	m := make(map[string]string, 0)
	lc := NewLRUCache(2, 3, 10*time.Second).LRU2(3)
	lc.Put("1", "1")
	m["1"] = "1"
	lc.Put("2", "2")
	m["2"] = "2"
	lc.Put("3", "3")
	m["3"] = "3"
	lc.Get("2") // l0 -> l1
	lc.Put("4", "4")
	m["4"] = "4"
	lc.Put("5", "5")
	m["5"] = "5"
	lc.Put("6", "6")
	m["6"] = "6"
	lc.Walk(func(key string, val *Value, ts int64) bool {
		if m[key] != (*val.I).(string) {
			t.Error("case failed")
		}
		delete(m, key)
		return true
	})
	if len(m) > 0 {
		fmt.Println(m)
		t.Error("case failed")
	}
}

func TestLRU2Cache(t *testing.T) {
	lc := NewLRUCache(1, 3, time.Second).LRU2(1)
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
	if _, ok := lc.Get("1"); !ok { // l0 -> l1
		t.Error("case 4 failed")
	}

	toCheck := "1"
	lc.Inspect(func(action int, key string, value *Value, ok int) {
		if action == DEL && value != nil && *(value.I) != toCheck {
			t.Error("case 4 failed")
		}
	})

	lc.Del("1")
	// del in l1

	if _, ok := lc.Get("1"); ok {
		t.Error("case 4 failed")
	}
	lc.Put("6", "6")
	lc.Put("7", "7")
	if _, ok := lc.Get("4"); ok {
		t.Error("case 4 failed")
	}

	// l0 -> l1 both exist
	lc.Put("1", "1")
	lc.Get("1") // l0 -> l1
	lc.Put("1", "2")

	// both del, return newest one
	toCheck = "2"
	lc.Del("1")

	if _, ok := lc.Get("1"); ok {
		t.Error("case 4 failed")
	}
}

func TestConcurrent(t *testing.T) {
	lc := NewLRUCache(4, 1, 2*time.Second)
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
}

func TestConcurrentLRU2(t *testing.T) {
	lc := NewLRUCache(4, 1, 2*time.Second).LRU2(1)
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
}

func TestInspect(t *testing.T) {
	lc := NewLRUCache(1, 3, 1*time.Second)
	lc.Inspect(func(action int, key string, value *Value, ok int) {
		if value != nil {
			fmt.Println(action, key, *(value.I), ok)
		} else {
			fmt.Println(action, key, ok)
		}
	})
	lc.Put("1", "1")
	lc.Put("1", "2")
	lc.Put("2", "2")
	lc.Put("3", "3")
	v, _ := lc.Get("2") // check reuse
	lc.Put("4", "4")
	lc.Put("5", "5")
	lc.Put("6", "6")
	if v != "2" {
		t.Error("case 3 failed")
	}
	lc.Get("10")
	lc.Del("6")
	lc.Del("10")
}
