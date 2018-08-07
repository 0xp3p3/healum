package handler

import (
	"bytes"
	"context"
	"encoding/json"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	organisation_proto "server/organisation-srv/proto/organisation"
	static_proto "server/static-srv/proto/static"
	"server/survey-srv/db"
	survey_proto "server/survey-srv/proto/survey"
	team_proto "server/team-srv/proto/team"
	user_proto "server/user-srv/proto/user"

	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
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

var survey = &survey_proto.Survey{
	// Id:          "111",
	Title:       "titles sample",
	OrgId:       "orgid",
	Tags:        []string{"tag1", "tag2"},
	Description: "description1",
	Creator:     &user_proto.User{Id: "userid"},
	IsTemplate:  true,
	Status:      survey_proto.SurveyStatus_DRAFT,
	Renders: []survey_proto.RenderTarget{
		survey_proto.RenderTarget_MOBILE, survey_proto.RenderTarget_WEB,
	},
	Setting: &static_proto.Setting{
		Visibility: static_proto.Visibility_PUBLIC,
	},
	Questions: []*survey_proto.Question{question},
	Shares:    []*user_proto.User{user1},
}

var question = &survey_proto.Question{
	Id:          "q111",
	Type:        survey_proto.QuestionType_FILE,
	Title:       "question",
	Description: "description",
}

var question_map = map[string]interface{}{
	"id":          "q111",
	"type":        survey_proto.QuestionType_FILE,
	"title":       "question",
	"description": "description",
	"fields": map[string]interface{}{
		"@type": "healum.com/proto/go.micro.srv.survey.TextQuestionField",
		"attributes": map[string]string{
			"description": "Hello world",
		},
		"validations": []survey_proto.Validation{},
		"attachment": map[string]interface{}{
			"type":   survey_proto.AttachmentType_VIDEO,
			"url":    "http://example.com/sample",
			"width":  320,
			"height": 240,
		},
	},
	"settings": map[string]interface{}{
		"@type":              "healum.com/proto/go.micro.srv.survey.TextQuestionSettings",
		"maximum_characters": 255,
		"multiple_rows":      true,
		"mandatory":          true,
		"image":              true,
		"video":              true,
		"variable":           true,
	},
}

var org1 = &organisation_proto.Organisation{
	Type: organisation_proto.OrganisationType_ROOT,
}

var user1 = &user_proto.User{
	Firstname: "david",
	Lastname:  "john",
	Tokens: []*user_proto.Token{
		{"11671c2e7da30e3c393813f60b327f9c2e2e08390761aa01e37ba5d3e6a617be", 1, "aaa"}, {"token_b", 2, "bbb"},
	},
}

var account1 = &account_proto.Account{
	Email:    "test" + common.Random(4) + "@email.com",
	Password: "pass1",
}

func initHandler() *SurveyService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	hdlr := &SurveyService{
		Broker:        nats_brker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		TeamClient:    team_proto.NewTeamServiceClient("go.micro.srv.team", cl),
	}
	return hdlr
}

func createSurvey(ctx context.Context, hdlr *SurveyService, t *testing.T) *survey_proto.Survey {
	// create org
	orgClient := organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl)
	user1.Id = ""
	account1.Email = "test" + common.Random(4) + "@email.com"
	rsp_org, err := orgClient.Create(ctx, &organisation_proto.CreateRequest{Organisation: org1, User: user1, Account: account1})
	if err != nil {
		t.Error(err)
		return nil
	}
	survey.OrgId = rsp_org.Data.Organisation.Id
	//	survey.Users  = []*user_proto.User{rsp_org.Data.User}
	survey.Shares = []*user_proto.User{rsp_org.Data.User}
	survey.Creator = rsp_org.Data.User

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

	// create survey
	req := &survey_proto.CreateRequest{
		Survey: survey,
		UserId: si.UserId,
		OrgId:  si.OrgId,
	}
	resp := &survey_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req, resp); err != nil {
		t.Error(err)
		return nil
	}
	return resp.Data.Survey
}

func TestAll(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
		return
	}

	req_all := &survey_proto.AllRequest{}
	resp_all := &survey_proto.AllResponse{}
	err := hdlr.All(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Surveys) == 0 {
		t.Error("Count does not match")
		return
	}

	t.Log(resp_all.Data.Surveys)
}

func TestNew(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	req := &survey_proto.NewRequest{}
	resp := &survey_proto.NewResponse{}
	res := hdlr.New(ctx, req, resp)
	if res != nil {
		t.Error(res)
		return
	}
}

func TestSurveyIsCreated(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	req_new := &survey_proto.NewRequest{}
	resp_new := &survey_proto.NewResponse{}
	res_new := hdlr.New(ctx, req_new, resp_new)
	if res_new != nil {
		t.Error(res_new)
	}

	survey.Id = resp_new.Data.SurveyId
	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
		return
	}
}

func TestSurveyRead(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
		return
	}

	req_read := &survey_proto.ReadRequest{Id: survey.Id}
	resp_read := &survey_proto.ReadResponse{}
	err := hdlr.Read(ctx, req_read, resp_read)
	if err != nil {
		t.Error(err)
		return
	}
	if resp_read.Data.Survey == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Survey.Id != survey.Id {
		t.Error("Id does not match")
		return
	}
}

func TestSurveyDelete(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
		return
	}

	req_del := &survey_proto.DeleteRequest{Id: survey.Id}
	resp_del := &survey_proto.DeleteResponse{}
	err := hdlr.Delete(ctx, req_del, resp_del)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSurveyCopy(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
		return
	}

	req_copy := &survey_proto.CopyRequest{SurveyId: survey.Id}
	resp_copy := &survey_proto.CopyResponse{}
	err := hdlr.Copy(ctx, req_copy, resp_copy)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestQuestions(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
		return
	}
	// questions query
	req_query := &survey_proto.QuestionsRequest{SurveyId: survey.Id}
	resp_query := &survey_proto.QuestionsResponse{}
	err := hdlr.Questions(ctx, req_query, resp_query)
	if err != nil {
		t.Error(err)
		return
	}
	if resp_query.Data.Questions[0].Id != survey.Questions[0].Id {
		t.Error("Id does not match")
		return
	}
}

func TestQuestionRef(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
		return
	}
	// question ref
	req_ref := &survey_proto.QuestionRefRequest{SurveyId: survey.Id, QuestionRef: survey.Questions[0].Id}
	resp_ref := &survey_proto.QuestionRefResponse{}
	err := hdlr.QuestionRef(ctx, req_ref, resp_ref)
	if err != nil {
		t.Error(err)
		return
	}
	if resp_ref.Data.Question.Id != survey.Questions[0].Id {
		t.Error("Id does not match")
		return
	}
}

func TestByCreator(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
		return
	}

	req_creator := &survey_proto.ByCreatorRequest{UserId: survey.Creator.Id}
	resp_creator := &survey_proto.ByCreatorResponse{}
	time.Sleep(1 * time.Second)
	err := hdlr.ByCreator(ctx, req_creator, resp_creator)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_creator.Data.Surveys) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_creator.Data.Surveys[0].Id != survey.Id {
		t.Error("Id does not match")
		return
	}
}

func TestLink(t *testing.T) {

}

func TestTempaltes(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
		return
	}

	req_temp := &survey_proto.TemplatesRequest{}
	resp_temp := &survey_proto.TemplatesResponse{}
	err := hdlr.Templates(ctx, req_temp, resp_temp)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_temp.Data.Surveys) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_temp.Data.Surveys[0].Id != survey.Id {
		t.Error("Id does not match")
		return
	}
}

func TestFilter(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
		return
	}

	req_filter := &survey_proto.FilterRequest{
		Status:       []survey_proto.SurveyStatus{survey_proto.SurveyStatus_DRAFT},
		Tags:         []string{"tag1", "tag3"},
		RenderTarget: []survey_proto.RenderTarget{survey_proto.RenderTarget_WEB, survey_proto.RenderTarget_MOBILE},
		Visibility:   []static_proto.Visibility{static_proto.Visibility_PUBLIC},
	}
	resp_filter := &survey_proto.FilterResponse{}
	err := hdlr.Filter(ctx, req_filter, resp_filter)
	if err != nil {
		t.Error(err)
	}
	if len(resp_filter.Data.Surveys) == 0 {
		t.Error("Count does not match")
		return
	}
}

func TestSearch(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())
	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
		return
	}

	req_search := &survey_proto.SearchRequest{
		Name:   survey.Title,
		OrgId:  survey.OrgId,
		Offset: 0,
		Limit:  10,
	}
	resp_search := &survey_proto.SearchResponse{}
	err := hdlr.Search(ctx, req_search, resp_search)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_search.Data.Surveys) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_search.Data.Surveys[0].Id != survey.Id {
		t.Error("Id does not match")
		return
	}
}

func TestShareSurvey(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
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

	req_share := &survey_proto.ShareSurveyRequest{
		Surveys: []*survey_proto.Survey{survey},
		Users:   []*user_proto.User{survey.Shares[0]},
		UserId:  si.UserId,
		OrgId:   si.OrgId,
	}

	rsp_share := &survey_proto.ShareSurveyResponse{}
	if err := hdlr.ShareSurvey(ctx, req_share, rsp_share); err != nil {
		t.Error(err)
		return
	}

}

func TestParsingSurvey(t *testing.T) {
	js := `{  
		"title":"new survey with questions",
		"creator":{  
		   "id":"4158873242954991778"
		},
		"org_id":"vzvG73i41edPNxJ7",
		"welcome":{  
		   "type":0,
		   "order":0,
		   "settings":{  
			  "showButton":true,
			  "buttonText":"",
			  "social_sharing_enabled":false,
			  "submit_mode":1,
			  "showTimeToAnswer":false
		   }
		},
		"thankyou":{  
		   "type":9,
		   "order":0,
		   "settings":{  
			  "showButton":true,
			  "buttonText":"",
			  "social_sharing_enabled":false,
			  "submit_mode":1,
			  "showTimeToAnswer":false
		   }
		},
		"renders":[  
	 
		],
		"tags":[  
	 
		],
		"setting":{  
		   "shareableLink":"http://hel.ly/2eA8KKW",
		   "linkSharingEnabled":true
		},
		"id":"GJqk5TMoSMfIXswT",
		"shares":[  
	 
		],
		"questions":[  
			{  
				"id":"ydbUs",
				"type":1,
				"order":1,
				"design":{  
					"bg_color":"fffff",
					"logo_url":"http://via.placeholder.com/300x300"
				},
				"title":"question 1?"
			},
			{  
				"id":"Hu7j6",
				"type":7,
				"order":2,
				"design":{  
					"bg_color":"fffff",
					"logo_url":"http://via.placeholder.com/300x300"
				},
				"settings":{  
					"@type":"healum.com/proto/go.micro.srv.survey.BinaryQuestionSettings",
					"mandatory":true,
					"image":false,
					"video":false,
					"buttonType":0
				},
				"title":"question 2?"
			}
		]
	 }`
	// var p map[string]interface{}
	// decoder := json.NewDecoder(bytes.NewReader([]byte(js)))
	// err := decoder.Decode(&p)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// t.Errorf("%+v", p)

	// // getting json string from json object
	// b, err := json.Marshal(p)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// t.Errorf("%+v", string(b))

	var pp survey_proto.Survey
	// getting Any object from json string
	if err := jsonpb.Unmarshal(strings.NewReader(js), &pp); err != nil {
		t.Error(err)
		return
	}

	t.Log(pp)
}

func TestParsingQuestion(t *testing.T) {
	b, err := json.Marshal(question_map)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("%+v", string(b))

	var pp survey_proto.Question
	// getting Any object from json string
	if err := jsonpb.Unmarshal(strings.NewReader(string(b)), &pp); err != nil {
		t.Error(err)
		return
	}
	t.Log(pp)
}

func TestAutocompleteSearch(t *testing.T) {
	initDb()
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	survey := createSurvey(ctx, hdlr, t)
	if survey == nil {
		return
	}

	req := &survey_proto.AutocompleteSearchRequest{"t"}
	rsp := &survey_proto.AutocompleteSearchResponse{}
	err := hdlr.AutocompleteSearch(ctx, req, rsp)
	if err != nil {
		t.Error(err)
		return
	}

	if len(rsp.Data.Response) == 0 {
		t.Error("Object count does not matched")
		return
	}

	t.Log(rsp.Data.Response)
}

func TestGetTopTags(t *testing.T) {
	TestAll(t)

	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	rsp := &survey_proto.GetTopTagsResponse{}
	if err := hdlr.GetTopTags(ctx, &survey_proto.GetTopTagsRequest{
		OrgId: survey.OrgId,
		N:     5,
	}, rsp); err != nil {
		t.Error(err)
		return
	}

	t.Log(rsp.Data.Tags)
}

func TestAutocompleteTags(t *testing.T) {
	TestAll(t)

	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	rsp := &survey_proto.AutocompleteTagsResponse{}
	if err := hdlr.AutocompleteTags(ctx, &survey_proto.AutocompleteTagsRequest{
		OrgId: survey.OrgId,
		Name:  "t",
	}, rsp); err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp.Data.Tags)
}
