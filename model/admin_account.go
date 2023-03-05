package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type AdminAccount struct {
	ID        int64 `gorm:"primary_key"`
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AdminAccountFlow struct {
}

var (
	adminAccountFlow *AdminAccountFlow
	adminAccountOnce sync.Once
)

func NewAdminAccountFlow() *AdminAccountFlow {
	adminAccountOnce.Do(func() {
		adminAccountFlow = new(AdminAccountFlow)
	})
	return adminAccountFlow
}

// QueryAdminPasswordByUsername 根据管理员用户名查询管理员密码以及ID, ID用来颁发token
func (*AdminAccountFlow) QueryAdminPasswordByUsername(username string) (int64, string, error) {
	var adminDAO AdminAccount
	if err := GetSysDB().Select("id", "password").Where("username = ?", username).Find(&adminDAO).Error; err != nil {
		log.Println(err)
		return 0, "", errors.New("查询用户错误")
	}
	return adminDAO.ID, adminDAO.Password, nil
}

// QueryAdminPasswordByUserID  根据管理员用户id查询管理员密码
func (*AdminAccountFlow) QueryAdminPasswordByUserID(userID int64) (string, error) {
	var adminDAO AdminAccount
	if err := GetSysDB().Select("password").Where("id = ?", userID).Find(&adminDAO).Error; err != nil {
		log.Println(err)
		return "", errors.New("查询用户错误")
	}
	return adminDAO.Password, nil
}

// QueryAdminExistByUsername 根据管理员用户名查询该用户名是否存在
func (*AdminAccountFlow) QueryAdminExistByUsername(username string) (bool, error) {
	var adminDAO AdminAccount
	if err := GetSysDB().Select("id", "password").Where("username = ?", username).Find(&adminDAO).Error; err != nil {
		log.Println(err)
		return false, errors.New("检查用户名错误")
	}
	if adminDAO.ID != 0 {
		return true, nil
	}
	return false, nil
}

func (*AdminAccountFlow) InsertAdminAccount(username, password string) error {
	adminDAO := &AdminAccount{
		Username: username,
		Password: password,
	}
	err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(adminDAO).Error
	})
	if err != nil {
		log.Println(err)
		return errors.New("插入用户错误")
	}
	return nil
}

func (*AdminAccountFlow) UpdateAdminPassword(userID int64, newPassword string) error {
	adminDAO := &AdminAccount{
		ID:       userID,
		Password: newPassword,
	}
	err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(adminDAO).Update("password", newPassword).Error
	})
	if err != nil {
		log.Println(err)
		return errors.New("更新密码错误")
	}
	return nil
}
