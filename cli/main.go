package main

import (
	"bufio"
	"fmt"
	"github.com/krishsinghhura/go-redis/resp"
	"net"
	"os"
	"strings"
)

func main() {
	// 1. Connect to our GoRedis Server
	host := "goredis.me:6379"
	if len(os.Args) > 1 {
		host = os.Args[1]
	}

	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Printf("❌ Failed to connect to %s: %v\n", host, err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("🚀 Connected to GoRedis at %s\n", host)
	fmt.Println("Type commands (e.g., SET name GoRedis) or 'exit' to quit.")

	reader := bufio.NewReader(os.Stdin)
	parser := resp.NewParser(conn)
	writer := resp.NewWriter(conn)

	for {
		fmt.Print("goredis> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye! 👋")
			break
		}

		if input == "" {
			continue
		}

		// 2. Turn the string into a RESP Array of BulkStrings
		parts := strings.Fields(input)
		args := make([]resp.Value, len(parts))
		for i, p := range parts {
			args[i] = resp.Value{Typ: "BulkString", Str: p}
		}
		command := resp.Value{Typ: "Array", Array: args}

		// 3. Send to Server
		err := writer.Write(command)
		if err != nil {
			fmt.Printf("❌ Error sending command: %v\n", err)
			break
		}

		// 4. Read Response from Server
		val, err := parser.ParseOne()
		if err != nil {
			fmt.Printf("❌ Error reading response: %v\n", err)
			break
		}

		// 5. Display Result
		printValue(val)
	}
}

func printValue(v resp.Value) {
	switch v.Typ {
	case "SimpleString":
		fmt.Println(v.Str)
	case "BulkString":
		if v.Str == "" {
			fmt.Println("(empty or nil)")
		} else {
			fmt.Printf("\"%s\"\n", v.Str)
		}
	case "Integer":
		fmt.Printf("(integer) %d\n", v.Num)
	case "Error":
		fmt.Printf("(error) %s\n", v.Str)
	case "Array":
		for i, item := range v.Array {
			fmt.Printf("%d) ", i+1)
			printValue(item)
		}
	default:
		fmt.Printf("%s\n", v.Str)
	}
}
