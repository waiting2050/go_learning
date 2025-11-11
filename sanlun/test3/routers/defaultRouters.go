package routers

import (
	"test3/controllers/default_"

	"github.com/gin-gonic/gin"
)

func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", default_.DefaultController{}.Index)
		defaultRouters.GET("/news", default_.DefaultController{}.News)
	}
}
