package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server/api/utils"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	"server/plan-srv/db"
	plan_proto "server/plan-srv/proto/plan"
	common_proto "server/static-srv/proto/common"
	static_proto "server/static-srv/proto/static"
	user_proto "server/user-srv/proto/user"
	"strconv"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

// var serverURL = "http://localhost:8080"
var planURL = "/server/plans"

var plan = &plan_proto.Plan{
	Id:          "111",
	Name:        "plan1",
	Description: "hello world",
	OrgId:       "orgid",
	TemplateId:  "template1",
	IsTemplate:  true,
	Status:      plan_proto.StatusEnum_DRAFT,
	Creator:     &user_proto.User{Id: "userid"},
	Goals: []*behaviour_proto.Goal{
		{Id: "1"}, {Id: "2"},
	},
	Days: map[string]*common_proto.DayItems{
		"1": &common_proto.DayItems{
			[]*common_proto.DayItem{
				{
					ContentId:        "content_id_0",
					CategoryId:       "category_id_0",
					CategoryIconSlug: "icon_slug",
					CategoryName:     "category_name",
				},
			},
		},
		"2": &common_proto.DayItems{
			[]*common_proto.DayItem{
				{
					ContentId:        "content_id_3",
					CategoryId:       "category_id_0",
					CategoryIconSlug: "icon_slug",
					CategoryName:     "category_name",
				},
				{
					ContentId:        "content_id_1",
					CategoryId:       "category_id_1",
					CategoryIconSlug: "icon_slug_1",
					CategoryName:     "category_name_1",
				},
				{
					ContentId:        "content_id_2",
					CategoryId:       "category_id_2",
					CategoryIconSlug: "icon_slug_2",
					CategoryName:     "category_name_2",
				},
			},
		},
	},
	Users:    []*user_proto.User{{Id: "userid"}},
	Shares:   []*user_proto.User{{Id: "userid"}},
	Duration: "P2DT",
	Tags:     []string{"tag111", "tag222", "tag333"},
	Setting: &static_proto.Setting{
		Visibility: static_proto.Visibility_PUBLIC,
	},
}

func initPlanDb() {
	cl := client.NewClient(
		client.Transport(nats_transport.NewTransport()),
		client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(5*time.Second),
		client.Retries(5))
	// ctx := common.NewTestContext(context.TODO())
	// db.DbPlanName = common.TestingName("healum_test")
	// db.DbPlanTable = common.TestingName("plan_test")
	// db.DbPlanItemTable = common.TestingName("plan_item_test")
	// db.DbPlanTodoTable = common.TestingName("plan_todo_test")
	// db.DbPlanGoalTable = common.TestingName("plan_goal_test")
	// db.DbPlanFilterTable = common.TestingName("plan_filter_test")
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func AllPlans(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+planURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

	}
	time.Sleep(time.Second)

	r := plan_proto.AllResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
	t.Log(r.Data.Plans)
}

func CreatePlan(plan *plan_proto.Plan, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return
	}
	plan.Creator.Id = userId

	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"plan": plan})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+planURL+"/plan?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	// if resp.StatusCode == http.StatusInternalServerError {
	// 	t.Skip("Skipping plan because already created")
	// }
	time.Sleep(2 * time.Second)
}

func ReadPlan(id string, t *testing.T) *plan_proto.Plan {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+planURL+"/plan/"+id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

	}
	time.Sleep(2 * time.Second)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	r := plan_proto.ReadResponse{}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		return nil
	}
	return r.Data.Plan
}

func DeletePlan(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a DELETE request.
	req, err := http.NewRequest("DELETE", serverURL+planURL+"/plan/"+id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

	}
	time.Sleep(time.Second)
}

func SearchPlans(search *plan_proto.SearchRequest, t *testing.T) {
	// Send a PUT request.
	jsonStr, err := json.Marshal(search)
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("POST", serverURL+planURL+"/search?session="+sessionId+"&offset=0&limit=20", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")
	}
	time.Sleep(time.Second)
}

func PlanTemplates(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+planURL+"/templates?session="+sessionId+"&offset=0&limit=20", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

	}
	time.Sleep(time.Second)
}

func PlanDrafts(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+planURL+"/drafts?session="+sessionId+"&offset=0&limit=20", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

	}
	time.Sleep(time.Second)
}

func PlanByCreator(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+planURL+"/creator/"+id+"?session="+sessionId+"&offset=0&limit=20", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

	}
	time.Sleep(time.Second)
}

func PlanFilters(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+planURL+"/filters/all?session="+sessionId+"&offset=0&limit=20", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

	}
	time.Sleep(time.Second)
}

func PlanTimeFilters(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	start := time.Now().Unix() - 10000
	end := time.Now().Unix() + 10000
	req, err := http.NewRequest("GET", serverURL+planURL+"/filter/time?start_date="+strconv.FormatInt(start, 10)+"&end_date="+strconv.FormatInt(end, 10)+"&session="+sessionId+"&offset=0&limit=20", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

	}
	time.Sleep(time.Second)
}

func PlanGoalFilters(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+planURL+"/filter/goal?filter=1,2&session="+sessionId+"&offset=0&limit=20", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

	}
	time.Sleep(time.Second)
}

func TestPlansAll(t *testing.T) {
	initPlanDb()

	CreatePlan(plan, t)
	AllPlans(t)
}

func TestPlanCreate(t *testing.T) {
	initPlanDb()

	CreatePlan(plan, t)
	time.Sleep(2 * time.Second)

	p := ReadPlan("111", t)
	if p == nil {
		t.Errorf("Plan does not matched")
		return
	}
	if p.Id != plan.Id {
		t.Errorf("Id does not matched")
		return
	}
	if p.Name != plan.Name {
		t.Errorf("Name does not matched")
		return
	}
}

func TestPlanDelete(t *testing.T) {
	initPlanDb()

	CreatePlan(plan, t)
	DeletePlan("111", t)
	p := ReadPlan("111", t)
	if p != nil {
		t.Errorf("Plan does not matched")
		return
	}
}

func TestPlanSearch(t *testing.T) {
	initPlanDb()

	CreatePlan(plan, t)
	search := &plan_proto.SearchRequest{
		Name:  "plan1",
		OrgId: "orgid",
	}
	SearchPlans(search, t)
}

func TestPlanTemplates(t *testing.T) {
	initPlanDb()

	CreatePlan(plan, t)
	PlanTemplates(t)
}

func TestPlanDrafts(t *testing.T) {
	initPlanDb()

	CreatePlan(plan, t)
	PlanDrafts(t)
}

func TestPlanByCreator(t *testing.T) {
	initPlanDb()

	CreatePlan(plan, t)
	PlanByCreator("userid", t)
}

func TestPlanFilters(t *testing.T) {
	initPlanDb()

	PlanFilters(t)
}

func TestPlanTimeFilters(t *testing.T) {
	initPlanDb()

	CreatePlan(plan, t)
	PlanTimeFilters(t)
}

func TestPlanGoalFilters(t *testing.T) {
	initPlanDb()

	CreatePlan(plan, t)
	PlanGoalFilters(t)
}

func TestErrReadPlan(t *testing.T) {
	initPlanDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+planURL+"/plan/999?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

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

func TestErrAllPlans(t *testing.T) {
	initPlanDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+planURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

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

func TestAutocompletePlanSearch(t *testing.T) {
	initPlanDb()

	CreatePlan(plan, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"title": "p"})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+planURL+"/plan/search/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := plan_proto.AutocompleteSearchResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Response) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Response)
}

func TestGetTopPlanTags(t *testing.T) {
	initPlanDb()

	CreatePlan(plan, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+planURL+"/tags/top/5?session="+sessionId, nil)

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := plan_proto.GetTopTagsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Tags) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Tags)
}

func TestAutocompletePlanTags(t *testing.T) {
	initPlanDb()

	CreatePlan(plan, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"name": "t"})
	if err != nil {
		t.Error(err)
		return
	}
	// log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+planURL+"/tags/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := plan_proto.AutocompleteTagsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Tags) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Tags)
}
