package handler

import (
	"context"
	"encoding/json"
	"server/activity-srv/db"
	activity_proto "server/activity-srv/proto/activity"
	"server/common"

	"github.com/micro/go-micro/broker"
	log "github.com/sirupsen/logrus"
)

type ActivityService struct {
	Broker broker.Broker
}

func (p *ActivityService) Create(ctx context.Context, req *activity_proto.CreateRequest, rsp *activity_proto.CreateResponse) error {
	log.Info("Received Activity.Create request")

	body, err := json.Marshal(req.Activity)
	if err != nil {
		return common.InternalServerError(common.ActivitySrv, p.Create, err, "parsing error")
	}
	if err := p.Broker.Publish(common.CREATE_ACTIVITY_CONTENT, &broker.Message{Body: body}); err != nil {
		return common.InternalServerError(common.ActivitySrv, p.Create, err, "subscribe error")
	}
	return nil
}

func (p *ActivityService) ListConfig(ctx context.Context, req *activity_proto.ListConfigRequest, rsp *activity_proto.ListConfigResponse) error {
	configs, err := db.ListConfig(ctx)
	if len(configs) == 0 || err != nil {
		return common.NotFound(common.ActivitySrv, p.ListConfig, err, "config not found")
	}
	rsp.Configs = configs
	return nil
}

func (p *ActivityService) CreateConfig(ctx context.Context, req *activity_proto.CreateConfigRequest, rsp *activity_proto.CreateConfigResponse) error {
	err := db.CreateConfig(ctx, req.Config)
	if err != nil {
		return common.InternalServerError(common.ActivitySrv, p.CreateConfig, err, "server error")
	}
	rsp.Config = req.Config
	return nil
}

func (p *ActivityService) ReadConfig(ctx context.Context, req *activity_proto.ReadConfigRequest, rsp *activity_proto.ReadConfigResponse) error {
	config, err := db.ReadConfig(ctx, req.Id)
	if config == nil || err != nil {
		return common.InternalServerError(common.ActivitySrv, p.ReadConfig, err, "server error")
	}
	rsp.Config = config
	return nil
}

func (p *ActivityService) DeleteConfig(ctx context.Context, req *activity_proto.DeleteConfigRequest, rsp *activity_proto.DeleteConfigResponse) error {
	if err := db.DeleteConfig(ctx, req.Id); err != nil {
		return common.InternalServerError(common.ActivitySrv, p.DeleteConfig, err, "server error")
	}
	return nil
}
