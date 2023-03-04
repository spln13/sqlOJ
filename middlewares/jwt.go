package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"sqlOJ/common"
	"time"
)

var jwtKey = []byte("key_spln")

type Claims struct {
	UserId   uint
	UserType uint
	jwt.RegisteredClaims
}

// ReleaseToken 颁发管理员专属token
// UserType -> 用户等级标识; 1 -> 学生; 2 -> 教师; 3 -> 管理员
func ReleaseToken(ID uint, authority uint) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId:   ID,
		UserType: authority,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: expirationTime},
			//ExpiresAt: expirationTime.Unix(),
			IssuedAt: &jwt.NumericDate{Time: time.Now()},
			Issuer:   "linan",
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析token
func ParseToken(tokenString string) (*Claims, bool) {
	token, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if token != nil {
		if key, ok := token.Claims.(*Claims); ok {
			if token.Valid {
				return key, true
			} else {
				return key, false
			}
		}
	}
	return nil, false
}

func StudentJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token") // token通过Query传递
		if tokenStr == "" {          // token通过Form传递
			tokenStr = c.PostForm("token")
		}
		if tokenStr == "" { // token通过json传递
			var jsonMap map[string]interface{}
			if err := c.BindJSON(&jsonMap); err != nil {
				// 处理解析错误
				c.AbortWithStatusJSON(http.StatusBadRequest, common.NewCommonResponse(400, "解析错误"))
				c.Abort()
				return
			}
			tokenStr = jsonMap["token"].(string)
			c.Set("jsonMap", &jsonMap)
		}
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusBadRequest, common.NewCommonResponse(401, "用户不存在"))
			c.Abort() //阻止执行
			return
		}
		//验证token
		tokenStruck, ok := ParseToken(tokenStr)
		if !ok {
			c.JSON(http.StatusBadRequest, common.NewCommonResponse(403, "token不正确"))
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt.Time.Unix() {
			c.JSON(http.StatusBadRequest, common.NewCommonResponse(402, "token过期"))
			c.Abort() //阻止执行
			return
		}
		c.Set("user_id", tokenStruck.UserId)
		c.Next()
	}
}

func TeacherJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		if tokenStr == "" { // token通过json传递
			var jsonMap map[string]interface{}
			if err := c.BindJSON(&jsonMap); err != nil {
				// 处理解析错误
				c.AbortWithStatusJSON(http.StatusBadRequest, common.NewCommonResponse(400, "解析错误"))
				c.Abort()
				return
			}
			tokenStr = jsonMap["token"].(string)
			c.Set("jsonMap", &jsonMap)
		}
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusBadRequest, common.NewCommonResponse(401, "用户不存在"))
			c.Abort() //阻止执行

			return
		}
		//验证token
		tokenStruck, ok := ParseToken(tokenStr)
		if !ok {
			c.JSON(http.StatusBadRequest, common.NewCommonResponse(403, "token不正确"))
			c.Abort() //阻止执行
			return
		}
		if tokenStruck.UserType < 2 {
			c.JSON(http.StatusBadRequest, common.NewCommonResponse(401, "无权限"))
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt.Time.Unix() {
			c.JSON(http.StatusBadRequest, common.NewCommonResponse(402, "token过期"))
			c.Abort() //阻止执行
			return
		}
		c.Set("user_id", tokenStruck.UserId)
		c.Next()
	}
}

func AdminJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		if tokenStr == "" { // token通过json传递
			var jsonMap map[string]interface{}
			if err := c.BindJSON(&jsonMap); err != nil {
				// 处理解析错误
				c.AbortWithStatusJSON(http.StatusBadRequest, common.NewCommonResponse(400, "解析错误"))
				c.Abort()
				return
			}
			tokenStr = jsonMap["token"].(string)
			c.Set("jsonMap", &jsonMap)
		}
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusBadRequest, common.NewCommonResponse(401, "用户不存在"))
			c.Abort() //阻止执行

			return
		}
		//验证token
		tokenStruck, ok := ParseToken(tokenStr)
		if !ok {
			c.JSON(http.StatusBadRequest, common.NewCommonResponse(403, "token不正确"))
			c.Abort() //阻止执行
			return
		}
		if tokenStruck.UserType < 3 {
			c.JSON(http.StatusBadRequest, common.NewCommonResponse(401, "无权限"))
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt.Time.Unix() {
			c.JSON(http.StatusBadRequest, common.NewCommonResponse(402, "token过期"))
			c.Abort() //阻止执行
			return
		}
		c.Set("user_id", tokenStruck.UserId)
		c.Next()
	}
}
