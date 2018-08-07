package handler

import (
	"context"
	"server/common"
	"server/response-srv/db"
	resp_proto "server/response-srv/proto/response"
	static_proto "server/static-srv/proto/static"

	log "github.com/sirupsen/logrus"
)

type ResponseService struct{}

func (p *ResponseService) Check(ctx context.Context, req *resp_proto.CheckRequest, rsp *resp_proto.CheckResponse) error {
	log.Info("Received Response.All request")
	// will return survey id and authentificationRequired
	survey, err := db.Check(ctx, req.ShortHash, req.OrgId, req.TeamId)
	if survey == nil || err != nil {
		return common.NotFound(common.ResponseSrv, p.Check, err, "response not found")
	}
	rsp.Data = &resp_proto.CheckResponse_Data{survey}
	return nil
}

func (p *ResponseService) All(ctx context.Context, req *resp_proto.AllRequest, rsp *resp_proto.AllResponse) error {
	log.Info("Received Response.All request")
	resps, err := db.All(ctx, req.SurveyId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(resps) == 0 || err != nil {
		return common.NotFound(common.ResponseSrv, p.All, err, "response not found")
	}
	rsp.Data = &resp_proto.ArrData{resps}
	return nil
}

func (p *ResponseService) Create(ctx context.Context, req *resp_proto.CreateRequest, rsp *resp_proto.CreateResponse) error {
	log.Info("Received Response.Create request")
	err := db.Create(ctx, req.SurveyId, req.Response)
	if err != nil {
		return common.InternalServerError(common.ResponseSrv, p.Create, err, "create error")
	}

	// delete repective pending from the pending table
	if err := db.RemovePendingSharedAction(ctx, req.SurveyId); err != nil {
		return common.InternalServerError(common.ResponseSrv, p.Create, err, "remove error")
	}
	// update shareX status
	if err := db.UpdateShareSurveyStatus(ctx, req.SurveyId, static_proto.ShareStatus_VIEWED); err != nil {
		return common.InternalServerError(common.ResponseSrv, p.Create, err, "update error")
	}
	rsp.Data = &resp_proto.Data{req.Response}
	return nil
}

func (p *ResponseService) UpdateState(ctx context.Context, req *resp_proto.UpdateStateRequest, rsp *resp_proto.UpdateStateResponse) error {
	log.Info("Received Response.UpdateState request")
	resp, err := db.UpdateState(ctx, req.SurveyId, req.Response.ResponseId, req.Response.State, req.OrgId, req.TeamId)
	if resp == nil || err != nil {
		return common.NotFound(common.ResponseSrv, p.UpdateState, err, "status not found")
	}
	rsp.ResponseId = resp.Id
	return nil
}

func (p *ResponseService) AllState(ctx context.Context, req *resp_proto.AllStateRequest, rsp *resp_proto.AllStateResponse) error {
	log.Info("Received Response.AllState request")
	resps, err := db.AllState(ctx, req.SurveyId, req.State, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(resps) == 0 || err != nil {
		return common.NotFound(common.ResponseSrv, p.AllState, err, "not found")
	}
	rsp.Data = &resp_proto.ArrData{resps}
	return nil
}

func (p *ResponseService) AllAggQuestion(ctx context.Context, req *resp_proto.AllAggQuestionRequest, rsp *resp_proto.AllAggQuestionResponse) error {
	log.Info("Received Response.AllAggQuestion request")
	resp, err := db.AllAggQuestion(ctx, req.SurveyId, req.OrgId, req.TeamId)
	if resp == nil || err != nil {
		return common.NotFound(common.ResponseSrv, p.AllAggQuestion, err, "not found")
	}
	rsp.Data = resp
	return nil
}

func (p *ResponseService) TimeFilter(ctx context.Context, req *resp_proto.TimeFilterRequest, rsp *resp_proto.TimeFilterResponse) error {
	log.Info("Received Response.TimeFilter request")
	resps, err := db.TimeFilter(ctx, req.SurveyId, req.From, req.To, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(resps) == 0 || err != nil {
		return common.NotFound(common.ResponseSrv, p.TimeFilter, err, "not found")
	}
	rsp.Data = &resp_proto.ArrData{resps}
	return nil
}

func (p *ResponseService) ReadStats(ctx context.Context, req *resp_proto.ReadStatsRequest, rsp *resp_proto.ReadStatsResponse) error {
	log.Info("Received Response.ReadStats request")
	stats, err := db.ReadStats(ctx, req.SurveyId)
	if stats == nil || err != nil {
		return common.NotFound(common.ResponseSrv, p.ReadStats, err, "not found")
	}
	rsp.Data = &resp_proto.ReadStatsResponse_Data{stats}
	return nil
}

func (p *ResponseService) ByUser(ctx context.Context, req *resp_proto.ByUserRequest, rsp *resp_proto.ByUserResponse) error {
	log.Info("Received Response.ByUser request")
	resps, err := db.ByUser(ctx, req.SurveyId, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(resps) == 0 || err != nil {
		return common.NotFound(common.ResponseSrv, p.ByUser, err, "not found")
	}
	rsp.Data = &resp_proto.ArrData{resps}
	return nil
}

func (p *ResponseService) ByAnyUser(ctx context.Context, req *resp_proto.ByAnyUserRequest, rsp *resp_proto.ByAnyUserResponse) error {
	log.Info("Received Response.ByAnyUser request")
	resps, err := db.ByAnyUser(ctx, req.SurveyId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(resps) == 0 || err != nil {
		return common.NotFound(common.ResponseSrv, p.ByAnyUser, err, "not found")
	}
	rsp.Data = &resp_proto.ArrData{resps}
	return nil
}

func (p *ResponseService) Update(ctx context.Context, req *resp_proto.UpdateRequest, rsp *resp_proto.UpdateResponse) error {
	log.Info("Received Response.Update request")
	return nil
}
