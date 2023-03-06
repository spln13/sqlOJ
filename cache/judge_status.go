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
		"status":     1,
		"answer":     answer,
		"submitTime": submitTime.Unix(),
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
