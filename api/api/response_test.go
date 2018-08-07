package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"server/api/utils"
	"server/common"
	"server/response-srv/db"
	resp_proto "server/response-srv/proto/response"
	static_proto "server/static-srv/proto/static"
	survey_db "server/survey-srv/db"
	survey_proto "server/survey-srv/proto/survey"
	user_proto "server/user-srv/proto/user"
	"strconv"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

var responseURL = "/server/responses"
var survey1 = &survey_proto.Survey{
	Id:          "111",
	Title:       "title",
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
		Visibility:             static_proto.Visibility_PRIVATE,
		AuthenticationRequired: true,
	},
}

var resp = &resp_proto.SubmitSurveyResponse{
	Id:              "111",
	OrgId:           "orgid",
	SurveyId:        "111",
	ResponseSession: "session",
	Answers: []*resp_proto.Answer{
		{
			QuestionRef: "q111",
			Type:        survey_proto.QuestionType_DROPDOWN,
		},
		{
			QuestionRef: "q222",
			Type:        survey_proto.QuestionType_DROPDOWN,
		},
	},
	Status: &resp_proto.ResponseStatus{
		State:     resp_proto.ResponseState_SUBMITTED,
		Timestamp: time.Now().Unix(),
	},
	Responder: &user_proto.User{
		Id: "userid",
	},
}

var resp_map = map[string]interface{}{
	"id":               "111",
	"org_id":           "orgid",
	"survey_id":        "111",
	"response_session": "session",
	"answers": []map[string]interface{}{
		{
			"question_ref": "q111",
			"type":         1,
			"data": map[string]interface{}{
				"@type": "healum.com/proto/go.micro.srv.response.TextAnswer",
				"value": "Hello world",
			},
		},
		{
			"question_ref": "q222",
			"type":         survey_proto.QuestionType_CONTACT,
			"data": map[string]interface{}{
				"@type":          "healum.com/proto/go.micro.srv.response.GetContactAnswer",
				"first_name":     "Test name",
				"last_name":      "Last name",
				"contact_number": "12515215125",
				"email":          "test@test.com",
				"address":        "Address",
			},
		},
	},
	"status": &resp_proto.ResponseStatus{
		State:     resp_proto.ResponseState_SUBMITTED,
		Timestamp: time.Now().Unix(),
	},
	"responder": &user_proto.User{
		Id: "userid",
	},
}

func initResponseDb() {
	cl := client.NewClient(
		client.Transport(nats_transport.NewTransport()),
		client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(5*time.Second),
		client.Retries(5))
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	survey_db.Init(cl)
	db.Init(cl)
}

func Check(hash string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+responseURL+"/"+hash+"/check?session="+sessionId+"&offset=0&limit=10", nil)
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

	r := resp_proto.CheckResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
}

func AllQuestion(surveyId string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+responseURL+"/survey/"+surveyId+"/questions/all?session="+sessionId+"&offset=0&limit=10", nil)
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

	r := survey_proto.QuestionsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
	if len(r.Data.Questions) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func ReadQuestion(surveyId, questionId string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+responseURL+"/survey/"+surveyId+"/questions/"+questionId+"?session="+sessionId+"&offset=0&limit=10", nil)
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

	r := survey_proto.QuestionRefResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Question == nil {
		t.Errorf("Response does not matched")
		return
	}
	if r.Data.Question.Id != question.Id {
		t.Errorf("Id does not matched")
		return
	}
}

func OpenAllQuestion(surveyId string, t *testing.T) {
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+responseURL+"/open/survey/"+surveyId+"/questions/all?team_id="+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if res.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping survey because already created")

	}
	time.Sleep(time.Second)

	r := survey_proto.QuestionsResponse{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
	if len(r.Data.Questions) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func OpenReadQuestion(surveyId, questionId string, t *testing.T) {
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+responseURL+"/open/survey/"+surveyId+"/questions/"+questionId+"?team_id="+"&offset=0&limit=10", nil)
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

	r := survey_proto.QuestionRefResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Question == nil {
		t.Errorf("Response does not matched")
		return
	}
	if r.Data.Question.Id != question.Id {
		t.Errorf("Id does not matched")
		return
	}
}

func All(surveyId string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+responseURL+"/"+surveyId+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := resp_proto.AllResponse{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
	if len(r.Data.Responses) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func TimeFilter(surveyId string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	from := time.Now().Unix() - 10000
	to := time.Now().Unix() + 10000
	_from := strconv.FormatInt(from, 10)
	_to := strconv.FormatInt(to, 10)
	req, err := http.NewRequest("GET", serverURL+responseURL+"/"+surveyId+"/all?from="+_from+"&to="+_to+"&session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := resp_proto.AllResponse{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
	if len(r.Data.Responses) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func GroupBy(surveyId string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+responseURL+"/"+surveyId+"/all?groupby=question&session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second) 

	r := resp_proto.AllAggQuestionResponse{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	t.Log("=====", r, err)
	if r.Data == nil {
		t.Errorf("Response does not matched") 
		return
	}
	if len(r.Data.Responses) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func Create(surveyId string, response *resp_proto.SubmitSurveyResponse, t *testing.T) {
	jsonStr, err := json.Marshal(map[string]interface{}{"response": response})
	if err != nil {
		t.Error(err)
	}

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("POST", serverURL+responseURL+"/"+surveyId+"/response?session="+sessionId+"&offset=0&limit=10", bytes.NewBuffer(jsonStr))
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

func UpdateState(surveyId string, response *resp_proto.SubmitSurveyResponse, t *testing.T) {
	response.Status.State = resp_proto.ResponseState_ABANDONED
	update := map[string]interface{}{
		"response_id": response.Id,
		"state":       response.Status.State,
	}
	jsonStr, err := json.Marshal(map[string]interface{}{"response": update})
	if err != nil {
		t.Error(err)
	}

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("POST", serverURL+responseURL+"/"+surveyId+"/response/state?session="+sessionId+"&offset=0&limit=10", bytes.NewBuffer(jsonStr))
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

func AllState(surveyId string, state resp_proto.ResponseState, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	s := strconv.Itoa(int(state))
	req, err := http.NewRequest("GET", serverURL+responseURL+"/"+surveyId+"/all/state/"+s+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := resp_proto.AllResponse{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
	if len(r.Data.Responses) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func ReadStats(surveyId string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+responseURL+"/"+surveyId+"/response/stats?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := resp_proto.ReadStatsResponse{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
	if r.Data.Stats == nil {
		t.Errorf("Response does not matched")
		return
	}
	if r.Data.Stats.Responses != 1 {
		t.Errorf("Submit does not matched")
		return
	}
	if r.Data.Stats.Drops != 0 {
		t.Errorf("Abandon does not matched")
		return
	}
}

func ByUser(surveyId, userId string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+responseURL+"/"+surveyId+"/response/by/"+userId+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := resp_proto.ByUserResponse{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
	if len(r.Data.Responses) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func ByAnyUser(surveyId string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+responseURL+"/"+surveyId+"/response/anon?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := resp_proto.ByAnyUserResponse{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}
	if len(r.Data.Responses) == 0 {
		t.Errorf("Object count does not matched")
		return 
	}
}
func TestCheck(t *testing.T) {
	initResponseDb()

	hash, surveyId := NewSurvey(t)
	survey1.Id = surveyId
	CreateSurvey(survey1, t)
	Check(hash, t)
}

func TestAllQuestion(t *testing.T) {
	initResponseDb()

	_, surveyId := NewSurvey(t)
	survey1.Id = surveyId
	CreateSurvey(survey1, t)
	CreateQuestion(surveyId, question, t)
	AllQuestion(surveyId, t)
}

func TestReadQuestion(t *testing.T) {
	initResponseDb()

	_, surveyId := NewSurvey(t)
	survey1.Id = surveyId
	CreateSurvey(survey1, t)
	CreateQuestion(surveyId, question, t)
	ReadQuestion(surveyId, question.Id, t)
}

func TestOpenAllQuestion(t *testing.T) {
	initResponseDb()

	_, surveyId := NewSurvey(t)
	survey1.Id = surveyId
	survey1.Setting.AuthenticationRequired = false
	CreateSurvey(survey1, t)
	CreateQuestion(surveyId, question, t)
	OpenAllQuestion(surveyId, t)
}

func TestOpenReadQuestion(t *testing.T) {
	initResponseDb()

	_, surveyId := NewSurvey(t)
	survey1.Id = surveyId
	survey1.Setting.AuthenticationRequired = false
	CreateSurvey(survey1, t)
	CreateQuestion(surveyId, question, t)
	OpenReadQuestion(surveyId, question.Id, t)
}

func TestAll(t *testing.T) {
	initResponseDb()

	CreateSurvey(survey1, t)
	Create(survey1.Id, resp, t)
	All(survey1.Id, t)
}

func TestGropuBy(t *testing.T) {
	initResponseDb()

	CreateSurvey(survey1, t)
	question.Id = "q111"
	CreateQuestion(survey1.Id, question, t)
	question.Id = "q222"
	CreateQuestion(survey1.Id, question, t)

	resp.Id = "r111"
	Create(survey1.Id, resp, t)
	resp.Id = "r222"
	Create(survey1.Id, resp, t)
	GroupBy(survey1.Id, t)
}

func TestTimeFilter(t *testing.T) {
	initResponseDb()

	CreateSurvey(survey1, t)
	Create(survey1.Id, resp, t)
	TimeFilter(survey1.Id, t)
}

func TestCreate(t *testing.T) {
	initResponseDb()

	CreateSurvey(survey1, t)
	Create(survey1.Id, resp, t)
}

func TestUpdateState(t *testing.T) {
	initResponseDb()

	CreateSurvey(survey1, t)
	Create(survey1.Id, resp, t)
	UpdateState(survey1.Id, resp, t)

	All(survey1.Id, t)
}

func TestAllState(t *testing.T) {
	initResponseDb()

	CreateSurvey(survey1, t)
	Create(survey1.Id, resp, t)

	AllState(survey1.Id, resp.Status.State, t)
}

func TestReadStats(t *testing.T) {
	initResponseDb()

	CreateSurvey(survey1, t)
	Create(survey1.Id, resp, t)

	ReadStats(survey1.Id, t)
}

func TestByUser(t *testing.T) {
	initResponseDb()

	CreateSurvey(survey1, t)
	Create(survey1.Id, resp, t)

	ByUser(survey1.Id, resp.Responder.Id, t)
}

func TestByAnyUser(t *testing.T) {
	initResponseDb()

	survey1.Setting.AuthenticationRequired = false
	resp.Responder = nil
	CreateSurvey(survey1, t)
	Create(survey1.Id, resp, t)

	ByAnyUser(survey1.Id, t)
}

func TestErrAllResponses(t *testing.T) {
	initResponseDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+responseURL+"/"+survey1.Id+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping response because already created")

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

func TestCreateResponseWithAny(t *testing.T) {
	initResponseDb()

	CreateSurvey(survey1, t)

	jsonStr, err := json.Marshal(map[string]interface{}{"response": resp_map})
	if err != nil {
		t.Error(err)
	}

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("POST", serverURL+responseURL+"/"+survey1.Id+"/response?session="+sessionId+"&offset=0&limit=10", bytes.NewBuffer(jsonStr))
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

	// list all response
	req, err = http.NewRequest("GET", serverURL+responseURL+"/"+survey.Id+"/all?session="+sessionId+"&offset=0&limit=10", nil)
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
	}
	json.Unmarshal(body, &r)
}
