package main

import (
	"context"
	"server/common"
	"server/sms-srv/handler"
	sms_proto "server/sms-srv/proto/sms"
	pubsub_proto "server/static-srv/proto/pubsub"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-os/config"
	"github.com/micro/go-os/config/source/file"
	"github.com/micro/go-os/metrics"
	_ "github.com/micro/go-plugins/broker/nats"
	"github.com/micro/go-plugins/metrics/telegraf"
	_ "github.com/micro/go-plugins/transport/nats"
	log "github.com/sirupsen/logrus"
)

func main() {
	configFile, _ := common.PopParameter("config")
	// Create a config instance
	conf := config.NewConfig(
		// poll every hour
		config.PollInterval(time.Hour),
		// use file as a config source
		config.WithSource(file.NewSource(config.SourceName(configFile))),
	)

	defer conf.Close()

	// create new metrics
	m := telegraf.NewMetrics(
		metrics.Namespace(conf.Get("service", "name").String("sms")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	common.SmsSrv = conf.Get("service", "sms").String("go.micro.srv.sms")
	version := conf.Get("service", "version").String("latest")
	descr := conf.Get("service", "description").String("Micro service")

	twilioToken := conf.Get("service", "TWILIO_AUTH_TOKEN").String("")
	twilioSid := conf.Get("service", "TWILIO_ACCOUNT_SID").String("")
	twilioFrom := conf.Get("service", "TWILIO_FROM").String("")

	service := micro.NewService(
		micro.Name(common.SmsSrv),
		micro.Version(version),
		micro.Metadata(map[string]string{"Description": descr}),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*10),
		micro.Flags(
			cli.BoolFlag{
				Name: "debug",
			},
			cli.StringFlag{
				Name: "twilio_auth_token",
			},
			cli.StringFlag{
				Name: "twilio_account_sid",
			},
			cli.StringFlag{
				Name: "twilio_from",
			},
		),
	)
	var debugEnabled bool
	service.Init(
		micro.Action(func(c *cli.Context) {
			debugEnabled = c.Bool("debug")
			twilioToken = c.String("twilio_auth_token")
			twilioSid = c.String("twilio_account_sid")
			twilioFrom = c.String("twilio_from")
		}),
	)
	if debugEnabled {
		log.SetLevel(log.DebugLevel)
	}

	brker := service.Client().Options().Broker
	brker.Connect()
	smsService := &handler.SmsService{
		Broker:      brker,
		TwilioToken: twilioToken,
		TwilioSid:   twilioSid,
		TwilioFrom:  twilioFrom,
	}
	if err := smsService.Subscribe(context.TODO(), &pubsub_proto.SubscribeRequest{common.SEND_SMS}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	sms_proto.RegisterSmsServiceHandler(service.Server(), smsService)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
