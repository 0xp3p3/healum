package api

import (
	"context"
	"errors"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	organisation_proto "server/organisation-srv/proto/organisation"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type OrganisationService struct {
	OrganisationClient organisation_proto.OrganisationServiceClient
	Auth               Filters
	Audit              AuditFilter
	ServerMetrics      metrics.Metrics
}

func (p OrganisationService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/organisations")

	audit := &audit_proto.Audit{
		ActionService:  common.OrganisationSrv,
		ActionResource: common.BASE + common.ORGANISATION_TYPE,
	}

	ws.Route(ws.GET("/all").To(p.AllOrganisations).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all orgs"))

	ws.Route(ws.POST("/organisation").To(p.CreateOrganisation).
		// Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create one org"))

	ws.Route(ws.PUT("/organisation").To(p.UpdateOrganisation).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Update one org"))

	ws.Route(ws.GET("/organisation/{org_id}").To(p.ReadOrganisation).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read one org"))

	ws.Route(ws.POST("/profile").To(p.CreateOrganisationProfile).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Created org profile"))

	ws.Route(ws.POST("/setting").To(p.CreateOrganisationSetting).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Created org setting"))

	ws.Route(ws.POST("/modules").To(p.UpdateModules).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("update modules"))

	ws.Route(ws.GET("/modules").To(p.GetModulesByOrg).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("get organisation modules"))

	restful.Add(ws)
}

/**
* @api {get} /server/organisations/all?session={session_id}&offset={offset}&limit={limit} List all organisations
* @apiVersion 0.1.0
* @apiName AllOrganisations
* @apiGroup Organisation
*
* @apiDescription List all organisations
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/organisations/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "organisations": [
*       {
*         "id": "111",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all notes successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The notes were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.note.AllNotes",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *OrganisationService) AllOrganisations(req *restful.Request, rsp *restful.Response) {
	log.Info("Received OrganisationService.All API request")

	req_org := new(organisation_proto.AllRequest)
	req_org.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_org.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_org.SortParameter = req.Attribute(SortParameter).(string)
	req_org.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	all_resp, err := p.OrganisationClient.All(ctx, req_org)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.All", "QueryError")
		return
	}

	all_resp.Code = http.StatusOK
	all_resp.Message = "Read all organisations succesfully"
	data := utils.MarshalAny(rsp, all_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/organisations/organisation?session={session_id} Create or update an organisation
* @apiVersion 0.1.0
* @apiName CreateOrganisation
* @apiGroup Organisation
*
* @apiDescription Create or update an organisation
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/organisations/organisation?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "organisation": {
*     "type": "ROOT",
*     "owner":  { User },
*   },
*   "account": { Account },
*   "user": { User },
*   "modules": [Module]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "organisation": {
*       "id": "111",
*       "type": "ROOT",
*       "owner":  { User },
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created organisation successfully"
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
*           "domain": "go.micro.srv.organisation.Create",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *OrganisationService) CreateOrganisation(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Organisation.Create API request")
	req_org := new(organisation_proto.CreateRequest)
	if err := utils.UnmarshalAny(req, rsp, req_org); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.Create", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.OrganisationClient.Create(ctx, req_org)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.Create", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created organisation successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {put} /server/organisations/organisation?session={session_id} Update an organisation
* @apiVersion 0.1.0
* @apiName CreateOrganisation
* @apiGroup UpdateOrganisation
*
* @apiDescription Update an organisation
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/organisations/organisation?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "organisation": {
*     "id": "111",
*     "name": "updated name",
*     "type": "ROOT"
*   },
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "organisation": {
*       "id": "111",
*       "type": "updated name",
*       "type": "ROOT",
*       "owner":  { User },
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Updated organisation successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, UpdateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "UpdateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.organisation.Update",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *OrganisationService) UpdateOrganisation(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Organisation.Update API request")
	req_org := new(organisation_proto.UpdateRequest)
	if err := utils.UnmarshalAny(req, rsp, req_org); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.Update", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.OrganisationClient.Update(ctx, req_org)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.Update", "UpdateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Updated organisation successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/organisations/organisation/{org_id}?session={session_id} View organisation detail
* @apiVersion 0.1.0
* @apiName ReadOrganisation
* @apiGroup Organisation
*
* @apiDescription View organisation detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/organisations/organisation/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "organisation": {
*       "id": "111",
*       "type": 1,
*       "owner":  { User },
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read organisation successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The organisation was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.organisation.Read",
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
*           "domain": "go.micro.srv.organisation.Read",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *OrganisationService) ReadOrganisation(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Organisation.Read API request")
	req_org := new(organisation_proto.ReadRequest)
	req_org.OrgId = req.PathParameter("org_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.OrganisationClient.Read(ctx, req_org)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.Read", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read organisation successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/organisations/profile?session={session_id} Update an organisation profile
* @apiVersion 0.1.0
* @apiName CreateOrganisationProfile
* @apiGroup Organisation
*
* @apiDescription Update an organisation profile
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/organisations/profile?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "profile": {
*     "id": "111",
*     "org_id": "orgid"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "profile": {
*       "id": "111",
*       "org_id": "orgid",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created profile successfully"
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
*           "domain": "go.micro.srv.organisation.Create",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *OrganisationService) CreateOrganisationProfile(req *restful.Request, rsp *restful.Response) {
	log.Info("Received OrganisationCreate API request")
	req_org := new(organisation_proto.CreateOrganisationProfileRequest)
	// err := req.ReadEntity(req_org)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.CreateOrganisationProfile", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_org); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.CreateOrganisationProfile", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.OrganisationClient.CreateOrganisationProfile(ctx, req_org)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.CreateOrganisationProfile", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created team successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/organisations/setting?session={session_id} Update an organisation setting
* @apiVersion 0.1.0
* @apiName CreateOrganisationSetting
* @apiGroup Organisation
*
* @apiDescription Update an organisation setting
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/organisations/setting?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "team": {
*     "id": "111",
*     "org_id": "orgid"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "team": {
*       "id": "111",
*       "org_id": "orgid",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created or updated setting successfully"
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
*           "domain": "go.micro.srv.organisation.CreateOrganisationSetting",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *OrganisationService) CreateOrganisationSetting(req *restful.Request, rsp *restful.Response) {
	log.Info("Received OrganisationCreate API request")
	req_org := new(organisation_proto.CreateOrganisationSettingRequest)
	err := req.ReadEntity(req_org)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.CreateOrganisationSetting", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.OrganisationClient.CreateOrganisationSetting(ctx, req_org)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.CreateOrganisationSetting", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created or updated setting successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/organisations/modules?session={session_id} Update an organisation modules
* @apiVersion 0.1.0
* @apiName UpdateModules
* @apiGroup Organisation
*
* @apiDescription Update an organisation modules
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/organisations/modules?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "modules": [
*     {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "org_id": "orgid",
*       "settings": "settings"
*     },
*     ... ...
*   ]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "modules": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "org_id": "orgid",
*         "settings": "settings",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Created or updated modules successfully"
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
*           "domain": "go.micro.srv.organisation.UpdateModules",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *OrganisationService) UpdateModules(req *restful.Request, rsp *restful.Response) {
	log.Info("Received UpdateModules API request")
	req_org := new(organisation_proto.UpdateModulesRequest)
	err := req.ReadEntity(req_org)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.UpdateModules", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.OrganisationClient.UpdateModules(ctx, req_org)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.UpdateModules", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created or updated modules successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/organisations/modules?session={session_id} Read an organisation modules
* @apiVersion 0.1.0
* @apiName GetModulesByOrg
* @apiGroup Organisation
*
* @apiDescription Read an organisation modules
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/organisations/modules?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "modules": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "org_id": "orgid",
*         "settings": "settings",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get modules by organisation successfully"
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
*           "domain": "go.micro.srv.organisation.GetModulesByOrg",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *OrganisationService) GetModulesByOrg(req *restful.Request, rsp *restful.Response) {
	log.Info("Received GetModulesByOrg API request")
	req_org := new(organisation_proto.GetModulesByOrgRequest)
	req_org.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.OrganisationClient.GetModulesByOrg(ctx, req_org)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.organisation.GetModulesByOrg", "CreateError")
		return
	} else if len(resp.Data.Modules) == 0 {
		utils.WriteErrorResponse(rsp, errors.New("NotFound"), "go.micro.srv.organisation.GetModulesByOrg", "NotFound")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get modules by organisation successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}
