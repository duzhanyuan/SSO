package myredis

import (
	"fmt"
	"gopkg.in/redis.v3"
	"service/consistent"
)

var (
	client       *redis.Client
	redisClients map[string]*redis.Client //virtualnode => client
	hasher       *consistent.Consistent
	nodesMapping map[string]string
)

func Init(url string) {
	client = redis.NewClient(&redis.Options{
		Addr:     url,
		DB:       0,
		PoolSize: 100,
	})
}

func InitCluster(addrs map[string]int64, nodes map[string]string) {
	redisClients = make(map[string]*redis.Client)
	for virtualNode, realNode := range nodes {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     realNode,
			DB:       addrs[realNode],
			PoolSize: 100,
		})
		if _, err := redisClient.Ping().Result(); err != nil {
			panic(err)
		}
		redisClients[virtualNode] = redisClient
	}
	hasher = consistent.New()
	for virtualNode := range nodes {
		hasher.Add(virtualNode)
	}
	nodesMapping = nodes
}

func Client() *redis.Client {
	return client
}

func ClusterClient(key string) *redis.Client {
	addr, err := hasher.Get(key)
	if err != nil {
		panic("redis handle nil")
	}
	fmt.Println(addr, "======", "key")
	handle, ok := redisClients[addr]
	if !ok {
		panic("redis handle nil")
	}
	return handle
}
