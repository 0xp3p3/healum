package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"server/common"
	db_proto "server/db-srv/proto/db"
	note_proto "server/note-srv/proto/note"
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

func noteToRecord(note *note_proto.Note) (string, error) {
	data, err := common.MarhalToObject(note)
	if err != nil {
		return "", err
	}

	common.FilterObject(data, "user", note.User)
	common.FilterObject(data, "creator", note.Creator)
	var creatorId string
	if note.Creator != nil {
		creatorId = note.Creator.Id
	}

	d := map[string]interface{}{
		"_key":       note.Id,
		"id":         note.Id,
		"created":    note.Created,
		"updated":    note.Updated,
		"name":       note.Title,
		"parameter1": note.OrgId,
		"parameter2": creatorId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToNote(r *db_proto.Record) (*note_proto.Note, error) {
	var p note_proto.Note
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func queryMerge() string {
	query := fmt.Sprintf(`
		LET user = (FOR p IN %v FILTER doc.data.user.id == p._key RETURN p.data)
		LET creator = (FOR p IN %v FILTER doc.data.creator.id == p._key RETURN p.data)
		RETURN MERGE_RECURSIVE(doc,{data:{user:user[0],creator:creator[0]}})`,
		common.DbUserTable, common.DbUserTable)
	return query
}

// All get all notes
func All(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*note_proto.Note, error) {
	var notes []*note_proto.Note
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbNoteTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbNoteTable)
	if err != nil {
		common.ErrorLog(common.NoteSrv, All, err, "RunQuery is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if note, err := recordToNote(r); err == nil {
			notes = append(notes, note)
		}
	}
	return notes, nil
}

// Creates a note
func Create(ctx context.Context, note *note_proto.Note) error {
	if len(note.Id) == 0 {
		note.Id = uuid.NewUUID().String()
	}
	if note.Created == 0 {
		note.Created = time.Now().Unix()
	}
	note.Updated = time.Now().Unix()

	record, err := noteToRecord(note)
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
		INTO %v`, note.Id, record, record, common.DbNoteTable)
	_, err = runQuery(ctx, q, common.DbNoteTable)
	return err
}

// Reads a note by ID
func Read(ctx context.Context, id, orgId, teamId string) (*note_proto.Note, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v 
		%s
		%s`, common.DbNoteTable, query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbNoteTable)
	if err != nil || len(resp.Records) == 0 {
		common.ErrorLog(common.NoteSrv, Read, err, "RunQuery is failed")
		return nil, err
	}
	data, err := recordToNote(resp.Records[0])
	return data, err
}

// Deletes a note by ID
func Delete(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbNoteTable, query, common.DbNoteTable)
	_, err := runQuery(ctx, q, common.DbNoteTable)
	return err
}

// Searches notes by name and/or ..., uses Elasticsearch middleware
func Search(ctx context.Context, name, orgId, teamId string, limit, offset, from, to int64, sortParameter, sortDirection string) ([]*note_proto.Note, error) {
	var notes []*note_proto.Note
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
		%s`, common.DbNoteTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbNoteTable)
	if err != nil {
		common.ErrorLog(common.NoteSrv, Search, err, "RunQuery is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if note, err := recordToNote(r); err == nil {
			notes = append(notes, note)
		}
	}
	return notes, nil
}

// ByCreator get all notes created by a particular team member
func ByCreator(ctx context.Context, id, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*note_proto.Note, error) {
	query := fmt.Sprintf(`FILTER doc.data.creator.id == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	var notes []*note_proto.Note
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbNoteTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbNoteTable)
	if err != nil {
		common.ErrorLog(common.NoteSrv, ByCreator, err, "RunQuery is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if note, err := recordToNote(r); err == nil {
			notes = append(notes, note)
		}
	}
	return notes, nil
}

// ByUser get all notes created for a user
func ByUser(ctx context.Context, id, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*note_proto.Note, error) {
	query := fmt.Sprintf(`FILTER doc.data.user.id == "%v"`, id)
	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	var notes []*note_proto.Note
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbNoteTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbNoteTable)
	if err != nil {
		common.ErrorLog(common.NoteSrv, ByUser, err, "RunQuery is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if note, err := recordToNote(r); err == nil {
			notes = append(notes, note)
		}
	}
	return notes, nil
}

// Filter ...
func Filter(ctx context.Context, category, tags []string, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*note_proto.Note, error) {
	var notes []*note_proto.Note

	// integrate status
	_categorys := common.QueryStringFromArray(category)
	_tags := common.QueryStringFromArray(tags)

	query := fmt.Sprintf(`FILTER (doc.data.category IN [%v] OR doc.data.tags ANY IN [%v])`, _categorys, _tags)
	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		%s`, common.DbNoteTable, query, sort_query, limit_query, queryMerge())

	resp, err := runQuery(ctx, q, common.DbNoteTable)
	if err != nil {
		common.ErrorLog(common.NoteSrv, Filter, err, "RunQuery is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if note, err := recordToNote(r); err == nil {
			notes = append(notes, note)
		}
	}
	return notes, nil
}
