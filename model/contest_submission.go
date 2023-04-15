package model

import (
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
	contestSubmissionDAO := ContestSubmission{
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
