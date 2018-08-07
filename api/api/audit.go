package api

import (
	"context"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type AuditService struct {
	AuditClient   audit_proto.AuditServiceClient
	Auth          Filters
	ServerMetrics metrics.Metrics
}

func (p AuditService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/audits")

	ws.Route(ws.POST("/read").To(p.FilterAudits).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		Doc("Logout the user"))

	restful.Add(ws)
}

func (p *AuditService) FilterAudits(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Audit.FilterAudits API request")

	req_audit := new(audit_proto.FilterAuditsRequest)
	req_audit.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_audit.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_audit.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_audit.SortParameter = req.Attribute(SortParameter).(string)
	req_audit.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.AuditClient.FilterAudits(ctx, req_audit)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.audit.FilterAudits", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read all audits successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}
