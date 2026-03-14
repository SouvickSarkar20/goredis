package server

import (
	"fmt"
	"io"
	"net"

	"goredis/cmd"
	"goredis/resp"
	"goredis/store"
)

func Start(port string) error {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Printf("Listening on %s\n", port)

	// Create our single global instance of the database store!
	db := store.NewStore()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		
		// Pass the shared database to each new client connection
		go handleConnection(conn, db)
	}
}

// handleConnection manages the lifecycle of a single client connection.
func handleConnection(conn net.Conn, db *store.Store) {
	// Let's print the remote address so we know who connected.
	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("New client connected: %s\n", remoteAddr)

	// 1. Ensure the connection is strictly closed when the waiter is done serving the table.
	// This prevents memory and file-descriptor leaks.
	defer conn.Close()

	// 2. We initialize our RESP writer for this client.
	w := resp.NewWriter(conn)
	
	// 3. We initialize our RESP parser for this client.
	// It wraps the raw connection and handles all the buffering and byte-reading for us!
	p := resp.NewParser(conn)

	// 4. The Waiter's infinite loop. Stay at the table until the client leaves.
	for {
		// ParseOne blocks until the client sends a complete RESP value (like an Array of Strings).
		value, err := p.ParseOne()
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client disconnected gracefully: %s\n", remoteAddr)
			} else {
				fmt.Printf("Error pulling command from client %s: %v\n", remoteAddr, err)
			}
			return
		}

		// 5. Let's print out what command we just parsed!
		// 'value' should be an Array where all the elements are BulkStrings.
		if value.Typ != "Array" {
			fmt.Println("Expected Array, got:", value.Typ)
			continue
		}

		// 6. Route the command using our Cmd package!
		// It will extract the command name, call the right handle function,
		// and use the Writer to send the response back.
		err = cmd.Handle(w, db, value)
		if err != nil {
			fmt.Printf("Error writing to client %s: %v\n", remoteAddr, err)
			return
		}
	}
}
