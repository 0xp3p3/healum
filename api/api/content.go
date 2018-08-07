package api

import (
	"context"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	content_proto "server/content-srv/proto/content"
	organisation_proto "server/organisation-srv/proto/organisation"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type ContentService struct {
	ContentClient      content_proto.ContentServiceClient
	Auth               Filters
	Audit              AuditFilter
	OrganisationClient organisation_proto.OrganisationServiceClient
	ServerMetrics      metrics.Metrics
}

func (p ContentService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/content")

	audit := &audit_proto.Audit{
		ActionService:  common.ContentSrv,
		ActionResource: common.BASE + common.CONTENT_TYPE,
	}

	ws.Route(ws.GET("/sources/all").To(p.AllSources).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all sources"))

	ws.Route(ws.POST("/source").To(p.CreateSource).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a source"))

	ws.Route(ws.GET("/source/{source_id}").To(p.ReadSource).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a source"))

	ws.Route(ws.DELETE("/source/{source_id}").To(p.DeleteSource).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a source"))

	ws.Route(ws.GET("/taxonomys/all").To(p.AllTaxonomys).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all taxonomys"))

	ws.Route(ws.POST("/taxonomy").To(p.CreateTaxonomy).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a taxonomy"))

	ws.Route(ws.GET("/taxonomy/{taxonomy_id}").To(p.ReadTaxonomy).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a taxonomy"))

	ws.Route(ws.DELETE("/taxonomy/{taxonomy_id}").To(p.DeleteTaxonomy).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a taxonomy"))

	ws.Route(ws.GET("/category/items/all").To(p.AllContentCategoryItems).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all contentCategoryItems"))

	ws.Route(ws.POST("/category/item").To(p.CreateContentCategoryItem).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a contentCategoryItem"))

	ws.Route(ws.GET("/category/item/{contentCategoryItem_id}").To(p.ReadContentCategoryItem).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a contentCategoryItem"))

	ws.Route(ws.DELETE("/category/item/{contentCategoryItem_id}").To(p.DeleteContentCategoryItem).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a contentCategoryItem"))

	ws.Route(ws.GET("/contents/all").To(p.AllContents).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all contents"))

	ws.Route(ws.POST("/content").To(p.CreateContent).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Doc("Create or update a content"))

	ws.Route(ws.GET("/content/{content_id}").To(p.ReadContent).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a content"))

	ws.Route(ws.DELETE("/content/{content_id}").To(p.DeleteContent).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a content"))

	ws.Route(ws.GET("/rules/all").To(p.AllContentRules).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all contentRules"))

	ws.Route(ws.POST("/rule").To(p.CreateContentRule).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a contentRule"))

	ws.Route(ws.GET("/rule/{contentRule_id}").To(p.ReadContentRule).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a contentRule"))

	ws.Route(ws.DELETE("/rule/{contentRule_id}").To(p.DeleteContentRule).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a contentRule"))

	ws.Route(ws.POST("/filter").To(p.FilterContent).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Filter contents"))

	ws.Route(ws.POST("/search").To(p.SearchContent).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search contents"))

	ws.Route(ws.POST("/share").To(p.ShareContent).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Share contents"))

	ws.Route(ws.GET("/user/shared/{user_id}").To(p.GetAllSharedContents).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all shared content with a particular user so far"))

	ws.Route(ws.GET("/recommendations/{user_id}").To(p.GetContentRecommendations).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get content recommendations"))

	ws.Route(ws.GET("/recommendations/{user_id}/filters").To(p.GetContentFiltersByPreference).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get content filters based on user preferences"))

	ws.Route(ws.POST("/recommendations/{user_id}/filter").To(p.FilterContentRecommendations).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Filter content recommendations"))

	ws.Route(ws.GET("/tags/top/{n}").To(p.GetTopContentTags).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Return top N tags for Content"))

	ws.Route(ws.POST("/tags/autocomplete").To(p.AutocompleteContentTags).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Autocomplete for tags for Content"))
	restful.Add(ws)
}

/**
* @api {get} /server/content/sources/all?session={session_id}&offset={offset}&limit={limit} List all sources
* @apiVersion 0.1.0
* @apiName AllSources
* @apiGroup Content
*
* @apiDescription List all sources
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/sources/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "sources": [
*       {
*         "id": "111",
*         "name": "title",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "url": ""http://www.example.com",
*         "icon_url": ""http://wnww.example.com/ico",
*         "tags": ["tag1", "tag2"],
*         "orge_id":"orgid",
*         "type": { ContentSourceType },
*         "attributionRequired": true,
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all sources successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The sources were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.AllSources",
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
*           "domain": "go.micro.srv.content.AllSources",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ContentService) AllSources(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.AllSources API request")
	req_source := new(content_proto.AllSourcesRequest)
	req_source.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_source.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_source.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_source.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_source.SortParameter = req.Attribute(SortParameter).(string)
	req_source.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.AllSources(ctx, req_source)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.AllSources", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all sources successfully"
	rsp.AddHeader("Content-Type", "sourcelication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/content/source?session={session_id} Create or update a source
* @apiVersion 0.1.0
* @apiName CreateSource
* @apiGroup Content
*
* @apiDescription Create or update a source
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/source?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "source": {
*     "id": "111",
*     "name": "title",
*     "description": "description",
*     "icon_slug": "iconslug",
*     "url": ""http://www.example.com",
*     "icon_url": ""http://wnww.example.com/ico",
*     "tags": ["tag1", "tag2"],
*     "orge_id":"orgid",
*     "type": { ContentSourceType },
*     "attributionRequired": true
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "source": {
*       "id": "111",
*       "name": "title",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "url": ""http://www.example.com",
*       "icon_url": ""http://wnww.example.com/ico",
*       "tags": ["tag1", "tag2"],
*       "orge_id":"orgid",
*       "type": { ContentSourceType },
*       "attributionRequired": true,
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created source successfully"
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
*           "domain": "go.micro.srv.content.CreateSource",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) CreateSource(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.CreateSource API request")
	req_source := new(content_proto.CreateSourceRequest)
	err := req.ReadEntity(req_source)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.CreateSource", "BindError")
		return
	}
	req_source.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_source.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.CreateSource(ctx, req_source)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.CreateSource", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created source successfully"
	rsp.AddHeader("Content-Type", "sourcelication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/content/source/{source_id}?session={session_id} View source detail
* @apiVersion 0.1.0
* @apiName ReadSource
* @apiGroup Content
*
* @apiDescription View source detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/source/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "source": {
*       "id": "111",
*       "name": "title",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "url": ""http://www.example.com",
*       "icon_url": ""http://wnww.example.com/ico",
*       "tags": ["tag1", "tag2"],
*       "orge_id":"orgid",
*       "type": { ContentSourceType },
*       "attributionRequired": true,
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read source successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The source was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.ReadSource",
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
*           "domain": "go.micro.srv.content.ReadSource",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *ContentService) ReadSource(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.ReadSource API request")
	req_source := new(content_proto.ReadSourceRequest)
	req_source.Id = req.PathParameter("source_id")
	req_source.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_source.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.ReadSource(ctx, req_source)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.ReadSource", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read source successfully"
	rsp.AddHeader("Content-Type", "sourcelication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/content/source/{source_id}?session={session_id} Delete a source
* @apiVersion 0.1.0
* @apiName DeleteSource
* @apiGroup Content
*
* @apiDescription Delete a source
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/source/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted source successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The source was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.DeleteSource",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) DeleteSource(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.DeleteSource API request")
	req_source := new(content_proto.DeleteSourceRequest)
	req_source.Id = req.PathParameter("source_id")
	req_source.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_source.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.DeleteSource(ctx, req_source)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.DeleteSource", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted source successfully"
	rsp.AddHeader("Content-Type", "sourcelication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/content/taxonomys/all?session={session_id}&offset={offset}&limit={limit} List all taxonomys
* @apiVersion 0.1.0
* @apiName AllTaxonomys
* @apiGroup Content
*
* @apiDescription List all taxonomys
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/taxonomys/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "taxonomys": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "org_id": "orgid",
*         "tags": ["tag1", "tag2"],
*         "weigth": 100,
*         "priority": 1,
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all taxonomys successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The taxonomys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.AllTaxonomys",
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
*           "domain": "go.micro.srv.content.AllTaxonomys",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ContentService) AllTaxonomys(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.AllTaxonomys API request")
	req_taxonomy := new(content_proto.AllTaxonomysRequest)
	req_taxonomy.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_taxonomy.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_taxonomy.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_taxonomy.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_taxonomy.SortParameter = req.Attribute(SortParameter).(string)
	req_taxonomy.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.AllTaxonomys(ctx, req_taxonomy)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.AllTaxonomys", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all taxonomys successfully"
	rsp.AddHeader("Content-Type", "taxonomylication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/content/taxonomy?session={session_id} Create or update a taxonomy
* @apiVersion 0.1.0
* @apiName CreateTaxonomy
* @apiGroup Content
*
* @apiDescription Create or update a taxonomy
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/taxonomy?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "taxonomy": {
*     "id": "111",
*     "name": "title",
*     "summary": "summary",
*     "description": "description",
*     "org_id": "orgid",
*     "tags": ["tag1", "tag2"],
*     "weigth": 100,
*     "priority": 1
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "taxonomy": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "org_id": "orgid",
*       "tags": ["tag1", "tag2"],
*       "weigth": 100,
*       "priority": 1,
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created taxonomy successfully"
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
*           "domain": "go.micro.srv.content.CreateTaxonomy",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) CreateTaxonomy(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.CreateTaxonomy API request")
	req_taxonomy := new(content_proto.CreateTaxonomyRequest)
	err := req.ReadEntity(req_taxonomy)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.CreateTaxonomy", "BindError")
		return
	}
	req_taxonomy.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_taxonomy.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.CreateTaxonomy(ctx, req_taxonomy)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.CreateTaxonomy", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created taxonomy successfully"
	rsp.AddHeader("Content-Type", "taxonomylication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/content/taxonomy/{taxonomy_id}?session={session_id} View taxonomy detail
* @apiVersion 0.1.0
* @apiName ReadTaxonomy
* @apiGroup Content
*
* @apiDescription View taxonomy detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/taxonomy/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "taxonomy": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "org_id": "orgid",
*       "tags": ["tag1", "tag2"],
*       "weigth": 100,
*       "priority": 1,
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read taxonomy successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The taxonomy was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.ReadTaxonomy",
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
*           "domain": "go.micro.srv.content.ReadTaxonomy",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *ContentService) ReadTaxonomy(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.ReadTaxonomy API request")
	req_taxonomy := new(content_proto.ReadTaxonomyRequest)
	req_taxonomy.Id = req.PathParameter("taxonomy_id")
	req_taxonomy.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_taxonomy.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.ReadTaxonomy(ctx, req_taxonomy)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.ReadTaxonomy", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read taxonomy successfully"
	rsp.AddHeader("Content-Type", "taxonomylication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/content/taxonomy/{taxonomy_id}?session={session_id} Delete a taxonomy
* @apiVersion 0.1.0
* @apiName DeleteTaxonomy
* @apiGroup Content
*
* @apiDescription Delete a taxonomy
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/taxonomy/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted taxonomy successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The taxonomy was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.DeleteTaxonomy",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) DeleteTaxonomy(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.DeleteTaxonomy API request")
	req_taxonomy := new(content_proto.DeleteTaxonomyRequest)
	req_taxonomy.Id = req.PathParameter("taxonomy_id")
	req_taxonomy.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_taxonomy.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.DeleteTaxonomy(ctx, req_taxonomy)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.DeleteTaxonomy", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted taxonomy successfully"
	rsp.AddHeader("Content-Type", "taxonomylication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/content/category/items/all?session={session_id}&offset={offset}&limit={limit} List all contentCategoryItems
* @apiVersion 0.1.0
* @apiName AllContentCategoryItems
* @apiGroup Content
*
* @apiDescription List all contentCategoryItems
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/category/items/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentCategoryItems": [
*       {
*         "id": "111",
*         "name": "title",
*         "name_slug": "nameSlug",
*         "icon_slug": "iconSlug",
*         "summary": "summary",
*         "description": "description",
*         "org_id": "orgid",
*         "tags": ["tag1", "tag2"],
*         "taxonomy": { Taxonomy },
*         "weight": 100,
*         "priority": 1,
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all contentCategoryItems successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentCategoryItems were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.AllContentCategoryItems",
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
*           "domain": "go.micro.srv.content.AllContentCategoryItems",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ContentService) AllContentCategoryItems(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.AllContentCategoryItems API request")
	req_contentCategoryItem := new(content_proto.AllContentCategoryItemsRequest)
	req_contentCategoryItem.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentCategoryItem.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_contentCategoryItem.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_contentCategoryItem.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_contentCategoryItem.SortParameter = req.Attribute(SortParameter).(string)
	req_contentCategoryItem.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.AllContentCategoryItems(ctx, req_contentCategoryItem)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.AllContentCategoryItems", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all contentCategoryItems successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "contentCategoryItemlication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/content/category/item?session={session_id} Create or update a contentCategoryItem
* @apiVersion 0.1.0
* @apiName CreateContentCategoryItem
* @apiGroup Content
*
* @apiDescription Create or update a contentCategoryItem
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/category/item?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "contentCategoryItem": {
*     "id": "111",
*     "name": "title",
*     "name_slug": "nameSlug",
*     "icon_slug": "iconSlug",
*     "summary": "summary",
*     "description": "description",
*     "org_id": "orgid",
*     "tags": ["tag1", "tag2"],
*     "taxonomy": { Taxonomy },
*     "weight": 100,
*     "priority": 1
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentCategoryItem": {
*       "id": "111",
*       "name": "title",
*       "name_slug": "nameSlug",
*       "icon_slug": "iconSlug",
*       "summary": "summary",
*       "description": "description",
*       "org_id": "orgid",
*       "tags": ["tag1", "tag2"],
*       "taxonomy": { Taxonomy },
*       "weight": 100,
*       "priority": 1,
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created contentCategoryItem successfully"
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
*           "domain": "go.micro.srv.content.CreateContentCategoryItem",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) CreateContentCategoryItem(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.CreateContentCategoryItem API request")
	req_contentCategoryItem := new(content_proto.CreateContentCategoryItemRequest)
	// err := req.ReadEntity(req_contentCategoryItem)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.CreateContentCategoryItem", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_contentCategoryItem); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.CreateContentCategoryItem", "BindError")
		return
	}
	req_contentCategoryItem.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentCategoryItem.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.CreateContentCategoryItem(ctx, req_contentCategoryItem)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.CreateContentCategoryItem", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created contentCategoryItem successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "contentCategoryItemlication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/content/category/item/{contentCategoryItem_id}?session={session_id} View contentCategoryItem detail
* @apiVersion 0.1.0
* @apiName ReadContentCategoryItem
* @apiGroup Content
*
* @apiDescription View contentCategoryItem detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/category/item/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentCategoryItem": {
*       "id": "111",
*       "name": "title",
*       "name_slug": "nameSlug",
*       "icon_slug": "iconSlug",
*       "summary": "summary",
*       "description": "description",
*       "org_id": "orgid",
*       "tags": ["tag1", "tag2"],
*       "taxonomy": { Taxonomy },
*       "weight": 100,
*       "priority": 1,
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read contentCategoryItem successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentCategoryItem was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.ReadContentCategoryItem",
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
*           "domain": "go.micro.srv.content.ReadContentCategoryItem",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *ContentService) ReadContentCategoryItem(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.ReadContentCategoryItem API request")
	req_contentCategoryItem := new(content_proto.ReadContentCategoryItemRequest)
	req_contentCategoryItem.Id = req.PathParameter("contentCategoryItem_id")
	req_contentCategoryItem.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentCategoryItem.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.ReadContentCategoryItem(ctx, req_contentCategoryItem)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.ReadContentCategoryItem", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read contentCategoryItem successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "contentCategoryItemlication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {delete} /server/content/category/item/{contentCategoryItem_id}?session={session_id} Delete a contentCategoryItem
* @apiVersion 0.1.0
* @apiName DeleteContentCategoryItem
* @apiGroup Content
*
* @apiDescription Delete a contentCategoryItem
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/category/item/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted contentCategoryItem successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentCategoryItem was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.DeleteContentCategoryItem",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) DeleteContentCategoryItem(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.DeleteContentCategoryItem API request")
	req_contentCategoryItem := new(content_proto.DeleteContentCategoryItemRequest)
	req_contentCategoryItem.Id = req.PathParameter("contentCategoryItem_id")
	req_contentCategoryItem.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentCategoryItem.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.DeleteContentCategoryItem(ctx, req_contentCategoryItem)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.DeleteContentCategoryItem", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted contentCategoryItem successfully"
	rsp.AddHeader("Content-Type", "contentCategoryItemlication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/content/contents/all?session={session_id}&offset={offset}&limit={limit} List all contents
* @apiVersion 0.1.0
* @apiName AllContents
* @apiGroup Content
*
* @apiDescription List all contents
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/contents/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contents": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": ["summary0", "summary1"],
*         "description": "description",
*         "org_id": "orgid",
*         "createdBy": { User },
*         "url": "http://www.example.com",
*         "author": "author",
*         "timestamp": 1517891917,
*         "text": "text",
*         "tags": [ ContentCategoryItem, ... ],
*         "type": { ContentType },
*         "source": { Source },
*         "category": { ContentCategory },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all contents successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contents were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.AllContents",
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
*           "domain": "go.micro.srv.content.AllContents",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ContentService) AllContents(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.AllContents API request")
	req_content := new(content_proto.AllContentsRequest)
	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_content.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_content.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_content.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_content.SortParameter = req.Attribute(SortParameter).(string)
	req_content.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.AllContents(ctx, req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.AllContents", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all contents successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/content/content?session={session_id} Create or update a content
* @apiVersion 0.1.0
* @apiName CreateContent
* @apiGroup Content
*
* @apiDescription Create or update a content
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/content?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "content": {
*     "id": "111",
*     "name": "title",
*     "summary": ["summary0", "summary1"],
*     "description": "description",
*     "org_id": "orgid",
*     "createdBy": { User },
*     "url": "http://www.example.com",
*     "author": "author",
*     "timestamp": 1517891917,
*     "text": "text",
*     "tags": [ {"id":"111"}, {"id":"222"}... ],
*     "type": { ContentType },
*     "source": { Source },
*     "category": { ContentCategory },
*     "recipe": { Recipe },
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "content": {
*       "id": "111",
*       "name": "title",
*       "summary": ["summary0", "summary1"],
*       "description": "description",
*       "org_id": "orgid",
*       "createdBy": { User },
*       "url": "http://www.example.com",
*       "author": "author",
*       "timestamp": 1517891917,
*       "text": "text",
*       "tags": [ ContentCategoryItem, ... ],
*       "type": { ContentType },
*       "source": { Source },
*       "category": { ContentCategory },
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created content successfully"
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
*           "domain": "go.micro.srv.content.CreateContent",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) CreateContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.CreateContent API request")
	req_content := new(content_proto.CreateContentRequest)
	if err := utils.UnmarshalAny(req, rsp, req_content); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.CreateContent", "BindError")
		return
	}
	req_content.UserId = req.Attribute(UserIdAttrName).(string)
	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_content.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.CreateContent(ctx, req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.CreateContent", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created content successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/content/content/{content_id}?session={session_id} View content detail
* @apiVersion 0.1.0
* @apiName ReadContent
* @apiGroup Content
*
* @apiDescription View content detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/content/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "content": {
        "id": "111",
*       "name": "title",
*       "summary": ["summary0", "summary1"],
*       "description": "description",
*       "org_id": "orgid",
*       "createdBy": { User },
*       "url": "http://www.example.com",
*       "author": "author",
*       "timestamp": 1517891917,
*       "text": "text",
*       "tags": [ ContentCategoryItem, ... ],
*       "type": { ContentType },
*       "source": { Source },
*       "category": { ContentCategory },
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read content successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The content was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.ReadContent",
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
*           "domain": "go.micro.srv.content.ReadContent",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *ContentService) ReadContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.ReadContent API request")
	req_content := new(content_proto.ReadContentRequest)
	req_content.Id = req.PathParameter("content_id")
	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_content.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.ReadContent(ctx, req_content)
	if err != nil || resp == nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.ReadContent", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read content successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {delete} /server/content/content/{content_id}?session={session_id} Delete a content
* @apiVersion 0.1.0
* @apiName DeleteContent
* @apiGroup Content
*
* @apiDescription Delete a content
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/content/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted content successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The content was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.DeleteContent",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) DeleteContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.DeleteContent API request")
	req_content := new(content_proto.DeleteContentRequest)

	req_content.Id = req.PathParameter("content_id")
	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_content.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.DeleteContent(ctx, req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.DeleteContent", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted content successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/content/rules/all?session={session_id}&offset={offset}&limit={limit} List all contentRules
* @apiVersion 0.1.0
* @apiName AllContentRules
* @apiGroup Content
*
* @apiDescription List all contentRules
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/rules/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentRules": [
*       {
*         "id": "111",
*         "org_id": "orgid",
*         "type": 1,
*         "source": { Source },
*         "sourceType": { SourceType },
*         "contentType": { ContentType },
*         "parentCategory": { ContentParentCategory },
*         "category": { ContentCategory },
*         "categoryItems": [ ContentCategoryItem, ContentCategoryItem, ... ],
*         "expression": { Expression },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all contentRules successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentRules were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.AllContentRules",
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
*           "domain": "go.micro.srv.content.AllContentRules",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ContentService) AllContentRules(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.AllContentRules API request")
	req_contentRule := new(content_proto.AllContentRulesRequest)
	req_contentRule.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentRule.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_contentRule.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_contentRule.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_contentRule.SortParameter = req.Attribute(SortParameter).(string)
	req_contentRule.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.AllContentRules(ctx, req_contentRule)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.AllContentRules", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all contentRules successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "contentRulelication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/content/rule?session={session_id} Create or update a contentRule
* @apiVersion 0.1.0
* @apiName CreateContentRule
* @apiGroup Content
*
* @apiDescription Create or update a contentRule
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/rule?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "contentRule": {
*     "id": "111",
*     "org_id": "orgid",
*     "type": 1,
*     "source": { Source },
*     "sourceType": { SourceType },
*     "contentType": { ContentType },
*     "parentCategory": { ContentParentCategory },
*     "category": { ContentCategory },
*     "categoryItems": [ ContentCategoryItem, ContentCategoryItem, ... ],
*     "expression": { Expression },
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentRule": {
*       "id": "111",
*       "org_id": "orgid",
*       "type": 1,
*       "source": { Source },
*       "sourceType": { SourceType },
*       "contentType": { ContentType },
*       "parentCategory": { ContentParentCategory },
*       "category": { ContentCategory },
*       "categoryItems": [ ContentCategoryItem, ContentCategoryItem, ... ],
*       "expression": { Expression },
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created contentRule successfully"
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
*           "domain": "go.micro.srv.content.CreateContentRule",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) CreateContentRule(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.CreateContentRule API request")
	req_contentRule := new(content_proto.CreateContentRuleRequest)
	// err := req.ReadEntity(req_contentRule)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.CreateContentRule", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_contentRule); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.CreateContentRule", "BindError")
		return
	}
	req_contentRule.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentRule.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.CreateContentRule(ctx, req_contentRule)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.CreateContentRule", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created contentRule successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "contentRulelication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/content/rule/{contentRule_id}?session={session_id} View contentRule detail
* @apiVersion 0.1.0
* @apiName ReadContentRule
* @apiGroup Content
*
* @apiDescription View contentRule detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/rule/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentRule": {
*       "id": "111",
*       "org_id": "orgid",
*       "type": 1,
*       "source": { Source },
*       "sourceType": { SourceType },
*       "contentType": { ContentType },
*       "parentCategory": { ContentParentCategory },
*       "category": { ContentCategory },
*       "categoryItems": [ ContentCategoryItem, ContentCategoryItem, ... ],
*       "expression": { Expression },
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read contentRule successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentRule was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.ReadContentRule",
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
*           "domain": "go.micro.srv.content.ReadContentRule",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *ContentService) ReadContentRule(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.ReadContentRule API request")
	req_contentRule := new(content_proto.ReadContentRuleRequest)
	req_contentRule.Id = req.PathParameter("contentRule_id")
	req_contentRule.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentRule.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.ReadContentRule(ctx, req_contentRule)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.ReadContentRule", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read contentRule successfully"
	rsp.AddHeader("Content-Type", "contentRulelication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/content/rule/{contentRule_id}?session={session_id} Delete a contentRule
* @apiVersion 0.1.0
* @apiName DeleteContentRule
* @apiGroup Content
*
* @apiDescription Delete a contentRule
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/rule/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted contentRule successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentRule was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.DeleteContentRule",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) DeleteContentRule(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.DeleteContentRule API request")
	req_contentRule := new(content_proto.DeleteContentRuleRequest)
	req_contentRule.Id = req.PathParameter("contentRule_id")
	req_contentRule.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentRule.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.DeleteContentRule(ctx, req_contentRule)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.DeleteContentRule", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted contentRule successfully"
	rsp.AddHeader("Content-Type", "contentRulelication/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/content/filter?session={session_id}&offset={offset}&limit={limit} Filter contents
* @apiVersion 0.1.0
* @apiName FilterContent
* @apiGroup Content
*
* @apiDescription Filter contents
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "sources": ["source1", ...],
*   "sourceTypes": ["sourceType1"],
*   "contentTypes": ["contentType1", ...],
*   "created_by": ["111", ...],
*   "contentParentCategories": ["contentParentCategory"],
*   "contentCategories": ["activity", ...],
*   "contentCategoryItems": ["contentCategoryItem"],
"	"type":["healum.com/proto/go.micro.srv.content.Recipe","healum.com/proto/go.micro.srv.content.Execise","healum.com/proto/go.micro.srv.content.Article","healum.com/proto/go.micro.srv.content.Video",],
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contents": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": ["summary0", "summary1"],
*         "description": "description",
*         "org_id": "orgid",
*         "createdBy": { User },
*         "url": "http://www.example.com",
*         "author": "author",
*         "timestamp": 1517891917,
*         "text": "text",
*         "tags": [ ContentCategoryItem, ... ],
*         "type": { ContentType },
*         "source": { Source },
*         "category": { ContentCategory },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all contents successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contents were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.AllContents",
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
*           "domain": "go.micro.srv.content.AllContents",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *ContentService) FilterContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.FilterContent API request")
	req_content := new(content_proto.FilterContentRequest)
	if err := utils.UnmarshalAny(req, rsp, req_content); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.FilterContent", "BindError")
		return
	}

	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_content.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_content.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_content.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_content.SortParameter = req.Attribute(SortParameter).(string)
	req_content.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.FilterContent(ctx, req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.FilterContent", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Filter contents successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/content/search?session={session_id}&offset={offset}&limit={limit} Search contents
* @apiVersion 0.1.0
* @apiName SearchContent
* @apiGroup Content
*
* @apiDescription Search contents
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/search?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "title": "title",
*   "description": "descript",
*   "summary": "summary"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contents": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": ["summary0", "summary1"],
*         "description": "description",
*         "org_id": "orgid",
*         "createdBy": { User },
*         "url": "http://www.example.com",
*         "author": "author",
*         "timestamp": 1517891917,
*         "text": "text",
*         "tags": [ ContentCategoryItem, ... ],
*         "type": { ContentType },
*         "source": { Source },
*         "category": { ContentCategory },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Search contents successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contents were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.SearchContent",
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
*           "domain": "go.micro.srv.content.SearchContent",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ContentService) SearchContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.SearchContent API request")
	req_search := new(content_proto.SearchContentRequest)
	if err := utils.UnmarshalAny(req, rsp, req_search); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.SearchContent", "BindError")
		return
	}
	req_search.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_search.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_search.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_search.Offset = req.Attribute(PaginateOffsetParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.SearchContent(ctx, req_search)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.SearchContent", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Search content successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/content/share?session={session_id} Share content with user
* @apiVersion 0.1.0
* @apiName ShareContent
* @apiGroup Content
*
* @apiDescription Share a different types of content with user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/share?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "contents": [ {Content}, ...],
*   "users": [ {User}, ...],
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "code": 200,
*     "message": "Shared all contents successfully"
*   }
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contents were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.AllContents",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) ShareContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.ShareContent API request")
	req_content := new(content_proto.ShareContentRequest)
	err := req.ReadEntity(req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.FilterContent", "BindError")
		return
	}
	req_content.UserId = req.Attribute(UserIdAttrName).(string)
	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.ShareContent(ctx, req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.ShareContent", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Filter contents successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/content/user/shared/{user_id}?session={session_id}&offset={offset}&limit={limit} Get all shared content with a particular user so far
* @apiVersion 0.1.0
* @apiName GetAllSharedContents
* @apiGroup Content
*
* @apiDescription API should return all content that has been shared with the normal user by any employee of the org
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/user/shared/{user_id}?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "shareContentUsers": [
*       {
*         "id": "111",
*         "content": { Content },
*         "user": { User },
*         "status": "SHARED",
*         "shared_by": { User },
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get all shared contents successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contents were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.AllContents",
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
*           "domain": "go.micro.srv.content.AllContents",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ContentService) GetAllSharedContents(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.GetAllSharedContents API request")
	req_content := new(content_proto.GetAllSharedContentsRequest)

	req_content.UserId = req.Attribute(UserIdAttrName).(string)
	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_content.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_content.Offset = req.Attribute(PaginateOffsetParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.GetAllSharedContents(ctx, req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.GetAllSharedContents", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get all shared contents successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/content/recommendations/{user_id}?session={session_id}&offset={offset}&limit={limit} Get content recommendations
* @apiVersion 0.1.0
* @apiName GetContentRecommendations
* @apiGroup Content
*
* @apiDescription API should return all content recommendations that are available for this normal user (not the employee) of the org where org_id = user.org_id
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/recommendations/{user_id}?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "recommendations": [
*       {
*         "id": "111",
*         "content": { Content },
*         "org_id": "orgid",
*         "user_id": "userid",
*         "tags": [{ContentCateogryItem},...],
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get content recommendations successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contents were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.GetContentRecommendations",
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
*           "domain": "go.micro.srv.content.GetContentRecommendations",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ContentService) GetContentRecommendations(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.GetContentRecommendations API request")
	req_content := new(content_proto.GetContentRecommendationsRequest)

	req_content.UserId = req.Attribute(UserIdAttrName).(string)
	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.GetContentRecommendations(ctx, req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.GetContentRecommendations", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get content recommendations successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/content/recommendations/{user_id}/filters?session={session_id} Get content filters based on user preferences
* @apiVersion 0.1.0
* @apiName GetContentFiltersByPreference
* @apiGroup Content
*
* @apiDescription This should return user.preferences.contentcategoryitems with their respective category
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/user/shared/{user_id}?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentCategoryItems": [
*       {
*         "id": "111",
*         "name": "title",
*         "name_slug": "nameSlug",
*         "icon_slug": "iconSlug",
*         "summary": "summary",
*         "description": "description",
*         "org_id": "orgid",
*         "tags": ["tag1", "tag2"],
*         "taxonomy": { Taxonomy },
*         "weight": 100,
*         "priority": 1,
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Filtered content category items successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contents were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.AllContents",
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
*           "domain": "go.micro.srv.content.AllContents",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ContentService) GetContentFiltersByPreference(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.GetContentFiltersByPreference API request")
	req_content := new(content_proto.GetContentFiltersByPreferenceRequest)

	req_content.UserId = req.Attribute(UserIdAttrName).(string)
	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.GetContentFiltersByPreference(ctx, req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.GetContentFiltersByPreference", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Filtered content category items successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/content/recommendations/{user_id}/filter?session={session_id}&offset={offset}&limit={limit} Get all shared content with a particular user so far
* @apiVersion 0.1.0
* @apiName FilterContentRecommendations
* @apiGroup Content
*
* @apiDescription API should return all content that has been shared with the normal user by any employee of the org
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/content/recommendations/{user_id}/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "contentCategoryItems": [
*     {
*       "id": "111",
*       "name": "title",
*       "name_slug": "nameSlug",
*       "icon_slug": "iconSlug",
*       "summary": "summary",
*       "description": "description",
*       "org_id": "orgid",
*       "tags": ["tag1", "tag2"],
*       "taxonomy": { Taxonomy },
*       "weight": 100,
*       "priority": 1
*     },
*     ... ...
*   ]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "image": "Content.Image",
*         "title": "Content.Title",
*         "author": "Content.Author",
*         "source": { Content.Source },
*         "content_id": "Content.Id",
*         "category_id": "Content.ContentCategory.Id",
*         "icon_lsug": "Content.ContentCategory.IconSlug",
*         "category_name": "Content.ContentCategory.Name"
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Filter content recommendations successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contents were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.AllContents",
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
*           "domain": "go.micro.srv.content.AllContents",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *ContentService) FilterContentRecommendations(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.FilterContentRecommendations API request")
	req_content := new(content_proto.FilterContentRecommendationsRequest)
	err := req.ReadEntity(req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.FilterContentRecommendations", "BindError")
		return
	}
	req_content.UserId = req.Attribute(UserIdAttrName).(string)
	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.FilterContentRecommendations(ctx, req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.FilterContentRecommendations", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Filter content recommendations successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/contents/tags/top/{n}?session={session_id} Return top N tags for Content
* @apiVersion 0.1.0
* @apiName GetTopContentTags
* @apiGroup Content
*
* @apiDescription For each of the following service we have return top N tags for content
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/contents/tags/top/5?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "tags": ["tag1","tag2","tag3",...]
*   },
*   "code": 200,
*   "message": "Get top content tags successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, SearchError.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.GetTopContentTags",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) GetTopContentTags(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.GetTopContentTags API request")
	req_content := new(content_proto.GetTopTagsRequest)
	n, _ := strconv.Atoi(req.PathParameter("n"))
	req_content.N = int64(n)
	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_content.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.GetTopTags(ctx, req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.GetTopContentTags", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Get top content tags successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/contents/tags/autocomplete?session={session_id} Autocomplete for tags for Content
* @apiVersion 0.1.0
* @apiName AutocompleteContentTags
* @apiGroup Content
*
* @apiDescription Autocomplete for tags for Content
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/contents/tags/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "name": "tag"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "tags": ["tag1","tag2","tag3",...]
*   },
*   "code": 200,
*   "message": "Autocomplete content tags successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, SearchError.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.content.AutocompleteContentTags",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ContentService) AutocompleteContentTags(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Content.AutocompleteContentTags API request")
	req_content := new(content_proto.AutocompleteTagsRequest)
	err := req.ReadEntity(req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.AutocompleteContentTags", "BindError")
		return
	}
	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)

	// req_content.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.AutocompleteTags(ctx, req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.AutocompleteContentTags", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Autocomplete content tags successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}
