package main

import (
	"github.com/micro/go-micro"
	"server/cloudkey-srv/db"
	"server/cloudkey-srv/storage"
	"server/cloudkey-srv/handler"
	cloudkey_proto "server/cloudkey-srv/proto/record"
	"time"

	"github.com/micro/go-os/config"
	"github.com/micro/go-os/config/source/file"
	"github.com/micro/go-os/config/source/os"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/nats"
	"log"
	"server/common"
	"server/cloudkey-srv/storage/gce"
	"github.com/micro/go-os/metrics"
	"github.com/micro/go-plugins/metrics/telegraf"
)

func main() {
	configcloudkey, _ := common.PopParameter("config")
	// Create a config instance
	conf := config.NewConfig(
		// poll every hour
		config.PollInterval(time.Hour),
		// use file as a config source
		config.WithSource(file.NewSource(config.SourceName(configcloudkey))),
	)

	defer conf.Close()

	// create new metrics
	m := telegraf.NewMetrics(
		metrics.Namespace(conf.Get("service", "name").String("file")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	name := conf.Get("service", "name").String("healum.srv.cloudkey")
	version := conf.Get("service", "version").String("latest")
	descr := conf.Get("service", "description").String("Micro service")

	gce.ProjectID = conf.Get("gce", "ProjectID").String("atlassian-1020")
	gce.BucketName = conf.Get("gce", "BucketName").String("go-project-test")

	service := micro.NewService(
		micro.Name(name),
		micro.Version(version),
		micro.Metadata(map[string]string{"Description": descr}),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second * 10),
	)
	service.Init(
		// init the storage
		micro.BeforeStart(func() error {
			return storage.Init()
		}))
	cloudkey_proto.RegisterCloudKeyServiceHandler(service.Server(), new(handler.CloudKeyService))
	db.Init(service.Client())

	configDynamic := config.NewConfig(
		// poll every hour
		config.PollInterval(time.Hour),
		// use config-srv as a config source
		config.WithSource(os.NewSource(config.SourceName("healum.srv.config"),
			config.SourceClient(service.Client()))),
	)

	defer configDynamic.Close()


	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
