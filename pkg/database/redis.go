/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

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
