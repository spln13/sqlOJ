package submission

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/cache"
	"sqlOJ/model"
	"strconv"
	"strings"
	"time"
)

type OneOneResponse struct {
	List       []OneOne `json:"list"`
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
}

type OneOne struct {
	Answer     string    `json:"answer"`
	Status     int       `json:"status"`
	OnChain    int       `json:"on_chain"`
	SubmitTime time.Time `json:"submit_time"`
}

func GetOneOneHandler(context *gin.Context) {
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
	exerciseIDStr := context.Query("exercise_id")
	// 将exerciseID转为int64
	exerciseID, err := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, OneOneResponse{
			List:       nil,
			StatusCode: 1,
			StatusMsg:  "解析题目信息错误",
		})
		return
	}
	userCacheMap, err := cache.GetUserJudgeStatus(userID, userType, exerciseID) // 查询缓存中的提交信息
	if err != nil {
		context.JSON(http.StatusInternalServerError, OneOneResponse{
			List:       nil,
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	var OneOneList []OneOne
	for key, value := range userCacheMap {
		keyArray := strings.Split(key, ":")
		submitTimeStr := keyArray[len(keyArray)-1]
		submitTimeUnix, err := strconv.ParseInt(submitTimeStr, 10, 64)
		if err != nil {
			log.Println(err)
		}
		submitTime := time.Unix(submitTimeUnix, 0)
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			log.Println(err)
		}
		oneOne := OneOne{
			Answer:     "",
			OnChain:    0,
			Status:     valueInt,
			SubmitTime: submitTime,
		}
		OneOneList = append(OneOneList, oneOne)
	}
	// 查询数据库中的提交信息
	submitHistoryArray, err := model.NewSubmitHistoryFlow().QueryThisExerciseUserSubmitHistory(userID, userType, exerciseID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, OneOneResponse{
			List:       nil,
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	for _, submitHistory := range submitHistoryArray {
		oneOne := OneOne{
			Answer:     submitHistory.StudentAnswer,
			Status:     submitHistory.Status,
			OnChain:    submitHistory.OnChain,
			SubmitTime: submitHistory.SubmitTime,
		}
		OneOneList = append(OneOneList, oneOne)
	}
	context.JSON(http.StatusOK, OneOneResponse{
		List:       OneOneList,
		StatusCode: 0,
		StatusMsg:  "",
	})
}
