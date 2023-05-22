package model

import (
	"errors"
	"log"
	"sync"
	"time"
)

type ContestClassAssociation struct {
	ID        int64 `gorm:"primary_key"`
	ContestID int64
	ClassID   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ContestClassAssociationFlow struct {
}

var (
	contestClassAssociationFlowDAO  *ContestClassAssociationFlow
	contestClassAssociationFlowOnce sync.Once
)

func NewContestClassAssociationFlow() *ContestClassAssociationFlow {
	contestClassAssociationFlowOnce.Do(func() {
		contestClassAssociationFlowDAO = new(ContestClassAssociationFlow)
	})
	return contestClassAssociationFlowDAO
}

func (*ContestClassAssociationFlow) InsertContestClassAssociation(contestID int64, ClassIDList []int64) error {
	for _, classID := range ClassIDList {
		contestClassAssociationDAO := ContestClassAssociation{ContestID: contestID, ClassID: classID}
		if err := GetSysDB().Create(&contestClassAssociationDAO).Error; err != nil {
			log.Println(err)
			return errors.New("插入竞赛班级关联信息错误")
		}
	}
	return nil
}

func (*ContestClassAssociationFlow) DeleteContestClassAssociation(contestID int64) error {
	err := GetSysDB().Delete(&ContestClassAssociation{}).Where("contest_id = ?", contestID).Error
	if err != nil {
		log.Println(err)
		return errors.New("删除竞赛班级关系错误")
	}
	return nil
}
