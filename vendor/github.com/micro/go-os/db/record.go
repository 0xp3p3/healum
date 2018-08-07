package db

import (
	"encoding/json"
	"time"
)

type record struct {
	id       string
	created  int64
	updated  int64
	metadata Metadata
	bytes    []byte
}

func newRecord(id string, md Metadata, data interface{}) Record {
	b, _ := json.Marshal(data)

	return &record{
		id:       id,
		metadata: md,
		created:  time.Now().Unix(),
		bytes:    b,
	}
}

func (r *record) Id() string {
	return r.id
}

func (r *record) Created() int64 {
	return r.created
}

func (r *record) Updated() int64 {
	return r.updated
}

func (r *record) Metadata() Metadata {
	return r.metadata
}

func (r *record) Bytes() []byte {
	return r.bytes
}

func (r *record) Scan(v interface{}) error {
	return json.Unmarshal(r.bytes, v)
}
