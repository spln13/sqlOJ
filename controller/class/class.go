package class

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/model"
	"strconv"
)

func CreateClassHandle(context *gin.Context) {
	userID, ok := context.MustGet("user_id").(int64)
	className := context.Query("name")
	teacherUsername := context.Query("username")
	if !ok {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "解析用户信息错误"))
		return
	}
	// 根据teacherUsername即教职工号查询对应教职工真实姓名
	realName, err := model.NewTeacherAccountFlow().QueryTeacherRealNameByUsername(teacherUsername)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	if err := model.NewClassFlow().InsertClass(className, teacherUsername, realName, userID); err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	context.JSON(http.StatusOK, common.NewCommonResponse(0, ""))
}

func AddStudentToClassHandle(context *gin.Context) {
	userID, ok := context.MustGet("user_id").(int64)
	if !ok {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "解析用户信息错误"))
		return
	}
	studentIDList := context.PostForm("student_id_list")
	classIDStr := context.PostForm("class_id")
	classID, err := strconv.ParseInt(classIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "获取classID"))
		return
	}
	className, err := model.NewClassFlow().QueryClassNameByClassID(classID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}
	// 将班级名插入学生信息表中
}
