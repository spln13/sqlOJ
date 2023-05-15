package main

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc"
	"log"
	"sqlOJ/controller/exercise"
	"sqlOJ/fabric"
	"sqlOJ/router"
)

func main() {
	server := router.InitServer()    // 初始化Gin服务器
	exercise.InitJudgeGoroutine(100) // 打开100个判题协程
	clientConnection, gw := fabric.InitFabricConnection()
	defer func(clientConnection *grpc.ClientConn) { // 关闭grpc连接
		err := clientConnection.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(clientConnection)
	defer func(gw *client.Gateway) { // 关闭gateway连接
		err := gw.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(gw)

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}

// GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn go run main.go
