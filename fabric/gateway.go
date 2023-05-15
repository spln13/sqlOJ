package fabric

import (
	"encoding/json"
	"fmt"
	"log"
	"sqlOJ/model"
)

type LedgerData struct {
	SubmissionID string
	UserID       string
	UserType     string
	ExerciseID   string
	ContestID    string
	Status       string // 提交状态
	Grade        string // 题目对应分数
	Number       string // 学生学号
}

var PendingTxQueue = make(chan LedgerData, 3000) // 上链队列 缓冲区大小3000

func PushOnChain() { // 上链协程
	for {
		data := <-PendingTxQueue // 从上链队列
		CreateSubmission(data.SubmissionID, data.UserID, data.UserType, data.ExerciseID, data.ContestID, data.Status, data.Grade, data.Number)
		// 将提交已上链的状态更新到MySQL, 区分题库提交和竞赛提交
		if data.ContestID == "0" { // 题库提交
			model.NewSubmitHistoryFlow().ModifySubmissionOnChain(data.SubmissionID) // 修改题库提交记录上链状态为已上链
		} else { // 竞赛提交
			model.NewContestSubmissionFlow().ModifyContestSubmissionOnChain(data.SubmissionID) // 修改竞赛提交记录上链状态为未上链
		}
	}
}

// CreateSubmission 将学生提交信息上链存储
// Submit a transaction synchronously, blocking until it has been committed to the ledger.

func CreateSubmission(submissionID, userID, userType, exerciseID, contestID, status, grade, number string) {
	// 将参数
	_, err := contract.SubmitTransaction("CreateSubmission", submissionID, userID, userType, exerciseID, contestID, status, grade, number)
	if err != nil {
		log.Println(err)
	}
}

type StudentScore struct {
	Number        int64   `json:"number"`
	ExerciseScore float64 `json:"exercise_score"`
	ContestScore  float64 `json:"contest_score"`
	Score         float64 `json:"score"`
}

// RatingStudents 调用评分链码获取学生做题情况
func RatingStudents() ([]StudentScore, error) {
	initLedger(contract)
	ratingResultBytes, err := contract.EvaluateTransaction("RatingStudents")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var ratingResult []StudentScore
	if err = json.Unmarshal(ratingResultBytes, &ratingResult); err != nil {
		log.Println(err)
		return nil, err
	}
	fmt.Println(ratingResult)
	return ratingResult, nil
}
