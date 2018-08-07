package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"server/common"
	sms_proto "server/sms-srv/proto/sms"
	pubsub_proto "server/static-srv/proto/pubsub"
	"testing"
	"time"

	"github.com/micro/go-micro/broker"
	nats_broker "github.com/micro/go-plugins/broker/nats"
)

func TestSubscribe(t *testing.T) {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()

	ctx := common.NewTestContext(context.TODO())
	hdlr := &SmsService{
		Broker:      nats_brker,
		TwilioToken: "f6b2a085c65d3c3981879abfb3219692",
		TwilioSid:   "ACaa3f9590a634f3276c392a34471cf975",
		TwilioFrom:  "+13152828170",
	}
	// subscribe
	req_subscribe := &pubsub_proto.SubscribeRequest{
		Channel: common.SEND_SMS,
	}
	rsp_subscribe := &pubsub_proto.SubscribeResponse{}
	if err := hdlr.Subscribe(ctx, req_subscribe, rsp_subscribe); err != nil {
		t.Error(err)
		return
	}
	// make the sms message
	// in the future, subscribe will be included user_id instead of phone
	message := &sms_proto.Subscribe{
		Phone:   "+8613042431402",
		Message: fmt.Sprintf(common.MSG_PASSCOD_SMS, "123456"),
	}
	body, err := json.Marshal(message)
	if err != nil {
		t.Error(err)
		return
	}

	// publish sms
	if err := nats_brker.Publish(common.SEND_SMS, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)
}

func TestSend(t *testing.T) {
	ctx := common.NewTestContext(context.TODO())
	hdlr := &SmsService{
		TwilioToken: "f6b2a085c65d3c3981879abfb3219692",
		TwilioSid:   "ACaa3f9590a634f3276c392a34471cf975",
		TwilioFrom:  "+13152828170",
	}

	req := &sms_proto.SendRequest{
		Phone:   "+8613042431402",
		Message: "Hello world!",
	}
	rsp := &sms_proto.SendResponse{}
	err := hdlr.Send(ctx, req, rsp)
	if err != nil {
		t.Error(err)
	}
}
