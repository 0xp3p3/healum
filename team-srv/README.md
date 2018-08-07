# Team Server

Team server is a microservice to store team.

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
	go get ./src/server/team-srv && team-srv -config ./src/server/config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222
	```

	OR as a docker container

	```shell
	docker run microhq/team-srv -config config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222 --registry_address=YOUR_REGISTRY_ADDRESS
	```
## The API
Team server implements the following RPC Methods

Team
- Create
- Read
- Update
- Delete
- All
- Search
- ByCreator

### Team.Create

```shell
curl -X POST \
  http://127.0.01:8080/server/teams/team \
  -H 'content-type: application/json' \
  -d '{
    "team": {
		"id":          "111",
		"name":        "team1",
		"org_id":      "orgid",
		"creator_id":  "222"
	}
}'
{
  "data": {
   "team": {
    "id": "111",
    "name": "team1",
    "created": 1518094075,
    "updated": 1518094075,
    "creator_id": "222",
    "org_id": "orgid"
   }
  }
 }
```

### Team.Read

```shell
curl -X GET \
  http://127.0.01:8080/server/teams/team/111 \
  -H 'content-type: application/json'
{
  "data": {
   "team": {
    "id": "111",
    "name": "team1",
    "created": 1518094075,
    "updated": 1518094075,
    "creator_id": "222",
    "org_id": "orgid"
   }
  }
 }
```

### Team.Update

```shell
curl -X PUT \
  http://localhost:8080/server/teams/team \
  -H 'content-type: application/json' 
  -d '{
    "team": {
        "id": "111",
		"name": "team2",
		"org_id": "org2",
		"creator_id": "222"
    }
}'
{}
```

### Team.Delete
```shell
curl -X DELETE \
  http://127.0.0.1:8080/server/teams/team/111 \
  -H 'content-type: application/json'
{}
```

### Team.All

```shell
curl -X GET \
  http://127.0.0.1:8080/server/teams/all \
  -H 'content-type: application/json'
{
  "data": {
   "teams": [
    {
     "id": "222",
     "name": "team1",
     "created": 1518094135,
     "updated": 1518094135,
     "creator_id": "222",
     "org_id": "orgid"
    },
    {
     "id": "111",
     "name": "team1",
     "created": 1518094075,
     "updated": 1518094075,
     "creator_id": "222",
     "org_id": "orgid"
    }
   ]
  }
 }
```

### Team.Search


```shell
curl -X POST \
  http://127.0.0.1:8080/server/teams/search \
  -H 'content-type: application/json' \
  -d '{
	"name":"team1"
}'
{
  "data": {
   "teams": [
    {
     "id": "222",
     "name": "team1",
     "created": 1518094135,
     "updated": 1518094135,
     "creator_id": "222",
     "org_id": "orgid"
    },
    {
     "id": "111",
     "name": "team1",
     "created": 1518094075,
     "updated": 1518094075,
     "creator_id": "222",
     "org_id": "orgid"
    }
   ]
  }
 }
```

### Team.Filter

```shell
curl -X GET \
  http://127.0.0.1:8080/server/teams/creator/222 \
  -H 'content-type: application/json'
{
  "data": {
   "teams": [
    {
     "id": "222",
     "name": "team1",
     "created": 1518094135,
     "updated": 1518094135,
     "creator_id": "222",
     "org_id": "orgid"
    },
    {
     "id": "111",
     "name": "team1",
     "created": 1518094075,
     "updated": 1518094075,
     "creator_id": "222",
     "org_id": "orgid"
    }
   ]
  }
 }
```