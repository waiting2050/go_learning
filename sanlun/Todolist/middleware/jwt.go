package middleware

import (
	"Todolist/models"
	"Todolist/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func (c *gin.Context)  {
		var code int
		code = 200
		token := c.GetHeader("Authorization")

		if token == "" {
			code = 404 // 未携带token
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				code = 401 // 无效或过期
			} else if time.Now().Unix() > claims.ExpiresAt.Time.Unix() {
				code = 401 // 过期
			} else {
				c.Set("user_id", claims.UserID) // 将解析出的user_id存储到上下文，供Controller使用
			}
		}

		if code != 200 {
			c.JSON(http.StatusUnauthorized, models.Response{
				Status: code,
				Msg: "鉴权失败，请重新登录",
			})
			c.Abort() // 阻止后续处理
			return
		}

		c.Next() // 继续后续处理
	}
}