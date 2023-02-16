[English README | 英文说明](README_en.md)

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
<p align="center">一款极简设计、高性能、并发安全、支持分布式一致性的轻量级内存缓存</p>

## 特性

- 🤏 代码量<300行、30s完成接入
- 🚀 高性能、极简设计、并发安全
- 🌈 支持`LRU` 和 [`LRU-2`](#LRU-2模式)两种模式
- 🦖 额外[小组件](#分布式一致性组件)支持分布式一致性

## 基准性能

> :snail: 代表很慢, :airplane: 代表快, :rocket: 代表非常快

> [👁️‍🗨️点我看用例](https://github.com/benchplus/gocache) [👁️‍🗨️点我看结果](https://benchplus.github.io/gocache/dev/bench/) （除了缓存命中率数值越低越好）

<table style="text-align: center">
   <tr>
      <td></td>
      <td><a href="https://github.com/allegro/bigcache">bigcache</a></td>
      <td><a href="https://github.com/FishGoddess/cachego">cachego</a></td>
      <td><a href="https://github.com/orca-zhang/ecache"><strong>ecache🌟</strong></a></td>
      <td><a href="https://github.com/coocood/freecache">freecache</a></td>
      <td><a href="https://github.com/bluele/gcache">gcache</a></td>
      <td><a href="https://github.com/patrickmn/go-cache">gocache</a></td>
   </tr>
   <tr>
      <td>PutInt</td>
      <td>:airplane:</td>
      <td></td>
      <td>:rocket:</td>
      <td>:rocket:</td>
      <td>:airplane:</td>
      <td>:airplane:</td>
   </tr>
   <tr>
      <td>GetInt</td>
      <td>:airplane:</td>
      <td>:airplane:</td>
      <td>:rocket:</td>
      <td></td>
      <td>:airplane:</td>
      <td>:airplane:</td>
   </tr>
   <tr>
      <td>Put1K</td>
      <td>:airplane:</td>
      <td>:airplane:</td>
      <td>:rocket:</td>
      <td>:rocket:</td>
      <td>:rocket:</td>
      <td>:airplane:</td>
   </tr>
   <tr>
      <td>Put1M</td>
      <td>:snail:</td>
      <td></td>
      <td>:rocket:</td>
      <td>:snail:</td>
      <td>:airplane:</td>
      <td>:airplane:</td>
   </tr>
   <tr>
      <td>PutTinyObject</td>
      <td>:airplane:</td>
      <td></td>
      <td>:rocket:</td>
      <td>:rocket:</td>
      <td>:airplane:</td>
      <td></td>
   </tr>
   <tr>
      <td>ChangeOutAllInt</td>
      <td>:airplane:</td>
      <td></td>
      <td>:rocket:</td>
      <td>:rocket:</td>
      <td>:airplane:</td>
      <td>:airplane:</td>
   </tr>
   <tr>
      <td>HeavyReadInt</td>
      <td>:rocket:</td>
      <td>:rocket:</td>
      <td>:rocket:</td>
      <td></td>
      <td></td>
      <td>:rocket:</td>
   </tr>
   <tr>
      <td>HeavyReadIntGC</td>
      <td>:airplane:</td>
      <td>:rocket:</td>
      <td>:rocket:</td>
      <td></td>
      <td>:airplane:</td>
      <td>:airplane:</td>
   </tr>
   <tr>
      <td>HeavyWriteInt</td>
      <td>:rocket:</td>
      <td>:airplane:</td>
      <td>:rocket:</td>
      <td>:rocket:</td>
      <td></td>
      <td>:airplane:</td>
   </tr>
   <tr>
      <td>HeavyWriteIntGC</td>
      <td>:rocket:</td>
      <td></td>
      <td>:airplane:</td>
      <td>:airplane:</td>
      <td></td>
      <td></td>
   </tr>
   <tr>
      <td>HeavyWrite1K</td>
      <td>:snail:</td>
      <td>:airplane:</td>
      <td>:rocket:</td>
      <td>:rocket:</td>
      <td></td>
      <td>:airplane:</td>
   </tr>
   <tr>
      <td>HeavyWrite1KGC</td>
      <td>:snail:</td>
      <td>:airplane:</td>
      <td>:rocket:</td>
      <td>:rocket:</td>
      <td></td>
      <td>:airplane:</td>
   </tr>
   <tr>
      <td>HeavyMixedInt</td>
      <td>:rocket:</td>
      <td>:airplane:</td>
      <td>:rocket:</td>
      <td></td>
      <td>:airplane:</td>
      <td>:rocket:</td>
   </tr>
   <tr>
    <td colspan="7">
      <a href="https://github.com/FishGoddess/cachego"><strong>FishGoddess/cachego</strong></a> 和 <a href="https://github.com/patrickmn/go-cache"><strong>patrickmn/go-cache</strong></a> 是简单的map+过期时间的实现，所以没有命中率测试
    </td>
   </tr>
   <tr>
    <td colspan="7">
      <a href="https://github.com/kpango/gache"><strong>kpango/gache</strong></a> & <a href="https://github.com/hlts2/gocache"><strong>hlts2/gocache</strong></a> 性能表现不是很好，所以从列表中剔除
    </td>
   </tr>
   <tr>
    <td colspan="7">
      <a href="https://github.com/patrickmn/go-cache"><strong>patrickmn/go-cache</strong></a> 是FIFO模式，其他的库都是LRU模式
    </td>
   </tr>
</table>

![](https://github.com/orca-zhang/ecache/raw/master/doc/benchmark.png)

> gc pause测试结果 [代码由`bigcache`提供](https://github.com/allegro/bigcache-bench)（数值越低越好）
![](https://github.com/orca-zhang/ecache/raw/master/doc/gc.png)

### 目前正在生产环境大流量验证中
- [`已验证`]公众号后台(几百QPS)：用户信息、订单信息、配置信息
- [`已验证`]推送系统(几万QPS)：可调整系统配置、信息去重、固定信息缓存
- [`已验证`]评论系统(几万QPS)：用户信息、分布式一致性组件

## 如何使用

#### 引入包（预计5秒）
``` go
import (
    "time"

    "github.com/orca-zhang/ecache"
)
```

#### 定义实例（预计5秒）
> 可以放置在任意位置（全局也可以），建议就近定义
``` go
var c = ecache.NewLRUCache(16, 200, 10 * time.Second)
```

#### 设置缓存（预计5秒）
``` go
c.Put("uid1", o) // `o`可以是任意变量，一般是对象指针，存放固定的信息，比如`*UserInfo`
```

#### 查询缓存（预计5秒）
``` go
if v, ok := c.Get("uid1"); ok {
    return v.(*UserInfo) // 不用类型断言，咱们自己控制类型
}
// 如果内存缓存没有查询到，下面再回源查redis/db
```

#### 删除缓存（预计5秒）
> 在信息发生变化的地方
``` go
c.Del("uid1")
```

#### 下载包（预计5秒）

> 非go modules模式：\
> sh>  ```go get -u github.com/orca-zhang/ecache```

> go modules模式：\
> sh>  ```go mod tidy && go mod download```

#### 运行吧
> 🎉 完美搞定 🚀 性能直接提升X倍！\
> sh>  ```go run <你的main.go文件>```

## 参数说明

- `NewLRUCache`
  - 第一个参数是桶的个数，用来分散锁的粒度，每个桶都会使用独立的锁，最大值为65535，支持65536个实例
    - 不用担心，随意设置一个就好，`ecache`会找一个合适的数字便于后面掩码计算
  - 第二个参数是每个桶所能容纳的item个数上限，最大值为65535
    - 意味着`ecache`全部写满的情况下，应该有`第一个参数 X 第二个参数`个item，最多能支持存储42亿个item
  - \[`可选`\]第三个参数是每个item的过期时间
    - `ecache`使用内部计时器提升性能，默认100ms精度，每秒校准
    - 不传或者传`0`，代表永久有效

## 最佳实践

- 支持任意类型的值
  - 提供`Put`/`PutInt64`/`PutBytes`三种方法，适应不同场景，需要与`Get`/`GetInt64`/`GetBytes`配对使用（后两种方法GC开销较小）
  - 复杂对象优先存放指针（注意⚠️一旦放进去不要再修改其字段，即使再拿出来也是，item有可能被其他人同时访问）
    - 如果需要修改，解决方案：取出字段每个单独赋值，或者用[copier做一次深拷贝后在副本上修改](#需要修改部分数据且用对象指针方式存储时)
    - 也可以存放对象（相对于直接存对象指针性能差一些，因为拿出去有拷贝）
    - 缓存的对象尽可能越往业务上层越大越好（节省内存拼装和组织时间）
- 如果不想因为类似遍历的请求把热数据刷掉，可以改用[`LRU-2`模式](#LRU-2模式)，可能有很少的损耗（💬 [什么是LRU-2](#什么是LRU-2)）
  - `LRU2`和`LRU`的大小设置分别为1/4和3/4效果较好
- 一个实例可以存储多种类型的对象，试试key格式化的时候加上前缀，用冒号分割
- 并发访问量大的场景，试试`256`、`1024`个桶，甚至更多
- 可以当作**缓冲队列**用于合并更新以减少刷盘次数（数据可以重建或容忍断电丢失的情况下）
  - 具体使用方式是[挂载`Inspector`](#注入监听器)监听驱逐事件
  - 终末或定时调用[`Walk`](#遍历所有元素)将数据刷到存储

## 特别场景

### 整型键、整型值和字节数组
``` go
// 整型键
c.Put(strconv.FormatInt(d, 10), o) // d为`int64`类型

// 整型值
c.PutInt64("uid1", int64(1))
if d, ok := c.GetInt64("uid1"); ok {
    // d为`int64`类型的1
}

// 字节数组
c.PutBytes("uid1", b)// b为`[]byte`类型
if b, ok := c.GetBytes("uid1"); ok {
    // b为`[]byte`类型
}
```

### LRU-2模式

- 💬 [什么是LRU-2](#什么是LRU-2)

> 直接在`NewLRUCache()`后面跟`.LRU2(<num>)`就好，参数`<num>`代表`LRU-2`热队列的item上限个数（每个桶）
``` go
var c = ecache.NewLRUCache(16, 200, 10 * time.Second).LRU2(1024)
```

### 空缓存哨兵（不存在的对象不用再回源）
``` go
// 设置的时候直接给`nil`就好
c.Put("uid1", nil)
```

``` go
// 读取的时候，也和正常差不多
if v, ok := c.Get("uid1"); ok {
  if v == nil { // 注意⚠️这里需要判断是不是空缓存哨兵
    return nil  // 是空缓存哨兵，那就返回没有信息或者也可以让`uid1`不出现在待回源列表里
  }
  return v.(*UserInfo)
}
// 如果内存缓存没有查询到，下面再回源查redis/db
```

### 需要修改部分数据，且用对象指针方式存储时

> 比如，我们从`ecache`中获取了`*UserInfo`类型的用户信息缓存`v`，需要修改其状态字段
``` go
import (
    "github.com/jinzhu/copier"
)
```

``` go
o := &UserInfo{}
copier.Copy(o, v) // 从`v`复制到`o`
o.Status = 1      // 修改副本的字段
```

### 注入监听器

``` go
// inspector - 可以用来做统计或者缓冲队列等
//   `action`:PUT, `status`: evicted=-1, updated=0, added=1
//   `action`:GET, `status`: miss=0, hit=1
//   `action`:DEL, `status`: miss=0, hit=1
//   `iface`/`bytes`只有在`status`不为0或者`action`为PUT时才不为nil
type inspector func(action int, key string, iface *interface{}, bytes []byte, status int)
```

- 使用方式
``` go
cache.Inspect(func(action int, key string, iface *interface{}, bytes []byte, status int) {
  // TODO: 实现你想做的事情
  //     监听器会根据注入顺序依次执行
  //     注意⚠️如果有耗时操作，尽量另开channel保证不阻塞当前协程

  // - 如何获取正确的值 -
  //   - `Put`:      `*iface`
  //   - `PutBytes`: `bytes`
  //   - `PutInt64`: `ecache.ToInt64(bytes)`
})
```

### 遍历所有元素

``` go
  // 只会遍历缓存中存在且未过期的项
  cache.Walk(func(key string, iface *interface{}, bytes []byte, expireAt int64) bool {
    // `key`是值，`iface`/`bytes`是值，`expireAt`是过期时间

    // - 如何获取正确的值 -
    //   - `Put`:      `*iface`
    //   - `PutBytes`: `bytes`
    //   - `PutInt64`: `ecache.ToInt64(bytes)`
    return true // 是否继续遍历
  })
```

## 统计缓存使用情况

> 实现超级简单，注入inspector后，每个操作只多了一次原子操作，具体看[代码](/stats/stats.go#L34)

##### 引入stats包
``` go
import (
    "github.com/orca-zhang/ecache/stats"
)
```

#### 绑定缓存实例
> 名称为自定义的池子名称，内部会按名称聚合\
> 注意⚠️绑定可以放在全局
``` go
var _ = stats.Bind("user", c)
var _ = stats.Bind("user", c0, c1, c2)
var _ = stats.Bind("token", caches...)
```

#### 获取统计信息
``` go
stats.Stats().Range(func(k, v interface{}) bool {
    fmt.Printf("stats: %s %+v\n", k, v) // k是池子名称，v是(*stats.StatsNode)类型
    // 其中统计了各种事件的次数，使用`HitRate`方法可以获得缓存命中率
    return true
})
```

## 分布式一致性组件

- 💬 [原理说明](#分布式一致性组件原理)

### 引入dist包
``` go
import (
    "github.com/orca-zhang/ecache/dist"
)
```

### 绑定缓存实例
> 名称为自定义的池子名称，内部会按名称聚合\
> 注意⚠️绑定可以放在全局，不依赖初始化
``` go
var _ = dist.Bind("user", c)
var _ = dist.Bind("user", c0, c1, c2)
var _ = dist.Bind("token", caches...)
```

### 绑定redis client
> 目前支持redigo和goredis，其他库可以自行实现dist.RedisCli接口，或者提issue给我

#### go-redis v7及以下版本
``` go
import (
    "github.com/orca-zhang/ecache/dist/goredis/v7"
)

dist.Init(goredis.Take(redisCli)) // redisCli是*redis.RedisClient类型
dist.Init(goredis.Take(redisCli, 100000)) // 第二个参数是channel缓冲区大小，不传默认100
```

#### go-redis v8及以上版本
``` go
import (
    "github.com/orca-zhang/ecache/dist/goredis"
)

dist.Init(goredis.Take(redisCli)) // redisCli是*redis.RedisClient类型
dist.Init(goredis.Take(redisCli, 100000)) // 第二个参数是channel缓冲区大小，不传默认100
```

#### redigo
> 注意⚠️`github.com/gomodule/redigo` 要求最低版本 `go 1.14`
``` go
import (
    "github.com/orca-zhang/ecache/dist/redigo"
)

dist.Init(redigo.Take(pool)) // pool是*redis.Pool类型
```

#### 主动通知所有节点、所有实例删除（包括本机）
> 当db的数据发生变化或者删除时调用\
> 发生错误时会降级成只处理本机所有实例（比如未初始化或者网络错误）
``` go
dist.OnDel("user", "uid1") // user是池子名称，uid1是要删除的key
```

## 使用[`lrucache`](http://github.com/orca-zhang/lrucache)的老用户升级指导

- 只需四步：
1. 引入包 `github.com/orca-zhang/lrucache` 改为 `github.com/orca-zhang/ecache`
2. `lrucache.NewSyncCache` 改为 `ecache.NewLRUCache`
3. 第3个参数从默认的单位秒改为`*time.Second`
4. `Delete`方法改为`Del`

# 不希望你白来

- 客官，既然来了，学点东西再走吧！
- 我想尽力让你明白`ecache`做了啥，以及为什么要这么做

## 什么是本地内存缓存

---
    L1 缓存引用 .................... 0.5 ns
    分支错误预测 ...................... 5 ns
    L2 缓存引用 ...................... 7 ns
    互斥锁/解锁 ...................... 25 ns
    主存储器引用 .................... 100 ns
    使用 Zippy 压缩 1K 字节 ........3,000 ns =   3 µs
    通过 1 Gbps 网络发送 2K 字节... 20,000 ns =  20 µs
    从内存中顺序读取 1 MB ........ 250,000 ns = 250 µs
    同一数据中心内的往返........... 500,000 ns = 0.5 ms
    发送数据包 加州<->荷兰 .... 150,000,000 ns = 150 ms

- 从上表可以看出，内存访问和网络访问(同数据中心)差不多是一千到一万倍的差距！
- 曾经遇到不止一个工程师：“缓存？上redis”，但我想说，redis不是万金油，某些程度上讲，用它还是噩梦（当然我说的是缓存一致性问题...😄）
- 因为内存操作非常快，相对于redis/db你基本可以忽略不计，比如现在有一个QPS是1000查询API，我们把结果缓存1秒，也就是1秒内不会请求redis/db，那回源次数降低到了1/1000（理想情况），意味着访问redis/db部分的性能提升了1000倍，听上去是不是很棒？
- 继续看，你会爱上她的！（当然也可能是他，亦或者是牠，ahaha）

### 使用场景，解决什么问题

- 高并发大流量场景
  - 缓存热点数据（比如人气比较高的直播间）
  - 突发QPS削峰（比如信息流中突发新闻）
  - 降低延迟和拥堵（比如短时间内频繁访问的页面）
- 节省成本
  - 单机场景（不部署redis、memcache也能快速提升QPS上限）
  - redis和db实例降配（能拦截大部分请求）
- 不怎么会变化的数据（写少读多）
  - 比如配置等（这类数据使用地方多，会有放大效应，很多时候可能会因为这些配置热key对redis/db实例的规格误判，需要单独为它们升配）
- 可以容忍短暂不一致的数据
  - 用户头像、昵称、商品库存(实际下单会在db再次检查)等
  - 修改的配置（过期时间10秒，那最多延迟10秒生效）
- 缓冲队列：合并更新以减少刷盘次数
  - 可以通过给查询打补丁来实现强一致（分布式情况下，需要在负载均衡层保证同用户/设备调度到同一节点）
  - 可以重建或容忍断电丢失的情况下

## 设计思路

> `ecache`是[`lrucache`](http://github.com/orca-zhang/lrucache)库的升级版本

- 最下层是用原生map和双链表实现的最基础`LRU`（最久未访问）
  - PS：我实现的其他版本（[go](https://github.com/orca-zhang/lrucache) / [c++](https://github.com/ez8-co/linked_hash) / [js](https://github.com/orca-zhang/ecache.js)）在leetcode都是超越100%的解法
- 第2层包了分桶策略、并发控制、过期控制（会自动选择2的幂次个桶，便于掩码计算）
- 第2.5层用很简单的方式实现了`LRU-2`能力，代码不超过20行，直接看源码（搜关键词`LRU-2`）

### 什么是LRU
- 最久未访问的优先驱逐
- 每次被访问，item会被刷新到队列的最前面
- 队列满后再次写入新item，优先驱逐队列最后面、也就是最久未访问的item

### 什么是LRU-2
- `LRU-K`是少于K次访问的用单独的`LRU`队列存放，超过K次的另外存放
- 主要优化的场景是比如一些遍历类型的查询，批量刷缓存以后，很容易把一些本来较热的item给驱逐掉
- 为了实现简单，我们这里实现的是`LRU-2`，也就是第2次访问就放到热队列里，并不记录访问次数
- 主要优化的是热key的缓存命中率
- 和mysql的[缓冲池lru算法](https://dev.mysql.com/doc/refman/5.7/en/innodb-buffer-pool.html)非常类似

### 分布式一致性组件原理

- 其实简单的利用了redis的pubsub功能
- 主动告知被缓存的信息有更新，广播到所有节点
- 某种意义上说，它只是缩小不一致时间窗口的一个方式（有网络延迟且不保证一定完成）
- 需要注意⚠️：
  - 尽量减少使用，适合用在写少读多`WORM(Write-Once-Read-Many)`的场景
    - redis性能毕竟不如内存，而且有广播类通信（写放大）
  - 以下场景会降级（时间窗口变大），但至少会保证当前节点的强一致性
    - redis不可用、网络错误
    - 消费goroutine panic
    - 存在未生效节点（灰度`canary`发布，或者发布过程中）的情况下，比如
      - 已使用`ecache`但首次添加此插件
      - 新加入缓存的数据或者新加的删除操作

### 关于性能

- 释放锁不用defer
- 不用异步清理（没意义，分散到写时驱逐更合理，不易抖动）
- 没有用内存容量来控制（单个item的大小一般都有预估大小，简单控制个数即可）
- 分桶策略，自动选择2的幂次个桶（分散锁竞争，2的幂次掩码操作更快）
- key用`string`类型（可扩展性强；语言内建支持引用，更省内存）
- 不用虚表头（虽然绕脑一些，但是有20%左右提升）
- 选择`LRU-2`实现`LRU-K`（实现简单，近乎没有额外损耗）
- 可以直接存指针（不用序列化，有些场景如果使用`[]byte`那优势大大降低）
- 使用内部计时器计时（默认100ms精度，每秒校准，剖析发现time.Now()产生临时对象导致GC耗时增加）
- 双链表用固定分配内存存储，用时间戳置0来标记删除，减少GC（并且同规格比`bigcache`节省内存50%以上）

#### 失败的优化尝试

- key由`string`改为`reflect.StringHeader`，结果：负优化
- 互斥锁改为读写锁，Get请求也会修改数据，访问违例，即使不改数据，结果：读写混合场景负优化
- 用`time.Timer`实现内部计时器，结果：触发不稳定，后直接用`time.Sleep`实现计时器
- 分布式一致性组件挂inspector自动同步更新和删除，结果：性能影响较大且需要特殊处理循环调用问题

### 关于GC优化

- 就像我在C++版性能剖析器里提到的[性能优化的几个层次](https://github.com/ez8-co/ezpp#性能优化的几个层次)，单从一个层次考虑性能并不高明
- 《第三层次》里有一句“没有比不存在的东西性能更快的了”（类似奥卡姆剃刀），能砍掉一定不要想着优化
- 比如为了减少GC大块分配内存，却提供`[]byte`的值存储，意味着可能需要序列化、拷贝（虽不在库的性能指标里，人家用还是要算，包括：GC、内存、CPU）
- 如果序列化的部分可以复用用在协议层拼接，能做到`ZeroCopy`，那也无可厚非，但实际分层以后，无法在协议层直接实现拼接，而`ecache`存储指针直接省了额外的部分
- 我想表达的并不是GC优化不重要，而更多应该结合场景，使用者额外损耗也需要考虑，而非宣称gc-free，结果用起来并非那样
- 我所崇尚的“暴力美学”是极简，缺陷率和代码量成正比，复杂的东西早晚会被淘汰，`KISS`才是王道
- `ecache`一共只有不到300行，千行bug率一定的情况下，它的bug不会多

## 常见问题
> 问：一个实例可以存储多种对象吗？
- 答：可以呀，比如加前缀格式化key就可以了（像用redis那样冒号分割），注意⚠️别搞错类型。

> 问：如何给不同item设置不同过期时间？
- 答：用多个缓存实例。（😄没想到吧）

> 问：如果有热热热热key问题怎么解决？
- 答：本身【本地内存缓存】就是用来扛住热key的，这里可以理解成是非常非常热的key（单节点几十万QPS），它们最大的问题是对单一bucket锁定次数过多，影响在同一个bucket的其他数据。那么可以这样：一是改用`LRU-2`不让类似遍历的请求把热数据刷掉，二是除了增加bucket，可以用多实例（同时写入相同的item）+读访问某一个（比如按访问用户uid hash）的方式，让热key有多个副本，不过删除（反写）的时候要注意多实例全部删除，适用于“写少读多`WORM(Write-Once-Read-Many)`”的场景，或者“写多读多”的场景可以把有变化的diff部分单独摘出来转化为“写少读多`WORM(Write-Once-Read-Many)`”的场景。

> 问：如果同一时间并发回源到DB查询同一个资源怎么优化？
- 答：可以使用[sync/singleflight](https://pkg.go.dev/golang.org/x/sync/singleflight)包，同时访问同一个资源时，只回源一次，防止热点数据把DB打爆的问题。

> 问：为什么不用虚表头方式处理双链表？太弱了吧！
- 答：2019-04-22泄漏的【[lrucache](http://github.com/orca-zhang/lrucache)】被人在V站上扒出来喷过，还真不是不会，现在的写法，虽然比pointer-to-pointer方式读起来绕脑，但是有20%左右的提升哈！（😄没想到吧）

## 相关文献

- [如何一步步提升Go内存缓存性能](https://my.oschina.net/u/5577511/blog/5438484)

## 致谢

感谢在开发过程中进行code review、勘误 & 提出宝贵建议的各位！（排名不分先后）

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
        <b>Leon Ding</b>
        <br />
        <sub><a href="https://mp.weixin.qq.com/mp/profile_ext?action=home&__biz=MzI3MzQwNjcyNg==&scene=124#wechat_redirect">[打码匠]</a></sub>
      </a>
    </td>
    <td align="center">
      <a href="https://github.com/Danceiny">
        <img src="https://avatars.githubusercontent.com/u/9427454?v=4" width="64px;" alt=""/>
        <br />
        <b>黄振</b>
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
        <b>水不要鱼</b>
        <br />
        <sub><a href="https://github.com/FishGoddess/cachego">[cachego]</a></sub>
      </a>
    </td>
  </tr>
</table>

## 赞助

通过成为赞助商来支持这个项目。 您的logo将显示在此处，并带有指向您网站的链接。 [[成为赞助商](https://opencollective.com/ecache#sponsor)]

<a href="https://opencollective.com/ecache/sponsor/0/website" target="_blank"><img src="https://opencollective.com/ecache/sponsor/0/avatar.svg"></a>
<a href="https://opencollective.com/ecache/sponsor/1/website" target="_blank"><img src="https://opencollective.com/ecache/sponsor/1/avatar.svg"></a>
<a href="https://opencollective.com/ecache/sponsor/2/website" target="_blank"><img src="https://opencollective.com/ecache/sponsor/2/avatar.svg"></a>
<a href="https://opencollective.com/ecache/sponsor/3/website" target="_blank"><img src="https://opencollective.com/ecache/sponsor/3/avatar.svg"></a>

## 贡献者

这个项目的存在要感谢所有做出贡献的人。

请给我们一个💖star💖来支持我们，谢谢。

并感谢我们所有的支持者！ 🙏

<a href="https://opencollective.com/ecache/backer/0/website?requireActive=false" target="_blank"><img src="https://opencollective.com/ecache/backer/0/avatar.svg?requireActive=false"></a>
<a href="https://opencollective.com/ecache/backer/1/website?requireActive=false" target="_blank"><img src="https://opencollective.com/ecache/backer/1/avatar.svg?requireActive=false"></a>
<a href="https://opencollective.com/ecache/backer/2/website?requireActive=false" target="_blank"><img src="https://opencollective.com/ecache/backer/2/avatar.svg?requireActive=false"></a>
<a href="https://opencollective.com/ecache/backer/3/website?requireActive=false" target="_blank"><img src="https://opencollective.com/ecache/backer/3/avatar.svg?requireActive=false"></a>
<a href="https://opencollective.com/ecache#backers" target="_blank"><img src="https://opencollective.com/ecache/contributors.svg?width=890" /></a>
