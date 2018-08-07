package elastic

import (
	"encoding/json"
	"fmt"
	elib "github.com/mattbaird/elastigo/lib"

	"github.com/micro/go-micro/registry"
	"math"
	"server/common"
	"server/db-srv/db"
	mdb "server/db-srv/proto/db"
	"sync"
	"time"
	"strings"
	"strconv"
	"errors"
)

type elasticsearchDriver struct{}

type elasticsearchDB struct {
	sync.RWMutex
	elasticIndex string
	elasticType  string
	connection   *elib.Conn
}

func init() {
	// Other drives should be added the same way (elasticsearch, redis, etc.)
	db.Drivers["elasticsearch"] = new(elasticsearchDriver)
}

func (d *elasticsearchDriver) NewDB(nodes ...*registry.Node) (db.DB, error) {

	if len(nodes) == 0 {
		return nil, db.ErrNotAvailable
	}
	c := elib.NewConn()
	c.Domain = nodes[0].Address
	c.Port = strconv.Itoa(nodes[0].Port)

	// test the connection
	_, err := c.Health()
	if err != nil {
		return nil, err
	}

	return &elasticsearchDB{
		connection: c,
	}, nil
}

func (d *elasticsearchDB) Init(mdb *mdb.Database) error {
	d.CreateDatabase(mdb.Name)
	d.elasticIndex = mdb.Name
	d.elasticType = mdb.Table

	return nil
}

func (d *elasticsearchDB) Close() error {
	d.Lock()
	defer d.Unlock()
	return nil
}

// Reads a generic record from a database. Records must be split across corresponding databases
func (d *elasticsearchDB) Read(id, extraId string) (*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()
	pl := &mdb.Record{}
	res, err := d.connection.Get(d.elasticIndex, d.elasticType, id, nil)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(*res.Source, pl)
	if pl == nil || pl.Id == "" {
		return nil, errors.New("not found")
	}

	return pl, nil
}

type hash map[string]interface{}

type Pin struct {
	Id string `json:"id"`
	//Location struct {
	//	Lat float64 `json:"lat"`
	//	Lon float64 `json:"lon"`
	//} `json:"location"`
	Location string `json:"location"`
}

// Creates a generic record in a database. Any datatype can be stored in the structure
func (d *elasticsearchDB) Create(r *mdb.Record) error {
	d.RLock()
	defer d.RUnlock()
	if r.Created == 0 {
		r.Created = time.Now().Unix()
	}
	r.Updated = time.Now().Unix()
	_, err := d.connection.Index(d.elasticIndex, d.elasticType, r.Id, nil, r)
	//mappingOpts := elib.MappingOptions{
	//	//Id:        elib.IdOptions{Index: "analyzed", Path: "id"},
	//	Properties: map[string]interface{}{
	//		// special properties that can't be expressed as tags
	//		"parameter21": map[string]interface{}{
	//			"type": "geo_point",
	//		},
	//
	//		"id": map[string]interface{}{
	//			"type": "text",
	//		},
	//	},
	//}
	//////options := `"pin": {
	//////               "properties": {
	//////                   "location": {
	//////                       "type": "geo_point"
	//////                   }
	//////               }
	//////           }`
	////pin := Pin{}
	////pin.Id = r.Id
	//err := d.connection.PutMapping(d.elasticIndex, d.elasticType, mdb.Record{}, mappingOpts)
	//fmt.Println(err)
	////_, err = d.connection.Index(d.elasticIndex, d.elasticType, r.Id, nil, r)
	////fmt.Println(err)
	////err = d.connection.PutMapping(d.elasticIndex, d.elasticType, *r, mappingOpts)
	////fmt.Println(err)
	////_, err = d.connection.Index(d.elasticIndex, d.elasticType, r.Id, nil, pin)
	//
	////err = d.connection.PutMapping(d.elasticIndex, d.elasticType, Pin{}, mappingOpts)
	return err
}

// U of CRUD for generic records
func (d *elasticsearchDB) Update(r *mdb.Record) error {
	d.RLock()
	defer d.RUnlock()
	if r.Created == 0 {
		r.Created = time.Now().Unix()
	}
	r.Updated = time.Now().Unix()
	_, err := d.connection.Index(d.elasticIndex, d.elasticType, r.Id, nil, r)
	return err
}

// D of CRUD for generic records
func (d *elasticsearchDB) Delete(id, extraId string) error {
	d.RLock()
	defer d.RUnlock()
	_, err := d.connection.Delete(d.elasticIndex, d.elasticType, id, nil)
	return err
}

// name and parameter are provided through name and parameter1 values of the metadata parameter. If they exists, the
// search run using the parameters. Otherwise it performs metadata-related search.
// Search returns all records if no search keys provided
// Unix timestemp interval or limit and offset work as expected
func (d *elasticsearchDB) Search(md map[string]string, from, to, limit, offset int64, reverse bool) ([]*mdb.Record, error) {
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
	query := "*"

	name := md["name"]
	delete(md, "name")

	paramter1 := md["parameter1"]
	delete(md, "parameter1")
	// search parameter
	paramter2 := md["parameter2"]
	delete(md, "parameter2")

	paramter3 := md["parameter3"]
	delete(md, "parameter3")

	_, hasAutocomlete := md[common.SearchableAutocompleteMeta]
	delete(md, common.SearchableAutocompleteMeta)

	if len(name) != 0 && len(paramter1) != 0 && len(paramter2) != 0 && len(paramter3) != 0 {
		query = fmt.Sprintf("name:\"%v\"  AND parameter1:\"%v\"  AND parameter2:\"%v\"  AND parameter3:\"%v\"", name, paramter1, paramter2, paramter3)
	} else if len(name) != 0 && len(paramter1) != 0 && len(paramter2) != 0 {
		query = fmt.Sprintf("name:\"%v\"  AND parameter1:\"%v\"  AND parameter2:\"%v\" ", name, paramter1, paramter2)
	} else if len(name) != 0 && len(paramter1) != 0 {
		query = fmt.Sprintf("name:\"%v\"  AND parameter1:\"%v\" ", name, paramter1)
	} else if len(name) != 0 && len(paramter2) != 0 {
		query = fmt.Sprintf("name:\"%v\" AND parameter2:\"%v\" ", name, paramter2)
	} else if len(paramter1) != 0 && len(paramter2) != 0 {
		query = fmt.Sprintf("parameter1:\"%v\" AND parameter2:\"%v\" ", paramter1, paramter2)
	} else if len(paramter2) != 0 && len(paramter3) != 0 {
		query = fmt.Sprintf("parameter2:\"%v\"  AND parameter3:\"%v\"", paramter2, paramter3)
		//query = paramter2 + "&" + paramter3
	} else if len(name) != 0 && len(paramter1) != 0 && len(paramter3) != 0 {
		//query = paramter1 + "&" + name + "&" + paramter3
		query = fmt.Sprintf("parameter1:\"%v\" AND name:\"%v\" AND parameter3:\"%v\"", paramter1, name, paramter3)
	} else if len(name) != 0 && len(paramter3) != 0 {
		//query = name + "&" + paramter3
		query = fmt.Sprintf("name:\"%v\" AND parameter3:\"%v\"", name, paramter3)
	} else if len(paramter1) != 0 && len(paramter3) != 0 {
		//query = paramter1 + "&" + paramter3
		query = fmt.Sprintf("parameter1:\"%v\" AND parameter3:\"%v\"", paramter1, paramter3)
	} else if len(paramter2) != 0 {
		query = fmt.Sprintf("parameter2:\"%v\"", paramter2)
		//query = paramter2
	} else if len(name) != 0 && len(paramter1) != 0 {
		//query = paramter1 + "&" + name
		query = fmt.Sprintf("parameter1:\"%v\" AND name:\"%v\"", paramter1, name)
	} else if len(name) != 0 {
		//query = name
		query = fmt.Sprintf("name:\"%v\"", name)
	} else if len(paramter1) != 0 {
		//query = paramter1
		query = fmt.Sprintf("parameter1:\"%v\"", paramter1)
	} else if len(paramter3) != 0 {
		//query = paramter3
		query = fmt.Sprintf("parameter3:\"%v\"", paramter3)
	} else if len(md) > 0 {
		// TODO custom map search?
	}

	var records []*mdb.Record
	if hasAutocomlete {
		elasticQuery := elib.Search(d.elasticIndex).Type(d.elasticType).Query(
			elib.Query().Search(query)).Size(size)
		out, err := elasticQuery.Result(d.connection)
		if err != nil {
			if strings.Contains(err.Error(), "No mapping found"){
				return records, nil
			}
			return nil, err
		}

		for _, hit := range out.Hits.Hits {
			var rec *mdb.Record
			if err := json.Unmarshal(*hit.Source, &rec); err != nil {
				return nil, err
			}
			records = append(records, rec)
		}
	} else {
		fromStr := fmt.Sprintf("%d", from)
		toStr := fmt.Sprintf("%d", to)
		//distanceFilter := elib.Filter().GeoDistance("100km", elib.NewGeoField("location", 32.3, 23.4))
		//fmt.Println(distanceFilter)
		//rangeOp := elib.Filter().And(elib.Filter().Range("created", fromStr, nil, toStr, nil, ""))
		//rangeOp := elib.Filter().Range("created", nil, nil, to, nil, "")
		//rangeOp1 := elib.CompoundFilter(rangeOp)
		//rangeOp1.Bool(string(elib.TEMBool))
		sorting := "asc"
		if reverse {
			sorting = "desc"
		}
		qryType := fmt.Sprintf(`
					{
					   "query":{
						  "bool":{
							 "must":[
								{
								   "query_string":{
									  "query": %q
								   }
								},
								{
								  "range": {
									"created": {
									  "gte": %s,
									  "lte": %s
									}
								  }
								}
							 ]
						  }
					   },
					   "sort":[
						  {
							 "created": {"order" :"%s"}
						  }
					   ]
					}
		`, query, fromStr, toStr, sorting)
		out, err := d.connection.Search(d.elasticIndex, d.elasticType, map[string]interface{} {"from" : offs, "size": size}, qryType)

		if err != nil {
			if strings.Contains(err.Error(), "No mapping found"){
				return records, nil
			}
			return nil, err
		}

		for _, hit := range out.Hits.Hits {
			var rec *mdb.Record
			if err := json.Unmarshal(*hit.Source, &rec); err != nil {
				return nil, err
			}
			records = append(records, rec)
		}
	}

	return records, nil
}

func (d *elasticsearchDB) RunQuery(query string) ([]*mdb.Record, error) {
	d.RLock()
	defer d.RUnlock()

	elasticQuery, err := d.connection.Search(d.elasticIndex, d.elasticIndex, nil, query);
	var records []*mdb.Record
	if err != nil {
		if strings.Contains(err.Error(), "No mapping found"){
			return records, nil
		}
		return nil, err
	}

	for _, hit := range elasticQuery.Hits.Hits {
		var rec *mdb.Record
		if err := json.Unmarshal(*hit.Source, &rec); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, nil
}

// A database must be created for every datatype (User, Auth, Room, etc.)
func (d *elasticsearchDB) DeleteDatabase(name string) error {
	d.Lock()
	defer d.Unlock()
	_, err := d.connection.DeleteIndex(name)
	return err
}

func (d *elasticsearchDB) CreateDatabase(name string) error {

	d.Lock()
	defer d.Unlock()
	res, _ := d.connection.IndicesExists(name)
	if !res {
		//mappingOpts := elib.MappingOptions{
		//	//Id:        elib.IdOptions{Index: "analyzed", Path: "id"},
		//	Properties: map[string]interface{}{
		//		// special properties that can't be expressed as tags
		//		"parameter21": map[string]interface{}{
		//			"type": "geo_point",
		//		},
		//
		//		"id": map[string]interface{}{
		//			"type": "text",
		//		},
		//	},
		//}
		//////options := `"pin": {
		//////               "properties": {
		//////                   "location": {
		//////                       "type": "geo_point"
		//////                   }
		//////               }
		//////           }`
		////pin := Pin{}
		////pin.Id = r.Id
		//_ = d.connection.PutMapping(d.elasticIndex, d.elasticType, mdb.Record{}, mappingOpts)
		d.connection.CreateIndex(name)
	}
	return nil
}
