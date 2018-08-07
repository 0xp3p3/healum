# Note Server

Note server is a microservice to store note.

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
	go get ./src/server/note-srv && note-srv -config ./src/server/config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222
	```

	OR as a docker container

	```shell
	docker run microhq/note-srv -config config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222 --registry_address=YOUR_REGISTRY_ADDRESS
	```

5. Generate API Document
    ```shell
    cd ./server/api/api
    apidoc -i api/ -o doc/ -t template/
    ```

## The API
Note server implements the following RPC Methods

Note
- Create
- Read
- Update
- Delete
- All
- Search
- Templates
- Drafts
- ByCreator
- ByUser
- Filter

### Note.Create

```shell
curl -X POST \
  http://127.0.01:8080/server/notes/note \
  -H 'content-type: application/json' \
  -d '{
    "note": {
		"id":          "333",
		"title":       "note1",
		"org_id":      "orgid",
		"description": "description1",
		"creator_id":  "222",
		"user_id": 	   "333",
		"category":	   "category1",
		"tags":        ["a","b","c"]
	}
}'
{
    "data": {
        "note": {
            "id": "111",
            "title": "note1",
            "description": "description1",
            "created": 1518097899,
            "updated": 1518097899,
            "user_id": "333",
            "creator_id": "222",
            "category": "category1",
            "tags": [
                "a",
                "b",
                "c"
            ],
            "org_id": "orgid"
        }
    }
}
```

### Note.Read

```shell
curl -X GET \
  http://127.0.01:8080/server/notes/note/111 \
  -H 'content-type: application/json'
{
    "data": {
        "note": {
            "id": "111",
            "title": "note1",
            "description": "description1",
            "created": 1518037481,
            "updated": 1518037481,
            "user_id": "333",
            "creator_id": "222",
            "category": "category1",
            "tags": [
                "a",
                "b",
                "c"
            ],
            "org_id": "orgid"
        }
    }
}
```

### Note.Update

```shell
curl -X PUT \
  http://localhost:8080/server/notes/note \
  -H 'content-type: application/json' 
  -d '{
    "note": {
        "id":          "111",
		"title":       "note2",
		"org_id":      "org2",
		"description": "description1",
		"creator_id":  "222",
		"user_id": 	   "333",
		"category":	   "category1",
		"tags":        ["a","b","c"]
    }
}'
{}
```

### Note.Delete
```shell
curl -X DELETE \
  http://127.0.0.1:8080/server/notes/note/111 \
  -H 'content-type: application/json'
{}
```

### Note.All

```shell
curl -X GET \
  http://127.0.0.1:8080/server/notes/all \
  -H 'content-type: application/json'
{
    "data": {
        "notes": [
            {
                "id": "111",
                "title": "note1",
                "description": "description1",
                "created": 1518037481,
                "updated": 1518037481,
                "user_id": "333",
                "creator_id": "222",
                "category": "category1",
                "tags": [
                    "a",
                    "b",
                    "c"
                ],
                "org_id": "orgid"
            }
        ]
    }
}
```

### Note.Search


```shell
curl -X POST \
  http://127.0.0.1:8080/server/notes/search \
  -H 'content-type: application/json' \
  -d '{
	"name":"note1"
}'
{
    "data": {
        "notes": [
            {
                "id": "111",
                "title": "note1",
                "description": "description1",
                "created": 1518037481,
                "updated": 1518037481,
                "user_id": "333",
                "creator_id": "222",
                "category": "category1",
                "tags": [
                    "a",
                    "b",
                    "c"
                ],
                "org_id": "orgid"
            }
        ]
    }
}
```

### Note.ByCreator
```shell
curl -X GET \
  http://127.0.0.1:8080/server/notes/creator/222 \
  -H 'content-type: application/json'
{
    "data": {
        "notes": [
            {
                "id": "111",
                "title": "note1",
                "description": "description1",
                "created": 1518037481,
                "updated": 1518037481,
                "user_id": "333",
                "creator_id": "222",
                "category": "category1",
                "tags": [
                    "a",
                    "b",
                    "c"
                ],
                "org_id": "orgid"
            }
        ]
    }
}
```

### Note.ByUser
```shell
curl -X GET \
  http://127.0.0.1:8080/server/notes/user/333 \
  -H 'content-type: application/json'
{
    "data": {
        "notes": [
            {
                "id": "111",
                "title": "note1",
                "description": "description1",
                "created": 1518037481,
                "updated": 1518037481,
                "user_id": "333",
                "creator_id": "222",
                "category": "category1",
                "tags": [
                    "a",
                    "b",
                    "c"
                ],
                "org_id": "orgid"
            }
        ]
    }
}
```

### Note.Filter

```shell
curl -X POST \
  http://127.0.0.1:8080/server/notes/filter \
  -H 'content-type: application/json' \
  -d '{
    "category": ["category1", "category2"],
    "tags": ["a", "c", "d"]
}'
{
    "data": {
        "notes": [
            {
                "id": "111",
                "title": "note1",
                "description": "description1",
                "created": 1518037481,
                "updated": 1518037481,
                "user_id": "333",
                "creator_id": "222",
                "category": "category1",
                "tags": [
                    "a",
                    "b",
                    "c"
                ],
                "org_id": "orgid"
            }
        ]
    }
}
```