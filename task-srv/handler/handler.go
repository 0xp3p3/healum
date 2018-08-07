package handler

import (
	"context"
	"server/common"
	"server/task-srv/db"
	task_proto "server/task-srv/proto/task"

	log "github.com/sirupsen/logrus"
)

type TaskService struct{}

func (p *TaskService) All(ctx context.Context, req *task_proto.AllRequest, rsp *task_proto.AllResponse) error {
	log.Info("Received Task.All request")
	tasks, err := db.All(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(tasks) == 0 || err != nil {
		return common.NotFound(common.TaskSrv, p.All, err, "not found")
	}
	rsp.Data = &task_proto.ArrData{tasks}
	return nil
}

func (p *TaskService) Create(ctx context.Context, req *task_proto.CreateRequest, rsp *task_proto.CreateResponse) error {
	log.Info("Received Task.Create request")
	if len(req.Task.Title) == 0 {
		return common.InternalServerError(common.TaskSrv, p.Create, nil, "task title empty")
	}
	if req.Task.Creator == nil {
		return common.InternalServerError(common.TaskSrv, p.Create, nil, "task creator empty")
	}

	err := db.Create(ctx, req.Task)
	if err != nil {
		return common.InternalServerError(common.TaskSrv, p.Create, err, "create error")
	}
	rsp.Data = &task_proto.Data{req.Task}

	return nil
}

func (p *TaskService) Read(ctx context.Context, req *task_proto.ReadRequest, rsp *task_proto.ReadResponse) error {
	log.Info("Received Task.Read request")
	task, err := db.Read(ctx, req.Id, req.OrgId, req.TeamId)
	if task == nil || err != nil {
		return common.NotFound(common.TaskSrv, p.Read, err, "not found")
	}
	rsp.Data = &task_proto.Data{task}
	return nil
}

func (p *TaskService) Delete(ctx context.Context, req *task_proto.DeleteRequest, rsp *task_proto.DeleteResponse) error {
	log.Info("Received Task.Delete request")
	if err := db.Delete(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.TaskSrv, p.Delete, err, "delete error")
	}
	return nil
}

func (p *TaskService) Search(ctx context.Context, req *task_proto.SearchRequest, rsp *task_proto.SearchResponse) error {
	log.Info("Received Task.Search request")
	tasks, err := db.Search(ctx, req.Name, req.OrgId, req.TeamId, req.Limit, req.Offset, req.From, req.To, req.SortParameter, req.SortDirection)
	if len(tasks) == 0 || err != nil {
		return common.NotFound(common.TaskSrv, p.Search, err, "not found")
	}
	rsp.Data = &task_proto.ArrData{tasks}
	return nil
}

func (p *TaskService) ByCreator(ctx context.Context, req *task_proto.ByCreatorRequest, rsp *task_proto.ByCreatorResponse) error {
	log.Info("Received Task.ByCreator request")
	tasks, err := db.ByCreator(ctx, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(tasks) == 0 || err != nil {
		return common.NotFound(common.TaskSrv, p.ByCreator, err, "not found")
	}
	rsp.Data = &task_proto.ArrData{tasks}
	return nil
}

func (p *TaskService) ByAssign(ctx context.Context, req *task_proto.ByAssignRequest, rsp *task_proto.ByAssignResponse) error {
	log.Info("Received Task.ByAssign request")
	tasks, err := db.ByAssign(ctx, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(tasks) == 0 || err != nil {
		return common.NotFound(common.TaskSrv, p.ByAssign, err, "not found")
	}
	rsp.Data = &task_proto.ArrData{tasks}
	return nil
}

func (p *TaskService) Filter(ctx context.Context, req *task_proto.FilterRequest, rsp *task_proto.FilterResponse) error {
	log.Info("Received Task.Filter request")

	tasks, err := db.Filter(ctx, req.Status, req.Category, req.Priority, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(tasks) == 0 || err != nil {
		return common.NotFound(common.TaskSrv, p.Filter, err, "not found")
	}
	rsp.Data = &task_proto.ArrData{tasks}
	return nil
}

func (p *TaskService) CountByUser(ctx context.Context, req *task_proto.CountByUserRequest, rsp *task_proto.CountByUserResponse) error {
	log.Info("Received Task.CountByUser request")
	resp, err := db.CountByUser(ctx, req.UserId, req.OrgId, req.TeamId)
	if resp == nil || err != nil {
		return common.NotFound(common.TaskSrv, p.CountByUser, err, "not found")
	}
	rsp.TaskCount = resp
	return nil
}

func (p *TaskService) Update(ctx context.Context, req *task_proto.UpdateRequest, rsp *task_proto.UpdateResponse) error {
	log.Info("Received Task.Update request")
	return nil
}
