package storage

import (
	"errors"
	mdb "server/cloudkey-srv/proto/record"
)

// Database Driver
type Driver interface {
	NewStorage() (ST, error)
}

// An initialised Storage
// Must call Init to load data structures
type ST interface {
	// Initialise the storage
	// If the database doesn't exist
	// will throw an error
	Init(mdb *mdb.Storage) error
	// Close the connection
	Close() error
	// Query commands
	CreateKey(orgid, cryptoKey string) (error)
	EncryptKey(orgid, cryptoKey, dek string) (string, error)
	DecryptKey(orgid, cryptoKey, encryptedDek string) (string, error)
}

type storage struct {
	defaultName string
	drivers     map[string]Driver
	driverKey   string
}

var (
	DefaultStorage *storage

	// Default database name
	DBServiceName = "kv"

	// Used to lookup the metadata for the driver
	DBDriverKey = "driver"

	// supported drivers: gce, aws
	Drivers = map[string]Driver{}

	// Errors
	ErrNotAvailable = errors.New("not available")
)

func NewStorage() *storage {
	return &storage{
		defaultName: DBServiceName,
		drivers:     Drivers,
		driverKey:   DBDriverKey,
	}
}

func (d *storage) name(db *mdb.Storage) string {

	if db != nil && len(db.Driver) != 0 {
		return db.Driver
	}
	for k, _ := range(d.drivers){
		return k
	}
	return ""
}


// looks up a registered Storage and its driver, creates a connection by the driver and prepares queries
func (d *storage) lookup(db *mdb.Storage) (ST, error) {
	driverName := d.name(db)
	// is the driver supported?
	dr, ok := d.drivers[driverName]
	if !ok {
		return nil, ErrNotAvailable
	}

	conn, err := dr.NewStorage()
	if err != nil {
		return nil, err
	}
	if err := conn.Init(db); err != nil {
		return nil, err
	}

	return conn, nil
}

func (d *storage) CreateKey(db *mdb.Storage, orgid, cryptoKey string) error {
	dr, err := d.lookup(db)
	if err != nil {
		return err
	}
	defer dr.Close()
	return dr.CreateKey(orgid, cryptoKey)
}

func (d *storage) EncryptKey(db *mdb.Storage, orgid, cryptoKey, dek string) (string, error) {
	dr, err := d.lookup(db)
	if err != nil {
		return "", err
	}
	defer dr.Close()
	return dr.EncryptKey(orgid, cryptoKey, dek)
}

func (d *storage) DecryptKey(db *mdb.Storage, orgid, cryptoKey, encryptedDek string) (string, error) {
	dr, err := d.lookup(db)
	if err != nil {
		return "", err
	}
	defer dr.Close()
	return dr.DecryptKey(orgid, cryptoKey, encryptedDek)
}

func Init() error {
	DefaultStorage = NewStorage()
	return nil
}

func CreateKey(db *mdb.Storage, orgid, cryptoKey string) error {
	return DefaultStorage.CreateKey(db, orgid, cryptoKey)
}

func EncryptKey(db *mdb.Storage, orgid, cryptoKey, dek string) (string, error) {
	return DefaultStorage.EncryptKey(db, orgid, cryptoKey, dek)
}

func DecryptKey(db *mdb.Storage, orgid, cryptoKey, encryptedDek string) (string, error) {
	return DefaultStorage.DecryptKey(db, orgid, cryptoKey, encryptedDek)
}
