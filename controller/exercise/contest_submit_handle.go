package exercise

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/cache"
	"strconv"
	"time"
)

func ContestSubmitHandle(context *gin.Context) {
	userID, ok1 := context.MustGet("user_id").(int64)
	userType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "解析用户身份错误"))
		return
	}
	contestIDStr := context.PostForm("contest_id")
	exerciseIDStr := context.PostForm("exercise_id")
	answer := context.PostForm("answer")
	contestID, err1 := strconv.ParseInt(contestIDStr, 10, 64)
	exerciseID, err2 := strconv.ParseInt(exerciseIDStr, 10, 64)
	userAgent := context.Request.UserAgent()
	if err1 != nil || err2 != nil {
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(1, "请求参数错误"))
		return
	}
	err := WriteContestMessage(userID, userType, exerciseID, contestID, answer, userAgent)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "写入判题队列错误"))
		return
	}
}

// WriteContestMessage 将做题数据写入channel，并且在cache中保存当前提交的判题状态
func WriteContestMessage(userID, userType, exerciseID, contestID int64, answer, userAgent string) error {
	message := SubmitMessage{
		UserID:     userID,
		UserType:   userType,
		ExerciseID: exerciseID,
		Answer:     answer,
		UserAgent:  userAgent,
		IsContest:  true,
		ContestID:  contestID,
		SubmitTime: time.Now(),
	}
	JudgeQueue <- message // 将判题数据写入channel
	err := cache.SetContestJudgeStatusPending(userID, userType, exerciseID, contestID, time.Now())
	return err
}
