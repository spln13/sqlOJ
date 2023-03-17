package model

import (
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type UserProblemStatus struct {
	ID         int64 `gorm:"primary_key"`
	UserID     int64
	ExerciseID int64
	Status     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UserProblemStatusFlow struct {
}

var (
	userProblemStatusFlow *UserProblemStatusFlow
	userProblemStatusOnce sync.Once
)

func NewUserProblemStatusFlow() *UserProblemStatusFlow {
	userProblemStatusOnce.Do(func() {
		userProblemStatusFlow = new(UserProblemStatusFlow)
	})
	return userProblemStatusFlow
}

func (*UserProblemStatusFlow) ModifyUserProblemStatus(userID, exerciseID int64, status int) {
	var userProblemStatusDAO UserProblemStatus
	if err := GetSysDB().Select("ID").Where("user_id = ? and exercise_id = ?", userID, exerciseID).Find(&userProblemStatusDAO).Error; err != nil {
		log.Println(err)
	}
	if userProblemStatusDAO.ID == 0 { // 其中没有记录
		userProblemStatusDAO := &UserProblemStatus{
			UserID:     userID,
			ExerciseID: exerciseID,
			Status:     status,
		}
		if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
			return tx.Create(userProblemStatusDAO).Error
		}); err != nil {
			log.Println(err)
		}
	} else {
		if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
			return tx.Where("user_id = ? and exercise_id = ?", userID, exerciseID).Update("status", status).Error
		}); err != nil {
			log.Println(err)
		}
	}
}
