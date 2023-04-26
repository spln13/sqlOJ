package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
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

func (*ContestSubmissionFlow) InsertContestSubmission(contestID, exerciseID, userID, userType int64, username, userAnswer, exerciseName, userAgent, contestName string, status int, submitTime time.Time) {
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
