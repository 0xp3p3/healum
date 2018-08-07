package main

import (
	"context"
	"server/common"
	"server/mob-push-srv/handler"
	mobpush_proto "server/mob-push-srv/proto/mobpush"
	pubsub_proto "server/static-srv/proto/pubsub"
	user_proto "server/user-srv/proto/user"
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
		metrics.Namespace(conf.Get("service", "name").String("mobpush")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	common.MobpushSrv = conf.Get("service", "mobpush").String("go.micro.srv.mobpush")
	version := conf.Get("service", "version").String("latest")
	descr := conf.Get("service", "description").String("Micro service")
	gorush_addr := conf.Get("service", "gorush_address").String("localhost:8088")

	service := micro.NewService(
		micro.Name(common.MobpushSrv),
		micro.Version(version),
		micro.Metadata(map[string]string{"Description": descr}),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*10),
		micro.Flags(
			cli.BoolFlag{
				Name: "debug",
			},
			cli.StringFlag{
				Name: "gorush_server_address",
			},
		),
	)
	var debugEnabled bool
	service.Init(
		micro.Action(func(c *cli.Context) {
			debugEnabled = c.Bool("debug")
			gorush_addr = c.String("gorush_server_address")
		}),
	)
	if debugEnabled {
		log.SetLevel(log.DebugLevel)
	}

	brker := service.Client().Options().Broker
	brker.Connect()
	mobpushService := &handler.MobpushService{
		Broker:     brker,
		PushUrl:    "http://" + gorush_addr + "/api/push",
		UserClient: user_proto.NewUserServiceClient("go.micro.srv.user", service.Client()),
	}
	mobpush_proto.RegisterMobpushServiceHandler(service.Server(), mobpushService)
	// subscribe in the first
	if err := mobpushService.Subscribe(context.TODO(), &pubsub_proto.SubscribeRequest{common.SEND_NOTIFICATION}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
