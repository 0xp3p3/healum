// Package db is an interface for abstracting away CRUD.
package db

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

type SearchOption func(*SearchOptions)

type Metadata map[string]interface{}

type Record interface {
	Id() string
	Created() int64
	Updated() int64
	Metadata() Metadata
	Bytes() []byte
	Scan(v interface{}) error
}

var (
	DefaultDatabase = "micro"
	DefaultTable    = "micro"
)

func NewDB(opts ...Option) DB {
	return newOS(opts...)
}

func NewRecord(id string, md Metadata, data interface{}) Record {
	return newRecord(id, md, data)
}
