package handler

import (
	"context"
	"silun/pkg/auth"
	"silun/pkg/config"
	"silun/pkg/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"
)

func RegisterRoutes(h *server.Hertz, db *gorm.DB, cfg *config.Config) {
	userService := service.NewUserService(db)
	userHandler := NewUserHandler(userService)

	videoService := service.NewVideoService(db)
	videoHandler := NewVideoHandler(videoService)

	interactionService := service.NewInteractionService(db)
	interactionHandler := NewInteractionHandler(interactionService)

	socialService := service.NewSocialService(db)
	socialHandler := NewSocialHandler(socialService)

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

	video := h.Group("/video")
	video.POST("/publish", videoHandler.PublishVideo)
	video.GET("/publish/list", videoHandler.GetPublishList)
	video.POST("/search", videoHandler.SearchVideo)
	video.GET("/popular", videoHandler.GetPopularVideos)

	like := h.Group("/like")
	like.POST("/action", auth.AuthMiddleware(), interactionHandler.LikeAction)
	like.GET("/list", interactionHandler.GetLikeList)

	comment := h.Group("/comment")
	comment.POST("/publish", auth.AuthMiddleware(), interactionHandler.PublishComment)
	comment.GET("/list", interactionHandler.GetCommentList)
	comment.POST("/delete", auth.AuthMiddleware(), interactionHandler.DeleteComment)

	relation := h.Group("/relation")
	relation.POST("/action", auth.AuthMiddleware(), socialHandler.FollowAction)
	relation.GET("/follow/list", socialHandler.GetFollowList)
	relation.GET("/follower/list", socialHandler.GetFollowerList)
	relation.GET("/friend/list", auth.AuthMiddleware(), socialHandler.GetFriendList)
}
