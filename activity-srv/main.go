package main

import (
	"context"
	"server/activity-srv/backend"
	"server/activity-srv/db"
	"server/activity-srv/handler"
	activity_proto "server/activity-srv/proto/activity"
	"server/common"
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
		metrics.Namespace(conf.Get("service", "name").String("activity")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	common.ActivitySrv = conf.Get("service", "name").String("go.mcro.srv.activity")
	version := conf.Get("service", "version").String("latest")
	descr := conf.Get("service", "description").String("Micro service")

	service := micro.NewService(
		micro.Name(common.ActivitySrv),
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
	activityService := &handler.ActivityService{
		Broker: brker,
	}
	activity_proto.RegisterActivityServiceHandler(service.Server(), activityService)

	backend.Init()
	//we need to initialise db-srv here for activity config
	db.Init(service.Client())
	initExtBackend(conf)
	// fetch all database from external apis
	go backend.FetchDatabase(activityService)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func initExtBackend(conf config.Config) {
	ctx := context.TODO()
	// read configs from the collection in the first
	configs, _ := db.ListConfig(ctx)

	if len(configs) == 0 {
		// init ext backends with configs
		backends := []*backend.ExtBackend{}
		if err := conf.Get("external").Scan(&backends); err != nil {
			return
		}
		// init every external apis
		for _, b := range backends {
			config := &activity_proto.Config{
				Name:         b.Name,
				AppURL:       b.AppURL,
				Next:         b.Next,
				Weight:       b.Weight,
				TimeInterval: b.TimeInterval,
				Enabled:      b.Enabled,
				Transform:    b.Transform,
			}
			db.CreateConfig(ctx, config)
		}
		// get all configs from activity-config collection again after storing with id and created datetime
		configs, _ = db.ListConfig(ctx)
	}

	for _, c := range configs {
		b := &backend.ExtBackend{
			ID:           c.Id,
			Name:         c.Name,
			AppURL:       c.AppURL,
			Next:         c.Next,
			Weight:       c.Weight,
			TimeInterval: c.TimeInterval,
			Enabled:      c.Enabled,
			Transform:    c.Transform,
		}
		// init backend with it's config information
		backend.InitExternal(b)
	}
}
