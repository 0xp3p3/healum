package handler

import (
	"context"
	"server/common"
	"server/static-srv/db"
	static_proto "server/static-srv/proto/static"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

var cl = client.NewClient(
	client.Transport(nats_transport.NewTransport()),
	client.Broker(nats_broker.NewBroker()),
	client.RequestTimeout(4*time.Second),
	client.Retries(5),
)

func initDb() {
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

var app = &static_proto.App{
	Id:          "111",
	Name:        "title",
	Summary:     "summary",
	Description: "description",
	IconSlug:    "iconslug",
	Image:       "images",
	Tags:        []string{"tag1", "tag2"},
	Platforms:   []*static_proto.Platform{platform},
}

var platform = &static_proto.Platform{
	Id:   "111",
	Name: "title",
	Url:  "url",
}

var wearable = &static_proto.Wearable{
	Id:          "111",
	Name:        "title",
	Summary:     "summary",
	Description: "description",
	Image:       "images",
}

var device = &static_proto.Device{
	Id:          "111",
	Name:        "title",
	Summary:     "summary",
	Description: "description",
	IconSlug:    "iconslug",
	Image:       "images",
	Tags:        []string{"tag1", "tag2"},
}

var marker = &static_proto.Marker{
	// Id:             "111",
	Name:           "title",
	Summary:        "summary",
	Description:    "description",
	OrgId:          "orgid",
	Apps:           []*static_proto.App{app},
	Wearables:      []*static_proto.Wearable{wearable},
	Devices:        []*static_proto.Device{device},
	TrackerMethods: []*static_proto.TrackerMethod{trackerMethod},
	Unit:           []string{`kgs`, `stones`, `%`},
	NameSlug:       "marker-slug",
}

var module = &static_proto.Module{
	Id:          "111",
	Name:        "title",
	Summary:     "summary",
	Description: "description",
}

var category = &static_proto.BehaviourCategory{
	Id:            "111",
	Name:          "title",
	Summary:       "summary",
	Description:   "description",
	OrgId:         "orgid",
	Aims:          []*static_proto.BehaviourCategoryAim{{Name: "name"}},
	MarkerDefault: marker,
	MarkerOptions: []*static_proto.Marker{marker},
}

var socialType = &static_proto.SocialType{
	Id:   "111",
	Name: "title",
	Url:  "http://www.example.com",
}

var notification = &static_proto.Notification{
	Id:                   "111",
	ModuleId:             "moduleId",
	Name:                 "title",
	Description:          "description",
	Target:               static_proto.NotificationTarget_FACILITATOR,
	NameSlug:             "nameSlug",
	IconSlug:             "iconSlug",
	NotificationReminder: 10,
	Unit:                 "mins",
}

var trackerMethod = &static_proto.TrackerMethod{
	Id:       "111",
	Name:     "title",
	NameSlug: "nameSlug",
	IconSlug: "iconSlug",
}

var behaviourCategoryAim = &static_proto.BehaviourCategoryAim{
	Id:       "111",
	Name:     "title",
	NameSlug: "nameSlug",
	IconSlug: "iconSlug",
}

var contentParentCategory = &static_proto.ContentParentCategory{
	Id:          "111",
	Name:        "title",
	Summary:     "summary",
	Description: "description",
	IconSlug:    "iconslug",
	OrgId:       "orgid",
	Tags:        []string{"tag1", "tag2"},
}

var contentCategory = &static_proto.ContentCategory{
	Id:          "111",
	Name:        "title",
	NameSlug:    "sample_slug",
	Summary:     "summary",
	Description: "description",
	IconSlug:    "iconslug",
	OrgId:       "orgid",
	Parent:      []*static_proto.ContentParentCategory{contentParentCategory},
	// Actions:     []*static_proto.Action{action},
	Tags:           []string{"tag1", "tag2"},
	TrackerMethods: []*static_proto.TrackerMethod{trackerMethod},
}

var contentType = &static_proto.ContentType{
	Id:                "111",
	Name:              "title",
	Description:       "description",
	ContentTypeString: "contentTypeString",
	Tags:              []string{"tag1", "tag2"},
}

var contentSourceType = &static_proto.ContentSourceType{
	Id:          "111",
	Name:        "title",
	Description: "description",
	Tags:        []string{"tag1", "tag2"},
}

var moduleTrigger = &static_proto.ModuleTrigger{
	Id:           "111",
	Name:         "title",
	Summary:      "summary",
	IconSlug:     "icon_slug",
	Type:         static_proto.TriggerType_TIME,
	Module:       module,
	ContentTypes: []*static_proto.TriggerContentType{triggerContentType},
}

var triggerContentType = &static_proto.TriggerContentType{
	Id:       "111",
	Name:     "title",
	NameSlug: "name_slug",
}

var setback = &static_proto.Setback{
	Id:          "111",
	OrgId:       "orgid",
	Name:        "title",
	Description: "description",
}

func createApp(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.App {
	// create platform
	platform := createPlatform(ctx, hdlr, t)
	if platform == nil {
		return nil
	}
	app.Platforms[0] = platform

	req_create := &static_proto.CreateAppRequest{App: app}
	resp_create := &static_proto.CreateAppResponse{}
	err := hdlr.CreateApp(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.App
}

func createPlatform(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.Platform {
	req_create := &static_proto.CreatePlatformRequest{Platform: platform}
	resp_create := &static_proto.CreatePlatformResponse{}
	err := hdlr.CreatePlatform(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.Platform
}

func createWearable(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.Wearable {
	req_create := &static_proto.CreateWearableRequest{Wearable: wearable}
	resp_create := &static_proto.CreateWearableResponse{}
	err := hdlr.CreateWearable(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.Wearable
}

func createDevice(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.Device {
	req_create := &static_proto.CreateDeviceRequest{Device: device}
	resp_create := &static_proto.CreateDeviceResponse{}
	err := hdlr.CreateDevice(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.Device
}

func createMarker(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.Marker {
	// create apps
	app := createApp(ctx, hdlr, t)
	if app == nil {
		return nil
	}
	marker.Apps[0] = app
	time.Sleep(time.Second)
	// create wearables
	wearable := createWearable(ctx, hdlr, t)
	if wearable == nil {
		return nil
	}
	marker.Wearables[0] = wearable
	time.Sleep(time.Second)
	// create devices
	device := createDevice(ctx, hdlr, t)
	if device == nil {
		return nil
	}
	marker.Devices[0] = device
	time.Sleep(time.Second)
	// create trackerMethod
	method := createTrackerMethod(ctx, hdlr, t)
	if method == nil {
		return nil
	}
	marker.TrackerMethods = []*static_proto.TrackerMethod{method}
	time.Sleep(time.Second)

	req_create := &static_proto.CreateMarkerRequest{Marker: marker}
	resp_create := &static_proto.CreateMarkerResponse{}
	err := hdlr.CreateMarker(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.Marker
}

func createModule(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.Module {
	req_create := &static_proto.CreateModuleRequest{Module: module}
	resp_create := &static_proto.CreateModuleResponse{}
	err := hdlr.CreateModule(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.Module
}

func createBehaviourCategory(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.BehaviourCategory {
	// create aim
	aim := createBehaviourCategoryAim(ctx, hdlr, t)
	if aim == nil {
		return nil
	}
	category.Aims[0] = aim
	// create marker default
	marker := createMarker(ctx, hdlr, t)
	if marker == nil {
		return nil
	}
	category.MarkerDefault = marker
	// create marker option
	category.MarkerOptions[0] = marker

	req_create := &static_proto.CreateBehaviourCategoryRequest{Category: category}
	resp_create := &static_proto.CreateBehaviourCategoryResponse{}
	err := hdlr.CreateBehaviourCategory(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}

	return resp_create.Data.Category
}

func createSocialType(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.SocialType {
	req_create := &static_proto.CreateSocialTypeRequest{SocialType: socialType}
	resp_create := &static_proto.CreateSocialTypeResponse{}
	err := hdlr.CreateSocialType(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.SocialType
}

func createNotification(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.Notification {
	req_create := &static_proto.CreateNotificationRequest{Notification: notification}
	resp_create := &static_proto.CreateNotificationResponse{}
	err := hdlr.CreateNotification(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.Notification
}

func createTrackerMethod(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.TrackerMethod {
	req_create := &static_proto.CreateTrackerMethodRequest{TrackerMethod: trackerMethod}
	resp_create := &static_proto.CreateTrackerMethodResponse{}
	err := hdlr.CreateTrackerMethod(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.TrackerMethod
}

func createBehaviourCategoryAim(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.BehaviourCategoryAim {
	req_create := &static_proto.CreateBehaviourCategoryAimRequest{BehaviourCategoryAim: behaviourCategoryAim}
	resp_create := &static_proto.CreateBehaviourCategoryAimResponse{}
	err := hdlr.CreateBehaviourCategoryAim(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.BehaviourCategoryAim
}

func createContentParentCategory(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.ContentParentCategory {
	req_create := &static_proto.CreateContentParentCategoryRequest{ContentParentCategory: contentParentCategory}
	resp_create := &static_proto.CreateContentParentCategoryResponse{}
	err := hdlr.CreateContentParentCategory(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.ContentParentCategory
}

func createContentCategory(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.ContentCategory {
	// create parent
	parentCategory := createContentParentCategory(ctx, hdlr, t)
	if parentCategory == nil {
		return nil
	}
	contentCategory.Parent[0] = parentCategory
	// create trackerMethod
	method := createTrackerMethod(ctx, hdlr, t)
	if method == nil {
		return nil
	}
	contentCategory.TrackerMethods[0] = method

	req_create := &static_proto.CreateContentCategoryRequest{ContentCategory: contentCategory}
	resp_create := &static_proto.CreateContentCategoryResponse{}
	err := hdlr.CreateContentCategory(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.ContentCategory
}

func createContentType(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.ContentType {
	req_create := &static_proto.CreateContentTypeRequest{ContentType: contentType}
	resp_create := &static_proto.CreateContentTypeResponse{}
	err := hdlr.CreateContentType(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.ContentType
}

func createContentSourceType(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.ContentSourceType {
	req_create := &static_proto.CreateContentSourceTypeRequest{ContentSourceType: contentSourceType}
	resp_create := &static_proto.CreateContentSourceTypeResponse{}
	err := hdlr.CreateContentSourceType(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.ContentSourceType
}

func createModuleTrigger(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.ModuleTrigger {
	// create module
	module := createModule(ctx, hdlr, t)
	if module == nil {
		return nil
	}
	moduleTrigger.Module = module
	// create events
	// create contentTypes
	contentType := createTriggerContentType(ctx, hdlr, t)
	if contentType == nil {
		return nil
	}
	moduleTrigger.ContentTypes[0] = contentType

	req_create := &static_proto.CreateModuleTriggerRequest{ModuleTrigger: moduleTrigger}
	resp_create := &static_proto.CreateModuleTriggerResponse{}
	err := hdlr.CreateModuleTrigger(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.ModuleTrigger
}

func createTriggerContentType(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.TriggerContentType {
	req_create := &static_proto.CreateTriggerContentTypeRequest{TriggerContentType: triggerContentType}
	resp_create := &static_proto.CreateTriggerContentTypeResponse{}
	err := hdlr.CreateTriggerContentType(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.TriggerContentType
}

func createSetback(ctx context.Context, hdlr *StaticService, t *testing.T) *static_proto.Setback {
	req_create := &static_proto.CreateSetbackRequest{Setback: setback}
	resp_create := &static_proto.CreateSetbackResponse{}
	err := hdlr.CreateSetback(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}

	return resp_create.Data.Setback
}

func TestAllApps(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	app := createApp(ctx, hdlr, t)
	if app == nil {
		return
	}

	req_all := &static_proto.AllAppsRequest{}
	resp_all := &static_proto.AllAppsResponse{}
	err := hdlr.AllApps(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Apps) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Apps[0].Id != app.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_all.Data.Apps)
}

func TestReadApp(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	app := createApp(ctx, hdlr, t)
	if app == nil {
		return
	}

	req_read := &static_proto.ReadAppRequest{Id: app.Id}
	resp_read := &static_proto.ReadAppResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadApp(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.App == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.App.Id != app.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_read.Data.App)
}

func TestDeleteApp(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	app := createApp(ctx, hdlr, t)
	if app == nil {
		return
	}

	req_del := &static_proto.DeleteAppRequest{Id: app.Id}
	resp_del := &static_proto.DeleteAppResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteApp(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadAppRequest{Id: app.Id}
	resp_read := &static_proto.ReadAppResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadApp(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllPlatforms(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	platform := createPlatform(ctx, hdlr, t)
	if platform == nil {
		return
	}

	req_all := &static_proto.AllPlatformsRequest{}
	resp_all := &static_proto.AllPlatformsResponse{}
	err := hdlr.AllPlatforms(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Platforms) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Platforms[0].Id != platform.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_all.Data.Platforms)
}

func TestReadPlatform(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	platform := createPlatform(ctx, hdlr, t)
	if platform == nil {
		return
	}

	req_read := &static_proto.ReadPlatformRequest{Id: platform.Id}
	resp_read := &static_proto.ReadPlatformResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadPlatform(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Platform == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Platform.Id != platform.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeletePlatform(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	platform := createPlatform(ctx, hdlr, t)
	if platform == nil {
		return
	}

	req_del := &static_proto.DeletePlatformRequest{Id: platform.Id}
	resp_del := &static_proto.DeletePlatformResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeletePlatform(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadPlatformRequest{Id: platform.Id}
	resp_read := &static_proto.ReadPlatformResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadPlatform(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllWearables(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	w := createWearable(ctx, hdlr, t)
	if w == nil {
		return
	}

	req_all := &static_proto.AllWearablesRequest{}
	resp_all := &static_proto.AllWearablesResponse{}
	err := hdlr.AllWearables(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Wearables) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Wearables[0].Id != w.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadWearable(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	w := createWearable(ctx, hdlr, t)
	if w == nil {
		return
	}

	req_read := &static_proto.ReadWearableRequest{Id: w.Id}
	resp_read := &static_proto.ReadWearableResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadWearable(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Wearable == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Wearable.Id != w.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteWearable(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	w := createWearable(ctx, hdlr, t)
	if w == nil {
		return
	}

	req_del := &static_proto.DeleteWearableRequest{Id: w.Id}
	resp_del := &static_proto.DeleteWearableResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteWearable(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadWearableRequest{Id: w.Id}
	resp_read := &static_proto.ReadWearableResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadWearable(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllDevices(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	device := createDevice(ctx, hdlr, t)
	if device == nil {
		return
	}

	req_all := &static_proto.AllDevicesRequest{}
	resp_all := &static_proto.AllDevicesResponse{}
	err := hdlr.AllDevices(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Devices) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Devices[0].Id != device.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadDevice(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	device := createDevice(ctx, hdlr, t)
	if device == nil {
		return
	}

	req_read := &static_proto.ReadDeviceRequest{Id: device.Id}
	resp_read := &static_proto.ReadDeviceResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadDevice(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Device == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Device.Id != device.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteDevice(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	device := createDevice(ctx, hdlr, t)
	if device == nil {
		return
	}

	req_del := &static_proto.DeleteDeviceRequest{Id: device.Id}
	resp_del := &static_proto.DeleteDeviceResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteDevice(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadDeviceRequest{Id: device.Id}
	resp_read := &static_proto.ReadDeviceResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadDevice(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllMarkers(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	m := createMarker(ctx, hdlr, t)
	if m == nil {
		return
	}
	time.Sleep(time.Second)

	req_all := &static_proto.AllMarkersRequest{
		SortParameter: "created",
		SortDirection: "DESC",
	}
	resp_all := &static_proto.AllMarkersResponse{}
	err := hdlr.AllMarkers(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}

	if len(resp_all.Data.Markers) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Markers[0].Id != m.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadMarker(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	m := createMarker(ctx, hdlr, t)
	if m == nil {
		return
	}

	req_read := &static_proto.ReadMarkerRequest{Id: m.Id}
	resp_read := &static_proto.ReadMarkerResponse{}
	time.Sleep(time.Second)
	err := hdlr.ReadMarker(ctx, req_read, resp_read)
	if err != nil {
		t.Error(err)
		return
	}
	if resp_read.Data.Marker == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Marker.Id != m.Id {
		t.Error("Id does not match")
		return
	}

	t.Log(resp_read.Data.Marker)
}

func TestDeleteMarker(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	m := createMarker(ctx, hdlr, t)
	if m == nil {
		return
	}

	req_del := &static_proto.DeleteMarkerRequest{Id: m.Id}
	resp_del := &static_proto.DeleteMarkerResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteMarker(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadMarkerRequest{Id: m.Id}
	resp_read := &static_proto.ReadMarkerResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadMarker(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestFilterMarker(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	m := createMarker(ctx, hdlr, t)
	if m == nil {
		return
	}

	req_filter := &static_proto.FilterMarkerRequest{
		TrackerMethods: []string{m.TrackerMethods[0].Id},
	}
	resp_filter := &static_proto.FilterMarkerResponse{}
	err := hdlr.FilterMarker(ctx, req_filter, resp_filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_filter.Data.Markers) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_filter.Data.Markers[0].TrackerMethods[0].Id != m.TrackerMethods[0].Id {
		t.Error("Tracker Id does not match")
		return
	}

	t.Log(resp_filter.Data.Markers)
}

func TestAllModules(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	module := createModule(ctx, hdlr, t)
	if module == nil {
		return
	}

	req_all := &static_proto.AllModulesRequest{}
	resp_all := &static_proto.AllModulesResponse{}
	err := hdlr.AllModules(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Modules) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Modules[0].Id != module.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadModule(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	module := createModule(ctx, hdlr, t)
	if module == nil {
		return
	}

	req_read := &static_proto.ReadModuleRequest{Id: module.Id}
	resp_read := &static_proto.ReadModuleResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadModule(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Module == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Module.Id != module.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteModule(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	module := createModule(ctx, hdlr, t)
	if module == nil {
		return
	}

	req_del := &static_proto.DeleteModuleRequest{Id: module.Id}
	resp_del := &static_proto.DeleteModuleResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteModule(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadModuleRequest{Id: module.Id}
	resp_read := &static_proto.ReadModuleResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadModule(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllBehaviourCategories(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	category := createBehaviourCategory(ctx, hdlr, t)
	if category == nil {
		return
	}
	req_all := &static_proto.AllBehaviourCategoriesRequest{}
	resp_all := &static_proto.AllBehaviourCategoriesResponse{}
	err := hdlr.AllBehaviourCategories(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Categories) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Categories[0].Id != category.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_all.Data.Categories)
}

func TestReadBehaviourCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	category := createBehaviourCategory(ctx, hdlr, t)
	if category == nil {
		return
	}
	req_read := &static_proto.ReadBehaviourCategoryRequest{Id: category.Id}
	resp_read := &static_proto.ReadBehaviourCategoryResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadBehaviourCategory(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Category == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Category.Id != category.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_read.Data.Category)
}

func TestDeleteBehaviourCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	category := createBehaviourCategory(ctx, hdlr, t)
	if category == nil {
		return
	}

	req_del := &static_proto.DeleteBehaviourCategoryRequest{Id: category.Id}
	resp_del := &static_proto.DeleteBehaviourCategoryResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteBehaviourCategory(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadBehaviourCategoryRequest{Id: category.Id}
	resp_read := &static_proto.ReadBehaviourCategoryResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadBehaviourCategory(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestFilterBehaviourCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	category := createBehaviourCategory(ctx, hdlr, t)
	if category == nil {
		return
	}

	req_filter := &static_proto.FilterBehaviourCategoryRequest{
		Markers: []string{category.MarkerOptions[0].Id},
	}
	resp_filter := &static_proto.FilterBehaviourCategoryResponse{}
	err := hdlr.FilterBehaviourCategory(ctx, req_filter, resp_filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_filter.Data.Categories) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_filter.Data.Categories[0].MarkerOptions[0].Id != category.MarkerOptions[0].Id {
		t.Error("Marker Id does not match")
		return
	}
	t.Log(resp_filter.Data.Categories)
}

func TestAllSocialTypes(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	socialType := createSocialType(ctx, hdlr, t)
	if socialType == nil {
		return
	}

	req_all := &static_proto.AllSocialTypesRequest{}
	resp_all := &static_proto.AllSocialTypesResponse{}
	err := hdlr.AllSocialTypes(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.SocialTypes) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.SocialTypes[0].Id != socialType.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadSocialType(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	socialType := createSocialType(ctx, hdlr, t)
	if socialType == nil {
		return
	}

	req_read := &static_proto.ReadSocialTypeRequest{Id: socialType.Id}
	resp_read := &static_proto.ReadSocialTypeResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadSocialType(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.SocialType == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.SocialType.Id != socialType.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_read.Data.SocialType)
}

func TestDeleteSocialType(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	socialType := createSocialType(ctx, hdlr, t)
	if socialType == nil {
		return
	}

	req_del := &static_proto.DeleteSocialTypeRequest{Id: socialType.Id}
	resp_del := &static_proto.DeleteSocialTypeResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteSocialType(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadSocialTypeRequest{Id: socialType.Id}
	resp_read := &static_proto.ReadSocialTypeResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadSocialType(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllNotifications(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	notification := createNotification(ctx, hdlr, t)
	if notification == nil {
		return
	}

	req_all := &static_proto.AllNotificationsRequest{}
	resp_all := &static_proto.AllNotificationsResponse{}
	err := hdlr.AllNotifications(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Notifications) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Notifications[0].Id != notification.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadNotification(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	notification := createNotification(ctx, hdlr, t)
	if notification == nil {
		return
	}

	req_read := &static_proto.ReadNotificationRequest{Id: notification.Id}
	resp_read := &static_proto.ReadNotificationResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadNotification(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Notification == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Notification.Id != notification.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteNotification(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	notification := createNotification(ctx, hdlr, t)
	if notification == nil {
		return
	}

	req_del := &static_proto.DeleteNotificationRequest{Id: notification.Id}
	resp_del := &static_proto.DeleteNotificationResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteNotification(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadNotificationRequest{Id: notification.Id}
	resp_read := &static_proto.ReadNotificationResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadNotification(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllTrackerMethods(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	trackerMethod := createTrackerMethod(ctx, hdlr, t)
	if trackerMethod == nil {
		return
	}

	req_all := &static_proto.AllTrackerMethodsRequest{}
	resp_all := &static_proto.AllTrackerMethodsResponse{}
	err := hdlr.AllTrackerMethods(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.TrackerMethods) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.TrackerMethods[0].Id != trackerMethod.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadTrackerMethod(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	trackerMethod := createTrackerMethod(ctx, hdlr, t)
	if trackerMethod == nil {
		return
	}

	req_read := &static_proto.ReadTrackerMethodRequest{Id: trackerMethod.Id}
	resp_read := &static_proto.ReadTrackerMethodResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadTrackerMethod(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.TrackerMethod == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.TrackerMethod.Id != trackerMethod.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteTrackerMethod(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	trackerMethod := createTrackerMethod(ctx, hdlr, t)
	if trackerMethod == nil {
		return
	}

	req_del := &static_proto.DeleteTrackerMethodRequest{Id: trackerMethod.Id}
	resp_del := &static_proto.DeleteTrackerMethodResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteTrackerMethod(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadTrackerMethodRequest{Id: trackerMethod.Id}
	resp_read := &static_proto.ReadTrackerMethodResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadTrackerMethod(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestFilterTrackerMethod(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	trackerMethod := createTrackerMethod(ctx, hdlr, t)
	if trackerMethod == nil {
		return
	}
	marker := createMarker(ctx, hdlr, t)
	if marker == nil {
		return
	}

	req_filter := &static_proto.FilterTrackerMethodRequest{
		Markers: []string{marker.Id},
	}
	resp_filter := &static_proto.FilterTrackerMethodResponse{}
	time.Sleep(time.Second)

	err := hdlr.FilterTrackerMethod(ctx, req_filter, resp_filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_filter.Data.TrackerMethods) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_filter.Data.TrackerMethods[0].Id != "111" {
		t.Error("Tracker Id does not match")
		return
	}
}

func TestAllBehaviourCategoryAims(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	behaviourCategoryAim := createBehaviourCategoryAim(ctx, hdlr, t)
	if behaviourCategoryAim == nil {
		return
	}

	req_all := &static_proto.AllBehaviourCategoryAimsRequest{}
	resp_all := &static_proto.AllBehaviourCategoryAimsResponse{}
	err := hdlr.AllBehaviourCategoryAims(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.BehaviourCategoryAims) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.BehaviourCategoryAims[0].Id != behaviourCategoryAim.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadBehaviourCategoryAim(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	behaviourCategoryAim := createBehaviourCategoryAim(ctx, hdlr, t)
	if behaviourCategoryAim == nil {
		return
	}

	req_read := &static_proto.ReadBehaviourCategoryAimRequest{Id: behaviourCategoryAim.Id}
	resp_read := &static_proto.ReadBehaviourCategoryAimResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadBehaviourCategoryAim(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.BehaviourCategoryAim == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.BehaviourCategoryAim.Id != behaviourCategoryAim.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteBehaviourCategoryAim(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	behaviourCategoryAim := createBehaviourCategoryAim(ctx, hdlr, t)
	if behaviourCategoryAim == nil {
		return
	}

	req_del := &static_proto.DeleteBehaviourCategoryAimRequest{Id: behaviourCategoryAim.Id}
	resp_del := &static_proto.DeleteBehaviourCategoryAimResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteBehaviourCategoryAim(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadBehaviourCategoryAimRequest{Id: behaviourCategoryAim.Id}
	resp_read := &static_proto.ReadBehaviourCategoryAimResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadBehaviourCategoryAim(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllContentParentCategories(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	contentParentCategory := createContentParentCategory(ctx, hdlr, t)
	if contentParentCategory == nil {
		return
	}

	req_all := &static_proto.AllContentParentCategoriesRequest{}
	resp_all := &static_proto.AllContentParentCategoriesResponse{}
	err := hdlr.AllContentParentCategories(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.ContentParentCategories) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.ContentParentCategories[0].Id != contentParentCategory.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadContentParentCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	contentParentCategory := createContentParentCategory(ctx, hdlr, t)
	if contentParentCategory == nil {
		return
	}

	req_read := &static_proto.ReadContentParentCategoryRequest{Id: contentParentCategory.Id}
	resp_read := &static_proto.ReadContentParentCategoryResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadContentParentCategory(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.ContentParentCategory == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.ContentParentCategory.Id != contentParentCategory.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteContentParentCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	contentParentCategory := createContentParentCategory(ctx, hdlr, t)
	if contentParentCategory == nil {
		return
	}

	req_del := &static_proto.DeleteContentParentCategoryRequest{Id: contentParentCategory.Id}
	resp_del := &static_proto.DeleteContentParentCategoryResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteContentParentCategory(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadContentParentCategoryRequest{Id: contentParentCategory.Id}
	resp_read := &static_proto.ReadContentParentCategoryResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadContentParentCategory(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllContentCategories(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	contentCategory := createContentCategory(ctx, hdlr, t)
	if contentCategory == nil {
		return
	}

	req_all := &static_proto.AllContentCategoriesRequest{}
	resp_all := &static_proto.AllContentCategoriesResponse{}
	err := hdlr.AllContentCategories(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.ContentCategories) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.ContentCategories[0].Id != contentCategory.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_all.Data.ContentCategories)
}

func TestReadContentCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	contentCategory := createContentCategory(ctx, hdlr, t)
	if contentCategory == nil {
		return
	}

	req_read := &static_proto.ReadContentCategoryRequest{Id: contentCategory.Id}
	resp_read := &static_proto.ReadContentCategoryResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadContentCategory(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.ContentCategory == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.ContentCategory.Id != contentCategory.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteContentCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	contentCategory := createContentCategory(ctx, hdlr, t)
	if contentCategory == nil {
		return
	}

	req_del := &static_proto.DeleteContentCategoryRequest{Id: contentCategory.Id}
	resp_del := &static_proto.DeleteContentCategoryResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteContentCategory(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadContentCategoryRequest{Id: contentCategory.Id}
	resp_read := &static_proto.ReadContentCategoryResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadContentCategory(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestReadContentCategoryByNameslug(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	createContentCategory(ctx, hdlr, t)

	req_read := &static_proto.ReadByNameslugRequest{NameSlug: "sample_slug"}
	resp_read := &static_proto.ReadContentCategoryByNameslugResponse{}
	time.Sleep(time.Second)

	res_read := hdlr.ReadContentCategoryByNameslug(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.ContentCategory == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.ContentCategory.NameSlug != "sample_slug" {
		t.Error("Nameslug does not match")
		return
	}
}
func TestAllContentTypes(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	contentType := createContentType(ctx, hdlr, t)
	if contentType == nil {
		return
	}

	req_all := &static_proto.AllContentTypesRequest{}
	resp_all := &static_proto.AllContentTypesResponse{}
	err := hdlr.AllContentTypes(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.ContentTypes) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.ContentTypes[0].Id != contentType.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadContentType(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	contentType := createContentType(ctx, hdlr, t)
	if contentType == nil {
		return
	}

	req_read := &static_proto.ReadContentTypeRequest{Id: contentType.Id}
	resp_read := &static_proto.ReadContentTypeResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadContentType(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.ContentType == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.ContentType.Id != contentType.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteContentType(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	contentType := createContentType(ctx, hdlr, t)
	if contentType == nil {
		return
	}

	req_del := &static_proto.DeleteContentTypeRequest{Id: contentType.Id}
	resp_del := &static_proto.DeleteContentTypeResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteContentType(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadContentTypeRequest{Id: contentType.Id}
	resp_read := &static_proto.ReadContentTypeResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadContentType(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllContentSourceTypes(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	contentSourceType := createContentSourceType(ctx, hdlr, t)
	if contentSourceType == nil {
		return
	}

	req_all := &static_proto.AllContentSourceTypesRequest{}
	resp_all := &static_proto.AllContentSourceTypesResponse{}
	err := hdlr.AllContentSourceTypes(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.ContentSourceTypes) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.ContentSourceTypes[0].Id != contentSourceType.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadContentSourceType(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	contentSourceType := createContentSourceType(ctx, hdlr, t)
	if contentSourceType == nil {
		return
	}

	req_read := &static_proto.ReadContentSourceTypeRequest{Id: contentSourceType.Id}
	resp_read := &static_proto.ReadContentSourceTypeResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadContentSourceType(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.ContentSourceType == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.ContentSourceType.Id != contentSourceType.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteContentSourceType(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	contentSourceType := createContentSourceType(ctx, hdlr, t)
	if contentSourceType == nil {
		return
	}

	req_del := &static_proto.DeleteContentSourceTypeRequest{Id: contentSourceType.Id}
	resp_del := &static_proto.DeleteContentSourceTypeResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteContentSourceType(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadContentSourceTypeRequest{Id: contentSourceType.Id}
	resp_read := &static_proto.ReadContentSourceTypeResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadContentSourceType(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllModuleTriggers(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	moduleTrigger := createModuleTrigger(ctx, hdlr, t)
	if moduleTrigger == nil {
		return
	}

	req_all := &static_proto.AllModuleTriggersRequest{}
	resp_all := &static_proto.AllModuleTriggersResponse{}
	err := hdlr.AllModuleTriggers(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.ModuleTriggers) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.ModuleTriggers[0].Id != moduleTrigger.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadModuleTrigger(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	moduleTrigger := createModuleTrigger(ctx, hdlr, t)
	if moduleTrigger == nil {
		return
	}

	req_read := &static_proto.ReadModuleTriggerRequest{Id: moduleTrigger.Id}
	resp_read := &static_proto.ReadModuleTriggerResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadModuleTrigger(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.ModuleTrigger == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.ModuleTrigger.Id != moduleTrigger.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteModuleTrigger(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	moduleTrigger := createModuleTrigger(ctx, hdlr, t)
	if moduleTrigger == nil {
		return
	}

	req_del := &static_proto.DeleteModuleTriggerRequest{Id: moduleTrigger.Id}
	resp_del := &static_proto.DeleteModuleTriggerResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteModuleTrigger(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadModuleTriggerRequest{Id: moduleTrigger.Id}
	resp_read := &static_proto.ReadModuleTriggerResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadModuleTrigger(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}
func TestFilterModuleTrigger(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	moduleTrigger := createModuleTrigger(ctx, hdlr, t)
	if moduleTrigger == nil {
		return
	}

	req_filter := &static_proto.FilterModuleTriggerRequest{
		Module:      []string{moduleTrigger.Id},
		TriggerType: []int64{int64(static_proto.TriggerType_TIME)},
	}
	resp_filter := &static_proto.FilterModuleTriggerResponse{}
	time.Sleep(2 * time.Second)
	err := hdlr.FilterModuleTrigger(ctx, req_filter, resp_filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_filter.Data.ModuleTriggers) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_filter.Data.ModuleTriggers[0].Module.Id != module.Id {
		t.Error("Tracker Id does not match")
		return
	}
}

func TestAllTriggerContentTypes(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	triggerContentType := createTriggerContentType(ctx, hdlr, t)
	if triggerContentType == nil {
		return
	}

	req_all := &static_proto.AllTriggerContentTypesRequest{}
	resp_all := &static_proto.AllTriggerContentTypesResponse{}
	err := hdlr.AllTriggerContentTypes(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.TriggerContentTypes) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.TriggerContentTypes[0].Id != triggerContentType.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadTriggerContentType(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	triggerContentType := createTriggerContentType(ctx, hdlr, t)
	if triggerContentType == nil {
		return
	}

	req_read := &static_proto.ReadTriggerContentTypeRequest{Id: triggerContentType.Id}
	resp_read := &static_proto.ReadTriggerContentTypeResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadTriggerContentType(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.TriggerContentType == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.TriggerContentType.Id != triggerContentType.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteTriggerContentType(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	triggerContentType := createTriggerContentType(ctx, hdlr, t)
	if triggerContentType == nil {
		return
	}

	req_del := &static_proto.DeleteTriggerContentTypeRequest{Id: triggerContentType.Id}
	resp_del := &static_proto.DeleteTriggerContentTypeResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteTriggerContentType(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadTriggerContentTypeRequest{Id: triggerContentType.Id}
	resp_read := &static_proto.ReadTriggerContentTypeResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadTriggerContentType(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestFilterTriggerContentType(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	triggerContentType := createTriggerContentType(ctx, hdlr, t)
	if triggerContentType == nil {
		return
	}

	req_filter := &static_proto.FilterTriggerContentTypeRequest{
		ModuleTrigger: []string{"111"},
	}
	resp_filter := &static_proto.FilterTriggerContentTypeResponse{}
	time.Sleep(2 * time.Second)
	err := hdlr.FilterTriggerContentType(ctx, req_filter, resp_filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_filter.Data.TriggerContentTypes) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_filter.Data.TriggerContentTypes[0].Id != triggerContentType.Id {
		t.Error("Tracker Id does not match")
		return
	}
}

func TestAllSetbacks(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	setback := createSetback(ctx, hdlr, t)
	if setback == nil {
		return
	}

	req_all := &static_proto.AllSetbacksRequest{}
	resp_all := &static_proto.AllSetbacksResponse{}
	err := hdlr.AllSetbacks(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Setbacks) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Setbacks[0].Id != setback.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadSetback(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	setback := createSetback(ctx, hdlr, t)
	if setback == nil {
		return
	}

	req_read := &static_proto.ReadSetbackRequest{Id: setback.Id}
	resp_read := &static_proto.ReadSetbackResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadSetback(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Setback == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Setback.Id != setback.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteSetback(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	setback := createSetback(ctx, hdlr, t)
	if setback == nil {
		return
	}

	req_del := &static_proto.DeleteSetbackRequest{Id: setback.Id}
	resp_del := &static_proto.DeleteSetbackResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteSetback(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
	}

	req_read := &static_proto.ReadSetbackRequest{Id: setback.Id}
	resp_read := &static_proto.ReadSetbackResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadSetback(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAutocompleteSetbackSearch(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	setback := createSetback(ctx, hdlr, t)
	if setback == nil {
		return
	}

	req_search := &static_proto.AutocompleteSetbackSearchRequest{"it"}
	resp_search := &static_proto.AllSetbacksResponse{}
	err := hdlr.AutocompleteSetbackSearch(ctx, req_search, resp_search)
	if err != nil {
		t.Error(err)
	}

	if len(resp_search.Data.Setbacks) == 0 {
		t.Error("Object count does not match")
		return
	}
}

func TestReadMarkerByNameslug(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(StaticService)

	createMarker(ctx, hdlr, t)

	req_read := &static_proto.ReadByNameslugRequest{NameSlug: marker.NameSlug}
	resp_read := &static_proto.ReadMarkerResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadMarkerByNameslug(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Marker == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Marker.Id != marker.Id {
		t.Error("Id does not match")
		return
	}
}
