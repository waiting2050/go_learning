package handler

import (
	"context"
	"silun/pkg/auth"
	"silun/pkg/config"
	"silun/pkg/database"
	"silun/pkg/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"
)

func RegisterRoutes(h *server.Hertz, db *gorm.DB, cfg *config.Config) {
	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService)

	public := h.Group("/")
	public.GET("/ping", func(ctx context.Context, c *app.RequestContext) {
		c.JSON(200, map[string]interface{}{
			"message": "pong",
		})
	})

	user := h.Group("/user")
	user.POST("/register", userHandler.Register)
	user.POST("/login", userHandler.Login)
	user.GET("/info", userHandler.GetUserInfo)
	user.POST("/avatar", userHandler.UploadAvatar)
}
