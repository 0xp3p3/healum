package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"server/account-srv/db"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	organisation_proto "server/organisation-srv/proto/organisation"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
	track_proto "server/track-srv/proto/track"
	user_db "server/user-srv/db"
	user_hdlr "server/user-srv/handler"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
	"github.com/pborman/uuid"
	"golang.org/x/crypto/bcrypt"
)

var account_email = &account_proto.Account{
	Email:    "email9@email.com",
	Password: "pass1",
}

var account_phone = &account_proto.Account{
	Phone:    "+447964290506",
	Passcode: "123456",
}

var user = &user_proto.User{
	OrgId:      "orgid",
	Firstname:  "david",
	Lastname:   "john",
	Tags:       []string{"a", "b", "c"},
	Preference: &user_proto.Preferences{},
	AvatarUrl:  "http://example.com",
	ContactDetails: []*user_proto.ContactDetail{
		{Id: "contact_detail_id"},
	},
	Addresses: []*static_proto.Address{{
		PostalCode: "111000",
	}},
	Tokens: []*user_proto.Token{
		{"11671c2e7da30e3c393813f60b327f9c2e2e08390761aa01e37ba5d3e6a617be", 1, "aaa"}, {"token_b", 2, "bbb"},
	},
}

var cl = client.NewClient(
	client.Transport(nats_transport.NewTransport()),
	client.Broker(nats_broker.NewBroker()),
	client.RequestTimeout(3*time.Second),
	client.Retries(5),
)

func initUserHandler() *user_hdlr.UserService {
	cl := client.NewClient(
		client.Transport(nats_transport.NewTransport()),
		client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(3*time.Second),
		client.Retries(5),
	)

	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()

	user_db.Init(cl)

	hdlr := &user_hdlr.UserService{
		Broker:        nats_brker,
		AccountClient: account_proto.NewAccountServiceClient("go.micro.srv.account", cl),
		TrackClient:   track_proto.NewTrackServiceClient("go.micro.srv.track", cl),
		KvClient:      kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		TeamClient:    team_proto.NewTeamServiceClient("go.micro.srv.team", cl),
		StaticClient:  static_proto.NewStaticServiceClient("go.micro.srv.static", cl),
	}
	return hdlr
}

func initHandler() *AccountService {
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)

	nats_brker := nats_broker.NewBroker()
	nats_brker.Init()
	nats_brker.Connect()
	return &AccountService{
		Broker:             nats_brker,
		KvClient:           kv_proto.NewKvServiceClient("go.micro.srv.kv", cl),
		UserClient:         user_proto.NewUserServiceClient("go.micro.srv.user", cl),
		TeamClient:         team_proto.NewTeamServiceClient("go.micro.srv.team", cl),
		OrganisationClient: organisation_proto.NewOrganisationServiceClient("go.micro.srv.organisation", cl),
	}
}

func initDb() {

}

func TestInitDb(t *testing.T) {
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)

	if err := db.Init(cl); err != nil {
		t.Error(err)
		return
	}
}

func createAccountWithPhone(ctx context.Context, hdlr *AccountService, t *testing.T) *user_proto.CreateResponse {
	user_hdlr := initUserHandler()

	account_phone.Phone = GeneratePasscode(8)
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account_email,
	}
	rsp_create := &user_proto.CreateResponse{}
	err := user_hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return rsp_create
}

func createAccountWithEmail(ctx context.Context, hdlr *AccountService, t *testing.T) *user_proto.CreateResponse {
	user_hdlr := initUserHandler()

	account_email.Email = "email" + common.Random(4) + "@email.com"
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account_email,
	}
	rsp_create := &user_proto.CreateResponse{}
	err := user_hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return nil
	}
	return rsp_create
}

func TestUUID(t *testing.T) {
	t.Log(uuid.NewUUID().String())
}

func TestPasscode(t *testing.T) {
	passcode := GeneratePasscode(6)
	t.Log(passcode)
}

func TestPasswordCompare(t *testing.T) {
	// encrypt
	salt := "CHO3VV5p0Y0E84rv" //common.Random(16)
	// h, err := bcrypt.GenerateFromPassword([]byte(x+salt+"pass1"), 10)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	pp := "JDJhJDEwJGt5b3BCZGt0OUx3b2xhUi5VNXh2WU9YMWQxeUJoUlljdHhwMHNEMzJQR2VQVVg2NVJzOFI2" //base64.StdEncoding.EncodeToString(h)

	// decypt
	s, err := base64.StdEncoding.DecodeString(pp)
	if err != nil {
		t.Error(err)
		return
	}

	if err := bcrypt.CompareHashAndPassword(s, []byte(x+salt+"pass2")); err != nil {
		t.Error(err)
		return
	}
}

// run kv-srv
func TestCreate(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	account_email.Email = "email" + common.Random(4) + "@email.com"
	req_create := &account_proto.CreateRequest{account_email}
	rsp_create := &account_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	if len(rsp_create.Token) == 0 {
		t.Error("not created")
		return
	}
}

func TestDoubleCreate(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	account_email.Email = "email" + common.Random(4) + "@email.com"
	req_create := &account_proto.CreateRequest{account_email}
	rsp_create := &account_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	if len(rsp_create.Token) == 0 {
		t.Error("not created")
		return
	}

	account_email.Id = ""
	if err := hdlr.Create(ctx, req_create, rsp_create); err == nil {
		t.Error("Not created account with email")
		return
	}
	t.Log(rsp_create.Account.Created)
}

func TestCreateWithPhone(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_create := &account_proto.CreateRequest{account_phone}
	rsp_create := &account_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	if len(rsp_create.Token) == 0 {
		t.Error("not created")
		return
	}
}

func TestDoubleCreateWithPhone(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_create := &account_proto.CreateRequest{account_phone}
	rsp_create := &account_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	if len(rsp_create.Token) == 0 {
		t.Error("not created")
		return
	}

	err = hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error("Not created account with phone")
		return
	}
}

func TestConfirmRegisterWithEmail(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	account_email.Email = "email" + common.Random(4) + "@email.com"
	req_create := &account_proto.CreateRequest{account_email}
	rsp_create := &account_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	req_confirm := &account_proto.ConfirmRegisterRequest{VerificationToken: rsp_create.Token, Password: "pass1"}
	rsp_confirm := &account_proto.ConfirmRegisterResponse{}
	err = hdlr.ConfirmRegister(ctx, req_confirm, rsp_confirm)
	if err != nil {
		t.Error(err)
		return
	}

	req_read := &account_proto.ReadRequest{&account_proto.Account{Id: rsp_create.Account.Id}}
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

func TestConfirmResend(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_create := &account_proto.CreateRequest{account_email}
	rsp_create := &account_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	if len(rsp_create.Token) == 0 {
		t.Error("Not created with email")
		return
	}

	req_resend := &account_proto.ConfirmResendRequest{Email: account_email.Email}
	rsp_resend := &account_proto.ConfirmResendResponse{}
	err = hdlr.ConfirmResend(ctx, req_resend, rsp_resend)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestConfirmResendWithPhone(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	//fix this to create a user and account dynamically
	// req_create := &account_proto.CreateRequest{account_phone}
	// rsp_create := &account_proto.CreateResponse{}
	// err := hdlr.Create(ctx, req_create, rsp_create)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// if len(rsp_create.Token) == 0 {
	// 	t.Error("Not created with email")
	// 	return
	// }

	req_resend := &account_proto.ConfirmResendRequest{Phone: account_phone.Phone}
	rsp_resend := &account_proto.ConfirmResendResponse{}
	err := hdlr.ConfirmResend(ctx, req_resend, rsp_resend)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestLogin(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	user_hdlr := initUserHandler()

	account_email.Email = "email" + common.Random(4) + "@email.com"
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account_email,
	}
	rsp_create := &user_proto.CreateResponse{}
	err := user_hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	rsp_token := &account_proto.ReadAccountTokenResponse{}
	if err := hdlr.ReadAccountToken(ctx, &account_proto.ReadAccountTokenRequest{AccountId: rsp_create.Data.Account.Id}, rsp_token); err != nil {
		t.Error(err)
		return
	}
	token := rsp_token.Token

	req_confirm := &account_proto.ConfirmRegisterRequest{VerificationToken: token, Password: "pass1"}
	rsp_confirm := &account_proto.ConfirmRegisterResponse{}
	err = hdlr.ConfirmRegister(ctx, req_confirm, rsp_confirm)
	if err != nil {
		t.Error(err)
		return
	}

	req_login := &account_proto.LoginRequest{Email: account_email.Email, Password: "pass1"}
	rsp_login := &account_proto.LoginResponse{}
	err = hdlr.Login(ctx, req_login, rsp_login)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDoubleLogin(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	rsp_create := createAccountWithEmail(ctx, hdlr, t)
	if rsp_create == nil {
		return
	}

	rsp_token := &account_proto.ReadAccountTokenResponse{}
	if err := hdlr.ReadAccountToken(ctx, &account_proto.ReadAccountTokenRequest{AccountId: rsp_create.Data.Account.Id}, rsp_token); err != nil {
		t.Error(err)
		return
	}
	token := rsp_token.Token

	req_confirm := &account_proto.ConfirmRegisterRequest{VerificationToken: token, Password: "pass1"}
	rsp_confirm := &account_proto.ConfirmRegisterResponse{}
	if err := hdlr.ConfirmRegister(ctx, req_confirm, rsp_confirm); err != nil {
		t.Error(err)
		return
	}

	Oneday = time.Second * 6

	req_login := &account_proto.LoginRequest{Email: rsp_create.Data.Account.Email, Password: "pass1"}
	rsp_login := &account_proto.LoginResponse{}
	if err := hdlr.Login(ctx, req_login, rsp_login); err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_login.Data)

	time.Sleep(time.Second * 3)
	req_login1 := &account_proto.LoginRequest{Email: rsp_create.Data.Account.Email, Password: "pass1"}
	rsp_login1 := &account_proto.LoginResponse{}
	if err := hdlr.Login(ctx, req_login1, rsp_login1); err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_login1.Data)
	if rsp_login.Data.Session.Id != rsp_login1.Data.Session.Id {
		t.Error("Session Id is not same.")
		return
	}

	time.Sleep(time.Second * 5)
	req_login2 := &account_proto.LoginRequest{Email: rsp_create.Data.Account.Email, Password: "pass1"}
	rsp_login2 := &account_proto.LoginResponse{}
	if err := hdlr.Login(ctx, req_login2, rsp_login2); err != nil {
		t.Error(err)
		return
	}
	t.Log(rsp_login2.Data)
	if rsp_login1.Data.Session.Id == rsp_login2.Data.Session.Id {
		t.Error("Session Id is not expired.")
		return
	}
}

func TestLoginWithPhone(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	user_hdlr := initUserHandler()

	account_phone.Phone = GeneratePasscode(8)
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account_phone,
	}
	rsp_create := &user_proto.CreateResponse{}
	err := user_hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	rsp_token := &account_proto.ReadAccountTokenResponse{}
	if err := hdlr.ReadAccountToken(ctx, &account_proto.ReadAccountTokenRequest{AccountId: rsp_create.Data.Account.Id}, rsp_token); err != nil {
		t.Error(err)
		return
	}
	token := rsp_token.Token

	req_confirm := &account_proto.ConfirmRegisterRequest{VerificationToken: token, Passcode: "123456"}
	rsp_confirm := &account_proto.ConfirmRegisterResponse{}
	err = hdlr.ConfirmRegister(ctx, req_confirm, rsp_confirm)
	if err != nil {
		t.Error(err)
		return
	}

	req_login := &account_proto.LoginRequest{Phone: account_phone.Phone, Passcode: "123456"}
	rsp_login := &account_proto.LoginResponse{}
	err = hdlr.Login(ctx, req_login, rsp_login)
	if err != nil {
		t.Error(err)
		return
	}

	// t.Error(rsp_login.Data)
	if rsp_login.Data.Session == nil {
		t.Error("Session does not matched")
		return
	}
}

func TestUpdateLockStatusByLoginWithPhone(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	account_phone.Phone = GeneratePasscode(8)
	req_create := &account_proto.CreateRequest{account_phone}
	rsp_create := &account_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	req_confirm := &account_proto.ConfirmRegisterRequest{VerificationToken: rsp_create.Token, Passcode: "123456"}
	rsp_confirm := &account_proto.ConfirmRegisterResponse{}
	err = hdlr.ConfirmRegister(ctx, req_confirm, rsp_confirm)
	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < 5; i++ {
		req_login := &account_proto.LoginRequest{Phone: account_phone.Phone, Passcode: "pass"}
		rsp_login := &account_proto.LoginResponse{}
		hdlr.Login(ctx, req_login, rsp_login)
		// if err != nil {
		// 	t.Error(err)
		// }

		// // t.Error(rsp_login.Data)
		// if rsp_login.Data.Session == nil {
		// 	t.Error("Session does not matched")
		// }

		time.Sleep(time.Second)
	}
}

func TestLogout(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_login := &account_proto.LoginRequest{Email: "email8@email.com", Password: "pass1"}
	rsp_login := &account_proto.LoginResponse{}
	if err := hdlr.Login(ctx, req_login, rsp_login); err != nil {
		t.Error(err)
		return
	}

	req_logout := &account_proto.LogoutRequest{SessionId: rsp_login.Data.Session.Id}
	rsp_logout := &account_proto.LogoutResponse{}
	if err := hdlr.Logout(ctx, req_logout, rsp_logout); err != nil {
		t.Error(err)
		return
	}
}

func TestRecoverPasswordWithEmail(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	account_email.Email = "email" + common.Random(4) + "@email.com"
	req_create := &user_proto.CreateRequest{User: user, Account: account_email}
	rsp_create, err := hdlr.UserClient.Create(ctx, req_create)
	if err != nil {
		t.Error(err)
		return
	}

	rsp_token := &account_proto.ReadAccountTokenResponse{}
	if err := hdlr.ReadAccountToken(ctx, &account_proto.ReadAccountTokenRequest{AccountId: rsp_create.Data.Account.Id}, rsp_token); err != nil {
		t.Error(err)
		return
	}
	token := rsp_token.Token

	req_confirm := &account_proto.ConfirmRegisterRequest{VerificationToken: token, Password: "pass1"}
	rsp_confirm := &account_proto.ConfirmRegisterResponse{}
	err = hdlr.ConfirmRegister(ctx, req_confirm, rsp_confirm)
	if err != nil {
		t.Error(err)
		return
	}

	req_recover := &account_proto.RecoverPasswordRequest{Email: account_email.Email}
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

	req_read := &account_proto.ReadRequest{&account_proto.Account{Id: rsp_create.Data.Account.Id}}
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
func TestRecoverPasswordWithPhone(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	account_phone.Phone = common.Random(8)
	req_create := &user_proto.CreateRequest{User: user, Account: account_phone}
	rsp_create, err := hdlr.UserClient.Create(ctx, req_create)
	if err != nil {
		t.Error(err)
		return
	}
	rsp_token := &account_proto.ReadAccountTokenResponse{}
	if err := hdlr.ReadAccountToken(ctx, &account_proto.ReadAccountTokenRequest{AccountId: rsp_create.Data.Account.Id}, rsp_token); err != nil {
		t.Error(err)
		return
	}
	token := rsp_token.Token

	req_confirm := &account_proto.ConfirmRegisterRequest{VerificationToken: token, Passcode: "123456"}
	rsp_confirm := &account_proto.ConfirmRegisterResponse{}
	err = hdlr.ConfirmRegister(ctx, req_confirm, rsp_confirm)
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

	req_read := &account_proto.ReadRequest{&account_proto.Account{Id: rsp_create.Data.Account.Id}}
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

func TestUpdatePassword(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	user_hdlr := initUserHandler()

	account_email.Email = "email" + common.Random(4) + "@email.com"
	req_create := &user_proto.CreateRequest{
		User:    user,
		Account: account_email,
	}
	rsp_create := &user_proto.CreateResponse{}
	err := user_hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	rsp_token := &account_proto.ReadAccountTokenResponse{}
	if err := hdlr.ReadAccountToken(ctx, &account_proto.ReadAccountTokenRequest{AccountId: rsp_create.Data.Account.Id}, rsp_token); err != nil {
		t.Error(err)
		return
	}
	token := rsp_token.Token

	req_confirm := &account_proto.ConfirmRegisterRequest{VerificationToken: token, Password: "pass1"}
	rsp_confirm := &account_proto.ConfirmRegisterResponse{}
	err = hdlr.ConfirmRegister(ctx, req_confirm, rsp_confirm)
	if err != nil {
		t.Error(err)
		return
	}

	req_recover := &account_proto.RecoverPasswordRequest{Email: account_email.Email}
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

	req_password := &account_proto.UpdatePasswordRequest{PasswordResetToken: rsp_recover.Data.PasswordResetToken, Password: "pass123"}
	rsp_password := &account_proto.UpdatePasswordResponse{}
	err = hdlr.UpdatePassword(ctx, req_password, rsp_password)
	if err != nil {
		t.Error(err)
		return
	}

	req_login := &account_proto.LoginRequest{Email: account_email.Email, Password: "pass123"}
	rsp_login := &account_proto.LoginResponse{}
	err = hdlr.Login(ctx, req_login, rsp_login)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestLock(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	account_email.Email = "email" + common.Random(4) + "@email.com"
	req_create := &account_proto.CreateRequest{account_email}
	rsp_create := &account_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	req_lock := &account_proto.LockRequest{rsp_create.Account.Id}
	rsp_lock := &account_proto.LockResponse{}
	err = hdlr.Lock(ctx, req_lock, rsp_lock)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestConfirmVerify(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	account_email.Email = "email" + common.Random(4) + "@email.com"
	req_create := &account_proto.CreateRequest{account_email}
	rsp_create := &account_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	req_verify := &account_proto.ConfirmVerifyRequest{Token: rsp_create.Token}
	rsp_verify := &account_proto.ConfirmVerifyResponse{}
	err = hdlr.ConfirmVerify(ctx, req_verify, rsp_verify)

	if err != nil {
		t.Error("Verify is failed")
		return
	}
}

func TestConfirmVerifyWithPhone(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	account_phone.Phone = GeneratePasscode(8)
	req_create := &account_proto.CreateRequest{account_phone}
	rsp_create := &account_proto.CreateResponse{}
	err := hdlr.Create(ctx, req_create, rsp_create)
	if err != nil {
		t.Error(err)
		return
	}

	req_verify := &account_proto.ConfirmVerifyRequest{Token: rsp_create.Token}
	rsp_verify := &account_proto.ConfirmVerifyResponse{}
	err = hdlr.ConfirmVerify(ctx, req_verify, rsp_verify)

	if err != nil {
		t.Error("Verify is failed")
		return
	}
}

func TestNoConfirmVerify(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()

	req_verify := &account_proto.ConfirmVerifyRequest{Token: "abcd"}
	rsp_verify := &account_proto.ConfirmVerifyResponse{}
	err := hdlr.ConfirmVerify(ctx, req_verify, rsp_verify)

	if err == nil {
		t.Error("Verify is failed")
		return
	}
}

func TestGetAccountStatus(t *testing.T) {
	initDb()
	ctx := common.NewTestContext(context.TODO())
	hdlr := initHandler()
	hdlr_user := initUserHandler()

	// login user
	req_login := &account_proto.LoginRequest{
		Email:    "email8@email.com",
		Password: "pass1",
	}
	resp_login := &account_proto.LoginResponse{}
	if err := hdlr.Login(ctx, req_login, resp_login); err != nil {
		t.Error("Login is failed")
		return
	}
	rsp_kv, err := hdlr.KvClient.ReadSession(ctx, &kv_proto.ReadSessionRequest{common.SESSION_INDEX, resp_login.Data.Session.Id})
	if err != nil {
		return
	}
	si := &account_proto.SessionInfo{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&si); err != nil {
		return
	}

	// create user & account
	account_email.Email = "email" + common.Random(4) + "@email.com"
	user.Id = ""
	user.OrgId = si.OrgId
	req_create := &user_proto.CreateRequest{
		Account: account_email,
		User:    user,
		OrgId:   si.OrgId,
		TeamId:  si.UserId,
	}
	rsp_create := &user_proto.CreateResponse{}
	if err := hdlr_user.Create(ctx, req_create, rsp_create); err != nil {
		t.Error(err)
		return
	}
	// set status
	req_set := &account_proto.SetAccountStatusRequest{
		UserId: rsp_create.Data.User.Id,
		Status: account_proto.AccountStatus_ACTIVE,
		OrgId:  si.OrgId,
		TeamId: si.UserId,
	}
	rsp_set := &account_proto.SetAccountStatusResponse{}

	if err := hdlr.SetAccountStatus(ctx, req_set, rsp_set); err != nil {
		t.Error(err)
		return
	}

	// get status
	req_status := &account_proto.GetAccountStatusRequest{
		UserId: rsp_create.Data.User.Id,
		OrgId:  si.OrgId,
		TeamId: si.UserId,
	}
	rsp_status := &account_proto.GetAccountStatusResponse{}
	if err := hdlr.GetAccountStatus(ctx, req_status, rsp_status); err != nil {
		t.Error(err)
		return
	}
	if rsp_status.Data.Account.Status != account_proto.AccountStatus_ACTIVE {
		t.Error("get status fail")
		return
	}
}

