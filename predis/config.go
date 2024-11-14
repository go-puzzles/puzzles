package predis

import (
	redis2 "github.com/go-puzzles/puzzles/dialer/redis"
	"github.com/gomodule/redigo/redis"
)

type RedisConf struct {
	Server   string `desc:"redis server name (default localhost:6379)"`
	Password string `desc:"redis server password"`
	Db       int    `desc:"redis db (default 0)"`
	MaxIdle  int    `desc:"redis maxIdle (default 100)"`
}

func (conf *RedisConf) DialRedisPool() *redis.Pool {
	return redis2.DialRedisPool(
		conf.Server,
		conf.Db,
		conf.MaxIdle,
		conf.Password,
	)
}

func (conf *RedisConf) SetDefault() {
	if conf.Server == "" && conf.MaxIdle == 0 {
		conf.Server = "localhost:6379"
		conf.MaxIdle = 100
	}
}
