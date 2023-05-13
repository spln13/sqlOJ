package fabric

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"sort"
	"strconv"
)

import (
	"encoding/json"
)

// SmartContract provides functions for managing a submission record
type SmartContract struct {
	contractapi.Contract
}

type Submission struct {
	UserID     int64 `json:"user_id"`
	UserType   int64 `json:"user_type"`
	Number     int64 `json:"number"`
	ExerciseID int64 `json:"exercise_id"`
	ContestID  int64 `json:"contest_id"`
	Grade      int64 `json:"grade"`
	Status     int64 `json:"status"`
}

// QuerySubmissionResult structure used for handling result of query
type QuerySubmissionResult struct {
	Key    string `json:"Key"`
	Record *Submission
}

// CreateSubmission adds a new submission to the world state with given details
func (s *SmartContract) CreateSubmission(ctx contractapi.TransactionContextInterface, submissionID, userID, userType, exerciseID, contestID, status, grade, number int64) error {
	submission := Submission{
		UserID:     userID,
		UserType:   userType,
		ExerciseID: exerciseID,
		ContestID:  contestID,
		Status:     status,
		Grade:      grade,
		Number:     number,
	}

	submissionAsBytes, _ := json.Marshal(submission)
	submissionIDStr := strconv.FormatInt(submissionID, 10)
	submissionKey := "submission_" + submissionIDStr
	return ctx.GetStub().PutState(submissionKey, submissionAsBytes)
}

// QuerySubmission returns the submission stored in the world state with given id
func (s *SmartContract) QuerySubmission(ctx contractapi.TransactionContextInterface, submissionID int64) (*Submission, error) {
	submissionIDStr := strconv.FormatInt(submissionID, 10)
	submissionKey := "submission_" + submissionIDStr

	submissionAsBytes, err := ctx.GetStub().GetState(submissionKey)

	if err != nil {
		return nil, fmt.Errorf("failed to read from world state. %s", err.Error())
	}

	if submissionAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", submissionKey)
	}

	submission := new(Submission)
	_ = json.Unmarshal(submissionAsBytes, submission)

	return submission, nil
}

// ChangeSubmissionStatus updates the status field of submission with given id in world state
func (s *SmartContract) ChangeSubmissionStatus(ctx contractapi.TransactionContextInterface, submissionID, newStatus int64) error {
	submission, err := s.QuerySubmission(ctx, submissionID)
	if err != nil {
		return err
	}

	submission.Status = newStatus
	submissionAsBytes, _ := json.Marshal(submission)
	submissionIDStr := strconv.FormatInt(submissionID, 10)
	submissionKey := "submission_" + submissionIDStr
	return ctx.GetStub().PutState(submissionKey, submissionAsBytes)
}

// QueryAllSubmission returns all submissions found in world state
func (s *SmartContract) QueryAllSubmission(ctx contractapi.TransactionContextInterface) ([]QuerySubmissionResult, error) {
	startKey := "submission_1"
	endKey := "submission_999999"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer func(resultsIterator shim.StateQueryIteratorInterface) {
		err := resultsIterator.Close()
		if err != nil {

		}
	}(resultsIterator)

	var results []QuerySubmissionResult

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		submission := new(Submission)
		_ = json.Unmarshal(queryResponse.Value, submission)

		queryResult := QuerySubmissionResult{Key: queryResponse.Key, Record: submission}
		results = append(results, queryResult)
	}

	return results, nil
}

func (s *SmartContract) RatingStudents(ctx contractapi.TransactionContextInterface) []byte {
	allSubmissionList, err := s.QueryAllSubmission(ctx) // 获取所有提交记录
	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return nil
	}
	userScoreMap := make(map[int64][]int64) // 用户得分Map
	// 规定userScoreMap的value中[exercise_score, contest_score]
	for _, submissionRecord := range allSubmissionList {
		submission := submissionRecord.Record
		number := submission.Number
		if submission.ContestID == 0 { // 题库中做题记录
			if len(userScoreMap[submission.Number]) > 0 { // 当前用户已有在规定userScore中有记录
				userScoreMap[number][0] += submission.Grade // 增加题库得分
			} else { // 当前用户还未在规定userScore中有记录
				userScoreMap[number] = append(userScoreMap[number], submission.Grade) // 加入题库成绩
				userScoreMap[number] = append(userScoreMap[number], 0)                // 加入竞赛成绩
			}

		} else { // 竞赛中做题记录
			if len(userScoreMap[submission.Number]) > 0 { // 当前用户已有在规定userScore中有记录
				userScoreMap[number][1] += submission.Grade // 增加竞赛得分
			} else { // 当前用户还未在规定userScore中有记录
				userScoreMap[number] = append(userScoreMap[number], 0)                // 加入题库成绩
				userScoreMap[number] = append(userScoreMap[number], submission.Grade) // 加入竞赛成绩
			}
		}
	}

	var numbers []int64 // 存储有提交记录的学生学号
	for key := range userScoreMap {
		numbers = append(numbers, key)
	}
	sort.Slice(numbers, func(i, j int) bool { // 对学号从小到大排序
		return numbers[i] < numbers[j]
	})
	var studentGrade [][]int64 // 用于存放学生学号和成绩的二维数组
	for _, number := range numbers {
		exerciseScore := userScoreMap[number][0]
		contestScore := userScoreMap[number][1]
		record := []int64{number, exerciseScore, contestScore} // 二维数组结构
		studentGrade = append(studentGrade, record)
	}
	studentGradeJSON, _ := json.Marshal(studentGrade)
	return studentGradeJSON
}

func InitChainCode() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
