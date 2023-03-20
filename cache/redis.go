package cache

import (
	"github.com/go-redis/redis"
)

type Config struct {
	Addr string `yaml:"Addr"`
}

var RedisClient *redis.Client

// 根据redis配置初始化一个客户端
func Init(redisConf Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisConf.Addr, // redis地址
	})

	//通过 *redis.Client.Ping() 来检查是否成功连接到了redis服务器
	if _, err := RedisClient.Ping().Result(); err != nil {
		return err
	}
	return nil
}
