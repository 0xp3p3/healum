# Go Server


## What's here?

- **db-srv** - database microservice. A proxy microservice to object storages backends (MySQL, InfluxDB, Redis, etc.)
- **cloudkey-srv** - key management service cloudkey-srv which is backed by Google Cloud Key Management Service (KMS).
- **organisation-srv** - service to store user organizations. Multi tenancy support
- **user-srv** - user and authentication microservice. User objects, sessions, etc.
- **api** - an external restful API.


### Prerequisites

Set GOPATH

Install Consul
[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

Run Consul
```
$ consul agent -dev -advertise=127.0.0.1
```

Install go-micro
[Quick start](https://github.com/micro/micro)

Install and run [NATS](https://github.com/nats-io/gnatsd)
```
$ go get github.com/nats-io/gnatsd
$ gnatsd
```

### How to run
Please, refer to a README section of each microservice for more details. 
Basically, it requires building and running each microservice with a few special command for some of the microservices.



or with Docker
```
$ docker-compose build
$ docker-compose up

# to stop it
$ docker-compose down
```

Copy-paste instruction for running a development configuration is available further in the readme.

CURRENT LIMITATION (OS X). Before you run Docker compose you have to set your docker machine IP.
```
$ docker-machine ip
``` 
If this IP is different from 127.0.0.1 you should change all the occurrences of 127.0.0.1 to you IP in files "docker-compose.yml", "db-srv/Dockerfile.db" and run docker-compose

### How to test
1. Run microservices as described in readme files
2. Go to ./src/server of you go path
3. Test services

```
$ go test ./...
```


### DevOps instructions

#### Update Docker images
1. Get access to healum Docker hub and login
```shell
$ docker login
```
2. Rebuild the Docker image with the new version (+1 to the old version)
```shell
$ docker build --pull -f Dockerfile.parent -t "healum/server:v2" .
```
3. Push the image
```shell
$ docker push healum/server:v2
```
4. Update the version on all the images of the container

#### Build micro with NATS

1. In $GOPATH/src/github.com/micro/micro/ create plugins.go
2. Add the following content ot plugins.go
```go
package main

import (
	// nats transport
	_ "github.com/micro/go-plugins/transport/nats"

	// nats broker
	_ "github.com/micro/go-plugins/broker/nats"
)
```
3. Rebuild micro
```shell
$ go build -i -o $GOPATH/bin/micro $GOPATH/src/github.com/micro/micro/main.go $GOPATH/src/github.com/micro/micro/plugins.go
```

#### Restart consul

```
$ consul agent -dev -advertise=127.0.0.1
$ ps -A | grep -m1 consul | awk '{print $1}' 
$ sudo kill -9 consul_id
$ consul agent -dev -advertise=127.0.0.1
```

#### Update dependencies
The vendoring package is committed with the repo. 
To update dependencies for each microservice (with installed glide):
```
remove vendor folder content
$ glide cache-clear
$ glide up -- quick
$ glide install --force
commit dependancies
``` 
Commit the vendor folder

#### Run protobuf

```
$ brew update && brew install --devel protobuf
$ go get -u github.com/micro/protobuf/{proto,protoc-gen-go}
$ protoc -I$GOPATH/src --go_out=plugins=micro:$GOPATH/src $GOPATH/src/server/user-srv/proto/user/user.proto  
```

#### Configure shell for Docker

```
$ docker-machine start
$ eval "$(docker-machine env default)" 
```

if there are problems with OS X Docker VM

```
$ docker-machine restart              # Restart the environment
$ eval $(docker-machine env default)  # Refresh your environment settings
```


#### Drop MySQL database

```
$ mysqladmin drop test_db --user=root
```

#### GOPATH
```
export GOPATH=$HOME
export PATH=$PATH:$HOME/bin
```

#### Run the API
```
# consul
consul agent -dev -advertise=127.0.0.1

# nats
gnatsd -DV

# db-srv
micro register service '{"name": "go.micro.db.mysql", "version": "0.0.1", "nodes": [{"id": "kv-1", "address": "127.0.0.1", "port": 3306, "metadata": {"driver": "mysql"}}]}'
micro register service '{"name": "go.micro.db.elasticsearch", "version": "0.0.1", "nodes": [{"id": "kv-2", "address": "127.0.0.1", "port": 9200, "metadata": {"driver": "elasticsearch"}}]}'
micro register service '{"name": "go.micro.db.redis", "version": "0.0.1", "nodes": [{"id": "kv-3", "address": "127.0.0.1", "port": 6379, "metadata": {"driver": "redis"}}]}'
micro register service '{"name": "go.micro.db.arangodb", "version": "0.0.1", "nodes": [{"id": "kv-4", "address": "127.0.0.1", "port": 8529, "metadata": {"driver": "arangodb"}}]}'
micro register service '{"name": "go.micro.db.influxdb", "version": "0.0.1", "nodes": [{"id": "kv-5", "address": "127.0.0.1", "port": 8086, "metadata": {"driver": "influxdb"}}]}'
go get server/db-srv && db-srv -config ./src/server/db-srv/config.json --database_service_namespace=go.micro.db --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222

# user-srv
go get server/user-srv && user-srv -config ./src/server/user-srv/config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222

# cloudkey-srv
go get server/cloudkey-srv && cloudkey-srv -config ./src/server/cloudkey-srv/config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222

# organisation-srv
go get server/organisation-srv && organisation-srv -config ./src/server/organisation-srv/config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222

# api
micro --api_handler=proxy api
go get server/api && api -config ./src/server/api/config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222

# metrics
influxd
telegraf -config ./src/server/metrics/telegraf.conf
```

## The API Demo
The purpose of the demo is a demonstration of the external API usage with underlying microservices interactions. During the demo newly created users will exchange messages in a room. 

Variables used here:

**API_ADDRESS** address of the API (example: 127.0.0.1 for OS X, 127.0.0.1 for Linux). 127.0.0.1 is used further

**FIRST_USER_SESSION_ID** - session ID of the first user on login (example: erLjzlWo1C9h9x3QnSScaZSX4YVNLXkb)

**SECOND_USER_SESSION_ID** - session ID of the second user on login (example: HD7hP3NRfQ8AFf1FpBUJ7Pcw7pRnOwFH)

**ROOM_ID** - ID of a room created (example: ZBzKPXMYm8FNFNQ5uFIh7uXQPvd7T8vl)


### Create two users
Create the first user
```
echo '{"user": {"email": "email@email.com", "orgid": "orgid"}, "password": "pass"}' | curl -d @- -H "Content-Type: application/json" http://127.0.0.1:8080/server/user/create
{
  "user": {
   "id": "1653185779643901076",
   "email": "email@email.com",
   "orgid": "orgid",
   "created": 1504688339,
   "updated": 1504688339
  }
}
```

If the user already exists, just create a user with another email.

Create the second user
```
echo '{"user": {"email": "email1@email.com", "orgid": "orgid"}, "password": "pass"}' | curl -d @- -H "Content-Type: application/json" http://127.0.0.1:8080/server/user/create
{
  "user": {
   "id": "3610389916378849322",
   "email": "email1@email.com",
   "created": 1469375164,
   "updated": 1469375164
  }
 }
```

### Authenticate users 
When users authenticated they get a session id and can perform operations on a server. The session id must be passed to queries so the server can identify authenticated users. There are a few authentication methods. The following method uses email and password.
Login for the first user. The "id" field contains FIRST_USER_SESSION_ID (without quotes):
```
echo '{"email": "email@email.com", "password": "pass", "orgid": "orgid"}' | curl -d @- -H "Content-Type: application/json" http://127.0.0.1:8080/server/user/login
{
  "session": {
   "id": "MmkqgyS7lHMjhLJ8K03TfyOea5VURv5Q",
   "email": "email@email.com",
   "orgid": "orgid",
   "created": 1504688398,
   "expires": 1505293198
  }
} 
```
Login for the second user. The "id" field contains SECOND_USER_SESSION_ID (without quotes)::
```
echo '{"email": "email1@email.com", "password": "pass"}' | curl -d @- -H "Content-Type: application/json" http://127.0.0.1:8080/server/user/login
{
  "session": {
   "id": "HD7hP3NRfQ8AFf1FpBUJ7Pcw7pRnOwFH",
   "email": "email1@email.com",
   "created": 1469375193,
   "expires": 1469979993
  }
 }
```
Save the session IDs of the users, because it is required for further commands. A new session id is generated everytime the user os logged in, old session ID is removed in that case.
