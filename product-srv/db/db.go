package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"server/common"
	db_proto "server/db-srv/proto/db"
	product_proto "server/product-srv/proto/product"
	static_proto "server/static-srv/proto/static"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/client"
	log "github.com/sirupsen/logrus"
)

type clientWrapper struct {
	Db_client db_proto.DBClient
}

var (
	ClientWrapper *clientWrapper
	ErrNotFound   = errors.New("not found")
)

// Storage for a db microservice client
func NewClientWrapper(serviceClient client.Client) *clientWrapper {
	cl := db_proto.NewDBClient("", serviceClient)

	return &clientWrapper{
		Db_client: cl,
	}
}

// Init initializes healum databases
func Init(serviceClient client.Client) error {
	ClientWrapper = NewClientWrapper(serviceClient)
	// if _, err := ClientWrapper.Db_client.Init(context.TODO(), &db_proto.InitRequest{}); err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	return nil
}

// RemoveDb removes healum database (for testing)
func RemoveDb(ctx context.Context, serviceClient client.Client) error {
	ClientWrapper = NewClientWrapper(serviceClient)
	if _, err := ClientWrapper.Db_client.RemoveDb(ctx, &db_proto.RemoveDbRequest{}); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func productToRecord(product *product_proto.Product) (string, error) {
	data, err := common.MarhalToObject(product)
	if err != nil {
		return "", err
	}

	common.FilterObject(data, "createdBy", product.CreatedBy)
	var createdById string
	if product.CreatedBy != nil {
		createdById = product.CreatedBy.Id
	}
	if len(product.Owners) > 0 {
		var arr []interface{}
		for _, item := range product.Owners {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["owners"] = arr
	} else {
		delete(data, "owners")
	}

	d := map[string]interface{}{
		"_key":       product.Id,
		"id":         product.Id,
		"created":    product.Created,
		"updated":    product.Updated,
		"name":       product.Name,
		"parameter1": product.OrgId,
		"parameter2": createdById,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToProduct(r *db_proto.Record) (*product_proto.Product, error) {
	var p product_proto.Product
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func serviceToRecord(service *product_proto.Service) (string, error) {
	data, err := common.MarhalToObject(service)
	if err != nil {
		return "", err
	}

	if len(service.Batches) > 0 {
		var arr []interface{}
		for _, item := range service.Batches {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["batches"] = arr
	} else {
		delete(data, "batches")
	}

	d := map[string]interface{}{
		"_key":       service.Id,
		"id":         service.Id,
		"created":    service.Created,
		"updated":    service.Updated,
		"name":       service.Name,
		"parameter1": service.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToService(r *db_proto.Record) (*product_proto.Service, error) {
	var p product_proto.Service
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func batchToRecord(batch *static_proto.Batch) (string, error) {
	data, err := common.MarhalToObject(batch)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":       batch.Id,
		"id":         batch.Id,
		"created":    batch.Created,
		"updated":    batch.Updated,
		"name":       batch.Name,
		"parameter1": batch.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToBatch(r *db_proto.Record) (*static_proto.Batch, error) {
	var p static_proto.Batch
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func runQuery(ctx context.Context, q string, table string) (*db_proto.RunQueryResponse, error) {
	return ClientWrapper.Db_client.RunQuery(ctx, &db_proto.RunQueryRequest{
		Database: &db_proto.Database{
			Name:     common.DbHealumName,
			Table:    table,
			Driver:   common.DbHealumDriver,
			Metadata: common.SearchableMetaMap,
		},
		Query: q,
	})
}

func queryProductMerge() string {
	query := fmt.Sprintf(`
		LET owners = (FILTER NOT_NULL(doc.data.owners) FOR p IN doc.data.owners FOR owner in %v FILTER p.id == owner._key RETURN owner.data)
		LET createdBy = (FOR user in %v FILTER doc.data.createdBy.id == user._key RETURN user.data)
		RETURN MERGE_RECURSIVE(doc,{data:{owners:owners,createdBy:createdBy[0]}})`, common.DbUserTable, common.DbUserTable)
	return query
}

func queryServiceMerge() string {
	query := fmt.Sprintf(`
		LET batches = (FILTER NOT_NULL(doc.data.batches) FOR p IN doc.data.batches FOR batch in %v FILTER p.id == batch._key RETURN batch.data)
		RETURN MERGE_RECURSIVE(doc,{data:{batches:batches}})`, common.DbBatchTable)
	return query
}

// AllProducts get all products
func AllProducts(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*product_proto.Product, error) {
	var products []*product_proto.Product
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)
	merge_query := queryProductMerge()

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbProductTable, query, sort_query, limit_query, merge_query)

	resp, err := runQuery(ctx, q, common.DbProductTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if product, err := recordToProduct(r); err == nil {
			products = append(products, product)
		}
	}
	return products, nil
}

// CreateProduct creates a product
func CreateProduct(ctx context.Context, product *product_proto.Product) error {
	if product.Created == 0 {
		product.Created = time.Now().Unix()
	}
	product.Updated = time.Now().Unix()
	record, err := productToRecord(product)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, product.Id, record, record, common.DbProductTable)
	_, err = runQuery(ctx, q, common.DbProductTable)
	return err
}

// ReadProduct reads a product by ID
func ReadProduct(ctx context.Context, id, orgId, teamId string) (*product_proto.Product, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	merge_query := queryProductMerge()

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s`, common.DbProductTable, query, merge_query)

	resp, err := runQuery(ctx, q, common.DbProductTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToProduct(resp.Records[0])
	return data, err
}

// DeleteProduct deletes a product by ID
func DeleteProduct(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbProductTable, query, common.DbProductTable)
	_, err := runQuery(ctx, q, common.DbProductTable)
	return err
}

// AutocompleteProduct autocomplete search products
func AutocompleteProduct(ctx context.Context, title string) ([]*product_proto.Product, error) {
	var products []*product_proto.Product
	// limit_query := common.QueryPaginate(offset, limit)
	// sort_query := common.QuerySort(sortParameter, sortDirection)
	merge_query := queryProductMerge()

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER LIKE(doc.name, "%v",true)
		%s`, common.DbProductTable, `%`+title+`%`, merge_query)

	resp, err := runQuery(ctx, q, common.DbProductTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if product, err := recordToProduct(r); err == nil {
			products = append(products, product)
		}
	}
	return products, nil
}

// AllServices get all services
func AllServices(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*product_proto.Service, error) {
	var services []*product_proto.Service
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)
	merge_query := queryServiceMerge()

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbServiceTable, query, sort_query, limit_query, merge_query)

	resp, err := runQuery(ctx, q, common.DbServiceTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if service, err := recordToService(r); err == nil {
			services = append(services, service)
		}
	}
	return services, nil
}

// CreateService creates a service
func CreateService(ctx context.Context, service *product_proto.Service) error {
	if service.Created == 0 {
		service.Created = time.Now().Unix()
	}
	service.Updated = time.Now().Unix()

	record, err := serviceToRecord(service)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, service.Id, record, record, common.DbServiceTable)
	_, err = runQuery(ctx, q, common.DbServiceTable)
	return err
}

// ReadService reads a service by ID
func ReadService(ctx context.Context, id, orgId, teamId string) (*product_proto.Service, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	merge_query := queryServiceMerge()

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s`, common.DbServiceTable, query, merge_query)

	resp, err := runQuery(ctx, q, common.DbServiceTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToService(resp.Records[0])
	return data, err
}

// DeleteService deletes a service by ID
func DeleteService(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbServiceTable, query, common.DbServiceTable)
	_, err := runQuery(ctx, q, common.DbServiceTable)
	return err
}

// AutocompleteService autocomplete search services
func AutocompleteService(ctx context.Context, title string) ([]*product_proto.Service, error) {
	var services []*product_proto.Service
	// limit_query := common.QueryPaginate(offset, limit)
	// sort_query := common.QuerySort(sortParameter, sortDirection)
	merge_query := queryServiceMerge()

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER LIKE(doc.name, "%v",true)
		%s`, common.DbServiceTable, `%`+title+`%`, merge_query)

	resp, err := runQuery(ctx, q, common.DbServiceTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if service, err := recordToService(r); err == nil {
			services = append(services, service)
		}
	}
	return services, nil
}

// AllBatches get all batches
func AllBatches(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.Batch, error) {
	var batches []*static_proto.Batch
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		RETURN doc`, common.DbBatchTable, query, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbBatchTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if batch, err := recordToBatch(r); err == nil {
			batches = append(batches, batch)
		}
	}
	return batches, nil
}

// CreateBatch creates a batch
func CreateBatch(ctx context.Context, batch *static_proto.Batch) error {
	if batch.Created == 0 {
		batch.Created = time.Now().Unix()
	}
	batch.Updated = time.Now().Unix()

	record, err := batchToRecord(batch)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, batch.Id, record, record, common.DbBatchTable)
	_, err = runQuery(ctx, q, common.DbBatchTable)
	return err
}

// ReadBatch reads a batch by ID
func ReadBatch(ctx context.Context, id, orgId, teamId string) (*static_proto.Batch, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbBatchTable, query)

	resp, err := runQuery(ctx, q, common.DbBatchTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToBatch(resp.Records[0])
	return data, err
}

// DeleteBatch deletes a batch by ID
func DeleteBatch(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbBatchTable, query, common.DbBatchTable)
	_, err := runQuery(ctx, q, common.DbBatchTable)
	return err
}
