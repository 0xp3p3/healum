package handler

import (
	"bytes"
	"context"
	account_proto "server/account-srv/proto/account"
	"server/common"
	kv_proto "server/kv-srv/proto/kv"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

var item = &kv_proto.Item{
	Key:        "key",
	Value:      []byte("hello world!"),
	Expiration: 0,
}

var account_email = &account_proto.Account{
	Email:    "email9@email.com",
	Password: "pass1",
}

func initHandler() *KvService {
	kvService := &KvService{
		Client: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
	return kvService
}

func TestPut(t *testing.T) {
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	// set key, value
	req_put := &kv_proto.PutRequest{item}
	rsp_put := &kv_proto.PutResponse{}
	err := hdlr.Put(ctx, req_put, rsp_put)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGet(t *testing.T) {
	hdlr := initHandler()

	hdlr.Client.Set("hi", "hello", 0)
	t.Log(hdlr.Client.Get("hi").Val())
}

func TestUnLock(t *testing.T) {
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	req_lock := &kv_proto.LockedRequest{Index: common.ACCOUNT_LOCKED_INDEX, AccountId: "test_id"}
	rsp_lock := &kv_proto.LockedResponse{}
	err := hdlr.Locked(ctx, req_lock, rsp_lock)
	if err != nil {
		t.Error("locked fail")
		return
	}

	// get ex with account id
	req_get := &kv_proto.IsLockedRequest{Index: common.ACCOUNT_LOCKED_INDEX, AccountId: "test_id"}
	rsp_get := &kv_proto.IsLockedResponse{}
	if err := hdlr.IsLocked(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if !rsp_get.Locked {
		t.Error("Value does not matched")
		return
	}

	req_unlock := &kv_proto.UnLockRequest{
		LockIndex: common.ACCOUNT_LOCKED_INDEX,
		AccountId: "test_id",
	}
	rsp_unlock := &kv_proto.UnLockResponse{}
	time.Sleep(time.Second)

	err = hdlr.UnLock(ctx, req_unlock, rsp_unlock)
	if err != nil {
		t.Error("unlocked fail")
		return
	}

	// get ex with account id
	if err := hdlr.IsLocked(ctx, req_get, rsp_get); err != nil {
		t.Error(err)
		return
	}

	if rsp_get.Locked {
		t.Error("Value does not matched")
		return
	}
}

func TestExpiry(t *testing.T) {
	hdlr := initHandler()

	// 3 seconds expire
	hdlr.Client.Set("3s", "hi, 3 secs", 3*time.Second)
	time.Sleep(time.Second)
	val, err := hdlr.Client.Get("3s").Bytes()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(string(val))
	// 3 seconds expire checking
	time.Sleep(5 * time.Second)
	val, err = hdlr.Client.Get("3s").Bytes()
	if err != nil {
		// t.Error(err.Error())
		t.Log("3 secs expire is successed")
	}
	// t.Log(string(val))

	// 3 mins expire checking
	hdlr.Client.Set("3m", "hi, 3 mins", 3*time.Minute)
	time.Sleep(time.Second)
	val, err = hdlr.Client.Get("3m").Bytes()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(string(val))

	// 3 hrs expire checking
	hdlr.Client.Set("3h", "hi, 3 hours", 3*time.Hour)
	time.Sleep(time.Second)
	val, err = hdlr.Client.Get("3h").Bytes()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(string(val))
}

func TestFull(t *testing.T) {
	hdlr := initHandler()
	ctx := common.NewTestContext(context.TODO())

	// set key, value
	req_put := &kv_proto.PutRequest{item}
	rsp_put := &kv_proto.PutResponse{}
	err := hdlr.Put(ctx, req_put, rsp_put)
	if err != nil {
		t.Error(err)
		return
	}

	// get key, value
	req_get := &kv_proto.GetRequest{item.Key}
	rsp_get := &kv_proto.GetResponse{}
	err = hdlr.Get(ctx, req_get, rsp_get)
	if err != nil {
		t.Error(err)
		return
	}

	if !bytes.Equal(rsp_get.Item.Value, item.Value) {
		t.Error("Value does not matched")
		return
	}

	// delete key
	req_del := &kv_proto.DelRequest{item.Key}
	rsp_del := &kv_proto.DelResponse{}
	err = hdlr.Del(ctx, req_del, rsp_del)
	if err != nil {
		t.Error(err)
		return
	}

	// check key
	err = hdlr.Get(ctx, req_get, rsp_get)
	if err == nil {
		t.Error(err)
		return
	}
}

func TestTxPipeline(t *testing.T) {
	hdlr := initHandler()

	_, err := hdlr.Client.TxPipelined(func(pipe redis.Pipeliner) error {
		pipe.SAdd("keyy", "valuee")
		pipe.Incr("keyy:valuee")
		return nil
	})

	if err != nil {
		t.Error(err)
		return
	}

	_, err = hdlr.Client.TxPipelined(func(pipe redis.Pipeliner) error {
		pipe.Get("keyy:valuee")
		pipe.ZAdd("keyy", redis.Z{3, "valuee"})
		return nil
	})
}
