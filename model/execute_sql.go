package model

import (
	"fmt"
	"log"
	"reflect"
)

// ExecuteRawSql 执行在练习数据库执行sql语句，结果保存在[]map[string]interface{}中返回
func ExecuteRawSql(sql string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	if err := GetExeDB().Raw(sql).Scan(&result).Error; err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}

// CompareModifySqlResult 执行更改类型的sql语句
// step1. 复制originTableName到tempTableName
// step2. 开启事务执行userSql语句将全表查询结果保存到userResult
// step3. 回滚
// step4. 开启事务执行expectedSql语句将全表查询结果保存到expectedResult
// step5. 回滚, 并删除临时表
// step6. 比较userResult与expectedResult返回结果
func CompareModifySqlResult(userSql string, expectedSql string, originTableName string, tempTableName string) int {
	// 创建临时表
	createTempTableSql := fmt.Sprintf("create table %s like %s", tempTableName, originTableName)
	insertDataSql := fmt.Sprintf(" insert into %s select * from %s", tempTableName, originTableName)
	queryTableSql := fmt.Sprintf("select * from %s", tempTableName)
	var userResult []map[string]interface{}
	var expectedResult []map[string]interface{}
	var statusCode int                                                // returnCode返回代码 1->AC 2->WA 3->RE
	if err := GetExeDB().Exec(createTempTableSql).Error; err != nil { // 复制表, 表名为tempTableName
		log.Println(err)
		statusCode = 3
	}
	if err := GetExeDB().Exec(insertDataSql).Error; err != nil { // 复制表, 表名为tempTableName
		log.Println(err)
		statusCode = 3
	}
	defer func() { // 函数退出前删除副本表
		dropTableSql := fmt.Sprintf("drop table %s", tempTableName)
		if err := GetExeDB().Exec(dropTableSql).Error; err != nil { // 删除副本表
			log.Println(err)
			statusCode = 3
		}
	}()
	tx1 := GetExeDB().Begin()                       // 开启手动事务
	if err := tx1.Exec(userSql).Error; err != nil { // 执行用户的sql语句
		log.Println(err)
		statusCode = 3
		return statusCode
	}
	if err := tx1.Raw(queryTableSql).Scan(&userResult).Error; err != nil { // 执行全表查询
		log.Println(err)
		statusCode = 3
		return statusCode
	}
	tx1.Rollback()                                      // 回滚
	tx2 := GetExeDB().Begin()                           // 开启手动事务
	if err := tx2.Exec(expectedSql).Error; err != nil { // 执行标准答案的sql语句
		log.Println(err)
		statusCode = 3
		return statusCode
	}
	if err := tx2.Raw(queryTableSql).Scan(&expectedResult).Error; err != nil { // 执行全表查询
		log.Println(err)
		statusCode = 3
		return statusCode
	}
	tx2.Rollback() // 回滚
	if reflect.DeepEqual(userResult, expectedResult) {
		statusCode = 1
	} else {
		statusCode = 2
	}
	return statusCode
}
