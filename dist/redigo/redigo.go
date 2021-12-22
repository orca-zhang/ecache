package dist

import (
	"github.com/gomodule/redigo/redis"
	"github.com/orca-zhang/cache/dist"
)

type RedigoCli struct {
	p *redis.Pool
}

// if the redis client is ready
func (g *RedigoCli) OK() bool {
	conn := g.p.Get()
	defer conn.Close()

	_, err := conn.Do("PING")
	return err == nil
}

// pub a key to channel
func (g *RedigoCli) Pub(channel, key string) error {
	conn := g.p.Get()
	defer conn.Close()

	_, err := conn.Do("PUBLISH", channel, key)
	return err
}

// sub a key from channel, callback uill tidy the local cache
func (g *RedigoCli) Sub(channel string, callback func(payload string)) error {
	conn := g.p.Get()
	defer conn.Close()

	psc := redis.PubSubConn{Conn: conn}
	_ = psc.Subscribe(channel)

	for {
		switch n := psc.Receive().(type) {
		case error:
			continue
		case redis.Message:
			callback(string(n.Data))
		}
	}
}

func Redigo(r *redis.Pool) dist.RedisCli {
	return &RedigoCli{p: r}
}
