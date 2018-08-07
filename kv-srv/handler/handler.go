package handler

import (
	"context"
	"fmt"
	kv_proto "server/kv-srv/proto/kv"
	"strings"
	"time"

	"server/common"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/jsonpb"
	google_protobuf1 "github.com/golang/protobuf/ptypes/struct"
	log "github.com/sirupsen/logrus"
)

type KvService struct {
	Client *redis.Client
}

func (p *KvService) Get(ctx context.Context, req *kv_proto.GetRequest, rsp *kv_proto.GetResponse) error {
	log.Info("Received Kv.Get request")
	val, err := p.Client.Get(req.Key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return common.NotFound(common.KvSrv, p.Get, err, "Key not found: "+req.Key)
		} else {
			return common.InternalServerError(common.KvSrv, p.Get, err, "get request error")
		}
	}

	if val == nil {
		return common.InternalServerError(common.KvSrv, p.Get, nil, "can't get val")
	}

	d, err := p.Client.TTL(req.Key).Result()
	if err != nil {
		return common.InternalServerError(common.KvSrv, p.Get, err, "parsing error")
	}

	rsp.Item = &kv_proto.Item{
		Key:        req.Key,
		Value:      val,
		Expiration: int64(d),
	}

	return nil
}

func (p *KvService) Put(ctx context.Context, req *kv_proto.PutRequest, rsp *kv_proto.PutResponse) error {
	log.Info("Received Kv.Put request")
	if err := p.Client.Set(req.Item.Key, req.Item.Value, time.Duration(req.Item.Expiration)).Err(); err != nil {
		return common.InternalServerError(common.KvSrv, p.Put, err, "put request error")
	}
	return nil
}

func (p *KvService) Del(ctx context.Context, req *kv_proto.DelRequest, rsp *kv_proto.DelResponse) error {
	log.Info("Received Kv.Del request")
	if err := p.Client.Del(req.Key).Err(); err != nil {
		return common.InternalServerError(common.KvSrv, p.Del, err, "delete request error")
	}
	return nil
}

func (p *KvService) GetEx(ctx context.Context, req *kv_proto.GetExRequest, rsp *kv_proto.GetExResponse) error {
	log.Info("Received Kv.GetEx request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	req_put := &kv_proto.GetRequest{req.Key}
	rsp_put := &kv_proto.GetResponse{}
	if err := p.Get(ctx, req_put, rsp_put); err != nil {
		return common.InternalServerError(common.KvSrv, p.GetEx, err, "get request error")
	}

	rsp.Item = rsp_put.Item
	return nil
}

func (p *KvService) PutEx(ctx context.Context, req *kv_proto.PutExRequest, rsp *kv_proto.PutExResponse) error {
	log.Info("Received Kv.PutEx request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	// req_put := &kv_proto.PutRequest{req.Item}
	// rsp_put := &kv_proto.PutResponse{}

	if err := p.Client.Set(req.Item.Key, req.Item.Value, time.Duration(req.Item.Expiration)).Err(); err != nil {
		return common.InternalServerError(common.KvSrv, p.PutEx, err, "put request error")
	}
	return nil
}

func (p *KvService) DelEx(ctx context.Context, req *kv_proto.DelExRequest, rsp *kv_proto.DelExResponse) error {
	log.Info("Received Kv.DelEx request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	req_put := &kv_proto.DelRequest{req.Key}
	rsp_put := &kv_proto.DelResponse{}
	if err := p.Del(ctx, req_put, rsp_put); err != nil {
		return common.InternalServerError(common.KvSrv, p.DelEx, err, "request error")
	}
	return nil
}

func (p *KvService) ConfirmToken(ctx context.Context, req *kv_proto.ConfirmTokenRequest, rsp *kv_proto.ConfirmTokenResponse) error {
	log.Info("Received Kv.ConfirmToken request")
	cmd := redis.NewStringCmd("SELECT", 0)
	p.Client.Process(cmd)

	rsp_get := &kv_proto.GetResponse{}
	if err := p.Get(ctx, &kv_proto.GetRequest{req.Key}, rsp_get); err != nil {
		return common.NotFound(common.KvSrv, p.ConfirmToken, err, "not found")
	}

	rsp.Value = string(rsp_get.Item.Value)
	return nil
}

func (p *KvService) RemoveToken(ctx context.Context, req *kv_proto.RemoveTokenRequest, rsp *kv_proto.RemoveTokenResponse) error {
	log.Info("Received Kv.RemoveToken request")
	cmd := redis.NewStringCmd("SELECT", common.VERIFICATION_TOKEN_INDEX)
	p.Client.Process(cmd)

	// remove token from redis after token
	rsp_del := &kv_proto.DelResponse{}
	log.Info("Removing confirmation token: ", req.Token)
	if err := p.Del(ctx, &kv_proto.DelRequest{Key: req.Token}, rsp_del); err != nil {
		return common.InternalServerError(common.KvSrv, p.ConfirmToken, err, "delete response fail")
	}
	//remove account_id for the token
	log.Info("Removing account: ", req.AccountId)
	if err := p.Del(ctx, &kv_proto.DelRequest{Key: req.AccountId}, rsp_del); err != nil {
		return err
	}

	return nil
}

func (p *KvService) Locked(ctx context.Context, req *kv_proto.LockedRequest, rsp *kv_proto.LockedResponse) error {
	log.Info("Received Kv.Locked request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	rsp_put := &kv_proto.PutExResponse{}
	err := p.PutEx(ctx, &kv_proto.PutExRequest{
		Index: req.Index,
		Item: &kv_proto.Item{
			Key:        req.AccountId,
			Value:      []byte("true"),
			Expiration: int64(24 * time.Hour),
		}}, rsp_put)
	return err
}

func (p *KvService) UnLock(ctx context.Context, req *kv_proto.UnLockRequest, rsp *kv_proto.UnLockResponse) error {
	log.Info("Received Kv.UnLock request")

	//remove key from locked index
	rsp_del := &kv_proto.DelExResponse{}
	err := p.DelEx(ctx, &kv_proto.DelExRequest{Index: req.LockIndex, Key: req.AccountId}, rsp_del)
	if err != nil {
		return err
	}

	//set auth failure to 0 in auth index
	err_auth_fail := p.DelEx(ctx, &kv_proto.DelExRequest{Index: req.AuthFailIndex, Key: req.AccountId}, rsp_del)
	if err_auth_fail != nil {
		return err_auth_fail
	}
	return nil
}

func (p *KvService) IsLocked(ctx context.Context, req *kv_proto.IsLockedRequest, rsp *kv_proto.IsLockedResponse) error {
	log.Info("Received Kv.IsLocked request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	rsp_get := &kv_proto.GetResponse{}
	err := p.Get(ctx, &kv_proto.GetRequest{req.AccountId}, rsp_get)
	if err != nil {
		rsp.Locked = false
	} else {
		rsp.Locked = true
	}
	return nil
}

func (p *KvService) AuthFailed(ctx context.Context, req *kv_proto.AuthFailedRequest, rsp *kv_proto.AuthFailedResponse) error {
	log.Info("Received Kv.AuthFailed request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	rsp_get := &kv_proto.GetResponse{}
	err := p.Get(ctx, &kv_proto.GetRequest{req.AccountId}, rsp_get)
	if err != nil {
		// init failed auth
		req_put := &kv_proto.PutRequest{&kv_proto.Item{
			Key:        req.AccountId,
			Value:      []byte("1"),
			Expiration: int64(time.Hour),
		}}
		rsp_put := &kv_proto.PutResponse{}
		if err := p.Put(ctx, req_put, rsp_put); err != nil {
			return common.InternalServerError(common.KvSrv, p.AuthFailed, err, "auth failed error")
		}
	}

	// increment
	val := p.Client.Incr(req.AccountId).Val()
	rsp.Failed = val
	return nil
}

func (p *KvService) ReadSession(ctx context.Context, req *kv_proto.ReadSessionRequest, rsp *kv_proto.ReadSessionResponse) error {
	log.Info("Received Kv.ReadSession request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	d, err := p.Client.Get(req.SessionId).Result()
	if err != nil {
		if err == redis.Nil {
			return common.Unauthorized(common.KvSrv, p.ReadSession, err, "Not Authorized")
		} else {
			return common.NotFound(common.KvSrv, p.ReadSession, err, "session not found")
		}
	}

	rsp.Value = d
	return nil
}

func (p *KvService) RemoveSession(ctx context.Context, req *kv_proto.RemoveSessionRequest, rsp *kv_proto.RemoveSessionResponse) error {
	log.Info("Received Kv.RemoveSession request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	rsp_del := &kv_proto.DelResponse{}
	return p.Del(ctx, &kv_proto.DelRequest{req.SessionId}, rsp_del)
}

func (p *KvService) IncTrackCount(ctx context.Context, req *kv_proto.IncTrackCountRequest, rsp *kv_proto.IncTrackCountResponse) error {
	log.Info("Received Kv.IncTrackCount request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	// increment
	val := p.Client.Incr(req.Key).Val()
	rsp.Count = val
	return nil
}

func (p *KvService) GetTrackCount(ctx context.Context, req *kv_proto.GetTrackCountRequest, rsp *kv_proto.GetTrackCountResponse) error {
	log.Info("Received Kv.GetTrackCount request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	val, err := p.Client.Get(req.Key).Int64()
	if err != nil {
		return common.NotFound(common.KvSrv, p.GetTrackCount, err, "track count not found")
	}
	rsp.Count = val
	return nil
}

func (p *KvService) SetTrackCount(ctx context.Context, req *kv_proto.SetTrackCountRequest, rsp *kv_proto.SetTrackCountResponse) error {
	log.Info("Received Kv.SetTrackCount request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	if err := p.Client.Set(req.Key, req.Count, 0).Err(); err != nil {
		return common.InternalServerError(common.KvSrv, p.SetTrackCount, err, "track count can't set")
	}
	return nil
}

func (p *KvService) IncBadgeCount(ctx context.Context, req *kv_proto.IncBadgeCountRequest, rsp *kv_proto.IncBadgeCountResponse) error {
	log.Info("Received Kv.IncBadgeCount request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	// increment
	val := p.Client.Incr(req.Key).Val()
	rsp.Count = val
	return nil
}

func (p *KvService) DecrBadgeCount(ctx context.Context, req *kv_proto.DecrBadgeCountRequest, rsp *kv_proto.DecrBadgeCountResponse) error {
	log.Info("Received Kv.DecrBadgeCount request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	// increment
	val := p.Client.Incr(req.Key).Val()
	rsp.Count = val
	return nil
}

func (p *KvService) GetBadgeCount(ctx context.Context, req *kv_proto.GetBadgeCountRequest, rsp *kv_proto.GetBadgeCountResponse) error {
	log.Info("Received Kv.GetBadgeCount request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	val, err := p.Client.Get(req.Key).Int64()
	if err != nil {
		return err
	}
	rsp.Count = val
	return nil
}

func (p *KvService) SetBadgeCount(ctx context.Context, req *kv_proto.SetBadgeCountRequest, rsp *kv_proto.SetBadgeCountResponse) error {
	log.Info("Received Kv.SetBadgeCount request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	if err := p.Client.Set(req.Key, req.Count, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (p *KvService) GetTrackValue(ctx context.Context, req *kv_proto.GetTrackValueRequest, rsp *kv_proto.GetTrackValueResponse) error {
	log.Info("Received Kv.GetTrackValue request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	d, err := p.Client.Get(req.Key).Bytes()
	if err != nil {
		return common.NotFound(common.KvSrv, p.GetTrackValue, err, "track value not found")
	}
	var v google_protobuf1.Value
	if err := jsonpb.Unmarshal(strings.NewReader(string(d)), &v); err != nil {
		return common.InternalServerError(common.KvSrv, p.GetTrackValue, err, "parsing error")
	}

	rsp.Value = &v
	return nil
}

func (p *KvService) TagsCloud(ctx context.Context, req *kv_proto.TagsCloudRequest, rsp *kv_proto.TagsCloudResponse) error {
	log.Info("Received Kv.TagsCloud request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	key := fmt.Sprintf("%v:%v:tags", req.OrgId, req.Object)
	for _, tag := range req.Tags {
		subKey := fmt.Sprintf("%v:%v", key, tag)
		_, err := p.Client.TxPipelined(func(pipe redis.Pipeliner) error {
			pipe.SAdd(key, tag)
			pipe.Incr(subKey)
			return nil
		})
		if err != nil {
			return nil
		}

		rkey := fmt.Sprintf("%v:range", key)
		if count, err := p.Client.Get(subKey).Float64(); err == nil {
			p.Client.ZAdd(rkey, redis.Z{count, tag})
		}
	}
	return nil
}

func (p *KvService) GetTopTags(ctx context.Context, req *kv_proto.GetTopTagsRequest, rsp *kv_proto.GetTopTagsResponse) error {
	log.Info("Received Kv.GetTopTags request")
	cmd := redis.NewStringCmd("SELECT", req.Index)
	p.Client.Process(cmd)

	rkey := fmt.Sprintf("%v:%v:tags:range", req.OrgId, req.Object)
	tags := p.Client.ZRevRange(rkey, 0, req.N).Val()
	rsp.Tags = tags
	return nil
}

// func (p *KvService) TxPipelined(fn func(redis.Pipeliner) error) error {
// 	log.Info("Received Kv.TxPipelined request")
// 	_, err := p.Client.TxPipelined(fn)
// 	if err != nil {
// 		return errors.InternalServerError("go.micro.srv.kv.TxPipelined", err.Error())
// 	}
// 	return nil
// }
