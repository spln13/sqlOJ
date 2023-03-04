package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type Class struct {
	ID           uint `gorm:"primary_key"`
	Name         string
	TeacherID    uint
	TeacherName  string
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

func (*ClassFlow) InsertClass(name string, teacherID uint, teacherName string) error {
	classDAO := &Class{Name: name, TeacherID: teacherID, TeacherName: teacherName}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(classDAO).Error
	}); err != nil {
		log.Println(err)
		return errors.New("保存班级信息错误")
	}
	return nil
}
