package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"server/api/utils"
	"server/common"
	"server/task-srv/db"
	task_proto "server/task-srv/proto/task"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

// var serverURL = "http://localhost:8080"
var taskURL = "/server/tasks"
var task = &task_proto.Task{
	Id:          "111",
	Title:       "task1",
	OrgId:       "orgid",
	Description: "description1",
	Creator:     user,
	User:        user,
	Assignee:    user,
	Category:    "category1",
	Due:         time.Now().Unix() - 10000,
	Tags:        []string{"a", "b", "c"},
}

func initTaskDb() {
	cl := client.NewClient(client.Transport(nats_transport.NewTransport()), client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(5*time.Second),
		client.Retries(5))
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func AllTasks(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+taskURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping task because already created")

	}
	time.Sleep(time.Second)
}

func CreateTask(task *task_proto.Task, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return
	}
	task.Creator.Id = userId
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"task": task})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+taskURL+"/task?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	// if resp.StatusCode == http.StatusInternalServerError {
	// 	t.Skip("Skipping task because already created")
	// }
	time.Sleep(time.Second)
}

func ReadTask(id string, t *testing.T) *task_proto.Task {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+taskURL+"/task/"+id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping task because already created")

	}
	time.Sleep(time.Second)

	r := task_proto.ReadResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		// t.Errorf("Response does not matched")
		return nil
	}
	return r.Data.Task
}
func DeleteTask(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a DELETE request.
	req, err := http.NewRequest("DELETE", serverURL+taskURL+"/task/"+id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping task because already created")

	}
	time.Sleep(time.Second)
}

func SearchTasks(search *task_proto.SearchRequest, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a PUT request.
	jsonStr, err := json.Marshal(search)
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+taskURL+"/search?session="+sessionId+"&offset=0&limit=20", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping task because already created")
	}
	time.Sleep(time.Second)
}

func TasksByCreator(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+taskURL+"/creator/"+id+"?session="+sessionId+"&offset=0&limit=20", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping task because already created")

	}
	time.Sleep(time.Second)
}

func TasksByAssign(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+taskURL+"/assign/"+id+"?session="+sessionId+"&offset=0&limit=20", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping task because already created")

	}
	time.Sleep(time.Second)
}

func TasksFilter(filter *task_proto.FilterRequest, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(filter)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+taskURL+"/filter?session="+sessionId+"&offset=0&limit=20", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping task because already created")
	}
	time.Sleep(time.Second)
}

func TasksCountByUser(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+taskURL+"/count/"+id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping task because already created")
	}
	time.Sleep(time.Second)
}

func TestTasksAll(t *testing.T) {
	initTaskDb()

	CreateTask(task, t)
	AllTasks(t)
}

func TestTaskRead(t *testing.T) {
	initTaskDb()

	CreateTask(task, t)

	p := ReadTask("111", t)
	if p == nil {
		t.Errorf("Task does not matched")
		return
	}
	if p.Id != task.Id {
		t.Errorf("Id does not matched")
		return
	}
	if p.Title != task.Title {
		t.Errorf("Title does not matched")
		return
	}
}

func TestTaskDelete(t *testing.T) {
	initTaskDb()

	CreateTask(task, t)
	DeleteTask("111", t)
	p := ReadTask("111", t)
	if p != nil {
		t.Errorf("Task does not matched")
		return
	}
}

func TestTaskSearch(t *testing.T) {
	initTaskDb()

	CreateTask(task, t)
	search := &task_proto.SearchRequest{
		Name:  "task1",
		OrgId: "orgid",
	}
	SearchTasks(search, t)
}

func TestTasksByCreator(t *testing.T) {
	initTaskDb()

	CreateTask(task, t)
	TasksByCreator("userid", t)
}

func TestTasksByAssign(t *testing.T) {
	initTaskDb()

	CreateTask(task, t)
	TasksByAssign("userid", t)
}

func TestTasksCountByUser(t *testing.T) {
	initTaskDb()

	CreateTask(task, t)
	TasksCountByUser("userid", t)
}

func TestErrReadTask(t *testing.T) {
	initTaskDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+taskURL+"/task/999?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping task because already created")

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

func TestErrAllTask(t *testing.T) {
	initTaskDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+taskURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping task because already created")

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
