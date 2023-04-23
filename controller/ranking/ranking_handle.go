package ranking

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/model"
	"sqlOJ/utils"
)

type Response struct {
	List []Item `json:"list"`
	utils.Response
}

type Item struct {
	Score    int64  `json:"score"`
	UserType int64  `json:"user_type"`
	Username string `json:"username"`
}

func GetRankingHandle(context *gin.Context) {
	rankingAPIList, err := model.NewScoreRecordFlow().GetRanking()
	if err != nil {
		context.JSON(http.StatusInternalServerError, Response{
			List:     nil,
			Response: utils.NewCommonResponse(1, err.Error()),
		})
		return
	}
	var itemList []Item
	for _, rankingAPI := range rankingAPIList {
		item := Item{
			Score:    rankingAPI.Score,
			UserType: rankingAPI.UserType,
			Username: rankingAPI.Username,
		}
		itemList = append(itemList, item)
	}
	context.JSON(http.StatusOK, Response{
		List:     itemList,
		Response: utils.NewCommonResponse(0, ""),
	})
}
