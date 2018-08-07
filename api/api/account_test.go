package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"server/account-srv/db"
	"server/account-srv/handler"
	account_proto "server/account-srv/proto/account"
	"server/api/utils"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	organisation_proto "server/organisation-srv/proto/organisation"
	team_proto "server/team-srv/proto/team"
	user_db "server/user-srv/db"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

var accountURL = "/server/account"

var account_email = &account_proto.Account{
	Email:    "email" + GenerateRand(3) + "@email.com",
	Password: "pass1",
}

var account_phone = &account_proto.Account{
	Phone:    "+8613042431402",
	Passcode: "123456",
}

var cl = client.NewClient(client.Transport(
	nats_transport.NewTransport()),
	client.Broker(nats_broker.NewBroker()),
	client.RequestTimeout(5*time.Second),
	client.Retries(5))

func initAccountDb() {
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	user_db.Init(cl)
	db.Init(cl)
}

func initAccountHandler() *handler.AccountService {
	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	hdlr := &handler.AccountService{
		Broker:             nats_brker,
		KvClient:           kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		UserClient:         user_proto.NewUserServiceClient("go.micro.srv.user", cl),
		TeamClient:         team_proto.NewTeamServiceClient("go.micro.srv.team", cl),
		OrganisationClient: organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl),
	}
	return hdlr
}

func GenerateRand(n int) string {
	letters := []rune("1234567890")
	rand.Seed(time.Now().UTC().UnixNano())
	randomString := make([]rune, n)
	for i := range randomString {
		randomString[i] = letters[rand.Intn(len(letters))]
	}
	return string(randomString)
}

func GetSessionId(email, pass string, t *testing.T) string {
	var jsonStr = []byte(fmt.Sprintf(`{"email": "%v", "password": "%v"}`, email, pass))
	req, err := http.NewRequest("POST", serverURL+accountURL+"/login", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return ""
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected response: %v, expected: %v", resp.StatusCode, http.StatusOK)
		body, _ := ioutil.ReadAll(resp.Body)
		t.Errorf("error body: %v", string(body))
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf(err.Error())
		return ""
	}
	var s = new(account_proto.LoginResponse)
	err = json.Unmarshal(body, &s)
	if err != nil {
		t.Errorf("error body: %v", string(body))
		return ""
	}
	time.Sleep(2 * time.Second)
	return s.Data.Session.Id
}

func CreateWithEmail(ctx context.Context, hdlr *handler.AccountService, t *testing.T) (*account_proto.Account, string) {
	account_email.Email = "email" + common.Random(4) + "@email.com"
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account_email,
	}
	rsp_create, err := hdlr.UserClient.Create(ctx, req_create)
	if err != nil {
		return nil, ""
	}

	rsp_token := &account_proto.ReadAccountTokenResponse{}
	if err := hdlr.ReadAccountToken(ctx, &account_proto.ReadAccountTokenRequest{AccountId: rsp_create.Data.Account.Id}, rsp_token); err != nil {
		t.Error(err)
		return nil, ""
	}
	token := rsp_token.Token

	return rsp_create.Data.Account, token

}

func CreateWithPhone(ctx context.Context, hdlr *handler.AccountService, t *testing.T) (*account_proto.Account, string) {
	account_phone.Phone = common.Random(8)
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account_phone,
	}
	rsp_create, err := hdlr.UserClient.Create(ctx, req_create)
	if err != nil {
		return nil, ""
	}

	rsp_token := &account_proto.ReadAccountTokenResponse{}
	if err := hdlr.ReadAccountToken(ctx, &account_proto.ReadAccountTokenRequest{AccountId: rsp_create.Data.Account.Id}, rsp_token); err != nil {
		t.Error(err)
		return nil, ""
	}
	token := rsp_token.Token

	return rsp_create.Data.Account, token

}

func TestGetSessionId(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	t.Log(sessionId)
}

func ConfirmToken(token string, t *testing.T) error {
	jsonStr, err := json.Marshal(map[string]interface{}{"verification_token": token})
	if err != nil {
		t.Error(err)
		return err
	}

	req, err := http.NewRequest("POST", serverURL+accountURL+"/confirm", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return err
	}
	time.Sleep(time.Second)

	// parsing response body
	r := account_proto.ConfirmRegisterResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return err
	}
	json.Unmarshal(body, &r)
	if r.Code != http.StatusOK {
		t.Errorf("Status does not matched")
		return err
	}

	return nil
}

func TestConfirmRegisterWithEmail(t *testing.T) {
	initAccountDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initAccountHandler()

	account_email.Email = "email" + common.Random(4) + "@email.com"
	account, token := CreateWithEmail(ctx, hdlr, t)
	jsonStr, err := json.Marshal(map[string]interface{}{"verification_token": token})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("POST", serverURL+accountURL+"/confirm", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(2 * time.Second)
	r := account_proto.ConfirmRegisterResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	t.Log(r)

	req_read := &account_proto.ReadRequest{&account_proto.Account{Id: account.Id}}
	rsp_read := &account_proto.ReadResponse{}
	err = hdlr.Read(ctx, req_read, rsp_read)
	if err != nil {
		t.Error(err)
		return
	}

	if rsp_read.Data.Account.Status != account_proto.AccountStatus_ACTIVE {
		t.Error("Status does not matched")
		return
	}

	if !rsp_read.Data.Account.Confirmed {
		t.Error("Confirmed does not matched")
		return
	}

}

func TestErrorConfirmRegisterWithEmail(t *testing.T) {
	initAccountDb()

	jsonStr, err := json.Marshal(map[string]interface{}{"verification_token": "token"})
	if err != nil {
		t.Error(err)
		return
	}
	req, err := http.NewRequest("POST", serverURL+accountURL+"/confirm", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(2 * time.Second)
	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	t.Log(r, r.Errors[0])
}

func TestConfirmResend(t *testing.T) {
	initAccountDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initAccountHandler()

	account_email.Email = "email" + common.Random(4) + "@email.com"
	account, _ := CreateWithEmail(ctx, hdlr, t)
	jsonStr, err := json.Marshal(map[string]interface{}{"email": account.Email})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("POST", serverURL+accountURL+"/confirm/resend", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(2 * time.Second)

	r := account_proto.ConfirmResendResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	t.Log(r)

	if r.Code != http.StatusOK {
		t.Error(err)
		return
	}
}

func TestConfirmResendWithPhone(t *testing.T) {
	initAccountDb()

	// create account
	ctx := common.NewTestContext(context.TODO())
	hdlr := initAccountHandler()
	account_phone.Phone = GenerateRandNumber(8)
	account, _ := CreateWithPhone(ctx, hdlr, t)
	if account == nil {
		return
	}
	jsonStr, err := json.Marshal(map[string]interface{}{"phone": account.Phone})
	if err != nil {
		t.Error(err)
		return
	}

	req, err := http.NewRequest("POST", serverURL+accountURL+"/confirm/resend", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(2 * time.Second)

	r := account_proto.ConfirmResendResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	t.Log(r)

	if r.Code != http.StatusOK {
		t.Error(err)
		return
	}
}

func TestLoginWithEmail(t *testing.T) {
	initAccountDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initAccountHandler()

	account_email.Email = "email" + common.Random(4) + "@email.com"
	_, token := CreateWithEmail(ctx, hdlr, t)

	jsonStr, err := json.Marshal(map[string]interface{}{"verification_token": token, "password": "pass1"})
	if err != nil {
		t.Error(err)
		return
	}
	// confirm user
	req, err := http.NewRequest("POST", serverURL+accountURL+"/confirm", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	// login
	jsonStr, err = json.Marshal(map[string]interface{}{"email": account_email.Email, "password": "pass1"})
	if err != nil {
		t.Error(err)
		return
	}
	req, err = http.NewRequest("POST", serverURL+accountURL+"/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := account_proto.LoginResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	t.Log(r)

	if r.Data.Session == nil {
		t.Errorf("Session does not matched")
		return
	}
}

func TestLogoutWithPhone(t *testing.T) {
	initAccountDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initAccountHandler()

	// create user
	account_phone.Phone = GenerateRandNumber(8)
	_, token := CreateWithPhone(ctx, hdlr, t)
	jsonStr, err := json.Marshal(map[string]interface{}{"verification_token": token, "passcode": "123456"})
	if err != nil {
		t.Error(err)
		return
	}

	// confirm user
	req, err := http.NewRequest("POST", serverURL+accountURL+"/confirm", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	jsonStr, err = json.Marshal(map[string]interface{}{"phone": account_phone.Phone, "passcode": "123456"})
	if err != nil {
		t.Error(err)
		return
	}

	// login user
	req, err = http.NewRequest("POST", serverURL+accountURL+"/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := account_proto.LoginResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if r.Data.Session == nil {
		t.Errorf("Session does not matched")
	}

	// logout user
	req, err = http.NewRequest("GET", serverURL+accountURL+"/logout?session="+r.Data.Session.Id, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)
}

func TestRecoverPassword(t *testing.T) {
	initAccountDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initAccountHandler()

	// create user
	account_phone.Phone = GenerateRandNumber(8)
	account, token := CreateWithPhone(ctx, hdlr, t)
	t.Log("account:", account, token)

	req_confirm := &account_proto.ConfirmRegisterRequest{VerificationToken: token, Passcode: "123456"}
	rsp_confirm := &account_proto.ConfirmRegisterResponse{}

	err := hdlr.ConfirmRegister(ctx, req_confirm, rsp_confirm)
	if err != nil {
		t.Error(err)
		return
	}

	jsonStr, err := json.Marshal(map[string]interface{}{"phone": account.Phone})
	if err != nil {
		t.Error(err)
		return
	}

	// recover password
	req, err := http.NewRequest("POST", serverURL+accountURL+"/pass/recover", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := account_proto.RecoverPasswordResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	t.Log(r)
	if len(r.Data.PasswordResetToken) == 0 {
		t.Errorf("Token does not matched")
		return
	}
}

func TestUpdatePassword(t *testing.T) {
	initAccountDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initAccountHandler()

	// create user
	account_phone.Phone = GenerateRandNumber(8)
	_, token := CreateWithPhone(ctx, hdlr, t)

	req_confirm := &account_proto.ConfirmRegisterRequest{VerificationToken: token, Passcode: "123456"}
	rsp_confirm := &account_proto.ConfirmRegisterResponse{}
	err := hdlr.ConfirmRegister(ctx, req_confirm, rsp_confirm)
	if err != nil {
		t.Error(err)
		return
	}

	req_recover := &account_proto.RecoverPasswordRequest{Phone: account_phone.Phone}
	rsp_recover := &account_proto.RecoverPasswordResponse{}
	err = hdlr.RecoverPassword(ctx, req_recover, rsp_recover)
	if err != nil {
		t.Error(err)
		return
	}
	if len(rsp_recover.Data.PasswordResetToken) == 0 {
		t.Error("Token does not matched")
		return
	}

	jsonStr, err := json.Marshal(map[string]interface{}{"password_reset_token": rsp_recover.Data.PasswordResetToken, "passcode": "654321"})
	if err != nil {
		t.Error(err)
		return
	}

	// update password
	req, err := http.NewRequest("POST", serverURL+accountURL+"/pass/update", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	req_login := &account_proto.LoginRequest{Phone: account_phone.Phone, Passcode: "654321"}
	rsp_login := &account_proto.LoginResponse{}
	err = hdlr.Login(ctx, req_login, rsp_login)
	if err != nil {
		t.Error(err)
		return
	}

	if rsp_login.Data.Session == nil {
		t.Error("Session does not matched")
		return
	}
}

func TestPassVerify(t *testing.T) {
	initAccountDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initAccountHandler()

	// create user
	account_phone.Phone = GenerateRandNumber(8)
	account, token := CreateWithPhone(ctx, hdlr, t)

	req_confirm := &account_proto.ConfirmRegisterRequest{VerificationToken: token, Passcode: "123456"}
	rsp_confirm := &account_proto.ConfirmRegisterResponse{}
	err := hdlr.ConfirmRegister(ctx, req_confirm, rsp_confirm)
	if err != nil {
		t.Error(err)
		return
	}

	jsonStr, err := json.Marshal(map[string]interface{}{"phone": account.Phone})
	if err != nil {
		t.Error(err)
		return
	}
	// pass verify
	req, err := http.NewRequest("POST", serverURL+accountURL+"/pass/recover", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r := account_proto.RecoverPasswordResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)
	if len(r.Data.PasswordResetToken) == 0 {
		t.Errorf("Token does not matched")
		return
	}

	jsonStr1, err := json.Marshal(map[string]interface{}{"token": r.Data.PasswordResetToken})
	if err != nil {
		t.Error(err)
		return
	}
	// pass verify
	req, err = http.NewRequest("POST", serverURL+accountURL+"/pass/verify", bytes.NewBuffer(jsonStr1))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r1 := account_proto.ConfirmVerifyResponse{}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r1)
	if r1.Code != http.StatusOK {
		t.Errorf("Token does not matched")
		return
	}
}

func TestNotPassVerify(t *testing.T) {
	initAccountDb()

	jsonStr1, err := json.Marshal(map[string]interface{}{"token": "abcd"})
	if err != nil {
		t.Error(err)
		return
	}
	// pass verify
	req, err := http.NewRequest("POST", serverURL+accountURL+"/pass/verify", bytes.NewBuffer(jsonStr1))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	time.Sleep(time.Second)

	r1 := account_proto.ConfirmVerifyResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r1)
	if r1.Code == http.StatusOK {
		t.Errorf("Token does not matched")
		return
	}
}
