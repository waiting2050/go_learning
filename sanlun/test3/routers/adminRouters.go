package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminRoutersInit(r *gin.Engine) {
	adminRouters := r.Group("/admin")
	{
		adminRouters.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "后台首页")
		})
		adminRouters.GET("/user", func(c *gin.Context) {
			c.String(http.StatusOK, "用户列表")
		})
		adminRouters.GET("/user/add", func(c *gin.Context) {
			c.String(http.StatusOK, "增加用户")
		})
		adminRouters.GET("/user/edit", func(c *gin.Context) {
			c.String(http.StatusOK, "修改用户")
		})
		adminRouters.GET("/article", func(c *gin.Context) {
			c.String(http.StatusOK, "新闻列表")
		})
	}
}
