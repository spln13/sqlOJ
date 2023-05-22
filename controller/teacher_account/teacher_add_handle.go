package teacher_account

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
)

func TeacherAddHandler(context *gin.Context) {
	username := context.Query("username")
	realName := context.Query("real_name")
	password, ok := context.MustGet("password_sha256").(string) // 获取到由中间件加密的密码
	if !ok {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "加密密码错误"))
		return
	}
	exist, err := model.NewTeacherAccountFlow().QueryTeacherExistByUsername(username)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "查询职工号错误"))
		return
	}
	if exist {
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(1, "职工号已存在"))
		return
	}
	userID, err := model.NewTeacherAccountFlow().InsertTeacherAccount(username, password, realName)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	err = model.NewScoreRecordFlow().InsertScoreRecord(userID, 2, username)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
}
