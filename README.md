[English README | è±æè¯´æ](README_en.md)

# ð¦ ecache
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
<p align="center">ä¸æ¬¾æç®è®¾è®¡ãé«æ§è½ãå¹¶åå®å¨ãæ¯æåå¸å¼ä¸è´æ§çè½»éçº§åå­ç¼å­</p>

## ç¹æ§

- ð¤ ä»£ç é<300è¡ã30så®ææ¥å¥
- ð é«æ§è½ãæç®è®¾è®¡ãå¹¶åå®å¨
- ð³ï¸âð æ¯æ`LRU` å [`LRU-2`](#LRU-2æ¨¡å¼)ä¸¤ç§æ¨¡å¼
- ð¦ é¢å¤[å°ç»ä»¶](#åå¸å¼ä¸è´æ§ç»ä»¶)æ¯æåå¸å¼ä¸è´æ§

## åºåæ§è½

> [ðï¸âð¨ï¸ç¹æçç¨ä¾](https://github.com/benchplus/gocache) [ðï¸âð¨ï¸ç¹æçç»æ](https://benchplus.github.io/gocache/dev/bench/) ï¼é¤äºç¼å­å½ä¸­çæ°å¼è¶ä½è¶å¥½ï¼
![](https://github.com/orca-zhang/ecache/raw/master/doc/benchmark.png)

> gc pauseæµè¯ç»æ [ä»£ç ç±`bigcache`æä¾](https://github.com/allegro/bigcache-bench)ï¼æ°å¼è¶ä½è¶å¥½ï¼
![](https://github.com/orca-zhang/ecache/raw/master/doc/gc.png)

### ç®åæ­£å¨çäº§ç¯å¢å¤§æµééªè¯ä¸­
- [`å·²éªè¯`]å¬ä¼å·åå°(å ç¾QPS)ï¼ç¨æ·ä¿¡æ¯ãè®¢åä¿¡æ¯ãéç½®ä¿¡æ¯
- [`éªè¯ä¸­`]æ¨éç³»ç»(å ä¸QPS)ï¼å¯è°æ´ç³»ç»éç½®ãä¿¡æ¯å»éãåºå®ä¿¡æ¯ç¼å­
- [`å¾ä¸çº¿`]è¯è®ºç³»ç»(å ä¸QPS)ï¼ç¨æ·ä¿¡æ¯ãåå¸å¼ä¸è´æ§ç»ä»¶

## å¦ä½ä½¿ç¨

#### å¼å¥åï¼é¢è®¡5ç§ï¼
``` go
import (
    "time"

    "github.com/orca-zhang/ecache"
)
```

#### å®ä¹å®ä¾ï¼é¢è®¡5ç§ï¼
> å¯ä»¥æ¾ç½®å¨ä»»æä½ç½®ï¼å¨å±ä¹å¯ä»¥ï¼ï¼å»ºè®®å°±è¿å®ä¹
``` go
var c = ecache.NewLRUCache(16, 200, 10 * time.Second)
```

#### è®¾ç½®ç¼å­ï¼é¢è®¡5ç§ï¼
``` go
c.Put("uid1", o) // `o`å¯ä»¥æ¯ä»»æåéï¼ä¸è¬æ¯å¯¹è±¡æéï¼å­æ¾åºå®çä¿¡æ¯ï¼æ¯å¦`*UserInfo`
```

#### æ¥è¯¢ç¼å­ï¼é¢è®¡5ç§ï¼
``` go
if v, ok := c.Get("uid1"); ok {
    return v.(*UserInfo) // ä¸ç¨ç±»åæ­è¨ï¼å±ä»¬èªå·±æ§å¶ç±»å
}
// å¦æåå­ç¼å­æ²¡ææ¥è¯¢å°ï¼ä¸é¢ååæºæ¥redis/db
```

#### å é¤ç¼å­ï¼é¢è®¡5ç§ï¼
> å¨ä¿¡æ¯åçååçå°æ¹
``` go
c.Del("uid1")
```

#### ä¸è½½åï¼é¢è®¡5ç§ï¼

> égo modulesæ¨¡å¼ï¼\
> sh>  ```go get -u github.com/orca-zhang/ecache```

> go modulesæ¨¡å¼ï¼\
> sh>  ```go mod tidy && go mod download```

#### è¿è¡å§
> ð å®ç¾æå® ð æ§è½ç´æ¥æåXåï¼\
> sh>  ```go run <ä½ çmain.goæä»¶>```

## åæ°è¯´æ

- `NewLRUCache`
  - ç¬¬ä¸ä¸ªåæ°æ¯æ¡¶çä¸ªæ°ï¼ç¨æ¥åæ£éçç²åº¦ï¼æ¯ä¸ªæ¡¶é½ä¼ä½¿ç¨ç¬ç«çéï¼æå¤§å¼ä¸º65535ï¼æ¯æ65536ä¸ªå®ä¾
    - ä¸ç¨æå¿ï¼éæè®¾ç½®ä¸ä¸ªå°±å¥½ï¼`ecache`ä¼æ¾ä¸ä¸ªåéçæ°å­ä¾¿äºåé¢æ©ç è®¡ç®
  - ç¬¬äºä¸ªåæ°æ¯æ¯ä¸ªæ¡¶æè½å®¹çº³çitemä¸ªæ°ä¸éï¼æå¤§å¼ä¸º65535
    - æå³ç`ecache`å¨é¨åæ»¡çæåµä¸ï¼åºè¯¥æ`ç¬¬ä¸ä¸ªåæ° X ç¬¬äºä¸ªåæ°`ä¸ªitemï¼æå¤è½æ¯æå­å¨42äº¿ä¸ªitem
  - \[`å¯é`\]ç¬¬ä¸ä¸ªåæ°æ¯æ¯ä¸ªitemçè¿ææ¶é´
    - `ecache`ä½¿ç¨åé¨è®¡æ¶å¨æåæ§è½ï¼é»è®¤100msç²¾åº¦ï¼æ¯ç§æ ¡å
    - ä¸ä¼ æèä¼ `0`ï¼ä»£è¡¨æ°¸ä¹ææ

## æä½³å®è·µ

- æ¯æä»»æç±»åçå¼
  - æä¾`Put`/`PutInt64`/`PutBytes`ä¸ç§æ¹æ³ï¼éåºä¸ååºæ¯ï¼éè¦ä¸`Get`/`GetInt64`/`GetBytes`éå¯¹ä½¿ç¨ï¼åä¸¤ç§æ¹æ³GCå¼éè¾å°ï¼
  - å¤æå¯¹è±¡ä¼åå­æ¾æéï¼æ³¨æâ ï¸ä¸æ¦æ¾è¿å»ä¸è¦åä¿®æ¹å¶å­æ®µï¼å³ä½¿åæ¿åºæ¥ä¹æ¯ï¼itemæå¯è½è¢«å¶ä»äººåæ¶è®¿é®ï¼
    - å¦æéè¦ä¿®æ¹ï¼è§£å³æ¹æ¡ï¼ååºå­æ®µæ¯ä¸ªåç¬èµå¼ï¼æèç¨[copieråä¸æ¬¡æ·±æ·è´åå¨å¯æ¬ä¸ä¿®æ¹](#éè¦ä¿®æ¹é¨åæ°æ®ä¸ç¨å¯¹è±¡æéæ¹å¼å­å¨æ¶)
    - ä¹å¯ä»¥å­æ¾å¯¹è±¡ï¼ç¸å¯¹äºç´æ¥å­å¯¹è±¡æéæ§è½å·®ä¸äºï¼å ä¸ºæ¿åºå»ææ·è´ï¼
    - ç¼å­çå¯¹è±¡å°½å¯è½è¶å¾ä¸å¡ä¸å±è¶å¤§è¶å¥½ï¼èçåå­æ¼è£åç»ç»æ¶é´ï¼
- å¦æä¸æ³å ä¸ºç±»ä¼¼éåçè¯·æ±æç­æ°æ®å·æï¼å¯ä»¥æ¹ç¨[`LRU-2`æ¨¡å¼](#LRU-2æ¨¡å¼)ï¼å¯è½æå¾å°çæèï¼ð¬ [ä»ä¹æ¯LRU-2](#ä»ä¹æ¯LRU-2)ï¼
  - `LRU2`å`LRU`çå¤§å°è®¾ç½®åå«ä¸º1/4å3/4ææè¾å¥½
- ä¸ä¸ªå®ä¾å¯ä»¥å­å¨å¤ç§ç±»åçå¯¹è±¡ï¼è¯è¯keyæ ¼å¼åçæ¶åå ä¸åç¼ï¼ç¨åå·åå²
- å¹¶åè®¿é®éå¤§çåºæ¯ï¼è¯è¯`256`ã`1024`ä¸ªæ¡¶ï¼çè³æ´å¤
- å¯ä»¥å½ä½**ç¼å²éå**ç¨äºåå¹¶æ´æ°ä»¥åå°å·çæ¬¡æ°ï¼æ°æ®å¯ä»¥éå»ºæå®¹å¿æ­çµä¸¢å¤±çæåµä¸ï¼
  - å·ä½ä½¿ç¨æ¹å¼æ¯[æè½½`Inspector`](#æ³¨å¥çå¬å¨)çå¬é©±éäºä»¶
  - ç»æ«æå®æ¶è°ç¨[`Walk`](#éåææåç´ )å°æ°æ®å·å°å­å¨

## ç¹å«åºæ¯

### æ´åé®ãæ´åå¼åå­èæ°ç»
``` go
// æ´åé®
c.Put(strconv.FormatInt(d, 10), o) // dä¸º`int64`ç±»å

// æ´åå¼
c.PutInt64("uid1", int64(1))
if d, ok := c.GetInt64("uid1"); ok {
    // dä¸º`int64`ç±»åç1
}

// å­èæ°ç»
c.PutBytes("uid1", b)// bä¸º`[]byte`ç±»å
if b, ok := c.GetBytes("uid1"); ok {
    // bä¸º`[]byte`ç±»å
}
```

### LRU-2æ¨¡å¼

- ð¬ [ä»ä¹æ¯LRU-2](#ä»ä¹æ¯LRU-2)

> ç´æ¥å¨`NewLRUCache()`åé¢è·`.LRU2(<num>)`å°±å¥½ï¼åæ°`<num>`ä»£è¡¨`LRU-2`ç­éåçitemä¸éä¸ªæ°ï¼æ¯ä¸ªæ¡¶ï¼
``` go
var c = ecache.NewLRUCache(16, 200, 10 * time.Second).LRU2(1024)
```

### ç©ºç¼å­å¨åµï¼ä¸å­å¨çå¯¹è±¡ä¸ç¨ååæºï¼
``` go
// è®¾ç½®çæ¶åç´æ¥ç»`nil`å°±å¥½
c.Put("uid1", nil)
```

``` go
// è¯»åçæ¶åï¼ä¹åæ­£å¸¸å·®ä¸å¤
if v, ok := c.Get("uid1"); ok {
  if v == nil { // æ³¨æâ ï¸è¿ééè¦å¤æ­æ¯ä¸æ¯ç©ºç¼å­å¨åµ
    return nil  // æ¯ç©ºç¼å­å¨åµï¼é£å°±è¿åæ²¡æä¿¡æ¯æèä¹å¯ä»¥è®©`uid1`ä¸åºç°å¨å¾åæºåè¡¨é
  }
  return v.(*UserInfo)
}
// å¦æåå­ç¼å­æ²¡ææ¥è¯¢å°ï¼ä¸é¢ååæºæ¥redis/db
```

### éè¦ä¿®æ¹é¨åæ°æ®ï¼ä¸ç¨å¯¹è±¡æéæ¹å¼å­å¨æ¶

> æ¯å¦ï¼æä»¬ä»`ecache`ä¸­è·åäº`*UserInfo`ç±»åçç¨æ·ä¿¡æ¯ç¼å­`v`ï¼éè¦ä¿®æ¹å¶ç¶æå­æ®µ
``` go
import (
    "github.com/jinzhu/copier"
)
```

``` go
o := &UserInfo{}
copier.Copy(o, v) // ä»`v`å¤å¶å°`o`
o.Status = 1      // ä¿®æ¹å¯æ¬çå­æ®µ
```

### æ³¨å¥çå¬å¨

``` go
// inspector - å¯ä»¥ç¨æ¥åç»è®¡æèç¼å²éåç­
//   `action`:PUT, `status`: evicted=-1, updated=0, added=1
//   `action`:GET, `status`: miss=0, hit=1
//   `action`:DEL, `status`: miss=0, hit=1
//   `iface`/`bytes`åªæå¨`status`ä¸ä¸º0æè`action`ä¸ºPUTæ¶æä¸ä¸ºnil
type inspector func(action int, key string, iface *interface{}, bytes []byte, status int)
```

- ä½¿ç¨æ¹å¼
``` go
cache.Inspect(func(action int, key string, iface *interface{}, bytes []byte, status int) {
  // TODO: å®ç°ä½ æ³åçäºæ
  //     çå¬å¨ä¼æ ¹æ®æ³¨å¥é¡ºåºä¾æ¬¡æ§è¡
  //     æ³¨æâ ï¸å¦ææèæ¶æä½ï¼å°½éå¦å¼channelä¿è¯ä¸é»å¡å½ååç¨

  // - å¦ä½è·åæ­£ç¡®çå¼ -
  //   - `Put`:      `*iface`
  //   - `PutBytes`: `bytes`
  //   - `PutInt64`: `ecache.ToInt64(bytes)`
})
```

### éåææåç´ 

``` go
  // åªä¼éåç¼å­ä¸­å­å¨ä¸æªè¿æçé¡¹
  cache.Walk(func(key string, iface *interface{}, bytes []byte, expireAt int64) bool {
    // `key`æ¯å¼ï¼`iface`/`bytes`æ¯å¼ï¼`expireAt`æ¯è¿ææ¶é´

    // - å¦ä½è·åæ­£ç¡®çå¼ -
    //   - `Put`:      `*iface`
    //   - `PutBytes`: `bytes`
    //   - `PutInt64`: `ecache.ToInt64(bytes)`
    return true // æ¯å¦ç»§ç»­éå
  })
```

## ç»è®¡ç¼å­ä½¿ç¨æåµ

> å®ç°è¶çº§ç®åï¼æ³¨å¥inspectoråï¼æ¯ä¸ªæä½åªå¤äºä¸æ¬¡åå­æä½ï¼å·ä½ç[ä»£ç ](/stats/stats.go#L34)

##### å¼å¥statså
``` go
import (
    "github.com/orca-zhang/ecache/stats"
)
```

#### ç»å®ç¼å­å®ä¾
> åç§°ä¸ºèªå®ä¹çæ± å­åç§°ï¼åé¨ä¼æåç§°èå\
> æ³¨æâ ï¸ç»å®å¯ä»¥æ¾å¨å¨å±
``` go
var _ = stats.Bind("user", c)
var _ = stats.Bind("user", c0, c1, c2)
var _ = stats.Bind("token", caches...)
```

#### è·åç»è®¡ä¿¡æ¯
``` go
stats.Stats().Range(func(k, v interface{}) bool {
    fmt.Printf("stats: %s %+v\n", k, v) // kæ¯æ± å­åç§°ï¼væ¯(*stats.StatsNode)ç±»å
    // å¶ä¸­ç»è®¡äºåç§äºä»¶çæ¬¡æ°ï¼ä½¿ç¨`HitRate`æ¹æ³å¯ä»¥è·å¾ç¼å­å½ä¸­ç
    return true
})
```

## åå¸å¼ä¸è´æ§ç»ä»¶

- ð¬ [åçè¯´æ](#åå¸å¼ä¸è´æ§ç»ä»¶åç)

### å¼å¥distå
``` go
import (
    "github.com/orca-zhang/ecache/dist"
)
```

### ç»å®ç¼å­å®ä¾
> åç§°ä¸ºèªå®ä¹çæ± å­åç§°ï¼åé¨ä¼æåç§°èå\
> æ³¨æâ ï¸ç»å®å¯ä»¥æ¾å¨å¨å±ï¼ä¸ä¾èµåå§å
``` go
var _ = dist.Bind("user", c)
var _ = dist.Bind("user", c0, c1, c2)
var _ = dist.Bind("token", caches...)
```

### ç»å®redis client
> ç®åæ¯æredigoågoredisï¼å¶ä»åºå¯ä»¥èªè¡å®ç°dist.RedisCliæ¥å£ï¼æèæissueç»æ

#### go-redis v7åä»¥ä¸çæ¬
``` go
import (
    "github.com/orca-zhang/ecache/dist/goredis/v7"
)

dist.Init(goredis.Take(redisCli)) // redisCliæ¯*redis.RedisClientç±»å
dist.Init(goredis.Take(redisCli, 100000)) // ç¬¬äºä¸ªåæ°æ¯channelç¼å²åºå¤§å°ï¼ä¸ä¼ é»è®¤100
```

#### go-redis v8åä»¥ä¸çæ¬
``` go
import (
    "github.com/orca-zhang/ecache/dist/goredis"
)

dist.Init(goredis.Take(redisCli)) // redisCliæ¯*redis.RedisClientç±»å
dist.Init(goredis.Take(redisCli, 100000)) // ç¬¬äºä¸ªåæ°æ¯channelç¼å²åºå¤§å°ï¼ä¸ä¼ é»è®¤100
```

#### redigo
> æ³¨æâ ï¸`github.com/gomodule/redigo` è¦æ±æä½çæ¬ `go 1.14`
``` go
import (
    "github.com/orca-zhang/ecache/dist/redigo"
)

dist.Init(redigo.Take(pool)) // poolæ¯*redis.Poolç±»å
```

#### ä¸»å¨éç¥ææèç¹ãææå®ä¾å é¤ï¼åæ¬æ¬æºï¼
> å½dbçæ°æ®åçååæèå é¤æ¶è°ç¨\
> åçéè¯¯æ¶ä¼éçº§æåªå¤çæ¬æºææå®ä¾ï¼æ¯å¦æªåå§åæèç½ç»éè¯¯ï¼
``` go
dist.OnDel("user", "uid1") // useræ¯æ± å­åç§°ï¼uid1æ¯è¦å é¤çkey
```

## ä½¿ç¨[`lrucache`](http://github.com/orca-zhang/lrucache)çèç¨æ·åçº§æå¯¼

- åªéåæ­¥ï¼
1. å¼å¥å `github.com/orca-zhang/lrucache` æ¹ä¸º `github.com/orca-zhang/ecache`
2. `lrucache.NewSyncCache` æ¹ä¸º `ecache.NewLRUCache`
3. ç¬¬3ä¸ªåæ°ä»é»è®¤çåä½ç§æ¹ä¸º`*time.Second`
4. `Delete`æ¹æ³æ¹ä¸º`Del`

# ä¸å¸æä½ ç½æ¥

- å®¢å®ï¼æ¢ç¶æ¥äºï¼å­¦ç¹ä¸è¥¿åèµ°å§ï¼
- ææ³å°½åè®©ä½ æç½`ecache`åäºå¥ï¼ä»¥åä¸ºä»ä¹è¦è¿ä¹å

## ä»ä¹æ¯æ¬å°åå­ç¼å­

---
    L1 ç¼å­å¼ç¨ .................... 0.5 ns
    åæ¯éè¯¯é¢æµ ...................... 5 ns
    L2 ç¼å­å¼ç¨ ...................... 7 ns
    äºæ¥é/è§£é ...................... 25 ns
    ä¸»å­å¨å¨å¼ç¨ .................... 100 ns
    ä½¿ç¨ Zippy åç¼© 1K å­è ........3,000 ns =   3 Âµs
    éè¿ 1 Gbps ç½ç»åé 2K å­è... 20,000 ns =  20 Âµs
    ä»åå­ä¸­é¡ºåºè¯»å 1 MB ........ 250,000 ns = 250 Âµs
    åä¸æ°æ®ä¸­å¿åçå¾è¿........... 500,000 ns = 0.5 ms
    åéæ°æ®å å å·<->è·å° .... 150,000,000 ns = 150 ms

- ä»ä¸è¡¨å¯ä»¥çåºï¼åå­è®¿é®åç½ç»è®¿é®(åæ°æ®ä¸­å¿)å·®ä¸å¤æ¯ä¸åå°ä¸ä¸åçå·®è·ï¼
- æ¾ç»éå°ä¸æ­¢ä¸ä¸ªå·¥ç¨å¸ï¼âç¼å­ï¼ä¸redisâï¼ä½ææ³è¯´ï¼redisä¸æ¯ä¸éæ²¹ï¼æäºç¨åº¦ä¸è®²ï¼ç¨å®è¿æ¯å©æ¢¦ï¼å½ç¶æè¯´çæ¯ç¼å­ä¸è´æ§é®é¢...ðï¼
- å ä¸ºåå­æä½éå¸¸å¿«ï¼ç¸å¯¹äºredis/dbä½ åºæ¬å¯ä»¥å¿½ç¥ä¸è®¡ï¼æ¯å¦ç°å¨æä¸ä¸ªQPSæ¯1000æ¥è¯¢APIï¼æä»¬æç»æç¼å­1ç§ï¼ä¹å°±æ¯1ç§åä¸ä¼è¯·æ±redis/dbï¼é£åæºæ¬¡æ°éä½å°äº1/1000ï¼çæ³æåµï¼ï¼æå³çè®¿é®redis/dbé¨åçæ§è½æåäº1000åï¼å¬ä¸å»æ¯ä¸æ¯å¾æ£ï¼
- ç»§ç»­çï¼ä½ ä¼ç±ä¸å¥¹çï¼ï¼å½ç¶ä¹å¯è½æ¯ä»ï¼äº¦æèæ¯ç ï¼ahahaï¼

### ä½¿ç¨åºæ¯ï¼è§£å³ä»ä¹é®é¢

- é«å¹¶åå¤§æµéåºæ¯
  - ç¼å­ç­ç¹æ°æ®ï¼æ¯å¦äººæ°æ¯è¾é«çç´æ­é´ï¼
  - çªåQPSåå³°ï¼æ¯å¦ä¿¡æ¯æµä¸­çªåæ°é»ï¼
  - éä½å»¶è¿åæ¥å µï¼æ¯å¦ç­æ¶é´åé¢ç¹è®¿é®çé¡µé¢ï¼
- èçææ¬
  - åæºåºæ¯ï¼ä¸é¨ç½²redisãmemcacheä¹è½å¿«éæåQPSä¸éï¼
  - redisådbå®ä¾ééï¼è½æ¦æªå¤§é¨åè¯·æ±ï¼
- ä¸æä¹ä¼ååçæ°æ®ï¼åå°è¯»å¤ï¼
  - æ¯å¦éç½®ç­ï¼è¿ç±»æ°æ®ä½¿ç¨å°æ¹å¤ï¼ä¼ææ¾å¤§æåºï¼å¾å¤æ¶åå¯è½ä¼å ä¸ºè¿äºéç½®ç­keyå¯¹redis/dbå®ä¾çè§æ ¼è¯¯å¤ï¼éè¦åç¬ä¸ºå®ä»¬åéï¼
- å¯ä»¥å®¹å¿ç­æä¸ä¸è´çæ°æ®
  - ç¨æ·å¤´åãæµç§°ãåååºå­(å®éä¸åä¼å¨dbåæ¬¡æ£æ¥)ç­
  - ä¿®æ¹çéç½®ï¼è¿ææ¶é´10ç§ï¼é£æå¤å»¶è¿10ç§çæï¼
- ç¼å²éåï¼åå¹¶æ´æ°ä»¥åå°å·çæ¬¡æ°
  - å¯ä»¥éè¿ç»æ¥è¯¢æè¡¥ä¸æ¥å®ç°å¼ºä¸è´ï¼åå¸å¼æåµä¸ï¼éè¦å¨è´è½½åè¡¡å±ä¿è¯åç¨æ·/è®¾å¤è°åº¦å°åä¸èç¹ï¼
  - å¯ä»¥éå»ºæå®¹å¿æ­çµä¸¢å¤±çæåµä¸

## è®¾è®¡æè·¯

> `ecache`æ¯[`lrucache`](http://github.com/orca-zhang/lrucache)åºçåçº§çæ¬

- æä¸å±æ¯ç¨åçmapååé¾è¡¨å®ç°çæåºç¡`LRU`ï¼æä¹æªè®¿é®ï¼
  - PSï¼æå®ç°çå¶ä»çæ¬ï¼[go](https://github.com/orca-zhang/lrucache) / [c++](https://github.com/ez8-co/linked_hash) / [js](https://github.com/orca-zhang/ecache.js)ï¼å¨leetcodeé½æ¯è¶è¶100%çè§£æ³
- ç¬¬2å±åäºåæ¡¶ç­ç¥ãå¹¶åæ§å¶ãè¿ææ§å¶ï¼ä¼èªå¨éæ©2çå¹æ¬¡ä¸ªæ¡¶ï¼ä¾¿äºæ©ç è®¡ç®ï¼
- ç¬¬2.5å±ç¨å¾ç®åçæ¹å¼å®ç°äº`LRU-2`è½åï¼ä»£ç ä¸è¶è¿20è¡ï¼ç´æ¥çæºç ï¼æå³é®è¯`LRU-2`ï¼

### ä»ä¹æ¯LRU
- æä¹æªè®¿é®çä¼åé©±é
- æ¯æ¬¡è¢«è®¿é®ï¼itemä¼è¢«å·æ°å°éåçæåé¢
- éåæ»¡ååæ¬¡åå¥æ°itemï¼ä¼åé©±ééåæåé¢ãä¹å°±æ¯æä¹æªè®¿é®çitem

### ä»ä¹æ¯LRU-2
- `LRU-K`æ¯å°äºKæ¬¡è®¿é®çç¨åç¬ç`LRU`éåå­æ¾ï¼è¶è¿Kæ¬¡çå¦å¤å­æ¾
- ä¸»è¦ä¼åçåºæ¯æ¯æ¯å¦ä¸äºéåç±»åçæ¥è¯¢ï¼æ¹éå·ç¼å­ä»¥åï¼å¾å®¹ææä¸äºæ¬æ¥è¾ç­çitemç»é©±éæ
- ä¸ºäºå®ç°ç®åï¼æä»¬è¿éå®ç°çæ¯`LRU-2`ï¼ä¹å°±æ¯ç¬¬2æ¬¡è®¿é®å°±æ¾å°ç­éåéï¼å¹¶ä¸è®°å½è®¿é®æ¬¡æ°
- ä¸»è¦ä¼åçæ¯ç­keyçç¼å­å½ä¸­ç
- åmysqlç[ç¼å²æ± lruç®æ³](https://dev.mysql.com/doc/refman/5.7/en/innodb-buffer-pool.html)éå¸¸ç±»ä¼¼

### åå¸å¼ä¸è´æ§ç»ä»¶åç

- å¶å®ç®åçå©ç¨äºredisçpubsubåè½
- ä¸»å¨åç¥è¢«ç¼å­çä¿¡æ¯ææ´æ°ï¼å¹¿æ­å°ææèç¹
- æç§æä¹ä¸è¯´ï¼å®åªæ¯ç¼©å°ä¸ä¸è´æ¶é´çªå£çä¸ä¸ªæ¹å¼ï¼æç½ç»å»¶è¿ä¸ä¸ä¿è¯ä¸å®å®æï¼
- éè¦æ³¨æâ ï¸ï¼
  - å°½éåå°ä½¿ç¨ï¼éåç¨å¨åå°è¯»å¤`WORM(Write-Once-Read-Many)`çåºæ¯
    - redisæ§è½æ¯ç«ä¸å¦åå­ï¼èä¸æå¹¿æ­ç±»éä¿¡ï¼åæ¾å¤§ï¼
  - ä»¥ä¸åºæ¯ä¼éçº§ï¼æ¶é´çªå£åå¤§ï¼ï¼ä½è³å°ä¼ä¿è¯å½åèç¹çå¼ºä¸è´æ§
    - redisä¸å¯ç¨ãç½ç»éè¯¯
    - æ¶è´¹goroutine panic
    - å­å¨æªçæèç¹ï¼ç°åº¦`canary`åå¸ï¼æèåå¸è¿ç¨ä¸­ï¼çæåµä¸ï¼æ¯å¦
      - å·²ä½¿ç¨`ecache`ä½é¦æ¬¡æ·»å æ­¤æä»¶
      - æ°å å¥ç¼å­çæ°æ®æèæ°å çå é¤æä½

### å³äºæ§è½

- éæ¾éä¸ç¨defer
- ä¸ç¨å¼æ­¥æ¸çï¼æ²¡æä¹ï¼åæ£å°åæ¶é©±éæ´åçï¼ä¸ææå¨ï¼
- æ²¡æç¨åå­å®¹éæ¥æ§å¶ï¼åä¸ªitemçå¤§å°ä¸è¬é½æé¢ä¼°å¤§å°ï¼ç®åæ§å¶ä¸ªæ°å³å¯ï¼
- åæ¡¶ç­ç¥ï¼èªå¨éæ©2çå¹æ¬¡ä¸ªæ¡¶ï¼åæ£éç«äºï¼2çå¹æ¬¡æ©ç æä½æ´å¿«ï¼
- keyç¨`string`ç±»åï¼å¯æ©å±æ§å¼ºï¼è¯­è¨åå»ºæ¯æå¼ç¨ï¼æ´çåå­ï¼
- ä¸ç¨èè¡¨å¤´ï¼è½ç¶ç»èä¸äºï¼ä½æ¯æ20%å·¦å³æåï¼
- éæ©`LRU-2`å®ç°`LRU-K`ï¼å®ç°ç®åï¼è¿ä¹æ²¡æé¢å¤æèï¼
- å¯ä»¥ç´æ¥å­æéï¼ä¸ç¨åºååï¼æäºåºæ¯å¦æä½¿ç¨`[]byte`é£ä¼å¿å¤§å¤§éä½ï¼
- ä½¿ç¨åé¨è®¡æ¶å¨è®¡æ¶ï¼é»è®¤100msç²¾åº¦ï¼æ¯ç§æ ¡åï¼åæåç°time.Now()äº§çä¸´æ¶å¯¹è±¡å¯¼è´GCèæ¶å¢å ï¼
- åé¾è¡¨ç¨åºå®åéåå­å­å¨ï¼ç¨æ¶é´æ³ç½®0æ¥æ è®°å é¤ï¼åå°GCï¼å¹¶ä¸åè§æ ¼æ¯`bigcache`èçåå­50%ä»¥ä¸ï¼

#### å¤±è´¥çä¼åå°è¯

- keyç±`string`æ¹ä¸º`reflect.StringHeader`ï¼ç»æï¼è´ä¼å
- äºæ¥éæ¹ä¸ºè¯»åéï¼Getè¯·æ±ä¹ä¼ä¿®æ¹æ°æ®ï¼è®¿é®è¿ä¾ï¼å³ä½¿ä¸æ¹æ°æ®ï¼ç»æï¼è¯»åæ··ååºæ¯è´ä¼å
- ç¨`time.Timer`å®ç°åé¨è®¡æ¶å¨ï¼ç»æï¼è§¦åä¸ç¨³å®ï¼åç´æ¥ç¨`time.Sleep`å®ç°è®¡æ¶å¨
- åå¸å¼ä¸è´æ§ç»ä»¶æinspectorèªå¨åæ­¥æ´æ°åå é¤ï¼ç»æï¼æ§è½å½±åè¾å¤§ä¸éè¦ç¹æ®å¤çå¾ªç¯è°ç¨é®é¢

### å³äºGCä¼å

- å°±åæå¨C++çæ§è½åæå¨éæå°ç[æ§è½ä¼åçå ä¸ªå±æ¬¡](https://github.com/ez8-co/ezpp#æ§è½ä¼åçå ä¸ªå±æ¬¡)ï¼åä»ä¸ä¸ªå±æ¬¡èèæ§è½å¹¶ä¸é«æ
- ãç¬¬ä¸å±æ¬¡ãéæä¸å¥âæ²¡ææ¯ä¸å­å¨çä¸è¥¿æ§è½æ´å¿«çäºâï¼ç±»ä¼¼å¥¥å¡å§ååï¼ï¼è½ç æä¸å®ä¸è¦æ³çä¼å
- æ¯å¦ä¸ºäºåå°GCå¤§ååéåå­ï¼å´æä¾`[]byte`çå¼å­å¨ï¼æå³çå¯è½éè¦åºååãæ·è´ï¼è½ä¸å¨åºçæ§è½ææ éï¼äººå®¶ç¨è¿æ¯è¦ç®ï¼åæ¬ï¼GCãåå­ãCPUï¼
- å¦æåºååçé¨åå¯ä»¥å¤ç¨ç¨å¨åè®®å±æ¼æ¥ï¼è½åå°`ZeroCopy`ï¼é£ä¹æ å¯åéï¼è`ecache`å­å¨æéç´æ¥çäºé¢å¤çé¨å
- ææ³è¡¨è¾¾çå¹¶ä¸æ¯GCä¼åä¸éè¦ï¼èæ´å¤åºè¯¥ç»ååºæ¯ï¼ä½¿ç¨èé¢å¤æèä¹éè¦èèï¼èéå®£ç§°gc-freeï¼ç»æç¨èµ·æ¥å¹¶éé£æ ·
- ææå´å°çâæ´åç¾å­¦âæ¯æç®ï¼ç¼ºé·çåä»£ç éææ­£æ¯ï¼å¤æçä¸è¥¿æ©æä¼è¢«æ·æ±°ï¼`KISS`ææ¯çé
- `ecache`ä¸å±åªæä¸å°300è¡ï¼åè¡bugçä¸å®çæåµä¸ï¼å®çbugä¸ä¼å¤

## å¸¸è§é®é¢
> é®ï¼ä¸ä¸ªå®ä¾å¯ä»¥å­å¨å¤ç§å¯¹è±¡åï¼
- ç­ï¼å¯ä»¥åï¼æ¯å¦å åç¼æ ¼å¼åkeyå°±å¯ä»¥äºï¼åç¨redisé£æ ·åå·åå²ï¼ï¼æ³¨æâ ï¸å«æéç±»åã

> é®ï¼å¦ä½ç»ä¸åitemè®¾ç½®ä¸åè¿ææ¶é´ï¼
- ç­ï¼ç¨å¤ä¸ªç¼å­å®ä¾ãï¼ðæ²¡æ³å°å§ï¼

> é®ï¼å¦ææç­ç­ç­ç­keyé®é¢æä¹è§£å³ï¼
- ç­ï¼æ¬èº«ãæ¬å°åå­ç¼å­ãå°±æ¯ç¨æ¥æç­keyçï¼è¿éå¯ä»¥çè§£ææ¯éå¸¸éå¸¸ç­çkeyï¼åèç¹å åä¸QPSï¼ï¼å®ä»¬æå¤§çé®é¢æ¯å¯¹åä¸bucketéå®æ¬¡æ°è¿å¤ï¼å½±åå¨åä¸ä¸ªbucketçå¶ä»æ°æ®ãé£ä¹å¯ä»¥è¿æ ·ï¼ä¸æ¯æ¹ç¨`LRU-2`ä¸è®©ç±»ä¼¼éåçè¯·æ±æç­æ°æ®å·æï¼äºæ¯é¤äºå¢å bucketï¼å¯ä»¥ç¨å¤å®ä¾ï¼åæ¶åå¥ç¸åçitemï¼+è¯»è®¿é®æä¸ä¸ªï¼æ¯å¦æè®¿é®ç¨æ·uid hashï¼çæ¹å¼ï¼è®©ç­keyæå¤ä¸ªå¯æ¬ï¼ä¸è¿å é¤ï¼ååï¼çæ¶åè¦æ³¨æå¤å®ä¾å¨é¨å é¤ï¼éç¨äºâåå°è¯»å¤`WORM(Write-Once-Read-Many)`âçåºæ¯ï¼æèâåå¤è¯»å¤âçåºæ¯å¯ä»¥ææååçdiffé¨ååç¬æåºæ¥è½¬åä¸ºâåå°è¯»å¤`WORM(Write-Once-Read-Many)`âçåºæ¯ã

> é®ï¼ä¸ºä»ä¹ä¸ç¨èè¡¨å¤´æ¹å¼å¤çåé¾è¡¨ï¼å¤ªå¼±äºå§ï¼
- ç­ï¼2019-04-22æ³æ¼çã[lrucache](http://github.com/orca-zhang/lrucache)ãè¢«äººå¨Vç«ä¸æåºæ¥å·è¿ï¼è¿çä¸æ¯ä¸ä¼ï¼ç°å¨çåæ³ï¼è½ç¶æ¯pointer-to-pointeræ¹å¼è¯»èµ·æ¥ç»èï¼ä½æ¯æ20%å·¦å³çæååï¼ï¼ðæ²¡æ³å°å§ï¼

## è´è°¢

æè°¢å¨å¼åè¿ç¨ä¸­è¿è¡code reviewãåè¯¯ & æåºå®è´µå»ºè®®çåä½ï¼ï¼æåä¸åååï¼

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
        <sub><a href="https://mp.weixin.qq.com/mp/profile_ext?action=home&__biz=MzI3MzQwNjcyNg==&scene=124#wechat_redirect">[æç å ]</a></sub>
      </a>
    </td>
    <td align="center">
      <a href="https://github.com/Danceiny">
        <img src="https://avatars.githubusercontent.com/u/9427454?v=4" width="64px;" alt=""/>
        <br />
        <b>é»æ¯</b>
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
        <b>æ°´ä¸è¦é±¼</b>
        <br />
        <sub><a href="https://github.com/FishGoddess/cachego">[cachego]</a></sub>
      </a>
    </td>
  </tr>
</table>

## èµå©

éè¿æä¸ºèµå©åæ¥æ¯æè¿ä¸ªé¡¹ç®ã æ¨çlogoå°æ¾ç¤ºå¨æ­¤å¤ï¼å¹¶å¸¦ææåæ¨ç½ç«çé¾æ¥ã [[æä¸ºèµå©å](https://opencollective.com/ecache#sponsor)]

<a href="https://opencollective.com/ecache/sponsor/0/website" target="_blank"><img src="https://opencollective.com/ecache/sponsor/0/avatar.svg"></a>
<a href="https://opencollective.com/ecache/sponsor/1/website" target="_blank"><img src="https://opencollective.com/ecache/sponsor/1/avatar.svg"></a>
<a href="https://opencollective.com/ecache/sponsor/2/website" target="_blank"><img src="https://opencollective.com/ecache/sponsor/2/avatar.svg"></a>
<a href="https://opencollective.com/ecache/sponsor/3/website" target="_blank"><img src="https://opencollective.com/ecache/sponsor/3/avatar.svg"></a>

## è´¡ç®è

è¿ä¸ªé¡¹ç®çå­å¨è¦æè°¢ææååºè´¡ç®çäººã

è¯·ç»æä»¬ä¸ä¸ªðstarðæ¥æ¯ææä»¬ï¼è°¢è°¢ã

å¹¶æè°¢æä»¬ææçæ¯æèï¼ ð

<a href="https://opencollective.com/ecache/backer/0/website?requireActive=false" target="_blank"><img src="https://opencollective.com/ecache/backer/0/avatar.svg?requireActive=false"></a>
<a href="https://opencollective.com/ecache/backer/1/website?requireActive=false" target="_blank"><img src="https://opencollective.com/ecache/backer/1/avatar.svg?requireActive=false"></a>
<a href="https://opencollective.com/ecache/backer/2/website?requireActive=false" target="_blank"><img src="https://opencollective.com/ecache/backer/2/avatar.svg?requireActive=false"></a>
<a href="https://opencollective.com/ecache/backer/3/website?requireActive=false" target="_blank"><img src="https://opencollective.com/ecache/backer/3/avatar.svg?requireActive=false"></a>
<a href="https://opencollective.com/ecache#backers" target="_blank"><img src="https://opencollective.com/ecache/contributors.svg?width=890" /></a>
