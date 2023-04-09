package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type ContestExercise struct {
	ID               int64 `gorm:"primary_key"`
	ContestID        int64
	OriginExerciseID int64
	EndAt            time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ContestExerciseFlow struct {
}

var (
	contestExerciseFlowDAO *ContestExerciseFlow
	contestExerciseOnce    sync.Once
)

func NewContestExerciseFlow() *ContestExerciseFlow {
	contestExerciseOnce.Do(func() {
		contestExerciseFlowDAO = new(ContestExerciseFlow)
	})
	return contestExerciseFlowDAO
}

// InsertContestExercise 向MySQL中插入副本表
func (*ContestExerciseFlow) InsertContestExercise(contestID int64, originExerciseIDList []int64, endAt time.Time) error {
	err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		for _, originExerciseID := range originExerciseIDList {
			contestExercise := ContestExercise{
				ContestID:        contestID,
				OriginExerciseID: originExerciseID,
				EndAt:            endAt,
			}
			if err := tx.Create(contestExercise).Error; err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return errors.New("生成副本题目错误")
	}
	return nil
}
