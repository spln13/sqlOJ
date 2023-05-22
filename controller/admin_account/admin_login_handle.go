package admin_account

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/middlewares"
	"sqlOJ/model"
	"sqlOJ/utils"
	"time"
)

func AdminLoginHandler(context *gin.Context) {
	username := context.Query("username")
	password, ok := context.MustGet("password_sha256").(string) // 获取到由中间件加密的密码
	if !ok {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "密码加密错误"))
		return
	}
	userID, passwordQuery, err := model.NewAdminAccountFlow().QueryAdminPasswordByUsername(username)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if userID == 0 {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "用户不存在"))
		return
	}
	if password != passwordQuery { // 密码错误
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "密码错误"))
		return
	}
	token, err := middlewares.ReleaseToken(userID, 3)
	log.Println(err)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "token颁发错误"))
	}
	// 设置cookie过期时间
	expires := time.Now().Add(7 * 24 * time.Hour)
	// 设置cookie
	context.SetCookie("token", token, int(expires.Unix()), "/", "localhost:8080", true, false)
	context.SetCookie("username", username, int(expires.Unix()), "/", "localhost:8080", true, false)
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
}
