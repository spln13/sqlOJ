package exercise

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/model"
)

type PublishExerciseData struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Answer      string  `json:"answer"`
	Grade       int     `json:"grade"`
	TableIDList []int64 `json:"table_id_list"`
}

// PublishExerciseHandle 完成发布题目功能
func PublishExerciseHandle(context *gin.Context) {
	publisherID, ok := context.MustGet("user_id").(int64) // 获取又JWT设置的user_id
	if !ok {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "解析用户token错误"))
		return
	}
	jsonMap := context.MustGet("jsonMap").(*map[string]interface{})
	name := (*jsonMap)["name"].(string)
	answer := (*jsonMap)["answer"].(string)
	description := (*jsonMap)["description"].(string)
	grade := int((*jsonMap)["grade"].(float64))
	visitable := int((*jsonMap)["grade"].(float64))
	tableIDInterfaceList := (*jsonMap)["table_id_list"].([]interface{})
	var tableIDList []int64
	for _, tableID := range tableIDInterfaceList {
		id := int64(tableID.(float64))
		tableIDList = append(tableIDList, id)
	}
	exerciseID, err := model.NewExerciseContentFlow().InsertExerciseContent(publisherID, name, answer, description, grade, visitable)
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
