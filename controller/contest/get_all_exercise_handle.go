package contest

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/model"
	"strconv"
)

type AllExerciseResponse struct {
	List []ExerciseInfo `json:"list"`
	common.Response
}

type ExerciseInfo struct {
	ExerciseID    int64  `json:"exercise_id"`
	ExerciseName  string `json:"exercise_name"`
	Grade         int    `json:"grade"`
	PassCount     int    `json:"pass_count"`
	PublisherName string `json:"publisher_name"`
	PublisherType int64  `json:"publisher_type"`
	Status        int    `json:"status"`
	SubmitCount   int    `json:"submit_count"`
}

func GetAllExerciseHandle(context *gin.Context) {
	userID, ok1 := context.MustGet("user_id").(int64)
	userType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusInternalServerError, AllExerciseResponse{
			List:     nil,
			Response: common.NewCommonResponse(1, "解析用户信息错误"),
		})
		return
	}
	contestIDStr := context.Query("contest_id")
	contestID, err := strconv.ParseInt(contestIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, AllExerciseResponse{
			List:     nil,
			Response: common.NewCommonResponse(1, "请求参数错误"),
		})
		return
	}
	exerciseIDList, err := model.NewContestExerciseAssociationFlow().GetExerciseIDListByContestID(contestID)
	if err != nil {
		context.JSON(http.StatusBadRequest, AllExerciseResponse{
			List:     nil,
			Response: common.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var exerciseInfoList []ExerciseInfo
	for _, exerciseID := range exerciseIDList {
		oneExercise, _ := model.NewExerciseContentFlow().GetOneExercise(exerciseID)
		problemStatusMap, _ := model.NewContestExerciseStatusFlow().QueryContestExerciseStatus(userID, userType, contestID)
		exerciseInfo := ExerciseInfo{
			ExerciseID:    oneExercise.ID,
			ExerciseName:  oneExercise.Name,
			Grade:         oneExercise.Grade,
			PassCount:     oneExercise.PassCount,
			PublisherName: oneExercise.PublisherName,
			PublisherType: oneExercise.PublisherType,
			SubmitCount:   oneExercise.SubmitCount,
			Status:        problemStatusMap[oneExercise.ID],
		}
		exerciseInfoList = append(exerciseInfoList, exerciseInfo)
	}
	context.JSON(http.StatusOK, AllExerciseResponse{
		List:     exerciseInfoList,
		Response: common.NewCommonResponse(0, ""),
	})
}
