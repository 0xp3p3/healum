package api

import (
	"context"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	organisation_proto "server/organisation-srv/proto/organisation"
	plan_proto "server/plan-srv/proto/plan"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type PlanService struct {
	PlanClient         plan_proto.PlanServiceClient
	Auth               Filters
	Audit              AuditFilter
	OrganisationClient organisation_proto.OrganisationServiceClient
	FilterMiddle       Filters
	ServerMetrics      metrics.Metrics
}

func (p PlanService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/plans")

	audit := &audit_proto.Audit{
		ActionService:  common.PlanSrv,
		ActionResource: common.BASE + common.PLAN_TYPE,
	}

	ws.Route(ws.GET("/all").To(p.AllPlans).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all plans"))

	ws.Route(ws.POST("/plan").To(p.CreatePlan).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("create data"))

	ws.Route(ws.GET("/plan/{plan_id}").To(p.ReadPlan).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("View Plan detail"))

	ws.Route(ws.DELETE("/plan/{plan_id}").To(p.DeletePlan).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("View Plan detail"))

	ws.Route(ws.POST("/search").To(p.SearchPlans).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Return searched plans"))

	ws.Route(ws.GET("/templates").To(p.Templates).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all templates"))

	ws.Route(ws.GET("/drafts").To(p.Drafts).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all plans where the status is draft"))

	ws.Route(ws.GET("/creator/{user_id}").To(p.ByCreator).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all plans created by a particular team member"))

	ws.Route(ws.GET("/filters/all").To(p.Filters).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all plan filters"))

	ws.Route(ws.GET("/filter/time").To(p.TimeFilters).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all plans by time period"))

	ws.Route(ws.GET("/filter/goal").To(p.GoalFilters).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all plans by goal category"))

	ws.Route(ws.POST("/plan/search/autocomplete").To(p.AutocompletePlanSearch).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete plan text"))

	ws.Route(ws.GET("/tags/top/{n}").To(p.GetTopPlanTags).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Return top N tags for Plan"))

	ws.Route(ws.POST("/tags/autocomplete").To(p.AutocompletePlanTags).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Autocomplete for tags for Plan"))

	restful.Add(ws)
}

/**
* @api {get} /server/plans/all?session={session_id}&offset={offset}&limit={limit} List all plans
* @apiVersion 0.1.0
* @apiName AllPlans
* @apiGroup Plan
*
* @apiDescription AllPlans
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "plans": [
*       {
*         "id": "111",
*         "title": "plan1",
*         "org_id": "orgid",
*         "description": "description1",
*         "created": 1517891917,
*         "updated": 1517891917,
*         "users":  [{ User }, ...],
*         "goals":  [{ Goal }, ...],
*         "duration": "P1Y2M3DT4H5M6S",
*         "start":  1517891917,
*         "end":  1537891917,
*         "endTimeUnspecified":  true,
*         "recurrence":  {"RRule": "rule"},
*         "creator":  { User },
*         "collaborators":  [{ User }, ...],
*         "shares":  [{ User }, ...],
*         "days": {
*           "1": {
*             "items": [
*               {
*                 "primary": {
*                    "id": "pi123"
*                 },
*                 "optional": [
*                   {
*                      "id": "pi123"
*                   }
*                 ]
*               }
*             ]
*           }
*         },
*         "isTemplate": true,
*         "templateId": "template_id",
*         "status": 1,
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all plans successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plans were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.plan.AllPlans",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) AllPlans(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.All API request")
	req_plan := new(plan_proto.AllRequest)
	req_plan.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_plan.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_plan.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_plan.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_plan.SortParameter = req.Attribute(SortParameter).(string)
	req_plan.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	all_resp, err := p.PlanClient.All(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.AllPlans", "QueryError")
		return
	}
	all_resp.Code = http.StatusOK
	all_resp.Message = "Read all plans succesfully"
	data := utils.MarshalAny(rsp, all_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/plans/plan?session={session_id} Create a plan
* @apiVersion 0.1.0
* @apiName CreatePlan
* @apiGroup Plan
*
* @apiDescription Create a plan
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/plan?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "plan": {
*     "collaborators": [],
*     "creator": {
*       "addresses": [],
*       "contactDetails": [],
*       "created": "0",
*       "dob": "0",
*       "firstname": "",
*       "gender": "Gender_NONE",
*       "id": "222",
*       "image": "",
*       "lastname": "",
*       "org_id": "orgid",
*       "tokens": []
*     },
*     "days": {
*       "1": {
*         "items": [
*           {
*             "categoryIconSlug": "",
*             "categoryId": "",
*             "categoryName": "",
*             "contentId": "",
*             "contentPicUrl": "",
*             "contentTitle": "",
*             "id": "day_item_001",
*             "options": [],
*             "post": {
*               "created": "0",
*               "creator": {
*                 "addresses": [],
*                 "contactDetails": [],
*                 "created": "0",
*                 "dob": "0",
*                 "firstname": "",
*                 "gender": "Gender_NONE",
*                 "id": "userid",
*                 "image": "",
*                 "lastname": "",
*                 "org_id": "orgid",
*                 "tokens": [],
*                 "updated": "0"
*               },
*               "id": "111",
*               "items": [],
*               "name": "todo1",
*               "org_id": "orgid",
*               "updated": "0"
*             },
*             "pre": {
*               "created": "0",
*               "creator": {
*                 "addresses": [],
*                 "contactDetails": [],
*                 "created": "0",
*                 "dob": "0",
*                 "firstname": "",
*                 "gender": "Gender_NONE",
*                 "id": "userid",
*                 "image": "",
*                 "lastname": "",
*                 "org_id": "orgid",
*                 "tokens": [],
*                 "updated": "0"
*               },
*               "id": "111",
*               "items": [],
*               "name": "todo1",
*               "org_id": "orgid",
*               "updated": "0"
*             },
*             "primary": false,
*             "time": ""
*           }
*         ]
*       },
*       "2": {
*         "items": [
*           {
*             "categoryIconSlug": "",
*             "categoryId": "",
*             "categoryName": "",
*             "contentId": "",
*             "contentPicUrl": "",
*             "contentTitle": "",
*             "id": "day_item_002",
*             "options": [],
*             "post": {
*               "created": "0",
*               "creator": {
*                 "addresses": [],
*                 "contactDetails": [],
*                 "created": "0",
*                 "dob": "0",
*                 "firstname": "",
*                 "gender": "Gender_NONE",
*                 "id": "userid",
*                 "image": "",
*                 "lastname": "",
*                 "org_id": "orgid",
*                 "tokens": [],
*                 "updated": "0"
*               },
*               "id": "111",
*               "items": [],
*               "name": "todo1",
*               "org_id": "orgid",
*               "updated": "0"
*             },
*             "pre": {
*               "created": "0",
*               "creator": {
*                 "addresses": [],
*                 "contactDetails": [],
*                 "created": "0",
*                 "dob": "0",
*                 "firstname": "",
*                 "gender": "Gender_NONE",
*                 "id": "userid",
*                 "image": "",
*                 "lastname": "",
*                 "org_id": "orgid",
*                 "tokens": [],
*                 "updated": "0"
*               },
*               "id": "111",
*               "items": [],
*               "name": "todo1",
*               "org_id": "orgid",
*               "updated": "0"
*             },
*             "primary": false,
*             "time": ""
*           }
*         ]
*       }
*     },
*     "description": "hello world",
*     "duration": "",
*     "end": "0",
*     "endTimeUnspecified": false,
*     "goals": [
*        {
*          "articles": [],
*          "category": null,
*          "challenges": [],
*          "completionApprovalRequired": false,
*          "created": "0",
*          "createdBy": null,
*          "description": "",
*          "duration": "",
*          "habits": [],
*          "id": "1",
*          "image": "",
*          "notifications": [],
*          "org_id": "orgid",
*          "setbacks": [],
*          "social": [],
*          "source": "",
*          "status": "Status_NONE",
*          "successCriterias": [],
*          "summary": "",
*          "tags": [],
*          "target": null,
*          "title": "",
*          "trackers": [],
*          "triggers": [],
*          "updated": "0",
*          "users": [],
*          "visibility": "Visibility_NONE"
*        },
*        {
*          "articles": [],
*          "category": null,
*          "challenges": [],
*          "completionApprovalRequired": false,
*          "created": "0",
*          "createdBy": null,
*          "description": "",
*          "duration": "",
*          "habits": [],
*          "id": "2",
*          "image": "",
*          "notifications": [],
*          "org_id": "orgid",
*          "setbacks": [],
*          "social": [],
*          "source": "",
*          "status": "Status_NONE",
*          "successCriterias": [],
*          "summary": "",
*          "tags": [],
*          "target": null,
*          "title": "",
*          "trackers": [],
*          "triggers": [],
*          "updated": "0",
*          "users": [],
*          "visibility": "Visibility_NONE"
*        }
*      ],
*      "id": "1212",
*      "isTemplate": true,
*      "itemsCount": "2",
*      "linkSharingEnabled": false,
*      "name": "plan1",
*      "org_id": "orgid",
*      "pic": "",
*      "recurrence": [],
*      "shares": [
*        {
*          "addresses": [],
*          "contactDetails": [],
*          "created": "0",
*          "dob": "0",
*          "firstname": "",
*          "gender": "Gender_NONE",
*          "id": "userid",
*          "image": "",
*          "lastname": "",
*          "org_id": "orgid",
*          "tokens": [],
*          "updated": "0"
*        }
*      ],
*      "start": "0",
*      "status": "DRAFT",
*      "templateId": "template1",
*      "updated": "1523428756",
*      "users": [
*         {
*           "addresses": [],
*           "contactDetails": [],
*           "created": "0",
*           "dob": "0",
*           "firstname": "",
*           "gender": "Gender_NONE",
*           "id": "userid",
*           "image": "",
*           "lastname": "",
*           "org_id": "orgid",
*           "tokens": [],
*           "updated": "0"
*         }
*       ],
*      "setting": {
*        "embeddingEnabled": false,
*        "linkSharingEnabled": false,
*        "notifications": [],
*        "shareableLink": "",
*        "social": [],
*        "visibility": "PUBLIC"
*      },
*     }
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "plan": {
*       "id": "111",
*       "title": "plan1",
*       "org_id": "orgid",
*       "description": "description1",
*       "created": 1517891917,
*       "updated": 1517891917,
*       "users":  [{ User }, ...],
*       "goals":  [{ Goal }, ...],
*       "duration":  100,
*       "start":  1517891917,
*       "end":  1537891917,
*       "endTimeUnspecified":  true,
*       "visibility":  1,
*       "recurrence":  {"RRule": "rule"},
*       "creator":  { User },
*       "collaborators":  [{ User }, ...],
*       "shares":  [{ User }, ...],
*       "days": {
*         "1": {
*           "items": [
*             {
*               "primary": {
*                  "id": "pi123"
*               },
*               "optional": [
*                 {
*                    "id": "pi123"
*                 }
*               ]
*             }
*           ]
*         }
*       },
*       "isTemplate": true,
*       "templateId": "template_id",
*       "linkSharingEnabled": true,
*       "embeddingEnabled": true,
*       "shareableLink": "http://example.com",
*       "status": 1,
*     },
*   },
*   "code": 200,
*   "message": "Created plan succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plans were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.plan.AllPlans",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) CreatePlan(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.Create API request")
	req_plan := new(plan_proto.CreateRequest)
	// err := req.ReadEntity(req_plan)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.CreatePlan", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_plan); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.CreatePlan", "BindError")
		return
	}
	req_plan.UserId = req.Attribute(UserIdAttrName).(string)
	req_plan.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_plan.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	create_resp, err := p.PlanClient.Create(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.CreatePlan", "CreateError")
		return
	}
	create_resp.Code = http.StatusOK
	create_resp.Message = "Created plan succesfully"
	data := utils.MarshalAny(rsp, create_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/plans/plan/{plan_id}?session={session_id} View plan detail
* @apiVersion 0.1.0
* @apiName ReadPlan
* @apiGroup Plan
*
* @apiDescription View plan detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/plan/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "plan": {
*       "id": "111",
*       "title": "plan1",
*       "org_id": "orgid",
*       "description": "description1",
*       "created": 1517891917,
*       "updated": 1517891917,
*       "users":  [{ User }, ...],
*       "goals":  [{ Goal }, ...],
*       "duration":  100,
*       "start":  1517891917,
*       "end":  1537891917,
*       "endTimeUnspecified":  true,
*       "recurrence":  {"RRule": "rule"},
*       "creator":  { User },
*       "collaborators":  [{ User }, ...],
*       "shares":  [{ User }, ...],
*       "days": {
*         "1": {
*           "items": [
*             {
*               "primary": {
*                  "id": "pi123"
*               },
*               "optional": [
*                 {
*                    "id": "pi123"
*                 }
*               ]
*             }
*           ]
*         }
*       },
*       "isTemplate": true,
*       "templateId": "template_id",
*       "status": 1,
*     }
*   },
*   "code": 200,
*   "message": "Read plan succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plan were not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.plan.ReadPlan",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) ReadPlan(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.Read API request")
	req_plan := new(plan_proto.ReadRequest)
	req_plan.Id = req.PathParameter("plan_id")
	req_plan.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_plan.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	read_resp, err := p.PlanClient.Read(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.ReadPlan", "ReadError")
		return
	}
	read_resp.Code = http.StatusOK
	read_resp.Message = "Read plan succesfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, read_resp)
}

/**
* @api {delete} /server/plans/plan/{plan_id}?session={session_id} Delete a plan
* @apiVersion 0.1.0
* @apiName DeletePlan
* @apiGroup Plan
*
* @apiDescription Delete a plan
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/plan/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted plan succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plan was not updated.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.plan.DeletePlan",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) DeletePlan(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.Delete API request")
	req_plan := new(plan_proto.DeleteRequest)
	req_plan.Id = req.PathParameter("plan_id")
	req_plan.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_plan.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	delete_resp, err := p.PlanClient.Delete(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.DeletePlan", "DeletError")
		return
	}

	delete_resp.Code = http.StatusOK
	delete_resp.Message = "Deleted plan succesfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, delete_resp)
}

/**
* @api {post} /server/plans/search?session={session_id}&offset={offset}&limit={limit} Search plans
* @apiVersion 0.1.0
* @apiName SearchPlans
* @apiGroup Plan
*
* @apiDescription SearchPlans
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/search?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "name": "plan1"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "plans": [
*       {
*         "id": "111",
*         "title": "plan1",
*         "org_id": "orgid",
*         "description": "description1",
*         "created": 1517891917,
*         "updated": 1517891917,
*         "users":  [{ User }, ...],
*         "goals":  [{ Goal }, ...],
*         "duration": "P1Y2M3DT4H5M6S",
*         "start":  1517891917,
*         "end":  1537891917,
*         "endTimeUnspecified":  true,
*         "recurrence":  {"RRule": "rule"},
*         "creator":  { User },
*         "collaborators":  [{ User }, ...],
*         "shares":  [{ User }, ...],
*         "days": {
*           "1": {
*             "items": [
*               {
*                 "primary": {
*                    "id": "pi123"
*                 },
*                 "optional": [
*                   {
*                      "id": "pi123"
*                   }
*                 ]
*               }
*             ]
*           }
*         },
*         "isTemplate": true,
*         "templateId": "template_id",
*         "status": 1,
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched plans successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plans were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "SearchError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.plan.SearchPlans",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) SearchPlans(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.Search API request")
	req_plan := new(plan_proto.SearchRequest)
	err := req.ReadEntity(req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.SearchPlans", "BindError")
		return
	}
	req_plan.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_plan.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_plan.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_plan.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_plan.SortParameter = req.Attribute(SortParameter).(string)
	req_plan.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	search_resp, err := p.PlanClient.Search(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.SearchPlans", "SearchError")
		return
	}
	search_resp.Code = http.StatusOK
	search_resp.Message = "Searched plans succesfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, search_resp)
}

/**
* @api {get} /server/plans/templates?session={session_id}&offset={offset}&limit={limit} Get all templates
* @apiVersion 0.1.0
* @apiName Templates
* @apiGroup Plan
*
* @apiDescription Get all templates - Send all plans where isTemplate = true
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/templates?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "plans": [
*       {
*         "id": "111",
*         "title": "plan1",
*         "org_id": "orgid",
*         "description": "description1",
*         "created": 1517891917,
*         "updated": 1517891917,
*         "users":  [{ User }, ...],
*         "goals":  [{ Goal }, ...],
*         "duration": "P1Y2M3DT4H5M6S",
*         "start":  1517891917,
*         "end":  1537891917,
*         "endTimeUnspecified":  true,
*         "recurrence":  {"RRule": "rule"},
*         "creator":  { User },
*         "collaborators":  [{ User }, ...],
*         "shares":  [{ User }, ...],
*         "days": {
*           "1": {
*             "items": [
*               {
*                 "primary": {
*                    "id": "pi123"
*                 },
*                 "optional": [
*                   {
*                      "id": "pi123"
*                   }
*                 ]
*               }
*             ]
*           }
*         },
*         "isTemplate": true,
*         "templateId": "template_id",
*         "status": 1,
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched plans successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plans were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.plan.Templates",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) Templates(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.Templates API request")
	req_plan := new(plan_proto.TemplatesRequest)
	req_plan.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_plan.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_plan.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_plan.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_plan.SortParameter = req.Attribute(SortParameter).(string)
	req_plan.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	temp_resp, err := p.PlanClient.Templates(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.Templates", "QueryError")
		return
	}
	temp_resp.Code = http.StatusOK
	temp_resp.Message = "Searched plans succesfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, temp_resp)
}

/**
* @api {get} /server/plans/drafts?session={session_id}&offset={offset}&limit={limit} Get all draft plans
* @apiVersion 0.1.0
* @apiName Drafts
* @apiGroup Plan
*
* @apiDescription Get all draft plans - Get all plans where the status is draft
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/drafts?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "plans": [
*       {
*         "id": "111",
*         "title": "plan1",
*         "org_id": "orgid",
*         "description": "description1",
*         "created": 1517891917,
*         "updated": 1517891917,
*         "users":  [{ User }, ...],
*         "goals":  [{ Goal }, ...],
*         "duration": "P1Y2M3DT4H5M6S",
*         "start":  1517891917,
*         "end":  1537891917,
*         "endTimeUnspecified":  true,
*         "recurrence":  {"RRule": "rule"},
*         "creator":  { User },
*         "collaborators":  [{ User }, ...],
*         "shares":  [{ User }, ...],
*         "days": {
*           "1": {
*             "items": [
*               {
*                 "primary": {
*                    "id": "pi123"
*                 },
*                 "optional": [
*                   {
*                      "id": "pi123"
*                   }
*                 ]
*               }
*             ]
*           }
*         },
*         "isTemplate": true,
*         "templateId": "template_id",
*         "status": 1,
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched plans successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plans were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.plan.Drafts",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) Drafts(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.Drafts API request")
	req_plan := new(plan_proto.DraftsRequest)
	req_plan.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_plan.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_plan.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_plan.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_plan.SortParameter = req.Attribute(SortParameter).(string)
	req_plan.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	draft_resp, err := p.PlanClient.Drafts(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.Drafts", "QueryError")
		return
	}
	draft_resp.Code = http.StatusOK
	draft_resp.Message = "Searched plans succesfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, draft_resp)
}

/**
* @api {get} /server/plans/creator/{user_id}?session={session_id}&offset={offset}&limit={limit} Get all plans created by a particular team member
* @apiVersion 0.1.0
* @apiName ByCreator
* @apiGroup Plan
*
* @apiDescription Get all plans created by a particular team member - Get all plans where createdBy = {userid}
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/creator/{userid}?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "plans": [
*       {
*         "id": "111",
*         "title": "plan1",
*         "org_id": "orgid",
*         "description": "description1",
*         "created": 1517891917,
*         "updated": 1517891917,
*         "users":  [{ User }, ...],
*         "goals":  [{ Goal }, ...],
*         "duration": "P1Y2M3DT4H5M6S",
*         "start":  1517891917,
*         "end":  1537891917,
*         "endTimeUnspecified":  true,
*         "recurrence":  {"RRule": "rule"},
*         "creator":  { User },
*         "collaborators":  [{ User }, ...],
*         "shares":  [{ User }, ...],
*         "days": {
*           "1": {
*             "items": [
*               {
*                 "primary": {
*                    "id": "pi123"
*                 },
*                 "optional": [
*                   {
*                      "id": "pi123"
*                   }
*                 ]
*               }
*             ]
*           }
*         },
*         "isTemplate": true,
*         "templateId": "template_id",
*         "status": 1,
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched plans successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plans were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.plan.ByCreator",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) ByCreator(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.ByCreator API request")
	req_plan := new(plan_proto.ByCreatorRequest)
	req_plan.UserId = req.PathParameter("user_id")
	req_plan.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_plan.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_plan.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_plan.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_plan.SortParameter = req.Attribute(SortParameter).(string)
	req_plan.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	creator_resp, err := p.PlanClient.ByCreator(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.ByCreator", "QueryError")
		return
	}
	creator_resp.Code = http.StatusOK
	creator_resp.Message = "Searched plans succesfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, creator_resp)
}

/**
* @api {get} /server/plans/filters/all?session={session_id}&offset={offset}&limit={limit} Get all plan filters
* @apiVersion 0.1.0
* @apiName Filters
* @apiGroup Plan
*
* @apiDescription Get all plan filters
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/filters/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "filters": [
*       {
*         "displayName": "display",
*         "filterSlug": "slug",
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched plans successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The filters were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.plan.Filters",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) Filters(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.Filter API request")
	req_plan := new(plan_proto.FiltersRequest)
	req_plan.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_plan.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_plan.SortParameter = req.Attribute(SortParameter).(string)
	req_plan.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	filters_resp, err := p.PlanClient.Filters(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.Filters", "QueryError")
		return
	}

	filters_resp.Code = http.StatusOK
	filters_resp.Message = "Filtered plans succesfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, filters_resp)
}

// TopFilters

/**
* @api {get} /server/plans/filter/time?start_date={startDate}&end_date={endDate}&session={session_id}&offset={offset}&limit={limit} Get all plans by time period
* @apiVersion 0.1.0
* @apiName TimeFilters
* @apiGroup Plan
*
* @apiDescription Get all plans by time period - Return all plans that were created between a specific time period
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/filter/time?start_date=1517791917&end_date=1517991917&session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "plans": [
*       {
*         "id": "111",
*         "title": "plan1",
*         "org_id": "orgid",
*         "description": "description1",
*         "created": 1517891917,
*         "updated": 1517891917,
*         "users":  [{ User }, ...],
*         "goals":  [{ Goal }, ...],
*         "duration": "P1Y2M3DT4H5M6S",
*         "start":  1517891917,
*         "end":  1537891917,
*         "endTimeUnspecified":  true,
*         "visibility":  1,
*         "recurrence":  {"RRule": "rule"},
*         "creator":  { User },
*         "collaborators":  [{ User }, ...],
*         "shares":  [{ User }, ...],
*         "days": {
*           "1": {
*             "items": [
*               {
*                 "primary": {
*                    "id": "pi123"
*                 },
*                 "optional": [
*                   {
*                      "id": "pi123"
*                   }
*                 ]
*               }
*             ]
*           }
*         },
*         "isTemplate": true,
*         "templateId": "template_id",
*         "linkSharingEnabled": true,
*         "embeddingEnabled": true,
*         "shareableLink": "http://example.com",
*         "status": 1,
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Filtered plans successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plans were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.plan.TimeFilters",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) TimeFilters(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.TimeFilters API request")
	req_plan := new(plan_proto.TimeFiltersRequest)
	req_plan.StartDate, _ = strconv.ParseInt(req.QueryParameter("start_date"), 10, 64)
	req_plan.EndDate, _ = strconv.ParseInt(req.QueryParameter("end_date"), 10, 64)
	req_plan.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_plan.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_plan.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_plan.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_plan.SortParameter = req.Attribute(SortParameter).(string)
	req_plan.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	filters_resp, err := p.PlanClient.TimeFilters(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.TimeFilters", "QueryError")
		return
	}
	filters_resp.Code = http.StatusOK
	filters_resp.Message = "Filtered plans succesfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, filters_resp)
}

// UserFilters, SuccessFilters, ConditionFilters

/**
* @api {get} /server/plans/filter/goal?filter=[comma_separate_list_of_goal_category]&session={session_id}&offset={offset}&limit={limit} Get all plans by goal category
* @apiVersion 0.1.0
* @apiName GoalFilters
* @apiGroup Plan
*
* @apiDescription Get all plans by goal category - Return all plans filtered by one or more goal category
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/filter/goal?filter="1,2"&session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "plans": [
*       {
*         "id": "111",
*         "title": "plan1",
*         "org_id": "orgid",
*         "description": "description1",
*         "created": 1517891917,
*         "updated": 1517891917,
*         "users":  [{ User }, ...],
*         "goals":  [{ Goal }, ...],
*         "duration":  "P1Y2M3DT4H5M6S",
*         "start":  1517891917,
*         "end":  1537891917,
*         "endTimeUnspecified":  true,
*         "recurrence":  {"RRule": "rule"},
*         "creator":  { User },
*         "collaborators":  [{ User }, ...],
*         "shares":  [{ User }, ...],
*         "days": {
*           "1": {
*             "items": [
*               {
*                 "primary": {
*                    "id": "pi123"
*                 },
*                 "optional": [
*                   {
*                      "id": "pi123"
*                   }
*                 ]
*               }
*             ]
*           }
*         },
*         "isTemplate": true,
*         "templateId": "template_id",
*         "status": 1,
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Filtered plans successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The plans were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.plan.GoalFilters",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) GoalFilters(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.GoalFilters API request")
	req_plan := new(plan_proto.GoalFiltersRequest)
	req_plan.Goals = req.QueryParameter("filter")
	req_plan.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_plan.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_plan.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_plan.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_plan.SortParameter = req.Attribute(SortParameter).(string)
	req_plan.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	goal_resp, err := p.PlanClient.GoalFilters(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.GoalFilters", "QueryError")
		return
	}
	goal_resp.Code = http.StatusOK
	goal_resp.Message = "Filtered plans succesfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, goal_resp)
}

/**
* @api {post} /server/plans/plan/search/autocomplete?session={session_id} autocomplete text search for plans
* @apiVersion 0.1.0
* @apiName AutocompletePlanSearch
* @apiGroup Plan
*
* @apiDescription Should return a list of plans based on text based search. This should not be paginated
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/plan/search/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "title": "p",
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "id": "111",
*         "title": "plan",
*         "org_id": "orgid",
*       },
*       {
*         "id": "222",
*         "title": "ptx",
*         "org_id": "orgid",
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read plans successfully"
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
*           "domain": "go.micro.srv.plan.AutocompletePlanSearch",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) AutocompletePlanSearch(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.AutocompletePlanSearch API request")
	req_search := new(plan_proto.AutocompleteSearchRequest)
	req_search.SortParameter = req.Attribute(SortParameter).(string)
	req_search.SortDirection = req.Attribute(SortDirection).(string)

	err := req.ReadEntity(req_search)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.AutocompletePlanSearch", "BindError")
		return
	}
	// req_search.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_search.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.PlanClient.AutocompleteSearch(ctx, req_search)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.AutocompletePlanSearch", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read plans successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/plans/tags/top/{n}?session={session_id} Return top N tags for Plan
* @apiVersion 0.1.0
* @apiName GetTopPlanTags
* @apiGroup Plan
*
* @apiDescription For each of the following service we have return top N tags for plan
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/tags/top/5?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "tags": ["tag1","tag2","tag3",...]
*   },
*   "code": 200,
*   "message": "Get top plan tags successfully"
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
*           "domain": "go.micro.srv.plan.GetTopPlanTags",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) GetTopPlanTags(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.GetTopPlanTags API request")
	req_plan := new(plan_proto.GetTopTagsRequest)
	n, _ := strconv.Atoi(req.PathParameter("n"))
	req_plan.N = int64(n)
	req_plan.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_plan.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.PlanClient.GetTopTags(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.GetTopPlanTags", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Get top plan tags successfully"
	rsp.AddHeader("Plan-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/plans/tags/autocomplete?session={session_id} Autocomplete for tags for Plan
* @apiVersion 0.1.0
* @apiName AutocompletePlanTags
* @apiGroup Plan
*
* @apiDescription Autocomplete for tags for Plan
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/plans/tags/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
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
*   "message": "Autocomplete plan tags successfully"
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
*           "domain": "go.micro.srv.plan.AutocompletePlanTags",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *PlanService) AutocompletePlanTags(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Plan.AutocompletePlanTags API request")
	req_plan := new(plan_proto.AutocompleteTagsRequest)
	err := req.ReadEntity(req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.AutocompletePlanTags", "BindError")
		return
	}
	req_plan.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_plan.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.PlanClient.AutocompleteTags(ctx, req_plan)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.plan.AutocompletePlanTags", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Autocomplete plan tags successfully"
	rsp.AddHeader("Plan-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}
