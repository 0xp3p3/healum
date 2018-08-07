package api

import (
	"github.com/emicklei/go-restful"
	"log"
	"github.com/micro/go-os/metrics"
)

// Returns status of metrics
type MetricsService struct {
	FilterMiddle Filters
	ServerMetrics metrics.Metrics
}

func (r MetricsService) Register() {
	ws := new(restful.WebService)

	ws.Path("/server/metrics")
	ws.Route(ws.GET("/").To(r.Varz).
		Filter(r.FilterMiddle.BasicAuthenticate).
		Doc("Shows current values of metrics"))

	restful.Add(ws)
}

// Metrics handler
func (r *MetricsService) Varz(req *restful.Request, rsp *restful.Response) {
	log.Print("Received Server.Varz API request")

	rsp.WriteEntity(map[string]string{
		"metrics": r.ServerMetrics.String(),
	})
}
