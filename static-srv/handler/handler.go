package handler

import (
	"context"
	"server/common"
	"server/static-srv/db"
	static_proto "server/static-srv/proto/static"

	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type StaticService struct{}

func (p *StaticService) AllApps(ctx context.Context, req *static_proto.AllAppsRequest, rsp *static_proto.AllAppsResponse) error {
	log.Info("Received Static.AllApps request")
	apps, err := db.AllApps(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(apps) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllApps, err, "not found")
	}
	rsp.Data = &static_proto.AppArrData{apps}
	return nil
}

func (p *StaticService) CreateApp(ctx context.Context, req *static_proto.CreateAppRequest, rsp *static_proto.CreateAppResponse) error {
	log.Info("Received Static.CreateApp request")
	if len(req.App.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateApp, nil, "app name empty")
	}
	if len(req.App.Id) == 0 {
		req.App.Id = uuid.NewUUID().String()
	}

	err := db.CreateApp(ctx, req.App)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateApp, err, "create error")
	}
	rsp.Data = &static_proto.AppData{req.App}
	return nil
}

func (p *StaticService) ReadApp(ctx context.Context, req *static_proto.ReadAppRequest, rsp *static_proto.ReadAppResponse) error {
	log.Info("Received Static.ReadApp request")
	app, err := db.ReadApp(ctx, req.Id, req.OrgId, req.TeamId)
	if app == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadApp, err, "not found")
	}
	rsp.Data = &static_proto.AppData{app}
	return nil
}

func (p *StaticService) DeleteApp(ctx context.Context, req *static_proto.DeleteAppRequest, rsp *static_proto.DeleteAppResponse) error {
	log.Info("Received Static.DeleteApp request")
	if err := db.DeleteApp(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteApp, nil, "delete error")
	}
	return nil
}

func (p *StaticService) AllPlatforms(ctx context.Context, req *static_proto.AllPlatformsRequest, rsp *static_proto.AllPlatformsResponse) error {
	log.Info("Received Static.AllPlatforms request")
	platforms, err := db.AllPlatforms(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(platforms) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllPlatforms, err, "platform not found")
	}
	rsp.Data = &static_proto.PlatformArrData{platforms}
	return nil
}

func (p *StaticService) CreatePlatform(ctx context.Context, req *static_proto.CreatePlatformRequest, rsp *static_proto.CreatePlatformResponse) error {
	log.Info("Received Static.CreatePlatform request")
	if len(req.Platform.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreatePlatform, nil, "platform name empty")
	}
	if len(req.Platform.Id) == 0 {
		req.Platform.Id = uuid.NewUUID().String()
	}

	err := db.CreatePlatform(ctx, req.Platform)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreatePlatform, err, "platform create error")
	}
	rsp.Data = &static_proto.PlatformData{req.Platform}
	return nil
}

func (p *StaticService) ReadPlatform(ctx context.Context, req *static_proto.ReadPlatformRequest, rsp *static_proto.ReadPlatformResponse) error {
	log.Info("Received Static.ReadPlatform request")
	platform, err := db.ReadPlatform(ctx, req.Id, req.OrgId, req.TeamId)
	if platform == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadPlatform, err, "platform not found")
	}
	rsp.Data = &static_proto.PlatformData{platform}
	return nil
}

func (p *StaticService) DeletePlatform(ctx context.Context, req *static_proto.DeletePlatformRequest, rsp *static_proto.DeletePlatformResponse) error {
	log.Info("Received Static.DeletePlatform request")
	if err := db.DeletePlatform(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeletePlatform, nil, "platform delete error")
	}
	return nil
}

func (p *StaticService) AllDevices(ctx context.Context, req *static_proto.AllDevicesRequest, rsp *static_proto.AllDevicesResponse) error {
	log.Info("Received Static.AllDevices request")
	devices, err := db.AllDevices(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(devices) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllDevices, err, "device not found")
	}
	rsp.Data = &static_proto.DeviceArrData{devices}
	return nil
}

func (p *StaticService) CreateDevice(ctx context.Context, req *static_proto.CreateDeviceRequest, rsp *static_proto.CreateDeviceResponse) error {
	log.Info("Received Static.CreateDevice request")
	if len(req.Device.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateDevice, nil, "device name empty")
	}
	if len(req.Device.Id) == 0 {
		req.Device.Id = uuid.NewUUID().String()
	}

	err := db.CreateDevice(ctx, req.Device)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateDevice, err, "create error")
	}
	rsp.Data = &static_proto.DeviceData{req.Device}
	return nil
}

func (p *StaticService) ReadDevice(ctx context.Context, req *static_proto.ReadDeviceRequest, rsp *static_proto.ReadDeviceResponse) error {
	log.Info("Received Static.ReadDevice request")
	device, err := db.ReadDevice(ctx, req.Id, req.OrgId, req.TeamId)
	if device == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadDevice, err, "device not found")
	}
	rsp.Data = &static_proto.DeviceData{device}
	return nil
}

func (p *StaticService) DeleteDevice(ctx context.Context, req *static_proto.DeleteDeviceRequest, rsp *static_proto.DeleteDeviceResponse) error {
	log.Info("Received Static.DeleteDevice request")
	if err := db.DeleteDevice(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteDevice, nil, "delete error")
	}
	return nil
}

func (p *StaticService) AllWearables(ctx context.Context, req *static_proto.AllWearablesRequest, rsp *static_proto.AllWearablesResponse) error {
	log.Info("Received Static.AllWearables request")
	wearables, err := db.AllWearables(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(wearables) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllWearables, err, "not found")
	}
	rsp.Data = &static_proto.WearableArrData{wearables}
	return nil
}

func (p *StaticService) CreateWearable(ctx context.Context, req *static_proto.CreateWearableRequest, rsp *static_proto.CreateWearableResponse) error {
	log.Info("Received Static.CreateWearable request")
	if len(req.Wearable.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateWearable, nil, "wearable name empty")
	}
	if len(req.Wearable.Id) == 0 {
		req.Wearable.Id = uuid.NewUUID().String()
	}

	err := db.CreateWearable(ctx, req.Wearable)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateWearable, err, "wearable create error")
	}
	rsp.Data = &static_proto.WearableData{req.Wearable}
	return nil
}

func (p *StaticService) ReadWearable(ctx context.Context, req *static_proto.ReadWearableRequest, rsp *static_proto.ReadWearableResponse) error {
	log.Info("Received Static.ReadWearable request")
	wearable, err := db.ReadWearable(ctx, req.Id, req.OrgId, req.TeamId)
	if wearable == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadWearable, err, "wearable not found")
	}
	rsp.Data = &static_proto.WearableData{wearable}
	return nil
}

func (p *StaticService) DeleteWearable(ctx context.Context, req *static_proto.DeleteWearableRequest, rsp *static_proto.DeleteWearableResponse) error {
	log.Info("Received Static.DeleteWearable request")
	if err := db.DeleteWearable(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteWearable, nil, "delete error")
	}
	return nil
}

func (p *StaticService) AllMarkers(ctx context.Context, req *static_proto.AllMarkersRequest, rsp *static_proto.AllMarkersResponse) error {
	log.Info("Received Static.AllMarkers request")
	markers, err := db.AllMarkers(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(markers) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllMarkers, err, "marker not found")
	}
	rsp.Data = &static_proto.MarkerArrData{markers}
	return nil
}

func (p *StaticService) CreateMarker(ctx context.Context, req *static_proto.CreateMarkerRequest, rsp *static_proto.CreateMarkerResponse) error {
	log.Info("Received Static.CreateMarker request")
	if len(req.Marker.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateMarker, nil, "marker name empty")
	}
	if len(req.Marker.Id) == 0 {
		req.Marker.Id = uuid.NewUUID().String()
	}

	err := db.CreateMarker(ctx, req.Marker)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateMarker, err, "marker create error")
	}
	rsp.Data = &static_proto.MarkerData{req.Marker}
	return nil
}

func (p *StaticService) ReadMarker(ctx context.Context, req *static_proto.ReadMarkerRequest, rsp *static_proto.ReadMarkerResponse) error {
	log.Info("Received Static.ReadMarker request")
	marker, err := db.ReadMarker(ctx, req.Id, req.OrgId, req.TeamId)
	if marker == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadMarker, err, "marker not found")
	}
	rsp.Data = &static_proto.MarkerData{marker}
	return nil
}

func (p *StaticService) DeleteMarker(ctx context.Context, req *static_proto.DeleteMarkerRequest, rsp *static_proto.DeleteMarkerResponse) error {
	log.Info("Received Static.DeleteMarker request")
	if err := db.DeleteMarker(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteMarker, err, "delete error")
	}
	return nil
}

func (p *StaticService) FilterMarker(ctx context.Context, req *static_proto.FilterMarkerRequest, rsp *static_proto.FilterMarkerResponse) error {
	log.Info("Received Static.FilterMarker request")
	markes, err := db.FilterMarker(ctx, req.TrackerMethods, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(markes) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.FilterMarker, err, "marker not found")
	}
	rsp.Data = &static_proto.MarkerArrData{markes}
	return nil
}

func (p *StaticService) AllModules(ctx context.Context, req *static_proto.AllModulesRequest, rsp *static_proto.AllModulesResponse) error {
	log.Info("Received Static.AllModules request")
	modules, err := db.AllModules(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(modules) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllModules, err, "module not found")
	}
	rsp.Data = &static_proto.ModuleArrData{modules}
	return nil
}

func (p *StaticService) CreateModule(ctx context.Context, req *static_proto.CreateModuleRequest, rsp *static_proto.CreateModuleResponse) error {
	log.Info("Received Static.CreateModule request")
	if len(req.Module.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateModule, nil, "module name empty")
	}
	if len(req.Module.Id) == 0 {
		req.Module.Id = uuid.NewUUID().String()
	}

	err := db.CreateModule(ctx, req.Module)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateModule, err, "create error")
	}
	rsp.Data = &static_proto.ModuleData{req.Module}
	return nil
}

func (p *StaticService) ReadModule(ctx context.Context, req *static_proto.ReadModuleRequest, rsp *static_proto.ReadModuleResponse) error {
	log.Info("Received Static.ReadModule request")
	module, err := db.ReadModule(ctx, req.Id, req.OrgId, req.TeamId)
	if module == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadModule, err, "module not found")
	}
	rsp.Data = &static_proto.ModuleData{module}
	return nil
}

func (p *StaticService) DeleteModule(ctx context.Context, req *static_proto.DeleteModuleRequest, rsp *static_proto.DeleteModuleResponse) error {
	log.Info("Received Static.DeleteModule request")
	if err := db.DeleteModule(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteModule, err, "delete error")
	}
	return nil
}

func (p *StaticService) AllBehaviourCategories(ctx context.Context, req *static_proto.AllBehaviourCategoriesRequest, rsp *static_proto.AllBehaviourCategoriesResponse) error {
	log.Info("Received Static.AllBehaviourCategories request")
	categories, err := db.AllBehaviourCategories(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, "created", "ASC")
	if len(categories) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllBehaviourCategories, err, "category not found")
	}
	rsp.Data = &static_proto.BehaviourCategoryArrData{categories}
	return nil
}

func (p *StaticService) CreateBehaviourCategory(ctx context.Context, req *static_proto.CreateBehaviourCategoryRequest, rsp *static_proto.CreateBehaviourCategoryResponse) error {
	log.Info("Received Static.CreateBehaviourCategory request")
	if len(req.Category.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateBehaviourCategory, nil, "category name empty")
	}
	if len(req.Category.Id) == 0 {
		req.Category.Id = uuid.NewUUID().String()
	}

	err := db.CreateBehaviourCategory(ctx, req.Category)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateBehaviourCategory, err, "create error")
	}
	rsp.Data = &static_proto.BehaviourCategoryData{req.Category}
	return nil
}

func (p *StaticService) ReadBehaviourCategory(ctx context.Context, req *static_proto.ReadBehaviourCategoryRequest, rsp *static_proto.ReadBehaviourCategoryResponse) error {
	log.Info("Received Static.ReadBehaviourCategory request")

	category, err := db.ReadBehaviourCategory(ctx, req.Id, req.OrgId, req.TeamId)
	if category == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadBehaviourCategory, err, "not found")
	}
	rsp.Data = &static_proto.BehaviourCategoryData{category}
	return nil
}

func (p *StaticService) DeleteBehaviourCategory(ctx context.Context, req *static_proto.DeleteBehaviourCategoryRequest, rsp *static_proto.DeleteBehaviourCategoryResponse) error {
	log.Info("Received Static.DeleteBehaviourCategory request")
	if err := db.DeleteBehaviourCategory(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteBehaviourCategory, nil, "delete error")
	}
	return nil
}

func (p *StaticService) FilterBehaviourCategory(ctx context.Context, req *static_proto.FilterBehaviourCategoryRequest, rsp *static_proto.FilterBehaviourCategoryResponse) error {
	log.Info("Received Static.FilterBehaviourCategory request")
	categories, err := db.FilterBehaviourCategory(ctx, req.Markers, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(categories) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.FilterBehaviourCategory, err, "not found")
	}
	rsp.Data = &static_proto.BehaviourCategoryArrData{categories}
	return nil
}

func (p *StaticService) AllSocialTypes(ctx context.Context, req *static_proto.AllSocialTypesRequest, rsp *static_proto.AllSocialTypesResponse) error {
	log.Info("Received Static.AllSocialTypes request")
	socialTypes, err := db.AllSocialTypes(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(socialTypes) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllSocialTypes, err, "not found")
	}
	rsp.Data = &static_proto.SocialTypeArrData{socialTypes}
	return nil
}

func (p *StaticService) CreateSocialType(ctx context.Context, req *static_proto.CreateSocialTypeRequest, rsp *static_proto.CreateSocialTypeResponse) error {
	log.Info("Received Static.CreateSocialType request")
	if len(req.SocialType.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateSocialType, nil, "social type name empty")
	}
	if len(req.SocialType.Id) == 0 {
		req.SocialType.Id = uuid.NewUUID().String()
	}

	err := db.CreateSocialType(ctx, req.SocialType)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateSocialType, err, "create error")
	}
	rsp.Data = &static_proto.SocialTypeData{req.SocialType}
	return nil
}

func (p *StaticService) ReadSocialType(ctx context.Context, req *static_proto.ReadSocialTypeRequest, rsp *static_proto.ReadSocialTypeResponse) error {
	log.Info("Received Static.ReadSocialType request")
	socialType, err := db.ReadSocialType(ctx, req.Id, req.OrgId, req.TeamId)
	if socialType == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadSocialType, err, "not found")
	}
	rsp.Data = &static_proto.SocialTypeData{socialType}
	return nil
}

func (p *StaticService) DeleteSocialType(ctx context.Context, req *static_proto.DeleteSocialTypeRequest, rsp *static_proto.DeleteSocialTypeResponse) error {
	log.Info("Received Static.DeleteSocialType request")
	if err := db.DeleteSocialType(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteSocialType, nil, "delete error")
	}
	return nil
}

func (p *StaticService) AllNotifications(ctx context.Context, req *static_proto.AllNotificationsRequest, rsp *static_proto.AllNotificationsResponse) error {
	log.Info("Received Static.AllNotifications request")
	notifications, err := db.AllNotifications(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(notifications) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllNotifications, err, "not found")
	}
	rsp.Data = &static_proto.NotificationArrData{notifications}
	return nil
}

func (p *StaticService) CreateNotification(ctx context.Context, req *static_proto.CreateNotificationRequest, rsp *static_proto.CreateNotificationResponse) error {
	log.Info("Received Static.CreateNotification request")
	if len(req.Notification.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateNotification, nil, "notification name empty")
	}
	if len(req.Notification.Id) == 0 {
		req.Notification.Id = uuid.NewUUID().String()
	}

	err := db.CreateNotification(ctx, req.Notification)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateNotification, err, "create error")
	}
	rsp.Data = &static_proto.NotificationData{req.Notification}
	return nil
}

func (p *StaticService) ReadNotification(ctx context.Context, req *static_proto.ReadNotificationRequest, rsp *static_proto.ReadNotificationResponse) error {
	log.Info("Received Static.ReadNotification request")
	notification, err := db.ReadNotification(ctx, req.Id, req.OrgId, req.TeamId)
	if notification == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadNotification, err, "not found")
	}
	rsp.Data = &static_proto.NotificationData{notification}
	return nil
}

func (p *StaticService) DeleteNotification(ctx context.Context, req *static_proto.DeleteNotificationRequest, rsp *static_proto.DeleteNotificationResponse) error {
	log.Info("Received Static.DeleteNotification request")
	if err := db.DeleteNotification(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteNotification, nil, "delete error")
	}
	return nil
}

func (p *StaticService) AllTrackerMethods(ctx context.Context, req *static_proto.AllTrackerMethodsRequest, rsp *static_proto.AllTrackerMethodsResponse) error {
	log.Info("Received Static.AllTrackerMethods request")
	trackerMethods, err := db.AllTrackerMethods(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(trackerMethods) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllTrackerMethods, err, "not found")
	}
	rsp.Data = &static_proto.TrackerMethodArrData{trackerMethods}
	return nil
}

func (p *StaticService) CreateTrackerMethod(ctx context.Context, req *static_proto.CreateTrackerMethodRequest, rsp *static_proto.CreateTrackerMethodResponse) error {
	log.Info("Received Static.CreateTrackerMethod request")
	if len(req.TrackerMethod.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateTrackerMethod, nil, "tracker method name empty")
	}
	if len(req.TrackerMethod.Id) == 0 {
		req.TrackerMethod.Id = uuid.NewUUID().String()
	}

	err := db.CreateTrackerMethod(ctx, req.TrackerMethod)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateTrackerMethod, err, "create error")
	}
	rsp.Data = &static_proto.TrackerMethodData{req.TrackerMethod}
	return nil
}

func (p *StaticService) ReadTrackerMethod(ctx context.Context, req *static_proto.ReadTrackerMethodRequest, rsp *static_proto.ReadTrackerMethodResponse) error {
	log.WithFields(log.Fields{"Id": req.Id, "NameSlug": req.NameSlug}).Info("Received Static.ReadTrackerMethod request for:")
	trackerMethod, err := db.ReadTrackerMethod(ctx, req.Id, req.NameSlug, req.OrgId, req.TeamId)
	if trackerMethod == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadTrackerMethod, err, "not found")
	}
	rsp.Data = &static_proto.TrackerMethodData{trackerMethod}
	return nil
}

func (p *StaticService) DeleteTrackerMethod(ctx context.Context, req *static_proto.DeleteTrackerMethodRequest, rsp *static_proto.DeleteTrackerMethodResponse) error {
	log.Info("Received Static.DeleteTrackerMethod request")
	if err := db.DeleteTrackerMethod(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteTrackerMethod, nil, "delete error")
	}
	return nil
}

func (p *StaticService) FilterTrackerMethod(ctx context.Context, req *static_proto.FilterTrackerMethodRequest, rsp *static_proto.FilterTrackerMethodResponse) error {
	log.Info("Received Static.FilterTrackerMethod request")
	trackerMethods, err := db.FilterTrackerMethod(ctx, req.Markers, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(trackerMethods) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.FilterTrackerMethod, err, "not found")
	}
	rsp.Data = &static_proto.TrackerMethodArrData{trackerMethods}
	return nil
}

func (p *StaticService) AllBehaviourCategoryAims(ctx context.Context, req *static_proto.AllBehaviourCategoryAimsRequest, rsp *static_proto.AllBehaviourCategoryAimsResponse) error {
	log.Info("Received Static.AllBehaviourCategoryAims request")
	behaviourCategoryAims, err := db.AllBehaviourCategoryAims(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(behaviourCategoryAims) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllBehaviourCategoryAims, err, "not found")
	}
	rsp.Data = &static_proto.BehaviourCategoryAimArrData{behaviourCategoryAims}
	return nil
}

func (p *StaticService) CreateBehaviourCategoryAim(ctx context.Context, req *static_proto.CreateBehaviourCategoryAimRequest, rsp *static_proto.CreateBehaviourCategoryAimResponse) error {
	log.Info("Received Static.CreateBehaviourCategoryAim request")
	if len(req.BehaviourCategoryAim.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateBehaviourCategoryAim, nil, "behaviour category aim name empty")
	}
	if len(req.BehaviourCategoryAim.Id) == 0 {
		req.BehaviourCategoryAim.Id = uuid.NewUUID().String()
	}

	err := db.CreateBehaviourCategoryAim(ctx, req.BehaviourCategoryAim)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateBehaviourCategoryAim, err, "create error")
	}
	rsp.Data = &static_proto.BehaviourCategoryAimData{req.BehaviourCategoryAim}
	return nil
}

func (p *StaticService) ReadBehaviourCategoryAim(ctx context.Context, req *static_proto.ReadBehaviourCategoryAimRequest, rsp *static_proto.ReadBehaviourCategoryAimResponse) error {
	log.Info("Received Static.ReadBehaviourCategoryAim request")
	behaviourCategoryAim, err := db.ReadBehaviourCategoryAim(ctx, req.Id, req.OrgId, req.TeamId)
	if behaviourCategoryAim == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadBehaviourCategoryAim, err, "not found")
	}
	rsp.Data = &static_proto.BehaviourCategoryAimData{behaviourCategoryAim}
	return nil
}

func (p *StaticService) DeleteBehaviourCategoryAim(ctx context.Context, req *static_proto.DeleteBehaviourCategoryAimRequest, rsp *static_proto.DeleteBehaviourCategoryAimResponse) error {
	log.Info("Received Static.DeleteBehaviourCategoryAim request")
	if err := db.DeleteBehaviourCategoryAim(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteBehaviourCategoryAim, nil, "delete error")
	}
	return nil
}

func (p *StaticService) AllContentParentCategories(ctx context.Context, req *static_proto.AllContentParentCategoriesRequest, rsp *static_proto.AllContentParentCategoriesResponse) error {
	log.Info("Received Static.AllContentParentCategories request")
	contentParentCategories, err := db.AllContentParentCategories(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(contentParentCategories) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllContentParentCategories, err, "not found")
	}
	rsp.Data = &static_proto.ContentParentCategoryArrData{contentParentCategories}
	return nil
}

func (p *StaticService) CreateContentParentCategory(ctx context.Context, req *static_proto.CreateContentParentCategoryRequest, rsp *static_proto.CreateContentParentCategoryResponse) error {
	log.Info("Received Static.CreateContentParentCategory request")
	if len(req.ContentParentCategory.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateContentParentCategory, nil, "category name empty")
	}
	if len(req.ContentParentCategory.Id) == 0 {
		req.ContentParentCategory.Id = uuid.NewUUID().String()
	}

	err := db.CreateContentParentCategory(ctx, req.ContentParentCategory)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateContentParentCategory, err, "create error")
	}
	rsp.Data = &static_proto.ContentParentCategoryData{req.ContentParentCategory}
	return nil
}

func (p *StaticService) ReadContentParentCategory(ctx context.Context, req *static_proto.ReadContentParentCategoryRequest, rsp *static_proto.ReadContentParentCategoryResponse) error {
	log.Info("Received Static.ReadContentParentCategory request")
	contentParentCategory, err := db.ReadContentParentCategory(ctx, req.Id, req.OrgId, req.TeamId)
	if contentParentCategory == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadContentParentCategory, err, "not found")
	}
	rsp.Data = &static_proto.ContentParentCategoryData{contentParentCategory}
	return nil
}

func (p *StaticService) DeleteContentParentCategory(ctx context.Context, req *static_proto.DeleteContentParentCategoryRequest, rsp *static_proto.DeleteContentParentCategoryResponse) error {
	log.Info("Received Static.DeleteContentParentCategory request")
	if err := db.DeleteContentParentCategory(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteContentParentCategory, nil, "delete error")
	}
	return nil
}

func (p *StaticService) AllContentCategories(ctx context.Context, req *static_proto.AllContentCategoriesRequest, rsp *static_proto.AllContentCategoriesResponse) error {
	log.Info("Received Static.AllContentCategories request")
	contentCategories, err := db.AllContentCategories(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(contentCategories) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllContentCategories, err, "not found")
	}
	rsp.Data = &static_proto.ContentCategoryArrData{contentCategories}
	return nil
}

func (p *StaticService) CreateContentCategory(ctx context.Context, req *static_proto.CreateContentCategoryRequest, rsp *static_proto.CreateContentCategoryResponse) error {
	log.Info("Received Static.CreateContentCategory request")
	if len(req.ContentCategory.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateContentCategory, nil, "content category name empty")
	}
	if len(req.ContentCategory.Id) == 0 {
		req.ContentCategory.Id = uuid.NewUUID().String()
	}
	err := db.CreateContentCategory(ctx, req.ContentCategory)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateContentCategory, err, "create error")
	}
	rsp.Data = &static_proto.ContentCategoryData{req.ContentCategory}
	return nil
}

func (p *StaticService) ReadContentCategory(ctx context.Context, req *static_proto.ReadContentCategoryRequest, rsp *static_proto.ReadContentCategoryResponse) error {
	log.Info("Received Static.ReadContentCategory request")
	contentCategory, err := db.ReadContentCategory(ctx, req.Id, req.OrgId, req.TeamId)
	if contentCategory == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadContentCategory, err, "not found")
	}
	rsp.Data = &static_proto.ContentCategoryData{contentCategory}
	return nil
}

func (p *StaticService) DeleteContentCategory(ctx context.Context, req *static_proto.DeleteContentCategoryRequest, rsp *static_proto.DeleteContentCategoryResponse) error {
	log.Info("Received Static.DeleteContentCategory request")
	if err := db.DeleteContentCategory(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteContentCategory, nil, "delete error")
	}
	return nil
}

func (p *StaticService) AllContentTypes(ctx context.Context, req *static_proto.AllContentTypesRequest, rsp *static_proto.AllContentTypesResponse) error {
	log.Info("Received Static.AllContentTypes request")
	contentTypes, err := db.AllContentTypes(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(contentTypes) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllContentTypes, err, "not found")
	}
	rsp.Data = &static_proto.ContentTypeArrData{contentTypes}
	return nil
}

func (p *StaticService) CreateContentType(ctx context.Context, req *static_proto.CreateContentTypeRequest, rsp *static_proto.CreateContentTypeResponse) error {
	log.Info("Received Static.CreateContentType request")
	if len(req.ContentType.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateContentType, nil, "content type name empty")
	}
	if len(req.ContentType.Id) == 0 {
		req.ContentType.Id = uuid.NewUUID().String()
	}

	err := db.CreateContentType(ctx, req.ContentType)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateContentType, err, "create error")
	}
	rsp.Data = &static_proto.ContentTypeData{req.ContentType}
	return nil
}

func (p *StaticService) ReadContentType(ctx context.Context, req *static_proto.ReadContentTypeRequest, rsp *static_proto.ReadContentTypeResponse) error {
	log.Info("Received Static.ReadContentType request")
	contentType, err := db.ReadContentType(ctx, req.Id, req.OrgId, req.TeamId)
	if contentType == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadContentType, err, "not found")
	}
	rsp.Data = &static_proto.ContentTypeData{contentType}
	return nil
}

func (p *StaticService) DeleteContentType(ctx context.Context, req *static_proto.DeleteContentTypeRequest, rsp *static_proto.DeleteContentTypeResponse) error {
	log.Info("Received Static.DeleteContentType request")
	if err := db.DeleteContentType(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteContentType, err, "delete error")
	}
	return nil
}

func (p *StaticService) AllContentSourceTypes(ctx context.Context, req *static_proto.AllContentSourceTypesRequest, rsp *static_proto.AllContentSourceTypesResponse) error {
	log.Info("Received Static.AllContentSourceTypes request")
	contentSourceTypes, err := db.AllContentSourceTypes(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(contentSourceTypes) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllContentSourceTypes, err, "not found")
	}
	rsp.Data = &static_proto.ContentSourceTypeArrData{contentSourceTypes}
	return nil
}

func (p *StaticService) CreateContentSourceType(ctx context.Context, req *static_proto.CreateContentSourceTypeRequest, rsp *static_proto.CreateContentSourceTypeResponse) error {
	log.Info("Received Static.CreateContentSourceType request")
	if len(req.ContentSourceType.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateContentSourceType, nil, "content source type name empty")
	}
	if len(req.ContentSourceType.Id) == 0 {
		req.ContentSourceType.Id = uuid.NewUUID().String()
	}

	err := db.CreateContentSourceType(ctx, req.ContentSourceType)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateContentSourceType, err, "create error")
	}
	rsp.Data = &static_proto.ContentSourceTypeData{req.ContentSourceType}
	return nil
}

func (p *StaticService) ReadContentSourceType(ctx context.Context, req *static_proto.ReadContentSourceTypeRequest, rsp *static_proto.ReadContentSourceTypeResponse) error {
	log.Info("Received Static.ReadContentSourceType request")
	contentSourceType, err := db.ReadContentSourceType(ctx, req.Id, req.OrgId, req.TeamId)
	if contentSourceType == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadContentSourceType, err, "not found")
	}
	rsp.Data = &static_proto.ContentSourceTypeData{contentSourceType}
	return nil
}

func (p *StaticService) DeleteContentSourceType(ctx context.Context, req *static_proto.DeleteContentSourceTypeRequest, rsp *static_proto.DeleteContentSourceTypeResponse) error {
	log.Info("Received Static.DeleteContentSourceType request")
	if err := db.DeleteContentSourceType(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteContentSourceType, nil, "delete error")
	}
	return nil
}

func (p *StaticService) AllModuleTriggers(ctx context.Context, req *static_proto.AllModuleTriggersRequest, rsp *static_proto.AllModuleTriggersResponse) error {
	log.Info("Received Static.AllModuleTriggers request")
	moduleTriggers, err := db.AllModuleTriggers(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(moduleTriggers) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllModuleTriggers, err, "not found")
	}
	rsp.Data = &static_proto.ModuleTriggerArrData{moduleTriggers}
	return nil
}

func (p *StaticService) CreateModuleTrigger(ctx context.Context, req *static_proto.CreateModuleTriggerRequest, rsp *static_proto.CreateModuleTriggerResponse) error {
	log.Info("Received Static.CreateModuleTrigger request")
	if len(req.ModuleTrigger.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateModuleTrigger, nil, "module trigger name empty")
	}
	if len(req.ModuleTrigger.Id) == 0 {
		req.ModuleTrigger.Id = uuid.NewUUID().String()
	}

	err := db.CreateModuleTrigger(ctx, req.ModuleTrigger)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateModuleTrigger, err, "create error")
	}
	rsp.Data = &static_proto.ModuleTriggerData{req.ModuleTrigger}
	return nil
}

func (p *StaticService) ReadModuleTrigger(ctx context.Context, req *static_proto.ReadModuleTriggerRequest, rsp *static_proto.ReadModuleTriggerResponse) error {
	log.Info("Received Static.ReadModuleTrigger request")
	moduleTrigger, err := db.ReadModuleTrigger(ctx, req.Id, req.OrgId, req.TeamId)
	if moduleTrigger == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadModuleTrigger, err, "not found")
	}
	rsp.Data = &static_proto.ModuleTriggerData{moduleTrigger}
	return nil
}

func (p *StaticService) DeleteModuleTrigger(ctx context.Context, req *static_proto.DeleteModuleTriggerRequest, rsp *static_proto.DeleteModuleTriggerResponse) error {
	log.Info("Received Static.DeleteModuleTrigger request")
	if err := db.DeleteModuleTrigger(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteModuleTrigger, nil, "delete error")
	}
	return nil
}

func (p *StaticService) FilterModuleTrigger(ctx context.Context, req *static_proto.FilterModuleTriggerRequest, rsp *static_proto.FilterModuleTriggerResponse) error {
	log.Info("Received Static.FilterModuleTrigger request")
	moduleTriggers, err := db.FilterModuleTrigger(ctx, req.Module, req.TriggerType, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(moduleTriggers) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.FilterModuleTrigger, err, "not found")
	}
	rsp.Data = &static_proto.ModuleTriggerArrData{moduleTriggers}
	return nil
}

func (p *StaticService) AllTriggerContentTypes(ctx context.Context, req *static_proto.AllTriggerContentTypesRequest, rsp *static_proto.AllTriggerContentTypesResponse) error {
	log.Info("Received Static.AllTriggerContentTypes request")
	triggerContentTypes, err := db.AllTriggerContentTypes(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(triggerContentTypes) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllTriggerContentTypes, err, "not found")
	}
	rsp.Data = &static_proto.TriggerContentTypeArrData{triggerContentTypes}
	return nil
}

func (p *StaticService) CreateTriggerContentType(ctx context.Context, req *static_proto.CreateTriggerContentTypeRequest, rsp *static_proto.CreateTriggerContentTypeResponse) error {
	log.Info("Received Static.CreateTriggerContentType request")
	if len(req.TriggerContentType.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateTriggerContentType, nil, "trigger content type name")
	}
	if len(req.TriggerContentType.Id) == 0 {
		req.TriggerContentType.Id = uuid.NewUUID().String()
	}

	err := db.CreateTriggerContentType(ctx, req.TriggerContentType)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateTriggerContentType, err, "create error")
	}
	rsp.Data = &static_proto.TriggerContentTypeData{req.TriggerContentType}
	return nil
}

func (p *StaticService) ReadTriggerContentType(ctx context.Context, req *static_proto.ReadTriggerContentTypeRequest, rsp *static_proto.ReadTriggerContentTypeResponse) error {
	log.Info("Received Static.ReadTriggerContentType request")
	triggerContentType, err := db.ReadTriggerContentType(ctx, req.Id, req.OrgId, req.TeamId)
	if triggerContentType == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadTriggerContentType, err, "not found")
	}
	rsp.Data = &static_proto.TriggerContentTypeData{triggerContentType}
	return nil
}

func (p *StaticService) DeleteTriggerContentType(ctx context.Context, req *static_proto.DeleteTriggerContentTypeRequest, rsp *static_proto.DeleteTriggerContentTypeResponse) error {
	log.Info("Received Static.DeleteTriggerContentType request")
	if err := db.DeleteTriggerContentType(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteTriggerContentType, nil, "delete error")
	}
	return nil
}

func (p *StaticService) FilterTriggerContentType(ctx context.Context, req *static_proto.FilterTriggerContentTypeRequest, rsp *static_proto.FilterTriggerContentTypeResponse) error {
	log.Info("Received Static.FilterTriggerContentType request")
	triggerContentTypes, err := db.FilterTriggerContentType(ctx, req.ModuleTrigger, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(triggerContentTypes) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.FilterTriggerContentType, err, "not found")
	}
	rsp.Data = &static_proto.TriggerContentTypeArrData{triggerContentTypes}
	return nil
}

//onl reays
func (p *StaticService) ReadContentCategoryByNameslug(ctx context.Context, req *static_proto.ReadByNameslugRequest, rsp *static_proto.ReadContentCategoryResponse) error {
	log.Info("Received Static.ReadContentCategoryByNameslug request")
	contentCategory, err := db.ReadContentCategoryByNameslug(ctx, req.NameSlug)
	if contentCategory == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadContentCategoryByNameslug, err, "not found")
	}
	rsp.Data = &static_proto.ContentCategoryData{contentCategory}
	return nil
}

//reads and create if no contentCategory found
//ReadContentCategoryByNameslugOrCreate This is for internal calls ONLY! where new content category is created if none exists, where the caller needs to create a new object based on this

func (p *StaticService) ReadContentCategoryByNameslugOrCreate(ctx context.Context, req *static_proto.ReadByNameslugRequest, rsp *static_proto.ReadContentCategoryResponse) error {
	log.Info("Received Static.ReadContentCategoryByNameslugOrCreate request")
	contentCategory, err := db.ReadContentCategoryByNameslug(ctx, req.NameSlug)

	if contentCategory == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadContentCategoryByNameslugOrCreate, err, "not found")
	}
	//if non contentCategory is found create a new one
	if contentCategory == nil {
		contentCategory = &static_proto.ContentCategory{NameSlug: req.NameSlug}
		contentCategory.Id = uuid.NewUUID().String()
		err := db.CreateContentCategory(ctx, contentCategory)
		if err != nil {
			return common.InternalServerError(common.StaticSrv, p.ReadContentCategoryByNameslugOrCreate, err, "create error")
		}
	}

	rsp.Data = &static_proto.ContentCategoryData{contentCategory}
	return nil
}

func (p *StaticService) CreateRole(ctx context.Context, req *static_proto.CreateRoleRequest, rsp *static_proto.CreateRoleResponse) error {
	log.Info("Received Static.CreateRole request")
	err := db.CreateRole(ctx, req.Role)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateRole, err, "create error")
	}
	rsp.Data = &static_proto.RoleData{req.Role}
	return nil
}

func (p *StaticService) ReadRoleByNameslug(ctx context.Context, req *static_proto.ReadRoleByNameslugRequest, rsp *static_proto.ReadRoleByNameslugResponse) error {
	log.Info("Received Static.ReadRoleByNameslug request")
	role, err := db.ReadRoleByNameslug(ctx, req.NameSlug)
	if role == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadRoleByNameslug, err, "not found")
	}
	rsp.Data = &static_proto.RoleData{role}
	return nil
}

func (p *StaticService) AllSetbacks(ctx context.Context, req *static_proto.AllSetbacksRequest, rsp *static_proto.AllSetbacksResponse) error {
	log.Info("Received Static.AllSetbacks request")
	setbacks, err := db.AllSetbacks(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(setbacks) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AllSetbacks, err, "not found")
	}
	rsp.Data = &static_proto.SetbackArrData{setbacks}
	return nil
}

func (p *StaticService) CreateSetback(ctx context.Context, req *static_proto.CreateSetbackRequest, rsp *static_proto.CreateSetbackResponse) error {
	log.Info("Received Static.CreateSetback request")
	if len(req.Setback.Name) == 0 {
		return common.InternalServerError(common.StaticSrv, p.CreateSetback, nil, "set back name empty")
	}
	// if len(req.Setback.OrgId) == 0 {
	// 	return errors.InternalServerError("go.micro.srv.static.Create", "OrgId cannot be empty")
	// }
	if len(req.Setback.Id) == 0 {
		req.Setback.Id = uuid.NewUUID().String()
	}

	err := db.CreateSetback(ctx, req.Setback)
	if err != nil {
		return common.InternalServerError(common.StaticSrv, p.CreateSetback, err, "create error")
	}
	rsp.Data = &static_proto.SetbackData{req.Setback}
	return nil
}

func (p *StaticService) ReadSetback(ctx context.Context, req *static_proto.ReadSetbackRequest, rsp *static_proto.ReadSetbackResponse) error {
	log.Info("Received Static.ReadSetback request")
	setback, err := db.ReadSetback(ctx, req.Id, req.OrgId, req.TeamId)
	if setback == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadSetback, err, "not found")
	}
	rsp.Data = &static_proto.SetbackData{setback}
	return nil
}

func (p *StaticService) DeleteSetback(ctx context.Context, req *static_proto.DeleteSetbackRequest, rsp *static_proto.DeleteSetbackResponse) error {
	log.Info("Received Static.DeleteSetback request")
	if err := db.DeleteSetback(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.StaticSrv, p.DeleteSetback, err, "delete error")
	}
	return nil
}

func (p *StaticService) AutocompleteSetbackSearch(ctx context.Context, req *static_proto.AutocompleteSetbackSearchRequest, rsp *static_proto.AllSetbacksResponse) error {
	log.Info("Received Static.AutocompleteSetbackSearch request")
	setbacks, err := db.AutocompleteSetbackSearch(ctx, req.Title)
	if len(setbacks) == 0 || err != nil {
		return common.NotFound(common.StaticSrv, p.AutocompleteSetbackSearch, err, "not found")
	}
	rsp.Data = &static_proto.SetbackArrData{setbacks}
	return nil
}

func (p *StaticService) ReadMarkerByNameslug(ctx context.Context, req *static_proto.ReadByNameslugRequest, rsp *static_proto.ReadMarkerResponse) error {
	log.Info("Received Static.ReadMarkerByNameslug request")
	marker, err := db.ReadMarkerByNameslug(ctx, req.NameSlug)
	if marker == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadMarkerByNameslug, err, "not found")
	}
	rsp.Data = &static_proto.MarkerData{marker}
	return nil
}

func (p *StaticService) ReadBehaviourCategoryAimByNameslug(ctx context.Context, req *static_proto.ReadByNameslugRequest, rsp *static_proto.ReadBehaviourCategoryAimResponse) error {
	log.Info("Received Static.ReadBehaviourCategoryAimByNameslug request")
	behaviourCategoryAim, err := db.ReadBehaviourCategoryAimByNameslug(ctx, req.NameSlug)
	if behaviourCategoryAim == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadBehaviourCategoryAimByNameslug, err, "not found")
	}
	rsp.Data = &static_proto.BehaviourCategoryAimData{behaviourCategoryAim}
	return nil
}

func (p *StaticService) ReadTrackerMethodByNameslug(ctx context.Context, req *static_proto.ReadByNameslugRequest, rsp *static_proto.ReadTrackerMethodResponse) error {
	log.Info("Received Static.ReadTrackerMethodByNameslug request")
	method, err := db.ReadTrackerMethodByNameslug(ctx, req.NameSlug)
	if method == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadTrackerMethodByNameslug, err, "not found")
	}
	rsp.Data = &static_proto.TrackerMethodData{method}
	return nil
}

func (p *StaticService) ReadAppByNameslug(ctx context.Context, req *static_proto.ReadByNameslugRequest, rsp *static_proto.ReadAppResponse) error {
	log.Info("Received Static.ReadAppByNameslug request")
	app, err := db.ReadAppByNameslug(ctx, req.NameSlug)
	if app == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadAppByNameslug, err, "not found")
	}
	rsp.Data = &static_proto.AppData{app}
	return nil
}

func (p *StaticService) ReadWearableByNameslug(ctx context.Context, req *static_proto.ReadByNameslugRequest, rsp *static_proto.ReadWearableResponse) error {
	log.Info("Received Static.ReadWearableByNameslug request")
	wearable, err := db.ReadWearableByNameslug(ctx, req.NameSlug)
	if wearable == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadWearableByNameslug, err, "not found")
	}
	rsp.Data = &static_proto.WearableData{wearable}
	return nil
}

func (p *StaticService) ReadDeviceByNameslug(ctx context.Context, req *static_proto.ReadByNameslugRequest, rsp *static_proto.ReadDeviceResponse) error {
	log.Info("Received Static.ReadDeviceByNameslug request")
	device, err := db.ReadDeviceByNameslug(ctx, req.NameSlug)
	if device == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadDeviceByNameslug, err, "not found")
	}
	rsp.Data = &static_proto.DeviceData{device}
	return nil
}

func (p *StaticService) UpdateApp(ctx context.Context, req *static_proto.UpdateAppRequest, rsp *static_proto.UpdateAppResponse) error {
	log.Info("Received Static.UpdateApp request")
	return nil
}

func (p *StaticService) UpdatePlatform(ctx context.Context, req *static_proto.UpdatePlatformRequest, rsp *static_proto.UpdatePlatformResponse) error {
	log.Info("Received Static.UpdatePlatform request")
	return nil
}

func (p *StaticService) UpdateDevice(ctx context.Context, req *static_proto.UpdateDeviceRequest, rsp *static_proto.UpdateDeviceResponse) error {
	log.Info("Received Static.UpdateDevice request")
	return nil
}

func (p *StaticService) UpdateWearable(ctx context.Context, req *static_proto.UpdateWearableRequest, rsp *static_proto.UpdateWearableResponse) error {
	log.Info("Received Static.UpdateWearable request")
	return nil
}

func (p *StaticService) UpdateMarker(ctx context.Context, req *static_proto.UpdateMarkerRequest, rsp *static_proto.UpdateMarkerResponse) error {
	log.Info("Received Static.UpdateMarker request")
	return nil
}

func (p *StaticService) UpdateModule(ctx context.Context, req *static_proto.UpdateModuleRequest, rsp *static_proto.UpdateModuleResponse) error {
	log.Info("Received Static.UpdateModule request")
	return nil
}

func (p *StaticService) UpdateBehaviourCategory(ctx context.Context, req *static_proto.UpdateBehaviourCategoryRequest, rsp *static_proto.UpdateBehaviourCategoryResponse) error {
	log.Info("Received Static.UpdateBehaviourCategory request")
	return nil
}

func (p *StaticService) UpdateSocialType(ctx context.Context, req *static_proto.UpdateSocialTypeRequest, rsp *static_proto.UpdateSocialTypeResponse) error {
	log.Info("Received Static.UpdateSocialType request")
	return nil
}

func (p *StaticService) UpdateNotification(ctx context.Context, req *static_proto.UpdateNotificationRequest, rsp *static_proto.UpdateNotificationResponse) error {
	log.Info("Received Static.UpdateNotification request")
	return nil
}

func (p *StaticService) UpdateTrackerMethod(ctx context.Context, req *static_proto.UpdateTrackerMethodRequest, rsp *static_proto.UpdateTrackerMethodResponse) error {
	log.Info("Received Static.UpdateTrackerMethod request")
	return nil
}

func (p *StaticService) UpdateBehaviourCategoryAim(ctx context.Context, req *static_proto.UpdateBehaviourCategoryAimRequest, rsp *static_proto.UpdateBehaviourCategoryAimResponse) error {
	log.Info("Received Static.UpdateBehaviourCategoryAim request")
	return nil
}

func (p *StaticService) UpdateContentParentCategory(ctx context.Context, req *static_proto.UpdateContentParentCategoryRequest, rsp *static_proto.UpdateContentParentCategoryResponse) error {
	log.Info("Received Static.UpdateContentParentCategory request")
	return nil
}

func (p *StaticService) UpdateContentCategory(ctx context.Context, req *static_proto.UpdateContentCategoryRequest, rsp *static_proto.UpdateContentCategoryResponse) error {
	log.Info("Received Static.UpdateContentCategory request")
	return nil
}

func (p *StaticService) UpdateContentType(ctx context.Context, req *static_proto.UpdateContentTypeRequest, rsp *static_proto.UpdateContentTypeResponse) error {
	log.Info("Received Static.UpdateContentType request")
	return nil
}

func (p *StaticService) UpdateContentSourceType(ctx context.Context, req *static_proto.UpdateContentSourceTypeRequest, rsp *static_proto.UpdateContentSourceTypeResponse) error {
	log.Info("Received Static.UpdateContentSourceType request")
	return nil
}

func (p *StaticService) UpdateModuleTrigger(ctx context.Context, req *static_proto.UpdateModuleTriggerRequest, rsp *static_proto.UpdateModuleTriggerResponse) error {
	log.Info("Received Static.UpdateModuleTrigger request")
	return nil
}

func (p *StaticService) UpdateTriggerContentType(ctx context.Context, req *static_proto.UpdateTriggerContentTypeRequest, rsp *static_proto.UpdateTriggerContentTypeResponse) error {
	log.Info("Received Static.UpdateTriggerContentType request")
	return nil
}

func (p *StaticService) UpdateRole(ctx context.Context, req *static_proto.UpdateRoleRequest, rsp *static_proto.UpdateRoleResponse) error {
	log.Info("Received Static.UpdateRole request")
	return nil
}

func (p *StaticService) UpdateSetback(ctx context.Context, req *static_proto.UpdateSetbackRequest, rsp *static_proto.UpdateSetbackResponse) error {
	log.Info("Received Static.UpdateSetback request")
	return nil
}

func (p *StaticService) UploadBehaviourCategoryAim(ctx context.Context, req *static_proto.UploadRequest, rsp *static_proto.UploadResponse) error {
	log.Info("Received Static.UploadBehaviourCategoryAim request")

	return nil
}

func (p *StaticService) UploadContentCategory(ctx context.Context, req *static_proto.UploadRequest, rsp *static_proto.UploadResponse) error {
	log.Info("Received Static.UploadContentCategory request")

	return nil
}

func (p *StaticService) UploadMarker(ctx context.Context, req *static_proto.UploadRequest, rsp *static_proto.UploadResponse) error {
	log.Info("Received Static.UploadMarker request")

	return nil
}

func (p *StaticService) UploadBehaviourCategory(ctx context.Context, req *static_proto.UploadRequest, rsp *static_proto.UploadResponse) error {
	log.Info("Received Static.UploadBehaviourCategory request")

	return nil
}

func (p *StaticService) UploadContentCategoryItem(ctx context.Context, req *static_proto.UploadRequest, rsp *static_proto.UploadResponse) error {
	log.Info("Received Static.UploadBehaviourCategory request")

	return nil
}

func (p *StaticService) UploadTrackerMethod(ctx context.Context, req *static_proto.UploadRequest, rsp *static_proto.UploadResponse) error {
	log.Info("Received Static.UploadBehaviourCategory request")

	return nil
}
