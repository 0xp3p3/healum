# Track Server

Track server is a microservice to store track.

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
	go get ./src/server/track-srv && track-srv -config ./src/server/config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222
	```

	OR as a docker container

	```shell
	docker run microhq/track-srv -config config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222 --registry_address=YOUR_REGISTRY_ADDRESS
	```
## The API
Track server implements the following RPC Methods

Track
- Create
- Read
- Update
- Delete
- All
- Search
- ByCreator

### Track.Create

```shell
curl -X POST \
  http://127.0.01:8080/server/tracks/track \
  -H 'content-type: application/json' \
  -d '{
    "track": {
		"id":          "111",
		"name":        "track1",
		"org_id":      "orgid",
		"creator_id":  "222"
	}
}'
{
  "data": {
   "track": {
    "id": "111",
    "name": "track1",
    "created": 1518094075,
    "updated": 1518094075,
    "creator_id": "222",
    "org_id": "orgid"
   }
  }
 }
```

### Track.Read

```shell
curl -X GET \
  http://127.0.01:8080/server/tracks/track/111 \
  -H 'content-type: application/json'
{
  "data": {
   "track": {
    "id": "111",
    "name": "track1",
    "created": 1518094075,
    "updated": 1518094075,
    "creator_id": "222",
    "org_id": "orgid"
   }
  }
 }
```

### Track.Update

```shell
curl -X PUT \
  http://localhost:8080/server/tracks/track \
  -H 'content-type: application/json' 
  -d '{
    "track": {
        "id": "111",
		"name": "track2",
		"org_id": "org2",
		"creator_id": "222"
    }
}'
{}
```

### Track.Delete
```shell
curl -X DELETE \
  http://127.0.0.1:8080/server/tracks/track/111 \
  -H 'content-type: application/json'
{}
```

### Track.All

```shell
curl -X GET \
  http://127.0.0.1:8080/server/tracks/all \
  -H 'content-type: application/json'
{
  "data": {
   "tracks": [
    {
     "id": "222",
     "name": "track1",
     "created": 1518094135,
     "updated": 1518094135,
     "creator_id": "222",
     "org_id": "orgid"
    },
    {
     "id": "111",
     "name": "track1",
     "created": 1518094075,
     "updated": 1518094075,
     "creator_id": "222",
     "org_id": "orgid"
    }
   ]
  }
 }
```

### Track.Search


```shell
curl -X POST \
  http://127.0.0.1:8080/server/tracks/search \
  -H 'content-type: application/json' \
  -d '{
	"name":"track1"
}'
{
  "data": {
   "tracks": [
    {
     "id": "222",
     "name": "track1",
     "created": 1518094135,
     "updated": 1518094135,
     "creator_id": "222",
     "org_id": "orgid"
    },
    {
     "id": "111",
     "name": "track1",
     "created": 1518094075,
     "updated": 1518094075,
     "creator_id": "222",
     "org_id": "orgid"
    }
   ]
  }
 }
```

### Track.Filter

```shell
curl -X GET \
  http://127.0.0.1:8080/server/tracks/creator/222 \
  -H 'content-type: application/json'
{
  "data": {
   "tracks": [
    {
     "id": "222",
     "name": "track1",
     "created": 1518094135,
     "updated": 1518094135,
     "creator_id": "222",
     "org_id": "orgid"
    },
    {
     "id": "111",
     "name": "track1",
     "created": 1518094075,
     "updated": 1518094075,
     "creator_id": "222",
     "org_id": "orgid"
    }
   ]
  }
 }
```