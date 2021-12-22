package goredis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/orca-zhang/cache/dist"
)

type GoRedisCli struct {
	ctx      context.Context
	redisCli *redis.Client
	chanSize int
}

// if the redis client is ready
func (g *GoRedisCli) OK() bool {
	_, err := g.redisCli.Ping(g.ctx).Result()
	return err == nil
}

// pub a payload to channel
func (g *GoRedisCli) Pub(channel, payload string) error {
	_, err := g.redisCli.Publish(g.ctx, channel, payload).Result()
	return err
}

// sub a payload from channel, callback uill tidy the local cache
func (g *GoRedisCli) Sub(channel string, callback func(payload string)) error {
	msgChan := g.redisCli.Subscribe(g.ctx, channel).ChannelSize(g.chanSize)
	for {
		if msg, _ := <-msgChan; msg != nil {
			callback(msg.Payload)
		} else {
			break
		}
	}
	return nil
}

func Take(r *redis.Client, size ...int) dist.RedisCli {
	s := 100 // default 100 messages
	if len(size) > 0 {
		s = size[0]
	}
	return &GoRedisCli{
		ctx:      context.TODO(),
		redisCli: r,
		chanSize: s,
	}
}
