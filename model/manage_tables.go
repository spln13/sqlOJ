package model

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// ExecSqlCreateTable 运行mysql命令，执行sql为文件，建表
func ExecSqlCreateTable(filePath string) error {
	cmd := exec.Command("mysql", "-uroot", "-pspln13spln", "sqloj_exercises")
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
