package utils

import (
	"github.com/micro/go-micro/client"
	micro_circuitbreaker "github.com/micro/go-plugins/wrapper/breaker/hystrix"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-os/config"
)

// Returns micro RPC client with a circuit breaker
func NewCircuitBreakerClient(cl client.Client, c config.Config) client.Client {
	// settings
	hystrix.DefaultTimeout = c.Get("hystrix", "DefaultTimeout").Int(hystrix.DefaultTimeout)
	hystrix.DefaultMaxConcurrent = c.Get("hystrix", "DefaultMaxConcurrent").Int(hystrix.DefaultMaxConcurrent)
	hystrix.DefaultVolumeThreshold = c.Get("hystrix", "DefaultVolumeThreshold").Int(hystrix.DefaultVolumeThreshold)
	hystrix.DefaultSleepWindow = c.Get("hystrix", "DefaultSleepWindow").Int(hystrix.DefaultSleepWindow)
	hystrix.DefaultErrorPercentThreshold = c.Get("hystrix", "DefaultErrorPercentThreshold").Int(hystrix.DefaultErrorPercentThreshold)

	wrapper := micro_circuitbreaker.NewClientWrapper()
	return wrapper(cl)
}
