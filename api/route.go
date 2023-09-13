package api

import (
	"github.com/gin-gonic/gin"
	"sub2clash/api/controller"
	"sub2clash/middleware"
)

func SetRoute(r *gin.Engine) {
	r.Use(middleware.ZapLogger())
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
}
