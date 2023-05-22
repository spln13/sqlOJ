package submission

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
	"strconv"
)

type ContestSubmissionDetailResponse struct {
	Answer string `json:"answer"`
	utils.Response
}

func ContestGetDetailHandler(context *gin.Context) {
	userID, ok1 := context.MustGet("user_id").(int64)
	userType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusBadRequest, ContestSubmissionDetailResponse{
			Response: utils.NewCommonResponse(1, "解析用户信息错误"),
		})
		return
	}
	submissionIDStr := context.Query("submission_id")
	submissionID, err := strconv.ParseInt(submissionIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, ContestSubmissionDetailResponse{
			Response: utils.NewCommonResponse(1, "请求参数错误"),
		})
		return
	}
	answer, err := model.NewContestSubmissionFlow().QueryOneSubmissionAnswer(userID, userType, submissionID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, ContestSubmissionDetailResponse{
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	context.JSON(http.StatusOK, ContestSubmissionDetailResponse{
		Answer:   answer,
		Response: utils.NewCommonResponse(0, ""),
	})
}
