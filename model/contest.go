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
	BeginAt       time.Time
	EndAt         time.Time
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
		BeginAt:       beginAt,
		EndAt:         endAt,
	}
	if err := GetSysDB().Create(contestDAO).Error; err != nil {
		log.Println(err)
		return 0, errors.New("创建竞赛错误")
	}
	return contestDAO.ID, nil
}

// GetAllContest 获取所有竞赛，按照开始时间降序
func (*ContestFlow) GetAllContest() ([]Contest, error) {
	var contestList []Contest
	err := GetSysDB().Model(&Contest{}).Select("id", "name", "publisher_name", "publisher_type", "begin_at", "end_at").Order("begin_at desc").Find(&contestList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("获取竞赛错误")
	}
	return contestList, nil
}

func (*ContestFlow) GetContestNameByID(contestID int64) string {
	var contestDAO Contest
	err := GetSysDB().Model(&Contest{}).Select("name").Where("id = ?", contestID).Find(&contestDAO).Error
	if err != nil {
		log.Println(err)
		return ""
	}
	return contestDAO.Name
}
