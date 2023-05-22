package admin_account

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sqlOJ/model"
)

func AdminAddDebugHandler(context *gin.Context) {
	username := context.Query("username")
	password, _ := context.MustGet("password_sha256").(string)
	if _, err := model.NewAdminAccountFlow().InsertAdminAccount(username, password); err != nil {
		log.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 1,
			"status_msg":  err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "",
	})
}
