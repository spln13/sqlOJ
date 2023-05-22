package class

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
)

type GetAllClassResponse struct {
	List []OneClass `json:"list"`
	utils.Response
}

type OneClass struct {
	ClassID      int64  `json:"class_id"`
	ClassName    string `json:"class_name"`
	StudentCount int    `json:"student_count"`
}

func GetAllClassHandler(context *gin.Context) {
	classDAOList, err := model.NewClassFlow().GetAllClass()
	if err != nil {
		context.JSON(http.StatusInternalServerError, GetAllClassResponse{
			List:     nil,
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var classList []OneClass
	for _, classDAO := range classDAOList {
		class := OneClass{
			ClassID:      classDAO.ID,
			ClassName:    classDAO.Name,
			StudentCount: classDAO.StudentCount,
		}
		classList = append(classList, class)
	}
	context.JSON(http.StatusOK, GetAllClassResponse{
		List:     classList,
		Response: utils.NewCommonResponse(0, ""),
	})
}
