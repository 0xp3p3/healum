package db

import (
	"errors"
	"log"
	db_proto "server/db-srv/proto/db"

	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type clientWrapper struct {
	Db_client db_proto.DBClient
}

var (
	ClientWrapper   *clientWrapper
	DbCloudKeyName  = "cloudkey"
	DbCloudKeyTable = "cloudkey"

	ErrNotFound = errors.New("not found")
)

// DbMapping
//

// Storage for a db microservice client
func NewClientWrapper(serviceClient client.Client) *clientWrapper {
	cl := db_proto.NewDBClient("",
		serviceClient)

	return &clientWrapper{
		Db_client: cl,
	}
}

func Init(serviceClient client.Client) {
	ClientWrapper = NewClientWrapper(serviceClient)
	_, err := ClientWrapper.Db_client.CreateDatabase(context.TODO(), &db_proto.CreateDatabaseRequest{
		&db_proto.Database{
			Name:  DbCloudKeyName,
			Table: DbCloudKeyTable,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}

func RemoveDb(ctx context.Context, serviceClient client.Client) {
	ClientWrapper = NewClientWrapper(serviceClient)
	_, err := ClientWrapper.Db_client.DeleteDatabase(ctx, &db_proto.DeleteDatabaseRequest{
		&db_proto.Database{
			Name:  DbCloudKeyName,
			Table: DbCloudKeyTable,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

}
