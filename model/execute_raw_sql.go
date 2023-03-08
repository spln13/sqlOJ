package model

import "log"

// ExecuteRawSql 执行在练习数据库执行sql语句，结果保存在[]map[string]interface{}中返回
func ExecuteRawSql(sql string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	if err := GetExeDB().Raw(sql).Scan(&result).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}
