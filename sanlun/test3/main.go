package main

import (
	"fmt"
	"test3/routers"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/golang/protobuf/ptypes/timestamp"
)

type Article struct {
	Title   string
	Content string
}

type UserInfo struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func UnixToTime(timestamp int) string {
	fmt.Println(timestamp)
	t := time.Unix(int64(timestamp), 0)

	return t.Format("2006-01-02 15:04:05")
}

func main() {
	r := gin.Default()

	// 自定义模板函数，要把这个函数放在加载模板前
	r.SetFuncMap(template.FuncMap{
		"UnixToTime": UnixToTime,
	})

	r.LoadHTMLGlob("templates/**/*.html")
	r.Static("/static", "./static")

	routers.AdminRoutersInit(r)
	routers.ApiRoutersInit(r)
	routers.DefaultRoutersInit(r)

	// // GET请求传值
	// r.GET("/", func(c *gin.Context) {
	// 	username := c.Query("username")
	// 	age := c.Query("age")
	// 	page := c.DefaultQuery("page", "1")

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"username": username,
	// 		"age":      age,
	// 		"page":     page,
	// 	})
	// })
	// r.GET("/article", func(c *gin.Context) {
	// 	id := c.DefaultQuery("id", "1")

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"msg": "新闻详情",
	// 		"id":  id,
	// 	})
	// })

	// // POST演示
	// r.GET("/user", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "default/user.html", gin.H{})
	// })
	// // 获取POST表单传过来的数据
	// r.POST("/doAddUser1", func(c *gin.Context) {
	// 	username := c.PostForm("username")
	// 	password := c.PostForm("password")
	// 	age := c.DefaultPostForm("age", "18")

	// 	c.JSON(http.StatusOK, gin.H{
	// 		"username": username,
	// 		"password": password,
	// 		"age":      age,
	// 	})
	// })

	// // 获取 GET POST传递的数据绑定到结构体
	// r.POST("/doAddUser2", func(c *gin.Context) {
	// 	user := &UserInfo{}
	// 	if err := c.ShouldBind(&user); err == nil {
	// 		c.JSON(http.StatusOK, user)
	// 	} else {
	// 		c.JSON(http.StatusOK, gin.H{
	// 			"err": err.Error(),
	// 		})
	// 	}
	// })

	// 前台
	// r.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "default/index.html", gin.H{
	// 		"title": "首页",
	// 		"score": 89,
	// 		"msg":   "msg",
	// 		"hobby": []string{"吃饭", "睡觉", "写代码"},
	// 		"newsList": []interface{}{
	// 			&Article{
	// 				Title:   "新闻标题111",
	// 				Content: "新闻详情111",
	// 			},
	// 			&Article{
	// 				Title:   "新闻标题222",
	// 				Content: "新闻详情222",
	// 			},
	// 		},
	// 		"testSlice": []string{},
	// 		"news": &Article{
	// 			Title:   "新闻标题",
	// 			Content: "新闻详情",
	// 		},
	// 		"date": 1629423555,
	// 	})
	// })
	// r.GET("/news", func(c *gin.Context) {
	// 	news := &Article{
	// 		Title:   "新闻标题",
	// 		Content: "新闻详情",
	// 	}
	// 	c.HTML(http.StatusOK, "default/news.html", gin.H{
	// 		"title": "新闻页面",
	// 		"news":  news,
	// 	})
	// })

	// // 动态路由传值
	// r.GET("/list/:cid", func(c *gin.Context) {
	// 	cid := c.Param("cid")
	// 	c.String(200, "%v", cid)
	// })

	// // 后台
	// r.GET("/admin", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "admin/index.html", gin.H{
	// 		"title": "后台首页",
	// 	})
	// })
	// r.GET("/admin/news", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "admin/news.html", gin.H{
	// 		"title": "新闻页面",
	// 	})
	// })

	r.Run()
}
