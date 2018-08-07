package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	account_proto "server/account-srv/proto/account"
	"server/common"
	product_proto "server/product-srv/proto/product"
	static_proto "server/static-srv/proto/static"
	"server/team-srv/db"
	team_proto "server/team-srv/proto/team"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

// var serverURL = "http://localhost:8080"
var teamURL = "/server/teams"

var team = &team_proto.Team{
	Id:          "111",
	Name:        "team1",
	Description: "hello world",
	Image:       "iamge001",
	Color:       "red",
	Products: []*product_proto.Product{
		{Id: "p111", Name: "product1"},
	},
	OrgId:     "orgid",
	CreatedBy: &user_proto.User{Id: "111"},
}

var account = &account_proto.Account{
	Email:    "email8@email.com",
	Password: "pass1",
}

var _account = &account_proto.Account{
	Email:    "email9@email.com",
	Password: "pass1",
}

var _user = &user_proto.User{
	Firstname: "David",
	Lastname:  "John",
	OrgId:     "orgid",
}

var employee = &team_proto.Employee{
	OrgId: "orgid",
	Role:  role,
	Teams: []*team_proto.Team{team},
	Profile: &team_proto.EmployeeProfile{
		OrgId: "orgid",
	},
}

var role = &static_proto.Role{
	OrgId: "orgid",
	Name:  "own",
}

func initTeamDb() {
	cl := client.NewClient(client.Transport(nats_transport.NewTransport()), client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(5*time.Second),
		client.Retries(5))
	// ctx := common.NewTestContext(context.TODO())
	// db.DbTeamName = common.TestingName("healum_test")
	// db.DbTeamTable = common.TestingName("team_test")
	// db.DbTeamDriver = "arangodb"
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func createTeam(team *team_proto.Team, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return
	}
	team.CreatedBy.Id = userId
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"team": team})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+teamURL+"/team?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)
}

func createTeamMember(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	_account.Email = "email" + common.Random(4) + "@email.com"
	req_create := &team_proto.CreateTeamMemberRequest{
		User:     _user,
		Account:  _account,
		Employee: employee,
	}

	jsonStr, err := json.Marshal(req_create)
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+teamURL+"/members/member?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func TestAllTeams(t *testing.T) {
	initTeamDb()

	createTeam(team, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+teamURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := team_proto.AllResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Error(r)
		t.Errorf("Object does not matched")
		return
	}
}

func TestReadTeam(t *testing.T) {
	initTeamDb()

	createTeam(team, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+teamURL+"/team/"+team.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := team_proto.ReadResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Team == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.Team.Id != team.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteTeam(t *testing.T) {
	initTeamDb()

	createTeam(team, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+teamURL+"/team/"+team.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+teamURL+"/team/"+team.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := team_proto.ReadResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestFilterTeams(t *testing.T) {
	initTeamDb()

	createTeam(team, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	filter := &team_proto.FilterRequest{
	// Product: []string{"product1"},
	}
	jsonStr, err := json.Marshal(filter)
	if err != nil {
		t.Error(err)
		return
	}
	// Send a POST request.
	req, err := http.NewRequest("POST", serverURL+teamURL+"/filter?session="+sessionId+"&offset=0&limit=10", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := team_proto.FilterResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Data.Teams) == 0 {
		t.Errorf("Object count is not matched")
		return
	}
}

func TestSearchTeams(t *testing.T) {
	initTeamDb()

	createTeam(team, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	search := &team_proto.SearchRequest{
		TeamName: "team1",
	}
	jsonStr, err := json.Marshal(search)
	if err != nil {
		t.Error(err)
		return
	}
	// Send a POST request.
	req, err := http.NewRequest("POST", serverURL+teamURL+"/search?session="+sessionId+"&offset=0&limit=10", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := team_proto.SearchResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Data.Teams) == 0 {
		t.Errorf("Object count is not matched")
	}
}

func TestAllTeamMember(t *testing.T) {
	initTeamDb()

	createTeam(team, t)
	createTeamMember(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	if len(sessionId) == 0 {
		t.Error("sessionId invalid")
		return
	}
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+teamURL+"/members/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := team_proto.AllTeamMemberResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Employees) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func TestReadTeamMember(t *testing.T) {
	initTeamDb()

	createTeam(team, t)
	createTeamMember(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return
	}

	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+teamURL+"/members/member/"+userId+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := team_proto.ReadTeamMemberResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.Employee == nil {
		t.Errorf("Object does not matched")
		return
	}
}

func TestFilterTeamMember(t *testing.T) {
	initTeamDb()

	createTeam(team, t)
	createTeamMember(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	filter := &team_proto.FilterTeamMemberRequest{
		Team: []string{team.Id},
	}
	jsonStr, err := json.Marshal(filter)
	if err != nil {
		t.Error(err)
		return
	}
	// Send a POST request.
	req, err := http.NewRequest("POST", serverURL+teamURL+"/members/filter?session="+sessionId+"&offset=0&limit=10", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := team_proto.FilterTeamMemberResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Data.Employees) == 0 {
		t.Errorf("Object count is not matched")
		return
	}
}
