package exercise

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sqlOJ/cache"
	"sqlOJ/common"
	"sqlOJ/model"
	"sync"
	"time"
)

// judge 负责从channel中读取SubmitMessage并进行判题处理
func judge() {
	for {
		message := <-JudgeQueue // 从判题队列中读取提交数据
		// channel 线程安全不用上锁
		userID := message.UserID
		userType := message.UserType
		exerciseID := message.ExerciseID
		submitTime := message.SubmitTime
		var wg sync.WaitGroup                                                         // 获取 SetJudgeStatusJudging 协程的完成信息
		go cache.SetJudgeStatusJudging(userID, userType, exerciseID, submitTime, &wg) // 设置cache中的判题状态
		wg.Add(1)
		answer := message.Answer
		userAgent := message.UserAgent
		expectedAnswer, expectedType := model.NewExerciseContentFlow().QueryAnswerTypeByExerciseID(exerciseID)
		// 获取参数

		// 获取用户名与题目名
		username := common.QueryUsername(userID, userType)
		exerciseName := model.NewExerciseContentFlow().QueryExerciseNameByExerciseID(exerciseID)
		equal, getType := checkSqlSyntax(answer, expectedAnswer)
		if equal { // 和标准答案相等，返回正确
			status := 1 // 答案正确
			// 插入做题记录表
			model.NewSubmitHistoryFlow().InsertSubmitHistory(userID, exerciseID, userType, status, answer, userAgent, username, exerciseName, submitTime)
			model.NewExerciseContentFlow().IncreasePassCountSubmitCount(exerciseID)
			wg.Wait() // 等待修改cache中状态的goroutine先完成，否则记录将不会被删去
			cache.DeleteSubmitStatus(userID, userType, exerciseID, submitTime)
			continue
		}
		if getType != expectedType { // sql语句类型不等，返回错误
			fmt.Println("getType:", getType)
			fmt.Println("expectedType:", expectedType)
			fmt.Println("getType != expectedType")
			status := 2 // 答案错误
			model.NewSubmitHistoryFlow().InsertSubmitHistory(userID, exerciseID, userType, status, answer, userAgent, username, exerciseName, submitTime)
			model.NewExerciseContentFlow().IncreaseSubmitCount(exerciseID)
			wg.Wait() // 等待修改cache中状态的goroutine先完成，否则记录将不会被删去
			cache.DeleteSubmitStatus(userID, userType, exerciseID, submitTime)
			continue
		}
		var status int
		if getType == 1 {
			status = selectJudge(answer, expectedAnswer, exerciseID)

		} else {
			status = modifyJudge(userID, exerciseID, submitTime, answer, expectedAnswer, getType)
		}
		model.NewSubmitHistoryFlow().InsertSubmitHistory(userID, exerciseID, userType, status, answer, userAgent, username, exerciseName, submitTime)
		if status == 1 { // 答案正确
			model.NewExerciseContentFlow().IncreasePassCountSubmitCount(exerciseID) // 自增提交总数和通过总数
			grade := model.NewExerciseContentFlow().QueryExerciseGrade(exerciseID)  // 获取当前习题的难度
			model.NewScoreRecordFlow().IncreaseScore(userID, userType, grade)       // 增加用户的积分
		} else { // 答案错误
			model.NewExerciseContentFlow().IncreaseSubmitCount(exerciseID) // 自增提交总数
		}
		model.NewUserProblemStatusFlow().ModifyUserProblemStatus(userID, exerciseID, userType, status) // 将用户做题数据写入用户做题表
		wg.Wait()
		cache.DeleteSubmitStatus(userID, userType, exerciseID, submitTime)
	}
}

// modifyJudge 负责评判 update, insert, delete 类型语句, 返回状态码 1->AC, 2->WA, 3->RE
func modifyJudge(userID int64, exerciseID int64, submitTime time.Time, userAnswer, expectedAnswer string, getType int) int {
	tempTableName := fmt.Sprintf("%d_%d_%d", userID, exerciseID, submitTime.Unix())                              // 生成用户临时表名
	modifiedUserSql, originUserTableName := replaceTableName(userAnswer, tempTableName, getType)                 // 将用户sql语句中的表名替换
	modifiedExpectedSql, originExpectedUserTableName := replaceTableName(expectedAnswer, tempTableName, getType) // 将答案sql语句中的表名替换
	if originUserTableName != originExpectedUserTableName {
		// 用户提交的答案修改的表名与标准答案表名不一样，直接返回错误
		return 2
	}
	// 对比答案sql语句和用户sql语句
	cacheResultBytes, ok := cache.GetExpectedResultBytes(exerciseID)
	if ok { // 缓存中存在答案
		statusCode := model.CompareModifySqlResultWithCache(modifiedUserSql, originExpectedUserTableName, tempTableName, cacheResultBytes)
		return statusCode
	}
	expectedResult, statusCode := model.CompareModifySqlResultWithoutCache(modifiedUserSql, modifiedExpectedSql, originExpectedUserTableName, tempTableName)
	if statusCode != 3 { // 若不发生错误
		cache.ExpectedResultCache(exerciseID, expectedResult) // 将答案sql语句执行结果写入cache]
	}
	return statusCode
}

// selectJudge 负责评判 select 类型数据, 返回状态码 1->AC, 2->WA, 3->RE
// 执行答案sql语句之前先查询cache中是否有对应缓存，若有则不执行sql; 若没有则执行sql, 并将结果写入cache
func selectJudge(userAnswer, expectedAnswer string, exerciseID int64) int {
	userResult, err := model.ExecuteRawSql(userAnswer)
	userResultBytes, err := json.Marshal(userResult) // 将查询结果转换成[]byte,与缓存中对应
	if err != nil {
		log.Println(err)
		return 3
	}
	expectedResultBytes, ok := cache.GetExpectedResultBytes(exerciseID)
	if !ok { // 在cache中没有查找到
		expectedResult, err := model.ExecuteRawSql(expectedAnswer) // 执行sql语句得到查询结果
		if err != nil {
			return 3 // RE
		}
		cache.ExpectedResultCache(exerciseID, expectedResult)   // 缓存查询结果
		expectedResultBytes, err = json.Marshal(expectedResult) // 将查询结果转换成[]byte,与缓存中对应
		if err != nil {
			log.Println(err)
			return 3
		}
	}
	if reflect.DeepEqual(userResultBytes, expectedResultBytes) { // 判断二者查询结果是否相等
		// 相等
		fmt.Println("相等")
		return 1
	}
	// 不等
	fmt.Println("不等")
	return 2
}

// InitJudgeGoroutine 开启n个判题协程从channel中读取SubmitMessage开始判题
func InitJudgeGoroutine(n int) {
	log.Println("Init Judge Goroutine...")
	for i := 0; i < n; i++ {
		go judge()
	}
	log.Println("Init Judge Goroutine finished.")
}
