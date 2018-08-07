package handler

import (
	"context"
	"encoding/json"
	"fmt"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	pubsub_proto "server/static-srv/proto/pubsub"
	"server/survey-srv/db"
	survey_proto "server/survey-srv/proto/survey"
	team_proto "server/team-srv/proto/team"
	user_proto "server/user-srv/proto/user"
	"strconv"

	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/nats"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type SurveyService struct {
	Broker        broker.Broker
	AccountClient account_proto.AccountServiceClient
	KvClient      kv_proto.KvServiceClient
	TeamClient    team_proto.TeamServiceClient
}

func (p *SurveyService) All(ctx context.Context, req *survey_proto.AllRequest, rsp *survey_proto.AllResponse) error {
	log.Info("Received Survey.All request")
	surveys, err := db.All(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(surveys) == 0 || err != nil {
		return common.NotFound(common.SurveySrv, p.All, err, "not found")
	}
	rsp.Data = &survey_proto.ArrData{surveys}
	return nil
}

func (p *SurveyService) New(ctx context.Context, req *survey_proto.NewRequest, rsp *survey_proto.NewResponse) error {
	log.Info("Received Survey.New request")
	m, err := db.New(ctx)
	if err != nil {
		return common.InternalServerError(common.SurveySrv, p.New, err, "server error")
	}
	rsp.Data = &survey_proto.NewResponse_Data{
		Hash:     m["unique_hash"],
		SurveyId: m["survey_id"],
	}
	return nil
}

func (p *SurveyService) Create(ctx context.Context, req *survey_proto.CreateRequest, rsp *survey_proto.CreateResponse) error {
	log.Info("Received Survey.Create request")
	if len(req.Survey.Title) == 0 {
		return common.InternalServerError(common.SurveySrv, p.Create, nil, "survey title empty")
	}
	if len(req.Survey.Id) == 0 {
		req.Survey.Id = uuid.NewUUID().String()
	}
	// create survey
	err := db.Create(ctx, req.Survey)
	if err != nil {
		return common.InternalServerError(common.SurveySrv, p.Create, err, "create error")
	}

	// share survey with user
	req_share := &survey_proto.ShareSurveyRequest{
		Surveys: []*survey_proto.Survey{req.Survey},
		Users:   req.Survey.Shares,
		UserId:  req.UserId,
		OrgId:   req.OrgId,
	}

	rsp_share := &survey_proto.ShareSurveyResponse{}
	if err := p.ShareSurvey(ctx, req_share, rsp_share); err != nil {
		return common.InternalServerError(common.SurveySrv, p.Create, err, "share error")
	}

	// create tags cloud
	if len(req.Survey.Tags) > 0 {
		if _, err := p.KvClient.TagsCloud(context.TODO(), &kv_proto.TagsCloudRequest{
			Index:  common.CLOUD_TAGS_INDEX,
			OrgId:  req.Survey.OrgId,
			Object: common.SURVEY,
			Tags:   req.Survey.Tags,
		}); err != nil {
			return common.InternalServerError(common.SurveySrv, p.Create, err, "tag error")
		}
	}

	rsp.Data = &survey_proto.Data{req.Survey}
	return nil
}

func (p *SurveyService) Read(ctx context.Context, req *survey_proto.ReadRequest, rsp *survey_proto.ReadResponse) error {
	log.Info("Received Survey.Read request")
	survey, err := db.Read(ctx, req.Id, req.OrgId, req.TeamId)
	if survey == nil || err != nil {
		return common.NotFound(common.SurveySrv, p.Read, err, "not found")
	}
	rsp.Data = &survey_proto.Data{survey}
	return nil
}

func (p *SurveyService) Delete(ctx context.Context, req *survey_proto.DeleteRequest, rsp *survey_proto.DeleteResponse) error {
	log.Info("Received Survey.Delete request")
	if err := db.Delete(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.SurveySrv, p.Delete, err, "delete error")
	}
	return nil
}

func (p *SurveyService) Copy(ctx context.Context, req *survey_proto.CopyRequest, rsp *survey_proto.CopyResponse) error {
	log.Info("Received Survey.Copy request")

	// read survey
	survey, err := db.Read(ctx, req.SurveyId, req.OrgId, req.TeamId)
	if survey == nil || err != nil {
		return common.NotFound(common.SurveySrv, p.Copy, err, "not found")
	}

	// create new id
	m, err := db.New(ctx)
	if err != nil {
		return common.InternalServerError(common.SurveySrv, p.Copy, err, "server error")
	}
	// create survey
	survey.Id = m["survey_id"]
	err = db.Create(ctx, survey)
	if err != nil {
		return common.InternalServerError(common.SurveySrv, p.Copy, err, "create error")
	}
	rsp.Data = &survey_proto.Data{survey}
	return nil
}

func (p *SurveyService) Questions(ctx context.Context, req *survey_proto.QuestionsRequest, rsp *survey_proto.QuestionsResponse) error {
	log.Info("Received Survey.Questions request")

	survey, err := db.Read(ctx, req.SurveyId, req.OrgId, req.TeamId)
	if err != nil {
		return err
	}

	questions, err := db.Questions(ctx, req.SurveyId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(questions) == 0 || err != nil {
		return common.NotFound(common.SurveySrv, p.Questions, err, "not found")
	}
	rsp.Data = &survey_proto.QuestionsResponse_ArrData{
		Welcome:   survey.Welcome,
		Thankyou:  survey.Thankyou,
		Questions: questions}
	return nil
}

func (p *SurveyService) QuestionRef(ctx context.Context, req *survey_proto.QuestionRefRequest, rsp *survey_proto.QuestionRefResponse) error {
	log.Info("Received Survey.QuestionRef request")
	question, err := db.QuestionRef(ctx, req.SurveyId, req.QuestionRef)
	if question == nil || err != nil {
		return common.NotFound(common.SurveySrv, p.QuestionRef, err, "not found")
	}
	rsp.Data = &survey_proto.QuestionRefResponse_Data{question}
	return nil
}

func (p *SurveyService) CreateQuestion(ctx context.Context, req *survey_proto.CreateQuestionRequest, rsp *survey_proto.CreateQuestionResponse) error {
	log.Info("Received Survey.CreateQuestion request")
	if len(req.SurveyId) == 0 {
		return common.InternalServerError(common.SurveySrv, p.CreateQuestion, nil, "survey id empty")
	}
	if req.Question == nil {
		return common.InternalServerError(common.SurveySrv, p.CreateQuestion, nil, "question empty")
	}

	err := db.CreateQuestion(ctx, req.SurveyId, req.Question)
	if err != nil {
		return common.InternalServerError(common.SurveySrv, p.CreateQuestion, err, "create error")
	}
	rsp.Data = &survey_proto.CreateQuestionResponse_Data{req.Question}
	return nil
}

func (p *SurveyService) ByCreator(ctx context.Context, req *survey_proto.ByCreatorRequest, rsp *survey_proto.ByCreatorResponse) error {
	log.Info("Received Survey.ByCreator request")
	surveys, err := db.ByCreator(ctx, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(surveys) == 0 || err != nil {
		return common.NotFound(common.SurveySrv, p.ByCreator, err, "not found")
	}
	rsp.Data = &survey_proto.ArrData{surveys}
	return nil
}

func (p *SurveyService) Link(ctx context.Context, req *survey_proto.LinkRequest, rsp *survey_proto.LinkResponse) error {
	log.Info("Received Survey.Link request")
	link, err := db.Link(ctx, req.SurveyId, req.OrgId, req.TeamId)
	if err != nil {
		return common.NotFound(common.SurveySrv, p.Link, err, "not found")
	}
	rsp.UrlLink = link
	return nil
}

func (p *SurveyService) Templates(ctx context.Context, req *survey_proto.TemplatesRequest, rsp *survey_proto.TemplatesResponse) error {
	log.Info("Received Survey.Templates request")
	surveys, err := db.Templates(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(surveys) == 0 || err != nil {
		return common.NotFound(common.SurveySrv, p.Templates, err, "not found")
	}
	rsp.Data = &survey_proto.ArrData{surveys}
	return nil
}

func (p *SurveyService) Filter(ctx context.Context, req *survey_proto.FilterRequest, rsp *survey_proto.FilterResponse) error {
	log.Info("Received Survey.Filter request")
	filters, err := db.Filter(ctx, req.Status, req.Tags, req.RenderTarget, req.Visibility, req.CreatedBy, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(filters) == 0 || err != nil {
		return common.NotFound(common.SurveySrv, p.Filter, err, "not found")
	}
	rsp.Data = &survey_proto.ArrData{filters}
	return nil
}

func (p *SurveyService) Search(ctx context.Context, req *survey_proto.SearchRequest, rsp *survey_proto.SearchResponse) error {
	log.Info("Received Survey.Search request")
	surveys, err := db.Search(ctx, req.Name, req.Description, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(surveys) == 0 || err != nil {
		return common.NotFound(common.SurveySrv, p.Search, err, "not found")
	}
	rsp.Data = &survey_proto.ArrData{surveys}
	return nil
}

func (p *SurveyService) ShareSurvey(ctx context.Context, req *survey_proto.ShareSurveyRequest, rsp *survey_proto.ShareSurveyResponse) error {
	log.Info("Received Survey.ShareSurvey request")

	if len(req.Surveys) == 0 {
		return common.InternalServerError(common.SurveySrv, p.ShareSurvey, nil, "survey empty")
	}
	if len(req.Users) == 0 {
		return common.InternalServerError(common.SurveySrv, p.ShareSurvey, nil, "user empty")
	}

	// checking valid sharedby (employee)
	req_employee := &team_proto.ReadEmployeeInfoRequest{req.UserId}
	rsp_employee, err := p.TeamClient.CheckValidEmployee(ctx, req_employee)
	if err != nil {
		return common.InternalServerError(common.SurveySrv, p.ShareSurvey, err, "CheckValidEmployee is failed")
	}
	if rsp_employee.Valid && rsp_employee.Employee != nil {
		userids, err := db.ShareSurvey(ctx, req.Surveys, req.Users, rsp_employee.Employee.User, req.OrgId)
		if err != nil {
			return common.InternalServerError(common.SurveySrv, p.ShareSurvey, err, "parsing error")
		}
		// send a notification to the users
		if len(userids) > 0 {
			message := fmt.Sprintf(common.MSG_NEW_SURVEY_SHARE, rsp_employee.Employee.User.Firstname)
			alert := &pubsub_proto.Alert{
				Title: fmt.Sprintf("New %v", common.SURVEY),
				Body:  message,
			}
			data := map[string]string{}
			//get current badge count here for user
			data[common.BASE+common.SURVEY_TYPE] = strconv.Itoa(len(req.Surveys))
			p.sendShareNotification(userids, message, alert, data)
	}
	}
	return nil
}

func (p *SurveyService) AutocompleteSearch(ctx context.Context, req *survey_proto.AutocompleteSearchRequest, rsp *survey_proto.AutocompleteSearchResponse) error {
	log.Info("Received Survey.AutocompleteSearch request")

	response, err := db.AutocompleteSearch(ctx, req.Title)
	if len(response) == 0 || err != nil {
		return common.NotFound(common.SurveySrv, p.AutocompleteSearch, err, "not found")
	}
	rsp.Data = &survey_proto.AutocompleteSearchResponse_Data{response}
	return nil
}

func (p *SurveyService) GetTopTags(ctx context.Context, req *survey_proto.GetTopTagsRequest, rsp *survey_proto.GetTopTagsResponse) error {
	log.Info("Received Survey.GetTopTags request")

	rsp_tags, err := p.KvClient.GetTopTags(ctx, &kv_proto.GetTopTagsRequest{
		Index:  common.CLOUD_TAGS_INDEX,
		N:      req.N,
		OrgId:  req.OrgId,
		Object: common.SURVEY,
	})
	if err != nil {
		return common.NotFound(common.SurveySrv, p.GetTopTags, err, "not found")
	}
	rsp.Data = &survey_proto.GetTopTagsResponse_Data{rsp_tags.Tags}
	return nil
}

func (p *SurveyService) AutocompleteTags(ctx context.Context, req *survey_proto.AutocompleteTagsRequest, rsp *survey_proto.AutocompleteTagsResponse) error {
	log.Info("Received Survey.AutocompleteTags request")

	tags, err := db.AutocompleteTags(ctx, req.OrgId, req.Name)
	if len(tags) == 0 || err != nil {
		return common.NotFound(common.SurveySrv, p.AutocompleteTags, err, "not found")
	}
	rsp.Data = &survey_proto.AutocompleteTagsResponse_Data{tags}
	return nil
}

func (p *SurveyService) WarmupCacheSurvey(ctx context.Context, req *survey_proto.WarmupCacheSurveyRequest, rsp *survey_proto.WarmupCacheSurveyResponse) error {
	log.Info("Received Survey.WarmupCacheSurvey request")

	var offset int64
	var limit int64
	offset = 0
	limit = 100

	for {
		items, err := db.All(ctx, "", "", offset, limit, req.SortParameter, req.SortDirection)
		if err != nil || len(items) == 0 {
			break
		}
		for _, item := range items {
			if len(item.Tags) > 0 {
				if _, err := p.KvClient.TagsCloud(ctx, &kv_proto.TagsCloudRequest{
					Index:  common.CLOUD_TAGS_INDEX,
					OrgId:  item.OrgId,
					Object: common.SURVEY,
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
func (p *SurveyService) sendShareNotification(userids []string, message string, alert *pubsub_proto.Alert, data map[string]string) error {
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

func (p *SurveyService) GetShareableSurveys(ctx context.Context, req *user_proto.GetShareableSurveyRequest, rsp *user_proto.GetShareableSurveyResponse) error {
	log.Info("Received Survey.GetShareableSurveys request")
	response, err := db.GetShareableSurveys(ctx, req.CreatedBy, req.Query, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, "", "")
	if err != nil {
		return err
	}
	rsp.Data = &user_proto.GetShareableSurveyResponse_Data{response}
	return nil
}

func (p *SurveyService) Update(ctx context.Context, req *survey_proto.UpdateRequest, rsp *survey_proto.UpdateResponse) error {
	log.Info("Received Survey.Update request")
	return nil
}
