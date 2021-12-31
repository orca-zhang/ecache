package ecache

import (
	"sync"
	"sync/atomic"
	"time"
)

func now() int64 {
	return atomic.LoadInt64(&clock)
}

var clock = time.Now().UnixNano()

func init() {
	go func() {
		for {
			// calibration every second
			atomic.StoreInt64(&clock, time.Now().UnixNano())
			for i := 0; i < 9; i++ {
				time.Sleep(100 * time.Millisecond)
				atomic.AddInt64(&clock, int64(100*time.Millisecond))
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

// node to store cache item
type node struct {
	p, n *node
	k    string
	v    interface{}
	ts   int64 // nano timestamp
}

// a data structure that is efficient to insert/fetch/delete cache items [both O(1) time complexity]
type cache struct {
	cap        int
	hmap       map[string]*node
	head, tail *node // not use pointer-to-pointer here, coz it's trade-off for performance
}

// create a new lru cache object
func create(cap int) *cache {
	return &cache{cap, make(map[string]*node, cap), nil, nil}
}

// put a cache item into lru cache, if added return 1, updated return 0
func (c *cache) put(k string, v interface{}, on inspector) int {
	if e, ok := c.hmap[k]; ok {
		e.v, e.ts = v, now()
		c._refresh(e)
		return 0
	}

	if c.cap <= 0 {
		return 0
	} else if len(c.hmap) >= c.cap {
		// transfer the tail item as the new item, then refresh
		on(PUT, c.tail.k, &c.tail.v, -1)
		delete(c.hmap, c.tail.k)
		c.tail.k, c.tail.v, c.tail.ts = k, v, now() // reuse to reduce gc
		c.hmap[k] = c.tail
		c._refresh(c.tail)
		return 1
	}

	e := &node{nil, c.head, k, v, now()}
	if len(c.hmap) <= 0 {
		c.tail = e
	} else {
		c.head.p = e
	}
	c.hmap[k], c.head = e, e
	return 1
}

// get value of key from lru cache with result
func (c *cache) get(k string) (*node, int) {
	if e, ok := c.hmap[k]; ok {
		c._refresh(e)
		return e, 1
	}
	return nil, 0
}

// delete item by key from lru cache
func (c *cache) del(k string) (*node, int) {
	if e, ok := c.hmap[k]; ok {
		delete(c.hmap, k)
		c._remove(e)
		return e, 1
	}
	return nil, 0
}

// calls f sequentially for each valid item in the lru cache
func (c *cache) walk(walker func(k string, v interface{}, ts int64) bool) {
	for i := c.head; i != nil; i = i.n {
		if !walker(i.k, i.v, i.ts) {
			return
		}
	}
}

func (c *cache) _refresh(e *node) {
	if e.p == nil { // head node
		return
	}
	e.p.n = e.n
	if e.n == nil { // tail node
		c.tail = e.p
	} else {
		e.n.p = e.p
	}
	e.p, e.n, c.head.p, c.head = nil, c.head, e, e
}

func (c *cache) _remove(e *node) {
	if e.p == nil { // head node
		c.head = e.n
	} else {
		e.p.n = e.n
	}
	if e.n == nil { // tail node
		c.tail = e.p
	} else {
		e.n.p = e.p
	}
}

// hashCode hashes a string to a unique hashcode. BKDR hash as default
func hashCode(s string) (hash int) {
	for i := 0; i < len(s); i++ {
		hash = hash*131 + int(s[i])
	}
	return hash
}

// Cache - concurrent cache structure
type Cache struct {
	locks      []sync.Mutex
	insts      [][2]*cache // level-0 for normal LRU, level-1 for LRU-2
	mask       int
	expiration time.Duration
	on         inspector
}

func nextPowOf2(cap int) int {
	if cap <= 1 {
		return 1
	}
	if cap&(cap-1) == 0 {
		return cap
	}
	cap |= cap >> 1
	cap |= cap >> 2
	cap |= cap >> 4
	cap |= cap >> 8
	cap |= cap >> 16
	return cap + 1
}

// NewLRUCache - create lru cache
// `bucketCnt` is buckets that shard items to reduce lock racing
// `capPerBkt` is length of each bucket, can store `capPerBkt * bucketCnt` count of items in Cache at most
// optional `expiration` is item alive time (and we only use lazy eviction here), default `0` stands for permanent
func NewLRUCache(bucketCnt int, capPerBkt int, expiration ...time.Duration) *Cache {
	size := nextPowOf2(bucketCnt)
	c := &Cache{make([]sync.Mutex, size), make([][2]*cache, size), size - 1, 0, func(int, string, *interface{}, int) {}}
	for i := range c.insts {
		c.insts[i][0] = create(capPerBkt)
	}
	if len(expiration) > 0 {
		c.expiration = expiration[0]
	}
	return c
}

// LRU2 - add LRU-2 support (especially LRU-2 that when item visited twice it moves to upper-level-cache)
// `capPerBkt` is length of each LRU-2 bucket, can store extra `capPerBkt * bucketCnt` count of items in Cache at most
func (c *Cache) LRU2(capPerBkt int) *Cache {
	for i := range c.insts {
		c.insts[i][1] = create(capPerBkt)
	}
	return c
}

// Put - put a item into cache
func (c *Cache) Put(key string, val interface{}) {
	idx := hashCode(key) & c.mask
	c.locks[idx].Lock()
	status := c.insts[idx][0].put(key, val, c.on)
	c.locks[idx].Unlock()
	c.on(PUT, key, &val, status)
}

// internal sub function that get item at specific level
func (c *Cache) get(key string, idx, level int) (*node, int) {
	if n, s := c.insts[idx][level].get(key); s > 0 {
		if c.expiration > 0 && now()-n.ts > int64(c.expiration) {
			// not necessary to remove the expired item here, otherwise will cause GC thrashing
			return nil, 0
		}
		return n, s
	}
	return nil, 0
}

// Get - get value of key from cache with result
func (c *Cache) Get(key string) (interface{}, bool) {
	idx := hashCode(key) & c.mask
	c.locks[idx].Lock()
	var n *node
	var s int
	if c.insts[idx][1] == nil { // (if LRU-2 mode not support, loss is little)
		// normal lru mode
		n, s = c.get(key, idx, 0)
	} else {
		// LRU-2 mode
		n, s = c.insts[idx][0].del(key)
		if s <= 0 {
			// re-find in level-1
			n, s = c.get(key, idx, 1)
		} else {
			// find in level-0, move to level-1
			c.insts[idx][1].put(key, n.v, c.on)
		}
	}
	if s <= 0 {
		c.locks[idx].Unlock()
		c.on(GET, key, nil, 0)
		return nil, false
	}
	c.locks[idx].Unlock()
	c.on(GET, key, &n.v, 1)
	return n.v, true
}

// Del - delete item by key from cache
func (c *Cache) Del(key string) {
	idx := hashCode(key) & c.mask
	c.locks[idx].Lock()
	n, s := c.insts[idx][0].del(key)
	if c.insts[idx][1] != nil { // (if LRU-2 mode not support, loss is little)
		n2, s2 := c.insts[idx][1].del(key)
		if n2 != nil && (n == nil || n.ts < n2.ts) { // callback latest added one if both exists
			n, s = n2, s2
		}
	}
	if s > 0 {
		c.on(DEL, key, &n.v, 1)
	} else {
		c.on(DEL, key, nil, 0)
	}
	c.locks[idx].Unlock()
}

// Walk - calls f sequentially for each valid item in the lru cache, return false to stop iteration for every bucket
func (c *Cache) Walk(walker func(k string, v interface{}, ts int64) bool) {
	for i := range c.insts {
		c.locks[i].Lock()
		c.insts[i][0].walk(walker)
		if c.insts[i][1] != nil {
			c.insts[i][1].walk(walker)
		}
		c.locks[i].Unlock()
	}
}

const (
	PUT = iota + 1
	GET
	DEL
)

// inspector - can be used to statistics cache hit/miss rate or other scenario like buffer queue
//   `action`:PUT, `status`: evicted=-1, updated=0, added=1
//   `action`:GET, `status`: miss=0, hit=1
//   `action`:DEL, `status`: miss=0, hit=1
//   `value` only valid when `status` is not 0 or `action` is PUT
type inspector func(action int, key string, value *interface{}, status int)

// Inspect - to inspect the actions
func (c *Cache) Inspect(insptr inspector) {
	old := c.on
	c.on = func(action int, key string, value *interface{}, status int) {
		old(action, key, value, status) // call as the declared order, old first
		insptr(action, key, value, status)
	}
}
