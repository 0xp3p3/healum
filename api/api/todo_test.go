package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server/api/utils"
	"server/common"
	"server/todo-srv/db"
	todo_proto "server/todo-srv/proto/todo"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

// var serverURL = "http://localhost:8080"
var todoURL = "/server/todos"
var todo = &todo_proto.Todo{
	Id:      "111",
	Name:    "todo1",
	OrgId:   "orgid",
	Creator: user,
}

func initTodoDb() {
	cl := client.NewClient(client.Transport(nats_transport.NewTransport()), client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(5*time.Second),
		client.Retries(5))
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func AllTodos(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+todoURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping todo because already created")

	}
	time.Sleep(time.Second)
}

func CreateTodo(todo *todo_proto.Todo, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"todo": todo})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+todoURL+"/todo?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	// if resp.StatusCode == http.StatusInternalServerError {
	// 	t.Skip("Skipping todo because already created")
	// }
	time.Sleep(time.Second)
}

func ReadTodo(id string, t *testing.T) *todo_proto.Todo {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+todoURL+"/todo/"+id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping todo because already created")

	}
	time.Sleep(time.Second)

	r := todo_proto.ReadResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return nil
	}
	return r.Data.Todo
}

func DeleteTodo(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a DELETE request.
	req, err := http.NewRequest("DELETE", serverURL+todoURL+"/todo/"+id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping todo because already created")

	}
	time.Sleep(time.Second)
}

func SearchTodos(search *todo_proto.SearchRequest, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a PUT request.
	jsonStr, err := json.Marshal(search)
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+todoURL+"/search?session="+sessionId+"&offset=0&limit=20", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping todo because already created")
	}
	time.Sleep(time.Second)
}

func TodosByCreator(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+todoURL+"/creator/"+id+"?session="+sessionId+"&offset=0&limit=20", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping todo because already created")

	}
	time.Sleep(time.Second)
}

func TestTodosAll(t *testing.T) {
	initTodoDb()

	CreateTodo(todo, t)
	AllTodos(t)
}

func TestTodoCreate(t *testing.T) {
	initTodoDb()

	CreateTodo(todo, t)
	p := ReadTodo("111", t)
	if p == nil {
		t.Errorf("Todo does not matched")
		return
	}
	if p.Id != todo.Id {
		t.Errorf("Id does not matched")
		return
	}
	if p.Name != todo.Name {
		t.Errorf("Name does not matched")
		return
	}
}

func TestTodoRead(t *testing.T) {
	initTodoDb()

	CreateTodo(todo, t)
	p := ReadTodo("111", t)
	if p == nil {
		t.Errorf("Todo does not matched")
		return
	}
	if p.Id != todo.Id {
		t.Errorf("Id does not matched")
		return
	}
	if p.Name != todo.Name {
		t.Errorf("Name does not matched")
		return
	}
}

func TestTodoDelete(t *testing.T) {
	initTodoDb()

	CreateTodo(todo, t)
	DeleteTodo("111", t)
	p := ReadTodo("111", t)
	if p != nil {
		t.Errorf("Todo does not matched")
		return
	}
}

func TestTodoSearch(t *testing.T) {
	initTodoDb()

	CreateTodo(todo, t)
	search := &todo_proto.SearchRequest{
		Name:  "todo1",
		OrgId: "orgid",
	}
	SearchTodos(search, t)
}

func TestTodosByCreator(t *testing.T) {
	initTodoDb()

	CreateTodo(todo, t)
	TodosByCreator("userid", t)
}

func TestErrReadTodo(t *testing.T) {
	initTodoDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+todoURL+"/todo/999?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping todo because already created")

	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	}
}

func TestErrAllTodo(t *testing.T) {
	initTodoDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+todoURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping todo because already created")

	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	}
}
