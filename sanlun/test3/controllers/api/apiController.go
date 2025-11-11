package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiController struct{}

func (con ApiController) Index (c *gin.Context) {
	c.String(http.StatusOK, "一个api接口")
}

func (con ApiController) Userlist (c *gin.Context) {
	c.String(http.StatusOK, "一个api接口-userlist")
}

func (con ApiController) Plist (c *gin.Context) {
	c.String(http.StatusOK, "一个api接口-plist")
}