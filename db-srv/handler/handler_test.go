package handler

import (
	mdb "server/db-srv/proto/db"
	"testing"
	"time"

	"server/common"
	database "server/db-srv/db"
	_ "server/db-srv/db/arangodb"
	_ "server/db-srv/db/elastic"
	_ "server/db-srv/db/influxdb"
	_ "server/db-srv/db/mysql"
	_ "server/db-srv/db/redis"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/mock"
	"github.com/micro/go-micro/selector"
	"golang.org/x/net/context"
)

var (
	TestDBName    = "test_db"
	TestDBTable   = "test_table"
	TestDbService = &registry.Service{
		Name:    "go.micro.db.test_db",
		Version: "1.0.3",
		Nodes: []*registry.Node{
			{
				Id:      "go.micro.db.test_db-1.0.3-345",
				Address: "127.0.0.1",
				Port:    3306,
				Metadata: map[string]string{
					"driver": "mysql",
				},
			},
			{
				Id:      "go.micro.db.test_db-1.0.3-346",
				Address: "127.0.0.1",
				Port:    9200,
				Metadata: map[string]string{
					"driver": "elasticsearch",
				},
			},
			{
				Id:      "go.micro.db.test_db-1.0.3-347",
				Address: "127.0.0.1",
				Port:    6379,
				Metadata: map[string]string{
					"driver": "redis",
				},
			},
			{
				Id:      "go.micro.db.test_db-1.0.3-348",
				Address: "127.0.0.1",
				Port:    8529,
				Metadata: map[string]string{
					"driver": "arangodb",
				},
			},
			{
				Id:      "go.micro.db.test_db-1.0.3-349",
				Address: "127.0.0.1",
				Port:    8086,
				Metadata: map[string]string{
					"driver": "influxdb",
				},
			},
		},
	}
)

func TestDatabseIsValidated(t *testing.T) {
	db := mdb.Database{
		Name:  TestDBName,
		Table: TestDBTable,
	}
	if err := validateDB("DB.Read", &db); err != nil {
		t.Error(err)
	}
}

func initDb(t *testing.T, driver string, md map[string]string) {
	reg := mock.NewRegistry()
	reg.Register(TestDbService)
	rs := selector.NewSelector(
		func(o *selector.Options) {
			o.Registry = reg
			o.Strategy = selector.RoundRobin
		})
	if "mock" != rs.Options().Registry.String() {
		t.Error("Must be a mock registry")
	}
	services, err := reg.ListServices()
	if err != nil {
		t.Error(err)
	}
	has_service := false
	for _, serv := range services {
		if serv.Name == "go.micro.db.test_db" {
			has_service = true
			break
		}
	}
	if !has_service {
		t.Error("MySQL service must present in the registry")
	}

	database.DefaultDB = database.NewDB(rs)

	ctx := common.NewTestContext(context.TODO())
	req := &mdb.DeleteDatabaseRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   driver,
			Metadata: md,
		},
	}

	resp := &mdb.DeleteDatabaseResponse{}
	hdlr := NewWrapper(new(DB))
	if err := hdlr.DeleteDatabase(ctx, req, resp); err != nil {
		t.Error("Delete database is failed", err)
	}
}

func TestDbCreated(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
}

func TestRecordCreated(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name" {
		t.Error("Cannot retreave a record")
	}
}

func TestRecordUpdated(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_upd := &mdb.UpdateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name1",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp_upd := &mdb.UpdateResponse{}

	res_upd := hdlr.Update(ctx, req_upd, resp_upd)
	if res_upd != nil {
		t.Error(res_upd)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name1" {
		t.Error("Cannot retreave a record")
	}
}

func TestRecordDeleted(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_del := &mdb.DeleteRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		"111",
		"test_param3",
	}

	resp_del := &mdb.DeleteResponse{}

	res_del := hdlr.Delete(ctx, req_del, resp_del)
	if res_del != nil {
		t.Error(res_del)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read == nil {
		t.Error(res_read)
	}
}

func TestRecordSearch(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		map[string]string{
			"test_key": "test_value",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordDistanceSearch1(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Lat:        40.7130,
			Lng:        -74.0064,
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
	time.Sleep(2 * time.Second)
	req_read := &mdb.SearchRequest{
		Database: &mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		Metadata: map[string]string{
			"distance": "10",
			"lat":      "40.7140",
			"lng":      "-74.0064",
		},
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordDistanceSearch2(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Lat:        40.7130,
			Lng:        -74.0064,
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
	time.Sleep(2 * time.Second)
	req_read := &mdb.SearchRequest{
		Database: &mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		Metadata: map[string]string{
			"distance": "0.00001",
			"lat":      "40.7142",
			"lng":      "-74.0064",
		},
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 0 {
		t.Error("Non empty result")
		return
	}
}

func TestRecordDistanceSearch3(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Lat:        40.7130,
			Lng:        -74.0064,
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
	time.Sleep(2 * time.Second)
	req_read := &mdb.SearchRequest{
		Database: &mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		Metadata: map[string]string{
			"distance": "0.0001",
			"lat":      "40.7130",
			"lng":      "-74.0064",
		},
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch1(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearchAll(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		map[string]string{
			"name":       "test_name",
			"parameter1": "test_param",
			"parameter2": "test_param2",
			"parameter3": "test_param3",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearchParam1(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		map[string]string{
			"parameter1": "test_param",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch2(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param2",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		map[string]string{
			"name":       "test_name",
			"parameter1": "test_param",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch3(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		map[string]string{
			"name": "test_name1",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 0 {
		t.Error("Non empty result")
		return
	}
}

func TestRecordSearchAsc(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		2,
		0,
		false,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 2 {
		t.Error("Not enough results")
		return
	}

	if resp_read.Records[0].Id != "111" {
		t.Error("Bad sorting")
		return
	}

	if resp_read.Records[1].Id != "222" {
		t.Error("Bad sorting")
		return
	}

}

func TestRecordSearchDesc(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		2,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 2 {
		t.Error("Not enough results")
		return
	}

	if resp_read.Records[0].Id != "222" {
		t.Error("Bad sorting")
		return
	}

	if resp_read.Records[1].Id != "111" {
		t.Error("Bad sorting")
		return
	}

}

func TestRecordSearchFrom(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Add(-time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		map[string]string{
			"name": "test_name",
		},
		time.Now().Add(-time.Minute).Unix(),
		0,
		2,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}

	if len(resp_read.Records) != 1 {
		t.Error("Bad results count")
		return
	}

	if resp_read.Records[0].Id != "222" {
		t.Error("Bad from filter")
		return
	}

}

func TestRecordSearchTo(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Add(-time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		map[string]string{
			"name": "test_name",
		},
		0,
		time.Now().Add(-time.Minute).Unix(),
		2,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}

	if len(resp_read.Records) != 1 {
		t.Error("Bad results count")
		return
	}

	if resp_read.Records[0].Id != "111" {
		t.Error("Bad to filter")
		return
	}

}

func TestRecordSearchEmptyParams(t *testing.T) {
	initDb(t, "mysql", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Add(-time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		Database: &mdb.Database{Name: TestDBName, Table: TestDBTable},
		Limit:    2,
		Offset:   0,
		Reverse:  true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}

	if len(resp_read.Records) != 2 {
		t.Error("Bad results count")
		return
	}

	if resp_read.Records[0].Id != "222" {
		t.Error("Bad to filter")
		return
	}

}

// elastic
func TestDbCreatedElastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
}

func TestRecordCreatedElastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name" {
		t.Error("Cannot retreave a record")
	}
}

func TestRecordUpdatedElastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_upd := &mdb.UpdateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name1",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp_upd := &mdb.UpdateResponse{}

	res_upd := hdlr.Update(ctx, req_upd, resp_upd)
	if res_upd != nil {
		t.Error(res_upd)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name1" {
		t.Error("Cannot retreave a record")
	}
}

func TestRecordDeletedElastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_del := &mdb.DeleteRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		"111",
		"test_param3",
	}

	resp_del := &mdb.DeleteResponse{}

	res_del := hdlr.Delete(ctx, req_del, resp_del)
	if res_del != nil {
		t.Error(res_del)
	}

	time.Sleep(time.Second * 2)

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read == nil {
		t.Error(res_read)
	}
}

func TestRecordDistanceSearch1Elastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		&mdb.Record{
			Id:         "1111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Lat:        40.7130,
			Lng:        -74.0064,
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
	time.Sleep(2 * time.Second)
	//req_read := &mdb.SearchRequest{
	//	Database:&mdb.Database{
	//		Name:  TestDBName,
	//		Table: TestDBTable,
	//		Driver: "elasticsearch",
	//	},
	//	Metadata: map[string]string{
	//		"distance": "10",
	//		"lat": "40.7140",
	//		"lng": "-74.0064",
	//	},
	//}
	//
	//resp_read := &mdb.SearchResponse{}
	//
	//res_read := hdlr.Search(ctx, req_read, resp_read)
	//if res_read != nil {
	//	t.Error(res_read)
	//	return
	//}
	//if len(resp_read.Records) == 0 {
	//	t.Error("Empty result")
	//	return
	//}
	//if resp_read.Records[0].Name != "test_name" {
	//	t.Error("Empty result")
	//	return
	//}
}

func TestRecordSearchElastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		map[string]string{
			"test_key": "test_value",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let elasticsearch finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearchParam1Param2Elastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		map[string]string{
			"parameter1": "test_param",
			"parameter2": "test_param2",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let elasticsearch finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch1Elastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		1,
		0,
		true,
	}
	// to let elasticsearch finish indexing
	time.Sleep(time.Second * 2)
	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearchParam1Elastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		map[string]string{
			"parameter2": "test_param2",
		},
		0,
		0,
		1,
		0,
		true,
	}
	// to let elasticsearch finish indexing
	time.Sleep(time.Second * 2)
	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch2Rlastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		map[string]string{
			"name":       "test_name",
			"parameter1": "test_param",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let elasticsearch finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch3Elastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "elasticsearch",
		},
		map[string]string{
			"name": "test_name1",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let elasticsearch finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 0 {
		t.Error("Non empty result")
		return
	}
}

func TestRecordSearchAscElastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	created := time.Now()
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		&mdb.Record{
			Id:         "111",
			Created:    created.Unix(),
			Updated:    created.Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		&mdb.Record{
			Id:         "222",
			Created:    created.Add(20 * time.Hour).Unix(),
			Updated:    created.Add(20 * time.Hour).Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		2,
		0,
		false,
	}

	resp_read := &mdb.SearchResponse{}

	// to let elasticsearch finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 2 {
		t.Error("Not enough results")
		return
	}

	if resp_read.Records[0].Id != "111" {
		t.Error("Bad sorting 111")
		return
	}

	if resp_read.Records[1].Id != "222" {
		t.Error("Bad sorting 222")
		return
	}

}

func TestRecordSearchDescElastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	created := time.Now()
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		&mdb.Record{
			Id:         "111",
			Created:    created.Unix(),
			Updated:    created.Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		&mdb.Record{
			Id:         "222",
			Created:    created.Add(20 * time.Hour).Unix(),
			Updated:    created.Add(20 * time.Hour).Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		2,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let elasticsearch finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 2 {
		t.Skip("Not enough results")
		return
	}

	if resp_read.Records[0].Id != "222" {
		t.Skip("Bad sorting")
		return
	}

	if resp_read.Records[1].Id != "111" {
		t.Skip("Bad sorting")
		return
	}

}

func TestRecordSearchFromElastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Add(-time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		map[string]string{
			"name": "test_name",
		},
		time.Now().Add(-time.Minute).Unix(),
		0,
		2,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let elasticsearch finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}

	if len(resp_read.Records) != 1 {
		t.Skip("Bad results count", len(resp_read.Records))
		return
	}

	if resp_read.Records[0].Id != "222" {
		t.Skip("Bad from filter")
		return
	}

}

func TestRecordSearchToElastic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Add(-time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		map[string]string{
			"name": "test_name",
		},
		0,
		time.Now().Add(-time.Minute).Unix(),
		2,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let elasticsearch finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}

	if len(resp_read.Records) != 1 {
		t.Skip("Bad results count")
		return
	}

	if resp_read.Records[0].Id != "111" {
		t.Skip("Bad to filter")
		return
	}

}

func TestRecordSearchEmptyParamsElsatic(t *testing.T) {
	initDb(t, "elasticsearch", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Add(-time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		Database: &mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "elasticsearch"},
		Limit:    2,
		Offset:   0,
		Reverse:  true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let elasticsearch finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}

	if len(resp_read.Records) != 2 {
		t.Error("Bad results count")
		return
	}

	if resp_read.Records[0].Id != "222" {
		t.Error("Bad to filter")
		return
	}

}

// redis
func TestDbCreatedRedis(t *testing.T) {
	initDb(t, "redis", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "redis",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
}

func TestRecordCreatedRedis(t *testing.T) {
	initDb(t, "redis", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "redis",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "redis",
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name" {
		t.Error("Cannot retreave a record")
	}
	if resp_read.Record.Metadata["test_key"] != "test_value" {
		t.Error("Cannot retreave a Metadata record")
	}
}

func TestRecordUpdatedRedis(t *testing.T) {
	initDb(t, "redis", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "redis",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_upd := &mdb.UpdateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "redis",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name1",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp_upd := &mdb.UpdateResponse{}

	res_upd := hdlr.Update(ctx, req_upd, resp_upd)
	if res_upd != nil {
		t.Error(res_upd)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "redis",
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name1" {
		t.Error("Cannot retreave a record")
	}
}

func TestRecordDeletedRedis(t *testing.T) {
	initDb(t, "redis", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "redis",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_del := &mdb.DeleteRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "redis",
		},
		"111",
		"test_param3",
	}

	resp_del := &mdb.DeleteResponse{}

	res_del := hdlr.Delete(ctx, req_del, resp_del)
	if res_del != nil {
		t.Error(res_del)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "redis",
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
	}
}

func TestDbDeletedRedis(t *testing.T) {
	initDb(t, "redis", map[string]string{})
	ctx := common.NewTestContext(context.TODO())

	hdlr := new(DB)

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "redis",
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
	}
}

// arangodb
func TestDbCreatedArangodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
}

func TestRecordCreatedArangodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name" {
		t.Error("Cannot retreave a record")
	}
	if resp_read.Record.Id != "111" {
		t.Error("No ID match")
	}
	if resp_read.Record.Metadata["test_key"] != "test_value" {
		t.Error("Cannot retreave a Metadata record")
	}
}

func TestRecordUpdatedArangodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_upd := &mdb.UpdateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name1",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp_upd := &mdb.UpdateResponse{}

	res_upd := hdlr.Update(ctx, req_upd, resp_upd)
	if res_upd != nil {
		t.Error(res_upd)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name1" {
		t.Error("Cannot retreave a record")
	}
}

func TestRecordDeletedArangodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_del := &mdb.DeleteRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		"111",
		"test_param3",
	}

	resp_del := &mdb.DeleteResponse{}

	res_del := hdlr.Delete(ctx, req_del, resp_del)
	if res_del != nil {
		t.Error(res_del)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read == nil {
		t.Error("Must be not found")
	}
}

func TestRecordSearchArangodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		map[string]string{
			"test_key": "test_value",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch1Arangodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch2Arangodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		map[string]string{
			"name":       "test_name",
			"parameter1": "test_param",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch3Arangodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "arangodb",
		},
		map[string]string{
			"name": "test_name1",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 0 {
		t.Error("Non empty result")
		return
	}
}

func TestRecordSearchAscArangodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value1",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Add(time.Minute).Unix(),
			Updated:    time.Now().Add(time.Minute).Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value2",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		2,
		0,
		false,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 2 {
		t.Error("Not enough results")
		return
	}

	if resp_read.Records[0].Id != "111" {
		t.Error("Bad sorting")
		return
	}

	if resp_read.Records[1].Id != "222" {
		t.Error("Bad sorting")
		return
	}
}

func TestRecordSearchDescArangodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Add(time.Minute).Unix(),
			Updated:    time.Now().Add(time.Minute).Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		2,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 2 {
		t.Error("Not enough results")
		return
	}

	if resp_read.Records[0].Id != "222" {
		t.Error("Bad sorting")
		return
	}

	if resp_read.Records[1].Id != "111" {
		t.Error("Bad sorting")
		return
	}

}

func TestRecordSearchFromArnagodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Add(-time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Add(time.Minute).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		map[string]string{
			"name": "test_name",
		},
		time.Now().Unix(),
		0,
		2,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}

	if len(resp_read.Records) != 1 {
		t.Errorf("Bad results count %v", len(resp_read.Records))
		return
	}

	if resp_read.Records[0].Id != "222" {
		t.Error("Bad from filter")
		return
	}

}

func TestRecordSearchToArnagodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Add(-time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		map[string]string{
			"name": "test_name",
		},
		0,
		time.Now().Add(-time.Minute).Unix(),
		2,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}

	if len(resp_read.Records) != 1 {
		t.Error("Bad results count")
		return
	}

	if resp_read.Records[0].Id != "111" {
		t.Error("Bad to filter")
		return
	}

}

func TestRecordSearchEmptyParamsArnagodb(t *testing.T) {
	initDb(t, "arangodb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Add(-time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		Database: &mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "arangodb"},
		Limit:    2,
		Offset:   0,
		Reverse:  true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}

	if len(resp_read.Records) != 2 {
		t.Error("Bad results count")
		return
	}

	if resp_read.Records[0].Id != "222" {
		t.Error("Bad to filter")
		return
	}

}

// searchable
func TestDbCreatedSearchable(t *testing.T) {
	md := map[string]string{common.SearchableMeta: ""}
	initDb(t, "mysql", md)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Metadata: md,
		},
		&mdb.Record{
			Id:         "1112",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}
	hdlr := NewWrapper(new(DB))
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
}

func TestRecordCreatedSearchable(t *testing.T) {
	md := map[string]string{common.SearchableMeta: ""}
	initDb(t, "mysql", md)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Metadata: md,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := NewWrapper(new(DB))
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
	// so the routine is finished
	time.Sleep(time.Second)

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "elasticsearch",
			Metadata: md,
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name" {
		t.Error("Cannot retreave a record")
	}
}

func TestRecordUpdatedSearchable(t *testing.T) {
	md := map[string]string{common.SearchableMeta: ""}
	initDb(t, "mysql", md)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Metadata: md,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := NewWrapper(new(DB))
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	// so the routine is finished
	time.Sleep(time.Second)

	req_upd := &mdb.UpdateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Metadata: md,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name1",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp_upd := &mdb.UpdateResponse{}

	res_upd := hdlr.Update(ctx, req_upd, resp_upd)
	if res_upd != nil {
		t.Error(res_upd)
	}

	// so the routine is finished
	time.Sleep(time.Second)

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "elasticsearch",
			Metadata: md,
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name1" {
		t.Error("Cannot retreave a record")
	}
}

func TestRecordDeletedSearchable(t *testing.T) {
	md := map[string]string{common.SearchableMeta: ""}
	initDb(t, "mysql", md)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Metadata: md,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := NewWrapper(new(DB))
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_del := &mdb.DeleteRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Metadata: md,
		},
		"111",
		"test_param3",
	}

	resp_del := &mdb.DeleteResponse{}

	res_del := hdlr.Delete(ctx, req_del, resp_del)
	if res_del != nil {
		t.Error(res_del)
	}

	// so the routine is finished
	time.Sleep(time.Second)

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:  TestDBName,
			Table: TestDBTable,
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read == nil {
		t.Error(res_read)
	}
}

func TestRecordSearchSearchable(t *testing.T) {
	md := map[string]string{common.SearchableMeta: ""}
	initDb(t, "mysql", md)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Metadata: md,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name1",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := NewWrapper(new(DB))
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	// so the routine is finished
	time.Sleep(time.Second * 2)

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Metadata: md,
		},
		map[string]string{
			"name": "test_name1",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name1" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearchAutocompleteSearchable(t *testing.T) {
	md := common.SearchableAutocompleteMetaMap
	initDb(t, "mysql", md)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Metadata: md,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name1",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := NewWrapper(new(DB))
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	// so the routine is finished
	time.Sleep(time.Second * 2)

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Metadata: md,
		},
		map[string]string{
			"name": "test_name1",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name1" {
		t.Error("Empty result")
		return
	}
}

// arangodb graph
func TestDbCreatedArangodbGraph(t *testing.T) {
	initDb(t, "arangodb", common.GraphMap)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
}

func TestRecordCreatedArangodbGraph(t *testing.T) {
	initDb(t, "arangodb", common.GraphMap)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name" {
		t.Error("Cannot retreave a record")
	}
	if resp_read.Record.Id != "111" {
		t.Error("No ID match")
	}
	if resp_read.Record.Metadata["test_key"] != "test_value" {
		t.Error("Cannot retreave a Metadata record")
	}
}

func TestRecordUpdatedArangodbGraph(t *testing.T) {
	initDb(t, "arangodb", common.GraphMap)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_upd := &mdb.UpdateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name1",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp_upd := &mdb.UpdateResponse{}

	res_upd := hdlr.Update(ctx, req_upd, resp_upd)
	if res_upd != nil {
		t.Error(res_upd)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name1" {
		t.Error("Cannot retreave a record")
	}
}

func TestRecordDeletedArangodbGraph(t *testing.T) {
	initDb(t, "arangodb", common.GraphMap)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_del := &mdb.DeleteRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		"111",
		"test_param3",
	}

	resp_del := &mdb.DeleteResponse{}

	res_del := hdlr.Delete(ctx, req_del, resp_del)
	if res_del != nil {
		t.Error(res_del)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read == nil {
		t.Error("Must be not found")
	}
}

func TestRecordSearch2ArangodbGraph(t *testing.T) {
	initDb(t, "arangodb", common.GraphMap)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch3ArangodbGraph(t *testing.T) {
	initDb(t, "arangodb", common.GraphMap)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		map[string]string{
			"parameter1": "test_param",
			"parameter3": "test_param3",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch4ArangodbGraph(t *testing.T) {
	initDb(t, "arangodb", common.GraphMap)
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param12",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:     TestDBName,
			Table:    TestDBTable,
			Driver:   "arangodb",
			Metadata: common.GraphMap,
		},
		map[string]string{
			"name":       "test_name",
			"parameter1": "test_param",
			"parameter3": "test_param3",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 1 {
		t.Error("bad number of results")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

// influxdb
func TestDbCreatedInfluxdb(t *testing.T) {
	initDb(t, "influxdb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
}

func TestRecordCreatedInfluxdb(t *testing.T) {
	initDb(t, "influxdb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
	}
	if resp_read.Record.Name != "test_name" {
		t.Error("Cannot retreave a record")
	}
}

func TestRecordDeletedInfluxdb(t *testing.T) {
	initDb(t, "influxdb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_del := &mdb.DeleteRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		"111",
		"test_param3",
	}

	resp_del := &mdb.DeleteResponse{}

	res_del := hdlr.Delete(ctx, req_del, resp_del)
	if res_del != nil {
		t.Error(res_del)
	}

	req_read := &mdb.ReadRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		"111",
		"test_param3",
	}

	resp_read := &mdb.ReadResponse{
		&mdb.Record{},
	}

	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read == nil {
		t.Error(res_read)
	}
}

func TestRecordSearchInfluxdb(t *testing.T) {
	initDb(t, "influxdb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		map[string]string{
			"test_key": "test_value",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let influxdb finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch1Influxdb(t *testing.T) {
	initDb(t, "influxdb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param2",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		1,
		0,
		true,
	}
	// to let influxdb finish indexing
	time.Sleep(time.Second * 2)
	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch2Influxdb(t *testing.T) {
	initDb(t, "influxdb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		map[string]string{
			"name":       "test_name",
			"parameter1": "test_param",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let influxdb finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) == 0 {
		t.Error("Empty result")
		return
	}
	if resp_read.Records[0].Name != "test_name" {
		t.Error("Empty result")
		return
	}
}

func TestRecordSearch3Influxdb(t *testing.T) {
	initDb(t, "influxdb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{
			Name:   TestDBName,
			Table:  TestDBTable,
			Driver: "influxdb",
		},
		map[string]string{
			"name": "test_name1",
		},
		0,
		0,
		1,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let influxdb finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 0 {
		t.Error("Non empty result")
		return
	}
}

func TestRecordSearchAscInfluxdb(t *testing.T) {
	initDb(t, "influxdb", map[string]string{})
	created := time.Now()
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		&mdb.Record{
			Id:         "111",
			Created:    created.Add(-40 * time.Hour).Unix(),
			Updated:    created.Add(-40 * time.Hour).Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		&mdb.Record{
			Id:         "222",
			Created:    created.Add(-20 * time.Hour).Unix(),
			Updated:    created.Add(-20 * time.Hour).Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key1": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		2,
		0,
		false,
	}

	resp_read := &mdb.SearchResponse{}

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 2 {
		t.Errorf("Not enough results, %v", len(resp_read.Records))
		return
	}

	if resp_read.Records[0].Id != "111" {
		t.Error("Bad sorting 111")
		return
	}

	if resp_read.Records[1].Id != "222" {
		t.Error("Bad sorting 222")
		return
	}

}

func TestRecordSearchDescInfluxdb(t *testing.T) {
	initDb(t, "influxdb", map[string]string{})
	created := time.Now()
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		&mdb.Record{
			Id:         "111",
			Created:    created.Add(-40 * time.Hour).Unix(),
			Updated:    created.Add(-40 * time.Hour).Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		&mdb.Record{
			Id:         "222",
			Created:    created.Add(-20 * time.Hour).Unix(),
			Updated:    created.Add(-20 * time.Hour).Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		map[string]string{
			"name": "test_name",
		},
		0,
		0,
		2,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let influxdb finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}
	if len(resp_read.Records) != 2 {
		t.Skip("Not enough results")
		return
	}

	if resp_read.Records[0].Id != "222" {
		t.Skip("Bad sorting")
		return
	}

	if resp_read.Records[1].Id != "111" {
		t.Skip("Bad sorting")
		return
	}

}

func TestRecordSearchFromInfluxdb(t *testing.T) {
	initDb(t, "influxdb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Add(-20 * time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Add(-10 * time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		map[string]string{
			"name": "test_name",
		},
		time.Now().Add(-15 * time.Hour).Unix(),
		0,
		2,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let influxdb finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}

	if len(resp_read.Records) != 1 {
		t.Skipf("Bad results count %v", len(resp_read.Records))
		return
	}

	if resp_read.Records[0].Id != "222" {
		t.Skip("Bad from filter")
		return
	}

}

func TestRecordSearchToInfluxdb(t *testing.T) {
	initDb(t, "influxdb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Add(-20 * time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Add(-10 * time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		map[string]string{
			"name": "test_name",
		},
		0,
		time.Now().Add(-15 * time.Hour).Unix(),
		2,
		0,
		true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let influxdb finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}

	if len(resp_read.Records) != 1 {
		t.Skip("Bad results count")
		return
	}

	if resp_read.Records[0].Id != "111" {
		t.Skip("Bad to filter")
		return
	}

}

func TestRecordSearchEmptyParamsInfluxdb(t *testing.T) {
	initDb(t, "influxdb", map[string]string{})
	ctx := common.NewTestContext(context.TODO())
	req := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		&mdb.Record{
			Id:         "111",
			Created:    time.Now().Add(-time.Hour).Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp := &mdb.CreateResponse{}

	hdlr := new(DB)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req1 := &mdb.CreateRequest{
		&mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		&mdb.Record{
			Id:         "222",
			Created:    time.Now().Unix(),
			Updated:    time.Now().Unix(),
			Name:       "test_name",
			Parameter1: "test_param",
			Parameter2: "test_param1",
			Parameter3: "test_param3",
			Metadata: map[string]string{
				"test_key": "test_value",
			},
		},
	}

	resp1 := &mdb.CreateResponse{}

	res1 := hdlr.Create(ctx, req1, resp1)
	if res1 != nil {
		t.Error(res1)
	}

	req_read := &mdb.SearchRequest{
		Database: &mdb.Database{Name: TestDBName, Table: TestDBTable, Driver: "influxdb"},
		Limit:    2,
		Offset:   0,
		Reverse:  true,
	}

	resp_read := &mdb.SearchResponse{}

	// to let influxdb finish indexing
	time.Sleep(time.Second * 2)

	res_read := hdlr.Search(ctx, req_read, resp_read)
	if res != nil {
		t.Error(res_read)
		return
	}

	if len(resp_read.Records) != 2 {
		t.Error("Bad results count")
		return
	}

	if resp_read.Records[0].Id != "222" {
		t.Error("Bad to filter")
		return
	}
}
