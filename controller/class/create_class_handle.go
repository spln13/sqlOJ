package class

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"os"
	"sqlOJ/model"
	"sqlOJ/utils"
)

type StudentInfo struct {
	Number string
	Name   string
}

func CreateClassHandle(context *gin.Context) {
	className := context.PostForm("name")
	classFile, err := context.FormFile("class_file")
	if err != nil {
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(1, "请求参数错"))
		return
	}
	// 查询名字是否重复
	valid, err := model.NewClassFlow().QueryClassNameValid(className)
	if err != nil {
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if !valid {
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(1, "班级名重复"))
		return
	}
	filePath := "./temp_excel_files/" + className + ".xlsx"               // excel文件保存路径
	if err := context.SaveUploadedFile(classFile, filePath); err != nil { // 将excel文件保存到本地
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "上传文件错误"))
		return
	}
	excelFile, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "解析文件错误"))
		return
	}
	defer func() { // defer关闭文件流
		if err := excelFile.Close(); err != nil {
			log.Println(err)
		}
	}()
	rows, err := excelFile.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "解析文件错误"))
		return
	}
	var studentInfoList []StudentInfo
	for _, row := range rows {
		numberStr := row[0]
		nameStr := row[1]
		studentInfo := StudentInfo{
			Number: numberStr,
			Name:   nameStr,
		}
		studentInfoList = append(studentInfoList, studentInfo)
	}
	if err := os.Remove(filePath); err != nil { // 删除excel文件
		log.Println(err)
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, "删除excel文件错误"))
		return
	}
	studentCount := len(studentInfoList)
	// 创建班级信息
	classID, err := model.NewClassFlow().CreateClass(className, studentCount)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}

	// 注册学生信息
	for _, studentInfo := range studentInfoList {
		password := passwordEncryption(studentInfo.Number)
		studentID, _ := model.NewStudentAccountFlow().CreateStudentAccount(studentInfo.Number, studentInfo.Name, className, password, classID)
		// 在ScoreRecord中插入学生记录
		_ = model.NewScoreRecordFlow().InsertScoreRecord(studentID, 1, studentInfo.Number)
	}
	context.JSON(http.StatusOK, utils.NewCommonResponse(0, ""))
}

// passwordEncryption 对明文密码使用SHA256进行加密
func passwordEncryption(password string) string {
	digest := sha256.New() // 对密码加密
	digest.Write([]byte(password))
	passwordSHA := hex.EncodeToString(digest.Sum(nil))
	return passwordSHA
}
