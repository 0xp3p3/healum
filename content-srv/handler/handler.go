package handler

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	account_proto "server/account-srv/proto/account"
	"server/common"
	"server/content-srv/db"
	content_proto "server/content-srv/proto/content"
	kv_proto "server/kv-srv/proto/kv"
	pubsub_proto "server/static-srv/proto/pubsub"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
	user_proto "server/user-srv/proto/user"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/nats"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

type ContentService struct {
	Broker        broker.Broker
	StaticClient  static_proto.StaticServiceClient
	AccountClient account_proto.AccountServiceClient
	KvClient      kv_proto.KvServiceClient
	TeamClient    team_proto.TeamServiceClient
}

func HashFromObject(obj interface{}) string {
	b, err := bson.Marshal(obj)
	if err != nil {
		log.WithField("err", err).Error("not serialized")
		return ""
	}
	hash := sha256.New()
	hash.Write(b)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (p *ContentService) AllSources(ctx context.Context, req *content_proto.AllSourcesRequest, rsp *content_proto.AllSourcesResponse) error {
	log.Info("Received Content.AllSources request")
	sources, err := db.AllSources(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(sources) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.AllSources, err, "not found")
	}
	rsp.Data = &content_proto.SourceArrData{sources}
	return nil
}

func (p *ContentService) CreateSource(ctx context.Context, req *content_proto.CreateSourceRequest, rsp *content_proto.CreateSourceResponse) error {
	log.Info("Received Content.CreateSource request")
	if len(req.Source.Name) == 0 {
		return common.Forbidden(common.ContentSrv, p.CreateSource, nil, "name empty")
	}
	if len(req.Source.Id) == 0 {
		req.Source.Id = uuid.NewUUID().String()
	}
	if err := db.CreateSource(ctx, req.Source); err != nil {
		return common.InternalServerError(common.ContentSrv, p.CreateSource, err, "create error")
	}
	rsp.Data = &content_proto.SourceData{req.Source}
	return nil
}

func (p *ContentService) ReadSource(ctx context.Context, req *content_proto.ReadSourceRequest, rsp *content_proto.ReadSourceResponse) error {
	log.Info("Received Content.ReadSource request")
	source, err := db.ReadSource(ctx, req.Id, req.OrgId, req.TeamId)
	if source == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.ReadSource, err, "not found")
	}
	rsp.Data = &content_proto.SourceData{source}
	return nil
}

func (p *ContentService) DeleteSource(ctx context.Context, req *content_proto.DeleteSourceRequest, rsp *content_proto.DeleteSourceResponse) error {
	log.Info("Received Content.DeleteSource request")
	if err := db.DeleteSource(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.ContentSrv, p.DeleteSource, err, "server error")
	}
	return nil
}

func (p *ContentService) AllTaxonomys(ctx context.Context, req *content_proto.AllTaxonomysRequest, rsp *content_proto.AllTaxonomysResponse) error {
	log.Info("Received Content.AllTaxonomys request")
	taxonomys, err := db.AllTaxonomys(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(taxonomys) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.AllTaxonomys, err, "not found")
	}
	rsp.Data = &content_proto.TaxonomyArrData{taxonomys}
	return nil
}

func (p *ContentService) CreateTaxonomy(ctx context.Context, req *content_proto.CreateTaxonomyRequest, rsp *content_proto.CreateTaxonomyResponse) error {
	log.Info("Received Content.CreateTaxonomy request")
	if len(req.Taxonomy.Name) == 0 {
		return common.Forbidden(common.ContentSrv, p.CreateTaxonomy, nil, "name empty")
	}
	if len(req.Taxonomy.Id) == 0 {
		req.Taxonomy.Id = uuid.NewUUID().String()
	}

	if err := db.CreateTaxonomy(ctx, req.Taxonomy); err != nil {
		return common.InternalServerError(common.ContentSrv, p.CreateTaxonomy, err, "create error")
	}
	rsp.Data = &content_proto.TaxonomyData{req.Taxonomy}
	return nil
}

func (p *ContentService) ReadTaxonomy(ctx context.Context, req *content_proto.ReadTaxonomyRequest, rsp *content_proto.ReadTaxonomyResponse) error {
	log.Info("Received Content.ReadTaxonomy request")
	taxonomy, err := db.ReadTaxonomy(ctx, req.Id, req.OrgId, req.TeamId)
	if taxonomy == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.ReadTaxonomy, err, "not found")
	}
	rsp.Data = &content_proto.TaxonomyData{taxonomy}
	return nil
}

func (p *ContentService) DeleteTaxonomy(ctx context.Context, req *content_proto.DeleteTaxonomyRequest, rsp *content_proto.DeleteTaxonomyResponse) error {
	log.Info("Received Content.DeleteTaxonomy request")
	if err := db.DeleteTaxonomy(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.ContentSrv, p.DeleteTaxonomy, err, "server error")
	}
	return nil
}

func (p *ContentService) AllContents(ctx context.Context, req *content_proto.AllContentsRequest, rsp *content_proto.AllContentsResponse) error {
	log.Info("Received Content.AllContents request")
	contents, err := db.AllContents(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(contents) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.AllContents, err, "not found")
	}
	rsp.Data = &content_proto.ContentArrData{contents}
	return nil
}

func (p *ContentService) CreateContent(ctx context.Context, req *content_proto.CreateContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreateContent request")
	if len(req.Content.Title) == 0 {
		return common.Forbidden(common.ContentSrv, p.CreateContent, nil, "title empty")
	}

	if err := db.CreateContent(ctx, req.Content); err != nil {
		return common.Forbidden(common.ContentSrv, p.CreateContent, err, "create error")
	}

	// share content with user
	if req.Content.Shares != nil && len(req.Content.Shares) > 0 {
		req_share := &content_proto.ShareContentRequest{
			Contents: []*content_proto.Content{req.Content},
			Users:    req.Content.Shares,
			UserId:   req.UserId,
			OrgId:    req.OrgId,
		}
		rsp_share := &content_proto.ShareContentResponse{}
		if err := p.ShareContent(ctx, req_share, rsp_share); err != nil {
			log.Error("content sharing is failed err:", err)
			return err
		}
	}

	// tags cloud
	if len(req.Content.Tags) > 0 {
		tags := []string{}
		for _, tag := range req.Content.Tags {
			tags = append(tags, tag.Name)
		}
		if _, err := p.KvClient.TagsCloud(context.TODO(), &kv_proto.TagsCloudRequest{
			Index:  common.CLOUD_TAGS_INDEX,
			OrgId:  req.Content.OrgId,
			Object: common.CONTENT,
			Tags:   tags,
		}); err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateContent, err, "tag error")
			return err
		}
	}

	rsp.Data = &content_proto.ContentData{req.Content}
	return nil
}

func (p *ContentService) CreateActivityContent(ctx context.Context, req *content_proto.CreateActivityContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreateActivityContent request")

	// get activity-category from static-srv
	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslugOrCreate(ctx, &static_proto.ReadByNameslugRequest{common.CONTENT_ACTIVITY_TYPE})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.CreateActivityContent, err, "category not found")
	}

	for _, activity := range req.Activitys {
		any, err := common.AnyFromObject(activity, common.CONTENT_ACTIVITY_TYPE)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateActivityContent, err, "any parsing error")
		}

		activity.Id = ""
		content := &content_proto.Content{
			Category: rsp_category.Data.ContentCategory,
			Item:     any,
			Hash:     HashFromObject(activity),
		}

		err = db.CreateContent(ctx, content)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateActivityContent, err, "create error")
		}
		// rsp.Data = &content_proto.ContentData{content}
	}
	return nil
}

func (p *ContentService) CreateRecipeContent(ctx context.Context, req *content_proto.CreateRecipeContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreateRecipeContent request")

	// get recipe-category from static-srv
	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslugOrCreate(ctx, &static_proto.ReadByNameslugRequest{common.CONTENT_RECIPE_TYPE})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.CreateRecipeContent, err, "category not found")
	}

	for _, recipe := range req.Recipes {
		any, err := common.AnyFromObject(recipe, common.CONTENT_RECIPE_TYPE)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateRecipeContent, err, "any parsing error")
		}
		//set all appropriat evalues here reqiure for the Content
		content := &content_proto.Content{
			Url:      recipe.Url,
			Image:    recipe.Image,
			Title:    recipe.Title,
			Author:   recipe.Source,
			Category: rsp_category.Data.ContentCategory,
			Item:     any,
			Hash:     HashFromObject(recipe),
		}
		err = db.CreateContent(ctx, content)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateRecipeContent, err, "create error")
		}
		// rsp.Data = &content_proto.ContentData{content}
	}
	return nil
}

func (p *ContentService) CreateExerciseContent(ctx context.Context, req *content_proto.CreateExerciseContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreateExerciseContent request")

	// get exercise-category from static-srv
	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslugOrCreate(ctx, &static_proto.ReadByNameslugRequest{common.CONTENT_EXERCISE_TYPE})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.CreateExerciseContent, err, "category not found")
	}

	for _, exercise := range req.Exercises {
		any, err := common.AnyFromObject(exercise, common.CONTENT_EXERCISE_TYPE)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateExerciseContent, err, "any parsing error")
		}
		content := &content_proto.Content{
			Category: rsp_category.Data.ContentCategory,
			Item:     any,
			Hash:     HashFromObject(exercise),
		}
		err = db.CreateContent(ctx, content)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateExerciseContent, err, "create error")
		}
		// rsp.Data = &content_proto.ContentData{content}
	}
	return nil
}

func (p *ContentService) CreateArticleContent(ctx context.Context, req *content_proto.CreateArticleContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreateArticleContent request")

	// get exercise-category from static-srv
	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslugOrCreate(ctx, &static_proto.ReadByNameslugRequest{common.CONTENT_ARTICLE_TYPE})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.CreateArticleContent, err, "category not found")
	}

	for _, article := range req.Articles {
		any, err := common.AnyFromObject(article, common.CONTENT_ARTICLE_TYPE)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateArticleContent, err, "any parsing error")
		}
		content := &content_proto.Content{
			Category: rsp_category.Data.ContentCategory,
			Item:     any,
			Hash:     HashFromObject(article),
		}
		err = db.CreateContent(ctx, content)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateArticleContent, err, "create error")
		}
		// rsp.Data = &content_proto.ContentData{content}
	}
	return nil
}

func (p *ContentService) CreatePlaceContent(ctx context.Context, req *content_proto.CreatePlaceContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreatePlaceContent request")

	// get place-category from static-srv
	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslugOrCreate(ctx, &static_proto.ReadByNameslugRequest{common.CONTENT_PLACE_TYPE})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.CreatePlaceContent, err, "category not found")
	}

	for _, place := range req.Places {
		any, err := common.AnyFromObject(place, common.CONTENT_PLACE_TYPE)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreatePlaceContent, err, "any parsing error")
		}
		content := &content_proto.Content{
			Category: rsp_category.Data.ContentCategory,
			Item:     any,
			Hash:     HashFromObject(place),
		}
		err = db.CreateContent(ctx, content)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreatePlaceContent, err, "create error")
		}
		// rsp.Data = &content_proto.ContentData{content}
	}
	return nil
}

func (p *ContentService) CreateWellbeingContent(ctx context.Context, req *content_proto.CreateWellbeingContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreateWellbeingContent request")

	// get wellbeing-category from static-srv
	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslugOrCreate(ctx, &static_proto.ReadByNameslugRequest{common.CONTENT_WELLBEING_TYPE})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.CreateWellbeingContent, err, "category not found")
	}

	for _, wellbeing := range req.Wellbeings {
		any, err := common.AnyFromObject(wellbeing, common.CONTENT_WELLBEING_TYPE)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateWellbeingContent, err, "any parsing error")
		}
		content := &content_proto.Content{
			Category: rsp_category.Data.ContentCategory,
			Item:     any,
			Hash:     HashFromObject(wellbeing),
		}
		err = db.CreateContent(ctx, content)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateWellbeingContent, err, "create error")
		}
		// rsp.Data = &content_proto.ContentData{content}
	}
	return nil
}

func (p *ContentService) CreateVideoContent(ctx context.Context, req *content_proto.CreateVideoContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreateVideoContent request")

	// get video-category from static-srv
	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslugOrCreate(ctx, &static_proto.ReadByNameslugRequest{common.CONTENT_VIDEO_TYPE})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.CreateVideoContent, err, "category not found")
	}

	for _, video := range req.Videos {
		any, err := common.AnyFromObject(video, common.CONTENT_VIDEO_TYPE)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateVideoContent, err, "any parsing error")
		}
		content := &content_proto.Content{
			Category: rsp_category.Data.ContentCategory,
			Item:     any,
			Hash:     HashFromObject(video),
		}
		err = db.CreateContent(ctx, content)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateVideoContent, err, "create error")
		}
		// rsp.Data = &content_proto.ContentData{content}
	}
	return nil
}

func (p *ContentService) CreateProductContent(ctx context.Context, req *content_proto.CreateProductContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreateProductContent request")

	// get product-category from static-srv
	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslugOrCreate(ctx, &static_proto.ReadByNameslugRequest{common.CONTENT_PRODUCT_TYPE})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.CreateProductContent, err, "category not found")
	}

	for _, product := range req.Products {
		any, err := common.AnyFromObject(product, common.CONTENT_PRODUCT_TYPE)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateProductContent, err, "any parsing error")
		}
		content := &content_proto.Content{
			Category: rsp_category.Data.ContentCategory,
			Item:     any,
			Hash:     HashFromObject(product),
		}
		err = db.CreateContent(ctx, content)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateProductContent, err, "create error")
		}
		// rsp.Data = &content_proto.ContentData{content}
	}
	return nil
}

func (p *ContentService) CreateServiceContent(ctx context.Context, req *content_proto.CreateServiceContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreateServiceContent request")

	// get service-category from static-srv
	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslugOrCreate(ctx, &static_proto.ReadByNameslugRequest{common.CONTENT_SERVICE_TYPE})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.CreateServiceContent, err, "category not found")
	}

	for _, service := range req.Services {
		any, err := common.AnyFromObject(service, common.CONTENT_SERVICE_TYPE)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateServiceContent, err, "any parsing error")
		}
		content := &content_proto.Content{
			Category: rsp_category.Data.ContentCategory,
			Item:     any,
			Hash:     HashFromObject(service),
		}
		err = db.CreateContent(ctx, content)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateServiceContent, err, "create error")
		}
		// rsp.Data = &content_proto.ContentData{content}
	}
	return nil
}

func (p *ContentService) CreateEventContent(ctx context.Context, req *content_proto.CreateEventContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreateEventContent request")

	// get event-category from static-srv
	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslugOrCreate(ctx, &static_proto.ReadByNameslugRequest{common.CONTENT_EVENT_TYPE})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.CreateEventContent, err, "category not found")
	}

	for _, event := range req.Events {
		if event.Created == 0 {
			event.Created = time.Now().Unix()
		}
		event.Updated = time.Now().Unix()
		any, err := common.AnyFromObject(event, common.CONTENT_EVENT_TYPE)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateEventContent, err, "any parsing error")
		}
		content := &content_proto.Content{
			Category: rsp_category.Data.ContentCategory,
			Item:     any,
			Hash:     HashFromObject(event),
		}
		err = db.CreateContent(ctx, content)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateEventContent, err, "create error")
		}
		// rsp.Data = &content_proto.ContentData{content}
	}
	return nil
}

func (p *ContentService) CreateResearchContent(ctx context.Context, req *content_proto.CreateResearchContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreateResearchContent request")

	// get research-category from static-srv
	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslugOrCreate(ctx, &static_proto.ReadByNameslugRequest{common.CONTENT_RESEARCH_TYPE})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.CreateResearchContent, err, "category not found")
	}

	for _, research := range req.Researchs {
		if research.Created == 0 {
			research.Created = time.Now().Unix()
		}
		research.Updated = time.Now().Unix()
		any, err := common.AnyFromObject(research, common.CONTENT_RESEARCH_TYPE)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateResearchContent, err, "any parsing error")
		}
		content := &content_proto.Content{
			Category: rsp_category.Data.ContentCategory,
			Item:     any,
			Hash:     HashFromObject(research),
		}
		err = db.CreateContent(ctx, content)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateResearchContent, err, "create error")
		}
		// rsp.Data = &content_proto.ContentData{content}
	}
	return nil
}

func (p *ContentService) CreateAppContent(ctx context.Context, req *content_proto.CreateAppContentRequest, rsp *content_proto.CreateContentResponse) error {
	log.Info("Received Content.CreateAppContent request")

	// get app-category from static-srv
	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslugOrCreate(ctx, &static_proto.ReadByNameslugRequest{common.CONTENT_APP_TYPE})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.CreateAppContent, err, "category not found")
	}

	for _, app := range req.Apps {
		if app.Created == 0 {
			app.Created = time.Now().Unix()
		}
		app.Updated = time.Now().Unix()
		any, err := common.AnyFromObject(app, common.CONTENT_APP_TYPE)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateAppContent, err, "any parsing error")
		}
		content := &content_proto.Content{
			Category: rsp_category.Data.ContentCategory,
			Item:     any,
			Hash:     HashFromObject(app),
		}
		err = db.CreateContent(ctx, content)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.CreateAppContent, err, "create error")
		}
		// rsp.Data = &content_proto.ContentData{content}
	}
	return nil
}

func (p *ContentService) CreateContentRecommendation(ctx context.Context, req *content_proto.CreateContentRecommendationRequest, rsp *content_proto.CreateContentRecommendationResponse) error {
	log.Info("Received Content.CreateContentRecommendation request")

	if err := db.CreateContentRecommendation(ctx, req.Recommendation); err != nil {
		return common.InternalServerError(common.ContentSrv, p.CreateContentRecommendation, err, "create error")
	}
	rsp.Data = &content_proto.CreateContentRecommendationResponse_Data{req.Recommendation}
	return nil
}

func (p *ContentService) ReadContent(ctx context.Context, req *content_proto.ReadContentRequest, rsp *content_proto.ReadContentResponse) error {
	log.Info("Received Content.ReadContent request")
	content, err := db.ReadContent(ctx, req.Id, req.OrgId, req.TeamId)
	if content == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.ReadContent, err, "content not found")
	}
	rsp.Data = &content_proto.ContentData{content}
	return nil
}

func (p *ContentService) DeleteContent(ctx context.Context, req *content_proto.DeleteContentRequest, rsp *content_proto.DeleteContentResponse) error {
	log.Info("Received Content.DeleteContent request")
	if err := db.DeleteContent(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.ContentSrv, p.DeleteContent, err, "delete error")
	}
	return nil
}

func (p *ContentService) AllContentCategoryItems(ctx context.Context, req *content_proto.AllContentCategoryItemsRequest, rsp *content_proto.AllContentCategoryItemsResponse) error {
	log.Info("Received Content.AllContentCategoryItems request")
	contentCategoryItems, err := db.AllContentCategoryItems(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(contentCategoryItems) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.AllContentCategoryItems, err, "content not found")
	}
	rsp.Data = &content_proto.ContentCategoryItemArrData{contentCategoryItems}
	return nil
}

func (p *ContentService) CreateContentCategoryItem(ctx context.Context, req *content_proto.CreateContentCategoryItemRequest, rsp *content_proto.CreateContentCategoryItemResponse) error {
	log.Info("Received Content.CreateContentCategoryItem request")
	if len(req.ContentCategoryItem.Name) == 0 {
		return common.BadRequest(common.ContentSrv, p.CreateContentCategoryItem, nil, "name empty")
	}
	if len(req.ContentCategoryItem.Id) == 0 {
		req.ContentCategoryItem.Id = uuid.NewUUID().String()
	}

	err := db.CreateContentCategoryItem(ctx, req.ContentCategoryItem)
	if err != nil {
		return common.InternalServerError(common.ContentSrv, p.CreateContentCategoryItem, err, "create error")
	}
	rsp.Data = &content_proto.ContentCategoryItemData{req.ContentCategoryItem}
	return nil
}

func (p *ContentService) ReadContentCategoryItem(ctx context.Context, req *content_proto.ReadContentCategoryItemRequest, rsp *content_proto.ReadContentCategoryItemResponse) error {
	log.Info("Received Content.ReadContentCategoryItem request")
	contentCategoryItem, err := db.ReadContentCategoryItem(ctx, req.Id, req.OrgId, req.TeamId)
	if contentCategoryItem == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.ReadContentCategoryItem, err, "category not found")
	}
	rsp.Data = &content_proto.ContentCategoryItemData{contentCategoryItem}
	return nil
}

func (p *ContentService) DeleteContentCategoryItem(ctx context.Context, req *content_proto.DeleteContentCategoryItemRequest, rsp *content_proto.DeleteContentCategoryItemResponse) error {
	log.Info("Received Content.DeleteContentCategoryItem request")
	if err := db.DeleteContentCategoryItem(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.ContentSrv, p.DeleteContentCategoryItem, err, "delete error")
	}
	return nil
}

func (p *ContentService) AllContentRules(ctx context.Context, req *content_proto.AllContentRulesRequest, rsp *content_proto.AllContentRulesResponse) error {
	log.Info("Received Content.AllContentRules request")
	contentRules, err := db.AllContentRules(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(contentRules) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.AllContentRules, err, "category not found")
	}
	rsp.Data = &content_proto.ContentRuleArrData{contentRules}
	return nil
}

func (p *ContentService) CreateContentRule(ctx context.Context, req *content_proto.CreateContentRuleRequest, rsp *content_proto.CreateContentRuleResponse) error {
	log.Info("Received Content.CreateContentRule request")
	if len(req.ContentRule.Id) == 0 {
		req.ContentRule.Id = uuid.NewUUID().String()
	}

	err := db.CreateContentRule(ctx, req.ContentRule)
	if err != nil {
		return common.InternalServerError(common.ContentSrv, p.CreateContentRule, err, "create error")
	}
	rsp.Data = &content_proto.ContentRuleData{req.ContentRule}
	return nil
}

func (p *ContentService) ReadContentRule(ctx context.Context, req *content_proto.ReadContentRuleRequest, rsp *content_proto.ReadContentRuleResponse) error {
	log.Info("Received Content.ReadContentRule request")
	contentRule, err := db.ReadContentRule(ctx, req.Id, req.OrgId, req.TeamId)
	if contentRule == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.ReadContentRule, err, "not found")
	}
	rsp.Data = &content_proto.ContentRuleData{contentRule}
	return nil
}

func (p *ContentService) DeleteContentRule(ctx context.Context, req *content_proto.DeleteContentRuleRequest, rsp *content_proto.DeleteContentRuleResponse) error {
	log.Info("Received Content.DeleteContentRule request")
	if err := db.DeleteContentRule(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.ContentSrv, p.DeleteContentRule, err, "server error")
	}
	return nil
}

func (p *ContentService) FilterContent(ctx context.Context, req *content_proto.FilterContentRequest, rsp *content_proto.FilterContentResponse) error {
	log.Info("Received Content.FilterContent request")
	contents, err := db.FilterContent(ctx, req)
	if len(contents) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.FilterContent, err, "content not found")
	}
	rsp.Data = &content_proto.ContentArrData{contents}
	return nil
}

func (p *ContentService) SearchContent(ctx context.Context, req *content_proto.SearchContentRequest, rsp *content_proto.SearchContentResponse) error {
	log.Info("Received Content.Search request")

	data, err := db.SearchContent(ctx, req.Title, req.Summary, req.Description, req.OrgId, req.TeamId, req.Offset, req.Limit)
	if err != nil {
		return err
	}

	rsp.Data = data
	return nil
}

func (p *ContentService) ShareContent(ctx context.Context, req *content_proto.ShareContentRequest, rsp *content_proto.ShareContentResponse) error {
	log.Info("Received Content.ShareContent request")

	if len(req.Contents) == 0 {
		return common.BadRequest(common.ContentSrv, p.ShareContent, nil, "contents empty")
	}
	if len(req.Users) == 0 {
		return common.BadRequest(common.ContentSrv, p.ShareContent, nil, "users empty")
	}

	// checking valid sharedby (employee)
	req_employee := &team_proto.ReadEmployeeInfoRequest{req.UserId}
	rsp_employee, err := p.TeamClient.CheckValidEmployee(ctx, req_employee)
	if err != nil {
		return common.InternalServerError(common.ContentSrv, p.ShareContent, err, "CheckValidEmployee is failed")
	}
	if rsp_employee.Valid && rsp_employee.Employee != nil {
		userids, err := db.ShareContent(ctx, req.Contents, req.Users, rsp_employee.Employee.User, req.OrgId)
		if err != nil {
			return common.InternalServerError(common.ContentSrv, p.ShareContent, err, "share error")
		}
		// send a notification to the users
		if len(userids) > 0 {
			message := fmt.Sprintf(common.MSG_NEW_CONTENT_SHARE, rsp_employee.Employee.User.Firstname)
			alert := &pubsub_proto.Alert{
				Title: fmt.Sprintf("New %v", common.CONTENT),
				Body:  message,
			}
			data := map[string]string{}
			//get current badge count here for user
			data[common.BASE+common.CONTENT_TYPE] = strconv.Itoa(len(req.Contents))
			p.sendShareNotification(userids, message, alert, data)
		}
	}
	return nil
}

func (p *ContentService) Subscribe(ctx context.Context, req *pubsub_proto.SubscribeRequest, rsp *pubsub_proto.SubscribeResponse) error {
	log.Info("Received Content.Subscribe request")

	_, err := p.Broker.Subscribe(req.Channel, func(pub broker.Publication) error {
		var err error
		switch req.Channel {
		case common.CREATE_ACTIVITY_CONTENT:
			msg := []*content_proto.Activity{}
			if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create activity content error")
			}
			go p.CreateActivityContent(ctx, &content_proto.CreateActivityContentRequest{Activitys: msg}, &content_proto.CreateContentResponse{})
		case common.CREATE_RECIPE_CONTENT:
			msg := []*content_proto.Recipe{}
			if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create recipe content error")
			}
			go p.CreateRecipeContent(ctx, &content_proto.CreateRecipeContentRequest{Recipes: msg}, &content_proto.CreateContentResponse{})
		case common.CREATE_EXCERCISE_CONTENT:
			msg := []*content_proto.Exercise{}
			if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create excercise content error")

			}
			go p.CreateExerciseContent(ctx, &content_proto.CreateExerciseContentRequest{Exercises: msg}, &content_proto.CreateContentResponse{})
		case common.CREATE_ARTICLE_CONTENT:
			msg := []*content_proto.Article{}
			if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create article content error")
			}
			go p.CreateArticleContent(ctx, &content_proto.CreateArticleContentRequest{Articles: msg}, &content_proto.CreateContentResponse{})
		case common.CREATE_PLACE_CONTENT:
			msg := []*content_proto.Place{}
			if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create place content error")
			}
			go p.CreatePlaceContent(ctx, &content_proto.CreatePlaceContentRequest{Places: msg}, &content_proto.CreateContentResponse{})
		case common.CREATE_WELLBEING_CONTENT:
			msg := []*content_proto.Wellbeing{}
			if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create wellbeing content error")
			}
			go p.CreateWellbeingContent(ctx, &content_proto.CreateWellbeingContentRequest{Wellbeings: msg}, &content_proto.CreateContentResponse{})
		case common.CREATE_CONTENT_RECOMMENDATION:
			msg := &content_proto.ContentRecommendation{}
			if err := jsonpb.Unmarshal(strings.NewReader(string(pub.Message().Body)), msg); err != nil {
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create content recommendation error")
			}
			go p.CreateContentRecommendation(ctx, &content_proto.CreateContentRecommendationRequest{Recommendation: msg}, &content_proto.CreateContentRecommendationResponse{})
		case common.CREATE_VIDEO_CONTENT:
			msg := []*content_proto.Video{}
			if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create video content error")
			}
			go p.CreateVideoContent(ctx, &content_proto.CreateVideoContentRequest{Videos: msg}, &content_proto.CreateContentResponse{})
		case common.CREATE_PRODUCT_CONTENT:
			msg := []*content_proto.Product{}
			if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create product content error")
			}
			go p.CreateProductContent(ctx, &content_proto.CreateProductContentRequest{Products: msg}, &content_proto.CreateContentResponse{})
		case common.CREATE_SERVICE_CONTENT:
			msg := []*content_proto.Service{}
			if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
				common.InternalServerError(common.ContentSrv, p.Subscribe, err, common.CREATE_SERVICE_CONTENT)
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create service content error")
			}
			go p.CreateServiceContent(ctx, &content_proto.CreateServiceContentRequest{Services: msg}, &content_proto.CreateContentResponse{})
		case common.CREATE_EVENT_CONTENT:
			msg := []*content_proto.Event{}
			if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create event content error")
			}
			go p.CreateEventContent(ctx, &content_proto.CreateEventContentRequest{Events: msg}, &content_proto.CreateContentResponse{})
		case common.CREATE_RESEARCH_CONTENT:
			msg := []*content_proto.Research{}
			if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create research content error")
			}
			go p.CreateResearchContent(ctx, &content_proto.CreateResearchContentRequest{Researchs: msg}, &content_proto.CreateContentResponse{})
		case common.CREATE_APP_CONTENT:
			msg := []*static_proto.App{}
			if err := json.Unmarshal(pub.Message().Body, &msg); err != nil {
				return common.InternalServerError(common.ContentSrv, p.Subscribe, err, "create app content error")
			}
			go p.CreateAppContent(ctx, &content_proto.CreateAppContentRequest{Apps: msg}, &content_proto.CreateContentResponse{})
		}

		return err
	})
	return err
}

func (p *ContentService) GetContentCategorys(ctx context.Context, req *content_proto.GetContentCategorysRequest, rsp *content_proto.GetContentCategorysResponse) error {
	log.Info("Received Content.GetContentCategorys request")

	categorys, err := db.GetContentCategorys(ctx)
	if len(categorys) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.GetContentCategorys, err, "content category not found")
	}
	rsp.Data = &content_proto.GetContentCategorysResponse_Data{categorys}
	return nil
}

func (p *ContentService) GetContentDetail(ctx context.Context, req *content_proto.GetContentDetailRequest, rsp *content_proto.GetContentDetailResponse) error {
	log.Info("Received Content.GetContentDetail request")
	data, err := db.GetContentDetail(ctx, req.ContentId)
	if data == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.GetContentDetail, err, "content detail not found")
	}
	rsp.Data = data
	return nil
}

func (p *ContentService) GetContentByCategory(ctx context.Context, req *content_proto.GetContentByCategoryRequest, rsp *content_proto.GetContentByCategoryResponse) error {
	log.Info("Received Content.GetContentByCategory request")

	contents, err := db.GetContentByCategory(ctx, req.CategoryId, req.Offset, req.Limit)
	if len(contents) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.GetContentByCategory, err, "content by category not found")
	}
	rsp.Data = &content_proto.GetContentByCategoryResponse_Data{contents}
	return nil
}

func (p *ContentService) GetFiltersForCategory(ctx context.Context, req *content_proto.GetFiltersForCategoryRequest, rsp *content_proto.GetFiltersForCategoryResponse) error {
	log.Info("Received Content.GetFiltersForCategory request")

	items, err := db.GetContentCategoryItemsByCategory(ctx, req.CategoryId)
	if len(items) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.GetFiltersForCategory, err, "content category item not found")
	}
	rsp.Data = &content_proto.ContentCategoryItemArrData{items}
	return nil
}

func (p *ContentService) FiltersAutocomplete(ctx context.Context, req *content_proto.FiltersAutocompleteRequest, rsp *content_proto.FiltersAutocompleteResponse) error {
	log.Info("Received Content.FiltersAutocomplete request")

	categorys, err := db.FilterCategoryAutocomplete(ctx, req.CategoryId, req.Name)
	if len(categorys) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.FiltersAutocomplete, err, "category not found")
	}
	rsp.Data = &content_proto.FiltersAutocompleteResponse_Data{categorys}
	return nil
}

func (p *ContentService) FilterContentInParticularCategory(ctx context.Context, req *content_proto.FilterContentInParticularCategoryRequest, rsp *content_proto.FilterContentInParticularCategoryResponse) error {
	log.Info("Received Content.FilterContentInParticularCategory request")

	contents, err := db.FilterContentInParticularCategory(ctx, req.CategoryId, req.ContentCategoryItems)
	if len(contents) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.FilterContentInParticularCategory, err, "contents not found")
	}
	rsp.Data = &content_proto.FilterContentInParticularCategoryResponse_Data{contents}
	return nil
}

func (p *ContentService) GetContentRecommendationByUser(ctx context.Context, req *content_proto.GetContentRecommendationByUserRequest, rsp *content_proto.GetContentRecommendationByUserResponse) error {
	log.Info("Received Content.GetContentRecommendationByUser request")

	recommendations, err := db.GetContentRecommendationByUser(ctx, req.UserId, req.OrgId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(recommendations) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.GetContentRecommendationByUser, err, "contents recommendation by user not found")
	}
	rsp.Data = &content_proto.GetContentRecommendationByUserResponse_Data{recommendations}
	return nil
}

func (p *ContentService) GetContentRecommendationByCategory(ctx context.Context, req *content_proto.GetContentRecommendationByCategoryRequest, rsp *content_proto.GetContentRecommendationByCategoryResponse) error {
	log.Info("Received Content.GetContentRecommendationByCategory request")

	recommendations, err := db.GetContentRecommendationByCategory(ctx, req.UserId, req.OrgId, req.CategoryId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(recommendations) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.GetContentRecommendationByCategory, err, "contents recommendation by category not found")
	}
	rsp.Data = &content_proto.GetContentRecommendationByCategoryResponse_Data{recommendations}
	return nil
}

func (p *ContentService) GetRandomItems(ctx context.Context, req *content_proto.GetRandomItemsRequest, rsp *content_proto.GetRandomItemsResponse) error {
	log.Info("Received Content.GetRandomItems request")

	contents, err := db.GetRandomItems(ctx, req.Count)
	if len(contents) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.GetRandomItems, err, "random item not found")
	}
	rsp.Data = &content_proto.ContentArrData{contents}
	return nil
}

func (p *ContentService) GetAllSharedContents(ctx context.Context, req *content_proto.GetAllSharedContentsRequest, rsp *content_proto.GetAllSharedContentsResponse) error {
	log.Info("Received Content.GetAllSharedContents request")

	sharedContents, err := db.GetAllSharedContents(ctx, req.UserId, req.OrgId, req.Offset, req.Limit)
	if len(sharedContents) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.GetAllSharedContents, err, "shared content not found")
	}
	rsp.Data = &content_proto.GetAllSharedContentsResponse_Data{sharedContents}
	return nil
}

func (p *ContentService) GetContentRecommendations(ctx context.Context, req *content_proto.GetContentRecommendationsRequest, rsp *content_proto.GetContentRecommendationsResponse) error {
	log.Info("Received Content.GetContentRecommendations request")

	recommendations, err := db.GetContentRecommendations(ctx, req.UserId, req.OrgId)
	if len(recommendations) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.GetContentRecommendations, err, "recommendation not found")
	}
	rsp.Data = &content_proto.GetContentRecommendationsResponse_Data{recommendations}
	return nil
}

func (p *ContentService) GetContentFiltersByPreference(ctx context.Context, req *content_proto.GetContentFiltersByPreferenceRequest, rsp *content_proto.GetContentFiltersByPreferenceResponse) error {
	log.Info("Received Content.GetContentFiltersByPreference request")

	contentCategoryItems, err := db.GetContentFiltersByPreference(ctx, req.UserId, req.OrgId)
	if len(contentCategoryItems) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.GetContentFiltersByPreference, err, "content category item not found")
	}
	rsp.Data = &content_proto.GetContentFiltersByPreferenceResponse_Data{contentCategoryItems}
	return nil
}

func (p *ContentService) FilterContentRecommendations(ctx context.Context, req *content_proto.FilterContentRecommendationsRequest, rsp *content_proto.FilterContentRecommendationsResponse) error {
	log.Info("Received Content.FilterContentRecommendations request")

	contentRecommendations, err := db.FilterContentRecommendations(ctx, req.UserId, req.OrgId, req.Items)
	if len(contentRecommendations) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.FilterContentRecommendations, err, "content recommendation not found")
	}
	rsp.Data = &content_proto.FilterContentRecommendationsResponse_Data{contentRecommendations}
	return nil
}

func (p *ContentService) GetTopTags(ctx context.Context, req *content_proto.GetTopTagsRequest, rsp *content_proto.GetTopTagsResponse) error {
	log.Info("Received Content.GetTopTags request")

	rsp_tags, err := p.KvClient.GetTopTags(ctx, &kv_proto.GetTopTagsRequest{
		Index:  common.CLOUD_TAGS_INDEX,
		N:      req.N,
		OrgId:  req.OrgId,
		Object: common.CONTENT,
	})
	if err != nil {
		return common.NotFound(common.ContentSrv, p.GetTopTags, err, "top tags not found")
	}
	rsp.Data = &content_proto.GetTopTagsResponse_Data{rsp_tags.Tags}
	return nil
}

func (p *ContentService) AutocompleteTags(ctx context.Context, req *content_proto.AutocompleteTagsRequest, rsp *content_proto.AutocompleteTagsResponse) error {
	log.Info("Received Content.AutocompleteTags request")

	tags, err := db.AutocompleteTags(ctx, req.OrgId, req.Name)
	if len(tags) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.AutocompleteTags, err, "tags not found")
	}
	rsp.Data = &content_proto.AutocompleteTagsResponse_Data{tags}
	return nil
}

func (p *ContentService) WarmupCacheContent(ctx context.Context, req *content_proto.WarmupCacheContentRequest, rsp *content_proto.WarmupCacheContentResponse) error {
	log.Info("Received Content.WarmupCacheContent request")

	var offset int64
	var limit int64
	offset = 0
	limit = 100

	for {
		items, err := db.AllContents(ctx, "", "", offset, limit, "", "")
		if err != nil || len(items) == 0 {
			break
		}
		for _, item := range items {
			if len(item.Tags) > 0 {
				tags := []string{}
				for _, tag := range item.Tags {
					if tag == nil {
						continue
					}
					tags = append(tags, tag.Name)
				}
				if len(tags) == 0 {
					continue
				}
				if _, err := p.KvClient.TagsCloud(ctx, &kv_proto.TagsCloudRequest{
					Index:  common.CLOUD_TAGS_INDEX,
					OrgId:  item.OrgId,
					Object: common.CONTENT,
					Tags:   tags,
				}); err != nil {
					log.Error("warmup cache err:", err)
				}
			}
		}
		offset += limit
	}

	return nil
}

func (p *ContentService) AutocompleteContentCategoryItem(ctx context.Context, req *content_proto.AutocompleteContentCategoryItemRequest, rsp *content_proto.AutocompleteContentCategoryItemResponse) error {
	log.Info("Received Content.AutocompleteContentCategoryItem request")

	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslug(ctx, &static_proto.ReadByNameslugRequest{NameSlug: req.NameSlug})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.AutocompleteContentCategoryItem, err, "category not found")
	}
	response, err := db.AutocompleteContentCategoryItem(ctx, rsp_category.Data.ContentCategory.Id, req.Name)
	if len(response) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.AutocompleteContentCategoryItem, err, "category item not found")
	}
	rsp.Data = &content_proto.AutocompleteContentCategoryItemResponse_Data{response}
	return nil
}

func (p *ContentService) AllContentCategoryItemByNameslug(ctx context.Context, req *content_proto.AllContentCategoryItemByNameslugRequest, rsp *content_proto.AllContentCategoryItemByNameslugResponse) error {
	log.Info("Received Content.AllContentCategoryItemByNameslug request")

	rsp_category, err := p.StaticClient.ReadContentCategoryByNameslug(ctx, &static_proto.ReadByNameslugRequest{NameSlug: req.NameSlug})
	if rsp_category == nil || err != nil {
		return common.NotFound(common.ContentSrv, p.AllContentCategoryItemByNameslug, err, "category not found")
	}
	response, err := db.AllContentCategoryItemByNameslug(ctx, rsp_category.Data.ContentCategory.Id)
	if len(response) == 0 || err != nil {
		return common.NotFound(common.ContentSrv, p.AllContentCategoryItemByNameslug, err, "category item not found")
	}
	rsp.Data = &content_proto.AllContentCategoryItemByNameslugResponse_Data{response}
	return nil
}

//FIXME: this is repeated in behaviour, plan, content and survey - combine to single function somewhere? Not sure where
func (p *ContentService) sendShareNotification(userids []string, message string, alert *pubsub_proto.Alert, data map[string]string) error {
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

func (p *ContentService) GetShareableContents(ctx context.Context, req *user_proto.GetShareableContentRequest, rsp *user_proto.GetShareableContentResponse) error {
	log.Info("Received Content.GetShareableContents request")
	response, err := db.GetShareableContents(ctx, req.CreatedBy, req.Type, req.Query, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, "", "")
	if err != nil {
		return err
	}
	rsp.Data = &user_proto.GetShareableContentResponse_Data{response}
	return nil
}

func (p *ContentService) UpdateSource(ctx context.Context, req *content_proto.UpdateSourceRequest, rsp *content_proto.UpdateSourceResponse) error {
	return nil
}

func (p *ContentService) UpdateTaxonomy(ctx context.Context, req *content_proto.UpdateTaxonomyRequest, rsp *content_proto.UpdateTaxonomyResponse) error {
	return nil
}

func (p *ContentService) UpdateContentCategoryItem(ctx context.Context, req *content_proto.UpdateContentCategoryItemRequest, rsp *content_proto.UpdateContentCategoryItemResponse) error {
	return nil
}

func (p *ContentService) UpdateContent(ctx context.Context, req *content_proto.UpdateContentRequest, rsp *content_proto.UpdateContentResponse) error {
	return nil
}

func (p *ContentService) UpdateContentRule(ctx context.Context, req *content_proto.UpdateContentRuleRequest, rsp *content_proto.UpdateContentRuleResponse) error {
	return nil
}

func (p *ContentService) ReadTaxonomyByNameslug(ctx context.Context, req *static_proto.ReadByNameslugRequest, rsp *content_proto.ReadTaxonomyResponse) error {
	log.Info("Received Content.ReadTaxonomyByNameslug request")
	taxonomy, err := db.ReadTaxonomyByNameslug(ctx, req.NameSlug)
	if taxonomy == nil || err != nil {
		return common.NotFound(common.StaticSrv, p.ReadTaxonomyByNameslug, err, "not found")
	}
	rsp.Data = &content_proto.TaxonomyData{taxonomy}
	return nil
}
