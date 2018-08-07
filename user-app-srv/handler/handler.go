package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	account_proto "server/account-srv/proto/account"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	content_proto "server/content-srv/proto/content"
	kv_proto "server/kv-srv/proto/kv"
	plan_proto "server/plan-srv/proto/plan"
	common_proto "server/static-srv/proto/common"
	static_proto "server/static-srv/proto/static"
	survey_proto "server/survey-srv/proto/survey"
	track_proto "server/track-srv/proto/track"
	"server/user-app-srv/db"
	userapp_proto "server/user-app-srv/proto/userapp"
	user_proto "server/user-srv/proto/user"
	"strconv"
	"time"

	duration "github.com/ChannelMeter/iso8601duration"
	"github.com/micro/go-micro/broker"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/golang/protobuf/ptypes"
)

type UserAppService struct {
	Broker          broker.Broker
	KvClient        kv_proto.KvServiceClient
	ContentClient   content_proto.ContentServiceClient
	BehaviourClient behaviour_proto.BehaviourServiceClient
	UserClient      user_proto.UserServiceClient
	TrackClient     track_proto.TrackServiceClient
	PlanClient      plan_proto.PlanServiceClient
	SurveyClient    survey_proto.SurveyServiceClient
	AccountClient   account_proto.AccountServiceClient
	StaticClient    static_proto.StaticServiceClient
}

func (p *UserAppService) CreateBookmark(ctx context.Context, req *userapp_proto.CreateBookmarkRequest, rsp *userapp_proto.CreateBookmarkResponse) error {
	log.Info("Received UserApp.CreateBookmark request")

	bookmark, err := db.CreateBookmark(ctx, req.OrgId, req.UserId, req.ContentId)
	if err != nil {
		return common.InternalServerError(common.UserappSrv, p.CreateBookmark, err, "create error")
	}
	// publish
	msg := &userapp_proto.ContentBookmarkStatusMessage{
		UserId:     req.UserId,
		ContentId:  req.ContentId,
		Bookmarked: true,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return common.InternalServerError(common.UserappSrv, p.CreateBookmark, err, "parsing error")
	}
	if err := p.Broker.Publish(common.CONTENT_BOOKMARKED, &broker.Message{Body: body}); err != nil {
		return common.NotFound(common.UserappSrv, p.CreateBookmark, err, fmt.Sprintf("%v Pub/Sub is failed", common.CONTENT_BOOKMARKED))
	}

	rsp.Data = &userapp_proto.CreateBookmarkResponse_Data{bookmark.Id}
	return nil
}

func (p *UserAppService) ReadBookmarkContents(ctx context.Context, req *userapp_proto.ReadBookmarkContentRequest, rsp *userapp_proto.ReadBookmarkContentResponse) error {
	log.Info("Received UserApp.ReadBookmarkContents request")
	bookmarkContents, err := db.ReadBookmarkContents(ctx, req.UserId)
	if len(bookmarkContents) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.ReadBookmarkContents, err, "ReadBookmarkContents query is failed")
	}
	rsp.Data = &userapp_proto.ReadBookmarkContentResponse_Data{bookmarkContents}
	return nil
}

func (p *UserAppService) ReadBookmarkContentCategorys(ctx context.Context, req *userapp_proto.ReadBookmarkContentCategorysRequest, rsp *userapp_proto.ReadBookmarkContentCategorysResponse) error {
	log.Info("Received UserApp.ReadBookmarkContentCategorys request")
	categorys, err := db.ReadBookmarkContentCategorys(ctx, req.UserId)
	if len(categorys) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.ReadBookmarkContentCategorys, err, "ReadBookmarkContentCategorys query is failed")
	}
	rsp.Data = &userapp_proto.ReadBookmarkContentCategorysResponse_Data{categorys}
	return nil
}

func (p *UserAppService) ReadBookmarkByCategory(ctx context.Context, req *userapp_proto.ReadBookmarkByCategoryRequest, rsp *userapp_proto.ReadBookmarkByCategoryResponse) error {
	log.Info("Received UserApp.ReadBookmarkByCategory request")
	bookmarkContents, err := db.ReadBookmarkByCategory(ctx, req.UserId, req.CategoryId)
	if len(bookmarkContents) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.ReadBookmarkByCategory, err, "ReadBookmarkByCategory query is failed")
	}
	rsp.Data = &userapp_proto.ReadBookmarkByCategoryResponse_Data{bookmarkContents}
	return nil
}

func (p *UserAppService) DeleteBookmark(ctx context.Context, req *userapp_proto.DeleteBookmarkRequest, rsp *userapp_proto.DeleteBookmarkResponse) error {
	log.Info("Received UserApp.DeleteBookmark request")

	// read bookmark to get content_id
	bookmark, err := db.ReadBookmark(ctx, req.BookmarkId)
	if err != nil || bookmark == nil {
		return common.NotFound(common.UserappSrv, p.DeleteBookmark, err, "ReadBookmark query is failed")
	}

	if err := db.DeleteBookmark(ctx, req.OrgId, req.UserId, req.BookmarkId); err != nil {
		return common.InternalServerError(common.UserappSrv, p.DeleteBookmark, err, "DeleteBookmark query is failed")
	}
	// publish
	msg := &userapp_proto.ContentBookmarkStatusMessage{
		UserId:     bookmark.User.Id,
		ContentId:  bookmark.Content.Id,
		Bookmarked: false,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return common.InternalServerError(common.UserappSrv, p.DeleteBookmark, err, "parsing error")
	}
	if err := p.Broker.Publish(common.CONTENT_UNBOOKMARKED, &broker.Message{Body: body}); err != nil {
		return common.InternalServerError(common.UserappSrv, p.DeleteBookmark, err, common.CONTENT_UNBOOKMARKED+" Pub/Sub is failed")
	}
	return nil
}

func (p *UserAppService) SearchBookmarks(ctx context.Context, req *userapp_proto.SearchBookmarkRequest, rsp *user_proto.GetShareableContentResponse) error {
	log.Info("Received UserApp.SearchBookmarks request")

	shares, err := db.SearchBookmarks(ctx, req.Title, req.Summary, req.Description, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit)
	if len(shares) == 0 || err != nil {
		common.ErrorLog(common.UserappSrv, common.GetFunctionName(p.GetSharedContent), err, "GetSharedContent query is failed")
		return err
	}
	rsp.Data = &user_proto.GetShareableContentResponse_Data{shares}
	return nil
}

func (p *UserAppService) GetSharedContent(ctx context.Context, req *userapp_proto.GetSharedContentRequest, rsp *userapp_proto.GetSharedContentResponse) error {
	log.Info("Received UserApp.GetSharedContent request")

	shares, err := db.GetSharedContent(ctx, req.UserId)
	if len(shares) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetSharedContent, err, "GetSharedContent query is failed")
	}
	rsp.Data = &userapp_proto.GetSharedContentResponse_Data{shares}
	return nil
}

func (p *UserAppService) GetSharedPlansForUser(ctx context.Context, req *userapp_proto.GetSharedPlanRequest, rsp *userapp_proto.GetSharedPlanResponse) error {
	log.Info("Received UserApp.GetSharedPlansForUser request")

	shares, err := db.GetSharedPlansForUser(ctx, req.UserId)
	if len(shares) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetSharedPlansForUser, err, "GetSharedPlansForUser query is failed")
	}
	rsp.Data = &userapp_proto.GetSharedPlanResponse_Data{shares}
	return nil
}

func (p *UserAppService) GetSharedSurveysForUser(ctx context.Context, req *userapp_proto.GetSharedSurveyRequest, rsp *userapp_proto.GetSharedSurveyResponse) error {
	log.Info("Received UserApp.GetSharedSurveysForUser request")

	shares, err := db.GetSharedSurveysForUser(ctx, req.UserId)
	if len(shares) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetSharedSurveysForUser, err, "GetSharedSurvey query is failed")
	}
	rsp.Data = &userapp_proto.GetSharedSurveyResponse_Data{shares}
	return nil
}

func (p *UserAppService) GetSharedGoalsForUser(ctx context.Context, req *userapp_proto.GetSharedGoalRequest, rsp *userapp_proto.GetSharedGoalResponse) error {
	log.Info("Received UserApp.GetSharedGoalsForUser request")

	shares, err := db.GetSharedGoalsForUser(ctx, req.UserId)
	if len(shares) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetSharedGoalsForUser, err, "GetSharedGoal query is failed")
	}
	rsp.Data = &userapp_proto.GetSharedGoalResponse_Data{shares}
	return nil
}

func (p *UserAppService) GetSharedChallengesForUser(ctx context.Context, req *userapp_proto.GetSharedChallengeRequest, rsp *userapp_proto.GetSharedChallengeResponse) error {
	log.Info("Received Challenge.GetSharedChallengesForUser request")

	shares, err := db.GetSharedChallengesForUser(ctx, req.UserId)
	if len(shares) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetSharedChallengesForUser, err, "GetSharedChallenge query is failed")
	}
	rsp.Data = &userapp_proto.GetSharedChallengeResponse_Data{shares}
	return nil
}

func (p *UserAppService) GetSharedHabitsForUser(ctx context.Context, req *userapp_proto.GetSharedHabitRequest, rsp *userapp_proto.GetSharedHabitResponse) error {
	log.Info("Received UserApp.GetSharedHabitsForUser request")

	shares, err := db.GetSharedHabitsForUser(ctx, req.UserId)
	if len(shares) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetSharedChallengesForUser, err, "GetSharedHabit query is failed")
	}
	rsp.Data = &userapp_proto.GetSharedHabitResponse_Data{shares}
	return nil
}

func (p *UserAppService) SignupToGoal(ctx context.Context, req *userapp_proto.SignupToGoalRequest, rsp *userapp_proto.SignupToGoalResponse) error {
	log.Info("Received UserApp.SignupToGoal request")

	// get count from REDIS
	var count int64
	key := common.GetTrackKey(req.UserId, "active_goal", "count")
	rsp_kv, err := p.KvClient.GetTrackCount(ctx, &kv_proto.GetTrackCountRequest{common.USERAPP_INDEX, key})
	if err != nil {
		// getting count from database
		count = 0
		goals, _ := db.GetJoinedGoals(ctx, req.UserId, true)
		if goals != nil {
			count = int64(len(goals))
		}
		if _, err := p.KvClient.SetTrackCount(ctx, &kv_proto.SetTrackCountRequest{common.USERAPP_INDEX, key, count}); err != nil {
			return common.InternalServerError(common.UserappSrv, p.SignupToGoal, err, "SetTrackCount is failed")
		}
	} else {
		count = rsp_kv.Count
	}

	// check validation with limit
	if count >= common.CURRENT_GOALS_JOINED {
		return common.InternalServerError(common.UserappSrv, p.SignupToGoal, nil, "Limit for goals joined reached")
	}

	// 19/06/2018
	// Current logic - goals needs to be shared with user, we fetch the shared goal and proceed forward
	// New Logic  - If the goal has not been shared with user - then the user is trying to join the goal through discover

	// get the shared goal by using behaviourclient
	rsp_goal, err := p.BehaviourClient.GetSharedGoal(ctx, &behaviour_proto.GetSharedGoalRequest{UserId: req.UserId, GoalId: req.GoalId})
	signingUpGoal := &behaviour_proto.Goal{}

	if err != nil {
		return common.NotFound(common.UserappSrv, p.SignupToGoal, err, "GetSharedGoal has failed")
		//this means that goal is not shared with the user and the user is trying to join through discovering the goal
		goal, err := p.BehaviourClient.ReadGoal(ctx, &behaviour_proto.ReadGoalRequest{GoalId: req.GoalId, OrgId: req.OrgId})
		if err != nil {
			return common.NotFound(common.UserappSrv, p.SignupToGoal, err, "ReadGoal has failed")
		}
		signingUpGoal = goal.Data.Goal
	} else if rsp_goal != nil {
		//this means goal is shared by a team member
		signingUpGoal = rsp_goal.Data.Goal.Goal
	}

	//TODO: throw an error NOT_FOUND here and exit
	if signingUpGoal == nil {
		return errors.New(err.Error())
	}

	// calculate duration from ISO8601
	dur, err := duration.FromString(signingUpGoal.Duration)
	if err != nil {
		return common.InternalServerError(common.UserappSrv, p.SignupToGoal, err, "Duration parsing error")
	}
	join_goal := &userapp_proto.JoinGoal{
		Id:     uuid.NewUUID().String(),
		Status: userapp_proto.ActionStatus_STARTED,
		Start:  time.Now().Unix(),
		End:    int64(time.Now().Add(dur.ToDuration()).Unix()),
		Target: signingUpGoal.Target,
	}

	is_new_signup, err := db.SignupToGoal(ctx, req.UserId, req.GoalId, join_goal)
	if err != nil {
		return common.InternalServerError(common.UserappSrv, p.SignupToGoal, err, "SignupToGoal query is failed")
	}
	if is_new_signup {
		// increment count only if this is a new goal that the user is signing upto (insert and not update)
		log.Info("User is joining a new goal: ", req.GoalId)
		if _, err := p.KvClient.IncTrackCount(ctx, &kv_proto.IncTrackCountRequest{common.USERAPP_INDEX, key}); err != nil {
			return common.InternalServerError(common.UserappSrv, p.SignupToGoal, err, "IncTrackCount is failed")
		}

		//only doing this for shared goals delete pending actions and update share goal status for shared goals only
		if rsp_goal != nil {
			log.Info("removing pending action for goal: ", req.GoalId)
			if err := p.RemovePendingSharedAction(ctx, req.GoalId, req.UserId); err != nil {
				return common.InternalServerError(common.UserappSrv, p.SignupToGoal, err, "RemovePendingSharedAction failed")
			}
			// update shareX status
			log.Info("updating shared goal status for goal: ", req.GoalId)
			if err := db.UpdateShareChallengeStatus(ctx, req.GoalId, static_proto.ShareStatus_VIEWED); err != nil {
				return common.InternalServerError(common.UserappSrv, p.SignupToGoal, err, "UpdateShareChallengeStatus query is failed")
			}
		}
	}

	rsp.Data = &userapp_proto.SignupToGoalResponse_Data{
		JoinGoal: join_goal,
		Count:    count,
	}
	return nil
}

func (p *UserAppService) GetCurrentJoinedGoals(ctx context.Context, req *userapp_proto.ListGoalRequest, rsp *userapp_proto.ListGoalResponse) error {
	log.Info("Received UserApp.GetCurrentJoinedGoals request")

	goals, err := db.GetJoinedGoals(ctx, req.UserId, true)
	if len(goals) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetCurrentJoinedGoals, err, "GetCurrentJoinedGoals query is failed")
	}
	rsp.Data = &userapp_proto.ListGoalResponse_Data{goals}
	return nil
}

func (p *UserAppService) GetAllJoinedGoals(ctx context.Context, req *userapp_proto.ListGoalRequest, rsp *userapp_proto.ListGoalResponse) error {
	log.Info("Received UserApp.GetAllJoinedGoals request")

	goals, err := db.GetJoinedGoals(ctx, req.UserId, false)
	if len(goals) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetAllJoinedGoals, err, "GetAllJoinedGoals query is failed")
	}
	rsp.Data = &userapp_proto.ListGoalResponse_Data{goals}
	return nil
}

func (p *UserAppService) SignupToChallenge(ctx context.Context, req *userapp_proto.SignupToChallengeRequest, rsp *userapp_proto.SignupToChallengeResponse) error {
	log.Info("Received UserApp.SignupToChallenge request for challenge: ", req.ChallengeId)

	// get count from REDIS
	var count int64
	key := common.GetTrackKey(req.UserId, "active_challenge", "count")
	rsp_kv, err := p.KvClient.GetTrackCount(ctx, &kv_proto.GetTrackCountRequest{common.USERAPP_INDEX, key})

	//not found in kv for this user, so fetch the joined challenge count from db and store it in kv for this user
	if err != nil {
		// getting count from database
		count = 0
		challenges, _ := db.GetJoinedChallenges(ctx, req.UserId, true)
		if challenges != nil {
			count = int64(len(challenges))
		}
		if _, err := p.KvClient.SetTrackCount(ctx, &kv_proto.SetTrackCountRequest{common.USERAPP_INDEX, key, count}); err != nil {
			return common.InternalServerError(common.UserappSrv, p.SignupToChallenge, err, "track cound set fail")
		}
	} else {
		count = rsp_kv.Count
	}
	log.Debug("signup challenge count:", count)
	// check validation with limit (the challenge join chount should be checked with >=)
	if count >= common.CURRENT_CHALLENGES_JOINED {
		return common.InternalServerError(common.UserappSrv, p.SignupToGoal, nil, "Limit current challenge joined")
	}

	// 19/06/2018
	// Current logic - challenge needs to be shared with user, we fetch the shared challenge and proceed forward
	// New Logic  - If the challenge has not been shared with user - then the user is trying to join the challenge through discover

	// get the shared challenge by using behaviourclient
	rsp_challenge, err := p.BehaviourClient.GetSharedChallenge(ctx, &behaviour_proto.GetSharedChallengeRequest{UserId: req.UserId, ChallengeId: req.ChallengeId})
	signingUpChallenge := &behaviour_proto.Challenge{}

	if err != nil {
		return common.NotFound(common.UserappSrv, p.SignupToChallenge, err, "GetSharedChallenge query is failed")
		//this means that challenge is not shared with the user and the user is trying to join through discovering the challenge
		challenge, err := p.BehaviourClient.ReadChallenge(ctx, &behaviour_proto.ReadChallengeRequest{ChallengeId: req.ChallengeId, OrgId: req.OrgId})
		if err != nil {
			return common.NotFound(common.UserappSrv, p.SignupToChallenge, err, "ReadChallenge query is failed")
		}
		signingUpChallenge = challenge.Data.Challenge
	} else if rsp_challenge != nil {
		signingUpChallenge = rsp_challenge.Data.Challenge.Challenge
	}

	//TODO: throw an error NOT_FOUND here and exit
	if signingUpChallenge == nil {
		return errors.New(err.Error())
	}

	// calculate duration from ISO8601
	dur, err := duration.FromString(signingUpChallenge.Duration)
	if err != nil {
		return errors.New("duration parsing error:" + err.Error())
	}
	join_challenge := &userapp_proto.JoinChallenge{
		Id:     uuid.NewUUID().String(),
		Status: userapp_proto.ActionStatus_STARTED,
		Start:  time.Now().Unix(),
		End:    int64(time.Now().Add(dur.ToDuration()).Unix()),
		Target: signingUpChallenge.Target,
	}

	is_new_signup, err := db.SignupToChallenge(ctx, req.UserId, req.ChallengeId, join_challenge)
	if err != nil {
		return common.InternalServerError(common.UserappSrv, p.SignupToChallenge, err, "SignupToChallenge query is failed")
	}
	// increment count only if this is a new challenge that the user is signing upto (insert and not update)
	if is_new_signup {
		log.Info("User is joining a new goal: ", req.ChallengeId)
		if _, err := p.KvClient.IncTrackCount(ctx, &kv_proto.IncTrackCountRequest{common.USERAPP_INDEX, key}); err != nil {
			return common.InternalServerError(common.UserappSrv, p.SignupToChallenge, err, "track count doesn't increase")
		}

		//only doing this for shared challenges delete pending actions and update share challenge status for shared challenges only
		if rsp_challenge != nil {
			log.Info("removing pending action for challenge: ", req.ChallengeId)
			if err := p.RemovePendingSharedAction(ctx, req.ChallengeId, req.UserId); err != nil {
				return common.InternalServerError(common.UserappSrv, p.SignupToChallenge, err, "UpdateShareChallengeStatus failed")
			}
			// update shareX status
			log.Info("updating shared challenge status for challenge: ", req.ChallengeId)
			if err := db.UpdateShareChallengeStatus(ctx, req.ChallengeId, static_proto.ShareStatus_VIEWED); err != nil {
				return common.InternalServerError(common.UserappSrv, p.SignupToChallenge, err, "UpdateShareChallengeStatus query is failed")
			}
		}
	}

	rsp.Data = &userapp_proto.SignupToChallengeResponse_Data{
		JoinChallenge: join_challenge,
		Count:         count,
	}
	return nil
}

func (p *UserAppService) GetCurrentJoinedChallenges(ctx context.Context, req *userapp_proto.ListChallengeRequest, rsp *userapp_proto.ListChallengeResponse) error {
	log.Info("Received UserApp.GetCurrentJoinedChallenges request")

	challenges, err := db.GetJoinedChallenges(ctx, req.UserId, true)
	if len(challenges) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetCurrentJoinedChallenges, err, "not found")
	}
	rsp.Data = &userapp_proto.ListChallengeResponse_Data{challenges}
	return nil
}

func (p *UserAppService) GetAllJoinedChallenges(ctx context.Context, req *userapp_proto.ListChallengeRequest, rsp *userapp_proto.ListChallengeResponse) error {
	log.Info("Received UserApp.GetJoinedChallenges request")

	challenges, err := db.GetJoinedChallenges(ctx, req.UserId, false)
	if len(challenges) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetAllJoinedChallenges, err, "not found")
	}
	rsp.Data = &userapp_proto.ListChallengeResponse_Data{challenges}
	return nil
}

func (p *UserAppService) SignupToHabit(ctx context.Context, req *userapp_proto.SignupToHabitRequest, rsp *userapp_proto.SignupToHabitResponse) error {
	log.Info("Received UserApp.SignupToHabit request")

	// get count from REDIS
	var count int64
	key := common.GetTrackKey(req.UserId, "active_habit", "count")
	rsp_kv, err := p.KvClient.GetTrackCount(ctx, &kv_proto.GetTrackCountRequest{common.USERAPP_INDEX, key})
	if err != nil {
		// getting count from database
		count = 0
		habits, _ := db.GetJoinedHabits(ctx, req.UserId, true)
		if habits != nil {
			count = int64(len(habits))
		}
		if _, err := p.KvClient.SetTrackCount(ctx, &kv_proto.SetTrackCountRequest{common.USERAPP_INDEX, key, count}); err != nil {
			return common.InternalServerError(common.UserappSrv, p.SignupToHabit, err, "track count fail")
		}
	} else {
		count = rsp_kv.Count
	}
	// check validation with limit
	if count >= common.CURRENT_HABITS_JOINED {
		return common.InternalServerError(common.UserappSrv, p.SignupToGoal, nil, "Limit current habit joined")
	}

	// 19/06/2018
	// Current logic - habits needs to be shared with user, we fetch the shared habit and proceed forward
	// New Logic  - If the habit has not been shared with user - then the user is trying to join the habit through discover

	// get the shared habit by using behaviourclient
	rsp_habit, err := p.BehaviourClient.GetSharedHabit(ctx, &behaviour_proto.GetSharedHabitRequest{UserId: req.UserId, HabitId: req.HabitId})
	signingUpHabit := &behaviour_proto.Habit{}

	if err != nil {
		log.Info("User is trying join habit through discovery: ", req.HabitId)
		return common.InternalServerError(common.UserappSrv, p.SignupToHabit, err, "GetSharedHabit query is failed")
		//this means that challenge is not shared with the user and the user is trying to join through discovering the challenge
		habit, err := p.BehaviourClient.ReadHabit(ctx, &behaviour_proto.ReadHabitRequest{HabitId: req.HabitId, OrgId: req.OrgId})
		if err != nil {
			return common.InternalServerError(common.UserappSrv, p.SignupToHabit, err, "ReadHabit query is failed")
		}
		signingUpHabit = habit.Data.Habit
	} else if rsp_habit != nil {
		signingUpHabit = rsp_habit.Data.Habit.Habit
	}

	//TODO: throw an error NOT_FOUND here and exit
	if signingUpHabit == nil {
		return errors.New(err.Error())
	}

	// calculate duration from ISO8601
	dur, err := duration.FromString(signingUpHabit.Duration)
	if err != nil {
		return errors.New("duration parsing error:" + err.Error())
	}
	join_habit := &userapp_proto.JoinHabit{
		Id:     uuid.NewUUID().String(),
		Status: userapp_proto.ActionStatus_STARTED,
		Start:  time.Now().Unix(),
		End:    int64(time.Now().Add(dur.ToDuration()).Unix()),
		Target: signingUpHabit.Target,
	}
	is_new_signup, err := db.SignupToHabit(ctx, req.UserId, req.HabitId, join_habit)
	if err != nil {
		return common.InternalServerError(common.UserappSrv, p.SignupToHabit, err, "SignupToHabit query is failed")
	}
	// increment count only if this is a new habit that the user is signing upto (insert and not update)
	if is_new_signup {
		log.Info("User is joining a new habit: ", req.HabitId)
		if _, err := p.KvClient.IncTrackCount(ctx, &kv_proto.IncTrackCountRequest{common.USERAPP_INDEX, key}); err != nil {
			return common.InternalServerError(common.UserappSrv, p.SignupToHabit, nil, "IncTrackCount is failed")
		}

		//only doing this for shared habits
		if rsp_habit != nil {
			log.Info("removing pending action for habit: ", req.HabitId)
			if err := p.RemovePendingSharedAction(ctx, req.HabitId, req.UserId); err != nil {
				return common.InternalServerError(common.UserappSrv, p.SignupToHabit, err, "RemovePendingSharedAction failed")
			}
			// update shareX status
			log.Info("updating shared habit status for habit: ", req.HabitId)
			if err := db.UpdateShareHabitStatus(ctx, req.HabitId, static_proto.ShareStatus_VIEWED); err != nil {
				return common.InternalServerError(common.UserappSrv, p.SignupToHabit, err, "UpdateShareHabitStatus query is failed")
			}
		}
	}

	rsp.Data = &userapp_proto.SignupToHabitResponse_Data{
		JoinHabit: join_habit,
		Count:     count,
	}
	return nil
}

func (p *UserAppService) GetCurrentJoinedHabits(ctx context.Context, req *userapp_proto.ListHabitRequest, rsp *userapp_proto.ListHabitResponse) error {
	log.Info("Received UserApp.GetCurrentJoinedHabits request")

	habits, err := db.GetJoinedHabits(ctx, req.UserId, true)
	if len(habits) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetCurrentJoinedHabits, err, "not found")
	}
	rsp.Data = &userapp_proto.ListHabitResponse_Data{habits}
	return nil
}

func (p *UserAppService) GetAllJoinedHabits(ctx context.Context, req *userapp_proto.ListHabitRequest, rsp *userapp_proto.ListHabitResponse) error {
	log.Info("Received UserApp.GetAllJoinedHabits request")

	habits, err := db.GetJoinedHabits(ctx, req.UserId, false)
	if len(habits) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetAllJoinedHabits, err, "not found")
	}
	rsp.Data = &userapp_proto.ListHabitResponse_Data{habits}
	return nil
}

func (p *UserAppService) ListMarkers(ctx context.Context, req *userapp_proto.ListMarkersRequest, rsp *userapp_proto.ListMarkersResponse) error {
	log.Info("Received UserApp.ListMarkers request")

	markers, err := db.ListMarkers(ctx, req.UserId)
	if len(markers) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.ListMarkers, err, "not found")
	}
	rsp.Data = &userapp_proto.ListMarkersResponse_Data{markers}
	return nil
}

func (p *UserAppService) GetPendingSharedActions(ctx context.Context, req *userapp_proto.GetPendingSharedActionsRequest, rsp *userapp_proto.GetPendingSharedActionsResponse) error {
	log.Info("Received UserApp.GetPendingSharedActions request")

	pendings, err := db.GetPendingSharedActions(ctx, req.UserId, req.OrgId, req.Offset, req.Limit, req.From, req.To, "created", "DESC")
	if len(pendings) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetPendingSharedActions, err, "not found")
	}

	for _, pending := range pendings {
		switch pending.Item.TypeUrl {
		case common.BASE + common.GOAL_TYPE:
			goal := &behaviour_proto.Goal{}
			if err := ptypes.UnmarshalAny(pending.Item, goal); err != nil {
				return err
			}

			rsp_goal, err := p.BehaviourClient.ReadGoal(ctx, &behaviour_proto.ReadGoalRequest{GoalId: goal.Id})
			if err == nil && rsp_goal != nil {
				pending.Title = rsp_goal.Data.Goal.Title
				pending.Image = rsp_goal.Data.Goal.Image
				pending.Summary = rsp_goal.Data.Goal.Summary
				pending.Duration = rsp_goal.Data.Goal.Duration
			}
		case common.BASE + common.CHALLENGE_TYPE:
			challenge := &behaviour_proto.Challenge{}
			if err := ptypes.UnmarshalAny(pending.Item, challenge); err != nil {
				return err
			}

			rsp_challenge, err := p.BehaviourClient.ReadChallenge(ctx, &behaviour_proto.ReadChallengeRequest{ChallengeId: challenge.Id})
			if err == nil && rsp_challenge != nil {
				pending.Title = rsp_challenge.Data.Challenge.Title
				pending.Image = rsp_challenge.Data.Challenge.Image
				pending.Summary = rsp_challenge.Data.Challenge.Summary
				pending.Duration = rsp_challenge.Data.Challenge.Duration
			}
		case common.BASE + common.HABIT_TYPE:
			habit := &behaviour_proto.Habit{}
			if err := ptypes.UnmarshalAny(pending.Item, habit); err != nil {
				return err
			}

			rsp_habit, err := p.BehaviourClient.ReadHabit(ctx, &behaviour_proto.ReadHabitRequest{HabitId: habit.Id})
			if err == nil && rsp_habit != nil {
				pending.Title = rsp_habit.Data.Habit.Title
				pending.Image = rsp_habit.Data.Habit.Image
				pending.Summary = rsp_habit.Data.Habit.Summary
				pending.Duration = rsp_habit.Data.Habit.Duration
			}

		case common.BASE + common.SURVEY_TYPE:
			survey := &survey_proto.Survey{}
			if err := ptypes.UnmarshalAny(pending.Item, survey); err != nil {
				return err
			}

			rsp_survey, err := p.SurveyClient.Read(ctx, &survey_proto.ReadRequest{Id: survey.Id})
			if err == nil && rsp_survey != nil {
				pending.Title = rsp_survey.Data.Survey.Title
				pending.Summary = rsp_survey.Data.Survey.Summary
				pending.Count = int32(len(rsp_survey.Data.Survey.Questions))
				//image : survey doesn't have an image
				//duration: survey doesn't have a duration, but can send duration taken to respond to question?
			}
		case common.BASE + common.PLAN_TYPE:
			//plan doens't have image, so set count
		case common.BASE + common.CONTENT_TYPE:
			content := &content_proto.Content{}
			if err := ptypes.UnmarshalAny(pending.Item, content); err != nil {
				common.ErrorLog(common.UserappSrv, common.GetFunctionName(p.GetPendingSharedActions), err, "Survey marshalling is failed")
				return err
			}
			rsp_content := &userapp_proto.ReadContentResponse{}
			err := p.GetContentDetail(ctx, &content_proto.ReadContentRequest{Id: content.Id}, rsp_content)
			if err == nil && rsp_content != nil {
				pending.Title = rsp_content.Data.Detail.Title
				pending.Image = rsp_content.Data.Detail.Image
				pending.Summary = rsp_content.Data.Detail.Summary
			}
		}
	}

	rsp.Data = &userapp_proto.GetPendingSharedActionsResponse_Data{pendings}
	return nil
}

func (p *UserAppService) GetGoalProgress(ctx context.Context, req *userapp_proto.GetGoalProgressRequest, rsp *userapp_proto.GetGoalProgressResponse) error {
	log.Info("Received UserApp.GetGoalProgress request")

	progress, err := db.GetGoalProgress(ctx, req.UserId)
	if len(progress) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetGoalProgress, err, "not found")
	}
	rsp.Data = &userapp_proto.GetGoalProgressResponse_Data{progress}
	return nil
}

func (p *UserAppService) GetDefaultMarkerHistory(ctx context.Context, req *track_proto.GetDefaultMarkerHistoryRequest, rsp *track_proto.GetDefaultMarkerHistoryResponse) error {
	log.Info("Received UserApp.GetDefaultMarkerHistory request")

	rsp_track, err := p.TrackClient.GetDefaultMarkerHistory(ctx, req)
	if rsp_track == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.GetDefaultMarkerHistory, err, "not found")
	}
	rsp.Data = rsp_track.Data
	return nil
}

func (p *UserAppService) GetCurrentChallengesWithCount(ctx context.Context, req *userapp_proto.GetCurrentChallengesWithCountRequest, rsp *userapp_proto.GetCurrentChallengesWithCountResponse) error {
	log.Info("Received UserApp.GetCurrentChallengesWithCount request")

	challenges, err := db.GetCurrentChallengesWithCount(ctx, req.UserId)
	if len(challenges) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetCurrentChallengesWithCount, err, "not found")
	}
	rsp.Data = &userapp_proto.GetCurrentChallengesWithCountResponse_Data{challenges}
	return nil
}

func (p *UserAppService) GetCurrentHabitsWithCount(ctx context.Context, req *userapp_proto.GetCurrentHabitsWithCountRequest, rsp *userapp_proto.GetCurrentHabitsWithCountResponse) error {
	log.Info("Received UserApp.GetCurrentHabitsWithCount request")

	habits, err := db.GetCurrentHabitsWithCount(ctx, req.UserId)
	if len(habits) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetCurrentHabitsWithCount, err, "not found")
	}
	rsp.Data = &userapp_proto.GetCurrentHabitsWithCountResponse_Data{habits}
	return nil
}

func (p *UserAppService) GetContentCategorys(ctx context.Context, req *content_proto.GetContentCategorysRequest, rsp *content_proto.GetContentCategorysResponse) error {
	log.Info("Received UserApp.GetContentCategorys request")

	rsp_content, err := p.ContentClient.GetContentCategorys(ctx, req)
	if rsp_content == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.GetContentCategorys, err, "not found")
	}
	// match the result
	rsp.Data = rsp_content.Data
	return nil
}

func (p *UserAppService) GetContentByCategory(ctx context.Context, req *content_proto.GetContentByCategoryRequest, rsp *content_proto.GetContentByCategoryResponse) error {
	log.Info("Received UserApp.GetContentByCategory request")

	rsp_content, err := p.ContentClient.GetContentByCategory(ctx, req)
	if rsp_content == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.GetContentByCategory, err, "not found")
	}
	rsp.Data = rsp_content.Data
	return nil
}

func (p *UserAppService) GetFiltersForCategory(ctx context.Context, req *content_proto.GetFiltersForCategoryRequest, rsp *content_proto.GetFiltersForCategoryResponse) error {
	log.Info("Received UserApp.GetFiltersForCategory request")

	rsp_content, err := p.ContentClient.GetFiltersForCategory(ctx, req)
	if rsp_content == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.GetFiltersForCategory, err, "not found")
	}
	rsp.Data = rsp_content.Data
	return nil
}

func (p *UserAppService) FiltersAutocomplete(ctx context.Context, req *content_proto.FiltersAutocompleteRequest, rsp *content_proto.FiltersAutocompleteResponse) error {
	log.Info("Received UserApp.FiltersAutocomplete request")

	rsp_content, err := p.ContentClient.FiltersAutocomplete(ctx, req)
	if rsp_content == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.FiltersAutocomplete, err, "not found")
	}
	rsp.Data = rsp_content.Data
	return nil
}

func (p *UserAppService) FilterContentInParticularCategory(ctx context.Context, req *content_proto.FilterContentInParticularCategoryRequest, rsp *content_proto.FilterContentInParticularCategoryResponse) error {
	log.Info("Received UserApp.FilterContentInParticularCategory request")

	rsp_content, err := p.ContentClient.FilterContentInParticularCategory(ctx, req)
	if rsp_content == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.FilterContentInParticularCategory, err, "not found")
	}
	rsp.Data = rsp_content.Data
	return nil
}

func (p *UserAppService) GetUserPreference(ctx context.Context, req *user_proto.ReadUserPreferenceRequest, rsp *user_proto.ReadUserPreferenceResponse) error {
	log.Info("Received UserApp.GetUserPreference request")

	rsp_user, err := p.UserClient.ReadUserPreference(ctx, req)
	if rsp_user == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.GetUserPreference, err, "not found")
	}
	rsp.Data = rsp_user.Data
	return nil
}

func (p *UserAppService) SaveUserPreference(ctx context.Context, req *user_proto.SaveUserPreferenceRequest, rsp *user_proto.SaveUserPreferenceResponse) error {
	log.Info("Received UserApp.SaveUserPreference request")

	rsp_user, err := p.UserClient.SaveUserPreference(ctx, req)
	if rsp_user == nil || err != nil {
		return common.InternalServerError(common.UserappSrv, p.SaveUserPreference, err, "save error")
	}
	rsp.Data = rsp_user.Data
	return nil
}

func (p *UserAppService) SaveUserDetails(ctx context.Context, req *userapp_proto.SaveUserDetailsRequest, rsp *userapp_proto.SaveUserDetailsResponse) error {
	log.Info("Received UserApp.SaveUserDetails request")

	user := &user_proto.User{
		Id:        req.UserId,
		OrgId:     req.OrgId,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Gender:    req.Gender,
		AvatarUrl: req.AvatarUrl,
		Dob:       req.Dob,
	}
	//get valid user
	req_user := &user_proto.UpdateRequest{User: user}
	rsp_user, err := p.UserClient.Update(ctx, req_user)
	if rsp_user == nil || err != nil {
		return common.InternalServerError(common.UserappSrv, p.SaveUserDetails, err, "update error")
	}

	if rsp_user.Data.User == nil {
		return errors.New("user_not_found")
	}

	rsp.Data = &userapp_proto.SaveUserDetailsResponse_Data{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Gender:    req.Gender,
		AvatarUrl: req.AvatarUrl,
		Dob:       req.Dob,
	}
	return nil
}

func (p *UserAppService) GetContentRecommendationByUser(ctx context.Context, req *content_proto.GetContentRecommendationByUserRequest, rsp *content_proto.GetContentRecommendationByUserResponse) error {
	log.Info("Received UserApp.GetContentRecommendationByUser request")

	rsp_content, err := p.ContentClient.GetContentRecommendationByUser(ctx, req)
	if rsp_content == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.GetContentRecommendationByUser, err, "not found")
	}
	rsp.Data = rsp_content.Data
	return nil
}

func (p *UserAppService) GetContentRecommendationByCategory(ctx context.Context, req *content_proto.GetContentRecommendationByCategoryRequest, rsp *content_proto.GetContentRecommendationByCategoryResponse) error {
	log.Info("Received UserApp.GetContentRecommendationByCategory request")

	rsp_content, err := p.ContentClient.GetContentRecommendationByCategory(ctx, req)
	if rsp_content == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.GetContentRecommendationByCategory, err, "not found")
	}
	rsp.Data = rsp_content.Data
	return nil
}

func (p *UserAppService) SaveRateForContent(ctx context.Context, req *userapp_proto.SaveRateForContentRequest, rsp *userapp_proto.SaveRateForContentResponse) error {
	log.Info("Received UserApp.SaveRateForContent request")

	contentRating := &userapp_proto.ContentRating{
		OrgId:     req.OrgId,
		UserId:    req.UserId,
		ContentId: req.ContentId,
		Rating:    req.Rating,
	}
	err := db.SaveRateForContent(ctx, contentRating)
	if err != nil {
		return common.InternalServerError(common.UserappSrv, p.SaveRateForContent, err, "save fail")
	}

	// publish
	msg := &userapp_proto.ContentRatingMessage{
		ContentId: req.ContentId,
		UserId:    req.UserId,
		Rating:    req.Rating,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return common.InternalServerError(common.UserappSrv, p.SaveRateForContent, err, "parsing error")
	}
	if err := p.Broker.Publish(common.CONTENT_RATING_UPDATED, &broker.Message{Body: body}); err != nil {
		return common.InternalServerError(common.UserappSrv, p.SaveRateForContent, err, "subscribe error")
	}
	rsp.Data = &userapp_proto.SaveRateForContentResponse_Data{contentRating}
	return nil
}

func (p *UserAppService) DislikeForContent(ctx context.Context, req *userapp_proto.DislikeForContentRequest, rsp *userapp_proto.DislikeForContentResponse) error {
	log.Info("Received UserApp.DislikeForContent request")

	contentDislike := &userapp_proto.ContentDislike{
		OrgId:     req.OrgId,
		UserId:    req.UserId,
		ContentId: req.ContentId,
	}
	err := db.DislikeForContent(ctx, contentDislike)
	if err != nil {
		return common.InternalServerError(common.UserappSrv, p.DislikeForContent, err, "dislike error")
	}

	// publish
	msg := &userapp_proto.ContentDislikeMessage{
		ContentId: req.ContentId,
		UserId:    req.UserId,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if err := p.Broker.Publish(common.CONTENT_DISLIKED, &broker.Message{Body: body}); err != nil {
		return common.InternalServerError(common.UserappSrv, p.DislikeForContent, err, "subscribe error")
	}
	rsp.Data = &userapp_proto.DislikeForContentResponse_Data{contentDislike}
	return nil
}

func (p *UserAppService) DislikeForSimilarContent(ctx context.Context, req *userapp_proto.DislikeForSimilarContentRequest, rsp *userapp_proto.DislikeForSimilarContentResponse) error {
	log.Info("Received UserApp.DislikeForSimilarContent request")

	contentDislikeSimilar := &userapp_proto.ContentDislikeSimilar{
		OrgId:     req.OrgId,
		UserId:    req.UserId,
		ContentId: req.ContentId,
		Tags:      req.Tags,
	}
	err := db.DislikeForSimilarContent(ctx, contentDislikeSimilar)
	if err != nil {
		return common.InternalServerError(common.UserappSrv, p.DislikeForSimilarContent, err, "dislike error")
	}

	// publish
	msg := &userapp_proto.ContentDislikeSimilarMessage{
		ContentId: req.ContentId,
		UserId:    req.UserId,
		Tags:      req.Tags,
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if err := p.Broker.Publish(common.CONTENT_DISLIKE_SIMILAR, &broker.Message{Body: body}); err != nil {
		return common.InternalServerError(common.UserappSrv, p.DislikeForSimilarContent, err, "subscribe error")
	}
	rsp.Data = &userapp_proto.DislikeForSimilarContentResponse_Data{contentDislikeSimilar}
	return nil
}

func (p *UserAppService) SaveUserFeedback(ctx context.Context, req *userapp_proto.SaveUserFeedbackRequest, rsp *userapp_proto.SaveUserFeedbackResponse) error {
	log.Info("Received UserApp.SaveUserFeedback request")

	feedback := &user_proto.UserFeedback{
		UserId:   req.UserId,
		OrgId:    req.OrgId,
		Feedback: req.Feedback,
		Rating:   req.Rating,
	}
	if err := db.SaveUserFeedback(ctx, feedback); err != nil {
		return common.NotFound(common.UserappSrv, p.JoinUserPlan, err, "not found")
	}
	rsp.Data = &userapp_proto.SaveUserFeedbackResponse_Data{feedback}
	return nil
}

func (p *UserAppService) JoinUserPlan(ctx context.Context, req *userapp_proto.JoinUserPlanRequest, rsp *userapp_proto.JoinUserPlanResponse) error {
	log.Info("Received UserApp.JoinUserPlan request")

	// get plan from plan-srv
	rsp_plan, err := p.PlanClient.Read(ctx, &plan_proto.ReadRequest{Id: req.PlanId})
	if rsp_plan == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.JoinUserPlan, err, "not found")
	}
	// create user-plan
	plan := rsp_plan.Data.Plan
	dur, err := duration.FromString(plan.Duration)
	if err != nil {
		return errors.New("duration parsing error:" + err.Error())
	}
	userplan := &userapp_proto.UserPlan{
		Name:        plan.Name,
		OrgId:       plan.OrgId,
		Pic:         plan.Pic,
		Description: plan.Description,
		TargetUser:  req.UserId,
		Plan:        plan,
		Goals:       plan.Goals,
		Duration:    plan.Duration,
		Start:       time.Now().Unix(),
		End:         time.Now().Add(time.Duration(dur.ToDuration())).Unix(),
		Creator:     plan.Creator,
		Days:        plan.Days,
		ItemsCount:  plan.ItemsCount,
	}
	if err := db.CreateUserPlan(ctx, req.UserId, req.PlanId, userplan, true); err != nil {
		return common.InternalServerError(common.UserappSrv, p.JoinUserPlan, err, "create error")
	}

	if rsp_plan != nil {
		if err := p.RemovePendingSharedAction(ctx, req.PlanId, req.UserId); err != nil {
			return common.InternalServerError(common.UserappSrv, p.JoinUserPlan, err, "can't remove")
		}

		log.Debug("updating shared plan status for plan: ", req.PlanId)
		// update share_plan_user status
		if err := db.UpdateSharePlanStatus(ctx, req.PlanId, static_proto.ShareStatus_VIEWED); err != nil {
			return common.InternalServerError(common.UserappSrv, p.JoinUserPlan, err, "share plan user update error")
		}
	}

	rsp.Data = &userapp_proto.JoinUserPlanResponse_Data{userplan}
	return nil
}

func (p *UserAppService) CreateUserPlan(ctx context.Context, req *userapp_proto.CreateUserPlanRequest, rsp *userapp_proto.CreateUserPlanResponse) error {
	log.Info("Received UserApp.CreateUserPlan request")

	// pick random itmes from content collection
	rsp_content, err := p.ContentClient.GetRandomItems(ctx, &content_proto.GetRandomItemsRequest{req.Days * req.ItemsPerDay})
	// make dayitems
	days := map[string]*common_proto.DayItems{}
	for i := 0; i < int(req.Days); i++ {
		dayitems := &common_proto.DayItems{[]*common_proto.DayItem{}}
		for j := 0; j < int(req.ItemsPerDay); j++ {
			index := i*int(req.ItemsPerDay) + j
			if index >= len(rsp_content.Data.Contents) {
				break
			}
			content := rsp_content.Data.Contents[index]
			dayitems.Items = append(dayitems.Items, &common_proto.DayItem{
				ContentId:        content.Id,
				ContentTitle:     content.Title,
				ContentPicUrl:    content.Image,
				CategoryId:       content.Category.Id,
				CategoryIconSlug: content.Category.IconSlug,
				CategoryName:     content.Category.Name,
			})
		}
		days[strconv.Itoa(i+1)] = dayitems
	}
	// get goal
	rsp_goal, err := p.BehaviourClient.ReadGoal(ctx, &behaviour_proto.ReadGoalRequest{GoalId: req.GoalId})
	if rsp_goal == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.CreateUserPlan, err, "not found")
	}
	// create user-plan
	d := duration.Duration{Days: int(req.Days)}
	userplan := &userapp_proto.UserPlan{
		Goals:      []*behaviour_proto.Goal{rsp_goal.Data.Goal},
		Duration:   d.String(),
		Start:      time.Now().Unix(),
		End:        time.Now().Add(time.Duration(time.Hour * 24 * time.Duration(req.Days))).Unix(),
		Creator:    &user_proto.User{Id: req.UserId},
		Days:       days,
		ItemsCount: req.Days * req.ItemsPerDay,
	}
	if err := db.CreateUserPlan(ctx, req.UserId, req.GoalId, userplan, false); err != nil {
		return common.InternalServerError(common.UserappSrv, p.CreateUserPlan, err, "userplan create error")
	}
	rsp.Data = &userapp_proto.CreateUserPlanResponse_Data{userplan}
	return nil
}

func (p *UserAppService) GetUserPlan(ctx context.Context, req *userapp_proto.GetUserPlanRequest, rsp *userapp_proto.GetUserPlanResponse) error {
	log.Info("Received UserApp.GetUserPlan request")

	userplan, err := db.GetUserPlan(ctx, req.UserId)
	if userplan == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.GetUserPlan, err, "not found")
	}
	rsp.Data = &userapp_proto.GetUserPlanResponse_Data{userplan}
	return nil
}

func (p *UserAppService) UpdateUserPlan(ctx context.Context, req *userapp_proto.UpdateUserPlanRequest, rsp *userapp_proto.UpdateUserPlanResponse) error {
	log.Info("Received UserApp.UpdateUserPlan request")

	return db.UpdateUserPlan(ctx, req.Id, req.OrgId, req.Goals, req.Days)
}

func (p *UserAppService) GetPlanItemsCountByCategory(ctx context.Context, req *userapp_proto.GetPlanItemsCountByCategoryRequest, rsp *userapp_proto.GetPlanItemsCountByCategoryResponse) error {
	log.Info("Received UserApp.GetPlanItemsCountByCategory request")

	// duration calculate after get user_plan
	userplan, err := db.GetUserPlanWithPlanId(ctx, req.PlanId)
	if userplan == nil || err != nil {
		log.Println("GetUserPlanWithPlanId err:", err)
		return common.NotFound(common.UserappSrv, p.GetPlanItemsCountByCategory, err, "GetPlanWithPlanId err")
	}
	dur, err := duration.FromString(userplan.Duration)
	if err != nil {
		return errors.New("duration parsing error:" + err.Error())
	}
	//
	response, err := db.GetPlanItemsCountByCategory(ctx, req.PlanId, strconv.Itoa(dur.Days))
	if len(response) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetPlanItemsCountByCategory, err, "GetPlanItemsCountByCategory err")
	}
	rsp.Data = &userapp_proto.GetPlanItemsCountByCategoryResponse_Data{response}
	return nil
}

func (p *UserAppService) GetPlanItemsCountByDay(ctx context.Context, req *userapp_proto.GetPlanItemsCountByDayRequest, rsp *userapp_proto.GetPlanItemsCountByDayResponse) error {
	log.Info("Received UserApp.GetPlanItemsCountByDay request")

	response, err := db.GetPlanItemsCountByDay(ctx, req.PlanId, req.DayNumber)
	if len(response) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetPlanItemsCountByDay, err, "not found")
	}
	rsp.Data = &userapp_proto.GetPlanItemsCountByDayResponse_Data{response}
	return nil
}

func (p *UserAppService) GetPlanItemsCountByCategoryAndDay(ctx context.Context, req *userapp_proto.GetPlanItemsCountByCategoryAndDayRequest, rsp *userapp_proto.GetPlanItemsCountByCategoryAndDayResponse) error {
	log.Info("Received UserApp.GetPlanItemsCountByCategoryAndDay request")

	// duration calculate after get user_plan
	userplan, err := db.GetUserPlanWithPlanId(ctx, req.PlanId)
	if userplan == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.GetPlanItemsCountByCategoryAndDay, err, "not found")
	}
	dur, err := duration.FromString(userplan.Duration)
	if err != nil {
		return errors.New("duration parsing error:" + err.Error())
	}
	//
	response, err := db.GetPlanItemsCountByCategoryAndDay(ctx, req.PlanId, strconv.Itoa(dur.Days))
	if len(response) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetPlanItemsCountByCategoryAndDay, err, "not found")
	}
	rsp.Data = &userapp_proto.GetPlanItemsCountByCategoryAndDayResponse_Data{response}
	return nil
}

func (p *UserAppService) Login(ctx context.Context, req *account_proto.LoginRequest, rsp *account_proto.LoginResponse) error {
	log.Info("Received UserApp.Login request")

	rsp_login, err := p.AccountClient.Login(ctx, req)
	if rsp_login == nil || err != nil {
		return common.InternalServerError(common.UserappSrv, p.Login, err, "login fail")
	}
	rsp.Data = rsp_login.Data
	return nil
}

func (p *UserAppService) Logout(ctx context.Context, req *account_proto.LogoutRequest, rsp *account_proto.LogoutResponse) error {
	log.Info("Received UserApp.Logout request")

	_, err := p.AccountClient.Logout(ctx, req)
	return err
}

func (p *UserAppService) AllGoalResponse(ctx context.Context, req *user_proto.AllGoalResponseRequest, rsp *user_proto.AllGoalResponseResponse) error {
	log.Info("Received UserApp.AllGoalResponse request")

	rsp_behaivour, err := p.BehaviourClient.AllGoalResponse(ctx, req)
	if rsp_behaivour == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.AllGoalResponse, err, "not found")
	}
	rsp.Data = rsp_behaivour.Data
	return nil
}

func (p *UserAppService) AllChallengeResponse(ctx context.Context, req *user_proto.AllChallengeResponseRequest, rsp *user_proto.AllChallengeResponseResponse) error {
	log.Info("Received UserApp.AllChallengeResponse request")

	rsp_behaivour, err := p.BehaviourClient.AllChallengeResponse(ctx, req)
	if rsp_behaivour == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.AllChallengeResponse, err, "not found")
	}
	rsp.Data = rsp_behaivour.Data
	return nil
}

func (p *UserAppService) AllHabitResponse(ctx context.Context, req *user_proto.AllHabitResponseRequest, rsp *user_proto.AllHabitResponseResponse) error {
	log.Info("Received UserApp.AllHabitResponse request")

	rsp_behaivour, err := p.BehaviourClient.AllHabitResponse(ctx, req)
	if rsp_behaivour == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.AllHabitResponse, err, "not found")
	}
	rsp.Data = rsp_behaivour.Data
	return nil
}

func (p *UserAppService) GetShareableContent(ctx context.Context, req *user_proto.GetShareableContentRequest, rsp *user_proto.GetShareableContentResponse) error {
	log.Info("Received UserApp.GetShareableContent request")

	rsp_content, err := p.ContentClient.GetShareableContents(ctx, req)
	if len(rsp_content.Data.Contents) == 0 || err != nil {
		return common.NotFound(common.UserappSrv, p.GetShareableContent, err, "not found")
	}
	rsp.Data = rsp_content.Data
	return nil
}

func (p *UserAppService) ReceivedItems(ctx context.Context, req *userapp_proto.ReceivedItemsRequest, rsp *userapp_proto.ReceivedItemsResponse) error {
	log.Info("Received UserApp.ReceivedItems request")
	if len(req.UserId) == 0 {
		return errors.New("user id could not be nil")
	}
	return db.ReceivedItems(ctx, req.UserId, req.Shared)
}

func (p *UserAppService) ReadUser(ctx context.Context, req *user_proto.ReadRequest, rsp *user_proto.ReadResponse) error {
	log.Info("Received UserApp.ReadUser request")
	rsp_user, err := p.UserClient.Read(ctx, req)
	if rsp_user == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.ReadUser, err, "Read User query is failed")
	}
	rsp.Data = rsp_user.Data
	return nil
}

func (p *UserAppService) AutocompleteContentCategoryItem(ctx context.Context, req *content_proto.AutocompleteContentCategoryItemRequest, rsp *content_proto.AutocompleteContentCategoryItemResponse) error {
	log.Info("Received UserApp.AutocompleteContentCategoryItem request")
	rsp_content, err := p.ContentClient.AutocompleteContentCategoryItem(ctx, req)
	if rsp_content == nil || err != nil {
		return common.InternalServerError(common.UserappSrv, p.AutocompleteContentCategoryItem, err, "AutocompleteContentCategoryItem query is failed")
	}
	rsp.Data = rsp_content.Data
	return nil
}

func (p *UserAppService) AllContentCategoryItemByNameslug(ctx context.Context, req *content_proto.AllContentCategoryItemByNameslugRequest, rsp *content_proto.AllContentCategoryItemByNameslugResponse) error {
	log.Info("Received UserApp.AllContentCategoryItemByNameslug request")
	rsp_content, err := p.ContentClient.AllContentCategoryItemByNameslug(ctx, req)
	if rsp_content == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.AllContentCategoryItemByNameslug, err, "not found")
	}
	rsp.Data = rsp_content.Data
	return nil
}

func (p *UserAppService) ReadMarkerByNameslug(ctx context.Context, req *static_proto.ReadByNameslugRequest, rsp *static_proto.ReadMarkerResponse) error {
	log.Info("Received UserApp.MarkerByNameslug request")
	rsp_marker, err := p.StaticClient.ReadMarkerByNameslug(ctx, req)
	if rsp_marker == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.ReadMarkerByNameslug, err, "not found")
	}
	rsp.Data = rsp_marker.Data
	return nil
}
func (p *UserAppService) GetGoalDetail(ctx context.Context, req *behaviour_proto.ReadGoalRequest, rsp *userapp_proto.ReadGoalResponse) error {
	log.Info("Received UserApp.GetGoalDetail request")
	rsp_goal, err := db.GetGoalDetail(ctx, req.GoalId, req.OrgId)
	if rsp_goal == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.GetGoalDetail, err, "not found")
	}
	rsp.Data = &userapp_proto.ReadGoalResponse_Data{rsp_goal}
	return nil
}
func (p *UserAppService) GetChallengeDetail(ctx context.Context, req *behaviour_proto.ReadChallengeRequest, rsp *userapp_proto.ReadChallengeResponse) error {
	log.Info("Received UserApp.GetChallengeDetail request")
	rsp_challenge, err := db.GetChallengeDetail(ctx, req.ChallengeId, req.OrgId)
	if rsp_challenge == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.GetChallengeDetail, err, "not found")
	}
	rsp.Data = &userapp_proto.ReadChallengeResponse_Data{rsp_challenge}
	return nil
}
func (p *UserAppService) GetHabitDetail(ctx context.Context, req *behaviour_proto.ReadHabitRequest, rsp *userapp_proto.ReadHabitResponse) error {
	log.Info("Received UserApp.GetHabitDetail request")
	rsp_habit, err := db.GetHabitDetail(ctx, req.HabitId, req.OrgId)
	if rsp_habit == nil || err != nil {
		return common.NotFound(common.UserappSrv, p.GetHabitDetail, err, "not found")
	}
	rsp.Data = &userapp_proto.ReadHabitResponse_Data{rsp_habit}
	return nil
}

func (p *UserAppService) GetContentDetail(ctx context.Context, req *content_proto.ReadContentRequest, rsp *userapp_proto.ReadContentResponse) error {
	log.Info("Received UserApp.GetContentDetail request")
	rsp_content, err := db.GetContentDetail(ctx, req.Id, req.UserId, req.OrgId)
	if err != nil {
		common.ErrorLog(common.UserappSrv, common.GetFunctionName(p.GetHabitDetail), err, "GetHabitDetail query is failed")
		return err
	}
	rsp.Data = &userapp_proto.ReadContentResponse_Data{rsp_content}
	return nil
}

func (p *UserAppService) RemovePendingSharedAction(ctx context.Context, itemId, userId string) error {
	log.Info("Received UserApp.RemovePendingSharedAction request")
	log.Debug("removing pending action for item.id : ", itemId)
	// delete repective pending from the pending table
	if err := db.RemovePendingSharedAction(ctx, itemId, userId); err != nil {
		return common.InternalServerError(common.UserappSrv, p.RemovePendingSharedAction, err, "RemovePendingSharedAction query is failed")
	}
	return nil
}
