package exercise

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
)

type AllExerciseResponse struct {
	List []AllExercise `json:"list"`
	utils.Response
}

type AllExercise struct {
	ExerciseID    int64  `json:"exercise_id"`
	ExerciseName  string `json:"exercise_name"`
	Grade         int    `json:"grade"`
	PassCount     int    `json:"pass_count"`
	PublisherName string `json:"publisher_name"`
	PublisherType int64  `json:"publisher_type"`
	SubmitCount   int    `json:"submit_count"`
	Status        int    `json:"status"`
}

// GetAllExerciseWithoutTokenHandler 获取题库中所有可见的题目条目, 未登录状态时请求
func GetAllExerciseWithoutTokenHandler(context *gin.Context) {
	exerciseContentArray, err := model.NewExerciseContentFlow().GetAllExercise()
	if err != nil {
		context.JSON(http.StatusInternalServerError, AllExerciseResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, "查询题库出错"),
		})
		return
	}
	var AllExerciseList []AllExercise
	for _, exerciseContent := range exerciseContentArray {
		publisherType := exerciseContent.PublisherType
		allExercise := AllExercise{
			ExerciseID:    exerciseContent.ID,
			ExerciseName:  exerciseContent.Name,
			Grade:         exerciseContent.Grade,
			PassCount:     exerciseContent.PassCount,
			PublisherName: exerciseContent.PublisherName,
			PublisherType: publisherType,
			SubmitCount:   exerciseContent.SubmitCount,
			Status:        0, // 未登录状态题目显示没有提交过
		}
		AllExerciseList = append(AllExerciseList, allExercise)
	}
	context.JSON(http.StatusOK, AllExerciseResponse{
		List:     AllExerciseList,
		Response: utils.NewCommonResponse(0, ""),
	})
}

// GetAllExerciseWithTokenHandler 登录用户获取所有题目信息
func GetAllExerciseWithTokenHandler(context *gin.Context) {
	userID, ok1 := context.MustGet("user_id").(int64)
	userType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusInternalServerError, AllExerciseResponse{
			Response: utils.NewCommonResponse(1, "解析用户参数错误"),
		})
		return
	}
	exerciseContentArray, err := model.NewExerciseContentFlow().GetAllExercise()
	if err != nil { // 获取所有可见的题目错误
		context.JSON(http.StatusInternalServerError, AllExerciseResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var AllExerciseWithTokenList []AllExercise
	problemStatusMap, err := model.NewUserProblemStatusFlow().QueryUserAllProblemStatus(userID, userType)
	// 获取用户所有提交过题目的状态
	if err != nil {
		context.JSON(http.StatusInternalServerError, AllExerciseResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	for _, exerciseContent := range exerciseContentArray {
		exerciseID := exerciseContent.ID
		status, ok := problemStatusMap[exerciseID]
		if !ok {
			status = 0 // 没有做过这一题
		}
		allExerciseWithToken := AllExercise{
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
	context.JSON(http.StatusOK, AllExerciseResponse{
		List:     AllExerciseWithTokenList,
		Response: utils.NewCommonResponse(0, ""),
	})
}
