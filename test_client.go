package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// This script simulates exactly what redis-cli does under the hood.
func main() {
	// 1. Open a raw TCP connection to the GoRedis server
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to GoRedis server at localhost:6379!")
	fmt.Println("This script will send exact RESP bytes over the network wire.")
	fmt.Println("-----------------------------------------------------------------")

	// 2. We construct the exact byte string that `redis-cli SET score 100` sends:
	// Array of 3 elements: SET, score, 100
	setRequest := "*3\r\n$3\r\nSET\r\n$5\r\nscore\r\n$3\r\n100\r\n"
	
	fmt.Printf("Sending SET command: %q\n", setRequest)
	_, err = conn.Write([]byte(setRequest))
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}

	// 3. Read the server's response
	reader := bufio.NewReader(conn)
	response, _ := reader.ReadString('\n')
	fmt.Printf("Server responded: %q\n", response)

	fmt.Println("-----------------------------------------------------------------")

	// 4. We construct the exact byte string that `redis-cli GET score` sends:
	// Array of 2 elements: GET, score
	getRequest := "*2\r\n$3\r\nGET\r\n$5\r\nscore\r\n"

	fmt.Printf("Sending GET command: %q\n", getRequest)
	_, err = conn.Write([]byte(getRequest))
	if err != nil {
		fmt.Println("Error writing:", err)
		return
	}

	// 5. Read the server's response
	response, _ = reader.ReadString('\n')
	response2, _ := reader.ReadString('\n')
	
	fmt.Printf("Server responded: %q%q\n", response, response2)
}
