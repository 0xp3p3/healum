package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"server/common"
	db_proto "server/db-srv/proto/db"
	organisation_proto "server/organisation-srv/proto/organisation"
	static_proto "server/static-srv/proto/static"
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
	if _, err := ClientWrapper.Db_client.InitDb(context.TODO(), &db_proto.InitDbRequest{}); err != nil {
		log.Fatal(err)
		return err
	}
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
		LET owner = (FOR owner in %v FILTER doc.data.owner.id == owner._key RETURN owner.data)
		LET parent = (FOR parent in %v FILTER doc.data.parent.id == parent._key RETURN parent.data)
		LET child = (FILTER NOT_NULL(doc.data.child) FOR p IN doc.data.child FOR child IN %v FILTER p.id == child._key RETURN child.data)
		LET modules = (
		    FILTER NOT_NULL(doc.data.modules) 
			FOR m,om in OUTBOUND doc._id %v
			return {
				id:m.id,
				name: m.data.name,
				icon_slug: m.data.icon_slug,
				name_slug: m.data.name_slug,
				description: m.data.description,
				summary: m.data.summary,
				display_name: NOT_NULL(om.data.display_name,m.data.name),
				path: m.data.path
			}
		)
		RETURN MERGE_RECURSIVE(doc,{data:{owner:owner[0],parent:parent[0], child:child,modules:modules}})`,
		common.DbUserTable,
		common.DbOrganisationTable,
		common.DbOrganisationTable,
		common.DbOrgModuleEdgeTable,
	)
	return query
}

func orgToRecord(org *organisation_proto.Organisation, options ...interface{}) (string, error) {
	data, err := common.MarhalToObject(org, options)
	if err != nil {
		return "", err
	}
	common.FilterObject(data, "owner", org.Owner)
	common.FilterObject(data, "parent", org.Parent)
	//child organisations
	if len(org.Child) > 0 {
		var arr []interface{}
		for _, item := range org.Child {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["child"] = arr
	} else {
		delete(data, "child")
	}

	//modules
	if len(org.Modules) > 0 {
		var arr []interface{}
		for _, item := range org.Modules {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["modules"] = arr
	} else {
		delete(data, "modules")
	}

	d := map[string]interface{}{
		"_key":    org.Id,
		"id":      org.Id,
		"name":    org.Name,
		"created": org.Created,
		"updated": org.Updated,
		"data":    data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToOrg(r *db_proto.Record) (*organisation_proto.Organisation, error) {
	var p organisation_proto.Organisation
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func profileToRecord(profile *organisation_proto.OrganisationProfile) (string, error) {
	data, err := common.MarhalToObject(profile)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    profile.Id,
		"id":      profile.Id,
		"created": profile.Created,
		"updated": profile.Updated,
		"data":    data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToProfile(r *db_proto.Record) (*organisation_proto.OrganisationProfile, error) {
	var p organisation_proto.OrganisationProfile
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func settingToRecord(setting *organisation_proto.OrganisationSetting) (string, error) {
	data, err := common.MarhalToObject(setting)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    setting.Id,
		"id":      setting.Id,
		"created": setting.Created,
		"updated": setting.Updated,
		"data":    data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToSetting(r *db_proto.Record) (*organisation_proto.OrganisationSetting, error) {
	var p organisation_proto.OrganisationSetting
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func organisationModuleToRecord(from, to string, module *organisation_proto.OrganisationModule) (string, error) {
	data, err := common.MarhalToObject(module)
	if err != nil {
		return "", err
	}

	if len(module.DisplayName) == 0 {
		delete(data, "display_name")
	}

	d := map[string]interface{}{
		"_key":       module.Id,
		"_from":      from,
		"_to":        to,
		"id":         module.Id,
		"parameter1": module.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func All(ctx context.Context, offset, limit int64, sortParameter, sortDirection string) ([]*organisation_proto.Organisation, error) {
	var orgs []*organisation_proto.Organisation

	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)
	merge_query := queryMerge()

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s`, common.DbOrganisationTable, sort_query, limit_query, merge_query)

	resp, err := runQuery(ctx, q, common.DbOrganisationTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if org, err := recordToOrg(r); err == nil {
			orgs = append(orgs, org)
		}
	}
	return orgs, nil
}

// Creates a organisation
func Create(ctx context.Context, org *organisation_proto.Organisation, modules []*static_proto.Module) error {
	// create organisation
	if len(org.Id) == 0 {
		org.Id = uuid.NewUUID().String()
	}
	if org.Created == 0 {
		org.Created = time.Now().Unix()
	}
	org.Updated = time.Now().Unix()
	if org.Type == 0 {
		org.Type = organisation_proto.OrganisationType_ROOT
	}

	//
	org.Modules = nil
	record, err := orgToRecord(org)
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
		INTO %v`, org.Id, record, record, common.DbOrganisationTable)
	_, err = runQuery(ctx, q, common.DbOrganisationTable)

	//updat moduless
	if err := UpdateModules(ctx, org.Id, modules); err != nil {
		return err
	}

	return nil
}

// Update update organisation
func Update(ctx context.Context, org *organisation_proto.Organisation) (*organisation_proto.Organisation, error) {
	org.Updated = time.Now().Unix()
	record, err := orgToRecord(org, true)
	if err != nil {
		common.InternalServerError(common.OrganisationSrv, Update, nil, "Organisation marshaling is invalid")
		return nil, err
	}
	if len(record) == 0 {
		return nil, errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v IN %v
		RETURN NEW`, org.Id, record, record, common.DbOrganisationTable)
	resp, err := runQuery(ctx, q, common.DbAccountTable)
	if err != nil || len(resp.Records) == 0 {
		common.NotFound(common.AccountSrv, Read, err, "Read query is failed")
		return nil, ErrNotFound
	}
	data, err := recordToOrg(resp.Records[0])
	return data, err
}

// Read reads a organisation by ID
func Read(ctx context.Context, orgId string) (*organisation_proto.Organisation, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, orgId)
	merge_query := queryMerge()
	q := fmt.Sprintf(`
		FOR doc IN %v 
		%s
		%s`, common.DbOrganisationTable, query, merge_query)

	resp, err := runQuery(ctx, q, common.DbOrganisationTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	org, err := recordToOrg(resp.Records[0])
	return org, err
}

// Creates a organisation profile
func CreateOrganisationProfile(ctx context.Context, profile *organisation_proto.OrganisationProfile) error {
	if len(profile.Id) == 0 {
		profile.Id = uuid.NewUUID().String()
	}
	if profile.Created == 0 {
		profile.Created = time.Now().Unix()
	}
	profile.Updated = time.Now().Unix()

	record, err := profileToRecord(profile)
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
		IN %v`, profile.Id, record, record, common.DbOrgProfileTable)
	_, err = runQuery(ctx, q, common.DbOrgProfileTable)
	if err != nil {
		return err
	}
	return nil
}

// Creates a organisation setting
func CreateOrganisationSetting(ctx context.Context, setting *organisation_proto.OrganisationSetting) error {
	if len(setting.Id) == 0 {
		setting.Id = uuid.NewUUID().String()
	}
	if setting.Created == 0 {
		setting.Created = time.Now().Unix()
	}
	setting.Updated = time.Now().Unix()

	record, err := settingToRecord(setting)
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
		IN %v`, setting.Id, record, record, common.DbOrgSettingTable)
	_, err = runQuery(ctx, q, common.DbOrgSettingTable)
	if err != nil {
		return err
	}
	return nil
}

func UpdateModules(ctx context.Context, orgId string, modules []*static_proto.Module) error {

	// delete module edge
	_from := fmt.Sprintf("%v/%v", common.DbOrganisationTable, orgId)
	q_del := fmt.Sprintf(`
		FOR doc IN %v
		FILTER doc._from == "%v"
		REMOVE doc IN %v`, common.DbOrgModuleEdgeTable, _from, common.DbOrgModuleEdgeTable)
	if _, err := runQuery(ctx, q_del, common.DbOrgModuleEdgeTable); err != nil {
		return err
	}

	// create modules edge table
	for _, module := range modules {
		organisation_module := &organisation_proto.OrganisationModule{
			Id:          uuid.NewUUID().String(),
			DisplayName: module.DisplayName,
			OrgId:       orgId,
		}

		_from := fmt.Sprintf("%v/%v", common.DbOrganisationTable, orgId)
		_to := fmt.Sprintf("%v/%v", common.DbModuleTable, module.Id)
		field := fmt.Sprintf(`{_from:"%v",_to:"%v"} `, _from, _to)

		record, err := organisationModuleToRecord(_from, _to, organisation_module)
		if err != nil {
			return err
		}
		q := fmt.Sprintf(`
			UPSERT %v
			INSERT %v
			UPDATE %v
			INTO %v`, field, record, record, common.DbOrgModuleEdgeTable)
		if _, err := runQuery(ctx, q, common.DbOrgModuleEdgeTable); err != nil {
			//log error here
			return err
		}
	}

	//update org with modules
	var arr []interface{}
	for _, item := range modules {
		arr = append(arr, map[string]string{
			"id": item.Id,
		})
	}

	body, err := json.Marshal(arr)
	if err != nil {
		return err
	}
	q := fmt.Sprintf(`
		FOR org IN %v
		UPDATE { _key: "%v" } WITH {data: {modules:%v}}
		IN %v`, common.DbOrganisationTable, orgId, string(body), common.DbOrganisationTable)

	if _, err := runQuery(ctx, q, common.DbOrganisationTable); err != nil {
		return err
	}
	return nil
}

func GetModulesByOrg(ctx context.Context, orgId string) ([]*static_proto.Module, error) {
	var modules []*static_proto.Module
	q := fmt.Sprintf(`
		FOR m, om IN OUTBOUND "%v/%v" %v
		return {data:{
				id:m.id,
				name: m.data.name,
				icon_slug: m.data.icon_slug,
				name_slug: m.data.name_slug,
				description: m.data.description,
				summary: m.data.summary,
				display_name: NOT_NULL(om.data.display_name, m.data.name),
				path: m.data.path
			}}`, common.DbOrganisationTable, orgId, common.DbOrgModuleEdgeTable)

	resp, err := runQuery(ctx, q, common.DbOrgModuleEdgeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if module, err := common.RecordToModule(r); err == nil {
			modules = append(modules, module)
		}
	}
	return modules, nil
}
