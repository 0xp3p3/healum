# Behaviour Server

Behaviour server is a microservice to store kv.

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Start a AarangoDB

4. Download and start the service

	```shell
	go get ./src/server/kv-srv && kv-srv -config ./src/server/config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222
	```

	OR as a docker container

	```shell
	docker run microhq/kv-srv -config config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222 --registry_address=YOUR_REGISTRY_ADDRESS
	```

5. Install Redis
	https://redis.io/topics/quickstart
	https://medium.com/@petehouston/install-and-config-redis-on-mac-os-x-via-homebrew-eb8df9a4f298

6. Requriemnt
	Install Redis
	Set server address and password in config.json of kv-srv
	```shell
	{
		"service": {
			"version": "0.0.1",
				"name": "go.micro.srv.kv",
			"description": "Task server is a microservice to store kv",
			"server": "127.0.0.1:6379",
			"password": ""
		}
	}
	```

7. Execution
	After all sets, check TestPut function in handler_test.go

	Open Terminal and run below commands to check test function

	```shell
	> redis-cli
	127.0.0.1:6379 > get key
	"hello world!"
	```	