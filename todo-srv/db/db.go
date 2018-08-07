package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"server/common"
	db_proto "server/db-srv/proto/db"
	todo_proto "server/todo-srv/proto/todo"
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
	query := fmt.Sprintf(`
		LET creator = (FOR p IN %v FILTER doc.data.creator.id == p._key RETURN p.data)
		LET items = (FILTER NOT_NULL(doc.data.items) FOR p IN %v FOR item IN doc.data.items FILTER item.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{creator:creator[0], items:items}})`,
		common.DbUserTable, common.DbTodoItemTable)
	return query
}

func todoToRecord(todo *todo_proto.Todo) (string, error) {
	data, err := common.MarhalToObject(todo)
	if err != nil {
		return "", err
	}

	common.FilterObject(data, "creator", todo.Creator)
	//items
	if len(todo.Items) > 0 {
		var arr []interface{}
		for _, item := range todo.Items {
			arr = append(arr, map[string]string{"id": item.Id})
		}
		data["items"] = arr
	} else {
		delete(data, "items")
	}

	d := map[string]interface{}{
		"_key":       todo.Id,
		"id":         todo.Id,
		"created":    todo.Created,
		"updated":    todo.Updated,
		"name":       todo.Name,
		"parameter1": todo.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToTodo(r *db_proto.Record) (*todo_proto.Todo, error) {
	var p todo_proto.Todo
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		common.ErrorLog(common.TodoSrv, recordToTodo, err, "Unmarshale is failed")
		return nil, err
	}
	return &p, nil
}

// All get all todos
func All(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*todo_proto.Todo, error) {
	var todos []*todo_proto.Todo
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbTodoTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbTodoTable)
	if err != nil {
		common.ErrorLog(common.TodoSrv, All, err, "All query is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if todo, err := recordToTodo(r); err == nil {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

// Creates a todo
func Create(ctx context.Context, todo *todo_proto.Todo) error {
	if len(todo.Id) == 0 {
		todo.Id = uuid.NewUUID().String()
	}
	if todo.Created == 0 {
		todo.Created = time.Now().Unix()
	}
	if todo.Updated == 0 {
		todo.Updated = time.Now().Unix()
	}

	record, err := todoToRecord(todo)
	if err != nil {
		common.ErrorLog(common.TodoSrv, Create, err, "Marshal is failed")
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		INTO %v`, todo.Id, record, record, common.DbTodoTable)
	_, err = runQuery(ctx, q, common.DbTodoTable)
	return err
}

// Reads a todo by ID
func Read(ctx context.Context, id, orgId, teamId string) (*todo_proto.Todo, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v 
		%s
		%s`, common.DbTodoTable, query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbTodoTable)
	if err != nil || len(resp.Records) == 0 {
		common.ErrorLog(common.TodoSrv, Read, err, "Read query is failed")
		return nil, err
	}

	data, err := recordToTodo(resp.Records[0])
	return data, err
}

// Deletes a todo by ID
func Delete(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbTodoTable, query, common.DbTodoTable)
	_, err := runQuery(ctx, q, common.DbTodoTable)
	return err
}

// Searches todos by name and/or ..., uses Elasticsearch middleware
func Search(ctx context.Context, name, orgId, teamId string, limit, offset, from, to int64, sortParameter, sortDirection string) ([]*todo_proto.Todo, error) {
	var todos []*todo_proto.Todo
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
		%s`, common.DbTodoTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbNoteTable)
	if err != nil {
		common.ErrorLog(common.TodoSrv, Search, err, "Search query is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if todo, err := recordToTodo(r); err == nil {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}

// ByCreator get all todos created by a particular team member
func ByCreator(ctx context.Context, id, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*todo_proto.Todo, error) {
	query := fmt.Sprintf(`FILTER doc.data.creator.id == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	var todos []*todo_proto.Todo
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbTodoTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbTodoTable)
	if err != nil {
		common.ErrorLog(common.TodoSrv, ByCreator, err, "ByCreator query is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if todo, err := recordToTodo(r); err == nil {
			todos = append(todos, todo)
		}
	}
	return todos, nil
}
