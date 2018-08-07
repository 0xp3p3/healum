package main

import (
	account_proto "server/account-srv/proto/account"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
	track_proto "server/track-srv/proto/track"
	"server/user-srv/db"
	"server/user-srv/handler"
	user_proto "server/user-srv/proto/user"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-os/config"
	"github.com/micro/go-os/config/source/file"
	"github.com/micro/go-os/config/source/os"
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
		metrics.Namespace(conf.Get("service", "name").String("user")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	// handler.RecaptchPublicKey = conf.Get("recaptcha", "public").String("")
	// handler.RecaptchPrivateKey = conf.Get("recaptcha", "private").String("")

	common.UserSrv = conf.Get("service", "name").String("go.micro.srv.user")
	version := conf.Get("service", "version").String("latest")
	descr := conf.Get("service", "description").String("Micro service")
	service := micro.NewService(
		micro.Name(common.UserSrv),
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
	userService := &handler.UserService{
		Broker:          brker,
		KvClient:        kv_proto.NewKvServiceClient("go.micro.srv.kv", service.Client()),
		AccountClient:   account_proto.NewAccountServiceClient("go.micro.srv.account", service.Client()),
		TrackClient:     track_proto.NewTrackServiceClient("go.micro.srv.track", service.Client()),
		TeamClient:      team_proto.NewTeamServiceClient("go.micro.srv.team", service.Client()),
		StaticClient:    static_proto.NewStaticServiceClient("go.micro.srv.static", service.Client()),
		BehaviourClient: behaviour_proto.NewBehaviourServiceClient("go.micro.srv.behaviour", service.Client()),
	}
	user_proto.RegisterUserServiceHandler(service.Server(), userService)
	db.Init(service.Client())

	configDynamic := config.NewConfig(
		// poll every hour
		config.PollInterval(time.Hour),
		// use config-srv as a config source
		config.WithSource(os.NewSource(config.SourceName("go.micro.srv.config"),
			config.SourceClient(service.Client()))),
	)

	defer configDynamic.Close()

	// common.ConfigValueWatcher(configDynamic, func(v config.Value) {
	// 	handler.RecaptchPublicKey = v.String("")
	// }, "recaptcha", "public")

	// common.ConfigValueWatcher(configDynamic, func(v config.Value) {
	// 	handler.RecaptchPrivateKey = v.String("")
	// }, "recaptcha", "private")

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

}
