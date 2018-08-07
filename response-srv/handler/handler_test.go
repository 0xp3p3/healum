package handler

import (
	"context"
	"encoding/json"
	"reflect"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	"server/response-srv/db"
	resp_proto "server/response-srv/proto/response"
	static_proto "server/static-srv/proto/static"
	survey_db "server/survey-srv/db"
	survey_hdlr "server/survey-srv/handler"
	survey_proto "server/survey-srv/proto/survey"
	user_proto "server/user-srv/proto/user"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
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
	survey_db.Init(cl)
	db.Init(cl)
}

var survey = &survey_proto.Survey{
	Id:          "111",
	Title:       "title",
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
		Visibility:             static_proto.Visibility_PRIVATE,
		AuthenticationRequired: true,
	},
	Shares: []*user_proto.User{{Id: "userid"}},
}

var response = &resp_proto.SubmitSurveyResponse{
	Id:              "111",
	OrgId:           "orgid",
	SurveyId:        "111",
	ResponseSession: "session",
	Answers: []*resp_proto.Answer{
		{
			QuestionRef: "q111",
			Type:        survey_proto.QuestionType_DROPDOWN,
			// Data: &google_protobuf.Any{
			// 	TypeUrl: "type.googleapis.com/go.micro.srv.response.GetContactAnswer",
			// 	Value:   []byte("\n\005xiang\022\004tian\032\006123456\"\013email@e.com*\007address"),
			// },
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

var question = &survey_proto.Question{
	Id:          "q111",
	Type:        survey_proto.QuestionType_FILE,
	Title:       "question",
	Description: "description",
}

var hash = ""

func initSurveyHandler() *survey_hdlr.SurveyService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	survey_hdlr := &survey_hdlr.SurveyService{
		Broker:        nats_brker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
	}
	return survey_hdlr
}

func createSurvey(ctx context.Context, t *testing.T) *survey_proto.Survey {
	// make new hash
	hdlr_survey := initSurveyHandler()
	req_new := &survey_proto.NewRequest{}
	resp_new := &survey_proto.NewResponse{}
	err := hdlr_survey.New(ctx, req_new, resp_new)
	if err != nil {
		t.Error(err)
		return nil
	}

	// create new survey with unique id
	hash = resp_new.Data.Hash
	survey.Id = resp_new.Data.SurveyId
	if err := survey_db.Create(ctx, survey); err != nil {
		t.Error(err)
		return nil
	}

	return survey
}

func TestCheck(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	survey := createSurvey(ctx, t)
	if survey == nil {
		return
	}
	time.Sleep(time.Second)

	hdlr := new(ResponseService)
	req := &resp_proto.CheckRequest{ShortHash: hash}
	resp := &resp_proto.CheckResponse{}
	res := hdlr.Check(ctx, req, resp)
	if res != nil {
		t.Error(res)
	}
}

func createResponse(ctx context.Context, hdlr *ResponseService, t *testing.T) *resp_proto.SubmitSurveyResponse {
	survey := createSurvey(ctx, t)
	if survey == nil {
		return nil
	}

	response.SurveyId = survey.Id
	req := &resp_proto.CreateRequest{SurveyId: survey.Id, Response: response}
	resp := &resp_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req, resp); err != nil {
		t.Error(err)
		return nil
	}

	return resp.Data.Response
}

func TestResponseIsCreate(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ResponseService)

	if response := createResponse(ctx, hdlr, t); response == nil {
		t.Error("Response is not created")
	}
}

func TestAll(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ResponseService)

	response := createResponse(ctx, hdlr, t)
	if response == nil {
		return
	}

	req := &resp_proto.AllRequest{SurveyId: response.SurveyId}
	res := &resp_proto.AllResponse{}
	if err := hdlr.All(ctx, req, res); err != nil {
		t.Error(err)
		return
	}

	if len(res.Data.Responses) == 0 {
		t.Error("Object count does not match")
		return
	}
	if res.Data.Responses[0].Id != response.Id {
		t.Error("Id does not match")
		return
	}
}

func TestUpdateState(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ResponseService)

	response := createResponse(ctx, hdlr, t)
	if response == nil {
		return
	}

	req := &resp_proto.UpdateStateRequest{
		Response: &resp_proto.UpdateStateRequest_Response{
			ResponseId: response.Id,
			State:      resp_proto.ResponseState_ABANDONED,
		},
		SurveyId: response.SurveyId,
	}
	res := &resp_proto.UpdateStateResponse{}
	if err := hdlr.UpdateState(ctx, req, res); err != nil {
		t.Error(err)
		return
	}

	req_all := &resp_proto.AllRequest{SurveyId: response.SurveyId}
	resp_all := &resp_proto.AllResponse{}
	if err := hdlr.All(ctx, req_all, resp_all); err != nil {
		t.Error(err)
	}

	if len(resp_all.Data.Responses) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Responses[0].Id != response.Id {
		t.Error("Id does not match")
		return
	}
	if resp_all.Data.Responses[0].Status.State != resp_proto.ResponseState_ABANDONED {
		t.Error("State is not updated")
		return
	}
}

func TestAllState(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ResponseService)

	response := createResponse(ctx, hdlr, t)
	if response == nil {
		return
	}

	req_all := &resp_proto.AllStateRequest{
		SurveyId: response.SurveyId,
		State:    response.Status.State,
	}
	resp_all := &resp_proto.AllStateResponse{}
	if err := hdlr.AllState(ctx, req_all, resp_all); err != nil {
		t.Error(err)
		return
	}

	if len(resp_all.Data.Responses) == 0 {
		t.Error("Object count does not match")
		return
	}
	if resp_all.Data.Responses[0].Id != response.Id {
		t.Error("Id does not match")
		return
	}
	if resp_all.Data.Responses[0].Status.State != response.Status.State {
		t.Error("State is not updated")
		return
	}
}

func TestTimeFilter(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ResponseService)
	response := createResponse(ctx, hdlr, t)
	if response == nil {
		return
	}

	req := &resp_proto.TimeFilterRequest{SurveyId: response.SurveyId,
		From: time.Now().Unix() - 10000,
		To:   time.Now().Unix() + 10000,
	}
	res := &resp_proto.TimeFilterResponse{}
	if err := hdlr.TimeFilter(ctx, req, res); err != nil {
		t.Error(err)
		return
	}

	if len(res.Data.Responses) == 0 {
		t.Error("Object count does not match")
		return
	}
	if res.Data.Responses[0].Id != response.Id {
		t.Error("Id does not match")
		return
	}
}

func TestReadStats(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ResponseService)
	response := createResponse(ctx, hdlr, t)
	if response == nil {
		return
	}

	req := &resp_proto.ReadStatsRequest{SurveyId: response.SurveyId}
	res := &resp_proto.ReadStatsResponse{}
	if err := hdlr.ReadStats(ctx, req, res); err != nil {
		t.Error(err)
		return
	}
	if res.Data.Stats.Responses != 1 {
		t.Error("Responses count does not match")
		return
	}
	if res.Data.Stats.Drops != 0 {
		t.Error("Drops count does not match")
		return
	}
}

func TestByUser(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ResponseService)
	response := createResponse(ctx, hdlr, t)
	if response == nil {
		return
	}

	req := &resp_proto.ByUserRequest{SurveyId: response.SurveyId, UserId: "userid"}
	res := &resp_proto.ByUserResponse{}
	if err := hdlr.ByUser(ctx, req, res); err != nil {
		t.Error(err)
		return
	}

	if len(res.Data.Responses) == 0 {
		t.Error("Object count does not match")
		return
	}
	if res.Data.Responses[0].Id != response.Id {
		t.Error("Id does not match")
		return
	}
}

func TestByAnyUser(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := new(ResponseService)
	response := createResponse(ctx, hdlr, t)
	if response == nil {
		return
	}

	req := &resp_proto.ByAnyUserRequest{SurveyId: response.SurveyId}
	res := &resp_proto.ByAnyUserResponse{}
	if err := hdlr.ByAnyUser(ctx, req, res); err != nil {
		t.Error(err)
		return
	}
	if len(res.Data.Responses) != 0 {
		t.Error("Object count does not match")
		return
	}
}

func TestGroupBy(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	survey := createSurvey(ctx, t)
	if survey == nil {
		return
	}

	hdlr_survey := initSurveyHandler()
	// create question
	question.Id = "q111"
	req_question := &survey_proto.CreateQuestionRequest{SurveyId: survey.Id, Question: question}
	resp_question := &survey_proto.CreateQuestionResponse{}
	time.Sleep(time.Second)
	err := hdlr_survey.CreateQuestion(ctx, req_question, resp_question)
	if err != nil {
		t.Error(err)
		return
	}

	// create question2
	question.Id = "q222"
	req_question = &survey_proto.CreateQuestionRequest{SurveyId: survey.Id, Question: question}
	resp_question = &survey_proto.CreateQuestionResponse{}
	time.Sleep(time.Second)
	err = hdlr_survey.CreateQuestion(ctx, req_question, resp_question)
	if err != nil {
		t.Error(err)
		return
	}

	response.Id = "r111"
	response.SurveyId = survey.Id
	hdlr := new(ResponseService)
	req := &resp_proto.CreateRequest{SurveyId: survey.Id, Response: response}
	res := &resp_proto.CreateResponse{}
	time.Sleep(time.Second)
	err = hdlr.Create(ctx, req, res)
	if err != nil {
		t.Error(err)
	}

	response.Id = "r222"
	req = &resp_proto.CreateRequest{SurveyId: survey.Id, Response: response}
	res = &resp_proto.CreateResponse{}
	time.Sleep(time.Second)
	err = hdlr.Create(ctx, req, res)
	if err != nil {
		t.Error(err)
	}

	req_group := &resp_proto.AllAggQuestionRequest{SurveyId: survey.Id, OrgId: "orgid"}
	res_group := &resp_proto.AllAggQuestionResponse{}
	err = hdlr.AllAggQuestion(ctx, req_group, res_group)
	if err != nil {
		t.Error(err)
	}
}

func TestDemo(t *testing.T) {
	name := proto.MessageName(&resp_proto.TextAnswer{})
	t.Log(name)
}

func TestAnything(t *testing.T) {
	t1 := &resp_proto.GetContactAnswer{
		FirstName:     "xiang",
		LastName:      "tian",
		ContactNumber: "123456",
		Email:         "email@e.com",
		Address:       "address",
	}
	serialized, err := ptypes.MarshalAny(t1)
	if err != nil {
		t.Fatal("could not serialize TextAnser")
	}
	// Blue was a great album by 3EB, before Cadgogan got kicked out
	// and Jenkins went full primadonna
	a := resp_proto.Answer{
		QuestionRef: "Queeesion",
		Type:        survey_proto.QuestionType_DROPDOWN,
		Data:        serialized,
		// &any.Any{
		// 	TypeUrl: "healum.com/proto/" + proto.MessageName(t1),
		// 	Value:   serialized,
		// },
	}
	// marshal to simulate going on the wire:
	serializedA, err := proto.Marshal(&a)
	if err != nil {
		t.Fatal("could not serialize anything")
	}
	// unmarshal to simulate coming off the wire
	var a2 resp_proto.Answer
	if err := proto.Unmarshal(serializedA, &a2); err != nil {
		t.Fatal("could not deserialize anything")
	}
	// unmarshal the timestamp
	var t2 resp_proto.GetContactAnswer
	if err := ptypes.UnmarshalAny(a2.Data, &t2); err != nil {
		t.Fatalf("Could not unmarshal timestamp from anything field: %s", err)
	}
	// Verify the values are as expected
	if !reflect.DeepEqual(t1, &t2) {
		t.Fatalf("Values don't match up:\n %+v \n %+v", t1, t2)
	}

	t.Logf("%+v", serialized)
	t.Logf("%+v", string(serializedA))
	t.Logf("%+v", t2)
}

func TestExample(t *testing.T) {
	t1 := &resp_proto.GetContactAnswer{}

	ex := map[string]string{
		"@type":          "healum.com/proto/go.micro.srv.response.GetContactAnswer",
		"first_name":     "xiang",
		"last_name":      "tian",
		"contact_number": "123456",
		"Email":          "email@e.com",
		"address":        "address",
	}

	serialized, _ := json.Marshal(ex)
	t.Log(string(serialized))

	// serialized, err := proto.Marshal(t1)
	// if err != nil {
	// 	t.Fatal("could not serialize TextAnser")
	// }
	// Blue was a great album by 3EB, before Cadgogan got kicked out
	// and Jenkins went full primadonna
	a := resp_proto.Answer{
		QuestionRef: "Queeesion",
		Type:        survey_proto.QuestionType_DROPDOWN,
		Data: &any.Any{
			TypeUrl: "healum.com/proto/" + proto.MessageName(t1),
			Value:   serialized,
		},
	}
	// marshal to simulate going on the wire:
	serializedA, err := proto.Marshal(&a)
	if err != nil {
		t.Fatal("could not serialize anything")
	}
	// unmarshal to simulate coming off the wire
	var a2 resp_proto.Answer
	if err := proto.Unmarshal(serializedA, &a2); err != nil {
		t.Fatal("could not deserialize anything")
	}
	// unmarshal the timestamp
	var t2 resp_proto.GetContactAnswer
	if err := ptypes.UnmarshalAny(a2.Data, &t2); err != nil {
		t.Fatalf("Could not unmarshal timestamp from anything field: %s", err)
	}
	// Verify the values are as expected
	if !reflect.DeepEqual(t1, &t2) {
		t.Fatalf("Values don't match up:\n %+v \n %+v", t1, t2)
	}

	t.Logf("%+v", serialized)
	t.Logf("%+v", string(serializedA))
	t.Logf("%+v", t2)
}

func TestAnyWithJson(t *testing.T) {
	var t1 any.Any

	raw1 := `{
		"@type":          "healum.com/proto/go.micro.srv.response.GetContactAnswer",
		"first_name":     "Test name",
		"last_name":      "Last name",
		"contact_number": "12515215125",
		"email":          "test@test.com",
		"address":        "Address"
	}`

	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &t1); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}

	dm := &resp_proto.GetContactAnswer{
		FirstName:     "Test name",
		LastName:      "Last name",
		ContactNumber: "12515215125",
		Email:         "test@test.com",
		Address:       "Address",
	}

	serialized, _ := proto.Marshal(dm)

	a := resp_proto.Answer{
		QuestionRef: "Queeesion",
		Type:        survey_proto.QuestionType_CONTACT,
		Data: &any.Any{
			TypeUrl: "healum.com/proto/go.micro.srv.response.GetContactAnswer",
			Value:   serialized,
		},
	}

	if !proto.Equal(&t1, a.Data) {
		t.Errorf("message contents not set correctly after unmarshalling JSON: got %s, wanted %s", t1, a.Data)
	}
}

func TestAnswerWithJson(t *testing.T) {

	var t2 resp_proto.Answer

	raw1 := `{
			"question_ref":"Queeesion",
			"type": 3,
			"data": {
				"@type":          "healum.com/proto/go.micro.srv.response.GetContactAnswer",
				"first_name":     "Test name",
				"last_name":      "Last name",
				"contact_number": "12515215125",
				"email":          "test@test.com",
				"address":        "Address"
			}
		}`

	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &t2); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
		return
	}

	dm := &resp_proto.GetContactAnswer{
		FirstName:     "Test name",
		LastName:      "Last name",
		ContactNumber: "12515215125",
		Email:         "test@test.com",
		Address:       "Address",
	}

	serialized, _ := proto.Marshal(dm)

	a := resp_proto.Answer{
		QuestionRef: "Queeesion",
		Type:        survey_proto.QuestionType_CONTACT,
		Data: &any.Any{
			TypeUrl: "healum.com/proto/go.micro.srv.response.GetContactAnswer",
			Value:   serialized,
		},
	}

	if !proto.Equal(&a, &t2) {
		t.Errorf("message contents not set correctly after unmarshalling JSON: got %s, wanted %s", &t2, &a)
		return
	}

	// t.Error(t2.Data)
	var c resp_proto.GetContactAnswer
	if err := ptypes.UnmarshalAny(t2.Data, &c); err != nil {
		t.Fatalf("Could not unmarshal timestamp from anything field: %s", err)
		return
	}
	// t.Errorf("%+v", c)

	// t.Errorf("+%v", string(b))
	marshaler := jsonpb.Marshaler{}
	js, err := marshaler.MarshalToString(&a)
	if err != nil {
		t.Errorf("an unexpected error occurred when marshaling any to JSON: %v", err)
		return
	}
	t.Log(js)
}
