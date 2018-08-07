package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	account_proto "server/account-srv/proto/account"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	organisation_proto "server/organisation-srv/proto/organisation"
	team_proto "server/team-srv/proto/team"
	user_proto "server/user-srv/proto/user"
	"strconv"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/nats"
)

const (
	// internal user_id retrieved from session key
	UserIdAttrName = "user_id"

	// internal org_id retrieved from session key
	OrgIdAttrName = "org_id"

	// internal team_id retrieved from session key
	TeamIdAttrName = "team_id"

	// internal employee_id retrieved from session key
	EmployeeIdAttrName = "employee_id"

	// Power level for the room
	UserPowerLevelAttrName = "power_id"

	// Power level for the organization
	OrgPowerLevelAttrName = "org_power_id"

	// Search username query parameter
	UserNameSearchParameter = "name"

	// Search username query parameter
	UserEmailSearchParameter = "email"

	// session auth parameter
	SessionParameter = "session"

	// Pagination from unix timestemp
	PaginateFromParameter = "from"

	// Pagination to unix timestemp
	PaginateToParameter = "to"

	// Pagination limit parameter
	PaginateLimitParameter = "limit"

	// Pagination offset parameter
	PaginateOffsetParameter = "offset"

	// Filter id parameter
	FilterParameter = "filter"

	// Search query parameter
	SearchParameter = "q"

	// Sort Parameter
	SortParameter = "sort_parameter"

	// Sort Direction
	SortDirection = "sort_direction"
)

type Filters struct {
	AccountClient      account_proto.AccountServiceClient
	TeamClient         team_proto.TeamServiceClient
	UserClient         user_proto.UserServiceClient
	OrganisationClient organisation_proto.OrganisationServiceClient
}

// Reads a session get-parameter and checks user session (if user has logged in)
func (r Filters) BasicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	sessionId := req.QueryParameter(SessionParameter)
	if len(sessionId) == 0 {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		// resp.WriteErrorString(401, "401: Not Authorized")
		utils.NoAuthorizedResponse(resp, errors.New("Not Authorized"), "basic.auth.error", "Invalid Session")
		return
	}
	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	session_resp, err := r.AccountClient.ReadSession(ctx, &account_proto.ReadSessionRequest{sessionId})
	if err != nil {
		// resp.AddHeader("Content-Type", "text/plain")
		// resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		utils.NoAuthorizedResponse(resp, err, "basic.auth.error", "Not Authorized")
		return
	}

	req.SetAttribute(UserIdAttrName, session_resp.UserId)
	req.SetAttribute(OrgIdAttrName, session_resp.OrgId)
	// log.Println(session_resp.UserId, session_resp.OrgId)
	chain.ProcessFilter(req, resp)
}

// // Returns closure with given action that checks the power level of the action.
// func (r Filters) RoomPowerLevel(action string) func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
// 	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
// 		powerlevel := req.Attribute(UserPowerLevelAttrName).(int64)
// 		ok, err := common.EnoughPowerLevel(action, powerlevel)
// 		if !ok && err != nil {
// 			resp.AddHeader("Content-Type", "text/plain")
// 			resp.WriteErrorString(http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		chain.ProcessFilter(req, resp)
// 	}
// }

// // Returns closure with given action that checks the power level of the action.
// func (r Filters) OrgPowerLevel(action string) func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
// 	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
// 		powerlevel := req.Attribute(OrgPowerLevelAttrName).(int64)
// 		ok, err := common.EnoughPowerLevel(action, powerlevel)
// 		if !ok && err != nil {
// 			resp.AddHeader("Content-Type", "text/plain")
// 			resp.WriteErrorString(http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		chain.ProcessFilter(req, resp)
// 	}
// }

// Gets user power level within the organisation
func (r Filters) OrganisationAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)

	orgId := req.Attribute(OrgIdAttrName).(string)
	if len(orgId) == 0 {
		utils.NoAuthorizedResponse(resp, errors.New("Not Organisation Authorized"), "org.auth.error", "Invalid Organisation")
		return
	}
	resp_org, err := r.OrganisationClient.ReadOrgInfo(ctx, &organisation_proto.ReadOrgInfoRequest{orgId})
	if err != nil || resp_org.OrgInfo == nil {
		utils.NoAuthorizedResponse(resp, errors.New("Not Organisation Authorized"), "org.auth.error", "Not Organisation Authorized")
		return
	}

	// will add power level later
	chain.ProcessFilter(req, resp)
}

// Gets user power level within the organization
// func (r Filters) OrganizationAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
// 	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)

// 	req_search := new(organisation_proto.SearchMembersRequest)
// 	req_search.Userid = req.Attribute(UserIdAttrName).(string)
// 	req_search.Orgid = req.Attribute(OrgIdAttrName).(string)
// 	resp_search, err := r.OrganisationClient.SearchMembers(ctx, req_search)
// 	if err != nil {
// 		resp.AddHeader("Content-Type", "text/plain")
// 		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	if len(resp_search.Members) > 0 {
// 		userlevel := resp_search.Members[0].Level

// 		req_search_team_member := new(organisation_proto.SearchTeamMembersRequest)
// 		req_search_team_member.Userid = req.Attribute(UserIdAttrName).(string)
// 		req_search_team_member.Orgid = req.Attribute(OrgIdAttrName).(string)
// 		resp_search_team_member, err := r.OrganisationClient.SearchTeamMembers(ctx, req_search_team_member)
// 		if err == nil && len(resp_search_team_member.TeamMembers) != 0 {
// 			req_read_team := new(organisation_proto.ReadTeamRequest)
// 			req_read_team.Id = resp_search_team_member.TeamMembers[0].Teamid
// 			resp_read_team, err := r.OrganisationClient.ReadTeam(ctx, req_read_team)
// 			if err != nil {
// 				if userlevel < resp_read_team.Team.Level {
// 					userlevel = resp_read_team.Team.Level
// 				}

// 			}
// 		}

// 		req.SetAttribute(OrgPowerLevelAttrName, userlevel)
// 	} else {
// 		req.SetAttribute(OrgPowerLevelAttrName, int64(common.PowerlevelAdmin))
// 	}

// 	chain.ProcessFilter(req, resp)
// }

// Check whether this is a valid employee
func (r Filters) EmployeeAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)

	//ReadEmployeeInfo - from kv
	resp_employee, err := r.TeamClient.ReadEmployeeInfo(ctx, &team_proto.ReadEmployeeInfoRequest{req.Attribute(UserIdAttrName).(string)})
	if err != nil || resp_employee.Employee.Employee == nil {
		utils.NoAuthorizedResponse(resp, errors.New("Not Authorized"), "employee.auth.error", "Not Authorized")
		return
	}

	// var e team_proto.Employee
	// if err := jsonpb.Unmarshal(strings.NewReader(resp_employee.Employee.Employee), &e); err != nil {
	// 	utils.NoAuthorizedResponse(resp, err, "employee.auth.error", "Not Authorized")
	// }
	req.SetAttribute(TeamIdAttrName, resp_employee.Employee.Employee.User.Id)

	chain.ProcessFilter(req, resp)
}

// Helper to extract int from string parameters
func parameterToValue(parameter string) int64 {
	if len(parameter) != 0 {
		val, err := strconv.Atoi(parameter)
		if err == nil && val > 0 {
			return int64(val)
		}
	}
	return 0
}

// Cuts pagination parameters: from, to, limit
func (r Filters) Paginate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	from := parameterToValue(req.QueryParameter(PaginateFromParameter))
	to := parameterToValue(req.QueryParameter(PaginateToParameter))
	limit := parameterToValue(req.QueryParameter(PaginateLimitParameter))
	offset := parameterToValue(req.QueryParameter(PaginateOffsetParameter))

	req.SetAttribute(PaginateFromParameter, from)
	req.SetAttribute(PaginateToParameter, to)
	req.SetAttribute(PaginateLimitParameter, limit)
	req.SetAttribute(PaginateOffsetParameter, offset)
	chain.ProcessFilter(req, resp)
}

// Cuts SortFilter parameters: sortParameter, sortDirection
func (r Filters) SortFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	sortParameter := req.QueryParameter(SortParameter)
	sortDirection := req.QueryParameter(SortDirection)

	req.SetAttribute(SortParameter, sortParameter)
	req.SetAttribute(SortDirection, sortDirection)
	chain.ProcessFilter(req, resp)
}

// // Decrypted key request
// func decryptedDataKey(orgid string, req *restful.Request, cloud_cl cloudkey_proto.CloudKeyServiceClient, org_cl organisation_proto.OrganisationServiceClient) (string, error) {
// 	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
// 	org_send, err := org_cl.Read(ctx, &organisation_proto.ReadRequest{
// 		Id: orgid,
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	cloudkey_send, err := cloud_cl.DecryptKey(ctx, &cloudkey_proto.DecryptKeyRequest{
// 		Orgid: orgid,
// 		Dek:   org_send.Organization.DataEncryptionKey,
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	return cloudkey_send.EncryptedDek, err
// }

type AuditFilter struct {
	Broker broker.Broker
	Audit  *audit_proto.Audit
}

func (r AuditFilter) Clone(audit *audit_proto.Audit) *AuditFilter {
	dest := &audit_proto.Audit{}
	CloneValue(audit, dest)
	clone := r
	clone.Audit = dest
	return &clone
}

// AuditFilter log all requests
func (r AuditFilter) AuditFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// get paramters
	data, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		return
	}
	audit := r.Audit
	audit.ActionName = req.Request.URL.Path
	audit.ActionTimestamp = time.Now().Unix()
	audit.ActionParameters = string(data)
	audit.ActionMethod = req.Request.Method
	audit.ActionMetaData = req.Request.RemoteAddr
	fmt.Println("process1")

	go func() {
		if body, err := json.Marshal(audit); err == nil {
			if err := r.Broker.Publish(common.AUDIT_ACTION, &broker.Message{Body: body}); err != nil {
				fmt.Println("audit_log publish is failed")
				return
			}
		}
	}()
	fmt.Println("process2")
	chain.ProcessFilter(req, resp)
}

func CloneValue(source interface{}, destin interface{}) {
	x := reflect.ValueOf(source)
	if x.Kind() == reflect.Ptr {
		starX := x.Elem()
		y := reflect.New(starX.Type())
		starY := y.Elem()
		starY.Set(starX)
		reflect.ValueOf(destin).Elem().Set(y.Elem())
	} else {
		destin = x.Interface()
	}
}
