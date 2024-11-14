package predis

import (
	"github.com/go-puzzles/puzzles/dialer"
	"github.com/gomodule/redigo/redis"
)

type RedisConf struct {
	Server   string `desc:"redis server name (default localhost:6379)"`
	Password string `desc:"redis server password"`
	Db       int    `desc:"redis db (default 0)"`
	MaxIdle  int    `desc:"redis maxIdle (default 100)"`
}

func (conf *RedisConf) DialRedisPool() *redis.Pool {
	return dialer.DialRedisPool(
		conf.Server,
		conf.Db,
		conf.MaxIdle,
		conf.Password,
	)
}

func (rc *RedisConf) SetDefault() {
	if rc.Server == "" && rc.MaxIdle == 0 {
		rc.Server = "localhost:6379"
		rc.MaxIdle = 100
	}
}
