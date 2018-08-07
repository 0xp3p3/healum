# Todo Server

Todo server is a microservice to store todo.

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
	go get ./src/server/todo-srv && todo-srv -config ./src/server/config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222
	```

	OR as a docker container

	```shell
	docker run microhq/todo-srv -config config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222 --registry_address=YOUR_REGISTRY_ADDRESS
	```
## The API
Todo server implements the following RPC Methods

Todo
- Create
- Read
- Update
- Delete
- All
- Search
- ByCreator

### Todo.Create

```shell
curl -X POST \
  http://127.0.01:8080/server/todos/todo \
  -H 'content-type: application/json' \
  -d '{
    "todo": {
		"id":          "111",
		"name":        "todo1",
		"org_id":      "orgid",
		"creator_id":  "222"
	}
}'
{
  "data": {
   "todo": {
    "id": "111",
    "name": "todo1",
    "created": 1518094075,
    "updated": 1518094075,
    "creator_id": "222",
    "org_id": "orgid"
   }
  }
 }
```

### Todo.Read

```shell
curl -X GET \
  http://127.0.01:8080/server/todos/todo/111 \
  -H 'content-type: application/json'
{
  "data": {
   "todo": {
    "id": "111",
    "name": "todo1",
    "created": 1518094075,
    "updated": 1518094075,
    "creator_id": "222",
    "org_id": "orgid"
   }
  }
 }
```

### Todo.Update

```shell
curl -X PUT \
  http://localhost:8080/server/todos/todo \
  -H 'content-type: application/json' 
  -d '{
    "todo": {
        "id": "111",
		"name": "todo2",
		"org_id": "org2",
		"creator_id": "222"
    }
}'
{}
```

### Todo.Delete
```shell
curl -X DELETE \
  http://127.0.0.1:8080/server/todos/todo/111 \
  -H 'content-type: application/json'
{}
```

### Todo.All

```shell
curl -X GET \
  http://127.0.0.1:8080/server/todos/all \
  -H 'content-type: application/json'
{
  "data": {
   "todos": [
    {
     "id": "222",
     "name": "todo1",
     "created": 1518094135,
     "updated": 1518094135,
     "creator_id": "222",
     "org_id": "orgid"
    },
    {
     "id": "111",
     "name": "todo1",
     "created": 1518094075,
     "updated": 1518094075,
     "creator_id": "222",
     "org_id": "orgid"
    }
   ]
  }
 }
```

### Todo.Search


```shell
curl -X POST \
  http://127.0.0.1:8080/server/todos/search \
  -H 'content-type: application/json' \
  -d '{
	"name":"todo1"
}'
{
  "data": {
   "todos": [
    {
     "id": "222",
     "name": "todo1",
     "created": 1518094135,
     "updated": 1518094135,
     "creator_id": "222",
     "org_id": "orgid"
    },
    {
     "id": "111",
     "name": "todo1",
     "created": 1518094075,
     "updated": 1518094075,
     "creator_id": "222",
     "org_id": "orgid"
    }
   ]
  }
 }
```

### Todo.Filter

```shell
curl -X GET \
  http://127.0.0.1:8080/server/todos/creator/222 \
  -H 'content-type: application/json'
{
  "data": {
   "todos": [
    {
     "id": "222",
     "name": "todo1",
     "created": 1518094135,
     "updated": 1518094135,
     "creator_id": "222",
     "org_id": "orgid"
    },
    {
     "id": "111",
     "name": "todo1",
     "created": 1518094075,
     "updated": 1518094075,
     "creator_id": "222",
     "org_id": "orgid"
    }
   ]
  }
 }
```