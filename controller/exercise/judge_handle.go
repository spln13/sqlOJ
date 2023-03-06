package exercise

import (
	"github.com/xwb1989/sqlparser"
	"sqlOJ/cache"
	"sqlOJ/model"
	"sync"
)

// Judge 负责从channel中读取SubmitMessage并进行判题处理
func judge() {
	for {
		message := <-JudgeQueue // 从判题队列中读取提交数据
		// channel 线程安全不用上锁
		userID := message.UserID
		userType := message.UserType
		exerciseID := message.ExerciseID
		var wg sync.WaitGroup
		go cache.SetJudgeStatusJudging(userID, userType, exerciseID, &wg) // 设置cache中的判题状态
		wg.Add(1)
		answer := message.Answer
		userAgent := message.UserAgent
		submitTime := message.SubmitTime
		expectedAnswer, expectedType := model.NewExerciseContentFlow().QueryAnswerTypeByExerciseID(exerciseID)
		equal, getType := checkSqlSyntax(answer, expectedAnswer)
		if equal { // 和标准答案相等，返回正确
			// TODO: 删除cache中对应的判题状态，并将判题数据转入数据库
			continue
		}
		if getType != expectedType { // sql语句类型不等，返回错误
			// TODO: 删除cache中对应的判题状态，并将判题数据转入数据库
			continue
		}

		wg.Wait() // 若cache中状态修改还未完成则等待完成，避免出错
		// TODO: 删除cache中对应的判题状态，并将判题数据转入数据库
	}
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
