package handler

import (
	"context"
	"server/common"
	"server/note-srv/db"
	note_proto "server/note-srv/proto/note"
	user_proto "server/user-srv/proto/user"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

var note = &note_proto.Note{
	Title:       "note1",
	OrgId:       "orgid",
	Description: "description1",
	Creator:     &user_proto.User{Id: "userid"},
	User:        &user_proto.User{Id: "userid"},
	Category:    "category1",
	Tags:        []string{"a", "b", "c"},
}

func initDb() {
	cl := client.NewClient(
		client.Transport(nats_transport.NewTransport()),
		client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(4*time.Second),
		client.Retries(5),
	)
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func createNote(ctx context.Context, hdlr *NoteService, t *testing.T) *note_proto.Note {
	req := &note_proto.CreateRequest{Note: note}
	resp := &note_proto.CreateResponse{}
	if err := hdlr.Create(ctx, req, resp); err != nil {
		t.Error(err)
		return nil
	}
	return resp.Data.Note
}

func TestAll(t *testing.T) {
	initDb()
	hdlr := new(NoteService)
	ctx := common.NewTestContext(context.TODO())
	note := createNote(ctx, hdlr, t)
	if note == nil {
		return
	}

	req_all := &note_proto.AllRequest{}
	resp_all := &note_proto.AllResponse{}
	time.Sleep(2 * time.Second)
	err := hdlr.All(ctx, req_all, resp_all)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_all.Data.Notes) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_all.Data.Notes[0].Id != note.Id {
		t.Error("Id does not match")
		return
	}
}
func TestNoteIsCreated(t *testing.T) {
	initDb()
	hdlr := new(NoteService)
	ctx := common.NewTestContext(context.TODO())
	note := createNote(ctx, hdlr, t)
	if note == nil {
		t.Error("Create is failed")
		return
	}
}

func TestNoteRead(t *testing.T) {
	initDb()
	hdlr := new(NoteService)
	ctx := common.NewTestContext(context.TODO())
	note := createNote(ctx, hdlr, t)
	if note == nil {
		return
	}

	req_read := &note_proto.ReadRequest{Id: note.Id}
	resp_read := &note_proto.ReadResponse{}
	if err := hdlr.Read(ctx, req_read, resp_read); err != nil {
		t.Error(err)
		return
	}
	if resp_read.Data.Note == nil {
		t.Error("Object could not be nil")
		return
	}
	if resp_read.Data.Note.Id != note.Id {
		t.Error("Id does not match")
		return
	}
}

func TestNoteDelete(t *testing.T) {
	initDb()
	hdlr := new(NoteService)
	ctx := common.NewTestContext(context.TODO())
	note := createNote(ctx, hdlr, t)
	if note == nil {
		return
	}

	req_del := &note_proto.DeleteRequest{Id: note.Id}
	resp_del := &note_proto.DeleteResponse{}
	if err := hdlr.Delete(ctx, req_del, resp_del); err != nil {
		t.Error(err)
		return
	}
}

func TestByCreator(t *testing.T) {
	initDb()
	hdlr := new(NoteService)
	ctx := common.NewTestContext(context.TODO())
	note := createNote(ctx, hdlr, t)
	if note == nil {
		return
	}

	req_creator := &note_proto.ByCreatorRequest{UserId: "userid"}
	resp_creator := &note_proto.ByCreatorResponse{}
	if err := hdlr.ByCreator(ctx, req_creator, resp_creator); err != nil {
		t.Error(err)
		return
	}
	if len(resp_creator.Data.Notes) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_creator.Data.Notes[0].Id != note.Id {
		t.Error("Id does not match")
		return
	}
}

func TestByUser(t *testing.T) {
	initDb()
	hdlr := new(NoteService)
	ctx := common.NewTestContext(context.TODO())
	note := createNote(ctx, hdlr, t)
	if note == nil {
		return
	}

	req_user := &note_proto.ByUserRequest{UserId: note.Creator.Id}
	resp_user := &note_proto.ByUserResponse{}
	err := hdlr.ByUser(ctx, req_user, resp_user)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_user.Data.Notes) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_user.Data.Notes[0].Id != note.Id {
		t.Error("Id does not match")
		return
	}
}

func TestFilter(t *testing.T) {
	initDb()
	hdlr := new(NoteService)
	ctx := common.NewTestContext(context.TODO())
	note := createNote(ctx, hdlr, t)
	if note == nil {
		return
	}

	req_filter := &note_proto.FilterRequest{
		Category: []string{},
		Tags:     []string{"a", "d"},
	}
	resp_filter := &note_proto.FilterResponse{}
	err := hdlr.Filter(ctx, req_filter, resp_filter)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_filter.Data.Notes) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_filter.Data.Notes[0].Id != note.Id {
		t.Error("Id does not match")
		return
	}
}

func TestSearch(t *testing.T) {
	initDb()
	hdlr := new(NoteService)
	ctx := common.NewTestContext(context.TODO())
	note := createNote(ctx, hdlr, t)
	if note == nil {
		return
	}

	req_search := &note_proto.SearchRequest{
		Name:   "note1",
		OrgId:  "orgid",
		Offset: 0,
		Limit:  10,
	}
	resp_search := &note_proto.SearchResponse{}
	time.Sleep(2 * time.Second)
	err := hdlr.Search(ctx, req_search, resp_search)
	if err != nil {
		t.Error(err)
		return
	}
	if len(resp_search.Data.Notes) == 0 {
		t.Error("Count does not match")
		return
	}
	if resp_search.Data.Notes[0].Id != note.Id {
		t.Error("Id does not match")
		return
	}
}
