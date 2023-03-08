# 判题更改型数据库查询语言

## Step1. 生成唯一临时表名
`userID_ExerciseID_submitTime`

## Step2. 创建临时表
使用gorm
```go
var tableName string = "source_table"
var newTableName string = "new_table"
result := db.Exec("CREATE TABLE " + newTableName + " LIKE " + tableName)
if result.Error != nil {
    panic(result.Error)
}
```

## Step3. 开启事务, 执行用户SQL
```go
func main() {
    db, err := gorm.Open("mysql", "user:password@tcp(host:port)/database")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // 开启事务
    tx := db.Begin()

    // 执行用户 SQL 语句
    result := tx.Exec("INSERT INTO users (name, age) VALUES (?, ?)", "user1", 20)
    if result.Error != nil {
        // 发生错误，回滚事务
        tx.Rollback()
        panic(result.Error)
    }

    // 回滚事务
    tx.Rollback()
}
```

## Step3. 开启事务，执行答案sql
> 优化方案: 将答案sql执行结果`[]map[string]interface{}`类型存于cache中设置5分钟过期时间，
> 执行答案sql前可查询cache，若未命中则执行sql
### 常规方案
同Step2
### 优化方案

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	// 连接 Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// 定义要存储的数据
	data := []map[string]interface{}{
		{
			"name": "Alice",
			"age":  20,
		},
		{
			"name": "Bob",
			"age":  30,
		},
	}

	// 将数据转换为 JSON 字符串
	dataJSON, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	// 存储数据并设置过期时间
	err = client.Set("mydata", dataJSON, 5*time.Minute).Err()
	if err != nil {
		panic(err)
	}

	// 从 Redis 中获取数据
	dataBytes, err := client.Get("mydata").Bytes()
	if err != nil {
		panic(err)
	}

	// 将 JSON 字符串解析为数据
	var retrievedData []map[string]interface{}
	err = json.Unmarshal(dataBytes, &retrievedData)
	if err != nil {
		panic(err)
	}

	// 输出数据
	fmt.Println(retrievedData)
}


```