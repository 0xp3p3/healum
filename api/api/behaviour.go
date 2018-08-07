package api

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	organisation_proto "server/organisation-srv/proto/organisation"
	static_proto "server/static-srv/proto/static"
	todo_proto "server/todo-srv/proto/todo"
	user_proto "server/user-srv/proto/user"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type BehaviourService struct {
	BehaviourClient    behaviour_proto.BehaviourServiceClient
	Auth               Filters
	Audit              AuditFilter
	ServerMetrics      metrics.Metrics
	OrganisationClient organisation_proto.OrganisationServiceClient
	StaticClient       static_proto.StaticServiceClient
}

func (p BehaviourService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/behaviours")

	audit := &audit_proto.Audit{
		ActionService:  common.BehaviourSrv,
		ActionResource: common.BASE + common.GOAL_TYPE,
	}

	ws.Route(ws.GET("/goals/all").To(p.AllGoals).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all goals"))

	ws.Route(ws.GET("/challenges/all").To(p.AllChallenges).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all challenges"))

	ws.Route(ws.GET("/habits/all").To(p.AllHabits).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all habits"))

	ws.Route(ws.POST("/goal").To(p.CreateGoal).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a goal"))

	ws.Route(ws.POST("/challenge").To(p.CreateChallenge).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a challenge"))

	ws.Route(ws.POST("/habit").To(p.CreateHabit).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a habit"))

	ws.Route(ws.GET("/goal/{goal_id}").To(p.ReadGoal).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a goal"))

	ws.Route(ws.GET("/challenge/{challenge_id}").To(p.ReadChallenge).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a challenge"))

	ws.Route(ws.GET("/habit/{habit_id}").To(p.ReadHabit).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a habit"))

	ws.Route(ws.DELETE("/goal/{goal_id}").To(p.DeleteGoal).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a goal"))

	ws.Route(ws.DELETE("/challenge/{challenge_id}").To(p.DeleteChallenge).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Doc("Delete a challenge"))

	ws.Route(ws.DELETE("/habit/{habit_id}").To(p.DeleteHabit).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a habit"))

	ws.Route(ws.POST("/filter").To(p.Filter).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Filter behaviours"))

	ws.Route(ws.POST("/search").To(p.Search).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search behaviours"))

	ws.Route(ws.POST("/search").To(p.Search).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search behaviours"))

	ws.Route(ws.POST("/goal/search/autocomplete").To(p.AutocompleteGoalSearch).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete goal text"))

	ws.Route(ws.POST("/challenge/search/autocomplete").To(p.AutocompleteChallengeSearch).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete challenge text"))

	ws.Route(ws.POST("/habit/search/autocomplete").To(p.AutocompleteHabitSearch).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete habit text"))

	ws.Route(ws.GET("/goals/tags/top/{n}").To(p.GetTopGoalTags).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Return top N tags for Goal"))

	ws.Route(ws.GET("/challenges/tags/top/{n}").To(p.GetTopChallengeTags).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Return top N tags for Challenge"))

	ws.Route(ws.GET("/habits/tags/top/{n}").To(p.GetTopHabitTags).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Return top N tags for Habit"))

	ws.Route(ws.POST("/goals/tags/autocomplete").To(p.AutocompleteGoalTags).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Autocomplete for tags for Goal"))

	ws.Route(ws.POST("/challenges/tags/autocomplete").To(p.AutocompleteChallengeTags).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Doc("Autocomplete for tags for Challenge"))

	ws.Route(ws.POST("/habits/tags/autocomplete").To(p.AutocompleteHabitTags).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Autocomplete for tags for Habit"))

	ws.Route(ws.POST("/goals/upload").To(p.UploadGoals).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Upload csv for Goals"))

	ws.Route(ws.POST("/challenges/upload").To(p.UploadChallenges).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Upload csv for Challenges"))

	ws.Route(ws.POST("/habits/upload").To(p.UploadHabits).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Upload csv for Habits"))

	restful.Add(ws)
}

/**
* @api {get} /server/behaviours/goals/all List all goals
* @apiVersion 0.1.0
* @apiName AllGoals
* @apiGroup Behaviour
*
* @apiDescription List all goals
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/goals/all
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "goals": [
*       {
*         "id": "g111",
*         "title": "g_title",
*         "org_id": "orgid",
*         "summary": "summary",
*         "description": "description",
*         "createdBy": { "id" : "userid", ...}
*         "status": 1,
*         "category": { "id" : "category111", ...}
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all goals successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The goals were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.behaviour.AllGoals",
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
*           "domain": "go.micro.srv.behaviour.AllGoals",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *BehaviourService) AllGoals(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.AllGoals API request")
	req_goal := new(behaviour_proto.AllGoalsRequest)
	req_goal.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_goal.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_goal.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_goal.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_goal.SortParameter = req.Attribute(SortParameter).(string)
	req_goal.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.AllGoals(ctx, req_goal)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AllGoals", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all goals successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/behaviours/challenges/all?session={session_id}&offset={offset}&limit={limit} List all challenges
* @apiVersion 0.1.0
* @apiName AllChallenges
* @apiGroup Behaviour
*
* @apiDescription List all challenges
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/challenges/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "challenges": [
*       {
*         "id": "c111",
*         "title": "c_title",
*         "org_id": "orgid",
*         "summary": "summary",
*         "description": "description",
*         "createdBy": { "id" : "userid", ...}
*         "status": 1,
*         "category": { "id" : "category222", ...}
*         "created": 1517891917,
*         "updated": 1517891917,
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all challenges successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenges were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.behaviour.AllChallenges",
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
*           "domain": "go.micro.srv.behaviour.AllChallenges",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *BehaviourService) AllChallenges(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.AllChallenges API request")
	req_challenge := new(behaviour_proto.AllChallengesRequest)
	req_challenge.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_challenge.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_challenge.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_challenge.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_challenge.SortParameter = req.Attribute(SortParameter).(string)
	req_challenge.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.AllChallenges(ctx, req_challenge)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AllChallenges", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read all challenges successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/behaviours/habits/all?session={session_id}&offset={offset}&limit={limit} List all habits
* @apiVersion 0.1.0
* @apiName AllHabits
* @apiGroup Behaviour
*
* @apiDescription List all habits
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/habits/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "habits": [
*       {
*         "id": "h111",
*         "title": "h_title",
*         "org_id": "orgid",
*         "summary": "summary",
*         "description": "description",
*         "createdBy": { "id" : "userid", ...}
*         "status": 1,
*         "category": { "id" : "category333", ...}
*         "created": 1517891917,
*         "updated": 1517891917,
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all habits successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The habits were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.behaviour.AllHabits",
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
*           "domain": "go.micro.srv.behaviour.AllHabits",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *BehaviourService) AllHabits(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.AllHabits API request")
	req_habit := new(behaviour_proto.AllHabitsRequest)
	req_habit.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_habit.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_habit.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_habit.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_habit.SortParameter = req.Attribute(SortParameter).(string)
	req_habit.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.AllHabits(ctx, req_habit)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AllHabits", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read all habits successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/behaviours/goal?session={session_id} Create or update a goal
* @apiVersion 0.1.0
* @apiName CreateGoal
* @apiGroup Behaviour
*
* @apiDescription Create or update a goal
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/goal?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "goal": {
*     "id": "g111",
*     "title": "g_title",
*     "org_id": "orgid",
*     "summary": "summary",
*     "description": "description",
*     "createdBy": { "id" : "userid", ...}
*     "status": 1,
*     "category": { "id" : "category111", ...}
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "goal": {
*       "id": "g111",
*       "title": "g_title",
*       "org_id": "orgid",
*       "summary": "summary",
*       "description": "description",
*       "createdBy": { "id" : "userid", ...}
*       "status": 1,
*       "category": { "id" : "category111", ...}
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created goal successfully"
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
*           "domain": "go.micro.srv.behaviour.CreateGoal",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) CreateGoal(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.CreateGoal API request")
	req_goal := new(behaviour_proto.CreateGoalRequest)
	// err := req.ReadEntity(req_goal)
	// if err != nil {
	// 	utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.CreateGoal", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_goal); err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.CreateGoal", "BindError")
		return
	}
	req_goal.UserId = req.Attribute(UserIdAttrName).(string)
	req_goal.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_goal.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.CreateGoal(ctx, req_goal)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.CreateGoal", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created goal successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/behaviours/challenge?session={session_id} Create or update a challenge
* @apiVersion 0.1.0
* @apiName CreateChallenge
* @apiGroup Behaviour
*
* @apiDescription Create or update a challenge
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/challenge?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "challenge": {
*     "id": "c111",
*     "title": "c_title",
*     "org_id": "orgid",
*     "summary": "summary",
*     "description": "description",
*     "createdBy": { "id" : "userid", ...}
*     "status": 1,
*     "category": { "id" : "category222", ...}
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "challenge": {
*       "id": "c111",
*       "title": "c_title",
*       "org_id": "orgid",
*       "summary": "summary",
*       "description": "description",
*       "createdBy": { "id" : "userid", ...}
*       "status": 1,
*       "category": { "id" : "category222", ...}
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created challenge successfully"
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
*           "domain": "go.micro.srv.behaviour.CreateChallenge",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) CreateChallenge(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.CreateChallenge API request")
	req_challenge := new(behaviour_proto.CreateChallengeRequest)
	// err := req.ReadEntity(req_challenge)
	// if err != nil {
	// 	utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.CreateChallenge", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_challenge); err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.CreateChallenge", "BindError")
		return
	}
	req_challenge.UserId = req.Attribute(UserIdAttrName).(string)
	req_challenge.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_challenge.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.CreateChallenge(ctx, req_challenge)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.CreateChallenge", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created challenge successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/behaviours/habit?session={session_id} Create or update a habit
* @apiVersion 0.1.0
* @apiName CreateHabit
* @apiGroup Behaviour
*
* @apiDescription Create or update a habit
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/habit?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "habit": {
*     "id": "h111",
*     "title": "h_title",
*     "org_id": "orgid",
*     "summary": "summary",
*     "description": "description",
*     "createdBy": { "id" : "userid", ...}
*     "status": 1,
*     "category": { "id" : "category333", ...}
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "habit": {
*       "id": "h111",
*       "title": "h_title",
*       "org_id": "orgid",
*       "summary": "summary",
*       "description": "description",
*       "createdBy": { "id" : "userid", ...}
*       "status": 1,
*       "category": { "id" : "category333", ...}
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created habit successfully"
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
*           "domain": "go.micro.srv.behaviour.CreateHabit",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) CreateHabit(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.CreateHabit API request")
	req_habit := new(behaviour_proto.CreateHabitRequest)
	// err := req.ReadEntity(req_habit)
	// if err != nil {
	// 	utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.CreateHabit", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_habit); err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.CreateHabit", "BindError")
		return
	}
	req_habit.UserId = req.Attribute(UserIdAttrName).(string)
	req_habit.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_habit.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.CreateHabit(ctx, req_habit)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.CreateHabit", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created habit successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/behaviours/goal/{goal_id}?session={session_id} View goal detail
* @apiVersion 0.1.0
* @apiName ReadGoal
* @apiGroup Behaviour
*
* @apiDescription View goal detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/goal/g111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "goal": {
*       "id": "g111",
*       "title": "g_title",
*       "org_id": "orgid",
*       "summary": "summary",
*       "description": "description",
*       "createdBy": { "id" : "userid", ...}
*       "status": 1,
*       "category": { "id" : "category111", ...}
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read goal successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The goal was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.behaviour.ReadGoal",
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
*           "domain": "go.micro.srv.behaviour.ReadGoal",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *BehaviourService) ReadGoal(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.ReadGoal API request")
	req_goal := new(behaviour_proto.ReadGoalRequest)
	req_goal.GoalId = req.PathParameter("goal_id")
	req_goal.OrgId = req.Attribute(OrgIdAttrName).(string)
	//req_goal.TeamId = req.Attribute(TeamIdAttrName).(string) - not used, can be removed

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.ReadGoal(ctx, req_goal)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.ReadGoal", "ReadError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read goal successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/behaviours/challenge/{challenge_id}?session={session_id} View challenge detail
* @apiVersion 0.1.0
* @apiName ReadChallenge
* @apiGroup Behaviour
*
* @apiDescription View challenge detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/challenge/c111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "challenge": {
*       "id": "c111",
*       "title": "c_title",
*       "org_id": "orgid",
*       "summary": "summary",
*       "description": "description",
*       "createdBy": { "id" : "userid", ...}
*       "status": 1,
*       "category": { "id" : "category222", ...}
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read challenge successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.behaviour.ReadChallenge",
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
*           "domain": "go.micro.srv.behaviour.ReadChallenge",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *BehaviourService) ReadChallenge(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.ReadChallenge API request")
	req_challenge := new(behaviour_proto.ReadChallengeRequest)
	req_challenge.ChallengeId = req.PathParameter("challenge_id")
	req_challenge.OrgId = req.Attribute(OrgIdAttrName).(string)
	//req_challenge.TeamId = req.Attribute(TeamIdAttrName).(string) - not used, can be removed

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.ReadChallenge(ctx, req_challenge)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.ReadChallenge", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read challenge successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/behaviours/habit/{habit_id}?session={session_id} View habit detail
* @apiVersion 0.1.0
* @apiName ReadHabit
* @apiGroup Behaviour
*
* @apiDescription View habit detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/habit/h111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "habit": {
*       "id": "h111",
*       "title": "h_title",
*       "org_id": "orgid",
*       "summary": "summary",
*       "description": "description",
*       "createdBy": { "id" : "userid", ...}
*       "status": 1,
*       "category": { "id" : "category333", ...}
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read habit successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The habit was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.behaviour.ReadHabit",
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
*           "domain": "go.micro.srv.behaviour.ReadHabit",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *BehaviourService) ReadHabit(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.ReadHabit API request")
	req_habit := new(behaviour_proto.ReadHabitRequest)
	req_habit.HabitId = req.PathParameter("habit_id")
	req_habit.OrgId = req.Attribute(OrgIdAttrName).(string)
	//req_habit.TeamId = req.Attribute(TeamIdAttrName).(string) - not used, can be removed

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.ReadHabit(ctx, req_habit)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.ReadHabit", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read habit successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {delete} /server/behaviours/goal/{goal_id}?session={session_id} Delete a goal
* @apiVersion 0.1.0
* @apiName DeleteGoal
* @apiGroup Behaviour
*
* @apiDescription Delete a goal
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/goal/g111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted goal successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The goal was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.behaviour.DeleteGoal",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) DeleteGoal(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.DeleteGoal API request")
	req_goal := new(behaviour_proto.DeleteGoalRequest)
	req_goal.GoalId = req.PathParameter("goal_id")
	req_goal.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_goal.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.DeleteGoal(ctx, req_goal)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.DeleteGoal", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted goal successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/behaviours/challenge/{challenge_id}?session={session_id} Delete a challenge
* @apiVersion 0.1.0
* @apiName DeleteChallenge
* @apiGroup Behaviour
*
* @apiDescription Delete a challenge
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/challenge/c111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted challenge successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The challenge was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.behaviour.DeleteChallenge",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) DeleteChallenge(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.DeleteChallenge API request")
	req_challenge := new(behaviour_proto.DeleteChallengeRequest)
	req_challenge.ChallengeId = req.PathParameter("challenge_id")
	req_challenge.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_challenge.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.DeleteChallenge(ctx, req_challenge)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.DeleteChallenge", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted challenge successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/behaviours/habit/{habit_id}?session={session_id} Delete a habit
* @apiVersion 0.1.0
* @apiName DeleteHabit
* @apiGroup Behaviour
*
* @apiDescription Delete a habit
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/habit/h111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted habit successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The habit was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.behaviour.DeleteHabit",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) DeleteHabit(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.DeleteHabit API request")
	req_habit := new(behaviour_proto.DeleteHabitRequest)
	req_habit.HabitId = req.PathParameter("habit_id")
	req_habit.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_habit.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.DeleteHabit(ctx, req_habit)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.DeleteHabit", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted habit successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/behaviours/filter?session={session_id}&offset={offset}&limit={limit} Filtering by optional fields
* @apiVersion 0.1.0
* @apiName Filter
* @apiGroup Behaviour
*
* @apiDescription Filter behaviours by one or more status, one or more type, one or more category
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "type": ["goal", "challenge", "habit"],
*   "status": [1],
*   "category": ["category111", "category222", "category333"],
*   "creator": ["userid"]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "goals": [
*       {
*         "id": "g111",
*         "title": "g_title",
*         "org_id": "orgid",
*         "summary": "summary",
*         "description": "description",
*         "createdBy": { "id" : "userid", ...}
*         "status": 1,
*         "category": { "id" : "category111", ...}
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ],
*     "challenges": [
*       {
*         "id": "c111",
*         "title": "c_title",
*         "org_id": "orgid",
*         "summary": "summary",
*         "description": "description",
*         "createdBy": { "id" : "userid", ...}
*         "status": 1,
*         "category": { "id" : "category222", ...}
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ],
*     "habits": [
*       {
*         "id": "h111",
*         "title": "h_title",
*         "org_id": "orgid",
*         "summary": "summary",
*         "description": "description",
*         "createdBy": { "id" : "userid", ...}
*         "status": 1,
*         "category": { "id" : "category333", ...}
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Filtered behaviours successfully"
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
*           "domain": "go.micro.srv.behaviour.Filter",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) Filter(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.Filter API request")
	req_filter := new(behaviour_proto.FilterRequest)
	// err := req.ReadEntity(req_filter)
	// if err != nil {
	// 	utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.Filter", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_filter); err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.Filter", "BindError")
		return
	}
	req_filter.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_filter.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_filter.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_filter.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_filter.SortParameter = req.Attribute(SortParameter).(string)
	req_filter.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.Filter(ctx, req_filter)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.Filter", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Filtered behaviours successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/behaviours/search?session={session_id}&offset={offset}&limit={limit} Simple Search behaviours
* @apiVersion 0.1.0
* @apiName Search
* @apiGroup Behaviour
*
* @apiDescription Simple Search behaviours - Return searched behaviours
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/search?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "name": "title",
*   "description": "descript",
*   "summary": "summary"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "goals": [
*       {
*         "id": "g111",
*         "title": "g_title",
*         "org_id": "orgid",
*         "summary": "summary",
*         "description": "description",
*         "createdBy": { "id" : "userid", ...}
*         "status": 1,
*         "category": { "id" : "category111", ...}
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ],
*     "challenges": [
*       {
*         "id": "c111",
*         "title": "c_title",
*         "org_id": "orgid",
*         "summary": "summary",
*         "description": "description",
*         "createdBy": { "id" : "userid", ...}
*         "status": 1,
*         "category": { "id" : "category222", ...}
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ],
*     "habits": [
*       {
*         "id": "h111",
*         "title": "h_title",
*         "org_id": "orgid",
*         "summary": "summary",
*         "description": "description",
*         "createdBy": { "id" : "userid", ...}
*         "status": 1,
*         "category": { "id" : "category333", ...}
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Filtered behaviours successfully"
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
*           "domain": "go.micro.srv.behaviour.Search",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) Search(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.Search API request")
	req_search := new(behaviour_proto.SearchRequest)
	// err := req.ReadEntity(req_search)
	// if err != nil {
	// 	utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.Search", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_search); err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.Search", "BindError")
		return
	}
	req_search.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_search.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_search.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_search.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_search.SortParameter = req.Attribute(SortParameter).(string)
	req_search.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.Search(ctx, req_search)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.Search", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Filtered behaviours successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/behaviours/goal/search/autocomplete?session={session_id} autocomplete text search for goals
* @apiVersion 0.1.0
* @apiName AutocompleteGoalSearch
* @apiGroup Behaviour
*
* @apiDescription Should return a list of challenges based on text based search. This should not be paginated
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/goal/search/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "title": "t",
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "id": "g111",
*         "title": "g_title",
*         "org_id": "orgid",
*       },
*       {
*         "id": "g222",
*         "title": "ttx",
*         "org_id": "orgid",
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read goals successfully"
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
*           "domain": "go.micro.srv.behaviour.AutocompleteGoalSearch",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) AutocompleteGoalSearch(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.AutocompleteGoalSearch API request")
	req_search := new(behaviour_proto.AutocompleteSearchRequest)
	err := req.ReadEntity(req_search)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AutocompleteGoalSearch", "BindError")
		return
	}

	// req_search.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_search.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.AutocompleteGoalSearch(ctx, req_search)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AutocompleteGoalSearch", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read goals successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/behaviours/challenge/search/autocomplete?session={session_id} autocomplete text search for challenges
* @apiVersion 0.1.0
* @apiName AutocompleteChallengeSearch
* @apiGroup Behaviour
*
* @apiDescription Should return a list of challenges based on text based search. This should not be paginated
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/challenge/search/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "title": "t",
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "id": "c111",
*         "title": "c_title",
*         "org_id": "orgid",
*       },
*       {
*         "id": "c222",
*         "title": "ttx",
*         "org_id": "orgid",
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read challenges successfully"
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
*           "domain": "go.micro.srv.behaviour.AutocompleteChallengeSearch",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) AutocompleteChallengeSearch(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.AutocompleteChallengeSearch API request")
	req_search := new(behaviour_proto.AutocompleteSearchRequest)
	err := req.ReadEntity(req_search)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AutocompleteChallengeSearch", "BindError")
		return
	}
	// req_search.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_search.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.AutocompleteChallengeSearch(ctx, req_search)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AutocompleteChallengeSearch", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read challenges successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/behaviours/habit/search/autocomplete?session={session_id} autocomplete text search for habits
* @apiVersion 0.1.0
* @apiName AutocompleteHabitSearch
* @apiGroup Behaviour
*
* @apiDescription Should return a list of habits based on text based search. This should not be paginated
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/habit/search/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "title": "t",
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": [
*       {
*         "id": "h111",
*         "title": "h_title",
*         "org_id": "orgid",
*       },
*       {
*         "id": "g222",
*         "title": "ttx",
*         "org_id": "orgid",
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read habits successfully"
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
*           "domain": "go.micro.srv.behaviour.AutocompleteHabitSearch",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) AutocompleteHabitSearch(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.AutocompleteHabitSearch API request")
	req_search := new(behaviour_proto.AutocompleteSearchRequest)
	err := req.ReadEntity(req_search)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AutocompleteHabitSearch", "BindError")
		return
	}
	// req_search.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_search.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.AutocompleteHabitSearch(ctx, req_search)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AutocompleteHabitSearch", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read habits successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/behaviours/goals/tags/top/{n}?session={session_id} Return top N tags for Goal
* @apiVersion 0.1.0
* @apiName GetTopGoalTags
* @apiGroup Behaviour
*
* @apiDescription For each of the following service we have return top N tags for goal
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/goals/tags/top/5?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "tags": ["tag1","tag2","tag3",...]
*   },
*   "code": 200,
*   "message": "Get top goal tags successfully"
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
*           "domain": "go.micro.srv.behaviour.GetTopGoalTags",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) GetTopGoalTags(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.GetTopGoalTags API request")
	req_behaviour := new(behaviour_proto.GetTopTagsRequest)
	n, _ := strconv.Atoi(req.PathParameter("n"))
	req_behaviour.N = int64(n)
	req_behaviour.Object = common.GOAL
	req_behaviour.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_behaviour.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.GetTopTags(ctx, req_behaviour)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.GetTopGoalTags", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Get top goal tags successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/behaviours/challenges/tags/top/{n}?session={session_id} Return top N tags for Challenge
* @apiVersion 0.1.0
* @apiName GetTopChallengeTags
* @apiGroup Behaviour
*
* @apiDescription For each of the following service we have return top N tags for Challenge
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/challenges/tags/top/5?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "tags": ["tag1","tag2","tag3",...]
*   },
*   "code": 200,
*   "message": "Get top challenge tags successfully"
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
*           "domain": "go.micro.srv.behaviour.GetTopChallengeTags",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) GetTopChallengeTags(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.GetTopChallengeTags API request")
	req_behaviour := new(behaviour_proto.GetTopTagsRequest)
	n, _ := strconv.Atoi(req.PathParameter("n"))
	req_behaviour.N = int64(n)
	req_behaviour.Object = common.CHALLENGE
	req_behaviour.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_behaviour.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.GetTopTags(ctx, req_behaviour)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.GetTopChallengeTags", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Get top challenge tags successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/behaviours/habits/tags/top/{n}?session={session_id} Return top N tags for Habit
* @apiVersion 0.1.0
* @apiName GetTopHabitTags
* @apiGroup Behaviour
*
* @apiDescription For each of the following service we have return top N tags for habit
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/habits/tags/top/5?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "tags": ["tag1","tag2","tag3",...]
*   },
*   "code": 200,
*   "message": "Get top habit tags successfully"
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
*           "domain": "go.micro.srv.behaviour.GetTopHabitTags",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) GetTopHabitTags(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.GetTopHabitTags API request")
	req_behaviour := new(behaviour_proto.GetTopTagsRequest)
	n, _ := strconv.Atoi(req.PathParameter("n"))
	req_behaviour.N = int64(n)
	req_behaviour.Object = common.HABIT
	req_behaviour.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_behaviour.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.GetTopTags(ctx, req_behaviour)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.GetTopHabitTags", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Get top habit tags successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/behaviours/goals/tags/autocomplete?session={session_id} Autocomplete for tags for Goal
* @apiVersion 0.1.0
* @apiName AutocompleteGoalTags
* @apiGroup Behaviour
*
* @apiDescription Autocomplete for tags for Goal
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/goals/tags/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
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
*   "message": "Autocomplete goal tags successfully"
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
*           "domain": "go.micro.srv.behaviour.AutocompleteGoalTags",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) AutocompleteGoalTags(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.AutocompleteGoalTags API request")
	req_behaviour := new(behaviour_proto.AutocompleteTagsRequest)
	err := req.ReadEntity(req_behaviour)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AutocompleteGoalTags", "BindError")
		return
	}
	req_behaviour.Object = common.GOAL
	req_behaviour.OrgId = req.Attribute(OrgIdAttrName).(string)

	// req_behaviour.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.AutocompleteTags(ctx, req_behaviour)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AutocompleteGoalTags", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Autocomplete goal tags successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/behaviours/challenges/tags/autocomplete?session={session_id} Autocomplete for tags for Challenge
* @apiVersion 0.1.0
* @apiName AutocompleteChallengeTags
* @apiGroup Behaviour
*
* @apiDescription Autocomplete for tags for Challenge
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/challenges/tags/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
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
*   "message": "Autocomplete challenge tags successfully"
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
*           "domain": "go.micro.srv.behaviour.AutocompleteChallengeTags",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) AutocompleteChallengeTags(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.AutocompleteChallengeTags API request")
	req_behaviour := new(behaviour_proto.AutocompleteTagsRequest)
	err := req.ReadEntity(req_behaviour)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AutocompleteChallengeTags", "BindError")
		return
	}
	req_behaviour.Object = common.CHALLENGE
	req_behaviour.OrgId = req.Attribute(OrgIdAttrName).(string)

	// req_behaviour.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.AutocompleteTags(ctx, req_behaviour)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AutocompleteChallengeTags", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Autocomplete challenge tags successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/behaviours/habits/tags/autocomplete?session={session_id} Autocomplete for tags for Habit
* @apiVersion 0.1.0
* @apiName AutocompleteHabitTags
* @apiGroup Behaviour
*
* @apiDescription Autocomplete for tags for Habit
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/habits/tags/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
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
*   "message": "Autocomplete habit tags successfully"
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
*           "domain": "go.micro.srv.behaviour.AutocompleteHabitTags",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) AutocompleteHabitTags(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.AutocompleteHabitTags API request")
	req_behaviour := new(behaviour_proto.AutocompleteTagsRequest)
	err := req.ReadEntity(req_behaviour)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AutocompleteHabitTags", "BindError")
		return
	}
	req_behaviour.Object = common.HABIT
	req_behaviour.OrgId = req.Attribute(OrgIdAttrName).(string)

	// req_behaviour.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.BehaviourClient.AutocompleteTags(ctx, req_behaviour)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.AutocompleteHabitTags", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Autocomplete habit tags successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/behaviours/goals/upload?session={session_id} Upload csv for Goals
* @apiVersion 0.1.0
* @apiName UploadGoals
* @apiGroup Behaviour
*
* @apiDescription Upload csv for Goals
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/goals/upload?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "upload_file": FileData
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Upload goal csv successfully"
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
*           "domain": "go.micro.srv.behaviour.UploadGoals",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) UploadGoals(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.UploadGoals API request")
	req.Request.ParseMultipartForm(32 << 20)
	file, _, err := req.Request.FormFile("upload_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	req_goal := new(behaviour_proto.CreateGoalRequest)
	req_goal.UserId = req.Attribute(UserIdAttrName).(string)
	req_goal.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_goal.TeamId = req.Attribute(TeamIdAttrName).(string)

	r := csv.NewReader(file)
	fields := map[int]string{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		// property parsing
		if len(fields) == 0 {
			for i, col := range row {
				fields[i] = col
			}
			continue
		}
		// row parsing
		for i, col := range row {
			goal := &behaviour_proto.Goal{
				OrgId:     req_goal.OrgId,
				CreatedBy: &user_proto.User{Id: req_goal.UserId},
				Target:    &static_proto.Target{},
				Todos: &todo_proto.Todo{
					Items: []*todo_proto.TodoItem{},
				},
			}
			switch fields[i] {
			case "title":
				goal.Title = col
			case "sumarry":
				goal.Summary = col
			case "description":
				goal.Description = col
			case "image":
				goal.Image = col
			case "target.aim":
				// read category with name_slug by col
				resp, err := p.StaticClient.ReadBehaviourCategoryAimByNameslug(ctx, &static_proto.ReadByNameslugRequest{NameSlug: col})
				if err != nil {
					continue
				}
				goal.Target.Aim = resp.Data.BehaviourCategoryAim
			case "target.marker":
				resp, err := p.StaticClient.ReadMarkerByNameslug(ctx, &static_proto.ReadByNameslugRequest{NameSlug: col})
				if err != nil {
					continue
				}
				goal.Target.Marker = resp.Data.Marker
			case "target.targetValue":
				v, _ := strconv.Atoi(col)
				goal.Target.TargetValue = int64(v)
			case "target.unit":
				goal.Target.Unit = col
			case "target.recurrence":
				goal.Target.Recurrence = []*static_proto.Recurrence{{col}}
			case "source":
				goal.Source = col
			case "duration":
				goal.Duration = col
			case "todos.items.title":
				goal.Todos.Items = append(goal.Todos.Items, &todo_proto.TodoItem{Title: col})
			}

			req_goal.Goal = goal
			_, err = p.BehaviourClient.CreateGoal(ctx, req_goal)
			if err != nil {
				utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.CreateGoal", "CreateGoalError")
				return
			}
		}
	}
	fmt.Println("finished goal created")
	resp := &behaviour_proto.UploadGoalsResponse{
		Code:    http.StatusOK,
		Message: "Upload goal csv successfully",
	}
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/behaviours/challenges/upload?session={session_id} Upload csv for Challenges
* @apiVersion 0.1.0
* @apiName UploadChallenges
* @apiGroup Behaviour
*
* @apiDescription Upload csv for Challenges
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/challenges/upload?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "upload_file": FileData
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Upload challenge csv successfully"
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
*           "domain": "go.micro.srv.behaviour.UploadChallenges",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) UploadChallenges(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.UploadChallenges API request")
	req.Request.ParseMultipartForm(32 << 20)
	file, _, err := req.Request.FormFile("upload_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	req_challenge := new(behaviour_proto.CreateChallengeRequest)
	req_challenge.UserId = req.Attribute(UserIdAttrName).(string)
	req_challenge.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_challenge.TeamId = req.Attribute(TeamIdAttrName).(string)

	r := csv.NewReader(file)
	fields := map[int]string{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		// property parsing
		if len(fields) == 0 {
			for i, col := range row {
				fields[i] = col
			}
			continue
		}
		// row parsing
		for i, col := range row {
			challenge := &behaviour_proto.Challenge{
				OrgId:     req_challenge.OrgId,
				CreatedBy: &user_proto.User{Id: req_challenge.UserId},
				Target:    &static_proto.Target{},
				Todos: &todo_proto.Todo{
					Items: []*todo_proto.TodoItem{},
				},
			}
			switch fields[i] {
			case "title":
				challenge.Title = col
			case "sumarry":
				challenge.Summary = col
			case "description":
				challenge.Description = col
			case "image":
				challenge.Image = col
			case "target.aim":
				// read category with name_slug by col
				resp, err := p.StaticClient.ReadBehaviourCategoryAimByNameslug(ctx, &static_proto.ReadByNameslugRequest{NameSlug: col})
				if err != nil {
					continue
				}
				challenge.Target.Aim = resp.Data.BehaviourCategoryAim
			case "target.marker":
				resp, err := p.StaticClient.ReadMarkerByNameslug(ctx, &static_proto.ReadByNameslugRequest{NameSlug: col})
				if err != nil {
					continue
				}
				challenge.Target.Marker = resp.Data.Marker
			case "target.targetValue":
				v, _ := strconv.Atoi(col)
				challenge.Target.TargetValue = int64(v)
			case "target.unit":
				challenge.Target.Unit = col
			case "target.recurrence":
				challenge.Target.Recurrence = []*static_proto.Recurrence{{col}}
			case "source":
				challenge.Source = col
			case "duration":
				challenge.Duration = col
			case "todos.items.title":
				challenge.Todos.Items = append(challenge.Todos.Items, &todo_proto.TodoItem{Title: col})
			}

			req_challenge.Challenge = challenge
			_, err = p.BehaviourClient.CreateChallenge(ctx, req_challenge)
			if err != nil {
				utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.CreateChallenge", "CreateChallengeError")
				return
			}
		}
	}
	fmt.Println("finished challenge created")
	resp := &behaviour_proto.UploadGoalsResponse{
		Code:    http.StatusOK,
		Message: "Upload challenge csv successfully",
	}
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/behaviours/habits/upload?session={session_id} Upload csv for Habits
* @apiVersion 0.1.0
* @apiName UploadHabits
* @apiGroup Behaviour
*
* @apiDescription Upload csv for Habits
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/behaviours/habits/upload?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "upload_file": FileData
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Upload habit csv successfully"
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
*           "domain": "go.micro.srv.behaviour.UploadHabits",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *BehaviourService) UploadHabits(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Behaviour.UploadHabits API request")

	req.Request.ParseMultipartForm(32 << 20)
	file, _, err := req.Request.FormFile("upload_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	req_habit := new(behaviour_proto.CreateHabitRequest)
	req_habit.UserId = req.Attribute(UserIdAttrName).(string)
	req_habit.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_habit.TeamId = req.Attribute(TeamIdAttrName).(string)

	r := csv.NewReader(file)
	fields := map[int]string{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		// property parsing
		if len(fields) == 0 {
			for i, col := range row {
				fields[i] = col
			}
			continue
		}
		// row parsing
		for i, col := range row {
			habit := &behaviour_proto.Habit{
				OrgId:     req_habit.OrgId,
				CreatedBy: &user_proto.User{Id: req_habit.UserId},
				Target:    &static_proto.Target{},
				Todos: &todo_proto.Todo{
					Items: []*todo_proto.TodoItem{},
				},
			}
			switch fields[i] {
			case "title":
				habit.Title = col
			case "sumarry":
				habit.Summary = col
			case "description":
				habit.Description = col
			case "image":
				habit.Image = col
			case "target.aim":
				// read category with name_slug by col
				resp, err := p.StaticClient.ReadBehaviourCategoryAimByNameslug(ctx, &static_proto.ReadByNameslugRequest{NameSlug: col})
				if err != nil {
					continue
				}
				habit.Target.Aim = resp.Data.BehaviourCategoryAim
			case "target.marker":
				resp, err := p.StaticClient.ReadMarkerByNameslug(ctx, &static_proto.ReadByNameslugRequest{NameSlug: col})
				if err != nil {
					continue
				}
				habit.Target.Marker = resp.Data.Marker
			case "target.targetValue":
				v, _ := strconv.Atoi(col)
				habit.Target.TargetValue = int64(v)
			case "target.unit":
				habit.Target.Unit = col
			case "target.recurrence":
				habit.Target.Recurrence = []*static_proto.Recurrence{{col}}
			case "source":
				habit.Source = col
			case "duration":
				habit.Duration = col
			case "todos.items.title":
				habit.Todos.Items = append(habit.Todos.Items, &todo_proto.TodoItem{Title: col})
			}

			req_habit.Habit = habit
			_, err = p.BehaviourClient.CreateHabit(ctx, req_habit)
			if err != nil {
				utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.CreateHabit", "CreateHabitError")
				return
			}
		}
	}
	fmt.Println("finished habit created")
	resp := &behaviour_proto.UploadGoalsResponse{
		Code:    http.StatusOK,
		Message: "Upload habit csv successfully",
	}
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}
