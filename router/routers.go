package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/cache"
	"sqlOJ/controller/admin_account"
	"sqlOJ/controller/class"
	"sqlOJ/controller/common"
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
	//server.GET("/register/", func(context *gin.Context) {
	//	context.HTML(http.StatusOK, "register.html", "")
	//})
	server.GET("/submission/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "my-submission.html", "")
	})
	server.GET("/migrate/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "migrate.html", "")
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
	server.GET("/submission/:submission_id", func(context *gin.Context) {
		context.HTML(http.StatusOK, "answer-detail.html", "")
	})
	server.GET("/contest/:contest_id", func(context *gin.Context) {
		context.HTML(http.StatusOK, "contest-problem-list.html", "")
	})
	server.GET("/exercise/my-submission/:exercise_id", func(context *gin.Context) {
		context.HTML(http.StatusOK, "exercise-my-submission.html", "")
	})
	server.GET("/contest/:contest_id/problem/:problem_id", func(context *gin.Context) {
		context.HTML(http.StatusOK, "contest-problem.html", "")
	})
	server.GET("/contest/:contest_id/my-submission", func(context *gin.Context) {
		context.HTML(http.StatusOK, "contest-my-submission.html", "")
	})
	server.GET("/contest/submission-detail/:submission_id", func(context *gin.Context) {
		context.HTML(http.StatusOK, "contest-answer-detail.html", "")
	})
	server.GET("/contest/:contest_id/problem/:problem_id/my-submission", func(context *gin.Context) {
		context.HTML(http.StatusOK, "contest-problem-my-submission.html", "")
	})
	server.GET("/contest/status/:contest_id", func(context *gin.Context) {
		context.HTML(http.StatusOK, "contest-status.html", "")
	})
	server.GET("/contest/submission/:contest_id", func(context *gin.Context) {
		context.HTML(http.StatusOK, "teacher-contest-submission.html", "")
	})

	server.GET("/teacher/exercise-answer/:exercise_id", func(context *gin.Context) {
		context.HTML(http.StatusOK, "teacher-exercise-answer.html", "")
	})
	server.GET("/download/score", middlewares.TeacherJWTMiddleware(), student_account.StudentScoreDownloadHandler) // 获取学生智能合约评分结果

	teacherHTMLGroup := server.Group("/teacher")
	{
		teacherHTMLGroup.GET("/upload-table/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-upload-table.html", "")
		})
		teacherHTMLGroup.GET("/publish-exercise/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-publish-exercise.html", "")
		})
		teacherHTMLGroup.GET("/publish-contest/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-publish-contest.html", "")
		})
		teacherHTMLGroup.GET("/submission/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-submission.html", "")
		})
		teacherHTMLGroup.GET("/tables/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-query-tables.html", "")
		})
		teacherHTMLGroup.GET("/students/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-students.html", "")
		})
		teacherHTMLGroup.GET("/exercises/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-exercises.html", "")
		})
		teacherHTMLGroup.GET("/contests/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-contests.html", "")
		})
		teacherHTMLGroup.GET("/class/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-class.html", "")
		})
		teacherHTMLGroup.GET("/create-class/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-create-class.html", "")
		})
		teacherHTMLGroup.GET("/migrate/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-migrate.html", "")
		})
		teacherHTMLGroup.GET("/score/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-student-score.html", "")
		})
		teacherHTMLGroup.GET("/add/", func(context *gin.Context) {
			context.HTML(http.StatusOK, "teacher-add.html", "")
		})
	}
	// api接口
	adminGroup := server.Group("/api/admin")
	{
		adminGroup.POST("/add-debug/", middlewares.PasswordEncryptionMiddleware(), admin_account.AdminAddDebugHandler)
		adminGroup.POST("/login/", middlewares.PasswordEncryptionMiddleware(), admin_account.AdminLoginHandler)                                                         // 系统管理员登录逻辑
		adminGroup.POST("/add/", middlewares.AdminJWTMiddleware(), middlewares.PasswordEncryptionMiddleware(), admin_account.AdminAddHandler)                           // 系统管理员手动添加新的管理员
		adminGroup.POST("/change-password/", middlewares.AdminJWTMiddleware(), middlewares.TwoPasswordEncryptionMiddleware(), admin_account.AdminChangePasswordHandler) // 系统管理员更改密码
	}
	teacherGroup := server.Group("/api/teacher")
	{
		teacherGroup.POST("/login/", middlewares.PasswordEncryptionMiddleware(), teacher_account.TeacherLoginHandler)                                                           // 老师登录
		teacherGroup.POST("/add/", middlewares.TeacherJWTMiddleware(), middlewares.PasswordEncryptionMiddleware(), teacher_account.TeacherAddHandler)                           // 管理员&老师手动添加老师账号
		teacherGroup.POST("/change-password/", middlewares.TwoPasswordEncryptionMiddleware(), middlewares.TeacherJWTMiddleware(), teacher_account.TeacherChangePasswordHandler) // 老师改密码
	}
	studentGroup := server.Group("/api/student")
	{
		studentGroup.POST("/login/", middlewares.PasswordEncryptionMiddleware(), student_account.StudentLoginHandler) // 学生登录接口
		//studentGroup.POST("/register/", middlewares.PasswordEncryptionMiddleware(), student_account.StudentRegisterHandler)                                                     // 学生注册接口
		studentGroup.POST("/change-password/", middlewares.StudentJWTMiddleware(), middlewares.TwoPasswordEncryptionMiddleware(), student_account.StudentChangePasswordHandler) // 学生改密码接口
		//studentGroup.POST("/email/send-code/", student_account.SendCodeHandler)                                                                                                 // 发送验证码
		studentGroup.GET("/get/all-students/", middlewares.TeacherJWTMiddleware(), student_account.GetAllStudentsHandler) // 获取所有学生信息
		studentGroup.POST("/reset/", middlewares.TeacherJWTMiddleware(), student_account.StudentResetHandler)             // 重置学生密码
		studentGroup.GET("/get/score/", middlewares.TeacherJWTMiddleware(), student_account.StudentScoreHandler)          // 获取学生智能合约评分结果
		studentGroup.DELETE("/delete/", middlewares.TeacherJWTMiddleware(), student_account.StudentDeleteHandler)         // 删除学生
	}
	classGroup := server.Group("/api/class")
	{
		classGroup.POST("/create/", middlewares.TeacherJWTMiddleware(), class.CreateClassHandler)            // 教师&管理员添加班级接口
		classGroup.POST("/add-student/", middlewares.TeacherJWTMiddleware(), class.AddStudentToClassHandler) // 教师添加学生进入班级接口
		classGroup.GET("/get/all-class", middlewares.TeacherJWTMiddleware(), class.GetAllClassHandler)       // 获取所有班级
	}
	exerciseGroup := server.Group("/api/exercise")
	{
		exerciseGroup.POST("/publish/exercise", middlewares.TeacherJWTMiddleware(), exercise.PublishExerciseHandler)          // 发布练习接口
		exerciseGroup.POST("/upload/table", middlewares.TeacherJWTMiddleware(), exercise.UploadTableHandler)                  // 发布练习表单接口
		exerciseGroup.POST("/submit/", middlewares.StudentJWTMiddleware(), exercise.SubmitHandler)                            // 处理提交习题接口
		exerciseGroup.GET("/get/all/without-token", exercise.GetAllExerciseWithoutTokenHandler)                               // 获取题库中所有可见的题目条目
		exerciseGroup.GET("/get/all/with-token", middlewares.StudentJWTMiddleware(), exercise.GetAllExerciseWithTokenHandler) //  登录用户获取题库中所有可见的题目
		exerciseGroup.GET("/get/one/", middlewares.StudentJWTMiddleware(), exercise.GetOneExerciseHandler)                    // 获取当前题目的题面
		exerciseGroup.GET("/get/all-tables/", middlewares.TeacherJWTMiddleware(), exercise.GetAllTableHandler)                // 获取所有数据表
		exerciseGroup.GET("/teacher/all-exercises/", middlewares.TeacherJWTMiddleware(), exercise.TeacherGetAllExercises)     // 教师获取题库所有题目
		exerciseGroup.GET("/teacher/answer/", middlewares.TeacherJWTMiddleware(), exercise.TeacherGetAnswer)                  // 教师获取题目答案
		exerciseGroup.DELETE("/delete/", middlewares.TeacherJWTMiddleware(), exercise.DeleteExerciseHandler)
		exerciseGroup.DELETE("/delete-table/", middlewares.TeacherJWTMiddleware(), exercise.DeleteTableHandler)
	}
	submissionGroup := server.Group("/api/submission")
	{
		submissionGroup.GET("/get/one-one/", middlewares.StudentJWTMiddleware(), submission.GetOneOneHandler)                   // 查询当前用户当前题目提交记录
		submissionGroup.GET("/get/one-all/", middlewares.StudentJWTMiddleware(), submission.GetOneAllHandler)                   // 查询当前用户所有提交记录
		submissionGroup.GET("/get/all-all/", middlewares.TeacherJWTMiddleware(), submission.GetAllAllHandler)                   // 获取所有提交记录
		submissionGroup.GET("/get/all-one/", middlewares.TeacherJWTMiddleware(), submission.GetAllOneHandler)                   // 获取当前题目所有用户的提交
		submissionGroup.GET("/get/answer-detail/", middlewares.StudentJWTMiddleware(), submission.GetAnswerDetailHandler)       // 获取提交详情信息
		submissionGroup.GET("/contest/get-all/", middlewares.TeacherJWTMiddleware(), submission.ContestGetAllSubmissionHandler) // 获取当前竞赛的所有提交
		//submissionGroup.GET("/contest/get-my/", middlewares.StudentJWTMiddleware(), middlewares.CheckContestGetAuthority(), submission.ContestGetMySubmissionHandler) // 获取当前用户在竞赛中的所有提交
		submissionGroup.GET("/contest/get-my/", middlewares.StudentJWTMiddleware(), submission.ContestGetMySubmissionHandler)      // test获取当前用户在竞赛中的所有提交
		submissionGroup.GET("/contest/get-exercise/", middlewares.TeacherJWTMiddleware(), submission.ContestGetOneExerciseHandler) // 获取当前竞赛当前题目所有提交
		submissionGroup.GET("/contest/detail/", middlewares.StudentJWTMiddleware(), submission.ContestGetDetailHandler)            // 获取竞赛提交记录对应的答案
		submissionGroup.GET("/contest/", middlewares.StudentJWTMiddleware(), submission.ContestGetUserExerciseHandler)             // 获取当前用户当前竞赛当前题目的所有提交

	}
	rankingGroup := server.Group("/api/ranking")
	{
		rankingGroup.GET("/get/list/", ranking.GetRankingHandler)   // 获取排行榜信息
		rankingGroup.GET("/get/min/", ranking.GetMinRankingHandler) // 获取min排行榜信息
	}
	contestGroup := server.Group("/api/contest")
	{
		contestGroup.POST("/create/", middlewares.TeacherJWTMiddleware(), contest.CreateContestHandler)                                          // 创建竞赛接口
		contestGroup.GET("/get/all/", contest.GetAllContestHandler)                                                                              // 获取所有竞赛接口
		contestGroup.GET("/get/contest/", middlewares.StudentJWTMiddleware(), middlewares.CheckContestGetAuthority(), contest.GetContestHandler) // 获取竞赛详情信息                                                                          // 获取所有竞赛接口
		//contestGroup.POST("/submit/", middlewares.StudentJWTMiddleware(), exercise.ContestSubmitHandler)          // test竞赛提交接口
		//contestGroup.GET("/get/all-exercise/", middlewares.StudentJWTMiddleware(), contest.GetAllExerciseHandler) // test获取竞赛中所有的题目
		contestGroup.GET("/status/", middlewares.TeacherJWTMiddleware(), contest.GetContestStatusHandler)                                                    // 获取竞赛的状态
		contestGroup.POST("/submit/", middlewares.StudentJWTMiddleware(), middlewares.CheckContestSubmitAuthority(), exercise.ContestSubmitHandler)          // 竞赛提交接口
		contestGroup.GET("/get/all-exercise/", middlewares.StudentJWTMiddleware(), middlewares.CheckContestGetAuthority(), contest.GetAllExerciseHandler)    // 获取竞赛中所有的题目
		contestGroup.GET("/get/check-authority/", middlewares.StudentJWTMiddleware(), middlewares.CheckContestGetAuthority(), contest.CheckAuthorityHandler) // 获取竞赛中所有的题目
		contestGroup.DELETE("/delete/", middlewares.TeacherJWTMiddleware(), contest.DeleteContestHandler)                                                    // 获取竞赛中所有的题目

	}
	server.GET("/api/get-type/", middlewares.StudentJWTMiddleware(), common.GetTypeHandler) // 获取用户类型
	return server
}
