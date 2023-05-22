package student_account

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
	"strconv"
)

func StudentResetHandler(context *gin.Context) {
	studentIDStr := context.Query("student_id")
	studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "请求参数错"))
		return
	}
	if err := model.NewStudentAccountFlow().ResetStudentPassword(studentID); err != nil {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, err.Error()))
		return
	}
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
}
