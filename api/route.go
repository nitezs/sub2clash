package api

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"sub2clash/api/controller"
	"sub2clash/config"
	"sub2clash/middleware"
)

//go:embed static
var staticFiles embed.FS

func SetRoute(r *gin.Engine) {
	r.Use(middleware.ZapLogger())

	// 使用内嵌的模板文件
	tpl, err := template.ParseFS(staticFiles, "static/*")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
	r.SetHTMLTemplate(tpl)

	r.GET(
		"/static/*filepath", func(c *gin.Context) {
			c.FileFromFS("static/"+c.Param("filepath"), http.FS(staticFiles))
		},
	)
	r.GET(
		"/", func(c *gin.Context) {
			version := config.Version
			if len(config.Version) > 7 {
				version = config.Version[:7]
			}
			c.HTML(
				200, "index.html", gin.H{
					"Version": version,
				},
			)
		},
	)
	r.GET(
		"/clash", func(c *gin.Context) {
			controller.SubmodHandler(c)
		},
	)
	r.GET(
		"/meta", func(c *gin.Context) {
			controller.SubHandler(c)
		},
	)
	r.POST(
		"/short", func(c *gin.Context) {
			controller.ShortLinkGenHandler(c)
		},
	)
	r.GET(
		"/s/:hash", func(c *gin.Context) {
			controller.ShortLinkGetHandler(c)
		},
	)
}
