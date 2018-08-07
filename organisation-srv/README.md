# Organisation Server

Organisation server is a microservice to store organisation.

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
	go get ./src/server/organisation-srv && organisation-srv -config ./src/server/config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222
	```

	OR as a docker container

	```shell
	docker run microhq/organisation-srv -config config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222 --registry_address=YOUR_REGISTRY_ADDRESS
	```

5. Generate API Document
    ```shell
    cd ./server/api/api
    apidoc -i api/ -o doc/ -t template/
    ```

## The API
Organisation server implements the following RPC Methods

Organisation
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

### Organisation.Create

```shell
curl -X POST \
  http://127.0.01:8080/server/organisations/organisation \
  -H 'content-type: application/json' \
  -d '{
    "organisation": {
		"id":          "333",
		"title":       "organisation1",
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
        "organisation": {
            "id": "111",
            "title": "organisation1",
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

### Organisation.Read

```shell
curl -X GET \
  http://127.0.01:8080/server/organisations/organisation/111 \
  -H 'content-type: application/json'
{
    "data": {
        "organisation": {
            "id": "111",
            "title": "organisation1",
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

### Organisation.Update

```shell
curl -X PUT \
  http://localhost:8080/server/organisations/organisation \
  -H 'content-type: application/json' 
  -d '{
    "organisation": {
        "id":          "111",
		"title":       "organisation2",
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

### Organisation.Delete
```shell
curl -X DELETE \
  http://127.0.0.1:8080/server/organisations/organisation/111 \
  -H 'content-type: application/json'
{}
```

### Organisation.All

```shell
curl -X GET \
  http://127.0.0.1:8080/server/organisations/all \
  -H 'content-type: application/json'
{
    "data": {
        "organisations": [
            {
                "id": "111",
                "title": "organisation1",
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

### Organisation.Search


```shell
curl -X POST \
  http://127.0.0.1:8080/server/organisations/search \
  -H 'content-type: application/json' \
  -d '{
	"name":"organisation1"
}'
{
    "data": {
        "organisations": [
            {
                "id": "111",
                "title": "organisation1",
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

### Organisation.ByCreator
```shell
curl -X GET \
  http://127.0.0.1:8080/server/organisations/creator/222 \
  -H 'content-type: application/json'
{
    "data": {
        "organisations": [
            {
                "id": "111",
                "title": "organisation1",
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

### Organisation.ByUser
```shell
curl -X GET \
  http://127.0.0.1:8080/server/organisations/user/333 \
  -H 'content-type: application/json'
{
    "data": {
        "organisations": [
            {
                "id": "111",
                "title": "organisation1",
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

### Organisation.Filter

```shell
curl -X POST \
  http://127.0.0.1:8080/server/organisations/filter \
  -H 'content-type: application/json' \
  -d '{
    "category": ["category1", "category2"],
    "tags": ["a", "c", "d"]
}'
{
    "data": {
        "organisations": [
            {
                "id": "111",
                "title": "organisation1",
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