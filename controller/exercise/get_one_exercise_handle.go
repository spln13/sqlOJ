package exercise

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/model"
	"strconv"
)

type OneExercise struct {
	Description   string `json:"description"`
	Grade         int    `json:"grade"`
	Name          string `json:"name"`
	PassCount     int    `json:"pass_count"`
	PublisherName string `json:"publisher_name"`
	PublisherType int64  `json:"publisher_type"`
	SubmitCount   int    `json:"submit_count"`
}

type OneExerciseResponse struct {
	OneExercise
	utils.Response
}

// GetOneExerciseHandle 获取当前题目的题目信息
func GetOneExerciseHandle(context *gin.Context) {
	exerciseIDStr := context.Query("exercise_id")
	exerciseID, err := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, OneExerciseResponse{
			OneExercise: OneExercise{},
			Response:    utils.NewCommonResponse(1, "解析题目信息错误"),
		})
		return
	}
	exerciseContent, err := model.NewExerciseContentFlow().GetOneExercise(exerciseID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, OneExerciseResponse{
			OneExercise: OneExercise{},
			Response:    utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	publisherID := exerciseContent.PublisherID
	publisherType := exerciseContent.PublisherType
	publisherName := utils.QueryUsername(publisherID, publisherType)
	oneExercise := OneExercise{
		Description:   exerciseContent.Description,
		Grade:         exerciseContent.Grade,
		Name:          exerciseContent.Name,
		PassCount:     exerciseContent.PassCount,
		PublisherName: publisherName,
		PublisherType: exerciseContent.PublisherType,
		SubmitCount:   exerciseContent.SubmitCount,
	}
	context.JSON(http.StatusOK, OneExerciseResponse{
		OneExercise: oneExercise,
		Response:    utils.NewCommonResponse(0, ""),
	})
}
