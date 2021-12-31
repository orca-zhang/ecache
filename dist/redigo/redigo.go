package redigo

import (
	"github.com/gomodule/redigo/redis"
	"github.com/orca-zhang/ecache/dist"
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

// pub a payload to channel
func (g *RedigoCli) Pub(channel, payload string) error {
	conn := g.p.Get()
	defer conn.Close()

	_, err := conn.Do("PUBLISH", channel, payload)
	return err
}

// sub a payload from channel, callback uill tidy the local cache
func (g *RedigoCli) Sub(channel string, callback func(payload string)) error {
	conn := g.p.Get()
	defer conn.Close()

	psc := redis.PubSubConn{Conn: conn}
	_ = psc.Subscribe(channel)

	for {
		switch msg := psc.Receive().(type) {
		case error:
			return msg
		case redis.Message:
			callback(string(msg.Data))
		}
	}
}

func Take(r *redis.Pool) dist.RedisCli {
	return &RedigoCli{p: r}
}
