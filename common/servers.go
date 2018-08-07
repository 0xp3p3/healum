package common

import (
	"reflect"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// all the server name.
var (
	AccountSrv      = "go.micro.srv.account"
	ActivitySrv     = "go.micro.srv.activity"
	BehaviourSrv    = "go.micro.srv.behaviour"
	ContentSrv      = "go.micro.srv.content"
	KvSrv           = "go.micro.srv.kv"
	MobpushSrv      = "go.micro.srv.mobpush"
	NoteSrv         = "go.micro.srv.note"
	OrganisationSrv = "go.micro.srv.organisation"
	PlanSrv         = "go.micro.srv.plan"
	ProductSrv      = "go.micro.srv.product"
	ResponseSrv     = "go.micro.srv.response"
	SmsSrv          = "go.micro.srv.sms"
	StaticSrv       = "go.micro.srv.static"
	SurveySrv       = "go.micro.srv.survey"
	TaskSrv         = "go.micro.srv.task"
	TeamSrv         = "go.micro.srv.team"
	TodoSrv         = "go.micro.srv.todo"
	TrackSrv        = "go.micro.srv.track"
	UserappSrv      = "go.micro.srv.userapp"
	UserSrv         = "go.micro.srv.user"
	DbSrv           = "go.micro.srv.db"
	EmailSrv        = "go.micro.srv.email"
	AuditSrv        = "go.micro.srv.audit"
)

// ErrorLog logs error
func ErrorLog(srv string, fun interface{}, err error, description string) {
	log.WithFields(log.Fields{
		"srv":  srv,
		"func": GetFunctionName(fun),
		"err":  err,
	}).Error(description)
}

// DebugLog logs debug
func DebugLog(srv string, fun interface{}, value interface{}, description string) {
	log.WithFields(log.Fields{
		"srv":  srv,
		"func": GetFunctionName(fun),
		"val":  value,
	}).Debug(description)
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
