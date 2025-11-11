package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct{
	BaseController
}

func (con UserController) Index(c *gin.Context) {
	//c.String(http.StatusOK, "用户列表--")
	con.success(c)
}

func (con UserController) Add(c *gin.Context) {
	c.String(http.StatusOK, "用户列表-add--")
}

func (con UserController) Edit(c *gin.Context) {
	c.String(http.StatusOK, "用户列表-edit--")
}
