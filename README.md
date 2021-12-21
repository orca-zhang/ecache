# 🦄 {cache}
<p align="center">
  <a href="#">
    <img src="https://github.com/orca-zhang/cache/raw/master/logo.svg">
  </a>
</p>

<p align="center">
  <a href="https://github.com/orca-zhang/cache/blob/master/LICENSE" alt="license MIT">
    <img src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat">
  </a>
  <a href="https://goreportcard.com/badge/github.com/orca-zhang/cache" alt="goreport A+">
    <img src="https://goreportcard.com/badge/github.com/orca-zhang/cache">
  </a>
  <a href="https://app.fossa.com/projects/git%2Bgithub.com%2Forca-zhang%2Fcache?ref=badge_shield" alt="FOSSA Status">
    <img src="https://app.fossa.com/api/projects/git%2Bgithub.com%2Forca-zhang%2Fcache.svg?type=shield"/>
  </a>
  <a href="https://orca-zhang.semaphoreci.com/projects/cache" alt="buiding pass">
    <img src="https://orca-zhang.semaphoreci.com/badges/cache.svg?style=shields">
  </a>
</p>
<p align="center">Extremely easy, fast, concurrency-safe and support distributed consistency.</p>

# 特性

- ✅ 代码量<300行、30s完成接入
- 🚀 高性能、极简设计、并发安全
- 🎉 同时支持`LRU` 和 `LFU-2`模式
- 🦖 额外小组件支持分布式一致性（WIP）

## 如何使用

- 引入包（预计5秒）
``` go
import (
    "time"

    "github.com/orca-zhang/cache"
)
```

- 定义实例（预计5秒）
> 可以放置在任意位置（全局也可以），建议就近定义
``` go
var c = cache.NewLRUCache(16, 200, 10 * time.Second)
```

- 设置缓存（预计5秒）
``` go
c.Put("uid1", o) // o可以是任意变量，一般是对象指针，存放固定的信息，比如*UserInfo
```

- 查询缓存（预计5秒）
``` go
if v, ok := c.Get("uid1"); ok {
    return v.(*UserInfo) // 不用类型断言，咱们自己控制类型
}
// 如果内存缓存没有查询到，下面再回源查redis/db
```

- 删除缓存（预计5秒）
> 在信息发生变化的地方
``` go
c.Del("uid1")
```

- 下载包（预计5秒）

> 非go modules模式：\
> sh>  ```go get -u github.com/orca-zhang/cache```

> go modules模式：\
> sh>  ```go mod tidy && go mod download```

- 运行吧
> 🎉 完美搞定 🚀 性能直接提升X倍！\
> sh>  ```go run <你的main.go文件>```

## 参数说明

- `NewLRUCache`
  - 第一个参数是桶的个数，用来分散锁的粒度，每个桶都会使用独立的锁
    - 不用担心，随意设置一个就好，`cache`会找一个等于或者略大于输入大小的2的幂次的数字，后面便于掩码计算
  - 第二个参数是每个桶所能容纳的item个数上限
    - 意味着`cache`全部写满的情况下，应该有`第一个参数✖️第二个参数`个item
  - 第三个参数是每个item的过期时间

## 最佳实践

- 复杂对象优先存放指针（注意⚠️：一旦放进去不要再修改其字段，即使拿出来也是，因为item有可能被其他人同时访问）
  - 如果需要修改，解决方案：取出字段每个单独赋值，或者用copier做一次深拷贝后在副本上修改
- 也可以存放对象（相对于上一个性能差一些，因为拿出去有拷贝）
- 缓存的对象尽可能越往业务上层越大越好（节省内存拼装和组织时间）
- 如果不想因为类似遍历的请求把热数据刷掉，可以改用[`LFU`模式](#LFU模式)，虽然有10%+的损耗（[什么是LFU](#什么是LFU)）
- 一个实例可以存储多种类型的对象，试试key格式化的时候加上前缀，用冒号分割
- 并发访问量大的场景，试试`256`、`1024`个桶，甚至更多

## 特别场景

- 空缓存哨兵（不存在的对象不用再回源）
``` go
// 设置的时候直接给`nil`就好
c.Put("uid1", nil)
```

``` go
// 读取的时候，也和正常差不多
if v, ok := c.Get("uid1"); ok {
  if v == nil { // 注意⚠️：这里需要判断是不是空缓存哨兵
    return nil // 是空缓存哨兵，那就返回没有信息或者也可以让`uid1`不出现在待回源列表里
  }
  return v.(*UserInfo)
}
// 如果内存缓存没有查询到，下面再回源查redis/db
```

- LFU模式（[什么是LFU](#什么是LFU)）
> 直接在`NewLRUCache()`后面跟`.LFU(<num>)`就好，参数`<num>`代表LFU热队列的item上限个数（每个桶）
``` go
var c = cache.NewLRUCache(16, 200, 10 * time.Second).LFU(1024)
```

# 不希望你白来

- 客官，既然来了，学点东西再走吧！
- 我想尽力让你明白`cache`做了啥，以及为什么要这么做

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
- 因为内存操作非常快，相对于redis/db你基本可以忽略不计，比如现在有一个查询接口，我们把结果缓存1秒，也就是1秒内不会请求redis/db，如果接口的qps是1000，那redis/db的请求数量降低到了1/1000（理想情况），意味着访问redis/db部分的性能提升了1000倍，听上去是不是很棒？
- 继续看，你会爱上她的！（当然也可能是他，亦或者是牠，ahaha）

### 使用场景，解决什么问题

- 高并发大流量场景
  - 缓存热点数据（比如人气最高的直播间）
  - 突发QPS削峰（比如信息流中突发新闻）
- 节省成本
  - 单机场景（不部署redis、memcache也能快速提升qps上限）
  - redis和db实例降配（能拦截大部分请求）
- 缓存不怎么会变化的数据（写少读多）
  - 比如配置等（这类数据使用地方多，会有放大效应，很多时候可能会因为这些配置热key对redis实例的规格误判，需要单独为它们升配）
- 缓存可以容忍短暂不一致的数据
  - 信息查询（用户头像、昵称、商品库存(实际下单会在db再次检查)等）
  - 配置延迟生效（过期时间10秒，那最多10秒生效）

## 设计思路

> 源自[lrucache](http://github.com/orca-zhang/lrucache)，`cache`是其升级版本

- 最下层是用原生map和存双链表的```node```实现的最基础`LRU`（最久未访问）
  - PS：我实现的其他版本（[go](https://github.com/orca-zhang/lrucache) / [C++](https://github.com/ez8-co/linked_hash) / [js](https://github.com/orca-zhang/cache.js)）在leetcode都是超越100%的解法
- 第2层包了分桶策略、并发控制、过期控制（会自动适配等于或者略大于输入大小的2的幂次个桶，便于掩码计算）
- 第2.5层用很简单的方式实现了`LFU`（最近最少使用）能力，代码不超过20行，直接看源码（搜关键词`lfu`）

### 什么是LRU
- 最久未访问的优先驱逐
- 每次被访问，item会被刷新到队列的最前面
- 队列满后再次写入新item，优先驱逐队列最后面、也就是最久未访问的item

### 什么是LFU
- 最近最少使用的优先被驱逐
- 这里就涉及到访问频率也就是次数的概念，其中一种是`LFU-K`，也就是少于K次访问的用单独的`LRU`队列存放，超过K次的另外存放
- 主要优化的场景是比如一些遍历类型的查询，批量刷缓存以后，很容易把一些本来较热的item给驱逐掉
- 为了实现简单，我们这里实现的是`LFU-2`，也就是第2次访问就放到热队列里，并不记录访问次数
- 主要优化的是热key的缓存命中率

### 分布式一致性组件（WIP - 开发中）

- 这里其实简单的利用了redis的pubsub功能
- 删除item的时候，它会通知到其他所有节点
- 某种意义上说，它只是缩小不一致时间窗口的一个方式（有网络延迟且不保证一定完成）
- 需要注意⚠️以下场景会降级（时间窗口变大）：
  - redis错误、网络错误、消费goroutine panic
  - 未在所有节点生效的情况下，比如
    - 已使用`cache`但首次添加此插件
    - 新增的item或者删除操作

### 关于性能

- 释放锁不用defer（单接口性能差20倍，看到有宣称`高性能`还用defer的，直接pass吧）
- 不用异步清理（没意义，分散到写时驱逐更合理，不易抖动）
- 没有用内存容量来控制（单个item的大小一般都有预估大小，简单控制个数即可）
- 分桶策略，自动选择2的幂次个桶（分散锁竞争，2的幂次掩码操作更快）
- key用string类型（可扩展性强；语言内建支持引用，更省内存）
- 不用虚表头（虽然绕脑一些，但是有20%左右提升）
- 选择`LFU-2`实现`LFU-K`（实现简单，近乎没有额外损耗）
- 没用整块内存（写满后复用以前的内存效果也很好，整块方式尝试过提升不大、但可读性大大降低）
- 可以直接存指针（不用序列化，如果使用`[]byte`那优势大大降低）

### 关于GC

- 就像我在C++版性能剖析器里提到的[性能优化的几个层次](https://github.com/ez8-co/ezpp#性能优化的几个层次)，单从一个层面考虑性能不高明
- 《第三层次》里有一句“没有比不存在的东西性能更快的了”（类似奥卡姆剃刀），能砍掉的东西一定不要想着优化
- 比如为了减少GC大块分配内存，却提供`[]byte`的值存储，意味着必须先序列化再拷贝，额外操作虽不在性能指标里，实际用还是要算
- 如果序列化的部分可以复用用在协议层拼接，能做到`ZeroCopy`，那也无可厚非，而`cache`存储指针直接省了额外的部分
- 有的library甚至还包了server，感觉忘了初心🤔️🤔️🤔️你的竞品是memcache、redis吗？谁会用？因为你快？
- 我所崇尚的“暴力美学”是极简，缺陷率和代码量成正比，复杂的东西早晚会被淘汰，`KISS`才是王道，`cache`一共只有不到300行，千行bug率一定的情况下，它的bug不会多

## 常见问题
> 问：一个实例可以存储多种对象吗？
- 答：可以呀，比如用前缀格式化key就可以了（像用redis那样，冒号分割），注意别搞错类型就好。

> 问：如何给不同item设置不同过期时间？
- 答：用多个缓存实例。（😄没想到吧）

> 问：如果有热热热热key问题怎么解决？
- 答：本身【本地内存缓存】就是用来抗热key的，这里可以理解成是非常非常热的key（单机几十万qps），它们最大的问题是对单一bucket锁定次数过多，影响在同一个bucket的其他数据。那么可以这样：一是改用`LFU-2`不让类似遍历的请求把热数据刷掉，二是除了增加bucket，可以用多实例（同时写入相同的item）+读随机访问某一个的方式，让热key有多个副本，不过删除（反写）的时候要注意多实例全部删除，适用于“写少读非常多”的场景，或者“写多读多”的场景可以把有变化的diff部分单独摘出来转化为“写少读多”的场景。

> 问：为什么不用虚表头方式处理双链表？太弱了吧！
- 答：2019-04-22泄漏的【[lrucache](http://github.com/orca-zhang/lrucache)】被人在V站上扒出来喷过，还真不是不会，现在的写法，虽然比pointer-to-pointer方式读起来绕脑，但是有20%左右的提升哈！（😄没想到吧）

> 问：为什么不提供int类型的key的接口？
- 答：考虑过，但是为了分布式一致性处理的简单，只提供string的接口看着也不错，用fmt.Sprint(i)也不麻烦。