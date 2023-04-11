package contest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/common"
	"sqlOJ/model"
	"time"
)

type Response struct {
	List []SignalContest `json:"list"`
	common.Response
}

type SignalContest struct {
	ContestID     int64     `json:"contest_id"`
	ContestName   string    `json:"contest_name"`
	PublisherName string    `json:"publisher_name"`
	PublisherType int64     `json:"publisher_type"`
	BeginAt       time.Time `json:"begin_at"`
	EndAt         time.Time `json:"end_at"`
}

func GetAllContestHandle(context *gin.Context) {
	contestDAOList, err := model.NewContestFlow().GetAllContest()
	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{
			List:     nil,
			Response: common.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var contestList []SignalContest
	for _, signalContest := range contestDAOList {
		signalContest := SignalContest{
			BeginAt:       signalContest.BeginAt,
			ContestID:     signalContest.ID,
			ContestName:   signalContest.Name,
			EndAt:         signalContest.EndAt,
			PublisherName: signalContest.PublisherName,
			PublisherType: signalContest.PublisherType,
		}
		contestList = append(contestList, signalContest)
	}
	context.JSON(http.StatusOK, Response{
		List:     contestList,
		Response: common.NewCommonResponse(0, ""),
	})
}
