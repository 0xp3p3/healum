package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"server/common"
	mobpush_proto "server/mob-push-srv/proto/mobpush"
	pubsub_proto "server/static-srv/proto/pubsub"
	user_proto "server/user-srv/proto/user"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/nats"
	log "github.com/sirupsen/logrus"
)

type MobpushService struct {
	Broker     broker.Broker
	PushUrl    string
	UserClient user_proto.UserServiceClient
}

func (p *MobpushService) generateNotifications(ctx context.Context, message *pubsub_proto.PublishBulkNotification) ([]*mobpush_proto.Notification, error) {
	log.Info("Received Mobpush.GenerateNotifications request")

	// make initial notification messages
	notifications := []*mobpush_proto.Notification{
		{
			Tokens:   []string{},
			Platform: 1,
			Message:  message.Notification.Message,
			Alert:    message.Notification.Alert,
			Data:     message.Notification.Data,
			Badge:    int32(len(message.Notification.Data)),
		},
		{
			Tokens:   []string{},
			Platform: 2,
			Message:  message.Notification.Message,
			Alert:    message.Notification.Alert,
			Data:     message.Notification.Data,
			Badge:    int32(len(message.Notification.Data)),
		},
	}
	// read device token with user id array
	req_tokens := user_proto.ReadTokensRequest{message.Notification.UserIds}
	rsp_tokens, err := p.UserClient.ReadTokens(ctx, &req_tokens)
	if err != nil {
		return notifications, err
	}
	for _, t := range rsp_tokens.Tokens {
		if t.Platform > 0 {
			notifications[t.Platform-1].Tokens = append(notifications[t.Platform-1].Tokens, t.DeviceToken)
			notifications[t.Platform-1].Topic = t.AppIdentifier
		}
	}
	return notifications, nil
}

func (p *MobpushService) Subscribe(ctx context.Context, req *pubsub_proto.SubscribeRequest, rsp *pubsub_proto.SubscribeResponse) error {
	_, err := p.Broker.Subscribe(req.Channel, func(pub broker.Publication) error {
		msg := &pubsub_proto.PublishBulkNotification{}
		decoder := json.NewDecoder(bytes.NewReader([]byte(pub.Message().Body)))
		err := decoder.Decode(&msg)
		if err != nil {
			return err
		}

		// creating notifications
		notifications, _ := p.generateNotifications(ctx, msg)
		req_push := &mobpush_proto.PushRequest{
			Notifications: notifications,
		}
		rsp_push := &mobpush_proto.PushResponse{}
		// push notification
		if err := p.Push(ctx, req_push, rsp_push); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (p *MobpushService) Push(ctx context.Context, req *mobpush_proto.PushRequest, rsp *mobpush_proto.PushResponse) error {
	log.Info("Received Mobpush.Push request")

	jsonStr, err := json.Marshal(req)
	if err != nil {
		return err
	}
	req_push, err := http.NewRequest("POST", p.PushUrl, bytes.NewBuffer(jsonStr))
	req_push.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req_push.Header)

	client := &http.Client{}
	rsp_push, err := client.Do(req_push)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(rsp_push.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(body, &rsp)

	return nil
}
