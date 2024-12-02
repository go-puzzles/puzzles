package goredis

type RedisConf struct {
	Server   string `desc:"redis server name (default localhost:6379)"`
	Db       int    `desc:"redis db (default 0)"`
	Username string `desc:"redis username"`
	Password string `desc:"redis server password"`
}

func (conf *RedisConf) DialRedisPool() *PuzzleRedisClient {
	return NewRedisClientWithAuth(conf.Server, conf.Db, conf.Username, conf.Password)
}

func (conf *RedisConf) SetDefault() {
	if conf.Server == "" {
		conf.Server = "localhost:6379"
	}

	if conf.Db == 0 {
		conf.Db = 0
	}
}
