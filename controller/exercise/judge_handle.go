package exercise

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"sqlOJ/cache"
	"sqlOJ/model"
	"sqlOJ/utils"
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
		answer := message.Answer
		userAgent := message.UserAgent
		isContest := message.IsContest // 是否是竞赛判题
		if isContest {
			contestID := message.ContestID
			contestJudge(userID, userType, exerciseID, contestID, submitTime, answer, userAgent)
		} else {
			exerciseJudge(userID, userType, exerciseID, submitTime, answer, userAgent)
		}

	}
}

// exerciseJudge 用于评测题库中提交的题目
func exerciseJudge(userID, userType, exerciseID int64, submitTime time.Time, answer, userAgent string) {
	var wg sync.WaitGroup                                                                 // 获取 SetExerciseJudgeStatusJudging 协程的完成信息
	go cache.SetExerciseJudgeStatusJudging(userID, userType, exerciseID, submitTime, &wg) // 设置cache中的判题状态
	wg.Add(1)

	expectedAnswer, expectedType := model.NewExerciseContentFlow().QueryAnswerTypeByExerciseID(exerciseID)
	// 获取参数

	// 获取用户名与题目名
	username := utils.QueryUsername(userID, userType)
	exerciseName := model.NewExerciseContentFlow().QueryExerciseNameByExerciseID(exerciseID)
	equal, getType := checkSqlSyntax(answer, expectedAnswer)
	if equal { // 和标准答案相等，返回正确
		status := 1 // 答案正确
		// 插入做题记录表
		model.NewExerciseContentFlow().IncreasePassCountSubmitCount(exerciseID)
		grade := model.NewExerciseContentFlow().QueryExerciseGrade(exerciseID) // 获取当前习题的难度
		submissionID := model.NewSubmitHistoryFlow().InsertSubmitHistory(userID, exerciseID, userType, status, answer, userAgent, username, exerciseName, submitTime)
		appendPendingTxQueue(userID, userType, submissionID, 0, exerciseID, status, grade) // 将提交记录append到上链队列
		thisExerciseStatus := model.NewUserProblemStatusFlow().QueryUserProblemStatus(userID, userType, exerciseID)
		if thisExerciseStatus > 1 { // wa&re&未提交
			model.NewScoreRecordFlow().IncreaseScore(userID, userType, grade) // 增加用户的积分
		}
		model.NewUserProblemStatusFlow().ModifyUserProblemStatus(userID, exerciseID, userType, status) // 将用户题目提交状态表中的状态设置为正确
		wg.Wait()                                                                                      // 等待修改cache中状态的goroutine先完成，否则记录将不会被删去
		cache.DeleteExerciseSubmitStatus(userID, userType, exerciseID, submitTime)
		return
	}
	if getType != expectedType { // sql语句类型不等，返回错误
		//fmt.Println("getType:", getType)
		//fmt.Println("expectedType:", expectedType)
		//fmt.Println("getType != expectedType")
		status := 2 // 答案错误
		submissionID := model.NewSubmitHistoryFlow().InsertSubmitHistory(userID, exerciseID, userType, status, answer, userAgent, username, exerciseName, submitTime)
		appendPendingTxQueue(userID, userType, submissionID, 0, exerciseID, status, 0) // 将提交记录append到上链队列
		model.NewExerciseContentFlow().IncreaseSubmitCount(exerciseID)
		model.NewUserProblemStatusFlow().ModifyUserProblemStatus(userID, exerciseID, userType, status) // 将用户题目提交状态表中的状态设置为错误
		wg.Wait()                                                                                      // 等待修改cache中状态的goroutine先完成，否则记录将不会被删去
		cache.DeleteExerciseSubmitStatus(userID, userType, exerciseID, submitTime)
		return
	}
	var status int
	if getType == 1 {
		status = selectJudge(answer, expectedAnswer, exerciseID)

	} else {
		status = modifyJudge(userID, exerciseID, submitTime, answer, expectedAnswer, getType)
	}
	submissionID := model.NewSubmitHistoryFlow().InsertSubmitHistory(userID, exerciseID, userType, status, answer, userAgent, username, exerciseName, submitTime)
	if status == 1 { // 答案正确
		model.NewExerciseContentFlow().IncreasePassCountSubmitCount(exerciseID)            // 自增提交总数和通过总数
		grade := model.NewExerciseContentFlow().QueryExerciseGrade(exerciseID)             // 获取当前习题的难度
		appendPendingTxQueue(userID, userType, submissionID, 0, exerciseID, status, grade) // 将提交记录append到上链队列
		thisExerciseStatus := model.NewUserProblemStatusFlow().QueryUserProblemStatus(userID, userType, exerciseID)
		if thisExerciseStatus > 1 { // wa&re&未提交
			model.NewScoreRecordFlow().IncreaseScore(userID, userType, grade) // 增加用户的积分
		}
	} else { // 答案错误
		model.NewExerciseContentFlow().IncreaseSubmitCount(exerciseID)                 // 自增提交总数
		appendPendingTxQueue(userID, userType, submissionID, 0, exerciseID, status, 0) // 将提交记录append到上链队列
	}
	model.NewUserProblemStatusFlow().ModifyUserProblemStatus(userID, exerciseID, userType, status) // 将用户做题数据写入用户做题表
	wg.Wait()
	cache.DeleteExerciseSubmitStatus(userID, userType, exerciseID, submitTime)
}

// contestJudge 用于评测竞赛中的题目
func contestJudge(userID, userType, exerciseID, contestID int64, submitTime time.Time, answer, userAgent string) {
	var wg sync.WaitGroup
	go cache.SetContestJudgeStatusJudging(userID, userType, exerciseID, contestID, submitTime, &wg)
	wg.Add(1)
	expectedAnswer, expectedType := model.NewExerciseContentFlow().QueryAnswerTypeByExerciseID(exerciseID)
	// 获取参数

	// 获取用户名, 题目名, 竞赛名
	username := utils.QueryUsername(userID, userType)
	exerciseName := model.NewExerciseContentFlow().QueryExerciseNameByExerciseID(exerciseID)
	contestName := model.NewContestFlow().GetContestNameByID(contestID)
	equal, getType := checkSqlSyntax(answer, expectedAnswer)
	if equal { // 和标准答案相等，返回正确
		status := 1 // 答案正确
		// 插入做题记录表
		submissionID := model.NewContestSubmissionFlow().InsertContestSubmission(contestID, exerciseID, userID, userType, username, answer, exerciseName, userAgent, contestName, status, submitTime)
		grade := model.NewExerciseContentFlow().QueryExerciseGrade(exerciseID)
		appendPendingTxQueue(userID, userType, submissionID, contestID, exerciseID, status, grade)
		model.NewExerciseContentFlow().IncreasePassCountSubmitCount(exerciseID)
		model.NewContestExerciseStatusFlow().ModifyContestExerciseStatus(userID, userType, exerciseID, contestID, status) // 将用户竞赛题目提交状态表中的状态设置为正确
		wg.Wait()                                                                                                         // 等待修改cache中状态的goroutine先完成，否则记录将不会被删去
		cache.DeleteContestSubmitStatus(userID, userType, exerciseID, contestID, submitTime)
		return
	}
	if getType != expectedType { // sql语句类型不等，返回错误
		status := 2 // 答案错误
		submissionID := model.NewContestSubmissionFlow().InsertContestSubmission(contestID, exerciseID, userID, userType, username, answer, exerciseName, userAgent, contestName, status, submitTime)
		appendPendingTxQueue(userID, userType, submissionID, contestID, exerciseID, status, 0)
		model.NewExerciseContentFlow().IncreaseSubmitCount(exerciseID)
		model.NewContestExerciseStatusFlow().ModifyContestExerciseStatus(userID, userType, exerciseID, contestID, status) // 更改提交状态
		wg.Wait()                                                                                                         // 等待修改cache中状态的goroutine先完成，否则记录将不会被删去
		cache.DeleteContestSubmitStatus(userID, userType, exerciseID, contestID, submitTime)
		return
	}
	var status int
	if getType == 1 { // select
		status = selectJudge(answer, expectedAnswer, exerciseID)

	} else { // update, insert, modify
		status = modifyJudge(userID, exerciseID, submitTime, answer, expectedAnswer, getType)
	}
	submissionID := model.NewContestSubmissionFlow().InsertContestSubmission(contestID, exerciseID, userID, userType, username, answer, exerciseName, userAgent, contestName, status, submitTime)
	grade := model.NewExerciseContentFlow().QueryExerciseGrade(exerciseID)
	appendPendingTxQueue(userID, userType, submissionID, contestID, exerciseID, status, grade)
	if status == 1 { // 答案正确
		model.NewExerciseContentFlow().IncreasePassCountSubmitCount(exerciseID) // 自增提交总数和通过总数
		//grade := model.NewExerciseContentFlow().QueryExerciseGrade(exerciseID)  // 获取当前习题的难度
		//model.NewScoreRecordFlow().IncreaseScore(userID, userType, grade)       // 增加用户的积分
	} else { // 答案错误
		model.NewExerciseContentFlow().IncreaseSubmitCount(exerciseID) // 自增提交总数
	}
	model.NewContestExerciseStatusFlow().ModifyContestExerciseStatus(userID, userType, exerciseID, contestID, status) // 更改提交状态
	wg.Wait()
	cache.DeleteContestSubmitStatus(userID, userType, exerciseID, contestID, submitTime)
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
		cache.ExpectedResultCache(exerciseID, expectedResult)   // 将查询结果写入缓存
		expectedResultBytes, err = json.Marshal(expectedResult) // 将查询结果转换成[]byte,与缓存中对应
		if err != nil {
			log.Println(err)
			return 3
		}
	}
	if reflect.DeepEqual(userResultBytes, expectedResultBytes) { // 判断二者查询结果是否相等
		return 1
	}
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
