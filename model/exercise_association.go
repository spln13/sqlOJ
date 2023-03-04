package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type ExerciseAssociation struct {
	ID         uint `gorm:"primary_key"`
	ExerciseID uint
	TableID    uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type ExerciseAssociationFlow struct {
}

var (
	exerciseAssociationFlow *ExerciseAssociationFlow
	exerciseAssociationOnce sync.Once
)

func NewExerciseAssociationFlow() *ExerciseAssociationFlow {
	exerciseAssociationOnce.Do(func() {
		exerciseAssociationFlow = new(ExerciseAssociationFlow)
	})
	return exerciseAssociationFlow
}

// InsertExerciseAssociation 将exerciseID与tableIDList的关系插入表中
func (*ExerciseAssociationFlow) InsertExerciseAssociation(exerciseID uint, tableIDList []uint) error {
	for _, tableID := range tableIDList {
		exerciseAssociationDAO := &ExerciseAssociation{ExerciseID: exerciseID, TableID: tableID}
		if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
			return tx.Create(exerciseAssociationDAO).Error
		}); err != nil {
			log.Println(err)
			return errors.New("添加关联表单失败")
		}
	}
	return nil
}
