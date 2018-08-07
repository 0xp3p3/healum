package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	db_proto "server/db-srv/proto/db"
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

func goalToRecord(goal *behaviour_proto.Goal) (string, error) {
	data, err := common.MarhalToObject(goal)
	if err != nil {
		return "", err
	}

	// category
	common.FilterObject(data, "category", goal.Category)
	// target
	if goal.Target != nil {
		item := goal.Target
		obj := map[string]interface{}{
			"targetValue": item.TargetValue,
			"unit":        item.Unit,
			"recurrence":  item.Recurrence,
		}
		if item.Aim != nil {
			obj["aim"] = map[string]string{"id": item.Aim.Id}
		}
		if item.Marker != nil {
			obj["marker"] = map[string]string{"id": item.Marker.Id}
		}
		data["target"] = obj
	}
	// createdby
	common.FilterObject(data, "createdBy", goal.CreatedBy)
	// trackers
	if len(goal.Trackers) > 0 {
		var arr []interface{}
		for _, item := range goal.Trackers {
			obj := map[string]interface{}{
				"frequency": item.Frequency,
				"until":     item.Until,
			}
			if item.Marker != nil {
				obj["marker"] = map[string]string{"id": item.Marker.Id}
			}
			if item.Method != nil {
				obj["method"] = map[string]string{"id": item.Method.Id}
			}
			arr = append(arr, obj)
		}
		data["trackers"] = arr
	} else {
		delete(data, "trackers")
	}
	// challenges
	if len(goal.Challenges) > 0 {
		var arr []interface{}
		for _, item := range goal.Challenges {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["challenges"] = arr
	} else {
		delete(data, "challenges")
	}
	// habits
	if len(goal.Habits) > 0 {
		var arr []interface{}
		for _, item := range goal.Habits {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["habits"] = arr
	} else {
		delete(data, "habits")
	}
	// triggers
	if len(goal.Triggers) > 0 {
		var arr []interface{}
		for _, item := range goal.Triggers {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["triggers"] = arr
	} else {
		delete(data, "triggers")
	}
	// setbacks
	if len(goal.Setbacks) > 0 {
		var arr []interface{}
		for _, item := range goal.Setbacks {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["setbacks"] = arr
	} else {
		delete(data, "setbacks")
	}
	// target users
	if len(goal.Users) > 0 {
		var arr []interface{}
		for _, item := range goal.Users {
			obj := map[string]interface{}{
				"currentValue":     item.CurrentValue,
				"expectedProgress": item.ExpectedProgress,
			}
			if item.User != nil {
				obj["user"] = map[string]string{"id": item.User.Id}
			}
			arr = append(arr, obj)
		}
		data["users"] = arr
	} else {
		delete(data, "users")
	}
	// todos
	common.FilterObject(data, "todos", goal.Todos)

	var createdById string
	if goal.CreatedBy != nil {
		createdById = goal.CreatedBy.Id
	}
	d := map[string]interface{}{
		"_key":       goal.Id,
		"id":         goal.Id,
		"created":    goal.Created,
		"updated":    goal.Updated,
		"name":       goal.Title,
		"parameter1": goal.OrgId,
		"parameter2": createdById,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToGoal(r *db_proto.Record) (*behaviour_proto.Goal, error) {
	var p behaviour_proto.Goal
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func challengeToRecord(challenge *behaviour_proto.Challenge) (string, error) {
	data, err := common.MarhalToObject(challenge)
	if err != nil {
		return "", err
	}

	// category
	common.FilterObject(data, "category", challenge.Category)
	// target
	if challenge.Target != nil {
		item := challenge.Target
		obj := map[string]interface{}{
			"targetValue": item.TargetValue,
			"unit":        item.Unit,
			"recurrence":  item.Recurrence,
		}
		if item.Aim != nil {
			obj["aim"] = map[string]string{"id": item.Aim.Id}
		}
		if item.Marker != nil {
			obj["marker"] = map[string]string{"id": item.Marker.Id}
		}
		data["target"] = obj
	}
	// createdby
	common.FilterObject(data, "createdBy", challenge.CreatedBy)
	// trackers
	if len(challenge.Trackers) > 0 {
		var arr []interface{}
		for _, item := range challenge.Trackers {
			obj := map[string]interface{}{
				"frequency": item.Frequency,
				"until":     item.Until,
			}
			if item.Marker != nil {
				obj["marker"] = map[string]string{"id": item.Marker.Id}
			}
			if item.Method != nil {
				obj["method"] = map[string]string{"id": item.Method.Id}
			}
			arr = append(arr, obj)
		}
		data["trackers"] = arr
	} else {
		delete(data, "trackers")
	}
	// habits
	if len(challenge.Habits) > 0 {
		var arr []interface{}
		for _, item := range challenge.Habits {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["habits"] = arr
	} else {
		delete(data, "habits")
	}
	// triggers
	if len(challenge.Triggers) > 0 {
		var arr []interface{}
		for _, item := range challenge.Triggers {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["triggers"] = arr
	} else {
		delete(data, "triggers")
	}
	// setbacks
	if len(challenge.Setbacks) > 0 {
		var arr []interface{}
		for _, item := range challenge.Setbacks {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["setbacks"] = arr
	} else {
		delete(data, "setbacks")
	}
	// target users
	if len(challenge.Users) > 0 {
		var arr []interface{}
		for _, item := range challenge.Users {
			obj := map[string]interface{}{
				"currentValue":     item.CurrentValue,
				"expectedProgress": item.ExpectedProgress,
			}
			if item.User != nil {
				obj["user"] = map[string]string{"id": item.User.Id}
			}
			arr = append(arr, obj)
		}
		data["users"] = arr
	} else {
		delete(data, "users")
	}
	// todos
	common.FilterObject(data, "todos", challenge.Todos)

	var createdById string
	if challenge.CreatedBy != nil {
		createdById = challenge.CreatedBy.Id
	}
	d := map[string]interface{}{
		"_key":       challenge.Id,
		"id":         challenge.Id,
		"created":    challenge.Created,
		"updated":    challenge.Updated,
		"name":       challenge.Title,
		"parameter1": challenge.OrgId,
		"parameter2": createdById,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToChallenge(r *db_proto.Record) (*behaviour_proto.Challenge, error) {
	var p behaviour_proto.Challenge
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func habitToRecord(habit *behaviour_proto.Habit) (string, error) {
	data, err := common.MarhalToObject(habit)
	if err != nil {
		return "", err
	}

	// category
	common.FilterObject(data, "category", habit.Category)
	// target
	if habit.Target != nil {
		item := habit.Target
		obj := map[string]interface{}{
			"targetValue": item.TargetValue,
			"unit":        item.Unit,
			"recurrence":  item.Recurrence,
		}
		if item.Aim != nil {
			obj["aim"] = map[string]string{"id": item.Aim.Id}
		}
		if item.Marker != nil {
			obj["marker"] = map[string]string{"id": item.Marker.Id}
		}
		data["target"] = obj
	}
	// createdby
	common.FilterObject(data, "createdBy", habit.CreatedBy)
	// trackers
	if len(habit.Trackers) > 0 {
		var arr []interface{}
		for _, item := range habit.Trackers {
			obj := map[string]interface{}{
				"frequency": item.Frequency,
				"until":     item.Until,
			}
			if item.Marker != nil {
				obj["marker"] = map[string]string{"id": item.Marker.Id}
			}
			if item.Method != nil {
				obj["method"] = map[string]string{"id": item.Method.Id}
			}
			arr = append(arr, obj)
		}
		data["trackers"] = arr
	} else {
		delete(data, "trackers")
	}

	// triggers
	if len(habit.Triggers) > 0 {
		var arr []interface{}
		for _, item := range habit.Triggers {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["triggers"] = arr
	} else {
		delete(data, "triggers")
	}
	// setbacks
	if len(habit.Setbacks) > 0 {
		var arr []interface{}
		for _, item := range habit.Setbacks {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["setbacks"] = arr
	} else {
		delete(data, "setbacks")
	}
	// target users
	if len(habit.Users) > 0 {
		var arr []interface{}
		for _, item := range habit.Users {
			obj := map[string]interface{}{
				"currentValue":     item.CurrentValue,
				"expectedProgress": item.ExpectedProgress,
			}
			if item.User != nil {
				obj["user"] = map[string]string{"id": item.User.Id}
			}
			arr = append(arr, obj)
		}
		data["users"] = arr
	} else {
		delete(data, "users")
	}
	// todos
	common.FilterObject(data, "todos", habit.Todos)

	var createdById string
	if habit.CreatedBy != nil {
		createdById = habit.CreatedBy.Id
	}
	d := map[string]interface{}{
		"_key":       habit.Id,
		"id":         habit.Id,
		"created":    habit.Created,
		"updated":    habit.Updated,
		"name":       habit.Title,
		"parameter1": habit.OrgId,
		"parameter2": createdById,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToHabit(r *db_proto.Record) (*behaviour_proto.Habit, error) {
	var p behaviour_proto.Habit
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func sharedGoalToRecord(from, to, orgId string, shared *behaviour_proto.ShareGoalUser, for_update bool) (string, error) {
	data, err := common.MarhalToObject(shared)
	if err != nil {
		return "", err
	}

	common.FilterObject(data, "goal", shared.Goal)
	common.FilterObject(data, "user", shared.User)
	common.FilterObject(data, "shared_by", shared.SharedBy)
	var sharedById string
	if shared.SharedBy != nil {
		sharedById = shared.SharedBy.Id
	}
	if for_update {
		//for update record, we want to delete the status as sharing it changes the current status for a user from current status to SHARED which we don't want
		delete(data, "status")
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

func recordToSharedGoal(r *db_proto.Record) (*behaviour_proto.ShareGoalUser, error) {
	var p behaviour_proto.ShareGoalUser
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func sharedChallengeToRecord(from, to, orgId string, shared *behaviour_proto.ShareChallengeUser, for_update bool) (string, error) {
	data, err := common.MarhalToObject(shared)
	if err != nil {
		return "", err
	}

	common.FilterObject(data, "challenge", shared.Challenge)
	common.FilterObject(data, "user", shared.User)
	common.FilterObject(data, "shared_by", shared.SharedBy)
	var sharedById string
	if shared.SharedBy != nil {
		sharedById = shared.SharedBy.Id
	}
	if for_update {
		//for update record, we want to delete the status as sharing it changes the current status for a user from current status to SHARED which we don't want
		delete(data, "status")
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

func recordToSharedChallenge(r *db_proto.Record) (*behaviour_proto.ShareChallengeUser, error) {
	var p behaviour_proto.ShareChallengeUser
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func sharedHabitToRecord(from, to, orgId string, shared *behaviour_proto.ShareHabitUser, for_update bool) (string, error) {
	data, err := common.MarhalToObject(shared)
	if err != nil {
		return "", err
	}

	common.FilterObject(data, "habit", shared.Habit)
	common.FilterObject(data, "user", shared.User)
	common.FilterObject(data, "shared_by", shared.SharedBy)
	var sharedById string
	if shared.SharedBy != nil {
		sharedById = shared.SharedBy.Id
	}
	if for_update {
		//for update record, we want to delete the status as sharing it changes the current status for a user from current status to SHARED which we don't want
		delete(data, "status")
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

func recordToGoalResponse(r *db_proto.Record) (*user_proto.GoalResponse, error) {
	var p user_proto.GoalResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToChallengeResponse(r *db_proto.Record) (*user_proto.ChallengeResponse, error) {
	var p user_proto.ChallengeResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToHabitResponse(r *db_proto.Record) (*user_proto.HabitResponse, error) {
	var p user_proto.HabitResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToSharedHabit(r *db_proto.Record) (*behaviour_proto.ShareHabitUser, error) {
	var p behaviour_proto.ShareHabitUser
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// AllGoals get all goals
func AllGoals(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*behaviour_proto.Goal, error) {
	var goals []*behaviour_proto.Goal
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET category = (FOR p IN %v FILTER doc.data.category.id == p._key RETURN p.data)
		LET target_aim = (FOR p IN %v FILTER doc.data.target.aim.id == p._key RETURN p.data)
		LET target_marker = (FOR p IN %v FILTER doc.data.target.marker.id == p._key RETURN p.data)
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		LET trackers = (
			FILTER NOT_NULL(doc.data.trackers)
			FOR tracker IN doc.data.trackers
				LET m = (FOR p IN %v FILTER tracker.marker.id == p._key RETURN p.data)
				LET t = (FOR p IN %v FILTER tracker.method.id == p._key RETURN p.data)
			RETURN MERGE_RECURSIVE(tracker, {
				marker: m[0],
				method: t[0]
			})
		)
		LET challenges = (
			FILTER NOT_NULL(doc.data.challenges)
			FOR c IN doc.data.challenges
			FOR p IN %v
			FILTER c.id == p._key RETURN p.data
		)
		LET habits = (
			FILTER NOT_NULL(doc.data.habits)
			FOR h IN doc.data.habits
			FOR p IN %v
			FILTER h.id == p._key RETURN p.data
		)
		LET triggers = (
			FILTER NOT_NULL(doc.data.triggers) 
			FOR t IN doc.data.triggers
			FOR p IN %v 
			FILTER t.id == p._key RETURN p.data
		)
		LET setbacks = (
			FILTER NOT_NULL(doc.data.setbacks) 
			FOR s IN doc.data.setbacks
			FOR p IN %v
			FILTER s.id == p._key RETURN p.data
		) 
		LET users = (
			FILTER NOT_NULL(doc.data.users)
			FOR target_user IN doc.data.users
				LET u = (FOR p IN %v FILTER target_user.user.id == p._key RETURN p.data)
			RETURN MERGE_RECURSIVE(target_user, {
				user: u[0]
			})
		)
		LET todo = (FOR p IN %v FILTER doc.data.todos.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{
		category:category[0],
		target:{aim:target_aim[0], marker:target_marker[0]},
		createdBy:createdBy[0],
		trackers: trackers,
		challenges: challenges,
		habits: habits,
		triggers:triggers,
		setbacks:setbacks,
		users:users,
		todos:todo[0]
		}})`, common.DbGoalTable, query, sort_query, limit_query,
		common.DbBehaviourCategoryTable, common.DbBehaviourCategoryAimTable, common.DbMarkerTable,
		common.DbUserTable, common.DbMarkerTable, common.DbTrackerMethodTable, common.DbChallengeTable,
		common.DbHabitTable, common.DbModuleTriggerTable, common.DbSetbackTable, common.DbUserTable,
		common.DbTodoTable,
	)

	resp, err := runQuery(ctx, q, common.DbGoalTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if goal, err := recordToGoal(r); err == nil {
			goals = append(goals, goal)
		}
	}
	return goals, nil
}

// AllChallenges get all challenges
func AllChallenges(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*behaviour_proto.Challenge, error) {
	var challenges []*behaviour_proto.Challenge
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET category = (FOR p IN %v FILTER doc.data.category.id == p._key RETURN p.data)
		LET target_aim = (FOR p IN %v FILTER doc.data.target.aim.id == p._key RETURN p.data)
		LET target_marker = (FOR p IN %v FILTER doc.data.target.marker.id == p._key RETURN p.data)
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		LET trackers = (
			FILTER NOT_NULL(doc.data.trackers)
			FOR tracker IN doc.data.trackers
				LET m = (FOR p IN %v FILTER tracker.marker.id == p._key RETURN p.data)
				LET t = (FOR p IN %v FILTER tracker.method.id == p._key RETURN p.data)
			RETURN MERGE_RECURSIVE(tracker, {
				marker: m[0],
				method: t[0]
			})
		)
		LET habits = (
			FILTER NOT_NULL(doc.data.habits)
			FOR h IN doc.data.habits
			FOR p IN %v
			FILTER h.id == p._key RETURN p.data
		)
		LET triggers = (
			FILTER NOT_NULL(doc.data.triggers) 
			FOR t IN doc.data.triggers
			FOR p IN %v 
			FILTER t.id == p._key RETURN p.data
		)
		LET setbacks = (
			FILTER NOT_NULL(doc.data.setbacks) 
			FOR s IN doc.data.setbacks
			FOR p IN %v
			FILTER s.id == p._key RETURN p.data
		) 
		LET users = (
			FILTER NOT_NULL(doc.data.users)
			FOR target_user IN doc.data.users
				LET u = (FOR p IN %v FILTER target_user.user.id == p._key RETURN p.data)
			RETURN MERGE_RECURSIVE(target_user, {
				user: u[0]
			})
		)
		LET todo = (FOR p IN %v FILTER doc.data.todos.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{
		category:category[0],
		target:{aim:target_aim[0], marker:target_marker[0]},
		createdBy:createdBy[0],
		trackers: trackers,
		habits: habits,
		triggers:triggers,
		setbacks:setbacks,
		users:users,
		todos:todo[0]
		}})`, common.DbChallengeTable, query, sort_query, limit_query,
		common.DbBehaviourCategoryTable, common.DbBehaviourCategoryAimTable, common.DbMarkerTable,
		common.DbUserTable, common.DbMarkerTable, common.DbTrackerMethodTable,
		common.DbHabitTable, common.DbModuleTriggerTable, common.DbSetbackTable, common.DbUserTable,
		common.DbTodoTable,
	)

	resp, err := runQuery(ctx, q, common.DbChallengeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if challenge, err := recordToChallenge(r); err == nil {
			challenges = append(challenges, challenge)
		}
	}
	return challenges, nil
}

// AllHabits get all habits
func AllHabits(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*behaviour_proto.Habit, error) {
	var habits []*behaviour_proto.Habit
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET category = (FOR p IN %v FILTER doc.data.category.id == p._key RETURN p.data)
		LET target_aim = (FOR p IN %v FILTER doc.data.target.aim.id == p._key RETURN p.data)
		LET target_marker = (FOR p IN %v FILTER doc.data.target.marker.id == p._key RETURN p.data)
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		LET trackers = (
			FILTER NOT_NULL(doc.data.trackers)
			FOR tracker IN doc.data.trackers
				LET m = (FOR p IN %v FILTER tracker.marker.id == p._key RETURN p.data)
				LET t = (FOR p IN %v FILTER tracker.method.id == p._key RETURN p.data)
			RETURN MERGE_RECURSIVE(tracker, {
				marker: m[0],
				method: t[0]
			})
		)
		LET triggers = (
			FILTER NOT_NULL(doc.data.triggers) 
			FOR t IN doc.data.triggers
			FOR p IN %v 
			FILTER t.id == p._key RETURN p.data
		)
		LET setbacks = (
			FILTER NOT_NULL(doc.data.setbacks) 
			FOR s IN doc.data.setbacks
			FOR p IN %v
			FILTER s.id == p._key RETURN p.data
		) 
		LET users = (
			FILTER NOT_NULL(doc.data.users)
			FOR target_user IN doc.data.users
				LET u = (FOR p IN %v FILTER target_user.user.id == p._key RETURN p.data)
			RETURN MERGE_RECURSIVE(target_user, {
				user: u[0]
			})
		)
		LET todo = (FOR p IN %v FILTER doc.data.todos.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{
		category:category[0],
		target:{aim:target_aim[0], marker:target_marker[0]},
		createdBy:createdBy[0],
		trackers: trackers,
		triggers:triggers,
		setbacks:setbacks,
		users:users,
		todos:todo[0]
		}})`, common.DbHabitTable, query, sort_query, limit_query,
		common.DbBehaviourCategoryTable, common.DbBehaviourCategoryAimTable, common.DbMarkerTable,
		common.DbUserTable, common.DbMarkerTable, common.DbTrackerMethodTable, common.DbModuleTriggerTable,
		common.DbSetbackTable, common.DbUserTable, common.DbTodoTable,
	)

	resp, err := runQuery(ctx, q, common.DbHabitTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if habit, err := recordToHabit(r); err == nil {
			habits = append(habits, habit)
		}
	}
	return habits, nil
}

// CreateGoal creates a goal
func CreateGoal(ctx context.Context, goal *behaviour_proto.Goal) error {
	if goal.Created == 0 {
		goal.Created = time.Now().Unix()
	}
	goal.Updated = time.Now().Unix()
	record, err := goalToRecord(goal)
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
		IN %v`, goal.Id, record, record, common.DbGoalTable)
	_, err = runQuery(ctx, q, common.DbGoalTable)
	return err
}

// CreateChallenge creates a challenge
func CreateChallenge(ctx context.Context, challenge *behaviour_proto.Challenge) error {
	if challenge.Created == 0 {
		challenge.Created = time.Now().Unix()
	}
	if challenge.Updated == 0 {
		challenge.Updated = time.Now().Unix()
	}
	record, err := challengeToRecord(challenge)
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
		IN %v`, challenge.Id, record, record, common.DbChallengeTable)
	_, err = runQuery(ctx, q, common.DbChallengeTable)
	return err
}

// CreateHabit creates a habit
func CreateHabit(ctx context.Context, habit *behaviour_proto.Habit) error {
	if habit.Created == 0 {
		habit.Created = time.Now().Unix()
	}
	if habit.Updated == 0 {
		habit.Updated = time.Now().Unix()
	}
	record, err := habitToRecord(habit)
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
		IN %v`, habit.Id, record, record, common.DbHabitTable)
	_, err = runQuery(ctx, q, common.DbHabitTable)
	return err
}

// ReadGoal reads a goal by ID
func ReadGoal(ctx context.Context, id, orgId, teamId string) (*behaviour_proto.Goal, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET category = (FOR p IN %v FILTER doc.data.category.id == p._key RETURN p.data)
		LET target_aim = (FOR p IN %v FILTER doc.data.target.aim.id == p._key RETURN p.data)
		LET target_marker = (FOR p IN %v FILTER doc.data.target.marker.id == p._key RETURN p.data)
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		LET trackers = (
			FILTER NOT_NULL(doc.data.trackers)
			FOR tracker IN doc.data.trackers
				LET m = (FOR p IN %v FILTER tracker.marker.id == p._key RETURN p.data)
				LET t = (FOR p IN %v FILTER tracker.method.id == p._key RETURN p.data)
			RETURN MERGE_RECURSIVE(tracker, {
				marker: m[0],
				method: t[0]
			})
		)
		LET challenges = (
			FILTER NOT_NULL(doc.data.challenges)
			FOR c IN doc.data.challenges
			FOR p IN %v
			FILTER c.id == p._key RETURN p.data
		)
		LET habits = (
			FILTER NOT_NULL(doc.data.habits)
			FOR h IN doc.data.habits
			FOR p IN %v
			FILTER h.id == p._key RETURN p.data
		)
		LET triggers = (
			FILTER NOT_NULL(doc.data.triggers) 
			FOR t IN doc.data.triggers
			FOR p IN %v 
			FILTER t.id == p._key RETURN p.data
		)
		LET setbacks = (
			FILTER NOT_NULL(doc.data.setbacks) 
			FOR s IN doc.data.setbacks
			FOR p IN %v
			FILTER s.id == p._key RETURN p.data
		) 
		LET users = (
			FILTER NOT_NULL(doc.data.users)
			FOR target_user IN doc.data.users
				LET u = (FOR p IN %v FILTER target_user.user.id == p._key RETURN p.data)
			RETURN MERGE_RECURSIVE(target_user, {
				user: u[0]
			})
		)
		LET todo = (FOR p IN %v FILTER doc.data.todos.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{
		category:category[0],
		target:{aim:target_aim[0], marker:target_marker[0]},
		createdBy:createdBy[0],
		trackers: trackers,
		challenges: challenges,
		habits: habits,
		triggers:triggers,
		setbacks:setbacks,
		users:users,
		todos:todo[0]
		}})`, common.DbGoalTable, query,
		common.DbBehaviourCategoryTable, common.DbBehaviourCategoryAimTable, common.DbMarkerTable,
		common.DbUserTable, common.DbMarkerTable, common.DbTrackerMethodTable, common.DbChallengeTable,
		common.DbHabitTable, common.DbModuleTriggerTable, common.DbSetbackTable, common.DbUserTable,
		common.DbTodoTable,
	)

	resp, err := runQuery(ctx, q, common.DbGoalTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToGoal(resp.Records[0])
	return data, err
}

// ReadChallenge reads a challenge by ID
func ReadChallenge(ctx context.Context, id, orgId, teamId string) (*behaviour_proto.Challenge, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET category = (FOR p IN %v FILTER doc.data.category.id == p._key RETURN p.data)
		LET target_aim = (FOR p IN %v FILTER doc.data.target.aim.id == p._key RETURN p.data)
		LET target_marker = (FOR p IN %v FILTER doc.data.target.marker.id == p._key RETURN p.data)
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		LET trackers = (
			FILTER NOT_NULL(doc.data.trackers)
			FOR tracker IN doc.data.trackers
				LET m = (FOR p IN %v FILTER tracker.marker.id == p._key RETURN p.data)
				LET t = (FOR p IN %v FILTER tracker.method.id == p._key RETURN p.data)
			RETURN MERGE_RECURSIVE(tracker, {
				marker: m[0],
				method: t[0]
			})
		)
		LET habits = (
			FILTER NOT_NULL(doc.data.habits)
			FOR h IN doc.data.habits
			FOR p IN %v
			FILTER h.id == p._key RETURN p.data
		)
		LET triggers = (
			FILTER NOT_NULL(doc.data.triggers) 
			FOR t IN doc.data.triggers
			FOR p IN %v 
			FILTER t.id == p._key RETURN p.data
		)
		LET setbacks = (
			FILTER NOT_NULL(doc.data.setbacks) 
			FOR s IN doc.data.setbacks
			FOR p IN %v
			FILTER s.id == p._key RETURN p.data
		) 
		LET users = (
			FILTER NOT_NULL(doc.data.users)
			FOR target_user IN doc.data.users
				LET u = (FOR p IN %v FILTER target_user.user.id == p._key RETURN p.data)
			RETURN MERGE_RECURSIVE(target_user, {
				user: u[0]
			})
		)
		LET todo = (FOR p IN %v FILTER doc.data.todos.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{
		category:category[0],
		target:{aim:target_aim[0], marker:target_marker[0]},
		createdBy:createdBy[0],
		trackers: trackers,
		habits: habits,
		triggers:triggers,
		setbacks:setbacks,
		users:users,
		todos:todo[0]
		}})`, common.DbChallengeTable, query,
		common.DbBehaviourCategoryTable, common.DbBehaviourCategoryAimTable, common.DbMarkerTable,
		common.DbUserTable, common.DbMarkerTable, common.DbTrackerMethodTable,
		common.DbHabitTable, common.DbModuleTriggerTable, common.DbSetbackTable, common.DbUserTable,
		common.DbTodoTable,
	)

	resp, err := runQuery(ctx, q, common.DbChallengeTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToChallenge(resp.Records[0])
	return data, err
}

// ReadHabit reads a habit by ID
func ReadHabit(ctx context.Context, id, orgId, teamId string) (*behaviour_proto.Habit, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET category = (FOR p IN %v FILTER doc.data.category.id == p._key RETURN p.data)
		LET target_aim = (FOR p IN %v FILTER doc.data.target.aim.id == p._key RETURN p.data)
		LET target_marker = (FOR p IN %v FILTER doc.data.target.marker.id == p._key RETURN p.data)
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		LET trackers = (
			FILTER NOT_NULL(doc.data.trackers)
			FOR tracker IN doc.data.trackers
				LET m = (FOR p IN %v FILTER tracker.marker.id == p._key RETURN p.data)
				LET t = (FOR p IN %v FILTER tracker.method.id == p._key RETURN p.data)
			RETURN MERGE_RECURSIVE(tracker, {
				marker: m[0],
				method: t[0]
			})
		)
		LET triggers = (
			FILTER NOT_NULL(doc.data.triggers) 
			FOR t IN doc.data.triggers
			FOR p IN %v 
			FILTER t.id == p._key RETURN p.data
		)
		LET setbacks = (
			FILTER NOT_NULL(doc.data.setbacks) 
			FOR s IN doc.data.setbacks
			FOR p IN %v
			FILTER s.id == p._key RETURN p.data
		) 
		LET users = (
			FILTER NOT_NULL(doc.data.users)
			FOR target_user IN doc.data.users
				LET u = (FOR p IN %v FILTER target_user.user.id == p._key RETURN p.data)
			RETURN MERGE_RECURSIVE(target_user, {
				user: u[0]
			})
		)
		LET todo = (FOR p IN %v FILTER doc.data.todos.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{
		category:category[0],
		target:{aim:target_aim[0], marker:target_marker[0]},
		createdBy:createdBy[0],
		trackers: trackers,
		triggers:triggers,
		setbacks:setbacks,
		users:users,
		todos:todo[0]
		}})`, common.DbHabitTable, query,
		common.DbBehaviourCategoryTable, common.DbBehaviourCategoryAimTable, common.DbMarkerTable,
		common.DbUserTable, common.DbMarkerTable, common.DbTrackerMethodTable, common.DbModuleTriggerTable,
		common.DbSetbackTable, common.DbUserTable, common.DbTodoTable,
	)

	resp, err := runQuery(ctx, q, common.DbHabitTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToHabit(resp.Records[0])
	return data, err
}

// DeleteGoal deletes a goal by ID
func DeleteGoal(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v		
		%s
		REMOVE doc IN %v`, common.DbGoalTable, query, common.DbGoalTable)
	_, err := runQuery(ctx, q, common.DbGoalTable)
	return err
}

// DeleteChallenge deletes a challenge by ID
func DeleteChallenge(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbChallengeTable, query, common.DbChallengeTable)
	_, err := runQuery(ctx, q, common.DbChallengeTable)
	return err
}

// DeleteHabit deletes a habit by ID
func DeleteHabit(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbHabitTable, query, common.DbHabitTable)
	_, err := runQuery(ctx, q, common.DbHabitTable)
	return err
}

func Filter(ctx context.Context, req *behaviour_proto.FilterRequest) (*behaviour_proto.FilterResponse_Data, error) {
	data := &behaviour_proto.FilterResponse_Data{}
	var goals []*behaviour_proto.Goal
	var challenges []*behaviour_proto.Challenge
	var habits []*behaviour_proto.Habit

	query := `FILTER`
	// search status
	if len(req.Status) > 0 {
		statuses := []string{}
		for _, s := range req.Status {
			statuses = append(statuses, fmt.Sprintf(`"%v"`, s))
		}
		query += fmt.Sprintf(" && doc.data.status IN [%v]", strings.Join(statuses[:], ","))
	}
	// search category
	if len(req.Category) > 0 {
		categories := common.QueryStringFromArray(req.Category)
		query += fmt.Sprintf(" && doc.data.category.id IN [%v]", categories)
	}
	// search creator
	if len(req.Creator) > 0 {
		creators := common.QueryStringFromArray(req.Creator)
		query += fmt.Sprintf(" && doc.data.createdBy.id IN [%v]", creators)
	}
	query = common.QueryAuth(query, req.OrgId, "")
	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)

	for _, t := range req.Type {
		switch t {
		case "goal":
			q := fmt.Sprintf(`
				FOR doc IN %v
				%s
				%s
				%s
				RETURN doc`, common.DbGoalTable, query, sort_query, limit_query)
			resp, err := runQuery(ctx, q, common.DbGoalTable)
			if err != nil {
				return nil, err
			}
			// parsing
			for _, r := range resp.Records {
				if goal, err := recordToGoal(r); err == nil {
					goals = append(goals, goal)
				}
			}
		case "challenge":
			q := fmt.Sprintf(`
				FOR doc IN %v
				%s
				%s
				%s
				RETURN doc`, common.DbChallengeTable, query, sort_query, limit_query)
			resp, err := runQuery(ctx, q, common.DbChallengeTable)
			if err != nil {
				return nil, err
			}
			// parsing
			for _, r := range resp.Records {
				if challenge, err := recordToChallenge(r); err == nil {
					challenges = append(challenges, challenge)
				}
			}
		case "habit":
			q := fmt.Sprintf(`
				FOR doc IN %v
				%s
				%s
				%s
				RETURN doc`, common.DbHabitTable, query, sort_query, limit_query)
			resp, err := runQuery(ctx, q, common.DbHabitTable)
			if err != nil {
				return nil, err
			}
			// parsing
			for _, r := range resp.Records {
				if habit, err := recordToHabit(r); err == nil {
					habits = append(habits, habit)
				}
			}
		}
	}

	data.Goals = goals
	data.Challenges = challenges
	data.Habits = habits
	return data, nil
}

func Search(ctx context.Context, req *behaviour_proto.SearchRequest) (*behaviour_proto.SearchResponse_Data, error) {
	data := &behaviour_proto.SearchResponse_Data{}
	var goals []*behaviour_proto.Goal
	var challenges []*behaviour_proto.Challenge
	var habits []*behaviour_proto.Habit

	query := `FILTER`
	if len(req.Name) > 0 {
		query += fmt.Sprintf(` && LIKE(doc.name, "%s",true)`, `%`+req.Name+`%`)
	}
	if len(req.Description) > 0 {
		query += fmt.Sprintf(` && LIKE(doc.data.description, "%v",true)`, `%`+req.Description+`%`)
	}
	if len(req.Summary) > 0 {
		query += fmt.Sprintf(` && LIKE(doc.data.summary, "%v",true)`, `%`+req.Summary+`%`)
	}
	query = common.QueryAuth(query, req.OrgId, "")
	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		RETURN doc`, common.DbGoalTable, query, sort_query, limit_query)
	resp, err := runQuery(ctx, q, common.DbGoalTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if goal, err := recordToGoal(r); err == nil {
			goals = append(goals, goal)
		}
	}

	q1 := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		RETURN doc`, common.DbChallengeTable, query, sort_query, limit_query)
	resp1, err := runQuery(ctx, q1, common.DbChallengeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp1.Records {
		if challenge, err := recordToChallenge(r); err == nil {
			challenges = append(challenges, challenge)
		}
	}

	q2 := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		RETURN doc`, common.DbHabitTable, query, sort_query, limit_query)
	resp2, err := runQuery(ctx, q2, common.DbHabitTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp2.Records {
		if habit, err := recordToHabit(r); err == nil {
			habits = append(habits, habit)
		}
	}
	data.Goals = goals
	data.Challenges = challenges
	data.Habits = habits
	return data, nil
}

func ShareGoal(ctx context.Context, goals []*behaviour_proto.Goal, users []*behaviour_proto.TargetedUser, sharedBy *user_proto.User, orgId string) ([]string, error) {
	userids := []string{}
	for _, goal := range goals {
		for _, user := range users {
			shared := &behaviour_proto.ShareGoalUser{
				Id:               uuid.NewUUID().String(),
				Goal:             goal,
				User:             user.User,
				Status:           static_proto.ShareStatus_SHARED,
				Updated:          time.Now().Unix(),
				Created:          time.Now().Unix(),
				SharedBy:         sharedBy,
				CurrentValue:     user.CurrentValue,
				ExpectedProgress: user.ExpectedProgress,
			}

			_from := fmt.Sprintf(`%v/%v`, common.DbGoalTable, goal.Id)
			_to := fmt.Sprintf(`%v/%v`, common.DbUserTable, user.User.Id)
			insert, err := sharedGoalToRecord(_from, _to, orgId, shared, false)
			update, err := sharedGoalToRecord(_from, _to, orgId, shared, true)
			if err != nil {
				return nil, err
			}
			if len(insert) == 0 || len(update) == 0 {
				return nil, errors.New("server serialization")
			}

			field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)
			q := fmt.Sprintf(`
				UPSERT %v
				INSERT %v
				UPDATE %v
				INTO %v
				RETURN {data:{user_id: OLD ? "" : NEW.data.user.id}}`, field, insert, update, common.DbShareGoalUserEdgeTable)

			resp, err := runQuery(ctx, q, common.DbShareGoalUserEdgeTable)
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

			any, err := common.FilteredAnyFromObject(common.GOAL_TYPE, goal.Id)
			if err != nil {
				log.Error("any from object err:", err)
				return nil, err
			}
			// save pending
			pending := &common_proto.Pending{
				Id:         uuid.NewUUID().String(),
				Created:    shared.Created,
				Updated:    shared.Updated,
				SharedBy:   sharedBy,
				SharedWith: user.User,
				Item:       any,
				OrgId:      orgId,
			}

			q1, err1 := common.SavePending(pending, goal.Id)
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

func ShareChallenge(ctx context.Context, challenges []*behaviour_proto.Challenge, users []*behaviour_proto.TargetedUser, sharedBy *user_proto.User, orgId string) ([]string, error) {
	userids := []string{}
	for _, challenge := range challenges {
		for _, user := range users {
			shared := &behaviour_proto.ShareChallengeUser{
				Id:               uuid.NewUUID().String(),
				Challenge:        challenge,
				User:             user.User,
				Status:           static_proto.ShareStatus_SHARED,
				Updated:          time.Now().Unix(),
				Created:          time.Now().Unix(),
				SharedBy:         sharedBy,
				CurrentValue:     user.CurrentValue,
				ExpectedProgress: user.ExpectedProgress,
			}

			_from := fmt.Sprintf(`%v/%v`, common.DbChallengeTable, challenge.Id)
			_to := fmt.Sprintf(`%v/%v`, common.DbUserTable, user.User.Id)
			insert, err := sharedChallengeToRecord(_from, _to, orgId, shared, false)
			update, err := sharedChallengeToRecord(_from, _to, orgId, shared, true)
			if err != nil {
				return nil, err
			}
			if len(insert) == 0 || len(update) == 0 {
				return nil, errors.New("server serialization")
			}

			field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)
			q := fmt.Sprintf(`
				UPSERT %v
				INSERT %v
				UPDATE %v
				INTO %v
				RETURN {data:{user_id: OLD ? "" : NEW.data.user.id}}`, field, insert, update, common.DbShareChallengeUserEdgeTable)

			resp, err := runQuery(ctx, q, common.DbShareChallengeUserEdgeTable)
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
			any, err := common.FilteredAnyFromObject(common.CHALLENGE_TYPE, challenge.Id)
			if err != nil {
				log.Error("any from object err:", err)
				return nil, err
			}
			// save pending
			pending := &common_proto.Pending{
				Id:         uuid.NewUUID().String(),
				Created:    shared.Created,
				Updated:    shared.Updated,
				SharedBy:   sharedBy,
				SharedWith: user.User,
				Item:       any,
				OrgId:      orgId,
			}

			q1, err1 := common.SavePending(pending, challenge.Id)
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

func ShareHabit(ctx context.Context, habits []*behaviour_proto.Habit, users []*behaviour_proto.TargetedUser, sharedBy *user_proto.User, orgId string) ([]string, error) {
	userids := []string{}
	for _, habit := range habits {
		for _, user := range users {
			shared := &behaviour_proto.ShareHabitUser{
				Id:               uuid.NewUUID().String(),
				Habit:            habit,
				User:             user.User,
				Status:           static_proto.ShareStatus_SHARED,
				Updated:          time.Now().Unix(),
				Created:          time.Now().Unix(),
				SharedBy:         sharedBy,
				CurrentValue:     user.CurrentValue,
				ExpectedProgress: user.ExpectedProgress,
			}

			_from := fmt.Sprintf(`%v/%v`, common.DbHabitTable, habit.Id)
			_to := fmt.Sprintf(`%v/%v`, common.DbUserTable, user.User.Id)
			insert, err := sharedHabitToRecord(_from, _to, orgId, shared, false)
			update, err := sharedHabitToRecord(_from, _to, orgId, shared, true)
			if err != nil {
				return nil, err
			}
			if len(insert) == 0 || len(update) == 0 {
				return nil, errors.New("server serialization")
			}

			field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)
			q := fmt.Sprintf(`
				UPSERT %v
				INSERT %v
				UPDATE %v
				INTO %v
				RETURN {data:{user_id: OLD ? "" : NEW.data.user.id}}`, field, insert, update, common.DbShareHabitUserEdgeTable)

			resp, err := runQuery(ctx, q, common.DbShareHabitUserEdgeTable)
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
			any, err := common.FilteredAnyFromObject(common.HABIT_TYPE, habit.Id)
			if err != nil {
				return nil, err
			}
			// save pending
			pending := &common_proto.Pending{
				Id:         uuid.NewUUID().String(),
				Created:    shared.Created,
				Updated:    shared.Updated,
				SharedBy:   sharedBy,
				SharedWith: user.User,
				Item:       any,
				OrgId:      orgId,
			}

			q1, err1 := common.SavePending(pending, habit.Id)
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

func AutocompleteGoalSearch(ctx context.Context, title string) ([]*static_proto.AutocompleteResponse, error) {
	response := []*static_proto.AutocompleteResponse{}
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER LIKE(doc.name, "%v",true)
		RETURN doc`, common.DbGoalTable, `%`+title+`%`)

	resp, err := runQuery(ctx, q, common.DbGoalTable)

	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		response = append(response, &static_proto.AutocompleteResponse{Id: r.Id, Title: r.Name, OrgId: r.Parameter1})
	}
	return response, nil
}

func AutocompleteChallengeSearch(ctx context.Context, title string) ([]*static_proto.AutocompleteResponse, error) {
	response := []*static_proto.AutocompleteResponse{}
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER LIKE(doc.name, "%v",true)
		RETURN doc`, common.DbChallengeTable, `%`+title+`%`)

	resp, err := runQuery(ctx, q, common.DbChallengeTable)

	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		response = append(response, &static_proto.AutocompleteResponse{Id: r.Id, Title: r.Name, OrgId: r.Parameter1})
	}
	return response, nil
}

func AutocompleteHabitSearch(ctx context.Context, title string) ([]*static_proto.AutocompleteResponse, error) {
	response := []*static_proto.AutocompleteResponse{}
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER LIKE(doc.name, "%v",true)
		RETURN doc`, common.DbHabitTable, `%`+title+`%`)

	resp, err := runQuery(ctx, q, common.DbHabitTable)

	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		response = append(response, &static_proto.AutocompleteResponse{Id: r.Id, Title: r.Name, OrgId: r.Parameter1})
	}
	return response, nil
}

func GetSharedGoal(ctx context.Context, userId, goalId string) (*behaviour_proto.ShareGoalUser, error) {
	query := fmt.Sprintf(`FILTER goal._key == "%v"`, goalId)

	q := fmt.Sprintf(`
		FOR goal,doc IN INBOUND "%v/%v" %v
		%v
		LET u = (FOR p IN %v FILTER doc.data.user.id == p._key RETURN p.data)
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {"data":{
			goal:goal.data,
			user:u[0],
			shared_by:sharedBy[0]
		}})`, common.DbUserTable, userId, common.DbShareGoalUserEdgeTable, query,
		common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareGoalUserEdgeTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, ErrNotFound
	}

	data, err := recordToSharedGoal(resp.Records[0])
	return data, err
}

func GetSharedChallenge(ctx context.Context, userId, challengeId string) (*behaviour_proto.ShareChallengeUser, error) {
	query := fmt.Sprintf(`FILTER challenge._key == "%v"`, challengeId)

	q := fmt.Sprintf(`
		FOR challenge,doc IN INBOUND "%v/%v" %v
		%v
		LET u = (FOR p IN %v FILTER doc.data.user.id == p._key RETURN p.data)
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {"data":{
			challenge:challenge.data,
			user:u[0],
			sharedBy:sharedBy[0]
		}})`, common.DbUserTable, userId, common.DbShareChallengeUserEdgeTable, query,
		common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareChallengeUserEdgeTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, ErrNotFound
	}

	return recordToSharedChallenge(resp.Records[0])
}

func GetSharedHabit(ctx context.Context, userId, habitId string) (*behaviour_proto.ShareHabitUser, error) {
	query := fmt.Sprintf(`FILTER habit._key == "%v"`, habitId)

	q := fmt.Sprintf(`
		FOR habit,doc IN INBOUND "%v/%v" %v
		%v
		LET u = (FOR p IN %v FILTER doc.data.user.id == p._key RETURN p.data)
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {"data":{
			habit:habit.data,
			user:u[0],
			shared_by:sharedBy[0]
		}})`, common.DbUserTable, userId, common.DbShareHabitUserEdgeTable, query,
		common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareHabitUserEdgeTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, ErrNotFound
	}

	return recordToSharedHabit(resp.Records[0])
}

func AutocompleteTags(ctx context.Context, orgId, object, name string) ([]string, error) {
	var tags []string

	var query string
	if len(orgId) > 0 {
		query = fmt.Sprintf(`FILTER doc.parameter1 == "%v"`, orgId)
	}

	var collection string
	switch object {
	case common.GOAL:
		collection = common.DbGoalTable
	case common.CHALLENGE:
		collection = common.DbChallengeTable
	case common.HABIT:
		collection = common.DbHabitTable
	}

	q := fmt.Sprintf(`
		LET tags = (FOR doc IN %v
		%v
		RETURN doc.data.tags)[**]
		FOR t IN tags
		FILTER LIKE(t,"%v",true)
		LET ret = {parameter1:t}
		RETURN DISTINCT ret
		`, collection, query, `%`+name+`%`)
	resp, err := runQuery(ctx, q, collection)

	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		tags = append(tags, r.Parameter1)
	}
	return tags, nil
}

func queryFilterSharedForUser(userId, shared_collection, joined_collection, org_query string) (string, string, string) {
	shared_query := ""
	joined_query := ""
	filter_shared_query := ""
	if len(userId) > 0 {
		shared_query = fmt.Sprintf(`LET shared = (
			FOR e, doc IN INBOUND "%v/%v" %v
			%s
			RETURN doc._from
		)`, common.DbUserTable, userId, shared_collection, org_query)
		joined_query = fmt.Sprintf(`LET joined = (
			FOR e, doc IN OUTBOUND "%v/%v" %v
			RETURN doc._to
		)`, common.DbUserTable, userId, joined_collection)
		filter_shared_query = "FILTER doc._id NOT IN shared && doc._id NOT IN joined"
	}
	return shared_query, joined_query, filter_shared_query
}

// AllGoalResponse get all goal responses
func AllGoalResponse(ctx context.Context, createdBy []string, search_term, userId, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*user_proto.GoalResponse, error) {
	var response []*user_proto.GoalResponse
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	shared_query, joined_query, filter_shared_query := queryFilterSharedForUser(userId, common.DbShareGoalUserEdgeTable, common.DbJoinGoalEdgeTable, query)

	filter_query := `FILTER`
	//filter by createdBy
	if len(createdBy) > 0 {
		createdBys := common.QueryStringFromArray(createdBy)
		filter_query += fmt.Sprintf(" && doc.data.createdBy.id IN [%v]", createdBys)
	}

	//filter by search term
	filter_query += common.QuerySharedResourceSearch(filter_query, search_term, "doc")

	q := fmt.Sprintf(`
		%s
		%s
		FOR doc IN %v
		%s
		%s
		%s
		%s
		%s
		LET target_aim = (FOR p IN %v FILTER doc.data.target.aim.id == p._key RETURN p.data)
		LET target_marker = (FOR p IN %v FILTER doc.data.target.marker.id == p._key RETURN p.data)
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.id,
			title:doc.name,
			image:doc.data.image,
			org_id: doc.data.org_id,
			summary: doc.data.summary,
			shared_by: {"id": createdBy[0].id, "firstname": createdBy[0].firstname, "lastname": createdBy[0].lastname, "avatar_url": createdBy[0].avatar_url},
			target:{
				aim:target_aim[0],
				marker:target_marker[0],
				targetValue:doc.data.target.targetValue,
				unit:doc.data.target.unit,
				recurrence:doc.data.target.recurrence
			},
			duration:doc.data.duration
		}}`,
		shared_query,
		joined_query,
		common.DbGoalTable,
		filter_shared_query,
		query,
		common.QueryClean(filter_query),
		sort_query, limit_query,
		common.DbBehaviourCategoryAimTable, common.DbMarkerTable,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbGoalTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if res, err := recordToGoalResponse(r); err == nil {
			response = append(response, res)
		}
	}
	return response, nil
}

// AllChallengeResponse get all challenge responses
func AllChallengeResponse(ctx context.Context, createdBy []string, search_term, userId, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*user_proto.ChallengeResponse, error) {
	var response []*user_proto.ChallengeResponse
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	shared_query, joined_query, filter_shared_query := queryFilterSharedForUser(userId, common.DbShareChallengeUserEdgeTable, common.DbJoinChallengeEdgeTable, query)

	filter_query := `FILTER`
	//filter by createdBy
	if len(createdBy) > 0 {
		createdBys := common.QueryStringFromArray(createdBy)
		filter_query += fmt.Sprintf(" && doc.data.createdBy.id IN [%v]", createdBys)
	}

	//filter by search term
	filter_query += common.QuerySharedResourceSearch(filter_query, search_term, "doc")

	q := fmt.Sprintf(`
		%s
		%s
		FOR doc IN %v
		%s
		%s
		%s
		%s
		%s
		LET target_aim = (FOR p IN %v FILTER doc.data.target.aim.id == p._key RETURN p.data)
		LET target_marker = (FOR p IN %v FILTER doc.data.target.marker.id == p._key RETURN p.data)
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.id,
			title:doc.name,
			image:doc.data.image,
			org_id: doc.data.org_id,
			summary: doc.data.summary,
			shared_by: {"id": createdBy[0].id, "firstname": createdBy[0].firstname, "lastname": createdBy[0].lastname, "avatar_url": createdBy[0].avatar_url},
			target:{
				aim:target_aim[0],
				marker:target_marker[0],
				targetValue:doc.data.target.targetValue,
				unit:doc.data.target.unit,
				recurrence:doc.data.target.recurrence
			},
			duration:doc.data.duration
		}}`,
		shared_query,
		joined_query,
		common.DbChallengeTable,
		filter_shared_query,
		query,
		common.QueryClean(filter_query),
		sort_query, limit_query,
		common.DbBehaviourCategoryAimTable, common.DbMarkerTable,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbChallengeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if res, err := recordToChallengeResponse(r); err == nil {
			response = append(response, res)
		}
	}
	return response, nil
}

// AllHabitResponse get all habit responses
func AllHabitResponse(ctx context.Context, createdBy []string, search_term, userId, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*user_proto.HabitResponse, error) {
	var response []*user_proto.HabitResponse
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	shared_query, joined_query, filter_shared_query := queryFilterSharedForUser(userId, common.DbShareHabitUserEdgeTable, common.DbJoinHabitEdgeTable, query)

	filter_query := `FILTER`
	//filter by createdBy
	if len(createdBy) > 0 {
		createdBys := common.QueryStringFromArray(createdBy)
		filter_query += fmt.Sprintf(" && doc.data.createdBy.id IN [%v]", createdBys)
	}

	//filter by search term
	filter_query += common.QuerySharedResourceSearch(filter_query, search_term, "doc")

	q := fmt.Sprintf(`
		%s
		%s
		FOR doc IN %v
		%s
		%s
		%s
		%s
		%s
		LET target_aim = (FOR p IN %v FILTER doc.data.target.aim.id == p._key RETURN p.data)
		LET target_marker = (FOR p IN %v FILTER doc.data.target.marker.id == p._key RETURN p.data)
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.id,
			title:doc.name,
			image:doc.data.image,
			org_id: doc.data.org_id,
			summary: doc.data.summary,
			shared_by: {"id": createdBy[0].id, "firstname": createdBy[0].firstname, "lastname": createdBy[0].lastname, "avatar_url": createdBy[0].avatar_url},
			target:{
				aim:target_aim[0],
				marker:target_marker[0],
				targetValue:doc.data.target.targetValue,
				unit:doc.data.target.unit,
				recurrence:doc.data.target.recurrence
			},
			duration:doc.data.duration
		}}`,
		shared_query,
		joined_query,
		common.DbHabitTable,
		filter_shared_query,
		query,
		common.QueryClean(filter_query),
		sort_query, limit_query,
		common.DbBehaviourCategoryAimTable, common.DbMarkerTable,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbHabitTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if res, err := recordToHabitResponse(r); err == nil {
			response = append(response, res)
		}
	}
	return response, nil
}
