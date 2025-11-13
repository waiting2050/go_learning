package routers

import (
	"fmt"
	"test4/middlewares"

	"github.com/gin-gonic/gin"
)

func AdminRoutersInit(r *gin.Engine) {
	adminRouters := r.Group("/admin", middlewares.InitMiddleWare)
	{
		adminRouters.GET("/", func(c *gin.Context) {
			username, _ := c.Get("username")
			fmt.Println(username)

			v, _ := username.(string)
			c.String(200, "后台首页" + v)
		})
		adminRouters.GET("/user", func(c *gin.Context) {
			c.String(200, "用户列表")
		})
		adminRouters.GET("/article", func(c *gin.Context) {
			c.String(200, "新闻列表")
		})
	}
}
