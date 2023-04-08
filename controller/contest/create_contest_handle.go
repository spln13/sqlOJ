package contest

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/cache"
	"sqlOJ/common"
	"sqlOJ/model"
	"time"
)

type CreateContestData struct {
	ContestName    string    `json:"contest_name"`
	BeginAt        time.Time `json:"begin_at"`
	EndAt          time.Time `json:"end_at"`
	ExerciseIDList []int64   `json:"exercise_id_list"`
	ClassIDList    []int64   `json:"class_id_list"`
}

func CreateContestHandle(context *gin.Context) {
	publisherID, ok1 := context.MustGet("user_id").(int64) // 获取又JWT设置的user_id
	publisherType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "解析用户token错误"))
		return
	}
	var createContestData CreateContestData
	if err := context.ShouldBindJSON(&createContestData); err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "解析请求参数错误"))
		return
	}
	publisherName := common.QueryUsername(publisherID, publisherType)
	// 在数据库中插入竞赛信息
	contestID, err := model.NewContestFlow().CreateContest(createContestData.ContestName, publisherName, publisherID, publisherType, createContestData.BeginAt, createContestData.EndAt)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	// 在数据库中插入竞赛和练习题引用关系
	err = model.NewContestExerciseAssociationFlow().InsertContestExerciseAssociation(contestID, createContestData.ExerciseIDList)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	// 在数据库中插入竞赛做题的班级
	err = model.NewContestClassAssociationFlow().InsertContestClassAssociation(contestID, createContestData.ClassIDList)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	// 查询班级中包括的学生ID
	studentIDList, err := model.NewStudentAccountFlow().QueryStudentIDByClassID(createContestData.ClassIDList)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	// 缓存竞赛禁止学生访问名单
	if err := cache.ContestForbidStudentCache(contestID, studentIDList, createContestData.BeginAt, createContestData.EndAt); err != nil {
		return
	}
	// TODO:需要更新权限设置方案
}
