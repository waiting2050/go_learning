package controllers

import (
	"Todolist/models"
	"Todolist/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var TodoService = new(service.TodoService)

const codeSuccess = 10000

// 增
func CreateTodo(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindBodyWithJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status: 40001,
			Msg:    "参数错误",
			Data:   nil,
		})
		return
	}

	userID, _ := c.Get("user_id")
	todo.UserID = userID.(uint)
	if err := TodoService.CreateTodo(&todo); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status: 50001,
			Msg:    "创建失败",
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status: codeSuccess,
		Msg:    "创建成功",
		Data:   todo,
	})
}

// 查（含分页、搜索、状态筛选）
func GetTodo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)

	pageNum, _ := strconv.Atoi(c.DefaultQuery("page_num", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	items, total, err := TodoService.GetTodos(uid, pageNum, pageSize, status, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status: 50002,
			Msg:    "查询失败",
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status: codeSuccess,
		Msg:    "查询成功",
		Data: models.DataList{
			Items: items,
			Total: total,
		},
	})
}

// 用来接收前端传来的目标状态
type statusForm struct {
	Status int `json:"status"`
}

// 改（单个）
func UpdateOneTodo(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var form statusForm
	if err := c.ShouldBindBodyWithJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status: 40001,
			Msg:    "参数错误",
		})
		return
	}

	if err := TodoService.UpdateOneTodo(id, userID.(uint), form.Status); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status: 40401,
			Msg:    "未找到任务或无权限",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status: codeSuccess,
		Msg:    "更新成功",
	})
}

// 改（批量）
func UpdateAllTodos(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var form statusForm
	if err := c.ShouldBindBodyWithJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status: 50003,
			Msg:    "参数错误",
		})
		return
	}

	// 批量更新
	if err := TodoService.UpdateAllTodos(userID.(uint), form.Status); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status: 50003,
			Msg:    "批量操作失败",
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status: codeSuccess,
		Msg:    "批量操作成功",
	})
}

// 删(单个)
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	if err := TodoService.DeleteTodo(id, userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status: 50004,
			Msg:    "删除失败",
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status: codeSuccess,
		Msg:    "删除成功",
	})
}

// 删（批量）
// type 参数说明: 1:已完成, 2:待办, 3:全部
func DeleteBatch(c *gin.Context) {
	userID, _ := c.Get("user_id")
	deleteType := c.Query("type")

	if err := TodoService.DeleteBatch(userID.(uint), deleteType); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status: 40002,
			Msg:    "参数错误或删除失败",
			Data:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status: codeSuccess,
		Msg:    "删除成功",
	})
}
