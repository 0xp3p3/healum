package utils

import (
	"github.com/juju/ratelimit"
	"github.com/micro/go-micro/client"
	micro_ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit"
)

const (
	// default ratelimit
	rateLimit = 100
)

// Returns micro RPC client with a rate limit
func NewRateLimitedClient(cl client.Client) client.Client {
	b := ratelimit.NewBucketWithRate(float64(rateLimit), int64(rateLimit))
	wrapper := micro_ratelimit.NewClientWrapper(b, false)
	return wrapper(cl)
}
