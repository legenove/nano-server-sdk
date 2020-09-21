package docs

import (
	"path"

	"github.com/legenove/nano-server-sdk/gincore"
	"github.com/legenove/nano-server-sdk/gincore/docs/views"
	"github.com/legenove/nano-server-sdk/servers"
	"github.com/legenove/utils"
)

func init() {
	if servers.Server.Doc {
		router := gincore.GetRouter()
		if utils.PathExists(path.Join(servers.Server.DocDir, "templates/")) {
			temp_path := path.Join(servers.Server.DocDir, "templates/")
			router.LoadHTMLGlob(temp_path + "/*")
		}
		if utils.PathExists(path.Join(servers.Server.DocDir, "static")) {
			static_path := path.Join(servers.Server.DocDir, "static")
			router.Static("/static", static_path)
		}
		// init routers
		router.GET("doc", views.TemplateDocApi)
		router.GET("doc/filedoc", views.TemplateDocFileApi)
		router.GET("doc/swagger/:name", views.TemplateDocSwaggerApi)
		router.GET("doc/proto/:name", views.TemplateDocProtoApi)
	}

}
