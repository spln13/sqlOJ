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

func StudentDeleteHandler(context *gin.Context) {
	studentIDStr := context.Query("student_id")
	studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "请求参数错"))
		return
	}
	// 删除 ranking , student_account, user_problem_status, contest_exercise_status
	err = model.NewScoreRecordFlow().DeleteRanking(studentID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	err = model.NewUserProblemStatusFlow().DeleteStudentStatus(studentID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	err = model.NewContestExerciseStatusFlow().DeleteContestStudentStatus(studentID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	err = model.NewStudentAccountFlow().DeleteStudentByID(studentID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
}
