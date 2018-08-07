package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	db_proto "server/db-srv/proto/db"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/client"
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

func auditToRecord(audit *audit_proto.Audit) (string, error) {
	data, err := common.MarhalToObject(audit)
	if err != nil {
		return "", err
	}

	d := map[string]interface{}{
		"_key":       audit.Id,
		"id":         audit.Id,
		"created":    audit.Created,
		"name":       audit.ActionName,
		"parameter1": audit.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToAudit(r *db_proto.Record) (*audit_proto.Audit, error) {
	var p audit_proto.Audit
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func AllAudits(ctx context.Context, orgId string, offset, limit int64, sortParameter, sortDirection string) ([]*audit_proto.Audit, error) {
	var audits []*audit_proto.Audit
	query := common.QueryAuth(`FILTER`, orgId, "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		RETURN doc`, common.DbAuditTable, query, sort_query, limit_query,
	)

	resp, err := runQuery(ctx, q, common.DbAuditTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if audit, err := recordToAudit(r); err == nil {
			audits = append(audits, audit)
		}
	}
	return audits, nil
}

func CreateAudit(ctx context.Context, audit *audit_proto.Audit) error {
	if audit.Created == 0 {
		audit.Created = time.Now().Unix()
	}
	record, err := auditToRecord(audit)
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
		IN %v`, audit.Id, record, record, common.DbAuditTable)
	_, err = runQuery(ctx, q, common.DbAuditTable)
	return err
}

func FilterAudits(ctx context.Context, req *audit_proto.FilterAuditsRequest) ([]*audit_proto.Audit, error) {
	var audits []*audit_proto.Audit

	query := `FILTER`
	if req.ActionType != audit_proto.ActionType_ActionType_NONE {
		query += fmt.Sprintf(` && doc.data.action_type == %d`, req.ActionType)
	}
	if len(req.ActionName) != 0 {
		query += fmt.Sprintf(` && doc.data.action_name == "%s"`, req.ActionName)
	}
	if len(req.ActionSourceUser) != 0 {
		query += fmt.Sprintf(` && doc.data.action_source_user == "%s"`, req.ActionSourceUser)
	}
	if len(req.ActionTargetUser) != 0 {
		query += fmt.Sprintf(` && doc.data.action_target_user == "%s"`, req.ActionTargetUser)
	}
	if req.ActionTimestamp > 0 {
		query += fmt.Sprintf(` && doc.data.action_timestamp == %d`, req.ActionTimestamp)
	}
	if len(req.ActionResource) != 0 {
		query += fmt.Sprintf(` && doc.data.action_resource == "%s"`, req.ActionResource)
	}
	if len(req.ActionService) != 0 {
		query += fmt.Sprintf(` && doc.data.action_service == "%s"`, req.ActionService)
	}
	if len(req.ActionMethod) != 0 {
		query += fmt.Sprintf(` && doc.data.action_method == "%s"`, req.ActionMethod)
	}
	if len(req.ActionMetaData) != 0 {
		query += fmt.Sprintf(` && doc.data.action_meta_data == "%s"`, req.ActionMetaData)
	}
	query = common.QueryAuth(query, req.OrgId, "")

	limit_query := common.QueryPaginate(req.Offset, req.Limit)
	sort_query := common.QuerySort(req.SortParameter, req.SortDirection)
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		RETURN doc`, common.DbAuditTable, query, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbAuditTable)
	if err != nil {
		common.ErrorLog(common.AuditSrv, FilterAudits, err, "RunQuery is failed")
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if audit, err := recordToAudit(r); err == nil {
			audits = append(audits, audit)
		}
	}
	return audits, nil
}
