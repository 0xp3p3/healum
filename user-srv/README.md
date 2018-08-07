# User Server

User server is a microservice to store user accounts and perform simple authentication.

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
	go get server/user-srv
	user-srv -config config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222
	```

	OR as a docker container

	```shell
	docker run microhq/user-srv -config config.json --broker=nats --transport=nats --broker_address=127.0.0.1:4222 --transport_address=127.0.0.1:4222 --registry_address=YOUR_REGISTRY_ADDRESS
	```
## The API
Account server implements the following RPC Methods

Account
- Create
- Read
- Update
- Delete
- Search
- UpdatePassword
- Login
- Logout
- RecaptchaCheck
- ReadSession
- CreateFilter
- ReadFilter
- CreateTick
- SearchTick
- ReadOrCreatePresence
- UpdateOrCreatePresence
- CreateContact
- ReadContact
- UpdateContact
- DeleteContact
- ListContact


### Account.Create
```shell
$ micro query go.micro.srv.user Account.Create '{"user":{"id": "ff3c06de-9e43-41c7-9bab-578f6b4ad32b", "username": "asim", "email": "asim@example.com"}, "password": "password1"}'
{}
```

### Account.Read
```shell
$ micro query go.micro.srv.user Account.Read '{"id": "ff3c06de-9e43-41c7-9bab-578f6b4ad32b"}'
{
	"user": {
		"created": 1.450816182e+09,
		"email": "asim@example.com",
		"id": "ff3c06de-9e43-41c7-9bab-578f6b4ad32b",
		"updated": 1.450816182e+09,
		"username": "asim"
	}
}
```

### Account.Update
```shell
$ micro query go.micro.srv.user Account.Update '{"user":{"id": "ff3c06de-9e43-41c7-9bab-578f6b4ad32b", "username": "asim", "email": "asim+update@example.com"}}'
{}
```

### Account.UpdatePassword
```shell
$ micro query go.micro.srv.user Account.UpdatePassword '{"userId": "ff3c06de-9e43-41c7-9bab-578f6b4ad32b", "oldPassword": "password1", "newPassword": "newpassword1", "confirmPassword": "newpassword1" }'
{}
```

### Account.Delete
```shell
$ micro query go.micro.srv.user Account.Delete '{"id": "ff3c06de-9e43-41c7-9bab-578f6b4ad32b"}'
{}
```

### Auth.Login
```shell
$ micro query go.micro.srv.user Account.Login '{"username": "asim", "password": "password1"}'
{
	"session": {
		"created": 1.450816852e+09,
		"expires": 1.451421652e+09,
		"id": "sr7UEBmIMg5hYOgiljnhrd4XLsnalNewBV9KzpZ9aD8w37b3jRmEujGtKGcGlXPg1yYoSHR3RLy66ugglw0tofTNGm57NrNYUHsFxfwuGC6pvCn8BecB7aEF6UxTyVFq",
		"username": "asim"
	}
}
```

### Auth.ReadSession
```shell
$ micro query go.micro.srv.user Account.ReadSession '{"sessionId": "sr7UEBmIMg5hYOgiljnhrd4XLsnalNewBV9KzpZ9aD8w37b3jRmEujGtKGcGlXPg1yYoSHR3RLy66ugglw0tofTNGm57NrNYUHsFxfwuGC6pvCn8BecB7aEF6UxTyVFq"}'
{
	"session": {
		"created": 1.450816852e+09,
		"expires": 1.451421652e+09,
		"id": "sr7UEBmIMg5hYOgiljnhrd4XLsnalNewBV9KzpZ9aD8w37b3jRmEujGtKGcGlXPg1yYoSHR3RLy66ugglw0tofTNGm57NrNYUHsFxfwuGC6pvCn8BecB7aEF6UxTyVFq",
		"username": "asim"
	}
}
```

### Auth.Logout
```shell
$ micro query go.micro.srv.user Account.Logout '{"sessionId": "sr7UEBmIMg5hYOgiljnhrd4XLsnalNewBV9KzpZ9aD8w37b3jRmEujGtKGcGlXPg1yYoSHR3RLy66ugglw0tofTNGm57NrNYUHsFxfwuGC6pvCn8BecB7aEF6UxTyVFq"}'
{}
```

### Account.CreateFilter
```shell
$ micro query go.micro.srv.user Account.CreateFilter '{"filter":{"userid": "asim", "includeTypes": ["m.room.*"]}}'
{"filterId": "sr7UEBmIMg5hYOgiljnhrd4XLsnalNewBV9KzpZ9"}
```

### Auth.ReadFilter
```shell
$ micro query go.micro.srv.user Account.ReadFilter '{"filterId": "sr7UEBmIMg5hYOgiljnhrd4XLsnalNewBV9KzpZ"}'
{
	"session": {
		"created": 1.450816852e+09,
		"expires": 1.451421652e+09,
		"id": "sr7UEBmIMg5hYOgiljnhrd4XLsnalNewBV9KzpZ9aD8w",
		"userid": "asim",
		"includeTypes": ["m.room.*"],
	}
}
```

### Account.CreateTick

```
micro query go.micro.srv.user Account.CreateTick
'{
    "tick": {
        "userid": "1",
        "measure_name": "measure",
        "measure_value": "measure_value",
    }
}'
```

### Account.SearchTick

```
micro query go.micro.srv.user Account.SearchTick '{"userid": "1", "measure_name":"measure_name"}'
'{
    "ticks": [
        {
            "userid": "1",
            "measure_name": "measure",
            "measure_value": "measure_value",
        }
    ]
}'
```

### Account.ReadOrCreatePresence
```
micro query go.micro.srv.user Account.ReadOrCreatePresence
'{
    "presence": {
        "orgid": "1",
        "presence": "presence",
    }
}'
'{
    "presence": {
        "orgid": "1",
        "presence": "presence",
    }
}'
```

### Account.UpdateOrCreatePresence

```
micro query go.micro.srv.user Account.UpdateOrCreatePresence
'{
    "presence": {
        "orgid": "1",
        "presence": "presence",
    }
}'
'{
    "presence": {
        "orgid": "1",
        "presence": "presence",
    }
}'
```

### Account.CreateContact
```
micro query go.micro.srv.user Account.CreateContact
```

### Account.ReadContact

```
micro query go.micro.srv.user Account.ReadContact
'{
    "contact": {
        "orgid": "1",
        "sender_id": "1",
        "receiver_id": "2",
    }
}'
```

### Account.UpdateContact

```
micro query go.micro.srv.user Account.UpdateContact
'{
    "contact": {
        "orgid": "1",
        "sender_id": "1",
        "receiver_id": "2",
    }
}'
'{
    "contact": {
        "orgid": "1",
        "sender_id": "1",
        "receiver_id": "2",
    }
}'
```

### Account.DeleteContact

```
micro query go.micro.srv.user Account.DeleteContact {"id": "1"}
```

### Account.ListContact

```
micro query go.micro.srv.user Account.ListContact {"orgid": "1", "sender_id": "1"}
'{
    "contacts": [
        {
            "orgid": "1",
            "sender_id": "1",
            "receiver_id": "2",
        }
    ]
}'
```
