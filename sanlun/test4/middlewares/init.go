package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func InitMiddleWare(c *gin.Context) {
	// 判断用户是否登录
	fmt.Println(time.Now())
	fmt.Println(c.Request.URL)

	c.Set("username", "张三")

	// 定义一个goroutine统计日志
	go func() {
		cCp := c.Copy()
		time.Sleep(5 * time.Second)
		fmt.Println("Done! in path " + cCp.Request.URL.Path)
	}()
}
