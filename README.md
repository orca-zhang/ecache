# 🦄 Cache

<p align="center">[cache] that support distributed consistency and very easy to use</p>
<p align="center">
  <a href="https://github.com/orca-zhang/cache/blob/master/LICENSE">
    <img src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat">
  </a>
</p>

## 特性

- [x] 代码量<300行，30s完成接入
- [x] 高性能、极简设计、并发安全
- [x] 同时支持`LRU` 和 `LFU-2`模式
- [x] 上市公司生产环境大流量验证
- [x] 额外小插件支持分布式一致性

## 如何使用

- 引入包（预计5s）
``` go
import (
    "github.com/orca-zhang/cache"
)
```

- 定义缓存实例子（预计5s）
> 可以放置在任意位置，建议就近定义
``` go
var c = cache.NewLRUCache(1, 200, 10 * time.Second)
```

- 设置缓存（预计5s）
``` go
c.Put("uid1", o) // o可以是任意变量，一般是对象指针，存放固定的信息，比如*UserInfo
```

- 查询缓存（预计5s）
``` go
if v, ok := c.Get("uid1"); ok {
    return v.(*UserInfo) // 不用类型断言，咱们自己控制类型
}
// 如果内存缓存没有查询到，下面再回源查redis/db
```

- 下包>>编译>>运行（此处预计若干秒）
> 搞定，性能直接提升X倍！
