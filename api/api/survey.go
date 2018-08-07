package api

import (
	"context"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	organisation_proto "server/organisation-srv/proto/organisation"
	survey_proto "server/survey-srv/proto/survey"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type SurveyService struct {
	SurveyClient       survey_proto.SurveyServiceClient
	Auth               Filters
	Audit              AuditFilter
	OrganisationClient organisation_proto.OrganisationServiceClient
	ServerMetrics      metrics.Metrics
}

func (p SurveyService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/surveys")

	audit := &audit_proto.Audit{
		ActionService:  common.SurveySrv,
		ActionResource: common.BASE + common.SURVEY_TYPE,
	}

	ws.Route(ws.GET("/all").To(p.AllSurveys).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all surveys"))

	ws.Route(ws.GET("/new").To(p.NewSurvey).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create unique id for new survey"))

	ws.Route(ws.POST("/survey/create").To(p.CreateSurvey).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("create data"))

	ws.Route(ws.GET("/survey/{survey_id}").To(p.ReadSurvey).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("View Survey detail"))

	ws.Route(ws.DELETE("/survey/{survey_id}").To(p.DeleteSurvey).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("View Survey detail"))

	ws.Route(ws.POST("/survey/copy").To(p.CopySurvey).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("View Survey detail"))

	ws.Route(ws.GET("/survey/{survey_id}/questions").To(p.Questions).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get survey questions"))

	ws.Route(ws.GET("/survey/{survey_id}/questions/{question_id}").To(p.QuestionRef).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get survey question by question_id "))

	ws.Route(ws.POST("/survey/{survey_id}/questions").To(p.CreateQuestion).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Add or Update a question in the survey - Accepts a new list of questions for a survey"))

	ws.Route(ws.GET("/creator/{user_id}").To(p.ByCreator).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all surveys created by a particular team member"))

	ws.Route(ws.GET("/survey/{survey_id}/link").To(p.Link).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all surveys created by a particular team member"))

	ws.Route(ws.GET("/templates").To(p.Templates).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all surveys created by a particular team member"))

	ws.Route(ws.POST("/filter").To(p.Filter).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all surveys created by a particular team member"))

	ws.Route(ws.POST("/search").To(p.SearchSurveys).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all surveys"))

	ws.Route(ws.POST("/survey/search/autocomplete").To(p.AutocompleteSurveySearch).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete survey text"))

	ws.Route(ws.GET("/tags/top/{n}").To(p.GetTopSurveyTags).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Return top N tags for Survey"))

	ws.Route(ws.POST("/tags/autocomplete").To(p.AutocompleteSurveyTags).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Autocomplete for tags for Survey"))

	restful.Add(ws)
}

// func fetchUserFromSurvey(ctx context.Context, survey *survey_proto.Survey, userClient user_proto.UserServiceClient) {
// 	if survey == nil {
// 		return
// 	}
// 	resp_user, err := userClient.Read(ctx, &user_proto.ReadRequest{survey.CreatorId})
// 	// fmt.Println("resp_user", resp_user)
// 	if err == nil {
// 		survey.Creator = resp_user.Data.User
// 	}
// }

/**
 * @api {get} /server/surveys/all?session={session_id}&offset={offset}&limit={limit} List all surveys
 * @apiVersion 0.1.0
 * @apiName ListAll
 * @apiGroup Survey
 *
 * @apiDescription This API endpoint should return a restricted set of information for quick access. It should return a lite version of the survey object
 *
 * @apiExample Example usage:
 * curl -i http://BASE_SERVER_URL/server/surveys/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
 *
 * @apiSuccessExample Success-Response:
 * HTTP/1.1 200 OK
 * {
 *   "data": {
 *     "surveys": [
 *       {
 *         "id": "111",
 *         "org_id": "orgid",
 *         "title": "survey",
 *         "description": "This is sample survey",
 *         "tags": "tags",
 *         "created": 1517891917,
 *         "updated": 1517891917,
 *         "creator_id": "userid",
 *         "creator": {User},
 *         "shares": [{User},{User}...],
 *         "renders": [],
 *         "setting": {
 *           "visibility": 0,
 *           "notifications": [],
 *           "social": [],
 *           "linkSharingEnabled": true,
 *           "embeddingEnabled": false,
 *           "authentificationRequried": true,
 *           "showCaptcha": true
 *         },
 *         "welcome": {DefaultQuestion},
 *         "thankyou": {DefaultQuestion},
 *         "questions": [{Question},{Question}...],
 *         "design": {
 *           "progress_bar_style": 0,
 *           "default_bg_color": "#45AC3C",
 *           "default_logo_url": "http://example.com/logo.png",
 *         },
 *         "status": 0,
 *         "isTemplate": true,
 *         "templateId": "templateId",
 *       },
 *       ... ...
 *     ]
 *   },
 *   "code": 200,
 *   "message": "Read all surveys successfully"
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
func (p *SurveyService) AllSurveys(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.All API request")
	req_survey := new(survey_proto.AllRequest)
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_survey.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_survey.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_survey.SortParameter = req.Attribute(SortParameter).(string)
	req_survey.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	all_resp, err := p.SurveyClient.All(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.AllSurveys", "QueryError")
		return
	}
	// fetch all users
	// for _, s := range all_resp.Data.Surveys {
	// 	fetchUserFromSurvey(ctx, s, p.Auth.UserClient)
	// }

	all_resp.Code = http.StatusOK
	all_resp.Message = "Read all surveys successfully"
	data := utils.MarshalAny(rsp, all_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
 * @api {get} /server/surveys/new?session={session_id} Create a unique survey id
 * @apiVersion 0.1.0
 * @apiName Create SurveyId
 * @apiGroup Survey
 *
 * @apiDescription When creating this new survey id for the new survey, we also need to create a short hash of the survey id by using a hashing algorithm and current timestamp (for creating unique hash) a mapping of this short hash and survey_id needs to be stored in arangodb as a map<k,v> = survey_hash<unique_hash,survey_id>. The unique_hash can't be more than 6 characters
 *
 * @apiExample Example usage:
 * curl -i http://BASE_SERVER_URL/server/surveys/new?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
 *
 * @apiSuccessExample Success-Response:
 * HTTP/1.1 200 OK
 * {
 *   "data": {
 *	 	"unique_hash": "4770GqJ",
 *		"survey_id": "TTLg2rYHkFe7GNXh"
 *   },
 *   "code": 200,
 *   "message": "Created survey successfully"
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
 *           "domain": "go.micro.srv.survey.NewSurvey",
 *           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
 *         }
 *       ]
 *     }
 */
func (p *SurveyService) NewSurvey(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.New API request")
	req_survey := new(survey_proto.NewRequest)
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	new_resp, err := p.SurveyClient.New(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.NewSurvey", "QueryError")
		return
	}

	new_resp.Code = http.StatusOK
	new_resp.Message = "Created survey successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, new_resp)
}

/**
* @api {post} /server/surveys/survey/create?session={session_id} Create or Update
* @apiVersion 0.1.0
* @apiName Create or Update
* @apiGroup Survey
*
* @apiDescription This API saves or updates the survey submitted in the POST body when saving the survey in the survey collections, please add a unique question_ref to each question in survey.questions[n] - this question_ref should be a unique hash id for this survey.questions[n]
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/surveys/survey/create?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "survey": {
*     "id": "111",
*     "org_id": "orgid",
*     "title": "new survey with questions",
*     "description": "This is sample survey",
*     "tags": "tags",
*     "creator_id": "userid",
*     "creator": {
*       "id":"4158873242954991778",
*       "orgid":"vzvG73i41edPNxJ7"
*     },
*     "shares": [{User},{User}...],
*     "renders": [],
*     "setting": {
*       "shareableLink":"http://hel.ly/2eA8KKW",
*       "linkSharingEnabled":true
*     },
*     "welcome":{
*       "type":0,
*       "order":0,
*       "settings":{
*         "showButton":true,
*         "buttonText":"",
*         "social_sharing_enabled":false,
*         "submit_mode":1,
*         "showTimeToAnswer":false
*       }
*     },
*     "thankyou": {
*       "type":9,
*       "order":0,
*       "settings":{
*         "showButton":true,
*         "buttonText":"",
*         "social_sharing_enabled":false,
*         "submit_mode":1,
*         "showTimeToAnswer":false
*       }
*     },
*     "questions": [
*       {
*         "id":"ydbUs",
*         "type":1,
*         "order":1,
*         "design":{
*           "bg_color":"fffff",
*           "logo_url":"http://via.placeholder.com/300x300"
*         },
*         "title":"question 1?"
*       },
*       {
*         "id":"Hu7j6",
*         "type":7,
*         "order":2,
*         "design":{
*           "bg_color":"fffff",
*           "logo_url":"http://via.placeholder.com/300x300"
*         },
*         "settings":{
*           "@type":"healum.com/proto/go.micro.srv.survey.BinaryQuestionSettings",
*           "mandatory":true,
*           "image":false,
*           "video":false,
*           "buttonType":0
*         },
*         "title":"question 2?"
*       }
*     ],
*     "design": {
*       "progress_bar_style": 0,
*       "default_bg_color": "#45AC3C",
*       "default_logo_url": "http://example.com/logo.png",
*     },
*     "status": 0,
*     "isTemplate": true,
*     "templateId": "templateId",
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "id": "111",
*     "org_id": "orgid",
*     "title": "new survey with questions",
*     "description": "This is sample survey",
*     "tags": "tags",
*     "created": 1517891917,
*     "updated": 1517891917,
*     "creator_id": "userid",
*     "creator": {
*       "id":"4158873242954991778",
*       "orgid":"vzvG73i41edPNxJ7"
*     },
*     "shares": [{User},{User}...],
*     "renders": [],
*     "setting": {
*       "shareableLink":"http://hel.ly/2eA8KKW",
*       "linkSharingEnabled":true
*     },
*     "welcome":{
*       "type":0,
*       "order":0,
*       "settings":{
*         "showButton":true,
*         "buttonText":"",
*         "social_sharing_enabled":false,
*         "submit_mode":1,
*         "showTimeToAnswer":false
*       }
*     },
*     "thankyou": {
*       "type":9,
*       "order":0,
*       "settings":{
*         "showButton":true,
*         "buttonText":"",
*         "social_sharing_enabled":false,
*         "submit_mode":1,
*         "showTimeToAnswer":false
*       }
*     },
*     "questions": [
*       {
*         "id":"ydbUs",
*         "type":1,
*         "order":1,
*         "design":{
*           "bg_color":"fffff",
*           "logo_url":"http://via.placeholder.com/300x300"
*         },
*         "title":"question 1?"
*       },
*       {
*         "id":"Hu7j6",
*         "type":7,
*         "order":2,
*         "design":{
*           "bg_color":"fffff",
*           "logo_url":"http://via.placeholder.com/300x300"
*         },
*         "settings":{
*           "@type":"healum.com/proto/go.micro.srv.survey.BinaryQuestionSettings",
*           "mandatory":true,
*           "image":false,
*           "video":false,
*           "buttonType":0
*         },
*         "title":"question 2?"
*       }
*     ],
*     "design": {
*       "progress_bar_style": 0,
*       "default_bg_color": "#45AC3C",
*       "default_logo_url": "http://example.com/logo.png",
*     },
*     "status": 0,
*     "isTemplate": true,
*     "templateId": "templateId",
*   },
*   "code": 200,
*   "message": "Created survey successfully"
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
func (p *SurveyService) CreateSurvey(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.Create API request")
	req_survey := new(survey_proto.CreateRequest)
	if err := utils.UnmarshalAny(req, rsp, req_survey); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.CreateSurvey", "BindError")
		return
	}
	req_survey.UserId = req.Attribute(UserIdAttrName).(string)
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	create_resp, err := p.SurveyClient.Create(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.CreateSurvey", "CreateError")
		return
	}

	create_resp.Code = http.StatusOK
	create_resp.Message = "Created survey successfully"
	data := utils.MarshalAny(rsp, create_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/surveys/survey/{survey_id}?session={session_id} View survey detail
* @apiVersion 0.1.0
* @apiName ReadSurvey
* @apiGroup Survey
*
* @apiDescription This API reads survey details with survey id.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/surveys/survey/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "survey": {
*       "id": "111",
*       "org_id": "orgid",
*       "title": "survey",
*       "description": "This is sample survey",
*       "tags": "tags",
*       "created": 1517891917,
*       "updated": 1517891917,
*       "creator_id": "userid",
*       "creator": {User},
*       "shares": [{User},{User}...],
*       "renders": [],
*       "setting": {
*         "visibility": 0,
*         "notifications": [],
*         "social": [],
*         "linkSharingEnabled": true,
*         "embeddingEnabled": false,
*         "authentificationRequried": true,
*         "showCaptcha": true
*       },
*       "welcome": {DefaultQuestion},
*       "thankyou": {DefaultQuestion},
*       "questions": [{Question},{Question}...],
*       "design": {
*         "progress_bar_style": 0,
*         "default_bg_color": "#45AC3C",
*         "default_logo_url": "http://example.com/logo.png",
*       },
*       "status": 0,
*       "isTemplate": true,
*       "templateId": "templateId",
*     }
*   },
*   "code": 200,
*   "message": "Read survey successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The survey was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.survey.Read",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *SurveyService) ReadSurvey(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.Read API request")
	req_survey := new(survey_proto.ReadRequest)
	req_survey.Id = req.PathParameter("survey_id")
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	read_resp, err := p.SurveyClient.Read(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.ReadSurvey", "ReadError")
		return
	}
	// fetchUserFromSurvey(ctx, read_resp.Data.Survey, p.Auth.UserClient)

	read_resp.Code = http.StatusOK
	read_resp.Message = "Read survey successfully"
	data := utils.MarshalAny(rsp, read_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {delete} /server/surveys/survey/{survey_id}?session={session_id} Delete a survey
* @apiVersion 0.1.0
* @apiName DeleteSurvey
* @apiGroup Survey
*
* @apiDescription This API delete existed survey with survey id.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/surveys/survey/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted survey successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The survey was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeletError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.survey.Delete",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *SurveyService) DeleteSurvey(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.Delete API request")
	req_survey := new(survey_proto.DeleteRequest)
	req_survey.Id = req.PathParameter("survey_id")
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	delete_resp, err := p.SurveyClient.Delete(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.DeleteSurvey", "DeletError")
		return
	}

	delete_resp.Code = http.StatusOK
	delete_resp.Message = "Deleted survey successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, delete_resp)
}

/**
* @api {post} /server/surveys/survey/copy?session={session_id} Copy a survey
* @apiVersion 0.1.0
* @apiName Copy
* @apiGroup Survey
*
* @apiDescription This API copies survey object of requested survey id
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/surveys/survey/copy?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "survey_id": "111"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "survey": {
*       "id": "TTLg2rYHkFe7GNXh",
*       "org_id": "orgid",
*       "title": "survey",
*       "description": "This is sample survey",
*       "tags": "tags",
*       "created": 1517891917,
*       "updated": 1517891917,
*       "creator_id": "userid",
*       "creator": {User},
*       "shares": [{User},{User}...],
*       "renders": [],
*       "setting": {
*         "visibility": 0,
*         "notifications": [],
*         "social": [],
*         "linkSharingEnabled": true,
*         "embeddingEnabled": false,
*         "authentificationRequried": true,
*         "showCaptcha": true
*       },
*       "welcome": {DefaultQuestion},
*       "thankyou": {DefaultQuestion},
*       "questions": [{Question},{Question}...],
*       "design": {
*         "progress_bar_style": 0,
*         "default_bg_color": "#45AC3C",
*         "default_logo_url": "http://example.com/logo.png",
*       },
*       "status": 0,
*       "isTemplate": true,
*       "templateId": "templateId",
*     }
*   },
*   "code": 200,
*   "message": "Read all surveys successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The survey was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CopyError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.survey.CopySurvey",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *SurveyService) CopySurvey(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.Copy API request")
	req_survey := new(survey_proto.CopyRequest)
	err := req.ReadEntity(req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.CopySurvey", "BindError")
		return
	}
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	copy_resp, err := p.SurveyClient.Copy(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.CopySurvey", "CopyError")
		return
	}
	// fetchUserFromSurvey(ctx, copy_resp.Data.Survey, p.Auth.UserClient)

	copy_resp.Code = http.StatusOK
	copy_resp.Message = "Copied survey successfully"
	data := utils.MarshalAny(rsp, copy_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/surveys/survey/{survey_id}/questions?session={session_id}&offset={offset}&limit={limit} Get survey questions
* @apiVersion 0.1.0
* @apiName ReadQuestions
* @apiGroup Survey
*
* @apiDescription Return a list of questions for a survey with id = surveyid.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/surveys/survey/111/questions?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "questions": [
*       {
*         "id":"Hu7j6",
*         "type":7,
*         "order":2,
*         "design":{
*           "bg_color":"fffff",
*           "logo_url":"http://via.placeholder.com/300x300"
*         },
*         "settings":{
*           "@type":"healum.com/proto/go.micro.srv.survey.BinaryQuestionSettings",
*           "mandatory":true,
*           "image":false,
*           "video":false,
*           "buttonType":0
*         },
*         "title":"question 2?"
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read survey successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The survey was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.survey.Questions",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *SurveyService) Questions(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.Questions API request")
	req_survey := new(survey_proto.QuestionsRequest)
	req_survey.SurveyId = req.PathParameter("survey_id")
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	questions_resp, err := p.SurveyClient.Questions(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.Questions", "QueryError")
		return
	}

	questions_resp.Code = http.StatusOK
	questions_resp.Message = "Read questions successfully"
	data := utils.MarshalAny(rsp, questions_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/surveys/survey/{survey_id}/question/{question_id}?session={session_id} Get survey question
* @apiVersion 0.1.0
* @apiName QuestionRef
* @apiGroup Survey
*
* @apiDescription This API Return a question by question_id
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/surveys/survey/111/question/q111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
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
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The survey was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QuestionRefError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.survey.QuestionRef",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *SurveyService) QuestionRef(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.QuestionRef API request")
	req_survey := new(survey_proto.QuestionRefRequest)
	req_survey.SurveyId = req.PathParameter("survey_id")
	req_survey.QuestionRef = req.PathParameter("question_id")
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	ref_resp, err := p.SurveyClient.QuestionRef(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.QuestionRef", "QuestionRefError")
		return
	}

	ref_resp.Code = http.StatusOK
	ref_resp.Message = "Read question successfully"
	data := utils.MarshalAny(rsp, ref_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/surveys/survey/{survey_id}/questions?session={session_id} Add or Update a question
* @apiVersion 0.1.0
* @apiName Add/Update Question
* @apiGroup Survey
*
* @apiDescription Add or Update a question in the survey - Accepts a new list of questions for a survey
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/surveys/survey/111/questions?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "question": {
*     "id": "q111",
*     "type": QuestionType,
*     "order": 100,
*     "title": "question",
*     "description": "description",
*     "design":  {
*       "progress_bar_style": 0,
*       "default_bg_color": "#45AC3C",
*       "default_logo_url": "http://example.com/logo.png",
*     },
*     "fields": google.protobuf.Any,
*     "settings": google.protobuf.Any
*   },
* }
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
*   "message": "Created question successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The question was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateQuestion",
*       "errors": [
*         {
*           "domain": "go.micro.srv.survey.CreateQuestion",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *SurveyService) CreateQuestion(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.CreateQuestion API request")
	req_question := new(survey_proto.CreateQuestionRequest)
	if err := utils.UnmarshalAny(req, rsp, req_question); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.CreateQuestion", "BindError")
		return
	}

	req_question.SurveyId = req.PathParameter("survey_id")
	req_question.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_question.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	create_resp, err := p.SurveyClient.CreateQuestion(ctx, req_question)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.CreateQuestion", "CreateQuestionError")
		return
	}

	create_resp.Code = http.StatusOK
	create_resp.Message = "Created question successfully"
	data := utils.MarshalAny(rsp, create_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
 * @api {get} /server/surveys/creator/{userid}?session={session_id}&offset={offset}&limit={limit} Get all surveys created by a particular team member
 * @apiVersion 0.1.0
 * @apiName ByCreator
 * @apiGroup Survey
 *
 * @apiDescription Get all surveys created by a particular team member - Get all surveys where createdBy = {userid}
 *
 * @apiExample Example usage:
 * curl -i http://BASE_SERVER_URL/server/surveys/creator/userid?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
 *
 * @apiSuccessExample Success-Response:
 * HTTP/1.1 200 OK
 * {
 *   "data": {
 *     "surveys": [
 *       {
 *         "id": "111",
 *         "org_id": "orgid",
 *         "title": "survey",
 *         "description": "This is sample survey",
 *         "tags": "tags",
 *         "created": 1517891917,
 *         "updated": 1517891917,
 *         "creator_id": "userid",
 *         "creator": {User},
 *         "shares": [{User},{User}...],
 *         "renders": [],
 *         "setting": {
 *           "visibility": 0,
 *           "notifications": [],
 *           "social": [],
 *           "linkSharingEnabled": true,
 *           "embeddingEnabled": false,
 *           "authentificationRequried": true,
 *           "showCaptcha": true
 *         },
 *         "welcome": {DefaultQuestion},
 *         "thankyou": {DefaultQuestion},
 *         "questions": [{Question},{Question}...],
 *         "design": {
 *           "progress_bar_style": 0,
 *           "default_bg_color": "#45AC3C",
 *           "default_logo_url": "http://example.com/logo.png",
 *         },
 *         "status": 0,
 *         "isTemplate": true,
 *         "templateId": "templateId",
 *       },
 *       ... ...
 *     ]
 *   },
 *   "code": 200,
 *   "message": "Read all surveys successfully"
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
 *           "domain": "go.micro.srv.survey.ByCreator",
 *           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
 *         }
 *       ]
 *     }
 */
func (p *SurveyService) ByCreator(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.ByCreator API request")
	req_survey := new(survey_proto.ByCreatorRequest)
	req_survey.UserId = req.PathParameter("user_id")
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_survey.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_survey.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_survey.SortParameter = req.Attribute(SortParameter).(string)
	req_survey.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	creator_resp, err := p.SurveyClient.ByCreator(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.ByCreator", "QueryError")
		return
	}
	// for _, s := range creator_resp.Data.Surveys {
	// 	fetchUserFromSurvey(ctx, s, p.Auth.UserClient)
	// }

	creator_resp.Code = http.StatusOK
	creator_resp.Message = "Read all surveys successfully"
	data := utils.MarshalAny(rsp, creator_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

func (p *SurveyService) Link(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.Link API request")
	req_survey := new(survey_proto.LinkRequest)
	req_survey.SurveyId = req.PathParameter("survey_id")
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	link_resp, err := p.SurveyClient.Link(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.Link", "QueryError")
		return
	}

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, link_resp)
}

/**
 * @api {get} /server/surveys/creator/{userid}?session={session_id}&offset={offset}&limit={limit} Get all templates
 * @apiVersion 0.1.0
 * @apiName Templates
 * @apiGroup Survey
 *
 * @apiDescription Get all templates - Get all surveys where isTemplate = true
 *
 * @apiExample Example usage:
 * curl -i http://BASE_SERVER_URL/server/surveys/templates?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
 *
 * @apiSuccessExample Success-Response:
 * HTTP/1.1 200 OK
 * {
 *   "data": {
 *     "surveys": [
 *       {
 *         "id": "111",
 *         "org_id": "orgid",
 *         "title": "survey",
 *         "description": "This is sample survey",
 *         "tags": "tags",
 *         "created": 1517891917,
 *         "updated": 1517891917,
 *         "creator_id": "userid",
 *         "creator": {User},
 *         "shares": [{User},{User}...],
 *         "renders": [],
 *         "setting": {
 *           "visibility": 0,
 *           "notifications": [],
 *           "social": [],
 *           "linkSharingEnabled": true,
 *           "embeddingEnabled": false,
 *           "authentificationRequried": true,
 *           "showCaptcha": true
 *         },
 *         "welcome": {DefaultQuestion},
 *         "thankyou": {DefaultQuestion},
 *         "questions": [{Question},{Question}...],
 *         "design": {
 *           "progress_bar_style": 0,
 *           "default_bg_color": "#45AC3C",
 *           "default_logo_url": "http://example.com/logo.png",
 *         },
 *         "status": 0,
 *         "isTemplate": true,
 *         "templateId": "templateId",
 *       },
 *       ... ...
 *     ]
 *   },
 *   "code": 200,
 *   "message": "Searched surveys successfully"
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
 *           "domain": "go.micro.srv.survey.Templates",
 *           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
 *         }
 *       ]
 *     }
 */
func (p *SurveyService) Templates(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.Templates API request")
	req_survey := new(survey_proto.TemplatesRequest)
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_survey.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_survey.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_survey.SortParameter = req.Attribute(SortParameter).(string)
	req_survey.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	templates_resp, err := p.SurveyClient.Templates(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.Templates", "QueryError")
		return
	}
	// for _, s := range templates_resp.Data.Surveys {
	// 	fetchUserFromSurvey(ctx, s, p.Auth.UserClient)
	// }

	templates_resp.Code = http.StatusOK
	templates_resp.Message = "Searched surveys successfully"
	data := utils.MarshalAny(rsp, templates_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
 * @api {post} /server/surveys/filter?session={session_id}&offset={offset}&limit={limit} Filter surveys
 * @apiVersion 0.1.0
 * @apiName Filter
 * @apiGroup Survey
 *
 * @apiDescription Filter surveys by one or more status, one or more priority, one or more category status, priority or category are optional fields and one or more values maybe available in post body
 *
 * @apiExample Example usage:
 * curl -i http://BASE_SERVER_URL/server/surveys/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
 *
 * @apiParamExample {json} Request-Example:
 * {
 *   "status": [SurveyStatus, ...],
 *   "tags": ["tags"],
 *   "renderTarget": [RenderTarget, ...],
 *   "visibility": [Visiblity, ...],
 *   "created_by": ["user_id","user_id",...],
 * }
 *
 *
 * @apiSuccessExample Success-Response:
 * HTTP/1.1 200 OK
 * {
 *   "data": {
 *     "surveys": [
 *       {
 *         "id": "111",
 *         "org_id": "orgid",
 *         "title": "survey",
 *         "description": "This is sample survey",
 *         "tags": "tags",
 *         "created": 1517891917,
 *         "updated": 1517891917,
 *         "creator_id": "userid",
 *         "creator": {User},
 *         "shares": [{User},{User}...],
 *         "renders": [],
 *         "setting": {
 *           "visibility": 0,
 *           "notifications": [],
 *           "social": [],
 *           "linkSharingEnabled": true,
 *           "embeddingEnabled": false,
 *           "authentificationRequried": true,
 *           "showCaptcha": true
 *         },
 *         "welcome": {DefaultQuestion},
 *         "thankyou": {DefaultQuestion},
 *         "questions": [{Question},{Question}...],
 *         "design": {
 *           "progress_bar_style": 0,
 *           "default_bg_color": "#45AC3C",
 *           "default_logo_url": "http://example.com/logo.png",
 *         },
 *         "status": 0,
 *         "isTemplate": true,
 *         "templateId": "templateId",
 *       },
 *       ... ...
 *     ]
 *   },
 *   "code": 200,
 *   "message": "Searched surveys successfully"
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
 *           "domain": "go.micro.srv.survey.Filter",
 *           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
 *         }
 *       ]
 *     }
 */
func (p *SurveyService) Filter(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.Filter API request")
	req_survey := new(survey_proto.FilterRequest)
	err := utils.UnmarshalAny(req, rsp, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.Filter", "BindError")
		return
	}
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_survey.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_survey.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_survey.SortParameter = req.Attribute(SortParameter).(string)
	req_survey.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	filter_resp, err := p.SurveyClient.Filter(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.Filter", "QueryError")
		return
	}
	// for _, s := range filter_resp.Data.Surveys {
	// 	fetchUserFromSurvey(ctx, s, p.Auth.UserClient)
	// }

	filter_resp.Code = http.StatusOK
	filter_resp.Message = "Searched surveys successfully"
	data := utils.MarshalAny(rsp, filter_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
 * @api {post} /server/surveys/search?session={session_id}&offset={offset}&limit={limit} Search surveys
 * @apiVersion 0.1.0
 * @apiName Search
 * @apiGroup Survey
 *
 * @apiDescription Search surveys - Return searched surveys
 *
 * @apiExample Example usage:
 * curl -i http://BASE_SERVER_URL/server/surveys/search?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
 *
 * @apiParamExample {json} Request-Example:
 * {
 *   "name": "title",
 *   "description": "description"
 * }
 *
 *
 * @apiSuccessExample Success-Response:
 * HTTP/1.1 200 OK
 * {
 *   "data": {
 *     "surveys": [
 *       {
 *         "id": "111",
 *         "org_id": "orgid",
 *         "title": "survey",
 *         "description": "This is sample survey",
 *         "tags": "tags",
 *         "created": 1517891917,
 *         "updated": 1517891917,
 *         "creator_id": "userid",
 *         "creator": {User},
 *         "shares": [{User},{User}...],
 *         "renders": [],
 *         "setting": {
 *           "visibility": 0,
 *           "notifications": [],
 *           "social": [],
 *           "linkSharingEnabled": true,
 *           "embeddingEnabled": false,
 *           "authentificationRequried": true,
 *           "showCaptcha": true
 *         },
 *         "welcome": {DefaultQuestion},
 *         "thankyou": {DefaultQuestion},
 *         "questions": [{Question},{Question}...],
 *         "design": {
 *           "progress_bar_style": 0,
 *           "default_bg_color": "#45AC3C",
 *           "default_logo_url": "http://example.com/logo.png",
 *         },
 *         "status": 0,
 *         "isTemplate": true,
 *         "templateId": "templateId",
 *       },
 *       ... ...
 *     ]
 *   },
 *   "code": 200,
 *   "message": "Searched surveys successfully"
 * }
 *
 * @apiError NoAuthorized 	Only authenticated users can access the data.
 * @apiError BadRequest   	The surveys were not found.
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 400 Bad Request
 *     {
 *       "code": 400,
 *       "message": "SearchError",
 *       "errors": [
 *         {
 *           "domain": "go.micro.srv.survey.SearchSurveys",
 *           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
 *         }
 *       ]
 *     }
 */
func (p *SurveyService) SearchSurveys(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.Search API request")
	req_survey := new(survey_proto.SearchRequest)
	err := req.ReadEntity(req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.SearchSurvey", "BindError")
		return
	}
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_survey.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_survey.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_survey.SortParameter = req.Attribute(SortParameter).(string)
	req_survey.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	search_resp, err := p.SurveyClient.Search(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.SearchSurvey", "SearchError")
		return
	}
	// for _, s := range search_resp.Data.Surveys {
	// 	fetchUserFromSurvey(ctx, s, p.Auth.UserClient)
	// }

	search_resp.Code = http.StatusOK
	search_resp.Message = "Searched surveys successfully"
	data := utils.MarshalAny(rsp, search_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/surveys/survey/search/autocomplete?session={session_id} autocomplete text search for surveys
* @apiVersion 0.1.0
* @apiName AutocompleteSurveySearch
* @apiGroup Survey
*
* @apiDescription Should return a list of surveys based on text based search. This should not be paginated
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/surveys/survey/search/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
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
*         "id": "111",
*         "title": "surveyh_title",
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
*   "message": "Read surveys successfully"
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
*           "domain": "go.micro.srv.survey.AutocompleteSurveySearch",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *SurveyService) AutocompleteSurveySearch(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.AutocompleteSurveySearch API request")
	req_search := new(survey_proto.AutocompleteSearchRequest)
	err := req.ReadEntity(req_search)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.AutocompleteSurveySearch", "BindError")
		return
	}
	// req_search.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_search.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.SurveyClient.AutocompleteSearch(ctx, req_search)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.AutocompleteSurveySearch", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read surveys successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/surveys/tags/top/{n}?session={session_id} Return top N tags for Survey
* @apiVersion 0.1.0
* @apiName GetTopSurveyTags
* @apiGroup Survey
*
* @apiDescription For each of the following service we have return top N tags for survey
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/surveys/tags/top/5?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "tags": ["tag1","tag2","tag3",...]
*   },
*   "code": 200,
*   "message": "Get top survey tags successfully"
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
*           "domain": "go.micro.srv.survey.GetTopSurveyTags",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *SurveyService) GetTopSurveyTags(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.GetTopSurveyTags API request")
	req_survey := new(survey_proto.GetTopTagsRequest)
	n, _ := strconv.Atoi(req.PathParameter("n"))
	req_survey.N = int64(n)
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.SurveyClient.GetTopTags(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.GetTopSurveyTags", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Get top survey tags successfully"
	rsp.AddHeader("Survey-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/surveys/tags/autocomplete?session={session_id} Autocomplete for tags for Survey
* @apiVersion 0.1.0
* @apiName AutocompleteSurveyTags
* @apiGroup Survey
*
* @apiDescription Autocomplete for tags for Survey
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/surveys/tags/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
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
*   "message": "Autocomplete survey tags successfully"
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
*           "domain": "go.micro.srv.survey.AutocompleteSurveyTags",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *SurveyService) AutocompleteSurveyTags(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Survey.AutocompleteSurveyTags API request")
	req_survey := new(survey_proto.AutocompleteTagsRequest)
	err := req.ReadEntity(req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.AutocompleteSurveyTags", "BindError")
		return
	}
	req_survey.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_survey.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.SurveyClient.AutocompleteTags(ctx, req_survey)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.survey.AutocompleteSurveyTags", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Autocomplete survey tags successfully"
	rsp.AddHeader("Survey-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}
