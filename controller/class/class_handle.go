package class

import (
	"github.com/gin-gonic/gin"
	"log"
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
	studentIDStringArray := context.PostFormArray("student_id_list")
	var studentIDList []int64
	for _, val := range studentIDStringArray {
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Println(err)
		}
		studentIDList = append(studentIDList, intVal)
	}
	classIDStr := context.PostForm("class_id")
	classID, err := strconv.ParseInt(classIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, "获取classID"))
		return
	}
	// 将班级班级信息更新到学生表中
	if err := model.NewStudentAccountFlow().UpdateStudentsClass(classID, studentIDList); err != nil {
		context.JSON(http.StatusInternalServerError, common.NewCommonResponse(1, err.Error()))
		return
	}

}
