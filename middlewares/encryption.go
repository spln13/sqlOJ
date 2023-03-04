package middlewares

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
)

// PasswordEncryptionMiddleware 使用SHA256算法对用户明文密码加密，向handles层发送加密后对密码进行后续处理
func PasswordEncryptionMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		password := context.Query("password") // 获取password
		if password == "" {
			password = context.PostForm("password")
		}
		digest := sha256.New() // 对密码加密
		digest.Write([]byte(password))
		passwordSHA := hex.EncodeToString(digest.Sum(nil))
		context.Set("password_sha256", passwordSHA) // 重写设置密码参数
		context.Next()                              // 放行
	}
}

// TwoPasswordEncryptionMiddleware 使用SHA256算法对用户明文密码加密
// 与PasswordEncryptionMiddleware不同的是此中间件仅当学生、教师、管理员更改密码时使用
// 参数中由`old_password`与`new_password`需要加密
func TwoPasswordEncryptionMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		oldPassword := context.Query("old_password")
		newPassword := context.Query("new_password")
		digest1 := sha256.New() // 对密码加密
		digest1.Write([]byte(oldPassword))
		OldPasswordSHA := hex.EncodeToString(digest1.Sum(nil))
		context.Set("old_password", OldPasswordSHA) // 重写设置密码参数
		digest2 := sha256.New()                     // 对密码加密
		digest2.Write([]byte(newPassword))
		NewPasswordSHA := hex.EncodeToString(digest2.Sum(nil))
		context.Set("new_password", NewPasswordSHA) // 重写设置密码参数
	}
}
