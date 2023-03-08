package exercise

import (
	"fmt"
	"log"
	"reflect"
	"sqlOJ/cache"
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

		equal, getType := checkSqlSyntax(answer, expectedAnswer)
		if equal { // 和标准答案相等，返回正确
			status := 0 // 答案正确
			// 插入做题记录表
			model.NewSubmitHistoryFlow().InsertSubmitHistory(userID, exerciseID, userType, status, answer, userAgent, submitTime)

			wg.Wait() // 等待修改cache中状态的goroutine先完成，否则记录将不会被删去
			cache.DeleteSubmitStatus(userID, userType, exerciseID, submitTime)
			continue
		}
		if getType != expectedType { // sql语句类型不等，返回错误
			status := 1 // 答案错误
			model.NewSubmitHistoryFlow().InsertSubmitHistory(userID, exerciseID, userType, status, answer, userAgent, submitTime)
			wg.Wait() // 等待修改cache中状态的goroutine先完成，否则记录将不会被删去
			cache.DeleteSubmitStatus(userID, userType, exerciseID, submitTime)
			continue
		}
		var status int
		if getType == 1 {
			status = selectJudge(answer, expectedAnswer)

		} else {
			status = modifyJudge(userID, exerciseID, submitTime, answer, expectedAnswer, getType)
		}

		model.NewSubmitHistoryFlow().InsertSubmitHistory(userID, exerciseID, userType, status, answer, userAgent, submitTime)
		wg.Wait()
		cache.DeleteSubmitStatus(userID, userType, exerciseID, submitTime)
	}
}

// modifyJudge 负责评判 update, insert, delete 类型语句, 返回状态码 1->AC, 2->WA, 3->RE
func modifyJudge(userID int64, exerciseID int64, submitTime time.Time, userAnswer, expectedAnswer string, getType int) int {
	tempTableName := fmt.Sprintf("%d_%d_%d", userID, exerciseID, submitTime.Unix()) // 生成临时表名
	modifiedUserSql, originUserTableName := replaceTableName(userAnswer, tempTableName, getType)
	modifiedExpectedSql, originExpectedUserTableName := replaceTableName(expectedAnswer, tempTableName, getType)
	if originUserTableName != originExpectedUserTableName {
		// 用户提交的答案修改的表名与标准答案表名不一样，直接返回错误
		return 2
	}
	statusCode := model.CompareModifySqlResult(modifiedUserSql, modifiedExpectedSql, originExpectedUserTableName, tempTableName)
	return statusCode
}

// selectJudge 负责评判 select 类型数据, 返回状态码 1->AC, 2->WA, 3->RE
func selectJudge(userAnswer, expectedAnswer string) int {
	userResult, err := model.ExecuteRawSql(userAnswer)
	if err != nil {
		return 3 // RE
	}
	expectedResult, err := model.ExecuteRawSql(expectedAnswer)
	if err != nil {
		return 3 // RE
	}
	if reflect.DeepEqual(userResult, expectedResult) { // 判断二者查询结果是否相等
		// 相等
		return 1
	}
	// 不等
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
