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
	Username      string
	ExerciseName  string
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

func (*SubmitHistoryFlow) InsertSubmitHistory(userID, exerciseID, userType int64, status int, studentAnswer, userAgent, username, exerciseName string, submitTime time.Time) {
	submitHistoryDAO := &SubmitHistory{
		UserID:        userID,
		ExerciseID:    exerciseID,
		UserType:      userType,
		Status:        status,
		Username:      username,
		ExerciseName:  exerciseName,
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
func (*SubmitHistoryFlow) QueryAllSubmitHistory() ([]SubmitHistory, error) {
	var submitHistory []SubmitHistory
	if err := GetSysDB().Model(&SubmitHistory{}).Omit("create_at", "update_at").Find(&submitHistory).Error; err != nil {
		log.Println(err)
		return nil, errors.New("查询提交记录错误")
	}
	return submitHistory, nil
}

// QueryThisExerciseSubmitHistory 查询当前习题所有的提交记录
func (*SubmitHistoryFlow) QueryThisExerciseSubmitHistory(exerciseID int64) ([]SubmitHistory, error) {
	var submitHistory []SubmitHistory
	if err := GetSysDB().Model(&SubmitHistory{}).Where("exercise_id = ?", exerciseID).Omit("create_at", "update_at", "exercise_id").Find(&submitHistory).Error; err != nil {
		log.Println(err)
		return nil, errors.New("查询提交记录错误")
	}
	return submitHistory, nil
}

// QueryThisExerciseUserSubmitHistory 查询userID, userType, exerciseID对应的提交记录
func (*SubmitHistoryFlow) QueryThisExerciseUserSubmitHistory(userID int64, userType int64, exerciseID int64) ([]SubmitHistory, error) {
	var submitHistory []SubmitHistory
	if err := GetSysDB().Model(&SubmitHistory{}).Select("student_answer", "status", "submit_time").Where("user_id = ? and exercise_id = ? and user_type = ?", userID, exerciseID, userType).Find(&submitHistory).Error; err != nil {
		log.Println(err)
		return nil, errors.New("查询提交记录错误")
	}
	return submitHistory, nil
}

func (*SubmitHistoryFlow) QueryThisUserSubmitHistory(userID int64) ([]SubmitHistory, error) {
	var submitHistory []SubmitHistory
	if err := GetSysDB().Model(&SubmitHistory{}).Select("student_answer", "exercise_id", "status", "submit_time", "user_agent").Where("user_id = ?", userID).Find(&submitHistory).Error; err != nil {
		log.Println(err)
		return nil, errors.New("查询提交记录错误")
	}
	return submitHistory, nil
}
