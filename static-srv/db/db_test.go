package db

import (
	"context"
	"fmt"
	"server/common"
	"testing"
	"time"

	"github.com/micro/go-micro/client"
	nats_broker "github.com/micro/go-plugins/broker/nats"
	nats_transport "github.com/micro/go-plugins/transport/nats"
)

var cl = client.NewClient(
	client.Transport(nats_transport.NewTransport()),
	client.Broker(nats_broker.NewBroker()),
	client.RequestTimeout(4*time.Second),
	client.Retries(5),
)

func TestRunQuery(t *testing.T) {
	Init(cl)
	ctx := common.NewTestContext(context.TODO())

	q := fmt.Sprintf(`
		UPSERT { _key: "f688666e-44b4-11e8-bcee-20c9d0453b15" } 
		INSERT {"_key":"f688666e-44b4-11e8-bcee-20c9d0453b15","created":1524240465,"data":{"apps":[{"created":"0","description":"description","iconSlug":"iconslug","id":"111","image":"images","name":"title","platforms":[{"created":"0","id":"111","name":"title","updated":"0","url":"url"}],"summary":"summary","tags":["tag1","tag2"],"updated":"0"}],"created":"1524240465","description":"description","devices":[{"created":"0","description":"description","iconSlug":"iconslug","id":"111","image":"images","name":"title","summary":"summary","tags":["tag1","tag2"],"updated":"0","url":""}],"iconSlug":"","id":"f688666e-44b4-11e8-bcee-20c9d0453b15","name":"title","org_id":"orgid","summary":"summary","trackerMethods":[{"created":"0","iconSlug":"iconSlug","id":"111","name":"title","nameSlug":"nameSlug","updated":"0"}],"unit":["kgs","stones","%%"],"updated":"1524240465","wearables":[{"created":"0","description":"description","iconSlug":"","id":"111","image":"images","name":"title","summary":"summary","tags":[],"updated":"0","url":""}]},"id":"f688666e-44b4-11e8-bcee-20c9d0453b15","name":"title","parameter1":"orgid","updated":1524240465} 
		UPDATE {"_key":"f688666e-44b4-11e8-bcee-20c9d0453b15","created":1524240465,"data":{"apps":[{"created":"0","description":"description","iconSlug":"iconslug","id":"111","image":"images","name":"title","platforms":[{"created":"0","id":"111","name":"title","updated":"0","url":"url"}],"summary":"summary","tags":["tag1","tag2"],"updated":"0"}],"created":"1524240465","description":"description","devices":[{"created":"0","description":"description","iconSlug":"iconslug","id":"111","image":"images","name":"title","summary":"summary","tags":["tag1","tag2"],"updated":"0","url":""}],"iconSlug":"","id":"f688666e-44b4-11e8-bcee-20c9d0453b15","name":"title","org_id":"orgid","summary":"summary","trackerMethods":[{"created":"0","iconSlug":"iconSlug","id":"111","name":"title","nameSlug":"nameSlug","updated":"0"}],"unit":["kgs","stones","%%"],"updated":"1524240465","wearables":[{"created":"0","description":"description","iconSlug":"","id":"111","image":"images","name":"title","summary":"summary","tags":[],"updated":"0","url":""}]},"id":"f688666e-44b4-11e8-bcee-20c9d0453b15","name":"title","parameter1":"orgid","updated":1524240465} 
		IN marker`)
	_, err := runQuery(ctx, q, common.DbMarkerTable)
	if err != nil {
		t.Error(err)
	}
}
