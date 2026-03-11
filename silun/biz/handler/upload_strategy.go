package handler

import (
	"context"
	"silun/biz/service"
	"silun/biz/utils"

	"github.com/cloudwego/hertz/pkg/app"
)

// UploadStrategyHandler 上传策略处理器
type UploadStrategyHandler struct {
	strategyService *service.UploadStrategyService
}

// NewUploadStrategyHandler 创建上传策略处理器
func NewUploadStrategyHandler(strategyService *service.UploadStrategyService) *UploadStrategyHandler {
	return &UploadStrategyHandler{strategyService: strategyService}
}

// GetUploadStrategy 获取上传策略建议
func (h *UploadStrategyHandler) GetUploadStrategy(ctx context.Context, c *app.RequestContext) {
	var req struct {
		FileName       string `form:"file_name" json:"file_name" binding:"required"`
		FileSize       int64  `form:"file_size" json:"file_size" binding:"required"`
		ContentType    string `form:"content_type" json:"content_type"`
		NetworkType    string `form:"network_type" json:"network_type"`
		UserPreference string `form:"user_preference" json:"user_preference"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		utils.Error(c, -1, "invalid request parameters")
		return
	}

	// 构建决策请求
	decisionReq := &service.UploadDecisionRequest{
		FileName:       req.FileName,
		FileSize:       req.FileSize,
		ContentType:    req.ContentType,
		NetworkType:    req.NetworkType,
		UserPreference: req.UserPreference,
	}

	// 如果网络类型未提供，使用默认值
	if decisionReq.NetworkType == "" {
		decisionReq.NetworkType = service.DetectNetworkType("")
	}

	// 获取决策
	decision := h.strategyService.DecideUploadStrategy(decisionReq)

	// 记录决策日志
	h.strategyService.LogUploadDecision(decisionReq, decision)

	utils.Success(c, decision)
}

// GetUploadRecommendation 获取上传建议
func (h *UploadStrategyHandler) GetUploadRecommendation(ctx context.Context, c *app.RequestContext) {
	fileSize, err := utils.ParseInt64(c.Query("file_size"))
	if err != nil || fileSize <= 0 {
		utils.Error(c, -1, "invalid file_size")
		return
	}

	fileName := c.Query("file_name")
	if fileName == "" {
		utils.Error(c, -1, "file_name is required")
		return
	}

	recommendation := h.strategyService.GetUploadRecommendation(fileSize, fileName)

	utils.Success(c, recommendation)
}
