package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	account_proto "server/account-srv/proto/account"
	behaviour_db "server/behaviour-srv/db"
	behaviour_hdlr "server/behaviour-srv/handler"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	content_db "server/content-srv/db"
	content_hdlr "server/content-srv/handler"
	content_proto "server/content-srv/proto/content"
	kv_proto "server/kv-srv/proto/kv"
	plan_db "server/plan-srv/db"
	plan_hdlr "server/plan-srv/handler"
	plan_proto "server/plan-srv/proto/plan"
	static_proto "server/static-srv/proto/static"
	track_proto "server/track-srv/proto/track"
	"server/user-app-srv/db"
	userapp_proto "server/user-app-srv/proto/userapp"
	user_proto "server/user-srv/proto/user"
	"strconv"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/broker"
	nats_broker "github.com/micro/go-plugins/broker/nats"
)

var userAppURL = "/server/user/app"

var preference = &user_proto.Preferences{
	OrgId:               "orgid",
	CurrentMeasurements: measurements,
	Conditions:          []*static_proto.ContentCategoryItem{contentCategoryItem},
	Allergies:           []*static_proto.ContentCategoryItem{contentCategoryItem},
	Food:                []*static_proto.ContentCategoryItem{contentCategoryItem},
	Cuisines:            []*static_proto.ContentCategoryItem{contentCategoryItem},
	Ethinicties:         []*static_proto.ContentCategoryItem{contentCategoryItem},
}

var user_plan *userapp_proto.UserPlan

func initUserAppDb() {
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)

	plan_db.Init(cl)
	content_db.Init(cl)
	behaviour_db.Init(cl)
}

func GetUserIdFromSession(sessionId string) (string, string) {
	ctx := common.NewTestContext(context.TODO())
	hdlr := initUserHandler()
	rsp_kv, err := hdlr.KvClient.ReadSession(ctx, &kv_proto.ReadSessionRequest{common.SESSION_INDEX, sessionId})
	if err != nil {
		return "", ""
	}
	si := &account_proto.SessionInfo{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&si); err != nil {
		return "", ""
	}
	return si.UserId, si.OrgId
}

func TestCreateBookmark(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	// create content
	createContent(content, t)

	jsonStr, err := json.Marshal(map[string]interface{}{"content_id": content.Id})
	if err != nil {
		t.Error(err)
		return
	}
	// create new user
	req, err := http.NewRequest("POST", serverURL+userAppURL+"/bookmark?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.CreateBookmarkResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
}

func TestReadBookmarkContents(t *testing.T) {
	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	// create new user
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/"+user.Id+"/bookmarks/all?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.ReadBookmarkContentResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestReadBookmarkContentCategorys(t *testing.T) {
	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	// create new user
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/userid/bookmarks/categorys?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.ReadBookmarkContentCategorysResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestReadBookmarkByCategory(t *testing.T) {
	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	// create new user
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/"+user.Id+"/"+content.Category.Id+"/bookmarks?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.ReadBookmarkByCategoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestDeleteBookmark(t *testing.T) {
	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	// create new user
	req, err := http.NewRequest("DELETE", serverURL+userAppURL+"/bookmark/"+"ac0dd6b6-ce29-42b9-8974-01e03da7b171"+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.DeleteBookmarkResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestGetSharedContent(t *testing.T) {
	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/"+user.Id+"/content/shared?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetSharedContentResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestGetSharedPlan(t *testing.T) {
	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/"+user.Id+"/plan/shared?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetSharedPlanResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}
func TestGetSharedSurvey(t *testing.T) {
	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/"+user.Id+"/survey/shared?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetSharedSurveyResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestGetSharedGoal(t *testing.T) {
	createGoal(goal, t)

	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/"+goal.Users[0].User.Id+"/goal/shared?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetSharedGoalResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)

	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestGetSharedChallenge(t *testing.T) {
	createChallenge(challenge, t)

	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/"+challenge.Users[0].User.Id+"/challenge/shared?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetSharedChallengeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestGetSharedHabit(t *testing.T) {
	createHabit(habit, t)

	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/"+habit.Users[0].User.Id+"/habit/shared?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetSharedHabitResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestSignupToGoal(t *testing.T) {
	createGoal(goal, t)

	initUserAppDb()

	// visit share and goal
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"goal_id": goal.Id,
		"user_id": goal.Users[0].User.Id,
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/goal/join?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.SignupToGoalResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
	if r.Data == nil {
		t.Error("Signup is failed")
		return
	}
}

func TestListAllGoal(t *testing.T) {
	TestSignupToGoal(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/goals/joined?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.ListGoalResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestListCurrentGoal(t *testing.T) {
	TestSignupToGoal(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/goals/current?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.ListGoalResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestSignupToChallenge(t *testing.T) {
	createChallenge(challenge, t)

	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"challenge_id": challenge.Id,
		"user_id":      user.Id,
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/challenge/join?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func TestListAllChallenge(t *testing.T) {
	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/challenges/joined?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.ListChallengeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestListCurrentChallenge(t *testing.T) {
	initUserAppDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/challenges/current?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.ListChallengeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestSignupToHabit(t *testing.T) {
	createHabit(habit, t)

	initUserAppDb()

	// visit share and habit

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"habit_id": habit.Id,
		"user_id":  user.Id,
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/habit/join?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

func TestListAllHabit(t *testing.T) {
	TestSignupToHabit(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/habits/joined?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.ListHabitResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestListCurrentHabit(t *testing.T) {
	TestSignupToHabit(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/habits/current?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.ListHabitResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestListMarkers(t *testing.T) {
	TestSignupToGoal(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/current/markers?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.ListMarkersResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestGetPendingSharedActions(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/pending?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetPendingSharedActionsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r.Data.Pendings)
}

func TestGetDefaultMarkerHistory(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/marker/default/history?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetDefaultMarkerHistoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestGetCurrentChallengesWithCount(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/challenge/current/count?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetCurrentChallengesWithCountResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestGetCurrentHabitsWithCount(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/habits/current/count?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetCurrentHabitsWithCountResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r)
}

func TestGetContentCategorys(t *testing.T) {
	createContent(content, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/content/categorys/all?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := content_proto.GetContentCategorysResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Data.Categorys) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestGetContentDetail(t *testing.T) {
	// createContent(content, t)
	// TestCreateBookmark(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/content/"+content.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := content_proto.GetContentDetailResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data.Content == nil {
		t.Error("Object does not matched")
		return
	}

	if !r.Data.Bookmarked {
		t.Error("Bookmarked does not matched")
		return
	}
}

func TestGetContentByCategory(t *testing.T) {
	createContent(content, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/content/category/"+content.Category.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := content_proto.GetContentByCategoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Data.Contents) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestGetFiltersForCategory(t *testing.T) {
	createContentCategoryItem(contentCategoryItem, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/content/category/"+contentCategoryItem.Category.Id+"/filters?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := content_proto.GetFiltersForCategoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Data.ContentCategoryItems) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestFiltersAutocomplete(t *testing.T) {
	createContentCategory(contentCategory, t)

	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"name": "it",
	})
	if err != nil {
		t.Error(err)
		return
	}

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("POST", serverURL+userAppURL+"/content/category/filters/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := content_proto.FiltersAutocompleteResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Data.ContentCategoryItems) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestFilterContentInParticularCategory(t *testing.T) {
	createContent(content, t)

	// Send a POST request.
	tags := []string{}
	for _, tag := range content.Tags {
		tags = append(tags, tag.Id)
	}
	jsonStr, err := json.Marshal(map[string]interface{}{
		"contentCategoryItems": tags,
	})
	if err != nil {
		t.Error(err)
		return
	}

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("POST", serverURL+userAppURL+"/content/category/"+content.Category.Id+"/filter?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := content_proto.FilterContentInParticularCategoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Data.Contents) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestGetUserPreference(t *testing.T) {
	initUserAppDb()
	// create user
	user := CreateUser(user, t)
	if user == nil {
		t.Error("create error")
		return
	}
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/"+user.Id+"/preferences?session="+sessionId, nil)
	log.Info(req)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	//time.Sleep(time.Second)

	r := user_proto.ReadUserPreferenceResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	log.Info(r)
	if r.Data.Preference == nil {
		t.Error("Object does not matched")
		return
	}
}
func TestSaveUserPreference(t *testing.T) {
	// create measurements
	TestAddMultipleMeasurements(t)
	time.Sleep(3 * time.Second)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"preference": preference,
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/preferences?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := user_proto.SaveUserPreferenceResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data.Preference == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestSaveUserDetails(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"firstname":  "first_name",
		"lastname":   "last_name",
		"avatar_url": "image.jpg",
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/details?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.SaveUserDetailsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestGetContentRecommendationByUser(t *testing.T) {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	// get user_id from session id
	ctx := common.NewTestContext(context.TODO())
	hdlr := initUserHandler()
	rsp_kv, err := hdlr.KvClient.ReadSession(ctx, &kv_proto.ReadSessionRequest{common.SESSION_INDEX, sessionId})
	if err != nil {
		t.Error("session reading error")
		return
	}
	si := &account_proto.SessionInfo{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&si); err != nil {
		log.Error("parsing error")
		return
	}

	obj := &content_proto.ContentRecommendation{
		OrgId:   "orgid",
		UserId:  si.UserId,
		Content: content,
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}
	// publish
	if err := nats_brker.Publish(common.CREATE_CONTENT_RECOMMENDATION, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	// get recommendations

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/content/recommendations/all?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := content_proto.GetContentRecommendationByUserResponse{}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Data.Recommendations) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestGetContentRecommendationByCategory(t *testing.T) {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	// get user_id from session id
	ctx := common.NewTestContext(context.TODO())
	hdlr := initUserHandler()
	rsp_kv, err := hdlr.KvClient.ReadSession(ctx, &kv_proto.ReadSessionRequest{common.SESSION_INDEX, sessionId})
	if err != nil {
		t.Error("session reading error")
		return
	}
	si := &account_proto.SessionInfo{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&si); err != nil {
		log.Error("parsing error")
		return
	}

	obj := &content_proto.ContentRecommendation{
		OrgId:   "orgid",
		UserId:  si.UserId,
		Content: content,
	}
	body, err := json.Marshal(obj)
	if err != nil {
		t.Error(err)
		return
	}
	// publish
	if err := nats_brker.Publish(common.CREATE_CONTENT_RECOMMENDATION, &broker.Message{Body: body}); err != nil {
		t.Error(err)
		return
	}
	time.Sleep(2 * time.Second)

	// get recommendations
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/content/recommendations/category/"+content.Category.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := content_proto.GetContentRecommendationByCategoryResponse{}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Data.Recommendations) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestSaveRateForContent(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"rating": 123,
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/content/"+content.Id+"/rating?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.SaveRateForContentResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data.ContentRating == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestDislikeForContent(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/content/"+content.Id+"/dislike?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.DislikeForContentResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data.ContentDislike == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestDislikeForSimilarContent(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"tags": []*static_proto.ContentCategoryItem{contentCategoryItem},
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/content/"+content.Id+"/dislike/similar?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.DislikeForSimilarContentResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data.ContentDislikeSimilar == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestSaveUserFeedback(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"feedback": "Very Nice",
		"rating":   222,
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/feedback?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.SaveUserFeedbackResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data.Feedback == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestJoinUserPlan(t *testing.T) {
	initUserAppDb()

	ctx := common.NewTestContext(context.TODO())
	// create plan
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	plan_hdlr := &plan_hdlr.PlanService{
		Broker:        nats_brker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
	}

	req_plan := &plan_proto.CreateRequest{
		Plan:   plan,
		UserId: "userid",
		OrgId:  "orgid",
	}
	rsp_plan := &plan_proto.CreateResponse{}
	err := plan_hdlr.Create(ctx, req_plan, rsp_plan)
	if err != nil {
		t.Error(err)
		return
	}

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"user_id": "userid",
		"plan_id": plan.Id,
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/plan/join?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.JoinUserPlanResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data.UserPlan == nil {
		t.Error("Object does not matched")
		return
	}

	user_plan = r.Data.UserPlan
}

func TestCreateUserPlan(t *testing.T) {
	initUserAppDb()

	ctx := common.NewTestContext(context.TODO())
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	// create 20 test contents
	content_hdlr := &content_hdlr.ContentService{
		Broker:        nats_brker,
		StaticClient:  static_proto.NewStaticServiceClient("go.micro.srv.static", cl),
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
	}

	for i := 0; i < 20; i++ {
		content.Id = "id_" + strconv.Itoa(i)
		content.Title = "title_" + strconv.Itoa(i)
		content.Category.Id = "category_" + strconv.Itoa(i/2)
		req_plan := &content_proto.CreateContentRequest{
			Content: content,
			UserId:  "userid",
			OrgId:   "orgid",
		}
		rsp_plan := &content_proto.CreateContentResponse{}
		err := content_hdlr.CreateContent(ctx, req_plan, rsp_plan)
		if err != nil {
			t.Error(err)
			return
		}
	}

	// create goal
	behaviour_hdlr := &behaviour_hdlr.BehaviourService{
		Broker:        nats_brker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
	}

	req_goal := &behaviour_proto.CreateGoalRequest{
		Goal:   goal,
		UserId: "userid",
		OrgId:  "orgid",
	}
	rsp_goal := &behaviour_proto.CreateGoalResponse{}
	err := behaviour_hdlr.CreateGoal(ctx, req_goal, rsp_goal)
	if err != nil {
		t.Error(err)
		return
	}

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"user_id":     "userid",
		"goal_id":     goal.Id,
		"days":        4,
		"itemsPerDay": 2,
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/plan/create?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.CreateUserPlanResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data.UserPlan == nil {
		t.Error("Object does not matched")
		return
	}
}

func TestGetUserPlan(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/plan/userid?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetUserPlanResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if r.Data.UserPlan == nil {
		t.Error("Object does not matched")
		return
	}
}
func TestUpdateUserPlan(t *testing.T) {
	TestJoinUserPlan(t)
	time.Sleep(2 * time.Second)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"id":     user_plan.Id,
		"org_id": "orgid123",
		"goals":  user_plan.Goals,
		"days":   user_plan.Days,
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/plan/update?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.UpdateUserPlanResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	// if r.Data.UserPlan == nil {
	// 	t.Error("Object does not matched")
	// 	return
	// }
}

func TestGetPlanItemsCountByCategory(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/plan/"+plan.Id+"/summary/count?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetPlanItemsCountByCategoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Data.ContentCount) == 0 {
		t.Error("Object count does not matched")
		return
	}
}
func TestGetPlanItemsCountByDay(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/plan/"+plan.Id+"/summary/2/count?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetPlanItemsCountByDayResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Data.ContentCount) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestGetPlanItemsCountByCategoryAndDay(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/plan/"+plan.Id+"/summary?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetPlanItemsCountByCategoryAndDayResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	if len(r.Data.ContentCount) == 0 {
		t.Error("Object count does not matched")
		return
	}
}

func TestUserappAllGoals(t *testing.T) {
	initBehaviourDb()

	createGoal(goal, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/goals/all?session="+sessionId+"&org_id="+goal.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := user_proto.AllGoalResponseResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	t.Log(r)

	if r.Data == nil {
		t.Errorf("Object does not matched")
		return
	}

	if len(r.Data.Goals) == 0 {
		t.Errorf("Object count does not matched")
		return
	}
}

func TestUserappAllChallenges(t *testing.T) {
	initBehaviourDb()

	createChallenge(challenge, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/challenges/all?session="+sessionId+"&org_id="+challenge.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := user_proto.AllChallengeResponseResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	t.Log(r)

	if r.Data == nil {
		t.Errorf("Object does not matched")
		return
	}

	if len(r.Data.Challenges) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestUserappAllHabits(t *testing.T) {
	initBehaviourDb()

	createHabit(habit, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/habits/all?session="+sessionId+"&org_id="+habit.OrgId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := user_proto.AllHabitResponseResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	t.Log(r)

	if r.Data == nil {
		t.Errorf("Object does not matched")
		return
	}

	if len(r.Data.Habits) == 0 {
		t.Errorf("Object count does not matched")
	}
}

func TestReceivedItems(t *testing.T) {
	initBehaviourDb()

	createGoal(goal, t)
	t.Log("userid:", goal.Users[0].User.Id)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/"+goal.Users[0].User.Id+"/goal/shared?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := userapp_proto.GetSharedGoalResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)

	}
	json.Unmarshal(body, &r)
	t.Log(r)

	// update shared to receive
	jsonStr, err := json.Marshal(map[string]interface{}{"shared": []*userapp_proto.SharedItem{{
		Type: common.GOAL_TYPE,
		Id:   r.Data.SharedGoals[0].Id,
	}}})
	if err != nil {
		t.Error(err)
		return
	}
	// Send a GET request.
	req, err = http.NewRequest("POST", serverURL+userAppURL+"/"+goal.Users[0].User.Id+"/shared?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r1 := userapp_proto.ReceivedItemsResponse{}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r1)
	t.Log(r1)
}

func TestAutocompleteContentCategoryItem(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	if len(sessionId) == 0 {
		return
	}
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"name": "ti",
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("POST", serverURL+userAppURL+"/content/category/sample_slug/items/autocomplete?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := content_proto.AutocompleteContentCategoryItemResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r.Data.Response)
}

func TestAllContentCategoryItemByName(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"name": "ti",
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/content/category/sample_slug/items/all?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := content_proto.AllContentCategoryItemByNameslugResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r.Data.Response)
}

func TestMarkerByNameslug(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{
		"name_slug": "marker-slug",
	})
	if err != nil {
		t.Error(err)
	}

	req, err := http.NewRequest("GET", serverURL+userAppURL+"/markers/marker-slug/marker?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := static_proto.ReadMarkerResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)

	t.Log(r.Data)
}

func TestGetGoalDetail(t *testing.T) {
	TestSignupToGoal(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	//Send a GET request
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/goal/"+goal.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error("unexpected errror is sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := userapp_proto.ReadGoalResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	t.Log(r)

	if r.Data == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Detail.GoalId != goal.Id {
		t.Errorf("Id does not matched")
		return
	}
}

func TestGetChallengeDetail(t *testing.T) {
	TestSignupToChallenge(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	//Send a GET request
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/challenge/"+challenge.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error("unexpected errror is sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := userapp_proto.ReadChallengeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	t.Log(r)

	if r.Data == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Detail.ChallengeId != challenge.Id {
		t.Errorf("Id does not matched")
		return
	}
}

func TestGetHabitDetail(t *testing.T) {
	TestSignupToHabit(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	//Send a GET request
	req, err := http.NewRequest("GET", serverURL+userAppURL+"/habit/"+habit.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error("unexpected errror is sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := userapp_proto.ReadHabitResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	t.Log(r)

	if r.Data == nil {
		t.Errorf("Object does not matched")
		return
	}
	if r.Data.Detail.HabitId != habit.Id {
		t.Errorf("Id does not matched")
		return
	}
}
