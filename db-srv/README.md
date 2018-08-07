# DB Server

**Experimental DB Server**

The DB server is an experimental proxy layer for backend databases. 

The DB uses the registry to find supported databases for proxying. It expects the metadata key "driver" to be set 
to the type of database. It should match one of the supported drivers.

Entities to be stored in a database include:
-  users
-  profiles
-  rooms
-  message history
-  contact list
-  settings
- ...

## Database Drivers

- mysql (mariadb)
- redis 
- elasticsearch
- arangodb
- influxdb



## Getting started

1. Install Consul

	Consul is the default registry/discovery for go-micro apps. It's however pluggable.
	[https://www.consul.io/intro/getting-started/install.html](https://www.consul.io/intro/getting-started/install.html)

2. Run Consul
	```
	$ consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul
	```

3. Start one of the supported databases (mysql, elasticsearch, ...) and register with the registry.
    Install MySQL
    
    ```
	$ brew install mysql
	```

    Run MySQL
    
    ```
	$ ln -sfv /usr/local/opt/mysql/*.plist  ~/Library/LaunchAgents
	$ launchctl load ~/Library/LaunchAgents/homebrew.mxcl.mysql.plist
	```

	```
	Example. Register location of the **kv** database hosted by mysql

	$ micro register service '{"name": "go.micro.db.mysql", "version": "0.0.1", "nodes": [{"id": "kv-1", "address": "127.0.0.1", "port": 3306, "metadata": {"driver": "mysql"}}]}'
	```
4. Download and start the service

	```shell
	$ go get server/db-srv
	db-srv -config config.json --database_service_namespace=go.micro.db --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222
	```

	OR as a docker container

	```shell
	docker run microhq/db-srv -config config.json --database_service_namespace=go.micro.db --registry_address=YOUR_REGISTRY_ADDRESS --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222
	```

## The API
DB server implements the following RPC Methods

DB
- Read
- Create
- Update
- Delete
- Search
- CreateDatabase
- DeleteDatabase


### DB.Create

```
micro query go.micro.srv.db DB.Create '{"database": {"name": "foo", "table": "bar"}, "record": {"id": "e7add322-e069-44c2-b920-c4fbfd62e6b5", "metadata": {"key": "value"}}}'
```

### DB.Read

```
micro query go.micro.srv.db DB.Read '{"database": {"name": "foo", "table": "bar"}, "id": "e7add322-e069-44c2-b920-c4fbfd62e6b5"}'

{
	"record": {
		"created": 1.454704366e+09,
		"id": "e7add322-e069-44c2-b920-c4fbfd62e6b5",
		"metadata": {
			"key": "value"
		},
		"updated": 1.454704366e+09
	}
}
```

### DB.Search

```
micro query go.micro.srv.db DB.Search '{"database": {"name": "foo", "table": "bar"}, "metadata": {"name": ""}}'

{
	"records": [
		{
			"created": 1.454704366e+09,
			"id": "e7add322-e069-44c2-b920-c4fbfd62e6b5",
			"metadata": {
				"key": "value"
			},
			"updated": 1.454704366e+09
		}
	]
}
```

### DB.Update

```
micro query go.micro.srv.db DB.Update '{"database": {"name": "foo", "table": "bar"}, "record": {"id": "e7add322-e069-44c2-b920-c4fbfd62e6b5", "metadata": {"key": "value", "key2": "value2"}}}'
```

### DB.Delete

```
micro query go.micro.srv.db DB.Delete '{"database": {"name": "foo", "table": "bar"}, "id": "e7add322-e069-44c2-b920-c4fbfd62e6b5"}'

```

### DB.CreateDatabase

```
micro query go.micro.srv.db DB.CreateDatabase '{"database": {"name": "foo", "table": "bar"}}'

```


### DB.DeleteDatabase

```
micro query go.micro.srv.db DB.DeleteDatabase '{"database": {"name": "foo", "table": "bar"}}'

```
