package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type ContestExerciseStatus struct {
	ID         int64
	ContestID  int64
	ExerciseID int64
	UserID     int64
	UserType   int64
	Status     int // 1->ac; 2->wa; 3->re
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ContestExerciseStatusFlow struct {
}

var (
	contestExerciseStatusFlowDAO  *ContestExerciseStatusFlow
	contestExerciseStatusFlowOnce sync.Once
)

func NewContestExerciseStatusFlow() *ContestExerciseStatusFlow {
	contestExerciseStatusFlowOnce.Do(func() {
		contestExerciseStatusFlowDAO = new(ContestExerciseStatusFlow)
	})
	return contestExerciseStatusFlowDAO
}

func (*ContestExerciseStatusFlow) ModifyContestExerciseStatus(userID, userType, exerciseID, contestID int64, status int) {
	var contestExerciseStatusDAO ContestExerciseStatus
	if err := GetSysDB().Model(&ContestExerciseStatus{}).Select("id", "status").Where("user_id = ? and user_type = ? and exercise_id = ? and contest_id = ?", userID, userType, exerciseID, contestID).Find(&contestExerciseStatusDAO).Error; err != nil {
		log.Println(err)
	}
	if contestExerciseStatusDAO.ID == 0 { // 其中没有记录，插入数据
		newContestExerciseStatusDAO := &ContestExerciseStatus{
			ContestID:  contestID,
			ExerciseID: exerciseID,
			UserID:     userID,
			UserType:   userType,
			Status:     status,
		}
		if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
			return tx.Create(newContestExerciseStatusDAO).Error
		}); err != nil {
			log.Println(err)
		}
		return
	}
	if contestExerciseStatusDAO.Status == 1 || contestExerciseStatusDAO.Status == status {
		// 此题AC过, 或者状态一致,不再更新
		return
	}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Where("user_id = ? and exercise_id = ? and contest_id = ? and user_type = ?", userID, exerciseID, contestID, userType).Update("status", status).Error
	}); err != nil {
		log.Println(err)
	}
}

type ContestExerciseStatusMin struct {
	ExerciseID int64
	Status     int
}

// QueryContestExerciseStatus 查询用户对应所有的做过的题以及状态
func (*ContestExerciseStatusFlow) QueryContestExerciseStatus(userID, userType, contestID int64) (map[int64]int, error) {
	var contestExerciseStatusMinList []ContestExerciseStatusMin
	err := GetSysDB().Model(&ContestExerciseStatus{}).Where("user_id = ? and user_type = ? and contest_id = ?", userID, userType, contestID).Find(&contestExerciseStatusMinList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("查询用户做题数据错误")
	}
	problemStatusMap := make(map[int64]int)
	for _, contestExerciseStatusMin := range contestExerciseStatusMinList {
		problemStatusMap[contestExerciseStatusMin.ExerciseID] = contestExerciseStatusMin.Status
	}
	return problemStatusMap, nil
}

func (*ContestExerciseStatusFlow) QueryStudentIDListByContestID(contestID int64) ([]int64, error) {
	var contestExerciseStatusList []ContestExerciseStatus
	err := GetSysDB().Model(&ContestExerciseStatus{}).Select("user_id").Where("contest_id = ? and user_type > 1", contestID).Find(&contestExerciseStatusList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("查询学生信息错误")
	}
	var studentIDList []int64
	for _, contestStatus := range contestExerciseStatusList {
		studentIDList = append(studentIDList, contestStatus.UserID)
	}
	return studentIDList, nil
}

func (*ContestExerciseStatusFlow) QueryStudentProblemStatusMap(userID, contestID int64) (map[int64]int, error) {
	var contestExerciseStatusList []ContestExerciseStatus
	err := GetSysDB().Model(&ContestExerciseStatus{}).Select("exercise_id", "status").
		Where("user_type = 1 and user_id = ? and contest_id = ?", userID, contestID).Find(&contestExerciseStatusList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("查询学生做题数据错误")
	}
	studentProblemStatusMap := make(map[int64]int)
	for _, contestStatus := range contestExerciseStatusList {
		studentProblemStatusMap[contestStatus.ExerciseID] = contestStatus.Status
	}
	return studentProblemStatusMap, nil
}
