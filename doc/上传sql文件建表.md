## 上传sql文件保存到本地
```go
func saveSQLFile(c *gin.Context) {
    file, err := c.FormFile("sql_file")
    if err != nil {
        c.AbortWithError(http.StatusBadRequest, err)
        return
    }

    err = c.SaveUploadedFile(file, "path/to/save/sql_file.sql")
    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "File saved successfully",
    })
}

```

## 运行该文件中的SQL语句来创建表
```go
func createTables(c *gin.Context) {
    db := c.MustGet("db").(*gorm.DB) // 获取 GORM 数据库连接

    // 读取 SQL 文件内容
    sqlFile, err := ioutil.ReadFile("path/to/save/sql_file.sql")
    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    // 执行 SQL 语句
    err = db.Exec(string(sqlFile)).Error
    if err != nil {
        c.AbortWithError(http.StatusInternalServerError, err)
        return
    }

    // 返回成功响应
    c.JSON(http.StatusOK, gin.H{
        "message": "Tables created successfully",
    })
}

```
## 在创建表之后删除SQL文件
```go
func deleteSQLFile() {
    err := os.Remove("path/to/save/sql_file.sql")
    if err != nil {
        // 处理错误
    }
}

```