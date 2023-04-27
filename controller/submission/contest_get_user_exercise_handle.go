package submission

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/cache"
	"sqlOJ/model"
	"sqlOJ/utils"
	"strconv"
	"strings"
	"time"
)

type ContestGetUserExerciseResponse struct {
	List []OneUserExerciseSubmission `json:"list"`
	utils.Response
}

type OneUserExerciseSubmission struct {
	OnChain    int       `json:"on_chain"`
	Status     int       `json:"status"`
	SubmitTime time.Time `json:"submit_time"`
}

func ContestGetUserExerciseHandle(context *gin.Context) {
	userID, ok1 := context.MustGet("user_id").(int64)
	userType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusBadRequest, ContestGetUserExerciseResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, "解析用户参数错误"),
		})
		return
	}
	contestIDStr := context.Query("contest_id")
	exerciseIDStr := context.Query("exercise_id")
	contestID, err1 := strconv.ParseInt(contestIDStr, 10, 64)
	exerciseID, err2 := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err1 != nil || err2 != nil {
		context.JSON(http.StatusBadRequest, ContestGetUserExerciseResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, "请求参数错误"),
		})
		return
	}
	// 查询缓存中的提交记录
	userCacheMap, err := cache.GetContestUserJudgeStatus(userID, userType, exerciseID, contestID)
	if err != nil {
		context.JSON(http.StatusBadRequest, ContestGetUserExerciseResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, "查询缓存错误"),
		})
		return
	}
	var oneUserExerciseSubmissionList []OneUserExerciseSubmission
	for key, value := range userCacheMap {
		keyArray := strings.Split(key, ":")
		submitTimeStr := keyArray[len(keyArray)-1]
		submitTimeUnix, _ := strconv.ParseInt(submitTimeStr, 10, 64)
		submitTime := time.Unix(submitTimeUnix, 0)
		valueInt, _ := strconv.Atoi(value)
		oneUserExerciseSubmission := OneUserExerciseSubmission{
			OnChain:    0,
			Status:     valueInt,
			SubmitTime: submitTime,
		}
		oneUserExerciseSubmissionList = append(oneUserExerciseSubmissionList, oneUserExerciseSubmission)
	}
	// 查询MySQL中的提交记录
	minContestSubmissionList, err := model.NewContestSubmissionFlow().QueryOneUserExerciseSubmission(userID, userType, contestID, exerciseID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, ContestGetUserExerciseResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	for _, minContestSubmission := range minContestSubmissionList {
		oneUserExerciseSubmission := OneUserExerciseSubmission{
			OnChain:    minContestSubmission.OnChain,
			Status:     minContestSubmission.Status,
			SubmitTime: minContestSubmission.SubmitTime,
		}
		oneUserExerciseSubmissionList = append(oneUserExerciseSubmissionList, oneUserExerciseSubmission)
	}
	context.JSON(http.StatusOK, ContestGetUserExerciseResponse{
		List:     oneUserExerciseSubmissionList,
		Response: utils.NewCommonResponse(0, ""),
	})
}
