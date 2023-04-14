package model

import (
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
