package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	content_proto "server/content-srv/proto/content"
	db_proto "server/db-srv/proto/db"
	common_proto "server/static-srv/proto/common"
	static_proto "server/static-srv/proto/static"
	track_proto "server/track-srv/proto/track"
	userapp_proto "server/user-app-srv/proto/userapp"
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

func bookmarkToRecord(orgId, userId, contentId string, bookmark *userapp_proto.Bookmark) (string, error) {
	_from := fmt.Sprintf(`%v/%v`, common.DbUserTable, userId)
	_to := fmt.Sprintf(`%v/%v`, common.DbContentTable, contentId)

	data, err := common.MarhalToObject(bookmark)
	if err != nil {
		return "", err
	}
	delete(data, "user")
	delete(data, "content")
	delete(data, "content_category")

	d := map[string]interface{}{
		"_key":       bookmark.Id,
		"_from":      _from,
		"_to":        _to,
		"id":         bookmark.Id,
		"created":    bookmark.Created,
		"parameter1": orgId,
		"parameter2": userId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToBookmark(r *db_proto.Record) (*userapp_proto.Bookmark, error) {
	var p userapp_proto.Bookmark
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToBookmarkContent(r *db_proto.Record) (*userapp_proto.BookmarkContent, error) {
	var p userapp_proto.BookmarkContent
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToSharedContent(r *db_proto.Record) (*content_proto.SharedContent, error) {
	var p content_proto.SharedContent
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToShareableContent(r *db_proto.Record) (*user_proto.ShareableContent, error) {
	var p user_proto.ShareableContent
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToContentDetail(r *db_proto.Record) (*userapp_proto.ContentDetail, error) {
	var p userapp_proto.ContentDetail
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToSharedPlan(r *db_proto.Record) (*userapp_proto.SharedPlan, error) {
	var p userapp_proto.SharedPlan
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToSharedSurvey(r *db_proto.Record) (*userapp_proto.SharedSurvey, error) {
	var p userapp_proto.SharedSurvey
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToSharedGoal(r *db_proto.Record) (*userapp_proto.SharedGoal, error) {
	var p userapp_proto.SharedGoal
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToSharedChallenge(r *db_proto.Record) (*userapp_proto.SharedChallenge, error) {
	var p userapp_proto.SharedChallenge
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToSharedHabit(r *db_proto.Record) (*userapp_proto.SharedHabit, error) {
	var p userapp_proto.SharedHabit
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToTrackMarker(r *db_proto.Record) (*track_proto.TrackMarker, error) {
	var p track_proto.TrackMarker
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func joingoalToRecord(from, to, userId, goalId string, join_goal *userapp_proto.JoinGoal) (string, error) {
	data, err := common.MarhalToObject(join_goal)
	if err != nil {
		return "", err
	}

	delete(data, "goal")
	delete(data, "user")
	delete(data, "target")

	d := map[string]interface{}{
		"_key":       join_goal.Id,
		"_from":      from,
		"_to":        to,
		"id":         join_goal.Id,
		"created":    join_goal.Start,
		"parameter1": userId,
		"parameter2": goalId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToJoinGoal(r *db_proto.Record) (*userapp_proto.JoinGoal, error) {
	var p userapp_proto.JoinGoal
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToGoalResponse(r *db_proto.Record) (*userapp_proto.JoinGoalResponse, error) {
	var p userapp_proto.JoinGoalResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func joinchallengeToRecord(from, to, userId, challengeId string, join_challenge *userapp_proto.JoinChallenge) (string, error) {
	data, err := common.MarhalToObject(join_challenge)
	if err != nil {
		return "", err
	}

	delete(data, "challenge")
	delete(data, "user")
	delete(data, "target")

	d := map[string]interface{}{
		"_key":       join_challenge.Id,
		"_from":      from,
		"_to":        to,
		"id":         join_challenge.Id,
		"created":    join_challenge.Start,
		"parameter1": userId,
		"parameter2": challengeId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToJoinChallenge(r *db_proto.Record) (*userapp_proto.JoinChallenge, error) {
	var p userapp_proto.JoinChallenge
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToChallengeResponse(r *db_proto.Record) (*userapp_proto.JoinChallengeResponse, error) {
	var p userapp_proto.JoinChallengeResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func joinhabitToRecord(from, to, userId, habitId string, join_habit *userapp_proto.JoinHabit) (string, error) {
	data, err := common.MarhalToObject(join_habit)
	if err != nil {
		return "", err
	}

	delete(data, "habit")
	delete(data, "user")
	delete(data, "target")

	d := map[string]interface{}{
		"_key":       join_habit.Id,
		"_from":      from,
		"_to":        to,
		"id":         join_habit.Id,
		"created":    join_habit.Start,
		"parameter1": userId,
		"parameter2": habitId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToJoinHabit(r *db_proto.Record) (*userapp_proto.JoinHabit, error) {
	var p userapp_proto.JoinHabit
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToHabitResponse(r *db_proto.Record) (*userapp_proto.JoinHabitResponse, error) {
	var p userapp_proto.JoinHabitResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToGoalProgressResponse(r *db_proto.Record) (*userapp_proto.GoalProgressResponse, error) {
	var p userapp_proto.GoalProgressResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToChallengeCountResponse(r *db_proto.Record) (*userapp_proto.ChallengeCountResponse, error) {
	var p userapp_proto.ChallengeCountResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToHabitCountResponse(r *db_proto.Record) (*userapp_proto.HabitCountResponse, error) {
	var p userapp_proto.HabitCountResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToPending(r *db_proto.Record) (*common_proto.PendingResponse, error) {
	var p common_proto.PendingResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToUserPlan(r *db_proto.Record) (*userapp_proto.UserPlan, error) {
	var p userapp_proto.UserPlan
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToCategoryCountResponse(r *db_proto.Record) (*userapp_proto.CategoryCountResponse, error) {
	var p userapp_proto.CategoryCountResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func contentRatingToRecord(from, to string, contentRating *userapp_proto.ContentRating) (string, error) {
	data, err := common.MarhalToObject(contentRating)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_from":      from,
		"_to":        to,
		"created":    contentRating.Created,
		"updated":    contentRating.Updated,
		"parameter1": contentRating.OrgId,
		"parameter2": contentRating.UserId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func contentDislikeToRecord(from, to string, contentDislike *userapp_proto.ContentDislike) (string, error) {
	data, err := common.MarhalToObject(contentDislike)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_from":      from,
		"_to":        to,
		"created":    contentDislike.Created,
		"updated":    contentDislike.Updated,
		"parameter1": contentDislike.OrgId,
		"parameter2": contentDislike.UserId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func contentDislikeSimilarToRecord(from, to string, contentDislikeSimilar *userapp_proto.ContentDislikeSimilar) (string, error) {
	data, err := common.MarhalToObject(contentDislikeSimilar)
	if err != nil {
		return "", err
	}
	// tags
	if len(contentDislikeSimilar.Tags) > 0 {
		var arr []interface{}
		for _, item := range contentDislikeSimilar.Tags {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["tags"] = arr
	} else {
		delete(data, "tags")
	}

	d := map[string]interface{}{
		"_from":      from,
		"_to":        to,
		"created":    contentDislikeSimilar.Created,
		"updated":    contentDislikeSimilar.Updated,
		"parameter1": contentDislikeSimilar.OrgId,
		"parameter2": contentDislikeSimilar.UserId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func feedbackToRecord(feedback *user_proto.UserFeedback) (string, error) {
	data, err := common.MarhalToObject(feedback)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":       feedback.Id,
		"id":         feedback.Id,
		"created":    feedback.Created,
		"updated":    feedback.Updated,
		"parameter1": feedback.OrgId,
		"parameter2": feedback.UserId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func userplanToRecord(userplan *userapp_proto.UserPlan) (string, error) {
	data, err := common.MarhalToObject(userplan)
	if err != nil {
		return "", err
	}

	// plan
	common.FilterObject(data, "plan", userplan.Plan)
	// goals
	if len(userplan.Goals) > 0 {
		var arr []interface{}
		for _, item := range userplan.Goals {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["goals"] = arr
	} else {
		delete(data, "goals")
	}
	//creator
	common.FilterObject(data, "creator", userplan.Creator)

	d := map[string]interface{}{
		"_key":       userplan.Id,
		"id":         userplan.Id,
		"created":    userplan.Created,
		"updated":    userplan.Updated,
		"parameter1": userplan.OrgId,
		"parameter2": userplan.Creator.Id,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToCategory(r *db_proto.Record) (*static_proto.BehaviourCategory, error) {
	var p static_proto.BehaviourCategory
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToGoalDetail(r *db_proto.Record) (*userapp_proto.GoalDetail, error) {
	var p userapp_proto.GoalDetail
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToChallengeDetail(r *db_proto.Record) (*userapp_proto.ChallengeDetail, error) {
	var p userapp_proto.ChallengeDetail
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToHabitDetail(r *db_proto.Record) (*userapp_proto.HabitDetail, error) {
	var p userapp_proto.HabitDetail
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

//if upsert results in an update => false / insert => true
func recordToCheckUpsertOperation(r *db_proto.Record) (bool, error) {
	var p userapp_proto.UpsertOperation
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return false, err
	}
	if p.Type == "update" {
		return false, nil
	} else if p.Type == "insert" {
		return true, nil
	}
	return false, nil
}

func CreateBookmark(ctx context.Context, orgId, userId, contentId string) (*userapp_proto.Bookmark, error) {
	bookmark := &userapp_proto.Bookmark{
		Id:      uuid.NewUUID().String(),
		Created: time.Now().Unix(),
	}
	record, err := bookmarkToRecord(orgId, userId, contentId, bookmark)
	if err != nil {
		return nil, err
	}
	if len(record) == 0 {
		return nil, errors.New("server serialization")
	}

	field := fmt.Sprintf(`{_from:"%v/%v",_to:"%v/%v"}`, common.DbUserTable, userId, common.DbContentTable, contentId)
	q := fmt.Sprintf(`
		UPSERT %v
		INSERT %v
		UPDATE %v
		INTO %v`, field, record, record, common.DbBookmarkEdgeTable)

	if _, err := runQuery(ctx, q, common.DbBookmarkEdgeTable); err != nil {
		return nil, err
	}
	return bookmark, nil
}

func ReadBookmarkContents(ctx context.Context, userId string) ([]*userapp_proto.BookmarkContent, error) {
	bookmarkContents := []*userapp_proto.BookmarkContent{}
	_from := fmt.Sprintf("%v/%v", common.DbUserTable, userId)
	q := fmt.Sprintf(`
		FOR content,doc IN OUTBOUND "%v" %v
		LET category = (FOR p IN %v FILTER content.data.category.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.id, 
			content_id:content.id, 
			image:content.data.image, 
			title:content.data.title, 
			category:category[0].name, 
			category_id:content.data.category.id}}`,
		_from, common.DbBookmarkEdgeTable, common.DbContentCategoryTable)

	resp, err := runQuery(ctx, q, common.DbBookmarkEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if b, err := recordToBookmarkContent(r); err == nil {
			bookmarkContents = append(bookmarkContents, b)
		}
	}
	return bookmarkContents, nil
}

func ReadBookmarkContentCategorys(ctx context.Context, userId string) ([]*userapp_proto.BookmarkContent, error) {
	categorys := []*userapp_proto.BookmarkContent{}
	_from := fmt.Sprintf("%v/%v", common.DbUserTable, userId)
	q := fmt.Sprintf(`
		FOR content,doc IN OUTBOUND "%v" %v
		LET category = (FOR p IN %v FILTER content.data.category.id == p._key RETURN p.data)
		RETURN DISTINCT {data:{
			category:category[0].name, 
			category_id:category[0].id}}`,
		_from, common.DbBookmarkEdgeTable, common.DbContentCategoryTable,
	)

	resp, err := runQuery(ctx, q, common.DbBookmarkEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if b, err := recordToBookmarkContent(r); err == nil {
			categorys = append(categorys, b)
		}
	}
	return categorys, nil
}

func ReadBookmarkByCategory(ctx context.Context, userId, categoryId string) ([]*userapp_proto.BookmarkContent, error) {
	bookmarkContents := []*userapp_proto.BookmarkContent{}
	_from := fmt.Sprintf("%v/%v", common.DbUserTable, userId)
	q := fmt.Sprintf(`
		FOR content, doc IN OUTBOUND "%v" %v
		FILTER content.data.category.id == "%v"
		LET category = (FOR p IN %v FILTER content.data.category.id == p._key RETURN p)
		RETURN {data:{
			id:doc.id, 
			content_id:content.id, 
			image:content.data.image, 
			title:content.data.title, 
			category:category[0].name, 
			category_id:content.data.category.id
		}}`, _from, common.DbBookmarkEdgeTable, categoryId, common.DbContentCategoryTable)

	resp, err := runQuery(ctx, q, common.DbBookmarkEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if b, err := recordToBookmarkContent(r); err == nil {
			bookmarkContents = append(bookmarkContents, b)
		}
	}
	return bookmarkContents, nil
}

func ReadBookmark(ctx context.Context, bookmarkId string) (*userapp_proto.Bookmark, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, bookmarkId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		FOR u IN %v FILTER u._id == doc._from
		FOR c IN %v FILTER c._id == doc._to
		LET category = (FOR p IN %v FILTER c.data.category.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc, {data:{
			user:u.data,
			content:c.data,
			contentCategory:category[0]
		}})`, common.DbBookmarkEdgeTable, query,
		common.DbUserTable, common.DbContentTable, common.DbContentCategoryTable,
	)

	resp, err := runQuery(ctx, q, common.DbBookmarkEdgeTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}
	// parsing
	b, err := recordToBookmark(resp.Records[0])
	if err != nil {
		return nil, err
	}
	return b, nil
}

func DeleteBookmark(ctx context.Context, orgid, userid, bookmarkId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, bookmarkId)
	query = common.QueryAuth(query, orgid, userid)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbBookmarkEdgeTable, query, common.DbBookmarkEdgeTable)
	_, err := runQuery(ctx, q, common.DbBookmarkEdgeTable)
	return err
}

func GetSharedContent(ctx context.Context, userId string) ([]*content_proto.SharedContent, error) {
	shares := []*content_proto.SharedContent{}

	q := fmt.Sprintf(`
		FOR content, doc IN INBOUND "%v/%v" %v
		FILTER doc.data.status == "%v"
		UPDATE doc WITH {updated:%v, data:{updated:%v, status:"%v"}} IN %v
		LET category = (FOR p IN %v FILTER content.data.category.id == p._key RETURN p.data)
		LET sharedBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			content_id:content.id,
			image:content.data.image,
			title:content.data.title,
			summary: content.data.summary,
			item:{"@type": content.data.item["@type"]},
			category:category[0].name,
			category_id:content.data.category.id,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url}}}`,
		common.DbUserTable, userId, common.DbShareContentUserEdgeTable,
		static_proto.ShareStatus_SHARED,
		time.Now().Unix(), time.Now().Unix(), static_proto.ShareStatus_RECEIVED, common.DbShareContentUserEdgeTable,
		common.DbContentCategoryTable,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareContentUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToSharedContent(r); err == nil {
			shares = append(shares, s)
		}
	}

	return shares, nil
}

func GetContentDetail(ctx context.Context, id, userId, orgId string) (*userapp_proto.ContentDetail, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	_to := fmt.Sprintf("%v/%v", common.DbContentTable, id)
	q := fmt.Sprintf(`
		FOR doc in %v
		%s
		LET category = (FOR p IN %v FILTER doc.data.category.id == p._key RETURN p.data)
		LET actions = (
			FOR p IN %v 
			FOR a IN category[0].actions
			FILTER a.id == p._key 
			RETURN p.data
		)
		LET bookmark  = (
			FOR u, edge IN ANY "%v" %v
			RETURN COUNT(edge)
		)
		LET rating = (
			FOR u, edge IN ANY "%v" %v
			FILTER u._key == "%v"
			RETURN edge.data.rating
		)
		LET type = (FOR p IN %v FILTER content.data.type.id == p._key RETURN p.data)
		LET source = (FOR p IN %v FILTER content.data.source.id == p._key RETURN p.data)
		LET sharedBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			image:doc.data.image,
			title:doc.data.title,
			summary: doc.data.summary,
			item:doc.data.item,
			category:category[0].name,
			category_id:doc.data.category.id,
			actions:actions,
			type:type[0],
			source:source[0],
			bookmarked:bookmark[0]>0,
			rating:rating[0],
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url}}}`,
		common.DbContentTable,
		query,
		common.DbContentCategoryTable,
		common.DbActionTable,
		_to, common.DbBookmarkEdgeTable,
		_to, common.DbContentRatingEdgeTable, userId,
		common.DbContentTypeTable,
		common.DbSourceTable,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareContentUserEdgeTable)
	if err != nil {
		return nil, err
	}

	data, err := recordToContentDetail(resp.Records[0])
	return data, err
}

func SearchBookmarks(ctx context.Context, title, summary, description string, userId, orgId, teamId string, offset, limit int64) ([]*user_proto.ShareableContent, error) {
	shares := []*user_proto.ShareableContent{}
	query := `FILTER`
	if len(title) > 0 {
		query += fmt.Sprintf(` && LIKE(doc.data.title, "%s",true)`, `%`+title+`%`)
	}
	if len(description) > 0 {
		query += fmt.Sprintf(` && LIKE(doc.data.description, "%v",true)`, `%`+description+`%`)
	}
	if len(summary) > 0 {
		query += fmt.Sprintf(` && LIKE(doc.data.summary, "%v",true)`, `%`+summary+`%`)
	}
	query = common.QueryAuth(query, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort("", "")

	q := fmt.Sprintf(`
		FOR doc in OUTBOUND "%v/%v" %v
		%s
		%s
		%s
		LET category = (FOR p IN %v FILTER content.data.category.id == p._key RETURN p.data)
		LET sharedBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			image:doc.data.image,
			title:doc.data.title,
			summary: doc.data.summary,
			item:{"@type": doc.data.item["@type"]},
			category:category[0].name,
			category_id:doc.data.category.id,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url}}}`,
		common.DbUserTable, userId, common.DbBookmarkEdgeTable,
		query,
		limit_query,
		sort_query,
		common.DbContentCategoryTable,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareContentUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToShareableContent(r); err == nil {
			shares = append(shares, s)
		}
	}

	return shares, nil
}

func GetSharedSurveysForUser(ctx context.Context, userId string) ([]*userapp_proto.SharedSurvey, error) {
	shares := []*userapp_proto.SharedSurvey{}
	q := fmt.Sprintf(`
		FOR survey, doc IN INBOUND "%v/%v" %v
		FILTER doc.data.status == "%v"
		UPDATE doc WITH {updated:%v, data:{updated:%v, status:"%v"}} IN %v
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			survey_id:survey.id,
			title:survey.data.title,
			count:doc.data.count,
			summary: survey.data.summary,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url}}}`,
		common.DbUserTable, userId, common.DbShareSurveyUserEdgeTable,
		static_proto.ShareStatus_SHARED,
		time.Now().Unix(), time.Now().Unix(), static_proto.ShareStatus_RECEIVED, common.DbShareSurveyUserEdgeTable,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareSurveyUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToSharedSurvey(r); err == nil {
			shares = append(shares, s)
		}
	}

	return shares, nil
}

func GetSharedPlansForUser(ctx context.Context, userId string) ([]*userapp_proto.SharedPlan, error) {
	shares := []*userapp_proto.SharedPlan{}
	q := fmt.Sprintf(`
		FOR plan, doc IN INBOUND "%v/%v" %v
		FILTER doc.data.status == "%v"
		UPDATE doc WITH {updated:%v, data:{updated:%v, status:"%v"}} IN %v
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			plan_id:plan.id,
			image:plan.data.pic,
			title:plan.data.name,
			duration:plan.data.duration,
			count:plan.data.items_count,
			summary: plan.data.summary,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url}}}`,
		common.DbUserTable, userId, common.DbSharePlanUserEdgeTable,
		static_proto.ShareStatus_SHARED,
		time.Now().Unix(), time.Now().Unix(), static_proto.ShareStatus_RECEIVED, common.DbSharePlanUserEdgeTable,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbSharePlanUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToSharedPlan(r); err == nil {
			shares = append(shares, s)
		}
	}

	return shares, nil
}

func GetSharedGoalsForUser(ctx context.Context, userId string) ([]*userapp_proto.SharedGoal, error) {
	shares := []*userapp_proto.SharedGoal{}
	q := fmt.Sprintf(`
		FOR goal, doc IN INBOUND "%v/%v" %v
		FILTER doc.data.status == "%v"
		UPDATE doc WITH {updated:%v, data:{updated:%v, status:"%v"}} IN %v
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			goal_id:goal.data.id,
			image:goal.data.image,
			title:goal.data.title,
			summary: goal.data.summary,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url},
			target:goal.data.target.targetValue,
			current:doc.data.currentValue,
			duration:goal.data.duration}}`,
		common.DbUserTable, userId, common.DbShareGoalUserEdgeTable,
		static_proto.ShareStatus_SHARED,
		time.Now().Unix(), time.Now().Unix(), static_proto.ShareStatus_RECEIVED, common.DbShareGoalUserEdgeTable,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareGoalUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToSharedGoal(r); err == nil {
			shares = append(shares, s)
		}
	}

	return shares, nil
}

func GetSharedChallengesForUser(ctx context.Context, userId string) ([]*userapp_proto.SharedChallenge, error) {
	shares := []*userapp_proto.SharedChallenge{}
	q := fmt.Sprintf(`
		FOR challenge, doc IN INBOUND "%v/%v" %v
		FILTER doc.data.status == "%v"
		UPDATE doc WITH {updated:%v, data:{updated:%v, status:"%v"}} IN %v
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			challenge_id:challenge.data.id,
			title:challenge.data.title,
			summary: challenge.data.summary,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url},
			target:challenge.data.target.targetValue,
			current:doc.data.currentValue,
			duration:challenge.data.duration}}`,
		common.DbUserTable, userId, common.DbShareChallengeUserEdgeTable,
		static_proto.ShareStatus_SHARED,
		time.Now().Unix(), time.Now().Unix(), static_proto.ShareStatus_RECEIVED, common.DbShareChallengeUserEdgeTable,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareChallengeUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToSharedChallenge(r); err == nil {
			shares = append(shares, s)
		}
	}

	return shares, nil
}

func GetSharedHabitsForUser(ctx context.Context, userId string) ([]*userapp_proto.SharedHabit, error) {
	shares := []*userapp_proto.SharedHabit{}
	q := fmt.Sprintf(`
		FOR habit, doc IN INBOUND "%v/%v" %v
		FILTER doc.data.status == "%v"
		UPDATE doc WITH {updated:%v, data:{updated:%v, status:"%v"}} IN %v
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			habit_id:habit.data.id,
			title:habit.data.title,
			summary: habit.data.summary,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url},
			target:habit.data.target.targetValue,
			current:doc.data.currentValue,
			duration:habit.data.duration}}`,
		common.DbUserTable, userId, common.DbShareHabitUserEdgeTable,
		static_proto.ShareStatus_SHARED,
		time.Now().Unix(), time.Now().Unix(), static_proto.ShareStatus_RECEIVED, common.DbShareHabitUserEdgeTable,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbShareHabitUserEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if s, err := recordToSharedHabit(r); err == nil {
			shares = append(shares, s)
		}
	}
	return shares, nil
}

func SignupToGoal(ctx context.Context, userId, goalId string, join_goal *userapp_proto.JoinGoal) (bool, error) {
	_from := fmt.Sprintf(`%v/%v`, common.DbUserTable, userId)
	_to := fmt.Sprintf(`%v/%v`, common.DbGoalTable, goalId)
	record, err := joingoalToRecord(_from, _to, userId, goalId, join_goal)
	if err != nil {
		return false, err
	}
	if len(record) == 0 {
		return false, errors.New("server serialization")
	}

	field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)
	q := fmt.Sprintf(`
		UPSERT %v
		INSERT %v
		UPDATE %v
		INTO %v
		RETURN {data:{type: OLD ? "update" : "insert"} }`, field, record, record, common.DbJoinGoalEdgeTable)

	resp, err := runQuery(ctx, q, common.DbJoinGoalEdgeTable)
	if err != nil {
		return false, err
	}

	// parsing to check whether this was an update (returns false) or insert (returns true)
	b, err := recordToCheckUpsertOperation(resp.Records[0])
	if err != nil {
		return false, err
	}
	return b, nil
}

func GetJoinedGoals(ctx context.Context, userId string, isCurrent bool) ([]*userapp_proto.JoinGoalResponse, error) {
	response := []*userapp_proto.JoinGoalResponse{}

	var query string
	if isCurrent {
		query = fmt.Sprintf(`FILTER doc.data.status == "%v" || doc.data.status == "%v"`, userapp_proto.ActionStatus_STARTED, userapp_proto.ActionStatus_IN_PROGRESS)
	}

	//TODO: The current value is not being set because it's not being fetched from the shared_doc
	// current:shared_doc.data.currentValue, (old query before change on 19/08/2018
	// solution is change the data model for join_goal_edge to save currentValue (which can be set by the user or from shared_doc)

	// shared doc query
	//		FOR shared_goal, shared_doc IN INBOUND "%v/%v" %v
	//		FILTER  goal._id == shared_goal._id

	q := fmt.Sprintf(`
		FOR goal, doc IN OUTBOUND "%v/%v" %v
		%v
		LET sharedBy = (FOR p IN %v FILTER p._key == goal.data.createdBy.id RETURN p.data)
		RETURN {data:{
					goal_id:goal.id,
					image:goal.data.image,
					title:goal.data.title,
					shared_by:CONCAT_SEPARATOR(" ", sharedBy[0].firstname, sharedBy[0].lastname),
					shared_by_image:sharedBy[0].avatar_url,
					current:goal.data.currentValue,
					target:goal.data.target.targetValue,
					duration:goal.data.duration}}`,
		common.DbUserTable, userId, common.DbJoinGoalEdgeTable,
		query,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbJoinGoalEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if j, err := recordToGoalResponse(r); err == nil {
			response = append(response, j)
		}
	}
	return response, nil
}

func SignupToChallenge(ctx context.Context, userId, challengeId string, join_challenge *userapp_proto.JoinChallenge) (bool, error) {
	_from := fmt.Sprintf(`%v/%v`, common.DbUserTable, userId)
	_to := fmt.Sprintf(`%v/%v`, common.DbChallengeTable, challengeId)
	record, err := joinchallengeToRecord(_from, _to, userId, challengeId, join_challenge)
	if err != nil {
		return false, err
	}
	if len(record) == 0 {
		return false, errors.New("server serialization")
	}

	field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)
	q := fmt.Sprintf(`
		UPSERT %v
		INSERT %v
		UPDATE %v
		INTO %v
		RETURN {data:{type: OLD ? "update" : "insert"} }`, field, record, record, common.DbJoinChallengeEdgeTable)

	resp, err := runQuery(ctx, q, common.DbJoinChallengeEdgeTable)
	if err != nil {
		return false, err
	}

	// parsing to check whether this was an update (returns false) or insert (returns true)
	b, err := recordToCheckUpsertOperation(resp.Records[0])
	if err != nil {
		return false, err
	}
	return b, nil
}

func GetJoinedChallenges(ctx context.Context, userId string, isCurrent bool) ([]*userapp_proto.JoinChallengeResponse, error) {
	response := []*userapp_proto.JoinChallengeResponse{}

	var query string
	if isCurrent {
		query = fmt.Sprintf(`FILTER doc.data.status == "%v" || doc.data.status == "%v"`, userapp_proto.ActionStatus_STARTED, userapp_proto.ActionStatus_IN_PROGRESS)
	}

	//TODO: The current value is not being set because it's not being fetched from the shared_doc
	// current:shared_doc.data.currentValue, (old query before change on 19/08/2018
	// solution is change the data model for join_challenge_edge to save currentValue (which can be set by the user or from shared_doc)

	// shared doc query
	//		FOR shared_challenge, shared_doc IN INBOUND "%v/%v" %v
	//		FILTER  challenge._id == shared_challenge._id

	q := fmt.Sprintf(`
		FOR challenge, doc IN OUTBOUND "%v/%v" %v
		%v
		LET sharedBy = (FOR p IN %v FILTER p._key == challenge.data.createdBy.id RETURN p.data)
		RETURN {data:{
					challenge_id:challenge.id,
					image:challenge.data.image,
					title:challenge.data.title,
					shared_by:CONCAT_SEPARATOR(" ", sharedBy[0].firstname, sharedBy[0].lastname),
					shared_by_image:sharedBy[0].avatar_url,
					current:challenge.data.currentValue,
					target:challenge.data.target.targetValue,
					duration:challenge.data.duration}}`,
		common.DbUserTable, userId, common.DbJoinChallengeEdgeTable,
		query,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbJoinChallengeEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if j, err := recordToChallengeResponse(r); err == nil {
			response = append(response, j)
		}
	}
	return response, nil
}

func SignupToHabit(ctx context.Context, userId, habitId string, join_habit *userapp_proto.JoinHabit) (bool, error) {
	_from := fmt.Sprintf(`%v/%v`, common.DbUserTable, userId)
	_to := fmt.Sprintf(`%v/%v`, common.DbHabitTable, habitId)
	record, err := joinhabitToRecord(_from, _to, userId, habitId, join_habit)
	if err != nil {
		return false, err
	}
	if len(record) == 0 {
		return false, errors.New("server serialization")
	}

	field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)
	q := fmt.Sprintf(`
		UPSERT %v
		INSERT %v
		UPDATE %v
		INTO %v
		RETURN {data:{type: OLD ? "update" : "insert"} }`, field, record, record, common.DbJoinHabitEdgeTable)

	resp, err := runQuery(ctx, q, common.DbJoinHabitEdgeTable)
	if err != nil {
		return false, err
	}

	// parsing to check whether this was an update (returns false) or insert (returns true)
	b, err := recordToCheckUpsertOperation(resp.Records[0])
	if err != nil {
		return false, err
	}
	return b, nil
}

func GetJoinedHabits(ctx context.Context, userId string, isCurrent bool) ([]*userapp_proto.JoinHabitResponse, error) {
	response := []*userapp_proto.JoinHabitResponse{}

	var query string
	if isCurrent {
		query = fmt.Sprintf(`FILTER doc.data.status == "%v" || doc.data.status == "%v"`, userapp_proto.ActionStatus_STARTED, userapp_proto.ActionStatus_IN_PROGRESS)
	}

	//TODO: The current value is not being set because it's not being fetched from the shared_doc
	// current:shared_doc.data.currentValue, (old query before change on 19/08/2018
	// solution is change the data model for join_habit_edge to save currentValue (which can be set by the user or from shared_doc)

	q := fmt.Sprintf(`
		FOR habit, doc IN OUTBOUND "%v/%v" %v
		%v
		LET sharedBy = (FOR p IN %v FILTER p._key == habit.data.createdBy.id RETURN p.data)
		RETURN {data:{
					habit_id:habit.id,
					image:habit.data.image,
					title:habit.data.title,
					shared_by:CONCAT_SEPARATOR(" ", sharedBy[0].firstname, sharedBy[0].lastname),
					shared_by_image:sharedBy[0].avatar_url,
					current:habit.data.currentValue,
					target:habit.data.target.targetValue,
					duration:habit.data.duration}}`,
		common.DbUserTable, userId, common.DbJoinHabitEdgeTable,
		query,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbJoinHabitEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if j, err := recordToHabitResponse(r); err == nil {
			response = append(response, j)
		}
	}
	return response, nil
}

func ListMarkers(ctx context.Context, userId string) ([]*userapp_proto.MarkerResponse, error) {
	response := []*userapp_proto.MarkerResponse{}

	query := fmt.Sprintf(`FILTER doc.data.status == "%v" || doc.data.status == "%v"`, userapp_proto.ActionStatus_STARTED, userapp_proto.ActionStatus_IN_PROGRESS)
	trackerMethodQuery := fmt.Sprintf(`LET trackerMethods = (
			FOR t IN OUTBOUND p %v
			OPTIONS {
				bfs: true,
				uniqueVertices: "global"
			}
			RETURN t.data
		)
		LET pret =(MERGE_RECURSIVE(p, {data:{
                trackerMethods: trackerMethods
        }}))`, common.DbMarkerTrackerEdgeTable)
	// for goals
	q := fmt.Sprintf(`
		FOR goal, doc IN OUTBOUND "%v/%v" %v
		%v
		FOR category IN %v FILTER goal.data.category.id == category._key
		LET default = (
			FOR p IN %v 
			FILTER category.data.markerDefault.id == p._key 
			%s
			RETURN pret.data
		)
		LET options = (
			FOR m IN category.data.markerOptions
			FOR p IN %v 
			FILTER m.id == p._key 
			%s
			RETURN pret.data
		)
		RETURN MERGE_RECURSIVE(category, {data:{
			markerDefault:default[0],
			markerOptions:options
			}})`,
		common.DbUserTable, userId, common.DbJoinGoalEdgeTable,
		query,
		common.DbBehaviourCategoryTable,
		common.DbMarkerTable, trackerMethodQuery, common.DbMarkerTable, trackerMethodQuery)
	resp, err := runQuery(ctx, q, common.DbJoinGoalEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if category, err := recordToCategory(r); err == nil {
			// marker default
			if category.MarkerDefault != nil {
				marker := &userapp_proto.MarkerResponse{
					IsDefault:      true,
					MarkerId:       category.MarkerDefault.Id,
					Name:           category.MarkerDefault.Name,
					IconSlug:       category.MarkerDefault.IconSlug,
					Unit:           category.MarkerDefault.Unit,
					TrackerMethods: category.MarkerDefault.TrackerMethods,
				}
				response = append(response, marker)
			}
			// marker opetions
			if category.MarkerOptions != nil {
				for _, m := range category.MarkerOptions {
					//log.Println(m)
					marker := &userapp_proto.MarkerResponse{
						IsDefault:      false,
						MarkerId:       m.Id,
						Name:           m.Name,
						IconSlug:       m.IconSlug,
						Unit:           m.Unit,
						TrackerMethods: m.TrackerMethods,
					}
					response = append(response, marker)
				}
			}
		} else {
			common.ErrorLog(common.UserappSrv, ListMarkers, err, "Category unmarshale is failed")
		}
	}

	// for challenge
	q = fmt.Sprintf(`
		FOR challenge, doc IN OUTBOUND "%v/%v" %v
		%v
		FOR category IN %v FILTER challenge.data.category.id == category._key
		LET default = (
			FOR p IN %v 
			FILTER category.data.markerDefault.id == p._key 
			%s
			RETURN pret.data
		)
		LET options = (
			FOR m IN category.data.markerOptions
			FOR p IN %v 
			FILTER m.id == p._key 
			%s
			RETURN pret.data
		)
		RETURN MERGE_RECURSIVE(category, {data:{
			markerDefault:default[0],
			markerOptions:options
			}})`,
		common.DbUserTable, userId, common.DbJoinChallengeEdgeTable,
		query,
		common.DbBehaviourCategoryTable,
		common.DbMarkerTable, trackerMethodQuery, common.DbMarkerTable, trackerMethodQuery)
	resp, err = runQuery(ctx, q, common.DbJoinChallengeEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if category, err := recordToCategory(r); err == nil {
			// marker default
			if category.MarkerDefault != nil {
				marker := &userapp_proto.MarkerResponse{
					IsDefault:      true,
					MarkerId:       category.MarkerDefault.Id,
					Name:           category.MarkerDefault.Name,
					IconSlug:       category.MarkerDefault.IconSlug,
					Unit:           category.MarkerDefault.Unit,
					TrackerMethods: category.MarkerDefault.TrackerMethods,
				}
				response = append(response, marker)
			}
			// marker opetions
			if category.MarkerOptions != nil {
				for _, m := range category.MarkerOptions {
					marker := &userapp_proto.MarkerResponse{
						IsDefault:      false,
						MarkerId:       m.Id,
						Name:           m.Name,
						IconSlug:       m.IconSlug,
						Unit:           m.Unit,
						TrackerMethods: m.TrackerMethods,
					}
					response = append(response, marker)
				}
			}
		}
	}

	return dedupResponse(response), nil
}

//remove duplicate markers from the response
func dedupResponse(response []*userapp_proto.MarkerResponse) []*userapp_proto.MarkerResponse {
	u := make([]*userapp_proto.MarkerResponse, 0, len(response))
	m := make(map[string]*userapp_proto.MarkerResponse)

	for _, marker := range response {
		if _, ok := m[marker.MarkerId]; !ok {
			m[marker.MarkerId] = marker
			u = append(u, marker)
		}
	}

	return u
}

func GetPendingSharedActions(ctx context.Context, userId, orgId string, offset, limit, from, to int64, sortParameter, sortDirection string) ([]*common_proto.PendingResponse, error) {
	// query := fmt.Sprintf(`FILTER doc.target_user == "%v"`, userId)
	query := fmt.Sprintf(`FILTER doc.parameter2 == "%v"`, userId)
	query = common.QueryAuth(query, orgId, "")

	if from > 0 {
		query += fmt.Sprintf(` && %v <= doc.created`, from)
	}
	if to > 0 {
		query += fmt.Sprintf(` && doc.created <= %v`, to)
	}

	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)
	response := []*common_proto.PendingResponse{}

	q := fmt.Sprintf(
		`FOR doc IN %v
		%s
		%s
		%s
		LET sharedBy = (FOR p IN %v FILTER doc.data.shared_by.id == p._key RETURN p.data)
		LET sharedWith = (FOR p IN %v FILTER doc.data.shared_with.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.data.id,
			org_id: doc.data.org_id,
			created: doc.created,
			shared_by: {"id": sharedBy[0].id, "firstname": sharedBy[0].firstname, "lastname": sharedBy[0].lastname, "avatar_url": sharedBy[0].avatar_url},
			item: doc.data.item}}`,
		common.DbPendingTable,
		query,
		limit_query,
		sort_query,
		common.DbUserTable, common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbPendingTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if p, err := recordToPending(r); err == nil {
			response = append(response, p)
		} else {
			log.Println(p, err)
		}
	}
	return response, nil
}

func GetGoalProgress(ctx context.Context, userId string) ([]*userapp_proto.GoalProgressResponse, error) {
	query := fmt.Sprintf(`FILTER doc.data.status == "%v" || doc.data.status == "%v"`, userapp_proto.ActionStatus_STARTED, userapp_proto.ActionStatus_IN_PROGRESS)
	response := []*userapp_proto.GoalProgressResponse{}

	q := fmt.Sprintf(`
		FOR goal, doc IN OUTBOUND "%v/%v" %v
		%s
		LET category = (
				FOR cat in %v
				filter cat._key == goal.data.category.id
				return cat.data )
		LET marker = (
				FOR m in %v
				filter m._key == category[0].markerDefault.id
				return m.data )
		LET user = (
				FOR u in %v
				filter u._key == "%v"
				return u )
		LET value = (
			FOR tm IN %v
			FILTER tm.data.marker.id == category[0].markerDefault.id && tm.data.user.id == "%v"
			SORT tm.created DESC
			LIMIT 0, 1
			RETURN tm.data.value )
		RETURN {data:{
				goal:{"id":goal.data.id, "title":goal.data.title},
				user: {"id":user[0].data.id},
				latestValue:value[0],
				target:goal.data.target.targetValue,
				unit: goal.data.target.unit}}`,
		common.DbUserTable, userId, common.DbJoinGoalEdgeTable,
		query,
		common.DbBehaviourCategoryTable,
		common.DbMarkerTable,
		common.DbUserTable, userId,
		common.DbTrackMarkerTable, userId)

	resp, err := runQuery(ctx, q, common.DbJoinGoalEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if p, err := recordToGoalProgressResponse(r); err == nil {
			response = append(response, p)
		} else {
			log.Println(p, err)
		}
	}
	return response, nil
}

func GetDefaultMarkerHistory(ctx context.Context, userId string, offset, limit, from, to int64, sortParameter, sortDirection string) ([]*track_proto.TrackMarker, error) {
	query := fmt.Sprintf(`FILTER doc.parameter1 == "%v"`, userId)
	query += ` && (doc.data.status == "STARTED" || doc.data.status == "IN_PROGRESS")`

	if from > 0 {
		query += fmt.Sprintf(` && %v <= doc.created`, from)
	}
	if to > 0 {
		query += fmt.Sprintf(` && doc.created <= %v`, to)
	}

	response := []*track_proto.TrackMarker{}
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection, "marker")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%v
		FOR share IN %v
		FILTER share._to == "%v/%v" && share._from == doc._to
		FOR marker IN %v
		FILTER marker.data.marker.id == share.data.goal.category.markerDefault.id
		%v
		%v
		RETURN marker`,
		common.DbJoinGoalEdgeTable, query, common.DbShareGoalUserEdgeTable, common.DbUserTable, userId, common.DbTrackMarkerTable, limit_query, sort_query)

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

func GetCurrentChallengesWithCount(ctx context.Context, userId string) ([]*userapp_proto.ChallengeCountResponse, error) {
	query := fmt.Sprintf(`FILTER doc.data.status == "%v" || doc.data.status == "%v"`, userapp_proto.ActionStatus_STARTED, userapp_proto.ActionStatus_IN_PROGRESS)

	response := []*userapp_proto.ChallengeCountResponse{}
	// 19/06/2018
	//TODO: This needs to fixed by adding currentValue to the join_challenge_user table
	// currently getting dummy value from challenge (which is only a placeholder value)
	// current = challenge.data.currentValue,
	q := fmt.Sprintf(
		`FOR challenge, doc IN OUTBOUND "%v/%v" %v
		%s
		LET date = ( FOR d IN %v
					COLLECT AGGREGATE
							minDate = MIN(d.created),
							maxDate = MAX(d.created)
					RETURN { minDate, maxDate }
		)

			LET count = (
				FOR th IN %v
				FILTER th.data.challenge.id == challenge._key 
						&& th.created >= date[0].minDate 
						&& th.created <= date[0].maxDate
						&& th.data.user.id == "%v"
				COLLECT WITH COUNT INTO length
				RETURN length
			)
		RETURN {data:{
						challenge_id:challenge._key,
						title:challenge.data.title,
						image:challenge.data.image,
						current:0,
						count:count[0],
						target:challenge.data.target.targetValue,
						duration:challenge.data.duration}}
		`,
		common.DbUserTable, userId, common.DbJoinChallengeEdgeTable,
		query,
		common.DbTrackChallengeTable,
		common.DbTrackChallengeTable, userId, userId)

	resp, err := runQuery(ctx, q, common.DbJoinChallengeEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if p, err := recordToChallengeCountResponse(r); err == nil {
			response = append(response, p)
		}
	}
	return response, nil
}

func GetCurrentHabitsWithCount(ctx context.Context, userId string) ([]*userapp_proto.HabitCountResponse, error) {
	query := fmt.Sprintf(`FILTER doc.data.status == "%v" || doc.data.status == "%v"`, userapp_proto.ActionStatus_STARTED, userapp_proto.ActionStatus_IN_PROGRESS)

	response := []*userapp_proto.HabitCountResponse{}
	// 19/06/2018
	//TODO: This needs to fixed by adding currentValue to the join_habit_user table
	// currently getting dummy value from habit (which is only a placeholder value)
	// current = habit.data.currentValue,
	q := fmt.Sprintf(
		`FOR habit, doc IN OUTBOUND "%v/%v" %v
		%s
		LET date = ( FOR d IN %v
					COLLECT AGGREGATE
							minDate = MIN(d.created),
							maxDate = MAX(d.created)
					RETURN { minDate, maxDate }
		)

			LET count = (
				FOR th IN %v
				FILTER th.data.habit.id == habit._key 
						&& th.created >= date[0].minDate 
						&& th.created <= date[0].maxDate
						&& th.data.user.id == "%v"
				COLLECT WITH COUNT INTO length
				RETURN length
			)
			RETURN {data:{
						habit_id:habit._key,
						title:habit.data.title,
						image:habit.data.image,
						current:0,
						count:count[0],
						target:habit.data.target.targetValue,
						duration:habit.data.duration}}	
		`,
		common.DbUserTable, userId, common.DbJoinHabitEdgeTable,
		query,
		common.DbTrackHabitTable,
		common.DbTrackHabitTable,
		userId, userId)

	resp, err := runQuery(ctx, q, common.DbJoinHabitEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if p, err := recordToHabitCountResponse(r); err == nil {
			response = append(response, p)
		}
	}
	return response, nil
}

func RemovePendingSharedAction(ctx context.Context, itemId, userId string) error {
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc.data.item.id == "%v" && doc.data.shared_with.id == "%v"
		REMOVE doc IN %v`, common.DbPendingTable, itemId, userId, common.DbPendingTable)

	_, err := runQuery(ctx, q, common.DbPendingTable)
	return err
}

func UpdateShareGoalStatus(ctx context.Context, goalId string, status static_proto.ShareStatus) error {
	_from := fmt.Sprintf(`%v/%v`, common.DbGoalTable, goalId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._from == "%v"
		UPDATE doc WITH {data:{status:"%v"}} IN %v`, common.DbShareGoalUserEdgeTable, _from, status, common.DbShareGoalUserEdgeTable)
	_, err := runQuery(ctx, q, common.DbShareGoalUserEdgeTable)
	return err
}

func UpdateShareChallengeStatus(ctx context.Context, challengeId string, status static_proto.ShareStatus) error {
	_from := fmt.Sprintf(`%v/%v`, common.DbChallengeTable, challengeId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._from == "%v"
		UPDATE doc WITH {data:{status:"%v"}} IN %v`, common.DbShareChallengeUserEdgeTable, _from, status, common.DbShareChallengeUserEdgeTable)
	_, err := runQuery(ctx, q, common.DbShareChallengeUserEdgeTable)
	return err
}

func UpdateShareHabitStatus(ctx context.Context, habitId string, status static_proto.ShareStatus) error {
	_from := fmt.Sprintf(`%v/%v`, common.DbHabitTable, habitId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._from == "%v"
		UPDATE doc WITH {data:{status:"%v"}} IN %v`, common.DbShareHabitUserEdgeTable, _from, status, common.DbShareHabitUserEdgeTable)
	_, err := runQuery(ctx, q, common.DbShareHabitUserEdgeTable)
	return err
}

func UpdateSharePlanStatus(ctx context.Context, planId string, status static_proto.ShareStatus) error {
	_from := fmt.Sprintf(`%v/%v`, common.DbPlanTable, planId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._from == "%v"
		UPDATE doc WITH {data:{status:"%v"}} IN %v`, common.DbSharePlanUserEdgeTable, _from, status, common.DbSharePlanUserEdgeTable)
	_, err := runQuery(ctx, q, common.DbSharePlanUserEdgeTable)
	return err
}

func SaveRateForContent(ctx context.Context, contentRating *userapp_proto.ContentRating) error {
	if contentRating.Created == 0 {
		contentRating.Created = time.Now().Unix()
	}
	contentRating.Updated = time.Now().Unix()

	_from := fmt.Sprintf(`%v/%v`, common.DbUserTable, contentRating.UserId)
	_to := fmt.Sprintf(`%v/%v`, common.DbContentTable, contentRating.ContentId)
	record, err := contentRatingToRecord(_from, _to, contentRating)
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
		INTO %v`, field, record, record, common.DbContentRatingEdgeTable)

	if _, err := runQuery(ctx, q, common.DbContentRatingEdgeTable); err != nil {
		return err
	}
	return nil
}

func DislikeForContent(ctx context.Context, contentDislike *userapp_proto.ContentDislike) error {
	if contentDislike.Created == 0 {
		contentDislike.Created = time.Now().Unix()
	}
	contentDislike.Updated = time.Now().Unix()

	_from := fmt.Sprintf(`%v/%v`, common.DbUserTable, contentDislike.UserId)
	_to := fmt.Sprintf(`%v/%v`, common.DbContentTable, contentDislike.ContentId)
	record, err := contentDislikeToRecord(_from, _to, contentDislike)
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
		INTO %v`, field, record, record, common.DbContentDislikeEdgeTable)

	if _, err := runQuery(ctx, q, common.DbContentDislikeEdgeTable); err != nil {
		return err
	}
	return nil
}

func DislikeForSimilarContent(ctx context.Context, contentDislikeSimilar *userapp_proto.ContentDislikeSimilar) error {
	if contentDislikeSimilar.Created == 0 {
		contentDislikeSimilar.Created = time.Now().Unix()
	}
	contentDislikeSimilar.Updated = time.Now().Unix()

	_from := fmt.Sprintf(`%v/%v`, common.DbUserTable, contentDislikeSimilar.UserId)
	_to := fmt.Sprintf(`%v/%v`, common.DbContentTable, contentDislikeSimilar.ContentId)
	record, err := contentDislikeSimilarToRecord(_from, _to, contentDislikeSimilar)
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
		INTO %v`, field, record, record, common.DbContentDislikeSimilarEdgeTable)

	if _, err := runQuery(ctx, q, common.DbContentDislikeSimilarEdgeTable); err != nil {
		return err
	}
	return nil
}

func SaveUserFeedback(ctx context.Context, feedback *user_proto.UserFeedback) error {
	if len(feedback.Id) == 0 {
		feedback.Id = uuid.NewUUID().String()
	}
	if feedback.Created == 0 {
		feedback.Created = time.Now().Unix()
	}
	feedback.Updated = time.Now().Unix()
	record, err := feedbackToRecord(feedback)
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
		INTO %v`, feedback.Id, record, record, common.DbUserFeedbackTable)

	if _, err := runQuery(ctx, q, common.DbUserFeedbackTable); err != nil {
		return err
	}
	return nil
}

func CreateUserPlan(ctx context.Context, userId, planId string, userplan *userapp_proto.UserPlan, edge bool) error {
	if len(userplan.Id) == 0 {
		userplan.Id = uuid.NewUUID().String()
	}
	if userplan.Created == 0 {
		userplan.Created = time.Now().Unix()
	}
	userplan.Updated = time.Now().Unix()

	if edge {
		_from := fmt.Sprintf(`%v/%v`, common.DbUserTable, userId)
		_to := fmt.Sprintf(`%v/%v`, common.DbPlanTable, planId)
		field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)

		q := fmt.Sprintf(`
			UPSERT %v
			INSERT %v
			UPDATE %v
			INTO %v`, field, field, field, common.DbUserPlanEdgeTable)
		if _, err := runQuery(ctx, q, common.DbUserPlanEdgeTable); err != nil {
			return err
		}
	}

	record, err := userplanToRecord(userplan)
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
		INTO %v`, userplan.Id, record, record, common.DbUserPlanTable)

	if _, err := runQuery(ctx, q, common.DbUserPlanTable); err != nil {
		return err
	}
	return nil
}

func GetUserPlan(ctx context.Context, userId string) (*userapp_proto.UserPlan, error) {
	query := fmt.Sprintf(`FILTER doc.data.targetUser == "%v"`, userId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbUserPlanTable, query)
	resp, err := runQuery(ctx, q, common.DbUserPlanTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}
	// parsing
	userplan, err := recordToUserPlan(resp.Records[0])
	if err != nil {
		return nil, err
	}
	return userplan, nil
}

func GetUserPlanWithPlanId(ctx context.Context, planId string) (*userapp_proto.UserPlan, error) {
	query := fmt.Sprintf(`FILTER doc.data.plan.id == "%v"`, planId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbUserPlanTable, query)

	resp, err := runQuery(ctx, q, common.DbUserPlanTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}
	// parsing
	userplan, err := recordToUserPlan(resp.Records[0])
	if err != nil {
		return nil, err
	}
	return userplan, nil
}

func UpdateUserPlan(ctx context.Context, id, orgId string, goals []*behaviour_proto.Goal, days map[string]*common_proto.DayItems) error {
	goals_body, err := json.Marshal(goals)
	if err != nil {
		return err
	}
	days_body, err := json.Marshal(days)
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc.id == "%v"
		UPDATE doc WITH {updated:%v, data:{org_id:"%v", goals:%v, days:%v, updated:%v}} IN %v
		RETURN NEW`, common.DbUserPlanTable, id, time.Now().Unix(), orgId, string(goals_body), string(days_body), time.Now().Unix(), common.DbUserPlanTable)
	_, err = runQuery(ctx, q, common.DbUserPlanTable)
	return err
}

func GetPlanItemsCountByCategory(ctx context.Context, planId, maxDay string) ([]*userapp_proto.CategoryCountResponse, error) {
	response := []*userapp_proto.CategoryCountResponse{}

	q := fmt.Sprintf(`
		LET ret = (
		FOR doc IN %v
		FILTER doc.data.plan.id == "%v"
		FOR i IN 1..%v
		RETURN doc.data.days[i].items
		)[**]
		FOR r IN ret
		FILTER r != null
		COLLECT category_id = r.categoryId INTO data
		RETURN {
			data: {
				"category_id":category_id,
				"category_icon_slug":data[0].r.categoryIconSlug,
				"category_name":data[0].r.categoryName,
				"item_count": COUNT(data)
			}
		}`, common.DbUserPlanTable, planId, maxDay)
	resp, err := runQuery(ctx, q, common.DbUserPlanTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if p, err := recordToCategoryCountResponse(r); err == nil {
			response = append(response, p)
		}
	}
	return response, nil
}

func GetPlanItemsCountByDay(ctx context.Context, planId, dayNumber string) ([]*userapp_proto.CategoryCountResponse, error) {
	response := []*userapp_proto.CategoryCountResponse{}

	q := fmt.Sprintf(`
		LET ret = (
		FOR doc IN %v
		FILTER doc.data.plan.id == "%v"
		RETURN doc.data.days["%v"].items
		)[**]
		FOR r IN ret
		FILTER r != null
		COLLECT category_id = r.categoryId INTO data
		RETURN {
			data: {
				"day_number":"%v",
				"category_id":category_id,
				"category_icon_slug":data[0].r.categoryIconSlug,
				"category_name":data[0].r.categoryName,
				"item_count": COUNT(data)
			}
		}`, common.DbUserPlanTable, planId, dayNumber, dayNumber)
	resp, err := runQuery(ctx, q, common.DbUserPlanTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if p, err := recordToCategoryCountResponse(r); err == nil {
			response = append(response, p)
		}
	}
	return response, nil
}

func GetPlanItemsCountByCategoryAndDay(ctx context.Context, planId, maxDay string) ([]*userapp_proto.CategoryCountResponse, error) {
	response := []*userapp_proto.CategoryCountResponse{}

	q := fmt.Sprintf(`
		LET ret = (
		FOR doc IN %v
		FILTER doc.data.plan.id == "%v"
		RETURN doc.data.days
		)
		LET days = (
		FOR r IN ret
		FOR i IN 1..%v
		FILTER r[i].items != null
		FOR item IN r[i].items
		LET day = MERGE_RECURSIVE(item, {"day_number":i})
		RETURN day
		)
		FOR day IN days
		COLLECT category_id = day.categoryId INTO data
		RETURN {
			data: {
				"day_number":TO_STRING(data[0].day.day_number),
				"category_id":category_id,
				"category_icon_slug":data[0].day.categoryIconSlug,
				"category_name":data[0].day.categoryName,
				"item_count": COUNT(data)
			}
		}`, common.DbUserPlanTable, planId, maxDay)
	resp, err := runQuery(ctx, q, common.DbUserPlanTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if p, err := recordToCategoryCountResponse(r); err == nil {
			response = append(response, p)
		} else {
			log.Println(err)
		}
	}
	return response, nil
}

func ReceivedItems(ctx context.Context, userId string, shared []*userapp_proto.SharedItem) error {
	for _, share := range shared {
		var collection string
		switch share.Type {
		case common.GOAL_TYPE:
			collection = common.DbShareGoalUserEdgeTable
		case common.CHALLENGE_TYPE:
			collection = common.DbShareChallengeUserEdgeTable
		case common.HABIT_TYPE:
			collection = common.DbShareHabitUserEdgeTable
		case common.PLAN_TYPE:
			collection = common.DbSharePlanUserEdgeTable
		case common.SURVEY_TYPE:
			collection = common.DbShareSurveyUserEdgeTable
		case common.CONTENT_TYPE:
			collection = common.DbShareContentUserEdgeTable
		}

		q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc.id == "%v" && doc._to == "%v/%v"
		UPDATE doc WITH {updated:%v, data:{updated:%v, status:"%v"}} IN %v`,
			collection, share.Id, common.DbUserTable, userId, time.Now().Unix(), time.Now().Unix(), static_proto.ShareStatus_RECEIVED, collection)
		_, err := runQuery(ctx, q, collection)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetGoalDetail reads a goal by ID
func GetGoalDetail(ctx context.Context, id, orgId string) (*userapp_proto.GoalDetail, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET category = (FOR p IN %v FILTER doc.data.category.id == p._key RETURN p.data)
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
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
		LET todo = (FOR p IN %v FILTER doc.data.todos.id == p._key RETURN p.data)

		RETURN {data:{
				goal_id:doc.data.id,
				title:doc.data.title,
				summary: doc.data.summary,
				description: doc.data.description,
				image:doc.data.image,
				shared_by: {"id": createdBy[0].id, "firstname": createdBy[0].firstname, "lastname": createdBy[0].lastname, "avatar_url": createdBy[0].avatar_url},
				target:doc.data.target,
				current:doc.data.currentValue,
				duration:doc.data.duration,
				source: doc.data.source,
				tags: doc.data.tags,
				challenges: challenges,
				habits: habits,
				category:category[0],
				todos:todo[0]
		}}`, common.DbGoalTable, query,
		common.DbBehaviourCategoryTable, common.DbUserTable, common.DbChallengeTable, common.DbHabitTable, common.DbTodoTable)

	resp, err := runQuery(ctx, q, common.DbGoalTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToGoalDetail(resp.Records[0])
	return data, err
}

// GetChallengeDetail reads a challenge by ID
func GetChallengeDetail(ctx context.Context, id, orgId string) (*userapp_proto.ChallengeDetail, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET category = (FOR p IN %v FILTER doc.data.category.id == p._key RETURN p.data)
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		LET habits = (
			FILTER NOT_NULL(doc.data.habits)
			FOR h IN doc.data.habits
			FOR p IN %v
			FILTER h.id == p._key RETURN p.data
		)
		LET todo = (FOR p IN %v FILTER doc.data.todos.id == p._key RETURN p.data)

		RETURN {data:{
				challenge_id:doc.data.id,
				title:doc.data.title,
				summary: doc.data.summary,
				description: doc.data.description,
				image:doc.data.image,
				shared_by: {"id": createdBy[0].id, "firstname": createdBy[0].firstname, "lastname": createdBy[0].lastname, "avatar_url": createdBy[0].avatar_url},
				target:doc.data.target,
				current:doc.data.currentValue,
				duration:doc.data.duration,
				source: doc.data.source,
				tags: doc.data.tags,
				habits: habits,
				category:category[0],
				todos:todo[0]
		}}`, common.DbChallengeTable, query, common.DbBehaviourCategoryTable, common.DbUserTable, common.DbHabitTable, common.DbTodoTable)

	resp, err := runQuery(ctx, q, common.DbChallengeTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToChallengeDetail(resp.Records[0])
	return data, err
}

// GetHabitDetail reads a habit by ID
func GetHabitDetail(ctx context.Context, id, orgId string) (*userapp_proto.HabitDetail, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET category = (FOR p IN %v FILTER doc.data.category.id == p._key RETURN p.data)
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)

		LET todo = (FOR p IN %v FILTER doc.data.todos.id == p._key RETURN p.data)

		RETURN {data:{
				habit_id:doc.data.id,
				title:doc.data.title,
				summary: doc.data.summary,
				description: doc.data.description,
				image:doc.data.image,
				shared_by: {"id": createdBy[0].id, "firstname": createdBy[0].firstname, "lastname": createdBy[0].lastname, "avatar_url": createdBy[0].avatar_url},
				target:doc.data.target,
				current:doc.data.currentValue,
				duration:doc.data.duration,
				source: doc.data.source,
				tags: doc.data.tags,
				category:category[0],
				todos:todo[0]
		}}`, common.DbHabitTable, query, common.DbBehaviourCategoryTable, common.DbUserTable, common.DbTodoTable)

	resp, err := runQuery(ctx, q, common.DbHabitTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToHabitDetail(resp.Records[0])
	return data, err
}
