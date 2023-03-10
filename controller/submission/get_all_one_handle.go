package submission

import "github.com/gin-gonic/gin"

type AllOneResponse struct {
	List       []AllOneList `json:"list"`
	StatusCode int64        `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
}

type AllOneList struct {
	Answer     string `json:"answer"`
	Status     int64  `json:"status"`
	SubmitTime string `json:"submit_time"`
	UserID     int64  `json:"user_id"`
	UserType   int64  `json:"user_type"`
	Username   string `json:"username"`
}

func GetAllOneHandle(context *gin.Context) {}
