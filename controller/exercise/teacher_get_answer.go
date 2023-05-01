package exercise

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"strconv"
)

type TeacherAnswerResponse struct {
	Answer     string `json:"answer"`
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func TeacherGetAnswer(context *gin.Context) {
	exerciseIDStr := context.Query("exercise_id")
	exerciseID, err := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, TeacherAnswerResponse{
			Answer:     "",
			StatusCode: 1,
			StatusMsg:  "请求参数错",
		})
		return
	}
	answer, _ := model.NewExerciseContentFlow().QueryAnswerTypeByExerciseID(exerciseID)
	context.JSON(http.StatusOK, TeacherAnswerResponse{
		Answer:     answer,
		StatusCode: 0,
		StatusMsg:  "",
	})
}
