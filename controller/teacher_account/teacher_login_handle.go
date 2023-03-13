package teacher_account

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/middlewares"
	"sqlOJ/model"
	"time"
)

func TeacherLoginHandle(context *gin.Context) {
	username := context.Query("username")
	password, ok := context.MustGet("password_sha256").(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "解析密码错误"))
		return
	}
	userID, passwordQuery, err := model.NewTeacherAccountFlow().QueryTeacherPasswordByUsername(username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	if password != passwordQuery { // 密码不匹配
		context.JSON(http.StatusBadRequest, common.NewCommonResponse(1, "密码错误"))
		return
	}
	// 验证通过，颁发token
	token, err := middlewares.ReleaseToken(userID, 3)
	log.Println(err)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "token颁发错误"))
		return
	}
	// 设置cookie过期时间
	expires := time.Now().Add(7 * 24 * time.Hour)
	// 设置cookie
	context.SetCookie("token", token, int(expires.Unix()), "/", "localhost:8080", true, false)
	context.SetCookie("username", username, int(expires.Unix()), "/", "localhost:8080", true, false)
	context.JSON(http.StatusOK, common.NewCommonResponse(0, ""))
}
