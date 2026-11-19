package ecache

import (
	"fmt"
	"github.com/zeebo/xxh3"
	"sync/atomic"
	"testing"
	"time"
)

var benchmarkGetSink atomic.Uint64

// Each worker gets a key in its own adjacent shard, avoiding true lock contention.
func benchmarkGetKeys() []string {
	keys := make([]string, 16)
	for i, left := 0, len(keys); left > 0; i++ {
		key := fmt.Sprintf("/files/preview/%d", i)
		shard := xxh3.HashString(key) & 15
		if keys[shard] == "" {
			keys[shard] = key
			left--
		}
	}
	return keys
}

func benchmarkGetCache(keys []string) *Cache {
	c := NewLRUCache(16, 64, time.Minute)
	for _, key := range keys {
		c.PutBytes(key, []byte{1})
	}
	return c
}

func BenchmarkGetSerial(b *testing.B) {
	keys := benchmarkGetKeys()
	c := benchmarkGetCache(keys)
	key := keys[0]
	var sum uint64
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, ok := c.GetBytes(key)
		if !ok {
			b.Fatal("cache miss")
		}
		sum += uint64(data[0])
	}
	benchmarkGetSink.Store(sum)
}

func BenchmarkGetShared(b *testing.B)      { benchmarkGetParallel(b, "shared") }
func BenchmarkGetShards(b *testing.B)      { benchmarkGetParallel(b, "shards") }
func BenchmarkGetDistributed(b *testing.B) { benchmarkGetParallel(b, "distributed") }

func benchmarkGetParallel(b *testing.B, mode string) {
	keys := benchmarkGetKeys()
	if mode == "distributed" {
		keys = make([]string, 256)
		for i := range keys {
			keys[i] = fmt.Sprintf("/files/preview/%d", i)
		}
	}
	c := benchmarkGetCache(keys)
	var worker atomic.Uint64
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		state := worker.Add(1)
		key := keys[(state-1)%uint64(len(keys))]
		if mode == "shared" {
			key = keys[0]
		}
		var sum uint64
		for pb.Next() {
			if mode == "distributed" {
				state = state*6364136223846793005 + 1442695040888963407
				key = keys[byte(state>>56)]
			}
			data, ok := c.GetBytes(key)
			if !ok {
				b.Fatal("cache miss")
			}
			sum += uint64(data[0])
		}
		benchmarkGetSink.Store(sum)
	})
}
