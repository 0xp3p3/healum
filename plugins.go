package main

import (
	// nats transport
	_ "github.com/micro/go-plugins/transport/nats"

	// nats broker
	_ "github.com/micro/go-plugins/broker/nats"
)