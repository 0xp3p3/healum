package api

import (
	"context"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	note_proto "server/note-srv/proto/note"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type NoteService struct {
	NoteClient    note_proto.NoteServiceClient
	Auth          Filters
	Audit         AuditFilter
	ServerMetrics metrics.Metrics
}

func (p NoteService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/notes")

	audit := &audit_proto.Audit{
		ActionService:  common.NoteSrv,
		ActionResource: common.BASE + common.NOTE_TYPE,
	}

	ws.Route(ws.GET("/all").To(p.AllNotes).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all notes"))

	ws.Route(ws.POST("/note").To(p.CreateNote).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("create data"))

	ws.Route(ws.GET("/note/{note_id}").To(p.ReadNote).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("View Note detail"))

	ws.Route(ws.DELETE("/note/{note_id}").To(p.DeleteNote).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("View Note detail"))

	ws.Route(ws.POST("/search").To(p.SearchNotes).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all notes"))

	ws.Route(ws.GET("/creator/{user_id}").To(p.ByCreator).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all notes where the status is draft"))

	ws.Route(ws.GET("/user/{user_id}").To(p.ByUser).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all notes where the status is draft"))

	ws.Route(ws.POST("/filter").To(p.Filter).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Filter notes by one or more category, one or more tags"))

	restful.Add(ws)
}

/**
* @api {get} /server/notes/all?session={session_id}&offset={offset}&limit={limit} List all notes
* @apiVersion 0.1.0
* @apiName AllNotes
* @apiGroup Note
*
* @apiDescription AllNotes
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/notes/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "notes": [
*       {
*         "id": "111",
*         "title": "note1",
*         "orgid": "orgid",
*         "description": "description1",
*         "creator":  { User },
*         "user": { User },
*         "tags": ["a","b","c"]
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
func (p *NoteService) AllNotes(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Note.All API request")
	req_note := new(note_proto.AllRequest)
	req_note.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_note.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_note.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_note.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_note.SortParameter = req.Attribute(SortParameter).(string)
	req_note.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	all_resp, err := p.NoteClient.All(ctx, req_note)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.note.AllNotes", "QueryError")
		return
	}
	all_resp.Code = http.StatusOK
	all_resp.Message = "Read all notes succesfully"
	data := utils.MarshalAny(rsp, all_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/notes/note?session={session_id} Create a note
* @apiVersion 0.1.0
* @apiName CreateNote
* @apiGroup Note
*
* @apiDescription Create a note
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/notes/note?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "note": {
*     "id": "111",
*     "title": "note1",
*     "orgid": "orgid",
*     "description": "description1",
*     "creator":  { User },
*     "user": { User },
*     "tags": ["a","b","c"]
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "note": {
*         "id": "111",
*         "title": "note1",
*         "orgid": "orgid",
*         "description": "description1",
*         "creator":  { User },
*         "user": { User },
*         "tags": ["a","b","c"]
*     }
*   },
*   "code": 200,
*   "message": "Created note succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The notes were not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.note.CreateNote",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *NoteService) CreateNote(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Note.Create API request")
	req_note := new(note_proto.CreateRequest)
	// err := req.ReadEntity(req_note)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.note.CreateNote", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_note); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.note.CreateNote", "BindError")
		return
	}
	req_note.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_note.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	create_resp, err := p.NoteClient.Create(ctx, req_note)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.note.CreateNote", "CreateError")
		return
	}
	create_resp.Code = http.StatusOK
	create_resp.Message = "Create note succesfully"
	data := utils.MarshalAny(rsp, create_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/notes/note/{note_id}?session={session_id} View note detail
* @apiVersion 0.1.0
* @apiName ReadNote
* @apiGroup Note
*
* @apiDescription View note detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/notes/note/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "note": {
*         "id": "111",
*         "title": "note1",
*         "orgid": "orgid",
*         "description": "description1",
*         "creator":  { User },
*         "user": { User },
*         "tags": ["a","b","c"]
*     }
*   },
*   "code": 200,
*   "message": "Read note succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The note were not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.note.ReadNote",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *NoteService) ReadNote(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Note.Read API request")
	req_note := new(note_proto.ReadRequest)
	req_note.Id = req.PathParameter("note_id")
	req_note.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_note.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	read_resp, err := p.NoteClient.Read(ctx, req_note)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.note.ReadNote", "ReadError")
		return
	}

	read_resp.Code = http.StatusOK
	read_resp.Message = "Read note succesfully"
	data := utils.MarshalAny(rsp, read_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {delete} /server/notes/note/{note_id}?session={session_id} Delete a note
* @apiVersion 0.1.0
* @apiName DeleteNote
* @apiGroup Note
*
* @apiDescription Delete a note
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/notes/note/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted note succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The note was not updated.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.note.DeleteNote",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *NoteService) DeleteNote(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Note.Delete API request")
	req_note := new(note_proto.DeleteRequest)
	req_note.Id = req.PathParameter("note_id")
	req_note.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_note.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	delete_resp, err := p.NoteClient.Delete(ctx, req_note)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.note.DeleteNote", "DeleteError")
		return
	}
	delete_resp.Code = http.StatusOK
	delete_resp.Message = "Deleted note succesfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, delete_resp)
}

/**
* @api {post} /server/notes/search?session={session_id}&offset={offset}&limit={limit} Search notes
* @apiVersion 0.1.0
* @apiName SearchNotes
* @apiGroup Note
*
* @apiDescription SearchNotes
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/notes/search?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "name": "note1"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "notes": [
*       {
*         "id": "111",
*         "title": "note1",
*         "orgid": "orgid",
*         "description": "description1",
*         "creator":  { User },
*         "user": { User },
*         "tags": ["a","b","c"]
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched notes successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The notes were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "SearchError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.note.SearchNotes",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *NoteService) SearchNotes(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Note.Search API request")
	req_note := new(note_proto.SearchRequest)
	err := req.ReadEntity(req_note)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.note.SearchNotes", "BindError")
		return
	}
	req_note.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_note.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_note.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_note.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_note.SortParameter = req.Attribute(SortParameter).(string)
	req_note.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	search_resp, err := p.NoteClient.Search(ctx, req_note)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.note.SearchNotes", "SearchError")
		return
	}

	search_resp.Code = http.StatusOK
	search_resp.Message = "Searched note succesfully"
	data := utils.MarshalAny(rsp, search_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/notes/user/{user_id}?session={session_id}&offset={offset}&limit={limit} Get all notes created by a particular team member
* @apiVersion 0.1.0
* @apiName ByUser
* @apiGroup Note
*
* @apiDescription Get all notes created by a particular team member - Get all notes where creator = {userid}
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/notes/user/userid?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "notes": [
*       {
*         "id": "111",
*         "title": "note1",
*         "orgid": "orgid",
*         "description": "description1",
*         "creator":  { User },
*         "user": { User },
*         "tags": ["a","b","c"]
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched notes succesfully"
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
*           "domain": "go.micro.srv.note.ByUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *NoteService) ByUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Note.ByUser API request")
	req_note := new(note_proto.ByUserRequest)
	req_note.UserId = req.PathParameter("user_id")
	req_note.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_note.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_note.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_note.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_note.SortParameter = req.Attribute(SortParameter).(string)
	req_note.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	byuser_resp, err := p.NoteClient.ByUser(ctx, req_note)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.note.ByUser", "QueryError")
		return
	}

	byuser_resp.Code = http.StatusOK
	byuser_resp.Message = "Searched notes succesfully"
	data := utils.MarshalAny(rsp, byuser_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/notes/creator/{user_id}?session={session_id}&offset={offset}&limit={limit} Get all notes created for a user
* @apiVersion 0.1.0
* @apiName ByCreator
* @apiGroup Note
*
* @apiDescription Get all notes created for a user - Get all notes where user = {userid}
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/notes/creator/userid?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "notes": [
*       {
*         "id": "111",
*         "title": "note1",
*         "orgid": "orgid",
*         "description": "description1",
*         "creator":  { User },
*         "user": { User },
*         "tags": ["a","b","c"]
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched notes succesfully"
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
*           "domain": "go.micro.srv.note.ByUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *NoteService) ByCreator(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Note.ByCreator API request")
	req_note := new(note_proto.ByCreatorRequest)
	req_note.UserId = req.PathParameter("user_id")
	req_note.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_note.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_note.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_note.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_note.SortParameter = req.Attribute(SortParameter).(string)
	req_note.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	bycreator_resp, err := p.NoteClient.ByCreator(ctx, req_note)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.note.ByCreator", "QueryError")
		return
	}

	bycreator_resp.Code = http.StatusOK
	bycreator_resp.Message = "Searched notes succesfully"
	data := utils.MarshalAny(rsp, bycreator_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/notes/filter?session={session_id}&offset={offset}&limit={limit} Filter notes
* @apiVersion 0.1.0
* @apiName Filter
* @apiGroup Note
*
* @apiDescription Filter notes by one or more status, one or more priority, one or more category status, priority or category are optional fields and one or more values maybe available in post body
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/notes/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "category": ["category1"],
*   "tags": ["a","b"]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "notes": [
*       {
*         "id": "111",
*         "title": "note1",
*         "orgid": "orgid",
*         "description": "description1",
*         "creator":  { User },
*         "user": { User },
*         "tags": ["a","b","c"]
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Filtered notes succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The notes were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "FilterError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.note.Filter",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *NoteService) Filter(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Note.Filter API request")
	req_note := new(note_proto.FilterRequest)
	err := req.ReadEntity(req_note)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.note.Filter", "BindError")
		return
	}
	req_note.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_note.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_note.SortParameter = req.Attribute(SortParameter).(string)
	req_note.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	filter_resp, err := p.NoteClient.Filter(ctx, req_note)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.note.Filter", "FilterError")
		return
	}

	filter_resp.Code = http.StatusOK
	filter_resp.Message = "Filtered notes succesfully"
	data := utils.MarshalAny(rsp, filter_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}
