package routers

import (
	"Todolist/controllers"
	"Todolist/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		// 用户模块（不需要登录即可访问）
		v1.POST("user/register", controllers.Register)
		v1.POST("user/login", controllers.Login)

		// 待办模块（需要登录）
		todo := v1.Group("/todo")
		todo.Use(middleware.JWT())
		{
			// 增
			todo.POST("", controllers.CreateTodo)

			// 查
			todo.GET("", controllers.GetTodo)

			// 改（单个）
			todo.PUT(":id", controllers.UpdateOneTodo)

			// 改（批量）
			todo.PUT("status/batch", controllers.UpdateAllTodos)

			// 删（单个）
			todo.DELETE("/:id", controllers.DeleteTodo)

			// 删（批量）
			todo.DELETE("batch", controllers.DeleteBatch)
		}
	}

	return r
}