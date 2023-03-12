package exercise

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/model"
)

type AllExerciseResponse struct {
	List []AllExercise `json:"list"`
	common.Response
}

type AllExercise struct {
	ExerciseID    int64  `json:"exercise_id"`
	ExerciseName  string `json:"exercise_name"`
	Grade         int    `json:"grade"`
	PassCount     int    `json:"pass_count"`
	PublisherName string `json:"publisher_name"`
	PublisherType int64  `json:"publisher_type"`
	SubmitCount   int    `json:"submit_count"`
}

// GetAllExerciseHandle 获取题库中所有可见的题目条目
func GetAllExerciseHandle(context *gin.Context) {
	exerciseContentArray, err := model.NewExerciseContentFlow().GetAllVisitableExercise()
	if err != nil {
		context.JSON(http.StatusInternalServerError, AllExerciseResponse{
			List:     nil,
			Response: common.NewCommonResponse(1, "查询题库出错"),
		})
		return
	}
	var AllExerciseList []AllExercise
	for _, exerciseContent := range exerciseContentArray {
		publisherID := exerciseContent.PublisherID
		publisherType := exerciseContent.PublisherType
		publisherName := common.QueryUsername(publisherID, publisherType)
		allExercise := AllExercise{
			ExerciseID:    exerciseContent.ID,
			ExerciseName:  exerciseContent.Name,
			Grade:         exerciseContent.Grade,
			PassCount:     exerciseContent.PassCount,
			PublisherName: publisherName,
			PublisherType: publisherType,
			SubmitCount:   exerciseContent.SubmitCount,
		}
		AllExerciseList = append(AllExerciseList, allExercise)
	}
	context.JSON(http.StatusOK, AllExerciseResponse{
		List:     AllExerciseList,
		Response: common.NewCommonResponse(0, ""),
	})
}
