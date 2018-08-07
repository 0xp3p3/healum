package handler

import (
	"context"
	"encoding/json"
	"fmt"
	account_proto "server/account-srv/proto/account"
	"server/behaviour-srv/db"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	pubsub_proto "server/static-srv/proto/pubsub"
	static_proto "server/static-srv/proto/static"
	team_proto "server/team-srv/proto/team"
	user_proto "server/user-srv/proto/user"
	"strconv"

	"github.com/micro/go-micro/broker"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type BehaviourService struct {
	Broker        broker.Broker
	AccountClient account_proto.AccountServiceClient
	StaticClient  static_proto.StaticServiceClient
	KvClient      kv_proto.KvServiceClient
	TeamClient    team_proto.TeamServiceClient
}

func (p *BehaviourService) AddSetbacks(ctx context.Context, setbacks []*static_proto.Setback) error {
	for _, setback := range setbacks {
		req_setback := &static_proto.CreateSetbackRequest{
			Setback: setback,
		}
		if _, err := p.StaticClient.CreateSetback(ctx, req_setback); err != nil {
			return common.InternalServerError(common.BehaviourSrv, p.AddSetbacks, err, "server error")
		}
	}
	return nil
}

func (p *BehaviourService) AllGoals(ctx context.Context, req *behaviour_proto.AllGoalsRequest, rsp *behaviour_proto.AllGoalsResponse) error {
	log.Info("Received Behaviour.AllGoals request")
	goals, err := db.AllGoals(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(goals) == 0 || err != nil {
		return common.NotFound(common.BehaviourSrv, p.AllGoals, err, "goal not found")
	}
	rsp.Data = &behaviour_proto.GoalArrData{goals}
	return nil
}

func (p *BehaviourService) CreateGoal(ctx context.Context, req *behaviour_proto.CreateGoalRequest, rsp *behaviour_proto.CreateGoalResponse) error {
	log.Info("Received Behaviour.CreateGoal request")
	if len(req.Goal.Title) == 0 {
		return common.InternalServerError(common.BehaviourSrv, p.CreateGoal, nil, "title empty")
	}
	if len(req.Goal.Id) == 0 {
		req.Goal.Id = uuid.NewUUID().String()
	}
	// create goal
	err := db.CreateGoal(ctx, req.Goal)
	if err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.CreateGoal, err, "server error")
	}
	// share goal with user
	if req.Goal.Users != nil && len(req.Goal.Users) > 0 {
		req_share := &behaviour_proto.ShareGoalRequest{
			Goals:  []*behaviour_proto.Goal{req.Goal},
			Users:  req.Goal.Users,
			UserId: req.UserId,
			OrgId:  req.OrgId,
		}
		rsp_share := &behaviour_proto.ShareGoalResponse{}
		if err := p.ShareGoal(ctx, req_share, rsp_share); err != nil {
			return common.InternalServerError(common.BehaviourSrv, p.CreateGoal, err, "share error")
		}
	}

	// save setbacks
	if err := p.AddSetbacks(ctx, req.Goal.Setbacks); err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.CreateGoal, err, "add error")
	}

	// add tags cloud
	if len(req.Goal.Tags) > 0 {
		if _, err := p.KvClient.TagsCloud(context.TODO(), &kv_proto.TagsCloudRequest{
			Index:  common.CLOUD_TAGS_INDEX,
			OrgId:  req.Goal.OrgId,
			Object: common.GOAL,
			Tags:   req.Goal.Tags,
		}); err != nil {
			return common.InternalServerError(common.BehaviourSrv, p.CreateGoal, err, "tag error")
		}
	}
	rsp.Data = &behaviour_proto.GoalData{req.Goal}
	return nil
}

func (p *BehaviourService) ReadGoal(ctx context.Context, req *behaviour_proto.ReadGoalRequest, rsp *behaviour_proto.ReadGoalResponse) error {
	log.Info("Received Behaviour.ReadGoal request")
	goal, err := db.ReadGoal(ctx, req.GoalId, req.OrgId, req.TeamId)
	if goal == nil || err != nil {
		return common.NotFound(common.BehaviourSrv, p.ReadGoal, err, "goal not found")
	}
	rsp.Data = &behaviour_proto.GoalData{goal}
	return nil
}

func (p *BehaviourService) DeleteGoal(ctx context.Context, req *behaviour_proto.DeleteGoalRequest, rsp *behaviour_proto.DeleteGoalResponse) error {
	log.Info("Received Behaviour.DeleteGoal request")
	if err := db.DeleteGoal(ctx, req.GoalId, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.DeleteGoal, err, "server error")
	}
	return nil
}

func (p *BehaviourService) AllHabits(ctx context.Context, req *behaviour_proto.AllHabitsRequest, rsp *behaviour_proto.AllHabitsResponse) error {
	log.Info("Received Behaviour.AllHabits request")
	habits, err := db.AllHabits(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(habits) == 0 || err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.AllHabits, err, "habit not found")
	}
	rsp.Data = &behaviour_proto.HabitArrData{habits}
	return nil
}

func (p *BehaviourService) CreateHabit(ctx context.Context, req *behaviour_proto.CreateHabitRequest, rsp *behaviour_proto.CreateHabitResponse) error {
	log.Info("Received Behaviour.CreateHabit request")
	if len(req.Habit.Title) == 0 {
		return common.InternalServerError(common.BehaviourSrv, p.CreateHabit, nil, "title empty")
	}
	if len(req.Habit.Id) == 0 {
		req.Habit.Id = uuid.NewUUID().String()
	}
	//create
	if err := db.CreateHabit(ctx, req.Habit); err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.CreateHabit, err, "create error")
	}
	// share habit with user
	if req.Habit.Users != nil && len(req.Habit.Users) > 0 {
		req_share := &behaviour_proto.ShareHabitRequest{
			Habits: []*behaviour_proto.Habit{req.Habit},
			Users:  req.Habit.Users,
			UserId: req.UserId,
			OrgId:  req.OrgId,
		}
		rsp_share := &behaviour_proto.ShareHabitResponse{}
		if err := p.ShareHabit(ctx, req_share, rsp_share); err != nil {
			return common.InternalServerError(common.BehaviourSrv, p.CreateHabit, err, "share error")
		}
	}

	// save setbacks
	if err := p.AddSetbacks(ctx, req.Habit.Setbacks); err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.CreateHabit, err, "add error")
	}

	// add tags cloud
	if len(req.Habit.Tags) > 0 {
		if _, err := p.KvClient.TagsCloud(context.TODO(), &kv_proto.TagsCloudRequest{
			Index:  common.CLOUD_TAGS_INDEX,
			OrgId:  req.Habit.OrgId,
			Object: common.HABIT,
			Tags:   req.Habit.Tags,
		}); err != nil {
			return common.InternalServerError(common.BehaviourSrv, p.CreateHabit, err, "tag error")
		}
	}

	rsp.Data = &behaviour_proto.HabitData{req.Habit}
	return nil
}

func (p *BehaviourService) ReadHabit(ctx context.Context, req *behaviour_proto.ReadHabitRequest, rsp *behaviour_proto.ReadHabitResponse) error {
	log.Info("Received Behaviour.ReadHabit request")
	habit, err := db.ReadHabit(ctx, req.HabitId, req.OrgId, req.TeamId)
	if habit == nil || err != nil {
		return common.NotFound(common.BehaviourSrv, p.ReadHabit, err, "habit not found")
	}
	rsp.Data = &behaviour_proto.HabitData{habit}
	return nil
}

func (p *BehaviourService) DeleteHabit(ctx context.Context, req *behaviour_proto.DeleteHabitRequest, rsp *behaviour_proto.DeleteHabitResponse) error {
	log.Info("Received Behaviour.DeleteHabit request")
	if err := db.DeleteHabit(ctx, req.HabitId, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.DeleteHabit, err, "server error")
	}
	return nil
}

func (p *BehaviourService) AllChallenges(ctx context.Context, req *behaviour_proto.AllChallengesRequest, rsp *behaviour_proto.AllChallengesResponse) error {
	log.Info("Received Behaviour.AllChallenges request")
	challenges, err := db.AllChallenges(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(challenges) == 0 || err != nil {
		return common.NotFound(common.BehaviourSrv, p.AllChallenges, err, "not found")
	}
	rsp.Data = &behaviour_proto.ChallengeArrData{challenges}
	return nil
}

func (p *BehaviourService) CreateChallenge(ctx context.Context, req *behaviour_proto.CreateChallengeRequest, rsp *behaviour_proto.CreateChallengeResponse) error {
	log.Info("Received Behaviour.CreateChallenge request")
	if len(req.Challenge.Title) == 0 {
		return common.InternalServerError(common.BehaviourSrv, p.CreateChallenge, nil, "title empty")
	}
	if len(req.Challenge.Id) == 0 {
		req.Challenge.Id = uuid.NewUUID().String()
	}
	// create
	err := db.CreateChallenge(ctx, req.Challenge)
	if err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.CreateChallenge, err, "server error")
	}
	// share habit with user
	if req.Challenge.Users != nil && len(req.Challenge.Users) > 0 {
		req_share := &behaviour_proto.ShareChallengeRequest{
			Challenges: []*behaviour_proto.Challenge{req.Challenge},
			Users:      req.Challenge.Users,
			UserId:     req.UserId,
			OrgId:      req.OrgId,
		}
		rsp_share := &behaviour_proto.ShareChallengeResponse{}
		if err := p.ShareChallenge(ctx, req_share, rsp_share); err != nil {
			return common.InternalServerError(common.BehaviourSrv, p.CreateChallenge, err, "share error")
		}
	}

	// save setbacks
	if err := p.AddSetbacks(ctx, req.Challenge.Setbacks); err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.CreateChallenge, err, "add error")
	}

	// add tags cloud
	if len(req.Challenge.Tags) > 0 {
		if _, err := p.KvClient.TagsCloud(context.TODO(), &kv_proto.TagsCloudRequest{
			Index:  common.CLOUD_TAGS_INDEX,
			OrgId:  req.Challenge.OrgId,
			Object: common.CHALLENGE,
			Tags:   req.Challenge.Tags,
		}); err != nil {
			return common.InternalServerError(common.BehaviourSrv, p.CreateChallenge, err, "tag error")
			return err
		}
	}

	rsp.Data = &behaviour_proto.ChallengeData{req.Challenge}
	return nil
}

func (p *BehaviourService) ReadChallenge(ctx context.Context, req *behaviour_proto.ReadChallengeRequest, rsp *behaviour_proto.ReadChallengeResponse) error {
	log.Info("Received Behaviour.ReadChallenge request")
	challenge, err := db.ReadChallenge(ctx, req.ChallengeId, req.OrgId, req.TeamId)
	if challenge == nil || err != nil {
		return common.NotFound(common.BehaviourSrv, p.ReadChallenge, err, "challenge not found")
	}
	rsp.Data = &behaviour_proto.ChallengeData{challenge}
	return nil
}

func (p *BehaviourService) DeleteChallenge(ctx context.Context, req *behaviour_proto.DeleteChallengeRequest, rsp *behaviour_proto.DeleteChallengeResponse) error {
	log.Info("Received Behaviour.DeleteChallenge request")
	if err := db.DeleteChallenge(ctx, req.ChallengeId, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.DeleteChallenge, err, "server error")
	}
	return nil
}

func (p *BehaviourService) Filter(ctx context.Context, req *behaviour_proto.FilterRequest, rsp *behaviour_proto.FilterResponse) error {
	log.Info("Received Behaviour.Filter request")
	data, err := db.Filter(ctx, req)
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.Filter, err, "not found")
	}
	rsp.Data = data
	return nil
}

func (p *BehaviourService) Search(ctx context.Context, req *behaviour_proto.SearchRequest, rsp *behaviour_proto.SearchResponse) error {
	log.Info("Received Behaviour.Search request")
	data, err := db.Search(ctx, req)
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.Search, err, "not found")
	}
	rsp.Data = data
	return nil
}

func (p *BehaviourService) ShareGoal(ctx context.Context, req *behaviour_proto.ShareGoalRequest, rsp *behaviour_proto.ShareGoalResponse) error {
	log.Info("Received Behaviour.ShareGoal request")

	if len(req.Goals) == 0 {
		return common.Forbidden(common.BehaviourSrv, p.ShareGoal, nil, "goals empty")
	}
	if len(req.Users) == 0 {
		return common.Forbidden(common.BehaviourSrv, p.ShareGoal, nil, "users empty")
	}

	// checking valid sharedby (employee)
	req_employee := &team_proto.ReadEmployeeInfoRequest{req.UserId}
	rsp_employee, err := p.TeamClient.CheckValidEmployee(ctx, req_employee)
	if err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.ShareGoal, err, "CheckValidEmployee is failed")
	}
	if rsp_employee.Valid && rsp_employee.Employee != nil {
		userids, err := db.ShareGoal(ctx, req.Goals, req.Users, rsp_employee.Employee.User, req.OrgId)
		if err != nil {
			return common.InternalServerError(common.BehaviourSrv, p.ShareGoal, err, "parsing error")
		}
		// send a notification to the users
		if len(userids) > 0 {
			message := fmt.Sprintf(common.MSG_NEW_GOAL_SHARE, rsp_employee.Employee.User.Firstname)
			alert := &pubsub_proto.Alert{
				Title: fmt.Sprintf("New %v", common.GOAL),
				Body:  message,
			}
			data := map[string]string{}
			//get current badge count here for user
			data[common.BASE+common.GOAL_TYPE] = strconv.Itoa(len(req.Goals))
			p.sendShareNotification(userids, message, alert, data)
		}
	}

	return nil
}

func (p *BehaviourService) ShareChallenge(ctx context.Context, req *behaviour_proto.ShareChallengeRequest, rsp *behaviour_proto.ShareChallengeResponse) error {
	log.Info("Received Behaviour.ShareChallenge request")

	if len(req.Challenges) == 0 {
		return common.Forbidden(common.BehaviourSrv, p.ShareChallenge, nil, "challenges empty")
	}
	if len(req.Users) == 0 {
		return common.Forbidden(common.BehaviourSrv, p.ShareChallenge, nil, "users empty")
	}
	// checking valid sharedby (employee)
	req_employee := &team_proto.ReadEmployeeInfoRequest{req.UserId}
	rsp_employee, err := p.TeamClient.CheckValidEmployee(ctx, req_employee)
	if err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.ShareChallenge, err, "CheckValidEmployee is failed")
	}
	if rsp_employee.Valid && rsp_employee.Employee != nil {
		userids, err := db.ShareChallenge(ctx, req.Challenges, req.Users, rsp_employee.Employee.User, req.OrgId)
		if err != nil {
			return common.InternalServerError(common.BehaviourSrv, p.ShareChallenge, err, "parsing error")
		}
		// send a notification to the users
		if len(userids) > 0 {
			message := fmt.Sprintf(common.MSG_NEW_CHALLENGE_SHARE, rsp_employee.Employee.User.Firstname)
			alert := &pubsub_proto.Alert{
				Title: fmt.Sprintf("New %v", common.CHALLENGE),
				Body:  message,
			}
			data := map[string]string{}
			//get current badge count here for user
			data[common.BASE+common.CHALLENGE_TYPE] = strconv.Itoa(len(req.Challenges))
			p.sendShareNotification(userids, message, alert, data)
		}
	}
	return nil
}

func (p *BehaviourService) ShareHabit(ctx context.Context, req *behaviour_proto.ShareHabitRequest, rsp *behaviour_proto.ShareHabitResponse) error {
	log.Info("Received Behaviour.ShareHabit request")

	if len(req.Habits) == 0 {
		return common.Forbidden(common.BehaviourSrv, p.ShareHabit, nil, "habits empty")
	}
	if len(req.Users) == 0 {
		return common.Forbidden(common.BehaviourSrv, p.ShareHabit, nil, "users empty")
	}
	// checking valid sharedby (employee)
	req_employee := &team_proto.ReadEmployeeInfoRequest{req.UserId}
	rsp_employee, err := p.TeamClient.CheckValidEmployee(ctx, req_employee)
	if err != nil {
		return common.InternalServerError(common.BehaviourSrv, p.ShareHabit, err, "CheckValidEmployee is failed")
	}
	if rsp_employee.Valid && rsp_employee.Employee != nil {
		userids, err := db.ShareHabit(ctx, req.Habits, req.Users, rsp_employee.Employee.User, req.OrgId)
		if err != nil {
			return common.InternalServerError(common.BehaviourSrv, p.ShareHabit, err, "habit share is failed")
		}
		// send a notification to the users
		if len(userids) > 0 {
			message := fmt.Sprintf(common.MSG_NEW_HABIT_SHARE, rsp_employee.Employee.User.Firstname)
			alert := &pubsub_proto.Alert{
				Title: fmt.Sprintf("New %v", common.HABIT),
				Body:  message,
			}
			data := map[string]string{}
			//get current badge count here for user
			data[common.BASE+common.HABIT_TYPE] = strconv.Itoa(len(req.Habits))
			p.sendShareNotification(userids, message, alert, data)
		}
	}
	return nil
}

func (p *BehaviourService) AutocompleteGoalSearch(ctx context.Context, req *behaviour_proto.AutocompleteSearchRequest, rsp *behaviour_proto.AutocompleteSearchResponse) error {
	log.Info("Received Behaviour.AutocompleteSearch request")
	response, err := db.AutocompleteGoalSearch(ctx, req.Title)
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.AutocompleteGoalSearch, err, "not found")
	}
	rsp.Data = &behaviour_proto.AutocompleteSearchResponse_Data{response}
	return nil
}

func (p *BehaviourService) AutocompleteChallengeSearch(ctx context.Context, req *behaviour_proto.AutocompleteSearchRequest, rsp *behaviour_proto.AutocompleteSearchResponse) error {
	log.Info("Received Behaviour.AutocompleteSearch request")
	response, err := db.AutocompleteChallengeSearch(ctx, req.Title)
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.AutocompleteChallengeSearch, err, "not found")
	}
	rsp.Data = &behaviour_proto.AutocompleteSearchResponse_Data{response}
	return nil
}

func (p *BehaviourService) AutocompleteHabitSearch(ctx context.Context, req *behaviour_proto.AutocompleteSearchRequest, rsp *behaviour_proto.AutocompleteSearchResponse) error {
	log.Info("Received Behaviour.AutocompleteSearch request")
	response, err := db.AutocompleteHabitSearch(ctx, req.Title)
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.AutocompleteChallengeSearch, err, "not found")
	}
	rsp.Data = &behaviour_proto.AutocompleteSearchResponse_Data{response}
	return nil
}

func (p *BehaviourService) GetSharedGoal(ctx context.Context, req *behaviour_proto.GetSharedGoalRequest, rsp *behaviour_proto.GetSharedGoalResponse) error {
	log.Info("Received Behaviour.GetSharedGoal request")

	goal, err := db.GetSharedGoal(ctx, req.UserId, req.GoalId)
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.GetSharedGoal, err, "not found")
	}
	rsp.Data = &behaviour_proto.GetSharedGoalResponse_Data{goal}
	return nil
}

func (p *BehaviourService) GetSharedChallenge(ctx context.Context, req *behaviour_proto.GetSharedChallengeRequest, rsp *behaviour_proto.GetSharedChallengeResponse) error {
	log.Info("Received Behaviour.GetSharedChallenge request")

	challenge, err := db.GetSharedChallenge(ctx, req.UserId, req.ChallengeId)
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.GetSharedChallenge, err, "not found")
	}
	rsp.Data = &behaviour_proto.GetSharedChallengeResponse_Data{challenge}
	return nil
}

func (p *BehaviourService) GetSharedHabit(ctx context.Context, req *behaviour_proto.GetSharedHabitRequest, rsp *behaviour_proto.GetSharedHabitResponse) error {
	log.Info("Received Behaviour.GetSharedHabit request")

	habit, err := db.GetSharedHabit(ctx, req.UserId, req.HabitId)
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.GetSharedHabit, err, "not found")
	}
	rsp.Data = &behaviour_proto.GetSharedHabitResponse_Data{habit}
	return nil
}

func (p *BehaviourService) GetTopTags(ctx context.Context, req *behaviour_proto.GetTopTagsRequest, rsp *behaviour_proto.GetTopTagsResponse) error {
	log.Info("Received Behaviour.GetTopTags request")

	rsp_tags, err := p.KvClient.GetTopTags(ctx, &kv_proto.GetTopTagsRequest{
		Index:  common.CLOUD_TAGS_INDEX,
		N:      req.N,
		OrgId:  req.OrgId,
		Object: req.Object,
	})
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.GetTopTags, err, "not found")
	}
	rsp.Data = &behaviour_proto.GetTopTagsResponse_Data{rsp_tags.Tags}
	return nil
}

func (p *BehaviourService) AutocompleteTags(ctx context.Context, req *behaviour_proto.AutocompleteTagsRequest, rsp *behaviour_proto.AutocompleteTagsResponse) error {
	log.Info("Received Behaviour.AutocompleteTags request")

	tags, err := db.AutocompleteTags(ctx, req.OrgId, req.Object, req.Name)
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.AutocompleteTags, err, "not found")
	}
	rsp.Data = &behaviour_proto.AutocompleteTagsResponse_Data{tags}
	return nil
}

func (p *BehaviourService) WarmupCacheBehaviour(ctx context.Context, req *behaviour_proto.WarmupCacheBehaviourRequest, rsp *behaviour_proto.WarmupCacheBehaviourResponse) error {
	log.Info("Received Behaviour.WarmupCacheBehaviour request")

	var offset int64
	var limit int64
	offset = 0
	limit = 100
	switch req.Object {
	case common.GOAL:
		for {
			items, err := db.AllGoals(ctx, "", "", offset, limit, "", "")
			if err != nil || len(items) == 0 {
				break
			}
			for _, item := range items {
				if len(item.Tags) > 0 {
					if _, err := p.KvClient.TagsCloud(ctx, &kv_proto.TagsCloudRequest{
						Index:  common.CLOUD_TAGS_INDEX,
						OrgId:  item.OrgId,
						Object: common.GOAL,
						Tags:   item.Tags,
					}); err != nil {
						return common.InternalServerError(common.BehaviourSrv, p.WarmupCacheBehaviour, err, "tag error")
					}
				}
			}
			offset += limit
		}
	case common.CHALLENGE:
		for {
			items, err := db.AllChallenges(ctx, "", "", offset, limit, "", "")
			if err != nil || len(items) == 0 {
				break
			}
			for _, item := range items {
				if len(item.Tags) > 0 {
					if _, err := p.KvClient.TagsCloud(ctx, &kv_proto.TagsCloudRequest{
						Index:  common.CLOUD_TAGS_INDEX,
						OrgId:  item.OrgId,
						Object: common.CHALLENGE,
						Tags:   item.Tags,
					}); err != nil {
						return common.InternalServerError(common.BehaviourSrv, p.WarmupCacheBehaviour, err, "tag error")
					}
				}
			}
			offset += limit
		}
	case common.HABIT:
		for {
			items, err := db.AllHabits(ctx, "", "", offset, limit, "", "")
			if err != nil || len(items) == 0 {
				break
			}
			for _, item := range items {
				if len(item.Tags) > 0 {
					if _, err := p.KvClient.TagsCloud(ctx, &kv_proto.TagsCloudRequest{
						Index:  common.CLOUD_TAGS_INDEX,
						OrgId:  item.OrgId,
						Object: common.HABIT,
						Tags:   item.Tags,
					}); err != nil {
						return common.InternalServerError(common.BehaviourSrv, p.WarmupCacheBehaviour, err, "tag error")
					}
				}
			}
			offset += limit
		}
	}

	return nil
}

func (p *BehaviourService) AllGoalResponse(ctx context.Context, req *user_proto.AllGoalResponseRequest, rsp *user_proto.AllGoalResponseResponse) error {
	log.Info("Received Behaviour.AllGoalResponse request")
	response, err := db.AllGoalResponse(ctx, req.CreatedBy, req.Query, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.AllGoalResponse, err, "not found")
	}
	rsp.Data = &user_proto.AllGoalResponseResponse_Data{response}
	return nil
}

func (p *BehaviourService) AllChallengeResponse(ctx context.Context, req *user_proto.AllChallengeResponseRequest, rsp *user_proto.AllChallengeResponseResponse) error {
	log.Info("Received Behaviour.AllChallengeResponse request")
	response, err := db.AllChallengeResponse(ctx, req.CreatedBy, req.Query, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.AllChallengeResponse, err, "not found")
	}
	rsp.Data = &user_proto.AllChallengeResponseResponse_Data{response}
	return nil
}

func (p *BehaviourService) AllHabitResponse(ctx context.Context, req *user_proto.AllHabitResponseRequest, rsp *user_proto.AllHabitResponseResponse) error {
	log.Info("Received Behaviour.AllHabitResponse request")
	response, err := db.AllHabitResponse(ctx, req.CreatedBy, req.Query, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if err != nil {
		return common.NotFound(common.BehaviourSrv, p.AllHabitResponse, err, "not found")
	}
	rsp.Data = &user_proto.AllHabitResponseResponse_Data{response}
	return nil
}

func (p *BehaviourService) UpdateGoal(ctx context.Context, req *behaviour_proto.UpdateGoalRequest, rsp *behaviour_proto.UpdateGoalResponse) error {
	log.Info("Received Behaviour.UpdateGoal request")

	return nil
}

func (p *BehaviourService) UpdateChallenge(ctx context.Context, req *behaviour_proto.UpdateChallengeRequest, rsp *behaviour_proto.UpdateChallengeResponse) error {
	log.Info("Received Behaviour.UpdateGoal request")

	return nil
}

func (p *BehaviourService) UpdateHabit(ctx context.Context, req *behaviour_proto.UpdateHabitRequest, rsp *behaviour_proto.UpdateHabitResponse) error {
	log.Info("Received Behaviour.UpdateGoal request")

	return nil
}

func (p *BehaviourService) UploadGoals(ctx context.Context, req *behaviour_proto.UploadGoalsRequest, rsp *behaviour_proto.UploadGoalsResponse) error {
	log.Info("Received Behaviour.UploadGoals request")

	// f, err := os.OpenFile("./upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer f.Close()
	// io.Copy(f, file)

	// response, err := db.UploadGoals(ctx, req.CreatedBy, req.Query, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	// if err != nil {
	// 	return common.NotFound(common.BehaviourSrv, p.UploadGoals, err, "not found")
	// }
	// rsp.Data = &user_proto.UploadGoalsResponse_Data{response}
	return nil
}

//FIXME: this is repeated in behaviour, plan, content and survey - combine to single function somewhere? Not sure where
func (p *BehaviourService) sendShareNotification(userids []string, message string, alert *pubsub_proto.Alert, data map[string]string) error {
	log.Info("Sending notification message for shared resource: ", message, userids)
	msg := &pubsub_proto.PublishBulkNotification{
		Notification: &pubsub_proto.BulkNotification{
			UserIds: userids,
			Message: message,
			Alert:   alert,
			Data:    data,
		},
	}
	if body, err := json.Marshal(msg); err == nil {
		if err := p.Broker.Publish(common.SEND_NOTIFICATION, &broker.Message{Body: body}); err != nil {
			return err
		}
	}
	return nil
}
