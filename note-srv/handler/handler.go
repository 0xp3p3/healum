package handler

import (
	"context"
	"server/common"
	"server/note-srv/db"
	note_proto "server/note-srv/proto/note"

	log "github.com/sirupsen/logrus"
)

type NoteService struct{}

func (p *NoteService) All(ctx context.Context, req *note_proto.AllRequest, rsp *note_proto.AllResponse) error {
	log.Info("Received Note.All request")
	notes, err := db.All(ctx, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(notes) == 0 || err != nil {
		return common.NotFound(common.NoteSrv, p.All, err, "note not found")
	}
	rsp.Data = &note_proto.ArrData{notes}
	return nil
}

func (p *NoteService) Create(ctx context.Context, req *note_proto.CreateRequest, rsp *note_proto.CreateResponse) error {
	log.Info("Received Note.Create request")
	if len(req.Note.Title) == 0 {
		return common.BadRequest(common.NoteSrv, p.Create, nil, "note title empty")
	}
	if req.Note.Creator == nil {
		return common.BadRequest(common.NoteSrv, p.Create, nil, "note creator empty")
	}

	err := db.Create(ctx, req.Note)
	if err != nil {
		return common.InternalServerError(common.NoteSrv, p.Create, err, "create error")
	}
	rsp.Data = &note_proto.Data{req.Note}
	return nil
}

func (p *NoteService) Read(ctx context.Context, req *note_proto.ReadRequest, rsp *note_proto.ReadResponse) error {
	log.Info("Received Note.Read request")
	note, err := db.Read(ctx, req.Id, req.OrgId, req.TeamId)
	if note == nil || err != nil {
		return common.NotFound(common.NoteSrv, p.Read, err, "note not found")
	}
	rsp.Data = &note_proto.Data{note}
	return nil
}

func (p *NoteService) Delete(ctx context.Context, req *note_proto.DeleteRequest, rsp *note_proto.DeleteResponse) error {
	log.Info("Received Note.Delete request")
	if err := db.Delete(ctx, req.Id, req.OrgId, req.TeamId); err != nil {
		return common.InternalServerError(common.NoteSrv, p.Delete, err, "delete error")
	}
	return nil
}

func (p *NoteService) Search(ctx context.Context, req *note_proto.SearchRequest, rsp *note_proto.SearchResponse) error {
	log.Info("Received Note.Search request")
	notes, err := db.Search(ctx, req.Name, req.OrgId, req.TeamId, req.Limit, req.Offset, req.From, req.To, req.SortParameter, req.SortDirection)
	if len(notes) == 0 || err != nil {
		return common.NotFound(common.NoteSrv, p.Search, err, "note not found")
	}
	rsp.Data = &note_proto.ArrData{notes}
	return nil
}

func (p *NoteService) ByCreator(ctx context.Context, req *note_proto.ByCreatorRequest, rsp *note_proto.ByCreatorResponse) error {
	log.Info("Received Note.ByCreator request")
	notes, err := db.ByCreator(ctx, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(notes) == 0 || err != nil {
		return common.NotFound(common.NoteSrv, p.ByCreator, err, "note not found")
	}
	rsp.Data = &note_proto.ArrData{notes}
	return nil
}

func (p *NoteService) ByUser(ctx context.Context, req *note_proto.ByUserRequest, rsp *note_proto.ByUserResponse) error {
	log.Info("Received Note.ByUser request")
	notes, err := db.ByUser(ctx, req.UserId, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(notes) == 0 || err != nil {
		return common.NotFound(common.NoteSrv, p.ByUser, err, "note not found")
	}
	rsp.Data = &note_proto.ArrData{notes}
	return nil
}

func (p *NoteService) Filter(ctx context.Context, req *note_proto.FilterRequest, rsp *note_proto.FilterResponse) error {
	log.Info("Received Note.Filter request")
	notes, err := db.Filter(ctx, req.Category, req.Tags, req.OrgId, req.TeamId, req.Offset, req.Limit, req.SortParameter, req.SortDirection)
	if len(notes) == 0 || err != nil {
		return common.NotFound(common.NoteSrv, p.Filter, err, "note not found")
	}
	rsp.Data = &note_proto.ArrData{notes}
	return nil
}

func (p *NoteService) Update(ctx context.Context, req *note_proto.UpdateRequest, rsp *note_proto.UpdateResponse) error {
	return nil
}
