package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type Class struct {
	ID           int64 `gorm:"primary_key"`
	Name         string
	StudentCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ClassFlow struct {
}

var (
	classFlow *ClassFlow
	classOnce sync.Once
)

func NewClassFlow() *ClassFlow {
	classOnce.Do(func() {
		classFlow = new(ClassFlow)
	})
	return classFlow
}

func (*ClassFlow) InsertClass(className string) error {
	classDAO := &Class{Name: className}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(classDAO).Error
	}); err != nil {
		log.Println(err)
		return errors.New("保存班级信息错误")
	}
	return nil
}

func (*ClassFlow) QueryClassNameByClassID(classID int64) (string, error) {
	var classDAO Class
	if err := GetSysDB().Select("name").Where("id = ?", classID).Find(&classDAO).Error; err != nil {
		log.Println(err)
		return "", errors.New("查询班级名错误")
	}
	return classDAO.Name, nil
}

func (*ClassFlow) IncreaseStudentCountInClass(classID int64, toAdd int) error {
	err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Class{}).Where("id = ?", classID).
			Update("student_count", gorm.Expr("student_count + ?", toAdd)).Error
		return err
	})
	if err != nil {
		log.Println(err)
		return errors.New("更新班级人数错误")
	}
	return nil
}

func (*ClassFlow) QueryClassIDNameMap() (map[int64]string, error) {
	var classDAOList []Class
	classIDNameMap := make(map[int64]string)
	err := GetSysDB().Model(&Class{}).Select("id", "name").Find(&classDAOList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("查询班级名错误")
	}
	for _, class := range classDAOList {
		classID := class.ID
		className := class.Name
		classIDNameMap[classID] = className
	}
	return classIDNameMap, nil
}
