package student_account

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/fabric"
	"sqlOJ/model"
	"sqlOJ/utils"
)

type StudentScoreResponse struct {
	List []StudentScoreInfo `json:"list"`
	utils.Response
}

type StudentScoreInfo struct {
	ContestScore  float64 `json:"contest_score"`
	ExerciseScore float64 `json:"exercise_score"`
	Number        int64   `json:"number"`
	RealName      string  `json:"real_name"`
	Score         float64 `json:"score"`
}

func StudentScoreHandle(context *gin.Context) {
	studentScoreArray, err := fabric.RatingStudents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, StudentScoreResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var studentScoreInfoList []StudentScoreInfo
	for _, studentScore := range studentScoreArray {
		realName, err := model.NewStudentAccountFlow().QueryStudentRealNameByNumber(studentScore.Number)
		if err != nil {
			context.JSON(http.StatusInternalServerError, StudentScoreResponse{
				List:     nil,
				Response: utils.NewCommonResponse(1, err.Error()),
			})
			return
		}
		studentScoreInfo := StudentScoreInfo{
			ContestScore:  studentScore.ContestScore,
			ExerciseScore: studentScore.ExerciseScore,
			Number:        studentScore.Number,
			RealName:      realName,
			Score:         studentScore.Score,
		}
		studentScoreInfoList = append(studentScoreInfoList, studentScoreInfo)
	}
	context.JSON(http.StatusOK, StudentScoreResponse{
		List:     studentScoreInfoList,
		Response: utils.NewCommonResponse(0, ""),
	})
}
