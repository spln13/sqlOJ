package contest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
	"strconv"
)

type OneContestStatusResponse struct {
	StatusList  [][]string `json:"status_list"`
	ProblemList []int64    `json:"problem_list"`
	utils.Response
}

func GetContestStatusHandle(context *gin.Context) {
	contestIDStr := context.Query("contest_id")
	contestID, err := strconv.ParseInt(contestIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, OneContestStatusResponse{
			StatusList:  nil,
			ProblemList: nil,
			Response:    utils.NewCommonResponse(1, "请求参数错误"),
		})
		return
	}
	// 首先使用contestID查询该场竞赛有多少道题目
	problemIDList, err := model.NewContestExerciseAssociationFlow().GetExerciseIDListByContestID(contestID)
	if err != nil {
		context.JSON(http.StatusBadRequest, OneContestStatusResponse{
			StatusList:  nil,
			ProblemList: nil,
			Response:    utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	// 获取竞赛状态中提交过学生学号, 忽略提交状态表中userType > 1的记录
	// 1. 首先查询参与竞赛的学生ID与学号的对应表map[int64]int64
	// 2. 遍历学生id, 获取当前学生题目id和状态的对应表map[int64]int64
	// 3. append到返回list中
	studentIDList, err := model.NewContestExerciseStatusFlow().QueryStudentIDListByContestID(contestID)
	if err != nil {
		context.JSON(http.StatusBadRequest, OneContestStatusResponse{
			StatusList:  nil,
			ProblemList: nil,
			Response:    utils.NewCommonResponse(1, err.Error()),
		})
		return
	}

	// 查询studentID对应的学号
	studentIDNumberMap, err := model.NewStudentAccountFlow().QueryStudentIDNumberMap(studentIDList)
	if err != nil {
		context.JSON(http.StatusBadRequest, OneContestStatusResponse{
			StatusList:  nil,
			ProblemList: nil,
			Response:    utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var statusList [][]string
	for _, studentID := range studentIDList {
		studentProblemStatusMap, err := model.NewContestExerciseStatusFlow().QueryStudentProblemStatusMap(studentID, contestID)
		if err != nil {
			context.JSON(http.StatusBadRequest, OneContestStatusResponse{
				StatusList:  nil,
				ProblemList: nil,
				Response:    utils.NewCommonResponse(1, err.Error()),
			})
			return
		}
		var oneStudentList []string
		studentNumber := studentIDNumberMap[studentID]
		studentNumberString := strconv.FormatInt(studentNumber, 10)
		oneStudentList = append(oneStudentList, studentNumberString)
		for _, problemID := range problemIDList {
			status := studentProblemStatusMap[problemID]
			statusString := strconv.Itoa(status)
			oneStudentList = append(oneStudentList, statusString)
		}
		statusList = append(statusList, oneStudentList)
	}
	context.JSON(http.StatusOK, OneContestStatusResponse{
		StatusList:  statusList,
		ProblemList: problemIDList,
		Response:    utils.NewCommonResponse(0, ""),
	})
}
