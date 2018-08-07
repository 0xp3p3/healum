package handler

import (
	"context"
	"server/common"
	"server/todo-srv/db"
	todo_proto "server/todo-srv/proto/todo"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

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

var todo = &todo_proto.Todo{
	Id:      "111",
	OrgId:   "orgid",
	Name:    "todo1",
	Creator: &user_proto.User{Id: "222"},
}

func createTodo(ctx context.Context, hdlr *TodoService, t *testing.T) *todo_proto.Todo {
	req := &todo_proto.CreateRequest{Todo: todo}
	resp := &todo_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req, resp); err != nil {
		t.Error(err)
		return nil
	}

	return resp.Data.Todo
}

func TestAll(t *testing.T) {
	initDb()
	hdlr := new(TodoService)
	ctx := common.NewTestContext(context.TODO())

	todo := createTodo(ctx, hdlr, t)
	if todo == nil {
		return
	}

	req_all := &todo_proto.AllRequest{
		SortParameter: "name",
		SortDirection: "ASC",
	}
	resp_all := &todo_proto.AllResponse{}
	err := hdlr.All(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Todos) == 0 {
		t.Error("Count does not match")
		return
	}
}

func TestTodoIsCreated(t *testing.T) {
	initDb()
	hdlr := new(TodoService)
	ctx := common.NewTestContext(context.TODO())

	todo := createTodo(ctx, hdlr, t)
	if todo == nil {
		t.Error("Create is failed")
		return
	}
}

func TestTodoRead(t *testing.T) {
	initDb()
	hdlr := new(TodoService)
	ctx := common.NewTestContext(context.TODO())

	todo := createTodo(ctx, hdlr, t)
	if todo == nil {
		t.Error("Create is failed")
		return
	}

	req_read := &todo_proto.ReadRequest{Id: todo.Id}
	resp_read := &todo_proto.ReadResponse{}
	if err := hdlr.Read(ctx, req_read, resp_read); err != nil {
		t.Error(err)
		return
	}
	if resp_read.Data.Todo == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Todo.Id != todo.Id {
		t.Error("Id does not match")
		return
	}
}

func TestTodoDelete(t *testing.T) {
	initDb()
	hdlr := new(TodoService)
	ctx := common.NewTestContext(context.TODO())

	todo := createTodo(ctx, hdlr, t)
	if todo == nil {
		t.Error("Create is failed")
		return
	}

	req_del := &todo_proto.DeleteRequest{Id: todo.Id}
	resp_del := &todo_proto.DeleteResponse{}
	if err := hdlr.Delete(ctx, req_del, resp_del); err != nil {
		t.Error(err)
	}
}

func TestByCreator(t *testing.T) {
	initDb()
	hdlr := new(TodoService)
	ctx := common.NewTestContext(context.TODO())
	todo := createTodo(ctx, hdlr, t)
	if todo == nil {
		t.Error("Create is failed")
		return
	}

	req_creator := &todo_proto.ByCreatorRequest{UserId: todo.Creator.Id}
	resp_creator := &todo_proto.ByCreatorResponse{}
	time.Sleep(2 * time.Second)
	err := hdlr.ByCreator(ctx, req_creator, resp_creator)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_creator.Data.Todos) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_creator.Data.Todos[0].Id != todo.Id {
		t.Error("Id does not match")
		return
	}
	// if resp_creator.Data.Todos[0].Creator.Id != todo.Creator.Id {
	// 	t.Error("IsTemplate does not match")
	// 	return
	// }
}

func TestSearch(t *testing.T) {
	initDb()
	hdlr := new(TodoService)
	ctx := common.NewTestContext(context.TODO())
	todo := createTodo(ctx, hdlr, t)
	if todo == nil {
		t.Error("Create is failed")
		return
	}

	req_search := &todo_proto.SearchRequest{
		Name:   todo.Name,
		OrgId:  todo.OrgId,
		Offset: 0,
		Limit:  10,
	}
	resp_search := &todo_proto.SearchResponse{}
	if err := hdlr.Search(ctx, req_search, resp_search); err != nil {
		t.Error(err)
		return
	}
	if len(resp_search.Data.Todos) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_search.Data.Todos[0].Id != todo.Id {
		t.Error("Id does not match")
		return
	}
}
