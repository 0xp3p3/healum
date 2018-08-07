package handler

import (
	"context"
	"server/audit-srv/db"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

var cl = client.NewClient(
	client.Transport(nats_transport.NewTransport()),
	client.Broker(nats_broker.NewBroker()),
	client.RequestTimeout(4*time.Second),
	client.Retries(5),
)

func initDb() {
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

var audit = &audit_proto.Audit{
	OrgId:            "orgid",
	ActionName:       "name1",
	ActionType:       audit_proto.ActionType_CREATED,
	ActionSourceUser: "source1",
	ActionTargetUser: "target1",
	ActionTimestamp:  123456,
	ActionResource:   "resource1",
	ActionParameters: "parameter1",
	ActionService:    "service1",
	ActionMethod:     "method1",
	ActionMetaData:   "meta1",
}

func initHandler() *AuditService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	hdlr := &AuditService{
		Broker: nats_brker,
	}
	return hdlr
}

func TestCreateAudit(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_create := &audit_proto.CreateAuditRequest{Audit: audit}
	rsp_create := &audit_proto.CreateAuditResponse{}
	if err := hdlr.CreateAudit(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}
}

func TestFilterAudits(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_filter := &audit_proto.FilterAuditsRequest{ActionName: "name1"}
	rsp_filter := &audit_proto.FilterAuditsResponse{}
	if err := hdlr.FilterAudits(ctx, req_filter, rsp_filter); err != nil {
		t.Error(err)
		return
	}

	if len(rsp_filter.Data.Audits) == 0 {
		t.Error("Object count is not matched")
		return
	}

	if rsp_filter.Data.Audits[0].ActionName != "name1" {
		t.Error("Name is not matched")
		return
	}
}
