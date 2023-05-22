package contest

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/cache"
	"sqlOJ/model"
	"sqlOJ/utils"
	"strconv"
	"time"
)

type CreateContestData struct {
	ContestName    string  `json:"contest_name"`
	BeginAt        string  `json:"begin_at"`
	EndAt          string  `json:"end_at"`
	ExerciseIDList []int64 `json:"exercise_id_list"`
	ClassIDList    []int64 `json:"class_id_list"`
}

func CreateContestHandler(context *gin.Context) {
	publisherID, ok1 := context.MustGet("user_id").(int64) // 获取又JWT设置的user_id
	publisherType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "解析用户token错误"))
		return
	}
	var createContestData CreateContestData
	if err := context.ShouldBindJSON(&createContestData); err != nil { // 解析请求参数
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "解析请求参数错误"))
		return
	}
	beginAt, err1 := time.Parse(time.RFC3339, createContestData.BeginAt)
	endAt, err2 := time.Parse(time.RFC3339, createContestData.EndAt)
	if err1 != nil || err2 != nil {
		log.Println(err1)
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(1, "解析请求参数错误"))
		return
	}
	publisherName := utils.QueryUsername(publisherID, publisherType)
	// 在数据库中插入竞赛信息
	contestID, err := model.NewContestFlow().CreateContest(createContestData.ContestName, publisherName, publisherID, publisherType, beginAt, endAt)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	// 在数据库中插入竞赛和练习题引用关系
	err = model.NewContestExerciseAssociationFlow().InsertContestExerciseAssociation(contestID, createContestData.ExerciseIDList)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	// 在数据库中插入竞赛做题的班级
	err = model.NewContestClassAssociationFlow().InsertContestClassAssociation(contestID, createContestData.ClassIDList)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	// 查询班级中包括的学生ID
	studentIDList, err := model.NewStudentAccountFlow().QueryStudentIDByClassID(createContestData.ClassIDList)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	// 缓存竞赛参与学生访问名单
	if err := cache.ContestStudentCache(contestID, studentIDList, beginAt, endAt); err != nil {
		return
	}
	context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(0, ""))
}

// DeleteContestHandler 删除竞赛
func DeleteContestHandler(context *gin.Context) {
	contestIDStr := context.Query("contest_id")
	contestID, err := strconv.ParseInt(contestIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(1, "请求参数错"))
		return
	}
	// 删除竞赛时需要删除所有跟此竞赛相关的数据表
	if err := model.NewContestExerciseAssociationFlow().DeleteContestExerciseAssociation(contestID); err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if err := model.NewContestClassAssociationFlow().DeleteContestClassAssociation(contestID); err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if err := model.NewContestExerciseStatusFlow().DeleteContestExerciseStatus(contestID); err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
	}
	if err := model.NewContestFlow().DeleteContest(contestID); err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
	return
}
