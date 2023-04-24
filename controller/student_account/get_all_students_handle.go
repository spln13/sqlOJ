package student_account

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
)

type AllStudentsResponse struct {
	List []StudentInfo `json:"list"`
	utils.Response
}

type StudentInfo struct {
	ClassID   int64  `json:"class_id"`
	ClassName string `json:"class_name"`
	Number    string `json:"number"`
	RealName  string `json:"real_name"`
	StudentID int64  `json:"student_id"`
	Username  string `json:"username"`
}

func GetAllStudentsHandle(context *gin.Context) {
	classIDNameMap, err := model.NewClassFlow().QueryClassIDNameMap()
	// classIDNameMap 是一个key为classID, value为className的hashMap
	if err != nil {
		context.JSON(http.StatusInternalServerError, AllStudentsResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	studentAccountDAOList, err := model.NewStudentAccountFlow().QueryAllStudent()
	if err != nil {
		context.JSON(http.StatusInternalServerError, AllStudentsResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var studentInfoList []StudentInfo
	for _, studentAccountDAO := range studentAccountDAOList {
		studentInfo := StudentInfo{
			ClassID:   studentAccountDAO.ClassID,
			ClassName: classIDNameMap[studentAccountDAO.ClassID],
			Number:    studentAccountDAO.Number,
			RealName:  studentAccountDAO.RealName,
			StudentID: studentAccountDAO.ID,
			Username:  studentAccountDAO.Username,
		}
		studentInfoList = append(studentInfoList, studentInfo)
	}
	context.JSON(http.StatusOK, AllStudentsResponse{
		List:     studentInfoList,
		Response: utils.NewCommonResponse(0, ""),
	})
}
