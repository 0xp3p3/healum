package main

import (
	"context"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	"server/organisation-srv/db"
	"server/organisation-srv/handler"
	organisation_proto "server/organisation-srv/proto/organisation"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
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
		metrics.Namespace(conf.Get("service", "name").String("organisation")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	common.OrganisationSrv = conf.Get("service", "organisation").String("go.micro.srv.organisation")
	version := conf.Get("service", "version").String("latest")
	descr := conf.Get("service", "description").String("Micro service")
	service := micro.NewService(
		micro.Name(common.OrganisationSrv),
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
	orgService := &handler.OrganisationService{
		Broker:        brker,
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", service.Client()),
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", service.Client()),
		TeamClient:    team_proto.NewTeamServiceClient("go.micro.srv.team", service.Client()),
		StaticClient:  static_proto.NewStaticServiceClient("go.micro.srv.static", service.Client()),
		UserClient:    user_proto.NewUserServiceClient("go.micro.srv.user", service.Client()),
	}
	organisation_proto.RegisterOrganisationServiceHandler(service.Server(), orgService)
	db.Init(service.Client())

	// warmup cache
	go orgService.WarmupCacheOrganisation(context.TODO(), &organisation_proto.WarmupCacheOrganisationRequest{}, &organisation_proto.WarmupCacheOrganisationResponse{})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
