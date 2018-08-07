package main

import (
	"server/account-srv/db"
	"server/account-srv/handler"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	organisation_proto "server/organisation-srv/proto/organisation"
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

func RetryFunc() {

}

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
		metrics.Namespace(conf.Get("service", "name").String("account")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	common.AccountSrv = conf.Get("service", "account").String("go.micro.srv.account")
	version := conf.Get("service", "version").String("latest")
	descr := conf.Get("service", "description").String("Micro service")

	// cl := client.NewClient(
	// 	// client.Transport(nats_transport.NewTransport()),
	// 	// client.Broker(nats_broker.NewBroker()),
	// 	client.RequestTimeout(3*time.Second),
	// 	client.Retries(5),
	// 	client.Retry(func(ctx context.Context, req client.Request, retryCount int, err error) (bool, error) {
	// 		log.Println("retry:", err)
	// 		return true, nil
	// 	}),
	// )
	service := micro.NewService(
		micro.Name(common.AccountSrv),
		micro.Version(version),
		micro.Metadata(map[string]string{"Description": descr}),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*10),
		// micro.Client(cl),
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
	accountService := &handler.AccountService{
		Broker:             brker,
		KvClient:           kv_proto.NewKvServiceClient("go.micro.srv.kv", service.Client()),
		UserClient:         user_proto.NewUserServiceClient("go.micro.srv.user", service.Client()),
		TeamClient:         team_proto.NewTeamServiceClient("go.micro.srv.team", service.Client()),
		OrganisationClient: organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", service.Client()),
	}
	account_proto.RegisterAccountServiceHandler(service.Server(), accountService)
	db.Init(service.Client())

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
