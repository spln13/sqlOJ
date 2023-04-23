package submission

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
	"strconv"
	"time"
)

type ContestMySubmissionResponse struct {
	MySubmissionList []MySubmission `json:"list"`
	utils.Response
}

type MySubmission struct {
	Answer       string    `json:"answer"`
	ExerciseID   int64     `json:"exercise_id"`
	ExerciseName string    `json:"exercise_name"`
	Status       int       `json:"status"`
	SubmitTime   time.Time `json:"submit_time"`
}

func ContestGetMySubmissionHandle(context *gin.Context) {
	userID, ok1 := context.MustGet("user_id").(int64)
	userType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusBadRequest, ContestMySubmissionResponse{
			MySubmissionList: nil,
			Response:         utils.NewCommonResponse(1, "解析用户信息错误"),
		})
		return
	}
	contestIDStr := context.Query("contest_id")
	contestID, err := strconv.ParseInt(contestIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, ContestMySubmissionResponse{
			MySubmissionList: nil,
			Response:         utils.NewCommonResponse(1, "请求参数错误"),
		})
		return
	}
	contestSubmissionList, err := model.NewContestSubmissionFlow().GetUserContestSubmission(userID, userType, contestID)
	if err != nil {
		context.JSON(http.StatusBadRequest, ContestMySubmissionResponse{
			MySubmissionList: nil,
			Response:         utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var MySubmissionList []MySubmission
	for _, contestSubmission := range contestSubmissionList {
		oneSubmission := MySubmission{
			Answer:       contestSubmission.UserAnswer,
			ExerciseID:   contestSubmission.ExerciseID,
			ExerciseName: contestSubmission.ExerciseName,
			Status:       contestSubmission.Status,
			SubmitTime:   contestSubmission.SubmitTime,
		}
		MySubmissionList = append(MySubmissionList, oneSubmission)
	}
	context.JSON(http.StatusOK, ContestMySubmissionResponse{
		MySubmissionList: MySubmissionList,
		Response:         utils.NewCommonResponse(0, ""),
	})
}
