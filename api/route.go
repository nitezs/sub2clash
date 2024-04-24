package api

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"sub2clash/api/handler"
	"sub2clash/constant"
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
			version := constant.Version
			if len(constant.Version) > 7 {
				version = constant.Version[:7]
			}
			c.HTML(
				200, "index.html", gin.H{
					"Version": version,
				},
			)
		},
	)
	r.GET("/clash", handler.SubmodHandler)
	r.GET("/meta", handler.SubHandler)
	r.GET("/s/:hash", handler.GetRawConfHandler)
	r.POST("/short", handler.GenerateLinkHandler)
	r.PUT("/short", handler.UpdateLinkHandler)
}
