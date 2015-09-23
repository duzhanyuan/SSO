package myredis

import (
	"gopkg.in/redis.v3"
)

var (
	client       *redis.Client
	redisClients map[int64]*redis.Client
)

func Init(url string) {
	client = redis.NewClient(&redis.Options{
		Addr:     url,
		DB:       0,
		PoolSize: 100,
	})
}

func InitCluster(urls map[int64]string) {
	redisClients = make(map[string]*redis.Client)
	for index, url := range urls {
		redisClient := redis.NewClient(&redis.Options{
			addr:     url,
			DB:       0,
			PoolSize: 100,
		})
		if _, err := redisClient.Ping().Result(); err != nil {
			panic(err)
		}
		redisClients[index] = redisClient
	}

}

func Client() *redis.Client {
	return client
}

func ClusterClient(index int64) *redis.Client {
	return redisClients[index]
}
