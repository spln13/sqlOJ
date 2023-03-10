package submission

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/model"
	"time"
)

type AllAllResponse struct {
	List []AllAll `json:"list"`
	common.Response
}

type AllAll struct {
	Answer       string    `json:"answer"`
	Status       int       `json:"status"`
	SubmitTime   time.Time `json:"submit_time"`
	UserID       int64     `json:"user_id"`
	UserType     int64     `json:"user_type"`
	Username     string    `json:"username"`
	ExerciseID   int64     `json:"exercise_id"`
	ExerciseName string    `json:"exercise_name"`
	UserAgent    string    `json:"user_agent"`
}

// GetAllAllHandle 获取所有的提交记录, 包括在cache中的
func GetAllAllHandle(context *gin.Context) {
	allSubmitHistory, err := model.NewSubmitHistoryFlow().QueryAllSubmitHistory()
	if err != nil {
		context.JSON(http.StatusInternalServerError, AllAllResponse{
			List:     nil,
			Response: common.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var allAllList []AllAll
	for _, submitHistory := range allSubmitHistory {
		userID := submitHistory.UserID
		userType := submitHistory.UserType
		var username string // 获取用户名
		if userType == 1 {  // 学生
			username = model.NewStudentAccountFlow().QueryStudentUsernameByUserID(userID)
		} else if userType == 2 { // 教师
			username = model.NewTeacherAccountFlow().QueryTeacherUsernameByUserID(userID)
		} else { // 管理员
			username = model.NewAdminAccountFlow().QueryAdminUsernameByUserID(userID)
		}
		exerciseID := submitHistory.ExerciseID
		exerciseName := model.NewExerciseContentFlow().QueryExerciseNameByExerciseID(exerciseID)
		allAll := AllAll{
			Answer:       submitHistory.StudentAnswer,
			Status:       submitHistory.Status,
			SubmitTime:   submitHistory.SubmitTime,
			UserAgent:    submitHistory.UserAgent,
			UserID:       userID,
			UserType:     userType,
			Username:     username,
			ExerciseID:   exerciseID,
			ExerciseName: exerciseName,
		}
		allAllList = append(allAllList, allAll)
	}
	context.JSON(http.StatusOK, AllAllResponse{
		List:     allAllList,
		Response: common.NewCommonResponse(0, ""),
	})
}
