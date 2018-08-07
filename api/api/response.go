package api

import (
	"context"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	resp_proto "server/response-srv/proto/response"
	survey_proto "server/survey-srv/proto/survey"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type ResponseService struct {
	ResponseClient resp_proto.ResponseServiceClient
	SurveyClient   survey_proto.SurveyServiceClient
	Auth           Filters
	Audit          AuditFilter
	ServerMetrics  metrics.Metrics
}

func (p ResponseService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/responses")

	audit := &audit_proto.Audit{
		ActionService:  common.ResponseSrv,
		ActionResource: common.BASE + common.RESPONSE_TYPE,
	}

	ws.Route(ws.GET("/{hash}/check").To(p.Check).
		Doc("Check response auth"))

	ws.Route(ws.GET("/survey/{survey_id}/questions/all").To(p.AllQuestion).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all survey questions"))

	ws.Route(ws.GET("/survey/{survey_id}/questions/{question_id}").To(p.ReadQuestion).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get survey question by question id"))

	ws.Route(ws.GET("/open/survey/{survey_id}/questions/all").To(p.OpenAllQuestion).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all survey questions (authentication NOT required)"))

	ws.Route(ws.GET("/open/survey/{survey_id}/questions/{question_id}").To(p.OpenReadQuestion).
		Filter(p.Auth.Paginate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get survey question by question id"))
	// maybe added groupby
	// maybe added certain timeperiod
	ws.Route(ws.GET("/{survey_id}/all").To(p.All).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Returns a list of submitted survey responses"))

	//TODO:need to create a response endpoint which doesn't require authentication
	ws.Route(ws.POST("/{survey_id}/response").To(p.Create).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Returns a created response"))

	ws.Route(ws.POST("/{survey_id}/response/state").To(p.UpdateState).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Update response state"))

	ws.Route(ws.GET("/{survey_id}/all/state/{response_state}").To(p.AllState).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all responses by response state"))

	ws.Route(ws.GET("/{survey_id}/response/stats").To(p.ReadStats).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Returns a list of submitted survey responses"))

	ws.Route(ws.GET("/{survey_id}/response/by/{user_id}").To(p.ByUser).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Returns a list of submitted survey responses"))

	ws.Route(ws.GET("/{survey_id}/response/anon").To(p.ByAnyUser).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Returns a list of submitted survey responses"))

	restful.Add(ws)
}

/**
* @api {get} /server/responses/{hash}/check?session={session_id} Check survey auth
* @apiVersion 0.1.0
* @apiName Check
* @apiGroup Response
*
* @apiDescription Check survey auth - checks whether for survey with short_hash (from map<short_hash,survey_id> requires authentication or not. Check should be made against authenticationRequired field of the survey
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/responses/2GGA1vK/check?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "survey": {
*       "id": "wZIElizIhP9LVobz",
*       "authenticationRequired": true
*     }
*   },
*   "code": 200,
*   "message": "Checked successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.response.Check",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ResponseService) Check(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Response.Check API request")
	req_resp := new(resp_proto.CheckRequest)
	req_resp.ShortHash = req.PathParameter("hash")
	//req_resp.OrgId = req.Attribute(OrgIdAttrName).(string)
	//req_resp.TeamId = req.Attribute(TeamIdAttrName).(string)
	// req_resp.Limit = req.Attribute(PaginateLimitParameter).(int64)
	// req_resp.Offset = req.Attribute(PaginateOffsetParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	res, err := p.ResponseClient.Check(ctx, req_resp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.Check", "QueryError")
		return
	}

	res.Code = http.StatusOK
	res.Message = "Checked successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, res)
}

/**
* @api {get} /server/responses/survey/{survey_id}/questions/all?session={session_id}&offset={offset}&limit={limit} Get all survey questions
* @apiVersion 0.1.0
* @apiName AllQuestion
* @apiGroup Response
*
* @apiDescription Get all survey questions (authentication required) - returns the array of survey.questions[n]. Use survey-srv to return an array of questions based on the survey_id.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/responses/survey/111/questions/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*	  "welcome": {DefaultQuestion},
*	  "thankyou": {DefaultQuestion},
*     "questions": [
*       {
*         "id": "q111",
*         "type": QuestionType,
*         "order": 100,
*         "title": "question",
*         "description": "description",
*         "design":  {
*           "progress_bar_style": 0,
*           "default_bg_color": "#45AC3C",
*           "default_logo_url": "http://example.com/logo.png",
*         },
*         "fields": google.protobuf.Any,
*         "settings": google.protobuf.Any
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all questions successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.response.AllQuestion",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ResponseService) AllQuestion(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Response.AllQuestion API request")
	req_resp := new(resp_proto.AllQuestionRequest)
	req_resp.SurveyId = req.PathParameter("survey_id")
	req_resp.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_resp.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_resp.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_resp.SortParameter = req.Attribute(SortParameter).(string)
	req_resp.SortDirection = req.Attribute(SortDirection).(string)

	req_survey := &survey_proto.QuestionsRequest{
		SurveyId: req_resp.SurveyId,
		OrgId:    req_resp.OrgId,
		Offset:   req_resp.Offset,
		Limit:    req_resp.Limit,
	}
	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	res, err := p.SurveyClient.Questions(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.AllQuestion", "QueryError")
		return
	}

	res.Code = http.StatusOK
	res.Message = "Read all questions successfully"
	data := utils.MarshalAny(rsp, res)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/responses/survey/{survey_id}/questions/{question_id}?session={session_id} Get survey question
* @apiVersion 0.1.0
* @apiName ReadQuestion
* @apiGroup Response
*
* @apiDescription Get survey question by question_id (authentication required) - returns questions with id = question_id. Use survey-srv to return an a question by id based on the question_id.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/responses/survey/111/questions/q111?session={session_id-}
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "question": {
*       "id": "q111",
*       "type": QuestionType,
*       "order": 100,
*       "title": "question",
*       "description": "description",
*       "design":  {
*         "progress_bar_style": 0,
*         "default_bg_color": "#45AC3C",
*         "default_logo_url": "http://example.com/logo.png",
*       },
*       "fields": google.protobuf.Any,
*       "settings": google.protobuf.Any
*     },
*   },
*   "code": 200,
*   "message": "Read question successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.response.ReadQuestion",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ResponseService) ReadQuestion(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Response.ReadQuestion API request")
	req_resp := new(resp_proto.ReadQuestionRequest)
	req_resp.SurveyId = req.PathParameter("survey_id")
	req_resp.QuestionId = req.PathParameter("question_id")
	req_resp.OrgId = req.Attribute(OrgIdAttrName).(string)

	req_survey := &survey_proto.QuestionRefRequest{
		SurveyId:    req_resp.SurveyId,
		QuestionRef: req_resp.QuestionId,
		OrgId:       req_resp.OrgId,
	}
	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	res, err := p.SurveyClient.QuestionRef(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.ReadQuestion", "QueryError")
		return
	}

	res.Code = http.StatusOK
	res.Message = "Read question successfully"
	data := utils.MarshalAny(rsp, res)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/responses/open/survey/{survey_id}/questions/all?session={session_id}&offset={offset}&limit={limit} Get all survey questions
* @apiVersion 0.1.0
* @apiName OpenAllQuestion
* @apiGroup Response
*
* @apiDescription Get all survey questions (authentication NOT required) - returns the array of survey.questions[n]. Use survey-srv to return an array of questions based on the survey_id.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/responses/open/survey/111/questions/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "questions": [
*       {
*         "id": "q111",
*         "type": QuestionType,
*         "order": 100,
*         "title": "question",
*         "description": "description",
*         "design":  {
*           "progress_bar_style": 0,
*           "default_bg_color": "#45AC3C",
*           "default_logo_url": "http://example.com/logo.png",
*         },
*         "fields": google.protobuf.Any,
*         "settings": google.protobuf.Any
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all questions successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.response.OpenAllQuestion",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ResponseService) OpenAllQuestion(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Response.OpenAllQuestion API request")
	req_resp := new(resp_proto.AllQuestionRequest)
	req_resp.SurveyId = req.PathParameter("survey_id")
	req_resp.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_resp.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_resp.SortParameter = req.Attribute(SortParameter).(string)
	req_resp.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)

	// fetch survey
	req_read := new(survey_proto.ReadRequest)
	req_read.Id = req_resp.SurveyId
	res_read, err := p.SurveyClient.Read(ctx, req_read)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.OpenAllQuestion", "SurveyError")
		return
	}
	if res_read.Data.Survey.Setting != nil && res_read.Data.Survey.Setting.AuthenticationRequired {
		res_read.Data = nil
		res_read.Code = http.StatusOK
		res_read.Message = "Fetched survey requires authentication"
		rsp.AddHeader("Content-Type", "application/json")
		rsp.WriteHeaderAndEntity(http.StatusOK, res_read)
		return
	}

	req_survey := &survey_proto.QuestionsRequest{
		SurveyId: req_resp.SurveyId,
		Offset:   req_resp.Offset,
		Limit:    req_resp.Limit,
	}
	res, err := p.SurveyClient.Questions(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.OpenAllQuestion", "QueryError")
		return
	}

	res.Code = http.StatusOK
	res.Message = "Read all questions successfully"
	data := utils.MarshalAny(rsp, res)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/responses/open/survey/{survey_id}/questions/{question_id}?session={session_id} Get survey question
* @apiVersion 0.1.0
* @apiName ReadQuestion
* @apiGroup Response
*
* @apiDescription Get survey question by question_id (authentication NOT required) - returns questions with id = question_id. Use survey-srv to return an a question by id based on the question_id.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/open/server/responses/survey/111/questions/q111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "question": {
*       "id": "q111",
*       "type": QuestionType,
*       "order": 100,
*       "title": "question",
*       "description": "description",
*       "design":  {
*         "progress_bar_style": 0,
*         "default_bg_color": "#45AC3C",
*         "default_logo_url": "http://example.com/logo.png",
*       },
*       "fields": google.protobuf.Any,
*       "settings": google.protobuf.Any
*     },
*   },
*   "code": 200,
*   "message": "Read question successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.response.ReadQuestion",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ResponseService) OpenReadQuestion(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Response.OpenReadQuestion API request")
	req_resp := new(resp_proto.ReadQuestionRequest)
	req_resp.SurveyId = req.PathParameter("survey_id")
	req_resp.QuestionId = req.PathParameter("question_id")
	// req_resp.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)
	// req_resp.Limit = req.Attribute(PaginateLimitParameter).(int64)
	// req_resp.Offset = req.Attribute(PaginateOffsetParameter).(int64)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	// fetch survey
	req_read := new(survey_proto.ReadRequest)
	req_read.Id = req_resp.SurveyId
	res_read, err := p.SurveyClient.Read(ctx, req_read)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.OpenReadQuestion", "SurveyError")
		return
	}
	if res_read.Data.Survey.Setting != nil && res_read.Data.Survey.Setting.AuthenticationRequired {
		res_read.Data = nil
		res_read.Code = http.StatusOK
		res_read.Message = "Fetched survey requires authentication"
		rsp.AddHeader("Content-Type", "application/json")
		rsp.WriteHeaderAndEntity(http.StatusOK, res_read)
		return
	}

	req_survey := &survey_proto.QuestionRefRequest{
		SurveyId:    req_resp.SurveyId,
		QuestionRef: req_resp.QuestionId,
	}
	res, err := p.SurveyClient.QuestionRef(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.OpenReadQuestion", "QueryError")
		return
	}

	res.Code = http.StatusOK
	res.Message = "Checked successfully"
	data := utils.MarshalAny(rsp, res)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/responses/{survey_id}/all?session={session_id}&offset={offset}&limit={limit} List all responses
* @apiVersion 0.1.0
* @apiName All
* @apiGroup Response
*
* @apiDescription List all responses - returns a list of submitted survey responses
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/open/server/responses/111/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "responses": [
*       {
*         "id": "111",
*         "org_id": "orgid",
*         "survey_id": "111",
*         "response_session": "sesssion",
*         "metadata":  { Metadata },
*         "responder": { User },
*         "answers": [{Answer}, {Answer}, ...],
*         "status": {
*           "state": ResponseState,
*           "timestamp": 1517891917
*         },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all responses successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.response.All",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

/**
* @api {get} /server/responses/{survey_id}/all?session={session_id}&offset={offset}&limit={limit}&from={timestamp}&to={timestamp} List all timeperiod responses
* @apiVersion 0.1.0
* @apiName AllTimeperiod
* @apiGroup Response
*
* @apiDescription List all responses between a certain timeperiod for a survey
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/open/server/responses/111/all??session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10&from=1517791917&to=1517991917
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "responses": [
*       {
*         "id": "111",
*         "org_id": "orgid",
*         "survey_id": "111",
*         "response_session": "sesssion",
*         "metadata":  { Metadata },
*         "responder": { User },
*         "answers": [{Answer}, {Answer}, ...],
*         "status": {
*           "state": ResponseState,
*           "timestamp": 1517891917
*         },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all responses successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.response.All",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

/**
* @api {get} /server/responses/{survey_id}/all?groupby=question?session={session_id}&offset={offset}&limit={limit} List all groupby responses
* @apiVersion 0.1.0
* @apiName AllGroupBy
* @apiGroup Response
*
* @apiDescription List all responses grouped by question - returns a list of GroupByQuestionResponses. All responses are being returned grouped by question for survey_id with survey.question[n]
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/open/server/responses/111/all?groupby=question&session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "org_id": "orgid",
*     "survey_id": "111",
*     "responses": [
*       {
*         "response_count": 8,
*         "skipped_count": 2,
*         "question_ref": "q111",
*         "type": QuestionType,
*         "answers": [{Answer}, {Answer}, ...]
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all responses successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.response.All",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ResponseService) All(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Response.All API request")
	req_resp := new(resp_proto.AllRequest)
	req_resp.SurveyId = req.PathParameter("survey_id")
	req_resp.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_resp.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_resp.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_resp.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_resp.SortParameter = req.Attribute(SortParameter).(string)
	req_resp.SortDirection = req.Attribute(SortDirection).(string)

	// spec groupby
	question := req.QueryParameter("groupby")
	// spec time period
	from := req.QueryParameter("from")
	to := req.QueryParameter("to")

	var ress interface{}
	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)

	if len(question) > 0 {
		req_group := &resp_proto.AllAggQuestionRequest{
			SurveyId: req_resp.SurveyId,
			OrgId:    req_resp.OrgId,
		}
		res, err := p.ResponseClient.AllAggQuestion(ctx, req_group)
		if err != nil {
			utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.AllGroupBy", "QueryError")
			return
		}
		res.Code = http.StatusOK
		res.Message = "Grouped by question successfully"
		ress = utils.MarshalAny(rsp, res)
	} else if len(from) > 0 && len(to) > 0 {
		f, _ := strconv.ParseInt(from, 10, 64)
		t, _ := strconv.ParseInt(to, 10, 64)
		req_time := &resp_proto.TimeFilterRequest{
			SurveyId: req_resp.SurveyId,
			From:     f,
			To:       t,
			OrgId:    req_resp.OrgId,
			TeamId:   req_resp.TeamId,
			Limit:    req_resp.Limit,
			Offset:   req_resp.Offset,
		}
		res, err := p.ResponseClient.TimeFilter(ctx, req_time)
		if err != nil {
			utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.AllTimeFilter", "QueryError")
			return
		}
		res.Code = http.StatusOK
		res.Message = "Time filtered all responses successfully"
		ress = utils.MarshalAny(rsp, res)
	} else {
		res, err := p.ResponseClient.All(ctx, req_resp)
		if err != nil {
			utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.All", "QueryError")
			return
		}
		res.Code = http.StatusOK
		res.Message = "Read all responses successfully"
		ress = utils.MarshalAny(rsp, res)
	}

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, ress)
}

/**
* @api {post} /server/responses/{survey_id}/response?session={session_id} Submit a response
* @apiVersion 0.1.0
* @apiName Submit
* @apiGroup Response
*
* @apiDescription Submit a response - This response needs to have a relation (edge) to the survey with survey_id and user with user_id (in case authentication is required).
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/responses/111/response?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "response": {
*     "id": "111",
*     "org_id": "orgid",
*     "survey_id": "111",
*     "response_session": "sesssion",
*     "metadata":  { Metadata },
*     "responder": { User },
*     "answers": [{Answer}, {Answer}, ...],
*     "status": {
*       "state": ResponseState,
*       "timestamp": 1517891917
*     }
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "response": {
*       "id": "111",
*       "org_id": "orgid",
*       "survey_id": "111",
*       "response_session": "sesssion",
*       "metadata":  { Metadata },
*       "responder": { User },
*       "answers": [{Answer}, {Answer}, ...],
*       "status": {
*         "state": ResponseState,
*         "timestamp": 1517891917
*       }
*     }
*   },
*   "code": 200,
*   "message": "Created response successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.survey.All",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ResponseService) Create(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Response.Create API request")
	req_resp := new(resp_proto.CreateRequest)
	// err := req.ReadEntity(req_resp)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.CreateSurvey", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_resp); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.Create", "BindError")
		return
	}

	req_resp.SurveyId = req.PathParameter("survey_id")
	//missing user_id here
	req_resp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	res, err := p.ResponseClient.Create(ctx, req_resp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.Create", "QueryError")
		return
	}

	res.Code = http.StatusOK
	res.Message = "Created response successfully"
	data := utils.MarshalAny(rsp, res)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/responses/{survey_id}/response/state?session={session_id} Update response state
* @apiVersion 0.1.0
* @apiName Update response state
* @apiGroup Response
*
* @apiDescription Update response state - returns response_id
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/responses/111/response/state?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "response": {
*     "id": "111",
*     "state": [ResponseState]
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "response_id": "111",
*   "code": 200,
*   "message": "Updated response successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.survey.UpdateState",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ResponseService) UpdateState(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Response.UpdateState request")
	req_resp := new(resp_proto.UpdateStateRequest)
	err := req.ReadEntity(req_resp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.UpdateState", "BindError")
		return
	}

	req_resp.SurveyId = req.PathParameter("survey_id")
	//missing user_id
	req_resp.OrgId = req.Attribute(OrgIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	res, err := p.ResponseClient.UpdateState(ctx, req_resp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.UpdateState", "QueryError")
		return
	}

	res.Code = http.StatusOK
	res.Message = "Updated response successfully"
	data := utils.MarshalAny(rsp, res)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/responses/{survey_id}/all/state/{response_state}?session={session_id}&offset={offset}&limit={limit} List all state responses
* @apiVersion 0.1.0
* @apiName AllState
* @apiGroup Response
*
* @apiDescription List all responses by response state - returns a list of responses by response state.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/responses/111/all/state/1?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "responses": [
*       {
*         "id": "111",
*         "org_id": "orgid",
*         "survey_id": "111",
*         "response_session": "sesssion",
*         "metadata":  { Metadata },
*         "responder": { User },
*         "answers": [{Answer}, {Answer}, ...],
*         "status": {
*           "state": ResponseState,
*           "timestamp": 1517891917
*         },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read state responses successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.response.ReadQuestion",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ResponseService) AllState(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Response.AllState API request")
	req_resp := new(resp_proto.AllStateRequest)
	req_resp.SurveyId = req.PathParameter("survey_id")
	s, _ := strconv.Atoi(req.PathParameter("response_state"))
	req_resp.State = resp_proto.ResponseState(s)

	req_resp.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_resp.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_resp.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_resp.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_resp.SortParameter = req.Attribute(SortParameter).(string)
	req_resp.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	res, err := p.ResponseClient.AllState(ctx, req_resp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.AllState", "QueryError")
		return
	}

	res.Code = http.StatusOK
	res.Message = "Read state responses successfully"
	data := utils.MarshalAny(rsp, res)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/responses/{survey_id}/all/states?session={session_id} Get response stats
* @apiVersion 0.1.0
* @apiName ReadStats
* @apiGroup Response
*
* @apiDescription Get response stats for a survey = survey_id
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/responses/111/all/states?session={session_id}
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "stats": {
*       "responses": 10,
*       "drops": 8
*     }
*   },
*   "code": 200,
*   "message": "Read stats successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.response.ReadStats",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ResponseService) ReadStats(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Response.ReadStats API request")
	req_resp := new(resp_proto.ReadStatsRequest)
	req_resp.SurveyId = req.PathParameter("survey_id")
	req_resp.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_resp.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	res, err := p.ResponseClient.ReadStats(ctx, req_resp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.ReadStats", "QueryError")
		return
	}

	res.Code = http.StatusOK
	res.Message = "Read stats successfully"
	data := utils.MarshalAny(rsp, res)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/responses/{survey_id}/response/by/{user_id}?session={session_id}&offset={offset}&limit={limit} Get all responses by a user
* @apiVersion 0.1.0
* @apiName ByUser
* @apiGroup Response
*
* @apiDescription Get all responses by a particular user - Get all responses for a survey (survey_id) where response.userid = {userid} where survey.authenticationrequired = true
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/responses/111/response/by/userid?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "responses": [
*       {
*         "id": "111",
*         "org_id": "orgid",
*         "survey_id": "111",
*         "response_session": "sesssion",
*         "metadata":  { Metadata },
*         "responder": { User },
*         "answers": [{Answer}, {Answer}, ...],
*         "status": {
*           "state": ResponseState,
*           "timestamp": 1517891917
*         },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all responses by user successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.response.ByUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ResponseService) ByUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Response.ByUser API request")
	req_resp := new(resp_proto.ByUserRequest)
	req_resp.SurveyId = req.PathParameter("survey_id")
	req_resp.UserId = req.PathParameter("user_id")

	req_resp.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_resp.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_resp.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_resp.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_resp.SortParameter = req.Attribute(SortParameter).(string)
	req_resp.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	res, err := p.ResponseClient.ByUser(ctx, req_resp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.ByUser", "QueryError")
		return
	}

	res.Code = http.StatusOK
	res.Message = "Read all responses by user successfully"
	data := utils.MarshalAny(rsp, res)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/responses/{survey_id}/response/anon?session={session_id}&offset={offset}&limit={limit} Get all responses by anonymous user
* @apiVersion 0.1.0
* @apiName ByAnyUser
* @apiGroup Response
*
* @apiDescription Get all responses by anonymous user - Get all responses for a survey (survey_id) where response.userid = null where survey.authenticationrequired = false
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/responses/111/response/anon?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "responses": [
*       {
*         "id": "111",
*         "org_id": "orgid",
*         "survey_id": "111",
*         "response_session": "sesssion",
*         "metadata":  { Metadata },
*         "responder": { User },
*         "answers": [{Answer}, {Answer}, ...],
*         "status": {
*           "state": ResponseState,
*           "timestamp": 1517891917
*         },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all responses by anonymous user successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The surveys were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.response.ByAnyUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *ResponseService) ByAnyUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Response.ByAnyUser API request")
	req_resp := new(resp_proto.ByAnyUserRequest)
	req_resp.SurveyId = req.PathParameter("survey_id")

	req_resp.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_resp.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_resp.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_resp.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_resp.SortParameter = req.Attribute(SortParameter).(string)
	req_resp.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	res, err := p.ResponseClient.ByAnyUser(ctx, req_resp)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.response.ByAnyUser", "QueryError")
		return
	}

	res.Code = http.StatusOK
	res.Message = "Read all responses by anonymous user successfully"
	data := utils.MarshalAny(rsp, res)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}
