# File Server

Activity service that returns driver related information with different backend 
e.g. 
	https://www.better.org.uk/odi/sessions.json
	https://fusionfunctions.azurewebsites.net/api/FusOpenFunction
	https://app.opensessions.io/api/rdpe/sessions
	https://ourparks.org.uk/getSessions

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Start a Arango db

4. Download and start the service

	```shell
	go get github.com/micro/activity-srv
	activity-srv"

	go get ./src/server/activity-srv && activity-srv -config ./src/server/activity-srv/config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222
	```

	OR as a docker container

	```shell
	docker run microhq/activity-srv --registry_address=YOUR_REGISTRY_ADDRESS
	```

## The API
Activity server implements the following RPC Methods

Data
- Query
