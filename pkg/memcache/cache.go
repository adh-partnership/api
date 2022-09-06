package memcache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

var Cache *ttlcache.Cache[string, interface{}]

func init() {
	Cache = ttlcache.New(
		ttlcache.WithTTL[string, interface{}](10 * time.Minute),
	)

	go Cache.Start()
}
