package main

import (
	"sqlOJ/controller/exercise"
	"sqlOJ/router"
)

func main() {
	server := router.InitServer()
	exercise.InitJudgeGoroutine(100) // 打开100个判题协程
	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
