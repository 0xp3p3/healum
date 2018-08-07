package backend

import (
	"context"
	"time"

	"server/activity-srv/db"
	"server/activity-srv/handler"
	activity_proto "server/activity-srv/proto/activity"
	content_proto "server/content-srv/proto/content"

	log "github.com/sirupsen/logrus"
)

// Backends with initialized credentials and weights
var (
	Backends map[string]Backend
	Weights  map[string]int64
)

// Backend is
type Backend interface {
	Query(source string) ([]*content_proto.Activity, error)
}

// Init initialize backend
func Init() {
	Backends = make(map[string]Backend)
	Weights = make(map[string]int64)
}

// QueryExternalAPI queries to external api
func QueryExternalAPI(service *handler.ActivityService, ext *ExtBackend, k string, done chan bool) {
	// go func() {
	// getting database from backend till finish fetching
	for {
		datas, err := ext.Query(k)
		if err != nil {
			log.WithField("err", err).Error("Finsihed external api because of response error")
			break
		}

		// after fetch from backend database and call create-srv handler to create activity-content
		for _, data := range datas {
			service.Create(context.TODO(), &activity_proto.CreateRequest{data}, &activity_proto.CreateResponse{})
		}

		if ext.Finished {
			log.Info("Finsihed external api:", k)
			ext.Finished = false
			break
		}
	}
	// save config with latest next url
	db.CreateConfig(context.Background(), &activity_proto.Config{
		Id:           ext.ID,
		Name:         ext.Name,
		AppURL:       ext.AppURL,
		Next:         ext.Next,
		Weight:       ext.Weight,
		TimeInterval: ext.TimeInterval,
		Enabled:      ext.Enabled,
		Transform:    ext.Transform,
	})
	// }()

	<-done
}

// FetchDatabase stores all database after fetch all datas from external apis
func FetchDatabase(service *handler.ActivityService) {
	// Query all the backends
	for key, val := range Backends {
		k := key
		v := val
		go func() {
			log.Info("Started fetch external api", k)
			for {
				ext := v.(*ExtBackend)
				if !ext.Enabled {
					continue
				}
				tickerChan := time.NewTicker(time.Minute * time.Duration(ext.TimeInterval)).C
				done := make(chan bool)
				go QueryExternalAPI(service, ext, k, done)

				select {
				case <-tickerChan:
					log.Info("Restarted fetch external api:", k)
					done <- true
					// remove old database in the collection
					// db.DeleteSource(context.Background(), ext.Name)
				}
			}
		}()
	}
}
