package api

import (
	"context"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	organisation_proto "server/organisation-srv/proto/organisation"
	team_proto "server/team-srv/proto/team"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type TeamService struct {
	TeamClient         team_proto.TeamServiceClient
	Auth               Filters
	Audit              AuditFilter
	OrganisationClient organisation_proto.OrganisationServiceClient
	ServerMetrics      metrics.Metrics
}

func (p TeamService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/teams")

	audit := &audit_proto.Audit{
		ActionService:  common.TeamSrv,
		ActionResource: common.BASE + common.TEAM_TYPE,
	}

	ws.Route(ws.GET("/all").To(p.All).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all teams"))

	ws.Route(ws.POST("/team").To(p.Create).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("create data"))

	ws.Route(ws.GET("/team/{team_id}").To(p.Read).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("View Team detail"))

	ws.Route(ws.DELETE("/team/{team_id}").To(p.Delete).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delte Team detail"))

	ws.Route(ws.POST("/filter").To(p.Filter).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Filter teams"))

	ws.Route(ws.POST("/search").To(p.Search).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search teams"))

	ws.Route(ws.GET("/members/all").To(p.AllTeamMember).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all team members"))

	ws.Route(ws.POST("/members/member").To(p.CreateTeamMember).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create team member"))

	ws.Route(ws.GET("/members/member/{user_id}").To(p.ReadTeamMember).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read team member"))

	ws.Route(ws.POST("/members/member/{user_id}/modules").To(p.CreateEmployeeModuleAccess).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read team member"))

	ws.Route(ws.POST("/members/filter").To(p.FilterTeamMember).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Filter team members"))

	ws.Route(ws.POST("/employee/{employee_id}/delete").To(p.DeleteEmployee).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete employee"))

	restful.Add(ws)
}

// func fetchUserFromTeam(ctx context.Context, team *team_proto.Team, userClient user.AccountClient) {
// 	if team == nil {
// 		return
// 	}
// 	resp_user, err := userClient.Read(ctx, &user.ReadRequest{
// 		Id:    team.User.Id,
// 		Orgid: team.OrgId,
// 	})
// 	if err == nil {
// 		team.User = resp_user.User
// 	}
// }

/**
* @api {get} /server/teams/all?session={session_id}&offset={offset}&limit={limit} List all teams
* @apiVersion 0.1.0
* @apiName All
* @apiGroup Team
*
* @apiDescription All
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/teams/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "teams": [
*       {
*         "id": "111",
*         "name": "team1",
*         "description": "hello world",
*         "image": "image001",
*         "color": "red",
*         "products": [{Product}, ...],
*         "org_id": "orgid",
*         "user":  { User },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all teams successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The teams were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.team.All",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TeamService) All(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Team.All API request")
	req_team := new(team_proto.AllRequest)
	req_team.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_team.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_team.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_team.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_team.SortParameter = req.Attribute(SortParameter).(string)
	req_team.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TeamClient.All(ctx, req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.All", "QueryError")
		return
	}
	// fetching user object
	// for _, team := range resp.Data.Teams {
	// 	fetchUserFromTeam(ctx, team, p.Auth.UserClient)
	// }

	resp.Code = http.StatusOK
	resp.Message = "Read all teams successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/teams/team?session={session_id} Create or update a team
* @apiVersion 0.1.0
* @apiName Create
* @apiGroup Team
*
* @apiDescription Create or update a team
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/teams/team?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "team": {
*     "id": "111",
*     "name": "team1",
*     "description": "hello world",
*     "image": "image001",
*     "color": "red",
*     "products": [{Product}, ...],
*     "org_id": "orgid",
*     "user":  { User },
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "team": {
*       "id": "111",
*       "name": "team1",
*       "description": "hello world",
*       "image": "image001",
*       "color": "red",
*       "products": [{Product}, ...],
*       "org_id": "orgid",
*       "user":  { User },
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created team successfully"
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
*           "domain": "go.micro.srv.team.Create",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TeamService) Create(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Team.Create API request")
	req_team := new(team_proto.CreateRequest)
	// err := req.ReadEntity(req_team)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.Create", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_team); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.Create", "BindError")
		return
	}
	req_team.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_team.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TeamClient.Create(ctx, req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.Create", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created team successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/teams/team/{team_id}?session={session_id} View team detail
* @apiVersion 0.1.0
* @apiName Read
* @apiGroup Team
*
* @apiDescription View team detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/teams/team/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "team": {
*       "id": "111",
*       "name": "team1",
*       "description": "hello world",
*       "image": "image001",
*       "color": "red",
*       "products": [{Product}, ...],
*       "org_id": "orgid",
*       "user":  { User },
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read team successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The team was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.team.Read",
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
*           "domain": "go.micro.srv.team.Read",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *TeamService) Read(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Team.Read API request")
	req_team := new(team_proto.ReadRequest)
	req_team.Id = req.PathParameter("team_id")
	req_team.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_team.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TeamClient.Read(ctx, req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.Read", "ReadError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read team successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {delete} /server/teams/team/{team_id}?session={session_id} Delete a team
* @apiVersion 0.1.0
* @apiName Delete
* @apiGroup Team
*
* @apiDescription Delete a team
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/teams/team/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted team successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The team was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.team.Delete",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TeamService) Delete(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Team.Delete API request")
	req_team := new(team_proto.DeleteRequest)
	req_team.Id = req.PathParameter("team_id")
	req_team.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_team.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TeamClient.Delete(ctx, req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.Delete", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted team successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/teams/filter?session={session_id}&offset={offset}&limit={limit} Filter teams
* @apiVersion 0.1.0
* @apiName Filter
* @apiGroup Team
*
* @apiDescription Filter teams
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/teams/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "product": ["product1"]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "teams": [
*       {
*         "id": "111",
*         "name": "team1",
*         "description": "hello world",
*         "image": "image001",
*         "color": "red",
*         "products": [{Product}, ...],
*         "org_id": "orgid",
*         "user":  { User },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Filtered teams successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The teams were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.team.Filter",
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
*           "domain": "go.micro.srv.team.Filter",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *TeamService) Filter(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Team.Filter API request")
	req_team := new(team_proto.FilterRequest)
	err := req.ReadEntity(req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.Filter", "BindError")
		return
	}

	req_team.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_team.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_team.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_team.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_team.SortParameter = req.Attribute(SortParameter).(string)
	req_team.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TeamClient.Filter(ctx, req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.Filter", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Filtered teams successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/teams/search?session={session_id}&offset={offset}&limit={limit} Search teams
* @apiVersion 0.1.0
* @apiName Search
* @apiGroup Team
*
* @apiDescription Search teams
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/teams/search?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "team_name": "team1",
*   "team_member": "teammember"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "teams": [
*       {
*         "id": "111",
*         "name": "team1",
*         "description": "hello world",
*         "image": "image001",
*         "color": "red",
*         "products": [{Product}, ...],
*         "org_id": "orgid",
*         "user":  { User },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Searched teams successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The teams were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.team.Search",
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
*           "domain": "go.micro.srv.team.Search",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *TeamService) Search(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Team.Search API request")
	req_team := new(team_proto.SearchRequest)
	err := req.ReadEntity(req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.Filter", "BindError")
		return
	}

	req_team.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_team.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_team.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_team.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_team.SortParameter = req.Attribute(SortParameter).(string)
	req_team.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TeamClient.Search(ctx, req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.Search", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Searched teams successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/teams/memebers/all?session={session_id}&offset={offset}&limit={limit} List all team members
* @apiVersion 0.1.0
* @apiName AllTeamMemeber
* @apiGroup Team
*
* @apiDescription Employess can be in a team or without a any team membership. This ENDPOINT should return all employees that are in a team (with team details) + all employees who are not in any team
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/teams/memebers/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "employees": [
*       {
*         "id": "111",
*         "org_id": "orgid",
*         "role":  { Role },
*         "profile":  { EmployeeProfile },
*         "teams":  { TeamMembership },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all employees successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The teams were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.team.All",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TeamService) AllTeamMember(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Team.AllTeamMember API request")
	req_team := new(team_proto.AllTeamMemberRequest)
	req_team.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_team.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_team.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_team.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_team.SortParameter = req.Attribute(SortParameter).(string)
	req_team.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TeamClient.AllTeamMember(ctx, req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.AllTeamMember", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all employees successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/teams/members/member?session={session_id} Create or update a team member
* @apiVersion 0.1.0
* @apiName CreateTeamMember
* @apiGroup Team
*
* @apiDescription Create team member
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/teams/members/member?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "user": { User },
*   "account": { Account },
*   "Employee": { Employee },
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Created employee successfully"
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
*           "domain": "go.micro.srv.team.CreateTeamMember",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TeamService) CreateTeamMember(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Team.CreateTeamMember API request")
	req_team := new(team_proto.CreateTeamMemberRequest)
	if err := utils.UnmarshalAny(req, rsp, req_team); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.CreateTeamMember", "BindError")
		return
	}
	req_team.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_team.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TeamClient.CreateTeamMember(ctx, req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.CreateTeamMember", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created team member successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/teams/members/member/{user_id}?session={session_id} View team member detail
* @apiVersion 0.1.0
* @apiName ReadTeamMember
* @apiGroup Team
*
* @apiDescription View team member detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/teams/members/member/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "employee": {
*       "id": "111",
*       "org_id": "orgid",
*       "role":  { Role },
*       "profile":  { EmployeeProfile },
*       "teams":  { TeamMembership },
*       "created": 1517891917,
*       "updated": 1517891917
*     },
*     "user": { User },
*   },
*   "code": 200,
*   "message": "Read team successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The team was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.team.Read",
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
*           "domain": "go.micro.srv.team.Read",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *TeamService) ReadTeamMember(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Team.ReadTeamMember API request")
	req_team := new(team_proto.ReadTeamMemberRequest)
	req_team.UserId = req.PathParameter("user_id")
	req_team.OrgId = req.Attribute(OrgIdAttrName).(string)
	//req_team.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TeamClient.ReadTeamMember(ctx, req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.ReadTeamMember", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read team member successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/teams/members/member/{user_id}/modules?session={session_id} create or update team member module access
* @apiVersion 0.1.0
* @apiName CreateEmployeeModuleAccess
* @apiGroup Team
*
* @apiDescription create or update team member module access
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/teams/members/member/111/modules?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "modules": [{module},{module}]
* }
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Update employee module access successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	Update employee module access failed.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.team.CreateEmployeeModuleAccess",
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
*           "domain": "go.micro.srv.team.CreateEmployeeModuleAccess",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *TeamService) CreateEmployeeModuleAccess(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Team.CreateEmployeeModuleAccess API request")
	req_employee_module_access := new(team_proto.CreateEmployeeModuleAccessRequest)
	err := req.ReadEntity(req_employee_module_access)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.CreateEmployeeModuleAccess", "BindError")
		return
	}

	req_employee_module_access.UserId = req.PathParameter("user_id")
	req_employee_module_access.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TeamClient.CreateEmployeeModuleAccess(ctx, req_employee_module_access)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.CreateEmployeeModuleAccess", "WriteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Update employee module access successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/teams/memebers/filter?session={session_id}&offset={offset}&limit={limit} Filtering by optional fields
* @apiVersion 0.1.0
* @apiName Filter
* @apiGroup Team
*
* @apiDescription Filter members by one or more teams
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/teams/memebers/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "team": ["team1"]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "teams": [
*       {
*         "id": "111",
*         "org_id": "orgid",
*         "role":  { Role },
*         "profile":  { EmployeeProfile },
*         "teams":  { TeamMembership },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Filtered team members successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The team members were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.team.FilterTeamMember",
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
*           "domain": "go.micro.srv.team.FilterTeamMember",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *TeamService) FilterTeamMember(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Team.FilterTeamMember API request")
	req_team := new(team_proto.FilterTeamMemberRequest)
	err := req.ReadEntity(req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.FilterTeamMember", "BindError")
		return
	}

	req_team.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_team.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_team.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_team.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_team.SortParameter = req.Attribute(SortParameter).(string)
	req_team.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TeamClient.FilterTeamMember(ctx, req_team)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.team.FilterTeamMember", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Filtered team members successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/teams/employee/{employee_id}/delete?session={session_id} Delete an employee
* @apiVersion 0.1.0
* @apiName DeleteEmployee
* @apiGroup Team
*
* @apiDescription Delete an employee
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/teams/employee/dsf89679-fsd234-s324d/delete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted employee successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The employee was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.team.DeleteEmployee",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TeamService) DeleteEmployee(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Team.DeleteEmployee API request")
	req_delete := new(team_proto.DeleteEmployeeRequest)
	req_delete.EmployeeId = req.PathParameter("employee_id")
	req_delete.UserId = req.Attribute(UserIdAttrName).(string)
	req_delete.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TeamClient.DeleteEmployee(ctx, req_delete)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.team.DeleteEmployee", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted employee successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}
