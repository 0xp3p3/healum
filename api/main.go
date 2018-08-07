package main

import (
	account_proto "server/account-srv/proto/account"
	"server/api/api"
	"server/api/utils"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	content_proto "server/content-srv/proto/content"
	note_proto "server/note-srv/proto/note"
	organisation_proto "server/organisation-srv/proto/organisation"
	plan_proto "server/plan-srv/proto/plan"
	product_proto "server/product-srv/proto/product"
	resp_proto "server/response-srv/proto/response"
	static_proto "server/static-srv/proto/static"
	survey_proto "server/survey-srv/proto/survey"
	task_proto "server/task-srv/proto/task"
	team_proto "server/team-srv/proto/team"
	todo_proto "server/todo-srv/proto/todo"
	track_proto "server/track-srv/proto/track"
	userapp_proto "server/user-app-srv/proto/userapp"
	user_proto "server/user-srv/proto/user"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-micro"
	"github.com/micro/go-os/config"
	"github.com/micro/go-os/config/source/file"
	"github.com/micro/go-os/metrics"
	_ "github.com/micro/go-plugins/broker/nats"
	"github.com/micro/go-plugins/metrics/telegraf"
	_ "github.com/micro/go-plugins/transport/nats"
	"github.com/micro/go-web"
	log "github.com/sirupsen/logrus"
)

func main() {
	configFile, _ := common.PopParameter("config")
	// Create a config instance
	config := config.NewConfig(
		// poll every hour
		config.PollInterval(time.Hour),
		// use file as a config source
		config.WithSource(file.NewSource(config.SourceName(configFile))),
	)

	defer config.Close()

	// create new metrics
	m := telegraf.NewMetrics(
		metrics.Namespace(config.Get("service", "name").String("api")),
		metrics.Collectors(
			// telegraf/statsd address
			common.MetricAddress(),
		),
	)
	defer m.Close()

	name := config.Get("service", "name").String("go.micro.api.server")
	version := config.Get("service", "version").String("latest")
	descr := config.Get("service", "description").String("Micro service")

	// cmd service
	cmd_service := micro.NewService(
		micro.Name(name),
		micro.Version(version),
		micro.Metadata(map[string]string{"Description": descr}),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*10),
	)

	cmd_service.Init()
	rateLimitedClient := utils.NewRateLimitedClient(cmd_service.Client())
	//securedClient := utils.NewCircuitBreakerClient(rateLimitedClient, config)
	// TODO return it back
	securedClient := rateLimitedClient

	// Auth
	auth_filter := api.Filters{
		AccountClient:      account_proto.NewAccountServiceClient("go.micro.srv.account", securedClient),
		UserClient:         user_proto.NewUserServiceClient("go.micro.srv.user", securedClient),
		OrganisationClient: organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", securedClient),
		TeamClient:         team_proto.NewTeamServiceClient("go.micro.srv.team", securedClient),
	}

	// Audit
	brker := cmd_service.Client().Options().Broker
	brker.Connect()
	audit_filter := api.AuditFilter{
		Broker: brker,
	}

	// User API
	user_service := api.UserService{
		UserClient:    user_proto.NewUserServiceClient("go.micro.srv.user", securedClient),
		UserAppClient: userapp_proto.NewUserAppServiceClient("go.micro.srv.userapp", securedClient),
		ContentClient: content_proto.NewContentServiceClient("go.micro.srv.content", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
	}
	user_service.Register()

	// Metrics API
	metrics_service := api.MetricsService{
		FilterMiddle:  auth_filter,
		ServerMetrics: m,
	}
	metrics_service.Register()

	// Timestamp API
	timestamp_service := api.TimestampService{}
	timestamp_service.Register()

	// Organization API
	// organization_service := api.OrganizationService{
	// 	OrganisationClient: organisation_proto.NewOrganisationClient("healum.srv.organization", securedClient),
	// 	UserClient:         user.NewAccountClient("go.micro.srv.user", securedClient),
	// 	FilterMiddle:       auth_filter,
	// 	ServerMetrics:      m,
	// }
	// organization_service.Register()

	// Plan API
	plan_service := api.PlanService{
		PlanClient:    plan_proto.NewPlanServiceClient("go.micro.srv.plan", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
	}
	plan_service.Register()

	// Task API
	task_service := api.TaskService{
		TaskClient:    task_proto.NewTaskServiceClient("go.micro.srv.task", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
	}
	task_service.Register()

	// Todo API
	todo_service := api.TodoService{
		TodoClient:    todo_proto.NewTodoServiceClient("go.micro.srv.todo", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
	}
	todo_service.Register()

	// Note API
	note_service := api.NoteService{
		NoteClient:    note_proto.NewNoteServiceClient("go.micro.srv.note", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
	}
	note_service.Register()

	// Survey API
	survey_service := api.SurveyService{
		SurveyClient:  survey_proto.NewSurveyServiceClient("go.micro.srv.survey", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
	}
	survey_service.Register()

	// Response API
	response_service := api.ResponseService{
		ResponseClient: resp_proto.NewResponseServiceClient("go.micro.srv.response", securedClient),
		SurveyClient:   survey_proto.NewSurveyServiceClient("go.micro.srv.survey", securedClient),
		Auth:           auth_filter,
		Audit:          audit_filter,
		ServerMetrics:  m,
	}
	response_service.Register()

	// Behaviour API
	behaviour_service := api.BehaviourService{
		BehaviourClient:    behaviour_proto.NewBehaviourServiceClient("go.micro.srv.behaviour", securedClient),
		Auth:               auth_filter,
		Audit:              audit_filter,
		ServerMetrics:      m,
		OrganisationClient: organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", securedClient),
		StaticClient:       static_proto.NewStaticServiceClient("go.micro.srv.static", securedClient),
	}
	behaviour_service.Register()

	// Static API
	static_service := api.StaticService{
		StaticClient:  static_proto.NewStaticServiceClient("go.micro.srv.static", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
		ContentClient: content_proto.NewContentServiceClient("go.micro.srv.content", securedClient),
	}
	static_service.Register()

	// Content API
	content_service := api.ContentService{
		ContentClient: content_proto.NewContentServiceClient("go.micro.srv.content", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
	}
	content_service.Register()

	// Account API
	account_service := api.AccountService{
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
	}
	account_service.Register()

	// Team API
	team_service := api.TeamService{
		TeamClient:    team_proto.NewTeamServiceClient("go.micro.srv.team", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
	}
	team_service.Register()

	// Track API
	track_service := api.TrackService{
		TrackClient:   track_proto.NewTrackServiceClient("go.micro.srv.track", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
	}
	track_service.Register()

	// Organisation API
	org_service := api.OrganisationService{
		OrganisationClient: organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", securedClient),
		Auth:               auth_filter,
		Audit:              audit_filter,
		ServerMetrics:      m,
	}
	org_service.Register()

	// UserApp API
	userapp_service := api.UserAppService{
		UserAppClient: userapp_proto.NewUserAppServiceClient("go.micro.srv.userapp", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
	}
	userapp_service.Register()

	// Product API
	product_service := api.ProductService{
		ProductClient: product_proto.NewProductServiceClient("go.micro.srv.product", securedClient),
		Auth:          auth_filter,
		Audit:         audit_filter,
		ServerMetrics: m,
	}
	product_service.Register()

	// Swagger service
	swagger_service := utils.SwaggerService{Config: config}
	swagger_service.Register()

	// Create service
	service := web.NewService(
		web.Name(name),
		web.Version(version),
		web.Metadata(map[string]string{"Description": descr}),
		web.Registry(cmd_service.Options().Registry),
		web.RegisterTTL(time.Minute),
		web.RegisterInterval(time.Second*10),
	)

	service.Init()
	log.SetLevel(log.DebugLevel)

	// Register Handler
	service.Handle("/", restful.DefaultContainer)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
