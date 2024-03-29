package exercise

import (
	"github.com/xwb1989/sqlparser"
	"sqlOJ/fabric"
	"sqlOJ/model"
	"strconv"
	"strings"
)

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

// replaceTableName 替换原sql语句中的表名为临时表名，并将替换后的sql语句与原表名返回
// 此时getType只可能3个取值: 2->Insert 3->Update 4->Delete
// 因为经过sqlparser处理处理，这里的sql语句语法正确，替换过程中不会出现数组越界等错误
func replaceTableName(sql string, tempTableName string, getType int) (string, string) {
	sqlSplit := strings.Split(sql, " ")
	var originTableName string
	if getType == 2 { // Insert
		for idx, str := range sqlSplit {
			if str == "into" {
				originTableName = sqlSplit[idx+1]
				sqlSplit[idx+1] = tempTableName
				break
			}
		}
	} else if getType == 3 { // Update
		for idx, str := range sqlSplit {
			if str == "update" {
				originTableName = sqlSplit[idx+1]
				sqlSplit[idx+1] = tempTableName
				break
			}
		}
	} else { // Delete
		for idx, str := range sqlSplit {
			if str == "from" {
				originTableName = sqlSplit[idx+1]
				sqlSplit[idx+1] = tempTableName
				break
			}
		}
	}
	modifiedSql := strings.Join(sqlSplit, " ")
	return modifiedSql, originTableName
}

func appendPendingTxQueue(userID, userType, submissionID, contestID, exerciseID int64, status, grade int) {
	// 首先将参数转换为string型
	var number string
	if userType != 1 { // 非学生用户提交
		number = "0"
	} else {
		number = model.NewStudentAccountFlow().QueryStudentNumberByID(userID)
	}
	userIDStr := strconv.FormatInt(userID, 10)
	userTypeStr := strconv.FormatInt(userType, 10)
	submissionIDStr := strconv.FormatInt(submissionID, 10)
	exerciseIDStr := strconv.FormatInt(exerciseID, 10)
	contestIDStr := strconv.FormatInt(contestID, 10)
	statusStr := strconv.Itoa(status)
	gradeStr := strconv.Itoa(grade)
	ledgerData := fabric.LedgerData{
		SubmissionID: submissionIDStr,
		UserID:       userIDStr,
		UserType:     userTypeStr,
		ExerciseID:   exerciseIDStr,
		ContestID:    contestIDStr,
		Status:       statusStr,
		Grade:        gradeStr,
		Number:       number,
	}
	fabric.PendingTxQueue <- ledgerData // 将需要上链的数据写入上链队列

}
