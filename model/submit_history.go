package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
	"time"
)

type SubmitHistory struct {
	ID            int64 `gorm:"primary_key"`
	UserID        int64
	ExerciseID    int64
	UserType      int64
	Status        int
	OnChain       int // 0: 待上链; 1: 已上链
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

func (*SubmitHistoryFlow) InsertSubmitHistory(userID, exerciseID, userType int64, status int, studentAnswer, userAgent, username, exerciseName string, submitTime time.Time) int64 {
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
	return submitHistoryDAO.ID
}

// QueryAllSubmitHistory 查询所有的提交记录
func (*SubmitHistoryFlow) QueryAllSubmitHistory() ([]SubmitHistory, error) {
	var submitHistory []SubmitHistory
	if err := GetSysDB().Model(&SubmitHistory{}).Omit("create_at", "update_at").Order("id desc").Find(&submitHistory).Error; err != nil {
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
	if err := GetSysDB().Model(&SubmitHistory{}).Select("student_answer", "status", "submit_time", "on_chain").Where("user_id = ? and exercise_id = ? and user_type = ?", userID, exerciseID, userType).Order("id desc").Find(&submitHistory).Error; err != nil {
		log.Println(err)
		return nil, errors.New("查询提交记录错误")
	}
	return submitHistory, nil
}

func (*SubmitHistoryFlow) QueryThisUserSubmitHistory(userID, userType int64) ([]SubmitHistory, error) {
	var submitHistory []SubmitHistory
	if err := GetSysDB().Model(&SubmitHistory{}).Select("id", "student_answer", "exercise_id", "status", "submit_time", "user_agent", "on_chain").Where("user_id = ? and user_type = ?", userID, userType).Order("id desc").Find(&submitHistory).Error; err != nil {
		log.Println(err)
		return nil, errors.New("查询提交记录错误")
	}
	return submitHistory, nil
}

func (*SubmitHistoryFlow) QuerySubmissionAnswer(submissionID int64) (int64, int64, string, error) {
	var submitHistory SubmitHistory
	err := GetSysDB().Model(&SubmitHistory{}).Select("user_id", "user_type", "student_answer").Where("id = ?", submissionID).Find(&submitHistory).Error
	if err != nil {
		return 0, 0, "", errors.New("查询提交信息错误")
	}
	return submitHistory.UserID, submitHistory.UserType, submitHistory.StudentAnswer, nil
}

func (*SubmitHistoryFlow) ModifySubmissionOnChain(submissionIDStr string) {
	submissionID, _ := strconv.ParseInt(submissionIDStr, 10, 64)
	err := GetSysDB().Model(&SubmitHistory{}).Where("id = ?", submissionID).Update("on_chain", 1).Error
	if err != nil {
		log.Println(err)
	}
}
