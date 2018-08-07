package utils

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/micro/go-os/config"
)

// Swagger is used for API documentation
type SwaggerService struct {
	Config      config.Config
	SwaggerPath string
}

// Registers Swagger service
func (u SwaggerService) Register() {
	u.SwaggerPath = u.Config.Get("swagger", "swagger.path").String("")

	// accept and respond in JSON unless told otherwise
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	// gzip if accepted
	restful.DefaultContainer.EnableContentEncoding(true)
	// faster router
	restful.DefaultContainer.Router(restful.CurlyRouter{})
	// no need to access body more than once
	restful.SetCacheReadEntity(false)

	// API Cross-origin requests
	// apiCors := props.GetBool("http.server.cors", false)

	addr := u.Config.Get("swagger", "http.server.host").String("localhost") + ":" + u.Config.Get("swagger", "http.server.port").String("8080")
	basePath := "http://" + addr

	// Register Swagger UI
	swagger.InstallSwaggerService(swagger.Config{
		WebServices:     restful.RegisteredWebServices(),
		WebServicesUrl:  basePath,
		ApiPath:         "/apidocs.json",
		SwaggerPath:     u.SwaggerPath,
		SwaggerFilePath: u.Config.Get("swagger", "swagger.file.path").String(""),
	})

}
