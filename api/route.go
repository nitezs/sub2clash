package api

import (
	"embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"sub2clash/api/controller"
	"sub2clash/config"
	"sub2clash/middleware"
)

//go:embed templates/*
var templates embed.FS

func SetRoute(r *gin.Engine) {
	r.Use(middleware.ZapLogger())
	// 使用内嵌的模板文件
	r.SetHTMLTemplate(template.Must(template.New("").ParseFS(templates, "templates/*")))
	r.GET(
		"/", func(c *gin.Context) {
			c.HTML(
				200, "index.html", gin.H{
					"Version": config.Version,
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
