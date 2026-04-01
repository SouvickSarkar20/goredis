package main

import (
	"fmt"
	"github.com/krishsinghhura/goredis"
)

func main() {
	// Connect to our cloud server (goredis.me:6379)
	// You can also connect to localhost:6379 if running locally
	client, err := goredis.NewClient("goredis.me:6379")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer client.Close()

	// 1. SET a value
	fmt.Println("Setting key: language to Go")
	err = client.Set("language", "Go")
	if err != nil {
		fmt.Printf("Set Error: %v\n", err)
		return
	}

	// 2. GET the value back
	val, err := client.Get("language")
	if err != nil {
		fmt.Printf("Get Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Success! Retrieved value: %s\n", val)

	// 3. HSET (Hash)
	fmt.Println("Setting hash user:1...")
	err = client.HSet("user:1", "name", "Krish")
	if err != nil {
		fmt.Printf("HSet Error: %v\n", err)
		return
	}

	// 4. HGET (Hash retrieval)
	fieldVal, err := client.HGet("user:1", "name")
	if err != nil {
		fmt.Printf("HGet Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Hash success! Field 'name' is: %s\n", fieldVal)
}
