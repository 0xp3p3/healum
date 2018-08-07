package api

import (
	"github.com/emicklei/go-restful"
	"log"
	"time"
)

// Returns current server unix timestamp ("since" token for events)
type TimestampService struct {
}

func (r TimestampService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/timestamp")
	ws.Route(ws.GET("/").To(r.Timestamp).
		Doc("Returns current server unix timestamp (since token for events)"))

	restful.Add(ws)
}

// Timestamp handler
func (r *TimestampService) Timestamp(req *restful.Request, rsp *restful.Response) {
	log.Print("Received Server.Timestamp API request")

	rsp.WriteEntity(map[string]int64{
		"timestamp": time.Now().Unix(),
	})
}
