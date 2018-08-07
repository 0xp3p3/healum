package main

import (
	"server/common"
	"server/kv-srv/handler"
	kv_proto "server/kv-srv/proto/kv"
	"time"

	"github.com/go-redis/redis"
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
		metrics.Namespace(conf.Get("service", "name").String("kv")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	common.KvSrv = conf.Get("service", "kv").String("go.micro.srv.kv")
	version := conf.Get("service", "version").String("latest")
	descr := conf.Get("service", "description").String("Micro service")
	service := micro.NewService(
		micro.Name(common.KvSrv),
		micro.Version(version),
		micro.Metadata(map[string]string{"Description": descr}),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*10),
		micro.Flags(
			cli.BoolFlag{
				Name: "debug",
			},
			cli.StringFlag{
				Name: "redis_server_address",
			},
			cli.StringFlag{
				Name: "redis_server_password",
			},
		),
	)
	var debugEnabled bool
	//there needs to be a better way to do this in production maybe by using prod flag?
	server := conf.Get("service", "server").String("127.0.0.1:6379")
	password := conf.Get("service", "password").String("")

	service.Init(
		micro.Action(func(c *cli.Context) {
			debugEnabled = c.Bool("debug")
			server = c.String("redis_server_address")
			password = c.String("redis_server_password")
		}),
	)
	if debugEnabled {
		log.SetLevel(log.DebugLevel)
	}

	kvService := &handler.KvService{
		Client: redis.NewClient(&redis.Options{
			Addr:     server,
			Password: password, // no password set
			DB:       0,        // use default DB
		}),
	}
	kv_proto.RegisterKvServiceHandler(service.Server(), kvService)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
