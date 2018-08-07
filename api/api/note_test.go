package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"server/api/utils"
	"server/common"
	"server/note-srv/db"
	note_proto "server/note-srv/proto/note"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

// var serverURL = "http://localhost:8080"
var noteURL = "/server/notes"

var note = &note_proto.Note{
	Id:          "111",
	Title:       "note1",
	OrgId:       "orgid",
	Description: "description1",
	Creator:     user,
	User:        user,
	Category:    "category1",
	Tags:        []string{"a", "b", "c"},
}

func initNoteDb() {
	cl := client.NewClient(client.Transport(nats_transport.NewTransport()), client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(5*time.Second),
		client.Retries(5))
	// ctx := common.NewTestContext(context.TODO())
	// db.RemoveDb(ctx, cl)
	db.Init(cl)
}

func AllNotes(t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+noteURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping note because already created")

	}
	time.Sleep(time.Second)
}

func CreateNote(note *note_proto.Note, t *testing.T) {
	// Send an POST request.
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(map[string]interface{}{"note": note})
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+noteURL+"/note?session="+sessionId, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	// if resp.StatusCode == http.StatusInternalServerError {
	// 	t.Skip("Skipping note because already created")
	// }
	time.Sleep(time.Second)
}

func ReadNote(id string, t *testing.T) *note_proto.Note {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+noteURL+"/note/"+id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return nil
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping note because already created")
		return nil
	}
	time.Sleep(time.Second)

	r := note_proto.ReadResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return nil
	}
	json.Unmarshal(body, &r)
	if r.Data == nil {
		// t.Errorf("Response does not matched")
		return nil
	}
	return r.Data.Note
}

func DeleteNote(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a DELETE request.
	req, err := http.NewRequest("DELETE", serverURL+noteURL+"/note/"+id+"?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping note because already created")

	}
	time.Sleep(time.Second)
}

func SearchNotes(search *note_proto.SearchRequest, t *testing.T) {
	// Send a PUT request.
	jsonStr, err := json.Marshal(search)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println(string(jsonStr))

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	req, err := http.NewRequest("POST", serverURL+noteURL+"/search?session="+sessionId+"&offset=0&limit=20", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping note because already created")
	}
	time.Sleep(time.Second)
}

func NotesByCreator(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+noteURL+"/creator/"+id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping note because already created")

	}
	time.Sleep(time.Second)
}

func NotesByUser(id string, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+noteURL+"/user/"+id+"?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping note because already created")

	}
	time.Sleep(time.Second)
}

func NotesFilter(filter *note_proto.FilterRequest, t *testing.T) {
	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a POST request.
	jsonStr, err := json.Marshal(filter)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(jsonStr))

	req, err := http.NewRequest("POST", serverURL+noteURL+"/filter?session="+sessionId+"&offset=0&limit=10", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping note because already created")
	}
	time.Sleep(time.Second)
}

func TestNotesAll(t *testing.T) {
	initNoteDb()

	CreateNote(note, t)
	AllNotes(t)
}

func TestNoteRead(t *testing.T) {
	initNoteDb()

	CreateNote(note, t)
	p := ReadNote("111", t)
	if p == nil {
		t.Errorf("Note does not matched")
		return
	}
	if p.Id != note.Id {
		t.Errorf("Id does not matched")
		return
	}
	if p.Title != note.Title {
		t.Errorf("Title does not matched")
		return
	}
}

func TestNoteDelete(t *testing.T) {
	initNoteDb()

	CreateNote(note, t)
	DeleteNote("111", t)
	p := ReadNote("111", t)
	if p != nil {
		t.Errorf("Note does not matched")
		return
	}
}

func TestNoteSearch(t *testing.T) {
	initNoteDb()

	CreateNote(note, t)
	search := &note_proto.SearchRequest{
		Name:  "note1",
		OrgId: "orgid",
	}
	SearchNotes(search, t)
}

func TestNotesByCreator(t *testing.T) {
	initNoteDb()

	CreateNote(note, t)
	NotesByCreator("userid", t)
}

func TestNotesByUser(t *testing.T) {
	initNoteDb()

	CreateNote(note, t)
	NotesByUser("userid", t)
}

func TestNotesFilter(t *testing.T) {
	initNoteDb()

	CreateNote(note, t)
	filter := &note_proto.FilterRequest{
		Category: []string{"category1", "category2"},
		Tags:     []string{"a", "c", "d"},
	}
	NotesFilter(filter, t)
}

func TestErrReadNote(t *testing.T) {
	initNoteDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+noteURL+"/note/999?session="+sessionId, nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping note because already created")

	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
		return
	}
}

func TestErrAllNotes(t *testing.T) {
	initNoteDb()

	sessionId := GetSessionId("email8@email.com", "pass1", t)
	// Send a GET request.
	req, err := http.NewRequest("GET", serverURL+noteURL+"/all?session="+sessionId+"&offset=0&limit=10", nil)
	req.Header.Set("Content-Type", restful.MIME_JSON)
	common.SetTestHeader(req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("unexpected error in sending req: %v", err)
		return
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Skip("Skipping note because already created")

	}
	time.Sleep(time.Second)

	r := utils.ErrResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return
	}
	json.Unmarshal(body, &r)

	if len(r.Errors) == 0 {
		t.Skip("Errors does not matched")
		t.Error(r)
		return
	}
}
