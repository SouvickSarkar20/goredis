package server

import (
	"fmt"
	"github.com/krishsinghhura/go-redis/cmd"
	"github.com/krishsinghhura/go-redis/persistence"
	"github.com/krishsinghhura/go-redis/resp"
	"github.com/krishsinghhura/go-redis/store"
	"github.com/krishsinghhura/go-redis/web"
	"io"
	"net"
	"os"
)

func Start(port string) error {

	listener, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Printf("Listening on %s\n", port)

	db := store.NewStore()

	os.MkdirAll("data", 0755) // create the data folder if it doesn't exist
	aof, err := persistence.OpenAOF("data/appendonly.aof", persistence.FsyncAlways)
	if err != nil {
		return fmt.Errorf("failed to open AOF: %w", err)
	}
	defer aof.Close()

	err = persistence.ReplayAOFTruncateTail("data/appendonly.aof", func(args []string) error {
		return cmd.ApplyAOFCommand(db, args)
	})

	if err != nil {
		return fmt.Errorf("failed to replay AOF: %w", err)
	}
	fmt.Println("AOF replay complete")

	cmd.SetAOF(aof)

	go func() {
		fmt.Println("Web UI starting on :8080")
		if err := web.StartHTTP(":8080", db); err != nil {
			fmt.Println("HTTP server error:", err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, db)
	}
}

func handleConnection(conn net.Conn, db *store.Store) {

	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("New client connected: %s\n", remoteAddr)

	defer conn.Close()

	w := resp.NewWriter(conn)
	p := resp.NewParser(conn)

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
