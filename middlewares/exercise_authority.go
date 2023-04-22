package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sqlOJ/cache"
	"sqlOJ/common"
)

// CheckExerciseAuthority 是一个用于检测用户是否有权限访问题库中的当前题目的中间件
// 参与竞赛的用户在竞赛期间无法访问被竞赛引用的题库原题, 该功能主要基于Redis实现
// 首先判断exerciseID对应的Set是否有竞赛, 若有则一次访问对应的竞赛, 看其是否过期, 再判断userID是否包含与contestID对应的Set, 若存在则Abort
// 若竞赛已过期, 则在exerciseID对应的Set中删除
func CheckExerciseAuthority() gin.HandlerFunc {
	return func(context *gin.Context) {
		userID, ok := context.MustGet("user_id").(int64)
		if !ok {
			context.JSON(http.StatusBadRequest, common.NewCommonResponse(403, "用户token错误"))
			context.Abort()
			return
		}
		exerciseIDStr := context.Query("exercise_id")
		contestIDStrList, err := cache.GetExerciseSetMember(exerciseIDStr)
		if err != nil {
			context.JSON(http.StatusInternalServerError, common.NewCommonResponse(403, err.Error()))
			context.Abort()
			return
		}
		for _, contestIDStr := range contestIDStrList {
			code, err := cache.CheckUserIDInContest(userID, contestIDStr)
			// code: 0->错误; 1->该键值不存在; 2->集合中存在; 3->集合中不存在
			if err != nil {
				context.JSON(http.StatusInternalServerError, common.NewCommonResponse(403, err.Error()))
				context.Abort()
				return
			}
			if code == 1 { // 键值不存在, 即竞赛已经结束
				// 删除exercise对应Set中的contestID
				if err := cache.DeleteContestIDInExercise(exerciseIDStr, contestIDStr); err != nil {
					context.JSON(http.StatusInternalServerError, common.NewCommonResponse(403, err.Error()))
					context.Abort()
					return
				}
				continue
			} else if code == 2 { // 集合中存在, 即学生此刻参与的竞赛有引用此题目
				context.JSON(http.StatusOK, common.NewCommonResponse(403, "竞赛进行中"))
				context.Abort()
				return
			}
		}
		context.Next()
	}
}

// CheckContestAuthority 检查用户是否有访问当前竞赛信息的权限
// 判断用户id是否在竞赛在Redis中对应的Set中
func CheckContestAuthority() gin.HandlerFunc {
	return func(context *gin.Context) {
		userID, ok1 := context.MustGet("user_id").(int64)
		userType, ok2 := context.MustGet("user_type").(int64)
		if !ok1 || !ok2 {
			context.JSON(http.StatusBadRequest, common.NewCommonResponse(403, "用户token错误"))
			context.Abort()
			return
		}
		if userType > 1 { // 教师或管理员, 可以直接访问
			context.Next()
		}
		exerciseIDStr := context.Query("contest_id")
		if exerciseIDStr == "" {
			exerciseIDStr = context.PostForm("contest_id")
		}
		code, err := cache.CheckUserIDInContest(userID, exerciseIDStr)
		if err != nil {
			context.JSON(http.StatusInternalServerError, common.NewCommonResponse(403, err.Error()))
			context.Abort()
			return
		}
		// code: 0->错误; 1->该键值不存在; 2->集合中存在; 3->集合中不存在
		fmt.Println(code)
		if code != 2 {
			context.JSON(http.StatusOK, common.NewCommonResponse(403, "您无权访问"))
			context.Abort()
			return
		}
		context.Next()
	}
}
