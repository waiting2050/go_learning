package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/golang/protobuf/ptypes/timestamp"
)

type Article struct {
	Title   string
	Content string
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

	// 前台
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "default/index.html", gin.H{
			"title": "首页",
			"score": 89,
			"msg":   "msg",
			"hobby": []string{"吃饭", "睡觉", "写代码"},
			"newsList": []interface{}{
				&Article{
					Title:   "新闻标题111",
					Content: "新闻详情111",
				},
				&Article{
					Title:   "新闻标题222",
					Content: "新闻详情222",
				},
			},
			"testSlice": []string{},
			"news": &Article{
				Title:   "新闻标题",
				Content: "新闻详情",
			},
			"date": 1629423555,
		})
	})
	r.GET("/news", func(c *gin.Context) {
		news := &Article{
			Title:   "新闻标题",
			Content: "新闻详情",
		}
		c.HTML(http.StatusOK, "default/news.html", gin.H{
			"title": "新闻页面",
			"news":  news,
		})
	})

	// 后台
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin/index.html", gin.H{
			"title": "后台首页",
		})
	})
	r.GET("/admin/news", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin/news.html", gin.H{
			"title": "新闻页面",
		})
	})

	r.Run()
}
