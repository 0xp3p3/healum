package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"server/common"
	db_proto "server/db-srv/proto/db"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
	user_proto "server/user-srv/proto/user"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/client"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type clientWrapper struct {
	Db_client db_proto.DBClient
}

var (
	ClientWrapper *clientWrapper
	ErrNotFound   = errors.New("not found")
)

// Storage for a db microservice client
func NewClientWrapper(serviceClient client.Client) *clientWrapper {
	cl := db_proto.NewDBClient("", serviceClient)

	return &clientWrapper{
		Db_client: cl,
	}
}

// Init initializes healum databases
func Init(serviceClient client.Client) error {
	ClientWrapper = NewClientWrapper(serviceClient)
	// if _, err := ClientWrapper.Db_client.Init(context.TODO(), &db_proto.InitRequest{}); err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	return nil
}

// RemoveDb removes healum database (for testing)
func RemoveDb(ctx context.Context, serviceClient client.Client) error {
	ClientWrapper = NewClientWrapper(serviceClient)
	if _, err := ClientWrapper.Db_client.RemoveDb(ctx, &db_proto.RemoveDbRequest{}); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func runQuery(ctx context.Context, q string, table string) (*db_proto.RunQueryResponse, error) {
	return ClientWrapper.Db_client.RunQuery(ctx, &db_proto.RunQueryRequest{
		Database: &db_proto.Database{
			Name:     common.DbHealumName,
			Table:    table,
			Driver:   common.DbHealumDriver,
			Metadata: common.SearchableMetaMap,
		},
		Query: q,
	})
}

func queryMerge() string {
	query := fmt.Sprintf(
		`LET createdBy = (FOR createdBy in %v FILTER doc.data.createdBy.id == createdBy._key RETURN createdBy.data)
		LET team_members= (FOR user,edge IN OUTBOUND CONCAT("%v/",doc.id) %v
		                    FOR role IN %v FILTER role._key == edge.parameter1
		                    RETURN {user:user.data, role:role.data})
		LET products = (FOR product IN OUTBOUND CONCAT("%v/",doc.id) %v FILTER NOT_NULL(product.data)  RETURN product.data)
		LET services = (FOR service IN OUTBOUND CONCAT("%v/",doc.id) %v FILTER NOT_NULL(service.data) RETURN service.data)
		RETURN MERGE_RECURSIVE(doc,{data:{createdBy:createdBy[0],team_members:team_members, products:products, services:services}})`,
		common.DbUserTable,
		common.DbTeamTable, common.DbTeamMembershipTable, common.DbRoleTable,
		common.DbTeamTable, common.DbTeamProductEdgeTable,
		common.DbTeamTable, common.DbTeamServiceEdgeTable,
	)
	return query
}

func queryMergeTeamMember(employeeProfileTable, userTable, organisationTable, teamTable, employeeModuleTable, organisationModuleTable string) string {
	query := fmt.Sprintf(`
		LET profile = (FOR profile in %v FILTER doc.data.profile.id == profile._key RETURN profile.data)
		LET user= (FOR user in %v FILTER doc._from == user._id RETURN user.data)
		LET organisation = (FOR organisation in %v FILTER doc._to == organisation._id RETURN organisation.data)
		LET teams = (FILTER NOT_NULL(doc.data.teams) FOR p IN doc.data.teams FOR team in %v FILTER p.id == team._key RETURN team.data)
		LET modules = (
		    FILTER NOT_NULL(doc.data.modules) 
			FOR m IN OUTBOUND doc._id %v
			FOR mo, om in OUTBOUND doc._to %v
			filter om._to == m._id
			return {
					id:m.id,
					name: m.data.name,
					icon_slug: m.data.icon_slug,
					name_slug: m.data.name_slug,
					description: m.data.description,
					summary: m.data.summary,
					display_name: NOT_NULL(om.data.display_name,m.data.name),
					path: m.data.path
				})
		RETURN MERGE_RECURSIVE(doc,{data:{profile:profile[0],user:user[0], organisation:organisation[0],teams:teams, modules:modules}})`,
		employeeProfileTable, userTable, organisationTable, teamTable, employeeModuleTable, organisationModuleTable)
	return query
}

func teamToRecord(team *team_proto.Team) (string, error) {
	data, err := common.MarhalToObject(team)
	if err != nil {
		return "", err
	}

	delete(data, "team_members")
	delete(data, "products")
	delete(data, "services")

	common.FilterObject(data, "createdBy", team.CreatedBy)
	var createdById string
	if team.CreatedBy != nil {
		createdById = team.CreatedBy.Id
	}

	d := map[string]interface{}{
		"_key":       team.Id,
		"id":         team.Id,
		"created":    team.Created,
		"updated":    team.Updated,
		"name":       team.Name,
		"parameter1": team.OrgId,
		"parameter2": createdById,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToTeam(r *db_proto.Record) (*team_proto.Team, error) {
	var p team_proto.Team
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}
func employeeToRecord(from, to string, employee *team_proto.Employee) (string, error) {
	data, err := common.MarhalToObject(employee)
	if err != nil {
		return "", err
	}

	common.FilterObject(data, "role", employee.Role)
	common.FilterObject(data, "profile", employee.Profile)
	common.FilterObject(data, "user", employee.User)
	common.FilterObject(data, "organisation", employee.Organisation)

	//teams
	if len(employee.Teams) > 0 {
		var arr []interface{}
		for _, item := range employee.Teams {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["teams"] = arr
	} else {
		delete(data, "teams")
	}
	//modules
	if len(employee.Modules) > 0 {
		var arr []interface{}
		for _, item := range employee.Modules {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["modules"] = arr
	} else {
		delete(data, "modules")
	}

	d := map[string]interface{}{
		"_key":       employee.Id,
		"_from":      from,
		"_to":        to,
		"id":         employee.Id,
		"created":    employee.Created,
		"updated":    employee.Updated,
		"parameter1": employee.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToEmployee(r *db_proto.Record) (*team_proto.Employee, error) {
	var p team_proto.Employee
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func employeeProfileToRecord(profile *team_proto.EmployeeProfile) (string, error) {
	data, err := common.MarhalToObject(profile)
	if err != nil {
		return "", err
	}

	d := map[string]interface{}{
		"_key":       profile.Id,
		"id":         profile.Id,
		"created":    profile.Created,
		"updated":    profile.Updated,
		"parameter1": profile.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToEmployeeProfile(r *db_proto.Record) (*team_proto.EmployeeProfile, error) {
	var p team_proto.EmployeeProfile
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func membershipToRecord(from, to, roleId string) (string, error) {
	d := map[string]interface{}{
		"_from":      from,
		"_to":        to,
		"parameter1": roleId,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToMemebership(r *db_proto.Record) (*team_proto.TeamMembership, error) {
	var p team_proto.TeamMembership
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// AllTeams get all teams
func All(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*team_proto.Team, error) {
	var teams []*team_proto.Team
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)
	merge_query := queryMerge()

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbTeamTable, query, sort_query, limit_query, merge_query)

	resp, err := runQuery(ctx, q, common.DbTeamTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if team, err := recordToTeam(r); err == nil {
			teams = append(teams, team)
		}
	}
	return teams, nil
}

// Create creates a team
func Create(ctx context.Context, team *team_proto.Team) error {
	if team.Created == 0 {
		team.Created = time.Now().Unix()
	}
	team.Updated = time.Now().Unix()

	record, err := teamToRecord(team)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 	
		IN %v`, team.Id, record, record, common.DbTeamTable)
	if _, err := runQuery(ctx, q, common.DbTeamTable); err != nil {
		return err
	}

	_from := fmt.Sprintf(`%v/%v`, common.DbTeamTable, team.Id)
	// store team_members edge
	for _, member := range team.TeamMembers {
		if member.Role != nil && member.User != nil {
			_from := fmt.Sprintf("%v/%v", common.DbTeamTable, team.Id)
			_to := fmt.Sprintf("%v/%v", common.DbUserTable, member.User.Id)

			record, err := membershipToRecord(_from, _to, member.Role.Id)
			if err != nil {
				return err
			}
			if len(record) == 0 {
				return errors.New("server serialization")
			}
			field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)
			q := fmt.Sprintf(`
			UPSERT %v
			INSERT %v
			UPDATE %v
			INTO %v`, field, record, record, common.DbTeamMembershipTable)
			if _, err := runQuery(ctx, q, common.DbTeamMembershipTable); err != nil {
				log.Println("membership edge create error:", err)
				return err
			}
		}
	}

	// store product edge
	for _, product := range team.Products {
		field := fmt.Sprintf(`{_from:"%v",_to:"%v/%v"} `, _from, common.DbProductTable, product.Id)
		q = fmt.Sprintf(`INSERT %v INTO %v`, field, common.DbTeamProductEdgeTable)
		if _, err := runQuery(ctx, q, common.DbTeamProductEdgeTable); err != nil {
			return err
		}
	}
	// store service edge
	for _, service := range team.Services {
		field := fmt.Sprintf(`{_from:"%v",_to:"%v/%v"} `, _from, common.DbServiceTable, service.Id)
		q = fmt.Sprintf(`INSERT %v INTO %v`, field, common.DbTeamServiceEdgeTable)
		if _, err := runQuery(ctx, q, common.DbTeamServiceEdgeTable); err != nil {
			return err
		}
	}

	return nil
}

// Read reads a team by ID
func Read(ctx context.Context, id, orgId, teamId string) (*team_proto.Team, error) {
	merge_query := queryMerge()
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s`, common.DbTeamTable, query, merge_query)

	resp, err := runQuery(ctx, q, common.DbTeamTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToTeam(resp.Records[0])
	return data, err
}

// Delete deletes a team by ID
func Delete(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbTeamTable, query, common.DbTeamTable)
	_, err := runQuery(ctx, q, common.DbTeamTable)
	return err
}

// Filter get all team filters
func Filter(ctx context.Context, req *team_proto.FilterRequest) ([]*team_proto.Team, error) {
	var teams []*team_proto.Team

	query := common.QueryAuth(`FILTER`, req.OrgId, "")
	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)
	merge_query := queryMerge()

	if len(req.Product) > 0 {
		g := common.QueryStringFromArray(req.Product)
		query += fmt.Sprintf(" && doc.data.products[*].id ANY IN [%v]", g)
	}

	q := fmt.Sprintf(`
		FOR doc IN %v
		%v
		%s
		%s
		%s`, common.DbTeamTable, query, sort_query, limit_query, merge_query)

	resp, err := runQuery(ctx, q, common.DbTeamTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if team, err := recordToTeam(r); err == nil {
			teams = append(teams, team)
		}
	}
	return teams, nil
}

// Searches teams by name and/or ..., uses Elasticsearch middleware
func Search(ctx context.Context, req *team_proto.SearchRequest) ([]*team_proto.Team, error) {
	var teams []*team_proto.Team

	query := common.QueryAuth(`FILTER`, req.OrgId, "")
	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)
	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	merge_query := queryMerge()

	if len(req.TeamName) > 0 {
		query += fmt.Sprintf(` && doc.name == "%v"`, req.TeamName)
	}
	if len(req.TeamMember) > 0 {
		//FIXME: add logic here
	}

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbTeamTable, query, sort_query, limit_query, merge_query)

	resp, err := runQuery(ctx, q, common.DbTeamTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if team, err := recordToTeam(r); err == nil {
			teams = append(teams, team)
		}
	}
	return teams, nil
}

func AllTeamMember(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*team_proto.Employee, error) {
	var employees []*team_proto.Employee

	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	merge_query := queryMergeTeamMember(common.DbEmployeeProfileTable, common.DbUserTable, common.DbOrganisationTable, common.DbTeamTable, common.DbEmployeeModuleEdgeTable, common.DbOrgModuleEdgeTable)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbEmployeeTable, query, sort_query, limit_query, merge_query)

	resp, err := runQuery(ctx, q, common.DbTeamTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if employee, err := recordToEmployee(r); err == nil {
			employees = append(employees, employee)
		}
	}
	return employees, nil
}

func ReadTeamMember(ctx context.Context, userId, orgId, teamId string) (*user_proto.User, *team_proto.Employee, error) {
	// read employee
	_from := fmt.Sprintf("%v/%v", common.DbUserTable, userId)
	query := fmt.Sprintf(`FILTER doc._from == "%v"`, _from)
	query = common.QueryAuth(query, orgId, "")
	merge_query := queryMergeTeamMember(common.DbEmployeeProfileTable, common.DbUserTable, common.DbOrganisationTable, common.DbTeamTable, common.DbEmployeeModuleEdgeTable, common.DbOrgModuleEdgeTable)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s`, common.DbEmployeeTable, query, merge_query)

	resp, err := runQuery(ctx, q, common.DbEmployeeTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, nil, err
	}

	employee, err := recordToEmployee(resp.Records[0])
	return employee.User, employee, err
}

func FilterTeamMember(ctx context.Context, req *team_proto.FilterTeamMemberRequest) ([]*team_proto.Employee, error) {
	var employees []*team_proto.Employee
	query := common.QueryAuth(`FILTER`, req.OrgId, req.TeamId)
	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)
	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	merge_query := queryMergeTeamMember(common.DbEmployeeProfileTable, common.DbUserTable, common.DbOrganisationTable, common.DbTeamTable, common.DbEmployeeModuleEdgeTable, common.DbOrgModuleEdgeTable)

	if len(req.Team) > 0 {
		g := common.QueryStringFromArray(req.Team)
		query += fmt.Sprintf(" && doc.data.teams[*].id ANY IN [%v]", g)
	}
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbEmployeeTable, query, sort_query, limit_query, merge_query)

	resp, err := runQuery(ctx, q, common.DbEmployeeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if employee, err := recordToEmployee(r); err == nil {
			employees = append(employees, employee)
		}
	}
	return employees, nil
}

func CreateEmployeeEdge(ctx context.Context, employee *team_proto.Employee, userId, orgId string) error {
	if len(employee.Id) == 0 {
		employee.Id = uuid.NewUUID().String()
	}
	if employee.Created == 0 {
		employee.Created = time.Now().Unix()
	}
	employee.Updated = time.Now().Unix()

	_from := fmt.Sprintf(`%v/%v`, common.DbUserTable, userId)
	_to := fmt.Sprintf(`%v/%v`, common.DbOrganisationTable, orgId)
	record, err := employeeToRecord(_from, _to, employee)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}
	field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)
	q := fmt.Sprintf(`
		UPSERT %v
		INSERT %v
		UPDATE %v
		IN %v`, field, record, record, common.DbEmployeeTable)

	if _, err := runQuery(ctx, q, common.DbEmployeeTable); err != nil {
		log.Println("employee create error", err)
		return err
	}

	return nil
}

func CreateEmployeeProfile(ctx context.Context, profile *team_proto.EmployeeProfile) error {
	// Create Employee profile entity
	if len(profile.Id) == 0 {
		profile.Id = uuid.NewUUID().String()
	}
	if profile.Created == 0 {
		profile.Created = time.Now().Unix()
	}
	profile.Updated = time.Now().Unix()

	record, err := employeeProfileToRecord(profile)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}
	q := fmt.Sprintf(`
			UPSERT { _key: "%v" } 
			INSERT %v 
			UPDATE %v 
			IN %v`, profile.Id, record, record, common.DbEmployeeProfileTable)
	if _, err := runQuery(ctx, q, common.DbUserTable); err != nil {
		log.Println("profile create error", err)
		return err
	}

	return nil
}

func CreateTeamMembership(ctx context.Context, employee *team_proto.Employee, user *user_proto.User) error {
	for _, team := range employee.Teams {
		_from := fmt.Sprintf("%v/%v", common.DbTeamTable, team.Id)
		_to := fmt.Sprintf("%v/%v", common.DbUserTable, user.Id)

		record, err := membershipToRecord(_from, _to, employee.Role.Id)
		if err != nil {
			return err
		}
		if len(record) == 0 {
			return err
		}
		field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)
		q := fmt.Sprintf(`
			UPSERT %v
			INSERT %v
			UPDATE %v
			INTO %v`, field, record, record, common.DbTeamMembershipTable)
		if _, err := runQuery(ctx, q, common.DbTeamMembershipTable); err != nil {
			log.Println("membership edge create error:", err)
			return err
		}
	}
	return nil
}

func CreateEmpoyeeModuleAccess(ctx context.Context, userId, orgId string, modules []*static_proto.Module) error {
	// delete module edge from employee module edge
	q_del := fmt.Sprintf(`
		FOR o,e IN OUTBOUND "%v/%v" %v
		FILTER e.parameter1 == "%v"
		FOR doc IN %v
		FILTER doc._from == e._id
		REMOVE doc IN %v`,
		common.DbUserTable, userId, common.DbEmployeeTable,
		orgId,
		common.DbEmployeeModuleEdgeTable,
		common.DbEmployeeModuleEdgeTable)
	if _, err := runQuery(ctx, q_del, common.DbOrgModuleEdgeTable); err != nil {
		return err
	}

	// create modules edge table
	for _, module := range modules {
		_to := fmt.Sprintf("%v/%v", common.DbModuleTable, module.Id)
		field := fmt.Sprintf(`{_from:e._id,_to:"%v"} `, _to)

		q := fmt.Sprintf(`
			FOR o,e IN OUTBOUND "%v/%v" %v
			FILTER e.parameter1 == "%v"
			UPSERT %v
			INSERT %v
			UPDATE %v
			INTO %v`,
			common.DbUserTable, userId, common.DbEmployeeTable,
			orgId,
			field,
			field,
			field,
			common.DbEmployeeModuleEdgeTable)
		if _, err := runQuery(ctx, q, common.DbEmployeeModuleEdgeTable); err != nil {
			//log error here
			return err
		}
	}

	//update empoyee with modules
	var arr []interface{}
	for _, item := range modules {
		arr = append(arr, map[string]string{
			"id": item.Id,
		})
	}

	body, err := json.Marshal(arr)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`
		FOR o,e IN OUTBOUND "%v/%v" %v
		FILTER e.parameter1 == "%v"
		UPDATE e WITH {data: {modules:%v}}
		IN %v`,
		common.DbUserTable, userId, common.DbEmployeeTable,
		orgId,
		string(body), common.DbEmployeeTable)

	if _, err := runQuery(ctx, q, common.DbEmployeeTable); err != nil {
		return err
	}
	return nil
}

func GetAccessibleModulesByEmployee(ctx context.Context, userId, orgId string) ([]*static_proto.Module, error) {
	var modules []*static_proto.Module
	q := fmt.Sprintf(`
		FOR o,e IN OUTBOUND "%v/%v" %v
		FILTER e.parameter1 == "%v"
		FOR m IN OUTBOUND e._id %v
		FOR mo, om in OUTBOUND o._id %v
		FILTER om._to == m._id
		return {data:{
				id:m.id,
				name: m.data.name,
				icon_slug: m.data.icon_slug,
				name_slug: m.data.name_slug,
				description: m.data.description,
				summary: m.data.summary,
				display_name: NOT_NULL(om.data.display_name,m.data.name),
				path: m.data.path
			}}`,
		common.DbUserTable, userId, common.DbEmployeeTable,
		orgId,
		common.DbEmployeeModuleEdgeTable,
		common.DbOrgModuleEdgeTable)

	resp, err := runQuery(ctx, q, common.DbEmployeeModuleEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if module, err := common.RecordToModule(r); err == nil {
			modules = append(modules, module)
		}
	}
	return modules, nil
}

func DeleteEmployee(ctx context.Context, employeeId, orgId string) error {
	queries := []string{}

	owner_query := fmt.Sprintf(`LET owner_id = (
		FOR org IN OUTBOUND "user/%v" %v
		RETURN org.data.owner.id
		)[0]`, employeeId, common.DbOrganisationTable)
	queries = append(queries, owner_query)

	goal_query := fmt.Sprintf(`LET goal = (
		FOR doc IN %v
		FILTER doc.data.createdBy.id == "%v"
		UPDATE doc._key WITH {data:{createdBy:{id:owner_id}}} IN %v
		RETURN OLD._key)`, common.DbGoalTable, employeeId, common.DbGoalTable)
	queries = append(queries, goal_query)

	challenge_query := fmt.Sprintf(`LET challenge = (
		FOR doc IN %v
		FILTER doc.data.createdBy.id == "%v"
		UPDATE doc._key WITH {data:{createdBy:{id:owner_id}}} IN %v
		RETURN OLD._key)`, common.DbChallengeTable, employeeId, common.DbChallengeTable)
	queries = append(queries, challenge_query)

	habit_query := fmt.Sprintf(`LET habit = (
		FOR doc IN %v
		FILTER doc.data.createdBy.id == "%v"
		UPDATE doc._key WITH {data:{createdBy:{id:owner_id}}} IN %v
		RETURN OLD._key)`, common.DbHabitTable, employeeId, common.DbHabitTable)
	queries = append(queries, habit_query)

	content_query := fmt.Sprintf(`LET content = (
		FOR doc IN %v
		FILTER doc.data.createdBy.id == "%v"
		UPDATE doc._key WITH {data:{createdBy:{id:owner_id}}} IN %v
		RETURN OLD._key)`, common.DbContentTable, employeeId, common.DbContentTable)
	queries = append(queries, content_query)

	share_goal_user := fmt.Sprintf(`LET share_goal_user = (
		FOR doc IN %v
		FILTER doc.data.shared_by.id == "%v"
		UPDATE doc._key WITH {data:{shared_by:{id:owner_id}}} IN %v
		RETURN OLD._key)`, common.DbShareGoalUserEdgeTable, employeeId, common.DbShareGoalUserEdgeTable)
	queries = append(queries, share_goal_user)

	share_challenge_user := fmt.Sprintf(`LET share_challenge_user = (
		FOR doc IN %v
		FILTER doc.data.shared_by.id == "%v"
		UPDATE doc._key WITH {data:{shared_by:{id:owner_id}}} IN %v
		RETURN OLD._key)`, common.DbShareChallengeUserEdgeTable, employeeId, common.DbShareChallengeUserEdgeTable)
	queries = append(queries, share_challenge_user)

	share_habit_user := fmt.Sprintf(`LET share_habit_user = (
		FOR doc IN %v
		FILTER doc.data.shared_by.id == "%v"
		UPDATE doc._key WITH {data:{shared_by:{id:owner_id}}} IN %v
		RETURN OLD._key)`, common.DbShareHabitUserEdgeTable, employeeId, common.DbShareHabitUserEdgeTable)
	queries = append(queries, share_habit_user)

	share_content_user := fmt.Sprintf(`LET share_content_user = (
		FOR doc IN %v
		FILTER doc.data.shared_by.id == "%v"
		UPDATE doc._key WITH {data:{shared_by:{id:owner_id}}} IN %v
		RETURN OLD._key)`, common.DbShareContentUserEdgeTable, employeeId, common.DbShareContentUserEdgeTable)
	queries = append(queries, share_content_user)

	share_survey_user := fmt.Sprintf(`LET survey = (FOR doc IN %v
		FILTER doc.data.shared_by.id == "%v"
		UPDATE doc._key WITH {data:{creator:{id:owner_id}}} IN %v
		RETURN OLD._key)`, common.DbShareSurveyUserEdgeTable, employeeId, common.DbShareSurveyUserEdgeTable)
	queries = append(queries, share_survey_user)

	team_membership := fmt.Sprintf(`LET team_membership = (
		FOR doc IN %v
		FILTER doc._to == "user/%v"
		REMOVE doc IN %v
		RETURN OLD._key)`, common.DbTeamMembershipTable, employeeId, common.DbTeamMembershipTable)
	queries = append(queries, team_membership)

	employee_profile := fmt.Sprintf(`LET employee_profile = (
		FOR e IN %v
		FILTER e._key == "%v"
		FOR ep IN e.employee_profile
		FOR p IN %v
			FILTER ep.id == p.id
			REMOVE p IN %v
	)`, common.DbEmployeeTable, employeeId, common.DbEmployeeProfileTable, common.DbEmployeeProfileTable)
	queries = append(queries, employee_profile)

	user_query := fmt.Sprintf(`LET user = (
		FOR doc IN %v
		FILTER doc._key == "%v"
		REMOVE doc IN %v
		)`, common.DbUserTable, employeeId, common.DbUserTable)
	queries = append(queries, user_query)

	account_edge_query := fmt.Sprintf(`LET account_id = (
		FOR doc IN %v
		FILTER doc._from == "user/%v"
		REMOVE doc IN %v
		RETURN OLD._to)`, common.DbUserAccountEdgeTable, employeeId, common.DbUserAccountEdgeTable)
	queries = append(queries, account_edge_query)

	account_query := fmt.Sprintf(`LET account = (
		FOR doc IN %v
		FILTER doc._id == account_id
		REMOVE doc IN account
		RETURN OLD._key)`, common.DbAccountTable)
	queries = append(queries, account_query)

	employee_query := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._from == "user/%v"
		REMOVE doc IN employee`, employeeId, common.DbEmployeeTable)
	queries = append(queries, employee_query)

	var q string
	for _, query := range queries {
		q = fmt.Sprintf(`%s
			%s`, q, query)
	}

	_, err := runQuery(ctx, q, common.DbUserTable)
	return err
}
