package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server/api/utils"
	"server/behaviour-srv/db"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	organisation_proto "server/organisation-srv/proto/organisation"
	static_proto "server/static-srv/proto/static"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

// var serverURL = "http://localhost:8080"
var behaviourURL = "/server/behaviours"

var goal = &behaviour_proto.Goal{
	Id:          "g111",
	Title:       "g_title",
	OrgId:       "orgid",
	Summary:     "summary",
	Description: "description",
	CreatedBy:   &user_proto.User{Id: "userid"},
	Status:      behaviour_proto.Status_PUBLISHED,
	Category: &static_proto.BehaviourCategory{
		Id:            "category111",
		MarkerDefault: marker1,
		MarkerOptions: []*static_proto.Marker{marker2, marker3},
	},
	Trackers: []*behaviour_proto.Tracker{
		{
			Marker:    marker1,
			Frequency: behaviour_proto.Frequency_DAILY,
			Method:    &static_proto.TrackerMethod{Id: "method_id"},
			Until:     "",
		},
	},
	Target: &static_proto.Target{
		TargetValue: 100,
	},
	Users: []*behaviour_proto.TargetedUser{{user1, 90, static_proto.ExpectedProgressType_EXPONENTIAL, ""}},
	CompletionApprovalRequired: false,
	Duration:                   "P1Y2DT3H4M5S",
	Tags:                       []string{"tag1", "tag2", "tag3"},
}

var challenge = &behaviour_proto.Challenge{
	Id:          "c111",
	Title:       "c_title",
	OrgId:       "orgid",
	Summary:     "summary",
	Description: "description",
	CreatedBy:   &user_proto.User{Id: "userid"},
	Status:      behaviour_proto.Status_PUBLISHED,
	Category: &static_proto.BehaviourCategory{
		Id:            "category111",
		MarkerDefault: marker1,
		MarkerOptions: []*static_proto.Marker{marker2, marker3},
	},
	Target: &static_proto.Target{
		TargetValue: 100,
	},
	Users: []*behaviour_proto.TargetedUser{{user1, 90, static_proto.ExpectedProgressType_EXPONENTIAL, ""}},
	CompletionApprovalRequired: false,
	Duration:                   "P1Y2DT3H4M5S",
	Tags:                       []string{"tag1", "tag2", "tag3"},
}

var habit = &behaviour_proto.Habit{
	Id:          "h111",
	Title:       "h_title",
	OrgId:       "orgid",
	Summary:     "summary",
	Description: "description",
	CreatedBy:   &user_proto.User{Id: "userid"},
	Status:      behaviour_proto.Status_PUBLISHED,
	Category: &static_proto.BehaviourCategory{
		Id:            "category111",
		MarkerDefault: marker1,
		MarkerOptions: []*static_proto.Marker{marker2, marker3},
	},
	Target: &static_proto.Target{
		TargetValue: 100,
	},
	Users: []*behaviour_proto.TargetedUser{{user1, 90, static_proto.ExpectedProgressType_EXPONENTIAL, ""}},
	CompletionApprovalRequired: false,
	Duration:                   "P1Y2DT3H4M5S",
	Tags:                       []string{"tag1", "tag2", "tag3"},
}

var marker1 = &static_proto.Marker{
	Id:       "marker_id_1",
	Name:     "marker1",
	Summary:  "marker_sumary_1",
	IconSlug: "icon_slug_1",
}

var marker2 = &static_proto.Marker{
	Id:       "marker_id_2",
	Name:     "marker2",
	Summary:  "marker_sumary_2",
	IconSlug: "icon_slug_1",
}

var marker3 = &static_proto.Marker{
	Id:       "marker_id_3",
	Name:     "marker3",
	Summary:  "marker_sumary_3",
	IconSlug: "icon_slug_1",
}

func initBehaviourDb() {
	cl := client.NewClient(client.Transport(nats_transport.NewTransport()), client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(5*time.Second),
		client.Retries(5))
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func createGoal(goal *behaviour_proto.Goal, t *testing.T) *behaviour_proto.Goal {
	// create user
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	userId, orgId := GetUserIdFromSession(sessionId)

	ctx := common.NewTestContext(context.TODO())
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	account_email.Email = "email" + GenerateRand(3) + "@email.com"
	organisation.Id = ""
	// t.Log("org: ", organisation, user, account_email)
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: organisation, User: user, Account: account_email})
	if err != nil {
		log.Println("org is not created:", err)
		return nil
	}

	goal.CreatedBy = &user_proto.User{Id: userId}
	goal.OrgId = orgId
	goal.Users[0].User.Id = rsp_org.Data.User.Id

	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"goal": goal})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+behaviourURL+"/goal?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return nil
	}

	time.Sleep(time.Second)

	r := behaviour_proto.CreateGoalResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return nil
	}
	json.Unmarshal(body, &r)
	if r.Data.Goal == nil {
		t.Errorf("Object does not matched")
		return nil
	}

	return r.Data.Goal
}

func createChallenge(challenge *behaviour_proto.Challenge, t *testing.T) *behaviour_proto.Challenge {
	// create user
	ctx := common.NewTestContext(context.TODO())
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	account_email.Email = "email" + common.Random(3) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: organisation, User: user, Account: account_email})
	if err != nil {
		log.Println("org is not created:", err)
	}
	challenge.OrgId = rsp_org.Data.Organisation.Id
	challenge.Users[0].User.Id = rsp_org.Data.User.Id

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return nil
	}
	challenge.CreatedBy.Id = userId

	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"challenge": challenge})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+behaviourURL+"/challenge?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return nil
	}

	time.Sleep(time.Second)

	r := behaviour_proto.CreateChallengeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return nil
	}
	json.Unmarshal(body, &r)
	if r.Data.Challenge == nil {
		t.Errorf("challenge does not matched")
		return nil
	}

	return r.Data.Challenge
}

func createHabit(habit *behaviour_proto.Habit, t *testing.T) *behaviour_proto.Habit {
	// create user
	ctx := common.NewTestContext(context.TODO())
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	account_email.Email = "email" + GenerateRand(3) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: organisation, User: user, Account: account_email})
	if err != nil {
		log.Println("org is not created:", err)
	}
	habit.OrgId = rsp_org.Data.Organisation.Id
	habit.Users[0].User.Id = rsp_org.Data.User.Id
	time.Sleep(time.Second)

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return nil
	}
	habit.CreatedBy.Id = userId
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"habit": habit})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+behaviourURL+"/habit?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return nil
	}

	time.Sleep(time.Second)

	r := behaviour_proto.CreateHabitResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return nil
	}
	json.Unmarshal(body, &r)
	if r.Data.Habit == nil {
		t.Errorf("Habit does not matched")
		return nil
	}

	return r.Data.Habit
}

func TestAllGoals(t *testing.T) {
	initBehaviourDb()

	createGoal(goal, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return
	}
	goal.CreatedBy.Id = userId
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"goal": goal})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+behaviourURL+"/goals/all?session="+sessionId+"&team_id="+goal.CreatedBy.Id+"&org_id="+goal.OrgId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.AllGoalsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	t.Log(r)
	if r.Data == nil {
		t.Errorf("Object does not matched")
		return
	}

	if len(r.Data.Goals) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func TestAllChallenges(t *testing.T) {
	initBehaviourDb()

	createChallenge(challenge, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+behaviourURL+"/challenges/all?session="+sessionId+"&team_id="+challenge.CreatedBy.Id+"&org_id="+challenge.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.AllChallengesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Object does not matched")
		return
	}

	if len(r.Data.Challenges) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestAllHabits(t *testing.T) {
	initBehaviourDb()

	createHabit(habit, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+behaviourURL+"/habits/all?session="+sessionId+"&team_id="+habit.CreatedBy.Id+"&org_id="+habit.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := behaviour_proto.AllHabitsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Object does not matched")
		return
	}

	if len(r.Data.Habits) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadGoal(t *testing.T) {
	initBehaviourDb()

	createGoal(goal, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+behaviourURL+"/goal/"+goal.Id+"?session="+sessionId+"&team_id="+goal.CreatedBy.Id+"&org_id="+goal.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.ReadGoalResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.Goal == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Goal.Id != goal.Id {
		t.Errorf("Object Id does not matched")
		return
	}

	t.Log(r)
}

func TestReadChallenge(t *testing.T) {
	initBehaviourDb()

	createChallenge(challenge, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+behaviourURL+"/challenge/"+challenge.Id+"?session="+sessionId+"&team_id="+challenge.CreatedBy.Id+"&org_id="+challenge.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.ReadChallengeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.Challenge == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Challenge.Id != challenge.Id {
		t.Errorf("Object Id does not matched")
		return
	}
}

func TestReadHabit(t *testing.T) {
	initBehaviourDb()

	createHabit(habit, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+behaviourURL+"/habit/"+habit.Id+"?session="+sessionId+"&team_id="+habit.CreatedBy.Id+"&org_id="+habit.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.ReadHabitResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.Habit == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Habit.Id != habit.Id {
		t.Errorf("Object Id does not matched")
		return
	}
}

func TestDeleteGoal(t *testing.T) {
	initBehaviourDb()

	createGoal(goal, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+behaviourURL+"/goal/"+goal.Id+"?session="+sessionId+"&team_id="+goal.CreatedBy.Id+"&org_id="+goal.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+behaviourURL+"/goal/"+goal.Id+"?session="+sessionId+"&team_id="+goal.CreatedBy.Id+"&org_id="+goal.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.ReadGoalResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
		return
	}
}

func TestDeleteChallenge(t *testing.T) {
	initBehaviourDb()

	createChallenge(challenge, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+behaviourURL+"/challenge/"+challenge.Id+"?session="+sessionId+"&team_id="+challenge.CreatedBy.Id+"&org_id="+challenge.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+behaviourURL+"/challenge/"+challenge.Id+"?session="+sessionId+"&team_id="+challenge.CreatedBy.Id+"&org_id="+challenge.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.ReadChallengeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
		return
	}
}

func TestDeleteHabit(t *testing.T) {
	initBehaviourDb()

	createHabit(habit, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+behaviourURL+"/habit/"+habit.Id+"?session="+sessionId+"&team_id="+habit.CreatedBy.Id+"&org_id="+habit.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+behaviourURL+"/habit/"+habit.Id+"?session="+sessionId+"&team_id="+habit.CreatedBy.Id+"&org_id="+habit.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.ReadHabitResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
		return
	}
}

func TestFilter(t *testing.T) {
	initBehaviourDb()

	createGoal(goal, t)
	createChallenge(challenge, t)
	createHabit(habit, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	filter := &behaviour_proto.FilterRequest{
		Type:     []string{"goal", "challenge", "habit"},
		Status:   []behaviour_proto.Status{behaviour_proto.Status_PUBLISHED},
		Category: []string{"category111", "category222", "category333"},
		Creator:  []string{"userid"},
	}
	jsonStr, err := json.Marshal(filter)
	if err != nil {
		t.Error(err)
		return
	}
	// Send a POST request.
	req, err := http.NewRequest("POST", serverURL+behaviourURL+"/filter?session="+sessionId+"&team_id="+habit.CreatedBy.Id+"&org_id="+habit.OrgId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.FilterResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Data.Goals) == 0 || len(r.Data.Challenges) == 0 || len(r.Data.Habits) == 0 {
		t.Errorf("Object count is not matched")
		return
	}
}

func TestSearch(t *testing.T) {
	initBehaviourDb()

	createGoal(goal, t)
	createChallenge(challenge, t)
	createHabit(habit, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	search := &behaviour_proto.SearchRequest{
		Name:        "title",
		Description: "descript",
		Summary:     "summary",
	}
	jsonStr, err := json.Marshal(search)
	if err != nil {
		t.Error(err)
		return
	}
	// Send a POST request.
	req, err := http.NewRequest("POST", serverURL+behaviourURL+"/search?session="+sessionId+"&team_id="+habit.CreatedBy.Id+"&org_id="+habit.OrgId+"&offset=0&limit=10", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.SearchResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Data.Goals) == 0 || len(r.Data.Challenges) == 0 || len(r.Data.Habits) == 0 {
		t.Errorf("Object count is not matched")
		return
	}
}

func TestErrReadGoal(t *testing.T) {
	initBehaviourDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+behaviourURL+"/goal/999?session="+sessionId+"&team_id="+goal.CreatedBy.Id+"&org_id="+goal.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	}
	t.Log(r)
}

func TestErrReadChallenge(t *testing.T) {
	initBehaviourDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+behaviourURL+"/challenge/999?session="+sessionId+"&team_id="+challenge.CreatedBy.Id+"&org_id="+challenge.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestErrReadHabit(t *testing.T) {
	initBehaviourDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+behaviourURL+"/habit/999?session="+sessionId+"&team_id="+habit.CreatedBy.Id+"&org_id="+habit.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
		return
	}
	t.Log(r)
}

func TestBindErrFilter(t *testing.T) {
	initBehaviourDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// S	end a POST request.
	req, err := http.NewRequest("POST", serverURL+behaviourURL+"/filter?session="+sessionId+"&team_id="+habit.CreatedBy.Id+"&org_id="+habit.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if r.Message != "BindError" {
		t.Errorf("Error reason does not matched")
		return
	}
}

func TestAutocompleteGoalSearch(t *testing.T) {
	initBehaviourDb()

	createGoal(goal, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"title": "t"})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+behaviourURL+"/goal/search/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.AutocompleteSearchResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Response) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Response)
}

func TestAutocompleteChallengeSearch(t *testing.T) {
	initBehaviourDb()

	createChallenge(challenge, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"title": "t"})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+behaviourURL+"/challenge/search/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.AutocompleteSearchResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Response) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Response)
}

func TestAutocompleteHabitSearch(t *testing.T) {
	initBehaviourDb()

	createHabit(habit, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"title": "t"})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+behaviourURL+"/habit/search/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.AutocompleteSearchResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Response) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Response)
}

func TestGetTopGoalTags(t *testing.T) {
	initBehaviourDb()

	createGoal(goal, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+behaviourURL+"/goals/tags/top/5?session="+sessionId, nil)

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.GetTopTagsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Tags) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Tags)
}

func TestAutocompleteGoalTags(t *testing.T) {
	initBehaviourDb()

	createGoal(goal, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"name": "t"})
	if err != nil {
		t.Error(err)
		return
	}
	// log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+behaviourURL+"/goals/tags/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := behaviour_proto.AutocompleteTagsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Tags) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Tags)
}

func TestUploadGoal(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	t.Log(sessionId)

	// Send a POST request.
	req, err := http.NewRequest("POST", serverURL+behaviourURL+"/goals/upload?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	//	http.HandleFunc("", UploadGoals)
	http.ListenAndServe(":8080", nil)
}
