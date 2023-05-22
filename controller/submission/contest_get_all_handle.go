package submission

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
	"strconv"
	"time"
)

type ContestAllSubmissionResponse struct {
	SubmissionList []Submission `json:"list"`
	utils.Response
}

type Submission struct {
	Answer       string    `json:"answer"`
	ExerciseID   int64     `json:"exercise_id"`
	ExerciseName string    `json:"exercise_name"`
	Status       int       `json:"status"`
	SubmitTime   time.Time `json:"submit_time"`
	UserAgent    string    `json:"user_agent"`
	UserID       int64     `json:"user_id"`
	UserType     int64     `json:"user_type"`
	Username     string    `json:"username"`
}

// ContestGetAllSubmissionHandler 获取一场竞赛中所有的提交
func ContestGetAllSubmissionHandler(context *gin.Context) {
	contestIDStr := context.Query("contest_id")
	contestID, err := strconv.ParseInt(contestIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, ContestAllSubmissionResponse{
			SubmissionList: nil,
			Response:       utils.NewCommonResponse(1, "请求参数错误"),
		})
		return
	}
	contestSubmissionList, err := model.NewContestSubmissionFlow().GetContestSubmissionByID(contestID)
	if err != nil {
		context.JSON(http.StatusBadRequest, ContestAllSubmissionResponse{
			SubmissionList: nil,
			Response:       utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var submissionList []Submission
	for _, contestSubmission := range contestSubmissionList {
		oneSubmission := Submission{
			Answer:       contestSubmission.UserAnswer,
			ExerciseID:   contestSubmission.ExerciseID,
			ExerciseName: contestSubmission.ExerciseName,
			Status:       contestSubmission.Status,
			SubmitTime:   contestSubmission.SubmitTime,
			UserAgent:    contestSubmission.UserAgent,
			UserID:       contestSubmission.UserID,
			UserType:     contestSubmission.UserType,
			Username:     contestSubmission.Username,
		}
		submissionList = append(submissionList, oneSubmission)
	}
	context.JSON(http.StatusOK, ContestAllSubmissionResponse{
		SubmissionList: submissionList,
		Response:       utils.NewCommonResponse(0, ""),
	})
}
