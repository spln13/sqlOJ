package model

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type ScoreRecord struct {
	ID        int64 `gorm:"primary_key"`
	UserID    int64
	UserType  int64
	Score     int64
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RankingAPI struct {
	Username string
	UserType int64
	Score    int64
}

type ScoreRecordFlow struct {
}

var (
	scoreTableFlow *ScoreRecordFlow
	scoreTableOnce sync.Once
)

func NewScoreRecordFlow() *ScoreRecordFlow {
	scoreTableOnce.Do(func() {
		scoreTableFlow = new(ScoreRecordFlow)
	})
	return scoreTableFlow
}

func (*ScoreRecordFlow) InsertScoreRecord(userID, userType int64, username string) error {
	scoreRecordDAO := ScoreRecord{
		Username: username,
		UserID:   userID,
		UserType: userType,
		Score:    0,
	}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Create(&scoreRecordDAO).Error
	}); err != nil {
		log.Println(err)
		return errors.New("创建用户得分记录错误")
	}
	return nil
}

func (*ScoreRecordFlow) IncreaseScore(userID, userType int64, grade int) {
	var score int64
	if grade == 1 { // easy
		score = 3
	} else if grade == 2 { // medium
		score = 5
	} else if grade == 3 { // hard
		score = 8
	} else {
		score = 0
	}
	if err := GetSysDB().Transaction(func(tx *gorm.DB) error {
		return tx.Model(&ScoreRecord{}).Where("user_id = ? and user_type = ?", userID, userType).Update("score", gorm.Expr("score + ?", score)).Error
	}); err != nil {
		log.Println(err)
	}
}

func (*ScoreRecordFlow) GetRanking() ([]RankingAPI, error) {
	var rankingAPIList []RankingAPI
	err := GetSysDB().Model(&ScoreRecord{}).Order("score desc").Find(&rankingAPIList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("获取排名错误")
	}
	return rankingAPIList, nil
}

// GetMinRanking 获取排行榜前五条记录
func (*ScoreRecordFlow) GetMinRanking() ([]RankingAPI, error) {
	var rankingAPIList []RankingAPI
	err := GetSysDB().Model(&ScoreRecord{}).Order("score desc").Limit(5).Find(&rankingAPIList).Error
	if err != nil {
		log.Println(err)
		return nil, errors.New("获取排名错误")
	}
	return rankingAPIList, nil
}

func (*ScoreRecordFlow) DeleteRanking(studentID int64) error {
	err := GetSysDB().Where("user_type = 1 and user_id = ?", studentID).Delete(&ScoreRecord{}).Error
	if err != nil {
		log.Println(err)
		return errors.New("删除学生分数信息错误")
	}
	return nil
}
