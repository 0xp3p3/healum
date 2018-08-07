package handler

import (
	"context"
	"encoding/json"
	"fmt"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	"server/plan-srv/db"
	plan_proto "server/plan-srv/proto/plan"
	pubsub_proto "server/static-srv/proto/pubsub"
	team_proto "server/team-srv/proto/team"
	"strconv"

	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/nats"
	log "github.com/sirupsen/logrus"
)

type PlanService struct {
	Broker        broker.Broker
	AccountClient account_proto.AccountServiceClient
	KvClient      kv_proto.KvServiceClient
	TeamClient    team_proto.TeamServiceClient
}

func (p *PlanService) All(ctx context.Context, req *plan_proto.AllRequest, rsp *plan_proto.AllResponse) error {
	log.Print("Received Plan.All request")
	plans, err := db.All(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(plans) == 0 || err != nil {
		return common.NotFound(common.PlanSrv, p.All, err, "plan not found")
	}
	rsp.Data = &plan_proto.ArrData{plans}
	return nil
}

func (p *PlanService) Create(ctx context.Context, req *plan_proto.CreateRequest, rsp *plan_proto.CreateResponse) error {
	log.Print("Received Plan.Create request")
	if len(req.Plan.Name) == 0 {
		return common.InternalServerError(common.PlanSrv, p.Create, nil, "plan name error")
	}
	if req.Plan.Creator == nil {
		return common.InternalServerError(common.PlanSrv, p.Create, nil, "plan creator error")
	}
	if len(req.Plan.OrgId) == 0 {
		return common.InternalServerError(common.PlanSrv, p.Create, nil, "plan id error")
	}
	// create
	err := db.Create(ctx, req.Plan)
	if err != nil {
		return common.InternalServerError(common.PlanSrv, p.Create, err, "create error")
	}
	// share plan with user
	req_share := &plan_proto.SharePlanRequest{
		Plans:  []*plan_proto.Plan{req.Plan},
		Users:  req.Plan.Shares,
		UserId: req.UserId,
		OrgId:  req.OrgId,
	}
	rsp_share := &plan_proto.SharePlanResponse{}
	if err := p.SharePlan(ctx, req_share, rsp_share); err != nil {
		return common.InternalServerError(common.PlanSrv, p.Create, err, "share error")
	}

	// create tags cloud
	if len(req.Plan.Tags) > 0 {
		if _, err := p.KvClient.TagsCloud(context.TODO(), &kv_proto.TagsCloudRequest{
			Index:  common.CLOUD_TAGS_INDEX,
			OrgId:  req.Plan.OrgId,
			Object: common.PLAN,
			Tags:   req.Plan.Tags,
		}); err != nil {
			return common.InternalServerError(common.PlanSrv, p.Create, err, "tags cloud error")
		}
	}

	rsp.Data = &plan_proto.Data{req.Plan}
	return nil
}

func (p *PlanService) Read(ctx context.Context, req *plan_proto.ReadRequest, rsp *plan_proto.ReadResponse) error {
	log.Print("Received Plan.Read request")
	plan, err := db.Read(ctx, req.Id, req.OrgId, req.TeamId)
	if plan == nil || err != nil {
		return common.NotFound(common.PlanSrv, p.Read, err, "plan not found")
	}
	rsp.Data = &plan_proto.Data{plan}
	return nil
}
func (p *PlanService) Delete(ctx context.Context, req *plan_proto.DeleteRequest, rsp *plan_proto.DeleteResponse) error {
	log.Print("Received Plan.Delete request")
	if err := db.Delete(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.PlanSrv, p.Delete, err, "delete error")
	}
	return nil
}

func (p *PlanService) Search(ctx context.Context, req *plan_proto.SearchRequest, rsp *plan_proto.SearchResponse) error {
	log.Print("Received Plan.Search request")
	plans, err := db.Search(ctx, req.Name, req.OrgId, req.TeamId, req.Offset, req.Limit, req.From, req.To, req.SortParameter, req.SortDirection)
	if len(plans) == 0 || err != nil {
		return common.NotFound(common.PlanSrv, p.Search, err, "plan not found")
	}
	rsp.Data = &plan_proto.ArrData{plans}
	return nil
}

func (p *PlanService) Templates(ctx context.Context, req *plan_proto.TemplatesRequest, rsp *plan_proto.TemplatesResponse) error {
	log.Print("Received Plan.Templates request")
	plans, err := db.Templates(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(plans) == 0 || err != nil {
		return common.NotFound(common.PlanSrv, p.Templates, err, "plan not found")
	}
	rsp.Data = &plan_proto.ArrData{plans}
	return nil
}

func (p *PlanService) Drafts(ctx context.Context, req *plan_proto.DraftsRequest, rsp *plan_proto.DraftsResponse) error {
	log.Print("Received Plan.Drafts request")
	plans, err := db.Drafts(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(plans) == 0 || err != nil {
		return common.NotFound(common.PlanSrv, p.Drafts, err, "plan not found")
	}
	rsp.Data = &plan_proto.ArrData{plans}
	return nil
}

func (p *PlanService) ByCreator(ctx context.Context, req *plan_proto.ByCreatorRequest, rsp *plan_proto.ByCreatorResponse) error {
	log.Print("Received Plan.ByCreator request")
	plans, err := db.ByCreator(ctx, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(plans) == 0 || err != nil {
		return common.NotFound(common.PlanSrv, p.ByCreator, err, "plan not found")
	}
	rsp.Data = &plan_proto.ArrData{plans}
	return nil
}

func (p *PlanService) Filters(ctx context.Context, req *plan_proto.FiltersRequest, rsp *plan_proto.FiltersResponse) error {
	log.Print("Received Plan.Filters request")
	filters, err := db.Filters(ctx, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(filters) == 0 || err != nil {
		return common.NotFound(common.PlanSrv, p.Filters, err, "plan not found")
	}
	rsp.Data = &plan_proto.FiltersResponse_Data{filters}
	return nil
}

func (p *PlanService) TopFilters(ctx context.Context, req *plan_proto.TopFiltersRequest, rsp *plan_proto.TopFiltersResponse) error {
	log.Print("Received Plan.TopFilters request")
	return nil
}

func (p *PlanService) TimeFilters(ctx context.Context, req *plan_proto.TimeFiltersRequest, rsp *plan_proto.TimeFiltersResponse) error {
	log.Print("Received Plan.TimeFilters request")
	plans, err := db.TimeFilters(ctx, req.StartDate, req.EndDate, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(plans) == 0 || err != nil {
		return common.NotFound(common.PlanSrv, p.TimeFilters, err, "plan not found")
	}
	rsp.Data = &plan_proto.ArrData{plans}
	return nil
}

func (p *PlanService) UseFilters(ctx context.Context, req *plan_proto.UseFiltersRequest, rsp *plan_proto.UseFiltersResponse) error {
	log.Print("Received Plan.UserFilters request")
	return nil
}

func (p *PlanService) SuccessFilters(ctx context.Context, req *plan_proto.SuccessFiltersRequest, rsp *plan_proto.SuccessFiltersResponse) error {
	log.Print("Received Plan.SuccessFilters request")
	return nil
}

func (p *PlanService) ConditionFilters(ctx context.Context, req *plan_proto.ConditionFiltersRequest, rsp *plan_proto.ConditionFiltersResponse) error {
	log.Print("Received Plan.ConditionFilters request")
	return nil
}

func (p *PlanService) GoalFilters(ctx context.Context, req *plan_proto.GoalFiltersRequest, rsp *plan_proto.GoalFiltersResponse) error {
	log.Print("Received Plan.GoalFilters request")
	plans, err := db.GoalFilters(ctx, req.Goals, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(plans) == 0 || err != nil {
		return common.NotFound(common.PlanSrv, p.GoalFilters, err, "plan not found")
	}
	rsp.Data = &plan_proto.ArrData{plans}
	return nil
}

func (p *PlanService) CreatePlanFilter(ctx context.Context, req *plan_proto.CreatePlanFilterRequest, rsp *plan_proto.CreatePlanFilterResponse) error {
	log.Print("Received Plan.CreatePlanFilter request")
	if len(req.Filter.DisplayName) == 0 {
		return common.BadRequest(common.PlanSrv, p.CreatePlanFilter, nil, "display name empty")
	}

	err := db.CreatePlanFilter(ctx, req.Filter)
	if err != nil {
		return common.InternalServerError(common.PlanSrv, p.CreatePlanFilter, nil, "create error")
	}
	rsp.Filter = req.Filter
	return nil
}

func (p *PlanService) SharePlan(ctx context.Context, req *plan_proto.SharePlanRequest, rsp *plan_proto.SharePlanResponse) error {
	log.Print("Received Plan.SharePlan request")

	if len(req.Plans) == 0 {
		return common.BadRequest(common.PlanSrv, p.SharePlan, nil, "plan empty")
	}
	if len(req.Users) == 0 {
		return common.BadRequest(common.PlanSrv, p.SharePlan, nil, "user empty")
	}

	// checking valid sharedby (employee)
	req_employee := &team_proto.ReadEmployeeInfoRequest{req.UserId}
	rsp_employee, err := p.TeamClient.CheckValidEmployee(ctx, req_employee)
	if err != nil {
		return common.InternalServerError(common.PlanSrv, p.SharePlan, err, "CheckValidEmployee is failed")
	}
	if rsp_employee.Valid && rsp_employee.Employee != nil {
		userids, err := db.SharePlan(ctx, req.Plans, req.Users, rsp_employee.Employee.User, req.OrgId)
		if err != nil {
			return common.InternalServerError(common.PlanSrv, p.SharePlan, err, "parsing error")
		}
		// send a notification to the users
		if len(userids) > 0 {
			message := fmt.Sprintf(common.MSG_NEW_PLAN_SHARE, rsp_employee.Employee.User.Firstname)
			alert := &pubsub_proto.Alert{
				Title: fmt.Sprintf("New %v", common.PLAN),
				Body:  message,
			}
			data := map[string]string{}
			//get current badge count here for user
			data[common.BASE + common.PLAN_TYPE] = strconv.Itoa(len(req.Plans))
			p.sendShareNotification(userids, message, alert, data)
	}
	}
	return nil
}

func (p *PlanService) AutocompleteSearch(ctx context.Context, req *plan_proto.AutocompleteSearchRequest, rsp *plan_proto.AutocompleteSearchResponse) error {
	log.Print("Received Plan.AutocompleteSearch request")

	response, err := db.AutocompleteSearch(ctx, req.Title, req.SortParameter, req.SortDirection)
	if len(response) == 0 || err != nil {
		return common.NotFound(common.PlanSrv, p.AutocompleteSearch, err, "not found")
	}
	rsp.Data = &plan_proto.AutocompleteSearchResponse_Data{response}
	return nil
}

func (p *PlanService) GetTopTags(ctx context.Context, req *plan_proto.GetTopTagsRequest, rsp *plan_proto.GetTopTagsResponse) error {
	log.Print("Received Survey.GetTopTags request")

	rsp_tags, err := p.KvClient.GetTopTags(ctx, &kv_proto.GetTopTagsRequest{
		Index:  common.CLOUD_TAGS_INDEX,
		N:      req.N,
		OrgId:  req.OrgId,
		Object: common.PLAN,
	})
	if err != nil {
		return common.NotFound(common.PlanSrv, p.GetTopTags, err, "not found")
	}
	rsp.Data = &plan_proto.GetTopTagsResponse_Data{rsp_tags.Tags}
	return nil
}

func (p *PlanService) AutocompleteTags(ctx context.Context, req *plan_proto.AutocompleteTagsRequest, rsp *plan_proto.AutocompleteTagsResponse) error {
	log.Print("Received Survey.AutocompleteTags request")

	tags, err := db.AutocompleteTags(ctx, req.OrgId, req.Name)
	if len(tags) == 0 || err != nil {
		return common.NotFound(common.PlanSrv, p.AutocompleteTags, err, "not found")
	}
	rsp.Data = &plan_proto.AutocompleteTagsResponse_Data{tags}
	return nil
}

func (p *PlanService) WarmupCache(ctx context.Context, req *plan_proto.WarmupCacheRequest, rsp *plan_proto.WarmupCacheResponse) error {
	log.Print("Received Survey.WarmupCache request")

	var offset int64
	var limit int64
	offset = 0
	limit = 100

	for {
		items, err := db.All(ctx, "", "", offset, limit, "", "")
		if err != nil || len(items) == 0 {
			break
		}
		for _, item := range items {
			if len(item.Tags) > 0 {
				if _, err := p.KvClient.TagsCloud(ctx, &kv_proto.TagsCloudRequest{
					Index:  common.CLOUD_TAGS_INDEX,
					OrgId:  item.OrgId,
					Object: common.PLAN,
					Tags:   item.Tags,
				}); err != nil {
					log.Println("warmup cache err:", err)
				}
			}
		}
		offset += limit
	}

	return nil
}

//FIXME: this is repeated in behaviour, plan, content and survey - combine to single function somewhere? Not sure where
func (p *PlanService) sendShareNotification(userids []string, message string, alert *pubsub_proto.Alert, data map[string]string) error {
	log.Info("Sending notification message for shared resource: ", message, userids)
	msg := &pubsub_proto.PublishBulkNotification{
		Notification: &pubsub_proto.BulkNotification{
		UserIds: userids,
		Message: message,
			Alert:   alert,
			Data:    data,
		},
	}
	if body, err := json.Marshal(msg); err == nil {
		if err := p.Broker.Publish(common.SEND_NOTIFICATION, &broker.Message{Body: body}); err != nil {
			return err
		}
	}
	return nil
}

func (p *PlanService) Update(ctx context.Context, req *plan_proto.UpdateRequest, rsp *plan_proto.UpdateResponse) error {
	log.Print("Received Plan.Update request")

	return nil
}
