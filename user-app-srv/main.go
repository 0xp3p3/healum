package main

import (
	account_proto "server/account-srv/proto/account"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	content_proto "server/content-srv/proto/content"
	kv_proto "server/kv-srv/proto/kv"
	plan_proto "server/plan-srv/proto/plan"
	static_proto "server/static-srv/proto/static"
	survey_proto "server/survey-srv/proto/survey"
	track_proto "server/track-srv/proto/track"
	"server/user-app-srv/db"
	"server/user-app-srv/handler"
	userapp_proto "server/user-app-srv/proto/userapp"
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
		metrics.Namespace(conf.Get("service", "name").String("userapp")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	common.UserappSrv = conf.Get("service", "userapp").String("go.micro.srv.userapp")
	version := conf.Get("service", "version").String("latest")
	descr := conf.Get("service", "description").String("Micro service")
	service := micro.NewService(
		micro.Name(common.UserappSrv),
		micro.Version(version),
		micro.Metadata(map[string]string{"Description": descr}),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*10),
		micro.Flags(
			cli.BoolFlag{
				Name: "debug",
			},
		),
	)
	var debugEnabled bool
	service.Init(
		micro.Action(func(c *cli.Context) {
			debugEnabled = c.Bool("debug")
		}),
	)
	if debugEnabled {
		log.SetLevel(log.DebugLevel)
	}

	brker := service.Client().Options().Broker
	brker.Connect()
	userappService := &handler.UserAppService{
		Broker:          brker,
		KvClient:        kv_proto.NewKvServiceClient("go.micro.srv.kv", service.Client()),
		ContentClient:   content_proto.NewContentServiceClient("go.micro.srv.content", service.Client()),
		BehaviourClient: behaviour_proto.NewBehaviourServiceClient("go.micro.srv.behaviour", service.Client()),
		UserClient:      user_proto.NewUserServiceClient("go.micro.srv.user", service.Client()),
		TrackClient:     track_proto.NewTrackServiceClient("go.micro.srv.track", service.Client()),
		PlanClient:      plan_proto.NewPlanServiceClient("go.micro.srv.plan", service.Client()),
		SurveyClient:    survey_proto.NewSurveyServiceClient("go.micro.srv.survey", service.Client()),
		AccountClient:   account_proto.NewAccountServiceClient("go.micro.srv.account", service.Client()),
		StaticClient:    static_proto.NewStaticServiceClient("go.micro.srv.static", service.Client()),
	}
	userapp_proto.RegisterUserAppServiceHandler(service.Server(), userappService)
	db.Init(service.Client())

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
