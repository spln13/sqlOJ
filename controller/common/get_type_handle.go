package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/utils"
)

type GetTypeResponse struct {
	Type int64 `json:"type"`
	utils.Response
}

// GetTypeHandle 获取token对应的用户类型
func GetTypeHandle(context *gin.Context) {
	userType, ok := context.MustGet("user_type").(int64)
	if !ok {
		context.JSON(http.StatusInternalServerError, GetTypeResponse{
			Type:     0,
			Response: utils.NewCommonResponse(1, "解析参数错误"),
		})
		return
	}
	context.JSON(http.StatusOK, GetTypeResponse{
		Type:     userType,
		Response: utils.NewCommonResponse(0, ""),
	})
}
