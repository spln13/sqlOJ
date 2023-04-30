package submission

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
	"time"
)

type AllAllResponse struct {
	List []AllAll `json:"list"`
	utils.Response
}

type AllAll struct {
	SubmissionID int64     `json:"submission_id"`
	Answer       string    `json:"answer"`
	Status       int       `json:"status"`
	SubmitTime   time.Time `json:"submit_time"`
	UserID       int64     `json:"user_id"`
	UserType     int64     `json:"user_type"`
	Username     string    `json:"username"`
	ExerciseID   int64     `json:"exercise_id"`
	ExerciseName string    `json:"exercise_name"`
	UserAgent    string    `json:"user_agent"`
	OnChain      int       `json:"on_chain"`
}

// GetAllAllHandle 获取所有的提交记录, 包括在cache中的
func GetAllAllHandle(context *gin.Context) {
	allSubmitHistory, err := model.NewSubmitHistoryFlow().QueryAllSubmitHistory()
	if err != nil {
		context.JSON(http.StatusInternalServerError, AllAllResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var allAllList []AllAll
	for _, submitHistory := range allSubmitHistory {
		userID := submitHistory.UserID
		userType := submitHistory.UserType
		username := utils.QueryUsername(userID, userType)
		exerciseID := submitHistory.ExerciseID
		exerciseName := model.NewExerciseContentFlow().QueryExerciseNameByExerciseID(exerciseID)
		allAll := AllAll{
			SubmissionID: submitHistory.ID,
			Answer:       submitHistory.StudentAnswer,
			Status:       submitHistory.Status,
			SubmitTime:   submitHistory.SubmitTime,
			UserAgent:    submitHistory.UserAgent,
			UserID:       userID,
			UserType:     userType,
			Username:     username,
			ExerciseID:   exerciseID,
			ExerciseName: exerciseName,
			OnChain:      submitHistory.OnChain,
		}
		allAllList = append(allAllList, allAll)
	}
	context.JSON(http.StatusOK, AllAllResponse{
		List:     allAllList,
		Response: utils.NewCommonResponse(0, ""),
	})
}
