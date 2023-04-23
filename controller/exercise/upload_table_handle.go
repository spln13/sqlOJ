package exercise

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"sqlOJ/model"
)

func UploadTableHandle(context *gin.Context) {
	publisherID, ok1 := context.MustGet("user_id").(int64)
	publisherType, ok2 := context.MustGet("user_type").(int64)
	if !ok1 || !ok2 {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "解析用户id错误"))
		return
	}
	name := context.PostForm("name")
	// 名字查重
	isExist, err := model.NewExerciseTableFlow().QueryExerciseTableExist(name)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if isExist {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "表名已重复"))
		return
	}
	description := context.PostForm("description")
	sqlFile, err := context.FormFile("sql_file")
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "获取文件错误"))
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
