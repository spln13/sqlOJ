package submission

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/model"
	"strconv"
	"time"
)

type AllOneResponse struct {
	List []AllOne `json:"list"`
	common.Response
}

type AllOne struct {
	Answer     string    `json:"answer"`
	Status     int       `json:"status"`
	SubmitTime time.Time `json:"submit_time"`
	UserID     int64     `json:"user_id"`
	UserType   int64     `json:"user_type"`
	Username   string    `json:"username"`
	UserAgent  string    `json:"user_agent"`
}

// GetAllOneHandle 获取当前题目所有用户的提交
func GetAllOneHandle(context *gin.Context) {
	exerciseIDStr := context.Query("exercise_id")
	exerciseID, err := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, AllAllResponse{
			List:     nil,
			Response: common.NewCommonResponse(1, "解析题目信息错误"),
		})
		return
	}
	var AllOneList []AllOne
	submitHistoryArray, err := model.NewSubmitHistoryFlow().QueryThisExerciseSubmitHistory(exerciseID)
	for _, submitHistory := range submitHistoryArray {
		userID := submitHistory.UserID
		userType := submitHistory.UserType
		var username string
		if userType == 1 { // 学生
			username = model.NewStudentAccountFlow().QueryStudentUsernameByUserID(userID)
		} else if userType == 2 { // 教师
			username = model.NewTeacherAccountFlow().QueryTeacherUsernameByUserID(userID)
		} else { // 管理员
			username = model.NewAdminAccountFlow().QueryAdminUsernameByUserID(userID)
		}
		allOne := AllOne{
			Answer:     submitHistory.StudentAnswer,
			Status:     submitHistory.Status,
			SubmitTime: submitHistory.SubmitTime,
			UserID:     userID,
			UserType:   userType,
			Username:   username,
			UserAgent:  submitHistory.UserAgent,
		}
		AllOneList = append(AllOneList, allOne)
	}
	context.JSON(http.StatusOK, AllOneResponse{
		List:     AllOneList,
		Response: common.NewCommonResponse(0, ""),
	})
}
