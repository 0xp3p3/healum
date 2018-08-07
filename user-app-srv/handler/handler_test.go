package handler

import (
	"bytes"
	"encoding/json"
	account_proto "server/account-srv/proto/account"
	behaviour_db "server/behaviour-srv/db"
	behaviour_hdlr "server/behaviour-srv/handler"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	content_db "server/content-srv/db"
	content_hdlr "server/content-srv/handler"
	content_proto "server/content-srv/proto/content"
	kv_proto "server/kv-srv/proto/kv"
	organisation_proto "server/organisation-srv/proto/organisation"
	plan_db "server/plan-srv/db"
	plan_hdlr "server/plan-srv/handler"
	plan_proto "server/plan-srv/proto/plan"
	static_db "server/static-srv/db"
	common_proto "server/static-srv/proto/common"
	static_proto "server/static-srv/proto/static"
	survey_db "server/survey-srv/db"
	survey_hdlr "server/survey-srv/handler"
	survey_proto "server/survey-srv/proto/survey"
	team_proto "server/team-srv/proto/team"
	todo_proto "server/todo-srv/proto/todo"
	track_proto "server/track-srv/proto/track"
	"server/user-app-srv/db"
	userapp_proto "server/user-app-srv/proto/userapp"
	user_db "server/user-srv/db"
	user_hdlr "server/user-srv/handler"
	user_proto "server/user-srv/proto/user"
	"strconv"
	"testing"
	"time"

	"context"

	duration "github.com/ChannelMeter/iso8601duration"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

var cl = client.NewClient(
	client.Transport(nats_transport.NewTransport()),
	client.Broker(nats_broker.NewBroker()),
	client.RequestTimeout(4*time.Second),
	client.Retries(5),
)

var account = &account_proto.Account{
	Email:    "email" + common.Random(4) + "@email.com",
	Password: "pass1",
	Status:   account_proto.AccountStatus_SUSPENDED,
}

var account_phone = &account_proto.Account{
	Phone:    "+8613042431402",
	Passcode: "123456",
}

var user = &user_proto.User{
	OrgId:      "orgid",
	Firstname:  "david",
	Lastname:   "john",
	Tags:       []string{"a", "b", "c"},
	Preference: &user_proto.Preferences{},
	AvatarUrl:  "http://example.com",
	ContactDetails: []*user_proto.ContactDetail{
		{Id: "contact_detail_id"},
	},
	Addresses: []*static_proto.Address{{
		PostalCode: "111000",
	}},
	Tokens: []*user_proto.Token{
		{"11671c2e7da30e3c393813f60b327f9c2e2e08390761aa01e37ba5d3e6a617be", 1, "aaa"}, {"token_b", 2, "bbb"},
	},
}

var content = &content_proto.Content{
	Title:       "title",
	Summary:     []string{"summary1"},
	Description: "description",
	OrgId:       "orgid",
	Image:       "http://example.com",
	CreatedBy:   &user_proto.User{Id: "userid"},
	Url:         "url",
	Author:      "author",
	Timestamp:   12345678,
	Tags:        []*static_proto.ContentCategoryItem{{Id: "111"}},
	Type:        &static_proto.ContentType{},
	Category: &static_proto.ContentCategory{
		Id:       "category111",
		Name:     "activity category",
		NameSlug: "acitivty",
		TrackerMethods: []*static_proto.TrackerMethod{
			{
				Id:       "tracker111",
				NameSlug: "count"},
		},
	},
}

var todo = &todo_proto.Todo{
	Id:      "111",
	Name:    "todo1",
	OrgId:   "orgid",
	Creator: &user_proto.User{Id: "userid"},
}

var plan = &plan_proto.Plan{
	Name:        "plan1",
	Description: "hello world",
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
}

var question = &survey_proto.Question{
	Id:          "q111",
	Type:        survey_proto.QuestionType_FILE,
	Title:       "question",
	Description: "description",
}

var survey = &survey_proto.Survey{
	Id:          "111",
	Title:       "title",
	OrgId:       "orgid",
	Tags:        []string{"tag1", "tag2"},
	Description: "description1",
	Creator:     &user_proto.User{Id: "userid"},
	IsTemplate:  true,
	Status:      survey_proto.SurveyStatus_DRAFT,
	Renders: []survey_proto.RenderTarget{
		survey_proto.RenderTarget_MOBILE, survey_proto.RenderTarget_WEB,
	},
	Setting: &static_proto.Setting{
		Visibility: static_proto.Visibility_PUBLIC,
	},
	Questions: []*survey_proto.Question{question},
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

var goal = &behaviour_proto.Goal{
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
	Users: []*behaviour_proto.TargetedUser{{user1, 90, static_proto.ExpectedProgressType_EXPONENTIAL, ""}},
	CompletionApprovalRequired: false,
	Duration:                   "P1Y2DT3H4M5S",
}

var challenge = &behaviour_proto.Challenge{
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
	Users: []*behaviour_proto.TargetedUser{{user1, 90, static_proto.ExpectedProgressType_EXPONENTIAL, ""}},
	CompletionApprovalRequired: false,
	Duration:                   "P1Y2DT3H4M5S",
}

var habit = &behaviour_proto.Habit{
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
	Users: []*behaviour_proto.TargetedUser{{user1, 90, static_proto.ExpectedProgressType_EXPONENTIAL, ""}},
	CompletionApprovalRequired: false,
	Duration:                   "P1Y2DT3H4M5S",
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

var org1 = &organisation_proto.Organisation{
	Type: organisation_proto.OrganisationType_ROOT,
}

var contentCategoryItem = &static_proto.ContentCategoryItem{
	Id:       "content_category_id",
	Name:     "sample",
	NameSlug: "name_slug",
	IconSlug: "icon_slug",
	Category: contentCategory,
}

var preference = &user_proto.Preferences{
	OrgId: "orgid",
	CurrentMeasurements: []*user_proto.Measurement{
		{
			Id:     "measure_id",
			UserId: "user_id",
			OrgId:  "org_id",
		},
	},
}

var contentParentCategory = &static_proto.ContentParentCategory{
	Id:          "111",
	Name:        "title",
	Summary:     "summary",
	Description: "description",
	IconSlug:    "iconslug",
	OrgId:       "orgid",
	Tags:        []string{"tag1", "tag2"},
}

var contentCategory = &static_proto.ContentCategory{
	Id:          "content_category_id",
	Name:        "title",
	Summary:     "summary",
	Description: "description",
	IconSlug:    "iconslug",
	OrgId:       "orgid",
	Parent:      []*static_proto.ContentParentCategory{contentParentCategory},
	Tags:        []string{"tag1", "tag2"},
}

var user_plan *userapp_proto.UserPlan

func initDb() {
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
	user_db.Init(cl)
	plan_db.Init(cl)
	content_db.Init(cl)
	survey_db.Init(cl)
	behaviour_db.Init(cl)
	static_db.Init(cl)
}

func initHandler() *UserAppService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	return &UserAppService{
		Broker:          nats_brker,
		KvClient:        kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		BehaviourClient: behaviour_proto.NewBehaviourServiceClient("go.micro.srv.behaviour", cl),
		ContentClient:   content_proto.NewContentServiceClient("go.micro.srv.content", cl),
		UserClient:      user_proto.NewUserServiceClient("go.micro.srv.user", cl),
		TrackClient:     track_proto.NewTrackServiceClient("go.micro.srv.track", cl),
		PlanClient:      plan_proto.NewPlanServiceClient("go.micro.srv.plan", cl),
		StaticClient:    static_proto.NewStaticServiceClient("go.micro.srv.static", cl),
		AccountClient:   account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
	}
}

func initBehaviourHandler() *behaviour_hdlr.BehaviourService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	behaviour_hdlr := &behaviour_hdlr.BehaviourService{
		Broker:        nats_brker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		StaticClient:  static_proto.NewStaticServiceClient("go.micro.srv.static", cl),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		TeamClient:    team_proto.NewTeamServiceClient("go.micro.srv.team", cl),
	}
	return behaviour_hdlr
}

func initContentHandler() *content_hdlr.ContentService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	content_hdlr := &content_hdlr.ContentService{
		Broker:        nats_brker,
		StaticClient:  static_proto.NewStaticServiceClient("go.micro.srv.static", cl),
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
	}
	return content_hdlr
}

func initPlanHandler() *plan_hdlr.PlanService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	plan_hdlr := &plan_hdlr.PlanService{
		Broker:        nats_brker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
	}
	return plan_hdlr
}

func initSurveyHandler() *survey_hdlr.SurveyService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	survey_hdlr := &survey_hdlr.SurveyService{
		Broker:        nats_brker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
	}
	return survey_hdlr
}

func createContent(ctx context.Context, hdlr *content_hdlr.ContentService, t *testing.T) *content_proto.Content {
	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return nil
	}

	// create content
	content.OrgId = rsp_org.Data.Organisation.Id
	content.CreatedBy = rsp_org.Data.User
	req_create := &content_proto.CreateContentRequest{
		Content: content,
		OrgId:   content.OrgId,
		TeamId:  content.CreatedBy.Id,
	}
	resp_create := &content_proto.CreateContentResponse{}
	if err := hdlr.CreateContent(ctx, req_create, resp_create); err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.Content
}

func createPlan(ctx context.Context, hdlr *plan_hdlr.PlanService, t *testing.T) *plan_proto.Plan {
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

	req := &plan_proto.CreateRequest{
		Plan:   plan,
		UserId: rsp_org.Data.User.Id,
		OrgId:  rsp_org.Data.Organisation.Id,
	}
	resp := &plan_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req, resp); err != nil {
		t.Error(err)
		return nil
	}

	return resp.Data.Plan
}

func createSurvey(ctx context.Context, hdlr *survey_hdlr.SurveyService, t *testing.T) *survey_proto.Survey {
	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	user1.Id = ""
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return nil
	}
	survey.OrgId = rsp_org.Data.Organisation.Id
	survey.Shares = []*user_proto.User{rsp_org.Data.User}
	survey.Creator = rsp_org.Data.User

	// create survey
	req := &survey_proto.CreateRequest{
		Survey: survey,
		UserId: rsp_org.Data.User.Id,
		OrgId:  rsp_org.Data.Organisation.Id,
	}
	resp := &survey_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req, resp); err != nil {
		t.Error(err)
		return nil
	}
	return resp.Data.Survey
}

func createGoal(ctx context.Context, hdlr *behaviour_hdlr.BehaviourService, t *testing.T) *behaviour_proto.Goal {
	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	user1.Id = ""
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return nil
	}
	goal.CreatedBy = rsp.Data.User
	goal.OrgId = rsp.Data.Organisation.Id
	goal.Users[0].User = rsp.Data.User

	// replace user1
	user1.Id = rsp.Data.User.Id

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

	req := &behaviour_proto.CreateGoalRequest{
		UserId: si.UserId,
		OrgId:  si.OrgId,
		TeamId: si.UserId,
		Goal:   goal,
	}
	resp := &behaviour_proto.CreateGoalResponse{}

	if err := hdlr.CreateGoal(ctx, req, resp); err != nil {
		log.Error("goal is not created:", err)
		return nil
	}

	return resp.Data.Goal
}

func createChallenge(ctx context.Context, hdlr *behaviour_hdlr.BehaviourService, t *testing.T) *behaviour_proto.Challenge {
	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	user1.Id = ""
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return nil
	}
	challenge.CreatedBy = rsp_org.Data.User
	challenge.OrgId = rsp_org.Data.Organisation.Id
	challenge.Users[0].User = rsp_org.Data.User

	// replace user1
	user1.Id = rsp_org.Data.User.Id

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

	req_create := &behaviour_proto.CreateChallengeRequest{
		UserId:    si.UserId,
		OrgId:     si.OrgId,
		TeamId:    si.UserId,
		Challenge: challenge,
	}
	resp_create := &behaviour_proto.CreateChallengeResponse{}
	if err := hdlr.CreateChallenge(ctx, req_create, resp_create); err != nil {
		log.Error("challenge is not created:", err)
		return nil
	}

	return resp_create.Data.Challenge
}

func createHabit(ctx context.Context, hdlr *behaviour_hdlr.BehaviourService, t *testing.T) *behaviour_proto.Habit {
	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	user1.Id = ""
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return nil
	}
	habit.CreatedBy = rsp_org.Data.User
	habit.OrgId = rsp_org.Data.Organisation.Id
	habit.Users[0].User = rsp_org.Data.User

	// replace user1
	user1.Id = rsp_org.Data.User.Id

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

	req_create := &behaviour_proto.CreateHabitRequest{
		UserId: si.UserId,
		OrgId:  si.OrgId,
		TeamId: si.UserId,
		Habit:  habit,
	}
	resp_create := &behaviour_proto.CreateHabitResponse{}
	if err := hdlr.CreateHabit(ctx, req_create, resp_create); err != nil {
		return nil
	}
	return resp_create.Data.Habit
}

func TestCreateBookmark(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create user
	rsp_user, err := hdlr.UserClient.Create(ctx, &user_proto.CreateRequest{User: user1, Account: nil})
	if err != nil {
		t.Error(err)
		return
	}
	// create bookmark
	req_bookmark := &userapp_proto.CreateBookmarkRequest{
		ContentId: content.Id,
		UserId:    rsp_user.Data.User.Id,
	}
	rsp_bookmark := &userapp_proto.CreateBookmarkResponse{}
	if err := hdlr.CreateBookmark(ctx, req_bookmark, rsp_bookmark); err != nil {
		t.Error(err)
		return
	}
}

func TestReadBookmarkContents(t *testing.T) {
	TestCreateBookmark(t)

	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// read bookmark contents
	req_read := &userapp_proto.ReadBookmarkContentRequest{
		UserId: "userid",
	}
	rsp_read := &userapp_proto.ReadBookmarkContentResponse{}
	if err := hdlr.ReadBookmarkContents(ctx, req_read, rsp_read); err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_read.Data.Bookmarks)
}

func TestReadBookmarkContentCategorys(t *testing.T) {
	TestCreateBookmark(t)

	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// read bookmark contents
	req_read := &userapp_proto.ReadBookmarkContentCategorysRequest{
		UserId: "userid",
	}
	rsp_read := &userapp_proto.ReadBookmarkContentCategorysResponse{}
	if err := hdlr.ReadBookmarkContentCategorys(ctx, req_read, rsp_read); err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_read.Data.Categorys)

	req_category := &userapp_proto.ReadBookmarkByCategoryRequest{
		UserId:     "userid",
		CategoryId: rsp_read.Data.Categorys[0].CategoryId,
	}
	rsp_category := &userapp_proto.ReadBookmarkByCategoryResponse{}
	if err := hdlr.ReadBookmarkByCategory(ctx, req_category, rsp_category); err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_read.Data.Categorys)
}

func TestReadBookmarkByCategory(t *testing.T) {
	TestCreateBookmark(t)

	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// read bookmark contents
	req_read := &userapp_proto.ReadBookmarkByCategoryRequest{
		UserId:     "userid",
		CategoryId: "category111",
	}
	rsp_read := &userapp_proto.ReadBookmarkByCategoryResponse{}
	if err := hdlr.ReadBookmarkByCategory(ctx, req_read, rsp_read); err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_read.Data.Bookmarks)
}

func TestDeleteBookmark(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create bookmark
	req_bookmark := &userapp_proto.CreateBookmarkRequest{
		ContentId: content.Id,
		UserId:    "userid",
	}
	rsp_bookmark := &userapp_proto.CreateBookmarkResponse{}
	if err := hdlr.CreateBookmark(ctx, req_bookmark, rsp_bookmark); err != nil {
		t.Error(err)
		return
	}

	// delete bookmark
	req_del := &userapp_proto.DeleteBookmarkRequest{BookmarkId: rsp_bookmark.Data.BookmarkId}
	rsp_del := &userapp_proto.DeleteBookmarkResponse{}
	if err := hdlr.DeleteBookmark(ctx, req_del, rsp_del); err != nil {
		t.Error(err)
		return
	}
}

func TestGetSharedContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	content := createContent(ctx, initContentHandler(), t)
	if content == nil {
		return
	}

	if err := content_db.ShareContent(ctx, []*content_proto.Content{content},
		[]*user_proto.User{content.CreatedBy},
		content.CreatedBy,
		content.OrgId); err != nil {
		t.Error(err)
		return
	}

	// getting shared content
	req_get := &userapp_proto.GetSharedContentRequest{content.CreatedBy.Id}
	rsp_get := &userapp_proto.GetSharedContentResponse{}
	if err := hdlr.GetSharedContent(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.SharedContents) == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Log(rsp_get.Data.SharedContents)
}

func TestGetSharedPlan(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create plan
	plan := createPlan(ctx, initPlanHandler(), t)
	if plan == nil {
		return
	}

	// share plan
	if err := plan_db.SharePlan(ctx, []*plan_proto.Plan{plan}, plan.Shares, plan.Creator, plan.OrgId); err != nil {
		t.Error(err)
		return
	}

	// getting shared plan
	req_get := &userapp_proto.GetSharedPlanRequest{plan.Shares[0].Id}
	rsp_get := &userapp_proto.GetSharedPlanResponse{}
	if err := hdlr.GetSharedPlansForUser(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.SharedPlans) == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Log(rsp_get.Data.SharedPlans)
}

func TestGetSharedSurvey(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create survey
	survey := createSurvey(ctx, initSurveyHandler(), t)
	if plan == nil {
		return
	}

	if err := survey_db.ShareSurvey(ctx, []*survey_proto.Survey{survey}, survey.Shares, survey.Creator, survey.OrgId); err != nil {
		t.Error(err)
		return
	}
	// getting shared survey
	req_get := &userapp_proto.GetSharedSurveyRequest{survey.Shares[0].Id}
	rsp_get := &userapp_proto.GetSharedSurveyResponse{}
	if err := hdlr.GetSharedSurveysForUser(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.SharedSurveys) == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Error(rsp_get.Data.SharedSurveys)
}

func TestGetSharedGoal(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createGoal(ctx, initBehaviourHandler(), t)
	if goal == nil {
		return
	}

	// share goal
	if err := behaviour_db.ShareGoal(ctx, []*behaviour_proto.Goal{goal},
		[]*behaviour_proto.TargetedUser{{User: goal.CreatedBy}},
		goal.CreatedBy, goal.OrgId); err != nil {
		t.Error(err)
		return
	}

	// getting shared goal
	req_get := &userapp_proto.GetSharedGoalRequest{goal.CreatedBy.Id}
	rsp_get := &userapp_proto.GetSharedGoalResponse{}
	if err := hdlr.GetSharedGoalsForUser(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.SharedGoals) == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Log(rsp_get.Data.SharedGoals)
}

func TestGetSharedChallenge(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	challenge := createChallenge(ctx, initBehaviourHandler(), t)
	if challenge == nil {
		return
	}

	// share challenge
	if err := behaviour_db.ShareChallenge(ctx, []*behaviour_proto.Challenge{challenge},
		[]*behaviour_proto.TargetedUser{{User: challenge.CreatedBy}},
		challenge.CreatedBy, challenge.OrgId); err != nil {
		t.Error(err)
		return
	}

	// getting shared challenge
	req_get := &userapp_proto.GetSharedChallengeRequest{challenge.CreatedBy.Id}
	rsp_get := &userapp_proto.GetSharedChallengeResponse{}
	if err := hdlr.GetSharedChallengesForUser(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.SharedChallenges) == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Log(rsp_get.Data.SharedChallenges)
}

func TestGetCurrentJoinedHabit(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	join := signupToHabit(ctx, hdlr, t)
	if join == nil {
		t.Error("SignupToHabit is failed")
		return
	}

	//GetCurrentJoinedHabit
	req_get := &userapp_proto.ListHabitRequest{
		UserId: join.User.Id,
	}
	rsp_get := &userapp_proto.ListHabitResponse{}
	err := hdlr.GetCurrentJoinedHabits(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_get.Data.Response)
}

func TestGetSharedHabit(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create habit
	habit := createHabit(ctx, initBehaviourHandler(), t)
	if habit == nil {
		return
	}

	// share habit
	if err := behaviour_db.ShareHabit(ctx, []*behaviour_proto.Habit{habit},
		[]*behaviour_proto.TargetedUser{{User: habit.CreatedBy}},
		habit.CreatedBy, habit.OrgId); err != nil {
		t.Error(err)
		return
	}

	// getting shared habit
	req_get := &userapp_proto.GetSharedHabitRequest{habit.CreatedBy.Id}
	rsp_get := &userapp_proto.GetSharedHabitResponse{}
	if err := hdlr.GetSharedHabitsForUser(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.SharedHabits) == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Log(rsp_get.Data.SharedHabits)
}

func TestGetAllJoinedHabit(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	join := signupToHabit(ctx, hdlr, t)
	if join == nil {
		t.Error("SignupToHabit is failed")
		return
	}

	//GetCurrentJoinedHabit
	req_get := &userapp_proto.ListHabitRequest{
		UserId: join.User.Id,
	}
	rsp_get := &userapp_proto.ListHabitResponse{}
	err := hdlr.GetAllJoinedHabits(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_get.Data.Response)
}

func TestCalcDuration(t *testing.T) {
	dur, err := duration.FromString("P1Y2DT3H4M5S")
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(dur.ToDuration())
}

func signupToGoal(ctx context.Context, hdlr *UserAppService, t *testing.T) *userapp_proto.JoinGoal {
	// create goal
	goal := createGoal(ctx, initBehaviourHandler(), t)
	if goal == nil {
		return nil
	}
	// share goal
	if err := behaviour_db.ShareGoal(ctx, []*behaviour_proto.Goal{goal},
		[]*behaviour_proto.TargetedUser{{User: goal.CreatedBy}},
		goal.CreatedBy, goal.OrgId); err != nil {
		t.Error(err)
		return nil
	}
	// signup to shared goal
	req_signup := &userapp_proto.SignupToGoalRequest{
		GoalId: goal.Id,
		UserId: goal.CreatedBy.Id,
		OrgId:  goal.OrgId,
	}
	rsp_signup := &userapp_proto.SignupToGoalResponse{}
	if err := hdlr.SignupToGoal(ctx, req_signup, rsp_signup); err != nil {
		t.Error(err)
		return nil
	}
	t.Log(rsp_signup.Data.JoinGoal)
	// add userid for testing
	rsp_signup.Data.JoinGoal.User = &user_proto.User{Id: goal.CreatedBy.Id}
	return rsp_signup.Data.JoinGoal
}

func TestSignupToGoal(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	if join := signupToGoal(ctx, hdlr, t); join == nil {
		t.Error("SignupToGoal is failed")
		return
	}
}

func TestGetCurrentJoinedGoals(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	join := signupToGoal(ctx, hdlr, t)
	if join == nil {
		t.Error("SignupToGoal is failed")
		return
	}

	//GetCurrentJoinedGoal
	req_get := &userapp_proto.ListGoalRequest{
		UserId: join.User.Id,
	}
	rsp_get := &userapp_proto.ListGoalResponse{}
	err := hdlr.GetCurrentJoinedGoals(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_get.Data.Response)
}

func TestGetAllJoinedGoals(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	join := signupToGoal(ctx, hdlr, t)
	if join == nil {
		t.Error("SignupToGoal is failed")
		return
	}

	//GetCurrentJoinedGoal
	req_get := &userapp_proto.ListGoalRequest{
		UserId: join.User.Id,
	}
	rsp_get := &userapp_proto.ListGoalResponse{}
	err := hdlr.GetAllJoinedGoals(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_get.Data.Response)
}

func signupToChallenge(ctx context.Context, hdlr *UserAppService, t *testing.T) *userapp_proto.JoinChallenge {
	// create challenge
	challenge := createChallenge(ctx, initBehaviourHandler(), t)
	if challenge == nil {
		return nil
	}
	// share challenge
	if err := behaviour_db.ShareChallenge(ctx, []*behaviour_proto.Challenge{challenge},
		[]*behaviour_proto.TargetedUser{{User: challenge.CreatedBy}},
		challenge.CreatedBy, challenge.OrgId); err != nil {
		t.Error(err)
		return nil
	}
	// signup to shared challenge
	req_signup := &userapp_proto.SignupToChallengeRequest{
		ChallengeId: challenge.Id,
		UserId:      challenge.CreatedBy.Id,
		OrgId:       goal.OrgId,
	}
	rsp_signup := &userapp_proto.SignupToChallengeResponse{}
	if err := hdlr.SignupToChallenge(ctx, req_signup, rsp_signup); err != nil {
		t.Error(err)
		return nil
	}
	t.Log(rsp_signup.Data.JoinChallenge)
	// add userid for testing
	rsp_signup.Data.JoinChallenge.User = &user_proto.User{Id: challenge.CreatedBy.Id}
	return rsp_signup.Data.JoinChallenge
}

func TestSignupToChallenge(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	if join := signupToChallenge(ctx, hdlr, t); join == nil {
		t.Error("SignupToChallenge is failed")
		return
	}
}

func signupToHabit(ctx context.Context, hdlr *UserAppService, t *testing.T) *userapp_proto.JoinHabit {
	// create habit
	habit := createHabit(ctx, initBehaviourHandler(), t)
	if habit == nil {
		return nil
	}
	// share habit
	if err := behaviour_db.ShareHabit(ctx, []*behaviour_proto.Habit{habit},
		[]*behaviour_proto.TargetedUser{{User: habit.CreatedBy}},
		habit.CreatedBy, habit.OrgId); err != nil {
		t.Error(err)
		return nil
	}
	// signup to shared habit
	req_signup := &userapp_proto.SignupToHabitRequest{
		HabitId: habit.Id,
		UserId:  habit.CreatedBy.Id,
		OrgId:   goal.OrgId,
	}
	rsp_signup := &userapp_proto.SignupToHabitResponse{}
	if err := hdlr.SignupToHabit(ctx, req_signup, rsp_signup); err != nil {
		t.Error(err)
		return nil
	}
	t.Log(rsp_signup.Data.JoinHabit)
	rsp_signup.Data.JoinHabit.User = &user_proto.User{Id: habit.CreatedBy.Id}
	return rsp_signup.Data.JoinHabit
}

func TestSignupToHabit(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	if join := signupToHabit(ctx, hdlr, t); join == nil {
		t.Error("SignupToHabit is failed")
		return
	}
}
func TestGetCurrentJoinedHabits(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	join := signupToHabit(ctx, hdlr, t)
	if join == nil {
		t.Error("SignupToHabit is failed")
		return
	}

	req_list := &userapp_proto.ListHabitRequest{habit.CreatedBy.Id}
	rsp_list := &userapp_proto.ListHabitResponse{}
	err := hdlr.GetCurrentJoinedHabits(ctx, req_list, rsp_list)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp_list.Data.Response)
}

func TestUUID(t *testing.T) {
	t.Log(uuid.NewUUID().String())
}

func TestListMarkers(t *testing.T) {
	initDb()
	// TestSignupToGoal(t)
	// TestSignupToChallenge(t)

	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// list markers request
	req_list := &userapp_proto.ListMarkersRequest{
		UserId: "userid",
	}
	rsp_list := &userapp_proto.ListMarkersResponse{}
	err := hdlr.ListMarkers(ctx, req_list, rsp_list)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp_list.Data.Markers)
}

func TestGetPendingSharedActions(t *testing.T) {
	initDb()
	// TestSignupToGoal(t)
	// TestSignupToChallenge(t)

	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// GetPendingSharedActions
	req_get := &userapp_proto.GetPendingSharedActionsRequest{}
	rsp_get := &userapp_proto.GetPendingSharedActionsResponse{}
	err := hdlr.GetPendingSharedActions(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp_get.Data.Pendings)
}

func TestGetGoalProgress(t *testing.T) {
	initDb()

	TestSignupToGoal(t)
	// TestSignupToChallenge(t)

	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

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

	// GetGoalProgress
	t.Log("user1.id:", user1.Id)
	req_get := &userapp_proto.GetGoalProgressRequest{
		UserId: user1.Id,
	}
	rsp_get := &userapp_proto.GetGoalProgressResponse{}
	if err := hdlr.GetGoalProgress(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp_get.Data.Response)
}

func TestGetDefaultMarkerHistory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	join := signupToGoal(ctx, hdlr, t)
	if join == nil {
		t.Error("SignupToGoal is failed")
		return
	}

	// ListCurrentChallengeWithCount
	req_get := &track_proto.GetDefaultMarkerHistoryRequest{
		UserId: join.User.Id,
		From:   100,
		To:     2000000000,
	}
	rsp_get := &track_proto.GetDefaultMarkerHistoryResponse{}
	err := hdlr.GetDefaultMarkerHistory(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp_get.Data.TrackMarkers)
}

func TestGetCurrentChallengesWithCount(t *testing.T) {
	initDb()

	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	join := signupToChallenge(ctx, hdlr, t)
	if join == nil {
		t.Error("SignupToChallenge is failed")
		return
	}

	// GetCurrentChallengesWithCount
	req_get := &userapp_proto.GetCurrentChallengesWithCountRequest{
		UserId: join.User.Id,
	}
	rsp_get := &userapp_proto.GetCurrentChallengesWithCountResponse{}
	err := hdlr.GetCurrentChallengesWithCount(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp_get.Data.Response)
}

func TestGetCurrentHabitsWithCount(t *testing.T) {
	initDb()

	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	join := signupToHabit(ctx, hdlr, t)
	if join == nil {
		t.Error("SignupToHabit is failed")
		return
	}

	// GetCurrentHabitsWithCount
	req_get := &userapp_proto.GetCurrentHabitsWithCountRequest{
		UserId: join.User.Id,
	}
	rsp_get := &userapp_proto.GetCurrentHabitsWithCountResponse{}
	err := hdlr.GetCurrentHabitsWithCount(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp_get.Data.Response)
}

func TestGetCurrentJoinedChallenge(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	join := signupToChallenge(ctx, hdlr, t)
	if join == nil {
		t.Error("SignupToChallenge is failed")
		return
	}

	//GetCurrentJoinedChallenge
	req_get := &userapp_proto.ListChallengeRequest{
		UserId: join.User.Id,
	}
	rsp_get := &userapp_proto.ListChallengeResponse{}
	err := hdlr.GetCurrentJoinedChallenges(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_get.Data.Response)
}

func TestGetAllJoinedChallenges(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	join := signupToChallenge(ctx, hdlr, t)
	if join == nil {
		t.Error("SignupToChallenge is failed")
		return
	}

	//GetCurrentJoinedChallenge
	req_get := &userapp_proto.ListChallengeRequest{
		UserId: join.User.Id,
	}
	rsp_get := &userapp_proto.ListChallengeResponse{}
	err := hdlr.GetAllJoinedChallenges(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_get.Data.Response)
}

func TestGetContentCategorys(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())

	// create content
	// content_hdlr := &content_hdlr.ContentService{}
	// req_content := &content_proto.CreateContentRequest{
	// 	Content: content,
	// }
	// rsp_content := &content_proto.CreateContentResponse{}
	// if err := content_hdlr.CreateContent(ctx, req_content, rsp_content); err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// get all
	hdlr := initHandler()
	req_all := &content_proto.GetContentCategorysRequest{}
	rsp_all := &content_proto.GetContentCategorysResponse{}
	if err := hdlr.GetContentCategorys(ctx, req_all, rsp_all); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_all.Data.Categorys) == 0 {
		t.Error("Object does not matched")
		return
	}
}

func TestGetContentDetail(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	// content_hdlr := &content_hdlr.ContentService{}
	// create content
	// req_content := &content_proto.CreateContentRequest{
	// 	Content: content,
	// }
	// rsp_content := &content_proto.CreateContentResponse{}
	// if err := content_hdlr.CreateContent(ctx, req_content, rsp_content); err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// create bookmark
	req_bookmark := &userapp_proto.CreateBookmarkRequest{
		ContentId: content.Id,
		UserId:    "userid",
	}
	rsp_bookmark := &userapp_proto.CreateBookmarkResponse{}
	if err := hdlr.CreateBookmark(ctx, req_bookmark, rsp_bookmark); err != nil {
		t.Error(err)
		return
	}

	// get
	req_get := &content_proto.GetContentDetailRequest{
		ContentId: content.Id,
	}
	rsp_get := &content_proto.GetContentDetailResponse{}
	if err := hdlr.GetContentDetail(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if rsp_get.Data.Content == nil {
		t.Error("Object does not matched")
		return
	}
	if !rsp_get.Data.Bookmarked {
		t.Error("Bookmarked does not matched")
		return
	}

	t.Log(rsp_get.Data)
}

func TestGetContentByCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())

	// create content
	// content_hdlr := &content_hdlr.ContentService{}
	// req_content := &content_proto.CreateContentRequest{
	// 	Content: content,
	// }
	// rsp_content := &content_proto.CreateContentResponse{}
	// if err := content_hdlr.CreateContent(ctx, req_content, rsp_content); err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// get
	hdlr := initHandler()
	req_get := &content_proto.GetContentByCategoryRequest{
		CategoryId: content.Category.Id,
	}
	rsp_get := &content_proto.GetContentByCategoryResponse{}
	if err := hdlr.GetContentByCategory(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.Contents) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestGetFiltersForCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())

	// create content category item
	content_hdlr := &content_hdlr.ContentService{}
	req_item := &content_proto.CreateContentCategoryItemRequest{
		ContentCategoryItem: contentCategoryItem,
	}
	rsp_item := &content_proto.CreateContentCategoryItemResponse{}
	if err := content_hdlr.CreateContentCategoryItem(ctx, req_item, rsp_item); err != nil {
		t.Error(err)
		return
	}

	// filter
	hdlr := initHandler()
	req_get := &content_proto.GetFiltersForCategoryRequest{CategoryId: contentCategoryItem.Category.Id}
	rsp_get := &content_proto.GetFiltersForCategoryResponse{}
	if err := hdlr.GetFiltersForCategory(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.ContentCategoryItems) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestFiltersAutocomplete(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())

	// create content category
	content_hdlr := &content_hdlr.ContentService{}
	req_item := &content_proto.CreateContentCategoryItemRequest{
		ContentCategoryItem: contentCategoryItem,
	}
	rsp_item := &content_proto.CreateContentCategoryItemResponse{}
	if err := content_hdlr.CreateContentCategoryItem(ctx, req_item, rsp_item); err != nil {
		t.Error(err)
		return
	}

	// category autocomplete
	hdlr := initHandler()
	req_autocomplete := &content_proto.FiltersAutocompleteRequest{
		CategoryId: contentCategoryItem.Category.Id,
		Name:       "am",
	}
	rsp_autocomplete := &content_proto.FiltersAutocompleteResponse{}
	if err := hdlr.FiltersAutocomplete(ctx, req_autocomplete, rsp_autocomplete); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_autocomplete.Data.ContentCategoryItems) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestFilterContentInParticularCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())

	// create content
	content_hdlr := &content_hdlr.ContentService{}
	req_content := &content_proto.CreateContentRequest{
		Content: content,
	}
	rsp_content := &content_proto.CreateContentResponse{}
	if err := content_hdlr.CreateContent(ctx, req_content, rsp_content); err != nil {
		t.Error(err)
		return
	}

	// filter
	hdlr := initHandler()
	req_filter := &content_proto.FilterContentInParticularCategoryRequest{
		CategoryId:           content.Category.Id,
		ContentCategoryItems: []string{content.Tags[0].Id},
	}
	rsp_filter := &content_proto.FilterContentInParticularCategoryResponse{}
	if err := hdlr.FilterContentInParticularCategory(ctx, req_filter, rsp_filter); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_filter.Data.Contents) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestGetUserPreference(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	//create user here
	user_hdlr := &user_hdlr.UserService{
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
	}
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := user_hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	req_save := &user_proto.SaveUserPreferenceRequest{
		Preference: preference,
		UserId:     rsp_create.Data.User.Id,
		OrgId:      rsp_create.Data.User.OrgId,
	}
	rsp_save := &user_proto.SaveUserPreferenceResponse{}
	err := hdlr.SaveUserPreference(ctx, req_save, rsp_save)
	if err != nil {
		t.Error(err)
		return
	}

	req_get := &user_proto.ReadUserPreferenceRequest{
		UserId: rsp_create.Data.User.Id,
		OrgId:  rsp_create.Data.User.OrgId,
	}

	rsp_get := &user_proto.ReadUserPreferenceResponse{}
	err = hdlr.GetUserPreference(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSaveUserPreference(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_get := &user_proto.SaveUserPreferenceRequest{
		Preference: preference,
	}
	rsp_get := &user_proto.SaveUserPreferenceResponse{}
	err := hdlr.SaveUserPreference(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSaveUserDetails(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	//create user here
	user_hdlr := &user_hdlr.UserService{
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
	}
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := user_hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	req_get := &userapp_proto.SaveUserDetailsRequest{
		UserId:    rsp_create.Data.User.Id,
		OrgId:     rsp_create.Data.User.OrgId,
		Firstname: "first_name",
		Lastname:  "last_name",
		AvatarUrl: "example.jpg",
		Dob:       112233,
	}
	rsp_get := &userapp_proto.SaveUserDetailsResponse{}
	err := hdlr.SaveUserDetails(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetContentRecommendationByUser(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	obj := &content_proto.ContentRecommendation{
		OrgId:   "orgid",
		UserId:  "userid",
		Content: content,
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}
	// publish
	if err := nats_brker.Publish(common.CREATE_CONTENT_RECOMMENDATION, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	// get recommend with user_id
	req_get := &content_proto.GetContentRecommendationByUserRequest{
		UserId: obj.UserId,
	}
	rsp_get := &content_proto.GetContentRecommendationByUserResponse{}
	if err := hdlr.GetContentRecommendationByUser(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.Recommendations) == 0 {
		t.Error("Object count does not matched")
	}
}

func TestGetContentRecommendationByCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	obj := &content_proto.ContentRecommendation{
		OrgId:   "orgid",
		UserId:  "userid",
		Content: content,
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}
	// publish
	if err := nats_brker.Publish(common.CREATE_CONTENT_RECOMMENDATION, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	req_get := &content_proto.GetContentRecommendationByCategoryRequest{
		CategoryId: content.Category.Id,
	}
	rsp_get := &content_proto.GetContentRecommendationByCategoryResponse{}
	if err := hdlr.GetContentRecommendationByCategory(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.Recommendations) == 0 {
		t.Error("Object count does not matched")
	}
}

func TestSaveRateForContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req := &userapp_proto.SaveRateForContentRequest{
		UserId:    user1.Id,
		OrgId:     user1.OrgId,
		ContentId: content.Id,
		Rating:    5,
	}
	rsp := &userapp_proto.SaveRateForContentResponse{}
	if err := hdlr.SaveRateForContent(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	if rsp.Data.ContentRating == nil {
		t.Error("Object does not matched")
	}
}

func TestDislikeForContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req := &userapp_proto.DislikeForContentRequest{
		UserId:    user1.Id,
		OrgId:     user1.OrgId,
		ContentId: content.Id,
	}
	rsp := &userapp_proto.DislikeForContentResponse{}
	if err := hdlr.DislikeForContent(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	if rsp.Data.ContentDislike == nil {
		t.Error("Object does not matched")
	}
}

func TestDislikeForSimilarContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req := &userapp_proto.DislikeForSimilarContentRequest{
		UserId:    user1.Id,
		OrgId:     user1.OrgId,
		ContentId: content.Id,
		Tags:      content.Tags,
	}
	rsp := &userapp_proto.DislikeForSimilarContentResponse{}
	if err := hdlr.DislikeForSimilarContent(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	if rsp.Data.ContentDislikeSimilar == nil {
		t.Error("Object does not matched")
	}
}

func TestSaveUserFeedback(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req := &userapp_proto.SaveUserFeedbackRequest{
		UserId:   user1.Id,
		OrgId:    user1.OrgId,
		Feedback: "very nice!!!",
		Rating:   9,
	}
	rsp := &userapp_proto.SaveUserFeedbackResponse{}
	if err := hdlr.SaveUserFeedback(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	if rsp.Data.Feedback == nil {
		t.Error("Object does not matched")
	}
}

func TestJoinUserPlan(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create plan
	plan_hdlr := &plan_hdlr.PlanService{
		Broker:        hdlr.Broker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
	}

	req_plan := &plan_proto.CreateRequest{
		Plan:   plan,
		UserId: "userid",
		OrgId:  "orgid",
	}
	rsp_plan := &plan_proto.CreateResponse{}
	err := plan_hdlr.Create(ctx, req_plan, rsp_plan)
	if err != nil {
		t.Error(err)
		return
	}

	// join user plan
	req_join := &userapp_proto.JoinUserPlanRequest{UserId: "userid", PlanId: plan.Id}
	rsp_join := &userapp_proto.JoinUserPlanResponse{}
	err = hdlr.JoinUserPlan(ctx, req_join, rsp_join)
	if err != nil {
		t.Error(err)
		return
	}

	if rsp_join.Data.UserPlan == nil {
		t.Error("Object does not matched")
		return
	}

	user_plan = rsp_join.Data.UserPlan
}

func TestCreateUserPlan(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create 20 test contents
	content_hdlr := &content_hdlr.ContentService{
		Broker:        hdlr.Broker,
		StaticClient:  static_proto.NewStaticServiceClient("go.micro.srv.static", cl),
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
	}

	for i := 0; i < 20; i++ {
		content.Id = "id_" + strconv.Itoa(i)
		content.Title = "title_" + strconv.Itoa(i)
		content.Category.Id = "category_" + strconv.Itoa(i/2)
		req_plan := &content_proto.CreateContentRequest{
			Content: content,
			UserId:  "userid",
			OrgId:   "orgid",
		}
		rsp_plan := &content_proto.CreateContentResponse{}
		err := content_hdlr.CreateContent(ctx, req_plan, rsp_plan)
		if err != nil {
			t.Error(err)
			return
		}
	}

	// create goal
	behaviour_hdlr := &behaviour_hdlr.BehaviourService{
		Broker:        hdlr.Broker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
	}

	req_goal := &behaviour_proto.CreateGoalRequest{
		Goal:   goal,
		UserId: "userid",
		OrgId:  "orgid",
	}
	rsp_goal := &behaviour_proto.CreateGoalResponse{}
	err := behaviour_hdlr.CreateGoal(ctx, req_goal, rsp_goal)
	if err != nil {
		t.Error(err)
		return
	}

	// create userplan
	req_create := &userapp_proto.CreateUserPlanRequest{
		UserId:      "userid",
		GoalId:      goal.Id,
		Days:        4,
		ItemsPerDay: 2,
	}
	rsp_create := &userapp_proto.CreateUserPlanResponse{}
	if err := hdlr.CreateUserPlan(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}
}

func TestGetUserPlan(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// read plan
	req_get := &userapp_proto.GetUserPlanRequest{UserId: "userid"}
	rsp_get := &userapp_proto.GetUserPlanResponse{}
	if err := hdlr.GetUserPlan(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if rsp_get.Data.UserPlan == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestUpdateUserPlan(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create userplan
	TestJoinUserPlan(t)
	time.Sleep(2 * time.Second)

	// update plan
	req_update := &userapp_proto.UpdateUserPlanRequest{
		Id:    user_plan.Id,
		OrgId: "orgid123",
		Goals: user_plan.Goals,
		Days:  user_plan.Days,
	}
	rsp_update := &userapp_proto.UpdateUserPlanResponse{}
	if err := hdlr.UpdateUserPlan(ctx, req_update, rsp_update); err != nil {
		t.Error(err)
		return
	}

	// if rsp_update.Data.UserPlan == nil {
	// 	t.Error("Object does not matched")
	// 	return
	// }
}

func TestGetPlanItemsCountByCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create userplan
	// TestJoinUserPlan(t)
	// time.Sleep(4 * time.Second)

	req_get := &userapp_proto.GetPlanItemsCountByCategoryRequest{
		PlanId: plan.Id,
	}
	rsp_get := &userapp_proto.GetPlanItemsCountByCategoryResponse{}
	if err := hdlr.GetPlanItemsCountByCategory(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.ContentCount) == 0 {
		t.Error("Object does not matched")
	}
}

func TestGetPlanItemsCountByDay(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create userplan
	// TestJoinUserPlan(t)
	// time.Sleep(4 * time.Second)

	req_get := &userapp_proto.GetPlanItemsCountByDayRequest{
		PlanId:    plan.Id,
		DayNumber: "2",
	}
	rsp_get := &userapp_proto.GetPlanItemsCountByDayResponse{}
	if err := hdlr.GetPlanItemsCountByDay(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.ContentCount) == 0 {
		t.Error("Object does not matched")
	}
}

func TestGetPlanItemsCountByCategoryAndDay(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create userplan
	// TestJoinUserPlan(t)
	// time.Sleep(4 * time.Second)

	req_get := &userapp_proto.GetPlanItemsCountByCategoryAndDayRequest{
		PlanId: plan.Id,
	}
	rsp_get := &userapp_proto.GetPlanItemsCountByCategoryAndDayResponse{}
	if err := hdlr.GetPlanItemsCountByCategoryAndDay(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.ContentCount) == 0 {
		t.Error("Object does not matched")
	}
}

func TestToDuration(t *testing.T) {
	duration, err := duration.FromString("P1DT")
	if err != nil {
		t.Error(err)
		return
	}
	t.Error(duration.ToDuration().Hours())
}

func TestReadMarkerByNameslug(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_read := &static_proto.ReadByNameslugRequest{NameSlug: "marker-slug"}
	resp := &static_proto.ReadMarkerResponse{}
	err := hdlr.ReadMarkerByNameslug(ctx, req_read, resp)
	if err != nil {
		t.Error(err)
		return
	}
	if resp.Data.Marker == nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestGetGoalDetail(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createGoal(ctx, initBehaviourHandler(), t)
	if goal == nil {
		t.Error("Goal is not created")
		return
	}
	resp := &userapp_proto.ReadGoalResponse{}
	if err := hdlr.GetGoalDetail(ctx, &behaviour_proto.ReadGoalRequest{
		GoalId: goal.Id, OrgId: goal.OrgId,
	}, resp); err != nil {
		t.Error(err)
		return
	}
	if resp.Data.Detail.GoalId != goal.Id {
		t.Error("Object does not matched")
		return
	}
}

func TestGetChallengeDetail(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	challenge := createChallenge(ctx, initBehaviourHandler(), t)
	if challenge == nil {
		t.Error("Challenge is not created")
		return
	}
	resp := &userapp_proto.ReadChallengeResponse{}
	if err := hdlr.GetChallengeDetail(ctx, &behaviour_proto.ReadChallengeRequest{
		ChallengeId: challenge.Id, OrgId: challenge.OrgId,
	}, resp); err != nil {
		t.Error(err)
		return
	}
	if resp.Data.Detail.ChallengeId != challenge.Id {
		t.Error("Object does not matched")
		return
	}
}

func TestGetHabitDetail(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	habit := createHabit(ctx, initBehaviourHandler(), t)
	if habit == nil {
		t.Error("Habit is not created")
		return
	}
	resp := &userapp_proto.ReadHabitResponse{}
	if err := hdlr.GetHabitDetail(ctx, &behaviour_proto.ReadHabitRequest{
		HabitId: habit.Id, OrgId: habit.OrgId,
	}, resp); err != nil {
		t.Error(err)
		return
	}
	if resp.Data.Detail.HabitId != habit.Id {
		t.Error("Object does not matched")
		return
	}
}
