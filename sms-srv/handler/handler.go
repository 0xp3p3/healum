package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	sms_proto "server/sms-srv/proto/sms"
	pubsub_proto "server/static-srv/proto/pubsub"

	"github.com/micro/go-micro/broker"
	"github.com/sfreiberg/gotwilio"
	log "github.com/sirupsen/logrus"
	phone_util "github.com/ttacon/libphonenumber"
)

type SmsService struct {
	Broker      broker.Broker
	TwilioToken string
	TwilioSid   string
	TwilioFrom  string
}

func (p *SmsService) Subscribe(ctx context.Context, req *pubsub_proto.SubscribeRequest, rsp *pubsub_proto.SubscribeResponse) error {
	log.Info("Received SmsService.Subscribe request")
	_, err := p.Broker.Subscribe(req.Channel, func(pub broker.Publication) error {
		message := sms_proto.Subscribe{}
		decoder := json.NewDecoder(bytes.NewReader([]byte(pub.Message().Body)))
		err := decoder.Decode(&message)
		if err != nil {
			return err
		}

		req_send := &sms_proto.SendRequest{
			Phone:   message.Phone,
			Message: message.Message,
		}
		rsp_send := &sms_proto.SendResponse{}
		// send message
		if err := p.Send(ctx, req_send, rsp_send); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (p *SmsService) Send(ctx context.Context, req *sms_proto.SendRequest, rsp *sms_proto.SendResponse) error {
	log.Info("Received SmsService.Send request")
	log.Debug("Using SID: ", p.TwilioSid)
	twilio := gotwilio.NewTwilioClient(p.TwilioSid, p.TwilioToken)
	//TODO: add a check for valid phone number before the text is sent?
	num, err := phone_util.Parse(req.Phone, "GB")
	if err != nil {
		log.Error(err)
		return err
	}

	from := p.TwilioFrom
	to := phone_util.Format(num, phone_util.E164)
	message := req.Message
	smsResponse, exc, err := twilio.SendSMS(from, to, message, "", "")
	if exc != nil {
		log.Error(exc.Message)
		return errors.New(exc.Message)
	}
	if err != nil {
		log.Error(err)
		return err
	}
	log.Debug(smsResponse)
	return nil
}
