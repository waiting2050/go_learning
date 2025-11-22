package models

import "gorm.io/gorm"

// Todo 代表一个待办事项
type Todo struct {
	gorm.Model        // 嵌入GORM模型，包含ID, CreatedAt, UpdatedAt, DeletedAt字段
	UserID     uint   `json:"user_id"`           // 关联用户ID
	Title      string `json:"title"`             // 待办事项标题
	Content    string `json:"content"`           // 待办事项内容
	Status     int    `json:"status"`            // 待办事项状态 (例如: 0-未开始, 1-已完成)
	StartTime  int64  `json:"start_time"`        // 开始时间戳
	EndTime    int64  `json:"end_time"`          // 结束时间戳
	Priority   int    `json:"priority"`          // 优先级 (例如: 0-低, 1-中, 2-高)
}

func (Todo) TableName() string {
	return "todo"
}
