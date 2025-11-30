package service

import (
	"Todolist/dao"
	"Todolist/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type TodoService struct{}

// CreateTodo 创建任务
// 逻辑：写入数据库 -> 删除该用户的缓存（确保下次查是新的）
func (s *TodoService) CreateTodo(todo *models.Todo) error {
	if err := dao.DB.Create(todo).Error; err != nil {
		return nil
	}

	s.clearCache(todo.UserID)
	return nil
}

// GetTodos 获取任务列表 (带 Redis 缓存)
// 逻辑：只缓存“第一页、无搜索关键词、无状态筛选”的热点数据，其他复杂查询直接走 DB
func (s *TodoService) GeTodos(userID uint, pageNum, pageSize int, status, keyword string) (interface{}, int64, error) {
	var todos []models.Todo
	var total int64

	CacheKey := fmt.Sprintf("todo_list_%d", userID)

	if keyword == "" && status == "" && pageNum == 1 {
		val, err := dao.RDB.Get(dao.Ctx, CacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(val), &todos); err == nil {
				return todos, int64(total), nil
			}
		}
	}

	query := dao.DB.Model(&models.Todo{}).Where("user_id = ?", userID)

	if keyword != "" {
		query = query.Where("title like ? or content like ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		statusINt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusINt)
	}

	offset := (pageNum - 1) * pageSize
	if err := query.Count(&total).Limit(pageSize).Offset(offset).Find(&todos).Error; err != nil {
		return nil, 0, err
	}

	if keyword == "" && status == "" && pageNum == 1 {
		data, _ := json.Marshal(todos)
		dao.RDB.Set(dao.Ctx, CacheKey, data, time.Hour)
	}

	return todos, int64(total), nil
}

func (s *TodoService) UpdateOneTodo(id string, userID uint, status int) error {
	var todo models.Todo
	if err := dao.DB.Where("id = ? and user_id = ?", id, userID).First(&todo).Error; err != nil {
		return err
	}

	if err := dao.DB.Model(&todo).Update("status = ?", status).Error; err != nil {
		return err
	}

	s.clearCache(userID)
	return nil
}

func (s *TodoService) UpdateAllTodos(userID uint, status int) error {
	if err := dao.DB.Model(&models.Todo{}).Where("user_id = ?", userID).Update("status", status).Error; err != nil {
		return err
	}

	s.clearCache(userID)
	return nil
}

func (s *TodoService) DeleteTodo(id string, userID uint) error {
	if err := dao.DB.Where("id = ? and user_id = ?", id, userID).Delete(&models.Todo{}).Error; err != nil {
		return err
	}

	s.clearCache(userID)
	return nil
}

func (s *TodoService) DeleteBatch(userID uint, deleteType string) error {
	query := dao.DB.Model(&models.Todo{}).Where("user_id = ?", userID)
	switch deleteType {
	case "1":
		query = query.Where("status = ?", 1)
	case "2":
		query = query.Where("status = ?", 2)
	case "3":
	default:
		return fmt.Errorf("参数错误")
	}

	if err := query.Delete(&models.Todo{}).Error; err != nil {
		return err
	}
	s.clearCache(userID)
	return nil
}

func (s *TodoService) clearCache(userID uint) {
	cacheKey := fmt.Sprintf("todo_list_%d", userID)
	dao.RDB.Del(dao.Ctx, cacheKey)
}
