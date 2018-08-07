package api

import (
	"context"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	product_proto "server/product-srv/proto/product"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type ProductService struct {
	ProductClient product_proto.ProductServiceClient
	Auth          Filters
	Audit         AuditFilter
	ServerMetrics metrics.Metrics
}

func (p ProductService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/products")

	audit := &audit_proto.Audit{
		ActionService:  common.ProductSrv,
		ActionResource: common.BASE + common.PRODUCT_TYPE,
	}

	ws.Route(ws.GET("/products/all").To(p.AllProducts).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all products"))

	ws.Route(ws.POST("/product").To(p.CreateProduct).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a product"))

	ws.Route(ws.GET("/product/{product_id}").To(p.ReadProduct).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a product"))

	ws.Route(ws.DELETE("/product/{product_id}").To(p.DeleteProduct).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a product"))

	ws.Route(ws.POST("/autocomplete/product").To(p.AutocompleteProduct).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Autocomplete search products"))

	ws.Route(ws.GET("/services/all").To(p.AllServices).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all services"))

	ws.Route(ws.POST("/service").To(p.CreateService).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a service"))

	ws.Route(ws.GET("/service/{service_id}").To(p.ReadService).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a service"))

	ws.Route(ws.DELETE("/service/{service_id}").To(p.DeleteService).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a service"))

	ws.Route(ws.POST("/autocomplete/service").To(p.AutocompleteService).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Autocomplete search services"))

	ws.Route(ws.GET("/batches/all").To(p.AllBatches).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all batches"))

	ws.Route(ws.POST("/batch").To(p.CreateBatch).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a batch"))

	ws.Route(ws.GET("/batch/{batch_id}").To(p.ReadBatch).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a batch"))

	ws.Route(ws.DELETE("/batch/{batch_id}").To(p.DeleteBatch).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a batch"))

	restful.Add(ws)
}

/**
* @api {get} /server/products/products/all?session={session_id}&offset={offset}&limit={limit} List all products
* @apiVersion 0.1.0
* @apiName AllProducts
* @apiGroup Product
*
* @apiDescription List all products
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/products/products/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "products": [
*       {
*         "id": "111",
*         "name": "title",
*         "orgid": "orgid",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all products successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The products were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.product.AllProducts",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.product.AllProducts",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ProductService) AllProducts(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Product.AllProducts API request")
	req_product := new(product_proto.AllProductsRequest)
	req_product.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_product.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_product.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_product.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_product.SortParameter = req.Attribute(SortParameter).(string)
	req_product.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.AllProducts(ctx, req_product)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.product.AllProducts", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all products successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/products/product?session={session_id}&offset={offset}&limit={limit} Create or update a product
* @apiVersion 0.1.0
* @apiName CreateProduct
* @apiGroup Product
*
* @apiDescription Create or update a product
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/products/product?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "product": {
*      "id": "111",
*      "name": "title",
*      "org_id": "orgid"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "product": {
*       "id": "111",
*       "name": "title",
*       "org_id": "orgid",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created product successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.product.CreateProduct",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ProductService) CreateProduct(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Product.CreateProduct API request")
	req_product := new(product_proto.CreateProductRequest)
	err := req.ReadEntity(req_product)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.product.CreateProduct", "BindError")
		return
	}
	req_product.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_product.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.CreateProduct(ctx, req_product)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.product.CreateProduct", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created product successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/products/product/{product_id}?session={session_id} View product detail
* @apiVersion 0.1.0
* @apiName ReadProduct
* @apiGroup Product
*
* @apiDescription View product detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/products/product/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "product": {
*       "id": "111",
*       "name": "title",
*       "org_id": "orgid",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read product successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The product was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.product.ReadProduct",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.product.ReadProduct",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *ProductService) ReadProduct(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Product.ReadProduct API request")
	req_product := new(product_proto.ReadProductRequest)
	req_product.Id = req.PathParameter("product_id")
	req_product.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_product.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.ReadProduct(ctx, req_product)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.product.ReadProduct", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read product successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/products/product/{product_id}?session={session_id} Delete a product
* @apiVersion 0.1.0
* @apiName DeleteProduct
* @apiGroup Product
*
* @apiDescription Delete a product
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/products/product/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted product successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The product was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.product.DeleteProduct",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ProductService) DeleteProduct(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Product.DeleteProduct API request")
	req_product := new(product_proto.DeleteProductRequest)
	req_product.Id = req.PathParameter("product_id")
	req_product.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_product.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.DeleteProduct(ctx, req_product)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.product.DeleteProduct", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted product successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/products/autocomplete/product?session={session_id}&offset={offset}&limit={limit} Autcomplete search product
* @apiVersion 0.1.0
* @apiName AutocompleteProduct
* @apiGroup Product
*
* @apiDescription Autcomplete search product
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/products/autocomplete/product?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "name": "ti"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "products": [
*       {
*         "id": "111",
*         "name": "title",
*         "orgid": "orgid",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Autocomplete searched product successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.product.AutocompleteProduct",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ProductService) AutocompleteProduct(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Product.AutocompleteProduct API request")
	req_product := new(product_proto.AutocompleteProductRequest)
	err := req.ReadEntity(req_product)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.product.AutocompleteProduct", "BindError")
		return
	}
	// req_product.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_product.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.AutocompleteProduct(ctx, req_product)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.product.AutocompleteProduct", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Autocomplete searched product successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/products/services/all?session={session_id}&offset={offset}&limit={limit} List all services
* @apiVersion 0.1.0
* @apiName AllServices
* @apiGroup Service
*
* @apiDescription List all services
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/products/services/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "services": [
*       {
*         "id": "111",
*         "name": "title",
*         "orgid": "orgid",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all services successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The services were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.service.AllServices",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.service.AllServices",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ProductService) AllServices(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Service.AllServices API request")
	req_service := new(product_proto.AllServicesRequest)
	req_service.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_service.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_service.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_service.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_service.SortParameter = req.Attribute(SortParameter).(string)
	req_service.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.AllServices(ctx, req_service)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.service.AllServices", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all services successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/services/service?session={session_id}&offset={offset}&limit={limit} Create or update a service
* @apiVersion 0.1.0
* @apiName CreateService
* @apiGroup Service
*
* @apiDescription Create or update a service
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/services/service?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "service": {
*      "id": "111",
*      "name": "title",
*      "org_id": "orgid"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "service": {
*       "id": "111",
*       "name": "title",
*       "org_id": "orgid",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created service successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.service.CreateService",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ProductService) CreateService(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Service.CreateService API request")
	req_service := new(product_proto.CreateServiceRequest)
	err := req.ReadEntity(req_service)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.service.CreateService", "BindError")
		return
	}
	req_service.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_service.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.CreateService(ctx, req_service)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.service.CreateService", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created service successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/services/service/{service_id}?session={session_id} View service detail
* @apiVersion 0.1.0
* @apiName ReadService
* @apiGroup Service
*
* @apiDescription View service detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/services/service/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "service": {
*       "id": "111",
*       "name": "title",
*       "org_id": "orgid",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read service successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The service was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.service.ReadService",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.service.ReadService",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *ProductService) ReadService(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Service.ReadService API request")
	req_service := new(product_proto.ReadServiceRequest)
	req_service.Id = req.PathParameter("service_id")
	req_service.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_service.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.ReadService(ctx, req_service)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.service.ReadService", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read service successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/services/service/{service_id}?session={session_id} Delete a service
* @apiVersion 0.1.0
* @apiName DeleteService
* @apiGroup Service
*
* @apiDescription Delete a service
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/services/service/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted service successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The service was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.service.DeleteService",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ProductService) DeleteService(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Service.DeleteService API request")
	req_service := new(product_proto.DeleteServiceRequest)
	req_service.Id = req.PathParameter("service_id")
	req_service.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_service.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.DeleteService(ctx, req_service)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.service.DeleteService", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted service successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/services/autocomplete/service?session={session_id}&offset={offset}&limit={limit} Autcomplete search service
* @apiVersion 0.1.0
* @apiName AutocompleteService
* @apiGroup Service
*
* @apiDescription Autcomplete search service
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/services/autocomplete/service?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "name": "ti"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "services": [
*       {
*         "id": "111",
*         "name": "title",
*         "orgid": "orgid",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Autocomplete searched service successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.service.AutocompleteService",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ProductService) AutocompleteService(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Service.AutocompleteService API request")
	req_service := new(product_proto.AutocompleteServiceRequest)
	err := req.ReadEntity(req_service)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.service.AutocompleteService", "BindError")
		return
	}
	// req_service.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_service.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.AutocompleteService(ctx, req_service)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.service.AutocompleteService", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Autocomplete searched service successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/products/batches/all?session={session_id}&offset={offset}&limit={limit} List all batches
* @apiVersion 0.1.0
* @apiName AllBatches
* @apiGroup Batch
*
* @apiDescription List all batches
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/products/batches/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "batches": [
*       {
*         "id": "111",
*         "name": "title",
*         "orgid": "orgid",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all batches successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The batches were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.batch.AllBatches",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.batch.AllBatches",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ProductService) AllBatches(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Batch.AllBatches API request")
	req_batch := new(product_proto.AllBatchesRequest)
	req_batch.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_batch.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_batch.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_batch.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_batch.SortParameter = req.Attribute(SortParameter).(string)
	req_batch.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.AllBatches(ctx, req_batch)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.batch.AllBatches", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all batches successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/batches/batch?session={session_id}&offset={offset}&limit={limit} Create or update a batch
* @apiVersion 0.1.0
* @apiName CreateBatch
* @apiGroup Batch
*
* @apiDescription Create or update a batch
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/batches/batch?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "batch": {
*      "id": "111",
*      "name": "title",
*      "org_id": "orgid"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "batch": {
*       "id": "111",
*       "name": "title",
*       "org_id": "orgid",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created batch successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.batch.CreateBatch",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ProductService) CreateBatch(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Batch.CreateBatch API request")
	req_batch := new(product_proto.CreateBatchRequest)
	err := req.ReadEntity(req_batch)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.batch.CreateBatch", "BindError")
		return
	}
	req_batch.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_batch.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.CreateBatch(ctx, req_batch)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.batch.CreateBatch", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created batch successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/batches/batch/{batch_id}?session={session_id} View batch detail
* @apiVersion 0.1.0
* @apiName ReadBatch
* @apiGroup Batch
*
* @apiDescription View batch detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/batches/batch/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "batch": {
*       "id": "111",
*       "name": "title",
*       "org_id": "orgid",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read batch successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The batch was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.batch.ReadBatch",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.batch.ReadBatch",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *ProductService) ReadBatch(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Batch.ReadBatch API request")
	req_batch := new(product_proto.ReadBatchRequest)
	req_batch.Id = req.PathParameter("batch_id")
	req_batch.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_batch.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.ReadBatch(ctx, req_batch)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.batch.ReadBatch", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read batch successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/batches/batch/{batch_id}?session={session_id} Delete a batch
* @apiVersion 0.1.0
* @apiName DeleteBatch
* @apiGroup Batch
*
* @apiDescription Delete a batch
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/batches/batch/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted batch successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The batch was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.batch.DeleteBatch",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ProductService) DeleteBatch(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Batch.DeleteBatch API request")
	req_batch := new(product_proto.DeleteBatchRequest)
	req_batch.Id = req.PathParameter("batch_id")
	req_batch.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_batch.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ProductClient.DeleteBatch(ctx, req_batch)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.batch.DeleteBatch", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted batch successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}
