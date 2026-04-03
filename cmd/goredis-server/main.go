package main

import (
	"fmt"
	"github.com/SouvickSarkar20/goredis/server"
	"os"
)

func main() {
	fmt.Println("Starting GoRedis server...")

	err := server.Start(":6379")
	if err != nil {
		fmt.Printf("Critical server error: %v\n", err)
		os.Exit(1)
	}
}
