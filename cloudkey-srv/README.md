# Cloud Key Server

Key management service cloudkey-srv which is backed by Google Cloud Key Management Service (KMS).

## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Start a mysql database

4. Download and start the service

	```shell
	go get github.com/micro/cloudkey-srv
	cloudkey-srv"
	```

	OR as a docker container

	```shell
	docker run microhq/cloudkey-srv --registry_address=YOUR_REGISTRY_ADDRESS
	```

## The API
Profile server implements the following RPC Methods

Record
- Read

### CloudKey.Read
```shell
micro query healum.srv.file CloudKey.Read '{
		"prgid": "1", 
	}
}' 
'{
	"filerecord": {
		"id": "1", 
		"userid": "user1",    
		"url": "http://example.com", 
	}
}' 
```
