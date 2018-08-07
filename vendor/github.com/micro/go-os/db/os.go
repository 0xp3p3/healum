package db

import (
	"fmt"

	db "github.com/micro/db-srv/proto/db"
	"github.com/micro/go-micro/client"

	"golang.org/x/net/context"
)

type os struct {
	opts Options
	c    db.DBClient
}

func newOS(opts ...Option) DB {
	var options Options
	for _, o := range opts {
		o(&options)
	}

	if options.Client == nil {
		options.Client = client.DefaultClient
	}

	if len(options.Database) == 0 {
		options.Database = DefaultDatabase
	}

	if len(options.Table) == 0 {
		options.Table = DefaultTable
	}

	return &os{
		opts: options,
		c:    db.NewDBClient("go.micro.srv.db", options.Client),
	}
}

func protoToRecord(r *db.Record) Record {
	if r == nil {
		return nil
	}

	metadata := map[string]interface{}{}

	for k, v := range r.Metadata {
		metadata[k] = v
	}

	return &record{
		id:       r.Id,
		created:  r.Created,
		updated:  r.Updated,
		metadata: metadata,
		bytes:    []byte(r.Bytes),
	}
}

func recordToProto(r Record) *db.Record {
	if r == nil {
		return nil
	}

	md := map[string]string{}

	for k, v := range r.Metadata() {
		md[k] = fmt.Sprintf("%v", v)
	}

	return &db.Record{
		Id:       r.Id(),
		Created:  r.Created(),
		Updated:  r.Updated(),
		Metadata: md,
		Bytes:    string(r.Bytes()),
	}
}

func (o *os) Close() error {
	return nil
}

func (o *os) Init(opts ...Option) error {
	// No reinits
	return nil
}

func (o *os) Options() Options {
	return o.opts
}

func (o *os) Read(id string) (Record, error) {
	rsp, err := o.c.Read(context.TODO(), &db.ReadRequest{
		Database: &db.Database{
			Name:  o.opts.Database,
			Table: o.opts.Table,
		},
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return protoToRecord(rsp.Record), nil
}

func (o *os) Create(r Record) error {
	_, err := o.c.Create(context.TODO(), &db.CreateRequest{
		Database: &db.Database{
			Name:  o.opts.Database,
			Table: o.opts.Table,
		},
		Record: recordToProto(r),
	})
	return err
}

func (o *os) Update(r Record) error {
	_, err := o.c.Update(context.TODO(), &db.UpdateRequest{
		Database: &db.Database{
			Name:  o.opts.Database,
			Table: o.opts.Table,
		},
		Record: recordToProto(r),
	})
	return err
}

func (o *os) Delete(id string) error {
	_, err := o.c.Delete(context.TODO(), &db.DeleteRequest{
		Database: &db.Database{
			Name:  o.opts.Database,
			Table: o.opts.Table,
		},
		Id: id,
	})
	return err
}

func (o *os) Search(opts ...SearchOption) ([]Record, error) {
	options := SearchOptions{
		Limit:  10,
		Offset: 0,
	}

	for _, o := range opts {
		o(&options)
	}

	metadata := map[string]string{}
	for k, v := range options.Metadata {
		metadata[k] = fmt.Sprintf("%v", v)
	}

	rsp, err := o.c.Search(context.TODO(), &db.SearchRequest{
		Database: &db.Database{
			Name:  o.opts.Database,
			Table: o.opts.Table,
		},
		Metadata: metadata,
		Limit:    options.Limit,
		Offset:   options.Offset,
	})
	if err != nil {
		return nil, err
	}

	var records []Record

	for _, r := range rsp.Records {
		records = append(records, protoToRecord(r))
	}

	return records, nil
}

func (o *os) String() string {
	return "os"
}
