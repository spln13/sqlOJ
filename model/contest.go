package model

import (
	"errors"
	"log"
	"sync"
	"time"
)

type Contest struct {
	ID            int64 `gorm:"primary_key"`
	Name          string
	PublisherID   int64
	PublisherType int64
	PublisherName string
	BeginTime     time.Time
	EndTime       time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ContestFlow struct {
}

var (
	contestFlowDAO  *ContestFlow
	contestFlowOnce sync.Once
)

func NewContestFlow() *ContestFlow {
	contestFlowOnce.Do(func() {
		contestFlowDAO = new(ContestFlow)
	})
	return contestFlowDAO
}

func (*ContestFlow) CreateContest(name, publisherName string, publisherID, publisherType int64, beginAt, endAt time.Time) (int64, error) {
	contestDAO := Contest{
		Name:          name,
		PublisherID:   publisherID,
		PublisherType: publisherType,
		PublisherName: publisherName,
		BeginTime:     beginAt,
		EndTime:       endAt,
	}
	if err := GetSysDB().Create(contestDAO).Error; err != nil {
		log.Println(err)
		return 0, errors.New("创建竞赛错误")
	}
	return contestDAO.ID, nil
}
