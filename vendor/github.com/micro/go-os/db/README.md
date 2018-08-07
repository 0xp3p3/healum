# DB [![GoDoc](https://godoc.org/github.com/micro/go-os?status.svg)](https://godoc.org/github.com/micro/go-os/db)

DB is a high level pluggable abstraction for CRUD via RPC.

The motivation is to create a DBaaS layer that 
allows RPC based proxying so that we can leverage go-micro and all the plugins. This allows auth, 
rate limiting, tracing and all the other things to be used. What we lose in database drivers we gain 
in not having to write CRUD a thousand times over.

<p align="center">
  <img src="https://github.com/micro/go-os/blob/master/doc/db.png" />
</p>

## Interface

Initial thoughts lie around a CRUD interface. The number of times 
one has to write CRUD on top of database libraries, having to think 
through schema and data modelling based on different databases is a 
pain. Going lower level than this doesn't pose any value.

Event sourcing can be tackled in a separate package.

```go
type DB interface {
        Close() error
        Init(...Option) error
        Options() Options
        Read(id string) (Record, error)
        Create(r Record) error
        Update(r Record) error
        Delete(id string) error
        Search(...SearchOption) ([]Record, error)
        String() string
}

type Option func(*Options)

type Metadata map[string]interface{}

type Record interface {
        Id() string
        Created() int64
        Updated() int64
        Metadata() Metadata
        Bytes() []byte
        Scan(v interface{}) error
}

func NewDB(opts ...Option) DB {
        return newOS(opts...)
}

func NewRecord(id string, md Metadata, data interface{}) Record {
        return newRecord(id, md, data)
}

```

## Supported Backends

- [DB Service](https://github.com/micro/db-srv)

## Usage

DB builds on go-micro/client and a backend RPC based db service that manages multiple backends like mysql, elasticsearch, etc. 
Doing this means go-micro/client wrappers can be used for rate limiting, auth, etc. 

```go
package main

import (
        "fmt"
        "github.com/micro/go-micro/cmd"
        "github.com/micro/go-os/db"
        "github.com/pborman/uuid"
)

func main() {
        // Create database instance
        database := db.NewDB(
                db.Database("foo"),
                db.Table("bar"),
        )
        
        // Create Thing type
        type Thing struct {
                Name string
        }
        
        // Create a new record
        record := db.NewRecord(
                // record id
                uuid.NewUUID().String(),
                // record metadata
                db.Metadata{"key": "value"},
                // record value
                &Thing{"dbthing"}),
        )
        
        fmt.Printf("Creating record: id: %s metadata: %+v bytes: %+v\n", record.Id(), record.Metadata(), string(record.Bytes()))
        
        // Create the record
        if err := database.Create(record); err != nil {
                fmt.Println(err)
                return
        }
        
        // Read the record back
        rec, err := database.Read(record.Id())
        if err != nil {
                fmt.Println(err)
                return
        }
        
        thing := new(Thing)
        
        // Scan into type Thing
        if err := rec.Scan(&thing); err != nil {
                fmt.Println("Error scanning read record", err)
                return
        }
        
        fmt.Printf("Read record: id: %s metadata: %+v bytes: %+v\n", rec.Id(), rec.Metadata(), thing)
        
        fmt.Println("Searching for metadata key:value")
        
        // Search using metadata
        records, err := database.Search(
                db.WithMetadata(db.Metadata{"key": "value"}), 
                db.WithLimit(10),
                db.WithOffset(0),
        )
        if err != nil {
                fmt.Println(err)
                return
        }
        
        for _, record := range records {
                thing := new(Thing)
        
                if err := record.Scan(&thing); err != nil {
                        fmt.Println("Error scanning record", record.Id(), err)
                        return
                }
        
                fmt.Printf("Record: id: %s metadata: %+v bytes: %+v\n", record.Id(), record.Metadata(), thing)
        }
        
        fmt.Println("Deleting", record.Id())
        
        // Delete the record
        if err := database.Delete(record.Id()); err != nil {
                fmt.Println(err)
                return
        }
}
```
