package handler

import (
	"context"
	"server/activity-srv/db"
	activity_proto "server/activity-srv/proto/activity"
	"server/common"
	content_proto "server/content-srv/proto/content"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

var cl = client.NewClient(
	client.Transport(nats_transport.NewTransport()),
	client.Broker(nats_broker.NewBroker()),
	client.RequestTimeout(3*time.Second),
	client.Retries(10),
)

func initDb() {
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func initHandler() *ActivityService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	return &ActivityService{
		Broker: nats_brker,
	}
}

func TestRemoveDb(t *testing.T) {

}

func TestActivityIsCreate(t *testing.T) {
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	req_create := &activity_proto.CreateRequest{
		Activity: &content_proto.Activity{
			Identifier: "identifier",
		},
	}
	rsp_create := &activity_proto.CreateResponse{}

	err := hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestConfigIsCreated(t *testing.T) {
	initDb()

	ctx := common.NewTestContext(context.TODO())

	err := db.CreateConfig(ctx, &activity_proto.Config{
		Name:         "better",
		AppURL:       "https://www.better.org.uk",
		Next:         "/odi/sessions.json?afterTimestamp=2015-08-05T11:13:09&afterId=5231587",
		Weight:       100,
		TimeInterval: 1,
		Enabled:      true,
		Transform: []*activity_proto.Transform{
			{"geo", "location/containedInPlace/geo", "location/geo"},
			{"object", "activity", "activity"},
			{"object", "offer", "offer"},
		},
	})

	if err != nil {
		t.Error(err)
		return
	}
}

func TestConfigRead(t *testing.T) {
	initDb()

	hdlr := new(ActivityService)
	ctx := common.NewTestContext(context.TODO())
	req := &activity_proto.CreateConfigRequest{
		&activity_proto.Config{
			Id:           "222",
			Name:         "better",
			AppURL:       "https://www.better.org.uk",
			Next:         "",
			Weight:       100,
			TimeInterval: 1,
			Enabled:      true,
		},
	}

	resp := &activity_proto.CreateConfigResponse{}
	err := hdlr.CreateConfig(ctx, req, resp)
	if err != nil {
		t.Error(err)
	}

	req_read := &activity_proto.ReadConfigRequest{
		"222",
	}
	resp_read := &activity_proto.ReadConfigResponse{}
	err = hdlr.ReadConfig(ctx, req_read, resp_read)
	if err != nil {
		t.Error(err)
		return
	}
	if resp_read.Config.Id != "222" {
		t.Error("Id does not match")
		return
	}
}

func TestConfigDelete(t *testing.T) {
	initDb()

	hdlr := new(ActivityService)
	ctx := common.NewTestContext(context.TODO())

	// create
	req_create := &activity_proto.CreateConfigRequest{
		&activity_proto.Config{
			Id:           "222",
			Name:         "better",
			AppURL:       "https://www.better.org.uk",
			Next:         "",
			Weight:       100,
			TimeInterval: 1,
			Enabled:      true,
		},
	}
	resp_create := &activity_proto.CreateConfigResponse{}
	err := hdlr.CreateConfig(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return
	}

	// delete
	req_delete := &activity_proto.DeleteConfigRequest{Id: "222"}
	resp_delete := &activity_proto.DeleteConfigResponse{}
	err = hdlr.DeleteConfig(ctx, req_delete, resp_delete)
	if err != nil {
		t.Error(err)
		return
	}
	req_read := &activity_proto.ReadConfigRequest{Id: "222"}
	resp_read := &activity_proto.ReadConfigResponse{}

	// read to check
	err = hdlr.ReadConfig(ctx, req_read, resp_read)
	if err != nil {
		t.Error(err)
		return
	}
	if resp_read.Config != nil {
		t.Error("Not deleted")
		return
	}
}
