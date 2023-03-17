package exercise

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/model"
)

type AllExerciseResponseWithoutToken struct {
	List []AllExerciseWithoutToken `json:"list"`
	common.Response
}

type AllExerciseResponseWithToken struct {
	List []AllExerciseWithToken `json:"list"`
	common.Response
}

type AllExerciseWithoutToken struct {
	ExerciseID    int64  `json:"exercise_id"`
	ExerciseName  string `json:"exercise_name"`
	Grade         int    `json:"grade"`
	PassCount     int    `json:"pass_count"`
	PublisherName string `json:"publisher_name"`
	PublisherType int64  `json:"publisher_type"`
	SubmitCount   int    `json:"submit_count"`
}

type AllExerciseWithToken struct {
	ExerciseID    int64  `json:"exercise_id"`
	ExerciseName  string `json:"exercise_name"`
	Grade         int    `json:"grade"`
	PassCount     int    `json:"pass_count"`
	PublisherName string `json:"publisher_name"`
	PublisherType int64  `json:"publisher_type"`
	SubmitCount   int    `json:"submit_count"`
	Status        int    `json:"status"`
}

// GetAllExerciseWithoutTokenHandle 获取题库中所有可见的题目条目, 未登录状态时请求
func GetAllExerciseWithoutTokenHandle(context *gin.Context) {
	exerciseContentArray, err := model.NewExerciseContentFlow().GetAllVisitableExercise()
	if err != nil {
		context.JSON(http.StatusInternalServerError, AllExerciseResponseWithoutToken{
			List:     nil,
			Response: common.NewCommonResponse(1, "查询题库出错"),
		})
		return
	}
	var AllExerciseList []AllExerciseWithoutToken
	for _, exerciseContent := range exerciseContentArray {
		publisherType := exerciseContent.PublisherType
		allExercise := AllExerciseWithoutToken{
			ExerciseID:    exerciseContent.ID,
			ExerciseName:  exerciseContent.Name,
			Grade:         exerciseContent.Grade,
			PassCount:     exerciseContent.PassCount,
			PublisherName: exerciseContent.PublisherName,
			PublisherType: publisherType,
			SubmitCount:   exerciseContent.SubmitCount,
		}
		AllExerciseList = append(AllExerciseList, allExercise)
	}
	context.JSON(http.StatusOK, AllExerciseResponseWithoutToken{
		List:     AllExerciseList,
		Response: common.NewCommonResponse(0, ""),
	})
}

// GetAllExerciseWithTokenHandle 登录用户获取所有题目信息
func GetAllExerciseWithTokenHandle(context *gin.Context) {
	userID, ok1 := context.MustGet("user_id").(int64)
	userType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusInternalServerError, AllExerciseResponseWithToken{
			Response: common.NewCommonResponse(1, "解析用户参数错误"),
		})
		return
	}
	exerciseContentArray, err := model.NewExerciseContentFlow().GetAllVisitableExercise()
	if err != nil {
		context.JSON(http.StatusInternalServerError, AllExerciseResponseWithoutToken{
			List:     nil,
			Response: common.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var AllExerciseWithTokenList []AllExerciseWithToken
	problemStatusMap, err := model.NewUserProblemStatusFlow().QueryUserProblemStatus(userID, userType)
	if err != nil {
		context.JSON(http.StatusInternalServerError, AllExerciseResponseWithToken{
			List:     nil,
			Response: common.NewCommonResponse(1, err.Error()),
		})
		return
	}
	for _, exerciseContent := range exerciseContentArray {
		exerciseID := exerciseContent.ID
		status := problemStatusMap[exerciseID]
		allExerciseWithToken := AllExerciseWithToken{
			ExerciseID:    exerciseID,
			ExerciseName:  exerciseContent.Name,
			Grade:         exerciseContent.Grade,
			PassCount:     exerciseContent.PassCount,
			PublisherName: exerciseContent.PublisherName,
			PublisherType: exerciseContent.PublisherType,
			SubmitCount:   exerciseContent.SubmitCount,
			Status:        status,
		}
		AllExerciseWithTokenList = append(AllExerciseWithTokenList, allExerciseWithToken)
	}
	context.JSON(http.StatusOK, AllExerciseResponseWithToken{
		List:     AllExerciseWithTokenList,
		Response: common.NewCommonResponse(0, ""),
	})
}
