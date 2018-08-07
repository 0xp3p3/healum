package handler

import (
	"context"
	"server/common"
	"server/todo-srv/db"
	todo_proto "server/todo-srv/proto/todo"

	log "github.com/sirupsen/logrus"
)

type TodoService struct{}

func (p *TodoService) All(ctx context.Context, req *todo_proto.AllRequest, rsp *todo_proto.AllResponse) error {
	log.Info("Received Todo.All request")
	todos, err := db.All(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(todos) == 0 || err != nil {
		return common.NotFound(common.TodoSrv, p.All, err, "not found")
	}
	rsp.Data = &todo_proto.ArrData{todos}
	return nil
}

func (p *TodoService) Create(ctx context.Context, req *todo_proto.CreateRequest, rsp *todo_proto.CreateResponse) error {
	log.Info("Received Todo.Create request")
	if len(req.Todo.Name) == 0 {
		return common.BadRequest(common.TodoSrv, p.Create, nil, "todo name empty")
	}
	if req.Todo.Creator == nil {
		return common.BadRequest(common.TodoSrv, p.Create, nil, "todo creator empty")
	}

	err := db.Create(ctx, req.Todo)
	if err != nil {
		return common.InternalServerError(common.TodoSrv, p.Create, err, "create error")
	}
	rsp.Data = &todo_proto.Data{req.Todo}
	return nil
}

func (p *TodoService) Read(ctx context.Context, req *todo_proto.ReadRequest, rsp *todo_proto.ReadResponse) error {
	log.Info("Received Todo.Read request")
	todo, err := db.Read(ctx, req.Id, req.OrgId, req.TeamId)
	if todo == nil || err != nil {
		return common.NotFound(common.TodoSrv, p.Read, err, "not found")
	}
	rsp.Data = &todo_proto.Data{todo}
	return nil
}

func (p *TodoService) Delete(ctx context.Context, req *todo_proto.DeleteRequest, rsp *todo_proto.DeleteResponse) error {
	log.Info("Received Todo.Delete request")
	if err := db.Delete(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.TodoSrv, p.Delete, err, "delete error")
	}
	return nil
}

func (p *TodoService) Search(ctx context.Context, req *todo_proto.SearchRequest, rsp *todo_proto.SearchResponse) error {
	log.Info("Received Todo.Search request")
	todos, err := db.Search(ctx, req.Name, req.OrgId, req.TeamId, req.Limit, req.Offset, req.From, req.To, req.SortParameter, req.SortDirection)
	if len(todos) == 0 || err != nil {
		return common.NotFound(common.TodoSrv, p.Search, err, "not found")
	}
	rsp.Data = &todo_proto.ArrData{todos}
	return nil
}

func (p *TodoService) ByCreator(ctx context.Context, req *todo_proto.ByCreatorRequest, rsp *todo_proto.ByCreatorResponse) error {
	log.Info("Received Todo.ByCreator request")
	todos, err := db.ByCreator(ctx, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(todos) == 0 || err != nil {
		return common.NotFound(common.TodoSrv, p.ByCreator, err, "not found")
	}
	rsp.Data = &todo_proto.ArrData{todos}
	return nil
}

func (p *TodoService) Update(ctx context.Context, req *todo_proto.UpdateRequest, rsp *todo_proto.UpdateResponse) error {
	log.Info("Received Todo.Update request")
	return nil
}
