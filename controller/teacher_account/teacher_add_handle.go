package teacher_account

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/model"
)

func TeacherAddHandle(context *gin.Context) {
	username := context.Query("username")
	realName := context.Query("real_name")
	password, ok := context.MustGet("password_sha256").(string) // 获取到由中间件加密的密码
	if !ok {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "加密密码错误"))
		return
	}
	exist, err := model.NewTeacherAccountFlow().QueryTeacherExistByUsername(username)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "查询职工号错误"))
		return
	}
	if exist {
		context.JSON(http.StatusBadRequest, common.NewCommonResponse(1, "职工号已存在"))
		return
	}
	if err := model.NewTeacherAccountFlow().InsertTeacherAccount(username, password, realName); err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	context.JSON(http.StatusOK, common.NewCommonResponse(0, ""))
}
