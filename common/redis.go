package common

const (
	VERIFICATION_TOKEN_INDEX = 0  // will store token after register, password reset
	AUTHENTIFICATION_INDEX   = 1  // will store count of auth failed
	SESSION_INDEX            = 2  // will store session after login
	ACCOUNT_LOCKED_INDEX     = 3  // will store locked account
	EMPLOYEE_INFO_INDEX      = 4  // will store employee info after login of an employee
	ORG_INFO_INDEX           = 5  // will store organisation info after create, warmcach
	SESSION_CONFIRM_INDEX    = 6  // will store session
	CLOUD_TAGS_INDEX         = 7  // will store tags for goal, challenge, hate, plan, content and survey tag info
	TRACK_INDEX              = 10 // will store track info for track-srv
	USERAPP_INDEX            = 11 // will store track info fror userapp-srv

	PLAN      = "plan"
	GOAL      = "goal"
	CHALLENGE = "challenge"
	HABIT     = "habit"
	SURVEY    = "survey"
	CONTENT   = "content"
)
