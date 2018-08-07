package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	account_proto "server/account-srv/proto/account"
	"server/api/utils"
	"server/common"
	"server/organisation-srv/db"
	organisation_proto "server/organisation-srv/proto/organisation"
	static_proto "server/static-srv/proto/static"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
)

// var serverURL = "http://localhost:8080"
var organisationURL = "/server/organisations"

var organisation = &organisation_proto.Organisation{
	Id:   "orgid",
	Name: "Test Organisation",
}

var profile = &organisation_proto.OrganisationProfile{}
var setting = &organisation_proto.OrganisationSetting{}

func initOrganisationDb() {
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func initHealumDb() {
	ctx := common.NewTestContext(context.TODO())
	db.RemoveDb(ctx, cl)

	db.Init(cl)
}

func createOrganisation(t *testing.T) *organisation_proto.Organisation {
	// create role
	ctx := common.NewTestContext(context.TODO())
	staticClient := static_proto.NewStaticServiceClient("go.micro.srv.static", cl)
	_, err := staticClient.CreateRole(ctx, &static_proto.CreateRoleRequest{
		&static_proto.Role{Name: "admin_role", NameSlug: "admin"},
	})

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	account.Email = "email" + common.Random(4) + "@ex.com"
	organisation.Id = organisation.Id + common.Random(4)
	user.OrgId = organisation.Id
	req_create := &organisation_proto.CreateRequest{
		Organisation: organisation,
		User:         user,
		Account:      account,
		Modules:      []*static_proto.Module{module},
	}

	jsonStr, err := json.Marshal(req_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	req, err := http.NewRequest("POST", serverURL+organisationURL+"/organisation?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return nil
	}

	r := organisation_proto.CreateResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Organisation == nil {
		t.Errorf("Object  does not matched")
		return nil
	}

	return r.Data.Organisation
}

func TestUpdateOrganisation(t *testing.T) {
	organisation := createOrganisation(t)
	if organisation == nil {
		return
	}

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	org_update := &organisation_proto.Organisation{
		Id:   organisation.Id,
		Name: "api_updated_name",
	}
	req_update := &organisation_proto.UpdateRequest{
		Organisation: org_update,
	}

	jsonStr, err := json.Marshal(req_update)
	if err != nil {
		t.Error(err)
		return
	}
	req, err := http.NewRequest("PUT", serverURL+organisationURL+"/organisation?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)
}

// to run this function
// api, org-srv, usr-srv, team-srv, kv-srv, static-srv, account-srv
func TestCreatOrganisationAndConfirmToken(t *testing.T) {
	// initHealumDb()

	ctx := common.NewTestContext(context.TODO())
	// check and create admin_role
	staticClient := static_proto.NewStaticServiceClient("go.micro.srv.static", cl)
	if _, err := staticClient.ReadRoleByNameslug(ctx, &static_proto.ReadRoleByNameslugRequest{"admin"}); err != nil {
		if _, err := staticClient.CreateRole(ctx, &static_proto.CreateRoleRequest{
			&static_proto.Role{Name: "admin_role", NameSlug: "admin"},
		}); err != nil {
			t.Error(err)
			return
		}
	}
	// create r equest body for org
	req_create := &organisation_proto.CreateRequest{
		Organisation: organisation,
		User:         user,
		Account:      account,
		Modules:      []*static_proto.Module{module},
	}
	jsonStr, err := json.Marshal(req_create)
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+organisationURL+"/organisation", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(2 * time.Second)
	// parsing response body
	r := organisation_proto.CreateResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Error(r)
		return
	}
	if r.Data.Organisation == nil {
		t.Errorf("Object does not matched")
		return
	}

	accountClient := account_proto.NewAccountServiceClient("go.micro.srv.account", cl)
	if _, err := accountClient.InternalConfirm(ctx, &account_proto.InternalConfirmRequest{
		AccountId: r.Data.Account.Id,
		Password:  "pass1",
	}); err != nil {
		t.Error(err)
	}

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	if len(sessionId) == 0 {
		t.Error("Session does not matched")
		return
	}
}

func TestAllOrganisations(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+organisationURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := organisation_proto.AllResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		t.Error(r)
		t.Errorf("Object count does not matched")
	}
}

func TestReadOrganisation(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+organisationURL+"/organisation/"+organisation.Id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := organisation_proto.ReadResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Organisation == nil {
		t.Errorf("Object  does not matched")
	}
}

// func TestUpdateModules(t *testing.T) {
// 	sessionId := GetSessionId("email8@email.com", "pass1", t)
// 	// Send a GET request.
// 	req, err := http.NewRequest("POST", serverURL+organisationURL+"/modules?session="+sessionId, nil)
// 	req.Header.Set("Content-Type", restful.MIME_JSON)
// 	common.SetTestHeader(req.Header)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		t.Errorf("unexpected error in sending req: %v", err)
// 	}
// 	time.Sleep(time.Second)

// 	r := organisation_proto.UpdateModulesResponse{}
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	json.Unmarshal(body, &r)
// 	if r.Data == nil {
// 		t.Errorf("Object  does not matched")
// 	}
// }

func TestCreateOrganisationProfile(t *testing.T) {
	req_create := &organisation_proto.CreateOrganisationProfileRequest{
		Profile: profile,
	}
	jsonStr, err := json.Marshal(req_create)
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	req, err := http.NewRequest("POST", serverURL+organisationURL+"/profile?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

	}
	time.Sleep(time.Second)

	r := organisation_proto.CreateOrganisationProfileResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Profile == nil {
		t.Errorf("Object  does not matched")
	}
}

func TestCreateOrganisationSetting(t *testing.T) {
	req_create := &organisation_proto.CreateOrganisationSettingRequest{
		Setting: setting,
	}
	jsonStr, err := json.Marshal(req_create)
	if err != nil {
		t.Error(err)
	}
	log.Println(string(jsonStr))

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	req, err := http.NewRequest("POST", serverURL+organisationURL+"/setting?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping plan because already created")

	}
	time.Sleep(time.Second)

	r := organisation_proto.CreateOrganisationSettingResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	json.Unmarshal(body, &r)
	if r.Data.Setting == nil {
		t.Errorf("Object  does not matched")
	}
}

func TestInvalidSession(t *testing.T) {
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+organisationURL+"/organisation/"+organisation.Id, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Code == http.StatusOK {
		t.Errorf("Response does not matched")
		return
	}
	t.Log("ok:", r)
}

func TestNotAuthorized(t *testing.T) {
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+organisationURL+"/organisation/"+organisation.Id+"?session=anystring", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Code == http.StatusOK {
		t.Errorf("Response does not matched")
		return
	}
	t.Log("ok:", r)
}
