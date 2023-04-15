package model

import (
	"errors"
	"log"
	"sync"
	"time"
)

type ContestExerciseAssociation struct {
	ID         int64 `gorm:"primary_key"`
	ContestID  int64
	ExerciseID int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ContestExerciseAssociationFlow struct {
}

var (
	contestExerciseAssociationFlowDAO  *ContestExerciseAssociationFlow
	contestExerciseAssociationFlowOnce sync.Once
)

func NewContestExerciseAssociationFlow() *ContestExerciseAssociationFlow {
	contestExerciseAssociationFlowOnce.Do(func() {
		contestExerciseAssociationFlowDAO = new(ContestExerciseAssociationFlow)
	})
	return contestExerciseAssociationFlowDAO
}

func (*ContestExerciseAssociationFlow) InsertContestExerciseAssociation(contestID int64, exerciseIDList []int64) error {
	for _, exerciseID := range exerciseIDList {
		contestExerciseAssociationDAO := ContestExerciseAssociation{ContestID: contestID, ExerciseID: exerciseID}
		if err := GetSysDB().Create(contestExerciseAssociationDAO).Error; err != nil {
			log.Println(err)
			return errors.New("插入竞赛题目关联错误")
		}
	}
	return nil
}

func (*ContestExerciseAssociationFlow) GetExerciseIDListByContestID(contestID int64) ([]int64, error) {
	var contestExerciseAssociationList []ContestExerciseAssociation
	err := GetSysDB().Select("exercise_id").Where("contest_id = ?", contestID).Find(&contestExerciseAssociationList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("查询竞赛中题目时错误")
	}
	exerciseIDList := make([]int64, len(contestExerciseAssociationList))
	for idx, contestExerciseAssociation := range contestExerciseAssociationList {
		exerciseIDList[idx] = contestExerciseAssociation.ExerciseID
	}
	return exerciseIDList, nil
}
