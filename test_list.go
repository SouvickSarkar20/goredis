package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// 1. Send LPUSH mylist apple banana
	// This tests pushing multiple values to the front at once!
	lpushRequest := "*4\r\n$5\r\nLPUSH\r\n$6\r\nmylist\r\n$5\r\napple\r\n$6\r\nbanana\r\n"
	
	fmt.Printf("Sending LPUSH command: %q\n", lpushRequest)
	conn.Write([]byte(lpushRequest))
	response, _ := reader.ReadString('\n')
	fmt.Printf("Server responded: %q\n\n", response)

	// 2. Send LPOP mylist
	// Since we pushed 'apple' then 'banana' to the front, 'banana' should be first out!
	lpopRequest := "*2\r\n$4\r\nLPOP\r\n$6\r\nmylist\r\n"

	fmt.Printf("Sending first LPOP command: %q\n", lpopRequest)
	conn.Write([]byte(lpopRequest))
	response, _ = reader.ReadString('\n')
	response2, _ := reader.ReadString('\n')
	fmt.Printf("Server responded: %q%q\n\n", response, response2)

	// 3. Send second LPOP mylist
	// This should return 'apple'
	fmt.Printf("Sending second LPOP command: %q\n", lpopRequest)
	conn.Write([]byte(lpopRequest))
	response, _ = reader.ReadString('\n')
	response2, _ = reader.ReadString('\n')
	fmt.Printf("Server responded: %q%q\n\n", response, response2)
	
	// 4. Send third LPOP mylist
	// List is empty, so it should return nil ($-1\r\n)
	fmt.Printf("Sending third LPOP command (empty list): %q\n", lpopRequest)
	conn.Write([]byte(lpopRequest))
	response, _ = reader.ReadString('\n')
	fmt.Printf("Server responded: %q\n\n", response)
}
