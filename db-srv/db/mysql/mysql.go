package mysql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/go-micro/registry"
	"math"
	"server/db-srv/db"
	mdb "server/db-srv/proto/db"
)

type mysqlDriver struct{}

type mysqlDB struct {
	url string

	sync.RWMutex
	name    string
	table   string
	conn    *sql.DB
	queries map[string]*sql.Stmt
}

var (
	DBUser = "root"
	DBPass = ""
)

func init() {
	// Other drives should be added the same way (elasticsearch, redis, etc.)
	db.Drivers["mysql"] = new(mysqlDriver)
}

func (d *mysqlDriver) NewDB(nodes ...*registry.Node) (db.DB, error) {
	if len(nodes) == 0 {
		return nil, db.ErrNotAvailable
	}
	url := fmt.Sprintf("tcp(%s:%d)/", nodes[0].Address, nodes[0].Port)

	// add credentials
	// TODO: take database credentials
	if len(DBUser) > 0 && len(DBPass) > 0 {
		url = fmt.Sprintf("%s:%s@%s", DBUser, DBPass, url)
	} else if len(DBUser) > 0 {
		url = fmt.Sprintf("%s@%s", DBUser, url)
	}

	// test the connection
	conn, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return &mysqlDB{
		url:     url,
		queries: make(map[string]*sql.Stmt),
	}, nil
}

func (d *mysqlDB) Init(mdb *mdb.Database) error {
	d.CreateDatabase(mdb.Name)
	// create connection
	dconn, err := sql.Open("mysql", d.url+mdb.Name)
	if err != nil {
		return err
	}

	d.conn = dconn
	if _, err = d.conn.Exec(fmt.Sprintf(mysqlSchema, mdb.Table)); err != nil {
		return err
	}

	for _, migration := range(migrations) {
		d.conn.Exec(fmt.Sprintf(migration, mdb.Table))
	}

	for query, statement := range mysqlQueries {
		prepared, err := d.conn.Prepare(fmt.Sprintf(statement, mdb.Name, mdb.Table))
		if err != nil {
			return err
		}
		d.queries[query] = prepared
	}
	d.name = mdb.Name
	d.table = mdb.Table

	return nil
}

func (d *mysqlDB) Close() error {
	d.Lock()
	defer d.Unlock()
	return d.conn.Close()
}

// Reads a generic record from a database. Records must be split across corresponding databases
func (d *mysqlDB) Read(id, extraId string) (*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()

	r := &mdb.Record{}

	var row *sql.Row
	if len(extraId) != 0 {
		row = d.queries["read"].QueryRow(id, extraId)
	} else {
		row = d.queries["readId"].QueryRow(id)
	}

	var meta []byte
	var extra float64
	if err := row.Scan(&r.Id, &r.Created, &r.Updated, &r.Name, &r.Parameter1, &r.Parameter2, &r.Parameter3, &r.Lat, &r.Lng, &meta, &extra); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal(meta, &r.Metadata); err != nil {
		return nil, err
	}

	return r, nil
}

// Creates a generic record in a database. Any datatype can be stored in the structure
func (d *mysqlDB) Create(r *mdb.Record) error {
	d.RLock()
	defer d.RUnlock()
	meta, err := json.Marshal(r.Metadata)
	if err != nil {
		return err
	}
	if r.Created == 0 {
		r.Created = time.Now().Unix()
	}
	r.Updated = time.Now().Unix()
	_, err = d.queries["create"].Exec(r.Id, r.Created, r.Updated, r.Name, r.Parameter1, r.Parameter2, r.Parameter3, r.Lat, r.Lng, string(meta))
	return err
}

// U of CRUD for generic records
func (d *mysqlDB) Update(r *mdb.Record) error {
	d.RLock()
	defer d.RUnlock()

	meta, err := json.Marshal(r.Metadata)
	if err != nil {
		return err
	}
	r.Updated = time.Now().Unix()

	_, err = d.queries["update"].Exec(r.Updated, r.Name, r.Parameter1, r.Parameter2, r.Parameter3, r.Lat, r.Lng, string(meta), r.Id)

	return nil
}

// D of CRUD for generic records
func (d *mysqlDB) Delete(id, extraId string) error {
	d.RLock()
	defer d.RUnlock()
	_, err := d.queries["delete"].Exec(id, extraId)
	return err
}

// name and parameter are provided through name and parameter1 values of the metadata parameter. If they exists, the
// search run using the parameters. Otherwise it performs metadata-related search.
// Search returns all records if no search keys provided
// Unix timestemp interval or limit and offset work as expected
func (d *mysqlDB) Search(md map[string]string, from, to, limit, offset int64, reverse bool) ([]*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()

	var rows *sql.Rows
	var err error

	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}
	// if from and to are not set, they are maximum values to filter properly
	if to == 0 {
		to = math.MaxInt32
	}
	order := "Asc"

	if reverse {
		order = "Desc"
	}
	name := md["name"]
	delete(md, "name")

	paramter1 := md["parameter1"]
	delete(md, "parameter1")

	// search parameter
	paramter2 := md["parameter2"]
	delete(md, "parameter2")

	// search parameter
	paramter3 := md["parameter3"]
	delete(md, "parameter3")

	// distance search
	distance := md["distance"]
	if distance == "0"{
		distance = "0.001"
	}
	delete(md, "distance")
	lat := md["lat"]
	delete(md, "lat")
	lng := md["lng"]
	delete(md, "lng")
	// search first
	if len(distance) != 0 && len(lat) != 0 && len(lng) != 0 {
		rows, err = d.queries["distanceSearch"].Query(lat, lng, lat, distance)
	} else if len(name) != 0 && len(paramter1) != 0 && len(paramter2) != 0 && len(paramter3) != 0 {
		rows, err = d.queries["sNameAndParameter1Parameter2Parameter3" + order].Query(from, to, paramter3, name, paramter1, paramter2, limit, offset)
	}  else if len(name) != 0 && len(paramter1) != 0 && len(paramter3) != 0 {
		rows, err = d.queries["sNameAndParameter1Parameter3" + order].Query(from, to, paramter3, name, paramter1, limit, offset)
	} else if len(paramter2) != 0 && len(paramter3) != 0 {
		// trailing %parameter2% is needed to search the message anywhere in the string
		paramter2 := "%" + paramter2 + "%"
		rows, err = d.queries["sParameter2Parameter3" + order].Query(from, to, paramter3, paramter2, limit, offset)
	} else if len(name) != 0 && len(paramter1) != 0 && len(paramter2) != 0 {
		paramter2 := "%" + paramter2 + "%"
		rows, err = d.queries["sNameAndParameter1Parameter2" + order].Query(from, to, paramter2, name, paramter1, limit, offset)
	} else if len(paramter2) != 0 {
		// trailing %parameter2% is needed to search the message anywhere in the string
		paramter2 := "%" + paramter2 + "%"
		rows, err = d.queries["sParameter2" + order].Query(from, to, paramter2, limit, offset)
	} else if len(name) != 0 && len(paramter1) != 0 && len(paramter2) != 0 {
		rows, err = d.queries["sNameAndParameter1Parameter2" + order].Query(from, to, paramter2, name, paramter1, limit, offset)
	} else if len(name) != 0 && len(paramter3) != 0 {
		rows, err = d.queries["sNameParameter3" + order].Query(from, to, paramter3, name, limit, offset)
	} else if len(name) != 0 && len(paramter2) != 0 {
		rows, err = d.queries["sNameParameter2" + order].Query(from, to, name, paramter2, limit, offset)
	} else if len(name) != 0 && len(paramter1) != 0 {
		rows, err = d.queries["sNameParameter1" + order].Query(from, to, name, paramter1, limit, offset)
	} else if len(paramter1) != 0 && len(paramter3) != 0 {
		rows, err = d.queries["sParameter1Parameter3" + order].Query(from, to, paramter3, paramter1, limit, offset)
		fmt.Println(paramter1, paramter3)
	} else if len(paramter1) != 0 && len(paramter2) != 0 {
		rows, err = d.queries["sParameter1Parameter2" + order].Query(from, to, paramter1, paramter2, limit, offset)
	} else if len(name) != 0 && len(paramter1) != 0 && len(paramter3) != 0 {
		rows, err = d.queries["sNameAndParameter1Parameter3" + order].Query(from, to, paramter3, name, paramter1, limit, offset)
	}  else if len(name) != 0 {
		rows, err = d.queries["sName" + order].Query(from, to, name, limit, offset)
	} else if len(paramter1) != 0 {
		rows, err = d.queries["sParameter1" + order].Query(from, to, paramter1, limit, offset)
	} else if len(paramter3) != 0 {
		rows, err = d.queries["searchParameter3" + order].Query(from, to, paramter3, limit, offset)
	} else if len(md) > 0 {
		// THIS IS SUPER CRUFT
		// TODO: DONT DO THIS
		// Note: Tried to use mariadb dynamic columns. They suck.
		var query string
		var args []interface{}

		// create statement for each key-val pair
		for k, v := range md {
			if len(query) == 0 {
				query += " "
			} else {
				query += "AND metadata like ? "
			}
			args = append(args, fmt.Sprintf(`%%"%s":"%s"%%`, k, v))
		}

		// append limit offset
		args = append(args, limit, offset)
		query += " limit ? offset ?"
		query = fmt.Sprintf(searchMetadataQ, d.name, d.table) + query
		// doe the query
		rows, err = d.conn.Query(query, args...)
	} else {
		rows, err = d.queries["search" + order].Query(from, to, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*mdb.Record
	for rows.Next() {
		r := &mdb.Record{}
		var meta []byte
		var extra float64
		if err := rows.Scan(&r.Id, &r.Created, &r.Updated, &r.Name, &r.Parameter1, &r.Parameter2, &r.Parameter3,  &r.Lat, &r.Lng, &meta, &extra); err != nil {
			if err == sql.ErrNoRows {
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		if err := json.Unmarshal(meta, &r.Metadata); err != nil {
			return nil, err
		}
		records = append(records, r)

	}
	if rows.Err() != nil {
		return nil, err
	}
	return records, nil
}

func (d *mysqlDB) RunQuery(query string) ([]*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()

	var rows *sql.Rows
	var err error

	var args []interface{}

	rows, err = d.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*mdb.Record
	for rows.Next() {
		r := &mdb.Record{}
		var meta []byte
		var extra float64
		if err := rows.Scan(&r.Id, &r.Created, &r.Updated, &r.Name, &r.Parameter1, &r.Parameter2, &r.Parameter3, &meta, &extra); err != nil {
			if err == sql.ErrNoRows {
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		if err := json.Unmarshal(meta, &r.Metadata); err != nil {
			return nil, err
		}
		records = append(records, r)

	}
	if rows.Err() != nil {
		return nil, err
	}
	return records, nil
}

// A database must be created for every datatype (User, Auth, Room, etc.)
func (d *mysqlDB) DeleteDatabase(name string) error {
	d.Lock()
	defer d.Unlock()
	conn, err := sql.Open("mysql", d.url)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create the database
	_, err = conn.Exec("DROP DATABASE IF EXISTS " + "`" + name + "`")
	return err
}

func (d *mysqlDB) CreateDatabase(name string) error {
	d.Lock()
	defer d.Unlock()
	// Create a conn to initialise the database
	conn, err := sql.Open("mysql", d.url)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create the database
	_, err = conn.Exec("CREATE DATABASE IF NOT EXISTS " + "`" + name + "`")
	return err
}
