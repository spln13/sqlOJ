package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/cache"
	"sqlOJ/controller/admin_account"
	"sqlOJ/controller/class"
	"sqlOJ/controller/contest"
	"sqlOJ/controller/exercise"
	"sqlOJ/controller/ranking"
	"sqlOJ/controller/student_account"
	"sqlOJ/controller/submission"
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
		context.HTML(http.StatusOK, "login.html", "")
	})
	server.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "home.html", "")
	})
	server.GET("/exercise/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "exercise.html", "")
	})
	server.GET("/contest/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "contest.html", "")
	})
	server.GET("/ranking/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "ranking.html", "")
	})
	server.GET("/register/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "register.html", "")
	})
	server.GET("/admin/login/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "admin-login.html", "")
	})
	server.GET("/teacher/login/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "teacher-login.html", "")
	})
	server.GET("/logout/", func(context *gin.Context) {
		context.SetCookie("token", "", -1, "/", "127.0.0.1:8080", true, false)
		context.SetCookie("username", "", -1, "/", "127.0.0.1:8080", true, false)
		context.Redirect(http.StatusFound, "/")
	})
	server.GET("/problem/:exercise_id", func(context *gin.Context) {
		context.HTML(http.StatusOK, "problem.html", "")
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
	classGroup := server.Group("/api/class")
	{
		classGroup.POST("/create/", middlewares.TeacherJWTMiddleware(), class.CreateClassHandle)            // 教师&管理员添加班级接口
		classGroup.POST("/add-student/", middlewares.TeacherJWTMiddleware(), class.AddStudentToClassHandle) // 教师添加学生进入班级接口
	}
	exerciseGroup := server.Group("/api/exercise")
	{
		exerciseGroup.POST("/publish/exercise", middlewares.TeacherJWTMiddleware(), exercise.PublishExerciseHandle)                             // 发布练习接口
		exerciseGroup.POST("/upload/table", middlewares.TeacherJWTMiddleware(), exercise.UploadTableHandle)                                     // 发布练习表单接口
		exerciseGroup.POST("/submit/", middlewares.StudentJWTMiddleware(), exercise.SubmitHandle)                                               // 处理提交习题接口
		exerciseGroup.GET("/get/all/without-token", exercise.GetAllExerciseWithoutTokenHandle)                                                  // 获取题库中所有可见的题目条目
		exerciseGroup.GET("/get/all/with-token", middlewares.StudentJWTMiddleware(), exercise.GetAllExerciseWithTokenHandle)                    //  登录用户获取题库中所有可见的题目
		exerciseGroup.GET("/get/one/", middlewares.StudentJWTMiddleware(), middlewares.CheckExerciseAuthority(), exercise.GetOneExerciseHandle) // 获取当前题目的题面
	}
	submissionGroup := server.Group("/api/submission")
	{
		submissionGroup.GET("/get/one-one/", middlewares.StudentJWTMiddleware(), submission.GetOneOneHandle)                                                      // 查询当前用户当前题目提交记录
		submissionGroup.GET("/get/one-all/", middlewares.StudentJWTMiddleware(), submission.GetOneAllHandle)                                                      // 查询当前用户所有提交记录
		submissionGroup.GET("/get/all-all/", middlewares.TeacherJWTMiddleware(), submission.GetAllAllHandle)                                                      // 获取所有提交记录
		submissionGroup.GET("/get/all-one/", middlewares.TeacherJWTMiddleware(), submission.GetAllOneHandle)                                                      // 获取当前题目所有用户的提交
		submissionGroup.GET("/contest/get-all/", middlewares.TeacherJWTMiddleware(), submission.ContestGetAllSubmissionHandle)                                    // 获取当前竞赛的所有提交
		submissionGroup.GET("/contest/get-my/", middlewares.StudentJWTMiddleware(), middlewares.CheckContestAuthority(), submission.ContestGetMySubmissionHandle) // 获取当前用户在竞赛中的所有提交
		submissionGroup.GET("/contest/get-exercise/", middlewares.TeacherJWTMiddleware(), submission.ContestGetOneExerciseHandle)                                 // 获取当前竞赛当前题目所有提交
	}
	rankingGroup := server.Group("/api/ranking")
	{
		rankingGroup.GET("/get/list/", ranking.GetRankingHandle) // 获取排行榜信息
	}
	contestGroup := server.Group("/api/contest")
	{
		contestGroup.POST("/create/", middlewares.TeacherJWTMiddleware(), contest.CreateContestHandle)                                                // 创建竞赛接口
		contestGroup.GET("/get/all/", contest.GetAllContestHandle)                                                                                    // 获取所有竞赛接口
		contestGroup.POST("/submit/", middlewares.StudentJWTMiddleware(), middlewares.CheckContestAuthority(), exercise.ContestSubmitHandle)          // 竞赛提交接口
		contestGroup.GET("/get/all-exercise/", middlewares.StudentJWTMiddleware(), middlewares.CheckContestAuthority(), contest.GetAllExerciseHandle) // 获取竞赛中所有的题目
		contestGroup.GET("/get/contest/", middlewares.StudentJWTMiddleware(), middlewares.CheckContestAuthority(), contest.GetContestHandle)          // 获取竞赛中所有的题目

	}
	return server
}
