package student_account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"sqlOJ/cache"
	"sqlOJ/config"
	"sqlOJ/model"
	"strconv"
	"time"
)

func SendCodeHandle(context *gin.Context) {
	emailAddr := context.Query("email")
	ok, err := cache.CheckEmailCodeSendTimeValid(emailAddr)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if !ok {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "请求间隔少于1分钟"))
		return
	}
	// 查询此邮箱是否已经被注册过
	exist, err := model.NewStudentAccountFlow().QueryStudentExistByEmail(emailAddr)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(1, err.Error()))
		return
	}
	if exist {
		context.JSON(http.StatusOK, utils.NewCommonResponse(1, "邮箱已被注册"))
		return
	}
	go SendCode(emailAddr)
	context.JSON(200, utils.NewCommonResponse(0, ""))
}

func SendCode(emailAddr string) {

	// Generate a random 6-digit code
	rand.Seed(time.Now().UnixNano())
	min := 100000                      // 最小值
	max := 999999                      // 最大值
	code := rand.Intn(max-min+1) + min // 生成随机数

	// Store the code with an expiration time of 5 minutes
	err := cache.EmailCodeCache(emailAddr, code, 5) // 设置5分钟过期时间, 使用redis缓存验证码

	// Create the email message
	subject := "BUCT SQL Online Judge Verification Code"
	body := "Your verification code is " + strconv.Itoa(code)
	qqMail := config.QQMailAccount
	qqPassword := config.QQMailPassword
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", qqMail, emailAddr, subject, body)

	// Set up the SMTP client
	auth := smtp.PlainAuth("", qqMail, qqPassword, "smtp.qq.com")
	host := "smtp.qq.com:587"

	// Send the email message
	err = smtp.SendMail(host, auth, qqMail, []string{emailAddr}, []byte(msg))
	if err != nil {
		log.Println(err)
		return
	}
}
