package submission

import "github.com/gin-gonic/gin"

type OneAllResponse struct {
	List       []OneAllList `json:"list"`
	StatusCode int64        `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
}

type OneAllList struct {
	Answer       string `json:"answer"`
	ExerciseID   string `json:"exercise_id"`
	ExerciseName string `json:"exercise_name"`
	Status       string `json:"status"`
	SubmitTime   string `json:"submit_time"`
}

func GetOneAllHandle(context *gin.Context) {}
