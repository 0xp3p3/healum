# Sidecar Broker

This is a broker plugin for the micro [sidecar](https://github.com/micro/micro/tree/master/car)

## Usage

Here's a simple usage guide

### Run Sidecar

```
go get github.com/micro/micro
```

```
micro sidecar
```

### Import and Flag plugin

```
import _ "github.com/micro/go-plugins/broker/sidecar"
```

```
go run main.go --broker=sidecar
```
