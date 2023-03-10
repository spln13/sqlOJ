package admin_account

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/model"
)

func AdminAddHandle(context *gin.Context) {
	username := context.Query("username")
	password, ok := context.MustGet("password_sha256").(string) // 获取到由中间件加密的密码
	if !ok {
		context.JSON(http.StatusOK, common.NewCommonResponse(1, "加密密码错误"))
		return
	}
	exist, err := model.NewAdminAccountFlow().QueryAdminExistByUsername(username)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "查询用户名错误"))
		return
	}
	if exist {
		context.JSON(http.StatusOK, common.NewCommonResponse(1, "用户名已存在"))
		return
	}
	if err := model.NewAdminAccountFlow().InsertAdminAccount(username, password); err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	context.JSON(http.StatusOK, common.NewCommonResponse(0, ""))
}
