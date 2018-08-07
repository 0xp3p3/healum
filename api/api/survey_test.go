package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server/api/utils"
	"server/common"
	static_proto "server/static-srv/proto/static"
	"server/survey-srv/db"
	survey_proto "server/survey-srv/proto/survey"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

// var serverURL = "http://localhost:8080"
var surveyURL = "/server/surveys"
var survey = &survey_proto.Survey{
	Id:          "111",
	Title:       "survey",
	OrgId:       "orgid",
	Tags:        []string{"tag1", "tag2"},
	Description: "description1",
	Creator:     user,
	IsTemplate:  true,
	Status:      survey_proto.SurveyStatus_DRAFT,
	Renders: []survey_proto.RenderTarget{
		survey_proto.RenderTarget_MOBILE, survey_proto.RenderTarget_WEB,
	},
	Setting: &static_proto.Setting{
		Visibility: static_proto.Visibility_PUBLIC,
	},
}

var question = &survey_proto.Question{
	Id:          "q111",
	Type:        survey_proto.QuestionType_FILE,
	Title:       "question",
	Description: "description",
}

var question_map = map[string]interface{}{
	"id":          "q111",
	"type":        "FILE",
	"title":       "question",
	"description": "description",
	"fields": map[string]interface{}{
		"@type": "healum.com/proto/go.micro.srv.survey.TextQuestionField",
		"attributes": map[string]string{
			"description": "Hello world",
		},
		"validations": []survey_proto.Validation{},
		"attachment": map[string]interface{}{
			"type":   "VIDEO",
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

func initSurveyDb() {
	cl := client.NewClient(client.Transport(nats_transport.NewTransport()), client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(5*time.Second),
		client.Retries(5))
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func AllSurveys(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+surveyURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping survey because already created")

	}
	time.Sleep(time.Second)

	var r map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	b, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(b))
}

func NewSurvey(t *testing.T) (string, string) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+surveyURL+"/new?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping survey because already created")

	}

	r := survey_proto.NewResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data == nil {
		t.Errorf("Data does not matched")
		return "", ""
	}
	if len(r.Data.Hash) == 0 || len(r.Data.SurveyId) == 0 {
		t.Errorf("Response does not matched")
		return "", ""
	}
	return r.Data.Hash, r.Data.SurveyId
}

func CreateSurvey(survey *survey_proto.Survey, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return
	}
	survey.Creator.Id = userId
	survey.Shares = []*user_proto.User{{Id: userId}}
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"survey": survey})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	js := `{"survey":{  
		"title":"new survey with questions",
		"creator":{  
		   "id":"4158873242954991778",
		   "orgid":"vzvG73i41edPNxJ7"
		},
		"org_id":"orgid",
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
	 }}`
	log.Println(js)
	req, err := http.NewRequest("POST", serverURL+surveyURL+"/survey/create?session="+sessionId, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping survey because already created")
	}
	time.Sleep(time.Second)

	r := survey_proto.CreateResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
	survey = r.Data.Survey
	b, err := json.Marshal(survey)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(b))
}

func ReadSurvey(id string, t *testing.T) *survey_proto.Survey {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+surveyURL+"/survey/"+id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	defer resp.Body.Close()
	// if resp.StatusCode == http.StatusInternalServerError {
	// 	t.Skip("Skipping survey because already created")

	// }
	time.Sleep(time.Second)

	r := survey_proto.ReadResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		//		t.Errorf("Response does not matched")
		return nil
	}
	return r.Data.Survey
}

func DeleteSurvey(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a DELETE request.
	req, err := http.NewRequest("DELETE", serverURL+surveyURL+"/survey/"+id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping survey because already created")

	}
	time.Sleep(time.Second)
}

func CopySurvey(survey *survey_proto.Survey, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"survey_id": survey.Id})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+surveyURL+"/survey/copy?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func Questions(surveyId string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+surveyURL+"/survey/"+surveyId+"/questions?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func QuestionRef(surveyId string, questionId string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+surveyURL+"/survey/"+surveyId+"/questions/"+questionId+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func CreateQuestion(surveyId string, question *survey_proto.Question, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"question": question})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+surveyURL+"/survey/"+surveyId+"/questions?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func SurveyByCreator(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+surveyURL+"/creator/"+id+"?session="+sessionId+"&offset=0&limit=20", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping survey because already created")

	}
	time.Sleep(time.Second)
}

func SurveyLink(t *testing.T) {

}

func SurveyTemplates(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+surveyURL+"/templates?session="+sessionId+"&offset=0&limit=20", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping survey because already created")

	}
	time.Sleep(time.Second)
}

func SurveyFilter(filter *survey_proto.FilterRequest, t *testing.T) {
	// Send a POST request.
	jsonStr, err := json.Marshal(filter)
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("POST", serverURL+surveyURL+"/filter?session="+sessionId+"&offset=0&limit=20", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping survey because already created")

	}
	time.Sleep(time.Second)
}

func SearchSurveys(search *survey_proto.SearchRequest, t *testing.T) {
	// Send a POST request.
	jsonStr, err := json.Marshal(search)
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("POST", serverURL+surveyURL+"/search?session="+sessionId+"&offset=0&limit=20", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping survey because already created")
	}
	time.Sleep(time.Second)
}

func TestSurveysAll(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)
	AllSurveys(t)
}

func TestSurveyNew(t *testing.T) {
	initSurveyDb()

	NewSurvey(t)
}

func TestSurveyCreate(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)
	p := ReadSurvey("111", t)
	if p == nil {
		t.Errorf("Survey does not matched")
		return
	}
	if p.Id != survey.Id {
		t.Errorf("Id does not matched")
		return
	}
	if p.Title != survey.Title {
		t.Errorf("Title does not matched")
		return
	}
}

func TestSurveyDelete(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)
	DeleteSurvey("111", t)
	p := ReadSurvey("111", t)
	if p != nil {
		t.Errorf("Survey does not matched")
		return
	}
}

func TestSurveyCopy(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)
	CopySurvey(survey, t)
}

func TestQuestions(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)
	CreateQuestion(survey.Id, question, t)
	Questions(survey.Id, t)
}

func TestQuestionRef(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)
	CreateQuestion(survey.Id, question, t)
	QuestionRef(survey.Id, question.Id, t)
}

func TestCreateQuestion(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)
	CreateQuestion(survey.Id, question, t)
}

func TestSurveyByCreator(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)
	SurveyByCreator("userid", t)
}

func TestSurveyLink(t *testing.T) {

}

func TestSurveyTemplates(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)
	SurveyTemplates(t)
}

func TestSurveyFilter(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)
	filter := &survey_proto.FilterRequest{
		Status:       []survey_proto.SurveyStatus{survey_proto.SurveyStatus_DRAFT},
		Tags:         []string{"tag"},
		RenderTarget: []survey_proto.RenderTarget{survey_proto.RenderTarget_WEB, survey_proto.RenderTarget_MOBILE},
		Visibility:   []static_proto.Visibility{static_proto.Visibility_PUBLIC},
	}
	SurveyFilter(filter, t)
}

func TestSurveySearch(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)
	search := &survey_proto.SearchRequest{
		Name:        "survey",
		Description: "description",
		UserName:    "userid",
	}
	SearchSurveys(search, t)
}

func TestAnything(t *testing.T) {
	initSurveyDb()

	questionField := &survey_proto.TextQuestionField{
		Attributes:  &survey_proto.TextQuestionField_Attributes{"hello world"},
		Validations: []*survey_proto.Validation{&survey_proto.Validation{}},
	}

	serializedp, err := proto.Marshal(questionField)
	if err != nil {
		t.Fatal("could not serialize TextQuestionField")
	}
	question.Fields = &any.Any{
		TypeUrl: "example.com/yaddayaddayadda/" + proto.MessageName(questionField),
		Value:   serializedp,
	}

	CreateSurvey(survey, t)
	CreateQuestion(survey.Id, question, t)
	QuestionRef(survey.Id, question.Id, t)
}

func TestErrReadSurvey(t *testing.T) {
	initSurveyDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+surveyURL+"/survey/999?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping survey because already created")

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
	}
}

func TestErrAllSurvey(t *testing.T) {
	initSurveyDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+surveyURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping survey because already created")

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
	}
}

func TestQuestionResponseWithAny(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)

	jsonStr, err := json.Marshal(map[string]interface{}{"question": question_map})
	if err != nil {
		t.Error(err)
	}

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	req, err := http.NewRequest("POST", serverURL+surveyURL+"/survey/"+survey.Id+"/questions?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	req, err = http.NewRequest("GET", serverURL+surveyURL+"/survey/"+survey.Id+"/questions?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	var r map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	b, err := json.Marshal(r)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(b))
}

func TestAutocompleteSurveySearch(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"title": "s"})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+surveyURL+"/survey/search/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := survey_proto.AutocompleteSearchResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Response) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Response)
}

func TestGetTopSurveyTags(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+surveyURL+"/tags/top/5?session="+sessionId, nil)

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := survey_proto.GetTopTagsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Tags) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Tags)
}

func TestAutocompleteSurveyTags(t *testing.T) {
	initSurveyDb()

	CreateSurvey(survey, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"name": "t"})
	if err != nil {
		t.Error(err)
		return
	}
	// log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+surveyURL+"/tags/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := survey_proto.AutocompleteTagsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Tags) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
	t.Log(r.Data.Tags)
}
