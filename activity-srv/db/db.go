package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	activity_proto "server/activity-srv/proto/activity"
	"server/common"
	db_proto "server/db-srv/proto/db"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/micro/go-micro/client"
	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

type clientWrapper struct {
	Db_client db_proto.DBClient
}

var (
	ClientWrapper *clientWrapper
	ErrNotFound   = errors.New("not found")
)

// DbMapping
//
// Data
// Id - Id
// Parameter1 - meta string
// Parameter3 - source
// Item - response real data
//
// Config
// Id - Id
// Parameter1 - meta string
// Parameter3 - source
//

// Storage for a db microservice client
func NewClientWrapper(serviceClient client.Client) *clientWrapper {
	cl := db_proto.NewDBClient("", serviceClient)

	return &clientWrapper{
		Db_client: cl,
	}
}

// Init initializes healum databases
func Init(serviceClient client.Client) error {
	ClientWrapper = NewClientWrapper(serviceClient)
	// if _, err := ClientWrapper.Db_client.Init(context.TODO(), &db_proto.InitRequest{}); err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	return nil
}

// RemoveDb removes healum database (for testing)
func RemoveDb(ctx context.Context, serviceClient client.Client) error {
	ClientWrapper = NewClientWrapper(serviceClient)
	if _, err := ClientWrapper.Db_client.RemoveDb(ctx, &db_proto.RemoveDbRequest{}); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func configToRecord(config *activity_proto.Config) (string, error) {
	d := map[string]interface{}{
		"_key": config.Id,
		"id":   config.Id,
		"name": config.Name,
		"data": config,
	}
	body, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(body), err
}

func recordToConfig(r *db_proto.Record) (*activity_proto.Config, error) {
	var p activity_proto.Config
	if err := jsonpb.Unmarshal(strings.NewReader(r.Parameter3), &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func runQuery(ctx context.Context, q string, table string) (*db_proto.RunQueryResponse, error) {
	return ClientWrapper.Db_client.RunQuery(ctx, &db_proto.RunQueryRequest{
		Database: &db_proto.Database{
			Name:     common.DbHealumName,
			Table:    table,
			Driver:   common.DbHealumDriver,
			Metadata: common.SearchableMetaMap,
		},
		Query: q,
	})
}

// ListConfig queries configs
func ListConfig(ctx context.Context) ([]*activity_proto.Config, error) {
	configs := []*activity_proto.Config{}

	q := fmt.Sprintf(`
			FOR doc IN %v
			RETURN doc`, common.DbActivityConfigTable)
	resp, err := runQuery(ctx, q, common.DbActivityConfigTable)
	if err != nil {
		return nil, err
	}
	// parsing
	for _, r := range resp.Records {
		if config, err := recordToConfig(r); err == nil {
			configs = append(configs, config)
		}
	}
	return configs, err
}

// CreateConfig creates config
func CreateConfig(ctx context.Context, config *activity_proto.Config) error {
	if len(config.Id) == 0 {
		config.Id = uuid.NewUUID().String()
	}
	record, err := configToRecord(config)
	if err != nil {
		return err
	}
	if len(record) == 0 {
		return errors.New("server serialization")
	}

	q := fmt.Sprintf(`
		UPSERT { _key: "%v" }
		INSERT %v
		UPDATE %v
		IN %v`, config.Id, record, record, common.DbActivityConfigTable)
	_, err = runQuery(ctx, q, common.DbActivityConfigTable)
	return err
}

// ReadConfig reads config with id
func ReadConfig(ctx context.Context, id string) (*activity_proto.Config, error) {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
		FOR doc IN %v
		%s
		RETURN doc`, common.DbActivityConfigTable, query)

	resp, err := runQuery(ctx, q, common.DbActivityConfigTable)
	if err != nil || len(resp.Records) == 0 {
		return nil, err
	}

	config, err := recordToConfig(resp.Records[0])
	return config, nil
}

// DeleteConfig delets config with id
func DeleteConfig(ctx context.Context, id string) error {
	query := fmt.Sprintf(`FILTER doc.id == "%v"`, id)

	q := fmt.Sprintf(`
			FOR doc IN %v
			%s
			REMOVE doc IN %v`, common.DbActivityConfigTable, query, common.DbActivityConfigTable)
	_, err := runQuery(ctx, q, common.DbActivityConfigTable)
	return err
}
