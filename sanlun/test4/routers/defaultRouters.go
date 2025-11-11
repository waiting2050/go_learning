package routers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func InitMiddleWare1(c *gin.Context) {
	start := time.Now().UnixNano()
	fmt.Println("1-我是一个中间件111")

	// 调用该请求的剩余程序
	c.Next()

	fmt.Println("2-我是一个中间件111")
	end := time.Now().UnixNano()

	fmt.Println(end - start)
}

func InitMiddleWare2(c *gin.Context) {
	start := time.Now().UnixNano()
	fmt.Println("1-我是一个中间件222")

	// 调用该请求的剩余程序
	c.Next()

	fmt.Println("2-我是一个中间件222")
	end := time.Now().UnixNano()

	fmt.Println(end - start)
}

func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", func(c *gin.Context) {
			fmt.Println("这是一个首页")
			c.HTML(http.StatusOK, "default/index.html", gin.H{
				"msg": "我是一个msg",
			})
		})
		defaultRouters.GET("/news", func(c *gin.Context) {
			c.String(200, "新闻")
		})
	}
}
