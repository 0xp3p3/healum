package common

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	db_proto "server/db-srv/proto/db"
	common_proto "server/static-srv/proto/common"
	static_proto "server/static-srv/proto/static"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
)

const (
	AudioContentType = "audio"
	VideoContentType = "video"
	ImageContentType = "image"

	SearchableMeta             = "searchable"
	SearchableAutocompleteMeta = "autocomplete"

	SearchAllFlag          = "all"
	SearchAutocompleteFlag = "autocomplete"

	GraphFlag = "graph"
	GraphName = "graphname"
	GraphFrom = "from"
	GraphTo   = "to"

	DefaultOrg = "default"
)

var (
	SearchableMetaMap map[string]string = map[string]string{
		SearchableMeta: "",
	}

	SearchableAutocompleteMetaMap map[string]string = map[string]string{
		SearchableMeta:             "",
		SearchableAutocompleteMeta: "",
	}

	GraphMap map[string]string = map[string]string{
		GraphFlag: "",
		GraphName: "",
		GraphFrom: "",
		GraphTo:   "",
	}
)

func NameFromContext(ctx context.Context, name string) string {

	if IsTestContext(ctx) {
		return TestingName(name)
	}
	return name
}

//This function returns a complete Any Object from an object
func AnyFromObject(msg proto.Message, typeurl string) (*google_protobuf.Any, error) {
	mapp, err := MarhalToObject(msg)
	if err != nil {
		return nil, err
	}

	mapp["@type"] = BASE + typeurl
	// getting json string from json object
	b, err := json.Marshal(mapp)
	if err != nil {
		return nil, err
	}
	// getting Any object from json string
	obj := &google_protobuf.Any{}
	if err := jsonpb.Unmarshal(strings.NewReader(string(b)), obj); err != nil {
		return nil, err
	}
	return obj, nil
}

//This function returns a filtered Any Object from an object (only @type and id of the object is returned)
func FilteredAnyFromObject(typeurl, id string) (*google_protobuf.Any, error) {
	mapp := map[string]string{"@type": BASE + typeurl, "id": id}
	// getting json string from json object
	b, err := json.Marshal(mapp)
	if err != nil {
		return nil, err
	}
	// getting Any object from json string
	obj := &google_protobuf.Any{}
	if err := jsonpb.Unmarshal(strings.NewReader(string(b)), obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func MarhalToObject(msg proto.Message, options ...interface{}) (map[string]interface{}, error) {
	marshaler := jsonpb.Marshaler{}
	if len(options) == 0 {
		marshaler.EmitDefaults = true
	}
	marshaler.OrigName = true
	// marshl object to json string
	js, err := marshaler.MarshalToString(msg)
	if err != nil {
		return nil, err
	}
	// getting json object from json string
	var data map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(js)))
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func PendingToRecord(pending *common_proto.Pending) (string, error) {
	//FIXME:change TempMarshalToObjectForPending  this after merge to MarhalToObject with EmitDefaults as false
	data, err := TempMarshalToObjectForPending(pending)
	if err != nil {
		return "", err
	}

	FilterObject(data, "shared_by", pending.SharedBy)
	FilterObject(data, "shared_with", pending.SharedWith)

	d := map[string]interface{}{
		"_key":       pending.Id,
		"id":         pending.Id,
		"created":    pending.Created,
		"parameter1": pending.OrgId,
		"parameter2": pending.SharedWith.Id,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func SavePending(pending *common_proto.Pending, itemId string) (string, error) {
	record, err := PendingToRecord(pending)
	if err != nil {
		return "", err
	}
	if len(record) == 0 {
		return "", errors.New("server serialization")
	}

	field := fmt.Sprintf(`{data:{item:{id:"%v"}}, parameter1:"%v", parameter2:"%v"} `, itemId, pending.OrgId, pending.SharedWith.Id)
	q := fmt.Sprintf(`
		UPSERT %v
		INSERT %v
		UPDATE %v
		INTO %v`, field, record, record, DbPendingTable)
	return q, nil
}

func GetTrackKey(user_id, obj_id, t string) string {
	return user_id + ":" + obj_id + ":" + t
}

func QueryAuth(query, orgId, teamId string) string {
	if len(orgId) != 0 {
		query += fmt.Sprintf(` && doc.parameter1 == "%s"`, orgId)
	}
	if len(teamId) != 0 {
		query += fmt.Sprintf(` && doc.parameter2 == "%s"`, teamId)
	}

	return QueryClean(query)
}

func QueryClean(query string) string {
	query = strings.Replace(query, `FILTER && `, `FILTER `, -1)
	query = strings.Replace(query, `FILTER || `, `FILTER `, -1)

	if query == `FILTER` {
		query = ""
	}
	return query
}

func QueryPaginate(offset, limit int64) string {
	if limit == 0 {
		limit = 10
	}
	offs := fmt.Sprintf("%d", offset)
	size := fmt.Sprintf("%d", limit)

	limit_query := fmt.Sprintf("LIMIT %s, %s", offs, size)
	return limit_query
}

func QuerySort(sortParameter, sortDirection string, document ...string) string {
	doc := "doc"
	if len(document) > 0 {
		doc = document[0]
	}

	if sortDirection != "ASC" && sortDirection != "DESC" {
		sortDirection = "DESC"
	}
	// default sort query
	if len(sortParameter) == 0 || len(sortDirection) == 0 {
		return fmt.Sprintf("SORT %v.data.updated DESC", doc)
	}

	return fmt.Sprintf(`SORT %v.data.%v %v`, doc, sortParameter, sortDirection)
}

func QueryStringFromArray(arr []string) string {
	strs := []string{}
	for _, s := range arr {
		strs = append(strs, `"`+s+`"`)
	}
	return strings.Join(strs[:], ",")
}

func FilterObject(data map[string]interface{}, key string, msg proto.Message) {
	if !reflect.ValueOf(msg).IsNil() {
		obj, err := MarhalToObject(msg)
		if obj != nil && err == nil {
			data[key] = map[string]string{
				"id": obj["id"].(string),
			}
			return
		}
	}
	delete(data, key)
}

func FilterArrayObject(data map[string]interface{}, key string, msg []proto.Message) {
	if len(msg) > 0 && !reflect.ValueOf(msg[0]).IsNil() {
		var arr []interface{}
		for _, item := range msg {
			obj, err := MarhalToObject(item)
			if obj != nil && err == nil {
				arr = append(arr, map[string]string{
					"id": obj["id"].(string),
				})
			}
		}
		data[key] = arr
	} else {
		delete(data, key)
	}
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

//FIXME:remove this after merge
func TempMarshalToObjectForPending(msg proto.Message) (map[string]interface{}, error) {
	marshaler := jsonpb.Marshaler{}
	marshaler.OrigName = true
	// marshl object to json string
	js, err := marshaler.MarshalToString(msg)
	if err != nil {
		return nil, err
	}
	// getting json object from json string
	var data map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(js)))
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

//This method return filter query for searching shared_by values for a specific term
//ALERT! uses doc.data. - changing the object name from doc to something else will cause error
func QuerySharedResourceSearch(filter_query string, term, doc string) string {
	if len(term) > 0 {
		filter_query += fmt.Sprintf(` && (LIKE(%v.name, "%s",true) || LIKE(%v.data.description, "%v",true) || LIKE(%v.data.summary, "%v",true))`, doc, `%`+term+`%`, doc, `%`+term+`%`, doc, `%`+term+`%`)
	}
	return QueryClean(filter_query)
}

//if upsert results in an update => "" / insert => user_id
func RecordToInsertedUserId(r *db_proto.Record) (string, error) {
	var p static_proto.SharedUserId
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return "", err
	}
	if len(p.UserId) > 0 {
		return p.UserId, nil
	}
	return "", nil
}

func RecordToModule(r *db_proto.Record) (*static_proto.Module, error) {
	var p static_proto.Module
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}
