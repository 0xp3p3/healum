package handler

import (
	"bytes"
	"context"
	"encoding/json"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	organisation_proto "server/organisation-srv/proto/organisation"
	product_proto "server/product-srv/proto/product"
	static_proto "server/static-srv/proto/static"
	"server/team-srv/db"
	team_proto "server/team-srv/proto/team"
	user_db "server/user-srv/db"
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
	// ctx := co	mmon.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
	user_db.Init(cl)
}

var team = &team_proto.Team{
	OrgId:       "orgid",
	Name:        "team1",
	Description: "hello world",
	Image:       "iamge001",
	Color:       "red",
	CreatedBy:   &user_proto.User{Id: "111"},
}

var account = &account_proto.Account{
	Email:    "email" + common.Random(4) + "@email.com",
	Password: "pass1",
}

var user = &user_proto.User{
	OrgId:     "orgid",
	Firstname: "David",
	Lastname:  "John",
	AvatarUrl: "http://example.com",
	Tokens: []*user_proto.Token{
		{"11671c2e7da30e3c393813f60b327f9c2e2e08390761aa01e37ba5d3e6a617be", 1, "aaa"}, {"token_b", 2, "bbb"},
	},
}

var employee = &team_proto.Employee{
	OrgId: "orgid",
	Role:  role,
	Profile: &team_proto.EmployeeProfile{
		OrgId: "orgid",
	},
	Teams: []*team_proto.Team{team},
}

var role = &static_proto.Role{
	OrgId: "orgid",
	Name:  "own",
}

var org1 = &organisation_proto.Organisation{
	Type: organisation_proto.OrganisationType_ROOT,
}

var product = &product_proto.Product{
	Name:  "product",
	OrgId: "orgid",
}

var service = &product_proto.Service{
	Name:  "service",
	OrgId: "orgid",
}

func initHandler() *TeamService {
	hdlr := &TeamService{
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		UserClient:    user_proto.NewUserServiceClient("go.micro.srv.user", cl),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
	}
	return hdlr
}

func createTeam(ctx context.Context, hdlr *TeamService, t *testing.T) *team_proto.Team {
	// c	reate org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	user.Id = ""
	account.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user, Account: account})
	if err != nil {
		t.Error(err)
		return nil
	}
	// team.TeamMembers = []*user_proto.User{rsp_org.Data.User}
	team.CreatedBy = rsp_org.Data.User

	// create role
	staticClient := static_proto.NewStaticServiceClient("go.micro.srv.static", cl)
	rsp_role, err := staticClient.CreateRole(ctx, &static_proto.CreateRoleRequest{
		&static_proto.Role{Name: "sample_role", NameSlug: "sample"},
	})
	if err != nil {
		t.Error(err)
		return nil
	}
	// create teammemeber
	team.TeamMembers = []*team_proto.TeamMember{
		{User: rsp_org.Data.User, Role: rsp_role.Data.Role},
	}

	// create product
	productClient := product_proto.NewProductServiceClient("go.micro.srv.product", cl)
	rsp_product, err := productClient.CreateProduct(ctx, &product_proto.CreateProductRequest{Product: product})
	if err != nil {
		t.Error(err)
		return nil
	}
	team.Products = []*product_proto.Product{rsp_product.Data.Product}
	// create service
	rsp_service, err := productClient.CreateService(ctx, &product_proto.CreateServiceRequest{Service: service})
	if err != nil {
		t.Error(err)
		return nil
	}
	team.Services = []*product_proto.Service{rsp_service.Data.Service}

	req := &team_proto.CreateRequest{
		Team:  team,
		OrgId: rsp_org.Data.Organisation.Id,
	}
	rsp := &team_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req, rsp); err != nil {
		t.Error(err)
		return nil
	}

	return rsp.Data.Team
}

func createTeamMember(ctx context.Context, hdlr *TeamService, t *testing.T) *team_proto.CreateTeamMemberResponse_Data {
	team := createTeam(ctx, hdlr, t)
	if team == nil {
		return nil
	}
	// create org
	user.Id = ""
	user.OrgId = "orgid"
	account.Email = "test" + common.Random(4) + "@email.com"

	req := &team_proto.CreateTeamMemberRequest{
		User:    user,
		Account: account,
		Employee: &team_proto.Employee{
			OrgId: "orgid",
			Role:  role,
			Teams: []*team_proto.Team{team},
		},
		OrgId: "orgid",
	}
	rsp := &team_proto.CreateTeamMemberResponse{}
	if err := hdlr.CreateTeamMember(ctx, req, rsp); err != nil {
		t.Error(err)
		return nil
	}
	return rsp.Data
}

func TestAll(t *testing.T) {
	initDb()
	hdlr := new(TeamService)
	ctx := common.NewTestContext(context.TODO())

	team := createTeam(ctx, hdlr, t)
	if team == nil {
		return
	}

	req_all := &team_proto.AllRequest{
		SortParameter: "name",
		SortDirection: "ASC",
	}
	resp_all := &team_proto.AllResponse{}
	err := hdlr.All(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Teams) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_all.Data.Teams[0].Id != team.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_all.Data.Teams)
}

func TestRead(t *testing.T) {
	initDb()
	hdlr := new(TeamService)
	ctx := common.NewTestContext(context.TODO())

	team := createTeam(ctx, hdlr, t)
	if team == nil {
		return
	}

	req_read := &team_proto.ReadRequest{Id: team.Id}
	rsp_read := &team_proto.ReadResponse{}
	err := hdlr.Read(ctx, req_read, rsp_read)
	if err != nil {
		t.Error(err)
		return
	}
	if rsp_read.Data.Team == nil {
		t.Error("Object could not be nil")
		return
	}
	if rsp_read.Data.Team.Id != team.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDelete(t *testing.T) {
	initDb()
	hdlr := new(TeamService)
	ctx := common.NewTestContext(context.TODO())

	team := createTeam(ctx, hdlr, t)
	if team == nil {
		return
	}

	req_del := &team_proto.DeleteRequest{Id: team.Id}
	rsp_del := &team_proto.DeleteResponse{}
	err := hdlr.Delete(ctx, req_del, rsp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &team_proto.ReadRequest{Id: team.Id}
	rsp_read := &team_proto.ReadResponse{}
	err = hdlr.Read(ctx, req_read, rsp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if rsp_read.Data != nil {
		t.Error("Not deleted")
		return
	}
}

func TestFilter(t *testing.T) {
	initDb()
	hdlr := new(TeamService)
	ctx := common.NewTestContext(context.TODO())

	team := createTeam(ctx, hdlr, t)
	if team == nil {
		return
	}
	req_filter := &team_proto.FilterRequest{
		// Product: []string{"product1"},
		OrgId:         team.OrgId,
		Offset:        0,
		Limit:         10,
		SortParameter: "name",
		SortDirection: "ASC",
	}
	rsp_filter := &team_proto.FilterResponse{}

	err := hdlr.Filter(ctx, req_filter, rsp_filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(rsp_filter.Data.Teams) == 0 {
		t.Error("Count does not match")
		return
	}
}

func TestSearch(t *testing.T) {
	initDb()
	hdlr := new(TeamService)
	ctx := common.NewTestContext(context.TODO())

	team := createTeam(ctx, hdlr, t)
	if team == nil {
		return
	}

	req_search := &team_proto.SearchRequest{
		OrgId:         team.OrgId,
		TeamName:      team.Name,
		Offset:        0,
		Limit:         10,
		SortParameter: "name",
		SortDirection: "ASC",
	}
	rsp_search := &team_proto.SearchResponse{}
	err := hdlr.Search(ctx, req_search, rsp_search)
	if err != nil {
		t.Error(err)
		return
	}
	if len(rsp_search.Data.Teams) == 0 {
		t.Error("Count does not match")
		return
	}
}

func TestCreateTeamMemeber(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	rsp := createTeamMember(ctx, hdlr, t)
	if rsp == nil {
		t.Error("Create team-member is failed")
		return
	}
}

func TestAllTeamMember(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	rsp_teammember := createTeamMember(ctx, hdlr, t)
	if rsp_teammember == nil {
		t.Error("Create team-member is failed")
		return
	}

	req := &team_proto.AllTeamMemberRequest{
		OrgId:  "orgid",
		Offset: 0,
		Limit:  10,
	}
	rsp := &team_proto.AllTeamMemberResponse{}
	if err := hdlr.AllTeamMember(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	if len(rsp.Data.Employees) == 0 {
		t.Error("Object count is not matched")
		return
	}
}

func TestReadTeamMemeber(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	rsp_teammember := createTeamMember(ctx, hdlr, t)
	if rsp_teammember == nil {
		t.Error("Create team-member is failed")
		return
	}

	req := &team_proto.ReadTeamMemberRequest{
		UserId: rsp_teammember.User.Id,
	}
	rsp := &team_proto.ReadTeamMemberResponse{}
	err := hdlr.ReadTeamMember(ctx, req, rsp)
	if err != nil {
		t.Error(err)
		return
	}
	if rsp.Data.Employee == nil {
		t.Error("Employee does not matched")
		return
	}
	if rsp.Data.User == nil {
		t.Error("User does not matched")
	}
}

func TestFilterTeamMember(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	rsp_teammember := createTeamMember(ctx, hdlr, t)
	if rsp_teammember == nil {
		return
	}

	req := &team_proto.FilterTeamMemberRequest{
		Team:          []string{rsp_teammember.Employee.Teams[0].Id},
		OrgId:         rsp_teammember.User.OrgId,
		Offset:        0,
		Limit:         10,
		SortParameter: "created",
		SortDirection: "DESC",
	}
	rsp := &team_proto.FilterTeamMemberResponse{}
	if err := hdlr.FilterTeamMember(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	if len(rsp.Data.Employees) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestReadEmployeeInfo(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

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

	req := &team_proto.ReadEmployeeInfoRequest{
		UserId: si.UserId,
	}
	rsp := &team_proto.ReadEmployeeInfoResponse{}
	if err := hdlr.ReadEmployeeInfo(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	if rsp.Employee == nil {
		t.Error("Object count does not matched")
		return
	}

}

func TestPutEmployeeInfo(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	rsp_teammember := createTeamMember(ctx, hdlr, t)
	if rsp_teammember == nil {
		return
	}
	// create org

	req := &team_proto.PutEmployeeInfoRequest{
		UserId: rsp_teammember.User.Id,
		OrgId:  rsp_teammember.User.OrgId,
	}
	rsp := &team_proto.PutEmployeeInfoResponse{}
	if err := hdlr.PutEmployeeInfo(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	if rsp.Data.Employee == nil {
		t.Error("Object count does not matched")
		return
	}

}

func TestCheckValidEmployee(t *testing.T) {
	initDb()

	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
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

	req := &team_proto.ReadEmployeeInfoRequest{
		UserId: si.UserId,
	}
	rsp := &team_proto.CheckValidEmployeeResponse{}
	if err := hdlr.CheckValidEmployee(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	if rsp.Valid == false {
		t.Error("Object count does not matched")
		return
	}

}
