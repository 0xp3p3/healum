# Plan Server

Plan server is a microservice to store plan.

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
	go get ./src/server/plan-srv && plan-srv -config ./src/server/config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222
	```

	OR as a docker container

	```shell
	docker run microhq/plan-srv -config config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222 --registry_address=YOUR_REGISTRY_ADDRESS
	```
## The API
Plan server implements the following RPC Methods

Plan
- Create
- Read
- Update
- Delete
- All
- Search
- Templates
- Drafts
- ByCreator
- Filters
- TimeFilters
- GoalFilters

### Plan.Create

```shell
curl -X POST \
  http://127.0.0.1:8080/server/plans/plan \
  -H 'content-type: application/json' \
  -d '{
    "plan": {
		"id":          "222",
		"name":        "plan1",
		"org_id":      "orgid",
		"description": "description1",
		"creator_id":  "userid",
		"templateId":  "template1",
		"isTemplate":  true,
		"status":      1,
		"goals": [{"id": "1"}, {"id": "2"}]
	}
}'
{
    "data": {
        "plan": {
            "id": "111",
            "name": "plan1",
            "org_id": "orgid",
            "description": "description1",
            "created": 1518138048,
            "updated": 1518138048,
            "goals": [
                {
                    "id": "1"
                },
                {
                    "id": "2"
                }
            ],
            "creator_id": "222",
            "isTemplate": true,
            "templateId": "template1",
            "status": 1
        }
    }
}
```

### Plan.Read

```shell
curl -X GET \
  http://127.0.01:8080/server/plans/plan/111 \
  -H 'content-type: application/json'
{
    "data": {
        "plan": {
            "id": "111",
            "name": "plan1",
            "org_id": "orgid",
            "description": "description1",
            "created": 1518138048,
            "updated": 1518138048,
            "goals": [
                {
                    "id": "1"
                },
                {
                    "id": "2"
                }
            ],
            "creator_id": "222",
            "isTemplate": true,
            "templateId": "template1",
            "status": 1
        }
    }
}
```

### Plan.Update

```shell
curl -X PUT \
  http://localhost:8080/server/plans/plan \
  -H 'content-type: application/json' 
  -d '{
    "plan": {
		"id":          "111",
		"name":        "plan2",
		"org_id":      "org2",
		"description": "description1",
		"creator_id":  "222",
		"templateId":  "template1",
		"isTemplate":  true,
		"status":      1,
		"goals": [{"id": "1"}, {"id": "2"}]
	}
}'
{}
```

### Plan.Delete
```shell
curl -X DELETE \
  http://127.0.0.1:8080/server/plans/plan/111 \
  -H 'content-type: application/json'
{}
```

### Plan.All

```shell
curl -X GET \
  http://127.0.0.1:8080/server/plans/all \
  -H 'content-type: application/json'
{
    "data": {
        "plans": [
            {
                "id": "666",
                "name": "plan1",
                "org_id": "orgid",
                "description": "description1",
                "created": 1518138048,
                "updated": 1518138048,
                "goals": [
                    {
                        "id": "1"
                    },
                    {
                        "id": "2"
                    }
                ],
                "creator_id": "222",
                "isTemplate": true,
                "templateId": "template1",
                "status": 1
            }
        ]
    }
}
```

### Plan.Search


```shell
curl -X POST \
  http://127.0.0.1:8080/server/plans/search \
  -H 'content-type: application/json' \
  -d '{
	"name":"plan1"
}'
{
    "data": {
        "plans": [
            {
                "id": "666",
                "name": "plan1",
                "org_id": "orgid",
                "description": "description1",
                "created": 1518138048,
                "updated": 1518138048,
                "goals": [
                    {
                        "id": "1"
                    },
                    {
                        "id": "2"
                    }
                ],
                "creator_id": "222",
                "isTemplate": true,
                "templateId": "template1",
                "status": 1
            }
        ]
    }
}
```

### Plan.Templates

```shell
curl -X GET \
  http://127.0.0.1:8080/server/plans/templates \
  -H 'content-type: application/json'
{
    "data": {
        "plans": [
            {
                "id": "666",
                "name": "plan1",
                "org_id": "orgid",
                "description": "description1",
                "created": 1518138048,
                "updated": 1518138048,
                "goals": [
                    {
                        "id": "1"
                    },
                    {
                        "id": "2"
                    }
                ],
                "creator_id": "222",
                "isTemplate": true,
                "templateId": "template1",
                "status": 1
            }
        ]
    }
}
```

### Plan.Drafts

```shell
curl -X GET \
  http://127.0.0.1:8080/server/plans/drafts \
  -H 'content-type: application/json'
{
    "data": {
        "plans": [
            {
                "id": "666",
                "name": "plan1",
                "org_id": "orgid",
                "description": "description1",
                "created": 1518138048,
                "updated": 1518138048,
                "goals": [
                    {
                        "id": "1"
                    },
                    {
                        "id": "2"
                    }
                ],
                "creator_id": "222",
                "isTemplate": true,
                "templateId": "template1",
                "status": 1
            }
        ]
    }
}
```

### Plan.ByCreator
```shell
curl -X GET \
  http://127.0.0.1:8080/server/plans/creator/222 \
  -H 'content-type: application/json'
{
    "data": {
        "plans": [
            {
                "id": "666",
                "name": "plan1",
                "org_id": "orgid",
                "description": "description1",
                "created": 1518138048,
                "updated": 1518138048,
                "goals": [
                    {
                        "id": "1"
                    },
                    {
                        "id": "2"
                    }
                ],
                "creator_id": "222",
                "isTemplate": true,
                "templateId": "template1",
                "status": 1
            }
        ]
    }
}
```

### Plan.Filters

```shell
curl -X GET \
  http://127.0.0.1:8080/server/plans/filters \
  -H 'content-type: application/json'
{}
```

### Plan.TimeFilters

```shell
curl -X GET \
  http://127.0.0.1:8080/server/plans/filter/time?start_date=1517790886&end_date=1517990886 \
  -H 'content-type: application/json'
{
    "data": {
        "plans": [
            {
                "id": "666",
                "name": "plan1",
                "org_id": "orgid",
                "description": "description1",
                "created": 1518138048,
                "updated": 1518138048,
                "goals": [
                    {
                        "id": "1"
                    },
                    {
                        "id": "2"
                    }
                ],
                "creator_id": "222",
                "isTemplate": true,
                "templateId": "template1",
                "status": 1
            }
        ]
    }
}
```

### Plan.GoalFilters

```shell
curl -X GET \
  http://127.0.0.1:8080/server/plans/filter/goal?filter=1,2 \
  -H 'content-type: application/json'
{
    "data": {
        "plans": [
            {
                "id": "666",
                "name": "plan1",
                "org_id": "orgid",
                "description": "description1",
                "created": 1518138048,
                "updated": 1518138048,
                "goals": [
                    {
                        "id": "1"
                    },
                    {
                        "id": "2"
                    }
                ],
                "creator_id": "222",
                "isTemplate": true,
                "templateId": "template1",
                "status": 1
            }
        ]
    }
}
```