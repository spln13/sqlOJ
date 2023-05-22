package submission

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
	"strconv"
	"time"
)

type ContestGetOneExerciseSubmissionResponse struct {
	OneExerciseSubmissionList []OneExerciseSubmission `json:"list"`
	utils.Response
}

type OneExerciseSubmission struct {
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

func ContestGetOneExerciseHandler(context *gin.Context) {
	contestIDStr := context.Query("contest_id")
	exerciseIDStr := context.Query("exercise_id")
	contestID, err1 := strconv.ParseInt(contestIDStr, 10, 64)
	exerciseID, err2 := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err1 != nil || err2 != nil {
		context.JSON(http.StatusBadRequest, ContestGetOneExerciseSubmissionResponse{
			OneExerciseSubmissionList: nil,
			Response:                  utils.NewCommonResponse(1, "请求参数错误"),
		})
		return
	}
	contestSubmissionList, err := model.NewContestSubmissionFlow().GetOneExerciseSubmission(contestID, exerciseID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, ContestGetOneExerciseSubmissionResponse{
			OneExerciseSubmissionList: nil,
			Response:                  utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var oneExerciseSubmissionList []OneExerciseSubmission
	for _, contestSubmission := range contestSubmissionList {
		oneExerciseSubmission := OneExerciseSubmission{
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
		oneExerciseSubmissionList = append(oneExerciseSubmissionList, oneExerciseSubmission)
	}
	context.JSON(http.StatusInternalServerError, ContestGetOneExerciseSubmissionResponse{
		OneExerciseSubmissionList: oneExerciseSubmissionList,
		Response:                  utils.NewCommonResponse(0, ""),
	})
	return
}
