package dist

import (
	"github.com/go-redis/redis/v7"
	"github.com/orca-zhang/cache/dist"
)

type GoRedisCli struct {
	redisCli *redis.Client
	chanSize int
}

// if the redis client is ready
func (g *GoRedisCli) OK() bool {
	_, err := g.redisCli.Ping().Result()
	return err == nil
}

// pub a key to channel
func (g *GoRedisCli) Pub(channel, key string) error {
	_, err := g.redisCli.Publish(channel, key).Result()
	return err
}

// sub a key from channel, callback uill tidy the local cache
func (g *GoRedisCli) Sub(channel string, callback func(payload string)) error {
	msgChan := g.redisCli.Subscribe(channel).ChannelSize(g.chanSize)

	for {
		msg, ok := <-msgChan
		if !ok {
			break
		}

		if msg != nil {
			callback(msg.Payload)
		}
	}
	return nil
}

func GoRedis(r *redis.Client, size ...int) dist.RedisCli {
	if r == nil {
		return nil
	}
	s := 100 // default 100 messages
	if len(size) > 0 {
		s = size[0]
	}
	return &GoRedisCli{
		redisCli: r,
		chanSize: s,
	}
}
