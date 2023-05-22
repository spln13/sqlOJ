package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type ExerciseContent struct {
	ID            int64 `gorm:"primary_key"`
	PublisherID   int64
	PublisherType int64
	PublisherName string
	Name          string
	Grade         int
	Answer        string
	Description   string
	SubmitCount   int
	PassCount     int
	Type          int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ExerciseContentFlow struct {
}

var (
	exerciseContentFlow *ExerciseContentFlow
	exerciseContentOnce sync.Once
)

func NewExerciseContentFlow() *ExerciseContentFlow {
	exerciseContentOnce.Do(func() {
		exerciseContentFlow = new(ExerciseContentFlow)
	})
	return exerciseContentFlow
}

func (*ExerciseContentFlow) InsertExerciseContent(publisherID int64, publisherType int64, name, answer, publisherName, description string, exeType, grade int) (int64, error) {
	exerciseContentDAO := &ExerciseContent{
		PublisherID:   publisherID,
		PublisherType: publisherType,
		PublisherName: publisherName,
		Name:          name,
		Grade:         grade,
		Answer:        answer,
		Description:   description,
		Type:          exeType,
	}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(exerciseContentDAO).Error
	}); err != nil {
		log.Println(err)
		return 0, errors.New("保存记录错误")
	}
	return exerciseContentDAO.ID, nil
}

// QueryAnswerTypeByExerciseID 根据题目ID查询答案
func (*ExerciseContentFlow) QueryAnswerTypeByExerciseID(exerciseID int64) (string, int) {
	var exerciseContentDAO ExerciseContent
	if err := GetSysDB().Select("id", "type", "answer").Where("id = ?", exerciseID).Find(&exerciseContentDAO); err != nil {
		log.Println(err)
	}
	if exerciseContentDAO.ID == 0 {
		log.Println("error: 未找到对应题目答案")
	}
	return exerciseContentDAO.Answer, exerciseContentDAO.Type
}

func (*ExerciseContentFlow) QueryExerciseNameByExerciseID(exerciseID int64) string {
	var exerciseContentDAO ExerciseContent
	if err := GetSysDB().Select("name").Where("id = ?", exerciseID).Find(&exerciseContentDAO).Error; err != nil {
		log.Println(err)
	}
	return exerciseContentDAO.Name
}

// GetAllExercise 获取当前数据库中所有可见的题目
func (*ExerciseContentFlow) GetAllExercise() ([]ExerciseContent, error) {
	var exerciseContentArray []ExerciseContent
	err := GetSysDB().Select("id, publisher_id", "publisher_type", "publisher_name", "name", "grade", "submit_count", "pass_count", "type").Find(&exerciseContentArray).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("获取题库数据错误")
	}
	return exerciseContentArray, nil
}

func (*ExerciseContentFlow) GetOneExercise(exerciseID int64) (ExerciseContent, error) {
	var exerciseContent ExerciseContent
	if err := GetSysDB().Select("id", "publisher_id", "publisher_type", "publisher_name", "name", "grade", "description", "submit_count", "pass_count", "created_at").Where("id = ?", exerciseID).Find(&exerciseContent).Error; err != nil {
		log.Println(err)
		return exerciseContent, errors.New("查询题目信息错误")
	}
	return exerciseContent, nil
}

// IncreasePassCountSubmitCount 答案正确时调用，将提交总数和通过总数自增
func (*ExerciseContentFlow) IncreasePassCountSubmitCount(exerciseID int64) {
	err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&ExerciseContent{}).Where("id = ?", exerciseID).Updates(map[string]interface{}{
			"submit_count": gorm.Expr("submit_count + ?", 1),
			"pass_count":   gorm.Expr("pass_count + ?", 1),
		}).Error
		return err
	})
	if err != nil {
		log.Println(err)
	}
}

// IncreaseSubmitCount 提交答案未通过时调用，将提交总数自增
func (*ExerciseContentFlow) IncreaseSubmitCount(exerciseID int64) {
	err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&ExerciseContent{}).Where("id = ?", exerciseID).Update("submit_count", gorm.Expr("submit_count + 1")).Error
		return err
	})
	if err != nil {
		log.Println(err)
	}
}

// QueryExerciseGrade 查询对应题目的难度等级
func (*ExerciseContentFlow) QueryExerciseGrade(exerciseID int64) int {
	var exerciseContentDAO ExerciseContent
	err := GetSysDB().Select("grade").Where("id = ?", exerciseID).Find(&exerciseContentDAO).Error
	if err != nil {
		log.Println(err)
	}
	return exerciseContentDAO.Grade
}

func (*ExerciseContentFlow) DeleteExerciseContent(exerciseID int64) error {
	err := GetSysDB().Delete(&ExerciseContent{}, exerciseID).Error
	if err != nil {
		log.Println(err)
		return errors.New("删除题目错误")
	}
	return nil
}
