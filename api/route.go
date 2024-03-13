package api

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"sub2clash/api/handler"
	"sub2clash/config"
	"sub2clash/middleware"

	"github.com/gin-gonic/gin"
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
			handler.SubmodHandler(c)
		},
	)
	r.GET(
		"/meta", func(c *gin.Context) {
			handler.SubHandler(c)
		},
	)
	r.GET(
		"/s/:hash", func(c *gin.Context) {
			handler.ShortLinkGetConfigHandler(c)
		},
	)
	r.GET(
		"/short", func(c *gin.Context) {
			handler.ShortLinkGetUrlHandler(c)
		})
	r.POST(
		"/short", func(c *gin.Context) {
			handler.ShortLinkGenHandler(c)
		},
	)
	r.PUT(
		"/short", func(c *gin.Context) {
			handler.ShortLinkUpdateHandler(c)
		},
	)
}
