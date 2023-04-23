package admin_account

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
)

// AdminChangePasswordHandle 将用户名对应的管理员旧密码更改为新密码
// 1. 比对输入旧密码与数据库中密码是否相同 2. 将数据库中密码更新为新密码
func AdminChangePasswordHandle(context *gin.Context) {
	UserID, ok := context.MustGet("user_id").(int64)
	if !ok {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "解析用户信息错误"))
		return
	}
	oldPassword, ok := context.MustGet("old_password").(string) // 获取用户旧密码
	if !ok {                                                    // 断言错误 一般不会发生
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "加密密码错误"))
		return
	}
	newPassword, ok := context.MustGet("new_password").(string) // 获取用户新密码
	if !ok {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "加密密码错误"))
		return
	}
	passwordQuery, err := model.NewAdminAccountFlow().QueryAdminPasswordByUserID(UserID) // 根据用户名获取数据库中密码
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "查询密码错误"))
		return
	}
	if passwordQuery != oldPassword { // 旧密码与数据库中密码验证失败
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "旧密码输入错误"))
		return
	}
	if err = model.NewAdminAccountFlow().UpdateAdminPassword(UserID, newPassword); err != nil { // 更新数据库中密码
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error())) // 更新失败
		return
	}
	// 更新成功
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
	return
}
