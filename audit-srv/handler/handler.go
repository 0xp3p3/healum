package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"server/audit-srv/db"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	pubsub_proto "server/static-srv/proto/pubsub"

	"github.com/micro/go-micro/broker"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type AuditService struct {
	Broker broker.Broker
}

// Subscribe get subscription from other services to create audit log
func (p *AuditService) Subscribe(ctx context.Context, req *pubsub_proto.SubscribeRequest, rsp *pubsub_proto.SubscribeResponse) error {
	_, err := p.Broker.Subscribe(req.Channel, func(pub broker.Publication) error {
		msg := &audit_proto.Audit{}
		decoder := json.NewDecoder(bytes.NewReader([]byte(pub.Message().Body)))
		err := decoder.Decode(&msg)
		if err != nil {
			return err
		}

		req_audit := &audit_proto.CreateAuditRequest{Audit: msg}
		rsp_audit := &audit_proto.CreateAuditResponse{}
		// create audit
		if err := p.CreateAudit(ctx, req_audit, rsp_audit); err != nil {
			return err
		}
		return nil
	})
	return err
}

// AllAudits reads all audits
func (p *AuditService) AllAudits(ctx context.Context, req *audit_proto.AllAuditsRequest, rsp *audit_proto.AllAuditsResponse) error {
	log.Info("Received Audit.AllAudits request")

	audits, err := db.AllAudits(ctx, req.OrgId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(audits) == 0 || err != nil {
		return common.NotFound(common.AuditSrv, p.AllAudits, err, "audit not found")
	}
	rsp.Data = &audit_proto.AuditArrData{audits}
	return nil
}

// CreateAudit creates audits from subsciption request
func (p *AuditService) CreateAudit(ctx context.Context, req *audit_proto.CreateAuditRequest, rsp *audit_proto.CreateAuditResponse) error {
	log.Info("Received Audit.CreateAudit request")
	if len(req.Audit.Id) == 0 {
		req.Audit.Id = uuid.NewUUID().String()
	}
	// create audit
	err := db.CreateAudit(ctx, req.Audit)
	if err != nil {
		return common.InternalServerError(common.AuditSrv, p.CreateAudit, err, "server error")
	}
	rsp.Data = &audit_proto.AuditData{req.Audit}
	return nil
}

// FilterAudits filters audits
func (p *AuditService) FilterAudits(ctx context.Context, req *audit_proto.FilterAuditsRequest, rsp *audit_proto.FilterAuditsResponse) error {
	log.Info("Received Audit.FilterAudits request")

	audits, err := db.FilterAudits(ctx, req)
	if len(audits) == 0 || err != nil {
		return common.NotFound(common.AuditSrv, p.AllAudits, err, "audit not found")
	}
	rsp.Data = &audit_proto.AuditArrData{audits}
	return nil
}
