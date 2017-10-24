package redis

import (
	"github.com/go-redis/redis"

	"github.com/prizem-io/gateway/config"
)

func Connect(config config.Configuration) *redis.Client {
	var options redis.Options
	config.UnmarshalKey("redis", &options)

	return redis.NewClient(&options)
}

func ConnectCluster(config config.Configuration) *redis.ClusterClient {
	var options redis.ClusterOptions
	config.UnmarshalKey("redisCluster", &options)

	return redis.NewClusterClient(&options)
}
