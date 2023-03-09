package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// ExpectedResultCache 缓存exerciseID对应的查询结果，减少数据库查询操作。
func ExpectedResultCache(exerciseID int64, result []map[string]interface{}) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	cacheKey := fmt.Sprintf("ExpectedResultCache:%d", exerciseID)
	err = rdb.Set(ctx, cacheKey, resultJSON, 5*time.Minute).Err()
	if err != nil {
		log.Println(err)
	}
}

// GetExpectedResult 从cache中读取exerciseID的结果返回. 成功bool类型返回true, 若不存在则bool类型返回值返回false，
func GetExpectedResult(exerciseID int64) ([]map[string]interface{}, bool) {
	cacheKey := fmt.Sprintf("ExpectedResultCache:%d", exerciseID)
	exists, err := rdb.Exists(ctx, cacheKey).Result()
	if err != nil {
		log.Println(err)
		return nil, false // 空结果，不存在
	}
	if exists == 0 { // 该key不存在
		return nil, false
	}
	resultBytes, err := rdb.Get(ctx, cacheKey).Bytes() // 从redis中获取数据
	if err != nil {
		log.Println(err)
	}
	// 将resultBytes还原成[]map[string]interface{}
	var retrievedData []map[string]interface{}
	err = json.Unmarshal(resultBytes, &retrievedData)
	if err != nil {
		log.Println(err)
	}
	return retrievedData, true
}
