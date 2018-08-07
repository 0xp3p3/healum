package main

import (
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	"server/team-srv/db"
	"server/team-srv/handler"
	team_proto "server/team-srv/proto/team"
	user_proto "server/user-srv/proto/user"
	organisation_proto "server/organisation-srv/proto/organisation"
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
		metrics.Namespace(conf.Get("service", "name").String("team")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	common.TeamSrv = conf.Get("service", "team").String("go.micro.srv.team")
	version := conf.Get("service", "version").String("latest")
	descr := conf.Get("service", "description").String("Micro service")
	service := micro.NewService(
		micro.Name(common.TeamSrv),
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

	teamSevice := &handler.TeamService{
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", service.Client()),
		UserClient:    user_proto.NewUserServiceClient("go.micro.srv.user", service.Client()),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", service.Client()),
		OrganisationClient: organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", service.Client()),
	}
	team_proto.RegisterTeamServiceHandler(service.Server(), teamSevice)
	db.Init(service.Client())

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
