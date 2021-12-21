package dist

import (
	"github.com/go-redis/redis"
)

type GoRedisCli struct {
	redisCli *redis.Client
}

// if the redis client is ready
func (g *GoRedisCli) OK() bool {
	pong, err := g.redisCli.Ping().Result()
	if err == nil && pong == "PONG" {
		return true
	}
	return false
}

// pub a key to channel
func (g *GoRedisCli) Pub(channel, key string) error {
	_, err := g.redisCli.Publish(channel, key).Result()
	return err
}

// sub a key from channel, callback uill tidy the local cache
func (g *GoRedisCli) Sub(channel string, callback func(payload string)) error {
	msgChan := g.redisCli.Subscribe(topic).Channel()

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

func GoRedis(r *redis.Client) RedisCli {
	if r == nil {
		return nil
	}
	return &GoRedisCli{
		redisCli: r,
	}
}
