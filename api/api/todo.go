package api

import (
	"context"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	todo_proto "server/todo-srv/proto/todo"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type TodoService struct {
	TodoClient    todo_proto.TodoServiceClient
	Auth          Filters
	Audit         AuditFilter
	ServerMetrics metrics.Metrics
}

func (p TodoService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/todos")

	audit := &audit_proto.Audit{
		ActionService:  common.TodoSrv,
		ActionResource: common.BASE + common.TODO_TYPE,
	}

	ws.Route(ws.GET("/all").To(p.AllTodos).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all todos"))

	ws.Route(ws.POST("/todo").To(p.CreateTodo).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("create data"))

	ws.Route(ws.GET("/todo/{todo_id}").To(p.ReadTodo).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("View Todo detail"))

	ws.Route(ws.DELETE("/todo/{todo_id}").To(p.DeleteTodo).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete Todo detail"))

	ws.Route(ws.POST("/search").To(p.SearchTodos).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search todos"))

	ws.Route(ws.GET("/creator/{user_id}").To(p.ByCreator).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all todos created by a particular team member"))

	restful.Add(ws)
}

/**
* @api {get} /server/todos/all?session={session_id}&offset={offset}&limit={limit} List all todos
* @apiVersion 0.1.0
* @apiName AllTodos
* @apiGroup Todo
*
* @apiDescription AllTodos
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/todos/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "todos": [
*       {
*         "id": "111",
*         "title": "todo1",
*         "orgid": "orgid",
*         "creatorId": "userId",
*         "creator":  { User },
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all todos successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The todos were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.todo.AllTodos",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TodoService) AllTodos(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Todo.All API request")
	req_todo := new(todo_proto.AllRequest)
	req_todo.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_todo.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_todo.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_todo.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_todo.SortParameter = req.Attribute(SortParameter).(string)
	req_todo.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	all_resp, err := p.TodoClient.All(ctx, req_todo)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.todo.AllTodos", "QueryError")
		return
	}
	all_resp.Code = http.StatusOK
	all_resp.Message = "Read all todos successfully"
	data := utils.MarshalAny(rsp, all_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/todos/todo?session={session_id} Create a todo
* @apiVersion 0.1.0
* @apiName CreateTodo
* @apiGroup Todo
*
* @apiDescription Create a todo
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/todos/todo?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "todo": {
*     "id": "111",
*     "title": "todo1",
*     "orgid": "orgid",
*     "creatorId": "userId"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "todo": {
*       "id": "111",
*       "title": "todo1",
*       "orgid": "orgid",
*       "creatorId": "userId",
*       "creator":  { User }
*     }
*   },
*   "code": 200,
*   "message": "Created todo succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The todos were not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.todo.CreateTodo",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TodoService) CreateTodo(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Todo.Create API request")
	req_todo := new(todo_proto.CreateRequest)
	// err := req.ReadEntity(req_todo)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.todo.CreateTodo", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_todo); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.todo.CreateTodo", "BindError")
		return
	}
	req_todo.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_todo.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	create_resp, err := p.TodoClient.Create(ctx, req_todo)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.todo.CreateTodo", "CreateError")
		return
	}
	create_resp.Code = http.StatusOK
	create_resp.Message = "Created todo successfully"
	data := utils.MarshalAny(rsp, create_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/todos/todo/{todo_id}?session={session_id} View todo detail
* @apiVersion 0.1.0
* @apiName ReadTodo
* @apiGroup Todo
*
* @apiDescription View todo detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/todos/todo/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "todo": {
*       "id": "111",
*       "title": "todo1",
*       "orgid": "orgid",
*       "creatorId": "userId",
*       "creator":  { User }
*     }
*   },
*   "code": 200,
*   "message": "Read todo succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The todo were not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.todo.ReadTodo",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TodoService) ReadTodo(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Todo.Read API request")
	req_todo := new(todo_proto.ReadRequest)
	req_todo.Id = req.PathParameter("todo_id")
	req_todo.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_todo.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	read_resp, err := p.TodoClient.Read(ctx, req_todo)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.todo.ReadTodo", "ReadError")
		return
	}
	read_resp.Code = http.StatusOK
	read_resp.Message = "Read todo successfully"
	data := utils.MarshalAny(rsp, read_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {delete} /server/todos/todo/{todo_id}?session={session_id} Delete a todo
* @apiVersion 0.1.0
* @apiName DeleteTodo
* @apiGroup Todo
*
* @apiDescription Delete a todo
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/todos/todo/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted todo succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The todo was not updated.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.todo.DeleteTodo",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TodoService) DeleteTodo(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Todo.Delete API request")
	req_todo := new(todo_proto.DeleteRequest)
	req_todo.Id = req.PathParameter("todo_id")
	req_todo.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_todo.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	delete_resp, err := p.TodoClient.Delete(ctx, req_todo)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.todo.DeleteTodo", "DeletError")
		return
	}

	delete_resp.Code = http.StatusOK
	delete_resp.Message = "Deleted todo successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, delete_resp)
}

/**
* @api {post} /server/todos/search?session={session_id}&offset={offset}&limit={limit} Search todos
* @apiVersion 0.1.0
* @apiName SearchTodos
* @apiGroup Todo
*
* @apiDescription SearchTodos
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/todos/search?session={session_id}&offset={offset}&limit={limit}
*
* @apiParamExample {json} Request-Example:
* {
*   "name": "todo1"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "todos": [
*       {
*         "id": "111",
*         "title": "todo1",
*         "orgid": "orgid",
*         "creatorId": "userId",
*         "creator":  { User },
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched todos successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The todos were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "SearchError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.todo.SearchTodos",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TodoService) SearchTodos(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Todo.Search API request")
	req_todo := new(todo_proto.SearchRequest)
	err := req.ReadEntity(req_todo)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.todo.SearchTodo", "BindError")
		return
	}
	req_todo.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_todo.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_todo.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_todo.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_todo.SortParameter = req.Attribute(SortParameter).(string)
	req_todo.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	search_resp, err := p.TodoClient.Search(ctx, req_todo)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.todo.SearchTodo", "SearchError")
		return
	}
	search_resp.Code = http.StatusOK
	search_resp.Message = "Searched todos successfully"
	data := utils.MarshalAny(rsp, search_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/todos/creator/{user_id}?session={session_id}&offset={offset}&limit={limit} Get all todos created for a user
* @apiVersion 0.1.0
* @apiName ByCreator
* @apiGroup Todo
*
* @apiDescription Get all todos created for a user - Get all todos where user = {userid}
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/todos/creator/userid?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "todos": [
*       {
*         "id": "111",
*         "title": "todo1",
*         "orgid": "orgid",
*         "creatorId": "userId",
*         "creator":  { User },
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched todos succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The todos were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.todo.ByUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TodoService) ByCreator(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Todo.ByCreator API request")
	req_todo := new(todo_proto.ByCreatorRequest)
	req_todo.UserId = req.PathParameter("user_id")
	req_todo.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_todo.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_todo.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_todo.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_todo.SortParameter = req.Attribute(SortParameter).(string)
	req_todo.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	creator_resp, err := p.TodoClient.ByCreator(ctx, req_todo)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.todo.ByCreator", "QueryError")
		return
	}
	creator_resp.Code = http.StatusOK
	creator_resp.Message = "Searched todos successfully"
	data := utils.MarshalAny(rsp, creator_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}
