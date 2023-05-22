package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type ExerciseAssociation struct {
	ID         int64 `gorm:"primary_key"`
	ExerciseID int64
	TableID    int64
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
func (*ExerciseAssociationFlow) InsertExerciseAssociation(exerciseID int64, tableIDList []int64) error {
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

// QueryAssociationExist 查找是否有题目与该数据表关联
func (*ExerciseAssociationFlow) QueryAssociationExist(tableID int64) (bool, error) {
	var exerciseAssociation ExerciseAssociation
	err := GetSysDB().Model(&ExerciseAssociation{}).Select("table_id = ?", tableID).Limit(1).Find(&exerciseAssociation).Error
	if err != nil {
		log.Println(err)
		return false, errors.New("查找题目数据表引用错误")
	}
	if exerciseAssociation.ID == 0 { // 不存在
		return false, nil
	}
	return true, nil
}

func (*ExerciseAssociationFlow) DeleteAssociation(exerciseID int64) error {
	err := GetSysDB().Delete(&ExerciseAssociationFlow{}).Where("exercise_id = ?", exerciseID).Error
	if err != nil {
		log.Println(err)
		return errors.New("删除关联关系错误")
	}
	return nil
}
