package handler

import (
	"bytes"
	"context"
	"encoding/json"
	account_proto "server/account-srv/proto/account"
	behaviour_db "server/behaviour-srv/db"
	behaviour_hdlr "server/behaviour-srv/handler"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	organisation_proto "server/organisation-srv/proto/organisation"
	"server/plan-srv/db"
	plan_proto "server/plan-srv/proto/plan"
	common_proto "server/static-srv/proto/common"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
	todo_proto "server/todo-srv/proto/todo"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

var cl = client.NewClient(
	client.Transport(nats_transport.NewTransport()),
	client.Broker(nats_broker.NewBroker()),
	client.RequestTimeout(4*time.Second),
	client.Retries(5),
)

func initDb() {
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
	behaviour_db.Init(cl)
}

var todo = &todo_proto.Todo{
	Id:      "111",
	Name:    "todo1",
	OrgId:   "orgid",
	Creator: &user_proto.User{Id: "userid"},
}

var plan = &plan_proto.Plan{
	Id:          "111",
	Name:        "plan1",
	OrgId:       "orgid",
	Description: "hello'world",
	TemplateId:  "template1",
	IsTemplate:  true,
	Status:      plan_proto.StatusEnum_DRAFT,
	Creator:     &user_proto.User{Id: "222"},
	Goals: []*behaviour_proto.Goal{
		{Id: "1", Title: "sample1"},
	},
	Days: map[string]*common_proto.DayItems{
		"1": &common_proto.DayItems{
			[]*common_proto.DayItem{
				{
					Id:   "day_item_001",
					Pre:  todo,
					Post: todo,
				},
			},
		},
		"2": &common_proto.DayItems{
			[]*common_proto.DayItem{
				{
					Id:   "day_item_002",
					Pre:  todo,
					Post: todo,
				},
			},
		},
	},
	Users:  []*user_proto.User{{Id: "userid"}},
	Shares: []*user_proto.User{{Id: "userid"}},
	Setting: &static_proto.Setting{
		Visibility: static_proto.Visibility_PUBLIC,
	},
	Tags: []string{"tag1", "tag2", "tag3", "tag4"},
}

var org1 = &organisation_proto.Organisation{
	Type: organisation_proto.OrganisationType_ROOT,
}

var user1 = &user_proto.User{
	Firstname: "david",
	Lastname:  "john",
	Tokens: []*user_proto.Token{
		{"11671c2e7da30e3c393813f60b327f9c2e2e08390761aa01e37ba5d3e6a617be", 1, "aaa"}, {"token_b", 2, "bbb"},
	},
}

var account1 = &account_proto.Account{
	Email:    "test" + common.Random(4) + "@email.com",
	Password: "pass1",
}

func initHandler() *PlanService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	hdlr := &PlanService{
		Broker:        nats_brker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		TeamClient:    team_proto.NewTeamServiceClient("go.micro.srv.team", cl),
	}
	return hdlr
}

func createPlan(ctx context.Context, hdlr *PlanService, t *testing.T) *plan_proto.Plan {
	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	user1.Id = ""
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return nil
	}
	plan.OrgId = rsp_org.Data.Organisation.Id
	plan.Users = []*user_proto.User{rsp_org.Data.User}
	plan.Creator = rsp_org.Data.User
	plan.Collaborators = []*user_proto.User{rsp_org.Data.User}
	plan.Shares = []*user_proto.User{rsp_org.Data.User}
	// creat goal
	goal := plan.Goals[0]
	goal.OrgId = rsp_org.Data.Organisation.Id
	goal.CreatedBy = rsp_org.Data.User
	goal.Users = []*behaviour_proto.TargetedUser{
		{User: rsp_org.Data.User},
	}
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	behaviour_hdlr := &behaviour_hdlr.BehaviourService{
		Broker:        nats_brker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		TeamClient:    team_proto.NewTeamServiceClient("go.micro.srv.team", cl),
	}

	// login user
	rsp_login, err := hdlr.AccountClient.Login(ctx, &account_proto.LoginRequest{
		Email:    "email8@email.com",
		Password: "pass1",
	})
	if err != nil {
		t.Error("Login is failed")
		return nil
	}
	rsp_kv, err := hdlr.KvClient.ReadSession(ctx, &kv_proto.ReadSessionRequest{common.SESSION_INDEX, rsp_login.Data.Session.Id})
	if err != nil {
		return nil
	}
	si := &account_proto.SessionInfo{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&si); err != nil {
		return nil
	}

	req_goal := &behaviour_proto.CreateGoalRequest{
		UserId: si.UserId,
		OrgId:  si.OrgId,
		Goal:   goal,
	}
	rsp_goal := &behaviour_proto.CreateGoalResponse{}
	if err := behaviour_hdlr.CreateGoal(ctx, req_goal, rsp_goal); err != nil {
		t.Error(err)
		return nil
	}

	req := &plan_proto.CreateRequest{
		Plan:   plan,
		UserId: si.UserId,
		OrgId:  si.OrgId,
	}
	resp := &plan_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req, resp); err != nil {
		t.Error(err)
		return nil
	}

	return resp.Data.Plan
}

func TestAll(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	plan := createPlan(ctx, hdlr, t)
	if plan == nil {
		return
	}

	req_all := &plan_proto.AllRequest{}
	resp_all := &plan_proto.AllResponse{}
	err := hdlr.All(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Plans) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_all.Data.Plans[0].Id != plan.Id {
		t.Error("Id does not match")
		return
	}
}

func TestPlanWithJSON(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	js := `{  
		"plan":{  
		   "id":"1212",
		   "isTemplate":false,
		   "users":[  
	 
		   ],
		   "pic":"http://via.placeholder.com/350x150",
		   "org_id":"vzvG73i41edPNxJ7",
		   "creatorId":"4158873242954991778",
		   "name":"asVDfdb",
		   "description":null,
		   "endTimeUnspecified":false,
		   "start":1522610055,
		   "end":1522782855,
		   "duration":"P1Y2DT3H4M5S",
		   "recurrence":[  
			  {  
				 "RRule":"FREQ=MONTHLY;UNTIL=20180403T191415Z;COUNT=2"
			  }
		   ],
		   "days":{  
			  "1":{  
				 "items":[]
			  },
			  "2":{  
				 "items":[]
			  }
		   },
		   "goals":[  
	 
		   ],
		   "collaborationEnabled":false,
		   "link_sharing_enabled":false,
		   "embeddingEnabled":false,
		   "embeddingEnabled":false,
		   "socials":[  
	 
		   ]
		}
	 }`

	req := plan_proto.CreateRequest{}
	// if err := jsonpb.Unmarshal(strings.NewReader(js), &req); err != nil {
	// 	t.Error(err)
	// 	return
	// }
	decoder := json.NewDecoder(bytes.NewReader([]byte(js)))
	err := decoder.Decode(&req)
	if err != nil {
		t.Error(err)
		return
	}
	req.UserId = "userid"
	req.Plan = plan
	resp := &plan_proto.CreateResponse{}
	res := hdlr.Create(ctx, &req, resp)
	if res != nil {
		t.Error(res)
		return
	}

	// t.Error(resp)

	req_read := &plan_proto.ReadRequest{Id: "111"}
	resp_read := &plan_proto.ReadResponse{}
	err = hdlr.Read(ctx, req_read, resp_read)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", resp_read.Data)
}

func TestPlanRead(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	req := &plan_proto.CreateRequest{Plan: plan, UserId: "userid"}
	resp := &plan_proto.CreateResponse{}
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_read := &plan_proto.ReadRequest{Id: "111"}
	resp_read := &plan_proto.ReadResponse{}
	time.Sleep(3 * time.Second)
	res_read := hdlr.Read(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Plan == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Plan.Id != "111" {
		t.Error("Id does not match")
		return
	}
	if resp_read.Data.Plan.Name != "plan1" {
		t.Error("Name does not match")
		return
	}
	if resp_read.Data.Plan.Description != "hello'world" {
		t.Error("Description does not match")
		return
	}

	t.Logf("%+v", resp_read)
}

func TestPlanDelete(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	req := &plan_proto.CreateRequest{Plan: plan, UserId: "userid"}
	resp := &plan_proto.CreateResponse{}

	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_del := &plan_proto.DeleteRequest{Id: "111"}
	resp_del := &plan_proto.DeleteResponse{}
	time.Sleep(2 * time.Second)
	res_del := hdlr.Delete(ctx, req_del, resp_del)
	if res_del != nil {
		t.Error(res)
	}
}

func TestTemplates(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	req := &plan_proto.CreateRequest{Plan: plan, UserId: "userid"}

	resp := &plan_proto.CreateResponse{}
	time.Sleep(2 * time.Second)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_temp := &plan_proto.TemplatesRequest{SortParameter: "name", SortDirection: "ASC"}
	resp_temp := &plan_proto.TemplatesResponse{}
	time.Sleep(2 * time.Second)
	err := hdlr.Templates(ctx, req_temp, resp_temp)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_temp.Data.Plans) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_temp.Data.Plans[0].Id != "111" {
		t.Error("Id does not match")
		return
	}
	if !resp_temp.Data.Plans[0].IsTemplate {
		t.Error("IsTemplate does not match")
		return
	}
}

func TestDrafts(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	req := &plan_proto.CreateRequest{Plan: plan, UserId: "userid"}
	resp := &plan_proto.CreateResponse{}
	time.Sleep(2 * time.Second)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_drafts := &plan_proto.DraftsRequest{SortParameter: "created", SortDirection: "DESC"}
	resp_drafts := &plan_proto.DraftsResponse{}
	time.Sleep(1 * time.Second)
	err := hdlr.Drafts(ctx, req_drafts, resp_drafts)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_drafts.Data.Plans) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_drafts.Data.Plans[0].Id != "111" {
		t.Error("Id does not match")
		return
	}
	if resp_drafts.Data.Plans[0].Status != plan_proto.StatusEnum_DRAFT {
		t.Error("IsTemplate does not match")
		return
	}
}

func TestByCreator(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	req := &plan_proto.CreateRequest{Plan: plan, UserId: "userid"}
	resp := &plan_proto.CreateResponse{}
	time.Sleep(1 * time.Second)
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_creator := &plan_proto.ByCreatorRequest{
		UserId: "222",
	}
	resp_creator := &plan_proto.ByCreatorResponse{}
	time.Sleep(1 * time.Second)
	err := hdlr.ByCreator(ctx, req_creator, resp_creator)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_creator.Data.Plans) == 0 {
		t.Error("Count does not match")
		return
	}
	// t.Log(res)
	// if resp_creator.Data.Plans[0].Creator.Id != "222" {
	// 	t.Error("Creator does not match")
	// 	return
	// }
}

func TestFilters(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	req := &plan_proto.CreatePlanFilterRequest{
		Filter: &plan_proto.PlanFilter{
			DisplayName: "planfilter1",
			FilterSlug:  "hello world",
		},
	}

	resp := &plan_proto.CreatePlanFilterResponse{}
	time.Sleep(1 * time.Second)
	res := hdlr.CreatePlanFilter(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_filters := &plan_proto.FiltersRequest{SortParameter: "name", SortDirection: "ASC"}
	resp_filters := &plan_proto.FiltersResponse{}
	time.Sleep(1 * time.Second)
	err := hdlr.Filters(ctx, req_filters, resp_filters)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_filters.Data.Filters) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_filters.Data.Filters[0].DisplayName != "planfilter1" {
		t.Error("Name does not match")
		return
	}
}

func TestTimeFilters(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	req := &plan_proto.CreateRequest{Plan: plan, UserId: "userid"}
	resp := &plan_proto.CreateResponse{}
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_timefilters := &plan_proto.TimeFiltersRequest{
		SortParameter: "name",
		SortDirection: "ASC",
		StartDate:     time.Now().Unix() - 10000,
		EndDate:       time.Now().Unix() + 10000,
	}
	resp_timefilters := &plan_proto.TimeFiltersResponse{}
	time.Sleep(1 * time.Second)
	err := hdlr.TimeFilters(ctx, req_timefilters, resp_timefilters)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_timefilters.Data.Plans) == 0 {
		t.Error("Count does not match")
		return
	}
}

func TestGoalFilters(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	req := &plan_proto.CreateRequest{Plan: plan, UserId: "userid"}
	resp := &plan_proto.CreateResponse{}
	res := hdlr.Create(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}

	req_goalfilters := &plan_proto.GoalFiltersRequest{
		Goals: "1,2",
	}
	resp_goalfilters := &plan_proto.GoalFiltersResponse{}
	time.Sleep(1 * time.Second)
	err := hdlr.GoalFilters(ctx, req_goalfilters, resp_goalfilters)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_goalfilters.Data.Plans) == 0 {
		t.Error("Count does not match")
		return
	}
}

func TestSearch(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	//
	plan := createPlan(ctx, hdlr, t)
	if plan == nil {
		return
	}

	req_search := &plan_proto.SearchRequest{
		Name:   plan.Name,
		OrgId:  plan.OrgId,
		Offset: 0,
		Limit:  10,
	}
	resp_search := &plan_proto.SearchResponse{}
	err := hdlr.Search(ctx, req_search, resp_search)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_search.Data.Plans) == 0 {
		t.Error("Count does not match")
		return
	}
	t.Log(resp_search.Data.Plans)
}

func TestSharePlan(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	plan := createPlan(ctx, hdlr, t)
	if plan == nil {
		return
	}

	// login user
	rsp_login, err := hdlr.AccountClient.Login(ctx, &account_proto.LoginRequest{
		Email:    "email8@email.com",
		Password: "pass1",
	})
	if err != nil {
		t.Error("Login is failed")
		return
	}
	rsp_kv, err := hdlr.KvClient.ReadSession(ctx, &kv_proto.ReadSessionRequest{common.SESSION_INDEX, rsp_login.Data.Session.Id})
	if err != nil {
		return
	}
	si := &account_proto.SessionInfo{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&si); err != nil {
		return
	}

	req_share := &plan_proto.SharePlanRequest{
		Plans:  []*plan_proto.Plan{plan},
		Users:  []*user_proto.User{plan.Shares[0]},
		UserId: si.UserId,
		OrgId:  si.OrgId,
	}

	rsp_share := &plan_proto.SharePlanResponse{}

	if err := hdlr.SharePlan(ctx, req_share, rsp_share); err != nil {
		t.Error(err)
		return
	}
}

func TestAutocompleteSearch(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	plan := createPlan(ctx, hdlr, t)
	if plan == nil {
		return
	}

	req := &plan_proto.AutocompleteSearchRequest{"p", "name", "ASC"}
	rsp := &plan_proto.AutocompleteSearchResponse{}
	err := hdlr.AutocompleteSearch(ctx, req, rsp)
	if err != nil {
		t.Error(err)
		return
	}

	if len(rsp.Data.Response) == 0 {
		t.Error("Object count does not matched")
		return
	}

	t.Log(rsp.Data.Response)
}

func TestGetTopTags(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	plan := createPlan(ctx, hdlr, t)
	if plan == nil {
		return
	}

	rsp := &plan_proto.GetTopTagsResponse{}
	if err := hdlr.GetTopTags(ctx, &plan_proto.GetTopTagsRequest{
		OrgId: plan.OrgId,
		N:     5,
	}, rsp); err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp.Data.Tags)
}

func TestAutocompleteTags(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	plan := createPlan(ctx, hdlr, t)
	if plan == nil {
		return
	}

	rsp := &plan_proto.AutocompleteTagsResponse{}
	if err := hdlr.AutocompleteTags(ctx, &plan_proto.AutocompleteTagsRequest{
		OrgId: plan.OrgId,
		Name:  "t",
	}, rsp); err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp.Data.Tags)
}
