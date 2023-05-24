package model

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sqlOJ/config"
)

// ExecSqlCreateTable 运行mysql命令，执行sql为文件，建表
func ExecSqlCreateTable(filePath string) error {
	cmd := exec.Command("mysql", "-uroot", "-p"+config.Password, "sqloj_exercises")
	input, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}
	cmd.Stdin = bytes.NewBuffer(input)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running command:", err)
		return err
	}
	fmt.Println(string(output))
	return nil
}

// ExecSqlDeleteTable 在练习数据库中删除tableName对应表名
func ExecSqlDeleteTable(tableName string) error {
	err := GetExeDB().Exec("drop table " + tableName).Error
	if err != nil {
		log.Println(err)
		return errors.New("删除数据表错误")
	}
	return nil
}
