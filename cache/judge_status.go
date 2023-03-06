package cache

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// SetJudgeStatusInQueue 将当前InQueue状态写入cache
func SetJudgeStatusInQueue(userID, userType, exerciseID int64, answer string, submitTime time.Time) error {
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
