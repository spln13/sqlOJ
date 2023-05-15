package model

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
	"time"
)

type StudentAccount struct {
	ID        int64  `gorm:"primary_key"`
	Number    string `gorm:"primary_key"`
	ClassID   int64
	ClassName string
	Username  string
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

func (*StudentAccountFlow) InsertStudentAccount(username, password, number, realName string) (int64, error) {
	studentAccountDAO := &StudentAccount{
		Number:   number,
		Username: username,
		Password: password,
		RealName: realName,
	}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(studentAccountDAO).Error
	}); err != nil {
		log.Println(err)
		return 0, errors.New("新增用户信息错误")
	}
	return studentAccountDAO.ID, nil
}

func (*StudentAccountFlow) CreateStudentAccount(number, realName, className, password string, classID int64) (int64, error) {
	studentAccountDAO := &StudentAccount{
		Number:    number,
		Username:  number,
		Password:  password,
		RealName:  realName,
		ClassName: className,
		ClassID:   classID,
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

func (*StudentAccountFlow) QueryStudentPasswordByEmail(email string) (int64, string, string, error) {
	var studentAccountDAO StudentAccount
	if err := GetSysDB().Select("id", "password", "username").Where("email = ?", email).Find(&studentAccountDAO).Error; err != nil {
		log.Println(err)
		return 0, "", "", errors.New("查询用户密码错误")
	}
	return studentAccountDAO.ID, studentAccountDAO.Password, studentAccountDAO.Username, nil
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

// UpdateStudentsClass 传入一个学生id列表，将学生的班级更新为 className
func (*StudentAccountFlow) UpdateStudentsClass(classID int64, studentIDList []int64) error {
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(&StudentAccount{}).Where("id in ?", studentIDList).Update("class_id", classID).Error
	}); err != nil {
		log.Println(err)
		return errors.New("更新学生班级属性错误")
	}
	return nil
}

func (*StudentAccountFlow) QueryStudentUsernameByUserID(userID int64) string {
	var studentAccountDAO StudentAccount
	if err := GetSysDB().Select("username").Where("id = ?", userID).Find(&studentAccountDAO).Error; err != nil {
		log.Println(err)
	}
	return studentAccountDAO.Username
}

type StudentClassAPI struct {
	ID int64
}

// QueryStudentIDByClassID 通过ClassIDList查询所有classID在其中的学生ID
func (*StudentAccountFlow) QueryStudentIDByClassID(classIDList []int64) ([]int64, error) {
	var studentClassAPIList []StudentClassAPI
	err := GetSysDB().Model(&StudentAccount{}).Where("class_id in ?", classIDList).Find(&studentClassAPIList).Error
	if err != nil {
		return nil, errors.New("查询学生ID错误")
	}
	var studentIDList []int64
	for _, studentClassAPI := range studentClassAPIList {
		studentIDList = append(studentIDList, studentClassAPI.ID)
	}
	return studentIDList, nil
}

func (*StudentAccountFlow) QueryAllStudent() ([]StudentAccount, error) {
	var studentAccountDAOList []StudentAccount
	err := GetSysDB().Model(&StudentAccount{}).Omit("password", "created_at", "updated_at").Order("number").Find(&studentAccountDAOList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("查询学生信息错误")
	}
	return studentAccountDAOList, nil
}

type StudentIDNumberStruct struct {
	ID     int64
	Number int64
}

func (*StudentAccountFlow) QueryStudentIDNumberMap(studentIDList []int64) (map[int64]int64, error) {
	var StudentIDNumberStructList []StudentIDNumberStruct
	err := GetSysDB().Model(&StudentAccount{}).Where("id in ?", studentIDList).Find(&StudentIDNumberStructList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("查询学生学号错误")
	}
	studentIDNumberMap := make(map[int64]int64)
	for _, studentIDNumberStruct := range StudentIDNumberStructList {
		studentIDNumberMap[studentIDNumberStruct.ID] = studentIDNumberStruct.Number
	}
	return studentIDNumberMap, nil
}

func (*StudentAccountFlow) ResetStudentPassword(studentID int64) error {
	var studentAccountDAO StudentAccount
	err := GetSysDB().Model(&StudentAccount{}).Select("number").Where("id = ?", studentID).Find(&studentAccountDAO).Error
	if err != nil {
		log.Println(err)
		return errors.New("获取学生信息错误")
	}
	digest := sha256.New() // 对密码加密
	digest.Write([]byte(studentAccountDAO.Number))
	passwordSHA := hex.EncodeToString(digest.Sum(nil))
	err = GetSysDB().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&StudentAccount{}).Where("id = ?", studentID).Update("password", passwordSHA).Error
		return err
	})
	if err != nil {
		log.Println(err)
		return errors.New("重置密码错误")
	}
	return nil
}

func (*StudentAccountFlow) QueryStudentRealNameByNumber(number int64) (string, error) {
	numberStr := strconv.FormatInt(number, 10)
	var studentAccount StudentAccount
	err := GetSysDB().Model(&StudentAccount{}).Select("real_name").Where("number = ?", numberStr).Find(&studentAccount).Error
	if err != nil {
		log.Println(err)
		return "", errors.New("查询学生真实姓名错误")
	}
	return studentAccount.RealName, nil
}

func (*StudentAccountFlow) QueryStudentNumberByID(studentID int64) string {
	var studentAccount StudentAccount
	err := GetSysDB().Model(&StudentAccount{}).Select("number").Where("id = ?", studentID).Find(&studentAccount).Error
	if err != nil {
		log.Println(err)
		return ""
	}
	return studentAccount.Number
}
