package api

import (
	"github.com/gin-gonic/gin"
	"sub/api/controller"
)

func SetRoute(r *gin.Engine) {
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
