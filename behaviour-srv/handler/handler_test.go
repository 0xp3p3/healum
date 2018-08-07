package handler

import (
	"bytes"
	"context"
	"encoding/json"
	account_proto "server/account-srv/proto/account"
	"server/behaviour-srv/db"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	content_proto "server/content-srv/proto/content"
	kv_proto "server/kv-srv/proto/kv"
	organisation_proto "server/organisation-srv/proto/organisation"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
	todo_proto "server/todo-srv/proto/todo"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
	log "github.com/sirupsen/logrus"
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
}

var goal = &behaviour_proto.Goal{
	Title:       "g_title",
	OrgId:       "orgid",
	Summary:     "summary",
	Description: "description",
	Status:      behaviour_proto.Status_PUBLISHED,
	Category:    &static_proto.BehaviourCategory{Id: "goal_category_id", Name: "sample"},
	CreatedBy:   &user_proto.User{Id: "userid"},
	Target: &static_proto.Target{
		Aim:         &static_proto.BehaviourCategoryAim{Id: "behaviour_category_aim_id", Name: "sample"},
		Marker:      &static_proto.Marker{Id: "marker_id", Name: "sample"},
		TargetValue: 100,
		Unit:        `%`,
	},
	Trackers: []*behaviour_proto.Tracker{
		{
			Marker:    &static_proto.Marker{Id: "marker_id", Name: "sample"},
			Frequency: behaviour_proto.Frequency_DAILY,
			Method:    &static_proto.TrackerMethod{Id: "method_id", Name: "sample"},
			Until:     "",
		},
	},
	Triggers:                   []*static_proto.ModuleTrigger{{Id: "triger_id", Name: "sample"}},
	CompletionApprovalRequired: false,
	Challenges:                 []*behaviour_proto.Challenge{challenge},
	Habits:                     []*behaviour_proto.Habit{habit},
	Users: []*behaviour_proto.TargetedUser{
		{user1, 100, static_proto.ExpectedProgressType_LINEAR, ""},
	},
	Tags:     []string{"tag1", "tag2", "tag3", "tag4"},
	Setbacks: []*static_proto.Setback{{Id: "setback_id", Name: "sample"}},
	Todos:    &todo_proto.Todo{Id: "todo_id", Name: "sample"},
}

var challenge = &behaviour_proto.Challenge{
	Title:       "c_title",
	OrgId:       "orgid",
	Summary:     "summary",
	Description: "description",
	Status:      behaviour_proto.Status_PUBLISHED,
	Category:    &static_proto.BehaviourCategory{Id: "challenge_category_id", Name: "sample"},
	CreatedBy:   &user_proto.User{Id: "userid"},
	Target: &static_proto.Target{
		Aim:         &static_proto.BehaviourCategoryAim{Id: "behaviour_category_aim_id", Name: "sample"},
		Marker:      &static_proto.Marker{Id: "marker_id", Name: "sample"},
		TargetValue: 100,
		Unit:        `%`,
	},
	Trackers: []*behaviour_proto.Tracker{
		{
			Marker:    &static_proto.Marker{Id: "marker_id", Name: "sample"},
			Frequency: behaviour_proto.Frequency_DAILY,
			Method:    &static_proto.TrackerMethod{Id: "method_id", Name: "sample"},
			Until:     "",
		},
	},
	Triggers:                   []*static_proto.ModuleTrigger{{Id: "triger_id", Name: "sample"}},
	Habits:                     []*behaviour_proto.Habit{habit},
	Articles:                   []*content_proto.Article{},
	CompletionApprovalRequired: false,
	Users: []*behaviour_proto.TargetedUser{
		{user1, 100, static_proto.ExpectedProgressType_LINEAR, ""},
	},
	Tags:     []string{"tag1", "tag2", "tag3", "tag4"},
	Setbacks: []*static_proto.Setback{{Id: "setback_id", Name: "sample"}},
	Todos:    &todo_proto.Todo{Id: "todo_id", Name: "sample"},
}

var habit = &behaviour_proto.Habit{
	Title:       "h_title",
	OrgId:       "orgid",
	Summary:     "summary",
	Description: "description",
	CreatedBy:   &user_proto.User{Id: "userid"},
	Status:      behaviour_proto.Status_PUBLISHED,
	Category:    &static_proto.BehaviourCategory{Id: "habit_category_id", Name: "sample"},
	Target: &static_proto.Target{
		Aim:         &static_proto.BehaviourCategoryAim{Id: "behaviour_category_aim_id", Name: "sample"},
		Marker:      &static_proto.Marker{Id: "marker_id", Name: "sample"},
		TargetValue: 100,
		Unit:        `%`,
	},
	Triggers: []*static_proto.ModuleTrigger{{Id: "triger_id", Name: "sample"}},
	Trackers: []*behaviour_proto.Tracker{
		{
			Marker:    &static_proto.Marker{Id: "marker_id", Name: "sample"},
			Frequency: behaviour_proto.Frequency_DAILY,
			Method:    &static_proto.TrackerMethod{Id: "method_id", Name: "sample"},
			Until:     "",
		},
	},
	CompletionApprovalRequired: false,
	Users: []*behaviour_proto.TargetedUser{
		{user1, 100, static_proto.ExpectedProgressType_LINEAR, ""},
	},
	Tags:     []string{"tag1", "tag2", "tag3", "tag4"},
	Setbacks: []*static_proto.Setback{{Id: "setback_id", Name: "sample"}},
	Todos:    &todo_proto.Todo{Id: "todo_id", Name: "sample"},
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

func initHandler() *BehaviourService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	hdlr := &BehaviourService{
		Broker:        nats_brker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		StaticClient:  static_proto.NewStaticServiceClient("go.micro.srv.static", cl),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		TeamClient:    team_proto.NewTeamServiceClient("go.micro.srv.team", cl),
	}
	return hdlr
}

func createGoal(ctx context.Context, hdlr *BehaviourService, t *testing.T) *behaviour_proto.Goal {
	// create behaviour category
	staticClient := static_proto.NewStaticServiceClient("go.micro.srv.static", cl)
	rsp_category, err := staticClient.CreateBehaviourCategory(ctx, &static_proto.CreateBehaviourCategoryRequest{
		Category: goal.Category,
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	goal.Category = rsp_category.Data.Category
	// create aim
	rsp_aim, err := staticClient.CreateBehaviourCategoryAim(ctx, &static_proto.CreateBehaviourCategoryAimRequest{
		BehaviourCategoryAim: goal.Target.Aim,
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	goal.Target.Aim = rsp_aim.Data.BehaviourCategoryAim
	// create marker
	rsp_marker, err := staticClient.CreateMarker(ctx, &static_proto.CreateMarkerRequest{
		Marker: goal.Target.Marker,
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	goal.Target.Marker = rsp_marker.Data.Marker
	// create method
	rsp_method, err := staticClient.CreateTrackerMethod(ctx, &static_proto.CreateTrackerMethodRequest{
		TrackerMethod: goal.Trackers[0].Method,
	})
	goal.Trackers[0].Marker = rsp_marker.Data.Marker
	goal.Trackers[0].Method = rsp_method.Data.TrackerMethod
	// create challenge
	challenge := createChallenge(ctx, hdlr, t)
	if challenge == nil {
		t.Error("challenge is error")
		return nil
	}
	goal.Challenges[0] = challenge
	// create habit
	habit := createHabit(ctx, hdlr, t)
	if habit == nil {
		t.Error("habit is error")
		return nil
	}
	goal.Habits[0] = habit
	// create setbacks
	rsp_setback, err := staticClient.CreateSetback(ctx, &static_proto.CreateSetbackRequest{Setback: goal.Setbacks[0]})
	if err != nil {
		t.Error(err)
		return nil
	}
	goal.Setbacks[0] = rsp_setback.Data.Setback

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

	goal.CreatedBy = &user_proto.User{Id: si.UserId}
	goal.OrgId = si.OrgId
	goal.Users[0].User = &user_proto.User{Id: si.UserId}
	req_create := &behaviour_proto.CreateGoalRequest{
		UserId: si.UserId,
		OrgId:  si.OrgId,
		Goal:   goal,
	}
	t.Log(si.UserId, si.OrgId)
	resp_create := &behaviour_proto.CreateGoalResponse{}

	if err := hdlr.CreateGoal(ctx, req_create, resp_create); err != nil {
		log.Error("goal is not created:", err)
		return nil
	}
	return resp_create.Data.Goal
}

func createChallenge(ctx context.Context, hdlr *BehaviourService, t *testing.T) *behaviour_proto.Challenge {
	// create behaviou category
	staticClient := static_proto.NewStaticServiceClient("go.micro.srv.static", cl)
	rsp_category, err := staticClient.CreateBehaviourCategory(ctx, &static_proto.CreateBehaviourCategoryRequest{
		Category: challenge.Category,
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	challenge.Category = rsp_category.Data.Category
	// create aim
	rsp_aim, err := staticClient.CreateBehaviourCategoryAim(ctx, &static_proto.CreateBehaviourCategoryAimRequest{
		BehaviourCategoryAim: challenge.Target.Aim,
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	challenge.Target.Aim = rsp_aim.Data.BehaviourCategoryAim
	// create marker
	rsp_marker, err := staticClient.CreateMarker(ctx, &static_proto.CreateMarkerRequest{
		Marker: challenge.Target.Marker,
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	challenge.Target.Marker = rsp_marker.Data.Marker
	// create method
	rsp_method, err := staticClient.CreateTrackerMethod(ctx, &static_proto.CreateTrackerMethodRequest{
		TrackerMethod: challenge.Trackers[0].Method,
	})
	challenge.Trackers[0].Marker = rsp_marker.Data.Marker
	challenge.Trackers[0].Method = rsp_method.Data.TrackerMethod

	// create habit
	habit := createHabit(ctx, hdlr, t)
	if habit == nil {
		return nil
	}
	challenge.Habits[0] = habit
	// create setbacks
	rsp_setback, err := staticClient.CreateSetback(ctx, &static_proto.CreateSetbackRequest{Setback: challenge.Setbacks[0]})
	if err != nil {
		t.Error(err)
		return nil
	}
	challenge.Setbacks[0] = rsp_setback.Data.Setback

	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	user1.Id = ""
	account1.Id = ""
	account1.Email = "test" + common.Random(4) + "@email.com"
	org1.Id = ""
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return nil
	}
	challenge.Users[0].User.Id = rsp_org.Data.User.Id
	challenge.OrgId = rsp_org.Data.Organisation.Id
	challenge.Users[0].User = rsp_org.Data.User

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
		Challenge: challenge,
	}
	resp_create := &behaviour_proto.CreateChallengeResponse{}
	if err := hdlr.CreateChallenge(ctx, req_create, resp_create); err != nil {
		return nil
	}
	return resp_create.Data.Challenge
}

func createHabit(ctx context.Context, hdlr *BehaviourService, t *testing.T) *behaviour_proto.Habit {
	// create behaviou category
	staticClient := static_proto.NewStaticServiceClient("go.micro.srv.static", cl)
	rsp_category, err := staticClient.CreateBehaviourCategory(ctx, &static_proto.CreateBehaviourCategoryRequest{
		Category: habit.Category,
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	habit.Category = rsp_category.Data.Category
	// create aim
	rsp_aim, err := staticClient.CreateBehaviourCategoryAim(ctx, &static_proto.CreateBehaviourCategoryAimRequest{
		BehaviourCategoryAim: habit.Target.Aim,
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	habit.Target.Aim = rsp_aim.Data.BehaviourCategoryAim
	// create marker
	rsp_marker, err := staticClient.CreateMarker(ctx, &static_proto.CreateMarkerRequest{
		Marker: habit.Target.Marker,
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	habit.Target.Marker = rsp_marker.Data.Marker
	// create method
	rsp_method, err := staticClient.CreateTrackerMethod(ctx, &static_proto.CreateTrackerMethodRequest{
		TrackerMethod: habit.Trackers[0].Method,
	})
	habit.Trackers[0].Marker = rsp_marker.Data.Marker
	habit.Trackers[0].Method = rsp_method.Data.TrackerMethod
	// create setbacks
	rsp_setback, err := staticClient.CreateSetback(ctx, &static_proto.CreateSetbackRequest{Setback: habit.Setbacks[0]})
	if err != nil {
		t.Error(err)
		return nil
	}
	habit.Setbacks[0] = rsp_setback.Data.Setback

	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	user1.Id = ""
	account1.Id = ""
	account1.Email = "test" + common.Random(4) + "@email.com"
	org1.Id = ""
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return nil
	}
	habit.Users[0].User.Id = rsp_org.Data.User.Id
	habit.OrgId = rsp_org.Data.Organisation.Id
	habit.Users[0].User = rsp_org.Data.User

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
		Habit:  habit,
	}
	resp_create := &behaviour_proto.CreateHabitResponse{}
	if err := hdlr.CreateHabit(ctx, req_create, resp_create); err != nil {
		return nil
	}
	return resp_create.Data.Habit
}

func TestAllGoals(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createGoal(ctx, hdlr, t)
	if goal == nil {
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

	req_all := &behaviour_proto.AllGoalsRequest{
		OrgId:         si.OrgId,
		TeamId:        si.UserId,
		SortParameter: "created",
		SortDirection: "DESC",
	}
	resp_all := &behaviour_proto.AllGoalsResponse{}
	err = hdlr.AllGoals(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Goals) == 0 {
		t.Error("Object count does not match")
		return
	}
}

func TestAllChallenges(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	challenge := createChallenge(ctx, hdlr, t)
	if challenge == nil {
		return
	}

	req_all := &behaviour_proto.AllChallengesRequest{}
	resp_all := &behaviour_proto.AllChallengesResponse{}
	err := hdlr.AllChallenges(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Challenges) == 0 {
		t.Error("Object count does not match")
		return
	}
}

func TestAllHabits(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	habit := createHabit(ctx, hdlr, t)
	if habit == nil {
		return
	}

	req_all := &behaviour_proto.AllHabitsRequest{}
	resp_all := &behaviour_proto.AllHabitsResponse{}
	err := hdlr.AllHabits(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Habits) == 0 {
		t.Error("Object count does not match")
		return
	}

	t.Log(resp_all.Data.Habits)
}

func TestReadGoal(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createGoal(ctx, hdlr, t)
	if goal == nil {
		return
	}

	req_read := &behaviour_proto.ReadGoalRequest{GoalId: goal.Id}
	resp_read := &behaviour_proto.ReadGoalResponse{}
	time.Sleep(2 * time.Second)
	res_read := hdlr.ReadGoal(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Goal == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Goal.Id != goal.Id {
		t.Error("Object Id does not matched")
		return
	}

	t.Log(resp_read.Data.Goal)
}

func TestReadChallenge(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	challenge := createChallenge(ctx, hdlr, t)
	if challenge == nil {
		return
	}

	req_read := &behaviour_proto.ReadChallengeRequest{ChallengeId: challenge.Id}
	resp_read := &behaviour_proto.ReadChallengeResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadChallenge(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Challenge == nil {
		t.Error("Object could not be nil")
		return
	}

	t.Log(resp_read.Data.Challenge)
}

func TestReadHabit(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	habit := createHabit(ctx, hdlr, t)
	if habit == nil {
		return
	}

	req_read := &behaviour_proto.ReadHabitRequest{HabitId: habit.Id}
	resp_read := &behaviour_proto.ReadHabitResponse{}
	time.Sleep(2 * time.Second)
	res_read := hdlr.ReadHabit(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Habit == nil {
		t.Error("Object could not be nil")
		return
	}

	t.Log(resp_read.Data.Habit)
}

func TestDeleteGoal(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createGoal(ctx, hdlr, t)
	if goal == nil {
		return
	}

	req_del := &behaviour_proto.DeleteGoalRequest{GoalId: goal.Id}
	resp_del := &behaviour_proto.DeleteGoalResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteGoal(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &behaviour_proto.ReadGoalRequest{GoalId: goal.Id}
	resp_read := &behaviour_proto.ReadGoalResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadGoal(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestDeleteChallenge(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	challenge := createChallenge(ctx, hdlr, t)
	if challenge == nil {
		return
	}

	req_del := &behaviour_proto.DeleteChallengeRequest{ChallengeId: challenge.Id}
	resp_del := &behaviour_proto.DeleteChallengeResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteChallenge(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &behaviour_proto.ReadChallengeRequest{ChallengeId: challenge.Id}
	resp_read := &behaviour_proto.ReadChallengeResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadChallenge(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data.Challenge != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestDeleteHabit(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	habit := createHabit(ctx, hdlr, t)
	if habit == nil {
		return
	}

	req_del := &behaviour_proto.DeleteHabitRequest{HabitId: habit.Id}
	resp_del := &behaviour_proto.DeleteHabitResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteHabit(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &behaviour_proto.ReadHabitRequest{HabitId: habit.Id}
	resp_read := &behaviour_proto.ReadHabitResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadHabit(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data.Habit != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestFilter(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createGoal(ctx, hdlr, t)
	if goal == nil {
		return
	}
	challenge := createChallenge(ctx, hdlr, t)
	if challenge == nil {
		return
	}
	habit := createHabit(ctx, hdlr, t)
	if habit == nil {
		return
	}

	req_filter := &behaviour_proto.FilterRequest{
		Type:     []string{"goal", "challenge", "habit"},
		Status:   []behaviour_proto.Status{behaviour_proto.Status_PUBLISHED},
		Category: []string{"category111", "category222", "category333"},
		Creator:  []string{"userid"},
	}
	resp_filter := &behaviour_proto.FilterResponse{}
	time.Sleep(time.Second)
	err := hdlr.Filter(ctx, req_filter, resp_filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_filter.Data.Goals) == 0 {
		t.Error("Goals count does not match")
		return
	}
	if len(resp_filter.Data.Challenges) == 0 {
		t.Error("Challenges count does not match")
		return
	}
	if len(resp_filter.Data.Habits) == 0 {
		t.Error("Habits count does not match")
		return
	}

	t.Log(resp_filter.Data)
}

func TestSearch(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createGoal(ctx, hdlr, t)
	if goal == nil {
		return
	}
	challenge := createChallenge(ctx, hdlr, t)
	if challenge == nil {
		return
	}
	habit := createHabit(ctx, hdlr, t)
	if habit == nil {
		return
	}

	req_search := &behaviour_proto.SearchRequest{
		Name:        "title",
		Summary:     "summary",
		Description: "description",
	}
	resp_search := &behaviour_proto.SearchResponse{}
	time.Sleep(time.Second)
	err := hdlr.Search(ctx, req_search, resp_search)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_search.Data.Goals) == 0 {
		t.Error("Goals count does not match")
		return
	}
	if len(resp_search.Data.Challenges) == 0 {
		t.Error("Challenges count does not match")
		return
	}
	if len(resp_search.Data.Habits) == 0 {
		t.Error("Habits count does not match")
		return
	}
}

func TestShareGoal(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createGoal(ctx, hdlr, t)
	if goal == nil {
		return
	}

	req_share := &behaviour_proto.ShareGoalRequest{
		Goals:  []*behaviour_proto.Goal{goal},
		Users:  goal.Users,
		UserId: goal.Users[0].User.Id,
		OrgId:  goal.Users[0].User.OrgId,
	}

	rsp_share := &behaviour_proto.ShareGoalResponse{}
	err := hdlr.ShareGoal(ctx, req_share, rsp_share)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAutocompleteGoalSearch(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createGoal(ctx, hdlr, t)
	if goal == nil {
		return
	}

	req := &behaviour_proto.AutocompleteSearchRequest{"t"}
	rsp := &behaviour_proto.AutocompleteSearchResponse{}
	err := hdlr.AutocompleteGoalSearch(ctx, req, rsp)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestShareChallenge(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_share := &behaviour_proto.ShareChallengeRequest{
		Challenges: []*behaviour_proto.Challenge{challenge},
		Users:      challenge.Users,
		UserId:     "userid",
		OrgId:      "orgid",
	}
	rsp_share := &behaviour_proto.ShareChallengeResponse{}
	err := hdlr.ShareChallenge(ctx, req_share, rsp_share)

	if err != nil {
		t.Error(err)
		return
	}
}

func TestAutocompleteChallengeSearch(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	challenge := createChallenge(ctx, hdlr, t)
	if challenge == nil {
		t.Error("create error")

		return
	}

	req := &behaviour_proto.AutocompleteSearchRequest{"t"}
	rsp := &behaviour_proto.AutocompleteSearchResponse{}
	err := hdlr.AutocompleteChallengeSearch(ctx, req, rsp)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp.Data.Response)
}

func TestShareHabit(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_share := &behaviour_proto.ShareHabitRequest{
		Habits: []*behaviour_proto.Habit{habit},
		Users:  habit.Users,
		UserId: "userid",
		OrgId:  "orgid",
	}

	rsp_share := &behaviour_proto.ShareHabitResponse{}
	err := hdlr.ShareHabit(ctx, req_share, rsp_share)

	if err != nil {
		t.Error(err)
		return
	}
}

func TestAutocompleteHabitSearch(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	habit := createHabit(ctx, hdlr, t)
	if habit == nil {
		t.Error("create error")
		return
	}

	req := &behaviour_proto.AutocompleteSearchRequest{"t"}
	rsp := &behaviour_proto.AutocompleteSearchResponse{}
	err := hdlr.AutocompleteHabitSearch(ctx, req, rsp)
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

	goal := createGoal(ctx, hdlr, t)
	if goal == nil {
		return
	}

	rsp := &behaviour_proto.GetTopTagsResponse{}
	if err := hdlr.GetTopTags(ctx, &behaviour_proto.GetTopTagsRequest{
		Object: common.GOAL,
		OrgId:  goal.OrgId,
		N:      5,
	}, rsp); err != nil {
		t.Error(err)
		return
	}
}

func TestAutocompleteTags(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	goal := createGoal(ctx, hdlr, t)
	if goal == nil {
		t.Error("create error")
		return
	}

	rsp := &behaviour_proto.AutocompleteTagsResponse{}
	if err := hdlr.AutocompleteTags(ctx, &behaviour_proto.AutocompleteTagsRequest{
		Object: common.GOAL,
		OrgId:  goal.OrgId,
		Name:   "t",
	}, rsp); err != nil {
		t.Error(err)
		return
	}
}
