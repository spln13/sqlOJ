package exercise

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/cache"
	"sqlOJ/common"
	"strconv"
	"time"
)

type SubmitMessage struct {
	UserID     int64
	UserType   int64
	ExerciseID int64
	Answer     string
	UserAgent  string
	SubmitTime time.Time
}

// SubmitHandle 处理用户提交题目的请求
func SubmitHandle(context *gin.Context) {
	userID, ok1 := context.MustGet("user_id").(int64)
	userType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(1, common.NewCommonResponse(1, "解析用户信息错误"))
		return
	}
	exerciseIDStr := context.PostForm("exercise_id")
	exerciseID, err := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "解析参数错误"))
		return
	}
	answer := context.PostForm("answer")
	userAgent := context.Request.UserAgent()
	ok, err := cache.CheckSubmitTimeValid(userID, userType, exerciseID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	if !ok { // 与上次发送未间隔3s
		context.JSON(http.StatusOK, common.NewCommonResponse(1, "请勿频繁发送"))
		return
	}
	if err := WriteMessage(userID, userType, exerciseID, answer, userAgent); err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "提交错误"))
		return
	}
	context.JSON(http.StatusOK, common.NewCommonResponse(0, "")) // 提交成功
}

var JudgeQueue = make(chan SubmitMessage, 2000) // 判题队列 缓冲区大小2000

// WriteMessage 将做题数据写入channel，并且在cache中保存当前提交的判题状态
func WriteMessage(userID, userType, exerciseID int64, answer, userAgent string) error {
	message := SubmitMessage{
		UserID:     userID,
		UserType:   userType,
		ExerciseID: exerciseID,
		Answer:     answer,
		UserAgent:  userAgent,
		SubmitTime: time.Now(),
	}
	JudgeQueue <- message // 将判题数据写入channel
	err := cache.SetJudgeStatusPending(userID, userType, exerciseID, time.Now())
	return err
}
