package submission

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"time"
)

type OneAllResponse struct {
	List       []OneAll `json:"list"`
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
}

type OneAll struct {
	SubmissionID int64     `json:"submission_id"`
	Answer       string    `json:"answer"`
	ExerciseID   int64     `json:"exercise_id"`
	ExerciseName string    `json:"exercise_name"`
	Status       int       `json:"status"`
	SubmitTime   time.Time `json:"submit_time"`
	UserAgent    string    `json:"user_agent"`
}

// GetOneAllHandle 查询当前用户所有提交记录
func GetOneAllHandle(context *gin.Context) {
	userID, ok1 := context.MustGet("user_id").(int64)
	userType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusInternalServerError, OneOneResponse{
			List:       nil,
			StatusCode: 1,
			StatusMsg:  "解析用户信息错误",
		})
		return
	}
	submitHistoryArray, err := model.NewSubmitHistoryFlow().QueryThisUserSubmitHistory(userID, userType)
	if err != nil {
		context.JSON(http.StatusInternalServerError, OneOneResponse{
			List:       nil,
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	var oneAllList []OneAll
	for _, submitHistory := range submitHistoryArray {
		exerciseID := submitHistory.ExerciseID
		exerciseName := model.NewExerciseContentFlow().QueryExerciseNameByExerciseID(exerciseID)
		oneAll := OneAll{
			SubmissionID: submitHistory.ID,
			Answer:       submitHistory.StudentAnswer,
			ExerciseID:   exerciseID,
			ExerciseName: exerciseName,
			Status:       submitHistory.Status,
			SubmitTime:   submitHistory.SubmitTime,
			UserAgent:    submitHistory.UserAgent,
		}
		oneAllList = append(oneAllList, oneAll)
	}
	context.JSON(http.StatusOK, OneAllResponse{
		List:       oneAllList,
		StatusCode: 0,
		StatusMsg:  "",
	})
}
