package exercise

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/xwb1989/sqlparser"
	"log"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
)

type PublishExerciseData struct {
	Answer      string  `json:"answer"`
	Description string  `json:"description"`
	Grade       int     `json:"grade"`
	Name        string  `json:"name"`
	TableIDList []int64 `json:"table_id_list"`
}

// PublishExerciseHandle 完成发布题目功能
func PublishExerciseHandle(context *gin.Context) {
	publisherID, ok1 := context.MustGet("user_id").(int64) // 获取由JWT设置的user_id
	publisherType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "解析用户token错误"))
		return
	}
	var publishExerciseData PublishExerciseData
	if err := context.ShouldBindJSON(&publishExerciseData); err != nil {
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(1, "请求参数错"))
		return
	}
	answer := publishExerciseData.Answer
	description := publishExerciseData.Description
	tableIDList := publishExerciseData.TableIDList
	name := publishExerciseData.Name
	grade := publishExerciseData.Grade
	publisherName := utils.QueryUsername(publisherID, publisherType)
	exeType, err := parseAnswer(answer)
	if err != nil {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if exeType == 0 {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "答案不合法"))
		return
	}
	exerciseID, err := model.NewExerciseContentFlow().InsertExerciseContent(publisherID, publisherType, name, answer, publisherName, description, exeType, grade)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if err := model.NewExerciseAssociationFlow().InsertExerciseAssociation(exerciseID, tableIDList); err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	// 自增练习表数据库中的关联数
	if err := model.NewExerciseTableFlow().IncreaseExerciseTableAssociationCount(tableIDList); err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "关联数据表错误"))
		return
	}
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
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
