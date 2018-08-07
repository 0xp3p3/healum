package handler

import (
	"context"
	"fmt"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	organisation_proto "server/organisation-srv/proto/organisation"
	static_proto "server/static-srv/proto/static"
	"server/team-srv/db"
	team_proto "server/team-srv/proto/team"
	user_proto "server/user-srv/proto/user"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type TeamService struct {
	AccountClient      account_proto.AccountServiceClient
	UserClient         user_proto.UserServiceClient
	KvClient           kv_proto.KvServiceClient
	OrganisationClient organisation_proto.OrganisationServiceClient
}

func (p *TeamService) All(ctx context.Context, req *team_proto.AllRequest, rsp *team_proto.AllResponse) error {
	log.Info("Received Team.All request")
	teams, err := db.All(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, "", "")
	if len(teams) == 0 || err != nil {
		return common.NotFound(common.TeamSrv, p.All, err, "not found")
	}
	rsp.Data = &team_proto.ArrData{teams}
	return nil
}

func (p *TeamService) Create(ctx context.Context, req *team_proto.CreateRequest, rsp *team_proto.CreateResponse) error {
	log.Info("Received Team.Create request")
	if len(req.Team.Name) == 0 {
		return common.BadRequest(common.TeamSrv, p.Create, nil, "team name empty")
	}
	if len(req.Team.Id) == 0 {
		req.Team.Id = uuid.NewUUID().String()
	}
	err := db.Create(ctx, req.Team)
	if err != nil {
		return common.InternalServerError(common.TeamSrv, p.Create, err, "create error")
	}
	rsp.Data = &team_proto.Data{req.Team}
	return nil
}

func (p *TeamService) Read(ctx context.Context, req *team_proto.ReadRequest, rsp *team_proto.ReadResponse) error {
	log.Info("Received Team.Read request")
	team, err := db.Read(ctx, req.Id, req.OrgId, req.TeamId)
	if team == nil || err != nil {
		return common.NotFound(common.TeamSrv, p.Read, err, "not found")
	}
	rsp.Data = &team_proto.Data{team}
	return nil
}

func (p *TeamService) Delete(ctx context.Context, req *team_proto.DeleteRequest, rsp *team_proto.DeleteResponse) error {
	log.Info("Received Team.Delete request")
	if err := db.Delete(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.TeamSrv, p.Delete, err, "delete error")
	}
	return nil
}

func (p *TeamService) Filter(ctx context.Context, req *team_proto.FilterRequest, rsp *team_proto.FilterResponse) error {
	log.Info("Received Team.Filter request")
	teams, err := db.Filter(ctx, req)
	if len(teams) == 0 || err != nil {
		return common.NotFound(common.TeamSrv, p.Filter, err, "not found")
	}
	rsp.Data = &team_proto.ArrData{teams}
	return nil
}

func (p *TeamService) Search(ctx context.Context, req *team_proto.SearchRequest, rsp *team_proto.SearchResponse) error {
	log.Info("Received Team.Search request")
	teams, err := db.Search(ctx, req)
	if len(teams) == 0 || err != nil {
		return common.NotFound(common.TeamSrv, p.Search, err, "not found")
	}
	rsp.Data = &team_proto.ArrData{teams}
	return nil
}

func (p *TeamService) AllTeamMember(ctx context.Context, req *team_proto.AllTeamMemberRequest, rsp *team_proto.AllTeamMemberResponse) error {
	log.Info("Received Team.AllTeamMember request")
	members, err := db.AllTeamMember(ctx, req.OrgId, "", req.Offset, req.Limit, "", "")
	if len(members) == 0 || err != nil {
		return common.NotFound(common.TeamSrv, p.AllTeamMember, err, "not found")
	}
	rsp.Data = &team_proto.AllTeamMemberResponse_Data{members}
	return nil
}

func (p *TeamService) CreateTeamMember(ctx context.Context, req *team_proto.CreateTeamMemberRequest, rsp *team_proto.CreateTeamMemberResponse) error {
	log.Info("Received Team.CreateTeamMember request")

	// check validation to create employee
	if req.User.OrgId != req.OrgId {
		return common.BadRequest(common.TeamSrv, p.CreateTeamMember, nil, "user's orgid invalid")

	}
	if req.Employee.OrgId != req.OrgId {
		return common.BadRequest(common.TeamSrv, p.CreateTeamMember, nil, "employee's orgid invalid")
	}

	//create user and account in user_srv
	//create relationship between user and account
	if req.User != nil {
		req_user := &user_proto.CreateRequest{User: req.User, Account: req.Account, TeamId: req.TeamId}
		rsp_user, err := p.UserClient.Create(ctx, req_user)
		if err != nil {
			return common.InternalServerError(common.TeamSrv, p.CreateTeamMember, err, "create error")
		}
		req.User = rsp_user.Data.User
		req.Account = rsp_user.Data.Account

		fmt.Println("req.User:", rsp_user.Data.User)
	}

	//create employee edge
	req_employee := &team_proto.CreateEmployeeEdgeRequest{
		Employee: req.Employee,
		OrgId:    req.Employee.OrgId,
		UserId:   req.User.Id,
	}
	rsp_employee := &team_proto.CreateEmployeeEdgeResponse{}
	if err := p.CreateEmployeeEdge(ctx, req_employee, rsp_employee); err != nil {
		return common.InternalServerError(common.TeamSrv, p.CreateTeamMember, err, "create employee edge error")
	}
	req.Employee = rsp_employee.Data.Employee

	// create employee profile
	profile := req.Employee.Profile
	if profile != nil {
		req_profile := &team_proto.CreateEmployeeProfileRequest{profile}
		rsp_profile := &team_proto.CreateEmployeeProfileResponse{}
		if err := p.CreateEmployeeProfile(ctx, req_profile, rsp_profile); err != nil {
			return common.InternalServerError(common.TeamSrv, p.CreateTeamMember, err, "create employee profile error")
		}
	}

	//create module access for this employee
	if len(req.Employee.Modules) > 0 {
		req_employee_modules := &team_proto.CreateEmployeeModuleAccessRequest{UserId: req.User.Id, OrgId: req.OrgId, Modules: req.Employee.Modules}
		if err := p.CreateEmployeeModuleAccess(ctx, req_employee_modules, &team_proto.CreateEmployeeModuleAccessResponse{}); err != nil {
			common.ErrorLog(common.TeamSrv, "Create", err, "CreateEmployeeModuleAccess failed")
		}
	}

	//create team membership
	if len(req.Employee.Teams) > 0 {
		req_membership := &team_proto.CreateTeamMembershipRequest{
			Employee: req.Employee,
			User:     req.User,
		}
		rsp_memebership := &team_proto.CreateTeamMembershipResponse{}
		if err := p.CreateTeamMembership(ctx, req_membership, rsp_memebership); err != nil {
			return common.InternalServerError(common.TeamSrv, p.CreateTeamMember, err, "create team membership error")
		}
	}

	rsp.Data = &team_proto.CreateTeamMemberResponse_Data{
		Employee: req.Employee,
		User:     req.User,
		Account:  req.Account,
	}
	return nil
}

func (p *TeamService) ReadTeamMember(ctx context.Context, req *team_proto.ReadTeamMemberRequest, rsp *team_proto.ReadTeamMemberResponse) error {
	log.Info("Received Team.ReadTeamMember request")
	user, employee, err := db.ReadTeamMember(ctx, req.UserId, req.OrgId, req.TeamId)

	if employee == nil || err != nil {
		return common.NotFound(common.TeamSrv, p.ReadTeamMember, err, "not found")
	}
	rsp.Data = &team_proto.ReadTeamMemberResponse_Data{User: user, Employee: employee}
	return nil
}

func (p *TeamService) FilterTeamMember(ctx context.Context, req *team_proto.FilterTeamMemberRequest, rsp *team_proto.FilterTeamMemberResponse) error {
	log.Info("Received Team.FilterTeamMember request")
	employees, err := db.FilterTeamMember(ctx, req)
	if len(employees) == 0 || err != nil {
		return common.NotFound(common.TeamSrv, p.FilterTeamMember, err, "not found")
	}
	rsp.Data = &team_proto.FilterTeamMemberResponse_Data{employees}
	return nil
}

func (p *TeamService) CreateEmployeeEdge(ctx context.Context, req *team_proto.CreateEmployeeEdgeRequest, rsp *team_proto.CreateEmployeeEdgeResponse) error {
	log.Info("Received Team.CreateEmployeeEdge request")
	err := db.CreateEmployeeEdge(ctx, req.Employee, req.UserId, req.OrgId)
	if err != nil {
		return common.InternalServerError(common.TeamSrv, p.CreateEmployeeEdge, err, "create employee edge error")
	}
	rsp.Data = &team_proto.CreateEmployeeEdgeResponse_Data{req.Employee}
	return nil
}

func (p *TeamService) CreateEmployeeProfile(ctx context.Context, req *team_proto.CreateEmployeeProfileRequest, rsp *team_proto.CreateEmployeeProfileResponse) error {
	log.Info("Received Team.CreateEmployeeProfile request")
	err := db.CreateEmployeeProfile(ctx, req.Profile)
	if err != nil {
		return common.InternalServerError(common.TeamSrv, p.CreateEmployeeProfile, err, "create employee profile error")
	}
	rsp.Data = &team_proto.CreateEmployeeProfileResponse_Data{req.Profile}
	return nil
}

func (p *TeamService) CreateTeamMembership(ctx context.Context, req *team_proto.CreateTeamMembershipRequest, rsp *team_proto.CreateTeamMembershipResponse) error {
	log.Info("Received Team.CreateEmployeeProfile request")
	err := db.CreateTeamMembership(ctx, req.Employee, req.User)
	if err != nil {
		return common.InternalServerError(common.TeamSrv, p.CreateTeamMembership, err, "create team membership error")
	}
	return nil
}

func (p *TeamService) Update(ctx context.Context, req *team_proto.UpdateRequest, rsp *team_proto.UpdateResponse) error {
	log.Info("Received Team.Update request")
	return nil
}
func (p *TeamService) CreateEmployeeModuleAccess(ctx context.Context, req *team_proto.CreateEmployeeModuleAccessRequest, rsp *team_proto.CreateEmployeeModuleAccessResponse) error {
	log.Info("Received Team.CreateEmployeeModuleAccess request for user: ", req.UserId)
	//get org modules
	resp_org, err := p.OrganisationClient.ReadOrgInfo(ctx, &organisation_proto.ReadOrgInfoRequest{req.OrgId})
	if err != nil {
		//log error here
	}
	em := []*static_proto.Module{}

	//check if e.modules part of org modules
	//FIXME:Improve the loops
	for _, m := range req.Modules {
		for _, om := range resp_org.OrgInfo.Modules {
			if m.Id == om.Id {
				em = append(em, m)
			}
		}
	}

	// give only access to org modules by creating employee module edge
	if err := db.CreateEmpoyeeModuleAccess(ctx, req.UserId, req.OrgId, em); err != nil {
		return err
	}
	return nil
}

func (p *TeamService) GetAccessibleModulesByEmployee(ctx context.Context, req *team_proto.GetAccessibleModulesByEmployeeRequest, rsp *team_proto.GetAccessibleModulesByEmployeeResponse) error {
	log.Info("Received Team.GetAccessibleModulesByEmployee request for user: ", req.UserId)
	modules, err := db.GetAccessibleModulesByEmployee(ctx, req.UserId, req.OrgId)
	if err != nil {
		return err
	}
	rsp.Data = &team_proto.GetAccessibleModulesByEmployeeResponse_Data{Modules: modules}
	return nil
}

//Read EmployeeInfo for EmployeeAuthenticate and validate shared_by user (the employee)
func (p *TeamService) ReadEmployeeInfo(ctx context.Context, req *team_proto.ReadEmployeeInfoRequest, rsp *team_proto.ReadEmployeeInfoResponse) error {
	log.Info("Received Team.ReadEmployeeInfo request")
	req_kv := &kv_proto.GetExRequest{common.EMPLOYEE_INFO_INDEX, req.UserId}
	rsp_kv, err := p.KvClient.GetEx(ctx, req_kv)
	if err != nil {
		return common.NotFound(common.TeamSrv, p.ReadEmployeeInfo, err, "GetEx is failed")
	}

	var ei team_proto.EmployeeInfo
	if err := jsonpb.Unmarshal(strings.NewReader(string(rsp_kv.Item.Value)), &ei); err != nil {
		return common.InternalServerError(common.TeamSrv, p.ReadEmployeeInfo, err, "Unmarshaller error")
	}
	rsp.Employee = &ei
	return nil
}

//store employee information in KV for in-memory checks using ReadEmployeeInfo
func (p *TeamService) PutEmployeeInfo(ctx context.Context, req *team_proto.PutEmployeeInfoRequest, rsp *team_proto.PutEmployeeInfoResponse) error {
	log.Info("Received Team.PutEmployeeInfo request")

	// get employee from team client
	log.Info("Read Employee details for user: ", req.UserId)
	_, employee, err := db.ReadTeamMember(ctx, req.UserId, req.OrgId, "")
	if err != nil {
		return common.NotFound(common.TeamSrv, p.PutEmployeeInfo, err, "ReadTeamMember query is failed")
	}

	//only store into redis if the user is employee
	if employee != nil {
		// create employee info
		ei := &team_proto.EmployeeInfo{
			UserId:   req.UserId,
			OrgId:    req.OrgId,
			Employee: employee,
		}
		//FIXME:Using EnumAsInts for now as default of EnumAsString is not working for jsonpb unmarshal - 01/07/2018
		marshaler := jsonpb.Marshaler{EmitDefaults: true, EnumsAsInts: true}
		employee_js, err := marshaler.MarshalToString(ei)
		if err != nil {
			return common.InternalServerError(common.TeamSrv, p.PutEmployeeInfo, err, "Marshalling error")
		}

		req_kv := &kv_proto.PutExRequest{
			Index: common.EMPLOYEE_INFO_INDEX,
			Item: &kv_proto.Item{
				Key:   employee.User.Id,
				Value: []byte(employee_js),
			},
		}
		if _, err := p.KvClient.PutEx(context.TODO(), req_kv); err != nil {
			return common.InternalServerError(common.TeamSrv, p.PutEmployeeInfo, err, "PutEx is failed")
		}
	}
	rsp.Data = &team_proto.PutEmployeeInfoResponse_Data{employee}
	return nil
}

//Read EmployeeInfo and return whether it's a valid employee or not
func (p *TeamService) CheckValidEmployee(ctx context.Context, req *team_proto.ReadEmployeeInfoRequest, rsp *team_proto.CheckValidEmployeeResponse) error {
	log.Info("Received Team.CheckValidEmployee request")

	// checking valid sharedby
	req_employee := &team_proto.ReadEmployeeInfoRequest{req.UserId}
	rsp_employee := &team_proto.ReadEmployeeInfoResponse{}
	err := p.ReadEmployeeInfo(ctx, req_employee, rsp_employee)
	if err != nil {
		return common.NotFound(common.TeamSrv, p.CheckValidEmployee, err, "ReadEmployeeInfo is failed")
	}

	// check shareWith
	if len(rsp_employee.Employee.UserId) > 0 && len(rsp_employee.Employee.OrgId) > 0 {
		// check shareBy
		if rsp_employee.Employee.Employee != nil {
			rsp.Valid = true
			rsp.Employee = rsp_employee.Employee.Employee
			return nil
		} else {
			rsp.Valid = false
			return nil
		}
	}
	return nil
}

func (p *TeamService) DeleteEmployee(ctx context.Context, req *team_proto.DeleteEmployeeRequest, rsp *team_proto.DeleteEmployeeResponse) error {
	log.Info("Received Team.DeleteEmployee request")

	req_employee := &team_proto.ReadEmployeeInfoRequest{req.UserId}
	rsp_employee := &team_proto.CheckValidEmployeeResponse{}
	err := p.CheckValidEmployee(ctx, req_employee, rsp_employee)
	if err != nil {
		return common.InternalServerError(common.ContentSrv, p.DeleteEmployee, err, "CheckValidEmployee is failed")
	}

	if rsp_employee.Valid && rsp_employee.Employee != nil {
		err := db.DeleteEmployee(ctx, req.EmployeeId, req.OrgId)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.DeleteEmployee, err, "delete error")
		}
	}
	return nil
}
