package handler

import (
	"context"
	"encoding/json"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	"server/organisation-srv/db"
	organisation_proto "server/organisation-srv/proto/organisation"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
	user_proto "server/user-srv/proto/user"
	"strings"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/nats"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type OrganisationService struct {
	Broker        broker.Broker
	KvClient      kv_proto.KvServiceClient
	AccountClient account_proto.AccountServiceClient
	TeamClient    team_proto.TeamServiceClient
	StaticClient  static_proto.StaticServiceClient
	UserClient    user_proto.UserServiceClient
}

func (p *OrganisationService) All(ctx context.Context, req *organisation_proto.AllRequest, rsp *organisation_proto.AllResponse) error {
	log.Info("Received Organisation.All request")
	orgs, err := db.All(ctx, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(orgs) == 0 || err != nil {
		return common.NotFound(common.OrganisationSrv, p.All, err, "organisation not found")
	}
	rsp.Data = &organisation_proto.AllResponse_Data{orgs}
	return nil
}

func (p *OrganisationService) Create(ctx context.Context, req *organisation_proto.CreateRequest, rsp *organisation_proto.CreateResponse) error {
	log.Info("Received Organisation.Create request")

	// porting one variable
	org := req.Organisation
	user := req.User
	account := req.Account

	if len(org.Id) == 0 {
		org.Id = uuid.NewUUID().String()
	} else if org.Id != user.OrgId {
		// check validation orgid of user
		return common.InternalServerError(common.OrganisationSrv, p.Create, nil, "match error")
	}

	// create org
	// create module edge table entry between organisation and module
	if err := db.Create(ctx, org, req.Modules); err != nil {
		return common.InternalServerError(common.OrganisationSrv, p.Create, err, "create error")
	}
	// update orgid for user
	user.OrgId = org.Id

	// get admin role
	rsp_role, err := p.StaticClient.ReadRoleByNameslug(ctx, &static_proto.ReadRoleByNameslugRequest{"admin"})
	if rsp_role == nil || err != nil {
		return common.NotFound(common.OrganisationSrv, p.Create, err, "role not found")
	}
	// Create a default team named Admin team for this organisation.
	teams := []*team_proto.Team{}
	// make team object
	req_team := &team_proto.CreateRequest{
		Team: &team_proto.Team{
			Name:  "admin",
			Id:    uuid.NewUUID().String(),
			OrgId: org.Id,
		},
	}
	rsp_team, err := p.TeamClient.Create(ctx, req_team) // call team create
	if err != nil {
		return common.InternalServerError(common.OrganisationSrv, p.Create, err, "team client create error")
	}
	teams = append(teams, rsp_team.Data.Team)

	// create team-member
	// create Employee edge table entry with the Organisation
	// owner by default gets access to all modules of an organisation
	employee := &team_proto.Employee{
		OrgId:   org.Id,
		Role:    rsp_role.Data.Role,
		Teams:   teams,
		Modules: req.Modules,
	}
	req_member := &team_proto.CreateTeamMemberRequest{
		Employee: employee,
		User:     user,
		Account:  account,
		OrgId:    org.Id,
	}
	rsp_member, err := p.TeamClient.CreateTeamMember(ctx, req_member)
	if err != nil {
		return common.InternalServerError(common.OrganisationSrv, p.Create, err, "team member create error")
	}
	user = rsp_member.Data.User
	account = rsp_member.Data.Account
	// remove sensitive data from account
	if account != nil {
		account.Passcode = ""
		account.Password = ""
	}

	// update org with create user by owner
	org.Owner = user
	if err := db.Create(ctx, org, req.Modules); err != nil {
		return common.InternalServerError(common.OrganisationSrv, p.Create, err, "team owner create error")
	}

	// put organisation info in REDIS
	go p.readAndPutOrg(ctx, org.Id)

	// remove sensitive data from user
	user = &user_proto.User{
		Id:        user.Id,
		OrgId:     user.OrgId,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		AvatarUrl: user.AvatarUrl,
		Gender:    user.Gender,
	}

	rsp.Data = &organisation_proto.CreateResponse_Data{
		Organisation: &organisation_proto.Organisation{Id: org.Id, Name: org.Name},
		Account:      &account_proto.Account{Id: account.Id},
		User:         &user_proto.User{Id: user.Id},
	}

	// publish organisation
	if body, err := json.Marshal(rsp.Data); err == nil {
		if err := p.Broker.Publish(common.ORGANISATION_CREATED, &broker.Message{Body: body}); err != nil {
			return common.InternalServerError(common.OrganisationSrv, p.Create, err, "subscribe error")
		}
	}

	return nil
}

func (p *OrganisationService) readAndPutOrg(ctx context.Context, orgId string) {
	log.Info("Received Organisation.readAndPutOrg request for Organisation: ", orgId)
	rsp_org, _ := db.Read(ctx, orgId)
	oi := &organisation_proto.OrgInfo{
		OrgId:   rsp_org.Id,
		Type:    rsp_org.Type,
		Owner:   rsp_org.Owner,
		Modules: rsp_org.Modules,
	}
	req_orginfo := &organisation_proto.PutOrgInfoRequest{OrgId: orgId, OrgInfo: oi}
	rsp_orginfo := &organisation_proto.PutOrgInfoResponse{}
	p.PutOrgInfo(context.TODO(), req_orginfo, rsp_orginfo)
}

//update employee modules when an organisation module access is removed
func (p *OrganisationService) updateAllEmployeesModuleAccess(ctx context.Context, orgId string) {
	log.Info("Received Organisation.updateEmployeesModuleAccess request for Organisation: ", orgId)
	// get all employes
	rsp_all_employees, err := p.TeamClient.AllTeamMember(ctx, &team_proto.AllTeamMemberRequest{OrgId: orgId, Limit: 1000})
	if err != nil {
		common.ErrorLog(common.OrganisationSrv, common.GetFunctionName(p.readAndPutOrg), err, "All Employee query is failed")
	}
	// get their existing modules
	for _, employee := range rsp_all_employees.Data.Employees {
		// create module accesss
		p.TeamClient.CreateEmployeeModuleAccess(ctx, &team_proto.CreateEmployeeModuleAccessRequest{UserId: employee.User.Id, Modules: employee.Modules})
	}
}

func (p *OrganisationService) Read(ctx context.Context, req *organisation_proto.ReadRequest, rsp *organisation_proto.ReadResponse) error {
	log.Info("Received Organisation.Read request")
	org, err := db.Read(ctx, req.OrgId)
	if org == nil || err != nil {
		return common.NotFound(common.OrganisationSrv, p.Read, err, "organization not found")
	}
	rsp.Data = &organisation_proto.OrgData{org}
	return nil
}

func (p *OrganisationService) CreateOrganisationProfile(ctx context.Context, req *organisation_proto.CreateOrganisationProfileRequest, rsp *organisation_proto.CreateOrganisationProfileResponse) error {
	log.Info("Received Organisation.CreateOrganisationProfile request")
	err := db.CreateOrganisationProfile(ctx, req.Profile)
	if err != nil {
		return common.InternalServerError(common.OrganisationSrv, p.CreateOrganisationProfile, err, "create organization profile error")
	}
	rsp.Data = &organisation_proto.CreateOrganisationProfileResponse_Data{req.Profile}
	return nil
}

func (p *OrganisationService) CreateOrganisationSetting(ctx context.Context, req *organisation_proto.CreateOrganisationSettingRequest, rsp *organisation_proto.CreateOrganisationSettingResponse) error {
	log.Info("Received Organisation.CreateOrganisationSetting request")
	err := db.CreateOrganisationSetting(ctx, req.Setting)
	if err != nil {
		return common.InternalServerError(common.OrganisationSrv, p.CreateOrganisationSetting, err, "create organization setting error")
	}
	rsp.Data = &organisation_proto.CreateOrganisationSettingResponse_Data{req.Setting}
	return nil
}

func (p *OrganisationService) ReadOrgInfo(ctx context.Context, req *organisation_proto.ReadOrgInfoRequest, rsp *organisation_proto.ReadOrgInfoResponse) error {
	log.Info("Received Organisation.ReadOrgInfo request")

	req_kv := &kv_proto.GetExRequest{common.ORG_INFO_INDEX, req.OrgId}
	rsp_kv, err := p.KvClient.GetEx(ctx, req_kv)
	if err != nil {
		return common.NotFound(common.OrganisationSrv, p.ReadOrgInfo, err, "organization info not found")
	}

	oi := &organisation_proto.OrgInfo{}
	if err := jsonpb.Unmarshal(strings.NewReader(string(rsp_kv.Item.Value)), oi); err != nil {
		return common.InternalServerError(common.OrganisationSrv, p.ReadOrgInfo, err, "parsing error")
	}
	rsp.OrgInfo = oi
	return nil
}

func (p *OrganisationService) PutOrgInfo(ctx context.Context, req *organisation_proto.PutOrgInfoRequest, rsp *organisation_proto.PutOrgInfoResponse) error {
	log.Info("Received Organisation.PutOrgInfo request")

	marshaler := jsonpb.Marshaler{EmitDefaults: true}
	oi_js, err := marshaler.MarshalToString(req.OrgInfo)
	if err != nil {
		return common.NotFound(common.OrganisationSrv, p.PutOrgInfo, err, "organization info not found")
	}
	req_kv := &kv_proto.PutExRequest{
		Index: common.ORG_INFO_INDEX,
		Item: &kv_proto.Item{
			Key:   req.OrgId,
			Value: []byte(oi_js),
		},
	}
	if _, err := p.KvClient.PutEx(context.TODO(), req_kv); err != nil {
		return common.InternalServerError(common.OrganisationSrv, p.PutOrgInfo, err, "organization update error")
	}
	return nil
}

func (p *OrganisationService) WarmupCacheOrganisation(ctx context.Context, req *organisation_proto.WarmupCacheOrganisationRequest, rsp *organisation_proto.WarmupCacheOrganisationResponse) error {
	log.Info("Received Organisation.WarmupCacheOrganisation request")

	//
	var offset int64
	var limit int64
	offset = 0
	limit = 100
	//the emptly for loop is automatically loop through offset and limit until no more orgs are found
	for {
		orgs, err := db.All(ctx, offset, limit, "", "")
		// throws error if there was error
		if err != nil {
			common.NotFound(common.OrganisationSrv, p.WarmupCacheOrganisation, err, "All organisation query is failed")
			break
		}
		// if no more organisations to be fetched from db, then break
		if len(orgs) == 0 {
			log.Info("All organisations fetched to warmupcache")
			break
		}
		for _, org := range orgs {
			// create org_info
			oi := &organisation_proto.OrgInfo{
				OrgId:   org.Id,
				Type:    org.Type,
				Owner:   org.Owner,
				Modules: org.Modules,
			}
			req_put := &organisation_proto.PutOrgInfoRequest{OrgId: org.Id, OrgInfo: oi}
			rsp_put := &organisation_proto.PutOrgInfoResponse{}
			if err := p.PutOrgInfo(ctx, req_put, rsp_put); err != nil {
				continue
			}
		}
		offset += limit
	}
	return nil
}

func (p *OrganisationService) UpdateModules(ctx context.Context, req *organisation_proto.UpdateModulesRequest, rsp *organisation_proto.UpdateModulesResponse) error {
	log.Info("Received Organisation.UpdateModules request")
	var wg sync.WaitGroup
	err := db.UpdateModules(ctx, req.OrgId, req.Modules)
	if err != nil {
		return common.InternalServerError(common.OrganisationSrv, p.UpdateModules, err, "module update error")
	}
	rsp.Data = &organisation_proto.UpdateModulesResponse_Data{req.Modules}
	wg.Add(1)
	go func() {
		log.Info("Updating Organisation Information in KV")
		p.readAndPutOrg(ctx, req.OrgId)
		wg.Done()
		log.Info("Updating Organisation Information in KV Completed")
	}()
	wg.Wait()

	log.Info("Update module access for all employees")
	//update module access for all employees
	go p.updateAllEmployeesModuleAccess(ctx, req.OrgId)

	return nil
}

func (p *OrganisationService) GetModulesByOrg(ctx context.Context, req *organisation_proto.GetModulesByOrgRequest, rsp *organisation_proto.GetModulesByOrgResponse) error {
	log.Info("Received Organisation.GetModules request")

	modules, err := db.GetModulesByOrg(ctx, req.OrgId)
	if len(modules) == 0 || err != nil {
		return common.BadRequest(common.OrganisationSrv, p.GetModulesByOrg, err, "module not found")
	}
	rsp.Data = &organisation_proto.GetModulesByOrgResponse_Data{modules}
	return nil
}

func (p *OrganisationService) Update(ctx context.Context, req *organisation_proto.UpdateRequest, rsp *organisation_proto.UpdateResponse) error {
	log.Info("Received Organisation.Update request")

	organisation, err := db.Update(ctx, req.Organisation)
	if organisation == nil || err != nil {
		return common.InternalServerError(common.OrganisationSrv, p.Update, err, "server error")
	}
	rsp.Data = &organisation_proto.UpdateResponse_Data{Organisation: organisation}
	return nil
}
