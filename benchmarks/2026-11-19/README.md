# ecache v1.2.1: shard-lock layout benchmark

Compared with published v1.2.0, the only runtime change is the shard-lock layout. Independent shard mutexes are now separated by an architecture-sized, power-of-two stride. APIs, hashing, LRU/LRU-2, expiration and inspector behavior are unchanged. Go 1.22 remains supported.

## Change

```go
const cacheLineSize = max(64, unsafe.Sizeof(cpu.CacheLinePad{}))

type cacheMutex struct {
    sync.Mutex
    _ [cacheLineSize - unsafe.Sizeof(sync.Mutex{})]byte
}
```

The lock slice and its allocation use `[]cacheMutex`. `unsafe` is used only for compile-time sizes, with no pointer operations. `x/sys/cpu` was already a transitive dependency; its version is unchanged. The 64-byte minimum also permits targets such as WebAssembly where CacheLinePad has size zero.

Subtracting the mutex size avoids the experimental 136-byte stride on arm64: the final stride is 128 bytes there and 64 bytes on amd64. The distance keeps independent locks from sharing a cache line; it does not remove contention when callers use the same shard. This increases lock storage by 120 bytes per shard on arm64: 1,920 bytes for 16 shards, excluding allocator rounding.

## Results

On Apple M4 Pro / macOS 26.6.2 / Go 1.27.4, random reads across 256 preloaded string keys improved from **79.31 to 25.40 ns/op at 4 CPUs (3.12×)** and **117.20 to 46.41 ns/op at 8 CPUs (2.53×)**. All measured hit paths allocate **0 B/op and 0 allocs/op**.

Serial and single-CPU results were essentially unchanged in this run; the single-CPU independent-shard median changed from 13.42 to 13.43 ns/op. Improvements depend on actual shard contention.

Five paired runs alternate baseline/candidate order. The table gives median [minimum, maximum] ns/op. A missing CPU suffix means GOMAXPROCS=1. Ratios greater than one favor the candidate. Timing runs did not overlap tests or builds.

| Benchmark | v1.2.0 median [min, max] | v1.2.1 median [min, max] | Old/new |
| --- | ---: | ---: | ---: |
| Distributed | 19.620 [18.890, 22.130] | 19.140 [18.870, 19.510] | 1.03× |
| Distributed-4 | 79.310 [59.530, 85.570] | 25.400 [22.930, 32.630] | 3.12× |
| Distributed-8 | 117.200 [99.540, 122.700] | 46.410 [45.060, 48.280] | 2.53× |
| Serial | 12.740 [12.410, 14.800] | 12.680 [12.520, 12.840] | 1.00× |
| Serial-4 | 12.810 [12.430, 13.140] | 12.750 [12.490, 12.900] | 1.00× |
| Serial-8 | 12.940 [12.780, 13.190] | 12.720 [12.420, 13.010] | 1.02× |
| Shards | 13.420 [13.200, 13.550] | 13.430 [13.260, 14.000] | 1.00× |
| Shards-4 | 64.740 [62.480, 73.960] | 3.515 [3.453, 3.552] | 18.42× |
| Shards-8 | 86.050 [79.170, 112.100] | 2.116 [2.109, 2.281] | 40.67× |
| Shared | 13.380 [13.200, 13.510] | 13.230 [13.090, 13.790] | 1.01× |
| Shared-4 | 136.000 [120.700, 154.600] | 120.800 [113.200, 127.000] | 1.13× |
| Shared-8 | 113.500 [100.800, 120.600] | 111.100 [99.400, 113.500] | 1.02× |

`Serial` repeatedly reads one key in one goroutine. `Shared` uses RunParallel to read the same key. `Shards` gives workers keys in distinct adjacent shards, isolating false sharing from true lock contention. `Distributed` randomly reads 256 preloaded string keys using independent per-worker PRNG state. IntegerShards and PairShards, where present, exercise ecache2 integer and two-integer keys on distinct shards. All caches use 16 shards, capacity 64 per shard, and a one-minute sliding TTL; setup and string formatting occur before timing. Every returned value contributes to a sum stored to an atomic sink after the loop.

These are hit-path microbenchmarks, not application-throughput or memory-footprint measurements. Independent-shard results deliberately isolate the contention problem. Filesystem I/O, misses, eviction, writes and other architectures were not performance-benchmarked.

## Reproduce

Take the baseline tag and candidate source in separate directories, keeping their own go.mod/go.sum. Put the same [cache_bench_test.go](../../cache_bench_test.go) in each, then compile in the corresponding directories:

```sh
GOTOOLCHAIN=go1.27.4 go test -mod=readonly -c -o /tmp/cache-VARIANT.test
```

Run the two binaries alternately five times, reversing order every other round:

```sh
/tmp/cache-VARIANT.test -test.run='^$' -test.bench='^BenchmarkGet' -test.benchmem -test.benchtime=300ms -test.cpu=1,4,8 -test.count=1
```

Raw outputs: [baseline](baseline.txt), [candidate](candidate.txt). Parsed observations: [measurements.json](measurements.json). Exact inputs, binary SHA-256, environment, commands and order: [metadata.json](metadata.json).

## Validation

- All original local cache, stats and in-process distribution tests passed, including the million-iteration concurrency cases: `go test -mod=readonly -count=1 -timeout=4m . ./stats ./dist`.
- Core functional tests plus all benchmark patterns passed race checks at CPU settings 1/4/8 and 1,000 iterations per pattern. Race timings are excluded from performance results.
- Go 1.22.12 core race checks passed using external linking on macOS; the existing LC_DYSYMTAB linker warning remains.
- Go 1.22.12 library builds passed for Linux amd64, Linux arm64 and JavaScript/WebAssembly; these were compile checks, not runtime performance tests.
- All three Redis adapters passed their original tests against an isolated local Redis server: `go test -mod=readonly -p=1 -count=1 -timeout=2m ./dist/goredis/... ./dist/redigo`. [Adapter outputs](redis-tests.txt).
- `go vet ./...`, module verification and diff whitespace checks passed.

The full suite was not run with race instrumentation; distribution helpers retain their existing synchronization behavior.
