package controllers

import (
	"Todolist/dao"
	"Todolist/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 增
func CreateTodo(c *gin.Context) {
	var todo models.Todo
	c.ShouldBindBodyWithJSON(&todo)

	userID, _ := c.Get("user_id")
	todo.UserID = userID.(uint)

	dao.DB.Create(&todo)
	c.JSON(http.StatusOK, models.Response{
		Status: http.StatusOK,
		Msg:    "创建成功",
		Data:   todo,
	})
}

// 查（含分页、搜索、状态筛选）
func GetTodo(c *gin.Context) {
	var todos []models.Todo
	var total int64

	userID, _ := c.Get("user_id")

	// 1.基础查询
	query := dao.DB.Model(&models.Todo{}).Where("user_id = ?", userID)

	// 2.关键词搜索
	keyword := c.Query("keyword")
	if keyword != "" {
		query = query.Where("title like ? or content like ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 3.状态筛选（0：未完成， 1：已完成， 不传：所有）
	status := c.Query("status")
	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}

	// 4.分页
	pageNum, _ := strconv.Atoi(c.DefaultQuery("page_num", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (pageNum - 1) * pageSize

	// 执行查询
	query.Count(&total).Limit(pageSize).Offset(offset).Find(&todos)

	c.JSON(http.StatusOK, models.Response{
		Status: http.StatusOK,
		Msg:    "查询成功",
		Data: models.DataList{
			Items: todos,
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
			Status: http.StatusBadRequest,
			Msg:    "参数错误",
		})
		return
	}

	var todo models.Todo
	// 校验UserID
	if err := dao.DB.Where("id = ? and user_id = ?", id, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Status: http.StatusNotFound,
			Msg:    "未找到任务",
		})
		return
	}

	err := dao.DB.Model(&todo).Update("status", form.Status).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status: http.StatusInternalServerError,
			Msg:    "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status: http.StatusOK,
		Msg:    "更新成功",
	})
}

// 改（批量）
func UpdateAllTodos(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var form statusForm
	if err := c.ShouldBindBodyWithJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Status: http.StatusBadRequest,
			Msg:    "参数错误",
		})
		return
	}

	// 批量更新
	result := dao.DB.Model(&models.Todo{}).Where("user_id = ?", userID).Update("status", form.Status)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status: http.StatusInternalServerError,
			Msg:    "批量更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status: http.StatusOK,
		Msg:    "批量操作成功",
		Data: models.DataList{
			Items: nil,
			Total: result.RowsAffected,
		},
	})
}

// 删(单个)
func DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	if err := dao.DB.Where("id = ? and user_id = ?", id, userID).Delete(&models.Todo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status: http.StatusInternalServerError,
			Msg:    "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status: http.StatusOK,
		Msg:    "删除成功",
	})
}

// 删（批量）
// type 参数说明: 1:已完成, 2:待办, 3:全部
func DeleteBatch(c *gin.Context) {
	userID, _ := c.Get("user_id")
	deleteType := c.Query("type")
	query := dao.DB.Where("user_id = ?", userID)

	switch deleteType {
	case "1": // 删除所有已完成
		query = query.Where("status = ?", 1)
	case "2": // 删除所有待办
		query = query.Where("status = ?", 0)
	case "3": // 删除所有
	default:
		c.JSON(http.StatusBadRequest, models.Response{
			Status: http.StatusBadRequest,
			Msg:    "参数错误，请指定type(1:已完成 2:待办 3:全部)",
		})
		return
	}

	if err := query.Delete(&models.Todo{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Status: http.StatusInternalServerError,
			Msg:    "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Status: http.StatusOK,
		Msg:    "删除成功",
	})
}
