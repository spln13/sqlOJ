package exercise

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
)

type AllTableResponse struct {
	List []TableInfo `json:"list"`
	utils.Response
}

type TableInfo struct {
	AssociationCount int    `json:"association_count"`
	Description      string `json:"description"`
	TableID          int64  `json:"table_id"`
	TableName        string `json:"table_name"`
}

func GetAllTableHandle(context *gin.Context) {
	var tableInfoList []TableInfo
	tableList, err := model.NewExerciseTableFlow().QueryAllTable()
	if err != nil {
		context.JSON(http.StatusInternalServerError, AllTableResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	for _, table := range tableList {
		tableInfo := TableInfo{
			AssociationCount: table.AssociationCount,
			Description:      table.Description,
			TableID:          table.ID,
			TableName:        table.Name,
		}
		tableInfoList = append(tableInfoList, tableInfo)
	}
	context.JSON(http.StatusOK, AllTableResponse{
		List:     tableInfoList,
		Response: utils.NewCommonResponse(0, ""),
	})
}
