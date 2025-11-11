package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

func (con BaseController) success(c *gin.Context) {
	c.String(http.StatusOK, "成功")
}

func (con BaseController) error(c *gin.Context) {
	c.String(http.StatusOK, "失败")
}
