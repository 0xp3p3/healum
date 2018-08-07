package redis

import (
	"encoding/json"
	"fmt"
	redis "gopkg.in/redis.v4"

	"github.com/micro/go-micro/registry"
	"hash/fnv"
	"server/db-srv/db"
	mdb "server/db-srv/proto/db"
	"strconv"
	"sync"
	"time"
)

type redisDriver struct{}

type redisDB struct {
	sync.RWMutex
	redisIndex string
	redisType  string
	url        string
	clients    map[string]*redis.Client
}

func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32()) & 16
}

func init() {
	// Other drives should be added the same way (elasticsearch, redis, etc.)
	db.Drivers["redis"] = new(redisDriver)
}

// records can pass expiryTime in the metadata
func expiryTime(r *mdb.Record) time.Duration {
	var res time.Duration = 0
	timeStr, ok := r.Metadata["expiryTime"]
	if !ok {
		return res
	}
	delete(r.Metadata, "expiryTime")
	parsed, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return res
	}
	return time.Duration(parsed)
}

func (d *redisDriver) NewDB(nodes ...*registry.Node) (db.DB, error) {

	if len(nodes) == 0 {
		return nil, db.ErrNotAvailable
	}
	url := fmt.Sprintf("%s:%d", nodes[0].Address, nodes[0].Port)

	return &redisDB{
		url:     url,
		clients: make(map[string]*redis.Client),
	}, nil
}

func (d *redisDB) Init(mdb *mdb.Database) error {
	d.CreateDatabase(mdb.Table)
	d.redisIndex = mdb.Name
	d.redisType = mdb.Table

	return nil
}

func (d *redisDB) Close() error {
	d.Lock()
	defer d.Unlock()
	client, ok := d.clients[d.redisType]
	if !ok {
		client.Close()
	}

	return nil
}

// Reads a generic record from a database. Records must be split across corresponding databases
func (d *redisDB) Read(id, extraId string) (*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()

	client, ok := d.clients[d.redisType]
	if !ok {
		return nil, db.ErrNotFound
	}

	r := client.Get(id)

	pl := &mdb.Record{}
	val := r.Val()
	if len(val) == 0 {
		return nil, nil
	}
	if err := json.Unmarshal([]byte(val), &pl); err != nil {
		return nil, err
	}

	return pl, nil
}

// Creates a generic record in a database. Any datatype can be stored in the structure
func (d *redisDB) Create(r *mdb.Record) error {
	d.RLock()
	defer d.RUnlock()
	if r.Created == 0 {
		r.Created = time.Now().Unix()
	}
	expiry := expiryTime(r)
	r.Updated = time.Now().Unix()
	res, err := json.Marshal(r)

	if err != nil {
		return err
	}
	client, ok := d.clients[d.redisType]
	if ok {
		client.Set(r.Id, string(res), expiry)
	}

	return nil
}

// U of CRUD for generic records
func (d *redisDB) Update(r *mdb.Record) error {
	d.RLock()
	defer d.RUnlock()
	if r.Created == 0 {
		r.Created = time.Now().Unix()
	}
	expiry := expiryTime(r)
	r.Updated = time.Now().Unix()
	res, err := json.Marshal(r)
	if err != nil {
		return err
	}
	client, ok := d.clients[d.redisType]
	if ok {
		client.Set(r.Id, string(res), expiry)
	}

	return nil
}

// D of CRUD for generic records
func (d *redisDB) Delete(id, extraId string) error {
	d.RLock()
	defer d.RUnlock()
	client, ok := d.clients[d.redisType]
	if ok {
		client.Del(id)
	}

	return nil
}

// name and parameter are provided through name and parameter1 values of the metadata parameter. If they exists, the
// search run using the parameters. Otherwise it performs metadata-related search.
// Search returns all records if no search keys provided
// Unix timestemp interval or limit and offset work as expected
func (d *redisDB) Search(md map[string]string, from, to, limit, offset int64, reverse bool) ([]*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()

	var records []*mdb.Record

	return records, nil
}

func (d *redisDB) RunQuery(query string) ([]*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()

	var records []*mdb.Record

	return records, nil
}

// A database must be created for every datatype (User, Auth, Room, etc.)
func (d *redisDB) DeleteDatabase(name string) error {
	d.Lock()
	defer d.Unlock()
	client, ok := d.clients[name]
	if !ok {
		client = redis.NewClient(&redis.Options{
			Addr:     d.url,
			Password: "",
			DB:       hash(name),
		})
	}
	client.FlushDb()
	return nil
}

func (d *redisDB) CreateDatabase(name string) error {
	d.Lock()
	defer d.Unlock()
	d.clients[name] = redis.NewClient(&redis.Options{
		Addr:     d.url,
		Password: "",
		DB:       hash(name),
	})
	return nil
}
