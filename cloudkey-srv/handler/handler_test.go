package handler

import (
	"fmt"
	"log"
	"server/cloudkey-srv/db"
	"server/cloudkey-srv/storage"
	_ "server/cloudkey-srv/storage/gce"
	"server/common"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudkms/v1"
)

func initDb() {
	cl := client.NewClient(
		client.Transport(nats_transport.NewTransport()),
		client.Broker(nats_broker.NewBroker()),
		client.RequestTimeout(3*time.Second),
		client.Retries(10),
	)
	ctx := common.NewTestContext(context.TODO())
	// db.DbCloudKeyName = common.TestingName("cloudkey")
	// db.DbCloudKeyTable = common.TestingName("cloudkey")

	db.RemoveDb(ctx, cl)
	db.Init(cl)
	storage.Init()
}

func TestKeyCreate(t *testing.T) {
	cl, err := google.DefaultClient(context.Background(), cloudkms.CloudPlatformScope)
	if err != nil {
		log.Fatalf("Unable to get default client: %v", err)
		return
	}
	cloudkmsService, err := cloudkms.New(cl)
	parent := "projects/locations/global/"
	fmt.Println(cloudkmsService.Projects.Locations.KeyRings.List(parent))
}
