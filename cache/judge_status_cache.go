package cache

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// SetJudgeStatusPending 将当前InQueue状态写入cache
func SetJudgeStatusPending(userID, userType, exerciseID int64, answer string, submitTime time.Time) error {
	statusHashMap := map[string]interface{}{
		"status":              1,
		"answer":              answer,
		"submitTime":          submitTime.Unix(),
		"validNextSubmitTime": submitTime.Add(5 * time.Second).Unix(),
	}

	key := fmt.Sprintf("%d:%d:%d", userID, userType, exerciseID)
	err := rdb.HMSet(ctx, key, statusHashMap).Err()
	if err != nil {
		log.Println(err)
		return errors.New("判题状态写入缓存错误")
	}
	return nil
}

// SetJudgeStatusJudging 将cache对应Hash表项中的status设置为2(judging)
func SetJudgeStatusJudging(userID, userType, exerciseID int64, wg *sync.WaitGroup) {
	key := fmt.Sprintf("%d:%d:%d", userID, userType, exerciseID)
	err := rdb.HSet(ctx, key, "status", 2).Err()
	if err != nil {
		log.Println(err)
	}
	wg.Done()
}

// CheckSubmitTimeValid 检查同一个用户在同一个题目的提交间隔是否合法
func CheckSubmitTimeValid(userID, userType, exerciseID int64) (bool, error) {
	key := fmt.Sprintf("%d:%d:%d", userID, userType, exerciseID)
	exists, err := rdb.HExists(ctx, key, "submitTime").Result()
	if err != nil {
		log.Println(err)
		return false, errors.New("查询缓存key错误")
	}
	if !exists {
		return true, nil
	}
	validNextSubmitTime, err := rdb.HGet(ctx, key, "validNextSubmitTime").Int64()
	if err != nil {
		log.Println(err)
		return false, errors.New("检查上次提交时间错误")
	}
	if validNextSubmitTime > time.Now().Unix() {
		return false, nil
	}
	return true, nil

}

// DeleteSubmitStatus 删除cache中对应的Status记录
func DeleteSubmitStatus(userID, userType, exerciseID int64) {
	key := fmt.Sprintf("%d:%d:%d", userID, userType, exerciseID)
	err := rdb.HDel(ctx, key, "status", "answer", "submitTime", "validNextSubmitTime").Err()
	if err != nil {
		log.Println(err)
	}
}
