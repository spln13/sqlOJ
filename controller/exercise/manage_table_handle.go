package exercise

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"sqlOJ/model"
	"sqlOJ/utils"
	"strconv"
)

func UploadTableHandler(context *gin.Context) {
	publisherID, ok1 := context.MustGet("user_id").(int64)
	publisherType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "解析用户id错误"))
		return
	}
	name := context.PostForm("name")
	description := context.PostForm("description")
	sqlFile, err := context.FormFile("sql_file")
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "获取文件错误"))
	}
	isExist, err := model.NewExerciseTableFlow().QueryExerciseTableExist(name) // 查看名字是否重复
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if isExist {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "表名已重复"))
		return
	}
	filePath := "./temp_sql_files/" + name + ".sql"
	if err = context.SaveUploadedFile(sqlFile, filePath); err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "上传文件错误"))
		return
	}

	if err := model.ExecSqlCreateTable(filePath); err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	// 删除sql文件
	if err := os.Remove(filePath); err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "删除sql文件错误"))
		return
	}
	// 将记录插入数据库
	if err := model.NewExerciseTableFlow().InsertExerciseTable(publisherID, publisherType, name, description); err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "更新系统数据库错误"))
		return
	}
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
	return
}

func DeleteTableHandler(context *gin.Context) {
	tableIDStr := context.Query("table_id")
	tableID, err := strconv.ParseInt(tableIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(1, "请求参数错误"))
		return
	}
	exist, err := model.NewExerciseAssociationFlow().QueryAssociationExist(tableID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if exist { // 如果存在关联, 删除失败
		context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
		return
	}
	// 根据tableID获取数据表的名字
	tableName, err := model.NewExerciseTableFlow().QueryTableNameByID(tableID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	err = model.ExecSqlDeleteTable(tableName) // 根据tableName执行sql语句删除sql_exercise中的数据表
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
}
