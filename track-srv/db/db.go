package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"server/common"
	db_proto "server/db-srv/proto/db"
	track_proto "server/track-srv/proto/track"
	userapp_proto "server/user-app-srv/proto/userapp"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/client"
	"github.com/pborman/uuid"
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

func trackGoalToRecord(trackGoal *track_proto.TrackGoal) (string, string, error) {
	data, err := common.MarhalToObject(trackGoal)
	if err != nil {
		return "", "", err
	}
	key := uuid.NewUUID().String()

	common.FilterObject(data, "user", trackGoal.User)
	common.FilterObject(data, "goal", trackGoal.Goal)
	var userId string
	if trackGoal.User != nil {
		userId = trackGoal.User.Id
	}

	d := map[string]interface{}{
		"_key":       key,
		"created":    trackGoal.Created,
		"parameter1": trackGoal.OrgId,
		"parameter2": userId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", "", err
	}
	return key, string(body), err
}

func recordToTrackGoal(r *db_proto.Record) (*track_proto.TrackGoal, error) {
	var p track_proto.TrackGoal
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func trackChallengeToRecord(trackChallenge *track_proto.TrackChallenge) (string, string, error) {
	data, err := common.MarhalToObject(trackChallenge)
	if err != nil {
		return "", "", err
	}
	key := uuid.NewUUID().String()

	common.FilterObject(data, "user", trackChallenge.User)
	common.FilterObject(data, "challenge", trackChallenge.Challenge)
	var userId string
	if trackChallenge.User != nil {
		userId = trackChallenge.User.Id
	}

	d := map[string]interface{}{
		"_key":       key,
		"created":    trackChallenge.Created,
		"parameter1": trackChallenge.OrgId,
		"parameter2": userId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", "", err
	}
	return key, string(body), err
}

func recordToTrackChallenge(r *db_proto.Record) (*track_proto.TrackChallenge, error) {
	var p track_proto.TrackChallenge
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func trackHabitToRecord(trackHabit *track_proto.TrackHabit) (string, string, error) {
	data, err := common.MarhalToObject(trackHabit)
	if err != nil {
		return "", "", err
	}
	key := uuid.NewUUID().String()

	common.FilterObject(data, "user", trackHabit.User)
	common.FilterObject(data, "habit", trackHabit.Habit)
	var userId string
	if trackHabit.User != nil {
		userId = trackHabit.User.Id
	}

	d := map[string]interface{}{
		"_key":       key,
		"created":    trackHabit.Created,
		"parameter1": trackHabit.OrgId,
		"parameter2": userId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", "", err
	}
	return key, string(body), err
}

func recordToTrackHabit(r *db_proto.Record) (*track_proto.TrackHabit, error) {
	var p track_proto.TrackHabit
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func trackContentToRecord(trackContent *track_proto.TrackContent) (string, string, error) {
	data, err := common.MarhalToObject(trackContent)
	if err != nil {
		return "", "", err
	}
	key := uuid.NewUUID().String()

	common.FilterObject(data, "user", trackContent.User)
	common.FilterObject(data, "content", trackContent.Content)
	var userId string
	if trackContent.User != nil {
		userId = trackContent.User.Id
	}

	d := map[string]interface{}{
		"_key":       key,
		"created":    trackContent.Created,
		"parameter1": trackContent.OrgId,
		"parameter2": userId,
		"data":       trackContent,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", "", err
	}
	return key, string(body), err
}

func recordToTrackContent(r *db_proto.Record) (*track_proto.TrackContent, error) {
	var p track_proto.TrackContent
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func trackMarkerToRecord(trackMarker *track_proto.TrackMarker) (string, error) {
	data, err := common.MarhalToObject(trackMarker)
	if err != nil {
		return "", err
	}

	common.FilterObject(data, "user", trackMarker.User)
	common.FilterObject(data, "marker", trackMarker.Marker)
	var userId string
	if trackMarker.User != nil {
		userId = trackMarker.User.Id
	}

	d := map[string]interface{}{
		"_key":       trackMarker.Id,
		"id":         trackMarker.Id,
		"created":    trackMarker.Created,
		"parameter1": trackMarker.OrgId,
		"parameter2": userId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToTrackMarker(r *db_proto.Record) (*track_proto.TrackMarker, error) {
	var p track_proto.TrackMarker
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func queryPaginate(offset, limit int64) (string, string) {
	if limit == 0 {
		limit = 10
	}
	offs := fmt.Sprintf("%d", offset)
	size := fmt.Sprintf("%d", limit)
	return offs, size
}

func CreateTrackGoal(ctx context.Context, trackGoal *track_proto.TrackGoal) error {
	key, record, err := trackGoalToRecord(trackGoal)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		INSERT %v 
		INTO %v`, record, common.DbTrackGoalTable)
	_, err = runQuery(ctx, q, common.DbTrackGoalTable)
	if err != nil {
		return err
	}

	field := fmt.Sprintf(`{_from:"%v/%v",_to:"%v/%v"} `, common.DbTrackGoalTable, key, common.DbGoalTable, trackGoal.Goal.Id)
	q = fmt.Sprintf(`
			INSERT %v
			INTO %v`, field, common.DbTrackGoalEdgeTable)
	_, err = runQuery(ctx, q, common.DbTrackGoalEdgeTable)
	if err != nil {
		return err
	}
	return nil
}

func GetGoalCount(ctx context.Context, user_id, goal_id string, from, to int64) (int64, error) {
	query := ""
	if len(user_id) > 0 {
		query += fmt.Sprintf(` && doc.data.user.id == "%v"`, user_id)
	}
	if len(goal_id) > 0 {
		query += fmt.Sprintf(` && doc.data.goal.id == "%v"`, goal_id)
	}
	q := fmt.Sprintf(`
		LET date = ( FOR doc IN %v
			COLLECT AGGREGATE
				minDate = MIN(doc.created),
				maxDate = MAX(doc.created)
			RETURN { minDate, maxDate }
		)
		FOR doc IN %v
		FILTER date[0].minDate < doc.created && doc.created < date[0].maxDate %v
			RETURN doc`, common.DbTrackGoalTable, common.DbTrackGoalTable, query)

	resp, err := runQuery(ctx, q, common.DbTrackGoalTable)
	if err != nil {
		return 0, err
	}
	return int64(len(resp.Records)), nil
}

func GetGoalHistory(ctx context.Context, req *track_proto.GetGoalHistoryRequest) ([]*track_proto.TrackGoal, error) {
	var goals []*track_proto.TrackGoal
	query := `FILTER`
	if len(req.GoalId) > 0 {
		query += fmt.Sprintf(` doc.data.goal.id == "%v"`, req.GoalId)
	}
	if req.From > 0 {
		query += fmt.Sprintf(` && %d < doc.created`, req.From)
	}
	if req.To > 0 {
		query += fmt.Sprintf(` && doc.created < %d`, req.To)
	}

	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%v
		%v
		%v
		LET goal = (FOR p IN %v FILTER p._key == doc.data.goal.id RETURN p.data)
		LET user = (FOR p IN %v FILTER p._key == doc.data.user.id RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			goal:goal[0],
			user:user[0]
		}})`, common.DbTrackGoalTable, query, sort_query, limit_query,
		common.DbGoalTable, common.DbUserTable,
	)
	resp, err := runQuery(ctx, q, common.DbTrackGoalTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if goal, err := recordToTrackGoal(r); err == nil {
			goals = append(goals, goal)
		}
	}
	return goals, nil
}

func CreateTrackChallenge(ctx context.Context, trackChallenge *track_proto.TrackChallenge) error {
	key, record, err := trackChallengeToRecord(trackChallenge)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		INSERT %v 
		INTO %v`, record, common.DbTrackChallengeTable)
	_, err = runQuery(ctx, q, common.DbTrackChallengeTable)
	if err != nil {
		return err
	}

	field := fmt.Sprintf(`{_from:"%v/%v",_to:"%v/%v"} `, common.DbTrackChallengeTable, key, common.DbChallengeTable, trackChallenge.Challenge.Id)
	q = fmt.Sprintf(`
			INSERT %v
			INTO %v`, field, common.DbTrackChallengeEdgeTable)
	_, err = runQuery(ctx, q, common.DbTrackChallengeEdgeTable)
	if err != nil {
		return err
	}
	return nil
}

func GetChallengeCount(ctx context.Context, user_id, challenge_id string, from, to int64) (int64, error) {
	query := ""
	if len(user_id) > 0 {
		query += fmt.Sprintf(` && doc.data.user.id == "%v"`, user_id)
	}
	if len(challenge_id) > 0 {
		query += fmt.Sprintf(` && doc.data.challenge.id == "%v"`, challenge_id)
	}
	q := fmt.Sprintf(`
		LET date = ( FOR doc IN %v
			COLLECT AGGREGATE
				minDate = MIN(doc.created),
				maxDate = MAX(doc.created)
			RETURN { minDate, maxDate }
		)
		FOR doc IN %v
		FILTER date[0].minDate < doc.created && doc.created < date[0].maxDate %v
		RETURN doc`, common.DbTrackChallengeTable, common.DbTrackChallengeTable, query)
	resp, err := runQuery(ctx, q, common.DbTrackChallengeTable)
	if err != nil {
		return 0, err
	}
	return int64(len(resp.Records)), nil
}

func GetChallengeHistory(ctx context.Context, req *track_proto.GetChallengeHistoryRequest) ([]*track_proto.TrackChallenge, error) {
	var challenges []*track_proto.TrackChallenge
	query := `FILTER`
	if len(req.ChallengeId) > 0 {
		query += fmt.Sprintf(` doc.data.challenge.id == "%v"`, req.ChallengeId)
	}
	if req.From > 0 {
		query += fmt.Sprintf(` && %d < doc.created`, req.From)
	}
	if req.To > 0 {
		query += fmt.Sprintf(` && doc.created < %d`, req.To)
	}

	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%v
		%v
		%v
		LET challenge = (FOR p IN %v FILTER p._key == doc.data.challenge.id RETURN p.data)
		LET user = (FOR p IN %v FILTER p._key == doc.data.user.id RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			challenge:challenge[0],
			user:user[0]
		}})`, common.DbTrackChallengeTable, query, sort_query, limit_query,
		common.DbChallengeTable, common.DbUserTable,
	)
	resp, err := runQuery(ctx, q, common.DbTrackChallengeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if challenge, err := recordToTrackChallenge(r); err == nil {
			challenges = append(challenges, challenge)
		}
	}
	return challenges, nil
}

func CreateTrackHabit(ctx context.Context, trackHabit *track_proto.TrackHabit) error {
	key, record, err := trackHabitToRecord(trackHabit)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		INSERT %v 
		INTO %v`, record, common.DbTrackHabitTable)
	_, err = runQuery(ctx, q, common.DbTrackHabitTable)
	if err != nil {
		return err
	}

	field := fmt.Sprintf(`{_from:"%v/%v",_to:"%v/%v"} `, common.DbTrackHabitTable, key, common.DbHabitTable, trackHabit.Habit.Id)
	q = fmt.Sprintf(`
			INSERT %v
			INTO %v`, field, common.DbTrackHabitEdgeTable)
	_, err = runQuery(ctx, q, common.DbTrackHabitEdgeTable)
	if err != nil {
		return err
	}
	return nil
}

func GetHabitCount(ctx context.Context, user_id, habit_id string, from, to int64) (int64, error) {
	query := ""
	if len(user_id) > 0 {
		query += fmt.Sprintf(` && doc.data.user.id == "%v"`, user_id)
	}
	if len(habit_id) > 0 {
		query += fmt.Sprintf(` && doc.data.habit.id == "%v"`, habit_id)
	}
	q := fmt.Sprintf(`
		LET date = ( FOR doc IN %v
			COLLECT AGGREGATE
				minDate = MIN(doc.created),
				maxDate = MAX(doc.created)
			RETURN { minDate, maxDate }
		)
		FOR doc IN %v
		FILTER date[0].minDate < doc.created && doc.created < date[0].maxDate %v
		RETURN doc`, common.DbTrackHabitTable, common.DbTrackHabitTable, query)
	resp, err := runQuery(ctx, q, common.DbTrackHabitTable)
	if err != nil {
		return 0, err
	}
	return int64(len(resp.Records)), nil
}

func GetHabitHistory(ctx context.Context, req *track_proto.GetHabitHistoryRequest) ([]*track_proto.TrackHabit, error) {
	var habits []*track_proto.TrackHabit
	query := `FILTER`
	if len(req.HabitId) > 0 {
		query += fmt.Sprintf(` doc.data.habit.id == "%v"`, req.HabitId)
	}
	if req.From > 0 {
		query += fmt.Sprintf(` && %d < doc.created`, req.From)
	}
	if req.To > 0 {
		query += fmt.Sprintf(` && doc.created < %d`, req.To)
	}

	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)
	limit_query := common.QueryPaginate(req.Offset, req.Limit)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%v
		%v
		%v
		LET habit = (FOR p IN %v FILTER p._key == doc.data.habit.id RETURN p.data)
		LET user = (FOR p IN %v FILTER p._key == doc.data.user.id RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			habit:habit[0],
			user:user[0]
		}})`, common.DbTrackHabitTable, query, sort_query, limit_query,
		common.DbHabitTable, common.DbHabitTable,
	)
	resp, err := runQuery(ctx, q, common.DbTrackHabitTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if habit, err := recordToTrackHabit(r); err == nil {
			habits = append(habits, habit)
		}
	}
	return habits, nil
}

func CreateTrackContent(ctx context.Context, trackContent *track_proto.TrackContent) error {
	key, record, err := trackContentToRecord(trackContent)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		INSERT %v 
		INTO %v`, record, common.DbTrackContentTable)
	_, err = runQuery(ctx, q, common.DbTrackContentTable)
	if err != nil {
		return err
	}

	field := fmt.Sprintf(`{_from:"%v/%v",_to:"%v/%v"} `, common.DbTrackContentTable, key, common.DbContentTable, trackContent.Content.Id)
	q = fmt.Sprintf(`
			INSERT %v
			INTO %v`, field, common.DbTrackContentEdgeTable)
	_, err = runQuery(ctx, q, common.DbTrackContentEdgeTable)
	if err != nil {
		return err
	}
	return nil
}

func GetContentCount(ctx context.Context, user_id, content_id string, from, to int64) (int64, error) {
	query := ""
	if len(user_id) > 0 {
		query += fmt.Sprintf(` && doc.data.user.id == "%v"`, user_id)
	}
	if len(content_id) > 0 {
		query += fmt.Sprintf(` && doc.data.content.id == "%v"`, content_id)
	}
	q := fmt.Sprintf(`
		LET date = ( FOR doc IN %v
			COLLECT AGGREGATE
				minDate = MIN(doc.created),
				maxDate = MAX(doc.created)
			RETURN { minDate, maxDate }
		)
		FOR doc IN %v
		FILTER date[0].minDate < doc.created && doc.created < date[0].maxDate %v
		RETURN doc`, common.DbTrackContentTable, common.DbTrackContentTable, query)

	resp, err := runQuery(ctx, q, common.DbTrackContentTable)
	if err != nil {
		return 0, err
	}
	return int64(len(resp.Records)), nil
}

func GetContentHistory(ctx context.Context, req *track_proto.GetContentHistoryRequest) ([]*track_proto.TrackContent, error) {
	var contents []*track_proto.TrackContent
	query := `FILTER`
	if len(req.ContentId) > 0 {
		query += fmt.Sprintf(` doc.data.content.id == "%v"`, req.ContentId)
	}
	if req.From > 0 {
		query += fmt.Sprintf(` && %d < doc.created`, req.From)
	}
	if req.To > 0 {
		query += fmt.Sprintf(` && doc.created < %d`, req.To)
	}

	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%v
		%v
		%v
		LET content = (FOR p IN %v FILTER p._key == doc.data.content.id RETURN p.data)
		LET user = (FOR p IN %v FILTER p._key == doc.data.user.id RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			content:content[0],
			user:user[0]
		}})`, common.DbTrackContentTable, query, sort_query, limit_query,
		common.DbContentTable, common.DbUserTable,
	)
	resp, err := runQuery(ctx, q, common.DbTrackContentTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if content, err := recordToTrackContent(r); err == nil {
			contents = append(contents, content)
		}
	}
	return contents, nil
}

func CreateTrackMarker(ctx context.Context, trackerMarker *track_proto.TrackMarker) error {
	record, err := trackMarkerToRecord(trackerMarker)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		INSERT %v 
		INTO %v`, record, common.DbTrackMarkerTable)
	_, err = runQuery(ctx, q, common.DbTrackMarkerTable)
	if err != nil {
		return err
	}
	return nil
}

func GetLastMarker(ctx context.Context, markerId, userId string) (*track_proto.TrackMarker, error) {
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc.data.user.id == "%v" && doc.data.marker.id == "%v"
		SORT doc.created DESC
		LIMIT 0, 1
		LET marker = (FOR p IN %v FILTER p._key == doc.data.marker.id RETURN p.data)
		LET user = (FOR p IN %v FILTER p._key == doc.data.user.id RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			marker:marker[0],
			user:user[0]
		}})`, common.DbTrackMarkerTable, userId, markerId,
		common.DbMarkerTable, common.DbUserTable,
	)
	resp, err := runQuery(ctx, q, common.DbTrackMarkerTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}
	data, err := recordToTrackMarker(resp.Records[0])
	return data, err
}

func GetMarkerHistory(ctx context.Context, markerId string, from, to, offset, limit int64, sortParameter, sortDirection string) ([]*track_proto.TrackMarker, error) {
	markers := []*track_proto.TrackMarker{}

	sort_query := common.QuerySort(sortParameter, sortDirection)
	limit_query := common.QueryPaginate(offset, limit)

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc.data.marker.id == "%v"
		%s
		%s
		LET marker = (FOR p IN %v FILTER p._key == doc.data.marker.id RETURN p.data)
		LET user = (FOR p IN %v FILTER p._key == doc.data.user.id RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			marker:marker[0],
			user:user[0]
		}})`, common.DbTrackMarkerTable, markerId, sort_query, limit_query,
		common.DbMarkerTable, common.DbUserTable,
	)
	resp, err := runQuery(ctx, q, common.DbTrackMarkerTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if marker, err := recordToTrackMarker(r); err == nil {
			markers = append(markers, marker)
		}
	}
	return markers, nil
}

func GetAllMarkerHistory(ctx context.Context, from, to, offset, limit int64, sortParameter, sortDirection string) ([]*track_proto.TrackMarker, error) {
	markers := []*track_proto.TrackMarker{}

	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		LET marker = (FOR p IN %v FILTER p._key == doc.data.marker.id RETURN p.data)
		LET user = (FOR p IN %v FILTER p._key == doc.data.user.id RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			marker:marker[0],
			user:user[0]
		}})`, common.DbTrackMarkerTable, sort_query, limit_query,
		common.DbMarkerTable, common.DbUserTable,
	)
	resp, err := runQuery(ctx, q, common.DbTrackMarkerTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if marker, err := recordToTrackMarker(r); err == nil {
			markers = append(markers, marker)
		}
	}
	return markers, nil
}

func GetDefaultMarkerHistory(ctx context.Context, userId string, offset, limit, from, to int64) ([]*track_proto.TrackMarker, error) {
	query := fmt.Sprintf(`FILTER joined.data.status == "%v" || joined.data.status == "%v"`, userapp_proto.ActionStatus_STARTED, userapp_proto.ActionStatus_IN_PROGRESS)

	if from > 0 {
		query += fmt.Sprintf(` && %v <= joined.created`, from)
	}
	if to > 0 {
		query += fmt.Sprintf(` && joined.created <= %v`, to)
	}
	response := []*track_proto.TrackMarker{}
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort("created", "")

	q := fmt.Sprintf(`
		FOR goal, joined IN OUTBOUND "%v/%v" %v
		%s
		LET category = (FOR category in %v 
			filter category._key == goal.data.category.id
			return category)
		FOR doc IN %v
		FILTER doc.data.marker.id == category[0].data.markerDefault.id && doc.data.user.id == "%v"
		%s
		%s
		RETURN doc`,
		common.DbUserTable, userId, common.DbJoinGoalEdgeTable, query, common.DbBehaviourCategoryTable, common.DbTrackMarkerTable, userId, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbTrackMarkerTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if p, err := recordToTrackMarker(r); err == nil {
			response = append(response, p)
		} else {
			log.Println(p, err)
		}
	}
	return response, nil
}

func UpdateJoinGoalStatus(ctx context.Context, goalId string, status userapp_proto.ActionStatus) error {
	_to := fmt.Sprintf(`%v/%v`, common.DbGoalTable, goalId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._to == "%v"
		UPDATE doc WITH {data:{status:"%v"}} IN %v`, common.DbJoinGoalEdgeTable, _to, status, common.DbJoinGoalEdgeTable)
	_, err := runQuery(ctx, q, common.DbJoinGoalEdgeTable)
	return err
}

func UpdateJoinChallengeStatus(ctx context.Context, challengeId string, status userapp_proto.ActionStatus) error {
	_to := fmt.Sprintf(`%v/%v`, common.DbChallengeTable, challengeId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._to == "%v"
		UPDATE doc WITH {data:{status:"%v"}} IN %v`, common.DbJoinChallengeEdgeTable, _to, status, common.DbJoinChallengeEdgeTable)
	_, err := runQuery(ctx, q, common.DbJoinChallengeEdgeTable)
	return err
}

func UpdateJoinHabitStatus(ctx context.Context, habitId string, status userapp_proto.ActionStatus) error {
	_to := fmt.Sprintf(`%v/%v`, common.DbHabitTable, habitId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._to == "%v"
		UPDATE doc WITH {data:{status:"%v"}} IN %v`, common.DbJoinHabitEdgeTable, _to, status, common.DbJoinHabitEdgeTable)
	_, err := runQuery(ctx, q, common.DbJoinHabitEdgeTable)
	return err
}
