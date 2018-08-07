package handler

import (
	"bytes"
	"context"
	"encoding/json"
	account_proto "server/account-srv/proto/account"
	behaviour_hdlr "server/behaviour-srv/handler"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	organisation_proto "server/organisation-srv/proto/organisation"
	product_proto "server/product-srv/proto/product"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
	track_proto "server/track-srv/proto/track"
	userapp_db "server/user-app-srv/db"
	"server/user-srv/db"
	user_proto "server/user-srv/proto/user"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	google_protobuf1 "github.com/golang/protobuf/ptypes/struct"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
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
		{user, 100, behaviour_proto.ExpectedProgressType_LINEAR, ""},
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
		{user, 100, behaviour_proto.ExpectedProgressType_LINEAR, ""},
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
		{user, 100, behaviour_proto.ExpectedProgressType_LINEAR, ""},
	},
	Duration: "P1Y2DT3H4M5S",
}

var org1 = &organisation_proto.Organisation{
	Type: organisation_proto.OrganisationType_ROOT,
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
	Conditions: []*static_proto.ContentCategoryItem{
		{
			Id:    "content_category_item_id1",
			OrgId: "org_id",
		}, {
			Id:    "content_category_item_id2",
			OrgId: "org_id",
		},
	},
	Allergies: []*static_proto.ContentCategoryItem{
		{
			Id:    "content_category_item_id3",
			OrgId: "org_id",
		}, {
			Id:    "content_category_item_id4",
			OrgId: "org_id",
		},
	},
	Ethinicties: []*static_proto.ContentCategoryItem{
		{
			Id:    "content_category_item_id5",
			OrgId: "org_id",
		}, {
			Id:    "content_category_item_id6",
			OrgId: "org_id",
		},
	},
	Food: []*static_proto.ContentCategoryItem{
		{
			Id:    "content_category_item_id7",
			OrgId: "org_id",
		}, {
			Id:    "content_category_item_id8",
			OrgId: "org_id",
		},
	},
	Cuisines: []*static_proto.ContentCategoryItem{
		{
			Id:    "content_category_item_id8",
			OrgId: "org_id",
		}, {
			Id:    "content_category_item_id10",
			OrgId: "org_id",
		},
	},
}

var batch = static_proto.Batch{
	Name: "sample_batch",
}

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

func initHandler() *UserService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()

	hdlr := &UserService{
		Broker:          nats_brker,
		AccountClient:   account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		TrackClient:     track_proto.NewTrackServiceClient("go.micro.srv.track", cl),
		KvClient:        kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		StaticClient:    static_proto.NewStaticServiceClient("go.micro.srv.static", cl),
		TeamClient:      team_proto.NewTeamServiceClient("go.micro.srv.team", cl),
		BehaviourClient: behaviour_proto.NewBehaviourServiceClient("go.micro.srv.behaviour", cl),
	}
	return hdlr
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

func TestCreateWithAccount(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)

	if err != nil {
		t.Error(err)
		return
	}
}

func TestCreateWithAccountPhone(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account_phone,
	}
	rsp_create := &user_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)

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

func TestCreateWithAccountAndPointOfContact(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_create_poc := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}

	rsp_create_poc := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create_poc, rsp_create_poc); err != nil {
		t.Error(err)
		return
	}

	user.PointOfContact = &user_proto.User{
		Id: rsp_create_poc.Data.User.Id,
	}
	user.Id = ""
	user.Preference = &user_proto.Preferences{}
	account.Email = "email" + common.Random(4) + "@email.com"
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}
}

//this test fails
func TestCreate(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	account.Email = "email" + common.Random(4) + "@email.com"
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}
}

func TestUpdate(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create object
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}
	user := &user_proto.User{
		Id:        rsp_create.Data.User.Id,
		OrgId:     rsp_create.Data.User.OrgId,
		Firstname: "first_name",
		Lastname:  "last_name",
		AvatarUrl: "example.jpg",
	}
	// update user
	req_update := &user_proto.UpdateRequest{
		User: user,
	}
	rsp_update := &user_proto.UpdateResponse{}
	if err := hdlr.Update(ctx, req_update, rsp_update); err != nil {
		t.Error(err)
		return
	}

	if rsp_update.Data == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestAll(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create batch
	productClient := product_proto.NewProductServiceClient("go.micro.srv.product", cl)
	rsp_batch, err := productClient.CreateBatch(ctx, &product_proto.CreateBatchRequest{Batch: &batch})
	if err != nil {
		t.Error(err)
		return
	}
	user.CurrentBatch = rsp_batch.Data.Batch
	account.Email = "email" + common.Random(4) + "@email.com"
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	req_all := &user_proto.AllRequest{}
	rsp_all := &user_proto.AllResponse{}
	if err := hdlr.All(ctx, req_all, rsp_all); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_all.Data.Users) == 0 {
		t.Error("Object count does not matched")
		return
	}

	t.Log(rsp_batch.Data.Batch)
}

func TestRead(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create object
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}
	// read user
	req_read := &user_proto.ReadRequest{
		UserId: rsp_create.Data.User.Id,
	}
	rsp_read := &user_proto.ReadResponse{}
	if err := hdlr.Read(ctx, req_read, rsp_read); err != nil {
		t.Error(err)
		return
	}

	if rsp_read.Data == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestReadByAccount(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	account.Email = "email" + common.Random(4) + "@email.com"
	// create object
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}
	// read user by account id TODO:(this needs a real id created during test)
	req_read := &user_proto.ReadByAccountRequest{
		AccountId: "a9b14682-7aea-11e8-bf1d-20c9d0453b15",
	}
	rsp_read := &user_proto.ReadByAccountResponse{}
	if err := hdlr.ReadByAccount(ctx, req_read, rsp_read); err != nil {
		t.Error(err)
		return
	}

	if rsp_read.Data == nil {
		t.Error("Object does not matched")
		return
	}
}

//this test fails
func TestFilter(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	// filter users
	req_filter := &user_proto.FilterRequest{
		Users: []string{rsp_create.Data.User.Id, "56ebe49b-5649-11e8-8c71-00155d4b0101", "2654e022-5647-11e8-9ff0-00155d4b0101"},
		OrgId: rsp_create.Data.User.OrgId,
	}
	rsp_filter := &user_proto.FilterResponse{}
	if err := hdlr.Filter(ctx, req_filter, rsp_filter); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_filter.Data.Users) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestUpdateTokens(t *testing.T) {
	initDb()

	ctx := common.NewTestContext(context.TODO())
	hdlr := &UserService{}

	tokens := []*user_proto.Token{
		{"abcd1234", 1, "aaaa"},
	}
	req_token := &user_proto.UpdateTokenRequest{"userid", tokens}
	rsp_token := &user_proto.UpdateTokenResponse{}
	if err := hdlr.UpdateTokens(ctx, req_token, rsp_token); err != nil {
		t.Error(err)
	}
}

func TestReadTokens(t *testing.T) {
	initDb()

	ctx := common.NewTestContext(context.TODO())
	hdlr := &UserService{}

	user_ids := []string{"userid", "8b291ae8-7a0c-11e8-b006-20c9d0453b15"}
	req_token := &user_proto.ReadTokensRequest{user_ids}
	rsp_token := &user_proto.ReadTokensResponse{}
	if err := hdlr.ReadTokens(ctx, req_token, rsp_token); err != nil {
		t.Error(err)
		return
	}
}

func TestReadUserPreference(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())

	// create user preference
	if err := db.SaveUserPreference(ctx, &user_proto.Preferences{
		OrgId:  "orgid",
		UserId: "userid",
	}); err != nil {
		t.Error(err)
		return
	}

	hdlr := initHandler()
	req_user := &user_proto.ReadUserPreferenceRequest{
		UserId: "userid",
		OrgId:  "orgid",
	}
	rsp_user := &user_proto.ReadUserPreferenceResponse{}
	if err := hdlr.ReadUserPreference(ctx, req_user, rsp_user); err != nil {
		t.Error(err)
		return
	}

	if rsp_user.Data.Preference == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestListUserFeedback(t *testing.T) {
	initDb()
	userapp_db.Init(cl)
	ctx := common.NewTestContext(context.TODO())

	// create user_feedback
	if err := userapp_db.SaveUserFeedback(ctx, &user_proto.UserFeedback{
		UserId:   "userid",
		OrgId:    "orgid",
		Feedback: "hello world",
	}); err != nil {
		t.Error(err)
		return
	}

	hdlr := initHandler()
	req_user := &user_proto.ListUserFeedbackRequest{
		UserId: "userid",
	}
	rsp_user := &user_proto.ListUserFeedbackResponse{}
	if err := hdlr.ListUserFeedback(ctx, req_user, rsp_user); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_user.Data.Feedbacks) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestFilterUser(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create user
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	// filter
	req_filter := &user_proto.FilterUserRequest{
		// Tags:       []string{"a", "c"},
		// Preference: &user_proto.Preferences{Id: user.Preference.Id},
		Status: account.Status,
	}
	rsp_filter := &user_proto.FilterUserResponse{}
	if err := hdlr.FilterUser(ctx, req_filter, rsp_filter); err != nil {
		t.Error(err)
		return
	}

	// data response
	if len(rsp_filter.Data.Response) == 0 {
		t.Error("Object count does not matched")
		return
	}
	t.Log(rsp_filter.Data.Response)
}

func TestDeleteUser(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	// create user
	account.Email = "email" + common.Random(4) + "@email.com"
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	// delete user
	req_del := &user_proto.DeleteRequest{UserId: rsp_create.Data.Account.Id}
	rsp_del := &user_proto.DeleteResponse{}

	time.Sleep(2 * time.Second)
	if err := hdlr.DeleteByUserId(ctx, req_del, rsp_del); err != nil {
		t.Error(err)
		return
	}
}

func TestSearchUser(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create user
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	// search
	req_search := &user_proto.SearchUserRequest{
		Name: "jo",
		// Gender:         user_proto.Gender_MALE,
		// Addresses: []*static_proto.Address{{PostalCode: "111000"}},
		ContactDetails: []*user_proto.ContactDetail{{Id: "contact_detail_id"}},
	}
	rsp_search := &user_proto.SearchUserResponse{}
	if err := hdlr.SearchUser(ctx, req_search, rsp_search); err != nil {
		t.Error(err)
		return
	}

	// data response
	if len(rsp_search.Data.Response) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestAutocompleteUser(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create user
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	// query
	req_auto := &user_proto.AutocompleteUserRequest{
		Name: "jo",
	}
	rsp_auto := &user_proto.AutocompleteUserResponse{}
	if err := hdlr.AutocompleteUser(ctx, req_auto, rsp_auto); err != nil {
		t.Error(err)
		return
	}

	// data response
	if len(rsp_auto.Data.Response) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestSetAccountStatus(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create user
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	// set status
	req_status := &account_proto.SetAccountStatusRequest{
		UserId: rsp_create.Data.User.Id,
		Status: account_proto.AccountStatus_ACTIVE,
	}
	rsp_status := &account_proto.SetAccountStatusResponse{}

	if err := hdlr.SetAccountStatus(ctx, req_status, rsp_status); err != nil {
		t.Error(err)
		return
	}
}

func TestResetUserPasswordWithEmail(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create user
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}

	// reset password
	req_reset := &account_proto.ResetUserPasswordRequest{
		UserId:   rsp_create.Data.User.Id,
		Password: "hello",
	}
	rsp_reset := &account_proto.ResetUserPasswordResponse{}

	if err := hdlr.ResetUserPassword(ctx, req_reset, rsp_reset); err != nil {
		t.Error(err)
		return
	}
}

func TestResetUserPasswordWithPhone(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create user
	account_phone.Phone = common.Random(8)
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account_phone,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
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

	// reset password
	req_reset := &account_proto.ResetUserPasswordRequest{
		UserId:   rsp_create.Data.User.Id,
		OrgId:    si.OrgId,
		TeamId:   si.UserId,
		Passcode: "12345",
	}
	rsp_reset := &account_proto.ResetUserPasswordResponse{}

	if err := hdlr.ResetUserPassword(ctx, req_reset, rsp_reset); err != nil {
		t.Error(err)
		return
	}
}

var measurements = []*user_proto.Measurement{
	{
		OrgId:  "orgid",
		Marker: &static_proto.Marker{},
		Method: &static_proto.TrackerMethod{},
		Unit:   "unit_test",
	},
}

func TestGetSharedResources(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	//	behaviour_hdlr := initBehaviourHandler()
	behaviourClient := behaviour_proto.NewBehaviourServiceClient("go.micro.srv.behaviour", cl)

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

	// create goal
	goal.CreatedBy = &user_proto.User{Id: si.UserId}
	goal.OrgId = si.OrgId
	goal.Users[0].User = &user_proto.User{Id: si.UserId}

	rsp_goal, err := behaviourClient.CreateGoal(ctx, &behaviour_proto.CreateGoalRequest{
		UserId: si.UserId,
		OrgId:  si.OrgId,
		Goal:   goal,
		TeamId: si.UserId,
	})
	if err != nil {
		t.Error(err)
		return
	}

	// create challenge
	challenge.CreatedBy = &user_proto.User{Id: si.UserId}
	challenge.OrgId = si.OrgId
	challenge.Users[0].User = &user_proto.User{Id: si.UserId}

	rsp_challenge, err := behaviourClient.CreateChallenge(ctx, &behaviour_proto.CreateChallengeRequest{
		UserId:    si.UserId,
		OrgId:     si.OrgId,
		Challenge: challenge,
	})
	if err != nil {
		t.Error(err)
		return
	}
	// create habit
	habit.CreatedBy = &user_proto.User{Id: si.UserId}
	habit.OrgId = si.OrgId
	habit.Users[0].User = &user_proto.User{Id: si.UserId}

	//	rsp_challenge.Data.Challenge.Id

	rsp_habit, err := behaviourClient.CreateHabit(ctx, &behaviour_proto.CreateHabitRequest{
		UserId: si.UserId,
		OrgId:  si.OrgId,
		Habit:  habit,
	})
	if err != nil {
		t.Error(err)
		return
	}

	// share goal
	req_share := &user_proto.GetSharedResourcesRequest{
		Type: []string{"healum.com/proto/go.micro.srv.behaviour.Goal",
			"healum.com/proto/go.micro.srv.behaviour.Challenge",
			"healum.com/proto/go.micro.srv.behaviour.Habit"},
		UserId:   si.UserId,
		OrgId:    si.OrgId,
		Status:   []static_proto.ShareStatus{static_proto.ShareStatus_SHARED},
		SharedBy: []string{si.UserId},
	}

	rsp_share := &user_proto.GetSharedResourcesResponse{}

	time.Sleep(2 * time.Second)
	if err := hdlr.GetSharedResources(ctx, req_share, rsp_share); err != nil {
		t.Error(err)
		return
	}

	is_goal := false
	is_challenge := false
	is_habit := false
	for _, s := range rsp_share.Data.SharedResources {
		if s.ResourceId == rsp_goal.Data.Goal.Id {
			is_goal = true
		}
		if s.ResourceId == rsp_challenge.Data.Challenge.Id {
			is_challenge = true
		}
		if s.ResourceId == rsp_habit.Data.Habit.Id {
			is_habit = true
		}
	}

	if !is_goal {
		t.Error("shared goal doesn't exist")
		return
	}
	if !is_challenge {
		t.Error("shared challenge doesn't exist")
		return
	}
	if !is_habit {
		t.Error("shared habit doesn't exist")
		return
	}
}

func TestSearchSharedResources(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	behaviourClient := behaviour_proto.NewBehaviourServiceClient("go.micro.srv.behaviour", cl)

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

	// create goal
	goal.CreatedBy = &user_proto.User{Id: si.UserId}
	goal.OrgId = si.OrgId
	goal.Users[0].User = &user_proto.User{Id: si.UserId}

	rsp_goal, err := behaviourClient.CreateGoal(ctx, &behaviour_proto.CreateGoalRequest{
		UserId: si.UserId,
		OrgId:  si.OrgId,
		Goal:   goal,
		TeamId: si.UserId,
	})
	if err != nil {
		t.Error(err)
		return
	}

	// create challenge
	challenge.CreatedBy = &user_proto.User{Id: si.UserId}
	challenge.OrgId = si.OrgId
	challenge.Users[0].User = &user_proto.User{Id: si.UserId}

	rsp_challenge, err := behaviourClient.CreateChallenge(ctx, &behaviour_proto.CreateChallengeRequest{
		UserId:    si.UserId,
		OrgId:     si.OrgId,
		Challenge: challenge,
	})
	if err != nil {
		t.Error(err)
		return
	}
	// create habit
	habit.CreatedBy = &user_proto.User{Id: si.UserId}
	habit.OrgId = si.OrgId
	habit.Users[0].User = &user_proto.User{Id: si.UserId}

	rsp_habit, err := behaviourClient.CreateHabit(ctx, &behaviour_proto.CreateHabitRequest{
		UserId: si.UserId,
		OrgId:  si.OrgId,
		Habit:  habit,
	})
	if err != nil {
		t.Error(err)
		return
	}

	// share goal
	req_share := &user_proto.SearchSharedResourcesRequest{
		Type: []string{"healum.com/proto/go.micro.srv.behaviour.Goal",
			"healum.com/proto/go.micro.srv.behaviour.Challenge",
			"healum.com/proto/go.micro.srv.behaviour.Habit"},
		UserId:   si.UserId,
		OrgId:    si.OrgId,
		Status:   []static_proto.ShareStatus{static_proto.ShareStatus_SHARED},
		SharedBy: []string{si.UserId},
	}

	rsp_share := &user_proto.SearchSharedResourcesResponse{}

	time.Sleep(2 * time.Second)
	if err := hdlr.SearchSharedResources(ctx, req_share, rsp_share); err != nil {
		t.Error(err)
		return
	}

	is_goal := false
	is_challenge := false
	is_habit := false
	for _, s := range rsp_share.Data.SharedResources {
		if s.ResourceId == rsp_goal.Data.Goal.Id {
			is_goal = true
		}
		if s.ResourceId == rsp_challenge.Data.Challenge.Id {
			is_challenge = true
		}
		if s.ResourceId == rsp_habit.Data.Habit.Id {
			is_habit = true
		}
	}

	if !is_goal {
		t.Error("shared goal doesn't exist")
		return
	}
	if !is_challenge {
		t.Error("shared challenge doesn't exist")
		return
	}
	if !is_habit {
		t.Error("shared habit doesn't exist")
		return
	}
}

func TestAddMultipleMeasurements(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create user
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)

	if err != nil {
		t.Error(err)
		return
	}
	// employee or normal user
	u := rsp_create.Data.User
	measurements[0].UserId = u.Id

	//create tracker method
	trackerMethod := &static_proto.TrackerMethod{
		Id:       "111",
		Name:     "title",
		NameSlug: "hcp",
		IconSlug: "iconSlug",
	}
	static_client := static_proto.NewStaticServiceClient("go.micro.srv.static", cl)

	req_create_trackmarkermethod := &static_proto.CreateTrackerMethodRequest{TrackerMethod: trackerMethod}
	resp_create_trackmarkermethod, err := static_client.CreateTrackerMethod(ctx, req_create_trackmarkermethod)
	if err != nil {
		t.Error(err)
		return
	}

	measurements[0].Method = resp_create_trackmarkermethod.Data.TrackerMethod

	// create marker
	rsp_marker, err := static_client.CreateMarker(ctx, &static_proto.CreateMarkerRequest{
		Marker: &static_proto.Marker{
			Name:           "test_marker",
			TrackerMethods: []*static_proto.TrackerMethod{resp_create_trackmarkermethod.Data.TrackerMethod},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	measurements[0].Marker = rsp_marker.Data.Marker

	// create value
	var v google_protobuf1.Value
	raw1 := `"hello world"`
	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &v); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}
	measurements[0].Value = &v
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

	// add measurements
	req_add := &user_proto.AddMultipleMeasurementsRequest{
		Measurements: measurements,
		UserId:       si.UserId,
		OrgId:        si.OrgId,
		TeamId:       si.UserId,
	}
	rsp_add := &user_proto.AddMultipleMeasurementsResponse{}
	if err := hdlr.AddMultipleMeasurements(ctx, req_add, rsp_add); err != nil {
		t.Error(err)
		return
	}
}

func TestGetAllMeasurementsHistory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create user
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)

	if err != nil {
		t.Error(err)
		return
	}
	// employee or normal user
	u := rsp_create.Data.User
	measurements[0].UserId = u.Id

	//create tracker method
	trackerMethod := &static_proto.TrackerMethod{
		Id:       "111",
		Name:     "title",
		NameSlug: "hcp",
		IconSlug: "iconSlug",
	}
	static_client := static_proto.NewStaticServiceClient("go.micro.srv.static", cl)

	req_create_trackmarkermethod := &static_proto.CreateTrackerMethodRequest{TrackerMethod: trackerMethod}
	resp_create_trackmarkermethod, err := static_client.CreateTrackerMethod(ctx, req_create_trackmarkermethod)
	if err != nil {
		t.Error(err)
		return
	}

	measurements[0].Method = resp_create_trackmarkermethod.Data.TrackerMethod

	// create marker
	rsp_marker, err := static_client.CreateMarker(ctx, &static_proto.CreateMarkerRequest{
		Marker: &static_proto.Marker{
			Name: "test_marker",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	measurements[0].Marker = rsp_marker.Data.Marker

	// create value
	var v google_protobuf1.Value
	raw1 := `"hello world"`
	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &v); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}
	measurements[0].Value = &v

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
	// add measurements
	req_add := &user_proto.AddMultipleMeasurementsRequest{
		Measurements: measurements,
		UserId:       si.UserId,
		OrgId:        si.OrgId,
		TeamId:       si.UserId,
	}
	rsp_add := &user_proto.AddMultipleMeasurementsResponse{}
	if err := hdlr.AddMultipleMeasurements(ctx, req_add, rsp_add); err != nil {
		t.Error(err)
		return
	}

	// test GetAllMeasurementsHistory
	rsp_history := &user_proto.GetAllMeasurementsHistoryResponse{}
	if err := hdlr.GetAllMeasurementsHistory(ctx, &user_proto.GetAllMeasurementsHistoryRequest{
		UserId: si.UserId,
		OrgId:  si.OrgId,
		TeamId: si.UserId,
	}, rsp_history); err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_history.Data.Measurements)
}

func TestGetMeasurementsHistory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	//create user
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}

	err := hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}
	//employee or normal user
	u := rsp_create.Data.User
	measurements[0].UserId = u.Id

	//create tracker method
	trackMethod := &static_proto.TrackerMethod{
		Id:       "111",
		Name:     "title",
		NameSlug: "hcp",
		IconSlug: "iconSlug",
	}
	static_client := static_proto.NewStaticServiceClient("go.micro.srv.static", cl)

	req_create_trackmarkermethod := &static_proto.CreateTrackerMethodRequest{TrackerMethod: trackMethod}
	resp_create_trackmarkermethod, err := static_client.CreateTrackerMethod(ctx, req_create_trackmarkermethod)
	if err != nil {
		t.Error(err)
		return
	}
	measurements[0].Method = resp_create_trackmarkermethod.Data.TrackerMethod
	// create marker
	rsp_marker, err := static_client.CreateMarker(ctx, &static_proto.CreateMarkerRequest{
		Marker: &static_proto.Marker{
			Name: "test_marker",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	measurements[0].Marker = rsp_marker.Data.Marker

	// create value
	var v google_protobuf1.Value
	raw1 := `"hello world"`
	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &v); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}
	measurements[0].Value = &v

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

	// add measurements
	req_add := &user_proto.AddMultipleMeasurementsRequest{
		Measurements: measurements,
		UserId:       si.UserId,
		OrgId:        si.OrgId,
		TeamId:       si.UserId,
	}
	rsp_add := &user_proto.AddMultipleMeasurementsResponse{}
	if err := hdlr.AddMultipleMeasurements(ctx, req_add, rsp_add); err != nil {
		t.Error(err)
		return
	}

	// test GetMeasurementsHistory
	rsp_history := &user_proto.GetMeasurementsHistoryResponse{}
	if err := hdlr.GetMeasurementsHistory(ctx, &user_proto.GetMeasurementsHistoryRequest{
		UserId:   si.UserId,
		OrgId:    si.OrgId,
		TeamId:   si.UserId,
		MarkerId: rsp_marker.Data.Marker.Id,
	}, rsp_history); err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_history.Data.Measurements)
}

func TestGetAllTrackedMarkers(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create user
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account,
	}
	rsp_create := &user_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)

	if err != nil {
		t.Error(err)
		return
	}
	// employee or normal user
	u := rsp_create.Data.User
	measurements[0].UserId = u.Id

	//create tracker method
	trackerMethod := &static_proto.TrackerMethod{
		Id:       "111",
		Name:     "title",
		NameSlug: "hcp",
		IconSlug: "iconSlug",
	}
	static_client := static_proto.NewStaticServiceClient("go.micro.srv.static", cl)

	req_create_trackmarkermethod := &static_proto.CreateTrackerMethodRequest{TrackerMethod: trackerMethod}
	resp_create_trackmarkermethod, err := static_client.CreateTrackerMethod(ctx, req_create_trackmarkermethod)
	if err != nil {
		t.Error(err)
		return
	}

	measurements[0].Method = resp_create_trackmarkermethod.Data.TrackerMethod

	// create marker
	rsp_marker, err := static_client.CreateMarker(ctx, &static_proto.CreateMarkerRequest{
		Marker: &static_proto.Marker{
			Name: "test_marker",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	measurements[0].Marker = rsp_marker.Data.Marker

	// create value
	var v google_protobuf1.Value
	raw1 := `"hello world"`
	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &v); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}
	measurements[0].Value = &v

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

	// add measurements
	req_add := &user_proto.AddMultipleMeasurementsRequest{
		Measurements: measurements,
		UserId:       si.UserId,
		OrgId:        si.OrgId,
		TeamId:       si.UserId,
	}
	rsp_add := &user_proto.AddMultipleMeasurementsResponse{}
	if err := hdlr.AddMultipleMeasurements(ctx, req_add, rsp_add); err != nil {
		t.Error(err)
		return
	}

	// test GetAllTrackedMarkers
	rsp_get := &user_proto.GetAllTrackedMarkersResponse{}
	if err := hdlr.GetAllTrackedMarkers(ctx, &user_proto.GetAllTrackedMarkersRequest{
		UserId: si.UserId,
		OrgId:  si.OrgId,
		TeamId: si.UserId,
	}, rsp_get); err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_get.Data.Markers)
}
