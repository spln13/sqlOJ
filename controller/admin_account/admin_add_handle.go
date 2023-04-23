package admin_account

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
)

func AdminAddHandle(context *gin.Context) {
	username := context.Query("username")
	password, ok := context.MustGet("password_sha256").(string) // 获取到由中间件加密的密码
	if !ok {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "加密密码错误"))
		return
	}
	exist, err := model.NewAdminAccountFlow().QueryAdminExistByUsername(username)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "查询用户名错误"))
		return
	}
	if exist {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "用户名已存在"))
		return
	}
	userID, err := model.NewAdminAccountFlow().InsertAdminAccount(username, password)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if err := model.NewScoreRecordFlow().InsertScoreRecord(userID, 3, username); err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
}
