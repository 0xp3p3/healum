package handler

import (
	"log"
	"server/common"
	"server/db-srv/db"
	mdb "server/db-srv/proto/db"

	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type DB struct{}

func validateDB(method string, d *mdb.Database) error {
	if d == nil {
		return errors.BadRequest("go.micro.srv.db."+method, "invalid database")
	}

	if len(d.Name) == 0 {
		return errors.BadRequest("go.micro.srv.db."+method, "database is blank")
	}
	if len(d.Table) == 0 {
		return errors.BadRequest("go.micro.srv.db."+method, "table is blank")
	}

	// TODO: check exists

	return nil
}

// InitDb initializes healum databases
func (d *DB) InitDb(ctx context.Context, req *mdb.InitDbRequest, rsp *mdb.InitDbResponse) error {
	if err := d.CreateDatabase(context.TODO(), &mdb.CreateDatabaseRequest{
		&mdb.Database{
			Name:   common.DbHealumName,
			Driver: "arangodb",
		},
	}, &mdb.CreateDatabaseResponse{}); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// RemoveDb removes healum database (for testing)
func (d *DB) RemoveDb(ctx context.Context, req *mdb.RemoveDbRequest, rsp *mdb.RemoveDbResponse) error {
	for _, cols := range common.DbHealum {
		if len(cols) == 0 {
			continue
		}

		if err := d.DeleteDatabase(ctx, &mdb.DeleteDatabaseRequest{
			Database: &mdb.Database{
				Name:   common.DbHealumName,
				Driver: common.DbHealumDriver,
			},
		}, &mdb.DeleteDatabaseResponse{}); err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) Read(ctx context.Context, req *mdb.ReadRequest, rsp *mdb.ReadResponse) error {
	if err := validateDB("DB.Read", req.Database); err != nil {
		return err
	}

	if len(req.Id) == 0 {
		return errors.BadRequest("go.micro.srv.db.DB.Read", "invalid id")
	}
	r, err := db.Read(req.Database, req.Id, req.Parameter3)
	if err != nil && err == db.ErrNotFound {
		return errors.NotFound("go.micro.srv.db.DB.Read", "not found")
	} else if err != nil {
		return errors.InternalServerError("go.micro.srv.db.DB.Read", err.Error())
	}

	rsp.Record = r

	return nil
}

func (d *DB) Create(ctx context.Context, req *mdb.CreateRequest, rsp *mdb.CreateResponse) error {

	if req.Record == nil {
		return errors.BadRequest("go.micro.srv.db.DB.Create", "invalid record")
	}

	if err := validateDB("DB.Create", req.Database); err != nil {
		return err
	}

	if len(req.Record.Id) == 0 {
		return errors.BadRequest("go.micro.srv.db.DB.Create", "invalid id")
	}
	if err := db.Create(req.Database, req.Record); err != nil {
		return errors.InternalServerError("go.micro.srv.db.DB.Create", err.Error())
	}

	return nil
}

func (d *DB) Update(ctx context.Context, req *mdb.UpdateRequest, rsp *mdb.UpdateResponse) error {
	if req.Record == nil {
		return errors.BadRequest("go.micro.srv.db.DB.Update", "invalid record")
	}

	if err := validateDB("DB.Update", req.Database); err != nil {
		return err
	}

	if len(req.Record.Id) == 0 {
		return errors.BadRequest("go.micro.srv.db.DB.Update", "invalid id")
	}

	if err := db.Update(req.Database, req.Record); err != nil && err == db.ErrNotFound {
		return errors.NotFound("go.micro.srv.db.DB.Update", "not found")
	} else if err != nil {
		return errors.InternalServerError("go.micro.srv.db.DB.Update", err.Error())
	}

	return nil
}

func (d *DB) Delete(ctx context.Context, req *mdb.DeleteRequest, rsp *mdb.DeleteResponse) error {
	if err := validateDB("DB.Delete", req.Database); err != nil {
		common.ErrorLog(common.DbSrv, d.Delete, err, "DB is invalid")
		return err
	}

	if len(req.Id) == 0 {
		common.ErrorLog(common.DbSrv, d.Delete, nil, "Id is invalid")
		return errors.BadRequest("go.micro.srv.db.DB.Delete", "invalid id")
	}

	if err := db.Delete(req.Database, req.Id, req.Parameter3); err != nil && err == db.ErrNotFound {
		return nil
	} else if err != nil {
		common.ErrorLog(common.DbSrv, d.Delete, err, "Delete query is failed")
		return errors.InternalServerError("go.micro.srv.db.DB.Delete", err.Error())
	}

	return nil
}

func (d *DB) Search(ctx context.Context, req *mdb.SearchRequest, rsp *mdb.SearchResponse) error {
	if err := validateDB("DB.Search", req.Database); err != nil {
		common.ErrorLog(common.DbSrv, d.Search, err, "DB is invalid")
		return err
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	if req.Offset < 0 {
		req.Offset = 0
	}

	r, err := db.Search(req.Database, req.Metadata, req.From, req.To, req.Limit, req.Offset, req.Reverse)
	if err != nil {
		common.ErrorLog(common.DbSrv, d.Search, err, "Search query is failed")
		return errors.InternalServerError("go.micro.srv.db.DB.Search", err.Error())
	}
	rsp.Records = r

	return nil
}

func (d *DB) RunQuery(ctx context.Context, req *mdb.RunQueryRequest, rsp *mdb.RunQueryResponse) error {
	if err := validateDB("DB.RunQuery", req.Database); err != nil {
		common.ErrorLog(common.DbSrv, d.RunQuery, err, "DB is invalid")
		return err
	}

	r, err := db.RunQuery(req.Database, req.Query)
	if err != nil {
		common.ErrorLog(common.DbSrv, d.RunQuery, err, "RunQuery is failed")
		return errors.InternalServerError("go.micro.srv.db.DB.RunQuery", err.Error())
	}
	rsp.Records = r

	return nil
}

func (d *DB) CreateDatabase(ctx context.Context, req *mdb.CreateDatabaseRequest, rsp *mdb.CreateDatabaseResponse) error {
	// if err := validateDB("DB.CreateDatabase", req.Database); err != nil {
	// 	common.ErrorLog(common.DbSrv, d.CreateDatabase, err, "DB is invalid")
	// 	return err
	// }

	if err := db.CreateDatabase(req.Database); err != nil {
		common.ErrorLog(common.DbSrv, d.CreateDatabase, err, "CreateDatbase is failed")
		return errors.InternalServerError("go.micro.srv.db.DB.CreateDatabase", err.Error())
	}
	return nil
}

func (d *DB) DeleteDatabase(ctx context.Context, req *mdb.DeleteDatabaseRequest, rsp *mdb.DeleteDatabaseResponse) error {
	// if err := validateDB("DB.DeleteDatabase", req.Database); err != nil {
	// 	common.ErrorLog(common.DbSrv, d.DeleteDatabase, err, "DB is invalid")
	// 	return err
	// }
	if !common.IsTestContext(ctx) {
		common.ErrorLog(common.DbSrv, d.DeleteDatabase, nil, "Context is invalid")
		return errors.BadRequest("go.micro.srv.db.DB.DeleteDatabase", "DeleteDatabase available only for testing contexts")
	}

	if err := db.DeleteDatabase(req.Database); err != nil {
		common.ErrorLog(common.DbSrv, d.DeleteDatabase, err, "DeleteDatabase is failed")
		return errors.InternalServerError("go.micro.srv.db.DB.DeleteDatabase", err.Error())
	}

	return nil
}
