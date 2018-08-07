package api

import (
	"context"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	organisation_proto "server/organisation-srv/proto/organisation"
	track_proto "server/track-srv/proto/track"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type TrackService struct {
	TrackClient        track_proto.TrackServiceClient
	Auth               Filters
	Audit              AuditFilter
	OrganisationClient organisation_proto.OrganisationServiceClient
	ServerMetrics      metrics.Metrics
}

func (p TrackService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/track")

	audit := &audit_proto.Audit{
		ActionService:  common.TrackSrv,
		ActionResource: common.BASE + common.TRACK_TYPE,
	}

	ws.Route(ws.POST("/goal/{goal_id}").To(p.CreateTrackGoal).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Created track goal"))

	ws.Route(ws.GET("/goal/{goal_id}/count").To(p.GetGoalCount).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get goal count"))

	ws.Route(ws.GET("/goal/{goal_id}/history").To(p.GetGoalHistory).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get goal history"))

	ws.Route(ws.POST("/challenge/{challenge_id}").To(p.CreateTrackChallenge).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Created track challenge"))

	ws.Route(ws.GET("/challenge/{challenge_id}/count").To(p.GetChallengeCount).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get challenge count"))

	ws.Route(ws.GET("/challenge/{challenge_id}/history").To(p.GetChallengeHistory).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get challenge history"))

	ws.Route(ws.POST("/habit/{habit_id}").To(p.CreateTrackHabit).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Created track habit"))

	ws.Route(ws.GET("/habit/{habit_id}/count").To(p.GetHabitCount).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get habit count"))

	ws.Route(ws.GET("/habit/{habit_id}/history").To(p.GetHabitHistory).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get habit history"))

	ws.Route(ws.POST("/content/{content_id}").To(p.CreateTrackContent).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Created track content"))

	ws.Route(ws.GET("/content/{content_id}/count").To(p.GetContentCount).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get content count"))

	ws.Route(ws.GET("/content/{content_id}/history").To(p.GetContentHistory).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get content history"))

	ws.Route(ws.POST("/marker/{marker_id}").To(p.CreateTrackMarker).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create track marker"))

	ws.Route(ws.GET("/marker/{marker_id}").To(p.GetLastMarker).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create track marker"))

	ws.Route(ws.GET("/marker/{marker_id}/history").To(p.GetMarkerHistory).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create track marker"))

	ws.Route(ws.GET("/marker/history/all").To(p.GetAllMarkerHistory).
		Filter(p.Auth.BasicAuthenticate).
		// Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create track marker"))
	restful.Add(ws)
}

/**
* @api {post} /server/track/goal/{goal_id}?session={session_id} Submit an event for a Goal
* @apiVersion 0.1.0
* @apiName CreateTrackGoal
* @apiGroup Track
*
* @apiDescription Submit an event for a Goal
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/goal/g111?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="
*
* @apiParamExample {json} Request-Example:
* {
*   "user": { User }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "track_goal": {
*       "user": { User },
*       "goal": { Goal },
*       "orgid": "orgid",
*       "created": 1517891917
*     },
*     "count": 3
*   },
*   "code": 200,
*   "message": "Created TrackGoal successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track goal was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.CreateTrackGoal",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) CreateTrackGoal(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.CreateTrackGoal API request")
	req_track := new(track_proto.CreateTrackGoalRequest)
	// err := req.ReadEntity(req_track)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackGoal", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_track); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackGoal", "BindError")
		return
	}
	req_track.GoalId = req.PathParameter("goal_id")
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.CreateTrackGoal(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackGoal", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created TrackGoal successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/track/goal/{goal_id}/count?session={session_id}&from={from}&to={to} Get goal count
* @apiVersion 0.1.0
* @apiName GetGoalCount
* @apiGroup Track
*
* @apiDescription Get goal count
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/goal/g111/count?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="&from=13245320&to=135435340
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "count": 3
*   },
*   "code": 200,
*   "message": "Created TrackGoal successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track goal was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.GetGoalCount",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) GetGoalCount(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.GetGoalCount API request")
	req_track := new(track_proto.GetGoalCountRequest)

	req_track.GoalId = req.PathParameter("goal_id")
	req_track.UserId = req.Attribute(UserIdAttrName).(string)
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_track.From = req.Attribute(PaginateFromParameter).(int64)
	req_track.To = req.Attribute(PaginateToParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.GetGoalCount(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackGoal", "CreateError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Get goal count successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/track/goal/{goal_id}/history?session={session_id}&from={from}&to={to}&offset={offset}&limit={limit} Get goal history
* @apiVersion 0.1.0
* @apiName GetGoalHistory
* @apiGroup Track
*
* @apiDescription Get goal history
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/goal/g111/count?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="&from=13245320&to=135435340&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "track_goals": [
*       {
*         "user": { User },
*         "goal": { Goal },
*         "orgid": "orgid",
*         "created": 1517891917
*     	},
*     	... ...
*     ],
*   },
*   "code": 200,
*   "message": "Created TrackGoal successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track goal was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.GetGoalHistory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) GetGoalHistory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.GetGoalHistory API request")
	req_track := new(track_proto.GetGoalHistoryRequest)

	req_track.GoalId = req.PathParameter("goal_id")
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_track.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_track.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_track.From = req.Attribute(PaginateFromParameter).(int64)
	req_track.To = req.Attribute(PaginateToParameter).(int64)
	req_track.SortParameter = req.Attribute(SortParameter).(string)
	req_track.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.GetGoalHistory(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.GetGoalHistory", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get goal history successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/track/challenge/{challenge_id}?session={session_id} Submit an event for a Challenge
* @apiVersion 0.1.0
* @apiName CreateTrackChallenge
* @apiGroup Track
*
* @apiDescription Submit an event for a Challenge
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/challenge/c111?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="
*
* @apiParamExample {json} Request-Example:
* {
*   "user": { User }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "track_challenge": {
*       "user": { User },
*       "challenge": { Challenge },
*       "orgid": "orgid",
*       "created": 1517891917
*     },
*     "count": 3
*   },
*   "code": 200,
*   "message": "Created TrackChallenge successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track challenge was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.CreateTrackChallenge",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) CreateTrackChallenge(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.CreateTrackChallenge API request")
	req_track := new(track_proto.CreateTrackChallengeRequest)
	// err := req.ReadEntity(req_track)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackChallenge", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_track); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackChallenge", "BindError")
		return
	}
	req_track.ChallengeId = req.PathParameter("challenge_id")
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.CreateTrackChallenge(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackChallenge", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created TrackChallenge successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/track/challenge/{challenge_id}/count?session={session_id}&from={from}&to={to} Get challenge count
* @apiVersion 0.1.0
* @apiName GetChallengeCount
* @apiGroup Track
*
* @apiDescription Get challenge count
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/challenge/c111/count?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="&from=13245320&to=135435340
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "count": 3
*   },
*   "code": 200,
*   "message": "Created TrackChallenge successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track challenge was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.GetChallengeCount",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) GetChallengeCount(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.GetChallengeCount API request")
	req_track := new(track_proto.GetChallengeCountRequest)

	req_track.ChallengeId = req.PathParameter("challenge_id")
	req_track.UserId = req.Attribute(UserIdAttrName).(string)
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_track.From = req.Attribute(PaginateFromParameter).(int64)
	req_track.To = req.Attribute(PaginateToParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.GetChallengeCount(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackChallenge", "CreateError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Get challenge count successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/track/challenge/{challenge_id}/history?session={session_id}&from={from}&to={to}&offset={offset}&limit={limit} Get challenge history
* @apiVersion 0.1.0
* @apiName GetChallengeHistory
* @apiGroup Track
*
* @apiDescription Get challenge history
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/challenge/c111/count?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="&from=13245320&to=135435340&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "track_challenges": [
*       {
*         "user": { User },
*         "challenge": { Challenge },
*         "orgid": "orgid",
*         "created": 1517891917
*     	},
*     	... ...
*     ],
*   },
*   "code": 200,
*   "message": "Created TrackChallenge successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track challenge was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.GetChallengeHistory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) GetChallengeHistory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.GetChallengeHistory API request")
	req_track := new(track_proto.GetChallengeHistoryRequest)

	req_track.ChallengeId = req.PathParameter("challenge_id")
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_track.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_track.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_track.From = req.Attribute(PaginateFromParameter).(int64)
	req_track.To = req.Attribute(PaginateToParameter).(int64)
	req_track.SortParameter = req.Attribute(SortParameter).(string)
	req_track.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.GetChallengeHistory(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.GetChallengeHistory", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get challenge history successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/track/habit/{habit_id}?session={session_id} Submit an event for a Habit
* @apiVersion 0.1.0
* @apiName CreateTrackHabit
* @apiGroup Track
*
* @apiDescription Submit an event for a Habit
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/habit/h111?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="
*
* @apiParamExample {json} Request-Example:
* {
*   "user": { User }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "track_habit": {
*       "user": { User },
*       "habit": { Habit },
*       "orgid": "orgid",
*       "created": 1517891917
*     },
*     "count": 3
*   },
*   "code": 200,
*   "message": "Created TrackHabit successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track habit was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.CreateTrackHabit",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) CreateTrackHabit(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.CreateTrackHabit API request")
	req_track := new(track_proto.CreateTrackHabitRequest)
	// err := req.ReadEntity(req_track)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackHabit", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_track); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackHabit", "BindError")
		return
	}

	req_track.HabitId = req.PathParameter("habit_id")
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.CreateTrackHabit(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackHabit", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created TrackHabit successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/track/habit/{habit_id}/count?session={session_id}&from={from}&to={to} Get habit count
* @apiVersion 0.1.0
* @apiName GetHabitCount
* @apiGroup Track
*
* @apiDescription Get habit count
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/habit/h111/count?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="&from=13245320&to=135435340
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "count": 3
*   },
*   "code": 200,
*   "message": "Created TrackHabit successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track habit was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.GetHabitCount",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) GetHabitCount(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.GetHabitCount API request")
	req_track := new(track_proto.GetHabitCountRequest)

	req_track.HabitId = req.PathParameter("habit_id")
	req_track.UserId = req.Attribute(UserIdAttrName).(string)
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_track.From = req.Attribute(PaginateFromParameter).(int64)
	req_track.To = req.Attribute(PaginateToParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.GetHabitCount(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackHabit", "CreateError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Get habit count successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/track/habit/{habit_id}/history?session={session_id}&from={from}&to={to}&offset={offset}&limit={limit} Get habit history
* @apiVersion 0.1.0
* @apiName GetHabitHistory
* @apiGroup Track
*
* @apiDescription Get habit history
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/habit/h111/count?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="&from=13245320&to=135435340&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "track_habits": [
*       {
*         "user": { User },
*         "habit": { Habit },
*         "orgid": "orgid",
*         "created": 1517891917
*     	},
*     	... ...
*     ],
*   },
*   "code": 200,
*   "message": "Created TrackHabit successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track habit was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.GetHabitHistory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) GetHabitHistory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.GetHabitHistory API request")
	req_track := new(track_proto.GetHabitHistoryRequest)

	req_track.HabitId = req.PathParameter("habit_id")
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_track.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_track.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_track.From = req.Attribute(PaginateFromParameter).(int64)
	req_track.To = req.Attribute(PaginateToParameter).(int64)
	req_track.SortParameter = req.Attribute(SortParameter).(string)
	req_track.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.GetHabitHistory(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.GetHabitHistory", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get habit history successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/track/content/{content_id}?session={session_id} Submit an event for a Content
* @apiVersion 0.1.0
* @apiName CreateTrackContent
* @apiGroup Track
*
* @apiDescription Submit an event for a Content
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/content/111?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="
*
* @apiParamExample {json} Request-Example:
* {
*   "user": { User }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "track_content": {
*       "user": { User },
*       "content": { Content },
*       "orgid": "orgid",
*       "created": 1517891917
*     },
*     "count": 3
*   },
*   "code": 200,
*   "message": "Created TrackContent successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track content was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.CreateTrackContent",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) CreateTrackContent(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.CreateTrackContent API request")
	req_track := new(track_proto.CreateTrackContentRequest)
	// err := req.ReadEntity(req_track)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackContent", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_track); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackContent", "BindError")
		return
	}

	req_track.ContentId = req.PathParameter("content_id")
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.CreateTrackContent(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackContent", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created TrackContent successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/track/content/{content_id}/count?session={session_id}&from={from}&to={to} Get content count
* @apiVersion 0.1.0
* @apiName GetContentCount
* @apiGroup Track
*
* @apiDescription Get content count
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/content/111/count?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="&from=13245320&to=135435340
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "count": 3
*   },
*   "code": 200,
*   "message": "Created TrackContent successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track content was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.GetContentCount",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) GetContentCount(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.GetContentCount API request")
	req_track := new(track_proto.GetContentCountRequest)

	req_track.ContentId = req.PathParameter("content_id")
	req_track.UserId = req.Attribute(UserIdAttrName).(string)
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_track.From = req.Attribute(PaginateFromParameter).(int64)
	req_track.To = req.Attribute(PaginateToParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.GetContentCount(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackContent", "CreateError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Get content count successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/track/content/{content_id}/history?session={session_id}&from={from}&to={to}&offset={offset}&limit={limit} Get content history
* @apiVersion 0.1.0
* @apiName GetContentHistory
* @apiGroup Track
*
* @apiDescription Get content history
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/content/111/count?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="&from=13245320&to=135435340&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "track_contents": [
*       {
*         "user": { User },
*         "content": { Content },
*         "orgid": "orgid",
*         "created": 1517891917
*     	},
*     	... ...
*     ],
*   },
*   "code": 200,
*   "message": "Created TrackContent successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track content was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.GetContentHistory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) GetContentHistory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.GetContentHistory API request")
	req_track := new(track_proto.GetContentHistoryRequest)

	req_track.ContentId = req.PathParameter("content_id")
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_track.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_track.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_track.From = req.Attribute(PaginateFromParameter).(int64)
	req_track.To = req.Attribute(PaginateToParameter).(int64)
	req_track.SortParameter = req.Attribute(SortParameter).(string)
	req_track.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.GetContentHistory(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.GetContentHistory", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get content history successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/track/marker/{marker_id}?session={session_id} Track any particular marker
* @apiVersion 0.1.0
* @apiName CreateTrackMarker
* @apiGroup Track
*
* @apiDescription Track a particular marker for a given user_id
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/marker/437c32b2-9dd7-410a-8a97-162e580d8a90?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="
*
* @apiParamExample {json} Request-Example: Manual Tracker Method
* {
*   "user_id": "userid",
*   "org_id": "orgid",
*   "marker_id": "markerid",
*   "tracker_method": {
*                        "id": "0df9459f-44aa-11e8-8879-66a036430288",
*                        "name": "Manual",
*                        "name_slug": "manual",
*                        "icon_slug": "manual-icon"
*                    },
*   "value": 3,
*   "unit": "gms"
* }
*
* @apiParamExample {json} Request-Example: Count Tracker Method
* {
*   "user_id": "userid",
*   "org_id": "orgid",
*   "marker_id": "markerid",
*   "tracker_method": {
*                        "id": "0df9459f-44aa-11e8-8879-66a036430288",
*                        "name": "Count",
*                        "name_slug": "count",
*                        "icon_slug": "count-icon"
*                    },
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "track_marker": {
*       "user": { User },
*       "org_id": "orgid",
*       "marker": { Marker },
*       "created": 1517891917,
*       "value": 3,
*       "unit": "unit sample"
*     }
*   },
*   "code": 200,
*   "message": "Created TrackMark successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track content was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.CreateTrackMarker",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) CreateTrackMarker(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.CreateTrackMarker API request")
	req_track := new(track_proto.CreateTrackMarkerRequest)
	// err := req.ReadEntity(req_track)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackMarker", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_track); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackMarker", "BindError")
		return
	}

	req_track.MarkerId = req.PathParameter("marker_id")
	req_track.UserId = req.Attribute(UserIdAttrName).(string)
	req_track.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.CreateTrackMarker(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.CreateTrackMarker", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created TrackMark successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/track/marker/{marker_id}?session={session_id} Get last tracked value for a particular marker
* @apiVersion 0.1.0
* @apiName GetLastMarker
* @apiGroup Track
*
* @apiDescription Get the most recent value for this marker that was tracked by this user_id
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/marker/437c32b2-9dd7-410a-8a97-162e580d8a90?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "value": 3
*   },
*   "code": 200,
*   "message": "Get value successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track content was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.GetLastMarker",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) GetLastMarker(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.GetLastMarker API request")
	req_track := new(track_proto.GetLastMarkerRequest)

	req_track.MarkerId = req.PathParameter("marker_id")
	req_track.UserId = req.Attribute(UserIdAttrName).(string)
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.GetLastMarker(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.GetLastMarker", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get value successfully"
	data := utils.MarshalAny(rsp, resp)

	// log.Infoln(resp, data)
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/track/marker/{marker_id}/history?session={session_id}&from={from}&to={to}&offset={offset}&limit={limit} Get marker tracking history
* @apiVersion 0.1.0
* @apiName GetMarkerHistory
* @apiGroup Track
*
* @apiDescription Get the history of records for a particular marker
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/marker/437c32b2-9dd7-410a-8a97-162e580d8a90/history?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="&from=13245320&to=135435340&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "track_contents": [
*       {
*         "user": { User },
*         "org_id": "orgid",
*         "marker": { Marker },
*         "created": 1517891917,
*         "value": 3,
*         "unit": "unit sample"
*     	},
*     	... ...
*     ],
*   },
*   "code": 200,
*   "message": "Get marker history successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track content was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.GetMarkerHistory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) GetMarkerHistory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.GetMarkerHistory API request")
	req_track := new(track_proto.GetMarkerHistoryRequest)

	req_track.MarkerId = req.PathParameter("marker_id")
	// req_track.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_track.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_track.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_track.From = req.Attribute(PaginateFromParameter).(int64)
	req_track.To = req.Attribute(PaginateToParameter).(int64)
	req_track.SortParameter = req.Attribute(SortParameter).(string)
	req_track.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.GetMarkerHistory(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.GetMarkerHistory", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get marker history successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/track/markers/history/all?session={session_id}&from={from}&to={to}&offset={offset}&limit={limit} Get all tracking history
* @apiVersion 0.1.0
* @apiName GetAllMarkerHistory
* @apiGroup Track
*
* @apiDescription Get the history of tracked records chronologically
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/track/markers/history/all?session="zqoFHcqIPwYbP2QvdfR_W0381FAI7k2HjOh7nGzNskE="&from=13245320&to=135435340&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "track_contents": [
*       {
*         "user": { User },
*         "org_id": "orgid",
*         "marker": { Marker },
*         "created": 1517891917,
*         "value": 3,
*         "unit": "unit sample"
*     	},
*     	... ...
*     ],
*   },
*   "code": 200,
*   "message": "Get markers all history successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The track content was not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.track.GetAllMarkerHistory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TrackService) GetAllMarkerHistory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Track.GetAllMarkerHistory API request")
	req_track := new(track_proto.GetAllMarkerHistoryRequest)

	req_track.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_track.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_track.From = req.Attribute(PaginateFromParameter).(int64)
	req_track.To = req.Attribute(PaginateToParameter).(int64)
	req_track.SortParameter = req.Attribute(SortParameter).(string)
	req_track.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.TrackClient.GetAllMarkerHistory(ctx, req_track)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.track.GetAllMarkerHistory", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Get markers all history successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}
