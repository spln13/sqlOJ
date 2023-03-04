package router

import (
	"github.com/gin-gonic/gin"
	"sqlOJ/cache"
	"sqlOJ/controller/admin_account"
	"sqlOJ/controller/exercise"
	"sqlOJ/controller/student_account"
	"sqlOJ/controller/teacher_account"
	"sqlOJ/middlewares"
	"sqlOJ/model"
)

func InitServer() *gin.Engine {
	model.InitDB()          // 初始化MySQL数据库链接
	cache.InitRedis()       // 初始化Redis
	server := gin.Default() // 初始化gin服务器
	server.Static("static", "./static")
	server.LoadHTMLGlob("template/*")

	// 返回HTML页面
	server.GET("/login/", func(context *gin.Context) {
		context.HTML(200, "login.html", "")
	})
	server.GET("/register/", func(context *gin.Context) {
		context.HTML(200, "register.html", "")
	})
	server.GET("/admin/login/", func(context *gin.Context) {
		context.HTML(200, "admin-login.html", "")
	})
	server.GET("/teacher/login/", func(context *gin.Context) {
		context.HTML(200, "teacher-login.html", "")
	})

	// api接口
	adminGroup := server.Group("/api/admin")
	{
		adminGroup.POST("/add-debug/", middlewares.PasswordEncryptionMiddleware(), admin_account.AdminAddDebugHandle)
		adminGroup.POST("/login/", middlewares.PasswordEncryptionMiddleware(), admin_account.AdminLoginHandle)                                                         // 系统管理员登录逻辑
		adminGroup.POST("/add/", middlewares.AdminJWTMiddleware(), middlewares.PasswordEncryptionMiddleware(), admin_account.AdminAddHandle)                           // 系统管理员手动添加新的管理员
		adminGroup.POST("/change-password/", middlewares.AdminJWTMiddleware(), middlewares.TwoPasswordEncryptionMiddleware(), admin_account.AdminChangePasswordHandle) // 系统管理员更改密码
	}
	teacherGroup := server.Group("/api/teacher")
	{
		teacherGroup.POST("/login/", middlewares.PasswordEncryptionMiddleware(), teacher_account.TeacherLoginHandle)                                                           // 老师登录
		teacherGroup.POST("/add/", middlewares.TeacherJWTMiddleware(), middlewares.PasswordEncryptionMiddleware(), teacher_account.TeacherAddHandle)                           // 管理员&老师手动添加老师账号
		teacherGroup.POST("/change-password/", middlewares.TwoPasswordEncryptionMiddleware(), middlewares.TeacherJWTMiddleware(), teacher_account.TeacherChangePasswordHandle) // 老师改密码
	}
	studentGroup := server.Group("/api/student")
	{
		studentGroup.POST("/login/", middlewares.PasswordEncryptionMiddleware(), student_account.StudentLoginHandle)                                                           // 学生登录接口
		studentGroup.POST("/register/", middlewares.PasswordEncryptionMiddleware(), student_account.StudentRegisterHandle)                                                     // 学生注册接口
		studentGroup.POST("/change-password/", middlewares.StudentJWTMiddleware(), middlewares.TwoPasswordEncryptionMiddleware(), student_account.StudentChangePasswordHandle) // 学生改密码接口
		studentGroup.POST("/email/send-code/", student_account.SendCodeHandle)
	}
	exerciseGroup := server.Group("/api/exercise")
	{
		exerciseGroup.POST("/publish/exercise", middlewares.TeacherJWTMiddleware(), exercise.PublishExerciseHandle)
		exerciseGroup.POST("/upload/table", middlewares.TeacherJWTMiddleware(), exercise.UploadTableHandle)
	}
	return server
}
