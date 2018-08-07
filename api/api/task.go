package api

import (
	"context"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	organisation_proto "server/organisation-srv/proto/organisation"
	task_proto "server/task-srv/proto/task"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type TaskService struct {
	TaskClient         task_proto.TaskServiceClient
	Auth               Filters
	Audit              AuditFilter
	OrganisationClient organisation_proto.OrganisationServiceClient
	FilterMiddle       Filters
	ServerMetrics      metrics.Metrics
}

func (p TaskService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/tasks")

	audit := &audit_proto.Audit{
		ActionService:  common.TaskSrv,
		ActionResource: common.BASE + common.TASK_TYPE,
	}

	ws.Route(ws.GET("/all").To(p.AllTasks).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all tasks"))

	ws.Route(ws.POST("/task").To(p.CreateTask).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("create data"))

	ws.Route(ws.GET("/task/{task_id}").To(p.ReadTask).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("View Task detail"))

	ws.Route(ws.DELETE("/task/{task_id}").To(p.DeleteTask).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("View Task detail"))

	ws.Route(ws.POST("/search").To(p.SearchTasks).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all tasks"))

	ws.Route(ws.GET("/creator/{user_id}").To(p.ByCreator).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all tasks created by a particular team member"))

	ws.Route(ws.GET("/assign/{user_id}").To(p.ByAssign).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Get all tasks assigned to a particular team member "))

	ws.Route(ws.POST("/filter").To(p.Filter).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Filter tasks by one or more status, one or more priority, one or more category status"))

	ws.Route(ws.GET("/count/{user_id}").To(p.CountByUser).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Return tasks count of expired tasks, assigned to the user tasks"))

	restful.Add(ws)
}

/**
* @api {get} /server/tasks/all?session={session_id}&offset={offset}&limit={limit} List all tasks
* @apiVersion 0.1.0
* @apiName AllTasks
* @apiGroup Task
*
* @apiDescription AllTasks
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/tasks/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "tasks": [
*       {
*         "id": "111",
*         "title": "task1",
*         "orgid": "orgid",
*         "description": "description1",
*         "creatorId": "userId",
*         "creator":  { User },
*         "assigneeId": "userId",
*         "assignee": { User },
*         "category": "category1",
*         "due": 1517891917,
*         "tags": ["a","b","c"]
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Read all tasks successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The tasks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.task.AllTasks",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TaskService) AllTasks(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Task.All API request")
	req_task := new(task_proto.AllRequest)
	req_task.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_task.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_task.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_task.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_task.SortParameter = req.Attribute(SortParameter).(string)
	req_task.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	all_resp, err := p.TaskClient.All(ctx, req_task)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.task.AllTodos", "QueryError")
		return
	}
	all_resp.Code = http.StatusOK
	all_resp.Message = "Read all tasks successfully"
	data := utils.MarshalAny(rsp, all_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/tasks/task?session={session_id} Create a task
* @apiVersion 0.1.0
* @apiName CreateTask
* @apiGroup Task
*
* @apiDescription Create a task
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/tasks/task?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "task": {
*     "id": "111",
*     "title": "task1",
*     "orgid": "orgid",
*     "description": "description1",
*     "creatorId": "userId",
*     "assigneeId": "userId",
*     "category": "category1",
*     "due": 1517891917,
*     "tags": ["a","b","c"]
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "task": {
*       "id": "111",
*       "title": "task1",
*       "orgid": "orgid",
*       "description": "description1",
*       "creatorId": "userId",
*       "creator":  { User },
*       "assigneeId": "userId",
*       "assignee": { User },
*       "category": "category1",
*       "due": 1517891917,
*       "tags": ["a","b","c"]
*     }
*   },
*   "code": 200,
*   "message": "Created task succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The tasks were not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.task.CreateTask",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TaskService) CreateTask(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Task.Create API request")
	req_task := new(task_proto.CreateRequest)
	// err := req.ReadEntity(req_task)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.task.CreateTask", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_task); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.task.CreateTask", "BindError")
		return
	}
	req_task.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_task.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	create_resp, err := p.TaskClient.Create(ctx, req_task)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.task.CreateTask", "CreateError")
		return
	}
	create_resp.Code = http.StatusOK
	create_resp.Message = "Created task successfully"
	data := utils.MarshalAny(rsp, create_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/tasks/task/{task_id}?session={session_id} View task detail
* @apiVersion 0.1.0
* @apiName ReadTask
* @apiGroup Task
*
* @apiDescription View task detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/tasks/task/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "task": {
*       "id": "111",
*       "title": "task1",
*       "orgid": "orgid",
*       "description": "description1",
*       "creatorId": "userId",
*       "creator":  { User },
*       "assigneeId": "userId",
*       "assignee": { User },
*       "category": "category1",
*       "due": 1517891917,
*       "tags": ["a","b","c"]
*     }
*   },
*   "code": 200,
*   "message": "Read task succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The task were not created.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.task.ReadTask",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TaskService) ReadTask(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Task.Read API request")
	req_task := new(task_proto.ReadRequest)
	req_task.Id = req.PathParameter("task_id")
	req_task.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_task.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	read_resp, err := p.TaskClient.Read(ctx, req_task)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.task.ReadTodo", "ReadError")
		return
	}
	read_resp.Code = http.StatusOK
	read_resp.Message = "Read task successfully"
	data := utils.MarshalAny(rsp, read_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {delete} /server/tasks/task/{task_id}?session={session_id} Delete a task
* @apiVersion 0.1.0
* @apiName DeleteTask
* @apiGroup Task
*
* @apiDescription Delete a task
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/tasks/task/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted task succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The task was not updated.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.task.DeleteTask",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TaskService) DeleteTask(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Task.Delete API request")
	req_task := new(task_proto.DeleteRequest)
	req_task.Id = req.PathParameter("task_id")
	req_task.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_task.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	delete_resp, err := p.TaskClient.Delete(ctx, req_task)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.task.DeleteTask", "DeleteError")
		return
	}

	delete_resp.Code = http.StatusOK
	delete_resp.Message = "Deleted task successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, delete_resp)
}

/**
* @api {post} /server/tasks/search?session={session_id}&offset={offset}&limit={limit} Search tasks
* @apiVersion 0.1.0
* @apiName SearchTasks
* @apiGroup Task
*
* @apiDescription SearchTasks
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/tasks/search?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "name": "task1"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "tasks": [
*       {
*         "id": "111",
*         "title": "task1",
*         "orgid": "orgid",
*         "description": "description1",
*         "creatorId": "userId",
*         "creator":  { User },
*         "assigneeId": "userId",
*         "assignee": { User },
*         "category": "category1",
*         "due": 1517891917,
*         "tags": ["a","b","c"]
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched tasks successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The tasks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "SearchError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.task.SearchTasks",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TaskService) SearchTasks(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Task.Search API request")
	req_task := new(task_proto.SearchRequest)
	err := req.ReadEntity(req_task)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.task.SearchTasks", "BindError")
		return
	}
	req_task.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_task.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_task.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_task.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_task.SortParameter = req.Attribute(SortParameter).(string)
	req_task.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	search_resp, err := p.TaskClient.Search(ctx, req_task)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.task.SearchTasks", "SearchError")
		return
	}
	search_resp.Code = http.StatusOK
	search_resp.Message = "Searched tasks successfully"
	data := utils.MarshalAny(rsp, search_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/tasks/creator/{user_id}?session={session_id} Get all tasks created for a user
* @apiVersion 0.1.0
* @apiName ByCreator
* @apiGroup Task
*
* @apiDescription Get all tasks created for a user - Get all tasks where user = {userid}
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/tasks/creator/userid?session={session_id}
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "tasks": [
*       {
*         "id": "111",
*         "title": "task1",
*         "orgid": "orgid",
*         "description": "description1",
*         "creatorId": "userId",
*         "creator":  { User },
*         "assigneeId": "userId",
*         "assignee": { User },
*         "category": "category1",
*         "due": 1517891917,
*         "tags": ["a","b","c"]
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched tasks succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The tasks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.task.ByUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TaskService) ByCreator(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Task.ByCreator API request")
	req_task := new(task_proto.ByCreatorRequest)
	req_task.UserId = req.PathParameter("user_id")
	req_task.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_task.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_task.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_task.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_task.SortParameter = req.Attribute(SortParameter).(string)
	req_task.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	creator_resp, err := p.TaskClient.ByCreator(ctx, req_task)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.task.ByCreator", "QueryError")
		return
	}
	creator_resp.Code = http.StatusOK
	creator_resp.Message = "Searched tasks successfully"
	data := utils.MarshalAny(rsp, creator_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/tasks/assign/{user_id}?session={session_id}&offset={offset}&limit={limit} Get all tasks assigned to a particular team member
* @apiVersion 0.1.0
* @apiName ByAssign
* @apiGroup Task
*
* @apiDescription Get all tasks assigned to a particular team member - Get all tasks where assignee = {userid}
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/tasks/assign/userid?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "tasks": [
*       {
*         "id": "111",
*         "title": "task1",
*         "orgid": "orgid",
*         "description": "description1",
*         "creatorId": "userId",
*         "creator":  { User },
*         "assigneeId": "userId",
*         "assignee": { User },
*         "category": "category1",
*         "due": 1517891917,
*         "tags": ["a","b","c"]
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Searched tasks succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The tasks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.task.ByUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TaskService) ByAssign(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Task.ByAssign API request")
	req_task := new(task_proto.ByAssignRequest)
	req_task.UserId = req.PathParameter("user_id")
	req_task.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_task.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_task.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_task.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_task.SortParameter = req.Attribute(SortParameter).(string)
	req_task.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	assign_resp, err := p.TaskClient.ByAssign(ctx, req_task)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.task.ByAssign", "QueryError")
		return
	}
	assign_resp.Code = http.StatusOK
	assign_resp.Message = "Searched tasks successfully"
	data := utils.MarshalAny(rsp, assign_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/tasks/filter?session={session_id}&offset={offset}&limit={limit} Filter tasks
* @apiVersion 0.1.0
* @apiName Filter
* @apiGroup Task
*
* @apiDescription Filter tasks by one or more status, one or more priority, one or more category status, priority or category are optional fields and one or more values maybe available in post body
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/tasks/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "status": [],
*   "category": ["category1"],
*   "priority": []
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "tasks": [
*       {
*         "id": "111",
*         "title": "task1",
*         "orgid": "orgid",
*         "description": "description1",
*         "creatorId": "userId",
*         "creator":  { User },
*         "assigneeId": "userId",
*         "assignee": { User },
*         "category": "category1",
*         "due": 1517891917,
*         "tags": ["a","b","c"]
*       },
*       ... ...
*      ]
*   },
*   "code": 200,
*   "message": "Filtered tasks succesfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The tasks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "FilterError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.task.Filter",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TaskService) Filter(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Task.Filter API request")
	req_task := new(task_proto.FilterRequest)
	req_task.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_task.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_task.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_task.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_task.SortParameter = req.Attribute(SortParameter).(string)
	req_task.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	filter_resp, err := p.TaskClient.Filter(ctx, req_task)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.task.Filter", "FilterError")
		return
	}
	filter_resp.Code = http.StatusOK
	filter_resp.Message = "Filtered tasks successfully"
	data := utils.MarshalAny(rsp, filter_resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/tasks/count/{user_id}?session={session_id}&offset={offset}&limit={limit} Get task Counts
* @apiVersion 0.1.0
* @apiName CountByUser
* @apiGroup Task
*
* @apiDescription Get task Counts - Return tasks count of expired tasks, assigned to the user tasks
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/tasks/count/userid?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "task_count": {
*       "expired": 8,
*       "assigned": 11
*     }
*   },
*   "code": 200,
*   "message": "Queried successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The tasks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.task.CountByUser",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *TaskService) CountByUser(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Task.CountByUser API request")
	req_task := new(task_proto.CountByUserRequest)
	req_task.UserId = req.PathParameter("user_id")
	req_task.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_task.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	count_resp, err := p.TaskClient.CountByUser(ctx, req_task)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.task.CountByUser", "QueryError")
		return
	}

	count_resp.Code = http.StatusOK
	count_resp.Message = "Queried successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, count_resp)
}
