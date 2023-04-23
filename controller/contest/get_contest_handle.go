package contest

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/model"
	"strconv"
	"time"
)

type OneContestResponse struct {
	OneContest
	utils.Response
}

type OneContest struct {
	ContestID   int64     `json:"contest_id"`
	ContestName string    `json:"contest_name"`
	EndAt       time.Time `json:"end_at"`
	BeginAt     time.Time `json:"begin_at"`
}

func GetContestHandle(context *gin.Context) {
	contestIDStr := context.Query("contest_id")
	contestID, err := strconv.ParseInt(contestIDStr, 10, 64)
	if err != nil {
		log.Println(err)
		context.JSON(http.StatusBadRequest, OneContestResponse{
			Response: utils.NewCommonResponse(1, "请求参数错误"),
		})
		return
	}
	contestDAO, err := model.NewContestFlow().GetContestInfo(contestID)
	if err != nil {
		context.JSON(http.StatusBadRequest, OneContestResponse{
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	oneContest := OneContest{
		ContestID:   contestDAO.ID,
		ContestName: contestDAO.Name,
		EndAt:       contestDAO.EndAt,
		BeginAt:     contestDAO.BeginAt,
	}
	context.JSON(http.StatusOK, OneContestResponse{
		OneContest: oneContest,
		Response:   utils.NewCommonResponse(0, ""),
	})
}
