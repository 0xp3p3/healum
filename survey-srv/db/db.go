package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"server/common"
	db_proto "server/db-srv/proto/db"
	common_proto "server/static-srv/proto/common"
	static_proto "server/static-srv/proto/static"
	survey_proto "server/survey-srv/proto/survey"
	user_proto "server/user-srv/proto/user"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/client"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
	hashids "github.com/speps/go-hashids"
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

func surveyToRecord(survey *survey_proto.Survey) (string, error) {
	data, err := common.MarhalToObject(survey)
	if err != nil {
		return "", err
	}

	var creatorId string
	if survey.Creator != nil {
		creatorId = survey.Creator.Id
	}
	common.FilterObject(data, "creator", survey.Creator)
	// filter shares
	if len(survey.Shares) > 0 {
		var arr []interface{}
		for _, item := range survey.Shares {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["shares"] = arr
	} else {
		delete(data, "shares")
	}
	// filter questions
	delete(data, "questions")

	d := map[string]interface{}{
		"_key":       survey.Id,
		"id":         survey.Id,
		"created":    survey.Created,
		"updated":    survey.Updated,
		"name":       survey.Title,
		"parameter1": survey.OrgId,
		"parameter2": creatorId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToSurvey(r *db_proto.Record) (*survey_proto.Survey, error) {
	var p survey_proto.Survey
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func questionToRecord(q *survey_proto.Question) (string, error) {
	data, err := common.MarhalToObject(q)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key": q.Id,
		"id":   q.Id,
		"name": q.Title,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToQuestion(r *db_proto.Record) (*survey_proto.Question, error) {
	var p survey_proto.Question
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func sharedToRecord(from, to, orgId string, shared *survey_proto.ShareSurveyUser) (string, error) {
	data, err := common.MarhalToObject(shared)
	if err != nil {
		return "", err
	}
	common.FilterObject(data, "survey", shared.Survey)
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

func recordToShared(r *db_proto.Record) (*survey_proto.ShareSurveyUser, error) {
	var p survey_proto.ShareSurveyUser
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func mapToRecord(m map[string]string) (string, error) {
	d := map[string]interface{}{
		"name":       m["unique_hash"],
		"parameter1": m["survey_id"],
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func surveyQuery(orgId, teamId string, offset, limit int64, sortParameter, sortDirection, filter_query string) string {
	org_query := "FILTER"
	org_query = common.QueryAuth(org_query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s
		LET creator = (FOR p IN %v FILTER p._key == doc.data.creator.id RETURN p.data)
		LET shares = (FILTER NOT_NULL(doc.data.shares) FOR u IN doc.data.shares FOR p IN %v FILTER u.id == p._key RETURN p.data)
		LET questions = (
			FOR q IN OUTBOUND doc._id %v
			RETURN q.data
		)
		RETURN MERGE_RECURSIVE(doc, {data:{
			creator:creator[0],
			shares:shares,
			questions:questions
		}})`,
		common.DbSurveyTable, org_query, filter_query, sort_query, limit_query,
		common.DbUserTable, common.DbUserTable,
		common.DbSurveyQuestionEdgeTable,
	)

	return q
}

func recordToShareableSurvey(r *db_proto.Record) (*user_proto.ShareableSurvey, error) {
	var p user_proto.ShareableSurvey
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func queryFilterSharedForUser(userId, shared_collection, org_query string) (string, string) {
	shared_query := ""
	filter_shared_query := ""
	if len(userId) > 0 {
		shared_query = fmt.Sprintf(`LET shared = (
			FOR e, doc IN INBOUND "%v/%v" %v
			%s
			RETURN doc._from
		)`, common.DbUserTable, userId, shared_collection, org_query)
		filter_shared_query = "FILTER doc._id NOT IN shared"
	}
	return shared_query, filter_shared_query
}

// All get all surveys
func All(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*survey_proto.Survey, error) {
	var surveys []*survey_proto.Survey
	q := surveyQuery(orgId, "", offset, limit, sortParameter, sortDirection, "")

	resp, err := runQuery(ctx, q, common.DbSurveyTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if survey, err := recordToSurvey(r); err == nil {
			surveys = append(surveys, survey)
		} else {
			common.ErrorLog(common.SurveySrv, All, err, "Survey unmarshalle is failed")
		}
	}
	return surveys, nil
}

// New create unique id
func New(ctx context.Context) (map[string]string, error) {
	survey_id := uuid.NewUUID().String()
	hd := hashids.NewData()
	hd.Salt = survey_id
	hd.MinLength = 6
	h, _ := hashids.NewWithData(hd)
	now := time.Now()
	unique_hash, err := h.Encode([]int{int(now.Unix())})
	if err != nil {
		return nil, err
	}

	m := map[string]string{}
	m["unique_hash"] = unique_hash
	m["survey_id"] = survey_id
	record, err := mapToRecord(m)
	if err != nil {
		return nil, err
	}

	q := fmt.Sprintf(`INSERT %v INTO %v`, record, common.DbSurveyHashTable)
	_, err = runQuery(ctx, q, common.DbSurveyHashTable)

	return m, err
}

// Creates a survey
func Create(ctx context.Context, survey *survey_proto.Survey) error {
	if len(survey.Id) == 0 {
		survey.Id = uuid.NewUUID().String()
	}
	if survey.Created == 0 {
		survey.Created = time.Now().Unix()
	}
	if survey.Updated == 0 {
		survey.Updated = time.Now().Unix()
	}
	if survey.Setting == nil {
		survey.Setting = &static_proto.Setting{
			AuthenticationRequired: true,
			ShowCaptcha:            true,
		}
	}

	record, err := surveyToRecord(survey)
	if err != nil {
		log.Println("surveyToRecord err:", err)
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, survey.Id, record, record, common.DbSurveyTable)
	_, err = runQuery(ctx, q, common.DbSurveyTable)
	if err != nil {
		log.Println("runQuery err:", err)
	}

	// after create question and make edge with survey.id
	for _, question := range survey.Questions {
		if err := CreateQuestion(ctx, survey.Id, question); err != nil {
			return err
		}
	}
	return err
}

// Reads a survey by ID
func Read(ctx context.Context, id, orgId, teamId string) (*survey_proto.Survey, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	q := surveyQuery(orgId, "", 0, 0, "", "", query)

	resp, err := runQuery(ctx, q, common.DbSurveyTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToSurvey(resp.Records[0])
	return data, err
}

// Deletes a survey by ID
func Delete(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbSurveyTable, query, common.DbSurveyTable)
	_, err := runQuery(ctx, q, common.DbSurveyTable)
	return err
}

// Questions ...
func Questions(ctx context.Context, surveyId string, offset, limit int64, sortParameter, sortDirection string) ([]*survey_proto.Question, error) {
	var questions []*survey_proto.Question
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN OUTBOUND "%v/%v" %v
		%v
		RETURN doc`,
		common.DbSurveyTable, surveyId, common.DbSurveyQuestionEdgeTable,
		sort_query,
	)

	resp, err := runQuery(ctx, q, common.DbSurveyQuestionTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	// parsing
	for _, r := range resp.Records {
		if q, err := recordToQuestion(r); err == nil {
			questions = append(questions, q)
		}
	}
	return questions, nil
}

// QuestionRef ...
func QuestionRef(ctx context.Context, surveyId, questionRef string) (*survey_proto.Question, error) {
	q := fmt.Sprintf(`
		FOR doc IN OUTBOUND "%v/%v" %v
		FILTER doc._key == "%v"
		RETURN doc`, common.DbSurveyTable, surveyId, common.DbSurveyQuestionEdgeTable, questionRef)

	resp, err := runQuery(ctx, q, common.DbSurveyQuestionTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToQuestion(resp.Records[0])
	return data, err
}

// CreateQuestion ...
func CreateQuestion(ctx context.Context, surveyId string, question *survey_proto.Question) error {
	if len(question.Id) == 0 {
		question.Id = uuid.NewUUID().String()
	}

	record, err := questionToRecord(question)
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
		IN %v`, question.Id, record, record, common.DbSurveyQuestionTable)
	_, err = runQuery(ctx, q, common.DbSurveyQuestionTable)
	if err != nil {
		return err
	}

	field := fmt.Sprintf(`{_from:"%v/%v",_to:"%v/%v"} `, common.DbSurveyTable, surveyId, common.DbSurveyQuestionTable, question.Id)
	q = fmt.Sprintf(`
		UPSERT %v
		INSERT %v
		UPDATE %v
		INTO %v`, field, field, field, common.DbSurveyQuestionEdgeTable)
	_, err = runQuery(ctx, q, common.DbSurveyQuestionTable)
	return err
}

func ByCreator(ctx context.Context, userId, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*survey_proto.Survey, error) {
	query := fmt.Sprintf(`FILTER doc.data.creator.id == "%v"`, userId)

	q := surveyQuery(orgId, teamId, offset, limit, sortParameter, sortDirection, query)

	var surveys []*survey_proto.Survey
	resp, err := runQuery(ctx, q, common.DbSurveyTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if survey, err := recordToSurvey(r); err == nil {
			surveys = append(surveys, survey)
		}
	}
	return surveys, nil
}

func Link(ctx context.Context, surveyId, orgId, teamId string) (string, error) {
	return "pending...", nil
}

// Templates ...
func Templates(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*survey_proto.Survey, error) {
	query := "FILTER doc.data.isTemplate == true"
	q := surveyQuery(orgId, "", offset, limit, sortParameter, sortDirection, query)

	var surveys []*survey_proto.Survey
	resp, err := runQuery(ctx, q, common.DbSurveyTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if survey, err := recordToSurvey(r); err == nil {
			surveys = append(surveys, survey)
		}
	}
	return surveys, nil
}

// Filter ...
func Filter(ctx context.Context, status []survey_proto.SurveyStatus, tags []string, renderTargets []survey_proto.RenderTarget, visibility []static_proto.Visibility, created_by []string, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*survey_proto.Survey, error) {
	query := "FILTER"
	if len(status) > 0 {
		t := []string{}
		for _, s := range status {
			t = append(t, fmt.Sprintf(`"%v"`, s))
		}
		query += fmt.Sprintf(" || doc.data.status IN [%v]", strings.Join(t[:], ","))
	}
	if len(tags) > 0 {
		g := common.QueryStringFromArray(tags)
		query += fmt.Sprintf(" || doc.data.tags ANY IN [%v]", g)
	}
	if len(renderTargets) > 0 {
		t := []string{}
		for _, s := range renderTargets {
			t = append(t, fmt.Sprintf(`"%v"`, s))
		}
		query += fmt.Sprintf(" || doc.data.renders ANY IN [%v]", strings.Join(t[:], ","))
	}
	if len(visibility) > 0 {
		t := []string{}
		for _, s := range visibility {
			t = append(t, fmt.Sprintf(`"%v"`, s))
		}
		query += fmt.Sprintf(" || doc.data.setting.visibility IN [%v]", strings.Join(t[:], ","))
	}
	// q := surveyQuery(orgId, teamId, offset, limit, sortParameter, sortDirection, query)

	if len(created_by) > 0 {
		t := []string{}
		for _, s := range created_by {
			t = append(t, fmt.Sprintf(`"%v"`, s))
		}
		query += fmt.Sprintf(" || doc.data.creator.id IN [%v]", strings.Join(t[:], ","))
	}

	q := surveyQuery(orgId, "", offset, limit, "", "", common.QueryClean(query))

	var surveys []*survey_proto.Survey
	resp, err := runQuery(ctx, q, common.DbSurveyTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if survey, err := recordToSurvey(r); err == nil {
			surveys = append(surveys, survey)
		}
	}
	return surveys, nil
}

// Searches surveys by name and/or ..., uses Elasticsearch middleware
func Search(ctx context.Context, name, description, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*survey_proto.Survey, error) {
	query := "FILTER"
	if len(name) > 0 {
		query += fmt.Sprintf(` && LIKE(doc.name, "%s",true)`, `%`+name+`%`)
	}
	if len(description) > 0 {
		query += fmt.Sprintf(` && doc.data.description == "%s"`, description)
	}
	q := surveyQuery(orgId, teamId, offset, limit, sortParameter, sortDirection, common.QueryClean(query))

	var surveys []*survey_proto.Survey
	resp, err := runQuery(ctx, q, common.DbSurveyTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if survey, err := recordToSurvey(r); err == nil {
			surveys = append(surveys, survey)
		}
	}
	return surveys, nil
}

func ShareSurvey(ctx context.Context, surveys []*survey_proto.Survey, users []*user_proto.User, sharedBy *user_proto.User, orgId string) ([]string, error) {
	userids := []string{}
	for _, survey := range surveys {
		for _, user := range users {
			shared := &survey_proto.ShareSurveyUser{
				Id:       uuid.NewUUID().String(),
				Survey:   survey,
				User:     user,
				Status:   static_proto.ShareStatus_SHARED,
				Updated:  time.Now().Unix(),
				Created:  time.Now().Unix(),
				SharedBy: sharedBy,
				Count:    int64(len(survey.Questions)),
			}

			_from := fmt.Sprintf(`%v/%v`, common.DbSurveyTable, survey.Id)
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
				RETURN {data:{user_id: OLD ? "" : NEW.data.user.id}}`, field, record, record, common.DbShareSurveyUserEdgeTable)

			resp, err := runQuery(ctx, q, common.DbShareSurveyUserEdgeTable)
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
			any, err := common.FilteredAnyFromObject(common.SURVEY_TYPE, survey.Id)
			if err != nil {
				return nil, err
			}
			// save pending
			pending := &common_proto.Pending{
				Id:         uuid.NewUUID().String(),
				OrgId:      orgId,
				Created:    shared.Created,
				Updated:    shared.Updated,
				SharedBy:   sharedBy,
				SharedWith: user,
				Item:       any,
			}

			q1, err1 := common.SavePending(pending, survey.Id)
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

func AutocompleteSearch(ctx context.Context, title string) ([]*static_proto.AutocompleteResponse, error) {
	response := []*static_proto.AutocompleteResponse{}
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER LIKE(doc.name, "%v",true)
		RETURN doc`, common.DbSurveyTable, `%`+title+`%`)

	resp, err := runQuery(ctx, q, common.DbSurveyTable)
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
		LET tags = (
			FOR doc IN %v
			%v
			RETURN doc.data.tags)[**]
		FOR t IN tags
		FILTER LIKE(t,"%v",true)
		LET ret = {parameter1:t}
		RETURN DISTINCT ret
		`, common.DbSurveyTable, query, `%`+name+`%`)
	resp, err := runQuery(ctx, q, common.DbSurveyTable)

	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		tags = append(tags, r.Parameter1)
	}
	return tags, nil
}

func GetShareableSurveys(ctx context.Context, createdBy []string, search_term, userId, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*user_proto.ShareableSurvey, error) {
	var response []*user_proto.ShareableSurvey
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	shared_query, filter_shared_query := queryFilterSharedForUser(userId, common.DbShareSurveyUserEdgeTable, query)

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
		FOR doc IN %v
		%s
		%s
		%s
		%s
		%s
		LET createdBy = (FOR p IN %v FILTER doc.data.createdBy.id == p._key RETURN p.data)
		RETURN {data:{
			id:doc.id,
			title:doc.name,
			org_id: doc.data.org_id,
			summary: doc.data.summary,
			shared_by: {"id": createdBy[0].id, "firstname": createdBy[0].firstname, "lastname": createdBy[0].lastname, "avatar_url": createdBy[0].avatar_url},
			count: doc.data.count
			response_time: doc.data.count * .1667
			
		}}`,
		shared_query,
		common.DbSurveyTable,
		filter_shared_query,
		query,
		common.QueryClean(filter_query),
		sort_query, limit_query,
		common.DbBehaviourCategoryAimTable, common.DbMarkerTable,
		common.DbUserTable,
	)

	resp, err := runQuery(ctx, q, common.DbSurveyTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if res, err := recordToShareableSurvey(r); err == nil {
			response = append(response, res)
		}
	}
	return response, nil
}
