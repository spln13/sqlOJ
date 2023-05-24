package fabric

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Submission describes basic details of what makes up a simple asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
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

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// submissions := []Submission{
	// 	{UserID: 1, UserType: 1, Number: 2019040471, ExerciseID: 1, ContestID: 0, Grade: 7, Status: 1},
	// 	{UserID: 2, UserType: 1, Number: 2019040472, ExerciseID: 1, ContestID: 1, Grade: 7, Status: 1},
	// 	{UserID: 3, UserType: 1, Number: 2019040473, ExerciseID: 1, ContestID: 1, Grade: 7, Status: 1},
	// 	{UserID: 4, UserType: 1, Number: 2019040474, ExerciseID: 2, ContestID: 0, Grade: 7, Status: 2},
	// 	{UserID: 5, UserType: 1, Number: 2019040475, ExerciseID: 2, ContestID: 0, Grade: 7, Status: 3},
	// 	{UserID: 6, UserType: 1, Number: 2019040476, ExerciseID: 3, ContestID: 0, Grade: 7, Status: 1},
	// 	{UserID: 7, UserType: 1, Number: 2019040471, ExerciseID: 2, ContestID: 1, Grade: 7, Status: 1},
	// 	{UserID: 7, UserType: 1, Number: 2019040471, ExerciseID: 2, ContestID: 0, Grade: 7, Status: 2},
	// 	{UserID: 7, UserType: 1, Number: 2019040471, ExerciseID: 2, ContestID: 0, Grade: 7, Status: 2},
	// }
	submissions := []Submission{
		{UserID: 1, UserType: 1, Number: 2019040471, ExerciseID: 1, ContestID: 0, Grade: 200, Status: 1},
		{UserID: 1, UserType: 1, Number: 2019040471, ExerciseID: 1, ContestID: 1, Grade: 150, Status: 1},
		{UserID: 2, UserType: 1, Number: 2019040472, ExerciseID: 1, ContestID: 0, Grade: 150, Status: 1},
		{UserID: 2, UserType: 1, Number: 2019040472, ExerciseID: 1, ContestID: 1, Grade: 150, Status: 1},
		{UserID: 3, UserType: 1, Number: 2019040473, ExerciseID: 1, ContestID: 0, Grade: 200, Status: 1},
		{UserID: 3, UserType: 1, Number: 2019040473, ExerciseID: 1, ContestID: 1, Grade: 200, Status: 1},
	}

	for idx, submission := range submissions {
		submissionJSON, err := json.Marshal(submission)
		if err != nil {
			return err
		}
		idxStr := strconv.Itoa(idx)
		submissionID := "submission_" + idxStr
		err = ctx.GetStub().PutState(submissionID, submissionJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
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
	var submissionKey string
	if contestID != 0 {
		submissionKey = "submission_contest_" + submissionIDStr
	} else {
		submissionKey = "submission_" + submissionIDStr
	}
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

// QueryAllSubmission returns all submissions found in world state
func (s *SmartContract) QueryAllSubmission(ctx contractapi.TransactionContextInterface) ([]QuerySubmissionResult, error) {
	//startKey := "submission_1"
	//endKey := "submission_999999"

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")

	if err != nil {
		return nil, err
	}
	defer func(resultsIterator shim.StateQueryIteratorInterface) {
		err := resultsIterator.Close()
		if err != nil {
			log.Println(err)
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

type StudentRating struct {
	Number        int64   `json:"number"`
	ExerciseScore float64 `json:"exercise_score"`
	ContestScore  float64 `json:"contest_score"`
	Score         float64 `json:"score"`
}

// RatingStudents 用于给出学生评分，评分细则见doc
func (s *SmartContract) RatingStudents(ctx contractapi.TransactionContextInterface) ([]StudentRating, error) {
	allSubmissionList, err := s.QueryAllSubmission(ctx) // 获取所有提交记录
	if err != nil {
		fmt.Printf("Error create chaincode: %s", err.Error())
		return nil, err
	}
	userScoreMap := make(map[int64][]int64) // 用户得分Map
	// 规定userScoreMap的value中[exercise_score, contest_score]
	for _, submissionRecord := range allSubmissionList {
		submission := submissionRecord.Record
		number := submission.Number
		if status := submission.Status; status != 1 { // 若提交状态不为正确则跳过
			continue
		}
		if userType := submission.UserType; userType != 1 { // 提交者不是学生则跳过
			continue
		}
		if submission.ContestID == 0 { // 题库中做题记录
			if len(userScoreMap[submission.Number]) > 0 { // 当前用户已有在规定userScore中有记录
				userScoreMap[number][0] += submission.Grade // 增加题库得分
			} else { // 当前用户还未在规定userScore中有记录
				userScoreMap[number] = append(userScoreMap[number], submission.Grade) // 加入题库成绩
				userScoreMap[number] = append(userScoreMap[number], 0)                // 加入竞赛成绩
			}
			// FIXME: 需要判断该提交记录是否以及加过分
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
	var studentGradeList []StudentRating // 用于存放学生学号和成绩的二维数组
	var maxExerciseScore int64 = 0       // 用于维护用户在题库中取得的最大得分
	var maxContestScore int64 = 0        // 用于维护用户在竞赛中取得的最大得分
	for _, number := range numbers {
		exerciseScore := userScoreMap[number][0]
		if exerciseScore > maxExerciseScore { // 维护用户在题库中取得的最大得分
			maxExerciseScore = exerciseScore
		}
		contestScore := userScoreMap[number][1]
		if contestScore > maxContestScore { // 维护用户在竞赛中取得的最大得分
			maxContestScore = contestScore
		}
		studentScore := StudentRating{
			Number:        number,
			ExerciseScore: float64(exerciseScore),
			ContestScore:  float64(contestScore),
			Score:         0,
		}
		studentGradeList = append(studentGradeList, studentScore)
	}
	// 已获取到maxExerciseScore, maxContestScore 根据文档中评分细则进行评分
	// max_grade = max(100, min(300, max_exercise_score))
	upperLimitExerciseGrade := max(100, min(300, maxExerciseScore)) // 最低100,最高300
	upperLimitContestGrade := max(100, min(300, maxContestScore))
	// exercise_score = (x / max_grade) * 100
	// contest_score = (x / max_grade) * 100

	for idx := range studentGradeList {
		studentGradeList[idx].ExerciseScore /= float64(upperLimitExerciseGrade)
		studentGradeList[idx].ContestScore /= float64(upperLimitContestGrade)
		studentGradeList[idx].ExerciseScore *= 100
		studentGradeList[idx].ContestScore *= 100
		// 智能合约给出的用户评分为: exercise_score * 0.3 + contest_score * 0.7
		studentGradeList[idx].Score = 0.3*studentGradeList[idx].ExerciseScore + 0.7*studentGradeList[idx].ContestScore
	}

	//studentGradeJSON, _ := json.Marshal(studentGradeList)
	return studentGradeList, nil
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
