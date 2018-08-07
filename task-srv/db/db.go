package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"server/common"
	db_proto "server/db-srv/proto/db"
	task_proto "server/task-srv/proto/task"
	"strconv"
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

func taskToRecord(task *task_proto.Task) (string, error) {
	data, err := common.MarhalToObject(task)
	if err != nil {
		return "", err
	}

	common.FilterObject(data, "user", task.User)
	common.FilterObject(data, "creator", task.Creator)
	common.FilterObject(data, "assignee", task.Assignee)
	var creatorId string
	if task.Creator != nil {
		creatorId = task.Creator.Id
	}

	d := map[string]interface{}{
		"_key":       task.Id,
		"id":         task.Id,
		"created":    task.Created,
		"updated":    task.Updated,
		"name":       task.Title,
		"parameter1": task.OrgId,
		"parameter2": creatorId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToTask(r *db_proto.Record) (*task_proto.Task, error) {
	var p task_proto.Task
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func queryMerge() string {
	query := fmt.Sprintf(`
		LET user = (FOR p IN %v FILTER doc.data.user.id == p._key RETURN p.data)
		LET creator = (FOR p IN %v FILTER doc.data.creator.id == p._key RETURN p.data)
		LET assignee = (FOR p IN %v FILTER doc.data.assignee.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{user:user[0],creator:creator[0],assignee:assignee[0]}})`,
		common.DbUserTable, common.DbUserTable, common.DbUserTable)
	return query
}

// All get all tasks
func All(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*task_proto.Task, error) {
	var tasks []*task_proto.Task
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbTaskTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbTaskTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if task, err := recordToTask(r); err == nil {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

// Creates a task
func Create(ctx context.Context, task *task_proto.Task) error {
	if len(task.Id) == 0 {
		task.Id = uuid.NewUUID().String()
	}
	if task.Created == 0 {
		task.Created = time.Now().Unix()
	}
	task.Updated = time.Now().Unix()

	record, err := taskToRecord(task)
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
		INTO %v`, task.Id, record, record, common.DbTaskTable)
	_, err = runQuery(ctx, q, common.DbTaskTable)
	return err
}

// Reads a task by ID
func Read(ctx context.Context, id, orgId, teamId string) (*task_proto.Task, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v 
		%s
		RETURN doc`, common.DbTaskTable, query)

	resp, err := runQuery(ctx, q, common.DbTaskTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToTask(resp.Records[0])
	return data, err
}

// Deletes a task by ID
func Delete(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbTaskTable, query, common.DbTaskTable)
	_, err := runQuery(ctx, q, common.DbTaskTable)
	return err
}

// Searches tasks by name and/or ..., uses Elasticsearch middleware
func Search(ctx context.Context, name, orgId, teamId string, limit, offset, from, to int64, sortParameter, sortDirection string) ([]*task_proto.Task, error) {
	var tasks []*task_proto.Task
	query := `FILTER`
	if len(name) > 0 {
		query += fmt.Sprintf(` doc.name == "%v"`, name)
	}
	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbTaskTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbTaskTable)
	if err != nil {
		common.ErrorLog(common.TaskSrv, Search, err, "Search query is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if task, err := recordToTask(r); err == nil {
			tasks = append(tasks, task)
		}
	}

	return tasks, nil
}

// ByCreator get all tasks created by a particular team member
func ByCreator(ctx context.Context, id, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*task_proto.Task, error) {
	query := fmt.Sprintf(`FILTER doc.data.creator.id == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	var tasks []*task_proto.Task
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbTaskTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbTaskTable)
	if err != nil {
		common.ErrorLog(common.TaskSrv, ByCreator, err, "ByCreator query is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if task, err := recordToTask(r); err == nil {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

// ByAssign get all tasks assigned to a particular team member
func ByAssign(ctx context.Context, id, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*task_proto.Task, error) {
	query := fmt.Sprintf(`FILTER doc.data.assignee.id == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	var tasks []*task_proto.Task
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbTaskTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbTaskTable)
	if err != nil {
		common.ErrorLog(common.TaskSrv, ByAssign, err, "ByAssign query is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if task, err := recordToTask(r); err == nil {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

// Filter ...
func Filter(ctx context.Context, status []task_proto.TaskStatus, category []string, priority []int64, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*task_proto.Task, error) {
	var tasks []*task_proto.Task

	// integrate status
	_status := []string{}
	for _, s := range status {
		_status = append(_status, fmt.Sprintf(`"%v"`, s))
	}
	s := strings.Join(_status[:], ",")

	c := common.QueryStringFromArray(category)
	_priority := []string{}
	for _, p := range priority {
		_priority = append(_priority, strconv.FormatInt(p, 10))
	}
	p := strings.Join(_priority[:], ",")

	query := fmt.Sprintf(`FILTER doc.data.status IN [%v] OR doc.data.category IN [%v] OR doc.data.priority IN [%v]`, s, c, p)
	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbTaskTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbTaskTable)
	if err != nil {
		common.ErrorLog(common.TaskSrv, Filter, err, "Filter query is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if task, err := recordToTask(r); err == nil {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func CountByUser(ctx context.Context, id, orgId, teamId string) (*task_proto.CountByUserResponse_TaskCount, error) {
	query1 := fmt.Sprintf(`FILTER doc.due < %v`, time.Now().Unix())
	query1 = common.QueryAuth(query1, orgId, teamId)
	query2 := fmt.Sprintf(`FILTER doc.data.assignee.id == "%v"`, id)
	query2 = common.QueryAuth(query2, orgId, teamId)

	count := &task_proto.CountByUserResponse_TaskCount{}
	q := fmt.Sprintf(`
		LET expired = (
			FOR doc IN %v
				%s
				COLLECT WITH COUNT INTO length
				RETURN length)
		LET assigned = (
			FOR doc IN %v
				%s
				COLLECT WITH COUNT INTO length
				RETURN length)
		RETURN {"parameter1":TO_STRING(expired[0]), "parameter2":TO_STRING(assigned[0])}`, common.DbTaskTable, query1, common.DbTaskTable, query2)

	resp, err := runQuery(ctx, q, common.DbTaskTable)
	if err != nil {
		common.ErrorLog(common.TaskSrv, CountByUser, err, "CountByUser query is failed")
		return nil, err
	}
	// parsing
	if len(resp.Records) > 0 {
		count.Expired, _ = strconv.ParseInt(resp.Records[0].Parameter1, 10, 64)
		count.Assigned, _ = strconv.ParseInt(resp.Records[0].Parameter2, 10, 64)
	}
	return count, err
}
