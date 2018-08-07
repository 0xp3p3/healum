package handler

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"server/account-srv/db"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	organisation_proto "server/organisation-srv/proto/organisation"
	sms_proto "server/sms-srv/proto/sms"
	team_proto "server/team-srv/proto/team"
	user_proto "server/user-srv/proto/user"
	"strings"
	"time"

	"github.com/micro/go-micro/broker"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	x      = "cruft123"
	digis  = "1234567890"
	Oneday = 24 * time.Hour
)

type AccountService struct {
	Broker             broker.Broker
	KvClient           kv_proto.KvServiceClient
	UserClient         user_proto.UserServiceClient
	TeamClient         team_proto.TeamServiceClient
	OrganisationClient organisation_proto.OrganisationServiceClient
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// GeneratePasscode return random digit string by n
func GeneratePasscode(n int) string {
	letters := []rune(digis)
	rand.Seed(time.Now().UTC().UnixNano())
	randomString := make([]rune, n)
	for i := range randomString {
		randomString[i] = letters[rand.Intn(len(letters))]
	}
	return string(randomString)
}

func (p *AccountService) Init() {

}

//removes old (unexpired tokens), generates new tokens for an account and stores them in REDIS
func (p *AccountService) PutVerification(ctx context.Context, account_id string, mode user_proto.ContactDetailType) (string, error) {
	log.Info("PutVerification token request received")

	// read token with account_id to check whether there are any existing tokens
	req_get := &kv_proto.GetExRequest{Index: common.VERIFICATION_TOKEN_INDEX, Key: account_id}
	rsp_get, err := p.KvClient.GetEx(ctx, req_get)

	if rsp_get != nil && err == nil {
		log.Info("Deleting existing tokens for account: ", account_id)
		va := account_proto.VerificationAccount{}
		if err := json.Unmarshal(rsp_get.Item.Value, &va); err != nil {
			common.InternalServerError(common.AccountSrv, p.PutVerification, nil, "Marshalling error")
			return "", err
		}

		//remove confirmation token from redis
		req_remove := &kv_proto.RemoveTokenRequest{AccountId: account_id, Token: va.Token}
		if _, err := p.KvClient.RemoveToken(ctx, req_remove); err != nil {
			common.ErrorLog(common.AccountSrv, common.GetFunctionName(p.PutVerification), err, "RemoveToken is failed")
			return "", err
		}
	}

	//No existing verification tokens were found
	log.Info("Generate new token for account: ", account_id)

	var token string
	// create token string
	if mode == user_proto.ContactDetailType_EMAIL {
		var err error
		token, err = GenerateRandomString(32)
		if err != nil {
			common.InternalServerError(common.AccountSrv, p.PutVerification, nil, "Token generation error")
			return "", err
		}
	} else if mode == user_proto.ContactDetailType_PHONE {
		token = GeneratePasscode(6)
	}

	// create verification token on redis db
	vt := &account_proto.VerificationToken{
		AccountId: account_id,
		ExpiresAt: time.Now().Add(Oneday).Unix(),
	}
	body, err := json.Marshal(vt)
	if err != nil {
		common.InternalServerError(common.AccountSrv, p.PutVerification, nil, "Marshalling error")
		return "", err
	}
	req_kv := &kv_proto.PutExRequest{
		Index: common.VERIFICATION_TOKEN_INDEX,
		Item: &kv_proto.Item{
			Key:        token,
			Value:      body,
			Expiration: int64(Oneday),
		},
	}
	if _, err := p.KvClient.PutEx(ctx, req_kv); err != nil {
		common.InternalServerError(common.AccountSrv, p.PutVerification, nil, "PutVerification token saving is failed")
		return "", err
	}

	// create accountid key, token value on redis db
	va := &account_proto.VerificationAccount{
		Token:     token,
		ExpiresAt: time.Now().Add(Oneday).Unix(),
	}
	body, err = json.Marshal(va)
	if err != nil {
		common.InternalServerError(common.AccountSrv, p.PutVerification, nil, "Marshalling error")
		return "", err
	}
	req_kv = &kv_proto.PutExRequest{
		Index: common.VERIFICATION_TOKEN_INDEX,
		Item: &kv_proto.Item{
			Key:        account_id,
			Value:      body,
			Expiration: int64(Oneday),
		},
	}
	if _, err := p.KvClient.PutEx(ctx, req_kv); err != nil {
		common.InternalServerError(common.AccountSrv, p.PutVerification, nil, "PutVerification account saving is failed")
		return "", err
	}

	return token, nil
}

func (p *AccountService) AuthFailed(ctx context.Context, account_id string) error {
	rsp, err := p.KvClient.AuthFailed(ctx, &kv_proto.AuthFailedRequest{common.AUTHENTIFICATION_INDEX, account_id})
	if err != nil {
		return err
	}
	if rsp.Failed >= 4 {
		if err := p.Lock(ctx, &account_proto.LockRequest{AccountId: account_id}, &account_proto.LockResponse{}); err != nil {
			return err
		}
	}

	return nil
}

func (p *AccountService) Read(ctx context.Context, req *account_proto.ReadRequest, rsp *account_proto.ReadResponse) error {
	log.Info("Received Account.Read request")
	account, err := db.Read(ctx, req.Account)
	if err != nil {
		return err
	}
	rsp.Data = &account_proto.AccountData{account}
	return nil
}

func (p *AccountService) Create(ctx context.Context, req *account_proto.CreateRequest, rsp *account_proto.CreateResponse) error {
	log.Info("Received Account.Create request")
	//TODO: refactor de-duplication to a separate internal method
	// checking duplicate account before create
	valid, err := db.DedupAccount(ctx, req.Account)
	if valid != nil && err == nil {
		if valid.Confirmed && valid.Status != account_proto.AccountStatus_INACTIVE {
			return common.InternalServerError(common.AccountSrv, p.Create, err, "Account is already existed")
		} else {
			//resend a new confirm verification token for a new existing account which is still INACTIVE
			log.Info("Responding with new token and existing account information for inactive account: ", valid.Id)

			// generate and save account token
			token, _, err := p.GenerateAndSaveToken(ctx, valid)
			if err != nil {
				return common.InternalServerError(common.AccountSrv, p.Create, err, "GenerateAndSaveToken is failed")
			}
			rsp.Token = token   // new token generated for existing INACTIVE/LOCKED account
			rsp.Account = valid // existing account
			return nil
		}
	}
	log.Info("Creating a new account")
	// create account
	if len(req.Account.Id) == 0 {
		req.Account.Id = uuid.NewUUID().String()
	}

	//FIXME:request shouldn't be coming with password or passcode. Remove it when request is fixed
	req.Account.Password = ""
	req.Account.Passcode = ""

	// create account in arango db
	if err := db.Create(ctx, req.Account); err != nil {
		return common.InternalServerError(common.AccountSrv, p.Create, err, "Create query is failed")
	}

	// generate and save account token
	token, _, err := p.GenerateAndSaveToken(ctx, req.Account)
	if err != nil {
		return common.InternalServerError(common.AccountSrv, p.Create, err, "GenerateAndSaveToken is failed")
	}

	rsp.Token = token
	rsp.Account = req.Account

	return nil
}

//internal function to generate a token and save it to redis
func (p *AccountService) GenerateAndSaveToken(ctx context.Context, account *account_proto.Account) (string, user_proto.ContactDetailType, error) {
	// separate mode
	var mode user_proto.ContactDetailType
	if len(account.Email) > 0 {
		mode = user_proto.ContactDetailType_EMAIL
	} else if len(account.Phone) > 0 {
		mode = user_proto.ContactDetailType_PHONE
	}
	// put verification to redis db
	token, err := p.PutVerification(ctx, account.Id, mode)
	if err != nil {
		common.InternalServerError(common.AccountSrv, p.GenerateAndSaveToken, err, "PutVerification is failed")
		return "", user_proto.ContactDetailType_ContactDetailType_NONE, err
	}
	return token, mode, nil
}

func (p *AccountService) Update(ctx context.Context, req *account_proto.UpdateRequest, rsp *account_proto.UpdateResponse) error {
	log.Info("Received Account.Update request")

	return nil
}

func (p *AccountService) GetPass(pass string) (string, string) {
	salt := common.Random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+pass), 10)
	if err != nil {
		return "", ""
	}
	pp := base64.StdEncoding.EncodeToString(h)
	return salt, pp
}

//This function takes the verification token, verifies the token, activates the account and sets the password
func (p *AccountService) ConfirmRegister(ctx context.Context, req *account_proto.ConfirmRegisterRequest, rsp *account_proto.ConfirmRegisterResponse) error {
	log.Info("Received Account.ConfirmRegister request")

	req_kv := &kv_proto.ConfirmTokenRequest{req.VerificationToken}
	rsp_kv, err := p.KvClient.ConfirmToken(ctx, req_kv)
	if err != nil {
		return err
	}
	// update status and confirmed flag
	vt := account_proto.VerificationToken{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&vt); err != nil {
		return err
	}

	// read account with id
	account, err := db.Read(ctx, &account_proto.Account{Id: vt.AccountId})
	if err != nil {
		return err
	}
	account.Status = account_proto.AccountStatus_ACTIVE
	account.Confirmed = true
	// password checking
	var salt string
	if len(account.Email) > 0 {
		salt, account.Password = p.GetPass(req.Password)
	} else if len(account.Phone) > 0 {
		salt, account.Passcode = p.GetPass(req.Passcode)
	}
	// create account with password
	if err := db.ConfirmAccountAndUpdatePass(ctx, account, salt); err != nil {
		common.ErrorLog(common.AccountSrv, common.GetFunctionName(p.ConfirmRegister), err, "ConfirmAccountAndUpdatePass is failed")
		return err
	}

	//remove confirmation token from redis
	req_remove := &kv_proto.RemoveTokenRequest{AccountId: vt.AccountId, Token: req.VerificationToken}
	if _, err := p.KvClient.RemoveToken(ctx, req_remove); err != nil {
		common.ErrorLog(common.AccountSrv, common.GetFunctionName(p.ConfirmRegister), err, "RemoveToken is failed")
		return err
	}
	return nil
}

func (p *AccountService) Login(ctx context.Context, req *account_proto.LoginRequest, rsp *account_proto.LoginResponse) error {
	log.Info("Received Account.Login request")
	// check validation request body
	if len(req.Email) == 0 && len(req.Phone) == 0 {
		return common.InternalServerError(common.AccountSrv, p.Login, nil, "server_error")
	}
	email := strings.ToLower(req.Email)

	// reading account from db
	account, err := db.Read(ctx, &account_proto.Account{Email: email, Phone: req.Phone})
	if err != nil {
		return common.NotFound(common.AccountSrv, p.Login, err, "Read query is failed")
	}

	if account.Status != account_proto.AccountStatus_ACTIVE {
		switch account.Status {
		case account_proto.AccountStatus_INACTIVE:
			return common.Forbidden(common.AccountSrv, p.Login, err, "inactive_account")
		case account_proto.AccountStatus_LOCKED:
			return common.Forbidden(common.AccountSrv, p.Login, err, "locked_account")
		case account_proto.AccountStatus_SUSPENDED:
			return common.Forbidden(common.AccountSrv, p.Login, err, "suspended_account")
		}
	}

	if !account.Confirmed {
		return common.Forbidden(common.AccountSrv, p.Login, err, "Not confirmed user login is failed")
	}

	// check lock status
	rsp_islocked, err := p.KvClient.IsLocked(ctx, &kv_proto.IsLockedRequest{common.ACCOUNT_LOCKED_INDEX, account.Id})
	if err != nil {
		log.Error("redis lock err:", err)
		return err
	}
	if rsp_islocked.Locked {
		return common.Forbidden(common.AccountSrv, p.Login, err, "Locked user login is failed")
	}

	// check with email and phone type
	var pass string
	var mode user_proto.ContactDetailType
	if len(req.Email) > 0 {
		pass = req.Password
		mode = user_proto.ContactDetailType_EMAIL // email mode
	} else if len(req.Phone) > 0 {
		pass = req.Passcode
		mode = user_proto.ContactDetailType_PHONE // phone mode
	}
	salt, secret, err := db.SaltAndPassword(ctx, account.Id, mode)
	if err != nil {
		log.Error("password-1 err:", err)
		return common.NotFound(common.AccountSrv, p.Login, err, "not_found")
	}
	s, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		log.Error("password-2 err:", err)
		return common.InternalServerError(common.AccountSrv, p.Login, err, "server_error")
	}
	// does it match?
	if err := bcrypt.CompareHashAndPassword(s, []byte(x+salt+pass)); err != nil {
		if err := p.AuthFailed(ctx, account.Id); err != nil {
			log.Error("password-3 err:", err)
		}

		log.Error("compare password err:", err)
		return common.Unauthorized(common.AccountSrv, p.Login, err, "access_denied")
	}

	//pass matches, reset lock counts if any
	log.Info("Resetting any authentication failure counts for account: ", account.Id)
	p.KvClient.UnLock(ctx, &kv_proto.UnLockRequest{AuthFailIndex: common.AUTHENTIFICATION_INDEX, LockIndex: common.ACCOUNT_LOCKED_INDEX, AccountId: account.Id})

	// fetch user from user-srv with accounid
	rsp_user, err := p.UserClient.ReadByAccount(ctx, &user_proto.ReadByAccountRequest{account.Id})

	if err != nil {
		return common.NotFound(common.AccountSrv, p.Login, err, "ReadByAccount query is failed")
	}
	if rsp_user.Data.User == nil {
		return common.NotFound(common.AccountSrv, p.Login, err, "User is not found")
	}

	// getting employee from employee edge
	userid := rsp_user.Data.User.Id
	orgid := rsp_user.Data.User.OrgId

	// update app_identifer + device_token + platform only if all three present in the request
	if len(req.AppIdentifier) > 0 && req.Platform > 0 && len(req.DeviceToken) > 0 {
		log.Info("Updating tokens for user: ", rsp_user.Data.User.Id)
		token := &user_proto.Token{DeviceToken: req.DeviceToken, Platform: req.Platform, UniqueIdentifier: req.UniqueIdentifier, AppIdentifier: req.AppIdentifier}
		tokens := map[string]*user_proto.Token{}
		if rsp_user.Data.User.Tokens != nil {
			tokens = rsp_user.Data.User.Tokens
		}
		tokens[req.AppIdentifier+":"+fmt.Sprint(req.Platform)] = token
		req_token := &user_proto.UpdateTokenRequest{userid, tokens}
		_, err := p.UserClient.UpdateTokens(ctx, req_token)
		if err != nil {
			common.InternalServerError(common.AccountSrv, p.Login, err, "UpdateTokens is failed")
			return err
		}
	}

	// checking valid organisation
	if rsp_org, err := p.OrganisationClient.Read(ctx, &organisation_proto.ReadRequest{orgid}); err == nil && rsp_org.Data.Organisation != nil {
		// _, err := p.KvClient.SetsOrgInfo(ctx, &kv_proto.SetsOrgInfoRequest{orgid})
		// create org_info
		oi := &organisation_proto.OrgInfo{
			OrgId:   rsp_org.Data.Organisation.Id,
			Type:    rsp_org.Data.Organisation.Type,
			Owner:   rsp_org.Data.Organisation.Owner,
			Modules: rsp_org.Data.Organisation.Modules,
		}
		if _, err := p.OrganisationClient.PutOrgInfo(ctx, &organisation_proto.PutOrgInfoRequest{OrgId: orgid, OrgInfo: oi}); err != nil {
			common.NotFound(common.AccountSrv, p.Login, err, "PutOrgInfo is failed")
			return err
		}
	} else {
		orgid = ""
	}

	//PutEmployeeInfo will return nil, if the user is not an employee
	if _, err := p.TeamClient.PutEmployeeInfo(ctx, &team_proto.PutEmployeeInfoRequest{UserId: userid, OrgId: orgid}); err != nil {
		return common.InternalServerError(common.AccountSrv, p.Login, err, "PutEmployeeInfo is failed")
	}

	var session string
	// check already existed session with account_id
	req_kv_get := &kv_proto.GetExRequest{
		Index: common.SESSION_CONFIRM_INDEX,
		Key:   account.Id,
	}
	if req_kv_get, err := p.KvClient.GetEx(context.TODO(), req_kv_get); err != nil {
		// create session
		if session, err = GenerateRandomString(32); err != nil {
			return err
		}
		// create session_confirm
		req_kv_session := &kv_proto.PutExRequest{
			Index: common.SESSION_CONFIRM_INDEX,
			Item: &kv_proto.Item{
				Key:        account.Id,
				Value:      []byte(session),
				Expiration: int64(Oneday),
			},
		}
		if _, err := p.KvClient.PutEx(ctx, req_kv_session); err != nil {
			return common.InternalServerError(common.AccountSrv, p.Login, err, "PutSessionInfo1 is failed")
		}
	} else {
		session = string(req_kv_get.Item.Value)
	}

	// update session_info with expire date again.
	si := &account_proto.SessionInfo{
		AccountId: account.Id,
		UserId:    userid,
		OrgId:     orgid,
		ExpiresAt: time.Now().Add(Oneday).Unix(),
	}
	body, err := json.Marshal(si)
	if err != nil {
		return err
	}
	req_kv := &kv_proto.PutExRequest{
		Index: common.SESSION_INDEX,
		Item: &kv_proto.Item{
			Key:        session,
			Value:      body,
			Expiration: int64(Oneday),
		},
	}
	if _, err := p.KvClient.PutEx(context.TODO(), req_kv); err != nil {
		return common.InternalServerError(common.AccountSrv, p.Login, err, "PutSessionInfo2 is failed")
	}

	// make response object
	rsp.Data = &account_proto.LoginResponse_Data{
		Session: &account_proto.Session{
			Id:        session,
			ExpiresAt: si.ExpiresAt,
		},
		AccountInfo: &account_proto.AccountInfo{
			UserId: userid,
			OrgId:  orgid,
		},
	}
	return nil
}

func (p *AccountService) Logout(ctx context.Context, req *account_proto.LogoutRequest, rsp *account_proto.LogoutResponse) error {
	log.Info("Received Account.Logout request")
	rsp_kv, err := p.KvClient.ReadSession(ctx, &kv_proto.ReadSessionRequest{common.SESSION_INDEX, req.SessionId})
	if err != nil {
		return err
	}
	si := &account_proto.SessionInfo{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&si); err != nil {
		return common.InternalServerError(common.AccountSrv, p.Logout, err, "Marshaller error")
	}
	log.WithField("account_id", si.AccountId).Warn("Get kv")

	//remove session
	log.Info("Removing session for user: ", si.AccountId)
	if _, err := p.KvClient.RemoveSession(ctx, &kv_proto.RemoveSessionRequest{common.SESSION_INDEX, req.SessionId}); err != nil {
		return common.InternalServerError(common.AccountSrv, p.Logout, err, "RemoveSession is failed")
	}
	if _, err := p.KvClient.RemoveSession(ctx, &kv_proto.RemoveSessionRequest{common.SESSION_CONFIRM_INDEX, si.AccountId}); err != nil {
		return common.InternalServerError(common.AccountSrv, p.Logout, err, "RemoveSession is failed")
	}

	//remove employee info
	log.Info("Removing employee info for user: ", req.UserId)
	req_kv := &kv_proto.DelExRequest{Index: common.EMPLOYEE_INFO_INDEX, Key: req.UserId}
	if _, err := p.KvClient.DelEx(context.TODO(), req_kv); err != nil {
		return common.InternalServerError(common.AccountSrv, p.Logout, err, "DelEx is failed")
	}

	return nil
}

//this function specifically responds to requests for resending confirmation token
func (p *AccountService) ConfirmResend(ctx context.Context, req *account_proto.ConfirmResendRequest, rsp *account_proto.ConfirmResendResponse) error {
	log.Info("Received Account.ConfirmResend request")
	// read account with 2 cases
	account, err := db.Read(ctx, &account_proto.Account{Email: strings.ToLower(req.Email), Phone: req.Phone})
	if err != nil {
		return common.NotFound(common.AccountSrv, p.ConfirmResend, err, "ConfirmResend is failed because of account read error")
	}

	// update status to inactive (for active account), and confirmed flag to false (where account is already confirmed)
	if err := db.UpdateStatus(ctx, account.Id, account_proto.AccountStatus_INACTIVE, false); err != nil {
		return common.InternalServerError(common.AccountSrv, p.ConfirmResend, err, "ConfirmResend is failed because of update status error")
	}

	//resend confirmation token
	_, err = p.ResendNewAccountConfirmationToken(ctx, account)
	return nil
}

//Function to recover/update password. Request is made by the user themselves
func (p *AccountService) RecoverPassword(ctx context.Context, req *account_proto.RecoverPasswordRequest, rsp *account_proto.RecoverPasswordResponse) error {
	log.Info("Received Account.RecoverPassword request")

	account, err := db.Read(ctx, &account_proto.Account{Email: strings.ToLower(req.Email), Phone: req.Phone})
	if err != nil {
		return common.NotFound(common.AccountSrv, p.RecoverPassword, err, "not_found")
	}
	if account.Status != account_proto.AccountStatus_ACTIVE {
		return common.InternalServerError(common.AccountSrv, p.RecoverPassword, err, "inactive_account")
	}
	if !account.Confirmed {
		return common.InternalServerError(common.AccountSrv, p.RecoverPassword, err, "not_confirmed")
	}

	// generate and save account token
	token, mode, err := p.GenerateAndSaveToken(ctx, account)
	if err != nil {
		return common.InternalServerError(common.AccountSrv, p.RecoverPassword, err, "GenerateAndSaveToken is failed")
	}

	//FIXME: This can be refactored into a single function and combined for NewAccountConfirmationToken, ResendNewAccountConfirmationToken
	rsp_user, rsp_org, err := p.ReadUserAndOrg(ctx, account.Id)
	if err != nil {
		return err
	}
	switch mode {
	case user_proto.ContactDetailType_EMAIL:
		//compose email message here
		// message := fmt.Sprintf(common.MSG_ACCOUNT_VERIFICATION_SMS, rsp_org.Data.Organisation.Name, rsp_user.Data.User.Firstname, req.Token, rsp_org.Data.Organisation.Name)
		// p.publishMessage(req.Mode, rsp_account.Data.Account.Phone, message)
		message := fmt.Sprintf(common.MSG_ACCOUNT_VERIFICATION_SMS, rsp_org.Data.Organisation.Name, rsp_user.Data.User.Firstname, token, rsp_org.Data.Organisation.Name)
		go p.publishMessage(int32(mode), account.Phone, message)

	case user_proto.ContactDetailType_PHONE:
		message := fmt.Sprintf(common.MSG_ACCOUNT_USER_PASSWORD_RESET_REQUEST_BY_USER_SMS, rsp_user.Data.User.Firstname, token, rsp_org.Data.Organisation.Name)
		go p.publishMessage(int32(mode), account.Phone, message)
	}
	rsp.Data = &account_proto.RecoverPasswordResponse_Data{PasswordResetToken: token}
	return nil
}

//Function to recover/update password. Password is reset by the employee
func (p *AccountService) ResetUserPassword(ctx context.Context, req *account_proto.ResetUserPasswordRequest, rsp *account_proto.ResetUserPasswordResponse) error {
	log.Info("Received Account.ResetUserPassword request")

	// read account with 2 cases
	account, err := db.GetAccountByUser(ctx, req.UserId)
	if err != nil {
		return common.NotFound(common.AccountSrv, p.ResetUserPassword, err, "ConfirmResend is failed because of account read error")
	}

	// checking for valid employee
	req_employee := &team_proto.ReadEmployeeInfoRequest{req.TeamId}
	rsp_employee, err := p.TeamClient.CheckValidEmployee(ctx, req_employee)
	if err != nil {
		return common.InternalServerError(common.AccountSrv, p.ResetUserPassword, err, "CheckValidEmployee is failed")
	}

	var mode user_proto.ContactDetailType
	//TODO:Generate random pass here if 'random' is indicated and move the random pass generation to server side
	//user account will now have a temporary password/passcode (old password/passcode is overwritten)
	err = p.updatePass(ctx, req.Password, req.Passcode, account.Id)
	if err != nil {
		return err
	}

	// generate and save account token
	token, mode, err := p.GenerateAndSaveToken(ctx, account)
	if err != nil {
		return common.InternalServerError(common.AccountSrv, p.ResetUserPassword, err, "GenerateAndSaveToken is failed")
	}

	//FIXME: This can be refactored into a single function and combined for NewAccountConfirmationToken, ResendNewAccountConfirmationToken
	rsp_user, rsp_org, err := p.ReadUserAndOrg(ctx, account.Id)
	if err != nil {
		return err
	}
	switch mode {
	case user_proto.ContactDetailType_EMAIL:
		//compose email message here
		// message := fmt.Sprintf(common.MSG_ACCOUNT_VERIFICATION_SMS, rsp_org.Data.Organisation.Name, rsp_user.Data.User.Firstname, req.Token, rsp_org.Data.Organisation.Name)
		// p.publishMessage(req.Mode, rsp_account.Data.Account.Phone, message)
	case user_proto.ContactDetailType_PHONE:
		message := fmt.Sprintf(common.MSG_ACCOUNT_USER_PASSWORD_RESET_BY_EMPLOYEE_SMS, rsp_user.Data.User.Firstname, rsp_org.Data.Organisation.Name, rsp_employee.Employee.User.Firstname, token, rsp_org.Data.Organisation.Name)
		return p.publishMessage(int32(mode), account.Phone, message)
	}

	return nil
}

//API request sent by the user to update the password or pascode with the password reset token
func (p *AccountService) UpdatePassword(ctx context.Context, req *account_proto.UpdatePasswordRequest, rsp *account_proto.UpdatePasswordResponse) error {
	log.Info("Received Account.UpdatePassword request")

	req_kv := &kv_proto.ConfirmTokenRequest{req.PasswordResetToken}
	rsp_kv, err := p.KvClient.ConfirmToken(ctx, req_kv)
	if err != nil {
		return common.InternalServerError(common.AccountSrv, p.UpdatePassword, err, "token_expire")
	}
	// update status and confirmed flag
	vt := account_proto.VerificationToken{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&vt); err != nil {
		return common.InternalServerError(common.AccountSrv, p.UpdatePassword, err, "parsing error")
	}

	// update password
	err = p.updatePass(ctx, req.Password, req.Passcode, vt.AccountId)
	if err != nil {
		return common.InternalServerError(common.AccountSrv, p.UpdatePassword, err, "update_pass")
	}

	//remove confirmation token from redis
	req_remove := &kv_proto.RemoveTokenRequest{AccountId: vt.AccountId, Token: req.PasswordResetToken}
	if _, err := p.KvClient.RemoveToken(ctx, req_remove); err != nil {
		common.ErrorLog(common.AccountSrv, common.GetFunctionName(p.UpdatePassword), err, "RemoveToken is failed")
		return err
	}

	return nil
}

func (p *AccountService) NewAccountConfirmationToken(ctx context.Context, req *account_proto.NewAccountConfirmationTokenRequest, rsp *account_proto.NewAccountConfirmationTokenResponse) error {
	log.Info("Recieved Account.NewAccountConfirmationToken request")

	//FIXME: This can be refactored into a single function and combined for NewAccountConfirmationToken, ResendNewAccountConfirmationToken
	rsp_user, rsp_org, err := p.ReadUserAndOrg(ctx, req.Account.Id)
	if err != nil {
		return err
	}
	switch req.Mode {
	case int32(user_proto.ContactDetailType_EMAIL):
		//TODO:compose email message here
		// message := fmt.Sprintf(common.MSG_ACCOUNT_VERIFICATION_SMS, rsp_org.Data.Organisation.Name, rsp_user.Data.User.Firstname, req.Token, rsp_org.Data.Organisation.Name)
		// p.publishMessage(req.Mode, rsp_account.Data.Account.Phone, message)
	case int32(user_proto.ContactDetailType_PHONE):
		message := fmt.Sprintf(common.MSG_ACCOUNT_VERIFICATION_SMS, rsp_org.Data.Organisation.Name, rsp_user.Data.User.Firstname, req.Token, rsp_org.Data.Organisation.Name)
		return p.publishMessage(req.Mode, req.Account.Phone, message)
	}
	return nil
}

//internal function to resend account confirmation token
func (p *AccountService) ResendNewAccountConfirmationToken(ctx context.Context, account *account_proto.Account) (string, error) {
	log.Info("Received Account.ResendConfirmationToken request")

	// generate and save account token
	token, mode, err := p.GenerateAndSaveToken(ctx, account)
	if err != nil {
		common.InternalServerError(common.AccountSrv, p.ResendNewAccountConfirmationToken, err, "GenerateAndSaveToken is failed")
		return "", err
	}

	//FIXME: This can be refactored into a single function and combined for NewAccountConfirmationToken, ResendNewAccountConfirmationToken
	rsp_user, rsp_org, err := p.ReadUserAndOrg(ctx, account.Id)
	if err != nil {
		return "", err
	}

	switch mode {
	case user_proto.ContactDetailType_EMAIL:
		// TODO:compose email message here
		// message := fmt.Sprintf(common.MSG_ACCOUNT_VERIFICATION_SMS, rsp_org.Data.Organisation.Name, rsp_user.Data.User.Firstname, req.Token, rsp_org.Data.Organisation.Name)
		// p.publishMessage(req.Mode, rsp_account.Data.Account.Phone, message)
	case user_proto.ContactDetailType_PHONE:
		message := fmt.Sprintf(common.MSG_ACCOUNT_VERIFICATION_RESEND_SMS, rsp_user.Data.User.Firstname, token, rsp_org.Data.Organisation.Name)
		return token, p.publishMessage(int32(mode), account.Phone, message)
	}
	return "", nil
}

//internal function to read and return user and org required in other functions
func (p *AccountService) ReadUserAndOrg(ctx context.Context, accountId string) (*user_proto.ReadByAccountResponse, *organisation_proto.ReadResponse, error) {
	// the following can be combined into a single query? instead of making 3 calls to the database
	// getting user
	rsp_user, err := p.UserClient.ReadByAccount(ctx, &user_proto.ReadByAccountRequest{accountId})
	if err != nil {
		common.NotFound(common.AccountSrv, p.ReadUserAndOrg, err, "ConfirmResend is failed because of user read error:"+accountId)
		return nil, nil, err
	}
	// getting organisation
	rsp_org, err := p.OrganisationClient.Read(ctx, &organisation_proto.ReadRequest{rsp_user.Data.User.OrgId})
	if err != nil {
		common.NotFound(common.AccountSrv, p.ReadUserAndOrg, err, "ConfirmResend is failed because of organisation read error")
		return nil, nil, err
	}
	///////////////////
	return rsp_user, rsp_org, nil
}

//internal function to update password or passcode consistently
func (p *AccountService) updatePass(ctx context.Context, password, passcode, accountId string) error {
	log.Info("Updating pass for account: ", accountId)
	//update password for Account.Email
	if len(password) > 0 {
		log.Info("Updating password for account: ", accountId)
		salt, pass, err := generateEncodedPass(ctx, password)
		if err == nil {
			if err := db.UpdatePassword(ctx, accountId, pass, salt); err != nil {
				common.ErrorLog(common.AccountSrv, common.GetFunctionName(p.updatePass), err, "UpdatePassword query is failed")
				return err
			}
		}
	} else if len(passcode) > 0 {
		//update passcode for Account.Phone
		log.Info("Updating passcode for account: ", accountId)
		salt, pass, err := generateEncodedPass(ctx, passcode)
		if err == nil {
			if err := db.UpdatePasscode(ctx, accountId, pass, salt); err != nil {
				common.ErrorLog(common.AccountSrv, common.GetFunctionName(p.updatePass), err, "UpdatePasscode query is failed")
				return err
			}
		}
	}

	//clear any existing active sessions for this account
	log.Info("Clearning active sessions for account: ", accountId)
	if err := p.clearActiveSession(ctx, accountId); err != nil {
		common.ErrorLog(common.AccountSrv, common.GetFunctionName(p.updatePass), err, "RemoveSession is failed")
		return err
	}
	return nil
}

func generateEncodedPass(ctx context.Context, pass string) (string, string, error) {
	salt := common.Random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+pass), 10)
	if err != nil {
		return "", "", common.InternalServerError(common.AccountSrv, generateEncodedPass, err, "parsing error")
	}
	return salt, base64.StdEncoding.EncodeToString(h), nil
}

func (p *AccountService) Lock(ctx context.Context, req *account_proto.LockRequest, rsp *account_proto.LockResponse) error {
	log.Info("Received Account.Lock request")

	// updated status to lock in arangodb
	if err := db.UpdateLockStatus(ctx, req.AccountId); err != nil {
		return common.InternalServerError(common.AccountSrv, p.Lock, err, "update lock status err")
	}

	_, err := p.KvClient.Locked(ctx, &kv_proto.LockedRequest{common.ACCOUNT_LOCKED_INDEX, req.AccountId})
	return err
}

func (p *AccountService) ReadSession(ctx context.Context, req *account_proto.ReadSessionRequest, rsp *account_proto.ReadSessionResponse) error {
	log.Info("Received Account.ReadSession request")
	rsp_kv, err := p.KvClient.ReadSession(ctx, &kv_proto.ReadSessionRequest{common.SESSION_INDEX, req.SessionId})
	if err != nil {
		return common.NotFound(common.AccountSrv, p.ReadSession, err, "not found")
	}

	//updating the session time using flowing window
	req_kv_session := &kv_proto.PutExRequest{
		Index: common.SESSION_INDEX,
		Item: &kv_proto.Item{
			Key:        req.SessionId,
			Value:      []byte(rsp_kv.Value),
			Expiration: int64(Oneday),
		},
	}
	if _, err := p.KvClient.PutEx(ctx, req_kv_session); err != nil {
		return common.InternalServerError(common.AccountSrv, p.ReadSession, err, "session_confirm putex error")
	}

	si := &account_proto.SessionInfo{}
	decoder := json.NewDecoder(bytes.NewReader([]byte(rsp_kv.Value)))
	if err := decoder.Decode(&si); err != nil {
		return common.InternalServerError(common.AccountSrv, p.ReadSession, err, "parsing error")
	}
	rsp.UserId = si.UserId
	rsp.OrgId = si.OrgId
	return nil
}

func (p *AccountService) InternalConfirm(ctx context.Context, req *account_proto.InternalConfirmRequest, rsp *account_proto.InternalConfirmResponse) error {
	log.Info("Received Account.InternalConfirm request")

	account, err := db.Read(ctx, &account_proto.Account{Id: req.AccountId})
	if err != nil {
		return common.NotFound(common.AccountSrv, p.InternalConfirm, err, "not found")
	}
	account.Status = account_proto.AccountStatus_ACTIVE
	account.Confirmed = true

	salt := common.Random(16)
	h, err := bcrypt.GenerateFromPassword([]byte(x+salt+req.Password), 10)
	if err != nil {
		return common.InternalServerError(common.AccountSrv, p.InternalConfirm, err, "parsing error")
	}
	pp := base64.StdEncoding.EncodeToString(h)
	account.Password = pp

	// create account with password
	return db.ConfirmAccountAndUpdatePass(ctx, account, salt)
}

func (p *AccountService) SetAccountStatus(ctx context.Context, req *account_proto.SetAccountStatusRequest, rsp *account_proto.SetAccountStatusResponse) error {
	log.Info("Received Account.SetAccountStatus request")

	//if account is being set to active again
	if req.Status == account_proto.AccountStatus_ACTIVE {
		rsp_account, err := db.GetAccountByUser(ctx, req.UserId)
		if err != nil {
			return common.InternalServerError(common.AccountSrv, p.SetAccountStatus, err, "account status get error")
		}
		//if current status is locked then remove locks in KV
		if rsp_account.Status == account_proto.AccountStatus_LOCKED {
			_, err := p.KvClient.UnLock(ctx, &kv_proto.UnLockRequest{AuthFailIndex: common.AUTHENTIFICATION_INDEX, LockIndex: common.ACCOUNT_LOCKED_INDEX, AccountId: rsp_account.Id})
			if err != nil {
				return common.InternalServerError(common.AccountSrv, p.SetAccountStatus, err, "account unlokcing error")
			}
		}
	}

	return db.SetAccountStatus(ctx, req.UserId, req.Status)
}

func (p *AccountService) GetAccountStatus(ctx context.Context, req *account_proto.GetAccountStatusRequest, rsp *account_proto.GetAccountStatusResponse) error {
	log.Info("Received Account.SetAccountStatus request")

	account_status, err := db.GetAccountStatus(ctx, req.UserId)
	if err != nil {
		return common.NotFound(common.AccountSrv, p.GetAccountStatus, err, "not found")
	}
	rsp.Data = &account_proto.GetAccountStatusResponse_Data{account_status}
	return nil
}

func (p *AccountService) ConfirmVerify(ctx context.Context, req *account_proto.ConfirmVerifyRequest, rsp *account_proto.ConfirmVerifyResponse) error {
	_, err := p.KvClient.GetEx(ctx, &kv_proto.GetExRequest{
		Index: common.VERIFICATION_TOKEN_INDEX,
		Key:   req.Token,
	})
	return err
}

func (p *AccountService) PassVerify(ctx context.Context, req *account_proto.ConfirmVerifyRequest, rsp *account_proto.ConfirmVerifyResponse) error {
	_, err := p.KvClient.GetEx(ctx, &kv_proto.GetExRequest{
		Index: common.VERIFICATION_TOKEN_INDEX,
		Key:   req.Token,
	})
	return err
}

func (p *AccountService) publishMessage(mode int32, target, message string) error {
	log.Info("Publishing message via communication channel")
	// eamil/phone pubsub
	var topic string
	var body []byte
	var err error
	switch mode {
	case int32(user_proto.ContactDetailType_EMAIL):
		topic = common.SEND_EMAIL
		//compose email message here
		body = nil
		return nil
	case int32(user_proto.ContactDetailType_PHONE):
		topic = common.SEND_SMS
		sms := &sms_proto.SendRequest{
			Phone:   target,
			Message: message,
		}
		body, err = json.Marshal(sms)
		if err != nil {
			return err
		}
		return p.Broker.Publish(topic, &broker.Message{Body: body})
	}
	return nil
}

// ReadAccountToken is an internal function that only to be for user test!
func (p *AccountService) ReadAccountToken(ctx context.Context, req *account_proto.ReadAccountTokenRequest, rsp *account_proto.ReadAccountTokenResponse) error {
	req_get := &kv_proto.GetExRequest{Index: common.VERIFICATION_TOKEN_INDEX, Key: req.AccountId}
	rsp_get, err := p.KvClient.GetEx(ctx, req_get)

	if err != nil {
		return common.NotFound(common.AccountSrv, p.ReadAccountToken, err, "redis server error")
	}
	// parsing redis data
	va := account_proto.VerificationAccount{}
	if err := json.Unmarshal(rsp_get.Item.Value, &va); err != nil {
		return common.NotFound(common.AccountSrv, p.ReadAccountToken, err, "parsing error")
	}
	rsp.Token = va.Token
	return nil
}

func (p *AccountService) clearActiveSession(ctx context.Context, account_id string) error {
	log.Info("Received Account.clearActiveSession request")

	// read session with account_id to check whether there are any existing session
	req_get := &kv_proto.GetExRequest{Index: common.SESSION_CONFIRM_INDEX, Key: account_id}
	rsp_get, err := p.KvClient.GetEx(ctx, req_get)

	if rsp_get != nil && err == nil {
		log.Info("Existing session found for account: ", account_id)
		session := string(rsp_get.Item.Value)

		//remove session
		log.Info("Removing session for account: ", account_id)
		if _, err := p.KvClient.RemoveSession(ctx, &kv_proto.RemoveSessionRequest{common.SESSION_INDEX, session}); err != nil {
			common.ErrorLog(common.AccountSrv, common.GetFunctionName(p.clearActiveSession), err, "RemoveSession is failed")
			return err
		}
		if _, err := p.KvClient.RemoveSession(ctx, &kv_proto.RemoveSessionRequest{common.SESSION_CONFIRM_INDEX, account_id}); err != nil {
			common.ErrorLog(common.AccountSrv, common.GetFunctionName(p.clearActiveSession), err, "RemoveSession is failed")
			return err
		}
	} else if err != nil {
		common.ErrorLog(common.AccountSrv, common.GetFunctionName(p.clearActiveSession), err, "No active session found")
		return err
	}
	return nil
}
