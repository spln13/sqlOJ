package model

import (
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type SubmitHistory struct {
	ID            int64 `gorm:"primary_key"`
	UserID        int64
	ExerciseID    int64
	UserType      int64
	Status        int
	StudentAnswer string
	UserAgent     string
	SubmitTime    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type SubmitHistoryFlow struct {
}

var (
	submitHistoryFlow *SubmitHistoryFlow
	submitHistoryOnce sync.Once
)

func NewSubmitHistoryFlow() *SubmitHistoryFlow {
	submitHistoryOnce.Do(func() {
		submitHistoryFlow = new(SubmitHistoryFlow)
	})
	return submitHistoryFlow
}

func (*SubmitHistoryFlow) InsertSubmitHistory(userID, exerciseID, userType int64, status int, studentAnswer, userAgent string, submitTime time.Time) {
	submitHistoryDAO := &SubmitHistory{
		UserID:        userID,
		ExerciseID:    exerciseID,
		UserType:      userType,
		Status:        status,
		UserAgent:     userAgent,
		StudentAnswer: studentAnswer,
		SubmitTime:    submitTime,
	}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(submitHistoryDAO).Error
	}); err != nil {
		log.Println(err)
	}
}
