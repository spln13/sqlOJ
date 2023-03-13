package student_account

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/middlewares"
	"sqlOJ/model"
	"strconv"
	"time"
)

func StudentLoginHandle(context *gin.Context) {
	usernameEmail := context.Query("username_email") // 用户名或者密码
	codeStr := context.Query("code")                 // code: 1-> 用户名登录; 2 -> 邮箱登录
	code, err := strconv.Atoi(codeStr)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "解析请求错误"))
		return
	}
	password, ok := context.MustGet("password_sha256").(string)
	if !ok {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "解析密码错误"))
		return
	}

	var (
		userID        int64
		passwordQuery string
		username      string
	)
	if code == 1 { // 使用用户名登录
		username = usernameEmail
		userID, passwordQuery, err = model.NewStudentAccountFlow().QueryStudentPasswordByUsername(usernameEmail)
		if err != nil {
			context.JSON(http.StatusBadRequest, common.NewCommonResponse(1, err.Error()))
			return
		}
	} else { // 使用邮箱登录
		userID, passwordQuery, username, err = model.NewStudentAccountFlow().QueryStudentPasswordByEmail(usernameEmail)
		if err != nil {
			context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
			return
		}
	}
	if userID == 0 { // 用户不存在
		context.JSON(http.StatusOK, common.NewCommonResponse(1, "用户不存在"))
		return
	}
	if password != passwordQuery {
		context.JSON(http.StatusOK, common.NewCommonResponse(1, "密码错误"))
		return
	}
	// 颁发token
	token, err := middlewares.ReleaseToken(userID, 1) // 学生等级为1
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
