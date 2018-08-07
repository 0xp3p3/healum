package arangodb

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	lib "github.com/solher/arangolite"
	lib_req "github.com/solher/arangolite/requests"

	"context"
	"math"
	"server/common"
	"server/db-srv/db"
	mdb "server/db-srv/proto/db"
	"sync"
	"time"

	"github.com/micro/go-micro/registry"
)

type arangodbDriver struct{}

type arangodbDB struct {
	sync.RWMutex
	dbCon            *lib.Database
	arangoUrl        string
	arangoName       string
	arangoCollection string
	graph            bool
	graphName        string
	graphFrom        string
	graphTo          string
}

var (
	DBUser = "root"
	DBPass = ""
)

type Node struct {
	lib.Document
}

type Record struct {
	Id         string            `json:"id,omitempty"`
	Created    int64             `json:"created,omitempty"`
	Updated    int64             `json:"updated,omitempty"`
	Name       string            `json:"name,omitempty"`
	Parameter1 string            `json:"parameter1,omitempty"`
	Parameter2 string            `json:"parameter2,omitempty"`
	Parameter3 string            `json:"parameter3,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	Data       interface{}       `json:"data,omitempty"`
}

func init() {
	// Other drives should be added the same way (arangodb, redis, etc.)
	db.Drivers["arangodb"] = new(arangodbDriver)
}

func (d *arangodbDriver) InitDb() {
}

//new database connection
func (d *arangodbDriver) NewDB(nodes ...*registry.Node) (db.DB, error) {
	if len(nodes) == 0 {
		return nil, db.ErrNotAvailable
	}
	//change this to SSL
	url := fmt.Sprintf("http://%s:%d", nodes[0].Address, nodes[0].Port)
	dbCon := lib.NewDatabase(
		lib.OptEndpoint(url),
		lib.OptBasicAuth(DBUser, DBPass),
	)

	return &arangodbDB{
		dbCon:     dbCon,
		arangoUrl: url,
	}, nil
}

func (d *arangodbDB) Init(db *mdb.Database) error {
	d.arangoName = db.Name
	d.CreateDatabase(db.Name)
	d.dbCon.Options(
		lib.OptDatabaseName(db.Name),
	)

	ctx := context.Background()

	for _, cols := range common.DbHealum {
		if len(cols) == 0 {
			continue
		}
		// general collection
		if len(cols) == 1 {
			d.dbCon.Run(ctx, nil, &lib_req.CreateCollection{Name: cols[0]})
		}

		// edge collection
		if len(cols) == 4 {
			if err := d.dbCon.Run(ctx, nil, &lib_req.GetGraph{Name: cols[1]}); err != nil {
				switch {
				case lib.IsErrNotFound(err):
					// create from and to collection if they don't exist
					d.dbCon.Run(ctx, nil, &lib_req.CreateCollection{Name: cols[2]}) // remove this after refactor
					d.dbCon.Run(ctx, nil, &lib_req.CreateCollection{Name: cols[3]}) // remove this after refactor

					// create edge collection
					d.dbCon.Run(ctx, nil, &lib_req.CreateCollection{Name: cols[0], Type: 3})

					// If graph does not exist, create a new one.
					edgeDefinitions := []lib_req.EdgeDefinition{
						{
							Collection: cols[0],
							From:       []string{cols[2]},
							To:         []string{cols[3]},
						},
					}
					d.dbCon.Run(ctx, nil, &lib_req.CreateGraph{Name: cols[1], EdgeDefinitions: edgeDefinitions})
				default:
					common.ErrorLog(common.DbSrv, d.Init, err, "Init is failed")
					log.Fatal(err)
				}
			}
		}
	}

	if d.graph {

	} else {
	}
	return nil
}

// func (d *arangodbDB) Init(db *mdb.Database) error {
// 	d.arangoName = db.Name
// 	d.arangoCollection = db.Table
// 	_, graph := db.Metadata[common.GraphFlag]
// 	if graph {
// 		d.graph = true
// 		d.graphName = db.Metadata[common.GraphName]
// 		d.graphFrom = db.Metadata[common.GraphFrom]
// 		d.graphTo = db.Metadata[common.GraphTo]
// 	}

// 	ctx := context.Background()

// 	d.CreateDatabase(db.Name)
// 	d.dbCon.Options(
// 		lib.OptDatabaseName(db.Name),
// 	)
// 	if d.graph {
// 		// d.dbCon.Run(ctx, nil, &lib_req.CreateCollection{Name: d.arangoCollection, Type: 3})
// 		// d.dbCon.Run(ctx, nil, &lib_req.CreateCollection{Name: "edges", Type: 3})
// 		// d.dbCon.Run(ctx, nil, &lib_req.CreateCollection{Name: "from", Type: 3})
// 		// d.dbCon.Run(ctx, nil, &lib_req.CreateCollection{Name: "to", Type: 3})
// 		// d.dbCon.Run(ctx, nil, &lib_req.CreateGraph{
// 		// 	Name:              d.arangoCollection,
// 		// 	OrphanCollections: []string{"from", "to"},
// 		// 	EdgeDefinitions: []lib_req.EdgeDefinition{
// 		// 		{
// 		// 			"edges",
// 		// 			[]string{"from"},
// 		// 			[]string{"to"},
// 		// 		},
// 		// 	},
// 		// })

// 		if err := d.dbCon.Run(ctx, nil, &lib_req.GetGraph{Name: d.graphName}); err != nil {
// 			switch {
// 			case lib.IsErrNotFound(err):
// 				// create from and to collection if they don't exist
// 				d.dbCon.Run(ctx, nil, &lib_req.CreateCollection{Name: d.graphFrom}) // remove this after refactor
// 				d.dbCon.Run(ctx, nil, &lib_req.CreateCollection{Name: d.graphTo})   // remove this after refactor

// 				// create edge collection
// 				d.dbCon.Run(ctx, nil, &lib_req.CreateCollection{Name: d.arangoCollection, Type: 3})

// 				// If graph does not exist, create a new one.
// 				edgeDefinitions := []lib_req.EdgeDefinition{
// 					{
// 						Collection: d.arangoCollection,
// 						From:       []string{d.graphFrom},
// 						To:         []string{d.graphTo},
// 					},
// 				}
// 				d.dbCon.Run(ctx, nil, &lib_req.CreateGraph{Name: d.graphName, EdgeDefinitions: edgeDefinitions})
// 			default:
// 				common.ErrorLog(common.DbSrv, d.Init, err, "Init is failed")
// 				log.Fatal(err)
// 			}
// 		}
// 	} else {
// 		d.dbCon.Run(ctx, nil, &lib_req.CreateCollection{Name: d.arangoCollection})
// 	}
// 	return nil
// }

func (d *arangodbDB) Close() error {
	d.Lock()
	defer d.Unlock()
	return nil
}

// Reads a generic record from a database. Records must be split across corresponding databases
func (d *arangodbDB) Read(id, orgid string) (*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()
	ctx := context.Background()

	nodes := []mdb.Record{}
	//	var records []byte
	var err error
	if d.graph {
		q := lib_req.NewAQL(`
		    FOR n
		    IN %s
		    FILTER n.id == "%s" && n.parameter3 == "%s"
		    RETURN n
		  `, "edges", id, orgid)
		err = d.dbCon.Run(ctx, &nodes, q)
		if err != nil {
			return nil, err
		}

	} else {
		q := lib_req.NewAQL(`
		    FOR n
		    IN %s
		    FILTER n.id == "%s" && n.parameter3 == "%s"
		    RETURN n
		  `, d.arangoCollection, id, orgid)

		err = d.dbCon.Run(ctx, &nodes, q)
		if err != nil {
			return nil, err
		}
	}

	if len(nodes) == 0 {
		return nil, db.ErrNotFound
	}

	return &nodes[0], nil
}

// Creates a generic record in a database. Any datatype can be stored in the structure
func (d *arangodbDB) Create(r *mdb.Record) error {
	d.RLock()
	defer d.RUnlock()
	if r.Created == 0 {
		r.Created = time.Now().Unix()
	}

	ctx := context.Background()

	r.Updated = time.Now().Unix()

	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	if d.graph {
		q := lib_req.NewAQL(`INSERT %v INTO %v`, fmt.Sprintf(`{_key: "%s"}`, r.Name), "from")
		_ = d.dbCon.Run(ctx, nil, q)
		q1 := lib_req.NewAQL(`INSERT %v INTO %v`, fmt.Sprintf(`{_key: "%s"}`, r.Parameter1), "to")
		_ = d.dbCon.Run(ctx, nil, q1)
		nobracketsData := data[1 : len(data)-1]
		q2 := lib_req.NewAQL(`INSERT %v INTO %v`, fmt.Sprintf(`{_from: "from/%s", _to: "to/%s", %s}`, r.Name, r.Parameter1, nobracketsData), "edges")
		err = d.dbCon.Run(ctx, nil, q2)
		return err

	} else {
		q := lib_req.NewAQL(`INSERT %v INTO %v`, string(data), d.arangoCollection)
		err = d.dbCon.Run(ctx, nil, q)
		return err
	}
	return nil
}

// U of CRUD for generic records
func (d *arangodbDB) Update(r *mdb.Record) error {
	d.RLock()
	defer d.RUnlock()
	if r.Created == 0 {
		r.Created = time.Now().Unix()
	}

	ctx := context.Background()
	r.Updated = time.Now().Unix()
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	if d.graph {
		q := lib_req.NewAQL(`INSERT %v INTO %v`, fmt.Sprintf(`{_key: "%s"}`, r.Name), "from")
		_ = d.dbCon.Run(ctx, nil, q)
		q1 := lib_req.NewAQL(`INSERT %v INTO %v`, fmt.Sprintf(`{_key: "%s"}`, r.Parameter1), "to")
		_ = d.dbCon.Run(ctx, nil, q1)
		nobracketsData := data[1 : len(data)-1]
		q2 := lib_req.NewAQL(`
			FOR n
		    	IN %s
		    	FILTER n.id == "%s"
			UPDATE n WITH %s IN %s`,
			"edges",
			r.Id,
			fmt.Sprintf(`{_from: "from/%s", _to: "to/%s", %s}`, r.Name, r.Parameter1, nobracketsData),
			"edges",
		)
		err = d.dbCon.Run(ctx, nil, q2)
		return err

	} else {
		q := lib_req.NewAQL(`
		    FOR n
		    IN %s
		    FILTER n.id == "%s"
		    UPDATE n WITH %s IN %s
		  `, d.arangoCollection, r.Id, string(data), d.arangoCollection)

		err = d.dbCon.Run(ctx, nil, q)
		if err != nil {
			return err
		}
	}

	return nil
}

// D of CRUD for generic records
func (d *arangodbDB) Delete(id, orgid string) error {
	d.RLock()
	defer d.RUnlock()

	ctx := context.Background()

	if d.graph {
		q2 := lib_req.NewAQL(`
			FOR n
		    	IN %s
		    	FILTER n.id == "%s" && n.parameter3 == "%s"
			REMOVE n IN %s`,
			"edges",
			id,
			orgid,
			"edges",
		)
		err := d.dbCon.Run(ctx, nil, q2)
		if err != nil {
			return err
		}

	} else {
		q := lib_req.NewAQL(`
		    FOR n
		    IN %s
		    FILTER n.id == "%s" && n.parameter3 == "%s"
		    REMOVE n IN %s
		  `, d.arangoCollection, id, orgid, d.arangoCollection)

		err := d.dbCon.Run(ctx, nil, q)
		if err != nil {
			return err
		}
	}

	return nil
}

// name and parameter are provided through name and parameter1 values of the metadata parameter. If they exists, the
// search run using the parameters. Otherwise it performs metadata-related search.
// Search returns all records if no search keys provided
// Unix timestemp interval or limit and offset work as expected
func (d *arangodbDB) Search(md map[string]string, from, to, limit, offset int64, reverse bool) ([]*mdb.Record, error) {
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

	size := fmt.Sprintf("%d", limit)
	offs := fmt.Sprintf("%d", offset)

	fromS := fmt.Sprintf("%d", from)
	toS := fmt.Sprintf("%d", to)
	query := fmt.Sprintf(`FILTER n.created >= %s && n.created <= %s`, fromS, toS)

	name := md["name"]
	delete(md, "name")

	paramter1 := md["parameter1"]
	delete(md, "parameter1")

	// search parameter
	paramter2 := md["parameter2"]
	delete(md, "parameter2")

	ctx := context.Background()

	nodes := []*mdb.Record{}
	var err error
	if d.graph {
		var q *lib_req.AQL
		if len(name) != 0 && len(paramter1) != 0 {
			q = lib_req.NewAQL(`
		   FOR x IN INTERSECTION ((FOR v, e, p IN 1..1 OUTBOUND "from/%s" edges RETURN e) , (FOR v, e, p IN 1..1 INBOUND "to/%s" edges RETURN e)) RETURN x`, name, paramter1)
		} else if len(name) != 0 {
			q = lib_req.NewAQL(`
		   FOR v, e, p IN 1..1 OUTBOUND "from/%s" edges RETURN e`, name)
		} else {
			q = lib_req.NewAQL(`
		   FOR v, e, p IN 1..1 INBOUND "to/%s" edges RETURN e`, paramter1)
		}

		err = d.dbCon.Run(ctx, &nodes, q)
		if err != nil {
			return nil, err
		}

	} else {
		if len(paramter2) != 0 {
			query = query + fmt.Sprintf(` && n.parameter2 == "%s"`, paramter2)
		} else if len(name) != 0 && len(paramter1) != 0 {
			query = query + fmt.Sprintf(` && n.name == "%s" && n.parameter1 == "%s"`, name, paramter1)
		} else if len(name) != 0 {
			query = query + fmt.Sprintf(` && n.name == "%s"`, name)
		} else if len(paramter1) != 0 {
			query = query + fmt.Sprintf(` && n.parameter1 == "%s"`, paramter1)
		} else if len(md) > 0 {
			// TODO all records for now
		}
		sortQ := "SORT n.created "
		if reverse {
			sortQ = sortQ + "DESC"
		}
		q := lib_req.NewAQL(`
		    FOR n
		    IN %s
		    %s
		    LIMIT %s, %s
		    %s
		    RETURN n
		  `, d.arangoCollection, query, offs, size, sortQ)

		fmt.Println(q)
		records := []*Record{}
		err = d.dbCon.Run(ctx, &records, q)
		if err != nil {
			return nil, err
		}
		fmt.Println("records:", records)

		for _, r := range records {
			node := &mdb.Record{
				Id:         r.Id,
				Created:    r.Created,
				Updated:    r.Updated,
				Name:       r.Name,
				Parameter1: r.Parameter1,
				Parameter2: r.Parameter2,
				// Parameter3: r.Parameter3,
				Metadata: r.Metadata,
			}
			if r.Data != nil {
				body, err := json.Marshal(r.Data)
				if err != nil {
					continue
				}
				node.Parameter3 = string(body)
			}
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

func (d *arangodbDB) RunQuery(query string) ([]*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()
	ctx := context.Background()
	query = strings.Replace(query, `'`, `\'`, -1)
	query = strings.Replace(query, `%`, `%%`, -1)
	fmt.Println(query)

	records := []*Record{}
	nodes := []*mdb.Record{}
	var err error

	q := lib_req.NewAQL(query)
	err = d.dbCon.Run(ctx, &records, q)
	if err != nil {
		common.ErrorLog(common.DbSrv, d.RunQuery, err, "RunQuery is failed")
		return nil, err
	}
	// fmt.Println("records:", records)

	for _, r := range records {
		// fmt.Println("record:", r)
		node := &mdb.Record{
			Id:         r.Id,
			Created:    r.Created,
			Updated:    r.Updated,
			Name:       r.Name,
			Parameter1: r.Parameter1,
			Parameter2: r.Parameter2,
			// Parameter3: r.Parameter3,
			Metadata: r.Metadata,
		}
		if r.Data != nil {
			body, err := json.Marshal(r.Data)
			if err != nil {
				common.ErrorLog(common.DbSrv, d.RunQuery, err, "Object marshaling is failed")
				continue
			}
			node.Parameter3 = string(body)
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// A database must be created for every datatype (User, Auth, Room, etc.)
func (d *arangodbDB) DeleteDatabase(name string) error {
	d.Lock()
	defer d.Unlock()
	ctx := context.Background()

	d.dbCon.Options(
		lib.OptDatabaseName("_system"),
	)
	err := d.dbCon.Run(ctx, nil, &lib_req.DropDatabase{
		Name: name,
	})
	return err
}

func (d *arangodbDB) CreateDatabase(name string) error {
	d.Lock()
	defer d.Unlock()
	ctx := context.Background()
	d.dbCon.Options(
		lib.OptDatabaseName("_system"),
	)
	d.dbCon.Run(ctx, nil, &lib_req.CreateDatabase{
		Name: name,
		Users: []map[string]interface{}{
			{"username": DBUser, "passwd": DBPass},
		},
	})
	d.dbCon.Options(
		lib.OptDatabaseName(d.arangoName),
	)

	return nil
}
