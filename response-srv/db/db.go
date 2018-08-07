package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"server/common"
	db_proto "server/db-srv/proto/db"
	resp_proto "server/response-srv/proto/response"
	static_proto "server/static-srv/proto/static"
	survey_proto "server/survey-srv/proto/survey"

	"context"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/client"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type clientWrapper struct {
	Db_client             db_proto.DBClient
	ResponseServiceClient resp_proto.ResponseServiceClient
}

var (
	ClientWrapper *clientWrapper
	ErrNotFound   = errors.New("not found")
)

// Storage for a db microservice client
func NewClientWrapper(serviceClient client.Client) *clientWrapper {
	cl := db_proto.NewDBClient("", serviceClient)
	cl1 := resp_proto.NewResponseServiceClient("", serviceClient)

	return &clientWrapper{
		Db_client:             cl,
		ResponseServiceClient: cl1,
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

func responseToRecord(resp *resp_proto.SubmitSurveyResponse) (string, error) {
	data, err := common.MarhalToObject(resp)
	if err != nil {
		return "", err
	}

	common.FilterObject(data, "responder", resp.Responder)

	d := map[string]interface{}{
		"_key":       resp.Id,
		"id":         resp.Id,
		"created":    resp.Created,
		"updated":    resp.Updated,
		"parameter1": resp.OrgId,
		"parameter2": resp.SurveyId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToResponse(r *db_proto.Record) (*resp_proto.SubmitSurveyResponse, error) {
	var p resp_proto.SubmitSurveyResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToGroup(r *db_proto.Record) (*resp_proto.GroupByQuestionResponse, error) {
	var p resp_proto.GroupByQuestionResponse
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToSurvey(r *db_proto.Record) (*survey_proto.Survey, error) {
	var p survey_proto.Survey
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func recordToStats(r *db_proto.Record) (*resp_proto.ReadStatsResponse_Data_Stats, error) {
	var p resp_proto.ReadStatsResponse_Data_Stats
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func queryMerge() string {
	query := fmt.Sprintf(`
		LET responder = (FOR p IN %v FILTER doc.data.responder.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{responder:{
			id: responder[0].id,
			firstname: responder[0].firstname,
			lastname: responder[0].lastname,
			avatar_url: responder[0].avatar_url
		}}})`,
		common.DbUserTable)
	return query
}

// Check survey status
func Check(ctx context.Context, shortHash, orgId, teamId string) (*resp_proto.CheckResponse_Data_Survey, error) {
	var survey resp_proto.CheckResponse_Data_Survey
	//
	q := fmt.Sprintf(`
		FOR h IN %v
		FILTER h.name == "%v" 
			FOR s IN %v
			FILTER s._key == h.parameter1
				RETURN s
		`, common.DbSurveyHashTable, shortHash, common.DbSurveyTable)

	resp, err := runQuery(ctx, q, common.DbSurveyTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}
	// parsing
	s, err := recordToSurvey(resp.Records[0])
	survey.Id = s.Id
	survey.AuthenticationRequired = s.Setting.AuthenticationRequired

	return &survey, nil
}

func All(ctx context.Context, surveyId, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*resp_proto.SubmitSurveyResponse, error) {
	var resps []*resp_proto.SubmitSurveyResponse
	query := common.QueryAuth(`FILTER`, orgId, surveyId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbResponseTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbResponseTable)
	if err != nil {
		common.ErrorLog(common.ResponseSrv, All, err, "All query is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if resp, err := recordToResponse(r); err == nil {
			resps = append(resps, resp)
		}
	}
	return resps, nil
}

func Create(ctx context.Context, surveyId string, response *resp_proto.SubmitSurveyResponse) error {
	response.SurveyId = surveyId
	if len(response.Id) == 0 {
		response.Id = uuid.NewUUID().String()
	}
	if response.Created == 0 {
		response.Created = time.Now().Unix()
	}
	response.Updated = time.Now().Unix()
	record, err := responseToRecord(response)
	if err != nil {
		common.ErrorLog(common.ResponseSrv, Create, err, "ResponseToRecord is failed")
		return err
	}
	if len(record) == 0 {
		common.ErrorLog(common.ResponseSrv, Create, err, "server serialization")
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, response.Id, record, record, common.DbResponseTable)
	if _, err := runQuery(ctx, q, common.DbResponseTable); err != nil {
		common.ErrorLog(common.ResponseSrv, Create, err, "DbResponseTable query is failed")
		return err
	}
	field := fmt.Sprintf(`{_from:"%v/%v",_to:"%v/%v"} `, common.DbSurveyTable, surveyId, common.DbResponseTable, response.Id)
	q = fmt.Sprintf(`
		UPSERT %v
		INSERT %v
		UPDATE %v
		INTO %v`, field, field, field, common.DbSurveyResponseEdgeTable)
	if _, err = runQuery(ctx, q, common.DbSurveyResponseEdgeTable); err != nil {
		common.ErrorLog(common.ResponseSrv, Create, err, "DbSurveyResponseEdgeTable query is failed")
		return err
	}
	return nil
}

func UpdateState(ctx context.Context, surveyId, responseId string, state resp_proto.ResponseState, orgId, teamId string) (*resp_proto.SubmitSurveyResponse, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, responseId)
	query = common.QueryAuth(query, orgId, surveyId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		UPDATE doc WITH {"data":{"status":{"state":"%v"}}} IN %v
		RETURN NEW`, common.DbResponseTable, query, state, common.DbResponseTable)

	resp, err := runQuery(ctx, q, common.DbResponseTable)
	if err != nil || len(resp.Records) == 0 {
		common.ErrorLog(common.ResponseSrv, UpdateState, err, "RunQuery is failed")
		return nil, err
	}
	// parsing
	data, err := recordToResponse(resp.Records[0])
	return data, err
}

func AllState(ctx context.Context, surveyId string, state resp_proto.ResponseState, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*resp_proto.SubmitSurveyResponse, error) {
	var resps []*resp_proto.SubmitSurveyResponse

	query := fmt.Sprintf(`FILTER doc.data.status.state == "%v"`, state)
	query = common.QueryAuth(query, orgId, surveyId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbResponseTable, query, sort_query, limit_query, queryMerge())
	resp, err := runQuery(ctx, q, common.DbResponseTable)
	if err != nil {
		common.ErrorLog(common.ResponseSrv, AllState, err, "RunQuery is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if resp, err := recordToResponse(r); err == nil {
			resps = append(resps, resp)
		}
	}
	return resps, nil
}

func AllAggQuestion(ctx context.Context, surveyId, orgId, teamId string) (*resp_proto.GroupByQuestionResponses, error) {
	res := &resp_proto.GroupByQuestionResponses{
		SurveyId: surveyId,
		OrgId:    orgId,
	}
	var groups []*resp_proto.GroupByQuestionResponse

	q := fmt.Sprintf(`
		LET arr = (
			FOR r IN %v
			FILTER r.parameter2 == "%v" && r.parameter1 == "%v"
			RETURN MERGE_RECURSIVE (r.data.answers[0], {created: r.created})
		)[**]
		LET questions = (FOR a IN arr
			RETURN {question_ref:a.question_ref, type:a.type}
		)
		LET refs = (FOR a IN questions
			RETURN DISTINCT a
		)
		FOR ref IN refs
		LET skipped_count = (
			FOR r IN arr
			FILTER r.question_ref == ref.question_ref && r.data == null
			COLLECT WITH COUNT INTO length
			RETURN length
		)[0]
		LET response_count = (
			FOR r IN arr
			FILTER r.question_ref == ref.question_ref && r.data != null
			COLLECT WITH COUNT INTO length
			RETURN length
		)[0]
		LET answers = (
			FOR r IN arr
			FILTER r.question_ref == ref.question_ref
			LET responder = (FOR p IN user FILTER p._key == r.user_id RETURN p.data)
			RETURN MERGE_RECURSIVE(r,{responder:{
                			id: responder[0].id,
                			firstname: responder[0].firstname,
                			lastname: responder[0].lastname,
                			avatar_url: responder[0].avatar_url
                		}})
		)
		RETURN {data:{
		question_ref:ref.question_ref,
		type:ref.type,
		response_count:response_count,
		skipped_count:skipped_count,
		answers:answers
		}}`, common.DbResponseTable, surveyId, orgId)
	resp, err := runQuery(ctx, q, common.DbResponseTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if resp, err := recordToGroup(r); err == nil {
			groups = append(groups, resp)
		}
	}
	res.Responses = groups
	return res, nil
}

func TimeFilter(ctx context.Context, surveyId string, from, to int64, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*resp_proto.SubmitSurveyResponse, error) {
	var resps []*resp_proto.SubmitSurveyResponse

	query := fmt.Sprintf("FILTER %v < doc.created && doc.created < %v", from, to)
	query = common.QueryAuth(query, orgId, surveyId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbResponseTable, query, sort_query, limit_query, queryMerge())
	resp, err := runQuery(ctx, q, common.DbResponseTable)
	if err != nil {
		common.ErrorLog(common.ResponseSrv, TimeFilter, err, "TimeFilter query is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if resp, err := recordToResponse(r); err == nil {
			resps = append(resps, resp)
		}
	}
	return resps, nil
}

func ReadStats(ctx context.Context, surveyId string) (*resp_proto.ReadStatsResponse_Data_Stats, error) {
	query := fmt.Sprintf(`FILTER doc.parameter2 == "%v"`, surveyId)
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET responses = (
			FILTER doc.data.status.state == "%v"
			COLLECT WITH COUNT INTO length
			RETURN length
		)
		LET drops = (
			FILTER doc.data.status.state == "%v"
			COLLECT WITH COUNT INTO length
			RETURN length
		)
		RETURN {
			"data":{
				"responses":responses[0], 
				"drops":drops[0]
			}
		}`, common.DbResponseTable, query, resp_proto.ResponseState_SUBMITTED, resp_proto.ResponseState_ABANDONED)
	resp, err := runQuery(ctx, q, common.DbResponseTable)
	if err != nil || len(resp.Records) == 0 {
		common.ErrorLog(common.ResponseSrv, ReadStats, err, "ReadStats query is failed")
		return nil, err
	}
	// parsing
	stats, err := recordToStats(resp.Records[0])
	return stats, err
}

func ByUser(ctx context.Context, surveyId, userId string, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*resp_proto.SubmitSurveyResponse, error) {
	var resps []*resp_proto.SubmitSurveyResponse

	query := fmt.Sprintf(`FILTER doc.data.responder.id == "%v"`, userId)
	query = common.QueryAuth(query, orgId, surveyId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		FOR s IN %v
		FILTER s._key == "%v" && s.data.setting.authenticationRequired == true
		FOR e IN %v
		FILTER e._from == s._id && e._to == doc._id
		%s
		%s
		%s`, common.DbResponseTable, query, common.DbSurveyTable, surveyId, common.DbSurveyResponseEdgeTable,
		sort_query, limit_query, queryMerge(),
	)
	resp, err := runQuery(ctx, q, common.DbResponseTable)
	if err != nil {
		common.ErrorLog(common.ResponseSrv, ByUser, err, "ByUser query is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if resp, err := recordToResponse(r); err == nil {
			resps = append(resps, resp)
		}
	}
	return resps, nil
}

func ByAnyUser(ctx context.Context, surveyId string, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*resp_proto.SubmitSurveyResponse, error) {
	var resps []*resp_proto.SubmitSurveyResponse

	query := `FILTER doc.data.responder == null`
	query = common.QueryAuth(query, orgId, surveyId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		FOR s IN %v
		FILTER s._key == "%v" && s.data.setting.authenticationRequired != true
		FOR e IN %v
		FILTER e._from == s._id && e._to == doc._id
		%s
		%s
		%s`, common.DbResponseTable, query, common.DbSurveyTable, surveyId, common.DbSurveyResponseEdgeTable,
		sort_query, limit_query, queryMerge(),
	)
	resp, err := runQuery(ctx, q, common.DbResponseTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if resp, err := recordToResponse(r); err == nil {
			resps = append(resps, resp)
		}
	}
	return resps, nil
}

func RemovePendingSharedAction(ctx context.Context, id string) error {
	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc.item.id == "%v"
		REMOVE doc IN %v`, common.DbPendingTable, id, common.DbPendingTable)

	_, err := runQuery(ctx, q, common.DbPendingTable)
	return err
}

func UpdateShareSurveyStatus(ctx context.Context, surveyId string, status static_proto.ShareStatus) error {
	_from := fmt.Sprintf(`%v/%v`, common.DbSurveyTable, surveyId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._from == "%v"
		UPDATE doc WITH {data:{status:"%v"}} IN %v`, common.DbShareSurveyUserEdgeTable, _from, status, common.DbShareSurveyUserEdgeTable)
	_, err := runQuery(ctx, q, common.DbShareSurveyUserEdgeTable)
	return err
}
