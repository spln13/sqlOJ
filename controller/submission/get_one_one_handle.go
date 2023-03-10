package submission

import "github.com/gin-gonic/gin"

type OneOneResponse struct {
	List       []OneOne `json:"list"`
	StatusCode int64    `json:"status_code"`
	StatusMsg  string   `json:"status_msg"`
}

type OneOne struct {
	Answer     string `json:"answer"`
	Status     int64  `json:"status"`
	SubmitTime string `json:"submit_time"`
}

func GetOneOneHandle(context *gin.Context) {

}
