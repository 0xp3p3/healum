package handler

import (
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
	static_db "server/static-srv/db"
	static_hdlr "server/static-srv/handler"
	static_proto "server/static-srv/proto/static"
	"server/track-srv/db"
	track_proto "server/track-srv/proto/track"
	userapp_db "server/user-app-srv/db"
	userapp_hdlr "server/user-app-srv/handler"
	userapp_proto "server/user-app-srv/proto/userapp"
	user_db "server/user-srv/db"
	user_hdlr "server/user-srv/handler"
	user_proto "server/user-srv/proto/user"
	"strings"
	"testing"
	"time"

	google_protobuf1 "github.com/golang/protobuf/ptypes/struct"
	"github.com/labstack/gommon/log"

	"context"

	"github.com/golang/protobuf/jsonpb"
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
	static_db.Init(cl)
	behaviour_db.Init(cl)
	content_db.Init(cl)
	userapp_db.Init(cl)
	user_db.Init(cl)

	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

var user = &user_proto.User{
	Firstname: "david",
	Lastname:  "john",
	OrgId:     "orgid",
	AvatarUrl: "http://example.com",
	Tokens: []*user_proto.Token{
		{"11671c2e7da30e3c393813f60b327f9c2e2e08390761aa01e37ba5d3e6a617be", 1, "aaa"}, {"token_b", 2, "bbb"},
	},
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

var goal = &behaviour_proto.Goal{
	Title:       "g_title",
	OrgId:       "orgid",
	Summary:     "summary",
	Description: "description",
	CreatedBy:   &user_proto.User{Id: "userid"},
	Status:      behaviour_proto.Status_PUBLISHED,
	Category:    &static_proto.BehaviourCategory{Id: "category111"},
	Trackers: []*behaviour_proto.Tracker{
		{
			Marker:    &static_proto.Marker{Id: "marker_id"},
			Frequency: behaviour_proto.Frequency_DAILY,
			Method:    &static_proto.TrackerMethod{Id: "method_id"},
			Until:     "",
		},
	},
	CompletionApprovalRequired: false,
	Users: []*behaviour_proto.TargetedUser{
		{user1, 100, static_proto.ExpectedProgressType_LINEAR, ""},
	},
	Duration: "P1Y2DT3H4M5S",
}

var challenge = &behaviour_proto.Challenge{
	Title:       "c_title",
	OrgId:       "orgid",
	Summary:     "summary",
	Description: "description",
	CreatedBy:   &user_proto.User{Id: "userid"},
	Status:      behaviour_proto.Status_PUBLISHED,
	Category:    &static_proto.BehaviourCategory{Id: "category222"},
	Trackers: []*behaviour_proto.Tracker{
		{
			Marker:    &static_proto.Marker{Id: "marker_id"},
			Frequency: behaviour_proto.Frequency_DAILY,
			Method:    &static_proto.TrackerMethod{Id: "method_id"},
			Until:     "",
		},
	},
	CompletionApprovalRequired: false,
	Users: []*behaviour_proto.TargetedUser{
		{user1, 100, static_proto.ExpectedProgressType_LINEAR, ""},
	},
	Duration: "P1Y2DT3H4M5S",
}

var habit = &behaviour_proto.Habit{
	Title:       "h_title",
	OrgId:       "orgid",
	Summary:     "summary",
	Description: "description",
	CreatedBy:   &user_proto.User{Id: "userid"},
	Status:      behaviour_proto.Status_PUBLISHED,
	Category:    &static_proto.BehaviourCategory{Id: "category333"},
	Trackers: []*behaviour_proto.Tracker{
		{
			Marker:    &static_proto.Marker{Id: "marker_id"},
			Frequency: behaviour_proto.Frequency_DAILY,
			Method:    &static_proto.TrackerMethod{Id: "method_id"},
			Until:     "",
		},
	},
	CompletionApprovalRequired: false,
	Users: []*behaviour_proto.TargetedUser{
		{user1, 100, static_proto.ExpectedProgressType_LINEAR, ""},
	},
	Duration: "P1Y2DT3H4M5S",
}

var content = &content_proto.Content{
	Id:      "111",
	OrgId:   "orgid",
	Title:   "content_title",
	Summary: []string{"summary1", "summary2"},
	Category: &static_proto.ContentCategory{
		Id:       "category111",
		Name:     "activity category",
		NameSlug: "acitivty",
		TrackerMethods: []*static_proto.TrackerMethod{
			{
				Id:       "tracker111",
				NameSlug: "count",
			},
		},
	},
}

var marker = &static_proto.Marker{
	Id:             "111",
	Name:           "title",
	Summary:        "summary",
	Description:    "description",
	OrgId:          "orgid",
	TrackerMethods: []*static_proto.TrackerMethod{trackerMethod_count},
}

var trackerMethod_count = &static_proto.TrackerMethod{
	Id:       "111",
	Name:     "title",
	NameSlug: "count",
	IconSlug: "iconSlug",
}

var trackerMethod_manual = &static_proto.TrackerMethod{
	Id:       "112",
	Name:     "title",
	NameSlug: "manual",
	IconSlug: "iconSlug",
}

var trackerMethod_photo = &static_proto.TrackerMethod{
	Id:       "113",
	Name:     "title",
	NameSlug: "photo",
	IconSlug: "iconSlug",
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

func initUserAppHandler() *userapp_hdlr.UserAppService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	return &userapp_hdlr.UserAppService{
		Broker:          nats_brker,
		KvClient:        kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		BehaviourClient: behaviour_proto.NewBehaviourServiceClient("go.micro.srv.behaviour", cl),
		ContentClient:   content_proto.NewContentServiceClient("go.micro.srv.content", cl),
		UserClient:      user_proto.NewUserServiceClient("go.micro.srv.user", cl),
		TrackClient:     track_proto.NewTrackServiceClient("go.micro.srv.track", cl),
	}
}

func initUserHandler() *user_hdlr.UserService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	return &user_hdlr.UserService{
		Broker:        nats_brker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		TrackClient:   track_proto.NewTrackServiceClient("go.micro.srv.track", cl),
	}
}

func createGoal(ctx context.Context, hdlr *behaviour_hdlr.BehaviourService, t *testing.T) *behaviour_proto.Goal {
	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	user1.Id = ""
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return nil
	}
	goal.CreatedBy = rsp_org.Data.User
	goal.OrgId = rsp_org.Data.Organisation.Id
	goal.Users[0].User = rsp_org.Data.User

	req_create := &behaviour_proto.CreateGoalRequest{
		UserId: rsp_org.Data.User.Id,
		OrgId:  goal.OrgId,
		Goal:   goal,
	}
	resp_create := &behaviour_proto.CreateGoalResponse{}

	if err := hdlr.CreateGoal(ctx, req_create, resp_create); err != nil {
		log.Error("goal is not created:", err)
		return nil
	}

	return resp_create.Data.Goal
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

	req_create := &behaviour_proto.CreateChallengeRequest{
		UserId:    rsp_org.Data.User.Id,
		OrgId:     challenge.OrgId,
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

	req_create := &behaviour_proto.CreateHabitRequest{
		UserId: rsp_org.Data.User.Id,
		OrgId:  habit.OrgId,
		Habit:  habit,
	}
	resp_create := &behaviour_proto.CreateHabitResponse{}
	if err := hdlr.CreateHabit(ctx, req_create, resp_create); err != nil {
		return nil
	}
	return resp_create.Data.Habit
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

	// create category
	rsp_category, err := hdlr.StaticClient.CreateContentCategory(ctx, &static_proto.CreateContentCategoryRequest{ContentCategory: content.Category})
	if err != nil {
		t.Error(err)
		return nil
	}
	content.Category = rsp_category.Data.ContentCategory

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

func createMarker(ctx context.Context, marker *static_proto.Marker, t *testing.T) *static_proto.Marker {
	hdlr := new(static_hdlr.StaticService)

	req_create := &static_proto.CreateMarkerRequest{Marker: marker}
	rsp_create := &static_proto.CreateMarkerResponse{}
	err := hdlr.CreateMarker(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return nil
	}

	return rsp_create.Data.Marker
}

func createTrackerMethod(ctx context.Context, method *static_proto.TrackerMethod, t *testing.T) *static_proto.TrackerMethod {
	hdlr := new(static_hdlr.StaticService)

	req_create := &static_proto.CreateTrackerMethodRequest{TrackerMethod: method}
	rsp_create := &static_proto.CreateTrackerMethodResponse{}
	err := hdlr.CreateTrackerMethod(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return nil
	}

	return rsp_create.Data.TrackerMethod
}

func initHandler() *TrackService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	return &TrackService{
		Broker:          nats_brker,
		KvClient:        kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		BehaviourClient: behaviour_proto.NewBehaviourServiceClient("go.micro.srv.behaviour", cl),
		ContentClient:   content_proto.NewContentServiceClient("go.micro.srv.content", cl),
		StaticClient:    static_proto.NewStaticServiceClient("go.micro.srv.static", cl),
		UserClient:      user_proto.NewUserServiceClient("go.micro.srv.user", cl),
	}
}

func createTrackGoal(ctx context.Context, hdlr *TrackService, t *testing.T) *behaviour_proto.Goal {
	goal := createGoal(ctx, initBehaviourHandler(), t)
	if goal == nil {
		return nil
	}

	if err := behaviour_db.ShareGoal(ctx, []*behaviour_proto.Goal{goal},
		[]*behaviour_proto.TargetedUser{{User: goal.CreatedBy}},
		goal.CreatedBy, goal.OrgId); err != nil {
		t.Error(err)
		return nil
	}

	// signup to shared goal
	userapp_hdlr := initUserAppHandler()
	req_signup := &userapp_proto.SignupToGoalRequest{
		GoalId: goal.Id,
		UserId: goal.CreatedBy.Id,
		OrgId:  goal.OrgId,
	}
	rsp_signup := &userapp_proto.SignupToGoalResponse{}
	if err := userapp_hdlr.SignupToGoal(ctx, req_signup, rsp_signup); err != nil {
		t.Error(err)
		return nil
	}

	req_create := &track_proto.CreateTrackGoalRequest{
		User:   goal.CreatedBy,
		GoalId: goal.Id,
		OrgId:  goal.OrgId,
	}
	resp_create := &track_proto.CreateTrackGoalResponse{}
	if err := hdlr.CreateTrackGoal(ctx, req_create, resp_create); err != nil {
		t.Error(err)
		return nil
	}

	if resp_create.Data.TrackGoal == nil {
		t.Error("Object does not matched")
		return nil
	}

	return goal
}
func TestCreateTrackGoal(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createTrackGoal(ctx, hdlr, t)
	if goal == nil {
		return
	}
}

func TestGetGoalCount(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createTrackGoal(ctx, hdlr, t)
	if goal == nil {
		return
	}

	req_get := &track_proto.GetGoalCountRequest{
		UserId: goal.CreatedBy.Id,
		GoalId: goal.Id,
	}
	resp_get := &track_proto.GetGoalCountResponse{}
	if err := hdlr.GetGoalCount(ctx, req_get, resp_get); err != nil {
		t.Error(err)
		return
	}

	if resp_get.Data.Count == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Log(resp_get.Data.Count)
}

// to do this, flushdb in Redis[10] databse
func TestSetGoalCount(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createTrackGoal(ctx, hdlr, t)
	if goal == nil {
		return
	}

	req_get := &track_proto.GetGoalCountRequest{
		UserId: goal.CreatedBy.Id,
		GoalId: goal.Id,
		From:   1511136024,
		To:     1631136024,
	}
	resp_get := &track_proto.GetGoalCountResponse{}
	err := hdlr.GetGoalCount(ctx, req_get, resp_get)
	if err != nil {
		t.Error(err)
		return
	}

	if resp_get.Data.Count == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Log(resp_get.Data.Count)
}

func TestGetGoalHistory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createTrackGoal(ctx, hdlr, t)
	if goal == nil {
		return
	}
	// getting history
	req_history := &track_proto.GetGoalHistoryRequest{
		GoalId: goal.Id,
		Offset: 0,
		Limit:  10,
		From:   1511136024,
		To:     1631136024,
	}
	resp_history := &track_proto.GetGoalHistoryResponse{}
	if err := hdlr.GetGoalHistory(ctx, req_history, resp_history); err != nil {
		t.Error(err)
		return
	}

	if len(resp_history.Data.TrackGoals) == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Log(resp_history.Data.TrackGoals)
}

func createTrackChallenge(ctx context.Context, hdlr *TrackService, t *testing.T) *behaviour_proto.Challenge {
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
	userapp_hdlr := initUserAppHandler()
	req_signup := &userapp_proto.SignupToChallengeRequest{
		ChallengeId: challenge.Id,
		UserId:      challenge.CreatedBy.Id,
	}
	rsp_signup := &userapp_proto.SignupToChallengeResponse{}
	if err := userapp_hdlr.SignupToChallenge(ctx, req_signup, rsp_signup); err != nil {
		t.Error(err)
		return nil
	}

	req_create := &track_proto.CreateTrackChallengeRequest{
		User:        challenge.CreatedBy,
		ChallengeId: challenge.Id,
		OrgId:       challenge.OrgId,
	}
	resp_create := &track_proto.CreateTrackChallengeResponse{}
	if err := hdlr.CreateTrackChallenge(ctx, req_create, resp_create); err != nil {
		t.Error(err)
		return nil
	}

	if resp_create.Data.TrackChallenge == nil {
		t.Error("Object does not matched")
		return nil
	}
	return challenge
}

func TestCreateTrackChallenge(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	challenge := createTrackChallenge(ctx, hdlr, t)
	if challenge == nil {
		return
	}
}

func TestGetChallengeCount(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	challenge := createTrackChallenge(ctx, hdlr, t)
	if challenge == nil {
		return
	}

	req_get := &track_proto.GetChallengeCountRequest{
		UserId:      challenge.CreatedBy.Id,
		ChallengeId: challenge.Id,
	}
	resp_get := &track_proto.GetChallengeCountResponse{}
	if err := hdlr.GetChallengeCount(ctx, req_get, resp_get); err != nil {
		t.Error(err)
		return
	}

	if resp_get.Data.Count == 0 {
		t.Error("Count does not matched")
		return
	}
}

// to do this, flushdb in Redis[10] databse
func TestSetChallengeCount(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	challenge := createTrackChallenge(ctx, hdlr, t)
	if challenge == nil {
		return
	}

	req_get := &track_proto.GetChallengeCountRequest{
		UserId:      challenge.CreatedBy.Id,
		ChallengeId: challenge.Id,
		From:        1511136024,
		To:          1531136024,
	}
	resp_get := &track_proto.GetChallengeCountResponse{}
	err := hdlr.GetChallengeCount(ctx, req_get, resp_get)
	if err != nil {
		t.Error(err)
		return
	}

	if resp_get.Data.Count == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Log(resp_get.Data.Count)
}

func TestGetChallengeHistory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	challenge := createTrackChallenge(ctx, hdlr, t)
	if challenge == nil {
		return
	}

	// getting history
	req_history := &track_proto.GetChallengeHistoryRequest{
		ChallengeId: challenge.Id,
		Offset:      0,
		Limit:       10,
		From:        1511136024,
		To:          1531136024,
	}
	resp_history := &track_proto.GetChallengeHistoryResponse{}
	if err := hdlr.GetChallengeHistory(ctx, req_history, resp_history); err != nil {
		t.Error(err)
		return
	}

	if len(resp_history.Data.TrackChallenges) == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Log(resp_history.Data.TrackChallenges)
}

func createTrackHabit(ctx context.Context, hdlr *TrackService, t *testing.T) *behaviour_proto.Habit {
	habit := createHabit(ctx, initBehaviourHandler(), t)
	if habit == nil {
		return nil
	}

	// signup to shared habit
	userapp_hdlr := initUserAppHandler()
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
	}
	rsp_signup := &userapp_proto.SignupToHabitResponse{}
	if err := userapp_hdlr.SignupToHabit(ctx, req_signup, rsp_signup); err != nil {
		t.Error(err)
		return nil
	}

	req_create := &track_proto.CreateTrackHabitRequest{
		User:    habit.CreatedBy,
		HabitId: habit.Id,
		OrgId:   habit.OrgId,
	}
	resp_create := &track_proto.CreateTrackHabitResponse{}
	err := hdlr.CreateTrackHabit(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}

	if resp_create.Data.TrackHabit == nil {
		t.Error("Object does not matched")
		return nil
	}

	return habit
}

func TestCreateTrackHabit(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	habit := createTrackHabit(ctx, hdlr, t)
	if habit == nil {
		return
	}
}
func TestGetHabitCount(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	habit := createTrackHabit(ctx, hdlr, t)
	if habit == nil {
		return
	}

	req_get := &track_proto.GetHabitCountRequest{
		UserId:  habit.CreatedBy.Id,
		HabitId: habit.Id,
	}
	resp_get := &track_proto.GetHabitCountResponse{}
	if err := hdlr.GetHabitCount(ctx, req_get, resp_get); err != nil {
		t.Error(err)
		return
	}

	if resp_get.Data.Count == 0 {
		t.Error("Count does not matched")
		return
	}
}

// to do this, flushdb in Redis[10] databse
func TestSetHabitCount(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	habit := createTrackHabit(ctx, hdlr, t)
	if habit == nil {
		return
	}

	req_get := &track_proto.GetHabitCountRequest{
		UserId:  habit.CreatedBy.Id,
		HabitId: habit.Id,
		From:    1511136024,
		To:      1531136024,
	}
	resp_get := &track_proto.GetHabitCountResponse{}
	err := hdlr.GetHabitCount(ctx, req_get, resp_get)
	if err != nil {
		t.Error(err)
		return
	}

	if resp_get.Data.Count == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Log(resp_get.Data.Count)
}

func TestGetHabitHistory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	habit := createTrackHabit(ctx, hdlr, t)
	if habit == nil {
		return
	}

	// getting history
	req_history := &track_proto.GetHabitHistoryRequest{
		HabitId:       habit.Id,
		Offset:        0,
		Limit:         10,
		From:          1511136024,
		To:            1531136024,
		SortParameter: "name",
		SortDirection: "ASC",
	}
	resp_history := &track_proto.GetHabitHistoryResponse{}
	if err := hdlr.GetHabitHistory(ctx, req_history, resp_history); err != nil {
		t.Error(err)
		return
	}

	if len(resp_history.Data.TrackHabits) == 0 {
		t.Error("Count does not matched")
		return
	}
}

func createTrackContent(ctx context.Context, hdlr *TrackService, t *testing.T) *content_proto.Content {
	content := createContent(ctx, initContentHandler(), t)
	if content == nil {
		return nil
	}

	req_create := &track_proto.CreateTrackContentRequest{
		User:      content.CreatedBy,
		ContentId: content.Id,
		OrgId:     content.OrgId,
	}
	resp_create := &track_proto.CreateTrackContentResponse{}
	err := hdlr.CreateTrackContent(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}

	if resp_create.Data.TrackContent == nil {
		t.Error("Object does not matched")
		return nil
	}
	return content
}

func TestCreateTrackContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	content := createTrackContent(ctx, hdlr, t)
	if content == nil {
		return
	}
}
func TestGetContentCount(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	content := createTrackContent(ctx, hdlr, t)
	if content == nil {
		return
	}

	req_get := &track_proto.GetContentCountRequest{
		UserId:    content.CreatedBy.Id,
		ContentId: content.Id,
	}
	resp_get := &track_proto.GetContentCountResponse{}
	if err := hdlr.GetContentCount(ctx, req_get, resp_get); err != nil {
		t.Error(err)
		return
	}

	if resp_get.Data.Count == 0 {
		t.Error("Count does not matched")
		return
	}
}

// to do this, flushdb in Redis[10] databse
func TestSetContentCount(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	content := createTrackContent(ctx, hdlr, t)
	if content == nil {
		return
	}

	req_get := &track_proto.GetContentCountRequest{
		UserId:    content.CreatedBy.Id,
		ContentId: content.Id,
		From:      1511136024,
		To:        1631136024,
	}
	resp_get := &track_proto.GetContentCountResponse{}
	err := hdlr.GetContentCount(ctx, req_get, resp_get)
	if err != nil {
		t.Error(err)
		return
	}

	if resp_get.Data.Count == 0 {
		t.Error("Count does not matched")
		return
	}

	t.Log(resp_get.Data.Count)
}

func TestGetContentHistory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	content := createTrackContent(ctx, hdlr, t)
	if content == nil {
		return
	}

	// getting history
	req_history := &track_proto.GetContentHistoryRequest{
		ContentId:     content.Id,
		Offset:        0,
		Limit:         10,
		From:          1511136024,
		To:            1631136024,
		SortParameter: "name",
		SortDirection: "ASC",
	}
	resp_history := &track_proto.GetContentHistoryResponse{}
	if err := hdlr.GetContentHistory(ctx, req_history, resp_history); err != nil {
		t.Error(err)
		return
	}

	if len(resp_history.Data.TrackContents) == 0 {
		t.Error("Count does not matched")
		return
	}
}

func TestCreateTrackMarkerUsingCountMethod(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return
	}

	method := createTrackerMethod(ctx, trackerMethod_count, t)
	if method == nil {
		return
	}
	marker := createMarker(ctx, marker, t)
	if marker == nil {
		return
	}

	var v google_protobuf1.Value
	raw1 := `3`

	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &v); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}

	req_create := &track_proto.CreateTrackMarkerRequest{
		UserId:        rsp_org.Data.User.Id,
		OrgId:         rsp_org.Data.Organisation.Id,
		MarkerId:      marker.Id,
		Value:         &v,
		TrackerMethod: method,
	}
	rsp_create := &track_proto.CreateTrackMarkerResponse{}

	if err := hdlr.CreateTrackMarker(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	if rsp_create.Data.TrackMarker == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestCreateTrackMarkerUsingManualMethod(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return
	}

	method := createTrackerMethod(ctx, trackerMethod_manual, t)
	if method == nil {
		return
	}
	marker := createMarker(ctx, marker, t)
	if marker == nil {
		return
	}

	var v google_protobuf1.Value
	raw1 := `3`

	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &v); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}

	req_create := &track_proto.CreateTrackMarkerRequest{
		UserId:        rsp_org.Data.User.Id,
		MarkerId:      marker.Id,
		Value:         &v,
		TrackerMethod: method,
	}
	rsp_create := &track_proto.CreateTrackMarkerResponse{}

	if err := hdlr.CreateTrackMarker(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	if rsp_create.Data.TrackMarker == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestGetLastMarker(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return
	}

	method := createTrackerMethod(ctx, trackerMethod_count, t)
	if method == nil {
		return
	}
	marker := createMarker(ctx, marker, t)
	if marker == nil {
		return
	}

	var v google_protobuf1.Value
	raw1 := `3`

	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &v); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}

	req_create := &track_proto.CreateTrackMarkerRequest{
		UserId:        rsp_org.Data.User.Id,
		MarkerId:      marker.Id,
		Value:         &v,
		TrackerMethod: method,
	}
	rsp_create := &track_proto.CreateTrackMarkerResponse{}
	if err := hdlr.CreateTrackMarker(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	if rsp_create.Data.TrackMarker == nil {
		t.Error("Object does not matched")
		return
	}

	req_read := &track_proto.GetLastMarkerRequest{MarkerId: marker.Id, UserId: rsp_org.Data.User.Id}
	rsp_read := &track_proto.GetLastMarkerResponse{}
	err = hdlr.GetLastMarker(ctx, req_read, rsp_read)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp_read.Data)
}

func TestGetLastMarkerWithManual(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return
	}

	method := createTrackerMethod(ctx, trackerMethod_count, t)
	if method == nil {
		return
	}
	marker := createMarker(ctx, marker, t)
	if marker == nil {
		return
	}

	var v google_protobuf1.Value
	raw1 := `"hello world"`

	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &v); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}

	req_create := &track_proto.CreateTrackMarkerRequest{
		UserId:        rsp_org.Data.User.Id,
		MarkerId:      marker.Id,
		Value:         &v,
		TrackerMethod: method,
	}
	rsp_create := &track_proto.CreateTrackMarkerResponse{}

	if err := hdlr.CreateTrackMarker(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	if rsp_create.Data.TrackMarker == nil {
		t.Error("Object does not matched")
		return
	}

	req_read := &track_proto.GetLastMarkerRequest{MarkerId: marker.Id, UserId: rsp_org.Data.User.Id}
	rsp_read := &track_proto.GetLastMarkerResponse{}
	if err := hdlr.GetLastMarker(ctx, req_read, rsp_read); err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp_read.Data)
}

func TestGetLastMarkerWithPhoto(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return
	}

	method := createTrackerMethod(ctx, trackerMethod_count, t)
	if method == nil {
		return
	}
	marker := createMarker(ctx, marker, t)
	if marker == nil {
		return
	}

	var v google_protobuf1.Value
	raw1 := `"http://example.com/image.jpg?30X30"`

	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &v); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}

	req_create := &track_proto.CreateTrackMarkerRequest{
		UserId:        rsp_org.Data.User.Id,
		MarkerId:      marker.Id,
		Value:         &v,
		TrackerMethod: method,
	}
	rsp_create := &track_proto.CreateTrackMarkerResponse{}
	if err := hdlr.CreateTrackMarker(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	if rsp_create.Data.TrackMarker == nil {
		t.Error("Object does not matched")
		return
	}

	req_read := &track_proto.GetLastMarkerRequest{MarkerId: marker.Id, UserId: rsp_org.Data.User.Id}
	rsp_read := &track_proto.GetLastMarkerResponse{}
	err = hdlr.GetLastMarker(ctx, req_read, rsp_read)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp_read.Data)
}

func TestGetMarkerHistory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	marker := createMarker(ctx, marker, t)
	if marker == nil {
		return
	}

	req_read := &track_proto.GetMarkerHistoryRequest{
		SortParameter: "name",
		SortDirection: "ASC",
		MarkerId:      marker.Id,
	}
	rsp_read := &track_proto.GetMarkerHistoryResponse{}
	err := hdlr.GetMarkerHistory(ctx, req_read, rsp_read)
	if err != nil {
		t.Error(err)
		return
	}

	if len(rsp_read.Data.TrackMarkers) == 0 {
		t.Error("Object count does not matched")
		return
	}
	t.Log(rsp_read.Data.TrackMarkers)
}

func TestGetAllMarkerHistory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_read := &track_proto.GetAllMarkerHistoryRequest{
		SortParameter: "name",
		SortDirection: "ASC",
	}
	rsp_read := &track_proto.GetAllMarkerHistoryResponse{}
	err := hdlr.GetAllMarkerHistory(ctx, req_read, rsp_read)
	if err != nil {
		t.Error(err)
		return
	}

	if len(rsp_read.Data.TrackMarkers) == 0 {
		t.Error("Object count does not matched")
		return
	}
	t.Log(rsp_read.Data.TrackMarkers)
}

func TestValueWithJson(t *testing.T) {
	var v google_protobuf1.Value
	raw1 := `"hello world"`

	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &v); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}

	t.Log(v)

	rr := track_proto.CreateTrackMarkerRequest{
		Value: &v,
	}
	t.Log(rr.Value)
	// vv := google_protobuf1.Value_NumberValue{4}
	// v.Value = &google_protobuf1.Value{}

	marshaler := jsonpb.Marshaler{}
	js, err := marshaler.MarshalToString(&v)
	if err != nil {
		return
	}
	t.Log(js)
}
