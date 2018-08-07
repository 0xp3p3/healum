package handler

import (
	"context"
	"encoding/json"
	"server/common"
	mobpush_proto "server/mob-push-srv/proto/mobpush"
	pubsub_proto "server/static-srv/proto/pubsub"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

const (
	ios     = 1
	android = 2
)

var cl = client.NewClient(
	client.Transport(nats_transport.NewTransport()),
	client.Broker(nats_broker.NewBroker()),
	client.RequestTimeout(4*time.Second),
	client.Retries(5),
)

func TestPublishMessage(t *testing.T) {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()

	ctx := common.NewTestContext(context.TODO())
	hdlr := &MobpushService{
		Broker:     nats_brker,
		PushUrl:    "http://gorush:8088/api/push",
		UserClient: user_proto.NewUserServiceClient("go.micro.srv.User", cl),
	}

	// subscribe
	req_subscribe := &pubsub_proto.SubscribeRequest{
		Channel: common.SEND_NOTIFICATION,
	}
	rsp_subscribe := &pubsub_proto.SubscribeResponse{}
	if err := hdlr.Subscribe(ctx, req_subscribe, rsp_subscribe); err != nil {
		t.Error(err)
		return
	}
	// make the notifications
	// in the future, subscribe will be included user_id instead of tokens and platform
	msg := &pubsub_proto.PublishBulkNotification{
		Notification: &pubsub_proto.BulkNotification{
			UserIds: []string{"userid"},
			Message: "Test with pub-sub from mob-push-srv",
		},
	}
	body, err := json.Marshal(msg)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := nats_brker.Publish(common.SEND_NOTIFICATION, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)
}

func TestPushNotification(t *testing.T) {
	ctx := common.NewTestContext(context.TODO())
	hdlr := &MobpushService{
		PushUrl: "http://gorush:8088/api/push",
	}
	// alert := &mobpush_proto.Alert{
	// 	Title:        "Hello World!",
	// 	Body:         "This is a test message",
	// 	ActionLocKey: "Open",
	// 	Action:       "Open",
	// }
	req_push := &mobpush_proto.PushRequest{
		[]*mobpush_proto.Notification{
			{
				Tokens:   []string{"45dd667b84361a5e83e3c048e20df8a1337dccbecb9206843f38698ca074fe3f"},
				Platform: 1,
				Message:  "Hello World again!",
				Topic:    "com.healum.prevo.Prevo",
			},
		},
	}
	rsp_push := &mobpush_proto.PushResponse{}
	err := hdlr.Push(ctx, req_push, rsp_push)
	if err != nil {
		t.Error(err)
		return
	}
	if !rsp_push.Success {
		t.Error("push failed")
		t.Error(rsp_push)
		return
	}
}

func TestGenerateNotifications(t *testing.T) {
	ctx := common.NewTestContext(context.TODO())
	hdlr := &MobpushService{
		PushUrl:    "http://gorush:8088/api/push",
		UserClient: user_proto.NewUserServiceClient("go.micro.srv.user", cl),
	}

	alert := &pubsub_proto.Alert{
		Title:        "Hello World!",
		Body:         "This is a test message",
		ActionLocKey: "Open",
		Action:       "Open",
	}

	_, err := hdlr.generateNotifications(ctx, []string{"81662f84-89dc-11e8-801c-00155d4b0101", "d98ea853-8b92-11e8-97d2-00155d4b0101"}, "test message", alert)
	if err != nil {
		t.Error(err)
		return
	}
}
