package student_account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"log"
	"net/http"
	"os"
	"sqlOJ/fabric"
	"sqlOJ/model"
	"sqlOJ/utils"
)

type StudentScoreResponse struct {
	List []StudentScoreInfo `json:"list"`
	utils.Response
}

type StudentScoreInfo struct {
	ContestScore  float64 `json:"contest_score"`
	ExerciseScore float64 `json:"exercise_score"`
	Number        int64   `json:"number"`
	RealName      string  `json:"real_name"`
	Score         float64 `json:"score"`
}

func StudentScoreHandle(context *gin.Context) {
	studentScoreArray, err := fabric.RatingStudents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, StudentScoreResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var studentScoreInfoList []StudentScoreInfo
	for _, studentScore := range studentScoreArray {
		realName, err := model.NewStudentAccountFlow().QueryStudentRealNameByNumber(studentScore.Number)
		if err != nil {
			context.JSON(http.StatusInternalServerError, StudentScoreResponse{
				List:     nil,
				Response: utils.NewCommonResponse(1, err.Error()),
			})
			return
		}
		studentScoreInfo := StudentScoreInfo{
			ContestScore:  studentScore.ContestScore,
			ExerciseScore: studentScore.ExerciseScore,
			Number:        studentScore.Number,
			RealName:      realName,
			Score:         studentScore.Score,
		}
		studentScoreInfoList = append(studentScoreInfoList, studentScoreInfo)
	}
	context.JSON(http.StatusOK, StudentScoreResponse{
		List:     studentScoreInfoList,
		Response: utils.NewCommonResponse(0, ""),
	})
}

func StudentScoreDownloadHandle(context *gin.Context) {
	studentScoreArray, err := fabric.RatingStudents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, StudentScoreResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	file := excelize.NewFile()
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	sheetName := "Sheet1"
	//headers := []string{"学号", "姓名", "竞赛得分", "题库得分", "总分"}
	_ = file.SetCellValue(sheetName, "A1", "学号")
	_ = file.SetCellValue(sheetName, "B1", "姓名")
	_ = file.SetCellValue(sheetName, "C1", "竞赛得分")
	_ = file.SetCellValue(sheetName, "D1", "题库得分")
	_ = file.SetCellValue(sheetName, "E1", "总分")
	// 从数据源中填充数据
	var studentScoreInfoList []StudentScoreInfo
	for _, studentScore := range studentScoreArray {
		realName, err := model.NewStudentAccountFlow().QueryStudentRealNameByNumber(studentScore.Number)
		if err != nil {
			context.JSON(http.StatusInternalServerError, StudentScoreResponse{
				List:     nil,
				Response: utils.NewCommonResponse(1, err.Error()),
			})
			return
		}
		studentScoreInfo := StudentScoreInfo{
			ContestScore:  studentScore.ContestScore,
			ExerciseScore: studentScore.ExerciseScore,
			Number:        studentScore.Number,
			RealName:      realName,
			Score:         studentScore.Score,
		}
		studentScoreInfoList = append(studentScoreInfoList, studentScoreInfo)
	}
	for i, scoreInfo := range studentScoreInfoList {
		rowIndex := i + 2
		_ = file.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), scoreInfo.Number)
		_ = file.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), scoreInfo.RealName)
		_ = file.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), scoreInfo.ContestScore)
		_ = file.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), scoreInfo.ExerciseScore)
		_ = file.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), scoreInfo.Score)
	}
	fileName := "./temp_excel_files/temp.xlsx"
	if err := file.SaveAs(fileName); err != nil {
		context.String(http.StatusInternalServerError, fmt.Sprintf("保存Excel文件失败：%v", err))
		return
	}
	// 设置响应头
	context.Header("Content-Description", "File Transfer")
	context.Header("Content-Disposition", "attachment; filename=学生成绩.xlsx")
	context.Header("Content-Type", "application/octet-stream")
	context.Header("Content-Transfer-Encoding", "binary")
	context.Header("Expires", "0")
	context.Header("Cache-Control", "must-revalidate")
	context.Header("Pragma", "public")

	// 发送文件给客户端
	context.File(fileName)
	// 删除临时文件
	if err := os.Remove(fileName); err != nil {
		log.Println("删除临时文件失败：", err)
	}
}
