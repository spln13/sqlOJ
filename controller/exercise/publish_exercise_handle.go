package exercise

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xwb1989/sqlparser"
	"log"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/model"
	"time"
)

type PublishExerciseData struct {
	Answer      string  `json:"answer"`
	Description string  `json:"description"`
	Grade       int     `json:"grade"`
	Name        string  `json:"name"`
	ShowAt      string  `json:"show_at"`
	TableIDList []int64 `json:"table_id_list"`
	Visitable   int     `json:"visitable"`
}

// PublishExerciseHandle 完成发布题目功能
func PublishExerciseHandle(context *gin.Context) {
	publisherID, ok1 := context.MustGet("user_id").(int64) // 获取又JWT设置的user_id
	publisherType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "解析用户token错误"))
		return
	}
	var publishExerciseData PublishExerciseData
	if err := context.ShouldBindJSON(&publishExerciseData); err != nil {
		context.JSON(http.StatusBadRequest, common.NewCommonResponse(1, "请求参数错"))
		return
	}
	answer := publishExerciseData.Answer
	description := publishExerciseData.Description
	showAtStr := publishExerciseData.ShowAt
	tableIDList := publishExerciseData.TableIDList
	visitable := publishExerciseData.Visitable
	name := publishExerciseData.Name
	grade := publishExerciseData.Grade
	// 时间格式
	layout := "2006-01-02 15:04:05"

	showAt, err := time.Parse(layout, showAtStr)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "解析公布时间错误"))
		return
	}
	exeType, err := parseAnswer(answer)
	if err != nil {
		context.JSON(http.StatusOK, common.NewCommonResponse(1, err.Error()))
		return
	}
	if exeType == 0 {
		context.JSON(http.StatusOK, common.NewCommonResponse(1, "答案不合法"))
		return
	}
	exerciseID, err := model.NewExerciseContentFlow().InsertExerciseContent(publisherID, publisherType, name, answer, description, exeType, grade, visitable, showAt)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	if err := model.NewExerciseAssociationFlow().InsertExerciseAssociation(exerciseID, tableIDList); err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	// 自增练习表数据库中的关联数
	if err := model.NewExerciseTableFlow().IncreaseExerciseTableAssociationCount(tableIDList); err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "关联数据表错误"))
		return
	}
	context.JSON(http.StatusOK, common.NewCommonResponse(0, ""))
	return
}

// parseAnswer 处理答案sql语句，判断其是否有语法错误, 并返回其类型
func parseAnswer(answer string) (int, error) {
	answerStmt, err := sqlparser.Parse(answer)
	if err != nil {
		log.Println(err)
		return 0, errors.New("答案语法错误")
	}
	var code int
	switch stmt := answerStmt.(type) {
	case *sqlparser.Select:
		code = 1
		_ = stmt
	case *sqlparser.Insert:
		code = 2
	case *sqlparser.Update:
		code = 3
	case *sqlparser.Delete:
		code = 4
	default:
		code = 0
	}
	return code, nil
}
