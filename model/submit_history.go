package model

import (
	"errors"
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

// QueryAllSubmitHistory 查询所有的提交记录
func QueryAllSubmitHistory() ([]SubmitHistory, error) {
	var submitHistory []SubmitHistory
	if err := GetSysDB().Model(&SubmitHistory{}).Omit("create_at", "update_at").Find(&submitHistory).Error; err != nil {
		log.Println(err)
		return nil, errors.New("查询提交记录错误")
	}
	return submitHistory, nil
}

// QueryThisExerciseSubmitHistory 查询当前习题所有的提交记录
func QueryThisExerciseSubmitHistory(exerciseID int64) ([]SubmitHistory, error) {
	var submitHistory []SubmitHistory
	if err := GetSysDB().Model(&SubmitHistory{}).Where("exercise_id = ?", exerciseID).Omit("create_at", "update_at", "exercise_id").Find(&submitHistory).Error; err != nil {
		log.Println(err)
		return nil, errors.New("查询提交记录错误")
	}
	return submitHistory, nil
}
