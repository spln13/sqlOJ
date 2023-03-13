package admin_account

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/middlewares"
	"sqlOJ/model"
	"time"
)

func AdminLoginHandle(context *gin.Context) {
	username := context.Query("username")
	password, ok := context.MustGet("password_sha256").(string) // 获取到由中间件加密的密码
	if !ok {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "密码加密错误"))
		return
	}
	userID, passwordQuery, err := model.NewAdminAccountFlow().QueryAdminPasswordByUsername(username)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	if userID == 0 {
		context.JSON(http.StatusOK, common.NewCommonResponse(1, "用户不存在"))
		return
	}
	if password != passwordQuery { // 密码错误
		context.JSON(http.StatusOK, common.NewCommonResponse(1, "密码错误"))
		return
	}
	token, err := middlewares.ReleaseToken(userID, 3)
	log.Println(err)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "token颁发错误"))
	}
	expires := time.Now().Add(24 * time.Hour) // 设置过期时间为 24 小时后
	// 设置第一个 cookie，名称为 "username"
	cookie1 := http.Cookie{
		Name:     "username",
		Value:    username,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,                    // 禁止通过 JavaScript 访问 cookie
		SameSite: http.SameSiteStrictMode, // 禁止跨站点请求伪造攻击
	}
	http.SetCookie(context.Writer, &cookie1)

	// 设置第二个 cookie，名称为 "token"
	cookie2 := http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expires,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(context.Writer, &cookie2)
	context.JSON(http.StatusOK, common.NewCommonResponse(0, ""))
}
