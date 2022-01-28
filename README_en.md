[Simplified Chinese README | 简体中文说明](README.md)

# 🦄 ecache
<p align="center">
  <a href="#">
    <img src="https://github.com/orca-zhang/ecache/raw/master/doc/logo.svg">
  </a>
</p>

<p align="center">
  <a href="/go.mod#L3" alt="go version">
    <img src="https://img.shields.io/badge/go%20version-%3E=1.11-brightgreen?style=flat"/>
  </a>
  <a href="https://goreportcard.com/badge/github.com/orca-zhang/ecache" alt="goreport">
    <img src="https://goreportcard.com/badge/github.com/orca-zhang/ecache">
  </a>
  <a href="https://orca-zhang.semaphoreci.com/projects/ecache" alt="buiding status">
    <img src="https://orca-zhang.semaphoreci.com/badges/ecache.svg?style=shields">
  </a>
  <a href="https://codecov.io/gh/orca-zhang/ecache" alt="codecov">
    <img src="https://codecov.io/gh/orca-zhang/ecache/branch/master/graph/badge.svg?token=F6LQbADKkq"/>
  </a>
  <a href="https://github.com/orca-zhang/ecache/blob/master/LICENSE" alt="license MIT">
    <img src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat">
  </a>
  <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Forca-zhang%2Fcache?ref=badge_shield" alt="FOSSA Status">
    <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Forca-zhang%2Fcache.svg?type=shield"/>
  </a>
  <a href="https://benchplus.github.io/gocache/dev/bench/" alt="continuous benchmark">
    <img src="https://img.shields.io/badge/benchmark-click--me-brightgreen.svg?style=flat"/>
  </a>
</p>
<p align="center">Extremely easy, ultra fast, concurrency-safe and support distributed consistency.</p>

## Features

- 🤏  Less than 300 lines, cost only ~30s to assemble
- 🚀  Extremely easy, ultra fast and  concurrency-safe
- 🏳️‍🌈  Support both `LRU` mode and  [`LRU-2`](#LRU-2-mode) mode inside
- 🦖  [Extra plugin](#Distributed-Consistency-Plugin) that support distributed consistency

## Benchmarks

> [👁️‍🗨️click me to see cases](https://github.com/benchplus/gocache) [👁️‍🗨️click me to see results](https://benchplus.github.io/gocache/dev/bench/) (the lower the better except cache hit rate)
![](https://github.com/orca-zhang/ecache/raw/master/doc/benchmark.png)

> gc pause test result [code provided by `bigcache`](https://github.com/allegro/bigcache-bench) (the lower the better)
![](https://github.com/orca-zhang/ecache/raw/master/doc/gc.png)

### Stablity validation in production environment
- [`Confirmed`]Official Account Backend(hundreds QPS), user & order info, configrations.
- [`WIP`]Push Platform(tens of thousands QPS), system configrations, deduplication, fixed info cache like app info and etc.
- [`TBD`]Comment Platform(tens of thousands QPS), user info and distributed consistency plugin for user avatar & nickname.

## How To Use

#### Import Package (almost 5s)
``` go
import (
    "time"

    "github.com/orca-zhang/ecache"
)
```

#### Definition (almost 5s)
> Can be placed in any position (global is also OK), it is recommended to define nearby
``` go
var c = ecache.NewLRUCache(16, 200, 10 * time.Second)
```

#### Put Item (almost 5s)
``` go
c.Put("uid1", o) // `o` can be any variable, generally an object pointer, storing fixed information, such as `*UserInfo`
```

#### Retrive Item (almost 5s)
``` go
if v, ok := c.Get("uid1"); ok {
    return v.(*UserInfo) // No type assertion, let's control the type by ourselves
}
// If it is not found in memory cache , go back to query the redis/db
```

#### Remove Item (almost 5s)
> when the original info was updated
``` go
c.Del("uid1")
```

#### Download Package (almost 5s)

> non-go modules mode:\
> sh>  ```go get -u github.com/orca-zhang/ecache```

> go modules mode:\
> sh>  ```go mod tidy && go mod download```

#### Fire
> 🎉 Finished. 🚀 Performance accelerated to X times! \
> sh>  ```go run <your-main.go file>```

## Instruction

- `NewLRUCache`
  - First parameter is the number of buckets, each bucket will use an independent lock, max to 65535(for 65536 buckets)
    - Don't worry, just set as you want, `ecache` will find a suitable number which is convenient for mask calculation later
  - Second parameter is the number of items that each bucket can hold, max to 65535
    - When `ecache` is full, there should be `first parameter X second parameter` item, can store max to 4.2 billion items
  - \[`Optional`\]Third parameter is the expiration time of each item
    - `ecache` uses internal counter to improve performance, default 100ms accuracy, calibration every second
    - No parameter or pass `0`, means permanent

## Best Practices

- Support any type of value
  - Provides `Put`/`PutInt64`/`PutBytes` three methods to adapt to different scenarios and need to be used in pairs with `Get`/`GetInt64`/`GetBytes` (the latter two methods have less GC cost)
  - Store pointers for complex objects (Note: ⚠️ Do not modify its fields once it is put in, even if it is taken out again, because the item may be accessed by other people at the same time)
    - If you need to modify, the solution: take out each individual assignment of the field, or use [copier to make a deep copy and modify on the copy](#need-to-modify-and-store-the-object-pointer)
    - Objects can also be stored directly (compared to the previous one, the performance is worse because there are copy operations when taken out)
    - The larger cached objects, the better, the upper level of the business, the better (save memory assembly and data organization time)
- If you don’t want to erase the hot data due to traversal requests, you can switch to [`LRU-2` mode](#LRU-2-mode), there may be very little loss (💬 [What Is LRU-2](#What-Is-LRU-2))
  - - The size of `LRU2` and `LRU` is set to 1/4 and 3/4, which may perform better。
- One instance can store multiple types of objects, try adding a prefix when formatting the key and separating it with a colon
- For scenes with large concurrent visits, try `256`, `1024` buckets, or even more
- Can be used as a **buffer queue** to merge updates to reduce disk flushes (data can be rebuilt or tolerate loss of power outage)
   - [Add an `Inspector`](#inject-an-inspector) to monitor the eviction event
   - At the end or intervally call [`Walk`](#fetch-all-items) to flush the data to storage

## Special Scenarios

### integer key, integer value and bytes value
``` go
// integer key
c.Put(strconv.FormatInt(d, 10), o) // d is type of `int64`

// integer value
c.PutInt64("uid1", int64(1))
if d, ok := c.GetInt64("uid1"); ok {
    // d is type of `int64` and value is 1
}

// bytes value
c.PutBytes("uid1", b)// b is type of `[]byte`
if b, ok := c.GetBytes("uid1"); ok {
    // b is type of `[]byte`
}
```

### LRU-2 mode

- 💬 [What Is LRU-2](#What-Is-LRU-2)

> Just follow `NewLRUCache()` directly with `.LRU2(<num>)`, and the parameter `<num>` represents the number of items in the `LRU-2` hot queue (per bucket)
``` go
var c = ecache.NewLRUCache(16, 200, 10 * time.Second).LRU2(1024)
```

### Empty cache sentry (non-existent objects do not need to query the source)
``` go
// Just give `nil` when put
c.Put("uid1", nil)
```

``` go
// When reading, it is almost like normal
if v, ok := c.Get("uid1"); ok {
  if v == nil { // Note:⚠️ it is necessary to judge whether it is empty
    return nil  // If it is empty, then return `nil` or you can prevent `uid1` from appearing in the list of source to be queried
  }
  return v.(*UserInfo)
}
// If the memory cache is miss, go back to query the redis/db
```

### Need to modify, and store the object pointer

> For example, we get the user information cache `v` of type `*UserInfo` from `ecache`, and need to modify its status field
``` go
import (
    "github.com/jinzhu/copier"
)
```

``` go
o := &UserInfo{}
copier.Copy(o, v) // Copy from `v` to `o`
o.Status = 1      // Modify the field of the copy
```

### Inject an inspector

``` go
// inspector - can be used to do statistics or buffer queues, etc.
// `action`:PUT, `status`: evicted=-1, updated=0, added=1
// `action`:GET, `status`: miss=0, hit=1
// `action`:DEL, `status`: miss=0, hit=1
// `iface`/`bytes` is not `nil` when `status` is not 0 or `action` is PUT
type inspector func(action int, key string, iface *interface{}, bytes []byte, status int)
```

- How to use
``` go
cache.Inspect(func(action int, key string, iface *interface{}, bytes []byte, status int) {
   // TODO: add what you want to do
   //     Inspector will be executed in sequence according to the injection order
   //     Note:⚠️ If there is a operation that takes a long time, try to transfer job to another channel to ensure not blocking current coroutine.
   
   // - how to fetch right value -
  //   - `Put`:      `*iface`
  //   - `PutBytes`: `bytes`
  //   - `PutInt64`: `ecache.ToInt64(bytes)`
})
```

### Fetch all items

``` go
  // only invalid items can be fetched
	cache.Walk(func(key string, iface *interface{}, bytes []byte, expireAt int64) bool {
    // `key` is key of item, `iface`/`bytes` is value of item, `expireAt` is the time that item expired

    // - how to fetch right value -
    //   - `Put`:      `*iface`
    //   - `PutBytes`: `bytes`
    //   - `PutInt64`: `ecache.ToInt64(bytes)`
	  return true // true stands for walk on
	})
```

## Cache Usage Statistics

> The implementation is super simple. After the inspector is injected, only one more atomic operation is added to each operation. See [details](/stats/stats.go#L26).

##### Import the `stats` package
``` go
import (
    "github.com/orca-zhang/ecache/stats"
)
```

#### Bind the cache instance
> The name is a custom pool name, which will be aggregated by name internally.\
> Note:⚠️ The binding can be placed in global scope.
``` go
var _ = stats.Bind("user", c)
var _ = stats.Bind("user", c0, c1, c2)
var _ = stats.Bind("token", caches...)
```

#### Get statistics
``` go
stats.Stats().Range(func(k, v interface{}) bool {
    fmt.Printf("stats: %s %+v\n", k, v) // k is name of pool, v is type of (*stats.StatsNode) that stats count of events use `HitRate` method can know cache hit rate
    return true
})
```

## Distributed Consistency Plugin

- 💬 [Principle Explanation](#Principle-of-Distributed-Consistency-Plugin)

### Import the `dist` package
``` go
import (
    "github.com/orca-zhang/ecache/dist"
)
```

### Bind cache instance
> The name is a custom pool name, which will be aggregated by name internally.\
> Note:⚠️ The binding can be placed in global scope and does not depend on initialization.
``` go
var _ = dist.Bind("user", c)
var _ = dist.Bind("user", c0, c1, c2)
var _ = dist.Bind("token", caches...)
```

### Bind redis client
> Currently `redigo` and `goredis` are supported, other libraries can implement the `dist.RedisCli` interface by yourselves, or you can submit an issue to me.

#### go-redis v7 and below
``` go
import (
    "github.com/orca-zhang/ecache/dist/goredis/v7"
)

dist.Init(goredis.Take(redisCli)) // redisCli is *redis.RedisClient type
dist.Init(goredis.Take(redisCli, 100000)) // Second parameter is size of channel buffer, default is 100 if not passed
```

#### go-redis v8 and above
``` go
import (
    "github.com/orca-zhang/ecache/dist/goredis"
)

dist.Init(goredis.Take(redisCli)) // redisCli is *redis.RedisClient type
dist.Init(goredis.Take(redisCli, 100000)) // Second parameter is size of channel buffer, default is 100 if not passed
```

#### redigo
> Note:⚠️ `github.com/gomodule/redigo` requires minimum version `go 1.14`
``` go
import (
    "github.com/orca-zhang/ecache/dist/redigo"
)

dist.Init(redigo.Take(pool)) // pool is of *redis.Pool type
```

#### Proactively notify all nodes and all instances to delete item (including local machine)
> Called when the data of db changes or is deleted\
> When error occurs, it will be downgraded to local operation (such as uninitialized or network error)
``` go
dist.OnDel("user", "uid1") // user is name of pool, uid1 is the key that want to be deleted
```

## Update guide for old [`lrucache`](http://github.com/orca-zhang/lrucache) fans

- Only four steps:
1. Import `github.com/orca-zhang/ecache` instead of `github.com/orca-zhang/lrucache`
2. `ecache.NewLRUCache` instead of `lrucache.NewSyncCache`
3. Third parameter should add unit `*time.Second`
4. `Delete` method replace to `Del`

# You won't leave empty-handed

- Guest officer, let's learn something before leaving!
- I want to try my best to make you understand what `ecache` did and why.

## What is local memory cache

---
    L1 cache reference ......................... 0.5 ns
    Branch mispredict ............................ 5 ns
    L2 cache reference ........................... 7 ns
    Mutex lock/unlock ........................... 25 ns
    Main memory reference ...................... 100 ns
    Compress 1K bytes with Zippy ............. 3,000 ns  =   3 µs
    Send 2K bytes over 1 Gbps network ....... 20,000 ns  =  20 µs
    Read 1 MB sequentially from memory ..... 250,000 ns  = 250 µs
    Round trip within same datacenter ...... 500,000 ns  = 0.5 ms
    Send packet CA<->Netherlands ....... 150,000,000 ns  = 150 ms

- As can be seen from the above table, the gap between memory access and network access (same as data center) is almost one to ten thousand times!
- I have encountered more than one engineer: "Cache? Use redis", but I want to say that redis is not a panacea. To some extent, using it is still a nightmare (of course I am talking about cache consistency issues...😄)
- Because the memory operation is very fast, it can be basically omitted when compared to redis/db. For example, there is a 1000QPS query API. If we cache the result for 1 second, which means that redis/db won't be requested within 1 second, then source query counts is reduced to 1/1000 (ideally), so the performance of accessing the redis/db part has been improved by 1000 times. Doesn't it sound great?
- Keep watching, you will fall in love with her! (Of course might be him or it, ahaha)

### Use Scenarios(problems to be solved)

- High concurrency and large traffic scenarios
   - Cache hotspot data (such as live broadcast rooms with high popularity)
   - Sudden QPS peak clipping (such as breaking news in the information stream)
   - Reduce latency and congestion (such as frequently visited pages in a short period of time)
- Cut costs
   - Stand-alone scenario (the QPS can be quickly increased without deploying redis or memcache)
   - Downgrade of redis and db instances (can intercept most requests)
- Persistent or semi-persistent data (write less and read more)
   - For example, configuration, etc. (This kind of data is used in many places, and there will be an amplification effect. In many cases, these configuration hot keys may lead to specification mis-upgrade of the redis/db instance)
- Inconsistent-tolerated data
   - Such as user avatar, nickname, product inventory (the actual order will be checked again in db), etc.
   - Modified configuration (expiration time is 10 seconds, then it will take effect with a maximum delay of 10 seconds)
- Buffer queue: merge updates to reduce disk flushes
   - Can achieve strong consistency by patching query with cache diff (in the case of distributed, it is necessary to ensure that the same user/device is balanced to the same node at the load balancing layer)
   - Data can be rebuilt or tolerate loss of power outage

## Design Ideas

> `ecache` is an upgraded version of the [`lrucache`](http://github.com/orca-zhang/lrucache) library

- Bottom layer is the most basic `LRU` implemented with native map and double-linked lists (the longest not visited)
   - PS: All other versions I implemented ([go](https://github.com/orca-zhang/lrucache) / [C++](https://github.com/ez8-co/linked_hash) / [js](https://github.com/orca-zhang/ecache.js)) in leetcode are solutions beats 100% submissions.
- Second layer includes bucketing strategy, concurrency control, and expiration control (it will automatically adapt to power-of-two buckets to facilitate mask calculation)
- The 2.5 layer implements the `LRU-2` ability in a very simple way, the code does not exceed 20 lines, directly look at the source code (search for the keyword `LRU-2`)

### What is LRU

- Evict longest visited item first.
- Each time it is accessed, the item will be refreshed to the first of the queue.
- Put a new item when the queue is full, the last item in the queue, that is the item has not been accessed for the longest time, will be evicted.

### What Is LRU-2

- `LRU-K` is that less than K visits items stored in a separate `LRU` queue, and additional queue for more than K visits
- The target scenario is that, for example, some traversal queries will evict some hot items that we need in the future.
- For the sake of simplicity, what we have implemented here is `LRU-2`, that is, the second visit is placed in the hot queue, and the count of visits is not recorded.
- It used to optimize the cache hit rate of hot keys.
- Very similar to [mysql's buffer pool lru algorithm](https://dev.mysql.com/doc/refman/5.7/en/innodb-buffer-pool.html).

### Principle of Distributed Consistency Plugin

- In fact, it simply uses the pubsub feature of redis
- Proactively inform that the cached information is updated and broadcast it to all the nodes
- In a sense, it is just a way to narrow the inconsistent time window (there is a network delay and it is not guaranteed to be completed)
- Pay Attention: ⚠️
   - Reduce the use even if necessary, suitable for the scenario where write less read more `WORM(Write-Once-Read-Many)`
     - Because redis's performance is not as good as memory after all, and there is broadcast communication (write amplification)
   - The following scenarios will be degraded (the time window becomes larger), but at least the strong consistency of the current node will be guaranteed
     - Redis is unavailable, network error
     - Consume goroutine panic
     - When not all nodes are ready (`canary` deployment, or in the process of deployment), such as
       - Already used `ecache` but added this plugin for the first time
       - Newly added cached data or newly added delete operation

### About performance

- No defer is needed to release the lock.
- No need to clean up asynchronously (clean-up is meaningless, it is more reasonable to disperse to eviction when writing, and it is not easy to GC thrashing).
- No memory capacity is used to control (the size of a single item generally has an estimated size, simply control the number).
- Bucket strategy, automatic selection of power-of-two buckets (reduce lock competition, power-of-two mask operation is faster).
- Use string type for key (strong scalability; built-in language support for reference, which saves memory).
- No virtual header for doubly-linked list (although it is a little bit around, but there is an increase of about 20%).
- Choose `LRU-2` to implement `LRU-K` (simple implementation, almost no additional loss).
- Store pointers directly (without serialization, the advantage is greatly reduced if you use `[]byte` in some scenarios).
- Use internal counter for timing (default 100ms accuracy, calibration per second, `pprof` found that time.Now() generates temporary objects, which leads to increased GC time consumption).
- -Double-linked list uses fixed allocation memory storage, uses zero timestamp to mark delete, reduces GC (and saves memory by more than 50% compared with `bigcache` in the same specification)

#### Failed optimization attempt

- The key is changed from string to `reflect.StringHeader`, result: negative optimization.
- The mutex lock is changed to a read-write lock, the Get request will also modify the data, and the access is illegal, even if the data is not changed, the result: negative optimization for read-write mixed scenarios.
- Use `time.Timer` implements the internal counter, the result: the trigger is unstable, use `time.Sleep` instead.
- Distributed consistency plugin that automatically updated and deleted by the inspector. The result: the performance decreased and the loop call problem needs to be specially dealt with.

### About GC optimization

- As I mentioned in the C++ version of the performance profiler [several levels of performance optimization](https://github.com/ez8-co/ezpp#性能优化的几个层次), consider at only one level is not good.
- The Third Level says, 'Nothing is faster than nothing' (similar to Occam's razor), you should not come up with optimization if you can remove it.
- For example, some library want to reduce GC by allocating large block of memory, but provides `[]byte` value storage, which means that it may need extra serialization and copy.
- If the serialized part can be reused in the protocol layer that `ZeroCopy` can be achieved is OK, but things go contrary to one's wishes, and the `ecache` storage pointer directly so that omit the extra cost.
- What I want to express is that GC optimization is really important, but more that it should be combined with the scene, and extra loss of client-end also needs to be considered, instead of claiming gc-free, the result is not that way.
- The violent aesthetics I advocate is minimalism, the defect rate is proportional to the amount of code, complex things will be eliminated sooner or later, and `KISS` is the true king.
- `ecache` has only less than 300 lines in total, and if the bug rate of thousand lines is fixed, there aren't many bugs in it.

## FAQ
> Q: Can an instance store multiple kind of objects?
- A: Yes, for example, you can format the key with a prefix (like redis keys separated by a colon), please pay attention to ⚠️not misusing the wrong type.

> Q: How to set different expiration times for different items?
- A: Use several cache instances. (😄did not expect?)

> Q: How to solve the problem of very-very-very hot key?
- A: The [local memory cache] is used for cache hot keys, so very-very-very hot keys here can be understood as single node hundreds of thousands of QPS, the biggest problem is that there are too many lock competitions on a single bucket, which affects other data in the same bucket. Then it can be like this: First, use `LRU-2` to prevent similar traversal requests from flushing hot data. Secondly, in addition to adding buckets, you can write multiple instances (write the same item at the same time) and read a certain one (for example, according to the access user uid hash), let the hot key have multiple copies. But when deleting (reverse writing), be careful to delete all instances of multiple instances, which is suitable for the scenario of "write less read more `WORM (Write-Once-Read-Many)`". " The scenario of “write more, read more” can extract the diff separately and turn it into a `WORM` scenario.

> Q: Why not deal with doubly-linked list in the way of virtual headers? It's bullshxt now!
- A: The leaked code [[lrucache](http://github.com/orca-zhang/lrucache)] has been challenged on the V2EX on 2019-04-22. It’s really not that I don't know to use virtual headers. Although it is more confusing to read than the pointer-to-pointer method, current way has an improvement of about 20%! (😄did not expect?)

## Thanks

Gratitude to them who performed code review, errata, and valuable suggestions during the development process! (names not listed in order)

<table>
  <tr>
    <td align="center">
      <a href="https://github.com/askuy">
        <img src="https://avatars.githubusercontent.com/u/14119383?v=4" width="64px;" alt=""/>
        <br />
        <b>askuy</b>
        <br />
        <sub><a href="https://github.com/gotomicro/ego">[ego]</a></sub>
      </a>
    </td>
    <td align="center">
      <a href="https://github.com/auula">
        <img src="https://avatars.githubusercontent.com/u/38412458?v=4" width="64px;" alt=""/>
        <br />
        <b>auula</b>
        <br />
        <sub><a href="https://mp.weixin.qq.com/mp/profile_ext?action=home&__biz=MzI3MzQwNjcyNg==&scene=124#wechat_redirect">[CodingSauce]</a></sub>
      </a>
    </td>
    <td align="center">
      <a href="https://github.com/Danceiny">
        <img src="https://avatars.githubusercontent.com/u/9427454?v=4" width="64px;" alt=""/>
        <br />
        <b>Danceiny</b>
        <br />
        <sub>&nbsp;</sub>
      </a>
    </td>
    <td align="center">
      <a href="https://github.com/IceCream01">
        <img src="https://avatars.githubusercontent.com/u/19547638?v=4" width="64px;" alt=""/>
        <br />
        <b>Ice</b>
        <br />
        <sub>&nbsp;</sub>
      </a>
    </td>
    <td align="center">
      <a href="https://github.com/FishGoddess">
        <img src="https://avatars.githubusercontent.com/u/36259784?v=4" width="64px;" alt=""/>
        <br />
        <b>FishGoddess</b>
        <br />
        <sub><a href="https://github.com/FishGoddess/cachego">[cachego]</a></sub>
      </a>
    </td>
  </tr>
</table>
