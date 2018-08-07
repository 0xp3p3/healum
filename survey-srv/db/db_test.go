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
		INSERT %v 
		IN %v`, `{"name":"hello\'s test"}`, common.DbSurveyTable)
	_, err := runQuery(ctx, q, common.DbSurveyTable)
	if err != nil {
		t.Error(err)
	}
}
