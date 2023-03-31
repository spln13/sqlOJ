package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

// UserProblemStatus 表示用户与对应题目的做题数据
type UserProblemStatus struct {
	ID         int64 `gorm:"primary_key"`
	UserID     int64
	ExerciseID int64
	UserType   int64
	Status     int // 1->ac; 2->wa
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UserProblemStatusMin struct {
	ExerciseID int64
	Status     int
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

func (*UserProblemStatusFlow) ModifyUserProblemStatus(userID, exerciseID, userType int64, status int) {
	var userProblemStatusDAO UserProblemStatus
	if err := GetSysDB().Select("ID").Where("user_id = ? and exercise_id = ?", userID, exerciseID).Find(&userProblemStatusDAO).Error; err != nil {
		log.Println(err)
	}
	if userProblemStatusDAO.ID == 0 { // 其中没有记录
		userProblemStatusDAO := &UserProblemStatus{
			UserID:     userID,
			ExerciseID: exerciseID,
			UserType:   userType,
			Status:     status,
		}
		if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
			return tx.Create(userProblemStatusDAO).Error
		}); err != nil {
			log.Println(err)
		}
	} else {
		if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
			return tx.Where("user_id = ? and exercise_id = ? and user_type = ?", userID, exerciseID, userType).Update("status", status).Error
		}); err != nil {
			log.Println(err)
		}
	}
}

// QueryUserProblemStatus 查询用户对应所有的做过的题以及状态
func (*UserProblemStatusFlow) QueryUserProblemStatus(userID, userType int64) (map[int64]int, error) {
	var userProblemStatusMinList []UserProblemStatusMin
	err := GetSysDB().Model(&UserProblemStatus{}).Where("user_id = ? and user_type = ?", userID, userType).Find(&userProblemStatusMinList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("查询用户做题数据错误")
	}
	problemStatusMap := make(map[int64]int)
	for _, userProblemStatusMin := range userProblemStatusMinList {
		problemStatusMap[userProblemStatusMin.ExerciseID] = userProblemStatusMin.Status
	}
	return problemStatusMap, nil
}
