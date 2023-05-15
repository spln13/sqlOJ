package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
	"time"
)

type ContestSubmission struct {
	ID           int64 `gorm:"primary_key"`
	ContestID    int64
	ExerciseID   int64
	UserID       int64
	UserType     int64
	Username     string
	UserAnswer   string
	ExerciseName string
	UserAgent    string
	ContestName  string
	Status       int
	OnChain      int
	SubmitTime   time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ContestSubmissionFlow struct {
}

var (
	contestSubmissionFlowDAO  *ContestSubmissionFlow
	contestSubmissionFlowOnce sync.Once
)

func NewContestSubmissionFlow() *ContestSubmissionFlow {
	contestSubmissionFlowOnce.Do(func() {
		contestSubmissionFlowDAO = new(ContestSubmissionFlow)
	})
	return contestSubmissionFlowDAO
}

func (*ContestSubmissionFlow) InsertContestSubmission(contestID, exerciseID, userID, userType int64, username, userAnswer, exerciseName, userAgent, contestName string, status int, submitTime time.Time) int64 {
	contestSubmissionDAO := &ContestSubmission{
		ExerciseID:   exerciseID,
		UserID:       userID,
		UserType:     userType,
		Username:     username,
		UserAnswer:   userAnswer,
		ExerciseName: exerciseName,
		UserAgent:    userAgent,
		ContestID:    contestID,
		ContestName:  contestName,
		Status:       status,
		SubmitTime:   submitTime,
	}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(contestSubmissionDAO).Error
	}); err != nil {
		log.Println(err)
	}
	return contestSubmissionDAO.ID
}

func (*ContestSubmissionFlow) GetContestSubmissionByID(contestID int64) ([]ContestSubmission, error) {
	var contestSubmissionList []ContestSubmission
	err := GetSysDB().Model(&ContestSubmission{}).Where("contest_id = ?", contestID).Omit("create_at", "update_at").Find(&contestSubmissionList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("查询竞赛提交错误")
	}
	return contestSubmissionList, nil
}

func (*ContestSubmissionFlow) GetUserContestSubmission(userID, userType, contestID int64) ([]ContestSubmission, error) {
	var contestSubmissionList []ContestSubmission
	err := GetSysDB().Model(&ContestSubmission{}).Where("contest_id = ? and user_id = ? and user_type = ?", contestID, userID, userType).Omit("create_at", "update_at").Find(&contestSubmissionList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("查询竞赛提交错误")
	}
	return contestSubmissionList, nil
}

func (*ContestSubmissionFlow) GetOneExerciseSubmission(contestID, exerciseID int64) ([]ContestSubmission, error) {
	var contestSubmissionList []ContestSubmission
	err := GetSysDB().Model(&ContestSubmission{}).Where("contest_id = ? and exercise_id = ?", contestID, exerciseID).Omit("create_at", "update_at").Find(&contestSubmissionList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("查询竞赛提交错误")
	}
	return contestSubmissionList, nil
}

func (*ContestSubmissionFlow) QueryOneSubmissionAnswer(userID, userType, submissionID int64) (string, error) {
	var contestSubmission ContestSubmission
	err := GetSysDB().Model(&ContestSubmission{}).Select("user_id", "user_type", "user_answer").
		Where("id = ?", submissionID).Find(&contestSubmission).Error
	if err != nil {
		return "", errors.New("查询提交答案错误")
	}
	if contestSubmission.UserID != userID || contestSubmission.UserType != userType {
		return "", errors.New("无权限")
	}
	return contestSubmission.UserAnswer, nil
}

type MinContestSubmission struct {
	OnChain    int
	Status     int
	SubmitTime time.Time
}

func (*ContestSubmissionFlow) QueryOneUserExerciseSubmission(userID, userType, contestID, exerciseID int64) ([]MinContestSubmission, error) {
	var minContestSubmissionList []MinContestSubmission
	err := GetSysDB().Model(&ContestSubmission{}).
		Where("user_id = ? and user_type = ? and contest_id = ? and exercise_id = ?", userID, userType, contestID, exerciseID).
		Find(&minContestSubmissionList).Error
	if err != nil {
		return nil, errors.New("查询提交记录错误")
	}
	return minContestSubmissionList, nil
}

func (*ContestSubmissionFlow) ModifyContestSubmissionOnChain(submissionIDStr string) {
	submissionID, _ := strconv.ParseInt(submissionIDStr, 10, 64)
	err := GetSysDB().Model(&ContestSubmission{}).Where("id = ?", submissionID).Update("on_chain", 1).Error
	if err != nil {
		log.Println(err)
	}
}
