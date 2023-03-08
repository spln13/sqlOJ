package exercise

import (
	"github.com/xwb1989/sqlparser"
	"reflect"
	"sqlOJ/cache"
	"sqlOJ/model"
	"sync"
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
			// 后续处理
		} else {
			status = modifyJudge(answer, expectedAnswer)
		}

		model.NewSubmitHistoryFlow().InsertSubmitHistory(userID, exerciseID, userType, status, answer, userAgent, submitTime)
		wg.Wait()
		cache.DeleteSubmitStatus(userID, userType, exerciseID, submitTime)
	}
}

// modifyJudge 负责评判 update, insert, delete 类型语句, 返回状态码 1->AC, 2->WA, 3->RE
func modifyJudge(userAnswer, expectedAnswer string) int {
	// TODO: 判题逻辑
	return 0
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

// checkSqlSyntax 检查用户提交的sql语句语法是否正确，并于标准答案(同样经过Parse)比对
// 并返回用户sql语句类型
func checkSqlSyntax(userAnswer string, expectedAnswer string) (bool, int) {
	userStmt, err := sqlparser.Parse(userAnswer)
	if err != nil {
		// userAnswer中有语法错误
		return false, 0
	}
	userStmtStr := sqlparser.String(userStmt)
	var code int
	switch stmt := userStmt.(type) {
	case *sqlparser.Select:
		code = 1
		_ = stmt
	case *sqlparser.Insert:
		code = 2
	case *sqlparser.Update:
		code = 3
	case *sqlparser.Delete:
		code = 4
	default:
		code = 0
	}
	return userStmtStr == expectedAnswer, code
}

// InitJudgeGoroutine 开启n个判题协程从channel中读取SubmitMessage开始判题
func InitJudgeGoroutine(n int) {
	for i := 0; i < n; i++ {
		go judge()
	}
}
