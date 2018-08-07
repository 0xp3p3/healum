package backend

import (
	"context"
	"server/activity-srv/handler"
	activity_proto "server/activity-srv/proto/activity"
	"server/common"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

func initDb() {
	Init()
	InitExternal(&ExtBackend{
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
}

func TestFetchDatabase(t *testing.T) {
	cl := client.NewClient(
		client.Transport(nats_transport.NewTransport()),
		client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(3*time.Second),
		client.Retries(10),
	)

	initDb()
	hdlr := new(handler.ActivityService)
	hdlr.Init(cl)

	b := Backends["better"]
	datas, err := b.Query("better")
	if err != nil {
		t.Error(err)
		return
	}

	ctx := common.NewTestContext(context.TODO())
	for _, data := range datas {
		rsp := &activity_proto.CreateResponse{}
		err := hdlr.Create(ctx, &activity_proto.CreateRequest{data}, rsp)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
