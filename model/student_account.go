package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type StudentAccount struct {
	ID        int64  `gorm:"primary_key"`
	Number    string `gorm:"primary_key"`
	ClassID   int64
	Username  string
	Email     string
	Password  string
	RealName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type StudentAccountFlow struct {
}

var (
	studentAccountFlow *StudentAccountFlow
	studentAccountOnce sync.Once
)

func NewStudentAccountFlow() *StudentAccountFlow {
	studentAccountOnce.Do(func() {
		studentAccountFlow = new(StudentAccountFlow)
	})
	return studentAccountFlow
}

func (*StudentAccountFlow) InsertStudentAccount(username, password, number, realName, email string) (int64, error) {
	studentAccountDAO := &StudentAccount{
		Number:   number,
		Username: username,
		Password: password,
		RealName: realName,
		Email:    email,
	}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(studentAccountDAO).Error
	}); err != nil {
		log.Println(err)
		return 0, errors.New("新增用户信息错误")
	}
	return studentAccountDAO.ID, nil
}

func (*StudentAccountFlow) QueryStudentPasswordByUsername(username string) (int64, string, error) {
	var studentAccountDAO StudentAccount
	if err := GetSysDB().Select("id", "password").Where("username = ?", username).Find(&studentAccountDAO).Error; err != nil {
		log.Println(err)
		return 0, "", errors.New("查询用户密码错误")
	}
	return studentAccountDAO.ID, studentAccountDAO.Password, nil
}

func (*StudentAccountFlow) QueryStudentPasswordByEmail(email string) (int64, string, error) {
	var studentAccountDAO StudentAccount
	if err := GetSysDB().Select("id", "password").Where("email = ?", email).Find(&studentAccountDAO).Error; err != nil {
		log.Println(err)
		return 0, "", errors.New("查询用户密码错误")
	}
	return studentAccountDAO.ID, studentAccountDAO.Password, nil
}

func (*StudentAccountFlow) QueryStudentPasswordByUserID(userID int64) (string, error) {
	var studentAccountDAO StudentAccount
	if err := GetSysDB().Select("password").Where("id = ?", userID).Find(&studentAccountDAO).Error; err != nil {
		log.Println(err)
		return "", errors.New("查询用户密码错误")
	}
	return studentAccountDAO.Password, nil
}

func (*StudentAccountFlow) UpdateStudentPasswordByUserID(userID int64, password string) error {
	studentAccountDAO := &StudentAccount{ID: userID, Password: password}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(&studentAccountDAO).Update("password", password).Error
	}); err != nil {
		log.Println(err)
		return errors.New("更新用户密码错误")
	}
	return nil
}

func (*StudentAccountFlow) QueryStudentExistByUsername(username string) (bool, error) {
	var studentAccountDAO StudentAccount
	if err := GetSysDB().Select("id").Where("username = ?", username).Find(&studentAccountDAO).Error; err != nil {
		log.Println(err)
		return false, errors.New("查询用户名存在信息错误")
	}
	if studentAccountDAO.ID != 0 {
		return true, nil
	}
	return false, nil
}

func (*StudentAccountFlow) QueryStudentExistByEmail(email string) (bool, error) {
	var studentAccountDAO StudentAccount
	if err := GetSysDB().Select("id").Where("email = ?", email).Find(&studentAccountDAO).Error; err != nil {
		log.Println(err)
		return false, errors.New("查询邮箱存在信息错误")
	}
	if studentAccountDAO.ID != 0 {
		return true, nil
	}
	return false, nil
}
