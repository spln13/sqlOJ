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
	Status     int // 1->ac; 2->wa; 3->re
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

// ModifyUserProblemStatus 修改user_problem_statuses表中用户与题目对应的数据
func (*UserProblemStatusFlow) ModifyUserProblemStatus(userID, exerciseID, userType int64, status int) {
	var userProblemStatusDAO UserProblemStatus
	if err := GetSysDB().Model(&UserProblemStatus{}).Select("id", "status").Where("user_id = ? and exercise_id = ? and user_type = ?", userID, exerciseID, userType).Find(&userProblemStatusDAO).Error; err != nil {
		log.Println(err)
	}
	if userProblemStatusDAO.ID == 0 { // 其中没有记录，插入数据
		newUserProblemStatusDAO := &UserProblemStatus{
			UserID:     userID,
			ExerciseID: exerciseID,
			UserType:   userType,
			Status:     status,
		}
		if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
			return tx.Create(newUserProblemStatusDAO).Error
		}); err != nil {
			log.Println(err)
		}
		return
	}
	// 其中有记录
	if userProblemStatusDAO.Status == 1 || userProblemStatusDAO.Status == status {
		// 此题AC过, 做错不再更新, 或者此题状态相同
		return
	}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(&UserProblemStatus{}).Where("user_id = ? and exercise_id = ? and user_type = ?", userID, exerciseID, userType).Update("status", status).Error
	}); err != nil {
		log.Println(err)
	}
}

// QueryUserAllProblemStatus 查询用户对应所有的做过的题以及状态
func (*UserProblemStatusFlow) QueryUserAllProblemStatus(userID, userType int64) (map[int64]int, error) {
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

// QueryUserProblemStatus 查询当前用户在当前题目的提交状态0:err; 1: ac; 2: wa; 3: re; 4: 未提交过
func (*UserProblemStatusFlow) QueryUserProblemStatus(userID, userType, exerciseID int64) int {
	var userProblemStatusDAO UserProblemStatus
	err := GetSysDB().Model(&UserProblemStatus{}).Select("id", "status").
		Where("user_id = ? and user_type = ? and exercise_id = ?", userID, userType, exerciseID).
		Find(&userProblemStatusDAO).Error
	if err != nil {
		log.Println(err)
		return 0
	}
	if userProblemStatusDAO.ID == 0 { // 未提交过
		return 4
	}
	return userProblemStatusDAO.Status
}
