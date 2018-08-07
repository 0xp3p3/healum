package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"server/common"
	db_proto "server/db-srv/proto/db"
	static_proto "server/static-srv/proto/static"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/client"
	"github.com/pborman/uuid"
)

type clientWrapper struct {
	Db_client db_proto.DBClient
}

var (
	ClientWrapper *clientWrapper
	ErrNotFound   = errors.New("not found")
)

// Storage for a db microservice client
func NewClientWrapper(serviceClient client.Client) *clientWrapper {
	cl := db_proto.NewDBClient("", serviceClient)

	return &clientWrapper{
		Db_client: cl,
	}
}

// Init initializes healum databases
func Init(serviceClient client.Client) error {
	ClientWrapper = NewClientWrapper(serviceClient)
	// if _, err := ClientWrapper.Db_client.Init(context.TODO(), &db_proto.InitRequest{}); err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	return nil
}

// RemoveDb removes healum database (for testing)
func RemoveDb(ctx context.Context, serviceClient client.Client) error {
	ClientWrapper = NewClientWrapper(serviceClient)
	if _, err := ClientWrapper.Db_client.RemoveDb(ctx, &db_proto.RemoveDbRequest{}); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func runQuery(ctx context.Context, q string, table string) (*db_proto.RunQueryResponse, error) {
	return ClientWrapper.Db_client.RunQuery(ctx, &db_proto.RunQueryRequest{
		Database: &db_proto.Database{
			Name:     common.DbHealumName,
			Table:    table,
			Driver:   common.DbHealumDriver,
			Metadata: common.SearchableMetaMap,
		},
		Query: q,
	})
}

func appToRecord(app *static_proto.App) (string, error) {
	data, err := common.MarhalToObject(app)
	if err != nil {
		return "", err
	}

	if len(app.Platforms) > 0 {
		var arr []interface{}
		for _, item := range app.Platforms {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["platforms"] = arr
	} else {
		delete(data, "platforms")
	}

	d := map[string]interface{}{
		"_key":    app.Id,
		"id":      app.Id,
		"created": app.Created,
		"updated": app.Updated,
		"name":    app.Name,
		// "parameter1": app.OrgId,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToApp(r *db_proto.Record) (*static_proto.App, error) {
	var p static_proto.App
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func platformToRecord(platform *static_proto.Platform) (string, error) {
	data, err := common.MarhalToObject(platform)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    platform.Id,
		"id":      platform.Id,
		"created": platform.Created,
		"updated": platform.Updated,
		"name":    platform.Name,
		// "parameter1": platform.OrgId,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToPlatform(r *db_proto.Record) (*static_proto.Platform, error) {
	var p static_proto.Platform
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func wearableToRecord(wearable *static_proto.Wearable) (string, error) {
	data, err := common.MarhalToObject(wearable)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    wearable.Id,
		"id":      wearable.Id,
		"created": wearable.Created,
		"updated": wearable.Updated,
		"name":    wearable.Name,
		// "parameter1": wearable.OrgId,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToWearable(r *db_proto.Record) (*static_proto.Wearable, error) {
	var p static_proto.Wearable
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func deviceToRecord(device *static_proto.Device) (string, error) {
	data, err := common.MarhalToObject(device)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    device.Id,
		"id":      device.Id,
		"created": device.Created,
		"updated": device.Updated,
		"name":    device.Name,
		// "parameter1": device.OrgId,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToDevice(r *db_proto.Record) (*static_proto.Device, error) {
	var p static_proto.Device
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func markerToRecord(marker *static_proto.Marker) (string, error) {
	data, err := common.MarhalToObject(marker)
	if err != nil {
		return "", err
	}

	if len(marker.Apps) > 0 {
		var arr []interface{}
		for _, item := range marker.Apps {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["apps"] = arr
	} else {
		delete(data, "apps")
	}

	if len(marker.Wearables) > 0 {
		var arr []interface{}
		for _, item := range marker.Wearables {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["wearables"] = arr
	} else {
		delete(data, "wearables")
	}

	if len(marker.Devices) > 0 {
		var arr []interface{}
		for _, item := range marker.Devices {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["devices"] = arr
	} else {
		delete(data, "devices")
	}

	delete(data, "trackerMethods")

	d := map[string]interface{}{
		"_key":       marker.Id,
		"id":         marker.Id,
		"created":    marker.Created,
		"updated":    marker.Updated,
		"name":       marker.Name,
		"parameter1": marker.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToMarker(r *db_proto.Record) (*static_proto.Marker, error) {
	var p static_proto.Marker
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func moduleToRecord(module *static_proto.Module) (string, error) {
	data, err := common.MarhalToObject(module)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    module.Id,
		"id":      module.Id,
		"created": module.Created,
		"updated": module.Updated,
		"name":    module.Name,
		"data":    data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToModule(r *db_proto.Record) (*static_proto.Module, error) {
	var p static_proto.Module
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func categoryToRecord(category *static_proto.BehaviourCategory) (string, error) {
	data, err := common.MarhalToObject(category)
	if err != nil {
		return "", err
	}
	if len(category.Aims) > 0 {
		var arr []interface{}
		for _, item := range category.Aims {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["aims"] = arr
	} else {
		delete(data, "aims")
	}

	common.FilterObject(data, "markerDefault", category.MarkerDefault)

	if len(category.MarkerOptions) > 0 {
		var arr []interface{}
		for _, item := range category.MarkerOptions {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["markerOptions"] = arr
	} else {
		delete(data, "markerOptions")
	}

	d := map[string]interface{}{
		"_key":       category.Id,
		"id":         category.Id,
		"created":    category.Created,
		"updated":    category.Updated,
		"name":       category.Name,
		"parameter1": category.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToCategory(r *db_proto.Record) (*static_proto.BehaviourCategory, error) {
	var p static_proto.BehaviourCategory
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func socialTypeToRecord(socialType *static_proto.SocialType) (string, error) {
	data, err := common.MarhalToObject(socialType)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    socialType.Id,
		"id":      socialType.Id,
		"created": socialType.Created,
		"updated": socialType.Updated,
		"name":    socialType.Name,
		// "parameter1": socialType.OrgId,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToSocialType(r *db_proto.Record) (*static_proto.SocialType, error) {
	var p static_proto.SocialType
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func notificationToRecord(notification *static_proto.Notification) (string, error) {
	data, err := common.MarhalToObject(notification)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    notification.Id,
		"id":      notification.Id,
		"created": notification.Created,
		"updated": notification.Updated,
		"name":    notification.Name,
		// "parameter1": notification.OrgId,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToNotification(r *db_proto.Record) (*static_proto.Notification, error) {
	var p static_proto.Notification
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func trackerMethodToRecord(trackerMethod *static_proto.TrackerMethod) (string, error) {
	data, err := common.MarhalToObject(trackerMethod)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    trackerMethod.Id,
		"id":      trackerMethod.Id,
		"created": trackerMethod.Created,
		"updated": trackerMethod.Updated,
		"name":    trackerMethod.Name,
		// "parameter1": trackerMethod.OrgId,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToTrackerMethod(r *db_proto.Record) (*static_proto.TrackerMethod, error) {
	var p static_proto.TrackerMethod
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func behaviourCategoryAimToRecord(behaviourCategoryAim *static_proto.BehaviourCategoryAim) (string, error) {
	data, err := common.MarhalToObject(behaviourCategoryAim)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    behaviourCategoryAim.Id,
		"id":      behaviourCategoryAim.Id,
		"created": behaviourCategoryAim.Created,
		"updated": behaviourCategoryAim.Updated,
		"name":    behaviourCategoryAim.Name,
		// "parameter1": behaviourCategoryAim.OrgId,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToBehaviourCategoryAim(r *db_proto.Record) (*static_proto.BehaviourCategoryAim, error) {
	var p static_proto.BehaviourCategoryAim
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func contentParentCategoryToRecord(contentParentCategory *static_proto.ContentParentCategory) (string, error) {
	data, err := common.MarhalToObject(contentParentCategory)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":       contentParentCategory.Id,
		"id":         contentParentCategory.Id,
		"created":    contentParentCategory.Created,
		"updated":    contentParentCategory.Updated,
		"name":       contentParentCategory.Name,
		"parameter1": contentParentCategory.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToContentParentCategory(r *db_proto.Record) (*static_proto.ContentParentCategory, error) {
	var p static_proto.ContentParentCategory
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func contentCategoryToRecord(contentCategory *static_proto.ContentCategory) (string, error) {
	data, err := common.MarhalToObject(contentCategory)
	if err != nil {
		return "", err
	}

	// matching contetnt category
	if len(contentCategory.Parent) > 0 {
		var arr []interface{}
		for _, item := range contentCategory.Parent {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["parent"] = arr
	} else {
		delete(data, "parent")
	}
	// matching contetnt category's actions
	if len(contentCategory.Actions) > 0 {
		var arr []interface{}
		for _, item := range contentCategory.Actions {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["actions"] = arr
	} else {
		delete(data, "actions")
	}
	// matching contetnt category's tracker methods
	if len(contentCategory.TrackerMethods) > 0 {
		var arr []interface{}
		for _, item := range contentCategory.TrackerMethods {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["trackerMethods"] = arr
	} else {
		delete(data, "trackerMethods")
	}

	d := map[string]interface{}{
		"_key":       contentCategory.Id,
		"id":         contentCategory.Id,
		"created":    contentCategory.Created,
		"updated":    contentCategory.Updated,
		"name":       contentCategory.Name,
		"parameter1": contentCategory.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToContentCategory(r *db_proto.Record) (*static_proto.ContentCategory, error) {
	var p static_proto.ContentCategory
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func contentTypeToRecord(contentType *static_proto.ContentType) (string, error) {
	data, err := common.MarhalToObject(contentType)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    contentType.Id,
		"id":      contentType.Id,
		"created": contentType.Created,
		"updated": contentType.Updated,
		"name":    contentType.Name,
		// "parameter1": contentType.OrgId,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToContentType(r *db_proto.Record) (*static_proto.ContentType, error) {
	var p static_proto.ContentType
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func contentSourceTypeToRecord(contentSourceType *static_proto.ContentSourceType) (string, error) {
	data, err := common.MarhalToObject(contentSourceType)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    contentSourceType.Id,
		"id":      contentSourceType.Id,
		"created": contentSourceType.Created,
		"updated": contentSourceType.Updated,
		"name":    contentSourceType.Name,
		// "parameter1": contentSourceType.OrgId,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToContentSourceType(r *db_proto.Record) (*static_proto.ContentSourceType, error) {
	var p static_proto.ContentSourceType
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func moduleTriggerToRecord(moduleTrigger *static_proto.ModuleTrigger) (string, error) {
	data, err := common.MarhalToObject(moduleTrigger)
	if err != nil {
		return "", err
	}

	common.FilterObject(data, "module", moduleTrigger.Module)
	if len(moduleTrigger.Events) > 0 {
		var arr []interface{}
		for _, item := range moduleTrigger.Events {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["events"] = arr
	} else {
		delete(data, "events")
	}
	if len(moduleTrigger.ContentTypes) > 0 {
		var arr []interface{}
		for _, item := range moduleTrigger.ContentTypes {
			arr = append(arr, map[string]string{
				"id": item.Id,
			})
		}
		data["contentTypes"] = arr
	} else {
		delete(data, "contentTypes")
	}

	d := map[string]interface{}{
		"_key":    moduleTrigger.Id,
		"id":      moduleTrigger.Id,
		"created": moduleTrigger.Created,
		"updated": moduleTrigger.Updated,
		"name":    moduleTrigger.Name,
		// "parameter1": moduleTrigger.OrgId,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToModuleTrigger(r *db_proto.Record) (*static_proto.ModuleTrigger, error) {
	var p static_proto.ModuleTrigger
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func triggerContentTypeToRecord(triggerContentType *static_proto.TriggerContentType) (string, error) {
	data, err := common.MarhalToObject(triggerContentType)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":    triggerContentType.Id,
		"id":      triggerContentType.Id,
		"created": triggerContentType.Created,
		"updated": triggerContentType.Updated,
		"name":    triggerContentType.Name,
		// "parameter1": triggerContentType.OrgId,
		"data": data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToTriggerContentType(r *db_proto.Record) (*static_proto.TriggerContentType, error) {
	var p static_proto.TriggerContentType
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func roleToRecord(role *static_proto.Role) (string, error) {
	data, err := common.MarhalToObject(role)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":       role.Id,
		"id":         role.Id,
		"created":    role.Created,
		"updated":    role.Updated,
		"name":       role.Name,
		"parameter1": role.NameSlug,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToRole(r *db_proto.Record) (*static_proto.Role, error) {
	var p static_proto.Role
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func setbackToRecord(setback *static_proto.Setback) (string, error) {
	data, err := common.MarhalToObject(setback)
	if err != nil {
		return "", err
	}
	d := map[string]interface{}{
		"_key":       setback.Id,
		"id":         setback.Id,
		"created":    setback.Created,
		"updated":    setback.Updated,
		"name":       setback.Name,
		"parameter1": setback.OrgId,
		"data":       data,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToSetback(r *db_proto.Record) (*static_proto.Setback, error) {
	var p static_proto.Setback
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func queryOrgAuth(orgId string) string {
	var query string
	if len(orgId) != 0 {
		query = fmt.Sprintf(` && doc.parameter1 == "%s"`, orgId)
	}
	return query
}

func queryTeamAuth(teamId string) string {
	var query string
	if len(teamId) != 0 {
		query = fmt.Sprintf(` && doc.parameter2 == "%s"`, teamId)
	}
	return query
}

func queryPaginate(offset, limit int64) (string, string) {
	if limit == 0 {
		limit = 10
	}
	offs := fmt.Sprintf("%d", offset)
	size := fmt.Sprintf("%d", limit)
	return offs, size
}

// AllApps get all apps
func AllApps(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.App, error) {
	var apps []*static_proto.App

	sort_query := common.QuerySort(sortParameter, sortDirection)
	limit_query := common.QueryPaginate(offset, limit)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		LET platforms = (
			FOR p IN doc.data.platforms
			FOR platform IN %v 
			FILTER platform._key == p.id
			RETURN platform.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			platforms:platforms
		}})`, common.DbAppTable, sort_query, limit_query, common.DbPlatformTable)

	resp, err := runQuery(ctx, q, common.DbAppTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if app, err := recordToApp(r); err == nil {
			apps = append(apps, app)
		}
	}
	return apps, nil
}

// CreateApp creates a app
func CreateApp(ctx context.Context, app *static_proto.App) error {
	if app.Created == 0 {
		app.Created = time.Now().Unix()
	}
	if app.Updated == 0 {
		app.Updated = time.Now().Unix()
	}
	record, err := appToRecord(app)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, app.Id, record, record, common.DbAppTable)
	_, err = runQuery(ctx, q, common.DbAppTable)
	return err
}

// ReadApp reads a app by ID
func ReadApp(ctx context.Context, id, orgId, teamId string) (*static_proto.App, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET platforms = (
			FOR p IN doc.data.platforms
			FOR platform IN %v 
			FILTER platform._key == p.id
			RETURN platform.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			platforms:platforms
		}})`, common.DbAppTable, query, common.DbPlatformTable)

	resp, err := runQuery(ctx, q, common.DbAppTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToApp(resp.Records[0])
	return data, err
}

// DeleteApp deletes a app by ID
func DeleteApp(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbAppTable, query, common.DbAppTable)
	_, err := runQuery(ctx, q, common.DbAppTable)
	return err
}

// AllPlatforms get all platforms
func AllPlatforms(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.Platform, error) {
	var platforms []*static_proto.Platform

	var limit_query string
	if offset != -1 && limit != -1 {
		offs, size := queryPaginate(offset, limit)
		limit_query = fmt.Sprintf("LIMIT %s, %s", offs, size)
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbPlatformTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbPlatformTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if platform, err := recordToPlatform(r); err == nil {
			platforms = append(platforms, platform)
		}
	}
	return platforms, nil
}

// CreatePlatform creates a platform
func CreatePlatform(ctx context.Context, platform *static_proto.Platform) error {
	if platform.Created == 0 {
		platform.Created = time.Now().Unix()
	}
	if platform.Updated == 0 {
		platform.Updated = time.Now().Unix()
	}
	record, err := platformToRecord(platform)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, platform.Id, record, record, common.DbPlatformTable)
	_, err = runQuery(ctx, q, common.DbPlatformTable)
	return err
}

// ReadPlatform reads a platform by ID
func ReadPlatform(ctx context.Context, id, orgId, teamId string) (*static_proto.Platform, error) {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbPlatformTable, query)

	resp, err := runQuery(ctx, q, common.DbPlatformTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToPlatform(resp.Records[0])
	return data, err
}

// DeletePlatform deletes a platform by ID
func DeletePlatform(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbPlatformTable, query, common.DbPlatformTable)
	_, err := runQuery(ctx, q, common.DbPlatformTable)
	return err
}

// AllWearables get all wearables
func AllWearables(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.Wearable, error) {
	var wearables []*static_proto.Wearable

	var limit_query string
	if offset != -1 && limit != -1 {
		offs, size := queryPaginate(offset, limit)
		limit_query = fmt.Sprintf("LIMIT %s, %s", offs, size)
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbWearableTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbWearableTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if wearable, err := recordToWearable(r); err == nil {
			wearables = append(wearables, wearable)
		}
	}
	return wearables, nil
}

// CreateWearable creates a wearable
func CreateWearable(ctx context.Context, wearable *static_proto.Wearable) error {
	if wearable.Created == 0 {
		wearable.Created = time.Now().Unix()
	}
	if wearable.Updated == 0 {
		wearable.Updated = time.Now().Unix()
	}
	record, err := wearableToRecord(wearable)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, wearable.Id, record, record, common.DbWearableTable)
	_, err = runQuery(ctx, q, common.DbWearableTable)
	return err
}

// ReadWearable reads a wearable by ID
func ReadWearable(ctx context.Context, id, orgId, teamId string) (*static_proto.Wearable, error) {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbWearableTable, query)

	resp, err := runQuery(ctx, q, common.DbWearableTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToWearable(resp.Records[0])
	return data, err
}

// DeleteWearable deletes a wearable by ID
func DeleteWearable(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbWearableTable, query, common.DbWearableTable)
	_, err := runQuery(ctx, q, common.DbWearableTable)
	return err
}

// AllDevices get all devices
func AllDevices(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.Device, error) {
	var devices []*static_proto.Device

	var limit_query string
	if offset != -1 && limit != -1 {
		offs, size := queryPaginate(offset, limit)
		limit_query = fmt.Sprintf("LIMIT %s, %s", offs, size)
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbDeviceTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbDeviceTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if device, err := recordToDevice(r); err == nil {
			devices = append(devices, device)
		}
	}
	return devices, nil
}

// CreateDevice creates a device
func CreateDevice(ctx context.Context, device *static_proto.Device) error {
	if device.Created == 0 {
		device.Created = time.Now().Unix()
	}
	if device.Updated == 0 {
		device.Updated = time.Now().Unix()
	}
	record, err := deviceToRecord(device)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, device.Id, record, record, common.DbDeviceTable)
	_, err = runQuery(ctx, q, common.DbDeviceTable)
	return err
}

// ReadDevice reads a device by ID
func ReadDevice(ctx context.Context, id, orgId, teamId string) (*static_proto.Device, error) {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbDeviceTable, query)

	resp, err := runQuery(ctx, q, common.DbDeviceTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToDevice(resp.Records[0])
	return data, err
}

// DeleteDevice deletes a device by ID
func DeleteDevice(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbDeviceTable, query, common.DbDeviceTable)
	_, err := runQuery(ctx, q, common.DbDeviceTable)
	return err
}

// AllMarkers get all markers
func AllMarkers(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.Marker, error) {
	var markers []*static_proto.Marker
	//FIXME:removing orgId from filtering for now - we need a better way to combine the general content catalogue with organisation content
	query := common.QueryAuth(`FILTER`, "", "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET apps = (
			FILTER NOT_NULL(doc.data.apps)
			FOR p IN doc.data.apps
			FOR app IN %v 
			FILTER p.id == app._key RETURN app.data
		)
		LET wearables = (
			FILTER NOT_NULL(doc.data.wearables)
			FOR p IN doc.data.wearables
			FOR wearable IN %v
			FILTER p.id == wearable._key RETURN wearable.data
		)
		LET devices = (
			FILTER NOT_NULL(doc.data.devices)
			FOR p IN doc.data.devices
			FOR device IN %v
			FILTER p.id == device._key RETURN device.data
		)
		LET trackerMethods = (
			FOR t IN OUTBOUND doc %v
			OPTIONS {
				bfs: true,
				uniqueVertices: "global"
			}
			RETURN t.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			apps:apps,
			wearables:wearables,
			devices:devices,
			trackerMethods:trackerMethods
		}})`, common.DbMarkerTable, query, sort_query, limit_query,
		common.DbAppTable, common.DbWearableTable, common.DbDeviceTable, common.DbMarkerTrackerEdgeTable)

	resp, err := runQuery(ctx, q, common.DbMarkerTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if marker, err := recordToMarker(r); err == nil {
			markers = append(markers, marker)
		}
	}
	return markers, nil
}

// CreateMarker creates a marker
func CreateMarker(ctx context.Context, marker *static_proto.Marker) error {
	if marker.Created == 0 {
		marker.Created = time.Now().Unix()
	}
	marker.Updated = time.Now().Unix()

	record, err := markerToRecord(marker)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, marker.Id, record, record, common.DbMarkerTable)
	_, err = runQuery(ctx, q, common.DbMarkerTable)
	if err != nil {
		common.ErrorLog(common.StaticSrv, CreateMarker, err, "Create runquery is failed")
		return err
	}

	_from := fmt.Sprintf("%v/%v", common.DbMarkerTable, marker.Id)
	// remove original edges
	q = fmt.Sprintf(`
		FOR doc IN %v 
		FILTER doc._from == "%v"
		REMOVE doc IN %v`, common.DbMarkerTrackerEdgeTable, _from, common.DbMarkerTrackerEdgeTable)
	if _, err := runQuery(ctx, q, common.DbMarkerTrackerEdgeTable); err != nil {
		common.ErrorLog(common.StaticSrv, CreateMarker, err, "Remove MarkerTrackerEdge failed")
	}
	for _, tracker := range marker.TrackerMethods {
		field := fmt.Sprintf(`{_from:"%v", _to:"%v/%v"}`, _from, common.DbTrackerMethodTable, tracker.Id)
		q = fmt.Sprintf(`INSERT %v INTO %v`, field, common.DbMarkerTrackerEdgeTable)
		_, err = runQuery(ctx, q, common.DbMarkerTrackerEdgeTable)
		if err != nil {
			return err
		}
	}
	return err
}

// ReadMarker reads a marker by ID
func ReadMarker(ctx context.Context, id, orgId, teamId string) (*static_proto.Marker, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET apps = (
			FILTER NOT_NULL(doc.data.apps)
			FOR p IN doc.data.apps
			FOR app IN %v 
			FILTER p.id == app._key RETURN app.data
		)
		LET wearables = (
			FILTER NOT_NULL(doc.data.wearables)
			FOR p IN doc.data.wearables
			FOR wearable IN %v
			FILTER p.id == wearable._key RETURN wearable.data
		)
		LET devices = (
			FILTER NOT_NULL(doc.data.devices)
			FOR p IN doc.data.devices
			FOR device IN %v
			FILTER p.id == device._key RETURN device.data
		)
		LET trackerMethods = (
			FOR t IN OUTBOUND doc %v
			OPTIONS {
				bfs: true,
				uniqueVertices: "global"
			}
			RETURN t.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			apps:apps,
			wearables:wearables,
			devices:devices,
			trackerMethods:trackerMethods
		}})`, common.DbMarkerTable, query,
		common.DbAppTable, common.DbWearableTable, common.DbDeviceTable, common.DbMarkerTrackerEdgeTable)

	resp, err := runQuery(ctx, q, common.DbMarkerTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToMarker(resp.Records[0])
	return data, err
}

// DeleteMarker deletes a marker by ID
func DeleteMarker(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query += queryOrgAuth(orgId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbMarkerTable, query, common.DbMarkerTable)
	_, err := runQuery(ctx, q, common.DbMarkerTable)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`FILTER doc._from == "%v/%v"`, common.DbMarkerTable, id)
	q = fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbMarkerTrackerEdgeTable, query, common.DbMarkerTrackerEdgeTable)
	_, err = runQuery(ctx, q, common.DbMarkerTrackerEdgeTable)
	return err
}

func FilterMarker(ctx context.Context, trackerMethods []string, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.Marker, error) {
	var markers []*static_proto.Marker
	query := `FILTER`
	// search creator
	if len(trackerMethods) > 0 {
		methods := common.QueryStringFromArray(trackerMethods)
		query += fmt.Sprintf(" && trackerMethods[*].id ANY IN [%v]", methods)
	}
	query = common.QueryAuth(query, orgId, teamId)
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		LET trackerMethods = (
			FOR t IN OUTBOUND doc %v
			OPTIONS {
				bfs: true,
				uniqueVertices: "global"
			}
			RETURN t.data
		)
		%s
		%s
		%s
		LET apps = (
			FILTER NOT_NULL(doc.data.apps)
			FOR p IN doc.data.apps
			FOR app IN %v 
			FILTER p.id == app._key RETURN app.data
		)
		LET wearables = (
			FILTER NOT_NULL(doc.data.wearables)
			FOR p IN doc.data.wearables
			FOR wearable IN %v
			FILTER p.id == wearable._key RETURN wearable.data
		)
		LET devices = (
			FILTER NOT_NULL(doc.data.devices)
			FOR p IN doc.data.devices
			FOR device IN %v
			FILTER p.id == device._key RETURN device.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			apps:apps,
			wearables:wearables,
			devices:devices,
			trackerMethods:trackerMethods
		}})`, common.DbMarkerTable, common.DbMarkerTrackerEdgeTable, query, sort_query, limit_query,
		common.DbAppTable, common.DbWearableTable, common.DbDeviceTable)

	resp, err := runQuery(ctx, q, common.DbMarkerTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if marker, err := recordToMarker(r); err == nil {
			markers = append(markers, marker)
		}
	}
	return markers, nil
}

// AllModules get all modules
func AllModules(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.Module, error) {
	var modules []*static_proto.Module

	var limit_query string
	if offset != -1 && limit != -1 {
		offs, size := queryPaginate(offset, limit)
		limit_query = fmt.Sprintf("LIMIT %s, %s", offs, size)
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbModuleTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbModuleTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if module, err := recordToModule(r); err == nil {
			modules = append(modules, module)
		}
	}
	return modules, nil
}

// CreateModule creates a module
func CreateModule(ctx context.Context, module *static_proto.Module) error {
	if module.Created == 0 {
		module.Created = time.Now().Unix()
	}
	if module.Updated == 0 {
		module.Updated = time.Now().Unix()
	}
	record, err := moduleToRecord(module)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, module.Id, record, record, common.DbModuleTable)
	_, err = runQuery(ctx, q, common.DbModuleTable)
	return err
}

// ReadModule reads a module by ID
func ReadModule(ctx context.Context, id, orgId, teamId string) (*static_proto.Module, error) {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbModuleTable, query)

	resp, err := runQuery(ctx, q, common.DbModuleTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToModule(resp.Records[0])
	return data, err
}

// DeleteModule deletes a module by ID
func DeleteModule(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)
	query += queryOrgAuth(orgId)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbModuleTable, query, common.DbModuleTable)
	_, err := runQuery(ctx, q, common.DbModuleTable)
	return err
}

// AllBehaviourCategories get all categories
func AllBehaviourCategories(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.BehaviourCategory, error) {
	var categories []*static_proto.BehaviourCategory
	//FIXME:Disabling filter for orgId and Employee temporariy
	//query := common.QueryAuth(`FILTER`, orgId, teamId)
	query := common.QueryAuth(`FILTER`, "", "")
	sort_query := common.QuerySort(sortParameter, sortDirection)
	limit_query := common.QueryPaginate(offset, limit)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET aims = (
			FILTER NOT_NULL(doc.data.aims)
			FOR p IN doc.data.aims
			FOR aim IN %v
			FILTER p.id == aim._key RETURN aim.data
		)
		LET default = (FOR d IN %v FILTER doc.data.markerDefault.id == d._key RETURN d.data)
		LET ops = (
			FILTER NOT_NULL(doc.data.markerOptions)
			FOR p IN doc.data.markerOptions
			FOR option IN %v
			FILTER p.id == option._key RETURN option.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			aims:aims,
			markerDefault:default[0],
			markerOptions:ops
		}})`, common.DbBehaviourCategoryTable, query, sort_query, limit_query,
		common.DbBehaviourCategoryAimTable, common.DbMarkerTable, common.DbMarkerTable,
	)

	resp, err := runQuery(ctx, q, common.DbBehaviourCategoryTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if category, err := recordToCategory(r); err == nil {
			categories = append(categories, category)
		}
	}
	return categories, nil
}

// CreateBehaviourCategory creates a category
func CreateBehaviourCategory(ctx context.Context, category *static_proto.BehaviourCategory) error {
	if category.Created == 0 {
		category.Created = time.Now().Unix()
	}
	category.Updated = time.Now().Unix()

	record, err := categoryToRecord(category)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, category.Id, record, record, common.DbBehaviourCategoryTable)
	_, err = runQuery(ctx, q, common.DbBehaviourCategoryTable)
	return err
}

// ReadBehaviourCategory reads a category by ID
func ReadBehaviourCategory(ctx context.Context, id, orgId, teamId string) (*static_proto.BehaviourCategory, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query = common.QueryAuth(query, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET aims = (
			FILTER NOT_NULL(doc.data.aims)
			FOR p IN doc.data.aims
			FOR aim IN %v
			FILTER p.id == aim._key RETURN aim.data
		)
		LET default = (FOR d IN %v FILTER doc.data.markerDefault.id == d._key RETURN d.data)
		LET ops = (
			FILTER NOT_NULL(doc.data.markerOptions)
			FOR p IN doc.data.markerOptions
			FOR option IN %v
			FILTER p.id == option._key RETURN option.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			aims:aims,
			markerDefault:default[0],
			markerOptions:ops
		}})`, common.DbBehaviourCategoryTable, query,
		common.DbBehaviourCategoryAimTable, common.DbMarkerTable, common.DbMarkerTable,
	)

	resp, err := runQuery(ctx, q, common.DbBehaviourCategoryTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToCategory(resp.Records[0])
	return data, err
}

// DeleteBehaviourCategory deletes a category by ID
func DeleteBehaviourCategory(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	query += queryOrgAuth(orgId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbBehaviourCategoryTable, query, common.DbBehaviourCategoryTable)
	_, err := runQuery(ctx, q, common.DbBehaviourCategoryTable)
	return err
}

func FilterBehaviourCategory(ctx context.Context, markers []string, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.BehaviourCategory, error) {
	var categories []*static_proto.BehaviourCategory

	query := `FILTER`
	// search creator
	if len(markers) > 0 {
		options := common.QueryStringFromArray(markers)
		query += fmt.Sprintf(" && ops[*].id ANY IN [%v]", options)
	}
	query = common.QueryAuth(query, orgId, "")
	sort_query := common.QuerySort(sortParameter, sortDirection)
	limit_query := common.QueryPaginate(offset, limit)

	q := fmt.Sprintf(`
		FOR doc IN %v
		LET ops = (
			FILTER NOT_NULL(doc.data.markerOptions)
			FOR p IN doc.data.markerOptions
			FOR option IN %v
			FILTER p.id == option._key RETURN option.data
		)
		%s
		%s
		%s
		LET aims = (
			FILTER NOT_NULL(doc.data.aims)
			FOR p IN doc.data.aims
			FOR aim IN %v
			FILTER p.id == aim._key RETURN aim.data
		)
		LET default = (FOR d IN %v FILTER doc.data.markerDefault.id == d._key RETURN d.data)
		RETURN MERGE_RECURSIVE(doc,{data:{
			aims:aims,
			markerDefault:default[0],
			markerOptions:ops
		}})`, common.DbBehaviourCategoryTable, common.DbMarkerTable,
		query, sort_query, limit_query,
		common.DbBehaviourCategoryAimTable, common.DbMarkerTable,
	)

	resp, err := runQuery(ctx, q, common.DbBehaviourCategoryTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if category, err := recordToCategory(r); err == nil {
			categories = append(categories, category)
		}
	}
	return categories, nil
}

// AllSocialTypes get all socialTypes
func AllSocialTypes(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.SocialType, error) {
	var socialTypes []*static_proto.SocialType
	var limit_query string
	if offset != -1 && limit != -1 {
		offs, size := queryPaginate(offset, limit)
		limit_query = fmt.Sprintf("LIMIT %s, %s", offs, size)
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbSocialTypeTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbSocialTypeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if socialType, err := recordToSocialType(r); err == nil {
			socialTypes = append(socialTypes, socialType)
		}
	}
	return socialTypes, nil
}

// CreateSocialType creates a socialType
func CreateSocialType(ctx context.Context, socialType *static_proto.SocialType) error {
	if socialType.Created == 0 {
		socialType.Created = time.Now().Unix()
	}
	if socialType.Updated == 0 {
		socialType.Updated = time.Now().Unix()
	}
	record, err := socialTypeToRecord(socialType)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, socialType.Id, record, record, common.DbSocialTypeTable)
	_, err = runQuery(ctx, q, common.DbSocialTypeTable)
	return err
}

// ReadSocialType reads a socialType by ID
func ReadSocialType(ctx context.Context, id, orgId, teamId string) (*static_proto.SocialType, error) {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbSocialTypeTable, query)

	resp, err := runQuery(ctx, q, common.DbSocialTypeTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToSocialType(resp.Records[0])
	return data, err
}

// DeleteSocialType deletes a socialType by ID
func DeleteSocialType(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbSocialTypeTable, query, common.DbSocialTypeTable)
	_, err := runQuery(ctx, q, common.DbSocialTypeTable)
	return err
}

// AllNotifications get all notifications
func AllNotifications(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.Notification, error) {
	var notifications []*static_proto.Notification
	var limit_query string
	if offset != -1 && limit != -1 {
		offs, size := queryPaginate(offset, limit)
		limit_query = fmt.Sprintf("LIMIT %s, %s", offs, size)
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbNotificationTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbNotificationTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if notification, err := recordToNotification(r); err == nil {
			notifications = append(notifications, notification)
		}
	}
	return notifications, nil
}

// CreateNotification creates a notification
func CreateNotification(ctx context.Context, notification *static_proto.Notification) error {
	if notification.Created == 0 {
		notification.Created = time.Now().Unix()
	}
	if notification.Updated == 0 {
		notification.Updated = time.Now().Unix()
	}
	record, err := notificationToRecord(notification)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, notification.Id, record, record, common.DbNotificationTable)
	_, err = runQuery(ctx, q, common.DbNotificationTable)
	return err
}

// ReadNotification reads a notification by ID
func ReadNotification(ctx context.Context, id, orgId, teamId string) (*static_proto.Notification, error) {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbNotificationTable, query)

	resp, err := runQuery(ctx, q, common.DbNotificationTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToNotification(resp.Records[0])
	return data, err
}

// DeleteNotification deletes a notification by ID
func DeleteNotification(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbNotificationTable, query, common.DbNotificationTable)
	_, err := runQuery(ctx, q, common.DbNotificationTable)
	return err
}

// AllTrackerMethods get all trackerMethods
func AllTrackerMethods(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.TrackerMethod, error) {
	var trackerMethods []*static_proto.TrackerMethod
	var limit_query string
	if offset != -1 && limit != -1 {
		offs, size := queryPaginate(offset, limit)
		limit_query = fmt.Sprintf("LIMIT %s, %s", offs, size)
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbTrackerMethodTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbTrackerMethodTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if trackerMethod, err := recordToTrackerMethod(r); err == nil {
			trackerMethods = append(trackerMethods, trackerMethod)
		}
	}
	return trackerMethods, nil
}

// CreateTrackerMethod creates a trackerMethod
func CreateTrackerMethod(ctx context.Context, trackerMethod *static_proto.TrackerMethod) error {
	if trackerMethod.Created == 0 {
		trackerMethod.Created = time.Now().Unix()
	}
	trackerMethod.Updated = time.Now().Unix()

	record, err := trackerMethodToRecord(trackerMethod)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, trackerMethod.Id, record, record, common.DbTrackerMethodTable)
	_, err = runQuery(ctx, q, common.DbTrackerMethodTable)
	return err
}

// ReadTrackerMethod reads a trackerMethod by ID
func ReadTrackerMethod(ctx context.Context, id, name_slug, orgId, teamId string) (*static_proto.TrackerMethod, error) {
	query := `FILTER`
	if len(id) != 0 {
		query += fmt.Sprintf(` || doc._key == "%s"`, id)
	}
	if len(name_slug) != 0 {
		query += fmt.Sprintf(` || doc.data.name_slug == "%s"`, name_slug)
	}
	query = strings.Replace(query, `FILTER && `, `FILTER `, -1)
	query = strings.Replace(query, `FILTER || `, `FILTER `, -1)

	if query == `FILTER` {
		query = ""
	}

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbTrackerMethodTable, query)

	resp, err := runQuery(ctx, q, common.DbTrackerMethodTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToTrackerMethod(resp.Records[0])
	return data, err
}

// DeleteTrackerMethod deletes a trackerMethod by ID
func DeleteTrackerMethod(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbTrackerMethodTable, query, common.DbTrackerMethodTable)
	_, err := runQuery(ctx, q, common.DbTrackerMethodTable)
	return err
}

// FilterTrackerMethod filter trackerMethod by marker
func FilterTrackerMethod(ctx context.Context, markers []string, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.TrackerMethod, error) {
	var trackerMethods []*static_proto.TrackerMethod
	query := `FILTER`
	// search creator
	if len(markers) > 0 {
		arr := common.QueryStringFromArray(markers)
		query += fmt.Sprintf(" && doc._key IN [%v]", arr)
	}
	query = common.QueryAuth(query, "", "")
	limit_query := common.QueryPaginate(offset, limit)
	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		LET methods = (
			FOR doc IN %v
			%s
			FOR t IN OUTBOUND doc %v
			OPTIONS {
				bfs: true,
				uniqueVertices: "global"
			}
			RETURN t
		)[**]
		FOR doc IN methods
		%s
		%s
		RETURN DISTINCT doc`, common.DbMarkerTable, query, common.DbMarkerTrackerEdgeTable,
		sort_query, limit_query,
	)

	resp, err := runQuery(ctx, q, common.DbMarkerTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if method, err := recordToTrackerMethod(r); err == nil {
			trackerMethods = append(trackerMethods, method)
		}
	}
	return trackerMethods, nil
}

// AllBehaviourCategoryAims get all behaviourCategoryAims
func AllBehaviourCategoryAims(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.BehaviourCategoryAim, error) {
	var behaviourCategoryAims []*static_proto.BehaviourCategoryAim

	sort_query := common.QuerySort(sortParameter, sortDirection)
	limit_query := common.QueryPaginate(offset, limit)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbBehaviourCategoryAimTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbBehaviourCategoryAimTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if behaviourCategoryAim, err := recordToBehaviourCategoryAim(r); err == nil {
			behaviourCategoryAims = append(behaviourCategoryAims, behaviourCategoryAim)
		}
	}
	return behaviourCategoryAims, nil
}

// CreateBehaviourCategoryAim creates a behaviourCategoryAim
func CreateBehaviourCategoryAim(ctx context.Context, behaviourCategoryAim *static_proto.BehaviourCategoryAim) error {
	if behaviourCategoryAim.Created == 0 {
		behaviourCategoryAim.Created = time.Now().Unix()
	}
	behaviourCategoryAim.Updated = time.Now().Unix()

	record, err := behaviourCategoryAimToRecord(behaviourCategoryAim)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, behaviourCategoryAim.Id, record, record, common.DbBehaviourCategoryAimTable)
	_, err = runQuery(ctx, q, common.DbBehaviourCategoryAimTable)
	return err
}

// ReadBehaviourCategoryAim reads a behaviourCategoryAim by ID
func ReadBehaviourCategoryAim(ctx context.Context, id, orgId, teamId string) (*static_proto.BehaviourCategoryAim, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbBehaviourCategoryAimTable, query)

	resp, err := runQuery(ctx, q, common.DbBehaviourCategoryAimTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToBehaviourCategoryAim(resp.Records[0])
	return data, err
}

// DeleteBehaviourCategoryAim deletes a behaviourCategoryAim by ID
func DeleteBehaviourCategoryAim(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbBehaviourCategoryAimTable, query, common.DbBehaviourCategoryAimTable)
	_, err := runQuery(ctx, q, common.DbBehaviourCategoryAimTable)
	return err
}

// AllContentParentCategories get all contentParentCategories
func AllContentParentCategories(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.ContentParentCategory, error) {
	var contentParentCategories []*static_proto.ContentParentCategory

	sort_query := common.QuerySort(sortParameter, sortDirection)
	limit_query := common.QueryAuth(`FILTER`, orgId, "")

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbContentParentCategoryTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbContentParentCategoryTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if contentParentCategory, err := recordToContentParentCategory(r); err == nil {
			contentParentCategories = append(contentParentCategories, contentParentCategory)
		}
	}
	return contentParentCategories, nil
}

// CreateContentParentCategory creates a contentParentCategory
func CreateContentParentCategory(ctx context.Context, contentParentCategory *static_proto.ContentParentCategory) error {
	if contentParentCategory.Created == 0 {
		contentParentCategory.Created = time.Now().Unix()
	}
	contentParentCategory.Updated = time.Now().Unix()

	record, err := contentParentCategoryToRecord(contentParentCategory)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, contentParentCategory.Id, record, record, common.DbContentParentCategoryTable)
	_, err = runQuery(ctx, q, common.DbContentParentCategoryTable)
	return err
}

// ReadContentParentCategory reads a contentParentCategory by ID
func ReadContentParentCategory(ctx context.Context, id, orgId, teamId string) (*static_proto.ContentParentCategory, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)
	// query = common.QueryAuth(query, orgId, teamId)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbContentParentCategoryTable, query)

	resp, err := runQuery(ctx, q, common.DbContentParentCategoryTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToContentParentCategory(resp.Records[0])
	return data, err
}

// DeleteContentParentCategory deletes a contentParentCategory by ID
func DeleteContentParentCategory(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbContentParentCategoryTable, query, common.DbContentParentCategoryTable)
	_, err := runQuery(ctx, q, common.DbContentParentCategoryTable)
	return err
}

// AllContentCategories get all contentCategories
func AllContentCategories(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.ContentCategory, error) {
	var contentCategories []*static_proto.ContentCategory

	sort_query := common.QuerySort(sortParameter, sortDirection)
	limit_query := common.QueryPaginate(offset, limit)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		LET parent = (
			FILTER NOT_NULL(doc.data.parent)
			FOR p IN doc.data.parent
			FOR parent IN %v
			FILTER p.id == parent._key RETURN parent.data
		)
		LET actions = (
			FILTER NOT_NULL(doc.data.actions)
			FOR p IN doc.data.actions
			FOR action IN %v
			FILTER p.id == action._key RETURN action.data
		)
		LET methods = (
			FILTER NOT_NULL(doc.data.trackerMethods)
			FOR p IN doc.data.trackerMethods
			FOR method IN %v
			FILTER p.id == method._key RETURN method.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			parent:parent,
			actions:actions,
			trackerMethods:methods
		}})`, common.DbContentCategoryTable, sort_query, limit_query,
		common.DbContentParentCategoryTable, common.DbActionTable, common.DbTrackerMethodTable,
	)

	resp, err := runQuery(ctx, q, common.DbContentCategoryTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if contentCategory, err := recordToContentCategory(r); err == nil {
			contentCategories = append(contentCategories, contentCategory)
		}
	}
	return contentCategories, nil
}

// CreateContentCategory creates a contentCategory
func CreateContentCategory(ctx context.Context, contentCategory *static_proto.ContentCategory) error {
	if contentCategory.Created == 0 {
		contentCategory.Created = time.Now().Unix()
	}
	contentCategory.Updated = time.Now().Unix()

	record, err := contentCategoryToRecord(contentCategory)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, contentCategory.Id, record, record, common.DbContentCategoryTable)
	_, err = runQuery(ctx, q, common.DbContentCategoryTable)
	return err
}

// ReadContentCategory reads a contentCategory by ID
func ReadContentCategory(ctx context.Context, id, orgId, teamId string) (*static_proto.ContentCategory, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET parent = (
			FILTER NOT_NULL(doc.data.parent)
			FOR p IN doc.data.parent
			FOR parent IN %v
			FILTER p.id == parent._key RETURN parent.data
		)
		LET actions = (
			FILTER NOT_NULL(doc.data.actions)
			FOR p IN doc.data.actions
			FOR action IN %v
			FILTER p.id == action._key RETURN action.data
		)
		LET methods = (
			FILTER NOT_NULL(doc.data.trackerMethods)
			FOR p IN doc.data.trackerMethods
			FOR method IN %v
			FILTER p.id == method._key RETURN method.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			parent:parent,
			actions:actions,
			trackerMethods:methods
		}})`, common.DbContentCategoryTable, query,
		common.DbContentParentCategoryTable, common.DbActionTable, common.DbTrackerMethodTable,
	)

	resp, err := runQuery(ctx, q, common.DbContentCategoryTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToContentCategory(resp.Records[0])
	return data, err
}

// DeleteContentCategory deletes a contentCategory by ID
func DeleteContentCategory(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbContentCategoryTable, query, common.DbContentCategoryTable)
	_, err := runQuery(ctx, q, common.DbContentCategoryTable)
	return err
}

// AllContentTypes get all contentTypes
func AllContentTypes(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.ContentType, error) {
	var contentTypes []*static_proto.ContentType

	var limit_query string
	if offset != -1 && limit != -1 {
		offs, size := queryPaginate(offset, limit)
		limit_query = fmt.Sprintf("LIMIT %s, %s", offs, size)
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbContentTypeTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbContentTypeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if contentType, err := recordToContentType(r); err == nil {
			contentTypes = append(contentTypes, contentType)
		}
	}
	return contentTypes, nil
}

// CreateContentType creates a contentType
func CreateContentType(ctx context.Context, contentType *static_proto.ContentType) error {
	if contentType.Created == 0 {
		contentType.Created = time.Now().Unix()
	}
	if contentType.Updated == 0 {
		contentType.Updated = time.Now().Unix()
	}
	record, err := contentTypeToRecord(contentType)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, contentType.Id, record, record, common.DbContentTypeTable)
	_, err = runQuery(ctx, q, common.DbContentTypeTable)
	return err
}

// ReadContentType reads a contentType by ID
func ReadContentType(ctx context.Context, id, orgId, teamId string) (*static_proto.ContentType, error) {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbContentTypeTable, query)

	resp, err := runQuery(ctx, q, common.DbContentTypeTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToContentType(resp.Records[0])
	return data, err
}

// DeleteContentType deletes a contentType by ID
func DeleteContentType(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbContentTypeTable, query, common.DbContentTypeTable)
	_, err := runQuery(ctx, q, common.DbContentTypeTable)
	return err
}

// AllContentSourceTypes get all contentSourceTypes
func AllContentSourceTypes(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.ContentSourceType, error) {
	var contentSourceTypes []*static_proto.ContentSourceType

	var limit_query string
	if offset != -1 && limit != -1 {
		offs, size := queryPaginate(offset, limit)
		limit_query = fmt.Sprintf("LIMIT %s, %s", offs, size)
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbContentSourceTypeTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbContentSourceTypeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if contentSourceType, err := recordToContentSourceType(r); err == nil {
			contentSourceTypes = append(contentSourceTypes, contentSourceType)
		}
	}
	return contentSourceTypes, nil
}

// CreateContentSourceType creates a contentSourceType
func CreateContentSourceType(ctx context.Context, contentSourceType *static_proto.ContentSourceType) error {
	if contentSourceType.Created == 0 {
		contentSourceType.Created = time.Now().Unix()
	}
	if contentSourceType.Updated == 0 {
		contentSourceType.Updated = time.Now().Unix()
	}
	record, err := contentSourceTypeToRecord(contentSourceType)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, contentSourceType.Id, record, record, common.DbContentSourceTypeTable)
	_, err = runQuery(ctx, q, common.DbContentSourceTypeTable)
	return err
}

// ReadContentSourceType reads a contentSourceType by ID
func ReadContentSourceType(ctx context.Context, id, orgId, teamId string) (*static_proto.ContentSourceType, error) {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbContentSourceTypeTable, query)

	resp, err := runQuery(ctx, q, common.DbContentSourceTypeTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToContentSourceType(resp.Records[0])
	return data, err
}

// DeleteContentSourceType deletes a contentSourceType by ID
func DeleteContentSourceType(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbContentSourceTypeTable, query, common.DbContentSourceTypeTable)
	_, err := runQuery(ctx, q, common.DbContentSourceTypeTable)
	return err
}

// AllModuleTriggers get all moduleTriggers
func AllModuleTriggers(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.ModuleTrigger, error) {
	var moduleTriggers []*static_proto.ModuleTrigger

	sort_query := common.QuerySort(sortParameter, sortDirection)
	limit_query := common.QueryPaginate(offset, limit)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		LET module = (FOR p IN %v FILTER p._key == doc.data.module.id RETURN p.data)
		LET events = (
			FILTER NOT_NULL(doc.data.events)
			FOR p IN doc.data.events
			FOR event IN %v
			FILTER p.id == event._key RETURN event.data
		)
		LET contentTypes = (
			FILTER NOT_NULL(doc.data.contentTypes)
			FOR p IN doc.data.contentTypes
			FOR contentType IN %v
			FILTER p.id == contentType._key RETURN contentType.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			module:module[0],
			events:events,
			contentTypes:contentTypes
		}})`, common.DbModuleTriggerTable, sort_query, limit_query,
		common.DbModuleTable, common.DbTriggerEventTable, common.DbTriggerContentTypeTable,
	)

	resp, err := runQuery(ctx, q, common.DbModuleTriggerTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if moduleTrigger, err := recordToModuleTrigger(r); err == nil {
			moduleTriggers = append(moduleTriggers, moduleTrigger)
		}
	}
	return moduleTriggers, nil
}

// CreateModuleTrigger creates a moduleTrigger
func CreateModuleTrigger(ctx context.Context, moduleTrigger *static_proto.ModuleTrigger) error {
	if moduleTrigger.Created == 0 {
		moduleTrigger.Created = time.Now().Unix()
	}
	if moduleTrigger.Updated == 0 {
		moduleTrigger.Updated = time.Now().Unix()
	}
	record, err := moduleTriggerToRecord(moduleTrigger)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, moduleTrigger.Id, record, record, common.DbModuleTriggerTable)
	_, err = runQuery(ctx, q, common.DbModuleTriggerTable)
	return err
}

// ReadModuleTrigger reads a moduleTrigger by ID
func ReadModuleTrigger(ctx context.Context, id, orgId, teamId string) (*static_proto.ModuleTrigger, error) {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET module = (FOR p IN %v FILTER p._key == doc.data.module.id RETURN p.data)
		LET events = (
			FILTER NOT_NULL(doc.data.events)
			FOR p IN doc.data.events
			FOR event IN %v
			FILTER p.id == event._key RETURN event.data
		)
		LET contentTypes = (
			FILTER NOT_NULL(doc.data.contentTypes)
			FOR p IN doc.data.contentTypes
			FOR contentType IN %v
			FILTER p.id == contentType._key RETURN contentType.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			module:module[0],
			events:events,
			contentTypes:contentTypes
		}})`, common.DbModuleTriggerTable, query,
		common.DbModuleTable, common.DbTriggerEventTable, common.DbTriggerContentTypeTable,
	)

	resp, err := runQuery(ctx, q, common.DbModuleTriggerTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToModuleTrigger(resp.Records[0])
	return data, err
}

// DeleteModuleTrigger deletes a moduleTrigger by ID
func DeleteModuleTrigger(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc._key == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		REMOVE doc IN %v`, common.DbModuleTriggerTable, query, common.DbModuleTriggerTable)
	_, err := runQuery(ctx, q, common.DbModuleTriggerTable)
	return err
}

// FilterModuleTrigger get all moduleTriggers
func FilterModuleTrigger(ctx context.Context, module []string, triggerType []int64, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.ModuleTrigger, error) {
	var moduleTriggers []*static_proto.ModuleTrigger

	query := `FILTER`
	// search creator
	if len(module) > 0 {
		arr := common.QueryStringFromArray(module)
		query += fmt.Sprintf(" || doc.data.module.id IN [%v]", arr)
	}
	if len(triggerType) > 0 {
		g := []string{}
		for _, t := range triggerType {
			g = append(g, strconv.FormatInt(t, 10))
		}
		query += fmt.Sprintf(" || doc.data.type IN [%v]", strings.Join(g[:], ","))
	}
	query = strings.Replace(query, `FILTER || `, `FILTER `, -1)
	if query == `FILTER` {
		query = ""
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)
	limit_query := common.QueryPaginate(offset, limit)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		%s
		LET module = (FOR p IN %v FILTER p._key == doc.data.module.id RETURN p.data)
		LET events = (
			FILTER NOT_NULL(doc.data.events)
			FOR p IN doc.data.events
			FOR event IN %v
			FILTER p.id == event._key RETURN event.data
		)
		LET contentTypes = (
			FILTER NOT_NULL(doc.data.contentTypes)
			FOR p IN doc.data.contentTypes
			FOR contentType IN %v
			FILTER p.id == contentType._key RETURN contentType.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			module:module[0],
			events:events,
			contentTypes:contentTypes
		}})`, common.DbModuleTriggerTable, query, sort_query, limit_query,
		common.DbModuleTable, common.DbTriggerEventTable, common.DbTriggerContentTypeTable,
	)

	resp, err := runQuery(ctx, q, common.DbModuleTriggerTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if moduleTrigger, err := recordToModuleTrigger(r); err == nil {
			moduleTriggers = append(moduleTriggers, moduleTrigger)
		}
	}
	return moduleTriggers, nil
}

// AllTriggerContentTypes get all triggerContentTypes
func AllTriggerContentTypes(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.TriggerContentType, error) {
	var triggerContentTypes []*static_proto.TriggerContentType

	var limit_query string
	if offset != -1 && limit != -1 {
		offs, size := queryPaginate(offset, limit)
		limit_query = fmt.Sprintf("LIMIT %s, %s", offs, size)
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbTriggerContentTypeTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbTriggerContentTypeTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if triggerContentType, err := recordToTriggerContentType(r); err == nil {
			triggerContentTypes = append(triggerContentTypes, triggerContentType)
		}
	}
	return triggerContentTypes, nil
}

// CreateTriggerContentType creates a triggerContentType
func CreateTriggerContentType(ctx context.Context, triggerContentType *static_proto.TriggerContentType) error {
	if triggerContentType.Created == 0 {
		triggerContentType.Created = time.Now().Unix()
	}
	if triggerContentType.Updated == 0 {
		triggerContentType.Updated = time.Now().Unix()
	}
	record, err := triggerContentTypeToRecord(triggerContentType)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, triggerContentType.Id, record, record, common.DbTriggerContentTypeTable)
	_, err = runQuery(ctx, q, common.DbTriggerContentTypeTable)
	return err
}

// ReadTriggerContentType reads a triggerContentType by ID
func ReadTriggerContentType(ctx context.Context, id, orgId, teamId string) (*static_proto.TriggerContentType, error) {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbTriggerContentTypeTable, query)

	resp, err := runQuery(ctx, q, common.DbTriggerContentTypeTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToTriggerContentType(resp.Records[0])
	return data, err
}

// DeleteTriggerContentType deletes a triggerContentType by ID
func DeleteTriggerContentType(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbTriggerContentTypeTable, query, common.DbTriggerContentTypeTable)
	_, err := runQuery(ctx, q, common.DbTriggerContentTypeTable)
	return err
}

// FilterTriggerContentType get all triggerContentTypes
func FilterTriggerContentType(ctx context.Context, moduleTrigger []string, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.TriggerContentType, error) {
	var triggerContentTypes []*static_proto.TriggerContentType

	query := `FILTER`
	// search creator
	if len(moduleTrigger) > 0 {
		arr := common.QueryStringFromArray(moduleTrigger)
		query += fmt.Sprintf(" && doc._key IN [%v]", arr)
	}
	// query += queryOrgAuth(orgId)
	query = strings.Replace(query, `FILTER && `, `FILTER `, -1)
	if query == `FILTER` {
		query = ""
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)
	limit_query := common.QueryPaginate(offset, limit)

	q := fmt.Sprintf(`
		LET contentTypes = (
			FOR doc IN %v
			%s
			FILTER NOT_NULL(doc.data.contentTypes)
			FOR p IN doc.data.contentTypes
			FOR contentType IN %v
			FILTER p.id == contentType._key 
			RETURN contentType
		)[**]
		FOR doc IN contentTypes
		%s
		%s
		RETURN DISTINCT doc`, common.DbModuleTriggerTable, query, common.DbTriggerContentTypeTable,
		sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbModuleTriggerTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if contentType, err := recordToTriggerContentType(r); err == nil {
			triggerContentTypes = append(triggerContentTypes, contentType)
		}
	}
	return triggerContentTypes, nil
}

func ReadContentCategoryByNameslug(ctx context.Context, nameSlug string) (*static_proto.ContentCategory, error) {
	query := fmt.Sprintf(`FILTER doc.data.name_slug == "%v"`, nameSlug)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET parent = (
			FILTER NOT_NULL(doc.data.parent)
			FOR p IN doc.data.parent
			FOR parent IN %v
			FILTER p.id == parent._key RETURN parent.data
		)
		LET actions = (
			FILTER NOT_NULL(doc.data.actions)
			FOR p IN doc.data.actions
			FOR action IN %v
			FILTER p.id == action._key RETURN action.data
		)
		LET methods = (
			FILTER NOT_NULL(doc.data.trackerMethods)
			FOR p IN doc.data.trackerMethods
			FOR method IN %v
			FILTER p.id == method._key RETURN method.data
		)
		RETURN MERGE_RECURSIVE(doc,{data:{
			parent:parent,
			actions:actions,
			trackerMethods:methods
		}})`, common.DbContentCategoryTable, query,
		common.DbContentParentCategoryTable, common.DbActionTable, common.DbTrackerMethodTable)

	resp, err := runQuery(ctx, q, common.DbContentCategoryTable)

	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToContentCategory(resp.Records[0])
	return data, err
}

// CreateRole creates a role
func CreateRole(ctx context.Context, role *static_proto.Role) error {
	if len(role.Id) == 0 {
		role.Id = uuid.NewUUID().String()
	}
	if role.Created == 0 {
		role.Created = time.Now().Unix()
	}
	role.Updated = time.Now().Unix()

	record, err := roleToRecord(role)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, role.Id, record, record, common.DbRoleTable)
	_, err = runQuery(ctx, q, common.DbTriggerContentTypeTable)
	return err
}

func ReadRoleByNameslug(ctx context.Context, nameSlug string) (*static_proto.Role, error) {
	query := fmt.Sprintf(`FILTER doc.parameter1 == "%v"`, nameSlug)
	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbRoleTable, query)

	resp, err := runQuery(ctx, q, common.DbRoleTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, ErrNotFound
	}
	role, err := recordToRole(resp.Records[0])
	return role, err
}

// AllSetbacks get all setbacks
func AllSetbacks(ctx context.Context, orgId, teamId string, offset, limit int64, sortParameter, sortDirection string) ([]*static_proto.Setback, error) {
	var setbacks []*static_proto.Setback

	var limit_query string
	if offset != -1 && limit != -1 {
		offs, size := queryPaginate(offset, limit)
		limit_query = fmt.Sprintf("LIMIT %s, %s", offs, size)
	}

	sort_query := common.QuerySort(sortParameter, sortDirection)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		%s
		RETURN doc`, common.DbSetbackTable, sort_query, limit_query)

	resp, err := runQuery(ctx, q, common.DbSetbackTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if setback, err := recordToSetback(r); err == nil {
			setbacks = append(setbacks, setback)
		}
	}
	return setbacks, nil
}

// CreateSetback creates a setback
func CreateSetback(ctx context.Context, setback *static_proto.Setback) error {
	if setback.Created == 0 {
		setback.Created = time.Now().Unix()
	}
	if setback.Updated == 0 {
		setback.Updated = time.Now().Unix()
	}
	record, err := setbackToRecord(setback)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { name: "%v", parameter1: "%v" } 
		INSERT %v 
		UPDATE %v 
		IN %v`, setback.Name, setback.OrgId, record, record, common.DbSetbackTable)
	_, err = runQuery(ctx, q, common.DbSetbackTable)
	return err
}

// ReadSetback reads a setback by ID
func ReadSetback(ctx context.Context, id, orgId, teamId string) (*static_proto.Setback, error) {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbSetbackTable, query)

	resp, err := runQuery(ctx, q, common.DbSetbackTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToSetback(resp.Records[0])
	return data, err
}

// DeleteSetback deletes a setback by ID
func DeleteSetback(ctx context.Context, id, orgId, teamId string) error {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc
		IN %v
		%s
		REMOVE doc IN %v`, common.DbSetbackTable, query, common.DbSetbackTable)
	_, err := runQuery(ctx, q, common.DbSetbackTable)
	return err
}

func AutocompleteSetbackSearch(ctx context.Context, title string) ([]*static_proto.Setback, error) {
	var setbacks []*static_proto.Setback

	q := fmt.Sprintf(`
		FOR doc IN %v
		FILTER LIKE(doc.name, "%v",true)
		RETURN doc`, common.DbSetbackTable, `%`+title+`%`)

	resp, err := runQuery(ctx, q, common.DbSetbackTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if setback, err := recordToSetback(r); err == nil {
			setbacks = append(setbacks, setback)
		}
	}
	return setbacks, nil
}

// ReadMarker reads a marker by name_slug
func ReadMarkerByNameslug(ctx context.Context, name_slug string) (*static_proto.Marker, error) {
	query := fmt.Sprintf(`FILTER doc.data.name_slug == "%v"`, name_slug)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		LET maps = (
			FOR e IN %v
				FILTER e._from == doc._id
			FOR i IN %v
				FILTER i._id == e._to
			RETURN i.data
		)
		LET ret = MERGE_RECURSIVE(doc, {"data":{"trackerMethods":maps}})
		RETURN ret`, common.DbMarkerTable, query, common.DbMarkerTrackerEdgeTable, common.DbTrackerMethodTable)

	resp, err := runQuery(ctx, q, common.DbMarkerTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToMarker(resp.Records[0])
	return data, err
}

// ReadBehaviourCategoryAimByNameslug returns categoryAim by name_slug
func ReadBehaviourCategoryAimByNameslug(ctx context.Context, name_slug string) (*static_proto.BehaviourCategoryAim, error) {
	query := fmt.Sprintf(`FILTER doc.data.name_slug == "%v"`, name_slug)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbBehaviourCategoryAimTable, query)

	resp, err := runQuery(ctx, q, common.DbBehaviourCategoryAimTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToBehaviourCategoryAim(resp.Records[0])
	return data, err
}

func ReadTrackerMethodByNameslug(ctx context.Context, name_slug string) (*static_proto.TrackerMethod, error) {
	query := fmt.Sprintf(`FILTER doc.data.name_slug == "%v"`, name_slug)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbTrackerMethodTable, query)

	resp, err := runQuery(ctx, q, common.DbTrackerMethodTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToTrackerMethod(resp.Records[0])
	return data, err
}

func ReadAppByNameslug(ctx context.Context, name_slug string) (*static_proto.App, error) {
	query := fmt.Sprintf(`FILTER doc.data.name_slug == "%v"`, name_slug)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbAppTable, query)

	resp, err := runQuery(ctx, q, common.DbAppTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToApp(resp.Records[0])
	return data, err
}

func ReadWearableByNameslug(ctx context.Context, name_slug string) (*static_proto.Wearable, error) {
	query := fmt.Sprintf(`FILTER doc.data.name_slug == "%v"`, name_slug)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbWearableTable, query)

	resp, err := runQuery(ctx, q, common.DbWearableTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToWearable(resp.Records[0])
	return data, err
}

func ReadDeviceByNameslug(ctx context.Context, name_slug string) (*static_proto.Device, error) {
	query := fmt.Sprintf(`FILTER doc.data.name_slug == "%v"`, name_slug)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbDeviceTable, query)

	resp, err := runQuery(ctx, q, common.DbDeviceTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	data, err := recordToDevice(resp.Records[0])
	return data, err
}
