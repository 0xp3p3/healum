package db

import (
	"errors"
	"log"
	"strings"

	mdb "server/db-srv/proto/db"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/selector"
)

// Database Driver
type Driver interface {
	NewDB(nodes ...*registry.Node) (DB, error)
}

// An initialised DB connection
// Must call Init to load Database/Table
type DB interface {
	// Initialise the database
	// If the database doesn't exist
	// will throw an error
	Init(mdb *mdb.Database) error
	// Close the connection
	Close() error
	// Query commands
	Read(id, extraId string) (*mdb.Record, error)
	Create(mdb *mdb.Record) error
	Update(mdb *mdb.Record) error
	Delete(id, extraId string) error
	Search(md map[string]string, from, to, limit, offset int64, reverse bool) ([]*mdb.Record, error)
	RunQuery(query string) ([]*mdb.Record, error)
	CreateDatabase(name string) error
	DeleteDatabase(name string) error
}

type db struct {
	selector    selector.Selector
	namespace   string
	defaultName string
	drivers     map[string]Driver
	driverKey   string
}

var (
	DefaultDB *db

	// Prefix for lookup in registry
	DBServiceNamespace = "go.micro.db"

	// Default database name
	DBServiceName = "kv"

	// Used to lookup the metadata for the driver
	DBDriverKey = "driver"

	// Default driver
	DefaultDriver = "mysql"

	// supported drivers: mysql, elastic, cassandra
	Drivers = map[string]Driver{}

	// Errors
	ErrNotFound     = errors.New("not found")
	ErrNotAvailable = errors.New("not available")
)

func NewDB(s selector.Selector) *db {
	return &db{
		selector:    s,
		namespace:   DBServiceNamespace,
		defaultName: DBServiceName,
		drivers:     Drivers,
		driverKey:   DBDriverKey,
	}
}

func (d *db) name(db *mdb.Database) string {
	name := strings.ToLower(db.Driver)
	if len(name) == 0 {
		name = DefaultDriver
	}
	if strings.HasPrefix(db.Name, "test_") {
		name = db.Name
	}
	// TODO: check by driver not by default name
	return strings.Join([]string{d.namespace, name}, ".")
}

// looks up a registered database and its driver, creates a connection by the driver and prepares queries
func (d *db) lookup(db *mdb.Database) (DB, error) {
	dbname := d.name(db)
	next, err := d.selector.Select(dbname)
	if err != nil {
		return nil, err
	}

	var id string
	// TODO: create a node list rather than connecting to one
	for {
		node, err := next()
		if err != nil {
			return nil, err
		}

		// seen all?
		if node.Id == id {
			return nil, ErrNotAvailable
		}

		id = node.Id

		// is the driver set?
		dv, ok := node.Metadata[d.driverKey]
		if !ok {
			continue
		}
		// does a database want to use the driver
		if len(db.Driver) != 0 {
			if strings.ToLower(db.Driver) != strings.ToLower(dv) {
				continue
			}
		} else {
			if strings.ToLower(DefaultDriver) != strings.ToLower(dv) {
				continue
			}
		}

		// is the driver supported?
		dr, ok := d.drivers[dv]
		if !ok {
			continue
		}
		conn, err := dr.NewDB(node)
		if err != nil {
			return nil, err
		}
		if err := conn.Init(db); err != nil {
			return nil, err
		}
		return conn, nil
	}

	return nil, ErrNotAvailable
}

func (d *db) Init() {
}

func (d *db) Read(db *mdb.Database, id, extraId string) (*mdb.Record, error) {
	dr, err := d.lookup(db)
	if err != nil {
		return nil, err
	}
	defer dr.Close()
	return dr.Read(id, extraId)
}

func (d *db) Create(db *mdb.Database, r *mdb.Record) error {
	dr, err := d.lookup(db)
	if err != nil {
		return err
	}
	defer dr.Close()
	return dr.Create(r)
}

func (d *db) Update(db *mdb.Database, r *mdb.Record) error {
	dr, err := d.lookup(db)
	if err != nil {
		return err
	}
	defer dr.Close()
	return dr.Update(r)
}

func (d *db) Delete(db *mdb.Database, id, extraId string) error {
	dr, err := d.lookup(db)
	if err != nil {
		return err
	}
	defer dr.Close()
	return dr.Delete(id, extraId)
}

func (d *db) Search(db *mdb.Database, md map[string]string, from, to, limit, offset int64, reverse bool) ([]*mdb.Record, error) {
	dr, err := d.lookup(db)
	if err != nil {
		return nil, err
	}
	defer dr.Close()
	return dr.Search(md, from, to, limit, offset, reverse)
}

func (d *db) RunQuery(db *mdb.Database, query string) ([]*mdb.Record, error) {
	dr, err := d.lookup(db)
	if err != nil {
		return nil, err
	}
	defer dr.Close()
	return dr.RunQuery(query)
}

func (d *db) CreateDatabase(db *mdb.Database) error {
	dr, err := d.lookup(db)
	if err != nil {
		return err
	}
	defer dr.Close()
	return dr.CreateDatabase(db.Name)
}

func (d *db) DeleteDatabase(db *mdb.Database) error {
	dr, err := d.lookup(db)
	if err != nil {
		return err
	}
	defer dr.Close()
	return dr.DeleteDatabase(db.Name)
}

func Init(s selector.Selector) error {
	DefaultDB = NewDB(s)
	return nil
}

func Read(db *mdb.Database, id, extraId string) (*mdb.Record, error) {
	return DefaultDB.Read(db, id, extraId)
}

func Create(db *mdb.Database, r *mdb.Record) error {
	return DefaultDB.Create(db, r)
}

func Update(db *mdb.Database, r *mdb.Record) error {
	return DefaultDB.Update(db, r)
}

func Delete(db *mdb.Database, id, extraId string) error {
	return DefaultDB.Delete(db, id, extraId)
}

func Search(db *mdb.Database, md map[string]string, from, to, limit, offset int64, reverse bool) ([]*mdb.Record, error) {
	return DefaultDB.Search(db, md, from, to, limit, offset, reverse)
}

func RunQuery(db *mdb.Database, query string) ([]*mdb.Record, error) {
	return DefaultDB.RunQuery(db, query)
}

func CreateDatabase(db *mdb.Database) error {
	log.Println(DefaultDB)

	return DefaultDB.CreateDatabase(db)
}

func DeleteDatabase(db *mdb.Database) error {
	return DefaultDB.DeleteDatabase(db)
}
