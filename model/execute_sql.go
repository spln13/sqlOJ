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

// CompareModifySqlResultWithoutCache 在答案sql执行结果没有缓存情况下执行更改类型的sql语句
// step1. 复制originTableName到UserTempTableName
// step2. 执行userSql语句将全表查询结果保存到userResult
// step3. 复制originTableName到AnswerTempTableName
// step4. 执行expectedSql语句将全表查询结果保存到expectedResult
// step5. 删除两个临时表
// step6. 比较userResult与expectedResult返回结果, 将expectedResult一并返回
func CompareModifySqlResultWithoutCache(userSql string, expectedSql string, originTableName string, userTempTableName string, answerTempTableName string) ([]map[string]interface{}, int) {
	// 创建临时表
	createUserTempTableSql := fmt.Sprintf("create table %s like %s", userTempTableName, originTableName)
	createAnswerTempTableSql := fmt.Sprintf("create table %s like %s", answerTempTableName, originTableName)
	dropTableSql := fmt.Sprintf("drop tables %s, %s", userTempTableName, answerTempTableName)
	var userResult []map[string]interface{}
	var expectedResult []map[string]interface{}
	var statusCode int
	defer func() { // 函数退出前删除两张临时表
		if err := GetExeDB().Exec(dropTableSql).Error; err != nil {
			log.Println(err)
		}
	}()
	// returnCode返回代码 1->AC 2->WA 3->RE
	if err := GetExeDB().Exec(createUserTempTableSql).Error; err != nil { // 复制表, 表名为userTempTableName
		log.Println(err)
	}
	if err := GetExeDB().Exec(createAnswerTempTableSql).Error; err != nil { // 复制表, 表名为answerTempTableName
		log.Println(err)
	}
	if err := GetExeDB().Exec(userSql).Error; err != nil { // 执行用户更改表的sql语句
		log.Println(err)
		statusCode = 3
		return nil, statusCode
	}
	queryUserTableSql := fmt.Sprintf("select * from %s", userTempTableName)
	if err := GetExeDB().Raw(queryUserTableSql).Scan(&userResult).Error; err != nil { // 执行全表查询，结果保存到userResult
		log.Println(err)
		statusCode = 3
		return nil, statusCode
	}
	if err := GetExeDB().Exec(expectedSql).Error; err != nil { // 执行答案更改表的sql语句
		log.Println(err)
		statusCode = 3
		return nil, statusCode
	}
	queryAnswerTableSql := fmt.Sprintf("select * from %s", answerTempTableName)
	if err := GetExeDB().Raw(queryAnswerTableSql).Scan(&expectedResult).Error; err != nil { // 执行全表查询，结果保存到expectedResult
		log.Println(err)
		statusCode = 3
		return nil, statusCode
	}
	if reflect.DeepEqual(userResult, expectedResult) {
		// 结果相等
		statusCode = 1
		return expectedResult, statusCode
	}
	// 结果不等
	statusCode = 2
	return expectedResult, statusCode
}

// CompareModifySqlResultWithCache 在有缓存情况下比对用户sql答案与标准答案
// step1. 复制originTableName到userTempTableName
// step2. 执行userSql语句将全表查询结果保存到userResult
// step3. 删除临时表
// step4. 比较userResult与expectedResult返回结果
func CompareModifySqlResultWithCache(userSql string, originTableName string, userTempTableName string, expectedResult []map[string]interface{}) int {
	createUserTempTableSql := fmt.Sprintf("create table %s like %s", userTempTableName, originTableName) // 定义创建用户临时表语句
	var statusCode int
	if err := GetExeDB().Exec(createUserTempTableSql).Error; err != nil { // 创建表错误
		log.Println(err)
		statusCode = 3
		return statusCode
	}
	dropTableSql := fmt.Sprintf("drop table %s", userTempTableName) // 删除表语句
	defer func() {                                                  // 函数退出前删除创建的临时表
		if err := GetExeDB().Exec(dropTableSql).Error; err != nil {
			log.Println(err)
		}
	}()
	if err := GetExeDB().Exec(userSql).Error; err != nil { // 执行用户更改表的sql语句
		log.Println(err)
		statusCode = 3
		return statusCode
	}
	var userResult []map[string]interface{}
	queryUserTableSql := fmt.Sprintf("select * from %s", userTempTableName)           // 查询全表语句
	if err := GetExeDB().Raw(queryUserTableSql).Scan(&userResult).Error; err != nil { // 执行全表查询，结果保存到userResult
		log.Println(err)
		statusCode = 3
		return statusCode
	}
	if reflect.DeepEqual(userResult, expectedResult) {
		// 结果相等
		return 1
	}
	// 结果不等
	return 2
}
