package influxdb

import (
	"fmt"
	"sync"
	"server/db-srv/db"
	"github.com/micro/go-micro/registry"
	"time"
	mdb "server/db-srv/proto/db"

	"github.com/influxdata/influxdb/client/v2"
	"log"
	"encoding/json"
	"math"
)

type influxdbDriver struct{}

type influxdbDB struct {
	sync.RWMutex
	influxdbDatabase string
	influxdbMeasure  string
	cl               client.Client
}

func init() {
	// Other drives should be added the same way (influxdb, redis, etc.)
	db.Drivers["influxdb"] = new(influxdbDriver)
}

func (d *influxdbDriver) NewDB(nodes ...*registry.Node) (db.DB, error) {
	if len(nodes) == 0 {
		return nil, db.ErrNotAvailable
	}

	url := fmt.Sprintf("http://%s:%d", nodes[0].Address, nodes[0].Port)

	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: url,
		Username: "root",
		Password: "root",
	})
	if err != nil {
		log.Fatalln("Error: ", err)
	}

	// test the connection
	_, _, err = c.Ping(2 * time.Second)
	if err != nil {
		return nil, err
	}

	return &influxdbDB{
		cl:     c,
	}, nil
}

func (d *influxdbDB) Init(mdb *mdb.Database) error {
	d.CreateDatabase(mdb.Name)

	d.influxdbDatabase = mdb.Name
	d.influxdbMeasure = mdb.Table

	return nil
}

func (d *influxdbDB) Close() error {
	d.Lock()
	defer d.Unlock()
	d.cl.Close()
	return nil
}

// Reads a generic record from a database. Records must be split across corresponding databases
func (d *influxdbDB) Read(id, extraId string) (*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()
	qText := fmt.Sprintf("SELECT * FROM %v WHERE id = '%v' AND parameter3 = '%v'", d.influxdbMeasure, id, extraId)
	q := client.Query{
		Command:  qText,
		Database: d.influxdbDatabase,
	}
	var res []client.Result
	if response, err := d.cl.Query(q); err == nil {
		if response.Error() != nil {
			return nil, response.Error()
		}
		res = response.Results
	} else {
		return nil, err
	}

	pl := &mdb.Record{}
	if len(res) == 0 {
		return nil, db.ErrNotFound
	}
	for _, r := range (res) {
		if len(r.Series) == 0 {
			return nil, db.ErrNotFound
		}
		for _, s := range (r.Series) {
			pl.Id = s.Values[0][1].(string)
			pl.Name = s.Values[0][3].(string)
			pl.Parameter1 = s.Values[0][4].(string)
			pl.Parameter2 = s.Values[0][5].(string)
			pl.Parameter3 = s.Values[0][6].(string)
			pl.Updated, _ = s.Values[0][7].(json.Number).Int64()
			metaStr := s.Values[0][2].(string)
			json.Unmarshal([]byte(metaStr), &pl.Metadata)
			return pl, nil
		}
	}
	return nil, nil
}

// Creates a generic record in a database. Any datatype can be stored in the structure
func (d *influxdbDB) Create(r *mdb.Record) error {
	d.RLock()
	defer d.RUnlock()
	if r.Created == 0 {
		r.Created = time.Now().Unix()
	}
	r.Updated = time.Now().Unix()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  d.influxdbDatabase,
		Precision: "s",
	})

	if err != nil {
		log.Fatalln("Error: ", err)
	}
	// Create a point and add to batch
	tags := map[string]string{
		"id": r.Id,
		"name1": r.Name,
		"parameter1": r.Parameter1,
		"parameter2": r.Parameter2,
		"parameter3": r.Parameter3,
	}
	meta, err := json.Marshal(r.Metadata)
	fields := map[string]interface{}{
		"updated":   r.Updated,
		"metadata": string(meta),
	}
	pt, err := client.NewPoint(d.influxdbMeasure, tags, fields, time.Unix(r.Created, 0))

	if err != nil {
		log.Fatalln("Error: ", err)
	}

	bp.AddPoint(pt)

	err = d.cl.Write(bp)

	return err
}

// U of CRUD for generic records
func (d *influxdbDB) Update(r *mdb.Record) error {
	d.RLock()
	defer d.RUnlock()
	// influxdb cannot update records
	return nil
}

// D of CRUD for generic records
func (d *influxdbDB) Delete(id, extraId string) error {
	d.RLock()
	defer d.RUnlock()
	qText := fmt.Sprintf("DROP SERIES FROM %v WHERE id = '%v' AND parameter3 = '%v'", d.influxdbMeasure, id, extraId)
	q := client.Query{
		Command:  qText,
		Database: d.influxdbDatabase,
	}
	if response, err := d.cl.Query(q); err == nil {
		if response.Error() != nil {
			return response.Error()
		}
	} else {
		return err
	}
	return nil
}

// name and parameter are provided through name and parameter1 values of the metadata parameter. If they exists, the
// search run using the parameters. Otherwise it performs metadata-related search.
// Search returns all records if no search keys provided
// Unix timestemp interval or limit and offset work as expected
func (d *influxdbDB) Search(md map[string]string, from, to, limit, offset int64, reverse bool) ([]*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()

	// if from and to are not set, they are maximum values to filter properly
	if to == 0 {
		to = math.MaxInt32
	}

	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}
	where := fmt.Sprintf("WHERE time >= '%v' AND time <= '%v' ", time.Unix(from, 0).Format("2006-01-02 15:04:05"), time.Unix(to, 0).Format("2006-01-02 15:04:05"))
	name := md["name"]
	delete(md, "name")

	paramter1 := md["parameter1"]
	delete(md, "parameter1")
	// search parameter
	paramter2 := md["parameter2"]
	delete(md, "parameter2")

	paramter3 := md["parameter3"]
	delete(md, "parameter3")

	if len(paramter2) != 0 && len(paramter3) != 0 {
		where += fmt.Sprintf(" AND parameter2 = '%v' AND parameter3 = '%v'", paramter2, paramter3)
	}else if len(paramter1) != 0 && len(name) != 0 && len(paramter3) != 0{
		where += fmt.Sprintf(" AND parameter1 = '%v' AND name1 = '%v' AND parameter3 = '%v'", paramter1, name, paramter3)
	}else if len(paramter1) != 0 && len(paramter3) != 0{
		where += fmt.Sprintf(" AND parameter1 = '%v' AND parameter3 = '%v'", paramter1, paramter3)
	}else if len(name) != 0 && len(paramter3) != 0 {
		where += fmt.Sprintf(" AND name1 = '%v' AND parameter3 = '%v'", name, paramter3)
	}else if len(paramter2) != 0 {
		where += fmt.Sprintf(" AND parameter2 = '%v'", paramter2)
	}else if len(paramter1) != 0 && len(name) != 0{
		where += fmt.Sprintf(" AND parameter1 = '%v' AND name1 = '%v'", paramter1, name)
	}else if len(paramter1) != 0{
		where += fmt.Sprintf(" AND parameter1 = '%v'", paramter1)
	}else if len(name) != 0 {
		where += fmt.Sprintf(" AND name1 = '%v'", name)
	}

	if (reverse){
		where += " ORDER BY time DESC"
	}

	qText := fmt.Sprintf("SELECT * FROM %v %v", d.influxdbMeasure, where)
	q := client.Query{
		Command:  qText,
		Database: d.influxdbDatabase,
	}
	var res []client.Result
	if response, err := d.cl.Query(q); err == nil {
		if response.Error() != nil {
			return nil, response.Error()
		}
		res = response.Results
	} else {
		return nil, err
	}
	var records []*mdb.Record

	if len(res) == 0 {
		return nil, db.ErrNotFound
	}
	for _, r := range (res) {
		for _, s := range (r.Series) {
			for _, v := range (s.Values) {
				pl := &mdb.Record{}
				pl.Id = v[1].(string)
				pl.Name = v[3].(string)
				pl.Parameter1 = v[4].(string)
				pl.Parameter2 = v[5].(string)
				pl.Parameter3 = v[6].(string)
				pl.Updated, _ = v[7].(json.Number).Int64()
				metaStr := v[2].(string)
				json.Unmarshal([]byte(metaStr), &pl.Metadata)
				records = append(records, pl)
			}
		}
	}
	return records, nil
}

func (d *influxdbDB) RunQuery(query string) ([]*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()

	q := client.Query{
		Command:  query,
		Database: d.influxdbDatabase,
	}
	var res []client.Result
	if response, err := d.cl.Query(q); err == nil {
		if response.Error() != nil {
			return nil, response.Error()
		}
		res = response.Results
	} else {
		return nil, err
	}
	var records []*mdb.Record

	if len(res) == 0 {
		return nil, db.ErrNotFound
	}
	for _, r := range (res) {
		for _, s := range (r.Series) {
			for _, v := range (s.Values) {
				pl := &mdb.Record{}
				pl.Id = v[1].(string)
				pl.Name = v[3].(string)
				pl.Parameter1 = v[4].(string)
				pl.Parameter2 = v[5].(string)
				pl.Parameter3 = v[6].(string)
				pl.Updated, _ = v[7].(json.Number).Int64()
				metaStr := v[2].(string)
				json.Unmarshal([]byte(metaStr), &pl.Metadata)
				records = append(records, pl)
			}
		}
	}
	return records, nil
}

// A database must be created for every datatype (User, Auth, Room, etc.)
func (d *influxdbDB) DeleteDatabase(name string) error {
	d.Lock()
	defer d.Unlock()
	q := client.NewQuery(fmt.Sprintf("DROP DATABASE %v", name), "", "")
	_, err := d.cl.Query(q)
	if err != nil {
		return err
	}
	return nil
}

func (d *influxdbDB) CreateDatabase(name string) error {
	d.Lock()
	defer d.Unlock()
	q := client.NewQuery(fmt.Sprintf("CREATE DATABASE %v", name), "", "")
	_, err := d.cl.Query(q)
	if err != nil {
		return err
	}
	return nil
}
