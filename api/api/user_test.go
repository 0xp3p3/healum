package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	account_proto "server/account-srv/proto/account"
	"server/api/utils"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	content_proto "server/content-srv/proto/content"
	kv_proto "server/kv-srv/proto/kv"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
	track_proto "server/track-srv/proto/track"
	userapp_db "server/user-app-srv/db"
	userapp_proto "server/user-app-srv/proto/userapp"
	"server/user-srv/db"
	user_hdlr "server/user-srv/handler"
	user_proto "server/user-srv/proto/user"
	"strings"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/protobuf/jsonpb"
	google_protobuf1 "github.com/golang/protobuf/ptypes/struct"
	nats_broker "github.com/micro/go-plugins/broker/nats"
)

var serverURL = "http://localhost:8080"
var userURL = "/server/users"

var user = &user_proto.User{
	OrgId:     "orgid",
	Firstname: "david",
	Lastname:  "john",
	AvatarUrl: "http://example.com",
	// Tokens: map[string]*user_proto.Token{
	// 	"a1": {"11671c2e7da30e3c393813f60b327f9c2e2e08390761aa01e37ba5d3e6a617be", 1, "a1", "aaa"},
	// 	"b1": {"token_b", 2, "b1", "bbb"},
	// },
}

func initUserHandler() *user_hdlr.UserService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()

	hdlr := &user_hdlr.UserService{
		Broker:          nats_brker,
		KvClient:        kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		AccountClient:   account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		TrackClient:     track_proto.NewTrackServiceClient("go.micro.srv.track", cl),
		TeamClient:      team_proto.NewTeamServiceClient("go.micro.srv.team", cl),
		StaticClient:    static_proto.NewStaticServiceClient("go.micro.srv.static", cl),
		BehaviourClient: behaviour_proto.NewBehaviourServiceClient("go.micro.srv.behaviour", cl),
	}
	return hdlr
}

func initUserDb() {
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func GenerateRandNumber(n int) string {
	letters := []rune("1234567890")
	rand.Seed(time.Now().UTC().UnixNano())
	randomString := make([]rune, n)
	for i := range randomString {
		randomString[i] = letters[rand.Intn(len(letters))]
	}
	return string(randomString)
}

func CreateUser(user *user_proto.User, t *testing.T) *user_proto.User {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	account.Email = "email" + GenerateRandNumber(4) + "@email.com"
	user.Id = ""
	jsonStr, err := json.Marshal(map[string]interface{}{"user": user, "account": account})
	if err != nil {
		t.Error(err)
		return nil
	}

	req, err := http.NewRequest("POST", serverURL+userURL+"/user?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return nil
	}
	// if resp.StatusCode == http.StatusInternalServerError {
	// 	t.Skip("Skipping user because already created")
	// }
	time.Sleep(2 * time.Second)

	r := user_proto.CreateResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return nil
	}
	json.Unmarshal(body, &r)
	user = r.Data.User
	log.Println("user:", user)
	return user
}

func TestUserIsCreated(t *testing.T) {
	// initUserDb()
	CreateUser(user, t)
}

func TestAllUsers(t *testing.T) {
	CreateUser(user, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	req, err := http.NewRequest("GET", serverURL+userURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := user_proto.AllResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}

	if len(r.Data.Users) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestReadUser(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	req, err := http.NewRequest("GET", serverURL+userURL+"/user/userid?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := user_proto.ReadResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}

	if r.Data.User.Id != user.Id {
		t.Error("Object id does not matched")
		return
	}
}

func TestFilterUsers(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"status": account_proto.AccountStatus_SUSPENDED})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("POST", serverURL+userURL+"/user/filter?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := user_proto.FilterUserResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}

	if len(r.Data.Response) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestNotOrgAuthorized(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	if len(sessionId) == 0 {
		t.Error("sessionId is invalid")
		return
	}
	// Send a POST request.
	account.Email = "email" + GenerateRandNumber(4) + "@email.com"
	user.Id = ""
	user.OrgId = "orgid"
	jsonStr, err := json.Marshal(map[string]interface{}{"user": user, "account": account})
	if err != nil {
		t.Error(err)
		return
	}
	// create new user
	req, err := http.NewRequest("POST", serverURL+userURL+"/user?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping user because already created")
	}
	time.Sleep(time.Second)

	r := user_proto.CreateResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	t.Log(r)

	hdlr := initUserHandler()
	ctx := common.NewTestContext(context.TODO())
	rsp_token := &user_proto.ReadAccountTokenResponse{}
	if err := hdlr.ReadAccountToken(ctx, &user_proto.ReadAccountTokenRequest{AccountId: r.Data.Account.Id}, rsp_token); err != nil {
		t.Error(err)
		return
	}
	token := rsp_token.Token
	t.Log("token:", token)

	if len(token) > 0 {
		if err := ConfirmToken(token, t); err != nil {
			t.Error("confirm failed")
			return
		}
	}
	// login with new user
	sessionId = GetSessionId(account.Email, "pass1", t)
	if len(sessionId) == 0 {
		t.Error("sessionId is invalid")
		return
	}
	req, err = http.NewRequest("GET", serverURL+userURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	e := utils.ErrResponse{}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &e)

	if e.Errors == nil {
		t.Error("Response does not matched")
		return
	}
	t.Log(e)
}

func TestNotEmployeeAuthorized(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	if len(sessionId) == 0 {
		t.Error("sessionId is invalid")
		return
	}
	// Send a POST request.
	account.Email = "email" + GenerateRandNumber(4) + "@email.com"
	user.Id = ""
	jsonStr, err := json.Marshal(map[string]interface{}{"user": user, "account": account})
	if err != nil {
		t.Error(err)
		return
	}
	// create new user
	req, err := http.NewRequest("POST", serverURL+userURL+"/user?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping user because already created")
	}
	time.Sleep(time.Second)

	r := user_proto.CreateResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	rsp_token := &user_proto.ReadAccountTokenResponse{}
	hdlr := initUserHandler()
	ctx := common.NewTestContext(context.TODO())
	if err := hdlr.ReadAccountToken(ctx, &user_proto.ReadAccountTokenRequest{AccountId: r.Data.Account.Id}, rsp_token); err != nil {
		t.Error(err)
		return
	}
	token := rsp_token.Token

	if len(token) > 0 {
		if err := ConfirmToken(token, t); err != nil {
			t.Error("confirm failed")
			return
		}
	}
	// login with new user
	sessionId = GetSessionId(account.Email, "pass1", t)
	if len(sessionId) == 0 {
		t.Error("sessionId is invalid")
		return
	}
	req, err = http.NewRequest("GET", serverURL+userURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	e := utils.ErrResponse{}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &e)

	if e.Errors == nil {
		t.Error("Response does not matched")
		return
	}
	t.Log(e)
}

func TestShareContentWithUser(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
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

	req, err := http.NewRequest("POST", serverURL+userURL+"/user/content/share?session="+sessionId, bytes.NewBuffer(jsonStr))
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
		return
	}
}

func TestReadUserPreference(t *testing.T) {
	initUserDb()
	ctx := common.NewTestContext(context.TODO())

	// create ContentCategoryItem
	createContentCategoryItem(contentCategoryItem, t)

	// create user preference
	if err := db.SaveUserPreference(ctx, &user_proto.Preferences{
		OrgId:       "orgid",
		UserId:      "userid",
		Conditions:  []*static_proto.ContentCategoryItem{contentCategoryItem},
		Allergies:   []*static_proto.ContentCategoryItem{contentCategoryItem},
		Food:        []*static_proto.ContentCategoryItem{contentCategoryItem},
		Cuisines:    []*static_proto.ContentCategoryItem{contentCategoryItem},
		Ethinicties: []*static_proto.ContentCategoryItem{contentCategoryItem},
	}); err != nil {
		t.Error(err)
		return
	}

	// login with new user
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+userURL+"/user/userid/preferences?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := user_proto.ReadUserPreferenceResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Preference == nil {
		t.Errorf("Response does not matched")
		return
	}
	t.Log(r)
}

func TestListUserFeedback(t *testing.T) {
	userapp_db.Init(cl)
	ctx := common.NewTestContext(context.TODO())
	// create user_feedback
	if err := userapp_db.SaveUserFeedback(ctx, &user_proto.UserFeedback{
		UserId:   "userid",
		OrgId:    "orgid",
		Feedback: "hello world",
	}); err != nil {
		t.Error(err)
		return
	}

	// login with new user
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+userURL+"/user/userid/feedback?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := user_proto.ListUserFeedbackResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.Feedbacks) == 0 {
		t.Errorf("Response count does not matched")
		return
	}
	t.Log(r)
}

func TestSearchUser(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"name": "jo"})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("POST", serverURL+userURL+"/user/search?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := user_proto.SearchUserResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}

	if len(r.Data.Response) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestAutocompleteUser(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"name": "jo"})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("POST", serverURL+userURL+"/user/search/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := user_proto.AutocompleteUserResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}

	if len(r.Data.Response) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestGetAccountStatus(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	account.Email = "email" + GenerateRandNumber(4) + "@email.com"
	user.Id = ""
	jsonStr, err := json.Marshal(map[string]interface{}{"user": user, "account": account})
	if err != nil {
		t.Error(err)
		return
	}
	// create user & account
	req, err := http.NewRequest("POST", serverURL+userURL+"/user?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(2 * time.Second)

	r := user_proto.CreateResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	user = r.Data.User

	// set status
	jsonStr, err = json.Marshal(map[string]interface{}{
		"status": account_proto.AccountStatus_ACTIVE,
	})
	if err != nil {
		t.Error(err)
		return
	}

	req, err = http.NewRequest("GET", serverURL+userURL+"/user/"+user.Id+"/account/status?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(2 * time.Second)

	r1 := account_proto.GetAccountStatusResponse{}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r1)

	if r1.Code != http.StatusOK {
		t.Error("Object not update")
	}
	t.Log(r)
}

func TestSetAccountStatus(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	account.Email = "email" + GenerateRandNumber(4) + "@email.com"
	user.Id = ""
	jsonStr, err := json.Marshal(map[string]interface{}{"user": user, "account": account})
	if err != nil {
		t.Error(err)
		return
	}
	// create user & account
	req, err := http.NewRequest("POST", serverURL+userURL+"/user?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(2 * time.Second)

	r := user_proto.CreateResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	user = r.Data.User

	// set status
	jsonStr, err = json.Marshal(map[string]interface{}{
		"status": account_proto.AccountStatus_ACTIVE,
	})
	if err != nil {
		t.Error(err)
		return
	}

	req, err = http.NewRequest("POST", serverURL+userURL+"/user/"+user.Id+"/account/status?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(2 * time.Second)

	r1 := account_proto.SetAccountStatusResponse{}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r1)

	if r1.Code != http.StatusOK {
		t.Error("Object not update")
	}
	t.Log(r)
}

func TestResetUserPassword(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	account.Email = "email" + GenerateRandNumber(4) + "@email.com"
	user.Id = ""
	jsonStr, err := json.Marshal(map[string]interface{}{"user": user, "account": account})
	if err != nil {
		t.Error(err)
		return
	}
	// create user & account
	req, err := http.NewRequest("POST", serverURL+userURL+"/user?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(2 * time.Second)

	r := user_proto.CreateResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	user = r.Data.User

	// set status
	jsonStr, err = json.Marshal(map[string]interface{}{
		"password": "hello",
	})
	if err != nil {
		t.Error(err)
		return
	}

	req, err = http.NewRequest("POST", serverURL+userURL+"/user/"+user.Id+"/account/pass/reset?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(2 * time.Second)

	r1 := account_proto.ResetUserPasswordResponse{}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r1)

	if r1.Code != http.StatusOK {
		t.Error("Object not update")
	}

	t.Log(r)
}

var measurements = []*user_proto.Measurement{
	{
		OrgId:  "orgid",
		Marker: &static_proto.Marker{},
		Method: &static_proto.TrackerMethod{NameSlug: "count"},
		Unit:   "unit_test",
	},
}

func TestAddMultipleMeasurements(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	// create marker
	m := createMarker(marker, t)

	userId, _ := GetUserIdFromSession(sessionId)

	if len(userId) == 0 {
		t.Error("userId error")
	}
	measurements[0].UserId = userId
	// log.Println("marker:", m)
	measurements[0].Marker = m
	measurements[0].Method = m.TrackerMethods[0]

	// create value
	var v google_protobuf1.Value
	raw1 := `"hello world"`
	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &v); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}
	measurements[0].Value = &v

	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"measurements": measurements,
	})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("POST", serverURL+userURL+"/user/"+userId+"/measurements/measurement?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := user_proto.AddMultipleMeasurementsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestGetSharedResources(t *testing.T) {
	// create goal
	goal.Id = ""
	rsp_goal := createGoal(goal, t)
	challenge.Id = ""
	rsp_challenge := createChallenge(challenge, t)
	habit.Id = ""
	rsp_habit := createHabit(habit, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
	}

	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"type": []string{"healum.com/proto/go.micro.srv.behaviour.Goal",
			"healum.com/proto/go.micro.srv.behaviour.Challenge",
			"healum.com/proto/go.micro.srv.behaviour.Habit"},
		"shared_by": []string{userId},
		"status":    []string{"SHARED"},
	})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("POST", serverURL+userURL+"/user/"+userId+"/shared?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := user_proto.GetSharedResourcesResponse{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	// t.Log(r)
	if r.Code != http.StatusOK {
		t.Error("Response is not matched")
		return
	}

	is_goal := false
	is_challenge := false
	is_habit := false

	for _, s := range r.Data.SharedResources {
		t.Log("goal id:", rsp_goal.Id, s.ResourceId)
		if rsp_goal.Id == s.ResourceId {
			is_goal = true
		}
		if rsp_challenge.Id == s.ResourceId {
			is_challenge = true
		}
		if rsp_habit.Id == s.ResourceId {
			is_habit = true
		}
	}

	if !is_goal {
		t.Error("shared goal doesn't exist")
		return
	}
	if !is_challenge {
		t.Error("shared challenge doesn't exist")
		return
	}
	if !is_habit {
		t.Error("shared habit doesn't exist")
		return
	}
}

func TestGetMeasurementsHistory(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
	}
	// create marker
	m := createMarker(marker, t)
	if m == nil {
		t.Error("marker create error")
		return
	}
	measurements[0].UserId = userId
	measurements[0].Marker = m
	measurements[0].Method = m.TrackerMethods[0]
	// create value
	var v google_protobuf1.Value
	raw1 := `"hello world"`
	if err := jsonpb.Unmarshal(strings.NewReader(raw1), &v); err != nil {
		t.Errorf("an unexpected error occurred when parsing into JSONPBUnmarshaler: %v", err)
	}
	measurements[0].Value = &v

	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"measurements": measurements,
	})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("POST", serverURL+userURL+"/user/"+userId+"/measurements/measurement?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)
	r := user_proto.AddMultipleMeasurementsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Code != http.StatusOK {
		t.Log(r)
		t.Error("Add multiple measurements is failed!")
		return
	}

	req, err = http.NewRequest("GET", serverURL+userURL+"/user/"+userId+"/measurements/"+m.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	hr := user_proto.GetMeasurementsHistoryResponse{}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &hr)

	if hr.Data == nil {
		t.Errorf("Response count not matched")
		return
	}
	t.Log(hr.Data.Measurements)
}

func TestGetAllMeasurementsHistory(t *testing.T) {
	TestAddMultipleMeasurements(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	userId, _ := GetUserIdFromSession(sessionId)
	req, err := http.NewRequest("GET", serverURL+userURL+"/user/"+userId+"/measurements/all?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := user_proto.GetAllMeasurementsHistoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if r.Data == nil {
		t.Errorf("Response count not matched")
		return
	}
	t.Log(r.Data.Measurements)
}

func TestGetAllTrackedMarkers(t *testing.T) {
	TestAddMultipleMeasurements(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return
	}
	// measurements[0].UserId = userId

	req, err := http.NewRequest("GET", serverURL+userURL+"/user/markers/"+measurements[0].UserId+"/all?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := user_proto.GetAllTrackedMarkersResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	// t.Error(r.Data.Measurements)

	if len(r.Data.Markers) == 0 {
		t.Errorf("Response count not matched")
		return
	}
}

func TestGetGoalProgress(t *testing.T) {
	// signup goal
	TestSignupToGoal(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+userURL+"/user/"+goal.Users[0].User.Id+"/goals/current/progress?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetGoalProgressResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}
