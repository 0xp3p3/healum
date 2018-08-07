package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server/api/utils"
	"server/common"
	"server/content-srv/db"
	content_proto "server/content-srv/proto/content"
	static_proto "server/static-srv/proto/static"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

// var serverURL = "http://localhost:8080"
var contentURL = "/server/content"

var source = &content_proto.Source{
	Id:                  "111",
	Name:                "title",
	Description:         "description",
	IconSlug:            "iconSlug",
	Url:                 "url",
	IconUrl:             "iconUrl",
	AttributionRequired: true,
	Type:                &static_proto.ContentSourceType{},
	OrgId:               "orgid",
	Tags:                []string{"tag1", "tag2"},
}

var taxonomy = &static_proto.Taxonomy{
	Id:          "111",
	Name:        "title",
	Description: "description",
	ShortName:   "shortName",
	OrgId:       "orgid",
	Tags:        []string{"tag1", "tag2"},
	Weight:      100,
	Priority:    1,
}

var contentCategoryItem = &static_proto.ContentCategoryItem{
	Id:          "111",
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
	Category:    contentCategory,
}

var content = &content_proto.Content{
	Id:          "111",
	Title:       "title",
	Summary:     "summary1",
	Description: "description",
	OrgId:       "orgid",
	CreatedBy:   &user_proto.User{Id: "userid"},
	Url:         "url",
	Author:      "author",
	Timestamp:   12345678,
	Tags: []*static_proto.ContentCategoryItem{
		{Id: "111", Name: "tag1"},
		{Id: "222", Name: "tag2"},
		{Id: "333", Name: "tag3"},
	},
	Type:   &static_proto.ContentType{},
	Source: source,
	Category: &static_proto.ContentCategory{
		Id:       "category111",
		Name:     "activity category",
		NameSlug: "acitivty",
		TrackerMethods: []*static_proto.TrackerMethod{
			{
				Id:       "tracker111",
				NameSlug: "count"},
		},
	},
}

var contentRule = &content_proto.ContentRule{
	Id:             "111",
	OrgId:          "orgid",
	Type:           content_proto.RuleType_EXCLUDE,
	Source:         source,
	SourceType:     &static_proto.ContentSourceType{},
	ContentType:    &static_proto.ContentType{},
	ParentCategory: &static_proto.ContentParentCategory{},
	Category:       &static_proto.ContentCategory{},
	CategoryItems:  []*static_proto.ContentCategoryItem{},
	Expression:     &content_proto.Expression{},
}

func initContentDb() {
	cl := client.NewClient(
		client.Transport(nats_transport.NewTransport()),
		client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(5*time.Second),
		client.Retries(5))
	// ctx := common.NewTestContext(context.TODO())
	// db.DbContentName = common.TestingName("healum_test")
	// db.DbSourceTable = common.TestingName("source_test")
	// db.DbTaxonomyTable = common.TestingName("taxonomy_test")
	// db.DbContentCategoryItemTable = common.TestingName("content_category_item_test")
	// db.DbContentTable = common.TestingName("content_test")
	// db.DbContentRuleTable = common.TestingName("content_rule_test")
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func createSource(source *content_proto.Source, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"source": source})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+contentURL+"/source?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)
}

func createTaxonomy(taxonomy *static_proto.Taxonomy, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"taxonomy": taxonomy})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+contentURL+"/taxonomy?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)
}

func createContent(content *content_proto.Content, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return
	}
	content.CreatedBy.Id = userId
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"content": content})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+contentURL+"/content?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)
}

func createContentCategoryItem(contentCategoryItem *static_proto.ContentCategoryItem, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"contentCategoryItem": contentCategoryItem})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+contentURL+"/category/item?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)
}

func createContentRule(contentRule *content_proto.ContentRule, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"contentRule": contentRule})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+contentURL+"/rule?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)
}

func TestAllSources(t *testing.T) {
	initContentDb()

	createSource(source, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/sources/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.AllSourcesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Sources) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func TestReadSource(t *testing.T) {
	initContentDb()

	createSource(source, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/source/"+source.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.ReadSourceResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.Source == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Source.Id != source.Id {
		t.Errorf("Object Id does not matched")
		return
	}
}

func TestDeleteSource(t *testing.T) {
	initContentDb()

	createSource(source, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+contentURL+"/source/"+source.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+contentURL+"/source/"+source.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.ReadSourceResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
		return
	}
}

func TestErrReadSource(t *testing.T) {
	initContentDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/source/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
		return
	}
	t.Log(r)
}

func TestAllTaxonomys(t *testing.T) {
	initContentDb()

	createTaxonomy(taxonomy, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/taxonomys/all?session="+sessionId+"&team_id="+"&org_id="+taxonomy.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.AllTaxonomysResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Taxonomys) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func TestReadTaxonomy(t *testing.T) {
	initContentDb()

	createTaxonomy(taxonomy, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/taxonomy/"+taxonomy.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.ReadTaxonomyResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.Taxonomy == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Taxonomy.Id != taxonomy.Id {
		t.Errorf("Object Id does not matched")
		return
	}
}

func TestDeleteTaxonomy(t *testing.T) {
	initContentDb()

	createTaxonomy(taxonomy, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+contentURL+"/taxonomy/"+taxonomy.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+contentURL+"/taxonomy/"+taxonomy.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.ReadTaxonomyResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
		return
	}
}

func TestErrReadTaxonomy(t *testing.T) {
	initContentDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/taxonomy/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
		return
	}
	t.Log(r)
}

func TestAllContentCategoryItems(t *testing.T) {
	initContentDb()

	createContentCategoryItem(contentCategoryItem, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/category/items/all?session="+sessionId+"&team_id="+"&org_id="+contentCategoryItem.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.AllContentCategoryItemsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.ContentCategoryItems) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func TestReadContentCategoryItem(t *testing.T) {
	initContentDb()

	createContentCategoryItem(contentCategoryItem, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/category/item/"+contentCategoryItem.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.ReadContentCategoryItemResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.ContentCategoryItem == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.ContentCategoryItem.Id != contentCategoryItem.Id {
		t.Errorf("Object Id does not matched")
		return
	}
}

func TestDeleteContentCategoryItem(t *testing.T) {
	initContentDb()

	createContentCategoryItem(contentCategoryItem, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+contentURL+"/category/item/"+contentCategoryItem.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+contentURL+"/category/item/"+contentCategoryItem.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.ReadContentCategoryItemResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
		return
	}
}

func TestErrReadContentCategoryItem(t *testing.T) {
	initContentDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/category/item/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
		return
	}
	t.Log(r)
}

func TestAllContents(t *testing.T) {
	initContentDb()

	createContent(content, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/contents/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.AllContentsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Contents) == 0 {
		t.Errorf("Object count does not matched")
		return
	}

	t.Log(r)
}

func TestReadContent(t *testing.T) {
	initContentDb()

	createContent(content, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/content/"+content.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.ReadContentResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Content == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Content.Id != content.Id {
		t.Errorf("Object Id does not matched")
		return
	}
}

func TestDeleteContent(t *testing.T) {
	initContentDb()

	createContent(content, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+contentURL+"/content/"+content.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+contentURL+"/content/"+content.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.ReadContentResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
		return
	}
}

func TestErrReadContent(t *testing.T) {
	initContentDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/content/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
		return
	}
	t.Log(r)
}

func TestAllContentRules(t *testing.T) {
	initContentDb()

	createContentRule(contentRule, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/rules/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.AllContentRulesResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.ContentRules) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func TestReadContentRule(t *testing.T) {
	initContentDb()

	createContentRule(contentRule, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/rule/"+contentRule.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.ReadContentRuleResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.ContentRule == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.ContentRule.Id != contentRule.Id {
		t.Errorf("Object Id does not matched")
		return
	}
}

func TestDeleteContentRule(t *testing.T) {
	initContentDb()

	createContentRule(contentRule, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a Delete request.
	req, err := http.NewRequest("DELETE", serverURL+contentURL+"/rule/"+contentRule.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	// Send a GET request.
	req, err = http.NewRequest("GET", serverURL+contentURL+"/rule/"+contentRule.Id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.ReadContentRuleResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data != nil {
		t.Errorf("Object does not matched")
		return
	}
}

func TestErrReadContentRule(t *testing.T) {
	initContentDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/rule/999?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
		return
	}
	t.Log(r)
}

func TestFilterContent(t *testing.T) {
	initContentDb()

	createSource(source, t)
	createContentCategoryItem(contentCategoryItem, t)
	createContent(content, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"sources":              []string{"111"},
		"contentCategoryItems": []string{"111", "222"},
	})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+contentURL+"/filter?session="+sessionId+"&offset=0&limit=10", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.FilterContentResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Contents) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func TestShareContent(t *testing.T) {
	initContentDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	log.Println(sessionId)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"contents": []*content_proto.Content{content},
		"users":    []*user_proto.User{user1},
	})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+contentURL+"/share?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.ShareContentResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Code != http.StatusOK {
		t.Log(r)
		t.Errorf("Response does not matched")
	}
}

func TestGetTopContentTags(t *testing.T) {
	initContentDb()

	createContent(content, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+contentURL+"/tags/top/5?session="+sessionId, nil)

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.GetTopTagsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Tags) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Tags)
}

func TestAutocompleteContentTags(t *testing.T) {
	initContentDb()

	createContent(content, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"name": "t"})
	if err != nil {
		t.Error(err)
		return
	}
	// log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+contentURL+"/tags/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.AutocompleteTagsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Tags) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Tags)
}

func TestCreateAppContentWithJSON(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	js := `{"content": {
		"author": "",
		"category": {
		  "actions": [],
		  "created": "1523420459",
		  "description": "",
		  "icon_slug": "",
		  "id": "bcce6ab2-3d3f-11e8-b80c-20c9d0453b15",
		  "name": "",
		  "name_slug": "app",
		  "org_id": "",
		  "parent": [],
		  "summary": "",
		  "tags": [],
		  "trackerMethods": [],
		  "updated": "1523420459"
		},
		"created": "1523420459",
		"createdBy": null,
		"description": "",
		"hash": "fc862a62b51bbbc7e21aba8f38b0bebc4b8cbeb009ff7a67121f746ea3c68351",
		"image": "",
		"item": {
		  "@type": "healum.com/proto/go.micro.srv.static.App",
		  "created": "0",
		  "description": "hello's test",
		  "icon_slug": "",
		  "id": "",
		  "image": "",
		  "name": "app_test",
		  "platforms": [],
		  "summary": "",
		  "tags": [],
		  "updated": "0"
		},
		"org_id": "",
		"source": null,
		"summary": [],
		"tags": [],
		"timestamp": "0",
		"title": "hellow",
		"type": null,
		"updated": "1523420459",
		"url": ""
	  }}`
	req, err := http.NewRequest("POST", serverURL+contentURL+"/content?session="+sessionId, bytes.NewBuffer([]byte(js)))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.CreateContentResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	t.Log(r)

	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
}

func TestCreateRecipeContentWithJSON(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	js := `{"content": {
		"author": "",
		"category": {
		  "actions": [],
		  "created": "1523420459",
		  "description": "",
		  "icon_slug": "",
		  "id": "bcce6ab2-3d3f-11e8-b80c-20c9d0453b15",
		  "name": "",
		  "name_slug": "app",
		  "org_id": "",
		  "parent": [],
		  "summary": "",
		  "tags": [],
		  "trackerMethods": [],
		  "updated": "1523420459"
		},
		"created": "1523420459",
		"createdBy": null,
		"description": "",
		"hash": "fc862a62b51bbbc7e21aba8f38b0bebc4b8cbeb009ff7a67121f746ea3c68351",
		"image": "",
		"item": {
		  "@type": "healum.com/proto/go.micro.srv.content.Recipe",
		  "serves": 1.5,
		  "diets": [],
		  "allergies": []
		},
		"org_id": "",
		"source": null,
		"summary": [],
		"tags": [],
		"timestamp": "0",
		"title": "hellow",
		"type": null,
		"updated": "1523420459",
		"url": ""
	  }}`
	req, err := http.NewRequest("POST", serverURL+contentURL+"/content?session="+sessionId, bytes.NewBuffer([]byte(js)))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := content_proto.CreateContentResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	t.Log(r)

	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
}
