package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type Class struct {
	ID              int64 `gorm:"primary_key"`
	Name            string
	TeacherUsername string // 教职工号
	TeacherName     string
	StudentCount    int
	CreateBy        int64 // 创建者的id
	CreatedAt       time.Time
	UpdatedAt       time.Time
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

func (*ClassFlow) InsertClass(className string, teacherUsername string, teacherName string, createBy int64) error {
	classDAO := &Class{Name: className, TeacherUsername: teacherUsername, TeacherName: teacherName, CreateBy: createBy}
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
