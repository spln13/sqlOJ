package main

import (
	"sqlOJ/router"
)

func main() {
	server := router.InitServer()
	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}
