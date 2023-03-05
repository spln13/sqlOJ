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
	Name          string
	Grade         int
	Answer        string
	Description   string
	SubmitCount   int
	PassCount     int
	Visitable     int
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

func (*ExerciseContentFlow) InsertExerciseContent(publisherID int64, publisherType int64, name string, answer string, description string, grade int, visitable int) (int64, error) {
	exerciseContentDAO := &ExerciseContent{
		PublisherID:   publisherID,
		PublisherType: publisherType,
		Name:          name,
		Grade:         grade,
		Answer:        answer,
		Description:   description,
		Visitable:     visitable,
	}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(exerciseContentDAO).Error
	}); err != nil {
		log.Println(err)
		return 0, errors.New("保存记录错误")
	}
	return exerciseContentDAO.ID, nil
}
