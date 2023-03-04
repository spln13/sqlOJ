package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type ExerciseTable struct {
	ID               uint `gorm:"primary_key"`
	PublishID        uint
	Name             string
	Description      string
	AssociationCount int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ExerciseTableFlow struct {
}

var (
	exerciseTableFlow *ExerciseTableFlow
	exerciseTableOnce sync.Once
)

func NewExerciseTableFlow() *ExerciseTableFlow {
	exerciseTableOnce.Do(func() {
		exerciseTableFlow = new(ExerciseTableFlow)
	})
	return exerciseTableFlow
}

func (*ExerciseTableFlow) InsertExerciseTable(publishID uint, name, description string) error {
	exerciseTableDAO := &ExerciseTable{
		PublishID:   publishID,
		Name:        name,
		Description: description,
	}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(exerciseTableDAO).Error
	}); err != nil {
		log.Println(err)
		return errors.New("保存记录错误")
	}
	return nil
}

// QueryExerciseTableExist 查询当前表名是否重复
func (*ExerciseTableFlow) QueryExerciseTableExist(name string) (bool, error) {
	var exerciseTableDAO ExerciseTable
	if err := GetSysDB().Select("id").Where("name = ?", name).Find(&exerciseTableDAO).Error; err != nil {
		log.Println(err)
		return false, errors.New("查询表名错误")
	}
	if exerciseTableDAO.ID == 0 {
		return false, nil
	}
	return true, nil
}

// IncreaseExerciseTableAssociationCount 自增exercise_tables中的association_count
func (*ExerciseTableFlow) IncreaseExerciseTableAssociationCount(tableIDList []uint) error {
	var errReturn error
	for _, tableID := range tableIDList {
		if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
			return tx.Model(ExerciseTable{}).Where("id = ?", tableID).Update("association_count", gorm.Expr("association_count + 1")).Error
		}); err != nil {
			log.Println(err)
			errReturn = err
		}
	}
	return errReturn
}
