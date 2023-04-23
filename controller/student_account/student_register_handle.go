package student_account

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/cache"
	"sqlOJ/model"
	"sqlOJ/utils"
)

func StudentRegisterHandle(context *gin.Context) {
	username := context.PostForm("username")
	number := context.PostForm("number")
	email := context.PostForm("email")
	realName := context.PostForm("real_name")
	code := context.PostForm("code")
	password, ok := context.MustGet("password_sha256").(string)
	if !ok { // 断言失败
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "加密密码存储错误"))
		return
	}
	// 校验用户名是否存在
	exist, err := model.NewStudentAccountFlow().QueryStudentExistByUsername(username)
	if err != nil { // 查询出错
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if exist {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "用户名已存在"))
		return
	}
	// 再次查询此邮箱是否已经被注册过
	exist, err = model.NewStudentAccountFlow().QueryStudentExistByEmail(email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if exist {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "邮箱已被注册"))
		return
	}
	ok, err = cache.VerifyEmailCode(email, code)
	if !ok {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, err.Error()))
		return
	}
	userID, err := model.NewStudentAccountFlow().InsertStudentAccount(username, password, number, realName, email)
	err = model.NewScoreRecordFlow().InsertScoreRecord(userID, 1, username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
}
