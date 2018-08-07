package main

import (
	"context"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	"server/plan-srv/db"
	"server/plan-srv/handler"
	plan_proto "server/plan-srv/proto/plan"
	team_proto "server/team-srv/proto/team"
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
		metrics.Namespace(conf.Get("service", "name").String("plan")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	common.PlanSrv = conf.Get("service", "plan").String("go.micro.srv.plan")
	version := conf.Get("service", "version").String("latest")
	descr := conf.Get("service", "description").String("Micro service")
	service := micro.NewService(
		micro.Name(common.PlanSrv),
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

	planService := &handler.PlanService{
		Broker:        service.Client().Options().Broker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", service.Client()),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", service.Client()),
		TeamClient:    team_proto.NewTeamServiceClient("go.micro.srv.team", service.Client()),
	}
	plan_proto.RegisterPlanServiceHandler(service.Server(), planService)
	db.Init(service.Client())

	go planService.WarmupCache(context.TODO(), &plan_proto.WarmupCacheRequest{}, &plan_proto.WarmupCacheResponse{})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
