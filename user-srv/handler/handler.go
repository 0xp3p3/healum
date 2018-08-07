package handler

import (
	"context"
	"encoding/json"
	account_proto "server/account-srv/proto/account"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	content_proto "server/content-srv/proto/content"
	kv_proto "server/kv-srv/proto/kv"
	static_proto "server/static-srv/proto/static"
	survey_proto "server/survey-srv/proto/survey"
	team_proto "server/team-srv/proto/team"
	track_proto "server/track-srv/proto/track"
	"server/user-srv/db"
	user_proto "server/user-srv/proto/user"

	"github.com/micro/go-micro/broker"
	_ "github.com/micro/go-plugins/broker/nats"
	_ "github.com/micro/go-plugins/transport/nats"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type UserService struct {
	Broker          broker.Broker
	KvClient        kv_proto.KvServiceClient
	AccountClient   account_proto.AccountServiceClient
	TrackClient     track_proto.TrackServiceClient
	TeamClient      team_proto.TeamServiceClient
	StaticClient    static_proto.StaticServiceClient
	BehaviourClient behaviour_proto.BehaviourServiceClient
	SurveyClient    survey_proto.SurveyServiceClient
	ContentClient   content_proto.ContentServiceClient
}

func (p *UserService) All(ctx context.Context, req *user_proto.AllRequest, rsp *user_proto.AllResponse) error {
	log.Info("Received User.All request")

	users, err := db.All(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(users) == 0 || err != nil {
		return common.NotFound(common.UserSrv, p.All, err, "not found")
	}
	rsp.Data = &user_proto.UserArrData{users}
	return nil
}

// Creates a user account
func (p *UserService) Create(ctx context.Context, req *user_proto.CreateRequest, rsp *user_proto.CreateResponse) error {
	log.Info("Received User.Create request")

	// create account
	var token string
	var mode user_proto.ContactDetailType
	//TODO:check wther the contactdetail type is primary contact or not
	if req.Account != nil {
		// create a new account
		if len(req.Account.Email) > 0 {
			mode = user_proto.ContactDetailType_EMAIL
		} else if len(req.Account.Phone) > 0 {
			mode = user_proto.ContactDetailType_PHONE
		}
		req_account := &account_proto.CreateRequest{req.Account}
		rsp_account, err := p.AccountClient.Create(ctx, req_account)
		if err != nil {
			return common.InternalServerError(common.UserSrv, p.Create, err, "query is failed")
		}
		token = rsp_account.Token
		req.Account = rsp_account.Account
	}

	// fetch user from user-srv with accounid
	user, err := db.ReadByAccount(ctx, req.Account.Id)
	if user == nil && err != nil {
		// create a new user after account creating, if there is no user for this account
		log.Info("Creating a new user for account: ", req.Account.Id)
		if len(req.User.Id) == 0 {
			req.User.Id = uuid.NewUUID().String()
		}

		err := db.Create(ctx, req.User, req.Account)
		user = req.User
		if err != nil {
			return common.InternalServerError(common.UserSrv, p.Create, err, "create error")
		}

		// checking for valid employee
		req_employee := &team_proto.ReadEmployeeInfoRequest{req.TeamId}
		rsp_employee, err := p.TeamClient.CheckValidEmployee(ctx, req_employee)
		if err != nil {
			return common.InternalServerError(common.UserSrv, p.Create, err, "CheckValidEmployee is failed")
		}

		if rsp_employee.Valid && rsp_employee.Employee != nil {
			//add measurements
			if req.User.Preference != nil && req.User.Preference.CurrentMeasurements != nil && len(req.User.Preference.CurrentMeasurements) > 0 {
				log.Info("Adding current measurements for user: ", req.User.Id)
				err = p.AddMeasurements(ctx, req.User.Id, req.OrgId, rsp_employee.Employee, req.User.Preference.CurrentMeasurements)
				if err != nil {
					return common.InternalServerError(common.UserSrv, p.Create, err, "AddMultipleMeasurements is failed")
				}
			}
		}
	}

	// send new confirm token message via email or phone
	if req.Account != nil {
		log.Info("Sending confirmation token for account: ", req.Account.Id)
		if _, err := p.AccountClient.NewAccountConfirmationToken(ctx, &account_proto.NewAccountConfirmationTokenRequest{
			Token:   token,
			Mode:    int32(mode),
			Account: &account_proto.Account{Id: req.Account.Id}}); err != nil {
			return common.InternalServerError(common.UserSrv, p.Create, err, "token error")
		}
	}

	rsp.Data = &user_proto.UserTokenData{
		User:    req.User,
		Account: &account_proto.Account{Id: req.Account.Id},
	}
	return nil
}

// udpate a users account
func (p *UserService) Update(ctx context.Context, req *user_proto.UpdateRequest, rsp *user_proto.UpdateResponse) error {
	log.Info("Received User.Update request")

	// update user
	err := db.Update(ctx, req.User)
	if err != nil {
		return common.InternalServerError(common.UserSrv, p.Update, err, "update error")
	}

	rsp.Data = &user_proto.UserData{
		User: req.User,
	}
	return nil
}

func (p *UserService) Read(ctx context.Context, req *user_proto.ReadRequest, rsp *user_proto.ReadResponse) error {
	log.Info("Received User.Read request")
	user, err := db.Read(ctx, req.UserId, req.OrgId)
	if user == nil || err != nil {
		return common.NotFound(common.UserSrv, p.Read, err, "not found")
	}
	rsp.Data = &user_proto.UserData{user}
	return nil
}

func (p *UserService) Filter(ctx context.Context, req *user_proto.FilterRequest, rsp *user_proto.FilterResponse) error {
	log.Info("Received User.Filter request")
	// filter users
	users, err := db.Filter(ctx, req)
	if len(users) == 0 || err != nil {
		return common.NotFound(common.UserSrv, p.Filter, err, "not found")
	}
	rsp.Data = &user_proto.UserArrData{users}
	return nil
}

func (p *UserService) Delete(ctx context.Context, req *user_proto.DeleteRequest, rsp *user_proto.DeleteResponse) error {
	log.Info("Received User.Delete request")
	//delete users by userId
	if err := db.Delete(ctx, req.UserId); err != nil {
		return common.InternalServerError(common.UserSrv, p.Delete, err, "delete error")
	}
	return nil
}

func (p *UserService) ReadByAccount(ctx context.Context, req *user_proto.ReadByAccountRequest, rsp *user_proto.ReadByAccountResponse) error {
	log.Info("Received User.ReadByAccount request")
	user, err := db.ReadByAccount(ctx, req.AccountId)
	if user == nil || err != nil {
		return common.NotFound(common.UserSrv, p.ReadByAccount, err, "not found")
	}
	rsp.Data = &user_proto.UserData{user}
	return nil
}

func (p *UserService) UpdateTokens(ctx context.Context, req *user_proto.UpdateTokenRequest, rsp *user_proto.UpdateTokenResponse) error {
	log.Info("Received User.UpdateToken request")
	if err := db.UpdateTokens(ctx, req.UserId, req.Tokens); err != nil {
		return common.InternalServerError(common.UserSrv, p.UpdateTokens, err, "update error")
	}
	return nil
}

func (p *UserService) ReadTokens(ctx context.Context, req *user_proto.ReadTokensRequest, rsp *user_proto.ReadTokensResponse) error {
	tokens, err := db.ReadTokens(ctx, req.UserIds)
	if len(tokens) == 0 || err != nil {
		return common.NotFound(common.UserSrv, p.ReadTokens, err, "not found")
	}
	rsp.Tokens = tokens
	return nil
}

func (p *UserService) SaveUserPreference(ctx context.Context, req *user_proto.SaveUserPreferenceRequest, rsp *user_proto.SaveUserPreferenceResponse) error {
	log.Info("Received User.SaveUserPreference request")

	if len(req.Preference.UserId) == 0 || req.Preference.UserId != req.UserId {
		req.Preference.UserId = req.UserId
	}

	if len(req.Preference.OrgId) == 0 || req.Preference.OrgId != req.OrgId {
		req.Preference.OrgId = req.OrgId
	}

	err := db.SaveUserPreference(ctx, req.Preference)
	if err != nil {
		return common.InternalServerError(common.UserSrv, p.SaveUserPreference, err, "user save error")
	}
	rsp.Data = &user_proto.SaveUserPreferenceResponse_Data{req.Preference}

	// reset userid & orgid
	req.Preference.UserId = req.UserId
	req.Preference.OrgId = req.OrgId

	obj := &user_proto.PreferencesMessage{
		OrgId:      req.OrgId,
		UserId:     req.UserId,
		Preference: req.Preference,
	}
	body, err := json.Marshal(obj)
	if err != nil {
		return common.InternalServerError(common.UserSrv, p.SaveUserPreference, err, "parsing error")
	}

	// publish
	if err := p.Broker.Publish(common.PREFERENCES_UPDATED, &broker.Message{Body: body}); err != nil {
		return common.InternalServerError(common.UserSrv, p.SaveUserPreference, err, "subscribe error")
	}
	return nil
}

func (p *UserService) ReadUserPreference(ctx context.Context, req *user_proto.ReadUserPreferenceRequest, rsp *user_proto.ReadUserPreferenceResponse) error {
	log.Info("Received User.ReadUserPreference request")

	preference, err := db.ReadUserPreference(ctx, req.OrgId, req.UserId)
	if preference == nil || err != nil {
		return common.NotFound(common.UserSrv, p.ReadUserPreference, err, "not found")
	}
	rsp.Data = &user_proto.ReadUserPreferenceResponse_Data{preference}
	return nil
}

func (p *UserService) ListUserFeedback(ctx context.Context, req *user_proto.ListUserFeedbackRequest, rsp *user_proto.ListUserFeedbackResponse) error {
	log.Info("Received User.ListUserFeedback request")

	feedbacks, err := db.ListUserFeedback(ctx, req.UserId)
	if len(feedbacks) == 0 || err != nil {
		return common.NotFound(common.UserSrv, p.ListUserFeedback, err, "not found")
	}
	rsp.Data = &user_proto.ListUserFeedbackResponse_Data{feedbacks}
	return nil
}

func (p *UserService) FilterUser(ctx context.Context, req *user_proto.FilterUserRequest, rsp *user_proto.FilterUserResponse) error {
	log.Info("Received User.FilterUser request")

	response, err := db.FilterUser(ctx, req)
	if len(response) == 0 || err != nil {
		return common.NotFound(common.UserSrv, p.FilterUser, err, "not found")
	}
	rsp.Data = &user_proto.FilterUserResponse_Data{response}
	return nil
}

func (p *UserService) SearchUser(ctx context.Context, req *user_proto.SearchUserRequest, rsp *user_proto.SearchUserResponse) error {
	log.Info("Received User.SearchUser request")

	response, err := db.SearchUser(ctx, req)
	if len(response) == 0 || err != nil {
		return common.NotFound(common.UserSrv, p.SearchUser, err, "not found")
	}
	rsp.Data = &user_proto.SearchUserResponse_Data{response}
	return nil
}

func (p *UserService) AutocompleteUser(ctx context.Context, req *user_proto.AutocompleteUserRequest, rsp *user_proto.AutocompleteUserResponse) error {
	log.Info("Received User.AutocompleteUser request")

	response, err := db.AutocompleteUser(ctx, req)
	if len(response) == 0 || err != nil {
		return common.NotFound(common.UserSrv, p.AutocompleteUser, err, "not found")
	}
	rsp.Data = &user_proto.AutocompleteUserResponse_Data{response}
	return nil
}

func (p *UserService) SetAccountStatus(ctx context.Context, req *account_proto.SetAccountStatusRequest, rsp *account_proto.SetAccountStatusResponse) error {
	log.Info("Received User.SetAccountStatus request")

	_, err := p.AccountClient.SetAccountStatus(ctx, req)
	if err != nil {
		return common.NotFound(common.UserSrv, p.SetAccountStatus, err, "not found")
	}
	return nil
}

func (p *UserService) GetAccountStatus(ctx context.Context, req *account_proto.GetAccountStatusRequest, rsp *account_proto.GetAccountStatusResponse) error {
	log.Info("Received User.GetAccountStatus request")

	response, err := p.AccountClient.GetAccountStatus(ctx, req)
	if err != nil {
		return common.NotFound(common.UserSrv, p.GetAccountStatus, err, "GetAccountStatus request is failed")
	}
	rsp.Data = response.Data
	return nil
}

func (p *UserService) ResetUserPassword(ctx context.Context, req *account_proto.ResetUserPasswordRequest, rsp *account_proto.ResetUserPasswordResponse) error {
	log.Info("Received User.ResetUserPassword request")

	_, err := p.AccountClient.ResetUserPassword(ctx, req)
	if err != nil {
		return common.NotFound(common.UserSrv, p.ResetUserPassword, err, "not found")
	}
	return nil
}

func (p *UserService) AddMultipleMeasurements(ctx context.Context, req *user_proto.AddMultipleMeasurementsRequest, rsp *user_proto.AddMultipleMeasurementsResponse) error {
	log.Info("Received User.AddMultipleMeasurements request")

	// fetch user for whom the measurement is being added
	//FIXME:Is this required? Not sure the purpose of this is. May we should do this for user validity check?
	user, err := db.Read(ctx, req.UserId, req.OrgId)
	if user == nil || err != nil {
		return common.NotFound(common.UserSrv, p.AddMultipleMeasurements, err, "not found")
	}

	if user == nil {
		return common.NotFound(common.UserSrv, p.AddMultipleMeasurements, err, "User not found")
	}

	// checking for valid employee
	req_employee := &team_proto.ReadEmployeeInfoRequest{req.TeamId}
	rsp_employee, err := p.TeamClient.CheckValidEmployee(ctx, req_employee)
	if err != nil {
		return common.InternalServerError(common.UserSrv, p.AddMultipleMeasurements, err, "CheckValidEmployee is failed")
	}

	if rsp_employee.Valid && rsp_employee.Employee != nil {
		//add measurements
		log.Info("Adding current measurements for user: ", user.Id)
		err = p.AddMeasurements(ctx, user.Id, req.OrgId, rsp_employee.Employee, req.Measurements)
		if err != nil {
			return common.InternalServerError(common.UserSrv, p.AddMultipleMeasurements, err, "AddMultipleMeasurements is failed")
		}
	}
	return nil
}

//internal function to add measurements
func (p *UserService) AddMeasurements(ctx context.Context, userId, orgId string, employee *team_proto.Employee, measurements []*user_proto.Measurement) error {
	log.Info("Received User.AddMeasurements request for user: ", userId)

	// create tracker marker
	for _, m := range measurements {
		log.Info("Adding measurements for marker: ", m.Marker.Id)
		rsp_track, err := p.TrackClient.CreateTrackMarker(ctx, &track_proto.CreateTrackMarkerRequest{
			Value:         m.Value,
			Unit:          m.Unit,
			MarkerId:      m.Marker.Id,
			UserId:        userId,
			TrackerMethod: m.Method,
			OrgId:         orgId,
		})
		if err != nil {
			return common.InternalServerError(common.UserSrv, p.AddMeasurements, err, "create tracker marker error")
		}
		m.OrgId = orgId
		//measured for
		m.UserId = userId
		//measured by
		m.MeasuredBy = employee.User
		//add method (if no id present)
		//get valid tracker_method if method.id is not present
		//FIXME:Refactor this so that createTrackerMarker returns the method used, in order to avoid multiple calls to the database for method
		if len(m.Method.Id) == 0 {
			req_tracker_method := &static_proto.ReadTrackerMethodRequest{NameSlug: m.Method.NameSlug}
			rsp_tracker_method, err := p.StaticClient.ReadTrackerMethod(ctx, req_tracker_method)
			if err != nil {
				return common.NotFound(common.UserSrv, p.AddMeasurements, err, "ReadTrackerMethod is failed")
			}
			m.Method = rsp_tracker_method.Data.TrackerMethod
		}
		_, err = db.AddMeasurement(ctx, m, rsp_track.Data.TrackMarker.Id)
		if err != nil {
			return common.InternalServerError(common.UserSrv, p.AddMeasurements, err, "AddMeasurements query is failed")
		}
	}
	return nil
}

func (p *UserService) GetAllMeasurementsHistory(ctx context.Context, req *user_proto.GetAllMeasurementsHistoryRequest, rsp *user_proto.GetAllMeasurementsHistoryResponse) error {
	log.Info("Received User.GetAllMeasurementsHistory request")

	measurements, err := db.GetAllMeasurementsHistory(ctx, req.UserId, req.OrgId, req.Offset, req.Limit)
	if len(measurements) == 0 || err != nil {
		return common.NotFound(common.UserSrv, p.GetAllMeasurementsHistory, err, "not found")
	}
	rsp.Data = &user_proto.GetAllMeasurementsHistoryResponse_Data{measurements}
	return nil
}

func (p *UserService) GetMeasurementsHistory(ctx context.Context, req *user_proto.GetMeasurementsHistoryRequest, rsp *user_proto.GetMeasurementsHistoryResponse) error {
	log.Info("Received User.GetMeasurementsHistory request")

	measurements, err := db.GetMeasurementsHistory(ctx, req.UserId, req.OrgId, req.MarkerId, req.Offset, req.Limit)
	if err != nil {
		common.NotFound(common.UserSrv, p.GetMeasurementsHistory, err, "GetMeasurementsHistory query is failed")
		return err
	}
	rsp.Data = &user_proto.GetMeasurementsHistoryResponse_Data{measurements}
	return nil
}

func (p *UserService) GetAllTrackedMarkers(ctx context.Context, req *user_proto.GetAllTrackedMarkersRequest, rsp *user_proto.GetAllTrackedMarkersResponse) error {
	log.Info("Received User.GetAllTrackedMarkers request")

	markers, err := db.GetAllTrackedMarkers(ctx, req.UserId, req.OrgId)
	if len(markers) == 0 || err != nil {
		return common.NotFound(common.UserSrv, p.GetAllTrackedMarkers, err, "not found")
	}
	rsp.Data = &user_proto.GetAllTrackedMarkersResponse_Data{markers}
	return nil
}

// ReadAccountToken is an internal function that only to be for user test!
func (p *UserService) ReadAccountToken(ctx context.Context, req *user_proto.ReadAccountTokenRequest, rsp *user_proto.ReadAccountTokenResponse) error {
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

//TODO: add plan
//get shared resources
func (p *UserService) GetSharedResources(ctx context.Context, req *user_proto.GetSharedResourcesRequest, rsp *user_proto.GetSharedResourcesResponse) error {
	log.Info("Received User.GetSharedResources request for user :", req.UserId)
	shared_resources := []*user_proto.SharedResourcesResponse{}
	goals, _ := common.InArray((common.BASE + common.GOAL_TYPE), req.Type)
	if goals {
		log.Info("get shared goals for user")
		rsp_goals, err := db.GetSharedGoalsForUser(ctx, req.Status, req.SharedBy, req.Query, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, "", "")
		if err != nil {
			return err
		}

		if len(rsp_goals) > 0 {
			for _, g := range rsp_goals {
				g.Type = common.BASE + common.GOAL_TYPE
				shared_resources = append(shared_resources, g)
			}
		}
	}

	challenges, _ := common.InArray((common.BASE + common.CHALLENGE_TYPE), req.Type)
	if challenges {
		log.Info("get shared challenges for user")
		rsp_challenges, err := db.GetSharedChallengesForUser(ctx, req.Status, req.SharedBy, req.Query, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, "", "")
		if err != nil {
			return err
		}

		if len(rsp_challenges) > 0 {
			for _, g := range rsp_challenges {
				g.Type = common.BASE + common.CHALLENGE_TYPE
				shared_resources = append(shared_resources, g)
			}
		}
	}

	habits, _ := common.InArray((common.BASE + common.HABIT_TYPE), req.Type)
	if habits {
		log.Info("get shared habits for user")
		rsp_habits, err := db.GetSharedHabitsForUser(ctx, req.Status, req.SharedBy, req.Query, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, "", "")
		if err != nil {
			return err
		}
		if len(rsp_habits) > 0 {
			for _, g := range rsp_habits {
				g.Type = common.BASE + common.HABIT_TYPE
				shared_resources = append(shared_resources, g)
			}
		}
	}

	surveys, _ := common.InArray((common.BASE + common.SURVEY_TYPE), req.Type)
	if surveys {
		log.Info("get shared surveys for user")
		rsp_surveys, err := db.GetSharedSurveysForUser(ctx, req.Status, req.SharedBy, req.Query, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, "", "")
		if err != nil {
			return err
		}
		if len(rsp_surveys) > 0 {
			for _, g := range rsp_surveys {
				g.Type = common.BASE + common.SURVEY_TYPE
				shared_resources = append(shared_resources, g)
			}
		}
	}

	contents, _ := common.InArray((common.BASE + common.CONTENT_TYPE), req.Type)
	if contents {
		log.Info("get shared contents for user")
		rsp_contents, err := db.GetSharedContentsForUser(ctx, req.Status, req.SharedBy, req.Query, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, "", "")
		if err != nil {
			return err
		}
		if len(rsp_contents) > 0 {
			for _, g := range rsp_contents {
				g.Type = common.BASE + common.CONTENT_TYPE
				shared_resources = append(shared_resources, g)
			}
		}
	}
	rsp.Data = &user_proto.GetSharedResourcesResponse_Data{SharedResources: shared_resources}
	return nil
}

//TODO: add plan
//fetches all the resources (goals, challenges, habits ... ) that can be shared with a specific user_id filtering out the ones that are already shared or joined (through discover)
func (p *UserService) GetAllShareableResources(ctx context.Context, req *user_proto.GetShareableResourcesRequest, rsp *user_proto.GetShareableResourcesResponse) error {
	log.Info("Received User.GetAllShareableResources request for user :", req.UserId)
	resources := []*user_proto.SharedResourcesResponse{}
	goals, _ := common.InArray((common.BASE + common.GOAL_TYPE), req.Type)
	if goals {
		log.Info("search goals for user: ", req.UserId)

		req_goals := &user_proto.AllGoalResponseRequest{UserId: req.UserId, OrgId: req.OrgId, TeamId: req.TeamId, Offset: req.Offset, Limit: req.Limit, CreatedBy: req.CreatedBy, Query: req.Query}
		rsp_goals, err := p.BehaviourClient.AllGoalResponse(ctx, req_goals)
		if err != nil {
			return err
		}
		if len(rsp_goals.Data.Goals) > 0 {
			for _, g := range rsp_goals.Data.Goals {
				resource := &user_proto.SharedResourcesResponse{
					Type:       common.BASE + common.GOAL_TYPE,
					OrgId:      g.OrgId,
					ResourceId: g.Id,
					Image:      g.Image,
					Title:      g.Title,
					Summary:    g.Summary,
					SharedBy:   g.SharedBy,
					Target:     g.Target,
					Duration:   g.Duration,
				}
				resources = append(resources, resource)
			}
		}
	}

	challenges, _ := common.InArray((common.BASE + common.CHALLENGE_TYPE), req.Type)
	if challenges {
		log.Info("get shareable challenges for user: ", req.UserId)

		req_challenges := &user_proto.AllChallengeResponseRequest{UserId: req.UserId, OrgId: req.OrgId, TeamId: req.TeamId, Offset: req.Offset, Limit: req.Limit, CreatedBy: req.CreatedBy, Query: req.Query}
		rsp_challenges, err := p.BehaviourClient.AllChallengeResponse(ctx, req_challenges)
		if err != nil {
			return err
		}
		if len(rsp_challenges.Data.Challenges) > 0 {
			for _, g := range rsp_challenges.Data.Challenges {
				resource := &user_proto.SharedResourcesResponse{
					Type:       common.BASE + common.CHALLENGE_TYPE,
					OrgId:      g.OrgId,
					ResourceId: g.Id,
					Image:      g.Image,
					Title:      g.Title,
					Summary:    g.Summary,
					SharedBy:   g.SharedBy,
					Target:     g.Target,
					Duration:   g.Duration,
				}
				resources = append(resources, resource)
			}
		}
	}

	habits, _ := common.InArray((common.BASE + common.HABIT_TYPE), req.Type)
	if habits {
		log.Info("get shareable habits for user: ", req.UserId)

		req_habits := &user_proto.AllHabitResponseRequest{UserId: req.UserId, OrgId: req.OrgId, TeamId: req.TeamId, Offset: req.Offset, Limit: req.Limit, CreatedBy: req.CreatedBy, Query: req.Query}
		rsp_habits, err := p.BehaviourClient.AllHabitResponse(ctx, req_habits)
		if err != nil {
			return err
		}
		if len(rsp_habits.Data.Habits) > 0 {
			for _, g := range rsp_habits.Data.Habits {
				resource := &user_proto.SharedResourcesResponse{
					Type:       common.BASE + common.HABIT_TYPE,
					OrgId:      g.OrgId,
					ResourceId: g.Id,
					Image:      g.Image,
					Title:      g.Title,
					Summary:    g.Summary,
					SharedBy:   g.SharedBy,
					Target:     g.Target,
					Duration:   g.Duration,
				}
				resources = append(resources, resource)
			}
		}
	}

	surveys, _ := common.InArray((common.BASE + common.SURVEY_TYPE), req.Type)
	if surveys {
		log.Info("get shareable surveys for user: ", req.UserId)

		req_surveys := &user_proto.GetShareableSurveyRequest{UserId: req.UserId, OrgId: req.OrgId, TeamId: req.TeamId, Offset: req.Offset, Limit: req.Limit, CreatedBy: req.CreatedBy, Query: req.Query}
		rsp_surveys, err := p.SurveyClient.GetShareableSurveys(ctx, req_surveys)
		if err != nil {
			return err
		}
		if len(rsp_surveys.Data.Surveys) > 0 {
			for _, g := range rsp_surveys.Data.Surveys {
				resource := &user_proto.SharedResourcesResponse{
					Type:         common.BASE + common.SURVEY_TYPE,
					OrgId:        g.OrgId,
					ResourceId:   g.Id,
					Title:        g.Title,
					Summary:      g.Summary,
					SharedBy:     g.SharedBy,
					Count:        g.Count,
					ResponseTime: g.ResponseTime,
				}
				resources = append(resources, resource)
			}
		}
	}

	article, _ := common.InArray((common.BASE + common.CONTENT_ARTICLE_TYPE), req.Type)
	video, _ := common.InArray((common.BASE + common.CONTENT_VIDEO_TYPE), req.Type)
	recipe, _ := common.InArray((common.BASE + common.CONTENT_RECIPE_TYPE), req.Type)
	exercise, _ := common.InArray((common.BASE + common.CONTENT_EXERCISE_TYPE), req.Type)
	if article || video || recipe || exercise {
		log.Info("get shareable contents for user: ", req.UserId)

		req_contents := &user_proto.GetShareableContentRequest{UserId: req.UserId, OrgId: req.OrgId, TeamId: req.TeamId, Offset: req.Offset, Limit: req.Limit}
		rsp_contents, err := p.ContentClient.GetShareableContents(ctx, req_contents)
		if err != nil {
			return err
		}
		if len(rsp_contents.Data.Contents) > 0 {
			for _, g := range rsp_contents.Data.Contents {
				resource := &user_proto.SharedResourcesResponse{
					Type:       g.Item.TypeUrl,
					OrgId:      g.OrgId,
					ResourceId: g.Id,
					Image:      g.Image,
					Title:      g.Title,
					Summary:    g.Summary,
					SharedBy:   g.SharedBy,
				}
				resources = append(resources, resource)
			}
		}
	}

	rsp.Data = &user_proto.GetShareableResourcesResponse_Data{Resources: resources}
	return nil
}

//TODO: add survey, plan and content
//share resources (goals, challenges, habits ... ) with a specific user_id
func (p *UserService) ShareResources(ctx context.Context, req *user_proto.ShareResourcesRequest, rsp *user_proto.ShareResourcesResponse) error {
	log.Info("Received User.ShareResources for user :", req.UserId)
	if req.Shares[common.BASE+common.GOAL_TYPE] != nil && len(req.Shares[common.BASE+common.GOAL_TYPE].Resource) > 0 {
		log.Info("share goals with user: ", req.UserId)

		for _, g := range req.Shares[common.BASE+common.GOAL_TYPE].Resource {
			goals := []*behaviour_proto.Goal{}
			users := []*behaviour_proto.TargetedUser{}
			goal := &behaviour_proto.Goal{
				Id: g.ResourceId,
			}
			goals = append(goals, goal)
			user := &behaviour_proto.TargetedUser{
				User:             &user_proto.User{Id: req.UserId},
				CurrentValue:     g.CurrentValue,
				ExpectedProgress: g.ExpectedProgress,
				Unit:             g.Unit,
			}
			users = append(users, user)
			//for sharegoal, we use the UserId of the sharer (employee) hence UserId = req.TeamId (employee)
			req_share_goal := &behaviour_proto.ShareGoalRequest{Goals: goals, UserId: req.TeamId, OrgId: req.OrgId, Users: users}
			if _, err := p.BehaviourClient.ShareGoal(ctx, req_share_goal); err != nil {
				return err
			}
		}
	}
	if req.Shares[common.BASE+common.CHALLENGE_TYPE] != nil && len(req.Shares[common.BASE+common.CHALLENGE_TYPE].Resource) > 0 {
		log.Info("share challenges with user: ", req.UserId)

		for _, c := range req.Shares[common.BASE+common.CHALLENGE_TYPE].Resource {
			challenges := []*behaviour_proto.Challenge{}
			users := []*behaviour_proto.TargetedUser{}
			challenge := &behaviour_proto.Challenge{
				Id: c.ResourceId,
			}
			challenges = append(challenges, challenge)
			user := &behaviour_proto.TargetedUser{
				User:             &user_proto.User{Id: req.UserId},
				CurrentValue:     c.CurrentValue,
				ExpectedProgress: c.ExpectedProgress,
				Unit:             c.Unit,
			}
			users = append(users, user)
			//for sharechallenge, we use the UserId of the sharer (employee) hence UserId = req.TeamId (employee)
			req_share_challenge := &behaviour_proto.ShareChallengeRequest{Challenges: challenges, UserId: req.TeamId, OrgId: req.OrgId, Users: users}
			if _, err := p.BehaviourClient.ShareChallenge(ctx, req_share_challenge); err != nil {
				return err
			}
		}
	}

	if req.Shares[common.BASE+common.HABIT_TYPE] != nil && len(req.Shares[common.BASE+common.HABIT_TYPE].Resource) > 0 {
		log.Info("share habits with user: ", req.UserId)

		for _, h := range req.Shares[common.BASE+common.HABIT_TYPE].Resource {
			habits := []*behaviour_proto.Habit{}
			users := []*behaviour_proto.TargetedUser{}
			habit := &behaviour_proto.Habit{
				Id: h.ResourceId,
			}
			habits = append(habits, habit)
			user := &behaviour_proto.TargetedUser{
				User:             &user_proto.User{Id: req.UserId},
				CurrentValue:     h.CurrentValue,
				ExpectedProgress: h.ExpectedProgress,
				Unit:             h.Unit,
			}
			users = append(users, user)
			//for sharehabit, we use the UserId of the sharer (employee) hence UserId = req.TeamId (employee)
			req_share_habit := &behaviour_proto.ShareHabitRequest{Habits: habits, UserId: req.TeamId, OrgId: req.OrgId, Users: users}
			if _, err := p.BehaviourClient.ShareHabit(ctx, req_share_habit); err != nil {
				return err
			}
		}
	}

	if req.Shares[common.BASE+common.SURVEY_TYPE] != nil && len(req.Shares[common.BASE+common.SURVEY_TYPE].Resource) > 0 {
		log.Info("share survey with user: ", req.UserId)

		surveys := []*survey_proto.Survey{}
		for _, h := range req.Shares[common.BASE+common.SURVEY_TYPE].Resource {
			survey := &survey_proto.Survey{
				Id: h.ResourceId,
			}
			surveys = append(surveys, survey)
		}
		users := []*user_proto.User{}
		users = append(users, &user_proto.User{Id: req.UserId})
		//for sharesurvey, we use the UserId of the sharer (employee) hence UserId = req.TeamId (employee)
		req_share_survey := &survey_proto.ShareSurveyRequest{Surveys: surveys, UserId: req.TeamId, OrgId: req.OrgId, Users: users}
		if _, err := p.SurveyClient.ShareSurvey(ctx, req_share_survey); err != nil {
			return err
		}
	}

	if req.Shares[common.BASE+common.CONTENT_TYPE] != nil && len(req.Shares[common.BASE+common.CONTENT_TYPE].Resource) > 0 {
		log.Info("share content with users: ", req.UserId)

		contents := []*content_proto.Content{}
		for _, h := range req.Shares[common.BASE+common.CONTENT_TYPE].Resource {
			content := &content_proto.Content{
				Id: h.ResourceId,
			}
			contents = append(contents, content)
		}
		users := []*user_proto.User{}
		users = append(users, &user_proto.User{Id: req.UserId})
		//for sharecontent, we use the UserId of the sharer (employee) hence UserId = req.TeamId (employee)
		req_share_content := &content_proto.ShareContentRequest{Contents: contents, UserId: req.TeamId, OrgId: req.OrgId, Users: users}
		if _, err := p.ContentClient.ShareContent(ctx, req_share_content); err != nil {
			return err
		}
	}
	return nil
}

//TODO: add plan and content
//share multiple resources (goals, challenges, habits ... ) with multiple users
func (p *UserService) ShareMultipleResources(ctx context.Context, req *user_proto.ShareMultipleResourcesRequest, rsp *user_proto.ShareResourcesResponse) error {
	log.Info("Received User.ShareMultipleResources with users", req.Users)

	targeted_users := []*behaviour_proto.TargetedUser{}
	users := []*user_proto.User{}
	if len(req.Users) > 0 {
		for _, u := range req.Users {
			targeted_user := &behaviour_proto.TargetedUser{
				User: &user_proto.User{Id: u},
			}

			targeted_users = append(targeted_users, targeted_user)
			users = append(users, &user_proto.User{Id: u})
		}

		if req.Shares[common.BASE+common.GOAL_TYPE] != nil && len(req.Shares[common.BASE+common.GOAL_TYPE].Resource) > 0 {
			log.Info("share goals with user: ", req.Users)
			goals := []*behaviour_proto.Goal{}
			for _, g := range req.Shares[common.BASE+common.GOAL_TYPE].Resource {
				goal := &behaviour_proto.Goal{
					Id: g,
				}
				goals = append(goals, goal)
			}

			//for sharegoal, we use the UserId of the sharer (employee) hence UserId = req.TeamId (employee)
			req_share_goal := &behaviour_proto.ShareGoalRequest{Goals: goals, UserId: req.TeamId, OrgId: req.OrgId, Users: targeted_users}
			if _, err := p.BehaviourClient.ShareGoal(ctx, req_share_goal); err != nil {
				return err
			}
		}
		if req.Shares[common.BASE+common.CHALLENGE_TYPE] != nil && len(req.Shares[common.BASE+common.CHALLENGE_TYPE].Resource) > 0 {
			log.Info("share challenges with user: ", req.Users)

			challenges := []*behaviour_proto.Challenge{}
			for _, c := range req.Shares[common.BASE+common.CHALLENGE_TYPE].Resource {
				challenge := &behaviour_proto.Challenge{
					Id: c,
				}
				challenges = append(challenges, challenge)
			}
			//for sharechallenge, we use the UserId of the sharer (employee) hence UserId = req.TeamId (employee)
			req_share_challenge := &behaviour_proto.ShareChallengeRequest{Challenges: challenges, UserId: req.TeamId, OrgId: req.OrgId, Users: targeted_users}
			if _, err := p.BehaviourClient.ShareChallenge(ctx, req_share_challenge); err != nil {
				return err
			}
		}

		if req.Shares[common.BASE+common.HABIT_TYPE] != nil && len(req.Shares[common.BASE+common.HABIT_TYPE].Resource) > 0 {
			log.Info("share habits with user: ", req.Users)

			habits := []*behaviour_proto.Habit{}
			for _, h := range req.Shares[common.BASE+common.HABIT_TYPE].Resource {
				habit := &behaviour_proto.Habit{
					Id: h,
				}
				habits = append(habits, habit)
			}
			//for sharehabit, we use the UserId of the sharer (employee) hence UserId = req.TeamId (employee)
			req_share_habit := &behaviour_proto.ShareHabitRequest{Habits: habits, UserId: req.TeamId, OrgId: req.OrgId, Users: targeted_users}
			if _, err := p.BehaviourClient.ShareHabit(ctx, req_share_habit); err != nil {
				return err
			}
		}

		if req.Shares[common.BASE+common.SURVEY_TYPE] != nil && len(req.Shares[common.BASE+common.SURVEY_TYPE].Resource) > 0 {
			log.Info("share survey with users: ", req.Users)

			surveys := []*survey_proto.Survey{}
			for _, h := range req.Shares[common.BASE+common.SURVEY_TYPE].Resource {
				survey := &survey_proto.Survey{
					Id: h,
				}
				surveys = append(surveys, survey)
			}
			//for sharesurvey, we use the UserId of the sharer (employee) hence UserId = req.TeamId (employee)
			req_share_survey := &survey_proto.ShareSurveyRequest{Surveys: surveys, UserId: req.TeamId, OrgId: req.OrgId, Users: users}
			if _, err := p.SurveyClient.ShareSurvey(ctx, req_share_survey); err != nil {
				return err
			}
		}

		if req.Shares[common.BASE+common.CONTENT_TYPE] != nil && len(req.Shares[common.BASE+common.CONTENT_TYPE].Resource) > 0 {
			log.Info("share content with users: ", req.Users)

			contents := []*content_proto.Content{}
			for _, h := range req.Shares[common.BASE+common.CONTENT_TYPE].Resource {
				content := &content_proto.Content{
					Id: h,
				}
				contents = append(contents, content)
			}
			//for sharecontent, we use the UserId of the sharer (employee) hence UserId = req.TeamId (employee)
			req_share_content := &content_proto.ShareContentRequest{Contents: contents, UserId: req.TeamId, OrgId: req.OrgId, Users: users}
			if _, err := p.ContentClient.ShareContent(ctx, req_share_content); err != nil {
				return err
			}
		}
	}
	return nil
}
