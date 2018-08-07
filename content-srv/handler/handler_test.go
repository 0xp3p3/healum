package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	account_proto "server/account-srv/proto/account"
	"server/common"
	"server/content-srv/db"
	content_proto "server/content-srv/proto/content"
	kv_proto "server/kv-srv/proto/kv"
	organisation_proto "server/organisation-srv/proto/organisation"
	pubsub_proto "server/static-srv/proto/pubsub"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
	user_db "server/user-srv/db"
	user_proto "server/user-srv/proto/user"
	"strconv"
	"testing"
	"time"

	"github.com/micro/go-micro/broker"
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

var source = &content_proto.Source{
	Name:                "title",
	Description:         "description",
	IconSlug:            "iconSlug",
	Url:                 "url",
	IconUrl:             "iconUrl",
	AttributionRequired: true,
	Type: &static_proto.ContentSourceType{
		Name: "contentSourceType",
	},
	OrgId: "orgid",
	Tags:  []string{"tag1", "tag2"},
}

var taxonomy = &static_proto.Taxonomy{
	Name:        "title",
	Description: "description",
	ShortName:   "shortName",
	OrgId:       "orgid",
	Tags:        []string{"tag1", "tag2"},
	Weight:      100,
	Priority:    1,
}

var contentCategoryItem = &static_proto.ContentCategoryItem{
	Name:        "title",
	NameSlug:    "nameSlug",
	IconSlug:    "iconSlug",
	Summary:     "summary",
	Description: "description",
	OrgId:       "orgid",
	Tags:        []string{"tag1", "tag2"},
	Taxonomy:    taxonomy,
	Weight:      100,
	Priority:    1,
	Category: &static_proto.ContentCategory{
		Name:     "sample",
		NameSlug: "sample_slug",
	},
}

var content = &content_proto.Content{
	Id:          "111",
	Title:       "title",
	Summary:     []string{"summary1"},
	Description: "description",
	OrgId:       "orgid",
	CreatedBy:   &user_proto.User{Firstname: "john", Lastname: "doe"},
	Url:         "url",
	Author:      "author",
	Timestamp:   12345678,
	Tags:        []*static_proto.ContentCategoryItem{},
	Type: &static_proto.ContentType{
		Name: "content_type_sample",
	},
	Source: source,
	Category: &static_proto.ContentCategory{
		Name:     "content_category",
		NameSlug: "content_category_slug",
		IconSlug: "icon_slug",
	},
	Setting: &static_proto.Setting{
		Visibility: static_proto.Visibility_PUBLIC,
	},
}

var recipeContent = &content_proto.Content{
	Id:          "111",
	Title:       "title",
	Summary:     []string{"summary1"},
	Description: "description",
	OrgId:       "orgid",
	CreatedBy:   &user_proto.User{Firstname: "john", Lastname: "doe"},
	Url:         "url",
	Author:      "author",
	Timestamp:   12345678,
	Tags:        []*static_proto.ContentCategoryItem{contentCategoryItem},
	Type: &static_proto.ContentType{
		Name: "content_type_sample",
	},
	Source: source,
	Category: &static_proto.ContentCategory{
		Id:      "111",
		Name:    "recipe",
		Summary: "this is recipe",
	},
}

var activityContent = &content_proto.Content{
	Id:          "111",
	Title:       "title",
	Summary:     []string{"summary1"},
	Description: "description",
	OrgId:       "orgid",
	CreatedBy:   &user_proto.User{Firstname: "john", Lastname: "doe"},
	Url:         "url",
	Author:      "author",
	Timestamp:   12345678,
	Tags:        []*static_proto.ContentCategoryItem{contentCategoryItem},
	Type: &static_proto.ContentType{
		Name: "content_type_sample",
	},
	Source: source,
	Category: &static_proto.ContentCategory{
		Id:      "111",
		Name:    "activity",
		Summary: "this is activity",
	},
}

var contentRule = &content_proto.ContentRule{
	Id:     "111",
	OrgId:  "orgid",
	Type:   content_proto.RuleType_EXCLUDE,
	Source: source,
	SourceType: &static_proto.ContentSourceType{
		Name: "sample_source_type",
	},
	ContentType: &static_proto.ContentType{
		Name: "sample_content_type",
	},
	ParentCategory: &static_proto.ContentParentCategory{
		Name: "sample_parent_category",
	},
	Category: &static_proto.ContentCategory{
		Name: "sample_content_category",
	},
	CategoryItems: []*static_proto.ContentCategoryItem{
		{
			Name:     "sample_content_category_item",
			NameSlug: "name_slug",
		},
	},
	Expression: &content_proto.Expression{
		Operator: content_proto.Operator_GTE,
		Min:      5,
		Max:      10,
	},
}

var org = &organisation_proto.Organisation{
	Type: organisation_proto.OrganisationType_NONE,
}

var account = &account_proto.Account{
	Email:    "email" + common.Random(4) + "@email.com",
	Password: "pass1",
}

var user = &user_proto.User{
	OrgId:     "orgid",
	Firstname: "david",
	Lastname:  "john",
	AvatarUrl: "http://example.com",
	Tokens: []*user_proto.Token{
		{"11671c2e7da30e3c393813f60b327f9c2e2e08390761aa01e37ba5d3e6a617be", 1, "aaa"}, {"token_b", 2, "bbb"},
	},
}

var user1 = &user_proto.User{
	OrgId:     "orgid",
	Firstname: "david",
	Lastname:  "john",
	AvatarUrl: "http://example.com",
	Tokens: []*user_proto.Token{
		{"11671c2e7da30e3c393813f60b327f9c2e2e08390761aa01e37ba5d3e6a617be", 1, "aaa"}, {"token_b", 2, "bbb"},
	},
}

var preference = &user_proto.Preferences{
	OrgId:  "orgid",
	UserId: "userid",
	CurrentMeasurements: []*user_proto.Measurement{
		{
			Id:     "measure_id",
			UserId: "userid",
			OrgId:  "orgid",
		},
	},
	Conditions: []*static_proto.ContentCategoryItem{
		{Name: "name_1", NameSlug: "name_slug_1"},
	},
	Allergies: []*static_proto.ContentCategoryItem{
		{Name: "name_2", NameSlug: "name_slug_2"},
	},
	Food: []*static_proto.ContentCategoryItem{
		{Name: "name_3", NameSlug: "name_slug_3"},
	},
	Cuisines: []*static_proto.ContentCategoryItem{
		{Name: "name_4", NameSlug: "name_slug_4"},
	},
	Ethinicties: []*static_proto.ContentCategoryItem{
		{Name: "name_5", NameSlug: "name_slug_5"},
	},
}

func initHandler() *ContentService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	hdlr := &ContentService{
		Broker:        nats_brker,
		StaticClient:  static_proto.NewStaticServiceClient("go.micro.srv.static", cl),
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		TeamClient:    team_proto.NewTeamServiceClient("go.micro.srv.team", cl),
	}
	return hdlr
}

func createSource(ctx context.Context, hdlr *ContentService, t *testing.T) *content_proto.Source {
	// create content source type
	rsp_static, err := hdlr.StaticClient.CreateContentSourceType(ctx, &static_proto.CreateContentSourceTypeRequest{ContentSourceType: source.Type})
	if err != nil {
		t.Error("Create ContentSourceType err:", err)
		return nil
	}
	source.Type = rsp_static.Data.ContentSourceType

	req_create := &content_proto.CreateSourceRequest{Source: source}
	resp_create := &content_proto.CreateSourceResponse{}
	if err := hdlr.CreateSource(ctx, req_create, resp_create); err != nil {
		t.Error("Create source err:", err)
		return nil
	}
	return resp_create.Data.Source
}

func createTaxonomy(ctx context.Context, hdlr *ContentService, t *testing.T) *static_proto.Taxonomy {
	req_create := &content_proto.CreateTaxonomyRequest{Taxonomy: taxonomy}
	resp_create := &content_proto.CreateTaxonomyResponse{}
	err := hdlr.CreateTaxonomy(ctx, req_create, resp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.Taxonomy
}

func createContentCategoryItem(ctx context.Context, hdlr *ContentService, t *testing.T) *static_proto.ContentCategoryItem {
	// create taxonomy
	tax := createTaxonomy(ctx, hdlr, t)
	if tax == nil {
		return nil
	}
	contentCategoryItem.Taxonomy = tax

	// create contentCategory
	rsp_static, err := hdlr.StaticClient.CreateContentCategory(ctx, &static_proto.CreateContentCategoryRequest{ContentCategory: contentCategoryItem.Category})
	if err != nil {
		t.Error("Create ContentCategory err:", err)
		return nil
	}
	contentCategoryItem.Category = rsp_static.Data.ContentCategory

	req_create := &content_proto.CreateContentCategoryItemRequest{ContentCategoryItem: contentCategoryItem}
	resp_create := &content_proto.CreateContentCategoryItemResponse{}
	if err := hdlr.CreateContentCategoryItem(ctx, req_create, resp_create); err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.ContentCategoryItem
}

func createContent(ctx context.Context, hdlr *ContentService, t *testing.T) *content_proto.Content {
	// login user
	rsp_login, err := hdlr.AccountClient.Login(ctx, &account_proto.LoginRequest{
		Email:    "email8@email.com",
		Password: "pass1",
	})
	if err != nil {
		t.Error("Login is failed")
		return nil
	}
	rsp_kv, err := hdlr.KvClient.ReadSession(ctx, &kv_proto.ReadSessionRequest{common.SESSION_INDEX, rsp_login.Data.Session.Id})
	if err != nil {
		return nil
	}
	si := &account_proto.SessionInfo{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&si); err != nil {
		return nil
	}

	content.CreatedBy = &user_proto.User{Id: si.UserId}
	// save contentCategoryItem
	tag := createContentCategoryItem(ctx, hdlr, t)
	content.Tags = []*static_proto.ContentCategoryItem{tag}
	// save contentType
	rsp_type, err := hdlr.StaticClient.CreateContentType(ctx, &static_proto.CreateContentTypeRequest{ContentType: content.Type})
	if err != nil {
		t.Error(err)
		return nil
	}
	content.Type = rsp_type.Data.ContentType
	// save save
	src := createSource(ctx, hdlr, t)
	if src == nil {
		return nil
	}
	content.Source = src
	// save contentCategoryId
	rsp_category, err := hdlr.StaticClient.CreateContentCategory(ctx, &static_proto.CreateContentCategoryRequest{ContentCategory: content.Category})
	if err != nil {
		t.Error(err)
		return nil
	}
	content.Category = rsp_category.Data.ContentCategory

	req_create := &content_proto.CreateContentRequest{
		Content: content,
		OrgId:   si.OrgId,
		UserId:  si.UserId,
	}
	resp_create := &content_proto.CreateContentResponse{}
	if err := hdlr.CreateContent(ctx, req_create, resp_create); err != nil {
		t.Error(err)
		return nil
	}

	return resp_create.Data.Content
}

func createContentRule(ctx context.Context, hdlr *ContentService, t *testing.T) *content_proto.ContentRule {
	// create source
	src := createSource(ctx, hdlr, t)
	if src == nil {
		return nil
	}
	contentRule.Source = src
	// create contentSourceTpe
	rsp_sourcetype, err := hdlr.StaticClient.CreateContentSourceType(ctx, &static_proto.CreateContentSourceTypeRequest{ContentSourceType: contentRule.SourceType})
	if err != nil {
		t.Error("Create ContentSourceType err:", err)
		return nil
	}
	contentRule.SourceType = rsp_sourcetype.Data.ContentSourceType
	// create contentType
	rsp_contenttype, err := hdlr.StaticClient.CreateContentType(ctx, &static_proto.CreateContentTypeRequest{ContentType: contentRule.ContentType})
	if err != nil {
		t.Error(err)
		return nil
	}
	contentRule.ContentType = rsp_contenttype.Data.ContentType
	// create contentParentCategory
	rsp_parentcategory, err := hdlr.StaticClient.CreateContentParentCategory(ctx, &static_proto.CreateContentParentCategoryRequest{ContentParentCategory: contentRule.ParentCategory})
	if err != nil {
		t.Error(err)
		return nil
	}
	contentRule.ParentCategory = rsp_parentcategory.Data.ContentParentCategory
	// create contentCategory
	rsp_contentcategory, err := hdlr.StaticClient.CreateContentCategory(ctx, &static_proto.CreateContentCategoryRequest{ContentCategory: contentRule.Category})
	if err != nil {
		t.Error("Create ContentCategory err:", err)
		return nil
	}
	contentRule.Category = rsp_contentcategory.Data.ContentCategory
	// create contentCategoryItems
	categoryitem := createContentCategoryItem(ctx, hdlr, t)
	if categoryitem == nil {
		return nil
	}
	contentRule.CategoryItems[0] = categoryitem

	req_create := &content_proto.CreateContentRuleRequest{ContentRule: contentRule}
	resp_create := &content_proto.CreateContentRuleResponse{}
	if err := hdlr.CreateContentRule(ctx, req_create, resp_create); err != nil {
		t.Error(err)
		return nil
	}
	return resp_create.Data.ContentRule
}

func TestAllSources(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	src := createSource(ctx, hdlr, t)
	if src == nil {
		return
	}

	req_all := &content_proto.AllSourcesRequest{}
	resp_all := &content_proto.AllSourcesResponse{}
	err := hdlr.AllSources(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}

	if len(resp_all.Data.Sources) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Sources[0].Id != src.Id {
		t.Error("Id does not match")
		return
	}

	t.Log(resp_all.Data.Sources)
}

func TestReadSource(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	src := createSource(ctx, hdlr, t)
	if src == nil {
		return
	}

	req_read := &content_proto.ReadSourceRequest{Id: src.Id}
	resp_read := &content_proto.ReadSourceResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadSource(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Source == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Source.Id != src.Id {
		t.Error("Id does not match")
		return
	}

	t.Log(resp_read.Data.Source)
}

func TestDeleteSource(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	src := createSource(ctx, hdlr, t)
	if src == nil {
		return
	}

	req_del := &content_proto.DeleteSourceRequest{Id: src.Id}
	resp_del := &content_proto.DeleteSourceResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteSource(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
		return
	}

	req_read := &content_proto.ReadSourceRequest{Id: src.Id}
	resp_read := &content_proto.ReadSourceResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadSource(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}

}

func TestAllTaxonomys(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	tax := createTaxonomy(ctx, hdlr, t)

	req_all := &content_proto.AllTaxonomysRequest{}
	resp_all := &content_proto.AllTaxonomysResponse{}
	err := hdlr.AllTaxonomys(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}

	if len(resp_all.Data.Taxonomys) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Taxonomys[0].Id != tax.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_all.Data.Taxonomys)
}

func TestReadTaxonomy(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	tax := createTaxonomy(ctx, hdlr, t)

	req_read := &content_proto.ReadTaxonomyRequest{Id: tax.Id}
	resp_read := &content_proto.ReadTaxonomyResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadTaxonomy(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Taxonomy == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Taxonomy.Id != tax.Id {
		t.Error("Id does not match")
		return
	}
}

func TestDeleteTaxonomy(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	tax := createTaxonomy(ctx, hdlr, t)

	req_del := &content_proto.DeleteTaxonomyRequest{Id: tax.Id}
	resp_del := &content_proto.DeleteTaxonomyResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteTaxonomy(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
		return
	}

	req_read := &content_proto.ReadTaxonomyRequest{Id: tax.Id}
	resp_read := &content_proto.ReadTaxonomyResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadTaxonomy(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllContentCategoryItems(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	categoryItem := createContentCategoryItem(ctx, hdlr, t)
	if categoryItem == nil {
		return
	}

	req_all := &content_proto.AllContentCategoryItemsRequest{}
	resp_all := &content_proto.AllContentCategoryItemsResponse{}
	err := hdlr.AllContentCategoryItems(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.ContentCategoryItems) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.ContentCategoryItems[0].Id != categoryItem.Id {
		t.Error("Id does not match")
		return
	}

	t.Log(resp_all.Data.ContentCategoryItems)
}

func TestReadContentCategoryItem(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	categoryItem := createContentCategoryItem(ctx, hdlr, t)
	if categoryItem == nil {
		return
	}

	req_read := &content_proto.ReadContentCategoryItemRequest{Id: categoryItem.Id}
	resp_read := &content_proto.ReadContentCategoryItemResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadContentCategoryItem(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.ContentCategoryItem == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.ContentCategoryItem.Id != categoryItem.Id {
		t.Error("Id does not match")
		return
	}

	t.Log(resp_read.Data.ContentCategoryItem)
}

func TestDeleteContentCategoryItem(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	categoryItem := createContentCategoryItem(ctx, hdlr, t)
	if categoryItem == nil {
		return
	}

	req_del := &content_proto.DeleteContentCategoryItemRequest{Id: categoryItem.Id}
	resp_del := &content_proto.DeleteContentCategoryItemResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteContentCategoryItem(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
		return
	}

	req_read := &content_proto.ReadContentCategoryItemRequest{Id: categoryItem.Id}
	resp_read := &content_proto.ReadContentCategoryItemResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadContentCategoryItem(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllContents(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	c := createContent(ctx, hdlr, t)
	if c == nil {
		return
	}
	req_all := &content_proto.AllContentsRequest{}
	resp_all := &content_proto.AllContentsResponse{}
	err := hdlr.AllContents(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}

	if len(resp_all.Data.Contents) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Contents[0].Id != c.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_all.Data.Contents)
}

func TestReadContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	c := createContent(ctx, hdlr, t)
	if c == nil {
		return
	}

	req_read := &content_proto.ReadContentRequest{Id: c.Id}
	resp_read := &content_proto.ReadContentResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadContent(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.Content == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Content.Id != c.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_read.Data.Content)
}

func TestDeleteContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	c := createContent(ctx, hdlr, t)

	req_del := &content_proto.DeleteContentRequest{Id: c.Id}
	resp_del := &content_proto.DeleteContentResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteContent(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
		return
	}

	req_read := &content_proto.ReadContentRequest{Id: c.Id}
	resp_read := &content_proto.ReadContentResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadContent(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
}

func TestAllContentRules(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	rule := createContentRule(ctx, hdlr, t)
	if rule == nil {
		return
	}
	req_all := &content_proto.AllContentRulesRequest{}
	resp_all := &content_proto.AllContentRulesResponse{}
	err := hdlr.AllContentRules(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}

	if len(resp_all.Data.ContentRules) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.ContentRules[0].Id != rule.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_all.Data.ContentRules)
}

func TestReadContentRule(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	rule := createContentRule(ctx, hdlr, t)

	req_read := &content_proto.ReadContentRuleRequest{Id: rule.Id}
	resp_read := &content_proto.ReadContentRuleResponse{}
	time.Sleep(time.Second)
	res_read := hdlr.ReadContentRule(ctx, req_read, resp_read)
	if res_read != nil {
		t.Error(res_read)
		return
	}
	if resp_read.Data.ContentRule == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.ContentRule.Id != rule.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_read.Data.ContentRule)
}

func TestDeleteContentRule(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	rule := createContentRule(ctx, hdlr, t)

	req_del := &content_proto.DeleteContentRuleRequest{Id: rule.Id}
	resp_del := &content_proto.DeleteContentRuleResponse{}
	time.Sleep(time.Second)
	err := hdlr.DeleteContentRule(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
		return
	}

	req_read := &content_proto.ReadContentRuleRequest{Id: rule.Id}
	resp_read := &content_proto.ReadContentRuleResponse{}
	time.Sleep(time.Second)
	err = hdlr.ReadContentRule(ctx, req_read, resp_read)
	if err == nil {
		t.Error(err)
		return
	}
	if resp_read.Data != nil {
		t.Error("Object could not be nil")
		return
	}
	t.Log(resp_read.Data)
}

func TestFilterContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	c := createContent(ctx, hdlr, t)
	if c == nil {
		return
	}
	req_filter := &content_proto.FilterContentRequest{
		Sources: []string{c.Source.Id},
	}
	resp_filter := &content_proto.FilterContentResponse{}
	err := hdlr.FilterContent(ctx, req_filter, resp_filter)
	if err != nil {
		t.Error(err)
		return
	}

	if len(resp_filter.Data.Contents) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_filter.Data.Contents[0].Id != c.Id {
		t.Error("Id does not match")
		return
	}
	t.Log(resp_filter.Data.Contents)
}

func TestFilterRecipeContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	content = recipeContent
	c := createContent(ctx, hdlr, t)
	if c == nil {
		return
	}

	req_filter := &content_proto.FilterContentRequest{
		ContentCategories: []string{"recipe"},
	}
	resp_filter := &content_proto.FilterContentResponse{}
	err := hdlr.FilterContent(ctx, req_filter, resp_filter)
	if err != nil {
		t.Error(err)
		return
	}

	if len(resp_filter.Data.Contents) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_filter.Data.Contents[0].Id != c.Id {
		t.Error("Id does not match")
		return
	}
}

func TestFilterActivityContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	content = activityContent
	c := createContent(ctx, hdlr, t)

	req_filter := &content_proto.FilterContentRequest{
		ContentCategories: []string{"activity"},
	}
	resp_filter := &content_proto.FilterContentResponse{}
	err := hdlr.FilterContent(ctx, req_filter, resp_filter)
	if err != nil {
		t.Error(err)
		return
	}

	if len(resp_filter.Data.Contents) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_filter.Data.Contents[0].Id != c.Id {
		t.Error("Id does not match")
		return
	}
}

func TestShareContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	c := createContent(ctx, hdlr, t)
	if c == nil {
		return
	}

	// login user
	rsp_login, err := hdlr.AccountClient.Login(ctx, &account_proto.LoginRequest{
		Email:    "email8@email.com",
		Password: "pass1",
	})
	if err != nil {
		t.Error("Login is failed")
		return
	}
	rsp_kv, err := hdlr.KvClient.ReadSession(ctx, &kv_proto.ReadSessionRequest{common.SESSION_INDEX, rsp_login.Data.Session.Id})
	if err != nil {
		return
	}
	si := &account_proto.SessionInfo{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&si); err != nil {
		return
	}

	req_share := &content_proto.ShareContentRequest{
		Contents: []*content_proto.Content{c},
		Users:    []*user_proto.User{&user_proto.User{Id: si.UserId}},
		UserId:   si.UserId,
		OrgId:    si.OrgId,
	}
	rsp_share := &content_proto.ShareContentResponse{}
	if err := hdlr.ShareContent(ctx, req_share, rsp_share); err != nil {
		t.Error(err)
		return
	}
}

func TestCreateActivityContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	if err := hdlr.Subscribe(ctx, &pubsub_proto.SubscribeRequest{common.CREATE_ACTIVITY_CONTENT}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	obj := []*content_proto.Activity{
		{
			Source: "activity_test",
		},
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := hdlr.Broker.Publish(common.CREATE_ACTIVITY_CONTENT, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	req_all := &content_proto.AllContentsRequest{}
	resp_all := &content_proto.AllContentsResponse{}
	err = hdlr.AllContents(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp_all.Data.Contents)
}

func TestCreateRecipeContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	if err := hdlr.Subscribe(ctx, &pubsub_proto.SubscribeRequest{common.CREATE_RECIPE_CONTENT}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	obj := []*content_proto.Recipe{
		{
			CookingMethod: "recipe_test",
		},
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := hdlr.Broker.Publish(common.CREATE_RECIPE_CONTENT, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	req_all := &content_proto.AllContentsRequest{}
	resp_all := &content_proto.AllContentsResponse{}
	err = hdlr.AllContents(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(resp_all.Data.Contents)
}

func TestCreateArticleContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	if err := hdlr.Subscribe(ctx, &pubsub_proto.SubscribeRequest{common.CREATE_ARTICLE_CONTENT}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	obj := []*content_proto.Article{
		{
			Text: "article_test",
		},
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := hdlr.Broker.Publish(common.CREATE_ARTICLE_CONTENT, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	req_all := &content_proto.AllContentsRequest{}
	resp_all := &content_proto.AllContentsResponse{}
	err = hdlr.AllContents(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(resp_all.Data.Contents)
}

func TestCreatePlaceContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	if err := hdlr.Subscribe(ctx, &pubsub_proto.SubscribeRequest{common.CREATE_PLACE_CONTENT}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	obj := []*content_proto.Place{
		{
			Type: "place_test",
		},
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := hdlr.Broker.Publish(common.CREATE_PLACE_CONTENT, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	req_all := &content_proto.AllContentsRequest{}
	resp_all := &content_proto.AllContentsResponse{}
	err = hdlr.AllContents(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	// if resp_all.Data.Contents[0].Category.NameSlug != "place" {
	// 	t.Error("Object could not be nil")
	// 	return
	// }
}

func TestCreateWellbeingContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	if err := hdlr.Subscribe(ctx, &pubsub_proto.SubscribeRequest{common.CREATE_WELLBEING_CONTENT}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	obj := []*content_proto.Wellbeing{
		{
			Source: "wellbeing_test",
		},
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := hdlr.Broker.Publish(common.CREATE_WELLBEING_CONTENT, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	req_all := &content_proto.AllContentsRequest{}
	resp_all := &content_proto.AllContentsResponse{}
	err = hdlr.AllContents(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Contents) == 0 {
		t.Error("Object count does not matched")
		return
	}
	// if resp_all.Data.Contents[0].Category.NameSlug != "wellbeing" {
	// 	t.Error("Object could not be nil")
	// 	return
	// }
}

func TestCreateContentRecommendation(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create content
	c := createContent(ctx, hdlr, t)

	if err := hdlr.Subscribe(ctx, &pubsub_proto.SubscribeRequest{common.CREATE_CONTENT_RECOMMENDATION}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
		return
	}

	obj := &content_proto.ContentRecommendation{
		OrgId:   "orgid",
		UserId:  "userid",
		Content: c,
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := hdlr.Broker.Publish(common.CREATE_CONTENT_RECOMMENDATION, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	// get recommendation by user
	req_get := &content_proto.GetContentRecommendationByUserRequest{
		UserId: "userid",
	}
	resp_get := &content_proto.GetContentRecommendationByUserResponse{}
	err = hdlr.GetContentRecommendationByUser(ctx, req_get, resp_get)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_get.Data.Recommendations) == 0 {
		t.Error("Object count does not matched")
		return
	}
	if resp_get.Data.Recommendations[0].ContentId != c.Id {
		t.Error("Object does not matched")
		return
	}

	// get recommend with category_id
	req_category := &content_proto.GetContentRecommendationByCategoryRequest{
		CategoryId: c.Category.Id,
	}
	rsp_category := &content_proto.GetContentRecommendationByCategoryResponse{}
	err = hdlr.GetContentRecommendationByCategory(ctx, req_category, rsp_category)
	if err != nil {
		t.Error(err)
		return
	}
	if len(rsp_category.Data.Recommendations) == 0 {
		t.Error("Object count does not matched")
		return
	}
	if rsp_category.Data.Recommendations[0].ContentId != c.Id {
		t.Error("Object does not matched")
		return
	}

}

func TestGetRandomItems(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// login user
	rsp_login, err := hdlr.AccountClient.Login(ctx, &account_proto.LoginRequest{
		Email:    "email8@email.com",
		Password: "pass1",
	})
	if err != nil {
		t.Error("Login is failed")
		return
	}
	rsp_kv, err := hdlr.KvClient.ReadSession(ctx, &kv_proto.ReadSessionRequest{common.SESSION_INDEX, rsp_login.Data.Session.Id})
	if err != nil {
		return
	}
	si := &account_proto.SessionInfo{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&si); err != nil {
		return
	}

	// count := 10
	for i := 0; i < 20; i++ {
		content.Id = "id_" + strconv.Itoa(i)
		content.Title = "title_" + strconv.Itoa(i)
		content.Category.Id = "category_" + strconv.Itoa(i/2)
		req := &content_proto.CreateContentRequest{
			Content: content,
			UserId:  si.UserId,
			OrgId:   si.OrgId,
			TeamId:  si.UserId,
		}
		rsp := &content_proto.CreateContentResponse{}
		if err := hdlr.CreateContent(ctx, req, rsp); err != nil {
			t.Error(err)
			return
		}
		if rsp.Data.Content == nil {
			t.Error("create create fail.")
		}
	}

	req_get := &content_proto.GetRandomItemsRequest{9}
	rsp_get := &content_proto.GetRandomItemsResponse{}
	err = hdlr.GetRandomItems(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}

	if len(rsp_get.Data.Contents) != 9 {
		t.Error("Object count does not matched")
		return
	}
	t.Log(rsp_get.Data)
}

func TestCreateVideoContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	if err := hdlr.Subscribe(ctx, &pubsub_proto.SubscribeRequest{common.CREATE_VIDEO_CONTENT}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	obj := []*content_proto.Video{
		{
			ContentUrl: "video_test",
		},
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := hdlr.Broker.Publish(common.CREATE_VIDEO_CONTENT, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	req_all := &content_proto.AllContentsRequest{}
	resp_all := &content_proto.AllContentsResponse{}
	err = hdlr.AllContents(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Contents) == 0 {
		t.Error("Object count does not matched")
		return
	}
	// if resp_all.Data.Contents[0].Category.NameSlug != "video" {
	// 	t.Error("Object could not be nil")
	// 	return
	// }
}

func TestCreateProductContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	if err := hdlr.Subscribe(ctx, &pubsub_proto.SubscribeRequest{common.CREATE_PRODUCT_CONTENT}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	obj := []*content_proto.Product{
		{
			Name: "product_test",
		},
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := hdlr.Broker.Publish(common.CREATE_PRODUCT_CONTENT, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	req_all := &content_proto.AllContentsRequest{}
	resp_all := &content_proto.AllContentsResponse{}
	err = hdlr.AllContents(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Contents) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestCreateServiceContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	if err := hdlr.Subscribe(ctx, &pubsub_proto.SubscribeRequest{common.CREATE_SERVICE_CONTENT}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	obj := []*content_proto.Service{
		{
			Name: "service_test",
		},
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := hdlr.Broker.Publish(common.CREATE_SERVICE_CONTENT, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(4 * time.Second)

	req_all := &content_proto.AllContentsRequest{}
	resp_all := &content_proto.AllContentsResponse{}
	err = hdlr.AllContents(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Contents) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestCreateEventContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	if err := hdlr.Subscribe(ctx, &pubsub_proto.SubscribeRequest{common.CREATE_EVENT_CONTENT}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	obj := []*content_proto.Event{
		{
			Name: "event_test",
		},
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := hdlr.Broker.Publish(common.CREATE_EVENT_CONTENT, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	req_all := &content_proto.AllContentsRequest{}
	resp_all := &content_proto.AllContentsResponse{}
	err = hdlr.AllContents(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Contents) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestCreateResearchContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	if err := hdlr.Subscribe(ctx, &pubsub_proto.SubscribeRequest{common.CREATE_RESEARCH_CONTENT}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	obj := []*content_proto.Research{
		{
			ArticleBody: "research_test",
		},
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := hdlr.Broker.Publish(common.CREATE_RESEARCH_CONTENT, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	req_all := &content_proto.AllContentsRequest{}
	resp_all := &content_proto.AllContentsResponse{}
	err = hdlr.AllContents(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Contents) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestCreateAppContent(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	if err := hdlr.Subscribe(ctx, &pubsub_proto.SubscribeRequest{common.CREATE_APP_CONTENT}, &pubsub_proto.SubscribeResponse{}); err != nil {
		log.Fatal(err)
	}

	obj := []*static_proto.App{
		{
			Name: "app_test",
		},
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}

	// publish
	if err := hdlr.Broker.Publish(common.CREATE_APP_CONTENT, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	req_all := &content_proto.AllContentsRequest{}
	resp_all := &content_proto.AllContentsResponse{}
	err = hdlr.AllContents(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Contents) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestGetAllSharedContents(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// create user
	userClient := user_proto.NewUserServiceClient("go.micro.srv.user", cl)
	account.Email = "email" + common.Random(4) + "@ex.com"
	rsp_user, err := userClient.Create(ctx, &user_proto.CreateRequest{
		User: user1, Account: account,
	})
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Second)
	// create employee
	if _, err := hdlr.StaticClient.CreateRole(ctx, &static_proto.CreateRoleRequest{
		&static_proto.Role{Name: "admin_role", NameSlug: "admin"},
	}); err != nil {
		t.Error(err)
		return
	}
	req_org := &organisation_proto.CreateRequest{
		Organisation: org,
		Account:      account,
		User:         user,
	}
	time.Sleep(time.Second)
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	rsp_org, err := orgClient.Create(ctx, req_org)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(time.Second)

	// login user
	rsp_login, err := hdlr.AccountClient.Login(ctx, &account_proto.LoginRequest{
		Email:    "email8@email.com",
		Password: "pass1",
	})
	if err != nil {
		t.Error("Login is failed")
		return
	}
	rsp_kv, err := hdlr.KvClient.ReadSession(ctx, &kv_proto.ReadSessionRequest{common.SESSION_INDEX, rsp_login.Data.Session.Id})
	if err != nil {
		return
	}
	si := &account_proto.SessionInfo{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&si); err != nil {
		return
	}

	req := &content_proto.CreateContentRequest{
		Content: content,
		UserId:  si.UserId,
		OrgId:   si.OrgId,
		TeamId:  si.UserId,
	}
	rsp := &content_proto.CreateContentResponse{}
	if err := hdlr.CreateContent(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}
	if rsp.Data.Content == nil {
		t.Error("create create fail.")
	}

	// c := createContent(ctx, hdlr, t)
	// if c == nil {
	// 	return
	// }

	if err = db.ShareContent(ctx, []*content_proto.Content{c},
		[]*user_proto.User{rsp_user.Data.User},
		rsp_org.Data.User,
		rsp_org.Data.Organisation.Id); err != nil {
		t.Error(err)
		return
	}

	req_share := &content_proto.GetAllSharedContentsRequest{
		UserId: rsp_user.Data.User.Id,
		OrgId:  rsp_org.Data.Organisation.Id,
	}
	rsp_share := &content_proto.GetAllSharedContentsResponse{}
	if err := hdlr.GetAllSharedContents(ctx, req_share, rsp_share); err != nil {
		t.Error(err)
		return
	}

	//t.Log(rsp.Data.ShareContentUsers)
}

func TestGetContentRecommendations(t *testing.T) {
	TestCreateContentRecommendation(t)

	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req := &content_proto.GetContentRecommendationsRequest{
		UserId: "userid",
		OrgId:  "orgid",
	}
	rsp := &content_proto.GetContentRecommendationsResponse{}
	if err := hdlr.GetContentRecommendations(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp.Data.Recommendations)
}

func TestGetContentFiltersByPreference(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	// save user preference
	user_db.Init(cl)
	if err := user_db.SaveUserPreference(ctx, preference); err != nil {
		t.Error(err)
		return
	}

	req := &content_proto.GetContentFiltersByPreferenceRequest{
		UserId: "userid",
		OrgId:  "orgid",
	}
	rsp := &content_proto.GetContentFiltersByPreferenceResponse{}

	if err := hdlr.GetContentFiltersByPreference(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp.Data.ContentCategoryItems)
}

func TestFilterContentRecommendations(t *testing.T) {
	// create content recommendation
	TestCreateContentRecommendation(t)

	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req := &content_proto.FilterContentRecommendationsRequest{
		UserId: "userid",
		OrgId:  "orgid",
		// Items: []*static_proto.ContentCategoryItem{{
		// 	Id: "category_id",
		// }},
	}
	rsp := &content_proto.FilterContentRecommendationsResponse{}

	if err := hdlr.FilterContentRecommendations(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp.Data.Response)
}

func TestGetTopTags(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	createContent(ctx, hdlr, t)

	rsp := &content_proto.GetTopTagsResponse{}
	if err := hdlr.GetTopTags(ctx, &content_proto.GetTopTagsRequest{
		OrgId: content.OrgId,
		N:     5,
	}, rsp); err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp.Data.Tags)
}

func TestAutocompleteTags(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	createContent(ctx, hdlr, t)

	rsp := &content_proto.AutocompleteTagsResponse{}
	if err := hdlr.AutocompleteTags(ctx, &content_proto.AutocompleteTagsRequest{
		OrgId: content.OrgId,
		Name:  "t",
	}, rsp); err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp.Data.Tags)
}

func TestGetContentByCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	content := createContent(ctx, hdlr, t)
	if content == nil {
		return
	}

	req := &content_proto.GetContentByCategoryRequest{CategoryId: content.Category.Id}
	rsp := &content_proto.GetContentByCategoryResponse{}
	if err := hdlr.GetContentByCategory(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}
	if rsp.Data.Contents[0].CategoryId != content.Category.Id {
		t.Error("CategoryId does not matched")
		return
	}
	t.Log(rsp.Data.Contents)
}

func TestGetFiltersForCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	item := createContentCategoryItem(ctx, hdlr, t)
	if item == nil {
		return
	}
	req := &content_proto.GetFiltersForCategoryRequest{CategoryId: item.Category.Id}
	rsp := &content_proto.GetFiltersForCategoryResponse{}
	err := hdlr.GetFiltersForCategory(ctx, req, rsp)
	if err != nil {
		t.Error(err)
		return
	}
	if len(rsp.Data.ContentCategoryItems) == 0 {
		t.Error("Object count does not matched")
		return
	}
	if rsp.Data.ContentCategoryItems[0].Category.Id != item.Category.Id {
		t.Error("Object Id does not matched")
		return
	}
	t.Log(rsp.Data.ContentCategoryItems)
}

func TestFiltersAutocomplete(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	item := createContentCategoryItem(ctx, hdlr, t)
	if item == nil {
		return
	}
	req := &content_proto.FiltersAutocompleteRequest{CategoryId: item.Category.Id, Name: item.Name[0:2]}
	rsp := &content_proto.FiltersAutocompleteResponse{}
	err := hdlr.FiltersAutocomplete(ctx, req, rsp)
	if err != nil {
		t.Error(err)
		return
	}
	if len(rsp.Data.ContentCategoryItems) == 0 {
		t.Error("Object count does not matched")
		return
	}
	if rsp.Data.ContentCategoryItems[0].Category.Id != item.Category.Id {
		t.Error("Object Id does not matched")
		return
	}
	t.Log(rsp.Data.ContentCategoryItems)
}

func TestFilterContentInParticularCategory(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	content := createContent(ctx, hdlr, t)
	if content == nil {
		return
	}
	// filter
	req := &content_proto.FilterContentInParticularCategoryRequest{
		CategoryId:           content.Category.Id,
		ContentCategoryItems: []string{content.Tags[0].Id},
	}
	rsp := &content_proto.FilterContentInParticularCategoryResponse{}
	if err := hdlr.FilterContentInParticularCategory(ctx, req, rsp); err != nil {
		t.Error(err)
		return
	}
	if len(rsp.Data.Contents) == 0 {
		t.Error("Object count does not matched")
		return
	}
	if rsp.Data.Contents[0].CategoryId != content.Category.Id {
		t.Error("Object Id does not matched")
		return
	}
	t.Log(rsp.Data.Contents)
}
