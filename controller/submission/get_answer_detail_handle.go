package submission

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"strconv"
)

type AnswerDetailResponse struct {
	Answer     string `json:"answer"`
	StatusCode int64  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func GetAnswerDetailHandler(context *gin.Context) {
	userID, ok1 := context.MustGet("user_id").(int64)
	userType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusBadRequest, AnswerDetailResponse{
			Answer:     "",
			StatusCode: 1,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	submissionIDStr := context.Query("submission_id")
	submissionID, err := strconv.ParseInt(submissionIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, AnswerDetailResponse{
			Answer:     "",
			StatusCode: 1,
			StatusMsg:  "请求参数错误",
		})
		return
	}
	queryUserID, queryUserType, answer, err := model.NewSubmitHistoryFlow().QuerySubmissionAnswer(submissionID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, AnswerDetailResponse{
			Answer:     "",
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	if queryUserID != userID || queryUserType != userType {
		context.JSON(http.StatusOK, AnswerDetailResponse{
			Answer:     "",
			StatusCode: 1,
			StatusMsg:  "无权访问",
		})
		return
	}
	context.JSON(http.StatusOK, AnswerDetailResponse{
		Answer:     answer,
		StatusCode: 0,
		StatusMsg:  "",
	})
}
