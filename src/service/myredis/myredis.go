package myredis

import (
	"fmt"
	"gopkg.in/redis.v3"
)

var (
	clusterSize  uint32
	client       *redis.Client
	redisClients map[uint32]*redis.Client
)

func Init(url string) {
	client = redis.NewClient(&redis.Options{
		Addr:     url,
		DB:       0,
		PoolSize: 100,
	})
}

func InitCluster(urls map[uint32]string) {
	tmp := len(urls)
	clusterSize = (uint32)(tmp)
	redisClients = make(map[uint32]*redis.Client)
	for index, url := range urls {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     url,
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

func ClusterClient(index uint32) *redis.Client {
	fmt.Println(index, "======", index%clusterSize)
	return redisClients[index%clusterSize]
}
