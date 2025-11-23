package controllers

import (
	"Todolist/dao"
	"Todolist/models"
	"Todolist/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status: http.StatusBadRequest,
			Msg:    "参数错误",
		})
	}

	var count int64
	dao.DB.Model(&models.User{}).Where("user_name = ?", user.UserName).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Status: http.StatusBadRequest,
			Msg:    "用户名已存在",
		})
		return
	}

	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.PasswordDigest), 10)
	user.PasswordDigest = string(bytes)
	dao.DB.Create(&user)
	c.JSON(http.StatusOK, models.Response{
		Status: http.StatusOK,
		Msg:    "注册成功",
	})
}

func Login(c *gin.Context) {
	var inputUser models.User
	if err := c.ShouldBindBodyWithJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status: http.StatusBadRequest,
			Msg:    "参数错误",
		})
		return
	}

	var dbUser models.User
	dao.DB.Where("user_name = ?", inputUser.UserName).First(&dbUser)
	if dbUser.ID == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Status: http.StatusBadRequest,
			Msg:    "用户不存在",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.PasswordDigest), []byte(inputUser.PasswordDigest)); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status: http.StatusBadRequest,
			Msg:    "密码错误",
		})
		return
	}

	token, _ := utils.GenerateToken(dbUser.ID)
	c.JSON(http.StatusOK, models.Response{
		Status: http.StatusOK,
		Msg:    "登录成功",
		Data: gin.H{
			"token": token,
		},
	})
}
