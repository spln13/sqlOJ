package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type TeacherAccount struct {
	ID        int64  `gorm:"primary_key"`
	Username  string // 此为教师职工号
	RealName  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TeacherAccountFlow struct {
}

var (
	teacherAccountFlow *TeacherAccountFlow
	teacherAccountOnce sync.Once
)

func NewTeacherAccountFlow() *TeacherAccountFlow {
	teacherAccountOnce.Do(func() {
		teacherAccountFlow = new(TeacherAccountFlow)
	})
	return teacherAccountFlow
}

func (*TeacherAccountFlow) QueryTeacherExistByUsername(username string) (bool, error) {
	var teacherAccountDAO TeacherAccount
	if err := GetSysDB().Select("id").Where("username = ?", username).Find(&teacherAccountDAO).Error; err != nil {
		log.Println(err)
		return false, errors.New("查询用户信息错误")
	}
	if teacherAccountDAO.ID != 0 { // 没有根据用户名找到
		return true, nil
	}
	return false, nil
}

func (*TeacherAccountFlow) InsertTeacherAccount(username, password, realName string) error {
	teacherAccountDAO := &TeacherAccount{Username: username, Password: password, RealName: realName}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(teacherAccountDAO).Error
	}); err != nil {
		log.Println(err.Error())
		return errors.New("保存教师信息错误")
	}
	return nil
}

// QueryTeacherPasswordByUsername 通过用户名查询教师密码
func (*TeacherAccountFlow) QueryTeacherPasswordByUsername(username string) (int64, string, error) {
	var teacherAccountDAO TeacherAccount
	if err := GetSysDB().Select("id", "password").Where("username = ?", username).Find(&teacherAccountDAO).Error; err != nil {
		log.Println(err)
		return 0, "", errors.New("查询密码错误")
	}
	if teacherAccountDAO.ID == 0 {
		return 0, "", errors.New("用户不存在")
	}
	return teacherAccountDAO.ID, teacherAccountDAO.Password, nil

}

func (*TeacherAccountFlow) QueryTeacherPasswordByUserID(userID int64) (string, error) {
	var teacherAccountDAO TeacherAccount
	if err := GetSysDB().Select("password").Where("id = ?", userID).Find(&teacherAccountDAO).Error; err != nil {
		log.Println(err)
		return "", errors.New("查询用户密码错误")
	}
	return teacherAccountDAO.Password, nil
}

func (*TeacherAccountFlow) UpdateTeacherPasswordByUserID(userID int64, password string) error {
	teacherAccountDAO := &TeacherAccount{ID: userID, Password: password}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(teacherAccountDAO).Update("password", password).Error
	}); err != nil {
		log.Println(err)
		return errors.New("更新密码错误")
	}
	return nil
}

func (*TeacherAccountFlow) QueryTeacherRealNameByUsername(username string) (string, error) {
	var teacherAccountDAO TeacherAccount
	if err := GetSysDB().Select("username").Where("username = ?", username).Find(&teacherAccountDAO).Error; err != nil {
		log.Println(err)
		return "", errors.New("查询教师姓名错误")
	}
	return teacherAccountDAO.RealName, nil
}
