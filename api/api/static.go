package api

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	content_proto "server/content-srv/proto/content"
	static_proto "server/static-srv/proto/static"
	"strconv"

	"github.com/gocarina/gocsv"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type StaticService struct {
	StaticClient  static_proto.StaticServiceClient
	Auth          Filters
	Audit         AuditFilter
	ServerMetrics metrics.Metrics
	ContentClient content_proto.ContentServiceClient
}

func (p StaticService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/static")

	audit := &audit_proto.Audit{
		ActionService: common.StaticSrv,
	}

	ws.Route(ws.GET("/apps/all").To(p.AllApps).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all apps"))

	ws.Route(ws.POST("/app").To(p.CreateApp).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a app"))

	ws.Route(ws.GET("/app/{app_id}").To(p.ReadApp).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a app"))

	ws.Route(ws.DELETE("/app/{app_id}").To(p.DeleteApp).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a app"))

	ws.Route(ws.GET("/platforms/all").To(p.AllPlatforms).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all platforms"))

	ws.Route(ws.POST("/platform").To(p.CreatePlatform).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a platform"))

	ws.Route(ws.GET("/platform/{platform_id}").To(p.ReadPlatform).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a platform"))

	ws.Route(ws.DELETE("/platform/{platform_id}").To(p.DeletePlatform).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a platform"))

	ws.Route(ws.GET("/wearables/all").To(p.AllWearables).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all wearables"))

	ws.Route(ws.POST("/wearable").To(p.CreateWearable).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a wearable"))

	ws.Route(ws.GET("/wearable/{wearable_id}").To(p.ReadWearable).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a wearable"))

	ws.Route(ws.DELETE("/wearable/{wearable_id}").To(p.DeleteWearable).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a wearable"))

	ws.Route(ws.GET("/devices/all").To(p.AllDevices).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all devices"))

	ws.Route(ws.POST("/device").To(p.CreateDevice).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a device"))

	ws.Route(ws.GET("/device/{device_id}").To(p.ReadDevice).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a device"))

	ws.Route(ws.DELETE("/device/{device_id}").To(p.DeleteDevice).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a device"))

	ws.Route(ws.GET("/markers/all").To(p.AllMarkers).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all markers"))

	ws.Route(ws.POST("/marker").To(p.CreateMarker).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Doc("Create or update a marker"))

	ws.Route(ws.GET("/marker/{marker_id}").To(p.ReadMarker).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a marker"))

	ws.Route(ws.DELETE("/marker/{marker_id}").To(p.DeleteMarker).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a marker"))

	ws.Route(ws.POST("/markers/filter").To(p.FilterMarker).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Filter marker"))

	ws.Route(ws.GET("/modules/all").To(p.AllModules).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all modules"))

	ws.Route(ws.POST("/module").To(p.CreateModule).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a module"))

	ws.Route(ws.GET("/module/{module_id}").To(p.ReadModule).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a module"))

	ws.Route(ws.DELETE("/module/{module_id}").To(p.DeleteModule).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a module"))

	ws.Route(ws.GET("/behaviour/categorys/all").To(p.AllBehaviourCategories).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all behaviourcategories"))

	ws.Route(ws.POST("/behaviour/category").To(p.CreateBehaviourCategory).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a behaviour category"))

	ws.Route(ws.GET("/behaviour/category/{category_id}").To(p.ReadBehaviourCategory).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a behaviour category"))

	ws.Route(ws.DELETE("/behaviour/category/{category_id}").To(p.DeleteBehaviourCategory).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a behaviour category"))

	ws.Route(ws.POST("/behaviour/categorys/filter").To(p.FilterBehaviourCategory).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Filter behaviour category"))

	ws.Route(ws.GET("/socialTypes/all").To(p.AllSocialTypes).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all socialTypes"))

	ws.Route(ws.POST("/socialType").To(p.CreateSocialType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a socialType"))

	ws.Route(ws.GET("/socialType/{socialType_id}").To(p.ReadSocialType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a socialType"))

	ws.Route(ws.DELETE("/socialType/{socialType_id}").To(p.DeleteSocialType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a socialType"))

	ws.Route(ws.GET("/notifications/all").To(p.AllNotifications).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all notifications"))

	ws.Route(ws.POST("/notification").To(p.CreateNotification).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a notification"))

	ws.Route(ws.GET("/notification/{notification_id}").To(p.ReadNotification).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a notification"))

	ws.Route(ws.DELETE("/notification/{notification_id}").To(p.DeleteNotification).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a notification"))

	ws.Route(ws.GET("/trackerMethods/all").To(p.AllTrackerMethods).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all trackerMethods"))

	ws.Route(ws.POST("/trackerMethod").To(p.CreateTrackerMethod).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a trackerMethod"))

	ws.Route(ws.GET("/trackerMethod/{trackerMethod_id}").To(p.ReadTrackerMethod).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a trackerMethod"))

	ws.Route(ws.DELETE("/trackerMethod/{trackerMethod_id}").To(p.DeleteTrackerMethod).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a trackerMethod"))

	ws.Route(ws.POST("/trackerMethod/filter").To(p.FilterTrackerMethod).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Filter trackerMethod"))

	ws.Route(ws.GET("/behaviourCategoryAims/all").To(p.AllBehaviourCategoryAims).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all behaviourCategoryAims"))

	ws.Route(ws.POST("/behaviourCategoryAim").To(p.CreateBehaviourCategoryAim).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a behaviourCategoryAim"))

	ws.Route(ws.GET("/behaviourCategoryAim/{behaviourCategoryAim_id}").To(p.ReadBehaviourCategoryAim).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a behaviourCategoryAim"))

	ws.Route(ws.DELETE("/behaviourCategoryAim/{behaviourCategoryAim_id}").To(p.DeleteBehaviourCategoryAim).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a behaviourCategoryAim"))

	ws.Route(ws.GET("/content/category/parents/all").To(p.AllContentParentCategories).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all contentParentCategories"))

	ws.Route(ws.POST("/content/category/parent").To(p.CreateContentParentCategory).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a contentParentCategory"))

	ws.Route(ws.GET("/content/category/parent/{contentParentCategory_id}").To(p.ReadContentParentCategory).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a contentParentCategory"))

	ws.Route(ws.DELETE("/content/category/parent/{contentParentCategory_id}").To(p.DeleteContentParentCategory).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a contentParentCategory"))

	ws.Route(ws.GET("/content/categorys/all").To(p.AllContentCategories).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all contentCategories"))

	ws.Route(ws.POST("/content/category").To(p.CreateContentCategory).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a contentCategory"))

	ws.Route(ws.GET("/content/category/{contentCategory_id}").To(p.ReadContentCategory).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a contentCategory"))

	ws.Route(ws.DELETE("/content/category/{contentCategory_id}").To(p.DeleteContentCategory).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a contentCategory"))

	ws.Route(ws.GET("/content/types/all").To(p.AllContentTypes).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all contentTypes"))

	ws.Route(ws.POST("/content/type").To(p.CreateContentType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a contentType"))

	ws.Route(ws.GET("/content/type/{contentType_id}").To(p.ReadContentType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a contentType"))

	ws.Route(ws.DELETE("/content/type/{contentType_id}").To(p.DeleteContentType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a contentType"))

	ws.Route(ws.GET("/content/source/types/all").To(p.AllContentSourceTypes).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all contentSourceTypes"))

	ws.Route(ws.POST("/content/source/type").To(p.CreateContentSourceType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a contentSourceType"))

	ws.Route(ws.GET("/content/source/type/{contentSourceType_id}").To(p.ReadContentSourceType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a contentSourceType"))

	ws.Route(ws.DELETE("/content/source/type/{contentSourceType_id}").To(p.DeleteContentSourceType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Doc("Delete a contentSourceType"))

	ws.Route(ws.GET("/module/triggers/all").To(p.AllModuleTriggers).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all moduleTriggers"))

	ws.Route(ws.POST("/module/trigger").To(p.CreateModuleTrigger).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a moduleTrigger"))

	ws.Route(ws.GET("/module/trigger/{moduleTrigger_id}").To(p.ReadModuleTrigger).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a moduleTrigger"))

	ws.Route(ws.DELETE("/module/trigger/{moduleTrigger_id}").To(p.DeleteModuleTrigger).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a moduleTrigger"))

	ws.Route(ws.POST("/module/triggers/filter").To(p.FilterModuleTrigger).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Filter marker"))

	ws.Route(ws.GET("/trigger/content/types/all").To(p.AllTriggerContentTypes).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all triggerContentTypes"))

	ws.Route(ws.POST("/trigger/content/type").To(p.CreateTriggerContentType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a triggerContentType"))

	ws.Route(ws.GET("/trigger/content/type/{triggerContentType_id}").To(p.ReadTriggerContentType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a triggerContentType"))

	ws.Route(ws.DELETE("/trigger/content/type/{triggerContentType_id}").To(p.DeleteTriggerContentType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a triggerContentType"))

	ws.Route(ws.POST("/trigger/content/types/filter").To(p.FilterTriggerContentType).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Filter marker"))

	ws.Route(ws.GET("/setbacks/all").To(p.AllSetbacks).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		Filter(p.Auth.Paginate).
		Filter(p.Auth.SortFilter).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("List all setbacks"))

	ws.Route(ws.POST("/setback").To(p.CreateSetback).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Create or update a setback"))

	ws.Route(ws.GET("/setback/{setback_id}").To(p.ReadSetback).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Read a setback"))

	ws.Route(ws.DELETE("/setback/{setback_id}").To(p.DeleteSetback).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Delete a setback"))

	ws.Route(ws.POST("/setback/search/autocomplete").To(p.AutocompleteSetbackSearch).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete setback"))

	ws.Route(ws.POST("/setback/search/autocomplete").To(p.AutocompleteSetbackSearch).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete setback"))

	ws.Route(ws.POST("/setback/search/autocomplete").To(p.AutocompleteSetbackSearch).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete setback"))

	ws.Route(ws.POST("/setback/search/autocomplete").To(p.AutocompleteSetbackSearch).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete setback"))

	ws.Route(ws.POST("/setback/search/autocomplete").To(p.AutocompleteSetbackSearch).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete setback"))

	ws.Route(ws.POST("/behaviourCategoryAim/upload").To(p.UploadBehaviourCategoryAim).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Upload behaviour category aim"))

	ws.Route(ws.POST("/content/category/upload").To(p.UploadContentCategory).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete setback"))

	ws.Route(ws.POST("/marker/upload").To(p.UploadMarker).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete setback"))

	ws.Route(ws.POST("/behaviour/category/upload").To(p.UploadBehaviourCategory).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete setback"))

	ws.Route(ws.POST("/content/category/item/upload").To(p.UploadContentCategoryItem).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete setback"))

	ws.Route(ws.POST("/trackerMethod/upload").To(p.UploadTrackerMethod).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.EmployeeAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Search autocomplete setback"))

	restful.Add(ws)
}

/**
* @api {get} /server/static/apps/all?session={session_id}&offset={offset}&limit={limit} List all apps
* @apiVersion 0.1.0
* @apiName AllApps
* @apiGroup Static
*
* @apiDescription List all apps
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/apps/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "apps": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "image": "image",
*         "tags": ["tag1", "tag2"],
*         "platforms": [ Platform, ... ],
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all apps successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The apps were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllApps",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllApps",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllApps(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllApps API request")
	req_app := new(static_proto.AllAppsRequest)
	req_app.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_app.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_app.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_app.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_app.SortParameter = req.Attribute(SortParameter).(string)
	req_app.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllApps(ctx, req_app)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllApps", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all apps successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/app?session={session_id} Create or update a app
* @apiVersion 0.1.0
* @apiName CreateApp
* @apiGroup Static
*
* @apiDescription Create or update a app
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/app?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "app": {
*     "id": "111",
*     "name": "title",
*     "summary": "summary",
*     "description": "description",
*     "icon_slug": "iconslug",
*     "image": "image",
*     "tags": ["tag1", "tag2"],
*     "platforms": [ Platform, ... ]
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "app": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "image": "image",
*       "tags": ["tag1", "tag2"],
*       "platforms": [ Platform, ... ],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created app successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateApp",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateApp(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateApp API request")
	req_app := new(static_proto.CreateAppRequest)
	err := req.ReadEntity(req_app)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateApp", "BindError")
		return
	}
	req_app.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_app.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateApp(ctx, req_app)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateApp", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created app successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/app/{app_id}?session={session_id} View app detail
* @apiVersion 0.1.0
* @apiName ReadApp
* @apiGroup Static
*
* @apiDescription View app detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/app/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "app": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "image": "image",
*       "tags": ["tag1", "tag2"],
*       "platforms": [ Platform, ... ],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read app successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The app was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadApp",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadApp",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadApp(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadApp API request")
	req_app := new(static_proto.ReadAppRequest)
	req_app.Id = req.PathParameter("app_id")
	req_app.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_app.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadApp(ctx, req_app)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadApp", "ReadError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read app successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/app/{app_id}?session={session_id} Delete a app
* @apiVersion 0.1.0
* @apiName DeleteApp
* @apiGroup Static
*
* @apiDescription Delete a app
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/app/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted app successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The app was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteApp",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteApp(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteApp API request")
	req_app := new(static_proto.DeleteAppRequest)
	req_app.Id = req.PathParameter("app_id")
	req_app.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_app.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteApp(ctx, req_app)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteApp", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted app successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/platforms/all?session={session_id}&offset={offset}&limit={limit} List all platforms
* @apiVersion 0.1.0
* @apiName AllPlatforms
* @apiGroup Static
*
* @apiDescription List all platforms
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/platforms/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "platforms": [
*       {
*         "id": "111",
*         "name": "title",
*         "url": "url",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all platforms successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The platforms were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllPlatforms",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllPlatforms",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllPlatforms(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllPlatforms API request")
	req_platform := new(static_proto.AllPlatformsRequest)
	req_platform.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_platform.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_platform.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_platform.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_platform.SortParameter = req.Attribute(SortParameter).(string)
	req_platform.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllPlatforms(ctx, req_platform)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllPlatforms", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all platforms successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/platform?session={session_id}&offset={offset}&limit={limit} Create or update a platform
* @apiVersion 0.1.0
* @apiName CreatePlatform
* @apiGroup Static
*
* @apiDescription Create or update a platform
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/platform?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "platform": {
*      "id": "111",
*      "name": "title",
*      "url": "url"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "platform": {
*       "id": "111",
*       "name": "title",
*       "url": "url",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created platform successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreatePlatform",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreatePlatform(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreatePlatform API request")
	req_platform := new(static_proto.CreatePlatformRequest)
	err := req.ReadEntity(req_platform)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreatePlatform", "BindError")
		return
	}
	req_platform.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_platform.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreatePlatform(ctx, req_platform)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreatePlatform", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created platform successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/platform/{platform_id}?session={session_id} View platform detail
* @apiVersion 0.1.0
* @apiName ReadPlatform
* @apiGroup Static
*
* @apiDescription View platform detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/platform/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "platform": {
*       "id": "111",
*       "name": "title",
*       "url": "url",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read platform successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The platform was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadPlatform",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadPlatform",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadPlatform(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadPlatform API request")
	req_platform := new(static_proto.ReadPlatformRequest)
	req_platform.Id = req.PathParameter("platform_id")
	req_platform.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_platform.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadPlatform(ctx, req_platform)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadPlatform", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read platform successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/platform/{platform_id}?session={session_id} Delete a platform
* @apiVersion 0.1.0
* @apiName DeletePlatform
* @apiGroup Static
*
* @apiDescription Delete a platform
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/platform/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted platform successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The platform was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeletePlatform",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeletePlatform(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeletePlatform API request")
	req_platform := new(static_proto.DeletePlatformRequest)
	req_platform.Id = req.PathParameter("platform_id")
	req_platform.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_platform.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeletePlatform(ctx, req_platform)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeletePlatform", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted platform successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/wearables/all?session={session_id}&offset={offset}&limit={limit} List all wearables
* @apiVersion 0.1.0
* @apiName AllWearables
* @apiGroup Static
*
* @apiDescription List all wearables
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/wearables/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "wearables": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "image": "image",
*         "url": "url",
*         "tags": ["tag1", "tag2"],
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all wearables successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The wearables were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllWearables",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllWearables",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllWearables(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllWearables API request")
	req_wearable := new(static_proto.AllWearablesRequest)
	req_wearable.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_wearable.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_wearable.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_wearable.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_wearable.SortParameter = req.Attribute(SortParameter).(string)
	req_wearable.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllWearables(ctx, req_wearable)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllWearables", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all wearables successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/wearable?session={session_id} Create or update a wearable
* @apiVersion 0.1.0
* @apiName CreateWearable
* @apiGroup Static
*
* @apiDescription Create or update a wearable
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/wearable?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "wearable": {
*     "id": "111",
*     "name": "title",
*     "summary": "summary",
*     "description": "description",
*     "icon_slug": "iconslug",
*     "image": "image",
*     "url": "url",
*     "tags": ["tag1", "tag2"]
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "wearable": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "image": "image",
*       "url": "url",
*       "tags": ["tag1", "tag2"],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created wearable successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateWearable",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateWearable(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateWearable API request")
	req_wearable := new(static_proto.CreateWearableRequest)
	err := req.ReadEntity(req_wearable)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateWearable", "BindError")
		return
	}
	req_wearable.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_wearable.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateWearable(ctx, req_wearable)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateWearable", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created wearable successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/wearable/{wearable_id}?session={session_id} View wearable detail
* @apiVersion 0.1.0
* @apiName ReadWearable
* @apiGroup Static
*
* @apiDescription View wearable detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/wearable/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "wearable": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "image": "image",
*       "url": "url",
*       "tags": ["tag1", "tag2"],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read wearable successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The wearable was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadWearable",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadWearable",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadWearable(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadWearable API request")
	req_wearable := new(static_proto.ReadWearableRequest)
	req_wearable.Id = req.PathParameter("wearable_id")
	req_wearable.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_wearable.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadWearable(ctx, req_wearable)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadWearable", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read wearable successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/wearable/{wearable_id}?session={session_id} Delete a wearable
* @apiVersion 0.1.0
* @apiName DeleteWearable
* @apiGroup Static
*
* @apiDescription Delete a wearable
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/wearable/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted wearable successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The wearable was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteWearable",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteWearable(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteWearable API request")
	req_wearable := new(static_proto.DeleteWearableRequest)
	req_wearable.Id = req.PathParameter("wearable_id")
	req_wearable.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_wearable.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteWearable(ctx, req_wearable)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteWearable", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted wearable successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/devices/all?session={session_id}&offset={offset}&limit={limit} List all devices
* @apiVersion 0.1.0
* @apiName AllDevices
* @apiGroup Static
*
* @apiDescription List all devices
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/devices/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "devices": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "image": "image",
*         "url": "url",
*         "tags": ["tag1", "tag2"]
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all devices successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The devices were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllDevices",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllDevices",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllDevices(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllDevices API request")
	req_device := new(static_proto.AllDevicesRequest)
	req_device.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_device.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_device.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_device.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_device.SortParameter = req.Attribute(SortParameter).(string)
	req_device.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllDevices(ctx, req_device)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllDevices", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all devices successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/device?session={session_id} Create or update a device
* @apiVersion 0.1.0
* @apiName CreateDevice
* @apiGroup Static
*
* @apiDescription Create or update a device
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/device?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "device": {
*     "id": "111",
*     "name": "title",
*     "summary": "summary",
*     "description": "description",
*     "icon_slug": "iconslug",
*     "image": "image",
*     "url": "url",
*     "tags": ["tag1", "tag2"]
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "device": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "image": "image",
*       "url": "url",
*       "tags": ["tag1", "tag2"],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created device successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateDevice",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateDevice(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateDevice API request")
	req_device := new(static_proto.CreateDeviceRequest)
	err := req.ReadEntity(req_device)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateDevice", "BindError")
		return
	}
	req_device.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_device.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateDevice(ctx, req_device)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateDevice", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created device successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/device/{device_id}?session={session_id} View device detail
* @apiVersion 0.1.0
* @apiName ReadDevice
* @apiGroup Static
*
* @apiDescription View device detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/device/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "device": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "image": "image",
*       "url": "url",
*       "tags": ["tag1", "tag2"],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read device successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The device was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadDevice",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadDevice",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadDevice(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadDevice API request")
	req_device := new(static_proto.ReadDeviceRequest)
	req_device.Id = req.PathParameter("device_id")
	req_device.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_device.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadDevice(ctx, req_device)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadDevice", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read device successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/device/{device_id}?session={session_id} Delete a device
* @apiVersion 0.1.0
* @apiName DeleteDevice
* @apiGroup Static
*
* @apiDescription Delete a device
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/device/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted device successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The device was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteDevice",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteDevice(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteDevice API request")
	req_device := new(static_proto.DeleteDeviceRequest)
	req_device.Id = req.PathParameter("device_id")
	req_device.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_device.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteDevice(ctx, req_device)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteDevice", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted device successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/markers/all?session={session_id}&offset={offset}&limit={limit} List all markers
* @apiVersion 0.1.0
* @apiName AllMarkers
* @apiGroup Static
*
* @apiDescription List all markers
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/markers/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "markers": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "org_id": "orgid",
*         "unit": "unit",
*         "apps": [ App, ...],
*         "wearables": [ Wearable, ...],
*         "devices": [ Device, ...],
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all markers successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The markers were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllMarkers",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllMarkers",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllMarkers(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllMarkers API request")
	req_marker := new(static_proto.AllMarkersRequest)
	req_marker.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_marker.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_marker.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_marker.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_marker.SortParameter = req.Attribute(SortParameter).(string)
	req_marker.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllMarkers(ctx, req_marker)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllMarkers", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all markers successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/marker?session={session_id} Create or update a marker
* @apiVersion 0.1.0
* @apiName CreateMarker
* @apiGroup Static
*
* @apiDescription Create or update a marker
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/marker?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "marker": {
*     "id": "111",
*     "name": "title",
*     "summary": "summary",
*     "description": "description",
*     "icon_slug": "iconslug",
*     "org_id": "orgid",
*     "unit": "unit",
*     "apps": [ App, ...],
*     "wearables": [ Wearable, ...],
*     "devices": [ Device, ...],
*     "trackerMethods": [ TrackerMethod, ...],
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "marker": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "org_id": "orgid",
*       "unit": "unit",
*       "apps": [ App, ...],
*       "wearables": [ Wearable, ...],
*       "devices": [ Device, ...],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created marker successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateMarker",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateMarker(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateMarker API request")
	req_marker := new(static_proto.CreateMarkerRequest)
	err := req.ReadEntity(req_marker)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateMarker", "BindError")
		return
	}
	req_marker.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_marker.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateMarker(ctx, req_marker)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateMarker", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created marker successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/marker/filter?session={session_id}&offset={offset}&limit={limit} List all markers
* @apiVersion 0.1.0
* @apiName FilterMarker
* @apiGroup Static
*
* @apiDescription Filter marker
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/marker/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "trackerMethods": ["111"]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "markers": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "org_id": "orgid",
*         "unit": "unit",
*         "apps": [ App, ...],
*         "wearables": [ Wearable, ...],
*         "devices": [ Device, ...],
*         "trackerMethods": [{Id: "111"}, ...],
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all markers successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The markers were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.FilterMarker",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.FilterMarker",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) FilterMarker(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.FilterMarker API request")
	req_marker := new(static_proto.FilterMarkerRequest)
	err := req.ReadEntity(req_marker)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.FilterMarker", "BindError")
		return
	}
	req_marker.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_marker.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_marker.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_marker.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_marker.SortParameter = req.Attribute(SortParameter).(string)
	req_marker.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.FilterMarker(ctx, req_marker)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.FilterMarker", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Filter markers successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/marker/{marker_id}?session={session_id} View marker detail
* @apiVersion 0.1.0
* @apiName ReadMarker
* @apiGroup Static
*
* @apiDescription View marker detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/marker/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "marker": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "org_id": "orgid",
*       "unit": "unit",
*       "apps": [ App, ...],
*       "wearables": [ Wearable, ...],
*       "devices": [ Device, ...],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read marker successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The marker was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadMarker",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadMarker",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadMarker(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadMarker API request")
	req_marker := new(static_proto.ReadMarkerRequest)
	req_marker.Id = req.PathParameter("marker_id")
	req_marker.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_marker.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadMarker(ctx, req_marker)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadMarker", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read marker successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/marker/{marker_id}?session={session_id}&offset={offset}&limit={limit} Delete a marker
* @apiVersion 0.1.0
* @apiName DeleteMarker
* @apiGroup Static
*
* @apiDescription Delete a marker
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/marker/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted marker successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The marker was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteMarker",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteMarker(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteMarker API request")
	req_marker := new(static_proto.DeleteMarkerRequest)
	req_marker.Id = req.PathParameter("marker_id")
	req_marker.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_marker.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteMarker(ctx, req_marker)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteMarker", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted marker successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/modules/all?session={session_id}&offset={offset}&limit={limit} List all modules
* @apiVersion 0.1.0
* @apiName AllModules
* @apiGroup Static
*
* @apiDescription List all modules
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/modules/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "modules": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "org_id": "orgid",
*         "settings": "settings",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all modules successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The modules were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllModules",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllModules",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllModules(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllModules API request")
	req_module := new(static_proto.AllModulesRequest)
	req_module.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_module.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_module.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_module.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_module.SortParameter = req.Attribute(SortParameter).(string)
	req_module.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllModules(ctx, req_module)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllModules", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all modules successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/module?session={session_id} Create or update a module
* @apiVersion 0.1.0
* @apiName CreateModule
* @apiGroup Static
*
* @apiDescription Create or update a module
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/module?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "module": {
*     "id": "111",
*     "name": "title",
*     "summary": "summary",
*     "description": "description",
*     "icon_slug": "iconslug",
*     "org_id": "orgid",
*     "settings": "settings"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "module": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "org_id": "orgid",
*       "settings": "settings",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created module successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateModule",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateModule(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateModule API request")
	req_module := new(static_proto.CreateModuleRequest)
	err := req.ReadEntity(req_module)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateModule", "BindError")
		return
	}
	req_module.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_module.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateModule(ctx, req_module)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateModule", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created module successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/module/{module_id}?session={session_id} View module detail
* @apiVersion 0.1.0
* @apiName ReadModule
* @apiGroup Static
*
* @apiDescription View module detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/module/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "module": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "org_id": "orgid",
*       "settings": "settings",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read module successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The module was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadModule",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadModule",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadModule(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadModule API request")
	req_module := new(static_proto.ReadModuleRequest)
	req_module.Id = req.PathParameter("module_id")
	req_module.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_module.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadModule(ctx, req_module)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadModule", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read module successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/module/{module_id}?session={session_id} Delete a module
* @apiVersion 0.1.0
* @apiName DeleteModule
* @apiGroup Static
*
* @apiDescription Delete a module
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/module/111?session={session_id}
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted module successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The module was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteModule",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteModule(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteModule API request")
	req_module := new(static_proto.DeleteModuleRequest)
	req_module.Id = req.PathParameter("module_id")
	req_module.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_module.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteModule(ctx, req_module)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteModule", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted module successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/behaviour/categorys/all?session={session_id}&offset={offset}&limit={limit} List all behaviour categories
* @apiVersion 0.1.0
* @apiName AllBehaviourCategories
* @apiGroup Static
*
* @apiDescription List all behaviour categories
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/behaviour/categorys/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "categories": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "org"_id: "orgid",
*         "image": "image",
*         "tags": ["tag1", "tag2"],
*         "aims": [ Aim, ... ],
*         "markerDefault": { Marker },
*         "markerOptions": [ Marker, ... ],
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all behaviour categories successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The behaviour categories were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllBehaviourCategories",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllBehaviourCategories",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllBehaviourCategories(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllBehaviourCategories API request")
	req_behaviourcategory := new(static_proto.AllBehaviourCategoriesRequest)
	req_behaviourcategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_behaviourcategory.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_behaviourcategory.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_behaviourcategory.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_behaviourcategory.SortParameter = req.Attribute(SortParameter).(string)
	req_behaviourcategory.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllBehaviourCategories(ctx, req_behaviourcategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllBehaviourCategories", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read all behaviour categories successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/behaviour/category?session={session_id} Create or update a behaviour category
* @apiVersion 0.1.0
* @apiName CreateBehaviourCategory
* @apiGroup Static
*
* @apiDescription Create or update a behaviour category
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/behaviour/category?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "category": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "org"_id: "orgid",
*       "image": "image",
*       "tags": ["tag1", "tag2"],
*       "aims": [ Aim, ... ],
*       "markerDefault": { Marker },
*       "markerOptions": [ Marker, ... ]
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "category": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "org"_id: "orgid",
*       "image": "image",
*       "tags": ["tag1", "tag2"],
*       "aims": [ Aim, ... ],
*       "markerDefault": { Marker },
*       "markerOptions": [ Marker, ... ],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created behaviour category successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateBehaviourCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateBehaviourCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateBehaviourCategory API request")
	req_behaviourcategory := new(static_proto.CreateBehaviourCategoryRequest)
	err := req.ReadEntity(req_behaviourcategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateBehaviourCategory", "BindError")
		return
	}
	req_behaviourcategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_behaviourcategory.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateBehaviourCategory(ctx, req_behaviourcategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateBehaviourCategory", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created behaviour category successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/behaviour/category/{behaviourcategory_id}?session={session_id} View behaviour category detail
* @apiVersion 0.1.0
* @apiName ReadBehaviourCategory
* @apiGroup Static
*
* @apiDescription View behaviour category detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/behaviour/category/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "category": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "org"_id: "orgid",
*       "image": "image",
*       "tags": ["tag1", "tag2"],
*       "aims": [ Aim, ... ],
*       "markerDefault": { Marker },
*       "markerOptions": [ Marker, ... ],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read behaviour category successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The behaviour category was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadBehaviourCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadBehaviourCategory",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadBehaviourCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadBehaviourCategory API request")
	req_behaviourcategory := new(static_proto.ReadBehaviourCategoryRequest)
	req_behaviourcategory.Id = req.PathParameter("behaviourcategory_id")
	req_behaviourcategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_behaviourcategory.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadBehaviourCategory(ctx, req_behaviourcategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadBehaviourCategory", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read behaviour category successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/behaviour/category/{behaviourcategory_id}?session={session_id} Delete a behaviour category
* @apiVersion 0.1.0
* @apiName DeleteBehaviourCategory
* @apiGroup Static
*
* @apiDescription Delete a behaviour category
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/behaviour/category/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted behaviour category successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The behaviour category was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteBehaviourCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteBehaviourCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteBehaviourCategory API request")
	req_behaviourcategory := new(static_proto.DeleteBehaviourCategoryRequest)
	req_behaviourcategory.Id = req.PathParameter("behaviourcategory_id")
	req_behaviourcategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_behaviourcategory.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteBehaviourCategory(ctx, req_behaviourcategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteBehaviourCategory", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted behaviour category successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/behaviour/category/filter?session={session_id}&offset={offset}&limit={limit} List all behaviour categories
* @apiVersion 0.1.0
* @apiName FilterBehaviourCategory
* @apiGroup Static
*
* @apiDescription Filter behaviour category
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/behaviour/category/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "markers": ["111"]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "categories": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "org"_id: "orgid",
*         "image": "image",
*         "tags": ["tag1", "tag2"],
*         "aims": [ Aim, ... ],
*         "markerDefault": { Marker },
*         "markerOptions": [ {Id: "111"}, ... ],
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Filter behaviour categories successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The behaviour categories were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.FilterBehaviourCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.FilterBehaviourCategory",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) FilterBehaviourCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.FilterBehaviourCategory API request")
	req_behaviourCategory := new(static_proto.FilterBehaviourCategoryRequest)
	err := req.ReadEntity(req_behaviourCategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.FilterBehaviourCategory", "BindError")
		return
	}
	req_behaviourCategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_behaviourCategory.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_behaviourCategory.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_behaviourCategory.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_behaviourCategory.SortParameter = req.Attribute(SortParameter).(string)
	req_behaviourCategory.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.FilterBehaviourCategory(ctx, req_behaviourCategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.FilterBehaviourCategory", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Filter behaviour categories successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/socialTypes/all?session={session_id}&offset={offset}&limit={limit} List all social types
* @apiVersion 0.1.0
* @apiName AllSocialTypes
* @apiGroup Static
*
* @apiDescription List all socialTypes
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/socialTypes/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "socialTypes": [
*       {
*         "id": "111",
*         "name": "title",
*         "url": "http://www.example.com",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all socialTypes successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The socialTypes were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllSocialTypes",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllSocialTypes",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllSocialTypes(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllSocialTypes API request")
	req_socialType := new(static_proto.AllSocialTypesRequest)
	req_socialType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_socialType.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_socialType.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_socialType.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_socialType.SortParameter = req.Attribute(SortParameter).(string)
	req_socialType.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllSocialTypes(ctx, req_socialType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllSocialTypes", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all socialTypes successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/socialType?session={session_id} Create or update a social type
* @apiVersion 0.1.0
* @apiName CreateSocialType
* @apiGroup Static
*
* @apiDescription Create or update a socialType
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/socialType?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "socialType": {
*       "id": "111",
*       "name": "title",
*       "url": "http://www.example.com"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "socialType": {
*       "id": "111",
*       "name": "title",
*       "url": "http://www.example.com",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created socialType successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateSocialType",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateSocialType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateSocialType API request")
	req_socialType := new(static_proto.CreateSocialTypeRequest)
	err := req.ReadEntity(req_socialType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateSocialType", "BindError")
		return
	}
	req_socialType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_socialType.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateSocialType(ctx, req_socialType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateSocialType", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created socialType successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/socialType/{socialType_id}?session={session_id} View social type detail
* @apiVersion 0.1.0
* @apiName ReadSocialType
* @apiGroup Static
*
* @apiDescription View socialType detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/socialType/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "socialType": {
*       "id": "111",
*       "name": "title",
*       "url": "http://www.example.com",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read socialType successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The socialType was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadSocialType",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadSocialType",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadSocialType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadSocialType API request")
	req_socialType := new(static_proto.ReadSocialTypeRequest)
	req_socialType.Id = req.PathParameter("socialType_id")
	req_socialType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_socialType.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadSocialType(ctx, req_socialType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadSocialType", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read socialType successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/socialType/{socialType_id}?session={session_id} Delete a SocialType
* @apiVersion 0.1.0
* @apiName DeleteSocialType
* @apiGroup Static
*
* @apiDescription Delete a SocialType
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/socialType/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted SocialType successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The socialType was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteSocialType",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteSocialType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteSocialType API request")
	req_socialType := new(static_proto.DeleteSocialTypeRequest)
	req_socialType.Id = req.PathParameter("socialType_id")
	req_socialType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_socialType.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteSocialType(ctx, req_socialType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteSocialType", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted SocialType successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/notifications/all?session={session_id}&offset={offset}&limit={limit} List all notifications
* @apiVersion 0.1.0
* @apiName AllNotifications
* @apiGroup Static
*
* @apiDescription List all notifications
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/notifications/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "notifications": [
*       {
*         "id": "111",
*         "module_id": "moduleId",
*         "name": "title",
*         "description": "description",
*         "target": { NotificationTarget },
*         "name_slug": "nameSlug",
*         "icon_slug": "iconSlug",
*         "notificationReminder": 10,
*         "unit": "mins",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all notifications successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The notifications were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllNotifications",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllNotifications",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllNotifications(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllNotifications API request")
	req_notification := new(static_proto.AllNotificationsRequest)
	req_notification.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_notification.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_notification.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_notification.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_notification.SortParameter = req.Attribute(SortParameter).(string)
	req_notification.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllNotifications(ctx, req_notification)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllNotifications", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all notifications successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/static/notification?session={session_id} Create or update a notification
* @apiVersion 0.1.0
* @apiName CreateNotification
* @apiGroup Static
*
* @apiDescription Create or update a notification
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/notification?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "notification": {
*       "id": "111",
*       "module_id": "moduleId",
*       "name": "title",
*       "description": "description",
*       "target": { NotificationTarget },
*       "name_slug": "nameSlug",
*       "icon_slug": "iconSlug",
*       "notificationReminder": 10,
*       "unit": "mins"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "notification": {
*       "id": "111",
*       "module_id": "moduleId",
*       "name": "title",
*       "description": "description",
*       "target": { NotificationTarget },
*       "name_slug": "nameSlug",
*       "icon_slug": "iconSlug",
*       "notificationReminder": 10,
*       "unit": "mins",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created notification successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateNotification",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateNotification(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateNotification API request")
	req_notification := new(static_proto.CreateNotificationRequest)
	// err := req.ReadEntity(req_notification)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateNotification", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_notification); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateNotification", "BindError")
		return
	}
	req_notification.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_notification.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateNotification(ctx, req_notification)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateNotification", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created notification successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/static/notification/{notification_id}?session={session_id} View notification detail
* @apiVersion 0.1.0
* @apiName ReadNotification
* @apiGroup Static
*
* @apiDescription View notification detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/notification/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "notification": {
*       "id": "111",
*       "module_id": "moduleId",
*       "name": "title",
*       "description": "description",
*       "target": { NotificationTarget },
*       "name_slug": "nameSlug",
*       "icon_slug": "iconSlug",
*       "notificationReminder": 10,
*       "unit": "mins",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read notification successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The notification was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadNotification",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadNotification",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadNotification(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadNotification API request")
	req_notification := new(static_proto.ReadNotificationRequest)
	req_notification.Id = req.PathParameter("notification_id")
	req_notification.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_notification.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadNotification(ctx, req_notification)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadNotification", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read notification successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {delete} /server/static/notification/{notification_id}?session={session_id} Delete a notification
* @apiVersion 0.1.0
* @apiName DeleteNotification
* @apiGroup Static
*
* @apiDescription Delete a notification
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/notification/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted notification successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The notification was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteNotification",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteNotification(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteNotification API request")
	req_notification := new(static_proto.DeleteNotificationRequest)
	req_notification.Id = req.PathParameter("notification_id")
	req_notification.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_notification.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteNotification(ctx, req_notification)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteNotification", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted notification successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/trackerMethods/all?session={session_id}&offset={offset}&limit={limit}  List all TrackerMethods
* @apiVersion 0.1.0
* @apiName AllTrackerMethods
* @apiGroup Static
*
* @apiDescription List all TrackerMethods
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/trackerMethods/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "trackerMethods": [
*       {
*         "id": "111",
*         "name": "title",
*         "name_slug": "nameSlug",
*         "icon_slug": "iconSlug",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all TtrackerMethods successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The TrackerMethods were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllTrackerMethods",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllTrackerMethods",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllTrackerMethods(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllTrackerMethods API request")
	req_trackerMethod := new(static_proto.AllTrackerMethodsRequest)
	req_trackerMethod.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_trackerMethod.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_trackerMethod.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_trackerMethod.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_trackerMethod.SortParameter = req.Attribute(SortParameter).(string)
	req_trackerMethod.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllTrackerMethods(ctx, req_trackerMethod)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllTrackerMethods", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all TrackerMethods successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/trackerMethod?session={session_id} Create or update a trackerMethod
* @apiVersion 0.1.0
* @apiName CreateTrackerMethod
* @apiGroup Static
*
* @apiDescription Create or update a trackerMethod
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/trackerMethod?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "trackerMethod": {
*     "id": "111",
*     "name": "title",
*     "name_slug": "nameSlug",
*     "icon_slug": "iconSlug"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "trackerMethod": {
*       "id": "111",
*       "name": "title",
*       "name_slug": "nameSlug",
*       "icon_slug": "iconSlug",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created TrackerMethod successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateTrackerMethod",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateTrackerMethod(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateTrackerMethod API request")
	req_trackerMethod := new(static_proto.CreateTrackerMethodRequest)
	err := req.ReadEntity(req_trackerMethod)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateTrackerMethod", "BindError")
		return
	}
	req_trackerMethod.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_trackerMethod.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateTrackerMethod(ctx, req_trackerMethod)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateTrackerMethod", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created trackerMethod successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/trackerMethod/{trackerMethod_id}?session={session_id} View TrackerMethod detail
* @apiVersion 0.1.0
* @apiName ReadTrackerMethod
* @apiGroup Static
*
* @apiDescription View TrackerMethod detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/trackerMethod/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "trackerMethod": {
*       "id": "111",
*       "name": "title",
*       "name_slug": "nameSlug",
*       "icon_slug": "iconSlug",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read TrackerMethod successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The TrackerMethod was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadTrackerMethod",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadTrackerMethod",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadTrackerMethod(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadTrackerMethod API request")
	req_trackerMethod := new(static_proto.ReadTrackerMethodRequest)
	req_trackerMethod.Id = req.PathParameter("trackerMethod_id")
	req_trackerMethod.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_trackerMethod.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadTrackerMethod(ctx, req_trackerMethod)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadTrackerMethod", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read trackerMethod successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/trackerMethod/{trackerMethod_id}?session={session_id} Delete a TrackerMethod
* @apiVersion 0.1.0
* @apiName DeleteTrackerMethod
* @apiGroup Static
*
* @apiDescription Delete a TrackerMethod
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/trackerMethod/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted TrackerMethod successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The TrackerMethod was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteTrackerMethod",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteTrackerMethod(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteTrackerMethod API request")
	req_trackerMethod := new(static_proto.DeleteTrackerMethodRequest)
	req_trackerMethod.Id = req.PathParameter("trackerMethod_id")
	req_trackerMethod.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_trackerMethod.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteTrackerMethod(ctx, req_trackerMethod)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteTrackerMethod", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted TrackerMethod successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/trackerMethod/filter?session={session_id}&offset={offset}&limit={limit} Filter trackerMethods
* @apiVersion 0.1.0
* @apiName FilterTrackerMethod
* @apiGroup Static
*
* @apiDescription Filter trackerMethod
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/trackerMethod/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "markers": ["111"]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "trackerMethods": [
*       {
*         "id": "111",
*         "name": "title",
*         "name_slug": "nameSlug",
*         "icon_slug": "iconSlug",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Filter trackerMethods successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The trackerMethods were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.FilterTrackerMethod",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.FilterTrackerMethod",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) FilterTrackerMethod(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.FilterTrackerMethod API request")
	req_trackerMethod := new(static_proto.FilterTrackerMethodRequest)
	err := req.ReadEntity(req_trackerMethod)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.FilterTrackerMethod", "BindError")
		return
	}
	req_trackerMethod.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_trackerMethod.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_trackerMethod.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_trackerMethod.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_trackerMethod.SortParameter = req.Attribute(SortParameter).(string)
	req_trackerMethod.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.FilterTrackerMethod(ctx, req_trackerMethod)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.FilterTrackerMethod", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Filter trackerMethods successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/behaviourCategoryAims/all?session={session_id}&offset={offset}&limit={limit}  List all behaviourCategoryAims
* @apiVersion 0.1.0
* @apiName AllBehaviourCategoryAims
* @apiGroup Static
*
* @apiDescription List all behaviourCategoryAims
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/behaviourCategoryAims/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "behaviourCategoryAims": [
*       {
*         "id": "111",
*         "name": "title",
*         "name_slug": "nameSlug",
*         "icon_slug": "iconslug",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all BehaviourCategoryAims successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The behaviourCategoryAims were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllBehaviourCategoryAims",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllBehaviourCategoryAims",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllBehaviourCategoryAims(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllBehaviourCategoryAims API request")
	req_behaviourCategoryAim := new(static_proto.AllBehaviourCategoryAimsRequest)
	req_behaviourCategoryAim.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_behaviourCategoryAim.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_behaviourCategoryAim.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_behaviourCategoryAim.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_behaviourCategoryAim.SortParameter = req.Attribute(SortParameter).(string)
	req_behaviourCategoryAim.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllBehaviourCategoryAims(ctx, req_behaviourCategoryAim)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllBehaviourCategoryAims", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read all behaviourCategoryAims successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/behaviourCategoryAim?session={session_id}&offset={offset}&limit={limit} Create or update a behaviourCategoryAim
* @apiVersion 0.1.0
* @apiName CreateBehaviourCategoryAim
* @apiGroup Static
*
* @apiDescription Create or update a behaviourCategoryAim
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/behaviourCategoryAim?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "behaviourCategoryAim": {
*       "id": "111",
*       "name": "title",
*       "name_slug": "nameSlug",
*       "icon_slug": "iconslug"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "behaviourCategoryAim": {
*       "id": "111",
*       "name": "title",
*       "name_slug": "nameSlug",
*       "icon_slug": "iconslug",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created BehaviourCategoryAim successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateBehaviourCategoryAim",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateBehaviourCategoryAim(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateBehaviourCategoryAim API request")
	req_behaviourCategoryAim := new(static_proto.CreateBehaviourCategoryAimRequest)
	err := req.ReadEntity(req_behaviourCategoryAim)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateBehaviourCategoryAim", "BindError")
		return
	}
	req_behaviourCategoryAim.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_behaviourCategoryAim.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateBehaviourCategoryAim(ctx, req_behaviourCategoryAim)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateBehaviourCategoryAim", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created BehaviourCategoryAim successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/behaviourCategoryAim/{behaviourCategoryAim_id}?session={session_id} View behaviourCategoryAim detail
* @apiVersion 0.1.0
* @apiName ReadBehaviourCategoryAim
* @apiGroup Static
*
* @apiDescription View behaviourCategoryAim detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/behaviourCategoryAim/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "behaviourCategoryAim": {
*       "id": "111",
*       "name": "title",
*       "name_slug": "nameSlug",
*       "icon_slug": "iconslug",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read BehaviourCategoryAim successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The BehaviourCategoryAim was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadBehaviourCategoryAim",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadBehaviourCategoryAim",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadBehaviourCategoryAim(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadBehaviourCategoryAim API request")
	req_behaviourCategoryAim := new(static_proto.ReadBehaviourCategoryAimRequest)
	req_behaviourCategoryAim.Id = req.PathParameter("behaviourCategoryAim_id")
	req_behaviourCategoryAim.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_behaviourCategoryAim.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadBehaviourCategoryAim(ctx, req_behaviourCategoryAim)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadBehaviourCategoryAim", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read BehaviourCategoryAims successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/behaviourCategoryAim/{behaviourCategoryAim_id}?session={session_id} Delete a behaviourCategoryAim
* @apiVersion 0.1.0
* @apiName DeleteBehaviourCategoryAim
* @apiGroup Static
*
* @apiDescription Delete a BehaviourCategoryAim
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/behaviourCategoryAim/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted BehaviourCategoryAim successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The BehaviourCategoryAim was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteBehaviourCategoryAim",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteBehaviourCategoryAim(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteBehaviourCategoryAim API request")
	req_behaviourCategoryAim := new(static_proto.DeleteBehaviourCategoryAimRequest)
	req_behaviourCategoryAim.Id = req.PathParameter("behaviourCategoryAim_id")
	req_behaviourCategoryAim.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_behaviourCategoryAim.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteBehaviourCategoryAim(ctx, req_behaviourCategoryAim)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteBehaviourCategoryAim", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted BehaviourCategoryAim successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/content/category/parents/all?session={session_id}&offset={offset}&limit={limit}  List all contentParentCategories
* @apiVersion 0.1.0
* @apiName AllContentParentCategories
* @apiGroup Static
*
* @apiDescription List all contentParentCategories
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/category/parents/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentParentCategories": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "org_id": "orgid",
*         "tags": ["tag1", "tag2"],
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all contentParentCategories successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentParentCategories were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllContentParentCategories",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllContentParentCategories",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllContentParentCategories(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllContentParentCategories API request")
	req_contentParentCategory := new(static_proto.AllContentParentCategoriesRequest)
	req_contentParentCategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentParentCategory.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_contentParentCategory.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_contentParentCategory.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_contentParentCategory.SortParameter = req.Attribute(SortParameter).(string)
	req_contentParentCategory.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllContentParentCategories(ctx, req_contentParentCategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllContentParentCategories", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all contentParentCategories successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/content/category/parent?session={session_id} Create or update a contentParentCategory
* @apiVersion 0.1.0
* @apiName CreateContentParentCategory
* @apiGroup Static
*
* @apiDescription Create or update a contentParentCategory
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/category/parent?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "contentParentCategory": {
*     "id": "111",
*     "name": "title",
*     "summary": "summary",
*     "description": "description",
*     "icon_slug": "iconslug",
*     "org_id": "orgid",
*     "tags": ["tag1", "tag2"]
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentParentCategory": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "org_id": "orgid",
*       "tags": ["tag1", "tag2"],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created contentParentCategory successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateContentParentCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateContentParentCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateContentParentCategory API request")
	req_contentParentCategory := new(static_proto.CreateContentParentCategoryRequest)
	err := req.ReadEntity(req_contentParentCategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateContentParentCategory", "BindError")
		return
	}
	req_contentParentCategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentParentCategory.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateContentParentCategory(ctx, req_contentParentCategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateContentParentCategory", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created contentParentCategory successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/content/category/parent/{contentParentCategory_id}?session={session_id} View contentParentCategory detail
* @apiVersion 0.1.0
* @apiName ReadContentParentCategory
* @apiGroup Static
*
* @apiDescription View contentParentCategory detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/category/parent/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentParentCategory": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "org_id": "orgid",
*       "tags": ["tag1", "tag2"],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read contentParentCategory successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentParentCategory was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadContentParentCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadContentParentCategory",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadContentParentCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadContentParentCategory API request")
	req_contentParentCategory := new(static_proto.ReadContentParentCategoryRequest)
	req_contentParentCategory.Id = req.PathParameter("contentParentCategory_id")
	req_contentParentCategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentParentCategory.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadContentParentCategory(ctx, req_contentParentCategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadContentParentCategory", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read contentParentCategory successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/content/category/parent/{contentParentCategory_id}?session={session_id} Delete a contentParentCategory
* @apiVersion 0.1.0
* @apiName DeleteContentParentCategory
* @apiGroup Static
*
* @apiDescription Delete a contentParentCategory
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/category/parent/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted contentParentCategory successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentParentCategory was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteContentParentCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteContentParentCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteContentParentCategory API request")
	req_contentParentCategory := new(static_proto.DeleteContentParentCategoryRequest)
	req_contentParentCategory.Id = req.PathParameter("contentParentCategory_id")
	req_contentParentCategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentParentCategory.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteContentParentCategory(ctx, req_contentParentCategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteContentParentCategory", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted contentParentCategory successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/content/categorys/all?session={session_id}&offset={offset}&limit={limit} List all contentCategories
* @apiVersion 0.1.0
* @apiName AllContentCategories
* @apiGroup Static
*
* @apiDescription List all contentCategories
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/categorys/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentCategories": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "description": "description",
*         "icon_slug": "iconslug",
*         "org_id": "orgid",
*         "tags": ["tag1", "tag2"],
*         "parent": [ ContentParentCategory, ... ],
*         "actions": [ Action, ... ],
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all contentCategories successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentCategories were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllContentCategories",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllContentCategories",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllContentCategories(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllContentCategories API request")
	req_contentCategory := new(static_proto.AllContentCategoriesRequest)
	req_contentCategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentCategory.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_contentCategory.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_contentCategory.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_contentCategory.SortParameter = req.Attribute(SortParameter).(string)
	req_contentCategory.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllContentCategories(ctx, req_contentCategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllContentCategories", "QueryError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read all contentCategories successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/content/category?session={session_id} Create or update a contentCategory
* @apiVersion 0.1.0
* @apiName CreateContentCategory
* @apiGroup Static
*
* @apiDescription Create or update a contentCategory
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/category?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "contentCategory": {
*      "id": "111",
*      "name": "title",
*      "summary": "summary",
*      "description": "description",
*      "icon_slug": "iconslug",
*      "org_id": "orgid",
*      "tags": ["tag1", "tag2"],
*      "parent": [ ContentParentCategory, ... ],
*      "actions": [ Action, ... ]
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentCategory": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "org_id": "orgid",
*       "tags": ["tag1", "tag2"],
*       "parent": [ ContentParentCategory, ... ],
*       "actions": [ Action, ... ],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created contentCategory successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateContentCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateContentCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateContentCategory API request")
	req_contentCategory := new(static_proto.CreateContentCategoryRequest)
	err := req.ReadEntity(req_contentCategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateContentCategory", "BindError")
		return
	}
	req_contentCategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentCategory.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateContentCategory(ctx, req_contentCategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateContentCategory", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created contentCategory successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/content/category/{contentCategory_id}?session={session_id} View contentCategory detail
* @apiVersion 0.1.0
* @apiName ReadContentCategory
* @apiGroup Static
*
* @apiDescription View contentCategory detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/category/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentCategory": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "description": "description",
*       "icon_slug": "iconslug",
*       "org_id": "orgid",
*       "tags": ["tag1", "tag2"],
*       "parent": [ ContentParentCategory, ... ],
*       "actions": [ Action, ... ],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read contentCategory successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentCategory was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadContentCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadContentCategory",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadContentCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadContentCategory API request")
	req_contentCategory := new(static_proto.ReadContentCategoryRequest)
	req_contentCategory.Id = req.PathParameter("contentCategory_id")
	req_contentCategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentCategory.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadContentCategory(ctx, req_contentCategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadContentCategory", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read contentCategory successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/content/category/{contentCategory_id}?session={session_id} Delete a contentCategory
* @apiVersion 0.1.0
* @apiName DeleteContentCategory
* @apiGroup Static
*
* @apiDescription Delete a contentCategory
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/category/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted contentCategory successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentCategory was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteContentCategory",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteContentCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteContentCategory API request")
	req_contentCategory := new(static_proto.DeleteContentCategoryRequest)
	req_contentCategory.Id = req.PathParameter("contentCategory_id")
	req_contentCategory.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentCategory.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteContentCategory(ctx, req_contentCategory)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteContentCategory", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted contentCategory successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/content/types/all?session={session_id}&offset={offset}&limit={limit}  List all contentTypes
* @apiVersion 0.1.0
* @apiName AllContentTypes
* @apiGroup Static
*
* @apiDescription List all contentTypes
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/types/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentTypes": [
*       {
*         "id": "111",
*         "name": "title",
*         "description": "description",
*         "contentTypeString": "contentTypeString",
*         "tags": ["tag1", "tag2"],
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all contentTypes successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentTypes were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllContentTypes",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllContentTypes",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllContentTypes(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllContentTypes API request")
	req_contentType := new(static_proto.AllContentTypesRequest)
	req_contentType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentType.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_contentType.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_contentType.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_contentType.SortParameter = req.Attribute(SortParameter).(string)
	req_contentType.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllContentTypes(ctx, req_contentType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllContentTypes", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all contentTypes successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/content/type?session={session_id} Create or update a contentType
* @apiVersion 0.1.0
* @apiName CreateContentType
* @apiGroup Static
*
* @apiDescription Create or update a contentType
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/type?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "contentType": {
*     "id": "111",
*     "name": "title",
*     "description": "description",
*     "contentTypeString": "contentTypeString",
*     "tags": ["tag1", "tag2"]
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentType": {
*       "id": "111",
*       "name": "title",
*       "description": "description",
*       "contentTypeString": "contentTypeString",
*       "tags": ["tag1", "tag2"],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created contentType successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateContentType",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateContentType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateContentType API request")
	req_contentType := new(static_proto.CreateContentTypeRequest)
	err := req.ReadEntity(req_contentType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateContentType", "BindError")
		return
	}
	req_contentType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentType.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateContentType(ctx, req_contentType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateContentType", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created contentType successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/content/type/{contentType_id}?session={session_id} View contentType detail
* @apiVersion 0.1.0
* @apiName ReadContentType
* @apiGroup Static
*
* @apiDescription View contentType detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/type/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentType": {
*       "id": "111",
*       "name": "title",
*       "description": "description",
*       "contentTypeString": "contentTypeString",
*       "tags": ["tag1", "tag2"],
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read contentType successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentType was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadContentType",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadContentType",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadContentType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadContentType API request")
	req_contentType := new(static_proto.ReadContentTypeRequest)
	req_contentType.Id = req.PathParameter("contentType_id")
	req_contentType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentType.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadContentType(ctx, req_contentType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadContentType", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read contentType successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/content/type/{contentType_id}?session={session_id} Delete a contentType
* @apiVersion 0.1.0
* @apiName DeleteContentType
* @apiGroup Static
*
* @apiDescription Delete a contentType
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/type/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted contentType successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentType was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteContentType",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteContentType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteContentType API request")
	req_contentType := new(static_proto.DeleteContentTypeRequest)
	req_contentType.Id = req.PathParameter("contentType_id")
	req_contentType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentType.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteContentType(ctx, req_contentType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteContentType", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted contentType successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/content/source/types/all?session={session_id}&offset={offset}&limit={limit}  List all contentSourceTypes
* @apiVersion 0.1.0
* @apiName AllContentSourceTypes
* @apiGroup Static
*
* @apiDescription List all contentSourceTypes
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/source/types/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentSourceTypes": [
*       {
*         "id": "111",
*         "name": "title",
*         "description": "description",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all contentSourceTypes successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentSourceTypes were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllContentSourceTypes",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllContentSourceTypes",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllContentSourceTypes(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllContentSourceTypes API request")
	req_contentSourceType := new(static_proto.AllContentSourceTypesRequest)
	req_contentSourceType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentSourceType.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_contentSourceType.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_contentSourceType.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_contentSourceType.SortParameter = req.Attribute(SortParameter).(string)
	req_contentSourceType.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllContentSourceTypes(ctx, req_contentSourceType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllContentSourceTypes", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all contentSourceTypes successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/content/source/type?session={session_id} Create or update a contentSourceType
* @apiVersion 0.1.0
* @apiName CreateContentSourceType
* @apiGroup Static
*
* @apiDescription Create or update a contentSourceType
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/source/type?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "contentSourceType": {
*     "id": "111",
*     "name": "title",
*     "description": "description"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentSourceType": {
*       "id": "111",
*       "name": "title",
*       "description": "description",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created contentSourceType successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateContentSourceType",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateContentSourceType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateContentSourceType API request")
	req_contentSourceType := new(static_proto.CreateContentSourceTypeRequest)
	err := req.ReadEntity(req_contentSourceType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateContentSourceType", "BindError")
		return
	}
	req_contentSourceType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentSourceType.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateContentSourceType(ctx, req_contentSourceType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateContentSourceType", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created contentSourceType successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/content/source/type/{contentSourceType_id}?session={session_id} View contentSourceType detail
* @apiVersion 0.1.0
* @apiName ReadContentSourceType
* @apiGroup Static
*
* @apiDescription View contentSourceType detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/source/type/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "contentSourceType": {
*       "id": "111",
*       "name": "title",
*       "description": "description",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read contentSourceType successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentSourceType was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadContentSourceType",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadContentSourceType",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadContentSourceType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadContentSourceType API request")
	req_contentSourceType := new(static_proto.ReadContentSourceTypeRequest)
	req_contentSourceType.Id = req.PathParameter("contentSourceType_id")
	req_contentSourceType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentSourceType.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadContentSourceType(ctx, req_contentSourceType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadContentSourceType", "ReadError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Read contentSourceType successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/content/source/type/{contentSourceType_id}?session={session_id} Delete a contentSourceType
* @apiVersion 0.1.0
* @apiName DeleteContentSourceType
* @apiGroup Static
*
* @apiDescription Delete a contentSourceType
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/content/source/type/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted contentSourceType successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The contentSourceType was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteContentSourceType",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteContentSourceType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteContentSourceType API request")
	req_contentSourceType := new(static_proto.DeleteContentSourceTypeRequest)
	req_contentSourceType.Id = req.PathParameter("contentSourceType_id")
	req_contentSourceType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_contentSourceType.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteContentSourceType(ctx, req_contentSourceType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteContentSourceType", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted contentSourceType successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/module/triggers/all?session={session_id}&offset={offset}&limit={limit}  List all moduleTriggers
* @apiVersion 0.1.0
* @apiName AllModuleTriggers
* @apiGroup Static
*
* @apiDescription List all moduleTriggers
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/module/triggers/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "moduleTriggers": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "icon_slug": "iconSlug",
*         "type": "EVENT",
*         "module": { Module },
*         "events": [{TriggerEvent}, ...],
*         "duration": ["P1Y2M3DT4H5M6S",...],
*         "delay": "",
*         "contentTypes": [{TriggerContentType}, ...],
*         "actionMethod": "CHAT",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all moduleTriggers successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The moduleTriggers were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllModuleTriggers",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllModuleTriggers",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllModuleTriggers(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllModuleTriggers API request")
	req_moduleTrigger := new(static_proto.AllModuleTriggersRequest)
	req_moduleTrigger.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_moduleTrigger.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_moduleTrigger.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_moduleTrigger.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_moduleTrigger.SortParameter = req.Attribute(SortParameter).(string)
	req_moduleTrigger.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllModuleTriggers(ctx, req_moduleTrigger)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllModuleTriggers", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all moduleTriggers successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {post} /server/static/module/trigger?session={session_id} Create or update a moduleTrigger
* @apiVersion 0.1.0
* @apiName CreateModuleTrigger
* @apiGroup Static
*
* @apiDescription Create or update a moduleTrigger
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/module/trigger?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "moduleTrigger": {
*     "id": "111",
*     "name": "title",
*     "summary": "summary",
*     "icon_slug": "iconSlug",
*     "type": "EVENT",
*     "module": { Module },
*     "events": [{TriggerEvent}, ...],
*     "duration": ["P1Y2M3DT4H5M6S",...],
*     "delay": "",
*     "contentTypes": [{TriggerContentType}, ...],
*     "actionMethod": "CHAT"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "moduleTrigger": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "icon_slug": "iconSlug",
*       "type": "EVENT",
*       "module": { Module },
*       "events": [{TriggerEvent}, ...],
*       "duration": ["P1Y2M3DT4H5M6S",...],
*       "delay": "",
*       "contentTypes": [{TriggerContentType}, ...],
*       "actionMethod": "CHAT",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created moduleTrigger successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateModuleTrigger",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateModuleTrigger(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateModuleTrigger API request")
	req_moduleTrigger := new(static_proto.CreateModuleTriggerRequest)
	// err := req.ReadEntity(req_moduleTrigger)
	// if err != nil {
	// 	utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateModuleTrigger", "BindError")
	// 	return
	// }
	if err := utils.UnmarshalAny(req, rsp, req_moduleTrigger); err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateModuleTrigger", "BindError")
		return
	}
	req_moduleTrigger.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_moduleTrigger.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateModuleTrigger(ctx, req_moduleTrigger)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateModuleTrigger", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created moduleTrigger successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/static/module/trigger/{moduleTrigger_id}?session={session_id} View moduleTrigger detail
* @apiVersion 0.1.0
* @apiName ReadModuleTrigger
* @apiGroup Static
*
* @apiDescription View moduleTrigger detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/module/trigger/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "moduleTrigger": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "icon_slug": "iconSlug",
*       "type": "EVENT",
*       "module": { Module },
*       "events": [{TriggerEvent}, ...],
*       "duration": ["P1Y2M3DT4H5M6S",...],
*       "delay": "",
*       "contentTypes": [{TriggerContentType}, ...],
*       "actionMethod": "CHAT",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read moduleTrigger successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The moduleTrigger was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadModuleTrigger",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadModuleTrigger",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadModuleTrigger(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadModuleTrigger API request")
	req_moduleTrigger := new(static_proto.ReadModuleTriggerRequest)
	req_moduleTrigger.Id = req.PathParameter("moduleTrigger_id")
	req_moduleTrigger.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_moduleTrigger.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadModuleTrigger(ctx, req_moduleTrigger)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadModuleTrigger", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read moduleTrigger successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {delete} /server/static/module/trigger/{moduleTrigger_id}?session={session_id} Delete a moduleTrigger
* @apiVersion 0.1.0
* @apiName DeleteModuleTrigger
* @apiGroup Static
*
* @apiDescription Delete a moduleTrigger
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/module/trigger/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted moduleTrigger successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The moduleTrigger was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteModuleTrigger",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteModuleTrigger(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteModuleTrigger API request")
	req_moduleTrigger := new(static_proto.DeleteModuleTriggerRequest)
	req_moduleTrigger.Id = req.PathParameter("moduleTrigger_id")
	req_moduleTrigger.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_moduleTrigger.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteModuleTrigger(ctx, req_moduleTrigger)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteModuleTrigger", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted moduleTrigger successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/module/triggers/filter?session={session_id}&offset={offset}&limit={limit} List all moduleTriggers
* @apiVersion 0.1.0
* @apiName AllModuleTriggers
* @apiGroup Static
*
* @apiDescription List all moduleTriggers
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/module/triggers/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "module": ["111"]
*   "triggerType": [1]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "moduleTriggers": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "icon_slug": "iconSlug",
*         "type": 1,
*         "module": { "id": 111, ... },
*         "triggerContentType": { TriggerContentType },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all moduleTriggers successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The moduleTriggers were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllModuleTriggers",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllModuleTriggers",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) FilterModuleTrigger(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.FilterModuleTrigger API request")
	req_moduleTrigger := new(static_proto.FilterModuleTriggerRequest)
	err := req.ReadEntity(req_moduleTrigger)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.FilterModuleTrigger", "BindError")
		return
	}
	req_moduleTrigger.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_moduleTrigger.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_moduleTrigger.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_moduleTrigger.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_moduleTrigger.SortParameter = req.Attribute(SortParameter).(string)
	req_moduleTrigger.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.FilterModuleTrigger(ctx, req_moduleTrigger)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.FilterModuleTrigger", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all moduleTriggers successfully"
	data := utils.MarshalAny(rsp, resp)

	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, data)
}

/**
* @api {get} /server/static/trigger/content/types/all?session={session_id}&offset={offset}&limit={limit} List all TriggerContentTypes
* @apiVersion 0.1.0
* @apiName AllTriggerContentTypes
* @apiGroup Static
*
* @apiDescription List all triggerContentTypes
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/trigger/content/types/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "triggerContentTypes": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "icon_slug": "iconSlug",
*         "type": 1,
*         "module": { Module },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all triggerContentTypes successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The triggerContentTypes were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllTriggerContentTypes",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllTriggerContentTypes",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllTriggerContentTypes(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllTriggerContentTypes API request")
	req_triggerContentType := new(static_proto.AllTriggerContentTypesRequest)
	req_triggerContentType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_triggerContentType.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_triggerContentType.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_triggerContentType.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_triggerContentType.SortParameter = req.Attribute(SortParameter).(string)
	req_triggerContentType.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllTriggerContentTypes(ctx, req_triggerContentType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllTriggerContentTypes", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all triggerContentTypes successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/trigger/content/type?session={session_id} Create or update a TriggerContentType
* @apiVersion 0.1.0
* @apiName CreateTriggerContentType
* @apiGroup Static
*
* @apiDescription Create or update a triggerContentType
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/trigger/content/type?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "triggerContentType": {
*     "id": "111",
*     "name": "title",
*     "summary": "summary",
*     "icon_slug": "iconSlug",
*     "type": 1,
*     "module": { Module }
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "triggerContentType": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "icon_slug": "iconSlug",
*       "type": 1,
*       "module": { Module },
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created triggerContentType successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateTriggerContentType",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateTriggerContentType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateTriggerContentType API request")
	req_triggerContentType := new(static_proto.CreateTriggerContentTypeRequest)
	err := req.ReadEntity(req_triggerContentType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateTriggerContentType", "BindError")
		return
	}
	req_triggerContentType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_triggerContentType.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateTriggerContentType(ctx, req_triggerContentType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateTriggerContentType", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created triggerContentType successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/trigger/content/type/{triggerContentType_id}?session={session_id} View TriggerContentType detail
* @apiVersion 0.1.0
* @apiName ReadTriggerContentType
* @apiGroup Static
*
* @apiDescription View triggerContentType detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/trigger/content/type/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "triggerContentType": {
*       "id": "111",
*       "name": "title",
*       "summary": "summary",
*       "icon_slug": "iconSlug",
*       "type": 1,
*       "module": { Module },
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read triggerContentType successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The triggerContentType was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadTriggerContentType",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadTriggerContentType",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadTriggerContentType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadTriggerContentType API request")
	req_triggerContentType := new(static_proto.ReadTriggerContentTypeRequest)
	req_triggerContentType.Id = req.PathParameter("triggerContentType_id")
	req_triggerContentType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_triggerContentType.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadTriggerContentType(ctx, req_triggerContentType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadTriggerContentType", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read triggerContentType successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/trigger/content/type/{triggerContentType_id}?session={session_id} Delete a TriggerContentType
* @apiVersion 0.1.0
* @apiName DeleteTriggerContentType
* @apiGroup Static
*
* @apiDescription Delete a triggerContentType
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/trigger/content/type/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted triggerContentType successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The triggerContentType was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteTriggerContentType",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteTriggerContentType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteTriggerContentType API request")
	req_triggerContentType := new(static_proto.DeleteTriggerContentTypeRequest)
	req_triggerContentType.Id = req.PathParameter("triggerContentType_id")
	req_triggerContentType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_triggerContentType.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteTriggerContentType(ctx, req_triggerContentType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteTriggerContentType", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted triggerContentType successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/trigger/content/types/filter?session={session_id}&offset={offset}&limit={limit} List all TriggerContentTypes
* @apiVersion 0.1.0
* @apiName AllTriggerContentTypes
* @apiGroup Static
*
* @apiDescription List all triggerContentTypes
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/trigger/content/types/filter?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiParamExample {json} Request-Example:
* {
*   "module": ["111"]
*   "triggerType": [1]
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "triggerContentTypes": [
*       {
*         "id": "111",
*         "name": "title",
*         "summary": "summary",
*         "icon_slug": "iconSlug",
*         "type": 1,
*         "module": { "id": 111, ... },
*         "triggerContentType": { TriggerContentType },
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all triggerContentTypes successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The triggerContentTypes were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllTriggerContentTypes",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllTriggerContentTypes",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) FilterTriggerContentType(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.FilterTriggerContentType API request")
	req_triggerContentType := new(static_proto.FilterTriggerContentTypeRequest)
	err := req.ReadEntity(req_triggerContentType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.FilterTriggerContentType", "BindError")
		return
	}
	req_triggerContentType.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_triggerContentType.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_triggerContentType.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_triggerContentType.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_triggerContentType.SortParameter = req.Attribute(SortParameter).(string)
	req_triggerContentType.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.FilterTriggerContentType(ctx, req_triggerContentType)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.FilterTriggerContentType", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all triggerContentTypes successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/setbacks/all?session={session_id}&offset={offset}&limit={limit}  List all setbacks
* @apiVersion 0.1.0
* @apiName AllSetbacks
* @apiGroup Static
*
* @apiDescription List all setbacks
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/setbacks/all?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="&offset=0&limit=10
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "setbacks": [
*       {
*         "id": "111",
*         "name": "title",
*         "description": "description",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Read all setbacks successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The setbacks were not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllSetbacks",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
* @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AllSetbacks",
*           "reason": "NotFound"
*         }
*       ]
*     }
 */
func (p *StaticService) AllSetbacks(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AllSetbacks API request")
	req_setback := new(static_proto.AllSetbacksRequest)
	req_setback.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_setback.TeamId = req.Attribute(TeamIdAttrName).(string)
	req_setback.Limit = req.Attribute(PaginateLimitParameter).(int64)
	req_setback.Offset = req.Attribute(PaginateOffsetParameter).(int64)
	req_setback.SortParameter = req.Attribute(SortParameter).(string)
	req_setback.SortDirection = req.Attribute(SortDirection).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AllSetbacks(ctx, req_setback)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AllSetbacks", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read all setbacks successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/setback?session={session_id} Create or update a setback
* @apiVersion 0.1.0
* @apiName CreateSetback
* @apiGroup Static
*
* @apiDescription Create or update a setback
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/setback?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "setback": {
*     "id": "111",
*     "name": "title",
*     "description": "description"
*   }
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "setback": {
*       "id": "111",
*       "name": "title",
*       "description": "description",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Created setback successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, CreateError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "CreateError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.CreateSetback",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) CreateSetback(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.CreateSetback API request")
	req_setback := new(static_proto.CreateSetbackRequest)
	err := req.ReadEntity(req_setback)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateSetback", "BindError")
		return
	}
	req_setback.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_setback.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.CreateSetback(ctx, req_setback)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.CreateSetback", "CreateError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Created setback successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/static/setback/{setback_id}?session={session_id} View setback detail
* @apiVersion 0.1.0
* @apiName ReadSetback
* @apiGroup Static
*
* @apiDescription View setback detail
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/setback/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "setback": {
*       "id": "111",
*       "name": "title",
*       "description": "description",
*       "created": 1517891917,
*       "updated": 1517891917
*     }
*   },
*   "code": 200,
*   "message": "Read setback successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The setback was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "ReadError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadSetback",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 * @apiErrorExample Not-Found:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "NotFound",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.ReadSetback",
*           "reason": "NotFound"
*         }
*       ]
*     }
*/
func (p *StaticService) ReadSetback(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.ReadSetback API request")
	req_setback := new(static_proto.ReadSetbackRequest)
	req_setback.Id = req.PathParameter("setback_id")
	req_setback.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_setback.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.ReadSetback(ctx, req_setback)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.ReadSetback", "ReadError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Read setback successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {delete} /server/static/setback/{setback_id}?session={session_id} Delete a setback
* @apiVersion 0.1.0
* @apiName DeleteSetback
* @apiGroup Static
*
* @apiDescription Delete a setback
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/setback/111?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Deleted setback successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	The setback was not found.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "DeleteError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.DeleteSetback",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) DeleteSetback(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.DeleteSetback API request")
	req_setback := new(static_proto.DeleteSetbackRequest)
	req_setback.Id = req.PathParameter("setback_id")
	req_setback.OrgId = req.Attribute(OrgIdAttrName).(string)
	req_setback.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.DeleteSetback(ctx, req_setback)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.DeleteSetback", "DeleteError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Deleted setback successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/static/setback/search/autocomplete?session={session_id} autocomplete text search for setbacks
* @apiVersion 0.1.0
* @apiName AutocompleteSetbackSearch
* @apiGroup Static
*
* @apiDescription It should use setback.name to filter and return list of setbacks
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/static/setback/search/autocomplete?session="qLN5aNAiway8h6pTPEEaz2YIdayRrdTMsoarQdbkulQ="
*
* @apiParamExample {json} Request-Example:
* {
*   "title": "p",
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "setbacks": [
*       {
*         "id": "111",
*         "name": "p1",
*         "description": "description",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       {
*         "id": "222",
*         "name": "p2",
*         "description": "description",
*         "created": 1517891917,
*         "updated": 1517891917
*       },
*       ... ...
*     ]
*   },
*   "code": 200,
*   "message": "Autocomplete setback search successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, SearchError.
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "QueryError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.static.AutocompleteSetbackSearch",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */
func (p *StaticService) AutocompleteSetbackSearch(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.AutocompleteSetbackSearch API request")

	req_search := new(static_proto.AutocompleteSetbackSearchRequest)
	err := req.ReadEntity(req_search)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AutocompleteSetbackSearch", "BindError")
		return
	}
	// req_search.OrgId = req.Attribute(OrgIdAttrName).(string)
	// req_search.TeamId = req.Attribute(TeamIdAttrName).(string)

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.StaticClient.AutocompleteSetbackSearch(ctx, req_search)
	if err != nil {
		utils.WriteErrorResponse(rsp, err, "go.micro.srv.static.AutocompleteSetbackSearch", "SearchError")
		return
	}
	resp.Code = http.StatusOK
	resp.Message = "Autocomplete setback search successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

func (p *StaticService) UploadBehaviourCategoryAim(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.UploadBehaviourCategoryAim API request")
	req.Request.ParseMultipartForm(32 << 20)
	file, _, err := req.Request.FormFile("upload_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	aims := []*static_proto.BehaviourCategoryAim{}
	if err := gocsv.Unmarshal(file, &aims); err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.static.UploadBehaviourCategoryAim", "File marshale errro")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	for _, aim := range aims {
		_, err = p.StaticClient.CreateBehaviourCategoryAim(ctx, &static_proto.CreateBehaviourCategoryAimRequest{BehaviourCategoryAim: aim})
		if err != nil {
			utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.static.UploadBehaviourCategoryAim", "CreadError")
			return
		}
	}

	resp := &static_proto.UploadResponse{
		Code:    http.StatusOK,
		Message: "Upload csv successfully",
	}
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

func (p *StaticService) UploadContentCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.UploadContentCategory API request")

	req.Request.ParseMultipartForm(32 << 20)
	file, _, err := req.Request.FormFile("upload_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	req_category := &static_proto.CreateContentCategoryRequest{
		OrgId: req.Attribute(OrgIdAttrName).(string),
	}

	r := csv.NewReader(file)
	fields := map[int]string{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		// property parsing
		if len(fields) == 0 {
			for i, col := range row {
				fields[i] = col
			}
			continue
		}
		// row parsing
		for i, col := range row {
			category := &static_proto.ContentCategory{
				OrgId:          req_category.OrgId,
				Tags:           []string{},
				Actions:        []*static_proto.Action{},
				TrackerMethods: []*static_proto.TrackerMethod{},
			}
			switch fields[i] {
			case "name":
				category.Name = col
			case "category":
				category.Summary = col
			case "description":
				category.Description = col
			case "tags":
				category.Tags = append(category.Tags, col)
			case "tracker_method.name_slug":
				resp, err := p.StaticClient.ReadTrackerMethodByNameslug(ctx, &static_proto.ReadByNameslugRequest{col})
				if err != nil {
					continue
				}
				category.TrackerMethods = append(category.TrackerMethods, resp.Data.TrackerMethod)
			}

			req_category.ContentCategory = category
			_, err := p.StaticClient.CreateContentCategory(ctx, &static_proto.CreateContentCategoryRequest{ContentCategory: category})
			if err != nil {
				utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.behaviour.UploadContentCategory", "CreateError")
				return
			}
		}
	}

	resp := &static_proto.UploadResponse{
		Code:    http.StatusOK,
		Message: "Upload csv successfully",
	}
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

func (p *StaticService) UploadMarker(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.UploadMarker API request")

	req.Request.ParseMultipartForm(32 << 20)
	file, _, err := req.Request.FormFile("upload_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// markers := []*static_proto.Marker{}
	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	req_marker := &static_proto.CreateMarkerRequest{
		OrgId: req.Attribute(OrgIdAttrName).(string),
	}

	r := csv.NewReader(file)
	fields := map[int]string{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		// property parsing
		if len(fields) == 0 {
			for i, col := range row {
				fields[i] = col
			}
			continue
		}
		// row parsing
		for i, col := range row {
			marker := &static_proto.Marker{
				OrgId:          req_marker.OrgId,
				Apps:           []*static_proto.App{},
				Wearables:      []*static_proto.Wearable{},
				Devices:        []*static_proto.Device{},
				TrackerMethods: []*static_proto.TrackerMethod{},
			}
			switch fields[i] {
			case "name":
				marker.Name = col
			case "summary":
				marker.Summary = col
			case "description":
				marker.Description = col
			case "unit":
				marker.Unit = append(marker.Unit, col)
			case "app.name_slug":
				resp, err := p.StaticClient.ReadAppByNameslug(ctx, &static_proto.ReadByNameslugRequest{col})
				if err != nil {
					continue
				}
				marker.Apps = append(marker.Apps, resp.Data.App)
			case "wearable.name_slug":
				resp, err := p.StaticClient.ReadWearableByNameslug(ctx, &static_proto.ReadByNameslugRequest{col})
				if err != nil {
					continue
				}
				marker.Wearables = append(marker.Wearables, resp.Data.Wearable)
			case "device.name_slug":
				resp, err := p.StaticClient.ReadDeviceByNameslug(ctx, &static_proto.ReadByNameslugRequest{col})
				if err != nil {
					continue
				}
				marker.Devices = append(marker.Devices, resp.Data.Device)
			case "method.name_slug":
				resp, err := p.StaticClient.ReadTrackerMethodByNameslug(ctx, &static_proto.ReadByNameslugRequest{col})
				if err != nil {
					continue
				}
				marker.TrackerMethods = append(marker.TrackerMethods, resp.Data.TrackerMethod)
			}

			req_marker.Marker = marker
			_, err := p.StaticClient.CreateMarker(ctx, req_marker)
			if err != nil {
				utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.markers.UploadMarker", "CreateMarkerError")
				return
			}
		}
	}

	resp := &static_proto.UploadResponse{
		Code:    http.StatusOK,
		Message: "Upload csv successfully",
	}
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

func (p *StaticService) UploadBehaviourCategory(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.UploadBehaviourCategoryAim API request")

	req.Request.ParseMultipartForm(32 << 20)
	file, _, err := req.Request.FormFile("upload_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	req_category := &static_proto.CreateBehaviourCategoryRequest{
		OrgId: req.Attribute(OrgIdAttrName).(string),
	}

	// categorys := []*static_proto.BehaviourCategory{}
	r := csv.NewReader(file)
	fields := map[int]string{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		// property parsing
		if len(fields) == 0 {
			for i, col := range row {
				fields[i] = col
			}
			continue
		}
		// row parsing
		for i, col := range row {
			category := &static_proto.BehaviourCategory{
				OrgId:         req_category.OrgId,
				Tags:          []string{},
				Aims:          []*static_proto.BehaviourCategoryAim{},
				MarkerDefault: &static_proto.Marker{},
				MarkerOptions: []*static_proto.Marker{},
			}
			switch fields[i] {
			case "name":
				category.Name = col
			case "summary":
				category.Summary = col
			case "description":
				category.Description = col
			case "tags":
				category.Tags = append(category.Tags, col)
			case "aim.name_slug":
				resp, err := p.StaticClient.ReadBehaviourCategoryAimByNameslug(ctx, &static_proto.ReadByNameslugRequest{col})
				if err != nil {
					continue
				}
				category.Aims = append(category.Aims, resp.Data.BehaviourCategoryAim)
			case "default.name_slug":
				resp, err := p.StaticClient.ReadMarkerByNameslug(ctx, &static_proto.ReadByNameslugRequest{col})
				if err != nil {
					continue
				}
				category.MarkerDefault = resp.Data.Marker
			case "options.name_slug":
				resp, err := p.StaticClient.ReadMarkerByNameslug(ctx, &static_proto.ReadByNameslugRequest{col})
				if err != nil {
					continue
				}
				category.MarkerOptions = append(category.MarkerOptions, resp.Data.Marker)
			}

			req_category.Category = category
			_, err = p.StaticClient.CreateBehaviourCategory(ctx, req_category)
			if err != nil {
				utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.markers.UploadBehaviourCategory", "CreateMarkerError")
				return
			}
		}
	}

	resp := &static_proto.UploadResponse{
		Code:    http.StatusOK,
		Message: "Upload csv successfully",
	}
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)

}

func (p *StaticService) UploadContentCategoryItem(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.UploadContentCategoryItem API request")

	req.Request.ParseMultipartForm(32 << 20)
	file, _, err := req.Request.FormFile("upload_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	req_item := &content_proto.CreateContentCategoryItemRequest{
		OrgId: req.Attribute(OrgIdAttrName).(string),
	}

	// items := []*static_proto.ContentCategoryItem{}
	r := csv.NewReader(file)
	fields := map[int]string{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		// property parsing
		if len(fields) == 0 {
			for i, col := range row {
				fields[i] = col
			}
			continue
		}
		// row parsing
		for i, col := range row {
			item := &static_proto.ContentCategoryItem{
				OrgId: req_item.OrgId,
				Tags:  []string{},
			}
			switch fields[i] {
			case "name":
				item.Name = col
			case "summary":
				item.Summary = col
			case "description":
				item.Description = col
			case "tags":
				item.Tags = append(item.Tags, col)
			case "taxonomy":
				resp, err := p.ContentClient.ReadTaxonomyByNameslug(ctx, &static_proto.ReadByNameslugRequest{col})
				if err != nil {
					continue
				}
				item.Taxonomy = resp.Data.Taxonomy
			case "weight":
				v, _ := strconv.Atoi(col)
				item.Weight = int64(v)
			case "priority":
				v, _ := strconv.Atoi(col)
				item.Priority = int64(v)
			case "category":
				resp, err := p.StaticClient.ReadContentCategoryByNameslug(ctx, &static_proto.ReadByNameslugRequest{col})
				if err != nil {
					continue
				}
				item.Category = resp.Data.ContentCategory
			}

			req_item.ContentCategoryItem = item
			_, err := p.ContentClient.CreateContentCategoryItem(ctx, req_item)
			if err != nil {
				utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.static.UploadContentCategoryItem", "CreateMarkerError")
				return
			}
		}
	}
	resp := &static_proto.UploadResponse{
		Code:    http.StatusOK,
		Message: "Upload csv successfully",
	}
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

func (p *StaticService) UploadTrackerMethod(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Static.UploadTrackerMethod API request")

	req.Request.ParseMultipartForm(32 << 20)
	file, _, err := req.Request.FormFile("upload_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	methods := []*static_proto.TrackerMethod{}
	if err := gocsv.Unmarshal(file, &methods); err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.static.UploadTrackerMethod", "Fild marshale errro")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	for _, method := range methods {
		_, err = p.StaticClient.CreateTrackerMethod(ctx, &static_proto.CreateTrackerMethodRequest{TrackerMethod: method})
		if err != nil {
			utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.static.UploadTrackerMethod", "CreateTrackerMethod")
			return
		}
	}

	resp := &static_proto.UploadResponse{
		Code:    http.StatusOK,
		Message: "Upload csv successfully",
	}
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}
