package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type ExerciseTable struct {
	ID               int64 `gorm:"primary_key"`
	PublisherID      int64
	PublisherType    int64
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

func (*ExerciseTableFlow) InsertExerciseTable(publisherID int64, publisherType int64, name, description string) error {
	exerciseTableDAO := &ExerciseTable{
		PublisherID:   publisherID,
		PublisherType: publisherType,
		Name:          name,
		Description:   description,
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
func (*ExerciseTableFlow) IncreaseExerciseTableAssociationCount(tableIDList []int64) error {
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

func (*ExerciseTableFlow) QueryAllTable() ([]ExerciseTable, error) {
	var tableList []ExerciseTable
	err := GetSysDB().Model(&ExerciseTable{}).Omit("created_at", "updated_at").Find(&tableList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("获取数据表信息错误")
	}
	return tableList, nil
}

func (*ExerciseTableFlow) QueryTableNameByID(tableID int64) (string, error) {
	var tableDAO ExerciseTable
	err := GetSysDB().Model(&ExerciseTable{}).Select("name").Where("id = ?", tableID).Find(&tableDAO).Error
	if err != nil {
		log.Println(err)
		return "", errors.New("查询数据表名称错误")
	}
	return tableDAO.Name, nil
}

func (*ExerciseTableFlow) DeleteTableByID(tableID int64) error {
	err := GetSysDB().Delete(&ExerciseTable{}, tableID).Error
	if err != nil {
		log.Println(err)
		return errors.New("删除数据表数据错误")
	}
	return nil
}
