package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController
}

func (con UserController) Index(c *gin.Context) {
	// c.String(200, "用户列表--")
	con.success(c)
}
func (con UserController) Add(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/useradd.html", gin.H{})
}
func (con UserController) Edit(c *gin.Context) {
	c.String(200, "用户列表-Edit------")
}
