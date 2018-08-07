package handler

import (
	"context"
	"server/common"
	"server/task-srv/db"
	task_proto "server/task-srv/proto/task"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

var task = &task_proto.Task{
	OrgId:    "orgid",
	Title:    "task1",
	User:     &user_proto.User{Id: "111"},
	Creator:  &user_proto.User{Id: "222"},
	Assignee: &user_proto.User{Id: "333"},
	Status:   task_proto.TaskStatus_COMPLETE,
	Category: "category1",
	Tags:     []string{"a", "b", "c"},
}

func initDb() {
	cl := client.NewClient(
		client.Transport(nats_transport.NewTransport()),
		client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(4*time.Second),
		client.Retries(5),
	)
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func createTask(ctx context.Context, hdlr *TaskService, t *testing.T) *task_proto.Task {

	req := &task_proto.CreateRequest{Task: task}
	resp := &task_proto.CreateResponse{}
	err := hdlr.Create(ctx, req, resp)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp.Data.Task
}

func TestAll(t *testing.T) {
	initDb()
	hdlr := new(TaskService)
	ctx := common.NewTestContext(context.TODO())

	task := createTask(ctx, hdlr, t)
	if task == nil {
		return
	}

	req_all := &task_proto.AllRequest{}
	resp_all := &task_proto.AllResponse{}
	if err := hdlr.All(ctx, req_all, resp_all); err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Tasks) == 0 {
		t.Error("Count does not match")
		return
	}
}

func TestTaskIsCreated(t *testing.T) {
	initDb()
	hdlr := new(TaskService)
	ctx := common.NewTestContext(context.TODO())

	task := createTask(ctx, hdlr, t)
	if task == nil {
		t.Error("Create is failed")
		return
	}
}

func TestTaskRead(t *testing.T) {
	initDb()
	hdlr := new(TaskService)
	ctx := common.NewTestContext(context.TODO())
	task := createTask(ctx, hdlr, t)
	if task == nil {
		return
	}

	req_read := &task_proto.ReadRequest{Id: task.Id}
	resp_read := &task_proto.ReadResponse{}
	if err := hdlr.Read(ctx, req_read, resp_read); err != nil {
		t.Error(err)
		return
	}
	if resp_read.Data.Task == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Task.Id != task.Id {
		t.Error("Id does not match")
		return
	}
}

func TestTaskDelete(t *testing.T) {
	initDb()
	hdlr := new(TaskService)
	ctx := common.NewTestContext(context.TODO())
	task := createTask(ctx, hdlr, t)
	if task == nil {
		return
	}

	req_del := &task_proto.DeleteRequest{Id: task.Id}
	resp_del := &task_proto.DeleteResponse{}
	if err := hdlr.Delete(ctx, req_del, resp_del); err != nil {
		t.Error(err)
		return
	}
}

func TestByCreator(t *testing.T) {
	initDb()
	hdlr := new(TaskService)
	ctx := common.NewTestContext(context.TODO())
	task := createTask(ctx, hdlr, t)
	if task == nil {
		return
	}

	req_creator := &task_proto.ByCreatorRequest{UserId: task.Creator.Id}
	resp_creator := &task_proto.ByCreatorResponse{}
	err := hdlr.ByCreator(ctx, req_creator, resp_creator)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_creator.Data.Tasks) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_creator.Data.Tasks[0].Id != task.Id {
		t.Error("Id does not match")
		return
	}
}

func TestByAssign(t *testing.T) {
	initDb()
	hdlr := new(TaskService)
	ctx := common.NewTestContext(context.TODO())
	task := createTask(ctx, hdlr, t)
	if task == nil {
		return
	}

	req_assign := &task_proto.ByAssignRequest{UserId: task.Assignee.Id}
	resp_assign := &task_proto.ByAssignResponse{}
	err := hdlr.ByAssign(ctx, req_assign, resp_assign)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_assign.Data.Tasks) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_assign.Data.Tasks[0].Id != task.Id {
		t.Error("Id does not match")
		return
	}
}

func TestFilter(t *testing.T) {
	initDb()
	hdlr := new(TaskService)
	ctx := common.NewTestContext(context.TODO())
	task := createTask(ctx, hdlr, t)
	if task == nil {
		return
	}

	req_filter := &task_proto.FilterRequest{
		Status:   []task_proto.TaskStatus{task_proto.TaskStatus_INPROGRESS, task_proto.TaskStatus_COMPLETE},
		Category: []string{"category1", "category2"},
		Priority: []int64{0, 1, 2},
	}
	resp_filter := &task_proto.FilterResponse{}
	if err := hdlr.Filter(ctx, req_filter, resp_filter); err != nil {
		t.Error(err)
		return
	}
	if len(resp_filter.Data.Tasks) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_filter.Data.Tasks[0].Id != task.Id {
		t.Error("Id does not match")
		return
	}
}

func TestCountByUser(t *testing.T) {
	initDb()
	hdlr := new(TaskService)
	ctx := common.NewTestContext(context.TODO())
	task := createTask(ctx, hdlr, t)
	if task == nil {
		return
	}
	req_count := &task_proto.CountByUserRequest{
		UserId: task.User.Id,
	}
	resp_count := &task_proto.CountByUserResponse{}
	if err := hdlr.CountByUser(ctx, req_count, resp_count); err != nil {
		t.Error(err)
		return
	}
	// if resp_count.TaskCount.Expired != 1 {
	// 	t.Error("Expired does not match")
	// 	return
	// }
	// if resp_count.TaskCount.Assigned != 1 {
	// 	t.Error("Expired does not match")
	// 	return
	// }
}

func TestSearch(t *testing.T) {
	initDb()
	hdlr := new(TaskService)
	ctx := common.NewTestContext(context.TODO())
	task := createTask(ctx, hdlr, t)
	if task == nil {
		return
	}

	req_search := &task_proto.SearchRequest{
		Name:   "task1",
		OrgId:  "orgid",
		Offset: 0,
		Limit:  10,
	}
	resp_search := &task_proto.SearchResponse{}
	if err := hdlr.Search(ctx, req_search, resp_search); err != nil {
		t.Error(err)
		return
	}
	if len(resp_search.Data.Tasks) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_search.Data.Tasks[0].Id != task.Id {
		t.Error("Id does not match")
		return
	}
}
