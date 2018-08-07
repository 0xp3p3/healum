package handler

import (
	"context"
	account_proto "server/account-srv/proto/account"
	"server/common"
	organisation_proto "server/organisation-srv/proto/organisation"
	"server/product-srv/db"
	product_proto "server/product-srv/proto/product"
	static_proto "server/static-srv/proto/static"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

var cl = client.NewClient(
	client.Transport(nats_transport.NewTransport()),
	client.Broker(nats_broker.NewBroker()),
	client.RequestTimeout(4*time.Second),
	client.Retries(5),
)

var product = &product_proto.Product{
	Name:  "product",
	OrgId: "orgid",
}

var service = &product_proto.Service{
	Name:  "service",
	OrgId: "orgid",
}

var batch = &static_proto.Batch{
	Name:  "batch",
	OrgId: "orgid",
}

var user = &user_proto.User{
	OrgId:     "orgid",
	Firstname: "David",
	Lastname:  "John",
	AvatarUrl: "http://example.com",
	Tokens: []*user_proto.Token{
		{"11671c2e7da30e3c393813f60b327f9c2e2e08390761aa01e37ba5d3e6a617be", 1, "aaa"}, {"token_b", 2, "bbb"},
	},
}

var role = &static_proto.Role{
	OrgId: "orgid",
	Name:  "own",
}

var org1 = &organisation_proto.Organisation{
	Type: organisation_proto.OrganisationType_ROOT,
}

var account = &account_proto.Account{
	Email:    "email" + common.Random(4) + "@email.com",
	Password: "pass1",
}

func initDb() {
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func createProduct(ctx context.Context, hdlr *ProductService, t *testing.T) *product_proto.Product {
	// create user
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	user.Id = ""
	account.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user, Account: account})
	if err != nil {
		t.Error(err)
		return nil
	}

	product.CreatedBy = rsp_org.Data.User
	product.Owners = []*user_proto.User{rsp_org.Data.User}
	product.OrgId = rsp_org.Data.Organisation.Id

	req_create := &product_proto.CreateProductRequest{Product: product}
	resp_create := &product_proto.CreateProductResponse{}
	if err := hdlr.CreateProduct(ctx, req_create, resp_create); err != nil {
		t.Error(err)
		return nil
	}

	return resp_create.Data.Product
}

func createService(ctx context.Context, hdlr *ProductService, t *testing.T) *product_proto.Service {
	batch := createBatch(ctx, hdlr, t)
	if batch == nil {
		return nil
	}
	service.Batches = []*static_proto.Batch{batch}

	req_create := &product_proto.CreateServiceRequest{Service: service}
	resp_create := &product_proto.CreateServiceResponse{}
	err := hdlr.CreateService(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}

	return resp_create.Data.Service
}

func createBatch(ctx context.Context, hdlr *ProductService, t *testing.T) *static_proto.Batch {
	req_create := &product_proto.CreateBatchRequest{Batch: batch}
	resp_create := &product_proto.CreateBatchResponse{}
	err := hdlr.CreateBatch(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}

	return resp_create.Data.Batch
}

func TestAllProducts(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ProductService)

	product := createProduct(ctx, hdlr, t)
	if product == nil {
		return
	}

	req_all := &product_proto.AllProductsRequest{}
	resp_all := &product_proto.AllProductsResponse{}
	err := hdlr.AllProducts(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Products) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Products[0].Id != product.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadProduct(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ProductService)

	product := createProduct(ctx, hdlr, t)
	if product == nil {
		return
	}

	req_read := &product_proto.ReadProductRequest{Id: product.Id}
	resp_read := &product_proto.ReadProductResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadProduct(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Product == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Product.Id != product.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteProduct(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ProductService)

	product := createProduct(ctx, hdlr, t)
	if product == nil {
		return
	}

	req_del := &product_proto.DeleteProductRequest{Id: product.Id}
	resp_del := &product_proto.DeleteProductResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteProduct(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &product_proto.ReadProductRequest{Id: product.Id}
	resp_read := &product_proto.ReadProductResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadProduct(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAutocompleteProduct(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ProductService)

	product := createProduct(ctx, hdlr, t)
	if product == nil {
		return
	}

	req_all := &product_proto.AutocompleteProductRequest{
		Title: "pro",
	}
	resp_all := &product_proto.AutocompleteProductResponse{}
	err := hdlr.AutocompleteProduct(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Products) == 0 {
		t.Error("Object count does not match")
		return
	}
}

func TestAllServices(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ProductService)

	service := createService(ctx, hdlr, t)
	if service == nil {
		return
	}

	req_all := &product_proto.AllServicesRequest{}
	resp_all := &product_proto.AllServicesResponse{}
	err := hdlr.AllServices(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Services) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Services[0].Id != service.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadService(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ProductService)

	service := createService(ctx, hdlr, t)
	if service == nil {
		return
	}

	req_read := &product_proto.ReadServiceRequest{Id: service.Id}
	resp_read := &product_proto.ReadServiceResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadService(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Service == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Service.Id != service.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteService(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ProductService)

	service := createService(ctx, hdlr, t)
	if service == nil {
		return
	}

	req_del := &product_proto.DeleteServiceRequest{Id: service.Id}
	resp_del := &product_proto.DeleteServiceResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteService(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
		return
	}

	req_read := &product_proto.ReadServiceRequest{Id: service.Id}
	resp_read := &product_proto.ReadServiceResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadService(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAutocompleteService(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ProductService)

	service := createService(ctx, hdlr, t)
	if service == nil {
		return
	}

	req_all := &product_proto.AutocompleteServiceRequest{
		Title: "ser",
	}
	resp_all := &product_proto.AutocompleteServiceResponse{}
	err := hdlr.AutocompleteService(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Services) == 0 {
		t.Error("Object count does not match")
		return
	}
}

func TestAllBatches(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ProductService)

	batch := createBatch(ctx, hdlr, t)
	if batch == nil {
		return
	}

	req_all := &product_proto.AllBatchesRequest{}
	resp_all := &product_proto.AllBatchesResponse{}
	err := hdlr.AllBatches(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Batches) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Batches[0].Id != batch.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadBatch(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ProductService)

	batch := createBatch(ctx, hdlr, t)
	if batch == nil {
		return
	}

	req_read := &product_proto.ReadBatchRequest{Id: batch.Id}
	resp_read := &product_proto.ReadBatchResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadBatch(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Batch == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Batch.Id != batch.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteBatch(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ProductService)

	batch := createBatch(ctx, hdlr, t)
	if batch == nil {
		return
	}

	req_del := &product_proto.DeleteBatchRequest{Id: batch.Id}
	resp_del := &product_proto.DeleteBatchResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteBatch(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
		return
	}

	req_read := &product_proto.ReadBatchRequest{Id: batch.Id}
	resp_read := &product_proto.ReadBatchResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadBatch(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}
