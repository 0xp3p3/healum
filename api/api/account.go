package api

import (
	"context"
	"net/http"
	account_proto "server/account-srv/proto/account"
	"server/api/utils"
	audit_proto "server/audit-srv/proto/audit"
	"server/common"

	"github.com/emicklei/go-restful"
	"github.com/micro/go-os/metrics"
	log "github.com/sirupsen/logrus"
)

// Event external API handler
type AccountService struct {
	AccountClient account_proto.AccountServiceClient
	Auth          Filters
	Audit         AuditFilter
	ServerMetrics metrics.Metrics
}

func (p AccountService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/account")

	audit := &audit_proto.Audit{
		ActionService:  common.AccountSrv,
		ActionResource: common.BASE + common.ACCOUNT_TYPE,
	}

	ws.Route(ws.POST("/confirm").To(p.ConfirmRegister).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Confirm user account"))

	ws.Route(ws.POST("/confirm/resend").To(p.ConfirmResend).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Resend account verification token"))

	ws.Route(ws.POST("/login").To(p.Login).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("User login"))

	ws.Route(ws.GET("/logout").To(p.Logout).
		Filter(p.Auth.BasicAuthenticate).
		Filter(p.Auth.OrganisationAuthenticate).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Logout the user"))

	ws.Route(ws.POST("/pass/recover").To(p.RecoverPassword).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Request password recovery"))

	ws.Route(ws.POST("/pass/update").To(p.UpdatePassword).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Update a user password"))

	ws.Route(ws.POST("/confirm/verify").To(p.ConfirmVerify).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Confirm verify token"))

	ws.Route(ws.POST("/pass/verify").To(p.PassVerify).
		// Filter(p.Audit.Clone(audit).AuditFilter).
		Doc("Pass verify token"))

	restful.Add(ws)
}

/**
* @api {post} /server/account/confirm Confirming registrations
* @apiVersion 0.1.0
* @apiName ConfirmRegister
* @apiGroup Account
*
* @apiDescription The functionality is to confirm their account and set the status of the account .
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/account/confirm
*
* @apiParamExample {json} Request-Email:
* {
*   "verification_token": "f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE="
*	"password":"some password"
* }
* @apiParamExample {json} Request-Phone:
* {
*   "verification_token": "654321"
*	"passcode":"123456"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Confirmed successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, QueryError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "BindError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.account.ConfirmRegister",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "BindError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.account.ConfirmRegister",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *AccountService) ConfirmRegister(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Account.ConfirmRegister API request")
	req_confirm := new(account_proto.ConfirmRegisterRequest)
	err := req.ReadEntity(req_confirm)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.ConfirmRegister", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.AccountClient.ConfirmRegister(ctx, req_confirm)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.ConfirmRegister", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Confirmed successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/account/confirm/resend Resend verification token
* @apiVersion 0.1.0
* @apiName ConfirmResend
* @apiGroup Account
*
* @apiDescription The functionality is to resend a verification token incase the verification token has been expired. The account needs to exist in the database for this to be a valid request.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/account/confirm/resend
*
* @apiParamExample {json} Request-Email:
* {
*   "email": "eamil8@email.com"
* }
*
* @apiParamExample {json} Request-Phone:
* {
*   "phone": "123-4567-890"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Resend token successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, QueryError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "BindError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.account.ConfirmResend",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *AccountService) ConfirmResend(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Account.ConfirmResend API request")
	req_resend := new(account_proto.ConfirmResendRequest)
	err := req.ReadEntity(req_resend)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.ConfirmResend", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.AccountClient.ConfirmResend(ctx, req_resend)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.ConfirmResend", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Resend token successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/account/login User login
* @apiVersion 0.1.0
* @apiName Login
* @apiGroup Account
*
* @apiDescription The functionality is to be able to login using password or code. The account needs to exist in the database for this to be a valid request.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/account/login
*
* @apiParamExample {json} Request-Email:
* {
*   "email": "email8@email.com",
*   "password": "pass1"
* }
*
* @apiParamExample {json} Request-Phone:
* {
*   "phone": "123-4567-890",
*   "passcode": "12345",
*   "device_token": "token_string",
*   "unique_identifier": "unique_id",
*   "app_identifier": "app.bundle.id",
*   "platform": 1
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "data": {
*     "session": {
*       "id": "f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE=",
*       "expires_at": 153252466
*     },
*     "user": { User }
*   },
*   "code": 200,
*   "message": "Login successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, QueryError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "BindError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.account.Login",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *AccountService) Login(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Account.Login API request")
	req_login := new(account_proto.LoginRequest)
	err := req.ReadEntity(req_login)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.Login", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.AccountClient.Login(ctx, req_login)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.Login", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Login successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {get} /server/account/logout Logout a user
* @apiVersion 0.1.0
* @apiName Logout
* @apiGroup Account
*
* @apiDescription The functionality is to be able to login using password or code. The account needs to exist in the database for this to be a valid request.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/account/logout?session={session_id}
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Logout successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, QueryError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "BindError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.account.Logout",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *AccountService) Logout(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Account.Logout API request")
	req_logout := new(account_proto.LogoutRequest)
	req_logout.SessionId = req.QueryParameter("session")
	req_logout.UserId = req.Attribute(UserIdAttrName).(string)
	req_logout.OrgId = req.Attribute(OrgIdAttrName).(string)
	log.WithField("session_id", req_logout.SessionId).Warn("Received session")

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.AccountClient.Logout(ctx, req_logout)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.Logout", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Logout successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/account/pass/recover Request password recovery token
* @apiVersion 0.1.0
* @apiName RecoverPassword
* @apiGroup Account
*
* @apiDescription The functionality is to request password or passcode recovery token. The account needs to exist in the database for this to be a valid request. This allows the user to reset forgotten password or passcode
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/account/pass/recover
*
* @apiParamExample {json} Request-Email:
* {
*   "email": "eamil8@email.com"
* }
*
* @apiParamExample {json} Request-Phone:
* {
*   "phone": "123-4567-890"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Sent token successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, QueryError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "BindError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.account.RecoverPassword",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *AccountService) RecoverPassword(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Account.RecoverPassword API request")
	req_recover := new(account_proto.RecoverPasswordRequest)
	err := req.ReadEntity(req_recover)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.RecoverPassword", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.AccountClient.RecoverPassword(ctx, req_recover)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.RecoverPassword", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Sent token successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/account/pass/update Update a user account's password or code
* @apiVersion 0.1.0
* @apiName UpdatePassword
* @apiGroup Account
*
* @apiDescription The functionality is to update account password or code. The account needs to exist in the database for this to be a valid request.
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/account/pass/update
*
* @apiParamExample {json} Request-Email:
* {
*   "password_reset_token": "f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE="
*   "password": "pass1"
* }
*
* @apiParamExample {json} Request-Phone:
* {
*   "password_reset_token": "f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE="
*   "passcode": "123456"
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Updated successfully"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, QueryError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "BindError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.account.UpdatePassword",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *AccountService) UpdatePassword(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Account.UpdatePassword API request")
	req_update := new(account_proto.UpdatePasswordRequest)
	err := req.ReadEntity(req_update)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.UpdatePassword", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.AccountClient.UpdatePassword(ctx, req_update)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.UpdatePassword", "QueryError")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Updated successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/account/confirm/verify Confirm verify token
* @apiVersion 0.1.0
* @apiName ConfirmVerify
* @apiGroup Account
*
* @apiDescription This doesn't confirm the account but just checks on redis if the confirmation token is valid (return 200 OK) or not (return error)
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/account/confirm/verify
*
* @apiParamExample {json} Request-Email:
* {
*   "token": "f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE="
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Token is valid"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, QueryError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "BindError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.account.ConfirmVerify",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *AccountService) ConfirmVerify(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Account.ConfirmVerify API request")
	req_resend := new(account_proto.ConfirmVerifyRequest)
	err := req.ReadEntity(req_resend)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.ConfirmVerify", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.AccountClient.ConfirmVerify(ctx, req_resend)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.ConfirmVerify", "Token is invalid")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Token is confirmed successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}

/**
* @api {post} /server/account/pass/verify Pass verify token
* @apiVersion 0.1.0
* @apiName PassVerify
* @apiGroup Account
*
* @apiDescription This doesn't confirm the password reset token but just checks on redis if the confirmation token is valid (return 200 OK) or not (return error)
*
* @apiExample Example usage:
* curl -i http://BASE_SERVER_URL/server/account/pass/verify
*
* @apiParamExample {json} Request-Email:
* {
*   "token": "f01ckVcMHLjgmsGXyKJbLdlovJyw-71C4HshATxe6tE="
* }
*
* @apiSuccessExample Success-Response:
* HTTP/1.1 200 OK
* {
*   "code": 200,
*   "message": "Token is valid"
* }
*
* @apiError NoAuthorized 	Only authenticated users can access the data.
* @apiError BadRequest   	BindError, QueryError
*
* @apiErrorExample Error-Response:
*     HTTP/1.1 400 Bad Request
*     {
*       "code": 400,
*       "message": "BindError",
*       "errors": [
*         {
*           "domain": "go.micro.srv.account.PassVerify",
*           "reason": "{\"id\":\"go.micro.client\",\"code\":500,\"detail\":\"none available\",\"status\":\"Internal Server Error\"}"
*         }
*       ]
*     }
 */

func (p *AccountService) PassVerify(req *restful.Request, rsp *restful.Response) {
	log.Info("Received Account.PassVerify API request")
	req_resend := new(account_proto.ConfirmVerifyRequest)
	err := req.ReadEntity(req_resend)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.PassVerify", "BindError")
		return
	}

	ctx := common.NewContextByHeader(context.TODO(), req.Request.Header)
	resp, err := p.AccountClient.PassVerify(ctx, req_resend)
	if err != nil {
		utils.WriteErrorResponseWithCode(rsp, err, "go.micro.srv.account.PassVerify", "Token is invalid")
		return
	}

	resp.Code = http.StatusOK
	resp.Message = "Token is confirmed successfully"
	rsp.AddHeader("Content-Type", "application/json")
	rsp.WriteHeaderAndEntity(http.StatusOK, resp)
}
