package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"
	"server/user-app-srv/db"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

// var serverURL = "http://localhost:8080"
var auditURL = "/server/audits"

var audit = &audit_proto.Audit{
	OrgId:      "orgid",
	ActionName: "test",
	ActionType: audit_proto.ActionType_CREATED,
}

func initAuidtDb() {
	cl := client.NewClient(client.Transport(nats_transport.NewTransport()), client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(5*time.Second),
		client.Retries(5))
	db.Init(cl)
}

func TestFilterAudits(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"action_name": "name1"})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("POST", serverURL+noteURL+"/note?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := audit_proto.FilterAuditsResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if r.Data == nil {
		t.Errorf("Response does not matched")
		return
	}

	if r.Data.Audits[0].ActionName != "name1" {
		t.Errorf("Object does not matched")
		return
	}

}
