package database

import "github.com/go-redis/redis/v8"

var Redis *redis.Client

type RedisOptions struct {
	// Sentinel config
	Sentinel      bool
	MasterName    string
	SentinelAddrs []string

	// Single Redis config
	Addr string

	// Common
	Password string
	DB       int
}

func ConnectRedis(options RedisOptions) {
	if options.Sentinel {
		Redis = redis.NewFailoverClient(buildSentinelOptions(options))
	} else {
		Redis = redis.NewClient(buildRedisOptions(options))
	}
}

func buildSentinelOptions(options RedisOptions) *redis.FailoverOptions {
	return &redis.FailoverOptions{
		MasterName:    options.MasterName,
		SentinelAddrs: options.SentinelAddrs,
		DB:            options.DB,
		Password:      options.Password,
	}
}

func buildRedisOptions(options RedisOptions) *redis.Options {
	return &redis.Options{
		Addr:     options.Addr,
		DB:       options.DB,
		Password: options.Password,
	}
}
