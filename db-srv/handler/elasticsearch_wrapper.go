package handler

import (
	mdb "server/db-srv/proto/db"

	"github.com/jinzhu/copier"

	"server/common"

	"golang.org/x/net/context"
)

const (
	elasticDriver = "elasticsearch"
)

// Wraps the handler with elasticsearch
type ElasticSearchWrapper struct {
	wrappedDb *DB
}

func NewWrapper(db *DB) *ElasticSearchWrapper {
	return &ElasticSearchWrapper{
		wrappedDb: db,
	}
}

//Checks if the data object is searchable in elastic
func isSearchable(db *mdb.Database) bool {
	// no extra elastic index for the data that is already in elastics
	if db.Driver == elasticDriver {
		return false
	}
	_, hasSearchable := db.Metadata[common.SearchableMeta]
	return hasSearchable
}

//Checks if the data object is searchable in elastic (autocomplete request)
func isAutocomlete(db *mdb.Database) bool {
	_, hasSearchable := db.Metadata[common.SearchableAutocompleteMeta]
	return hasSearchable
}

// Init initialize handled functions
func (d *ElasticSearchWrapper) InitDb(ctx context.Context, req *mdb.InitDbRequest, rsp *mdb.InitDbResponse) error {
	return d.wrappedDb.InitDb(ctx, req, rsp)
}

// RemoveDb remove handled functions
func (d *ElasticSearchWrapper) RemoveDb(ctx context.Context, req *mdb.RemoveDbRequest, rsp *mdb.RemoveDbResponse) error {
	return d.wrappedDb.RemoveDb(ctx, req, rsp)
}

// Wrapped handled functions
func (d *ElasticSearchWrapper) Read(ctx context.Context, req *mdb.ReadRequest, rsp *mdb.ReadResponse) error {
	return d.wrappedDb.Read(ctx, req, rsp)
}

func (d *ElasticSearchWrapper) Create(ctx context.Context, req *mdb.CreateRequest, rsp *mdb.CreateResponse) error {

	if isSearchable(req.Database) {
		// async index update
		go func() {
			reqCopy := new(mdb.CreateRequest)
			copier.Copy(reqCopy, req)
			reqCopy.Database = new(mdb.Database)
			copier.Copy(reqCopy.Database, req.Database)
			reqCopy.Database.Driver = elasticDriver
			rspCopy := new(mdb.CreateResponse)
			d.wrappedDb.Create(ctx, reqCopy, rspCopy)
		}()
	}
	return d.wrappedDb.Create(ctx, req, rsp)
}

func (d *ElasticSearchWrapper) Update(ctx context.Context, req *mdb.UpdateRequest, rsp *mdb.UpdateResponse) error {
	if isSearchable(req.Database) {
		// async index update
		go func() {
			reqCopy := new(mdb.UpdateRequest)
			copier.Copy(reqCopy, req)
			reqCopy.Database = new(mdb.Database)
			copier.Copy(reqCopy.Database, req.Database)
			reqCopy.Database.Driver = elasticDriver
			rspCopy := new(mdb.UpdateResponse)
			d.wrappedDb.Update(ctx, reqCopy, rspCopy)
		}()
	}

	return d.wrappedDb.Update(ctx, req, rsp)
}

func (d *ElasticSearchWrapper) Delete(ctx context.Context, req *mdb.DeleteRequest, rsp *mdb.DeleteResponse) error {
	if isSearchable(req.Database) {
		// async index update
		go func() {
			reqCopy := new(mdb.DeleteRequest)
			copier.Copy(reqCopy, req)
			reqCopy.Database = new(mdb.Database)
			copier.Copy(reqCopy.Database, req.Database)
			reqCopy.Database.Driver = elasticDriver
			rspCopy := new(mdb.DeleteResponse)
			d.wrappedDb.Delete(ctx, reqCopy, rspCopy)
		}()
	}

	return d.wrappedDb.Delete(ctx, req, rsp)
}

func (d *ElasticSearchWrapper) Search(ctx context.Context, req *mdb.SearchRequest, rsp *mdb.SearchResponse) error {
	// search in elastic
	if isSearchable(req.Database) {
		// async index update
		reqCopy := new(mdb.SearchRequest)
		copier.Copy(reqCopy, req)
		reqCopy.Database = new(mdb.Database)
		copier.Copy(reqCopy.Database, req.Database)
		reqCopy.Database.Driver = elasticDriver
		if isAutocomlete(req.Database) {
			// adds autocomplete flag
			reqCopy.Metadata[common.SearchableAutocompleteMeta] = ""
		}
		return d.wrappedDb.Search(ctx, reqCopy, rsp)
	}

	return d.wrappedDb.Search(ctx, req, rsp)
}

func (d *ElasticSearchWrapper) RunQuery(ctx context.Context, req *mdb.RunQueryRequest, rsp *mdb.RunQueryResponse) error {
	if isSearchable(req.Database) {
		// async index update
		go func() {
			reqCopy := new(mdb.RunQueryRequest)
			copier.Copy(reqCopy, req)
			reqCopy.Database = new(mdb.Database)
			copier.Copy(reqCopy.Database, req.Database)
			reqCopy.Database.Driver = elasticDriver
			rspCopy := new(mdb.RunQueryResponse)
			d.wrappedDb.RunQuery(ctx, reqCopy, rspCopy)
		}()
	}

	return d.wrappedDb.RunQuery(ctx, req, rsp)
}

func (d *ElasticSearchWrapper) CreateDatabase(ctx context.Context, req *mdb.CreateDatabaseRequest, rsp *mdb.CreateDatabaseResponse) error {
	if isSearchable(req.Database) {
		reqCopy := new(mdb.CreateDatabaseRequest)
		copier.Copy(reqCopy, req)
		reqCopy.Database = new(mdb.Database)
		copier.Copy(reqCopy.Database, req.Database)
		reqCopy.Database.Driver = elasticDriver
		rspCopy := new(mdb.CreateDatabaseResponse)
		d.wrappedDb.CreateDatabase(ctx, reqCopy, rspCopy)
	}
	return d.wrappedDb.CreateDatabase(ctx, req, rsp)
}

func (d *ElasticSearchWrapper) DeleteDatabase(ctx context.Context, req *mdb.DeleteDatabaseRequest, rsp *mdb.DeleteDatabaseResponse) error {
	if isSearchable(req.Database) {
		reqCopy := new(mdb.DeleteDatabaseRequest)
		copier.Copy(reqCopy, req)
		reqCopy.Database = new(mdb.Database)
		copier.Copy(reqCopy.Database, req.Database)
		reqCopy.Database.Driver = elasticDriver
		rspCopy := new(mdb.DeleteDatabaseResponse)
		d.wrappedDb.DeleteDatabase(ctx, reqCopy, rspCopy)
	}
	return d.wrappedDb.DeleteDatabase(ctx, req, rsp)
}
