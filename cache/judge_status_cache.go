package cache

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// SetExerciseJudgeStatusPending 将当前InQueue状态写入cache
func SetExerciseJudgeStatusPending(userID, userType, exerciseID int64, submitTime time.Time) error {
	statusHashMap := map[string]interface{}{
		"status": 4, // 在队列中状态码是4
	}
	// 设置下次合法提交时间为3秒后
	nextSubmitTimeKey := fmt.Sprintf("%d:%d:%d_next_submit_time", userID, userType, exerciseID)
	err := rdb.Set(ctx, nextSubmitTimeKey, submitTime.Add(5*time.Second).Unix(), time.Duration(3)*time.Second).Err()
	if err != nil {
		log.Println(err)
		return errors.New("缓存提交时间错误")
	}
	hashKey := fmt.Sprintf("%d:%d:%d:%d", userID, userType, exerciseID, submitTime.Unix())
	err = rdb.HSet(ctx, hashKey, statusHashMap).Err()
	if err != nil {
		log.Println(err)
		return errors.New("判题状态写入缓存错误")
	}
	return nil
}

// SetExerciseJudgeStatusJudging 将cache对应Hash表项中的status设置为5(judging)
func SetExerciseJudgeStatusJudging(userID, userType, exerciseID int64, submitTime time.Time, wg *sync.WaitGroup) {
	key := fmt.Sprintf("%d:%d:%d:%d", userID, userType, exerciseID, submitTime.Unix())
	err := rdb.HSet(ctx, key, "status", 5).Err() // judging中状态码是5
	if err != nil {
		log.Println(err)
	}
	wg.Done()
}

// SetContestJudgeStatusPending 将当前InQueue状态写入cache
func SetContestJudgeStatusPending(userID, userType, exerciseID, contestID int64, submitTime time.Time) error {
	statusHashMap := map[string]interface{}{
		"status": 4, // 在队列中状态码是4
	}
	// 设置下次合法提交时间为3秒后
	nextSubmitTimeKey := fmt.Sprintf("%d:%d:%d_next_submit_time", userID, userType, exerciseID)
	err := rdb.Set(ctx, nextSubmitTimeKey, submitTime.Add(5*time.Second).Unix(), time.Duration(3)*time.Second).Err()
	if err != nil {
		log.Println(err)
		return errors.New("缓存提交时间错误")
	}
	hashKey := fmt.Sprintf("%d:%d:%d:%d:%d", userID, userType, exerciseID, contestID, submitTime.Unix())
	err = rdb.HSet(ctx, hashKey, statusHashMap).Err()
	if err != nil {
		log.Println(err)
		return errors.New("判题状态写入缓存错误")
	}
	return nil
}

// SetContestJudgeStatusJudging 将cache对应Hash表项中的status设置为2(judging)
func SetContestJudgeStatusJudging(userID, userType, exerciseID, contestID int64, submitTime time.Time, wg *sync.WaitGroup) {
	key := fmt.Sprintf("%d:%d:%d:%d:%d", userID, userType, exerciseID, contestID, submitTime.Unix())
	err := rdb.HSet(ctx, key, "status", 5).Err() // judging中状态码是5
	if err != nil {
		log.Println(err)
	}
	wg.Done()
}

// CheckSubmitTimeValid 检查同一个用户在同一个题目的提交间隔是否合法
func CheckSubmitTimeValid(userID, userType, exerciseID int64) (bool, error) {
	key := fmt.Sprintf("%d:%d:%d_next_submit_time", userID, userType, exerciseID)
	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		log.Println(err)
		return false, errors.New("查询缓存key错误")
	}
	if exists == 0 {
		return true, nil
	}
	validNextSubmitTime, err := rdb.Get(ctx, key).Int64()
	if err != nil {
		log.Println(err)
		return false, errors.New("检查上次提交时间错误")
	}
	if validNextSubmitTime > time.Now().Unix() {
		return false, nil
	}
	return true, nil
}

// DeleteExerciseSubmitStatus 删除cache中对应的Status记录
func DeleteExerciseSubmitStatus(userID, userType, exerciseID int64, submitTime time.Time) {
	key := fmt.Sprintf("%d:%d:%d:%d", userID, userType, exerciseID, submitTime.Unix())
	err := rdb.HDel(ctx, key, "status").Err()
	if err != nil {
		log.Println(err)
	}
}

// DeleteContestSubmitStatus 删除cache中contest对应的Status记录
func DeleteContestSubmitStatus(userID, userType, exerciseID, contestID int64, submitTime time.Time) {
	key := fmt.Sprintf("%d:%d:%d:%d:%d", userID, userType, exerciseID, contestID, submitTime.Unix())
	err := rdb.HDel(ctx, key, "status").Err()
	if err != nil {
		log.Println(err)
	}
}

// GetUserJudgeStatus 获得指定userID与userType与exerciseID的redis结果
func GetUserJudgeStatus(userID int64, userType int64, exerciseID int64) (map[string]string, error) {
	pattern := fmt.Sprintf("%d:%d:%d:*", userID, userType, exerciseID)
	keys, err := rdb.Keys(ctx, pattern).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	kvMap := make(map[string]string)
	for _, key := range keys {
		value, err := rdb.HGet(ctx, key, "status").Result()
		if err != nil {
			log.Println(err)
			return nil, errors.New("查询提交状态缓存错误")
		}
		kvMap[key] = value
	}
	return kvMap, nil
}

// GetContestUserJudgeStatus 获得指定userID与userType与exerciseID的redis结果
func GetContestUserJudgeStatus(userID, userType, exerciseID, contestID int64) (map[string]string, error) {
	pattern := fmt.Sprintf("%d:%d:%d:%d:*", userID, userType, exerciseID, contestID)
	keys, err := rdb.Keys(ctx, pattern).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	kvMap := make(map[string]string)
	for _, key := range keys {
		value, err := rdb.HGet(ctx, key, "status").Result()
		if err != nil {
			log.Println(err)
			return nil, errors.New("查询提交状态缓存错误")
		}
		kvMap[key] = value
	}
	return kvMap, nil
}

// GetAllKeyValue 获取缓存中所有的key-value返回map[string]string
func GetAllKeyValue() (map[string]string, error) {
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	kvMap := make(map[string]string)
	for _, key := range keys {
		value, err := rdb.HGet(ctx, key, "status").Result()
		if err != nil {
			log.Println(err)
			return nil, errors.New("查询提交状态缓存错误")
		}
		kvMap[key] = value
	}
	return kvMap, nil
}
