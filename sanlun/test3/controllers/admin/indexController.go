package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IndexController struct{}

func (con IndexController) Index(c *gin.Context) {
	c.String(http.StatusOK, "用户列表--")
}
