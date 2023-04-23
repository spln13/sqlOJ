package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

var jwtKey = []byte("key_spln")

type Claims struct {
	UserId   int64
	UserType int64
	jwt.RegisteredClaims
}

// ReleaseToken 颁发管理员专属token
// UserType -> 用户等级标识; 1 -> 学生; 2 -> 教师; 3 -> 管理员
func ReleaseToken(ID int64, authority int64) (string, error) {
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
		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.NewCommonResponse(402, "token不存在"))
			c.Abort()
			return
		}
		//验证token
		tokenStruck, ok := ParseToken(tokenStr)
		if !ok {
			c.JSON(http.StatusBadRequest, utils.NewCommonResponse(403, "token不正确"))
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt.Time.Unix() {
			// token超时, 清空token
			c.SetCookie("token", "", -1, "/", "localhost:8080", true, false)
			c.SetCookie("username", "", -1, "/", "localhost:8080", true, false)

			c.JSON(http.StatusBadRequest, utils.NewCommonResponse(402, "token过期"))
			c.Abort() //阻止执行
			return
		}
		c.Set("user_id", tokenStruck.UserId)
		c.Set("user_type", tokenStruck.UserType)
		c.Next()
	}
}

func TeacherJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.NewCommonResponse(402, "token不存在"))
			c.Abort()
			return
		}
		//验证token
		tokenStruck, ok := ParseToken(tokenStr)
		if !ok {
			c.JSON(http.StatusBadRequest, utils.NewCommonResponse(403, "token不正确"))
			c.Abort() //阻止执行
			return
		}
		if tokenStruck.UserType < 2 {
			c.JSON(http.StatusBadRequest, utils.NewCommonResponse(401, "无权限"))
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt.Time.Unix() {
			c.JSON(http.StatusBadRequest, utils.NewCommonResponse(402, "token过期"))
			c.Abort() //阻止执行
			return
		}
		c.Set("user_id", tokenStruck.UserId)
		c.Set("user_type", tokenStruck.UserType)
		c.Next()
	}
}

func AdminJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.NewCommonResponse(402, "token不存在"))
			c.Abort()
			return
		}
		//验证token
		tokenStruck, ok := ParseToken(tokenStr)
		if !ok {
			c.JSON(http.StatusBadRequest, utils.NewCommonResponse(403, "token不正确"))
			c.Abort() //阻止执行
			return
		}
		if tokenStruck.UserType < 3 {
			c.JSON(http.StatusBadRequest, utils.NewCommonResponse(401, "无权限"))
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt.Time.Unix() {
			c.JSON(http.StatusBadRequest, utils.NewCommonResponse(402, "token过期"))
			c.Abort() //阻止执行
			return
		}
		c.Set("user_id", tokenStruck.UserId)
		c.Set("user_type", tokenStruck.UserType)
		c.Next()
	}
}
