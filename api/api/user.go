package api

import (
	"context"
	"net/http"
	account_proto "server/account-srv/proto/account"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	content_proto "server/content-srv/proto/content"
	userapp_proto "server/user-app-srv/proto/userapp"
	user_proto "server/user-srv/proto/user"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

type UserService struct {
	UserClient    user_proto.UserServiceClient
	UserAppClient userapp_proto.UserAppServiceClient
	ContentClient content_proto.ContentServiceClient
	Auth          Filters
	Audit         AuditFilter
	ServerMetrics metrics.Metrics
}

func (u UserService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/users")

	audit := &audit_proto.Audit{
		ActionService:  common.UserSrv,
		ActionResource: common.BASE + common.USER_TYPE,
	}

	ws.Route(ws.GET("/all").To(u.AllUsers).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		Filter(u.Auth.SortFilter).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List all users"))

	ws.Route(ws.POST("/user").To(u.CreateUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Create a user"))

	ws.Route(ws.POST("/user/share/multiple").To(u.ShareMultipleResources).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Share a resource/s this user"))

	ws.Route(ws.GET("/user/{user_id}").To(u.ReadUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Read a user"))

	// ws.Route(ws.POST("/filter").To(u.FilterUser).
	// 	Filter(u.Auth.BasicAuthenticate).
	// 	Filter(u.Auth.OrganisationAuthenticate).
	// 	Filter(u.Auth.EmployeeAuthenticate).
	// 	Filter(u.Auth.Paginate).
	// 	Doc("Filter user"))

	ws.Route(ws.POST("/user/content/share").To(u.ShareContent).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Share contents"))

	ws.Route(ws.GET("/user/{user_id}/preferences").To(u.ReadUserPreference).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get user preferences"))

	ws.Route(ws.GET("/user/{user_id}/feedback").To(u.ListUserFeedback).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("List user feedback"))

	ws.Route(ws.POST("/user/filter").To(u.FilterUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		Filter(u.Auth.SortFilter).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Filtering user by optional fields"))

	ws.Route(ws.POST("/user/search").To(u.SearchUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Search user"))

	ws.Route(ws.POST("/user/search/autocomplete").To(u.AutocompleteUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.SortFilter).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Autocomplete users"))

	ws.Route(ws.POST("/user/{user_id}/account/status").To(u.SetAccountStatus).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Set account status"))

	ws.Route(ws.GET("/user/{user_id}/account/status").To(u.GetAccountStatus).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Set account status"))

	ws.Route(ws.POST("/user/{user_id}/account/pass/reset").To(u.ResetUserPassword).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Reset / Update account password or passcode"))

	ws.Route(ws.POST("/user/{user_id}/measurements/measurement").To(u.AddMultipleMeasurements).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Add multiple measurements for a user"))

	ws.Route(ws.GET("/user/{user_id}/measurements/{marker_id}").To(u.GetMeasurementsHistory).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get measurements history for this user for specific marker"))

	ws.Route(ws.GET("/user/{user_id}/measurements/all").To(u.GetAllMeasurementsHistory).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get all measurements history for this user"))

	ws.Route(ws.GET("/user/{user_id}/markers/all").To(u.GetAllTrackedMarkers).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get all markers that are being tracked for this user"))

	ws.Route(ws.POST("/user/{user_id}/shared").To(u.GetSharedResources).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get resources shared with this user"))

	ws.Route(ws.POST("/user/{user_id}/shared/search").To(u.SearchSharedResources).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get resources shared with this user"))

	ws.Route(ws.POST("/user/{user_id}/share").To(u.ShareResources).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Share a resource/s this user"))

	ws.Route(ws.POST("/user/{user_id}/share/all").To(u.GetAllShareableResources).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get all shareable resources for this user"))

	ws.Route(ws.POST("/user/{user_id}/share/search").To(u.SearchShareableResources).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Search all shareable resources for this user"))

	ws.Route(ws.GET("/user/{user_id}/goals/current/progress").To(u.GetGoalProgress).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		Filter(u.Auth.Paginate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Get all markers that are being tracked for this user"))

	ws.Route(ws.POST("/user/{user_id}/delete").To(u.DeleteUser).
		Filter(u.Auth.BasicAuthenticate).
		Filter(u.Auth.OrganisationAuthenticate).
		Filter(u.Auth.EmployeeAuthenticate).
		// Filter(u.Audit.Clone(audit).AuditFilter).
		Doc("Deiete user by userid"))

	restful.Add(ws)
}

/**
* @api {get} /server/users/all List all user
* @apiVersion 0.1.0
* @apiName Users
* @apiGroup User
*
* @apiDescription List all users
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/users/all
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "users": [
*       {
*         "id": "userid",
*         "org_id": "orgid",
*         "firstname": "david",
*         "lastname": "john",
*         "image": "http://example.com",
*         "gender": "MALE",
*         "contact_details": { ContactDetail }
*         "addresses": [ Address, ... ],
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all users successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The users were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.AllUsers",
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
*           "domain": "go.micro.srv.user.AllUsers",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) AllUsers(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.AllUsers API request")
	req_user := new(user_proto.AllRequest)
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_user.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_user.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_user.SortParameter = req.Attribute(SortParameter).(string)
	req_user.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.All(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.AllUsers", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read all users successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/users/user?session={session_id} Create or update a user
* @apiVersion 0.1.0
* @apiName CreateUser
* @apiGroup User
*
* @apiDescription Create or update a user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "user": {
*     "org_id": "orgid",
*     "firstname": "david",
*     "lastname": "john",
*     "image": "http://example.com",
*     "gender": "MALE",
*     "contact_details": { ContactDetail }
*     "addresses": [ Address, ... ]
*   },
*   "account": {
*     "email": "email8@email.com",
*     "password": "pass1",
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "user": {
*       "id": "userid",
*       "org_id": "orgid",
*       "firstname": "david",
*       "lastname": "john",
*       "image": "http://example.com",
*       "gender": "MALE",
*       "contact_details": { ContactDetail }
*       "addresses": [ Address, ... ],
*       "created": 1517891917,
*       "updated": 1517891917
*     },
*	  "account" :{
*		"id":"accountid"
*	   }
*   },
*   "code": 200,
*   "message": "Created user successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, QueryError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.CreateUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserService) CreateUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.CreateUser API request")
	req_user := new(user_proto.CreateRequest)
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)

	if err := utils.UnmarshalAny(req, rsp, req_user); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.CreateUser", "BindError")
		return
	}
	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.Create(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.CreateUser", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created user successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/users/user/{user_id}?session={session_id} View user detail
* @apiVersion 0.1.0
* @apiName ReadUser
* @apiGroup User
*
* @apiDescription View user detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "user": {
*       "id": "userid",
*       "org_id": "orgid",
*       "firstname": "david",
*       "lastname": "john",
*       "image": "http://example.com",
*       "gender": "MALE",
*       "contact_details": { ContactDetail }
*       "addresses": [ Address, ... ],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read user successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.ReadUser",
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
*           "domain": "go.micro.srv.user.ReadUser",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) ReadUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.ReadUser API request")
	req_user := new(user_proto.ReadRequest)
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.Read(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.ReadUser", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read user successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/users/user/content/share?session={session_id} Get all shared content with a particular user so far
* @apiVersion 0.1.0
* @apiName ShareContent
* @apiGroup User
*
* @apiDescription API should return all content that has been shared with the normal user by any employee of the org
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/content/share?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
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
func (p *UserService) ShareContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.ShareContent API request")
	req_content := new(content_proto.ShareContentRequest)
	req_content.UserId = req.Attribute(UserIdAttrName).(string)
	req_content.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_content.TeamId = req.Attribute(TeamIdAttrName).(string)
	// err := req.ReadEntity(req_content)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.ShareContent", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_content); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.ShareContent", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.ContentClient.ShareContent(ctx, req_content)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.ShareContent", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Filter contents successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/users/user/{user_id}/preferences?session={session_id} Get user preferences
* @apiVersion 0.1.0
* @apiName ReadUserPreference
* @apiGroup User
*
* @apiDescription API should allow an employee to get user preferences
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/preferences?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "preference": {
*       "allergies": [
*         {
*           "category": null,
*           "created": "0",
*           "description": "",
*           "icon_slug": "",
*           "id": "",
*           "name": "name_2",
*           "name_slug": "name_slug_2",
*           "orgId": "",
*           "priority": "0",
*           "summary": "",
*           "tags": [],
*           "taxonomy": null,
*           "updated": "0",
*           "weight": "0"
*         }
*       ],
*       "conditions": [
*         {
*           "category": null,
*           "created": "0",
*           "description": "",
*           "icon_slug": "",
*           "id": "",
*           "name": "name_1",
*           "name_slug": "name_slug_1",
*           "orgId": "",
*           "priority": "0",
*           "summary": "",
*           "tags": [],
*           "taxonomy": null,
*           "updated": "0",
*           "weight": "0"
*         }
*       ],
*       "created": "1523701578",
*       "cuisines": [
*         {
*           "category": null,
*           "created": "0",
*           "description": "",
*           "icon_slug": "",
*           "id": "",
*           "name": "name_4",
*           "name_slug": "name_slug_4",
*           "orgId": "",
*           "priority": "0",
*           "summary": "",
*           "tags": [],
*           "taxonomy": null,
*           "updated": "0",
*           "weight": "0"
*         }
*       ],
*       "currentMeasurements": [
*         {
*           "created": "0",
*           "id": "measure_id",
*           "marker": null,
*           "measuredBy": null,
*           "method": null,
*           "orgId": "orgid",
*           "unit": "",
*           "updated": "0",
*           "userId": "userid",
*           "value": null
*         }
*       ],
*       "ethinicties": [
*         {
*           "category": null,
*           "created": "0",
*           "description": "",
*           "icon_slug": "",
*           "id": "",
*           "name": "name_5",
*           "name_slug": "name_slug_5",
*           "orgId": "",
*           "priority": "0",
*           "summary": "",
*           "tags": [],
*           "taxonomy": null,
*           "updated": "0",
*           "weight": "0"
*         }
*       ],
*       "food": [
*         {
*           "category": null,
*           "created": "0",
*           "description": "",
*           "icon_slug": "",
*           "id": "",
*           "name": "name_3",
*           "name_slug": "name_slug_3",
*           "orgId": "",
*           "priority": "0",
*           "summary": "",
*           "tags": [],
*           "taxonomy": null,
*           "updated": "0",
*           "weight": "0"
*         }
*       ],
*       "id": "449c6ec2-3fce-11e8-8085-20c9d0453b15",
*       "orgId": "orgid",
*       "updated": "1523701578",
*       "userId": "userid"
*     }
*   },
*   "code": 200,
*   "message": "Read user successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.ReadUserPreference",
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
*           "domain": "go.micro.srv.user.ReadUserPreference",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) ReadUserPreference(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.ReadUserPreference API request")
	req_user := new(user_proto.ReadUserPreferenceRequest)
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.ReadUserPreference(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.ReadUserPreference", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read user successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/users/user/{user_id}/feedback?session={session_id} List user feedback
* @apiVersion 0.1.0
* @apiName ListUserFeedback
* @apiGroup User
*
* @apiDescription API should allow an employee to get user feedback
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/feedback?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "feedbacks": {
*       "created": "1523718102",
*       "feedback": "very nice!!!",
*       "id": "be2150b6-3ff4-11e8-be21-20c9d0453b15",
*       "orgId": "orgid",
*       "rating": "9",
*       "updated": "1523718102",
*       "userId": "userid"
*     }
*   },
*   "code": 200,
*   "message": "List user feedback successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.ListUserFeedback",
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
*           "domain": "go.micro.srv.user.ListUserFeedback",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) ListUserFeedback(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.ListUserFeedback API request")
	req_user := new(user_proto.ListUserFeedbackRequest)
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.ListUserFeedback(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.ListUserFeedback", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "List user feedback successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/users/filter?session={session_id}&offset={offset}&limit={limit} Filtering user by optional fields
* @apiVersion 0.1.0
* @apiName Filter
* @apiGroup User
*
* @apiDescription Filter users by one or more of the optional fields as stated below
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "currentBatch": { Batch },
*   "status": Status,
*   "tags": ["tag111", "tag222", "tag333"],
*   "preferences": { Preferences }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "id": "userid",
*         "org_id": "orgid",
*         "firstname": "david",
*         "lastname": "john",
*         "avatar_url": "http://example.com",
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Filtered users successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, QueryError.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.FilterUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserService) FilterUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.Filter API request")
	req_filter := new(user_proto.FilterUserRequest)
	err := req.ReadEntity(req_filter)
	if err != nil {
		log.Infoln(err)
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.Filter", "BindError")
		return
	}
	req_filter.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_filter.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_filter.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_filter.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_filter.SortParameter = req.Attribute(SortParameter).(string)
	req_filter.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.FilterUser(ctx, req_filter)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.Filter", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Filtered users successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/users/user/search?session={session_id}&offset={offset}&limit={limit} Search user
* @apiVersion 0.1.0
* @apiName SearchUser
* @apiGroup User
*
* @apiDescription Simple Search users - Return searched list of users, should be paginated. This is different from autocomplete in terms of the POST body request that you receive
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/search?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "name": { Batch },
*   "gender": Status,
*   "dob": ["tag111", "tag222", "tag333"],
*   "addresses": [{Address}, ...],
*   "contact_details": [{ContactDetail}, ...]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "id": "userid",
*         "org_id": "orgid",
*         "firstname": "david",
*         "lastname": "john",
*         "avatar_url": "http://example.com",
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Searched user successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.SearchUser",
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
*           "domain": "go.micro.srv.user.SearchUser",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) SearchUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.SearchUser API request")
	req_user := new(user_proto.SearchUserRequest)
	err := req.ReadEntity(req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.Filter", "BindError")
		return
	}
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_user.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_user.Offset = req.Attribute(PaginateOffsetParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.SearchUser(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.SearchUser", "ReadError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Searched user successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/users/user/search/autocomplete?session={session_id} Autocomplete users
* @apiVersion 0.1.0
* @apiName AutocompleteUser
* @apiGroup User
*
* @apiDescription Should return a list of surveys based on text based search. This should not be paginated
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/search/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "name": "sub_string"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "id": "userid",
*         "org_id": "orgid",
*         "firstname": "david",
*         "lastname": "john",
*         "avatar_url": "http://example.com",
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Autocomplete users successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.AutocompleteUser",
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
*           "domain": "go.micro.srv.user.AutocompleteUser",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) AutocompleteUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.AutocompleteUser API request")
	req_user := new(user_proto.AutocompleteUserRequest)
	err := req.ReadEntity(req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.AutocompleteUser", "BindError")
		return
	}
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_user.SortParameter = req.Attribute(SortParameter).(string)
	req_user.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.AutocompleteUser(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.AutocompleteUser", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Autocomplete users successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/users/user/{user_id}/account/status?session={session_id} Set account status
* @apiVersion 0.1.0
* @apiName SetAccountStatus
* @apiGroup User
*
* @apiDescription API should allow an employee to update account status (lock, suspend etc.) of the normal user of an org
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/account/status?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "status": "SUSPEND"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Read user successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.SetAccountStatus",
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
*           "domain": "go.micro.srv.user.SetAccountStatus",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) SetAccountStatus(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.SetAccountStatus API request")
	req_user := new(account_proto.SetAccountStatusRequest)
	if err := utils.UnmarshalAny(req, rsp, req_user); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.SetAccountStatus", "BindError")
		return
	}
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)
	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.SetAccountStatus(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.SetAccountStatus", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Account status updated successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/users/user/{user_id}/account/status?session={session_id} Get  account status
* @apiVersion 0.1.0
* @apiName GetAccountStatus
* @apiGroup User
*
* @apiDescription API should allow an employee to get account status (lock, suspend etc.) of the normal user of an org
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/account/status?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "status": "SUSPEND"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*     "code": 200,
*     "data": {
*         "account": {
*             "status": "ACTIVE"
*         }
*     },
*     "message": "Account status returned successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.GetAccountStatus",
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
*           "domain": "go.micro.srv.user.GetAccountStatus",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) GetAccountStatus(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.GetAccountStatus API request")
	req_user := new(account_proto.GetAccountStatusRequest)

	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.GetAccountStatus(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.GetAccountStatus", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Account status returned successfully"

	data := utils.MarshalAny(rsp, resp)
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/users/user/{user_id}/account/pass/reset?session={session_id} Reset / Update account password or passcode
* @apiVersion 0.1.0
* @apiName ResetUserPassword
* @apiGroup User
*
* @apiDescription API should allow an employee to reset and send password or passcode to the normal user of an org
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/account/pass/reset?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example: With password
* {
*   "password": "password"
* }
*
* @apiParamExample {json} Request-Example: With passcode
* {
*   "passcode": "12345"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Reset password successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.ResetUserPassword",
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
*           "domain": "go.micro.srv.user.ResetUserPassword",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) ResetUserPassword(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.ResetUserPassword API request")
	req_user := new(account_proto.ResetUserPasswordRequest)
	err := req.ReadEntity(req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.ResetUserPassword", "BindError")
		return
	}
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)
	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.ResetUserPassword(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.ResetUserPassword", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Reset password successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/users/user/{user_id}/measurements/measurement?session={session_id} Add multiple measurements for a user
* @apiVersion 0.1.0
* @apiName AddMultipleMeasurements
* @apiGroup User
*
* @apiDescription save multiple measurements for a user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/measurements/measurement?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "measurements": [
*     {
*       "user_id": "userid",
*       "org_id": "orgid",
*       "marker": { Marker },
*       "value": "value",
*       "unit": "unit",
*     },
*     ... ...
*   ]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": {
*       "user_id": "userid",
*       "org_id": "orgid",
*       "measuredBy": { User },
*       "track_marker_id": "track_marker_id",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Add multiple measurements successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.AddMultipleMeasurements",
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
*           "domain": "go.micro.srv.user.AddMultipleMeasurements",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) AddMultipleMeasurements(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.AddMultipleMeasurements API request")
	req_user := new(user_proto.AddMultipleMeasurementsRequest)
	if err := utils.UnmarshalAny(req, rsp, req_user); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.content.SaveUserPreference", "BindError")
		return
	}
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.AddMultipleMeasurements(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.AddMultipleMeasurements", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Add multiple measurements successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/users/user/{user_id}/measurements/{marker_id}?session={session_id}&offset={offset}&limit={limit} Get all measurements history for this user
* @apiVersion 0.1.0
* @apiName GetMeasurementsHistory
* @apiGroup User
*
* @apiDescription This endpoint should be paginated and return historical measurements / measurement for this user for a specific marker
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/measurements/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "measurements": {
*       "user_id": "userid",
*       "org_id": "orgid",
*       "marker": { Marker },
*       "method": { TrackerMethod },
*       "measuredBy": { User },
*       "value": "string value"
*       "unit": "unit",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Get measurements history successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.GetMeasurementsHistory",
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
*           "domain": "go.micro.srv.user.GetMeasurementsHistory",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) GetMeasurementsHistory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.GetMeasurementsHistory API request")
	req_user := new(user_proto.GetMeasurementsHistoryRequest)
	req_user.UserId = req.PathParameter("user_id")
	req_user.MarkerId = req.PathParameter("marker_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.GetMeasurementsHistory(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.GetMeasurementsHistory", "ReadError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Get measurements history successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/users/user/{user_id}/measurements/all?session={session_id}&offset={offset}&limit={limit} Get all measurements history for this user
* @apiVersion 0.1.0
* @apiName GetAllMeasurementsHistory
* @apiGroup User
*
* @apiDescription This endpoint should be paginated and return historical measurements / measurement for this user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/measurements/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "measurements": {
*       "user_id": "userid",
*       "org_id": "orgid",
*       "marker": { Marker },
*       "method": { TrackerMethod },
*       "measuredBy": { User },
*       "value": "string value"
*       "unit": "unit",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Get all measurements history successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.GetAllMeasurementsHistory",
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
*           "domain": "go.micro.srv.user.GetAllMeasurementsHistory",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) GetAllMeasurementsHistory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.GetAllMeasurementsHistory API request")
	req_user := new(user_proto.GetAllMeasurementsHistoryRequest)
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.GetAllMeasurementsHistory(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.GetAllMeasurementsHistory", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get all measurements history successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/users/user/{user_id}/markers/all?session={session_id} Get all markers that are being tracked for this user
* @apiVersion 0.1.0
* @apiName GetAllTrackedMarkers
* @apiGroup User
*
* @apiDescription return a list of markers through all the measurements on the user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/markers/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "markers": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "org_id": "orgid",
*         "unit": "unit",
*         "apps": [ App, ...],
*         "wearables": [ Wearable, ...],
*         "devices": [ Device, ...],
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read user markers successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.GetAllTrackedMarkers",
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
*           "domain": "go.micro.srv.user.GetAllTrackedMarkers",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) GetAllTrackedMarkers(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.GetAllTrackedMarkers API request")
	req_user := new(user_proto.GetAllTrackedMarkersRequest)
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.GetAllTrackedMarkers(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.GetAllTrackedMarkers", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read user markers successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/users/user/{user_id}/shared?session={session_id} Get shared resources with this user
* @apiVersion 0.1.0
* @apiName GetSharedResources
* @apiGroup User
*
* @apiDescription return a list of resources shared with this user (shared goals, challenges, habits, surveys, content)
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/shared?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*	"type":["healum.com/proto/go.micro.srv.behaviour.Goal","healum.com/proto/go.micro.srv.behaviour.Challenge","healum.com/proto/go.micro.srv.behaviour.Habit"],
*	"status":["VIEWED","SHARED","RECEIVED","ACTIONED"],
*	"shared_by":["user_id","user_id"]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*     "code": 200,
*     "data": {
*         "shared_resources": [
*             {
*               "duration": "P0D",
*               "id": "6339f167-6993-11e8-ae1b-66a036430288",
*               "image": "iamge",
*               "resource_id": "aeda839e-5b46-11e8-b6ba-66a036430288",
*               "shared_by": "user name",
*               "title": "test",
*               "type": "healum.com/proto/go.micro.srv.behaviour.Goal",
* 				"shared_by_image":"some_image",
* 				"current":1,
* 				"target": 19,
* 				"duration" P20,
* 				"count":12
*             },
* 			...
*         ]
*     },
*     "message": "Returned list of shared resources successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.GetSharedResources",
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
*           "domain": "go.micro.srv.user.GetSharedResources",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) GetSharedResources(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.GetSharedResources API request")
	req_user := new(user_proto.GetSharedResourcesRequest)

	if err := utils.UnmarshalAny(req, rsp, req_user); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.GetSharedResources", "BindError")
		return
	}
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_user.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_user.Offset = req.Attribute(PaginateOffsetParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.GetSharedResources(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.GetSharedResources", "ReadError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Returned list of shared resources successfully"
	data := utils.MarshalAny(rsp, resp)
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/users/user/{user_id}/shared/search?session={session_id} Search shared resources with this user
* @apiVersion 0.1.0
* @apiName SearchSharedResources
* @apiGroup User
*
* @apiDescription return a list of resources shared with this user (shared goals, challenges, habits, surveys, content)
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/shared/search?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*	"type":["healum.com/proto/go.micro.srv.behaviour.Goal","healum.com/proto/go.micro.srv.behaviour.Challenge","healum.com/proto/go.micro.srv.behaviour.Habit"],
*	"status":["VIEWED","SHARED","RECEIVED","ACTIONED"],
*	"shared_by":["user_id","user_id"],
*	"query":"search term"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*     "code": 200,
*     "data": {
*         "shared_resources": [
*             {
*               "duration": "P0D",
*               "id": "6339f167-6993-11e8-ae1b-66a036430288",
*               "image": "iamge",
*               "resource_id": "aeda839e-5b46-11e8-b6ba-66a036430288",
*               "shared_by": "user name",
*               "title": "test",
*               "type": "healum.com/proto/go.micro.srv.behaviour.Goal",
* 				"shared_by_image":"some_image",
* 				"current":1,
* 				"target": 19,
* 				"duration" P20,
* 				"count":12
*             },
* 			...
*         ]
*     },
*     "message": "Returned list of shared resources successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.SearchSharedResources",
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
*           "domain": "go.micro.srv.user.SearchSharedResources",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) SearchSharedResources(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.SearchSharedResources API request")
	req_user := new(user_proto.GetSharedResourcesRequest)

	if err := utils.UnmarshalAny(req, rsp, req_user); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.SearchSharedResources", "BindError")
		return
	}
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_user.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_user.Offset = req.Attribute(PaginateOffsetParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.GetSharedResources(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.SearchSharedResources", "ReadError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Returned list of shared resources successfully"
	data := utils.MarshalAny(rsp, resp)
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/users/user/{user_id}/goals/current/progress?session={session_id} Get current goal progress of a particular user
* @apiVersion 0.1.0
* @apiName GetGoalProgress
* @apiGroup User
*
* @apiDescription Return progress on the current goal
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=/goals/current/progress?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "goal": { Goal },
*         "user": { User },
*         "latestValue": 4,
*         "target": 100,
*         "unit": Kgs
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Get goal progress sucessfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.GetGoalProgress",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserService) GetGoalProgress(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.GetGoalProgress API request")
	req_userapp := new(userapp_proto.GetGoalProgressRequest)
	req_userapp.UserId = req.PathParameter("user_id")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserAppClient.GetGoalProgress(ctx, req_userapp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.GetGoalProgress", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get goal progress sucessfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/users/user/{user_id}/share/all?session={session_id}&offset={offset}&limit={limit} List all resources available for sharing
* @apiVersion 0.1.0
* @apiName GetAllShareableResources
* @apiGroup user
*
* @apiDescription List all shareable resources
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/{user_id}/share/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*	"type":["healum.com/proto/go.micro.srv.behaviour.Goal","healum.com/proto/go.micro.srv.behaviour.Challenge","healum.com/proto/go.micro.srv.behaviour.Habit"],
*	"created_by":["user_id","user_id"]
* }

* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "resources": [
*       {
*         "id": "g111",
*         "title": "g_title",
*         "image": "http://image.com",
*         "summary": "summary",
*         "createdby": "david john"
*         "createdby_pic": "http://image.com/image",
*         "target": { Target }
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all shareable resources successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The resources were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.GetAllShareableResources",
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
*           "domain": "go.micro.srv.user.GetAllShareableResources",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) GetAllShareableResources(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.GetShareableResources API request")
	req_user := new(user_proto.GetShareableResourcesRequest)

	if err := utils.UnmarshalAny(req, rsp, req_user); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.GetAllShareableResources", "BindError")
		return
	}

	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_user.Offset = req.Attribute(PaginateOffsetParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.GetAllShareableResources(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.GetAllShareableResources", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all shareable resources successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/users/user/{user_id}/share/search?session={session_id}&offset={offset}&limit={limit} Search all resources available for sharing
* @apiVersion 0.1.0
* @apiName SearchShareableResources
* @apiGroup user
*
* @apiDescription Search all shareable resources
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/{user_id}/share/search?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*	"type":["healum.com/proto/go.micro.srv.behaviour.Goal","healum.com/proto/go.micro.srv.behaviour.Challenge","healum.com/proto/go.micro.srv.behaviour.Habit"],
*	"created_by":["user_id","user_id"],
*	"query":"test"
* }

* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "resources": [
*       {
*         "id": "g111",
*         "title": "g_title",
*         "image": "http://image.com",
*         "summary": "summary",
*         "createdby": "david john"
*         "createdby_pic": "http://image.com/image",
*         "target": { Target }
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all shareable resources successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The resources were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.SearchShareableResources",
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
*           "domain": "go.micro.srv.user.SearchShareableResources",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *UserService) SearchShareableResources(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.SearchShareableResources API request")
	req_user := new(user_proto.GetShareableResourcesRequest)

	if err := utils.UnmarshalAny(req, rsp, req_user); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.SearchShareableResources", "BindError")
		return
	}

	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_user.Offset = req.Attribute(PaginateOffsetParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.GetAllShareableResources(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.GetAllShareableResources", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Searched shareable resources successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/users/user/{user_id}/share?session={session_id}&offset={offset}&limit={limit} Share resources with a specific user
* @apiVersion 0.1.0
* @apiName ShareResources
* @apiGroup user
*
* @apiDescription Share resources with a specific user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/{user_id}/share?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*    "shares":{
*       "healum.com/proto/go.micro.srv.behaviour.Challenge":{
*          "resource":[
*       {
*                "resource_id":"f4d7357a-630b-11e8-809d-66a036430288",
*                "unit":"Kgs",
*                "currentValue":10,
*                "expectedProgress":"LINEAR"
*       },
*             {
*                "resource_id":"11f5cf5d-6271-11e8-809d-66a036430288",
*                "unit":"Kgs",
*                "currentValue":7,
*                "expectedProgress":"BELL"
*             }
*     ]
*   },
*       "healum.com/proto/go.micro.srv.behaviour.Goal":{
*          "resource":[
*             {
*                "resource_id":"f4d7357a-630b-11e8-809d-66a036430288",
*                "unit":"Kgs",
*                "currentValue":10,
*                "expectedProgress":"LINEAR"
*             },
*             {
*                "resource_id":"11f5cf5d-6271-11e8-809d-66a036430288",
*                "unit":"Kgs",
*                "currentValue":7,
*                "expectedProgress":"BELL"
*             }
*          ]
*       }
*    }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Shared resources successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The resources were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.ShareResources",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserService) ShareResources(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.ShareResources API request")
	req_user := new(user_proto.ShareResourcesRequest)

	if err := utils.UnmarshalAny(req, rsp, req_user); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.ShareResources", "BindError")
		return
	}

	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.ShareResources(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.ShareResources", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Shared resources successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/users/user/share/multiple?session={session_id} Share multiple resources with multiple user
* @apiVersion 0.1.0
* @apiName ShareMultipleResources
* @apiGroup user
*
* @apiDescription Share multiple resources with multiple users
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/user/share/multiple?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*    "shares":{
*       "healum.com/proto/go.micro.srv.behaviour.Challenge":{
*          "resource":["f4d7357a-630b-11e8-809d-66a036430288","f4d7357a-630b-11e8-809d-66a036430288","f4d7357a-630b-11e8-809d-66a036430288"]
*       },
*       "healum.com/proto/go.micro.srv.behaviour.Goal":{
*          "resource":["f4d7357a-630b-11e8-809d-66a036430288","f4d7357a-630b-11e8-809d-66a036430288","f4d7357a-630b-11e8-809d-66a036430288"]
*       }
*    },
*	"users":["f4d3557a-630b-11e8-809d-66a036430288","f4235357a-630b-11e8-809d-66a036430288","f42127a-630b-11e8-809d-66a036430288"]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Shared multiple resources successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The resources were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.ShareMultipleResources",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserService) ShareMultipleResources(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.ShareMultipleResources API request")
	req_user := new(user_proto.ShareMultipleResourcesRequest)

	if err := utils.UnmarshalAny(req, rsp, req_user); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.ShareMultipleResources", "BindError")
		return
	}

	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.ShareMultipleResources(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.user.ShareMultipleResources", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Shared multiple resources successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/users/{user_id}/delete?session={session_id} Delete a user
* @apiVersion 0.1.0
* @apiName UserUser
* @apiGroup User
*
* @apiDescription Delete a user
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/users/dsf89679-fsd234-s324d/delete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted user successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The user was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.user.DeleteUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *UserService) DeleteUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received User.DeleteUser API request")
	req_user := new(user_proto.DeleteRequest)
	req_user.UserId = req.PathParameter("user_id")
	req_user.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_user.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.UserClient.Delete(ctx, req_user)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.user.DeleteUser", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted user successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}
