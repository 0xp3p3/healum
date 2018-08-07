package handler

import (
	"context"
	account_proto "server/account-srv/proto/account"
	"server/common"
	content_proto "server/content-srv/proto/content"
	kv_proto "server/kv-srv/proto/kv"
	"server/organisation-srv/db"
	organisation_proto "server/organisation-srv/proto/organisation"
	static_db "server/static-srv/db"
	static_hdlr "server/static-srv/handler"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
	user_proto "server/user-srv/proto/user"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/jinzhu/copier"
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
	static_db.Init(cl)
}

var account = &account_proto.Account{
	Email:    "email" + common.Random(4) + "@email.com",
	Password: "pass1",
}

var user = &user_proto.User{
	// Id:        "user_id",
	Firstname: "David",
	Lastname:  "John",
	AvatarUrl: "http://example.com",
	Tokens: []*user_proto.Token{
		{"11671c2e7da30e3c393813f60b327f9c2e2e08390761aa01e37ba5d3e6a617be", 1, "aaa"}, {"token_b", 2, "bbb"},
	},
}

var org = &organisation_proto.Organisation{
	Type: organisation_proto.OrganisationType_NONE,
}

var module = &static_proto.Module{
	Id:       "111",
	Name:     "module1",
	IconSlug: "icon_slug",
}

var org_profile = &organisation_proto.OrganisationProfile{}

var org_setting = &organisation_proto.OrganisationSetting{}

func createOrganisation(ctx context.Context, hdlr *OrganisationService, t *testing.T) *organisation_proto.Organisation {
	// create role
	_, err := hdlr.StaticClient.CreateRole(ctx, &static_proto.CreateRoleRequest{
		&static_proto.Role{Name: "admin_role", NameSlug: "admin"},
	})
	if err != nil {
		t.Error(err)
		return nil
	}

	req := &organisation_proto.CreateRequest{
		Organisation: org,
		Account:      account,
		User:         user,
		Modules:      []*static_proto.Module{module},
	}
	rsp := &organisation_proto.CreateResponse{}
	err = hdlr.Create(ctx, req, rsp)
	if err != nil {
		t.Error(err)
		return nil
	}

	return rsp.Data.Organisation
}

func createModule(ctx context.Context, t *testing.T) *static_proto.Module {
	hdlr := new(static_hdlr.StaticService)

	req_create := &static_proto.CreateModuleRequest{Module: module}
	rsp_create := &static_proto.CreateModuleResponse{}
	err := hdlr.CreateModule(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return rsp_create.Data.Module
}

func initHandler() *OrganisationService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	hdlr := &OrganisationService{
		Broker:        nats_brker,
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		TeamClient:    team_proto.NewTeamServiceClient("go.micro.srv.team", cl),
		StaticClient:  static_proto.NewStaticServiceClient("go.micro.srv.static", cl),
		UserClient:    user_proto.NewUserServiceClient("go.micro.srv.user", cl),
	}
	return hdlr
}

func TestAll(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	// is created organisation
	org := createOrganisation(ctx, hdlr, t)
	if org == nil {
		return
	}

	req_all := &organisation_proto.AllRequest{
		SortParameter: "created",
		SortDirection: "DESC",
	}
	resp_all := &organisation_proto.AllResponse{}
	err := hdlr.All(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if resp_all.Data.Organisations[0].Id != org.Id {
		t.Error("Object does not match")
		return
	}
}

func TestRead(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	// is created organisation
	org := createOrganisation(ctx, hdlr, t)
	if org == nil {
		return
	}

	req_read := &organisation_proto.ReadRequest{OrgId: org.Id}
	rsp_read := &organisation_proto.ReadResponse{}
	err := hdlr.Read(ctx, req_read, rsp_read)
	if err != nil {
		t.Error(err)
		return
	}
	if rsp_read.Data.Organisation == nil {
		t.Error("Object could not be nil")
		return
	}
	if rsp_read.Data.Organisation.Id != org.Id {
		t.Error("Id does not match")
		return
	}
}

func TestCreateOrganisationProfile(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	req := &organisation_proto.CreateOrganisationProfileRequest{
		Profile: org_profile,
	}
	rsp := &organisation_proto.CreateOrganisationProfileResponse{}
	err := hdlr.CreateOrganisationProfile(ctx, req, rsp)
	if err != nil {
		t.Error(err)
	}

	if len(rsp.Data.Profile.Id) == 0 {
		t.Error("Id does not matched")
		return
	}
}

func TestCreateOrganisationSetting(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	req := &organisation_proto.CreateOrganisationSettingRequest{
		Setting: org_setting,
	}
	rsp := &organisation_proto.CreateOrganisationSettingResponse{}
	err := hdlr.CreateOrganisationSetting(ctx, req, rsp)
	if err != nil {
		t.Error(err)
	}

	if len(rsp.Data.Setting.Id) == 0 {
		t.Error("Id does not matched")
		return
	}
}

func TestUpdateModulesByOrg(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	//create module
	module := createModule(ctx, t)
	if module == nil {
		t.Error()
		return
	}
	// is created organisation
	if err := createOrganisation(ctx, hdlr, t); err != nil {
		t.Error(err)
		return
	}

	t.Log("orgId: ", org.Id)
	req_put := &organisation_proto.UpdateModulesRequest{OrgId: org.Id, Modules: []*static_proto.Module{module}}
	rsp_put := &organisation_proto.UpdateModulesResponse{}

	if err := hdlr.UpdateModules(ctx, req_put, rsp_put); err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp_put.Data.Modules)

}

func TestReadOrgInfo(t *testing.T) {
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	org.Id = "orgid"
	req_put := &organisation_proto.PutOrgInfoRequest{OrgId: org.Id, Org: org}
	rsp_put := &organisation_proto.PutOrgInfoResponse{}
	if err := hdlr.PutOrgInfo(ctx, req_put, rsp_put); err != nil {
		t.Error(err)
		return
	}

	req_read := &organisation_proto.ReadOrgInfoRequest{org.Id}
	rsp_read := &organisation_proto.ReadOrgInfoResponse{}
	if err := hdlr.ReadOrgInfo(ctx, req_read, rsp_read); err != nil {
		t.Error(err)
		return
	}

	if rsp_read.OrgInfo == nil {
		t.Error("Object does not matched")
	}

	t.Log(rsp_read.OrgInfo)
}

func TestUnmarshal(t *testing.T) {
	logJson := `{"org_id":"orgid","type":"NONE","owner":{"id":"","org_id":"orgid","created":"0","updated":"0","firstname":"David","lastname":"John","image":"","gender":"MALE","dob":"0","contactDetails":[],"addresses":[],"tokens":[]}}`
	// um := jsonpb.Unmarshaler{}
	oi := organisation_proto.OrgInfo{}
	err := jsonpb.Unmarshal(strings.NewReader(logJson), &oi)
	// assert.NoError(t, err)

	oi_temp := &organisation_proto.OrgInfo{}
	copier.Copy(oi_temp, &oi)
	t.Log(oi_temp)
	t.Log(err)
}

func TestGetModulesByOrg(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	// create modules
	module := createModule(ctx, t)
	if module == nil {
		return
	}

	// is created organisation
	org.Modules = []*static_proto.Module{module}
	org := createOrganisation(ctx, hdlr, t)
	if org == nil {
		return
	}

	req_get := &organisation_proto.GetModulesByOrgRequest{OrgId: org.Id}
	rsp_get := &organisation_proto.GetModulesByOrgResponse{}
	if err := hdlr.GetModulesByOrg(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp_get.Data.Modules)
}

func TestUpdate(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	// is created organisation
	org := createOrganisation(ctx, hdlr, t)
	if org == nil {
		return
	}

	org.Name = "updated org"
	org.Locations = []*content_proto.Place{{Id: "111", Name: "place1"}}
	req := &organisation_proto.UpdateRequest{Organisation: org}
	rsp := &organisation_proto.UpdateResponse{}
	if err := hdlr.Update(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	if rsp.Data.Organisation.Id != org.Id {
		t.Error("organisation id is not matched")
		return
	}

	if rsp.Data.Organisation.Name != org.Name {
		t.Error("name does not matched")
		return
	}
}
