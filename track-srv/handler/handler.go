package handler

import (
	"context"
	"fmt"
	behaviour_proto "server/behaviour-srv/proto/behaviour"
	"server/common"
	content_proto "server/content-srv/proto/content"
	kv_proto "server/kv-srv/proto/kv"
	static_proto "server/static-srv/proto/static"
	"server/track-srv/db"
	track_proto "server/track-srv/proto/track"
	userapp_proto "server/user-app-srv/proto/userapp"
	user_proto "server/user-srv/proto/user"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"
	google_protobuf1 "github.com/golang/protobuf/ptypes/struct"
	"github.com/micro/go-micro/broker"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type TrackService struct {
	Broker          broker.Broker
	KvClient        kv_proto.KvServiceClient
	BehaviourClient behaviour_proto.BehaviourServiceClient
	ContentClient   content_proto.ContentServiceClient
	StaticClient    static_proto.StaticServiceClient
	UserClient      user_proto.UserServiceClient
}

// func (p *TrackService) GetTrackKey(user_id, obj_id, t string) string {
// 	return user_id + ":" + obj_id + ":" + t
// }

func (p *TrackService) CreateTrackGoal(ctx context.Context, req *track_proto.CreateTrackGoalRequest, rsp *track_proto.CreateTrackGoalResponse) error {
	log.Info("Received Track.CreateTrackGoal request")
	// fetch goal from behaviour-srv
	rsp_goal, err := p.BehaviourClient.ReadGoal(ctx, &behaviour_proto.ReadGoalRequest{
		GoalId: req.GoalId,
		OrgId:  req.OrgId,
		TeamId: req.TeamId,
	})
	if rsp_goal == nil || err != nil {
		return common.NotFound(common.TrackSrv, p.CreateTrackGoal, err, "not found")
	}
	// create track goal
	trackGoal := &track_proto.TrackGoal{
		User:    req.User,
		OrgId:   req.OrgId,
		Goal:    rsp_goal.Data.Goal,
		Created: time.Now().Unix(),
	}
	if err := db.CreateTrackGoal(ctx, trackGoal); err != nil {
		return common.InternalServerError(common.TrackSrv, p.CreateTrackGoal, err, "create error")
	}
	// increment track
	rsp_inc, err := p.KvClient.IncTrackCount(ctx, &kv_proto.IncTrackCountRequest{common.TRACK_INDEX, common.GetTrackKey(req.User.Id, rsp_goal.Data.Goal.Id, "goal")})
	if err != nil {
		return common.NotFound(common.TrackSrv, p.CreateTrackGoal, err, "track count inc error")
	}
	// update status of join goal
	if err := db.UpdateJoinGoalStatus(ctx, req.GoalId, userapp_proto.ActionStatus_IN_PROGRESS); err != nil {
		return common.InternalServerError(common.TrackSrv, p.CreateTrackGoal, err, "update error")
	}

	rsp.Data = &track_proto.CreateTrackGoalResponse_Data{trackGoal, rsp_inc.Count}
	return nil
}

func (p *TrackService) GetGoalCount(ctx context.Context, req *track_proto.GetGoalCountRequest, rsp *track_proto.GetGoalCountResponse) error {
	log.Info("Received Track.GetGoalCount request")
	key := common.GetTrackKey(req.UserId, req.GoalId, "goal")
	var count int64
	rsp_get, err := p.KvClient.GetTrackCount(ctx, &kv_proto.GetTrackCountRequest{common.TRACK_INDEX, key})
	if rsp_get.Count == 0 || err != nil {
		// fetching count from db
		c, err := db.GetGoalCount(ctx, req.UserId, req.GoalId, req.From, req.To)
		if err != nil {
			return common.NotFound(common.TrackSrv, p.GetGoalCount, err, "not found")
		}
		if _, err := p.KvClient.SetTrackCount(ctx, &kv_proto.SetTrackCountRequest{common.TRACK_INDEX, key, count}); err != nil {
			return common.InternalServerError(common.TrackSrv, p.GetGoalCount, err, "track count set error")
		}
		count = c
	} else {
		count = rsp_get.Count
	}
	rsp.Data = &track_proto.GetGoalCountResponse_Data{count}
	return nil
}

func (p *TrackService) GetGoalHistory(ctx context.Context, req *track_proto.GetGoalHistoryRequest, rsp *track_proto.GetGoalHistoryResponse) error {
	log.Info("Received Track.GetGoalHistory request")
	trackGoals, err := db.GetGoalHistory(ctx, req)
	if len(trackGoals) == 0 || err != nil {
		return common.NotFound(common.TrackSrv, p.GetGoalHistory, err, "not found")
	}
	rsp.Data = &track_proto.GetGoalHistoryResponse_Data{trackGoals}
	return nil
}

func (p *TrackService) CreateTrackChallenge(ctx context.Context, req *track_proto.CreateTrackChallengeRequest, rsp *track_proto.CreateTrackChallengeResponse) error {
	log.Info("Received Track.CreateTrackChallenge request")
	// fetch challenge from behaviour-srv
	rsp_challenge, err := p.BehaviourClient.ReadChallenge(ctx, &behaviour_proto.ReadChallengeRequest{
		ChallengeId: req.ChallengeId,
		OrgId:       req.OrgId,
		TeamId:      req.TeamId,
	})
	if err != nil {
		return common.NotFound(common.TrackSrv, p.CreateTrackChallenge, err, "not found")
	}
	// create track challenge
	trackChallenge := &track_proto.TrackChallenge{
		User:      req.User,
		OrgId:     req.OrgId,
		Challenge: rsp_challenge.Data.Challenge,
		Created:   time.Now().Unix(),
	}
	// log.Println(trackChallenge)
	if err := db.CreateTrackChallenge(ctx, trackChallenge); err != nil {
		return common.InternalServerError(common.TrackSrv, p.CreateTrackChallenge, err, "create error")
	}
	// increment track
	rsp_inc, err := p.KvClient.IncTrackCount(ctx, &kv_proto.IncTrackCountRequest{common.TRACK_INDEX, common.GetTrackKey(req.User.Id, rsp_challenge.Data.Challenge.Id, "challenge")})
	if err != nil {
		return common.NotFound(common.TrackSrv, p.CreateTrackChallenge, err, "track count inc error")
	}
	// update status of join challenge
	if err := db.UpdateJoinChallengeStatus(ctx, req.ChallengeId, userapp_proto.ActionStatus_IN_PROGRESS); err != nil {
		return common.InternalServerError(common.TrackSrv, p.CreateTrackChallenge, err, "update error")
	}

	rsp.Data = &track_proto.CreateTrackChallengeResponse_Data{trackChallenge, rsp_inc.Count}
	return nil
}

func (p *TrackService) GetChallengeCount(ctx context.Context, req *track_proto.GetChallengeCountRequest, rsp *track_proto.GetChallengeCountResponse) error {
	log.Info("Received Track.GetChallengeCount request")
	key := common.GetTrackKey(req.UserId, req.ChallengeId, "challenge")
	var count int64
	rsp_get, err := p.KvClient.GetTrackCount(ctx, &kv_proto.GetTrackCountRequest{common.TRACK_INDEX, key})
	if err != nil {
		// fetching count from db
		count, err = db.GetChallengeCount(ctx, req.UserId, req.ChallengeId, req.From, req.To)
		if count == 0 || err != nil {
			return common.NotFound(common.TrackSrv, p.GetChallengeCount, err, "not found")
		}
		if _, err := p.KvClient.SetTrackCount(ctx, &kv_proto.SetTrackCountRequest{common.TRACK_INDEX, key, count}); err != nil {
			return common.InternalServerError(common.TrackSrv, p.GetChallengeCount, err, "track count set error")
		}
	} else {
		count = rsp_get.Count
	}
	rsp.Data = &track_proto.GetChallengeCountResponse_Data{count}
	return nil
}

func (p *TrackService) GetChallengeHistory(ctx context.Context, req *track_proto.GetChallengeHistoryRequest, rsp *track_proto.GetChallengeHistoryResponse) error {
	log.Info("Received Track.GetChallengeHistory request")
	trackChallenges, err := db.GetChallengeHistory(ctx, req)
	if len(trackChallenges) == 0 || err != nil {
		return common.NotFound(common.TrackSrv, p.GetChallengeHistory, err, "not found")
	}
	rsp.Data = &track_proto.GetChallengeHistoryResponse_Data{trackChallenges}
	return nil
}

func (p *TrackService) CreateTrackHabit(ctx context.Context, req *track_proto.CreateTrackHabitRequest, rsp *track_proto.CreateTrackHabitResponse) error {
	log.Info("Received Track.CreateTrackHabit request")
	// fetch habit from behaviour-srv
	rsp_habit, err := p.BehaviourClient.ReadHabit(ctx, &behaviour_proto.ReadHabitRequest{
		HabitId: req.HabitId,
		OrgId:   req.OrgId,
		TeamId:  req.TeamId,
	})
	if err != nil {
		return common.NotFound(common.TrackSrv, p.CreateTrackHabit, err, "not found")
	}
	// create track habit
	trackHabit := &track_proto.TrackHabit{
		User:    req.User,
		OrgId:   req.OrgId,
		Habit:   rsp_habit.Data.Habit,
		Created: time.Now().Unix(),
	}
	// log.Println(trackHabit)
	if err := db.CreateTrackHabit(ctx, trackHabit); err != nil {
		return common.InternalServerError(common.TrackSrv, p.CreateTrackHabit, err, "create error")
	}
	// increment track
	rsp_inc, err := p.KvClient.IncTrackCount(ctx, &kv_proto.IncTrackCountRequest{common.TRACK_INDEX, common.GetTrackKey(req.User.Id, rsp_habit.Data.Habit.Id, "habit")})
	if err != nil {
		return common.NotFound(common.TrackSrv, p.CreateTrackHabit, err, "track count inc error")
	}
	// update status of join habit
	if err := db.UpdateJoinHabitStatus(ctx, req.HabitId, userapp_proto.ActionStatus_IN_PROGRESS); err != nil {
		return common.InternalServerError(common.TrackSrv, p.CreateTrackHabit, err, "update error")
	}

	rsp.Data = &track_proto.CreateTrackHabitResponse_Data{trackHabit, rsp_inc.Count}
	return nil
}

func (p *TrackService) GetHabitCount(ctx context.Context, req *track_proto.GetHabitCountRequest, rsp *track_proto.GetHabitCountResponse) error {
	log.Info("Received Track.GetHabitCount request")
	key := common.GetTrackKey(req.UserId, req.HabitId, "habit")
	var count int64
	rsp_get, err := p.KvClient.GetTrackCount(ctx, &kv_proto.GetTrackCountRequest{common.TRACK_INDEX, key})
	if err != nil {
		// fetching count from db
		count, err = db.GetHabitCount(ctx, req.UserId, req.HabitId, req.From, req.To)
		if count == 0 || err != nil {
			return common.NotFound(common.TrackSrv, p.GetHabitCount, err, "not found")
		}
		if _, err := p.KvClient.SetTrackCount(ctx, &kv_proto.SetTrackCountRequest{common.TRACK_INDEX, key, count}); err != nil {
			return common.InternalServerError(common.TrackSrv, p.GetHabitCount, err, "track count set error")
		}
	} else {
		count = rsp_get.Count
	}
	rsp.Data = &track_proto.GetHabitCountResponse_Data{count}
	return nil
}

func (p *TrackService) GetHabitHistory(ctx context.Context, req *track_proto.GetHabitHistoryRequest, rsp *track_proto.GetHabitHistoryResponse) error {
	log.Info("Received Track.GetHabitHistory request")
	trackHabits, err := db.GetHabitHistory(ctx, req)
	if len(trackHabits) == 0 || err != nil {
		return common.NotFound(common.TrackSrv, p.GetHabitHistory, err, "not found")
	}
	rsp.Data = &track_proto.GetHabitHistoryResponse_Data{trackHabits}
	return nil
}

func (p *TrackService) CreateTrackContent(ctx context.Context, req *track_proto.CreateTrackContentRequest, rsp *track_proto.CreateTrackContentResponse) error {
	log.Info("Received Track.CreateTrackContent request")
	// fetch content from behaviour-srv
	rsp_content, err := p.ContentClient.ReadContent(ctx, &content_proto.ReadContentRequest{
		Id:     req.ContentId,
		OrgId:  req.OrgId,
		TeamId: req.TeamId,
	})
	if rsp_content == nil || err != nil {
		return common.NotFound(common.TrackSrv, p.CreateTrackContent, err, "not found")
	}
	content := rsp_content.Data.Content
	if content == nil {
		return common.BadRequest(common.TrackSrv, p.CreateTrackContent, nil, "content empty")
	}
	// handler logic
	if content.Category == nil {
		return common.BadRequest(common.TrackSrv, p.CreateTrackContent, nil, "content category empty")
	}
	if content.Category.TrackerMethods == nil || len(content.Category.TrackerMethods) == 0 {
		return common.BadRequest(common.TrackSrv, p.CreateTrackContent, nil, "category tracker method empty")
	}
	// check count name_slug
	trackerMethod := content.Category.TrackerMethods[0]
	// create track content
	trackContent := &track_proto.TrackContent{
		User:    req.User,
		OrgId:   req.OrgId,
		Content: content,
		Created: time.Now().Unix(),
	}
	// log.Println(trackContent)
	if err := db.CreateTrackContent(ctx, trackContent); err != nil {
		return common.InternalServerError(common.TrackSrv, p.CreateTrackContent, err, "create error")
	}
	// increment track
	rsp_inc, err := p.KvClient.IncTrackCount(ctx, &kv_proto.IncTrackCountRequest{common.TRACK_INDEX, common.GetTrackKey(req.User.Id, rsp_content.Data.Content.Id, trackerMethod.NameSlug)})
	if err != nil {
		return common.NotFound(common.TrackSrv, p.CreateTrackContent, err, "not found")
	}
	rsp.Data = &track_proto.CreateTrackContentResponse_Data{trackContent, rsp_inc.Count}
	return nil
}

func (p *TrackService) GetContentCount(ctx context.Context, req *track_proto.GetContentCountRequest, rsp *track_proto.GetContentCountResponse) error {
	log.Info("Received Track.GetContentCount request")
	// fetch content from behaviour-srv
	rsp_content, err := p.ContentClient.ReadContent(ctx, &content_proto.ReadContentRequest{
		Id:     req.ContentId,
		OrgId:  req.OrgId,
		TeamId: req.TeamId,
	})
	if err != nil {
		return common.NotFound(common.TrackSrv, p.GetContentCount, err, "not found")
	}
	content := rsp_content.Data.Content
	if content == nil {
		return common.BadRequest(common.TrackSrv, p.GetContentCount, err, "content empty")
	}
	// handler logic
	if content.Category == nil {
		return common.BadRequest(common.TrackSrv, p.GetContentCount, err, "category empty")
	}
	if content.Category.TrackerMethods == nil || len(content.Category.TrackerMethods) == 0 {
		return common.BadRequest(common.TrackSrv, p.GetContentCount, err, "tracker method empty")
	}
	// check count name_slug
	trackerMethod := content.Category.TrackerMethods[0]
	key := common.GetTrackKey(req.UserId, req.ContentId, trackerMethod.NameSlug)

	var count int64
	rsp_get, err := p.KvClient.GetTrackCount(ctx, &kv_proto.GetTrackCountRequest{common.TRACK_INDEX, key})
	if err != nil {
		// fetching count from db
		count, err = db.GetContentCount(ctx, req.UserId, req.ContentId, req.From, req.To)
		if count == 0 || err != nil {
			return common.NotFound(common.TrackSrv, p.GetContentCount, err, "not found")
		}
		if _, err := p.KvClient.SetTrackCount(ctx, &kv_proto.SetTrackCountRequest{common.TRACK_INDEX, key, count}); err != nil {
			return common.InternalServerError(common.TrackSrv, p.GetContentCount, err, "set track count error")
		}
	} else {
		count = rsp_get.Count
	}
	rsp.Data = &track_proto.GetContentCountResponse_Data{count}
	return nil
}

func (p *TrackService) GetContentHistory(ctx context.Context, req *track_proto.GetContentHistoryRequest, rsp *track_proto.GetContentHistoryResponse) error {
	log.Info("Received Track.GetContentHistory request")
	trackContents, err := db.GetContentHistory(ctx, req)
	if len(trackContents) == 0 || err != nil {
		return common.NotFound(common.TrackSrv, p.GetContentHistory, err, "not found")
	}
	rsp.Data = &track_proto.GetContentHistoryResponse_Data{trackContents}
	return nil
}

func (p *TrackService) SubscribeMarker(ctx context.Context, req *track_proto.SubscribeMarkerRequest, rsp *track_proto.SubscribeMarkerResponse) error {
	log.Info("Received Track.SubscribeMarker request")

	_, err := p.Broker.Subscribe(req.Channel, func(pub broker.Publication) error {
		// save track-marker value
		return nil
	})
	return err
}

func (p *TrackService) CreateTrackMarker(ctx context.Context, req *track_proto.CreateTrackMarkerRequest, rsp *track_proto.CreateTrackMarkerResponse) error {
	log.Info("Received Track.CreateTrackMarker request")

	//get valid marker
	req_marker := &static_proto.ReadMarkerRequest{Id: req.MarkerId}
	rsp_marker, err := p.StaticClient.ReadMarker(ctx, req_marker)
	if err != nil {
		return common.NotFound(common.TrackSrv, p.CreateTrackMarker, err, "not found")
	}
	if rsp_marker.Data.Marker == nil {
		return common.BadRequest(common.TrackSrv, p.CreateTrackMarker, err, "data marker empty")
	}
	//get valid user
	req_user := &user_proto.ReadRequest{UserId: req.UserId, OrgId: req.OrgId}
	rsp_user, err := p.UserClient.Read(ctx, req_user)
	if rsp_user == nil && err != nil {
		return common.NotFound(common.TrackSrv, p.CreateTrackMarker, err, "not found")
	}

	if rsp_user.Data.User == nil {
		return common.BadRequest(common.TrackSrv, p.CreateTrackMarker, err, "data user empty")
	}

	//get valid tracker_method
	log.Debug("req.TrackerMethod:", req.TrackerMethod)
	req_tracker_method := &static_proto.ReadTrackerMethodRequest{Id: req.TrackerMethod.Id, NameSlug: req.TrackerMethod.NameSlug}
	rsp_tracker_method, err := p.StaticClient.ReadTrackerMethod(ctx, req_tracker_method)
	if err != nil {
		return common.NotFound(common.TrackSrv, p.CreateTrackMarker, err, "not found")
	}
	if rsp_tracker_method.Data.TrackerMethod == nil {
		return common.BadRequest(common.TrackSrv, p.CreateTrackMarker, err, "tracker method empty")
	}

	trackMarker := &track_proto.TrackMarker{}
	trackerMethod := rsp_tracker_method.Data.TrackerMethod

	marshaler := jsonpb.Marshaler{}

	log.Info("tracking marker using mode: ", trackerMethod.NameSlug)
	switch trackerMethod.NameSlug {
	case "count":
		// increment track
		if _, err := p.KvClient.IncTrackCount(ctx, &kv_proto.IncTrackCountRequest{common.TRACK_INDEX, common.GetTrackKey(req.UserId, req.MarkerId, trackerMethod.NameSlug)}); err != nil {
			return common.NotFound(common.TrackSrv, p.CreateTrackMarker, err, "track count inc error")
		}
		//create the request to store count value in the db for count mode
		trackMarker = &track_proto.TrackMarker{
			Id:      uuid.NewUUID().String(),
			User:    rsp_user.Data.User,
			Marker:  rsp_marker.Data.Marker,
			Created: time.Now().Unix(),
			OrgId:   req.OrgId,
		}
	case "manual", "hcp":
		//create the request to store value to store in the db for manual mode
		log.Info("tracking marker", req.MarkerId)
		js, err := marshaler.MarshalToString(req.Value)
		if err != nil {
			return common.InternalServerError(common.TrackSrv, p.CreateTrackMarker, err, "parsing error")
		}
		if req.Value == nil {
			return common.BadRequest(common.TrackSrv, p.CreateTrackMarker, nil, "value empty")
		}

		item := &kv_proto.Item{
			Key:   common.GetTrackKey(req.UserId, req.MarkerId, trackerMethod.NameSlug),
			Value: []byte(js),
		}
		if _, err := p.KvClient.PutEx(ctx, &kv_proto.PutExRequest{common.TRACK_INDEX, item}); err != nil {
			return common.InternalServerError(common.TrackSrv, p.CreateTrackMarker, err, "update error")
		}

		//create the request to store value in the db for manual mode
		trackMarker = &track_proto.TrackMarker{
			Id:      uuid.NewUUID().String(),
			User:    rsp_user.Data.User,
			Marker:  rsp_marker.Data.Marker,
			Created: time.Now().Unix(),
			Value:   req.Value,
			Unit:    req.Unit,
			OrgId:   req.OrgId,
		}

	case "photo":
		//create the request to store value to store in the db for manual mode
		js, err := marshaler.MarshalToString(req.Value)
		if err != nil {
			return common.InternalServerError(common.TrackSrv, p.CreateTrackMarker, err, "parsing error")
		}
		if req.Value == nil {
			return common.BadRequest(common.TrackSrv, p.CreateTrackMarker, nil, "value empty")
		}
		if err := p.Broker.Publish(common.UPLOAD_FILE, &broker.Message{Body: []byte(js)}); err != nil {
			return common.InternalServerError(common.TrackSrv, p.CreateTrackMarker, err, "upload file subcribe error")
		}
		// subscrib logic is missing
		_, err = p.Broker.Subscribe(common.FILE_UPLOADED, func(pub broker.Publication) error {
			return common.InternalServerError(common.TrackSrv, p.CreateTrackMarker, err, "file upload subscribe error")
		})
		if err != nil {
			return common.InternalServerError(common.TrackSrv, p.CreateTrackMarker, err, "subscribe error")
		}
	}

	//create document in the track collection
	if err := db.CreateTrackMarker(ctx, trackMarker); err != nil {
		return common.InternalServerError(common.TrackSrv, p.CreateTrackMarker, err, "create error")
	}
	rsp.Data = &track_proto.CreateTrackMarkerResponse_Data{trackMarker}
	return nil
}

func (p *TrackService) GetLastMarker(ctx context.Context, req *track_proto.GetLastMarkerRequest, rsp *track_proto.GetLastMarkerResponse) error {
	log.Info("Received Track.GetLastMarker request")

	track_marker, err := db.GetLastMarker(ctx, req.MarkerId, req.UserId)
	if track_marker == nil || err != nil {
		return common.NotFound(common.TrackSrv, p.GetLastMarker, err, "not found")
	}

	if len(track_marker.Marker.TrackerMethods) > 0 {
		trackerMethod := track_marker.Marker.TrackerMethods[0]
		key := common.GetTrackKey(track_marker.User.Id, track_marker.Marker.Id, trackerMethod.NameSlug)
		value := google_protobuf1.Value{}
		switch trackerMethod.NameSlug {
		case "count":
			rsp_get, err := p.KvClient.GetTrackCount(ctx, &kv_proto.GetTrackCountRequest{common.TRACK_INDEX, key})

			count := 0
			if err == nil {
				count = int(rsp_get.Count)
			} else {
				if _, err := p.KvClient.SetTrackCount(ctx, &kv_proto.SetTrackCountRequest{common.TRACK_INDEX, key, 0}); err != nil {
					return common.InternalServerError(common.TrackSrv, p.GetLastMarker, err, "track count set error")
				}
			}
			// unmarshal to object
			if err := jsonpb.Unmarshal(strings.NewReader(fmt.Sprintf("%v", count)), &value); err != nil {
				return common.InternalServerError(common.TrackSrv, p.GetLastMarker, err, "parsing error")
			}
			// making response
			rsp.Data = &track_proto.GetLastMarkerResponse_Data{&value}
		case "manual", "photo":
			rsp_get, err := p.KvClient.GetTrackValue(ctx, &kv_proto.GetTrackValueRequest{common.TRACK_INDEX, key})
			if err != nil {
				marshaler := jsonpb.Marshaler{}
				js, err := marshaler.MarshalToString(track_marker.Value)
				if err != nil {
					return common.InternalServerError(common.TrackSrv, p.GetLastMarker, err, "parsing error")
				}
				item := &kv_proto.Item{
					Key:   common.GetTrackKey(track_marker.User.Id, req.MarkerId, trackerMethod.NameSlug),
					Value: []byte(js),
				}
				if _, err := p.KvClient.PutEx(ctx, &kv_proto.PutExRequest{common.TRACK_INDEX, item}); err != nil {
					return common.InternalServerError(common.TrackSrv, p.GetLastMarker, err, "update error")
				}
				// unmarshal value
				if err := jsonpb.Unmarshal(strings.NewReader(js), &value); err != nil {
					return common.InternalServerError(common.TrackSrv, p.GetLastMarker, err, "parsing error")
				}
				rsp.Data = &track_proto.GetLastMarkerResponse_Data{&value}
			} else {
				rsp.Data = &track_proto.GetLastMarkerResponse_Data{rsp_get.Value}
			}
		}
	}

	return nil
}

func (p *TrackService) GetMarkerHistory(ctx context.Context, req *track_proto.GetMarkerHistoryRequest, rsp *track_proto.GetMarkerHistoryResponse) error {
	log.Info("Received Track.GetMarkerHistory request")

	markers, err := db.GetMarkerHistory(ctx, req.MarkerId, req.From, req.To, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(markers) == 0 || err != nil {
		return common.NotFound(common.TrackSrv, p.GetMarkerHistory, err, "not found")
	}
	rsp.Data = &track_proto.GetMarkerHistoryResponse_Data{markers}
	return nil
}

func (p *TrackService) GetAllMarkerHistory(ctx context.Context, req *track_proto.GetAllMarkerHistoryRequest, rsp *track_proto.GetAllMarkerHistoryResponse) error {
	log.Info("Received Track.GetAllMarkerHistory request")

	markers, err := db.GetAllMarkerHistory(ctx, req.From, req.To, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(markers) == 0 || err != nil {
		return common.NotFound(common.TrackSrv, p.GetAllMarkerHistory, err, "not found")
	}
	rsp.Data = &track_proto.GetAllMarkerHistoryResponse_Data{markers}
	return nil
}

func (p *TrackService) GetDefaultMarkerHistory(ctx context.Context, req *track_proto.GetDefaultMarkerHistoryRequest, rsp *track_proto.GetDefaultMarkerHistoryResponse) error {
	log.Info("Received Track.GetDefaultMarkerHistory request")

	markers, err := db.GetDefaultMarkerHistory(ctx, req.UserId, req.Offset, req.Limit, req.From, req.To)
	if len(markers) == 0 || err != nil {
		return common.NotFound(common.TrackSrv, p.GetDefaultMarkerHistory, err, "not found")
	}
	rsp.Data = &track_proto.GetDefaultMarkerHistoryResponse_Data{markers}
	return nil
}
