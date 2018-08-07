package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	account_proto "server/account-srv/proto/account"
	"server/common"
	db_proto "server/db-srv/proto/db"
	static_proto "server/static-srv/proto/static"
	user_proto "server/user-srv/proto/user"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
	query := fmt.Sprintf(`
		LET poc = (FOR poc in %v FILTER doc.data.pointOfContact.id == poc._key RETURN poc.data)
		LET b = (FOR b IN OUTBOUND doc._id %v RETURN b.data)
		LET pref = (FOR pref in %v FILTER doc.data.preference.id == pref._key RETURN pref.data)
		RETURN MERGE_RECURSIVE(doc,{data:{pointOfContact:poc[0],currentBatch:b[0], preference:pref[0]}})`,
		common.DbUserTable, common.DbUserBatchEdgeTable, common.DbPreferenceTable)
	return query
}

//This method return filter query for filtering status values
//ALERT! uses doc.data.status - changing the object name from doc to something else will cause error
func querySharedResourceStatus(filter_query string, status []static_proto.ShareStatus) string {
	if len(status) > 0 {
		statuses := []string{}
		for _, s := range status {
			statuses = append(statuses, fmt.Sprintf(`"%v"`, s))
		}
		filter_query += fmt.Sprintf(" && doc.data.status IN [%v]", strings.Join(statuses[:], ","))
	}
	return common.QueryClean(filter_query)
}

//This method return filter query for filtering shared_by values
//ALERT! uses doc.data.shared_by - changing the object name from doc to something else will cause error
func querySharedResourceSharedBy(filter_query string, sharedBy []string) string {
	if len(sharedBy) > 0 {
		sharedByIds := []string{}
		for _, s := range sharedBy {
			sharedByIds = append(sharedByIds, fmt.Sprintf(`"%v"`, s))
		}
		filter_query += fmt.Sprintf(" && doc.data.shared_by.id IN [%v]", strings.Join(sharedByIds[:], ","))
	}
	return common.QueryClean(filter_query)
}

func userToRecord(user *user_proto.User, options ...interface{}) (string, error) {
	data, err := common.MarhalToObject(user, options)
	if err != nil {
		return "", err
	}
	common.FilterObject(data, "pointOfContact", user.PointOfContact)
	delete(data, "currentBatch")
	common.FilterObject(data, "preference", user.Preference)

	d := map[string]interface{}{
		"_key":       user.Id,
		"id":         user.Id,
		"created":    user.Created,
		"updated":    user.Updated,
		"name":       user.Lastname,
		"parameter1": user.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToUser(r *db_proto.Record) (*user_proto.User, error) {
	var p user_proto.User
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}
func preferenceToRecord(p *user_proto.Preferences) (string, error) {
	data, err := common.MarhalToObject(p)
	if err != nil {
		return "", err
	}

	//currentMeasurements
	delete(data, "currentMeasurements")

	//conditions
	if len(p.Conditions) > 0 {
		var arr []interface{}
		for _, item := range p.Conditions {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["conditions"] = arr
	} else {
		delete(data, "conditions")
	}
	//allergies
	if len(p.Allergies) > 0 {
		var arr []interface{}
		for _, item := range p.Allergies {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["allergies"] = arr
	} else {
		delete(data, "allergies")
	}
	//food
	if len(p.Food) > 0 {
		var arr []interface{}
		for _, item := range p.Food {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["food"] = arr
	} else {
		delete(data, "food")
	}
	//cuisines
	if len(p.Cuisines) > 0 {
		var arr []interface{}
		for _, item := range p.Cuisines {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["cuisines"] = arr
	} else {
		delete(data, "cuisines")
	}
	//ethinicties
	if len(p.Ethinicties) > 0 {
		var arr []interface{}
		for _, item := range p.Ethinicties {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["ethinicties"] = arr
	} else {
		delete(data, "ethinicties")
	}

	d := map[string]interface{}{
		"_key":       p.Id,
		"id":         p.Id,
		"created":    p.Created,
		"updated":    p.Updated,
		"parameter1": p.OrgId,
		"parameter2": p.UserId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToPreference(r *db_proto.Record) (*user_proto.Preferences, error) {
	var p user_proto.Preferences
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToFeedback(r *db_proto.Record) (*user_proto.UserFeedback, error) {
	var p user_proto.UserFeedback
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToSearchResponse(r *db_proto.Record) (*user_proto.SearchResponse, error) {
	var p user_proto.SearchResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func userMeasurementEdgeToRecord(from, to string, p *user_proto.UserMeasurementEdge) (string, error) {
	data, err := common.MarhalToObject(p)
	if err != nil {
		return "", err
	}

	common.FilterObject(data, "measuredBy", p.MeasuredBy)
	common.FilterObject(data, "method", p.Method)

	d := map[string]interface{}{
		"_key":       p.Id,
		"_from":      from,
		"_to":        to,
		"id":         p.Id,
		"created":    p.Created,
		"updated":    p.Updated,
		"parameter1": p.OrgId,
		"parameter2": p.UserId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToMeasurement(r *db_proto.Record) (*user_proto.Measurement, error) {
	var p user_proto.Measurement
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToMarker(r *db_proto.Record) (*static_proto.Marker, error) {
	var p static_proto.Marker
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToSharedResource(r *db_proto.Record) (*user_proto.SharedResourcesResponse, error) {
	var p user_proto.SharedResourcesResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func All(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*user_proto.User, error) {
	var users []*user_proto.User
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)
	merge_query := queryMerge()
	q := fmt.Sprintf(`
	LET employees = (
		FOR doc IN %v
		%s
		RETURN doc._from
	)
	FOR doc IN %v
	FILTER doc._id NOT IN employees
	%s
	%s
	%s
	%s`,
		common.DbEmployeeTable, query, common.DbUserTable, query, sort_query, limit_query, merge_query)
	resp, err := runQuery(ctx, q, common.DbUserTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if user, err := recordToUser(r); err == nil {
			users = append(users, user)
		}
	}
	return users, nil
}

func Create(ctx context.Context, user *user_proto.User, account *account_proto.Account) error {
	// create user entity
	if user.Created == 0 {
		user.Created = time.Now().Unix()
	}
	user.Updated = time.Now().Unix()
	if user.Gender == 0 {
		user.Gender = static_proto.Gender_OTHER
	}

	//FIXME: maby we should only add preference when a normal user is being created - don'tneed this for employee
	// if preference is NOT provided, add an empty preference object and save it
	if user.Preference == nil {
		user.Preference = &user_proto.Preferences{
			Id:      uuid.NewUUID().String(),
			Created: time.Now().Unix(),
			Updated: time.Now().Unix(),
			UserId:  user.Id,
			OrgId:   user.OrgId,
		}
	}

	// if preference is provided, create new preference id
	if user.Preference != nil && len(user.Preference.Id) == 0 {
		user.Preference.Id = uuid.NewUUID().String()
		user.Preference.Created = time.Now().Unix()
		user.Preference.Updated = time.Now().Unix()
		user.Preference.UserId = user.Id
		user.Preference.OrgId = user.OrgId
	}
	user.Preference.OrgId = user.OrgId
	user.Preference.UserId = user.Id

	record, err := userToRecord(user)
	if err != nil {
		common.ErrorLog(common.UserSrv, Create, err, "Record parsing error")
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}
	// convert preference record
	record_preference, err := preferenceToRecord(user.Preference)
	if err != nil {
		common.ErrorLog(common.UserSrv, Create, err, "Preference record parsing error")
		return err
	}
	if len(record_preference) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
			INSERT %v IN %v
			UPSERT { _key: "%v" } 
			INSERT %v 
			UPDATE %v 
			IN %v`, record_preference, common.DbPreferenceTable, user.Id, record, record, common.DbUserTable)
	if _, err := runQuery(ctx, q, common.DbUserTable); err != nil {
		common.ErrorLog(common.UserSrv, Create, err, "User collection query running is failed")
		return err
	}

	// check validation user-account
	if account != nil && len(account.Id) > 0 {
		log.Debug("userid:", user.Id)
		log.Debug("accountid:", account.Id)
		q = fmt.Sprintf(`
		FOR doc IN OUTBOUND "%v/%v" %v
		RETURN doc`, common.DbUserTable, user.Id, common.DbUserAccountEdgeTable)
		if resp, err := runQuery(ctx, q, common.DbUserAccountEdgeTable); err == nil && len(resp.Records) > 0 {
			common.ErrorLog(common.UserSrv, Create, nil, "User has already one account")
			return errors.New("User has already one account")
		} else {
			// Create relationship between user and account
			field := fmt.Sprintf(`{_from:"%v/%v",_to:"%v/%v"} `, common.DbUserTable, user.Id, common.DbAccountTable, account.Id)
			record := fmt.Sprintf(`{_from:"%v/%v",_to:"%v/%v", parameter1:"%v",parameter2:"%v"} `, common.DbUserTable, user.Id, common.DbAccountTable, account.Id, user.Id, account.Id)
			q = fmt.Sprintf(`
			UPSERT %v
			INSERT %v
			UPDATE %v
			INTO %v`, field, record, record, common.DbUserAccountEdgeTable)
			if _, err := runQuery(ctx, q, common.DbUserAccountEdgeTable); err != nil {
				common.ErrorLog(common.UserSrv, Create, nil, "DbUserAccountEdgeTable query running is failed")
				return err
			}
		}
	}

	if len(user.OrgId) > 0 {
		// normal user
		field := fmt.Sprintf(`{_from:"%v/%v",_to:"%v/%v"} `, common.DbUserTable, user.Id, common.DbOrganisationTable, user.OrgId)
		q = fmt.Sprintf(`
		UPSERT %v
		INSERT %v
		UPDATE %v
		INTO %v`, field, field, field, common.DbUserOrgEdgeTable)
		if _, err := runQuery(ctx, q, common.DbUserOrgEdgeTable); err != nil {
			common.ErrorLog(common.UserSrv, Create, nil, "DbUserOrgEdgeTable query running is failed")
			return err
		}
	}

	if user.CurrentBatch != nil {
		batch := user.CurrentBatch
		field := fmt.Sprintf(`{_from:"%v/%v",_to:"%v/%v"} `, common.DbUserTable, user.Id, common.DbBatchTable, batch.Id)
		q = fmt.Sprintf(`
			UPSERT %v
			INSERT %v
			UPDATE %v
			INTO %v`, field, field, field, common.DbUserBatchEdgeTable)
		if _, err := runQuery(ctx, q, common.DbUserBatchEdgeTable); err != nil {
			common.ErrorLog(common.UserSrv, Create, nil, "DbUserBatchEdgeTable query running is failed")
			return err
		}
	}

	return nil
}
func Update(ctx context.Context, user *user_proto.User) error {
	if len(user.Id) == 0 {
		err := errors.New("Can't update without user_id")
		common.ErrorLog(common.UserSrv, Update, err, "Can't update without user_id")
		return err
	}
	user.Updated = time.Now().Unix()
	doNotEmitdefaults := true
	record, err := userToRecord(user, doNotEmitdefaults)
	if err != nil {
		common.ErrorLog(common.UserSrv, Update, err, "Record parsing error")
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}
	q := fmt.Sprintf(`
		FOR u IN %v
		UPDATE %v
		IN %v`, common.DbUserTable, record, common.DbUserTable)
	if _, err := runQuery(ctx, q, common.DbUserTable); err != nil {
		common.ErrorLog(common.UserSrv, Update, err, "Update User collection query running is failed")
		return err
	}

	return nil
}

func Read(ctx context.Context, userId, orgId string) (*user_proto.User, error) {
	merge_query := queryMerge()
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, userId)
	query = common.QueryAuth(query, orgId, "")
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s`, common.DbUserTable, query, merge_query)
	resp, err := runQuery(ctx, q, common.DbUserTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToUser(resp.Records[0])
	return data, err
}

func Filter(ctx context.Context, req *user_proto.FilterRequest) ([]*user_proto.User, error) {
	var users []*user_proto.User
	query := common.QueryAuth(`FILTER`, req.OrgId, "")
	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)
	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	merge_query := queryMerge()

	// search name
	if len(req.Users) > 0 {
		users := common.QueryStringFromArray(req.Users)
		query += fmt.Sprintf(" && doc.data.id IN [%v]", users)
	}

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbUserTable, query, sort_query, limit_query, merge_query)

	resp, err := runQuery(ctx, q, common.DbUserTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if user, err := recordToUser(r); err == nil {
			users = append(users, user)
		}
	}
	return users, nil
}

func Delete(ctx context.Context, userId string) error {
	queries := []string{}
	user_query := fmt.Sprintf(`LET user = (
		FOR doc IN %v
		FILTER doc._key == "%v"
		REMOVE doc IN %v
		RETURN OLD._key)`, common.DbUserTable, userId, common.DbUserTable)
	queries = append(queries, user_query)

	user_account_edge_query := fmt.Sprintf(`LET user_account_edge = (
		FOR doc IN %v
		FILTER doc._from == "user/%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbUserAccountEdgeTable, userId, common.DbUserAccountEdgeTable)
	queries = append(queries, user_account_edge_query)

	user_preference_query := fmt.Sprintf(`LET user_parameter2 = (
		FOR doc IN %v
		FILTER doc.parameter2 == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbPreferenceTable, userId, common.DbPreferenceTable)
	queries = append(queries, user_preference_query)

	share_goal_user_query := fmt.Sprintf(`LET share_goal_user = (
		FOR doc IN %v
		FILTER doc.data.shared_by == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbShareGoalUserEdgeTable, userId, common.DbShareGoalUserEdgeTable)
	queries = append(queries, share_goal_user_query)

	track_goal_query := fmt.Sprintf(`LET track_goal = (
		FOR doc IN %v
		FILTER doc.data.user.id == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbTrackGoalTable, userId, common.DbTrackGoalTable)
	queries = append(queries, track_goal_query)

	join_goal_edge_query := fmt.Sprintf(`LET join_goal_edge = (
		FOR doc IN %v
		FILTER doc._from == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbJoinGoalEdgeTable, userId, common.DbJoinGoalEdgeTable)
	queries = append(queries, join_goal_edge_query)

	share_challenge_user_query := fmt.Sprintf(`LET share_challenge_user = (
		FOR doc IN %v
		FILTER doc.data.shared_by == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbShareChallengeUserEdgeTable, userId, common.DbShareChallengeUserEdgeTable)
	queries = append(queries, share_challenge_user_query)

	track_challenge_query := fmt.Sprintf(`LET track_challenge = (
		FOR doc IN %v
		FILTER doc.data.user.id == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbTrackChallengeTable, userId, common.DbTrackChallengeTable)
	queries = append(queries, track_challenge_query)

	join_challenge_edge_query := fmt.Sprintf(`LET join_challenge_edge = (
		FOR doc IN %v
		FILTER doc._from == "user/%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbJoinChallengeEdgeTable, userId, common.DbJoinChallengeEdgeTable)
	queries = append(queries, join_challenge_edge_query)

	share_habit_user_query := fmt.Sprintf(`LET share_habit_user = (
		FOR doc IN %v
		FILTER doc.data.shared_by.id == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbShareHabitUserEdgeTable, userId, common.DbShareHabitUserEdgeTable)
	queries = append(queries, share_habit_user_query)

	track_habit_query := fmt.Sprintf(`LET track_habit = (
		FOR doc IN %v
		FILTER doc.data.user.id == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbTrackHabitTable, userId, common.DbTrackHabitTable)
	queries = append(queries, track_habit_query)

	join_habit_edge_query := fmt.Sprintf(`LET join_habit_edge = (
		FOR doc IN %v
		FILTER doc._from == "user/%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbJoinHabitEdgeTable, userId, common.DbJoinHabitEdgeTable)
	queries = append(queries, join_habit_edge_query)

	share_plan_user_query := fmt.Sprintf(`LET share_plan_user = (
		FOR doc IN %v
		FILTER doc.data.shared_by.id == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbSharePlanUserEdgeTable, userId, common.DbSharePlanUserEdgeTable)
	queries = append(queries, share_plan_user_query)

	share_content_user_query := fmt.Sprintf(`LET share_content_user = (
		FOR doc IN %v
		FILTER doc.data.shared_by.id == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbShareContentUserEdgeTable, userId, common.DbShareContentUserEdgeTable)
	queries = append(queries, share_content_user_query)

	track_marker_query := fmt.Sprintf(`LET track_marker = (
		FOR doc IN %v
		FILTER doc.data.user.id == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbTrackMarkerTable, userId, common.DbTrackMarkerTable)
	queries = append(queries, track_marker_query)

	track_goal_edge_query := fmt.Sprintf(`LET track_goal_edge = (
		FOR doc IN %v
		FILTER doc._from == "track_goal/%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbTrackGoalEdgeTable, userId, common.DbTrackGoalEdgeTable)
	queries = append(queries, track_goal_edge_query)

	track_challenge_edge_query := fmt.Sprintf(`LET track_challenge_edge = (
		FOR doc IN %v
		FILTER doc._from == "track_challenge/%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbTrackChallengeEdgeTable, userId, common.DbTrackChallengeEdgeTable)
	queries = append(queries, track_challenge_edge_query)

	track_habit_edge_query := fmt.Sprintf(`LET track_habit_edge = (
		FOR doc IN %v
		FILTER doc._from == "track_habit/%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbTrackHabitEdgeTable, userId, common.DbTrackHabitEdgeTable)
	queries = append(queries, track_habit_edge_query)

	track_content_edge_query := fmt.Sprintf(`LET track_content_edge = (
		FOR doc IN %v
		FILTER doc._from == "track_content/%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbTrackContentEdgeTable, userId, common.DbTrackContentEdgeTable)
	queries = append(queries, track_content_edge_query)

	note_query := fmt.Sprintf(`LET note = (
		FOR doc IN %v
		FILTER doc.data.user.id == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbNoteTable, userId, common.DbNoteTable)
	queries = append(queries, note_query)

	pending_query := fmt.Sprintf(`LET pending = (
		FOR doc IN %v
		FILTER doc.data.sharedWith.id == "%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbPendingTable, userId, common.DbPendingTable)
	queries = append(queries, pending_query)

	share_survey_user_query := fmt.Sprintf(`LET share_survey_user = (
		FOR doc IN %v
		FILTER doc._from == "survey/%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbShareSurveyUserEdgeTable, userId, common.DbShareSurveyUserEdgeTable)
	queries = append(queries, share_survey_user_query)

	user_measurement_edge_query := fmt.Sprintf(`LET user_measurement_edge = (
		FOR doc IN %v
		FILTER doc._from == "user/%v"
		REMOVE doc IN %v
		RETURN OLD._key`, common.DbUserMeasurementEdgeTable, userId, common.DbUserMeasurementEdgeTable)
	queries = append(queries, user_measurement_edge_query)

	var q string
	for _, query := range queries {
		q = fmt.Sprintf(`%s
			%s`, q, query)
	}

	_, err := runQuery(ctx, q, common.DbUserTable)
	return err
}

func ReadByAccount(ctx context.Context, accountId string) (*user_proto.User, error) {
	merge_query := queryMerge()
	q := fmt.Sprintf(`
			FOR doc IN INBOUND "%v/%v" %v
			%s`, common.DbAccountTable, accountId, common.DbUserAccountEdgeTable, merge_query)
	resp, err := runQuery(ctx, q, common.DbUserTable)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if len(resp.Records) == 0 {
		return nil, ErrNotFound
	}

	data, err := recordToUser(resp.Records[0])
	fmt.Println(err)
	return data, err
}

func UpdateTokens(ctx context.Context, userId string, tokens map[string]*user_proto.Token) error {
	body, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc.id == "%v"
		UPDATE doc WITH {data:{tokens:%v}} IN %v`, common.DbUserTable, userId, string(body), common.DbUserTable)
	_, err = runQuery(ctx, q, common.DbUserTable)
	return err
}

func ReadTokens(ctx context.Context, userIds []string) ([]*user_proto.Token, error) {
	var query string
	var tokens []*user_proto.Token
	// search users
	if len(userIds) > 0 {
		ids := []string{}
		for _, u := range userIds {
			ids = append(ids, `"`+u+`"`)
		}
		query = fmt.Sprintf(`FILTER doc.data.id IN [%v]`, strings.Join(ids[:], ","))
	}

	q := fmt.Sprintf(`
			FOR doc IN %v
			%v
			RETURN {data:{tokens: doc.data.tokens}}`, common.DbUserTable, query)

	resp, err := runQuery(ctx, q, common.DbUserTable)
	if err != nil || resp.Records[0] == nil {
		return nil, err
	}

	// parsing
	for _, r := range resp.Records {
		if user, err := recordToUser(r); err == nil {
			for _, t := range user.Tokens {
				tokens = append(tokens, t)
			}
		}
	}

	return tokens, nil
}

func SaveUserPreference(ctx context.Context, preference *user_proto.Preferences) error {
	// create user entity
	if len(preference.Id) == 0 {
		preference.Id = uuid.NewUUID().String()
	}
	if preference.Created == 0 {
		preference.Created = time.Now().Unix()
	}
	preference.Updated = time.Now().Unix()

	record, err := preferenceToRecord(preference)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { parameter1:"%v", parameter2:"%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, preference.OrgId, preference.UserId, record, record, common.DbPreferenceTable)

	if _, err := runQuery(ctx, q, common.DbUserTable); err != nil {
		return err
	}
	return nil
}

func ReadUserPreference(ctx context.Context, orgId, userId string) (*user_proto.Preferences, error) {
	query := fmt.Sprintf(`FILTER doc.parameter2 == "%v"`, userId)
	query = common.QueryAuth(query, orgId, "")
	log.Info(query)
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET currentMeasurements = ( 
			FOR measurement, ume IN OUTBOUND "%v/%v" %v
			RETURN {
				id:ume.data.id,
				org_id:ume.parameter1,
				user_id:ume.parameter2,
				created:ume.created,
				updated:ume.updated,
				marker:measurement.data.marker,
				method:ume.data.method,
				measuredBy:ume.data.measuredBy,
				value:measurement.data.value,
				unit:measurement.data.unit
		})
		LET conditions = ( 
			FILTER NOT_NULL(doc.data.conditions)
			FOR item IN doc.data.conditions
			FOR p IN %v
			FILTER item.id == p._key RETURN p.data)
		LET allergies = ( 
			FILTER NOT_NULL(doc.data.allergies)
			FOR item IN doc.data.allergies
			FOR p IN %v
			FILTER item.id == p._key RETURN p.data)
		LET food = ( 
			FILTER NOT_NULL(doc.data.food)
			FOR item IN doc.data.food
			FOR p IN %v
			FILTER item.id == p._key RETURN p.data)
		LET cuisines = ( 
			FILTER NOT_NULL(doc.data.cuisines)
			FOR item IN doc.data.cuisines
			FOR p IN %v
			FILTER item.id == p._key RETURN p.data)
		LET ethinicties = ( 
			FILTER NOT_NULL(doc.data.ethinicties)
			FOR item IN doc.data.ethinicties
			FOR p IN %v
			FILTER item.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{
			currentMeasurements:currentMeasurements,
			conditions:conditions,
			allergies:allergies,
			food:food,
			cuisines:cuisines,
			ethinicties:ethinicties
		}})`, common.DbPreferenceTable, query,
		common.DbUserTable, userId, common.DbUserMeasurementEdgeTable,
		common.DbContentCategoryItemTable,
		common.DbContentCategoryItemTable,
		common.DbContentCategoryItemTable,
		common.DbContentCategoryItemTable,
		common.DbContentCategoryItemTable)
	resp, err := runQuery(ctx, q, common.DbPreferenceTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToPreference(resp.Records[0])
	return data, err
}

func ListUserFeedback(ctx context.Context, userId string) ([]*user_proto.UserFeedback, error) {
	feedbacks := []*user_proto.UserFeedback{}

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc.parameter2 == "%v"
		RETURN doc`, common.DbUserFeedbackTable, userId)

	resp, err := runQuery(ctx, q, common.DbUserFeedbackTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if feedback, err := recordToFeedback(r); err == nil {
			feedbacks = append(feedbacks, feedback)
		}
	}
	return feedbacks, nil
}

func FilterUser(ctx context.Context, req *user_proto.FilterUserRequest) ([]*user_proto.SearchResponse, error) {
	response := []*user_proto.SearchResponse{}

	query := `FILTER`
	// add query of tags
	if len(req.Tags) > 0 {
		tags := common.QueryStringFromArray(req.Tags)
		query += fmt.Sprintf(" || doc.data.tags ANY IN [%v]", tags)
	}
	// add query of preference
	if req.Preference != nil {
		query += fmt.Sprintf(` || doc.data.preference.id == "%v"`, req.Preference.Id)
	}
	// add query of status
	if req.Status != 0 {
		query += fmt.Sprintf(` || account.data.status == "%v"`, req.Status)
	}
	query = strings.Replace(query, `FILTER || `, `FILTER `, -1)
	if query == `FILTER` {
		query = ""
	}
	query = common.QueryAuth(query, req.OrgId, "")
	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	q := fmt.Sprintf(`
		FOR doc IN %v
		FOR account IN OUTBOUND doc._id %v
		%v
		%s
		RETURN {data:{
			user_id: doc.id,
			org_id: doc.data.org_id,
			firstname: doc.data.firstname,
			lastname: doc.data.lastname,
			avatar_url: doc.data.avatar_url
		}}`, common.DbUserTable, common.DbUserAccountEdgeTable, query, limit_query)

	resp, err := runQuery(ctx, q, common.DbUserTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if res, err := recordToSearchResponse(r); err == nil {
			response = append(response, res)
		}
	}
	return response, nil
}

func SearchUser(ctx context.Context, req *user_proto.SearchUserRequest) ([]*user_proto.SearchResponse, error) {
	response := []*user_proto.SearchResponse{}
	search_query := `FILTER`
	org_query := `FILTER`
	org_query = common.QueryAuth(org_query, req.OrgId, "")
	// add query of tags
	if len(req.Name) > 0 {
		search_query += fmt.Sprintf(` || (LIKE(doc.data.firstname, "%s",true) || LIKE(doc.data.lastname, "%s",true))`, `%`+req.Name+`%`, `%`+req.Name+`%`)
	}
	// add query of preference
	if req.Gender != 0 {
		search_query += fmt.Sprintf(` || doc.data.gender == "%v"`, req.Gender)
	}
	// add query of status
	if len(req.Addresses) > 0 {
		addrs := []string{}
		for _, s := range req.Addresses {
			addrs = append(addrs, `"`+s.PostalCode+`"`)
		}
		search_query += fmt.Sprintf(" || doc.data.addresses[*].postalCode ANY IN [%v]", strings.Join(addrs[:], ","))
	}
	// add query of contactdetail
	if len(req.ContactDetails) != 0 {
		contacts := []string{}
		for _, s := range req.ContactDetails {
			contacts = append(contacts, `"`+s.Id+`"`)
		}
		search_query += fmt.Sprintf(" || doc.data.contactDetails[*].id ANY IN [%v]", strings.Join(contacts[:], ","))
	}
	search_query = common.QueryAuth(search_query, req.OrgId, "")
	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	q := fmt.Sprintf(`
		LET employees = (
			FOR doc IN %v
			%s
			RETURN doc._from
		)
		FOR doc IN %v
		FILTER doc._id NOT IN employees
		%v
		%s
		RETURN {data:{
			user_id: doc.id,
			org_id: doc.data.org_id,
			firstname: doc.data.firstname,
			lastname: doc.data.lastname,
			avatar_url: doc.data.avatar_url
		}}`,
		common.DbEmployeeTable, org_query,
		common.DbUserTable, search_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbUserTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if res, err := recordToSearchResponse(r); err == nil {
			response = append(response, res)
		}
	}
	return response, nil
}

func AutocompleteUser(ctx context.Context, req *user_proto.AutocompleteUserRequest) ([]*user_proto.SearchResponse, error) {
	response := []*user_proto.SearchResponse{}
	search_query := `FILTER`
	org_query := `FILTER`
	org_query = common.QueryAuth(org_query, req.OrgId, "")
	search_query = fmt.Sprintf(`FILTER LIKE(doc.data.firstname, "%s",true) || LIKE(doc.data.lastname, "%s",true)`, `%`+req.Name+`%`, `%`+req.Name+`%`)
	search_query = common.QueryAuth(search_query, req.OrgId, "")
	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)

	q := fmt.Sprintf(`
		LET employees = (
			FOR doc IN %v
			%s
			RETURN doc._from
		)
		FOR doc IN %v
		FILTER doc._id NOT IN employees
		%s
		%s
		%s
		RETURN {data:{
			user_id: doc.id,
			org_id: doc.data.org_id,
			firstname: doc.data.firstname,
			lastname: doc.data.lastname,
			avatar_url: doc.data.avatar_url
		}}`,
		common.DbEmployeeTable, org_query,
		common.DbUserTable, search_query, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbUserTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if res, err := recordToSearchResponse(r); err == nil {
			response = append(response, res)
		}
	}
	return response, nil
}

func AddMeasurement(ctx context.Context, measurement *user_proto.Measurement, track_marker_id string) (*user_proto.UserMeasurementEdge, error) {
	// save edge
	_from := fmt.Sprintf(`%v/%v`, common.DbUserTable, measurement.UserId)
	_to := fmt.Sprintf(`%v/%v`, common.DbTrackMarkerTable, track_marker_id)

	p := &user_proto.UserMeasurementEdge{
		Id:            uuid.NewUUID().String(),
		UserId:        measurement.UserId,
		OrgId:         measurement.OrgId,
		Created:       time.Now().Unix(),
		Updated:       time.Now().Unix(),
		Method:        measurement.Method,
		MeasuredBy:    measurement.MeasuredBy,
		TrackMarkerId: track_marker_id,
	}
	record, err := userMeasurementEdgeToRecord(_from, _to, p)
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
		INTO %v`, field, record, record, common.DbUserMeasurementEdgeTable)
	if _, err := runQuery(ctx, q, common.DbUserMeasurementEdgeTable); err != nil {
		return nil, err
	}
	return p, nil
}

//TODO: add sort parmater and direction
func GetAllMeasurementsHistory(ctx context.Context, userId, orgId string, offset, limit int64) ([]*user_proto.Measurement, error) {
	measurements := []*user_proto.Measurement{}
	query := `FILTER`
	query = common.QueryAuth(query, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort("", "")
	q := fmt.Sprintf(`
		FOR tc, doc IN OUTBOUND "%v/%v" %v
		%s
		%s
		%s
		LET marker = (FOR marker in %v FILTER tc.data.marker.id == marker._key RETURN marker.data)
		LET method = (FOR method in %v FILTER tc.data.method.id == method._key RETURN method.data)
		LET measuredBy = (FOR measuredBy in %v FILTER doc.data.measuredBy.id == measuredBy._key RETURN measuredBy.data)
		RETURN {data:{
				id:doc.id,
				org_id:doc.parameter1,
				user_id:doc.parameter2,
				created:doc.created,
				updated:doc.updated,
				marker:marker[0],
				method:method[0],
				measuredBy:measuredBy[0],
				value:tc.data.value,
				unit:tc.data.unit
		}}`,
		common.DbUserTable, userId, common.DbUserMeasurementEdgeTable,
		query,
		limit_query,
		sort_query,
		common.DbMarkerTable,
		common.DbTrackerMethodTable,
		common.DbUserTable)

	resp, err := runQuery(ctx, q, common.DbUserMeasurementEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if m, err := recordToMeasurement(r); err == nil {
			measurements = append(measurements, m)
		}
	}
	return measurements, nil
}

func GetMeasurementsHistory(ctx context.Context, userId, orgId, markerId string, offset, limit int64) ([]*user_proto.Measurement, error) {
	measurements := []*user_proto.Measurement{}
	query := fmt.Sprintf(`FILTER tc.data.marker.id == "%v"`, markerId)
	query = common.QueryAuth(query, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort("", "")
	q := fmt.Sprintf(`
		FOR tc, doc IN OUTBOUND "%v/%v" %v
		%s
		%s
		%s
		LET marker = (FOR marker in %v FILTER tc.data.marker.id == marker._key RETURN marker.data)
		LET method = (FOR method in %v FILTER tc.data.method.id == method._key RETURN method.data)
		LET measuredBy = (FOR measuredBy in %v FILTER doc.data.measuredBy.id == measuredBy._key RETURN measuredBy.data)
		RETURN {data:{
				id:doc.id,
				org_id:doc.parameter1,
				user_id:doc.parameter2,
				created:doc.created,
				updated:doc.updated,
				marker:marker[0],
				method:method[0],
				measuredBy:measuredBy[0],
				value:tc.data.value,
				unit:tc.data.unit
		}}`,
		common.DbUserTable, userId, common.DbUserMeasurementEdgeTable,
		query,
		limit_query,
		sort_query,
		common.DbMarkerTable,
		common.DbTrackerMethodTable,
		common.DbUserTable)

	resp, err := runQuery(ctx, q, common.DbUserMeasurementEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if m, err := recordToMeasurement(r); err == nil {
			measurements = append(measurements, m)
		}
	}
	return measurements, nil
}

func GetAllTrackedMarkers(ctx context.Context, userId, orgId string) ([]*static_proto.Marker, error) {
	markers := []*static_proto.Marker{}
	query := `FILTER`
	query = common.QueryAuth(query, orgId, "")
	q := fmt.Sprintf(`
		FOR tc, doc IN OUTBOUND "%v/%v" %v
		%s
		LET marker = (FOR m in %v FILTER tc.data.marker.id == m._key RETURN m.data)
		RETURN DISTINCT {data: marker[0]}`,
		common.DbUserTable, userId, common.DbUserMeasurementEdgeTable,
		query,
		common.DbMarkerTable)

	resp, err := runQuery(ctx, q, common.DbUserMeasurementEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if m, err := recordToMarker(r); err == nil {
			markers = append(markers, m)
		}
	}
	return markers, nil
}

func GetSharedGoalsForUser(ctx context.Context, status []static_proto.ShareStatus, sharedBy []string, search_term, userId, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*user_proto.SharedResourcesResponse, error) {
	shares := []*user_proto.SharedResourcesResponse{}
	doc_name := `goal`
	org_query := `FILTER`
	org_query = common.QueryAuth(org_query, orgId, "")

	filter_query := `FILTER`
	// filter by status
	filter_query = querySharedResourceStatus(filter_query, status)

	//filter by shared_by
	filter_query = querySharedResourceSharedBy(filter_query, sharedBy)

	//filter by search term
	filter_query = common.QuerySharedResourceSearch(filter_query, search_term, doc_name)

	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR %v, doc IN INBOUND "%v/%v" %v
		%s
		%s
		%s
		%s
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			resource_id:goal.data.id,
			image:goal.data.image,
			title:goal.data.title,
			org_id: goal.data.org_id,
			summary: goal.data.summary,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url},
			target:goal.data.target,
			current:doc.data.currentValue,
			duration:goal.data.duration}}`,
		doc_name,
		common.DbUserTable, userId, common.DbShareGoalUserEdgeTable,
		org_query,
		common.QueryClean(filter_query),
		limit_query,
		sort_query,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareGoalUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToSharedResource(r); err == nil {
			shares = append(shares, s)
		}
	}

	return shares, nil
}

func GetSharedChallengesForUser(ctx context.Context, status []static_proto.ShareStatus, sharedBy []string, search_term, userId, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*user_proto.SharedResourcesResponse, error) {
	shares := []*user_proto.SharedResourcesResponse{}
	doc_name := `challenge`
	org_query := `FILTER`
	org_query = common.QueryAuth(org_query, orgId, "")

	filter_query := `FILTER`
	// filter by status
	filter_query = querySharedResourceStatus(filter_query, status)

	//filter by shared_by
	filter_query = querySharedResourceSharedBy(filter_query, sharedBy)

	//filter by search term
	filter_query = common.QuerySharedResourceSearch(filter_query, search_term, doc_name)

	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR %v, doc IN INBOUND "%v/%v" %v
		%s
		%s
		%s
		%s
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			resource_id:challenge.data.id,
			image:challenge.data.image,
			title:challenge.data.title,
			org_id: challenge.data.org_id,
			summary: challenge.data.summary,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url},
			target:challenge.data.target,
			current:doc.data.currentValue,
			duration:challenge.data.duration}}`,
		doc_name,
		common.DbUserTable, userId, common.DbShareChallengeUserEdgeTable,
		org_query,
		common.QueryClean(filter_query),
		limit_query,
		sort_query,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareChallengeUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToSharedResource(r); err == nil {
			shares = append(shares, s)
		}
	}

	return shares, nil
}

func GetSharedHabitsForUser(ctx context.Context, status []static_proto.ShareStatus, sharedBy []string, search_term, userId, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*user_proto.SharedResourcesResponse, error) {
	shares := []*user_proto.SharedResourcesResponse{}
	doc_name := `habit`
	org_query := `FILTER`
	org_query = common.QueryAuth(org_query, orgId, "")

	filter_query := `FILTER`
	// filter by status
	filter_query = querySharedResourceStatus(filter_query, status)

	//filter by shared_by
	filter_query = querySharedResourceSharedBy(filter_query, sharedBy)

	//filter by search term
	filter_query = common.QuerySharedResourceSearch(filter_query, search_term, doc_name)

	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR %v, doc IN INBOUND "%v/%v" %v
		%s
		%s
		%s
		%s
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			resource_id:habit.data.id,
			image:habit.data.image,
			title:habit.data.title,
			org_id: habit.data.org_id,
			summary: habit.data.summary,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url},
			target:habit.data.target,
			current:doc.data.currentValue}}`,
		doc_name,
		common.DbUserTable, userId, common.DbShareHabitUserEdgeTable,
		org_query,
		common.QueryClean(filter_query),
		limit_query,
		sort_query,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareHabitUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToSharedResource(r); err == nil {
			shares = append(shares, s)
		}
	}
	return shares, nil
}
func GetSharedSurveysForUser(ctx context.Context, status []static_proto.ShareStatus, sharedBy []string, search_term, userId, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*user_proto.SharedResourcesResponse, error) {
	shares := []*user_proto.SharedResourcesResponse{}
	doc_name := `survey`
	org_query := `FILTER`
	org_query = common.QueryAuth(org_query, orgId, "")

	filter_query := `FILTER`
	// filter by status
	filter_query = querySharedResourceStatus(filter_query, status)

	//filter by shared_by
	filter_query = querySharedResourceSharedBy(filter_query, sharedBy)

	//filter by search term
	filter_query = common.QuerySharedResourceSearch(filter_query, search_term, doc_name)

	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR %v, doc IN INBOUND "%v/%v" %v
		%s
		%s
		%s
		%s
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			resource_id:survey.data.id,
			title:survey.data.title,
			org_id: survey.data.org_id,
			summary: survey.data.summary,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url},
			count:survey.data.count}}`,
		doc_name,
		common.DbUserTable, userId, common.DbShareSurveyUserEdgeTable,
		org_query,
		common.QueryClean(filter_query),
		limit_query,
		sort_query,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareSurveyUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToSharedResource(r); err == nil {
			shares = append(shares, s)
		}
	}
	return shares, nil
}
func GetSharedContentsForUser(ctx context.Context, status []static_proto.ShareStatus, sharedBy []string, search_term, userId, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*user_proto.SharedResourcesResponse, error) {
	shares := []*user_proto.SharedResourcesResponse{}
	doc_name := `content`
	org_query := `FILTER`
	org_query = common.QueryAuth(org_query, orgId, "")

	filter_query := `FILTER`
	// filter by status
	filter_query = querySharedResourceStatus(filter_query, status)

	//filter by shared_by
	filter_query = querySharedResourceSharedBy(filter_query, sharedBy)

	//filter by search term
	filter_query = common.QuerySharedResourceSearch(filter_query, search_term, doc_name)

	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR %v, doc IN INBOUND "%v/%v" %v
		%s
		%s
		%s
		%s
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			resource_id:content.data.id,
			image:content.data.image,
			title:content.data.title,
			org_id: content.data.org_id,
			summary: content.data.summary,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url}
			}}`,
		doc_name,
		common.DbUserTable, userId, common.DbShareContentUserEdgeTable,
		org_query,
		common.QueryClean(filter_query),
		limit_query,
		sort_query,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareContentUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToSharedResource(r); err == nil {
			shares = append(shares, s)
		}
	}
	return shares, nil
}
