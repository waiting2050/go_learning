package controllers

import (
	"Todolist/dao"
	"Todolist/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindBodyWithJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest,models.Response{
			Status: http.StatusBadRequest,
			Msg: "请求参数错误",
			Data: nil,
		})
		return
	}

	if err := dao.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status: http.StatusInternalServerError,
			Msg: "创建失败",
			Data: nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status: http.StatusOK,
		Msg: "创建成功",
		Data: todo,
	})
}