package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server/api/utils"
	"server/common"
	"server/static-srv/db"
	static_proto "server/static-srv/proto/static"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

// var serverURL = "http://localhost:8080"
var staticURL = "/server/static"

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
}

var module = &static_proto.Module{
	Id:          "111",
	Name:        "title",
	Summary:     "summary",
	Description: "description",
}

var category = &static_proto.BehaviourCategory{
	Id:          "111",
	Name:        "title",
	Summary:     "summary",
	Description: "description",
	OrgId:       "orgid",
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
	NameSlug:             "hcp",
	IconSlug:             "iconSlug",
	NotificationReminder: 10,
	Unit:                 "mins",
}

var trackerMethod = &static_proto.TrackerMethod{
	Id:       "111",
	Name:     "title",
	NameSlug: "hcp",
	IconSlug: "iconSlug",
}

var behaviourCategoryAim = &static_proto.BehaviourCategoryAim{
	Id:       "111",
	Name:     "title",
	NameSlug: "hcp",
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
	Summary:     "summary",
	Description: "description",
	IconSlug:    "iconslug",
	OrgId:       "orgid",
	Parent:      []*static_proto.ContentParentCategory{contentParentCategory},
	Tags:        []string{"tag1", "tag2"},
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
	Name:        "title",
	OrgId:       "orgid",
	Description: "description",
}

func initStaticDb() {
	cl := client.NewClient(client.Transport(nats_transport.NewTransport()), client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(5*time.Second),
		client.Retries(5))
	// ctx := common.NewTestContext(context.TODO())
	// db.DbStaticName = common.TestingName("healum_test")
	// db.DbAppTable = common.TestingName("app_test")
	// db.DbPlatformTable = common.TestingName("platform_test")
	// db.DbWearableTable = common.TestingName("wearable_test")
	// db.DbDeviceTable = common.TestingName("device_test")
	// db.DbMarkerTable = common.TestingName("marker_test")
	// db.DbModuleTable = common.TestingName("module_test")
	// db.DbBehaviourCategoryTable = common.TestingName("behaviour_category_test")
	// db.DbSocialTypeTable = common.TestingName("social_type_test")
	// db.DbNotificationTable = common.TestingName("notification_test")
	// db.DbTrackerMethodTable = common.TestingName("tracker_method_test")
	// db.DbBehaviourCategoryAimTable = common.TestingName("behaviour_category_aim_test")
	// db.DbMarkerTrackerEdgeTable = common.TestingName("marker_tracker_method_edge_test")
	// db.DbModuleTriggerTable = common.TestingName("module_trigger_test")
	// db.DbStaticDriver = "arangodb"
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func createApp(app *static_proto.App, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"app": app})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/app?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createPlatform(platform *static_proto.Platform, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"platform": platform})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/platform?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createWearable(wearable *static_proto.Wearable, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"wearable": wearable})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/wearable?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createDevice(device *static_proto.Device, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"device": device})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/device?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createMarker(marker *static_proto.Marker, t *testing.T) *static_proto.Marker {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"marker": marker})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/marker?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return nil
	}
	time.Sleep(time.Second)

	r := static_proto.CreateMarkerResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data == nil {
		t.Errorf("Response does not matched")
		return nil
	}

	marker = r.Data.Marker
	return marker
}

func createModule(module *static_proto.Module, t *testing.T) *static_proto.Module {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"module": module})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/module?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.CreateModuleResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data == nil {
		t.Errorf("Response does not matched")
		return nil
	}

	module = r.Data.Module
	return module
}

func createCategory(category *static_proto.BehaviourCategory, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"category": category})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/behaviour/category?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createSocialType(socialType *static_proto.SocialType, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"socialType": socialType})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/socialType?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createNotification(notification *static_proto.Notification, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"notification": notification})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/notification?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createTrackerMethod(trackerMethod *static_proto.TrackerMethod, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"trackerMethod": trackerMethod})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/trackerMethod?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createBehaviourCategoryAim(behaviourCategoryAim *static_proto.BehaviourCategoryAim, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"behaviourCategoryAim": behaviourCategoryAim})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/behaviourCategoryAim?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createContentParentCategory(contentParentCategory *static_proto.ContentParentCategory, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"contentParentCategory": contentParentCategory})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/content/category/parent?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createContentCategory(contentCategory *static_proto.ContentCategory, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"contentCategory": contentCategory})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/content/category?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createContentType(contentType *static_proto.ContentType, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"contentType": contentType})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/content/type?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createContentSourceType(contentSourceType *static_proto.ContentSourceType, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"contentSourceType": contentSourceType})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/content/source/type?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createModuleTrigger(moduleTrigger *static_proto.ModuleTrigger, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"moduleTrigger": moduleTrigger})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/module/trigger?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createTriggerContentType(triggerContentType *static_proto.TriggerContentType, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"triggerContentType": triggerContentType})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/trigger/content/type?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func createSetback(setback *static_proto.Setback, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"setback": setback})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/setback?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func TestAllApps(t *testing.T) {
	initStaticDb()

	createApp(app, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/apps/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllAppsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Apps) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadApp(t *testing.T) {
	initStaticDb()

	createApp(app, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/app/"+app.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadAppResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.App == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.App.Id != app.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteApp(t *testing.T) {
	initStaticDb()

	createApp(app, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/app/"+app.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/app/"+app.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadAppResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadApp(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/app/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllPlatforms(t *testing.T) {
	initStaticDb()

	createPlatform(platform, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/platforms/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllPlatformsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Platforms) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadPlatform(t *testing.T) {
	initStaticDb()

	createPlatform(platform, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/platform/"+platform.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadPlatformResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Platform == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.Platform.Id != platform.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeletePlatform(t *testing.T) {
	initStaticDb()

	createPlatform(platform, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/platform/"+platform.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/platform/"+platform.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadPlatformResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadPlatform(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/platform/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllWearables(t *testing.T) {
	initStaticDb()

	createWearable(wearable, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/wearables/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllWearablesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Wearables) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadWearable(t *testing.T) {
	initStaticDb()

	createWearable(wearable, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/wearable/"+wearable.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadWearableResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Wearable == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.Wearable.Id != wearable.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteWearable(t *testing.T) {
	initStaticDb()

	createWearable(wearable, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/wearable/"+wearable.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/wearable/"+wearable.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadWearableResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadWearable(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/wearable/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllDevices(t *testing.T) {
	initStaticDb()

	createDevice(device, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/devices/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllDevicesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Devices) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadDevice(t *testing.T) {
	initStaticDb()

	createDevice(device, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/device/"+device.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadDeviceResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Device == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.Device.Id != device.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteDevice(t *testing.T) {
	initStaticDb()

	createDevice(device, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/device/"+device.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/device/"+device.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadDeviceResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadDevice(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/device/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllMarkers(t *testing.T) {
	initStaticDb()

	createMarker(marker, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/markers/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllMarkersResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Markers) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadMarker(t *testing.T) {
	initStaticDb()

	marker := createMarker(marker, t)
	if marker == nil {
		t.Error("Marker invalid")
		return
	}

	t.Log("marker:", marker, marker.Id)
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/marker/"+marker.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadMarkerResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Marker == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.Marker.Id != marker.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteMarker(t *testing.T) {
	initStaticDb()

	marker := createMarker(marker, t)
	if marker == nil {
		t.Error("Marker invalid")
		return
	}

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/marker/"+marker.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/marker/"+marker.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := static_proto.ReadMarkerResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
		return
	}
}

func TestErrReadMarker(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/marker/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllModules(t *testing.T) {
	initStaticDb()

	createModule(module, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/modules/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllModulesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Modules) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadModule(t *testing.T) {
	initStaticDb()

	createModule(module, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/module/"+module.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadModuleResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Module == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.Module.Id != module.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteModule(t *testing.T) {
	initStaticDb()

	createModule(module, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/module/"+module.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/module/"+module.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadModuleResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadModule(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/module/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllBehaviourCategories(t *testing.T) {
	initStaticDb()

	createCategory(category, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/behaviour/categorys/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllBehaviourCategoriesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Categories) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadBehaviourCategory(t *testing.T) {
	initStaticDb()

	createCategory(category, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/behaviour/category/"+category.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadBehaviourCategoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.Category == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Category.Id != category.Id {
		t.Errorf("Object Id does not matched")
		return
	}
}

func TestDeleteBehaviourCategory(t *testing.T) {
	initStaticDb()

	createCategory(category, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/behaviour/category/"+category.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/behaviour/category/"+category.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadBehaviourCategoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadBehaviourCategory(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/behaviour/category/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllSocialTypes(t *testing.T) {
	initStaticDb()

	createSocialType(socialType, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/socialTypes/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllSocialTypesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.SocialTypes) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadSocialType(t *testing.T) {
	initStaticDb()

	createSocialType(socialType, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/socialType/"+socialType.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadSocialTypeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.SocialType == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.SocialType.Id != socialType.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteSocialType(t *testing.T) {
	initStaticDb()

	createSocialType(socialType, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/socialType/"+socialType.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/socialType/"+socialType.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadSocialTypeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadSocialType(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/socialType/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllNotifications(t *testing.T) {
	initStaticDb()

	createNotification(notification, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/notifications/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllNotificationsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Notifications) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadNotification(t *testing.T) {
	initStaticDb()

	createNotification(notification, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/notification/"+notification.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadNotificationResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Notification == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.Notification.Id != notification.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteNotification(t *testing.T) {
	initStaticDb()

	createNotification(notification, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/notification/"+notification.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/notification/"+notification.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadNotificationResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadNotification(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/notification/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllTrackerMethods(t *testing.T) {
	initStaticDb()

	createTrackerMethod(trackerMethod, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/trackerMethods/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllTrackerMethodsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.TrackerMethods) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadTrackerMethod(t *testing.T) {
	initStaticDb()

	createTrackerMethod(trackerMethod, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/trackerMethod/"+trackerMethod.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadTrackerMethodResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.TrackerMethod == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.TrackerMethod.Id != trackerMethod.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteTrackerMethod(t *testing.T) {
	initStaticDb()

	createTrackerMethod(trackerMethod, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/trackerMethod/"+trackerMethod.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/trackerMethod/"+trackerMethod.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadTrackerMethodResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestFilterTrackerMethod(t *testing.T) {
	initStaticDb()

	createTrackerMethod(trackerMethod, t)
	createMarker(marker, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"markers": []string{"111"}})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+staticURL+"/trackerMethod/filter?session="+sessionId+"&offset=0&limit=10", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.FilterTrackerMethodResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data.TrackerMethods[0].Id != trackerMethod.Id {
		t.Errorf("TrackerMethod does not matched")
		return
	}
}

func TestErrReadTrackerMethod(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/trackerMethod/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllBehaviourCategoryAims(t *testing.T) {
	initStaticDb()

	createBehaviourCategoryAim(behaviourCategoryAim, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/behaviourCategoryAims/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllBehaviourCategoryAimsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.BehaviourCategoryAims) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadBehaviourCategoryAim(t *testing.T) {
	initStaticDb()

	createBehaviourCategoryAim(behaviourCategoryAim, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/behaviourCategoryAim/"+behaviourCategoryAim.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadBehaviourCategoryAimResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.BehaviourCategoryAim == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.BehaviourCategoryAim.Id != behaviourCategoryAim.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteBehaviourCategoryAim(t *testing.T) {
	initStaticDb()

	createBehaviourCategoryAim(behaviourCategoryAim, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/behaviourCategoryAim/"+behaviourCategoryAim.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/behaviourCategoryAim/"+behaviourCategoryAim.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadBehaviourCategoryAimResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadBehaviourCategoryAim(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/behaviourCategoryAim/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllContentParentCategories(t *testing.T) {
	initStaticDb()

	createContentParentCategory(contentParentCategory, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/content/category/parents/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllContentParentCategoriesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.ContentParentCategories) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadContentParentCategory(t *testing.T) {
	initStaticDb()

	createContentParentCategory(contentParentCategory, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/content/category/parent/"+contentParentCategory.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadContentParentCategoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.ContentParentCategory == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.ContentParentCategory.Id != contentParentCategory.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteContentParentCategory(t *testing.T) {
	initStaticDb()

	createContentParentCategory(contentParentCategory, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/content/category/parent/"+contentParentCategory.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/content/category/parent/"+contentParentCategory.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadContentParentCategoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadContentParentCategory(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/content/category/parent/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllContentCategories(t *testing.T) {
	initStaticDb()

	createContentCategory(contentCategory, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/content/categorys/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllContentCategoriesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.ContentCategories) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadContentCategory(t *testing.T) {
	initStaticDb()

	createContentCategory(contentCategory, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/content/category/"+contentCategory.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadContentCategoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.ContentCategory == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.ContentCategory.Id != contentCategory.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteContentCategory(t *testing.T) {
	initStaticDb()

	createContentCategory(contentCategory, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/content/category/"+contentCategory.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/content/category/"+contentCategory.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadContentCategoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadContentCategory(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/content/category/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllContentTypes(t *testing.T) {
	initStaticDb()

	createContentType(contentType, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/content/types/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllContentTypesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.ContentTypes) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadContentType(t *testing.T) {
	initStaticDb()

	createContentType(contentType, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/content/type/"+contentType.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadContentTypeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.ContentType == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.ContentType.Id != contentType.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteContentType(t *testing.T) {
	initStaticDb()

	createContentType(contentType, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/content/type/"+contentType.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/content/type/"+contentType.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadContentTypeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadContentType(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/content/type/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllContentSourceTypes(t *testing.T) {
	initStaticDb()

	createContentSourceType(contentSourceType, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/content/source/types/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllContentSourceTypesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.ContentSourceTypes) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadContentSourceType(t *testing.T) {
	initStaticDb()

	createContentSourceType(contentSourceType, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/content/source/type/"+contentSourceType.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadContentSourceTypeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.ContentSourceType == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.ContentSourceType.Id != contentSourceType.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteContentSourceType(t *testing.T) {
	initStaticDb()

	createContentSourceType(contentSourceType, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/content/source/type/"+contentSourceType.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/content/source/type/"+contentSourceType.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadContentSourceTypeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadContentSourceType(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/content/source/type/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllModuleTriggers(t *testing.T) {
	initStaticDb()

	createModuleTrigger(moduleTrigger, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/module/triggers/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllModuleTriggersResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.ModuleTriggers) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadModuleTrigger(t *testing.T) {
	initStaticDb()

	createModuleTrigger(moduleTrigger, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/module/trigger/"+moduleTrigger.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadModuleTriggerResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.ModuleTrigger == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.ModuleTrigger.Id != moduleTrigger.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteModuleTrigger(t *testing.T) {
	initStaticDb()

	createModuleTrigger(moduleTrigger, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/module/trigger/"+moduleTrigger.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/module/trigger/"+moduleTrigger.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadModuleTriggerResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadModuleTrigger(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/module/trigger/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllTriggerContentTypes(t *testing.T) {
	initStaticDb()

	createTriggerContentType(triggerContentType, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/trigger/content/types/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllTriggerContentTypesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.TriggerContentTypes) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadTriggerContentType(t *testing.T) {
	initStaticDb()

	createTriggerContentType(triggerContentType, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/trigger/content/type/"+triggerContentType.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadTriggerContentTypeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.TriggerContentType == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.TriggerContentType.Id != triggerContentType.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteTriggerContentType(t *testing.T) {
	initStaticDb()

	createTriggerContentType(triggerContentType, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/trigger/content/type/"+triggerContentType.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/trigger/content/type/"+triggerContentType.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadTriggerContentTypeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadTriggerContentType(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/trigger/content/type/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAllSetbacks(t *testing.T) {
	initStaticDb()

	createSetback(setback, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/setbacks/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.AllSetbacksResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Setbacks) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReadSetback(t *testing.T) {
	initStaticDb()

	createSetback(setback, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/setback/"+setback.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadSetbackResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Setback == nil {
		t.Errorf("Object does not matched")
	}
	if r.Data.Setback.Id != setback.Id {
		t.Errorf("Object Id does not matched")
	}
}

func TestDeleteSetback(t *testing.T) {
	initStaticDb()

	createSetback(setback, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+staticURL+"/setback/"+setback.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+staticURL+"/setback/"+setback.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadSetbackResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
	}
}

func TestErrReadSetback(t *testing.T) {
	initStaticDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+staticURL+"/setback/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
	} else {
		t.Log(r)
	}
}

func TestAutocompleteSetbackSearch(t *testing.T) {
	initStaticDb()

	createSetback(setback, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	jsonStr, err := json.Marshal(map[string]interface{}{"title": "it"})
	if err != nil {
		t.Error(err)
		return
	}
	// Send a POST request.
	req, err := http.NewRequest("POST", serverURL+staticURL+"/setback/search/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := static_proto.AllSetbacksResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Setbacks) == 0 {
		t.Errorf("Object count does not matched")
	}
}
