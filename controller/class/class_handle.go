package class

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
	"strconv"
)

func CreateClassHandle(context *gin.Context) {
	className := context.Query("name")
	// 根据teacherUsername即教职工号查询对应教职工真实姓名
	if err := model.NewClassFlow().InsertClass(className); err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
}

func AddStudentToClassHandle(context *gin.Context) {
	studentIDStringArray := context.PostFormArray("student_id_list")
	var studentIDList []int64
	for _, val := range studentIDStringArray {
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Println(err)
			context.JSON(http.StatusBadRequest, utils.NewCommonResponse(1, "请求参数错误"))
			return
		}
		studentIDList = append(studentIDList, intVal)
	}
	classIDStr := context.PostForm("class_id")
	classID, err := strconv.ParseInt(classIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "获取classID"))
		return
	}
	studentCount := len(studentIDList)
	err = model.NewClassFlow().IncreaseStudentCountInClass(classID, studentCount)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	// 将班级班级信息更新到学生表中
	if err := model.NewStudentAccountFlow().UpdateStudentsClass(classID, studentIDList); err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}

	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
}
