package main

import (
	"context"
	"silun/biz/auth"
	"silun/biz/handler"
	"silun/biz/model"
	"silun/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"
)

func registerRoutes(h *server.Hertz, db *gorm.DB, cfg *model.Config) {
	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService)

	videoService := service.NewVideoService(db)
	videoHandler := handler.NewVideoHandler(videoService)

	interactionService := service.NewInteractionService(db)
	interactionHandler := handler.NewInteractionHandler(interactionService)

	socialService := service.NewSocialService(db)
	socialHandler := handler.NewSocialHandler(socialService)

	// 分片上传服务
	uploadService := service.NewUploadService(db)
	uploadHandler := handler.NewUploadHandler(uploadService, videoService)

	// 上传策略服务
	strategyService := service.NewUploadStrategyService(nil)
	strategyHandler := handler.NewUploadStrategyHandler(strategyService)

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
	user.PUT("/avatar/upload", auth.AuthMiddleware(), userHandler.UploadAvatar)

	video := h.Group("/video")
	video.POST("/publish", auth.AuthMiddleware(), videoHandler.PublishVideo)
	video.GET("/publish/list", videoHandler.GetPublishList)
	video.POST("/search", videoHandler.SearchVideo)
	video.GET("/popular", videoHandler.GetPopularVideos)

	// 上传策略接口（公开，用于前端决策）
	strategy := h.Group("/upload/strategy")
	strategy.GET("/decide", strategyHandler.GetUploadStrategy)
	strategy.GET("/recommendation", strategyHandler.GetUploadRecommendation)

	// 分片上传接口
	upload := h.Group("/upload")
	upload.POST("/init", auth.AuthMiddleware(), uploadHandler.InitUpload)
	upload.POST("/chunk", auth.AuthMiddleware(), uploadHandler.UploadChunk)
	upload.GET("/status", auth.AuthMiddleware(), uploadHandler.GetUploadStatus)
	upload.POST("/merge", auth.AuthMiddleware(), uploadHandler.MergeChunks)
	upload.POST("/cancel", auth.AuthMiddleware(), uploadHandler.CancelUpload)

	like := h.Group("/like")
	like.POST("/action", auth.AuthMiddleware(), interactionHandler.LikeAction)
	like.GET("/list", interactionHandler.GetLikeList)

	comment := h.Group("/comment")
	comment.POST("/publish", auth.AuthMiddleware(), interactionHandler.PublishComment)
	comment.GET("/list", interactionHandler.GetCommentList)
	comment.POST("/delete", auth.AuthMiddleware(), interactionHandler.DeleteComment)

	relation := h.Group("/relation")
	relation.POST("/action", auth.AuthMiddleware(), socialHandler.FollowAction)

	h.GET("/following/list", socialHandler.GetFollowList)
	h.GET("/follower/list", socialHandler.GetFollowerList)
	h.GET("/friends/list", auth.AuthMiddleware(), socialHandler.GetFriendList)
}
