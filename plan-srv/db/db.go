package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"server/common"
	db_proto "server/db-srv/proto/db"
	plan_proto "server/plan-srv/proto/plan"
	common_proto "server/static-srv/proto/common"
	static_proto "server/static-srv/proto/static"
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

func planToRecord(plan *plan_proto.Plan) (string, error) {
	data, err := common.MarhalToObject(plan)
	if err != nil {
		return "", err
	}
	// users
	if len(plan.Users) > 0 {
		var arr []interface{}
		for _, item := range plan.Users {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["users"] = arr
	} else {
		delete(data, "users")
	}
	// goals
	if len(plan.Goals) > 0 {
		var arr []interface{}
		for _, item := range plan.Goals {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["goals"] = arr
	} else {
		delete(data, "goals")
	}
	// creator
	var creatorId string
	if plan.Creator != nil {
		creatorId = plan.Creator.Id
	}
	common.FilterObject(data, "creator", plan.Creator)
	// collaborators
	if len(plan.Collaborators) > 0 {
		var arr []interface{}
		for _, item := range plan.Collaborators {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["collaborators"] = arr
	} else {
		delete(data, "collaborators")
	}
	// shares
	if len(plan.Shares) > 0 {
		var arr []interface{}
		for _, item := range plan.Shares {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["shares"] = arr
	} else {
		delete(data, "shares")
	}

	d := map[string]interface{}{
		"_key":       plan.Id,
		"id":         plan.Id,
		"created":    plan.Created,
		"updated":    plan.Updated,
		"name":       plan.Name,
		"parameter1": plan.OrgId,
		"parameter2": creatorId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToPlan(r *db_proto.Record) (*plan_proto.Plan, error) {
	var p plan_proto.Plan
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func filterToRecord(filter *plan_proto.PlanFilter) (string, error) {
	d := map[string]interface{}{
		"name":       filter.DisplayName,
		"parameter1": filter.FilterSlug,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToFilter(r *db_proto.Record) (*plan_proto.PlanFilter, error) {
	var filter plan_proto.PlanFilter
	// r.Parameter3 = strings.Replace(r.Parameter3, `'`, `"`, -1)
	// decoder := json.NewDecoder(bytes.NewReader([]byte(r.Parameter3)))
	// err := decoder.Decode(&plan)
	// if err != nil {
	// return nil, err
	// }
	filter.DisplayName = r.Name
	filter.FilterSlug = r.Parameter1
	return &filter, nil
}

func sharedToRecord(from, to, orgId string, shared *plan_proto.SharePlanUser) (string, error) {
	data, err := common.MarhalToObject(shared)
	if err != nil {
		return "", err
	}
	common.FilterObject(data, "plan", shared.Plan)
	common.FilterObject(data, "user", shared.User)
	common.FilterObject(data, "shared_by", shared.SharedBy)
	var sharedById string
	if shared.SharedBy != nil {
		sharedById = shared.SharedBy.Id
	}
	d := map[string]interface{}{
		"_from":      from,
		"_to":        to,
		"created":    shared.Created,
		"updated":    shared.Updated,
		"parameter1": orgId,
		"parameter2": sharedById,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToShared(r *db_proto.Record) (*plan_proto.SharePlanUser, error) {
	var p plan_proto.SharePlanUser
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// All get all plans
func All(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*plan_proto.Plan, error) {
	var plans []*plan_proto.Plan
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s	
		LET users = (FILTER NOT_NULL(doc.data.users) FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET goals = (FILTER NOT_NULL(doc.data.goals) FOR g IN doc.data.goals FOR p IN %v FILTER g.id == p._key RETURN p.data)
		LET creator = (FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET collaborators = (FILTER NOT_NULL(doc.data.collaborators) FOR u IN doc.data.collaborators FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET shares = (FILTER NOT_NULL(doc.data.shares) FOR u IN doc.data.shares FOR p IN %v FILTER u.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			users:users,
			goals:goals,
			creator:creator[0],
			collaborators:collaborators,
			shares:shares
		}})`, common.DbPlanTable, query, sort_query, limit_query,
		common.DbUserTable, common.DbGoalTable, common.DbUserTable, common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbPlanTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if plan, err := recordToPlan(r); err == nil {
			plans = append(plans, plan)
		}
	}
	return plans, nil
}

// Creates a plan
func Create(ctx context.Context, plan *plan_proto.Plan) error {
	if len(plan.Id) == 0 {
		plan.Id = uuid.NewUUID().String()
	}
	if plan.Created == 0 {
		plan.Created = time.Now().Unix()
	}
	plan.Updated = time.Now().Unix()
	// calc items_count
	items_count := 0
	if plan.Days != nil {
		for _, v := range plan.Days {
			items_count += len(v.Items)
		}
	}
	plan.ItemsCount = int64(items_count)

	record, err := planToRecord(plan)
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
		INTO %v`, plan.Id, record, record, common.DbPlanTable)
	_, err = runQuery(ctx, q, common.DbPlanTable)
	return err
}

// Reads a plan by ID
func Read(ctx context.Context, id, orgId, teamId string) (*plan_proto.Plan, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v 
		%s 
		LET users = (FILTER NOT_NULL(doc.data.users) FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET goals = (FILTER NOT_NULL(doc.data.goals) FOR g IN doc.data.goals FOR p IN %v FILTER g.id == p._key RETURN p.data)
		LET creator = (FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET collaborators = (FILTER NOT_NULL(doc.data.collaborators) FOR u IN doc.data.collaborators FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET shares = (FILTER NOT_NULL(doc.data.shares) FOR u IN doc.data.shares FOR p IN %v FILTER u.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			users:users,
			goals:goals,
			creator:creator[0],
			collaborators:collaborators,
			shares:shares
		}})`, common.DbPlanTable, query,
		common.DbUserTable, common.DbGoalTable, common.DbUserTable, common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbPlanTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToPlan(resp.Records[0])
	return data, err
}

// Deletes a plan by ID
func Delete(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbPlanTable, query, common.DbPlanTable)
	_, err := runQuery(ctx, q, common.DbPlanTable)
	return err
}

// Searches plans by name and/or ..., uses Elasticsearch middleware
func Search(ctx context.Context, name, orgid, teamId string, offset, limit, from, to int64, sortParameter, sortDirection string) ([]*plan_proto.Plan, error) {
	var plans []*plan_proto.Plan

	query := fmt.Sprintf(`FILTER doc.name == "%v"`, name)
	query = common.QueryAuth(query, orgid, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET users = (FILTER NOT_NULL(doc.data.users) FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET goals = (FILTER NOT_NULL(doc.data.goals) FOR g IN doc.data.goals FOR p IN %v FILTER g.id == p._key RETURN p.data)
		LET creator = (FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET collaborators = (FILTER NOT_NULL(doc.data.collaborators) FOR u IN doc.data.collaborators FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET shares = (FILTER NOT_NULL(doc.data.shares) FOR u IN doc.data.shares FOR p IN %v FILTER u.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			users:users,
			goals:goals,
			creator:creator[0],
			collaborators:collaborators,
			shares:shares
		}})`, common.DbPlanTable, query, sort_query, limit_query,
		common.DbUserTable, common.DbGoalTable, common.DbUserTable, common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbPlanTable)
	if err != nil {
		return nil, err
	}

	// parsing
	for _, r := range resp.Records {
		if plan, err := recordToPlan(r); err == nil {
			plans = append(plans, plan)
		}
	}
	return plans, nil
}

// Templates get all templates - Send all plans where isTemplate = true
func Templates(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*plan_proto.Plan, error) {
	query := `FILTER doc.data.isTemplate == true`
	query = common.QueryAuth(query, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	var plans []*plan_proto.Plan
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET users = (FILTER NOT_NULL(doc.data.users) FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET goals = (FILTER NOT_NULL(doc.data.goals) FOR g IN doc.data.goals FOR p IN %v FILTER g.id == p._key RETURN p.data)
		LET creator = (FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET collaborators = (FILTER NOT_NULL(doc.data.collaborators) FOR u IN doc.data.collaborators FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET shares = (FILTER NOT_NULL(doc.data.shares) FOR u IN doc.data.shares FOR p IN %v FILTER u.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			users:users,
			goals:goals,
			creator:creator[0],
			collaborators:collaborators,
			shares:shares
		}})`, common.DbPlanTable, query, sort_query, limit_query,
		common.DbUserTable, common.DbGoalTable, common.DbUserTable, common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbPlanTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if plan, err := recordToPlan(r); err == nil {
			plans = append(plans, plan)
		}
	}
	return plans, nil
}

// Drafts get all draft plans - Get all plans where the status is draft
func Drafts(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*plan_proto.Plan, error) {
	query := fmt.Sprintf(`FILTER doc.data.status == "%v"`, plan_proto.StatusEnum_DRAFT)
	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	var plans []*plan_proto.Plan
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET users = (FILTER NOT_NULL(doc.data.users) FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET goals = (FILTER NOT_NULL(doc.data.goals) FOR g IN doc.data.goals FOR p IN %v FILTER g.id == p._key RETURN p.data)
		LET creator = (FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET collaborators = (FILTER NOT_NULL(doc.data.collaborators) FOR u IN doc.data.collaborators FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET shares = (FILTER NOT_NULL(doc.data.shares) FOR u IN doc.data.shares FOR p IN %v FILTER u.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			users:users,
			goals:goals,
			creator:creator[0],
			collaborators:collaborators,
			shares:shares
		}})`, common.DbPlanTable, query, sort_query, limit_query,
		common.DbUserTable, common.DbGoalTable, common.DbUserTable, common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbPlanTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if plan, err := recordToPlan(r); err == nil {
			plans = append(plans, plan)
		}
	}
	return plans, nil
}

// ByCreator get all plans created by a particular team member - Get all plans where createdBy = {userid}
func ByCreator(ctx context.Context, id, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*plan_proto.Plan, error) {
	query := fmt.Sprintf(`FILTER doc.data.creator.id == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	var plans []*plan_proto.Plan
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET users = (FILTER NOT_NULL(doc.data.users) FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET goals = (FILTER NOT_NULL(doc.data.goals) FOR g IN doc.data.goals FOR p IN %v FILTER g.id == p._key RETURN p.data)
		LET creator = (FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET collaborators = (FILTER NOT_NULL(doc.data.collaborators) FOR u IN doc.data.collaborators FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET shares = (FILTER NOT_NULL(doc.data.shares) FOR u IN doc.data.shares FOR p IN %v FILTER u.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			users:users,
			goals:goals,
			creator:creator[0],
			collaborators:collaborators,
			shares:shares
		}})`, common.DbPlanTable, query, sort_query, limit_query,
		common.DbUserTable, common.DbGoalTable, common.DbUserTable, common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbPlanTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if plan, err := recordToPlan(r); err == nil {
			plans = append(plans, plan)
		}
	}
	return plans, nil
}

// Filters get all plan filters
func Filters(ctx context.Context, offset, limit int64, sortParameter, sortDirection string) ([]*plan_proto.PlanFilter, error) {
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	var filters []*plan_proto.PlanFilter
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbPlanFilterTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbPlanFilterTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if filter, err := recordToFilter(r); err == nil {
			filters = append(filters, filter)
		}
	}
	return filters, nil
}

// TimeFilters get all plans by time period - Return all plans that were created between a specific time period
func TimeFilters(ctx context.Context, start, end int64, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*plan_proto.Plan, error) {
	query := fmt.Sprintf(`FILTER %d < doc.created AND doc.created < %d`, start, end)
	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	var plans []*plan_proto.Plan
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET users = (FILTER NOT_NULL(doc.data.users) FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET goals = (FILTER NOT_NULL(doc.data.goals) FOR g IN doc.data.goals FOR p IN %v FILTER g.id == p._key RETURN p.data)
		LET creator = (FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET collaborators = (FILTER NOT_NULL(doc.data.collaborators) FOR u IN doc.data.collaborators FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET shares = (FILTER NOT_NULL(doc.data.shares) FOR u IN doc.data.shares FOR p IN %v FILTER u.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			users:users,
			goals:goals,
			creator:creator[0],
			collaborators:collaborators,
			shares:shares
		}})`, common.DbPlanTable, query, sort_query, limit_query,
		common.DbUserTable, common.DbGoalTable, common.DbUserTable, common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbPlanTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if plan, err := recordToPlan(r); err == nil {
			plans = append(plans, plan)
		}
	}
	return plans, nil
}

// GoalFilters get all plans by goal category - Return all plans filtered by one or more goal category
func GoalFilters(ctx context.Context, goals string, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*plan_proto.Plan, error) {
	var plans []*plan_proto.Plan
	arr := strings.Split(goals, ",")
	g := common.QueryStringFromArray(arr)
	query := fmt.Sprintf(`FILTER doc.data.goals[*].id ANY IN [%v]`, g)

	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc  IN %v
		%s
		%s
		%s
		LET users = (FILTER NOT_NULL(doc.data.users) FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET goals = (FILTER NOT_NULL(doc.data.goals) FOR g IN doc.data.goals FOR p IN %v FILTER g.id == p._key RETURN p.data)
		LET creator = (FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET collaborators = (FILTER NOT_NULL(doc.data.collaborators) FOR u IN doc.data.collaborators FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET shares = (FILTER NOT_NULL(doc.data.shares) FOR u IN doc.data.shares FOR p IN %v FILTER u.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			users:users,
			goals:goals,
			creator:creator[0],
			collaborators:collaborators,
			shares:shares
		}})`, common.DbPlanTable, query, sort_query, limit_query,
		common.DbUserTable, common.DbGoalTable, common.DbUserTable, common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbPlanTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if plan, err := recordToPlan(r); err == nil {
			plans = append(plans, plan)
		}
	}
	return plans, nil
}

// CreatePlanFilter create plan filter for test
func CreatePlanFilter(ctx context.Context, filter *plan_proto.PlanFilter) error {
	record, err := filterToRecord(filter)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`INSERT %v INTO %v`, record, common.DbPlanFilterTable)
	_, err = runQuery(ctx, q, common.DbPlanFilterTable)
	return err
}

func SharePlan(ctx context.Context, plans []*plan_proto.Plan, users []*user_proto.User, sharedBy *user_proto.User, orgId string) ([]string, error) {
	userids := []string{}
	for _, plan := range plans {
		for _, user := range users {
			shared := &plan_proto.SharePlanUser{
				Id:       uuid.NewUUID().String(),
				Status:   static_proto.ShareStatus_SHARED,
				Updated:  time.Now().Unix(),
				Created:  time.Now().Unix(),
				Plan:     plan,
				User:     user,
				SharedBy: sharedBy,
			}

			_from := fmt.Sprintf(`%v/%v`, common.DbPlanTable, plan.Id)
			_to := fmt.Sprintf(`%v/%v`, common.DbUserTable, user.Id)
			record, err := sharedToRecord(_from, _to, orgId, shared)
			if err != nil {
				return nil, err
			}
			if len(record) == 0 {
				return nil, errors.New("server serialization")
			}

			field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)
			q := fmt.Sprintf(`
				UPSERT %v
				INSERT %v
				UPDATE %v
				INTO %v
				RETURN {data:{user_id: OLD ? "" : NEW.data.user.id}}`, field, record, record, common.DbSharePlanUserEdgeTable)

			resp, err := runQuery(ctx, q, common.DbSharePlanUserEdgeTable)
			if err != nil {
				return nil, err
			}

			// parsing to check whether this was an update (returns nothing) or insert (returns inserted user_id)
			b, err := common.RecordToInsertedUserId(resp.Records[0])
			if err != nil {
				return nil, err
			}
			if len(b) > 0 {
				userids = append(userids, b)
			}

			// save pending
			any, err := common.FilteredAnyFromObject(common.PLAN_TYPE, plan.Id)
			if err != nil {
				return nil, err
			}
			// save pending
			pending := &common_proto.Pending{
				Id:         uuid.NewUUID().String(),
				Created:    shared.Created,
				Updated:    shared.Updated,
				SharedBy:   sharedBy,
				SharedWith: user,
				Item:       any,
				OrgId:      orgId,
			}

			q1, err1 := common.SavePending(pending, plan.Id)
			if err1 != nil {
				return nil, err1
			}

			if _, err := runQuery(ctx, q1, common.DbPendingTable); err != nil {
				return nil, err
			}
		}
	}
	return userids, nil
}

func AutocompleteSearch(ctx context.Context, title string, sortParameter, sortDirection string) ([]*static_proto.AutocompleteResponse, error) {
	response := []*static_proto.AutocompleteResponse{}
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER LIKE(doc.name, "%v",true)
		%s
		LET users = (FILTER NOT_NULL(doc.data.users) FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET goals = (FILTER NOT_NULL(doc.data.goals) FOR g IN doc.data.goals FOR p IN %v FILTER g.id == p._key RETURN p.data)
		LET creator = (FOR u IN doc.data.users FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET collaborators = (FILTER NOT_NULL(doc.data.collaborators) FOR u IN doc.data.collaborators FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET shares = (FILTER NOT_NULL(doc.data.shares) FOR u IN doc.data.shares FOR p IN %v FILTER u.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			users:users,
			goals:goals,
			creator:creator[0],
			collaborators:collaborators,
			shares:shares
		}})`, common.DbPlanTable, `%`+title+`%`, sort_query,
		common.DbUserTable, common.DbGoalTable, common.DbUserTable, common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbPlanTable)

	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		response = append(response, &static_proto.AutocompleteResponse{Id: r.Id, Title: r.Name, OrgId: r.Parameter1})
	}
	return response, nil
}

func AutocompleteTags(ctx context.Context, orgId, name string) ([]string, error) {
	var tags []string

	var query string
	if len(orgId) > 0 {
		query = fmt.Sprintf(`FILTER doc.parameter1 == "%v"`, orgId)
	}

	q := fmt.Sprintf(`
		LET tags = (FOR doc IN %v
		%v
		RETURN doc.data.tags)[**]
		FOR t IN tags
		FILTER LIKE(t,"%v",true)
		LET ret = {parameter1:t}
		RETURN DISTINCT ret
		`, common.DbPlanTable, query, `%`+name+`%`)
	resp, err := runQuery(ctx, q, common.DbPlanTable)

	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		tags = append(tags, r.Parameter1)
	}
	return tags, nil
}
