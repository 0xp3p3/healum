package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server/common"
	"server/track-srv/db"
	track_proto "server/track-srv/proto/track"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

var trackURL = "/server/track"

func initTrackDb() {
	cl := client.NewClient(
		client.Transport(nats_transport.NewTransport()),
		client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(5*time.Second),
		client.Retries(5))
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

var user1 = &user_proto.User{
	Id:        "userid",
	OrgId:     "orgid",
	Firstname: "david",
	Lastname:  "john",
	AvatarUrl: "http://example.com",
	Tokens: map[string]*user_proto.Token{
		"a1": {"11671c2e7da30e3c393813f60b327f9c2e2e08390761aa01e37ba5d3e6a617be", 1, "a1", "aaa"},
		"b1": {"token_b", 2, "b1", "bbb"},
	},
}

func TestCreateTrackGoal(t *testing.T) {
	createGoal(goal, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)

	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return
	}
	user1.Id = userId
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"user": user1})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+trackURL+"/goal/"+goal.Id+"?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.CreateTrackGoalResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.TrackGoal == nil {
		t.Errorf("Object  does not matched")
	}
}

func TestGetGoalCount(t *testing.T) {
	TestCreateTrackGoal(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+trackURL+"/goal/"+goal.Id+"/count?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetGoalCountResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Count == 0 {
		t.Errorf("Count  does not matched")
	}

	t.Log(r.Data.Count)
}

// flushdb before test
// select 10
// flushdb
func TestSetGoalCount(t *testing.T) {
	// initTrackDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+trackURL+"/goal/"+goal.Id+"/count?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetGoalCountResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	t.Log(r)
	if r.Data.Count == 0 {
		t.Errorf("Count  does not matched")
		return
	}

	t.Log(r.Data.Count)
}

func TestGetGoalHistory(t *testing.T) {
	TestCreateTrackGoal(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+trackURL+"/goal/"+goal.Id+"/history?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetGoalHistoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.TrackGoals) == 0 {
		t.Errorf("Count  does not matched")
	}
}

func TestCreateTrackChallenge(t *testing.T) {
	createChallenge(challenge, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return
	}
	user1.Id = userId
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"user": user1})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+trackURL+"/challenge/"+challenge.Id+"?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.CreateTrackChallengeResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.TrackChallenge == nil {
		t.Errorf("Object  does not matched")
	}
}

func TestGetChallengeCount(t *testing.T) {
	TestCreateTrackChallenge(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+trackURL+"/challenge/"+challenge.Id+"/count?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetChallengeCountResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.Count == 0 {
		t.Errorf("Count  does not matched")
		return
	}

	t.Log(r.Data.Count)
}

// flushdb before test
// select 10
// flushdb
func TestSetChallengeCount(t *testing.T) {
	// initTrackDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+trackURL+"/challenge/"+challenge.Id+"/count?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetChallengeCountResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	t.Log(r)
	if r.Data.Count == 0 {
		t.Errorf("Count  does not matched")
	}

	t.Log(r.Data.Count)
}

func TestGetChallengeHistory(t *testing.T) {
	TestCreateTrackChallenge(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+trackURL+"/challenge/"+challenge.Id+"/history?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetChallengeHistoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.TrackChallenges) == 0 {
		t.Errorf("Count  does not matched")
	}
}

func TestCreateTrackHabit(t *testing.T) {
	createHabit(habit, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return
	}
	user1.Id = userId
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"user": user1})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+trackURL+"/habit/"+habit.Id+"?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.CreateTrackHabitResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.TrackHabit == nil {
		t.Errorf("Object  does not matched")
	}
}

func TestGetHabitCount(t *testing.T) {
	TestCreateTrackHabit(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+trackURL+"/habit/"+habit.Id+"/count?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetHabitCountResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.Count == 0 {
		t.Errorf("Count  does not matched")
		return
	}

	t.Log(r.Data.Count)
}

// flushdb before test
// select 10
// flushdb
func TestSetHabitCount(t *testing.T) {
	// initTrackDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+trackURL+"/habit/"+habit.Id+"/count?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetHabitCountResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	t.Log(r)
	if r.Data.Count == 0 {
		t.Errorf("Count  does not matched")
	}

	t.Log(r.Data.Count)
}

func TestGetHabitHistory(t *testing.T) {
	TestCreateTrackHabit(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+trackURL+"/habit/"+habit.Id+"/history?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetHabitHistoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.TrackHabits) == 0 {
		t.Errorf("Count  does not matched")
	}
}

func TestCreateTrackContent(t *testing.T) {
	createContent(content, t)
	createContentCategoryItem(contentCategoryItem, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	userId, _ := GetUserIdFromSession(sessionId)
	if len(userId) == 0 {
		t.Error("userId error")
		return
	}
	user1.Id = userId
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"user": user1})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+trackURL+"/content/"+content.Id+"?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.CreateTrackContentResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.TrackContent == nil {
		t.Errorf("Object  does not matched")
		return
	}
}

func TestGetContentCount(t *testing.T) {
	TestCreateTrackContent(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+trackURL+"/content/"+content.Id+"/count?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetContentCountResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.Count == 0 {
		t.Errorf("Count  does not matched")
		return
	}

	t.Log(r.Data.Count)
}

// flushdb before test
// select 10
// flushdb
func TestSetContentCount(t *testing.T) {
	// initTrackDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+trackURL+"/content/"+content.Id+"/count?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetContentCountResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	t.Log(r)
	if r.Data.Count == 0 {
		t.Errorf("Count  does not matched")
	}

	t.Log(r.Data.Count)
}

func TestGetContentHistory(t *testing.T) {
	TestCreateTrackContent(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+trackURL+"/content/"+content.Id+"/history?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetContentHistoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.TrackContents) == 0 {
		t.Errorf("Count  does not matched")
	}
}

func TestCreateTrackMarker(t *testing.T) {
	createTrackerMethod(trackerMethod, t)
	createContent(content, t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"value": 3, "unit": "sample"})
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	// Send a POST request
	req, err := http.NewRequest("POST", serverURL+trackURL+"/marker/"+marker.Id+"?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.CreateTrackMarkerResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.TrackMarker == nil {
		t.Errorf("Object does not matched")
	}
}

// to run, behaviour-srv, static-srv, content-srv
func TestGetLastMarker(t *testing.T) {
	TestCreateTrackMarker(t)

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+trackURL+"/marker/"+marker.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(body))
}

func TestGetMarkerHistory(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+trackURL+"/marker/"+marker.Id+"/history?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetMarkerHistoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.TrackMarkers) == 0 {
		t.Errorf("Object count does not matched")
	}
	t.Log(r)
}

func TestGetAllMarkerHistory(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("GET", serverURL+trackURL+"/marker/history/all?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := track_proto.GetAllMarkerHistoryResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if len(r.Data.TrackMarkers) == 0 {
		t.Errorf("Object count does not matched")
	}
	t.Log(r)
}
