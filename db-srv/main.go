package main

import (
	"context"
	"log"
	"server/common"
	"server/db-srv/db"
	_ "server/db-srv/db/arangodb"
	_ "server/db-srv/db/elastic"
	_ "server/db-srv/db/influxdb"
	_ "server/db-srv/db/mysql"
	_ "server/db-srv/db/redis"
	"server/db-srv/handler"
	proto "server/db-srv/proto/db"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/selector"
	"github.com/micro/go-os/config"
	"github.com/micro/go-os/config/source/file"
	"github.com/micro/go-os/metrics"
	_ "github.com/micro/go-plugins/broker/nats"
	"github.com/micro/go-plugins/metrics/telegraf"
	_ "github.com/micro/go-plugins/transport/nats"
)

func main() {
	configFile, _ := common.PopParameter("config")
	// Create a config instance
	conf := config.NewConfig(
		// poll every hour
		config.PollInterval(time.Hour),
		// use file as an initial config source
		config.WithSource(file.NewSource(config.SourceName(configFile))),
	)

	defer conf.Close()

	// create new metrics
	m := telegraf.NewMetrics(
		metrics.Namespace(conf.Get("service", "name").String("db")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	common.DbSrv = conf.Get("service", "name").String("go.micro.srv.db")
	version := conf.Get("service", "version").String("latest")
	descr := conf.Get("service", "description").String("Micro service")

	service := micro.NewService(
		micro.Name(common.DbSrv),
		micro.Version(version),
		micro.Metadata(map[string]string{"Description": descr}),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*10),

		micro.Flags(
			cli.StringFlag{
				Name:   "database_service_namespace",
				EnvVar: "DATABASE_SERVICE_NAMESPACE",
				Usage:  "The namespace used when looking up databases in registry e.g go.micro.db",
			},
		),

		micro.Action(func(c *cli.Context) {
			if len(c.String("database_service_namespace")) > 0 {
				db.DBServiceNamespace = c.String("database_service_namespace")
			}
		}),
	)

	service.Init(
		// init the db
		micro.BeforeStart(func() error {
			sel := service.Client().Options().Selector
			sel.Init(selector.SetStrategy(selector.RoundRobin))
			return db.Init(sel)
		}),
	)

	db.DefaultDB = db.NewDB(service.Client().Options().Selector)
	dbService := &handler.DB{}
	// init db
	if err := dbService.InitDb(context.TODO(), &proto.InitDbRequest{}, &proto.InitDbResponse{}); err != nil {
		log.Fatal(err)
	}
	proto.RegisterDBHandler(service.Server(), handler.NewWrapper(dbService))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
