// Package discovery is an interface for scalable service discovery.
package discovery

import (
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-os/discovery"
)

func init() {
	cmd.DefaultRegistries["os"] = NewRegistry
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	return discovery.NewRegistry(opts...)
}
